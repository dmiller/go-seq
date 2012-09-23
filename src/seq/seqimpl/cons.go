package seqimpl

import (
	"hash"
	"seq"
	"seq/sequtils"
)

// Cons implements an immutable cons cell (sort of CAR/CDR, from the beginning of Lisp)
//
// Zero value is (nil) (list of one element, nil), nil metadata, hash not cached.
// We use 0 as an indicator that hash is not cached.
type Cons struct {
	first interface{}
	more  seq.Seq
	AMeta
	hash uint32
}

// Cons needs to implement the seq interfaces: 
//    Obj, Meta, Seq, Sequential, PersistentCollection, Seqable, IHashEq
//   Also, Equatable and Hashable
//
// I'm not sure yet if I'll be doing IHashEq
// Also, Sequential is a marker interface that I haven't figured out how to translate
//    because I can't figure out a significant use of it in the Clojure code

// interface Meta is covered by the AMeta embedding

// c-tors

func NewCons(first interface{}, more seq.Seq) *Cons {
	return &Cons{first: first, more: more}
}

func NewConsM(meta seq.PersistentMap, first interface{}, more seq.Seq) *Cons {
	nc := &Cons{AMeta: AMeta{meta}, first: first, more: more}
	return nc
}

// interface seq.Obj

func (c *Cons) WithMeta(meta seq.PersistentMap) seq.Obj {
	return NewConsM(meta, c.first, c.more)
}

// interface seq.Seqable

func (c *Cons) Seq() seq.Seq {
	return c
}

// interface seq.PersistentCollection

func (c *Cons) Count() int {
	return 1 + sequtils.Count(c.more)
}

func (c *Cons) Cons(o interface{}) seq.PersistentCollection {
	return c.SCons(o)
}

func (c *Cons) Empty() seq.PersistentCollection {
	return CachedEmptyList
}

func (c *Cons) Equiv(o interface{}) bool {
	if c == o {
		return true
	}

	if os, ok := o.(seq.Seqable); ok {
		return sequtils.SeqEquiv(c, os.Seq())
	}

	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
	return false
}

// interface seq.Seq

func (c *Cons) First() interface{} {
	return c.first
}

func (c *Cons) Next() seq.Seq {
	return c.More().Seq()
}

func (c *Cons) More() seq.Seq {
	if c.more == nil {
		return CachedEmptyList
	}

	return c.more
}

func (c *Cons) SCons(o interface{}) seq.Seq {
	return &Cons{first: o, more: c}
}

// interfaces Equatable, Hashable

func (c *Cons) Equals(o interface{}) bool {
	if c == o {
		return true
	}

	if os, ok := o.(seq.Seqable); ok {
		return sequtils.SeqEquals(c, os.Seq())
	}

	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
	return false
}

func (c *Cons) Hash() uint32 {
	if c.hash == 0 {
		c.hash = sequtils.HashSeq(c)
	}

	return c.hash
}

func (c *Cons) AddHash(h hash.Hash) {
	if c.hash == 0 {
		c.hash = sequtils.HashSeq(c)
	}

	sequtils.AddHashUint64(h, uint64(c.hash))
}
