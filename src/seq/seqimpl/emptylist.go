package seqimpl

import (
	"seq"
	"seq/sequtils"
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

// interface seq.Obj

func (e *EmptyList) WithMeta(meta seq.PersistentMap) seq.Obj {
	e2 := &EmptyList{}
	e2.meta = meta
	return e2	
}

// interface seq.Seqable

func (e *EmptyList) Seq() seq.Seq {
	return nil
}

// interface seq.PersistentCollection

func (e *EmptyList) Count() int {
	return 0
}

func (e *EmptyList) Cons(o interface{}) seq.PersistentCollection {
	return e.SCons(o)
}

func (e *EmptyList) Empty() seq.PersistentCollection {
	return e
}

func (e *EmptyList) Equiv(o interface{}) bool {
 	if e == o {
 		return true
 	}
	
	return sequtils.Equiv(e,o)
}

// interface seq.Seq

func (e *EmptyList) 	First() interface{} {
	return nil
}

func (e *EmptyList) 	Next() seq.Seq {
	return nil
}

func (e *EmptyList) 	More() seq.Seq {
	return e
}

func (e *EmptyList) 	SCons(o interface{}) seq.Seq {
	// TODO: really, this needs to return a PersistentList of one element.
	// Fix when we have a true PersistentList
	return &Cons{first: o, more: e} 
}


// interface Counted

func (e *EmptyList) implementsCounted() {
}


// PersistentStack

func (e *EmptyList)	Peek() interface{} {
	return nil
}

func (e *EmptyList)	Pop() seq.PersistentStack {
	// in Clojure, popping throws an exception
	// should we add another return value?
	// For the moment, just return nil
	return nil
}