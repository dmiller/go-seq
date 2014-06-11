// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
)

const (
	// RetryLimit is the number of times to retry a transaction in case of a conflict.
	RetryLimit = 10000

	// LockWaitMsecs is the number of milliseconds to wait for a lock
	LockWaitMsecs = 100

	// BargeWaitNanos is the number of nanoseconds old another transaction must be before we 'barge' it.
	BargeWaitNanos = 10*1000000
)

const (
	// RUNNING transaction state
	RUNNING = iota
	COMMITTING
	RETRY
	KILLED
	COMMITTED
)

// Info represents the current state of an STM transaction
type Info struct {
	status int8
	startPoint int64
	latch CountDownLatch
}

func newInfo(status int8, startPoint int64) *Info {
	return &Info{status:status,startPoint:startPoint,latch:NewCountDownLatch(1)}
}

func (info *Info) isRunning() bool {
	s = info.status
	return s == RUNNING || s == COMMITTING
}


// lastPoint is the current point
// Used to provide a total ordering on transactions
// for the purpose of determining preference on transactions
// where there are conflicts.
// Transactions consume a point for init, for each retry,
// and on commit if writing
var lastPoint int64



// Tx provides STM transaction semantics for Agents and Refs
type Tx struct {
	// The state of the transaction
	info *Info

	// The point at the start of the current retry (or first try).
	readPoint int64

	// The point at the start of the transaction.
	startPoint int64

	// The system ticks at the start of the transaction
	startTime int64

    /// <summary>
    /// Cached retry exception.
    /// </summary>
    //readonly RetryEx _retryex = new RetryEx();

	// Agent actions pending on this thread
	actions []Action

	// Ref assignments made in this transaction (both sets and commutes)
	vals map[Ref] interface{}

	// Refs that have been set in this transaction
	sets map[Ref] bool

	// Ref commutes that have been made in this transaction
	commutes map[Ref] []CFn

	// Refs holding read locks
	ensures map[Ref] bool
}


// Point manipulation

// Get a new read point value
func (tx *StmTx) getReadPoint() {
	tx.readPoint = Atomic.AddInt64(&lastPoint,1)
}

func(tx *StmTx) getCommitPoint() int64 {
	return Atomic.AddInt64(&lastPoint,1)
}


// Actions



/*

        #region Actions

        /// <summary>
        /// Stop this transaction.
        /// </summary>
        /// <param name="status">The new status.</param>
        void Stop(int status)
        {
            if (_info != null)
            {
                lock (_info)
                {
                    _info.Status.set(status);
                    _info.Latch.CountDown();
                }
                _info = null;
                _vals.Clear();
                _sets.Clear();
                _commutes.Clear();
                // Java commented out: _actions.Clear();
            }
        }

        void TryWriteLock(Ref r)
        {
            try
            {
                if (!r.TryEnterWriteLock(LockWaitMsecs))
                    throw _retryex;
            }
            catch (ThreadInterruptedException )
            {
                throw _retryex;
            }
        }

        void ReleaseIfEnsured(Ref r)
        {
            if (_ensures.Contains(r))
            {
                _ensures.Remove(r);
                r.ExitReadLock();
            }
        }


        object BlockAndBail(Info refinfo)
        {
            //stop prior to blocking
            Stop(RETRY);
            try
            {
                refinfo.Latch.Await(LockWaitMsecs);
            }
            catch (ThreadInterruptedException)
            {
                //ignore
            }
            throw _retryex;
        }


        /// <summary>
        /// Lock a ref.
        /// </summary>
        /// <param name="r">The ref to lock.</param>
        /// <returns>The most recent value of the ref.</returns>
        object Lock(Ref r)
        {
            // can't upgrade read lock, so release it.
            ReleaseIfEnsured(r);

            bool unlocked = true;
            try
            {
                TryWriteLock(r);
                unlocked = false;

                if (r.CurrentValPoint() > _readPoint)
                    throw _retryex;

                Info refinfo = r.TInfo;

                // write lock conflict
                if (refinfo != null && refinfo != _info && refinfo.IsRunning)
                {
                    if (!Barge(refinfo))
                    {
                        r.ExitWriteLock();
                        unlocked = true;
                        return BlockAndBail(refinfo);
                    }
                }

                r.TInfo = _info;
                return r.TryGetVal();
            }
            finally
            {
                if (!unlocked)
                {
                    r.ExitWriteLock();
                }
            }
        }

        /// <summary>
        /// Kill this transaction.
        /// </summary>
        void Abort()
        {
            Stop(KILLED);
            throw new AbortException();
        }

        /// <summary>
        /// Determine if sufficient clock time has elapsed to barge another transaction.
        /// </summary>
        /// <returns><value>true</value> if enough time has elapsed; <value>false</value> otherwise.</returns>
        private bool BargeTimeElapsed()
        {
            return Environment.TickCount - _startTime > BargeWaitTicks;
        }

        /// <summary>
        /// Try to barge a conflicting transaction.
        /// </summary>7
        /// <param name="refinfo">The info on the other transaction.</param>
        /// <returns><value>true</value> if we killed the other transaction; <value>false</value> otherwise.</returns>
        private bool Barge(Info refinfo)
        {
            bool barged = false;
            // if this transaction is older
            //   try to abort the other
            if (BargeTimeElapsed() && _startPoint < refinfo.StartPoint)
            {
                barged = refinfo.Status.compareAndSet(RUNNING, KILLED);
                if (barged)
                    refinfo.Latch.CountDown();
            }
            return barged;
        }

        /// <summary>
        /// Get the transaction running on this thread (throw exception if no transaction). 
        /// </summary>
        /// <returns>The running transaction.</returns>
        public static LockingTransaction GetEx()
        {
            LockingTransaction t = _transaction;
            if (t == null || t._info == null)
                throw new InvalidOperationException("No transaction running");
            return t;
        }

        /// <summary>
        /// Get the transaction running on this thread (or null if no transaction).
        /// </summary>
        /// <returns>The running transaction if there is one, else <value>null</value>.</returns>
        static internal LockingTransaction GetRunning()
        {
            LockingTransaction t = _transaction;
            if (t == null || t._info == null)
                return null;
            return t;
        }

        /// <summary>
        /// Is there a transaction running on this thread?
        /// </summary>
        /// <returns><value>true</value> if there is a transaction running on this thread; <value>false</value> otherwise.</returns>
        /// <remarks>Initial lowercase in name for core.clj compatibility.</remarks>
        [System.Diagnostics.CodeAnalysis.SuppressMessage("Microsoft.Naming", "CA1709:IdentifiersShouldBeCasedCorrectly")]
        public static bool isRunning()
        {
            return GetRunning() != null;
        }

        /// <summary>
        /// Invoke a function in a transaction
        /// </summary>
        /// <param name="fn">The function to invoke.</param>
        /// <returns>The value computed by the function.</returns>
        /// <remarks>Initial lowercase in name for core.clj compatibility.</remarks>
        [System.Diagnostics.CodeAnalysis.SuppressMessage("Microsoft.Naming", "CA1709:IdentifiersShouldBeCasedCorrectly")]
        public static object runInTransaction(IFn fn)
        {
            // TODO: This can be called on something more general than  an IFn.
            // We can could define a delegate for this, probably use ThreadStartDelegate.
            // Should still have a version that takes IFn.
            LockingTransaction t = _transaction;
            if (t == null)
                _transaction = t = new LockingTransaction();

            if (t._info != null)
                return fn.invoke();

            return t.Run(fn);
        }

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


        /// <summary>
        /// Start a transaction and invoke a function.
        /// </summary>
        /// <param name="fn">The function to invoke.</param>
        /// <returns>The value computed by the function.</returns>
        object Run(IFn fn)
        {
            // TODO: Define an overload called on ThreadStartDelegate or something equivalent.

            bool done = false;
            object ret = null;
            List<Ref> locked = new List<Ref>();
            List<Notify> notify = new List<Notify>();

            for (int i = 0; !done && i < RetryLimit; i++)
            {
                try
                {
                    GetReadPoint();
                    if (i == 0)
                    {
                        _startPoint = _readPoint;
                        _startTime = Environment.TickCount;
                    }

                    _info = new Info(RUNNING, _startPoint);
                    ret = fn.invoke();

                    // make sure no one has killed us before this point,
                    // and can't from now on
                    if (_info.Status.compareAndSet(RUNNING, COMMITTING))
                    {
                        foreach (KeyValuePair<Ref, List<CFn>> pair in _commutes)
                        {
                            Ref r = pair.Key;
                            if (_sets.Contains(r))
                                continue;

                            bool wasEnsured = _ensures.Contains(r);
                            // can't upgrade read lock, so release
                            ReleaseIfEnsured(r);
                            TryWriteLock(r);
                            locked.Add(r);

                            if (wasEnsured && r.CurrentValPoint() > _readPoint )
                                throw _retryex;

                            Info refinfo = r.TInfo;
                            if ( refinfo != null && refinfo != _info && refinfo.IsRunning)
                            {
                                if (!Barge(refinfo))
                                {
                                    throw _retryex;
                                }
                            }
                            object val = r.TryGetVal();
                            _vals[r] = val;
                            foreach (CFn f in pair.Value)
                                _vals[r] = f.Fn.applyTo(RT.cons(_vals[r], f.Args));
                        }
                        foreach (Ref r in _sets)
                        {
                            TryWriteLock(r);
                            locked.Add(r);
                        }
                        // validate and enqueue notifications
                        foreach (KeyValuePair<Ref, object> pair in _vals)
                        {
                            Ref r = pair.Key;
                            r.Validate(pair.Value);
                        }

                        // at this point, all values calced, all refs to be written locked
                        // no more client code to be called
                        int msecs = System.Environment.TickCount;
                        long commitPoint = GetCommitPoint();
                        foreach (KeyValuePair<Ref, object> pair in _vals)
                        {
                            Ref r = pair.Key;
                            object oldval = r.TryGetVal();
                            object newval = pair.Value;
                          
                            r.SetValue(newval, commitPoint, msecs);
                            if (r.getWatches().count() > 0)
                                notify.Add(new Notify(r, oldval, newval));
                        }

                        done = true;
                        _info.Status.set(COMMITTED);
                    }
                }
                catch (RetryEx)
                {
                    // eat this so we retry rather than fall out
                }
                catch (Exception ex)
                {
                    if (ContainsNestedRetryEx(ex))
                    {
                        // Wrapped exception, eat it.
                    }
                    else
                    {
                        throw;
                    }
                }
                finally
                {
                    for (int k = locked.Count - 1; k >= 0; --k)
                    {
                        locked[k].ExitWriteLock();
                    }
                    locked.Clear();
                    foreach (Ref r in _ensures)
                        r.ExitReadLock();
                    _ensures.Clear();
                    Stop(done ? COMMITTED : RETRY);
                    try
                    {
                        if (done) // re-dispatch out of transaction
                        {
                            foreach (Notify n in notify)
                            {
                                n._ref.NotifyWatches(n._oldval, n._newval);
                            }
                            foreach (Agent.Action action in _actions)
                            {
                                Agent.DispatchAction(action);
                            }
                        }
                    }
                    finally
                    {
                        notify.Clear();
                        _actions.Clear();
                    }
                }
            }
            if (!done)
                throw new InvalidOperationException("Transaction failed after reaching retry limit");
            return ret;
        }

        /// <summary>
        /// Determine if the exception wraps a <see cref="RetryEx">RetryEx</see> at some level.
        /// </summary>
        /// <param name="ex">The exception to test.</param>
        /// <returns><value>true</value> if there is a nested  <see cref="RetryEx">RetryEx</see>; <value>false</value> otherwise.</returns>
        /// <remarks>Needed because sometimes our retry exceptions get wrapped.  You do not want to know how long it took to track down this problem.</remarks>
        private static bool ContainsNestedRetryEx(Exception ex)
        {
            for (Exception e = ex; e != null; e = e.InnerException)
                if (e is RetryEx)
                    return true;
            return false;
        }

        /// <summary>
        /// Add an agent action sent during the transaction to a queue.
        /// </summary>
        /// <param name="action">The action that was sent.</param>
        internal void Enqueue(Agent.Action action)
        {
            _actions.Add(action);
        }

        /// <summary>
        /// Get the value of a ref most recently set in this transaction (or prior to entering).
        /// </summary>
        /// <param name="r"></param>
        /// <param name="tvals"></param>
        /// <returns>The value.</returns>
        internal object DoGet(Ref r)
        {
            if (!_info.IsRunning)
                throw _retryex;
            if (_vals.ContainsKey(r))
            {
                return _vals[r];
            }
            try
            {
                r.EnterReadLock();
                if (r.TVals == null)
                    throw new InvalidOperationException(r.ToString() + " is not bound.");
                Ref.TVal ver = r.TVals;
                do
                {
                    if (ver.Point <= _readPoint)
                    {
                        return ver.Val;
                    }
                } while ((ver = ver.Prior) != r.TVals);
            }
            finally
            {
                r.ExitReadLock();
            }
            // no version of val precedes the read point
            r.AddFault();
            throw _retryex;
        }

        /// <summary>
        /// Set the value of a ref inside the transaction.
        /// </summary>
        /// <param name="r">The ref to set.</param>
        /// <param name="val">The value.</param>
        /// <returns>The value.</returns>
        internal object DoSet(Ref r, object val)
        {
            if (!_info.IsRunning)
                throw _retryex;
            if (_commutes.ContainsKey(r))
                throw new InvalidOperationException("Can't set after commute");
            if (!_sets.Contains(r))
            {
                _sets.Add(r);
                Lock(r);
            }
            _vals[r] = val;
            return val;
        }

        /// <summary>
        /// Touch a ref.  (Lock it.)
        /// </summary>
        /// <param name="r">The ref to touch.</param>
        internal void DoEnsure(Ref r)
        {
            if (!_info.IsRunning)
                throw _retryex;
            if (_ensures.Contains(r))
                return;

            Lock(r);

            // someone completed a write after our shapshot
            if (r.CurrentValPoint() > _readPoint)
            {
                r.ExitReadLock();
                throw _retryex;
            }

            Info refinfo = r.TInfo;

            // writer exists
            if (refinfo != null && refinfo.IsRunning)
            {
                r.ExitReadLock();
                if (refinfo != _info)  // not us, ensure is doomed
                    BlockAndBail(refinfo);
            }
            else
                _ensures.Add(r);
        }


        /// <summary>
        /// Post a commute on a ref in this transaction.
        /// </summary>
        /// <param name="r">The ref.</param>
        /// <param name="fn">The commuting function.</param>
        /// <param name="args">Additional arguments to the function.</param>
        /// <returns>The computed value.</returns>
        internal object DoCommute(Ref r, IFn fn, ISeq args)
        {
            if (!_info.IsRunning)
                throw _retryex;
            if (!_vals.ContainsKey(r))
            {
                object val = null;
                try
                {
                    r.EnterReadLock();
                    val = r.TryGetVal();
                }
                finally
                {
                    r.ExitReadLock();
                }
                _vals[r] = val;
            }
            List<CFn> fns;
            if (! _commutes.TryGetValue(r, out fns))
                _commutes[r] = fns = new List<CFn>();
            fns.Add(new CFn(fn, args));
            object ret = fn.applyTo(RT.cons(_vals[r], args));
            _vals[r] = ret;

            return ret;
        }

        #endregion
    }
}

        /// <summary>
        /// The transaction running on the current thread.  (Thread-local.)
        /// </summary>
        //[ThreadStatic]
        //private static LockingTransaction _transaction;

        #region supporting classes

        /// <summary>
        /// Pending call of a function on arguments.
        /// </summary>
        class CFn
        {
            #region Data

            /// <summary>
            ///  The function to be called.
            /// </summary>
            readonly IFn _fn;

            /// <summary>
            ///  The function to be called.
            /// </summary>
            public IFn Fn
            {
                get { return _fn; }
            }

            /// <summary>
            /// The arguments to the function.
            /// </summary>
            readonly ISeq _args;

            /// <summary>
            /// The arguments to the function.
            /// </summary>
            public ISeq Args
            {
                get { return _args; }
            }

            #endregion

            #region C-tors

            /// <summary>
            /// Construct one.
            /// </summary>
            /// <param name="fn">The function to invoke.</param>
            /// <param name="args">The arguments to invoke the function on.</param>
            public CFn(IFn fn, ISeq args)
            {
                _fn = fn;
                _args = args;
            }

            #endregion
        }
        /// <summary>
        /// Exception thrown when a retry is necessary.
        /// </summary>
        [Serializable]
        public class RetryEx : Exception
        {
            #region C-tors

            public RetryEx()
            {
            }

            public RetryEx(String message)
                : base(message)
            {
            }

            public RetryEx(String message, Exception innerException)
                : base(message, innerException)
            {
            }

            protected RetryEx(SerializationInfo info, StreamingContext context)
                : base(info, context)
            {
            }

            #endregion

        }

        /// <summary>
        /// Exception thrown when a transaction has been aborted.
        /// </summary>
        [Serializable]
        public class AbortException : Exception
        {
            #region C-tors

            public AbortException()
            {
            }

            public AbortException(String message)
                : base(message)
            {
            }

            public AbortException(String message, Exception innerException)
                : base(message, innerException)
            {
            }

            protected AbortException(SerializationInfo info, StreamingContext context)
                : base(info, context)
            {
            }

            #endregion
        }
*/