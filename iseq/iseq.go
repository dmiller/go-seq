// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package iseq contains the public interfaces for the Clojure sequence library (ported to #golang)
package iseq

import ()

// Equivable is the interface for types that support testing for equivalence.
type Equivable interface {

	// Returns true if this equivalent to the given object.
	Equiv(o interface{}) bool
}

// Hashable is the interface for types that support computing a hash code.
// Two items that are Equiv should have the same hash code.
type Hashable interface {

	// Hashable extends Equivable
	Equivable

	// Returns a hash code for the thing.
	Hash() uint32
}

// Sequable is the interface for types that can produce an iseq.Seq
type Seqable interface {

	// Returns a sequence that can be iterated across.
	Seq() Seq
}

// PCollection is the most basic interface for types that implement
// an immutable, persistent collection.
// A PCollection is Seqable.
// To access the elements of the PCollection c, you must call c.seq().
type PCollection interface {
	Seqable
	Equivable

	// Returns the number of items in the collection.
	// Not guaranteed to be O(1).
	Count() int

	// Returns a new PCollection with an element added to this collection.
	// Which end the element is added to is collection-specific.
	Cons(o interface{}) PCollection

	// Returns an empty PCollection of the same type (if possible).
	Empty() PCollection
}

// Seq is the interface for sequential access to a collection.
// A Seq is itself a PCollection.
// A non-null Seq has at least one element.
type Seq interface {
	PCollection

	// Returns the first item in the sequence.
	// Calls seq on its argument.
	// If coll is nil, returns nil.
	First() interface{}

	// Returns a seq of the items after the first.
	// Calls seq on its argument.
	// If there are no more items, returns nil.
	Next() Seq

	// Returns a possibly empty seq of the items after the first.
	// Calls seq on its argument.
	More() Seq

	// Returns a new Seq with o added.  Type-specific version of PCollection.Cons.
	ConsS(o interface{}) Seq
}

// The Lookup interface supports looking up a value by key.
type Lookup interface {

	// Returns the value associated with the given key, or nil if the key is not present.
	ValAt(key interface{}) interface{}

	// Returns the value associated with the given key, or the provided default value if the key is not present.
	ValAtD(key interface{}, notFound interface{}) interface{}
}

// A MapEntry is an immutable key/value pair.
type MapEntry interface {

	// The key
	Key() interface{}

	// The value
	Val() interface{}
}

// Counted is the interface implemented by a collection to indicate it provides a constant-time count method.
type Counted interface {

	// Returns the number of elements in the collection, in constant time.
	Count1() int
}

// An Associative is a persistent, immutable collection supporting key/value lookup.
type Associative interface {
	PCollection
	Lookup

	// Returns true if there is an entry for the given key, false otherwise.
	ContainsKey(key interface{}) bool

	// Returns a MapEntry with the key/val for the given key, nil if key is not present.
	EntryAt(key interface{}) MapEntry

	// Returns a new Associative with key associated with value.
	Assoc(key interface{}, val interface{}) Associative
}

// A PMap is a persistent, immutable map (key/value) collection.
type PMap interface {
	Associative
	Counted

	// Returns a (new) PMap with key/val added.
	// Type-specific version of Associative.Assoc.
	AssocM(key interface{}, val interface{}) PMap

	// Returns a (possibly new) PMap with no entry for key.
	Without(key interface{}) PMap

	// Returns a (new) PMap with the key/val added.
	// Type-specific version of PCollection.Cons.
	ConsM(e MapEntry) PMap
}

// A Meta represents an object that can have metadata attached.
type Meta interface {

	// Returns the attached metadata
	Meta() PMap
}

// An MetaW is a Meta that also supports creating a copy with new metadata
// This was originally clojure.lang.Obj
type MetaW interface {
	Meta
	WithMeta(meta PMap) MetaW
}

// An Indexed collection supports direct access to the n-th item in the collection.
type Indexed interface {
	Counted

	// Returns the i-th entry, or nil if i is out of bounds.
	Nth(i int) interface{}

	// Returns the i-th entry, or a default value if i is out of bounds.
	NthD(i int, notFound interface{}) interface{}

	// Returns (i-th entry,nil), or (nil,error) if i is out of bounds.
	NthE(i int) (interface{}, error)
}

// An IndexedSeq is a Counted Seq that has a notion of the index of its first element in the context of a parent collection.
type IndexedSeq interface {
	Seq
	Counted

	// Returns the index of the first element of this seq, relative to its parent.
	Index() int
}

// A Reversible is a collection that supports iterating through its items in reverse order.
type Reversible interface {

	// Returns a Seq which is the reverse.
	Rseq() Seq
}

// A PStackOps has stack operations of peek and pop.
type PStackOps interface {

	// Returns the element on top of the stack or nil if empty
	Peek() interface{}

	// Returns a new stack with the top element removed.
	Pop() PStack
}

// A PStack is a persistent, immutable collection support stack operations.
type PStack interface {
	PCollection
	PStackOps
}

// A PList is a persistent, immutable list (a PCollection with stack operations)
type PList interface {
	PCollection
	PStackOps
}

// A PVector is a persistent, immutable vector.
// A PVector is a PCollection that supports lookup by an int index and supports stack operations.
type PVector interface {
	Associative
	PStackOps
	Reversible
	Indexed

	ConsV(interface{}) PVector
	AssocN(i int, val interface{}) PVector
}

// A PSet is a persistent, immutable collection of unique elements.
type PSet interface {
	PCollection
	Counted
	Disjoin(key interface{}) PSet
	Contains(key interface{}) bool
	// Get ?? do we need ??

}

// A Chunk is used internally to efficiently sequence through collections.
type Chunk interface {
	Indexed
	DropFirst() Chunk
}

// A Comparer supports comparing itself to other objects.
type Comparer interface {
	Compare(y interface{}) int
}

// A CompareFn compares two objects for <, = , >
type CompareFn func(interface{}, interface{}) int

// A Sorted collection maintains its entry in sorted order given by a comparator function.
type Sorted interface {
	Comparator() CompareFn
	EntryKey(entry interface{}) interface{}
	SeqA(ascending bool) Seq
	SeqFrom(key interface{}, ascending bool) Seq
}
