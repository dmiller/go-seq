// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
)

// Cons implements an immutable cons cell
// (Think CAR/CDR, for the old-timers)
//
// Zero value is (nil) (list of one element, namely, nil),
// nil metadata, hash not cached.
// We use 0 as an indicator that hash is not cached.
type Cons struct {
	first interface{}
	more  iseq.Seq
	AMeta
	hash uint32
}

// Cons needs to implement the iseq interfaces:
//   Meta, MetaW, Seq, Sequential, PCollection, Seqable
//   Also, Equivable and Hashable
//
// Sequential is a marker interface that I haven't figured out how to translate
//    because I can't figure out a significant use of it in the Clojure code
//
// interface Meta is covered by the AMeta embedding

// c-tors

// NewCons returns a new Cons cell
func NewCons(first interface{}, more iseq.Seq) *Cons {
	return &Cons{first: first, more: more}
}

// NewConsM returns a new Cons cell with metadata attached
func NewConsM(meta iseq.PMap, first interface{}, more iseq.Seq) *Cons {
	return &Cons{AMeta: AMeta{meta}, first: first, more: more}
}

// interface iseq.MetaW

// WithMeta returns a new Cons representing an iseq.PMap with new metadata attached
func (c *Cons) WithMeta(meta iseq.PMap) iseq.MetaW {
	return NewConsM(meta, c.first, c.more)
}

// interface iseq.Seqable

// Seq returns an iseq.Seq for this Cons (namely, itself)
func (c *Cons) Seq() iseq.Seq {
	return c
}

// interface iseq.PCollection

// Count returns the number of items in the Cons, viewed as a collection.
func (c *Cons) Count() int {
	return 1 + sequtil.Count(c.more)
}

// Cons a new Cons with the given item added to the front.
func (c *Cons) Cons(o interface{}) iseq.PCollection {
	return c.ConsS(o)
}

// Empty returns an empty collection.
// A Cons cannot be empty, so it cannot return the same type, as is preferred.
// Cons returns an EmptyList
func (c *Cons) Empty() iseq.PCollection {
	// A Cons cannot be empty, so we need to return something else.
	return CachedEmptyList
}

// interface iseq.Seq

// First returns the first item in the Cons, its head.
func (c *Cons) First() interface{} {
	return c.first
}

// Next returns the tail of the Cons, the seq of the cons after the first.
func (c *Cons) Next() iseq.Seq {
	return c.More().Seq()
}

// More returns a possibly empty seq of the items after the first.
func (c *Cons) More() iseq.Seq {
	if c.more == nil {
		return CachedEmptyList
	}

	return c.more
}

// ConsS returns a new Cons with a value added to the front.
func (c *Cons) ConsS(o interface{}) iseq.Seq {
	return &Cons{first: o, more: c}
}

// interfaces Equivable, Hashable

// Equiv returns true if this Cons is eqivalent to the given object, treated as an iseq.Seqable.
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

// Hash computes a Hash value for the Cons, treated as a sequence.
func (c *Cons) Hash() uint32 {
	if c.hash == 0 {
		c.hash = sequtil.HashSeq(c)
	}

	return c.hash
}
