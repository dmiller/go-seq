// Package iseq contains the public interfaces for the Clojure sequence library
package iseq

import (
	"hash"
)

type Equatable interface {
	Equals(o interface{}) bool
}

type Hashable interface {
	Equatable
	Hash() uint32
	AddHash(hash.Hash)
}

type Seqable interface {
	Seq() Seq
}

type PersistentCollection interface {
	Seqable
	Count() int
	Cons(o interface{}) PersistentCollection
	Empty() PersistentCollection
	Equiv(o interface{}) bool
}

type Seq interface {
	PersistentCollection
	First() interface{}
	Next() Seq
	More() Seq
	SCons(o interface{}) Seq
}

type Lookup interface {
	ValAt(key interface{}) interface{}
	ValAtD(key interface{}, notFound interface{}) interface{}
}

type MapEntry interface {
	Key() interface{}
	Val() interface{}
}

type Counted interface {
	Count1() int
}

type Associative interface {
	PersistentCollection
	Lookup
	ContainsKey(key interface{}) bool
	EntryAt(key interface{}) MapEntry
	Assoc(key interface{}, val interface{}) Associative
}

type PersistentMap interface {
	Associative
	Counted
	AssocM(key interface{}, val interface{}) PersistentMap
	AssocEx(key interface{}, val interface{}) (result PersistentMap, ok bool)
	Without(key interface{}) (result PersistentMap, same bool)
	ConsM(e MapEntry) PersistentMap
}

type Meta interface {
	Meta() PersistentMap
}

type Obj interface {
	Meta
	WithMeta(meta PersistentMap) Obj
}

type Indexed interface {
	Counted
	Nth(i int) interface{}
	NthD(i int, notFound interface{}) interface{}
}

type IndexedSeq interface {
	Seq
	Counted
	Index() int
}

type Reversible interface {
	Rseq() Seq
}

type PersistentStack interface {
	Peek() interface{}
	Pop() PersistentStack
}

type PersistentList interface {
	PersistentStack
	PersistentCollection
}

type PersistentVector interface {
	Associative
	PersistentStack
	Reversible
	Indexed
	ConsV(interface{}) PersistentVector
	AssocN(i int, val interface{}) PersistentVector
}

type PersistentSet interface {
	PersistentCollection
	Counted
	Disjoin(key interface{}) PersistentSet
	Contains(key interface{}) bool
	// Get ?? do we need ??

}

type Chunk interface {
	Indexed
	DropFirst() Chunk
}
