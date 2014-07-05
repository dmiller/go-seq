// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
)

// EmptyList implements an empty iseq.PList
// Used as a canonical empty list for other types that cannot represent empty.
type EmptyList struct {
	AMeta
}

var (
	// CachedEmptyList is a cached EmptyList.  No need to have more than one.
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

// WithMeta returns a new empty list with the given metadata attached
func (e *EmptyList) WithMeta(meta iseq.PMap) iseq.MetaW {
	return &EmptyList{AMeta: AMeta{meta}}
}

// interface iseq.Seqable

// Seq returns nil (since an EmptyList is an iseq.Seq with no values)
func (e *EmptyList) Seq() iseq.Seq {
	return nil
}

// interface iseq.PCollection

// Count returns 0 (the number of elements in an EmptyList)
func (e *EmptyList) Count() int {
	return 0
}

// Cons returns a PList whose one element is the given item
func (e *EmptyList) Cons(o interface{}) iseq.PCollection {
	return e.ConsS(o)
}

// Empty returns an EmptyList (namely, this EmptyList itself)
func (e *EmptyList) Empty() iseq.PCollection {
	return e
}

// interface iseq.Seq

// First returns nil (no first element in an EmptyList)
func (e *EmptyList) First() interface{} {
	return nil
}

// Next returns nil (no remaining element in an EmptyList)
func (e *EmptyList) Next() iseq.Seq {
	return nil
}

// More returns this EmptyList itself
func (e *EmptyList) More() iseq.Seq {
	return e
}

// ConsS returns a PList whose one element is the given item
func (e *EmptyList) ConsS(o interface{}) iseq.Seq {
	return NewPList1(o)
}

// interface Counted

// Count1 returns 0 (the length of an EmptyList)
func (e *EmptyList) Count1() int {
	return 0
}

// PStack

// Peek returns nil (no first element in an EmptyList)
func (e *EmptyList) Peek() interface{} {
	return nil
}

// Pop returns nil (nothing left in an EmptyList)
func (e *EmptyList) Pop() iseq.PStack {
	// in Clojure, popping throws an exception
	// should we add another return value?
	// For the moment, just return nil
	return nil
}

// interfaces Equivable, Hashable

// Equiv checks if the argument is an empty sequence.
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

// Hash computes a hash code for an EmptyList (all the same)
func (e *EmptyList) Hash() uint32 {
	return hashCode
}
