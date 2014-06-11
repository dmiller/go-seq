// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
)

// EmptyList implements an empty iseq.PList
// Used as a canonical empty list for other types that cannot represent empty.
type EmptyList struct {
	AMeta
}

var (
	CachedEmptyList = &EmptyList{}
)

// EmptyList needs to implement the iseq interfaces:
//   Meta, MetaW, Seq, Sequential, PCollection, PStack, PList, Seqable, Counted
//   Also, Equivable and Hashable
//
// Sequential is a marker interface that I haven't figured out how to translate
//    because I can't figure out a significant use of it in the Clojure code
// PList is just PStack + PCollection, so nothing added

// interface Meta is covered by the AMeta embedding

// interface iseq.MetaW

func (e *EmptyList) WithMeta(meta iseq.PMap) iseq.Obj {
	return &EmptyList{AMeta: AMeta{meta}}
}

// interface iseq.Seqable

func (e *EmptyList) Seq() iseq.Seq {
	return nil
}

// interface iseq.PCollection

func (e *EmptyList) Count() int {
	return 0
}

func (e *EmptyList) Cons(o interface{}) iseq.PCollection {
	return e.SCons(o)
}

func (e *EmptyList) Empty() iseq.PCollection {
	return e
}

// interface iseq.Seq

func (e *EmptyList) First() interface{} {
	return nil
}

func (e *EmptyList) Next() iseq.Seq {
	return nil
}

func (e *EmptyList) More() iseq.Seq {
	return e
}

func (e *EmptyList) ConsS(o interface{}) iseq.Seq {
	// TODO: really, this needs to return a PList of one element.
	// Fix when we have a true PList
	return &Cons{first: o, more: e}
}

// interface Counted

func (e *EmptyList) Count1() int {
	return 0
}

// PStack

func (e *EmptyList) Peek() interface{} {
	return nil
}

func (e *EmptyList) Pop() iseq.PStack {
	// in Clojure, popping throws an exception
	// should we add another return value?
	// For the moment, just return nil
	return nil
}

// interfaces Equivable, Hashable

func (e *EmptyList) Equiv(o interface{}) bool {
	if e == o {
		return true
	}

	if s, ok := o.(iseq.Seqable); ok {
		return s.Seq() == nil
	}

	return false
}

// TODO: figure out a standard hash code for empty sequences?
var hashCode uint32 = 1337

func (c *EmptyList) Hash() uint32 {
	return hashCode
}
