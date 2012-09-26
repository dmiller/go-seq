package seq

import (
	"hash"
	"iseq"
	"sequtil"
)

// EmptyList implements an empty PersistentList
type EmptyList struct {
	AMeta
}

var (
	CachedEmptyList = &EmptyList{}
)

// Cons needs to implement the seq interfaces: 
//    Obj, Meta, Seq, Sequential, PersistentCollection, PersistentStack, PersistentList, Seqable, Counted
//   Also, Equatable and Hashable
//
// Sequential is a marker interface that I haven't figured out how to translate
//    because I can't figure out a significant use of it in the Clojure code
// PersistentList is just PersistentStack + PersistentCollection, so nothing added

// interface Meta is covered by the AMeta embedding

// interface iseq.Obj

func (e *EmptyList) WithMeta(meta iseq.PersistentMap) iseq.Obj {
	e2 := &EmptyList{}
	e2.meta = meta
	return e2
}

// interface iseq.Seqable

func (e *EmptyList) Seq() iseq.Seq {
	return nil
}

// interface iseq.PersistentCollection

func (e *EmptyList) Count() int {
	return 0
}

func (e *EmptyList) Cons(o interface{}) iseq.PersistentCollection {
	return e.SCons(o)
}

func (e *EmptyList) Empty() iseq.PersistentCollection {
	return e
}

func (e *EmptyList) Equiv(o interface{}) bool {
	if e == o {
		return true
	}

	// TODO: deal with other sequence types
	if s, ok := o.(iseq.Seqable); ok {
		return s.Seq() == nil
	}

	return false
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

func (e *EmptyList) SCons(o interface{}) iseq.Seq {
	// TODO: really, this needs to return a PersistentList of one element.
	// Fix when we have a true PersistentList
	return &Cons{first: o, more: e}
}

// interface Counted

func (e *EmptyList) Count1() int {
	return 0
}

// PersistentStack

func (e *EmptyList) Peek() interface{} {
	return nil
}

func (e *EmptyList) Pop() iseq.PersistentStack {
	// in Clojure, popping throws an exception
	// should we add another return value?
	// For the moment, just return nil
	return nil
}

// interfaces Equatable, Hashable

func (c *EmptyList) Equals(o interface{}) bool {
	return c.Equiv(o)
}

var hashCode uint32 = 37

func (c *EmptyList) Hash() uint32 {
	return hashCode
}

func (c *EmptyList) AddHash(h hash.Hash) {
	sequtil.AddHashUint64(h, uint64(hashCode))
}
