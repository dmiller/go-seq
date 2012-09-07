package seqimpl

import (
	"seq"
	"seq/sequtils"
	)

// Cons implements an immutable cons cell (sort of CAR/CDR, from the beginning of Lisp)
type Cons struct {
	AMeta
	first interface{}
	more seq.Seq
}

// Cons needs to implement the seq interfaces: 
//    Obj, Meta, Seq, Sequential, PersistentCollection, Seqable, IHashEq
//   Also, Equatable and Hashable
//
// I'm not sure yet if I'll be doing IHashEq
// Also, Sequential is a marker interface that I haven't figured out how to translate
//    because I can't figure out a significant use of it in the Clojure code

// interface Meta is covered by the AMeta embedding

// interface seq.Obj

func (c *Cons) WithMeta(meta seq.PersistentMap) seq.Obj {
	nc := &Cons{first: c.first, more: c.more}
	nc.meta = meta
	return nc	
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
	return CachedEmptyList;
}

func (c *Cons) Equiv(o interface{}) bool {
 	if c == o {
 		return true
 	}
	
 	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
 	os, ok := o.(seq.Seqable)

 	if !ok {
 		return false
 	}

	return sequtils.SeqEquiv(c,os.Seq());
}

// interface seq.Seq

func (c *Cons) 	First() interface{} {
	return c.first
}

func (c *Cons) 	Next() seq.Seq {
	return c.More().Seq()
}

func (c *Cons) 	More() seq.Seq {
	if c.more == nil {
		return CachedEmptyList;
	}

	return c.more;
}

func (c *Cons) 	SCons(o interface{}) seq.Seq {
	return &Cons{first: o, more: c} 
}
