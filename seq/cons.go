// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
	"hash"
)

// Cons implements an immutable cons cell (think CAR/CDR, for the old-timers)
//
// Zero value is (nil) (list of one element, nil), nil metadata, hash not cached.
// We use 0 as an indicator that hash is not cached.
type Cons struct {
	first interface{}
	more  iseq.Seq
	AMeta
	hash uint32
}

// Cons needs to implement the iseq interfaces: 
//    Obj, Meta, Seq, Sequential, PCollection, Seqable, IHashEq
//   Also, Equatable and Hashable
//
// I'm not sure yet if I'll be doing IHashEq
// Also, Sequential is a marker interface that I haven't figured out how to translate
//    because I can't figure out a significant use of it in the Clojure code

// interface Meta is covered by the AMeta embedding

// c-tors

func NewCons(first interface{}, more iseq.Seq) *Cons {
	return &Cons{first: first, more: more}
}

func NewConsM(meta iseq.PMap, first interface{}, more iseq.Seq) *Cons {
	nc := &Cons{AMeta: AMeta{meta}, first: first, more: more}
	return nc
}

// interface iseq.Obj

func (c *Cons) WithMeta(meta iseq.PMap) iseq.Obj {
	return NewConsM(meta, c.first, c.more)
}

// interface iseq.Seqable

func (c *Cons) Seq() iseq.Seq {
	return c
}

// interface iseq.PCollection

func (c *Cons) Count() int {
	return 1 + sequtil.Count(c.more)
}

func (c *Cons) Cons(o interface{}) iseq.PCollection {
	return c.SCons(o)
}

func (c *Cons) Empty() iseq.PCollection {
	return CachedEmptyList
}

func (c *Cons) Equiv(o interface{}) bool {
	if c == o {
		return true
	}

	if os, ok := o.(iseq.Seqable); ok {
		return sequtil.SeqEquiv(c, os.Seq())
	}

	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
	return false
}

// interface iseq.Seq

func (c *Cons) First() interface{} {
	return c.first
}

func (c *Cons) Next() iseq.Seq {
	return c.More().Seq()
}

func (c *Cons) More() iseq.Seq {
	if c.more == nil {
		return CachedEmptyList
	}

	return c.more
}

func (c *Cons) SCons(o interface{}) iseq.Seq {
	return &Cons{first: o, more: c}
}

// interfaces Equatable, Hashable

func (c *Cons) Equals(o interface{}) bool {
	if c == o {
		return true
	}

	if os, ok := o.(iseq.Seqable); ok {
		return sequtil.SeqEquals(c, os.Seq())
	}

	// TODO: handle built-in 'sequable' things such as arrays, slices, strings
	return false
}

func (c *Cons) Hash() uint32 {
	if c.hash == 0 {
		c.hash = sequtil.HashSeq(c)
	}

	return c.hash
}

func (c *Cons) AddHash(h hash.Hash) {
	if c.hash == 0 {
		c.hash = sequtil.HashSeq(c)
	}

	sequtil.AddHashUint64(h, uint64(c.hash))
}
