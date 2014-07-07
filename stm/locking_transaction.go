// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// retryLimit is the number of times to retry a transaction in case of a conflict.
	retryLimit = 10000

	// LockWaitMsecs is the number of milliseconds to wait for a lock
	lockWaitMsecs = 100 * time.Millisecond

	// BargeWaitNanos is the number of nanoseconds old another transaction must be before we 'barge' it.
	bargeWaitNanos = 10 * time.Nanosecond
)

const (
	// RUNNING transaction state
	txRunning uint32 = iota
	txCommitting
	txRetry
	txKilled
	txCommitted
)

// TxInfo represents the current state of an STM transaction
type TxInfo struct {
	status     uint32
	startPoint uint64
	latch      *CountDownLatch
	lock       sync.Mutex
}

func newTxInfo(status uint32, startPoint uint64) *TxInfo {
	return &TxInfo{status: status, startPoint: startPoint, latch: NewCountDownLatch(1)}
}

func (info *TxInfo) isRunning() bool {
	s := atomic.LoadUint32(&info.status)
	return s == txRunning || s == txCommitting
}

// lastPoint is the current point
// Used to provide a total ordering on transactions
// for the purpose of determining preference on transactions
// where there are conflicts.
// Transactions consume a point for init, for each retry,
// and on commit if writing
var lastPoint = new(IDGenerator)

// A CFn is a function suitable for calling as a commute on a ref
type CFn func(interface{}, ...interface{}) interface{}

// A CFnCall is a pending call of a function on arguments.
// Used to store commute calls on refs.
type CFnCall struct {
	fn   CFn
	args []interface{}
}

// A TxFn is a function suitable for calling in a transaction
type TxFn func(*Tx) interface{}

// Cached error to use in panics to signal a retry
var retryError error = errors.New("Retry")

// Tx provides STM transaction semantics for Agents and Refs
type Tx struct {
	// The state of the transaction
	info *TxInfo

	// The point at the start of the current retry (or first try).
	readPoint uint64

	// The point at the start of the transaction.
	startPoint uint64

	// The system time at the start of the transaction
	startTime time.Time

	// Agent actions pending on this thread
	//actions []Action

	// Ref assignments made in this transaction (both sets and commutes)
	vals map[*Ref]interface{}

	// Refs that have been set in this transaction
	sets map[*Ref]bool

	// Ref commutes that have been made in this transaction
	commutes map[*Ref][]*CFnCall

	// Refs holding read locks
	ensures map[*Ref]bool
}

func NewTx() *Tx {
	return &Tx{
		info:     nil,
		vals:     make(map[*Ref]interface{}),
		sets:     make(map[*Ref]bool),
		commutes: make(map[*Ref][]*CFnCall),
		ensures:  make(map[*Ref]bool),
	}
}

// Point manipulation

// Get a new read point value
func (tx *Tx) getReadPoint() {
	tx.readPoint = lastPoint.Next()
}

func getCommitPoint() uint64 {
	return lastPoint.Next()
}

// Actions

// Stop this transaction
func (tx *Tx) Stop(s uint32) {
	if tx.info == nil {
		return
	}

	tx.info.setStatus(s, true)
	tx.info = nil
	tx.vals = make(map[*Ref]interface{})
	tx.sets = make(map[*Ref]bool)
	tx.commutes = make(map[*Ref][]*CFnCall)
	// Java commented out: _actions.Clear();
	// Note that tx.ensured is not cleared

}

func (t *TxInfo) setStatus(s uint32, countDown bool) {

	t.lock.Lock()
	defer t.lock.Unlock()
	atomic.StoreUint32(&t.status, s)
	if countDown {
		t.latch.CountDown()
	}
}

func tryWriteLock(r *Ref) {
	if !r.tryEnterWriteLock(lockWaitMsecs) {
		panic(retryError)
	}
}

func (tx *Tx) releaseIfEnsured(r *Ref) {
	if _, ok := tx.ensures[r]; ok {
		delete(tx.ensures, r)
		r.exitReadLock()
	}
}

func (tx *Tx) blockAndBail(refinfo *TxInfo) interface{} {
	// stop prior to blocking
	tx.Stop(txRetry)
	refinfo.latch.Await(lockWaitMsecs)
	panic(retryError)
}

func (tx *Tx) lockRef(r *Ref) interface{} {
	// can't upgrade a read lock, so release it
	// TODO: Determine if this is true
	tx.releaseIfEnsured(r)

	locked := false
	defer func() {
		if locked {
			r.exitWriteLock()
		}
	}()

	tryWriteLock(r)
	locked = true

	if r.currValPoint() > tx.readPoint {
		panic(retryError)
	}

	refinfo := r.tinfo

	// write lock conflict
	if refinfo != nil && refinfo != tx.info && refinfo.isRunning() {
		if !tx.barge(refinfo) {
			r.exitWriteLock()
			locked = false
			return tx.blockAndBail(refinfo)
		}
	}

	r.tinfo = tx.info
	return r.tryGetVal()
}

// Barging

// Determine if sufficient clock time has elapsed to barge another transaction
// Returns true if enough time elapsed, false otherwise
func (tx *Tx) bargeTimeElapsed() bool {
	return time.Now().Sub(tx.startTime) > bargeWaitNanos
}

// Try to barge a conflicting transation
func (tx *Tx) barge(refinfo *TxInfo) bool {
	barged := false

	// if this transation is older, try to abort the other
	if tx.bargeTimeElapsed() && tx.startPoint < refinfo.startPoint {
		barged = atomic.CompareAndSwapUint32(&refinfo.status, txRunning, txKilled)
		if barged {
			refinfo.latch.CountDown()
		}
	}

	return barged
}

// Start a transaction and invoke a function, passing it the transaction.
// Returns the value computed by the function.
func RunInTransaction(fn TxFn) (interface{}, error) {
	tx := NewTx()
	return tx.Run(fn)
}

func (tx *Tx) Run(fn TxFn) (interface{}, error) {

	done := false
	locked := make([]*Ref, 0, 10)
	// notify := make([]*Notify)

	defer func() {
		for k := len(locked) - 1; k >= 0; k-- {
			locked[k].exitWriteLock()
		}

		locked = nil
		for r, _ := range tx.ensures {
			r.exitReadLock()
		}
		tx.ensures = nil
		if done {
			tx.Stop(txCommitted)
		} else {
			tx.Stop(txRetry)
		}
		if done {
			// do notifies and agent actions, if we every implement
		}

	}()

	for i := 0; !done && i < retryLimit; i++ {
		ret, err := tx.tryRun(i, fn, &locked)
		if err == nil {
			done = true
			return ret, nil
		}
	}

	return nil, errors.New("Transaction failed after reaching retry limit")
}

// One iteration of the Run loop
// Split out so that we can catch a retry panic
func (tx *Tx) tryRun(i int, fn TxFn, locked *[]*Ref) (ret interface{}, err error) {

	ret, err = nil, nil

	defer func() {
		r := recover()
		if r == retryError {
			ret, err = nil, retryError
		} else if r != nil {
			panic(r)
		}
	}()

	tx.getReadPoint()
	if i == 0 {
		tx.startPoint = tx.readPoint
		tx.startTime = time.Now()
	}

	tx.info = newTxInfo(txRunning, tx.startPoint)
	ret = fn(tx)

	// make sure no one has killed us before this point, and can't from now on
	if atomic.CompareAndSwapUint32(&tx.info.status, txRunning, txCommitting) {
		for r, calls := range tx.commutes {
			if _, ok := tx.sets[r]; ok {
				continue
			}
			_, wasEnsured := tx.ensures[r]
			tx.releaseIfEnsured(r)
			tryWriteLock(r)
			*locked = append(*locked, r)
			if wasEnsured && r.currValPoint() > tx.readPoint {
				panic(retryError)
			}

			refInfo := r.tinfo
			if refInfo != nil && refInfo != tx.info && refInfo.isRunning() {
				if !tx.barge(refInfo) {
					panic(retryError)
				}
			}
			val := r.tryGetVal()
			tx.vals[r] = val
			for _, call := range calls {
				tx.vals[r] = call.fn(tx.vals[r], call.args...)
			}
		}

		for r, _ := range tx.sets {
			tryWriteLock(r)
			*locked = append(*locked, r)
		}

		// if we do validations for refs, it goes here

		// at this point,
		//    all values are calculated,
		//    all refs to be written are locked
		//    no more client code to be called
		commitPoint := getCommitPoint()
		for r, newV := range tx.vals {
			//oldV := r.tryGetVal()
			r.setValue(newV, commitPoint)
			// todo: call notifies
		}
		atomic.StoreUint32(&tx.info.status, txCommitted)
	}

	return
}

// Get the value of a Ref (most recently sent in this transaction or value prior to entering)
func (tx *Tx) doGet(r *Ref) interface{} {
	if !tx.info.isRunning() {
		panic(retryError)
	}
	if v, ok := tx.vals[r]; ok {
		return v
	}
	r.enterReadLock()
	defer r.exitReadLock()
	if r.tvals == nil {
		panic(fmt.Errorf("%v is not bound", r))
	}
	ver := r.tvals
	for {
		if ver.point <= tx.readPoint {
			return ver.val
		}
		ver = ver.prior
		if ver == r.tvals {
			break
		}
	}

	// no version of val precedes the read point
	r.addFault()
	panic(retryError)
}

// Set the value of a Ref inside the transaction
func (tx *Tx) doSet(r *Ref, v interface{}) interface{} {
	if !tx.info.isRunning() {
		panic(retryError)
	}
	if _, ok := tx.commutes[r]; ok {
		panic(errors.New("Can't set after commute"))
	}
	if _, ok := tx.sets[r]; !ok {
		tx.sets[r] = true
		tx.lockRef(r)
	}
	tx.vals[r] = v
	return v
}

func (tx *Tx) doEnsure(r *Ref) {
	if !tx.info.isRunning() {
		panic(retryError)
	}
	if _, ok := tx.ensures[r]; ok {
		return
	}
	r.enterReadLock()

	// someone completed a write after our shapshot
	if r.currValPoint() > tx.readPoint {
		r.exitReadLock()
		panic(retryError)
	}

	refInfo := r.tinfo

	// writer exists
	if refInfo != nil && refInfo.isRunning() {
		r.exitReadLock()
		if refInfo != tx.info {
			// not us, ensure is doomed
			tx.blockAndBail(refInfo)
		}
	} else {
		tx.ensures[r] = true
	}
}

func (tx *Tx) GetAndStoreRefVal(r *Ref) {
	if _, ok := tx.vals[r]; !ok {
		var val interface{}
		r.enterReadLock()
		defer r.exitReadLock()
		val = r.tryGetVal()
		tx.vals[r] = val
	}
}

// Post a commute on a ref into this transaction
func (tx *Tx) doCommute(r *Ref, fn CFn, args ...interface{}) interface{} {
	if !tx.info.isRunning() {
		panic(retryError)
	}
	tx.GetAndStoreRefVal(r)

	calls, ok := tx.commutes[r]
	if !ok {
		calls = make([]*CFnCall, 0, 5)
	} else {
	}
	calls = append(calls, &CFnCall{fn: fn, args: args})
	tx.commutes[r] = calls

	ret := fn(tx.vals[r], args...)
	tx.vals[r] = ret

	return ret
}

// Kill this transaction
func (tx *Tx) Abort() {
	tx.Stop(txKilled)
	panic(errors.New("Transaction aborted")) // handle some other way?
}

/*

   class Notify
   {
       public readonly Ref _ref;
       public readonly object _oldval;
       public readonly object _newval;

       public Notify(Ref r, object oldval, object newval)
       {
           _ref = r;
           _oldval = oldval;
           _newval = newval;
       }
   }


*/
