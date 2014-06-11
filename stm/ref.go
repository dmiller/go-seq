// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
)


// A Ref holds a value that can be updated in an STM transaction.
type Ref {

}

func (r *Ref) Deref(tx *Tx) interface{} {

}
        /// <summary>
        /// Gets the (immutable) value the reference is holding.
        /// </summary>
        /// <returns>The value</returns>
        public override object deref()
        {
            LockingTransaction t = LockingTransaction.GetRunning();
            if (t == null)
            {
                object ret = currentVal();
                //Console.WriteLine("Thr {0}, {1}: No-trans get => {2}", Thread.CurrentThread.ManagedThreadId,DebugStr(), ret);
                return ret;
            }
            return t.DoGet(this);
        }

        object currentVal()
        {
            try
            {
                _lock.EnterReadLock();
                if (_tvals != null)
                    return _tvals.Val;
                throw new InvalidOperationException(String.Format("{0} is unbound.", ToString()));
            }
            finally
            {
                _lock.ExitReadLock();
            }
        }

/*
        #region Nested classes

        /// <summary>
        /// Represents the value of reference on a thread at particular point in time.
        /// </summary>
        public sealed class TVal
        {
            #region Data
 
            /// The value.
            object _val;

            /// The transaction commit/read point at which this value was set.
            long _point;


            /// The prior <see cref="TVal">TVal</see>.
            /// <remarks>Implements a doubly-linked circular list.</remarks>
            TVal _prior;

            /// The next  <see cref="TVal">TVal</see>.
            /// <remarks>Implements a doubly-linked circular list.</remarks>
            TVal _next;

            #endregion

            #region Ctors

            /// Construct a TVal, linked to a previous TVal.
            public TVal(object val, long point, int msecs, TVal prior)
            {
                _val = val;
                _point = point;
                _prior = prior;
                _next = _prior._next;
                _prior._next = this;
                _next._prior = this;
            }

            /// Construct a TVal, linked to itself.
            public TVal(object val, long point, int msecs)
            {
                _val = val;
                _point = point;
                _prior = this;
                _next = this;
            }

            #endregion

            #region other

            /// Set the value/point.
            public void SetValue(object val, long point, int msecs)
            {
                _val = val;
                _point = point;
            }

            #endregion
        }

        #endregion

        #region Data

        /// Values at points in time for this reference.
        TVal _tvals;

        /// Number of faults for the reference.
        readonly AtomicInteger _faults;

        /// Reader/writer lock for the reference.
        readonly ReaderWriterLockSlim _lock;

        /// Info on the transaction locking this ref.
        LockingTransaction.Info _tinfo;

        /// An id uniquely identifying this reference.
        readonly long _id;


        volatile int _minHistory = 0;
        volatile int _maxHistory = 10;
        // Same for min
        public Ref setMaxHistory(int maxHistory)
        {
            _maxHistory = maxHistory;
            return this;
        }

        /// Used to generate unique ids.
        static readonly AtomicLong _ids = new AtomicLong();

        bool _disposed = false;

        #endregion

        #region C-tors & factory methods

        ///  Construct a ref with given initial value.
        public Ref(object initVal)
            : this(initVal, null)
        {
        }
        ///  Construct a ref with given initial value and metadata.
        public Ref(object initval, IPersistentMap meta)
            : base(meta)
        {
            _id = _ids.getAndIncrement();
            _faults = new AtomicInteger();
            _lock = new ReaderWriterLockSlim(LockRecursionPolicy.NoRecursion);
            _tvals = new TVal(initval, 0, System.Environment.TickCount);
        }

        #endregion

        #region History counts

        public int getHistoryCount()
        {
            try
            {
                EnterWriteLock();
                return HistCount();
            }
            finally
            {
                ExitWriteLock();
            }
        }

        int HistCount()
        {
            if (_tvals == null)
                return 0;
            else
            {
                int count = 0;
                for (TVal tv = _tvals.Next; tv != _tvals; tv = tv.Next)
                    count++;
                return count;
            }
        }

        #endregion       

        #region IDeref Members

        /// Gets the (immutable) value the reference is holding.
        public override object deref()
        {
            LockingTransaction t = LockingTransaction.GetRunning();
            if (t == null)
            {
                object ret = currentVal();
                //Console.WriteLine("Thr {0}, {1}: No-trans get => {2}", Thread.CurrentThread.ManagedThreadId,DebugStr(), ret);
                return ret;
            }
            return t.DoGet(this);
        }

        object currentVal()
        {
            try
            {
                _lock.EnterReadLock();
                if (_tvals != null)
                    return _tvals.Val;
                throw new InvalidOperationException(String.Format("{0} is unbound.", ToString()));
            }
            finally
            {
                _lock.ExitReadLock();
            }
        }

        #endregion

        #region  Interface for LockingTransaction

        /// Get the read lock.
        internal void EnterReadLock()
        {
            _lock.EnterReadLock();
        }

        /// Release the read lock.
        internal void ExitReadLock()
        {
            _lock.ExitReadLock();
        }

        /// Get the write lock.
        internal void EnterWriteLock()
        {
            _lock.EnterWriteLock();
        }


        /// Get the write lock.
        internal bool TryEnterWriteLock(int msecTimeout)
        {
            return _lock.TryEnterWriteLock(msecTimeout);
        }

        /// Release the write lock.
        internal void ExitWriteLock()
        {
            _lock.ExitWriteLock();
        }

        /// Add to the fault count.
        public void AddFault()
        {
            _faults.incrementAndGet();
        }

        /// Get the read/commit point associated with the current value.
        public long CurrentValPoint()
        {
            return _tvals != null ? _tvals.Point : -1;
        }

        /// Try to get the value (else null).
        public object TryGetVal()
        {
            return _tvals == null ? null : _tvals.Val;
        }

        /// Set the value.
        internal void SetValue(object val, long commitPoint, int msecs)
        {
            int hcount = HistCount();

            if (_tvals == null)
                _tvals = new TVal(val, commitPoint, msecs);
            else if ( (_faults.get() > 0 && hcount < _maxHistory) || hcount < _minHistory )
            {
                _tvals = new TVal(val, commitPoint, msecs, _tvals);
                _faults.set(0);
            }
            else
            {
                _tvals = _tvals.Next;
                _tvals.SetValue(val, commitPoint, msecs);
            }
        }

        #endregion

        #region Ref operations

        /// Set the value (must be in a transaction).
        public object set(object val)
        {
            return LockingTransaction.GetEx().DoSet(this, val);
        }

        /// Apply a commute to the reference. (Must be in a transaction.)
        public object commute(IFn fn, ISeq args)
        {
            return LockingTransaction.GetEx().DoCommute(this, fn, args);
        }

        /// Change to a computed value.
        public object alter(IFn fn, ISeq args)
        {
            LockingTransaction t = LockingTransaction.GetEx();
            return t.DoSet(this, fn.applyTo(RT.cons(t.DoGet(this), args)));
        }

        /// Touch the reference.  (Add to the tracking list in the current transaction.)
        [System.Diagnostics.CodeAnalysis.SuppressMessage("Microsoft.Naming", "CA1709:IdentifiersShouldBeCasedCorrectly", MessageId = "touch")]
        public void touch()
        {
            LockingTransaction.GetEx().DoEnsure(this);
        }

        #endregion

        #region object overrides

        public override bool Equals(object obj)
        {
            if (ReferenceEquals(this, obj))
                return true;

            Ref r = obj as Ref;
            if (r == null)
                return false;

            return _id == r._id;
        }

        public override int GetHashCode()
        {
            return _id.GetHashCode();
        }
        #endregion

*/
