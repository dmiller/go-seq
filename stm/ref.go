// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// A tval holds the value for a Ref in a transaction
type tval struct {

	// the value
	val interface{}

	// the transaction commit/read point at which this value was set
	point uint64

	// the prior tval
	// with next, implements a doubly-linked circular list
	prior *tval

	// the next tval
	// with prior, implements a doubly-linked circular list
	next *tval
}

func newTvalPrior(val interface{}, point uint64, prior *tval) *tval {
	t := &tval{val: val, point: point, prior: prior, next: prior.next}
	prior.next = t
	t.next.prior = t
	return t
}

func newTval(val interface{}, point uint64) *tval {
	t := &tval{val: val, point: point}
	t.prior = t
	t.next = t
	return t
}

func (t *tval) setValue(val interface{}, point uint64) {
	t.val = val
	t.point = point
}

const (
	// DefaultMaxHistory is the default max history for Refs
	DefaultMaxHistory = 10
)

// A Ref holds a value that can be updated in an STM transaction.
type Ref struct {

	// Values at points in time for this Ref
	// Doubly-linked list, size controlled by history limit
	tvals *tval

	// Number of faults for the reference
	faults uint32

	lock sync.RWMutex

	minHistory uint
	maxHistory uint

	// TXInfo on the transaction locking this ref.
	tinfo *TxInfo

	id uint64
}

// id generator for Refs
var refIds = new(IDGenerator)

// Factories

func NewRef(val interface{}) *Ref {
	return &Ref{
		id:         refIds.Next(),
		maxHistory: DefaultMaxHistory,
		tvals:      newTval(val, 0)}
}

// Getting values

// Gets the value for the reference in the transaction
func (r *Ref) Deref(tx *Tx) interface{} {
	if tx == nil {
		return r.currentVal()
	}
	return tx.doGet(r)
}

func (r *Ref) currentVal() interface{} {
	r.enterReadLock()
	defer r.exitReadLock()
	if r.tvals == nil {
		panic(fmt.Errorf("%v is unbound", r))
	}
	return r.tvals.val
}

// history count, limits

func (r *Ref) SetMaxHistory(m uint) *Ref {
	r.maxHistory = m
	return r
}

func (r *Ref) SetMinHistory(m uint) *Ref {
	r.minHistory = m
	return r
}

func (r *Ref) MinHistory() uint {
	return r.minHistory
}

func (r *Ref) MaxHistory() uint {
	return r.maxHistory
}

func (r *Ref) HistoryCount() uint {
	r.enterWriteLock()
	defer r.exitWriteLock()
	return r.calcHistoryCount()
}

func (r *Ref) calcHistoryCount() uint {
	if r.tvals == nil {
		return 0
	}

	count := uint(0)
	for tv := r.tvals.next; tv != r.tvals; tv = tv.next {
		count = count + 1
	}
	return count
}

// Lock management

func (r *Ref) enterReadLock() {
	r.lock.RLock()
}

func (r *Ref) exitReadLock() {
	r.lock.RUnlock()
}

func (r *Ref) enterWriteLock() {
	r.lock.Lock()
}

func (r *Ref) exitWriteLock() {
	r.lock.Unlock()
}

func (r *Ref) tryEnterWriteLock(dur time.Duration) bool {
	lockChan := make(chan bool, 1)
	toChan := time.After(dur)
	timedOut := int32(0)

	go func() {
		r.lock.Lock()
		lockChan <- true
		defer func() {
			if atomic.LoadInt32(&timedOut) == 1 {
				r.lock.Unlock()
			}
		}()
	}()

	select {
	case <-toChan:
		atomic.StoreInt32(&timedOut, 1)
		return false
	case <-lockChan:
		return true
	}
}

// Fault management

// Add to the fault count
func (r *Ref) addFault() {
	atomic.AddUint32(&r.faults, 1)
}

func (r *Ref) getFault() uint32 {
	return atomic.LoadUint32(&r.faults)
}

func (r *Ref) setFault(v uint32) {
	atomic.StoreUint32(&r.faults, v)
}

// Value management

// Get the read/commit point associated with the current value
func (r *Ref) currValPoint() uint64 {
	if r.tvals == nil {
		return 0
	}
	return r.tvals.point
}

// Get current value (else null if no current value)
func (r *Ref) tryGetVal() interface{} {
	if r.tvals == nil {
		return nil
	}
	return r.tvals.val
}

// Set the value
func (r *Ref) setValue(val interface{}, commitPoint uint64) {
	hcnt := r.calcHistoryCount()

	if r.tvals == nil {
		r.tvals = newTval(val, commitPoint)
	} else if (r.getFault() > 0 && hcnt < r.maxHistory) || hcnt < r.minHistory {
		r.tvals = newTvalPrior(val, commitPoint, r.tvals)
		r.setFault(0)
	} else {
		r.tvals = r.tvals.next
		r.tvals.setValue(val, commitPoint)
	}
}

// public interface for Ref

func (r *Ref) Set(tx *Tx, val interface{}) interface{} {
	return tx.doSet(r, val)
}

func (r *Ref) Commute(tx *Tx, fn CFn, args ...interface{}) interface{} {
	return tx.doCommute(r, fn, args...)
}

func (r *Ref) Alter(tx *Tx, fn CFn, args ...interface{}) interface{} {
	return tx.doSet(r, fn(tx.doGet(r), args...))
}

func (r *Ref) Touch(tx *Tx) {
	tx.doEnsure(r)
}
