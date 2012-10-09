// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"hash"
	"iseq"
	"sequtil"
)

// PList implements a persistent immutable list
//
// Zero value is not valid!
// We use 0 as an indicator that hash is not cached.
type PList struct {
	first interface{}
	rest  iseq.PList
	count int
	AMeta
	hash uint32
}

// PList needs to implement the iseq interfaces: 
//    Obj, Meta, Seq, Sequential, PList (= PCollection + PStack), Seqable, Counted, IHashEq
//   Also, Equatable and Hashable
//
// I'm not sure yet if I'll be doing IHashEq
// Also, Sequential is a marker interface that I haven't figured out how to translate
//    because I can't figure out a significant use of it in the Clojure code

// interface Meta is covered by the AMeta embedding

// c-tors

func NewPList1(first interface{}) *PList {
	return &PList{first: first, count: 1}
}

func NewPList1N(first interface{}, rest iseq.PList, count int) *PList {
	return &PList{first: first, rest: rest, count: count}
}

func NewPListMeta1N(meta iseq.PMap, first interface{}, rest iseq.PList, count int) *PList {
	return &PList{AMeta: AMeta{meta}, first: first, rest: rest, count: count}
}

func NewPListFromSlice(init []interface{}) *PList {
	var ret *PList
	count := 0
	for i := len(init) - 1; i >= 0; i-- {
		count++
		ret = NewPList1N(init[i], ret, count)
	}
	return ret
}

// interface iseq.Obj

func (p *PList) WithMeta(meta iseq.PMap) iseq.Obj {
	return NewPListMeta1N(meta, p.first, p.rest, p.count)
}

// interface iseq.Seqable

func (p *PList) Seq() iseq.Seq {
	return p
}

// interface iseq.PCollection

func (p *PList) Count() int {
	return p.count
}

func (p *PList) Cons(o interface{}) iseq.PCollection {
	return p.SCons(o)
}

func (p *PList) Empty() iseq.PCollection {
	return CachedEmptyList.WithMeta(p.meta).(iseq.PCollection)
}

func (p *PList) Equiv(o interface{}) bool {
	if p == o {
		return true
	}

	if os, ok := o.(iseq.Seqable); ok {
		return sequtil.SeqEquiv(p, os.Seq())
	}

	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
	return false
}

// interface iseq.Seq

func (p *PList) First() interface{} {
	return p.first
}

func (p *PList) Next() iseq.Seq {
	if p.count == 1 {
		return nil
	}
	return p.rest.Seq()
}

func (p *PList) More() iseq.Seq {
	s := p.Next()
	if s == nil {
		return CachedEmptyList
	}
	return s
}

func (p *PList) SCons(o interface{}) iseq.Seq {
	return NewPListMeta1N(p.meta, o, p, p.count+1)
}

// interface Counted

func (p *PList) Count1() int {
	return p.count
}

// PStack

func (p *PList) Peek() interface{} {
	return p.first
}

func (p *PList) Pop() iseq.PStack {
	if p.rest == nil {
		return CachedEmptyList.WithMeta(p.meta).(iseq.PStack)
	}
	return p.rest
}

// interfaces Equatable, Hashable

func (p *PList) Equals(o interface{}) bool {
	if p == o {
		return true
	}

	if os, ok := o.(iseq.Seqable); ok {
		return sequtil.SeqEquals(p, os.Seq())
	}

	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
	return false
}

func (p *PList) Hash() uint32 {
	if p.hash == 0 {
		p.hash = sequtil.HashSeq(p)
	}

	return p.hash
}

func (p *PList) AddHash(h hash.Hash) {
	if p.hash == 0 {
		p.hash = sequtil.HashSeq(p)
	}

	sequtil.AddHashUint64(h, uint64(p.hash))
}

/*
   #region IReduce Members

     /// <summary>
     /// Reduce the collection using a function.
     /// </summary>
     /// <param name="f">The function to apply.</param>
     /// <returns>The reduced value</returns>
     public object reduce(IFn f)
     {
         object ret = first();
         for (ISeq s = next(); s != null; s = s.next())
             ret = f.invoke(ret, s.first());
         return ret;
     }

     /// <summary>
     /// Reduce the collection using a function.
     /// </summary>
     /// <param name="f">The function to apply.</param>
     /// <param name="start">An initial value to get started.</param>
     /// <returns>The reduced value</returns>
     public object reduce(IFn f, object start)
     {
         object ret = f.invoke(start, first());
         for (ISeq s = next(); s != null; s = s.next())
             ret = f.invoke(ret, s.first());
         return ret;
     }

     #endregion
*/
