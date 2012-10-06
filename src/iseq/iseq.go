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

type PCollection interface {
	Seqable
	Count() int
	Cons(o interface{}) PCollection
	Empty() PCollection
	Equiv(o interface{}) bool
}

type Seq interface {
	PCollection
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
	PCollection
	Lookup
	ContainsKey(key interface{}) bool
	EntryAt(key interface{}) MapEntry
	Assoc(key interface{}, val interface{}) Associative
}

type PMap interface {
	Associative
	Counted
	AssocM(key interface{}, val interface{}) PMap
	AssocEx(key interface{}, val interface{}) (result PMap, ok bool)
	Without(key interface{}) (result PMap, same bool)
	ConsM(e MapEntry) PMap
}

type Meta interface {
	Meta() PMap
}

type Obj interface {
	Meta
	WithMeta(meta PMap) Obj
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

type PStack interface {
	Peek() interface{}
	Pop() PStack
}

type PList interface {
	PStack
	PCollection
}

type PVector interface {
	Associative
	PStack
	Reversible
	Indexed
	ConsV(interface{}) PVector
	AssocN(i int, val interface{}) PVector
}

type PSet interface {
	PCollection
	Counted
	Disjoin(key interface{}) PSet
	Contains(key interface{}) bool
	// Get ?? do we need ??

}

type Chunk interface {
	Indexed
	DropFirst() Chunk
}
