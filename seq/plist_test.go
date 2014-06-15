// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
	"testing"
)

func TestNewPlist1(t *testing.T) {
	p := NewPList1("abc")

	if p.Meta() != nil {
		t.Error("NewPlist1 ctor should have nil meta")
	}

	if p.First() != "abc" {
		t.Error("NewPlist1 ctor did not initialize first")
	}

	if p.Next() != nil {
		t.Error("NewPlist1 ctor did not initialize next to nil")
	}

	if p.Count() != 1 {
		t.Error("NewPlist1 ctor did not initialize count properly")
	}
}

func TestNewPListFromSlice(t *testing.T) {
	base := []interface{}{"abc", "def", 2, 3}
	p := NewPListFromSlice(base)

	if p.Meta() != nil {
		t.Error("NewPListFromSlice ctor should have nil meta")
	}

	if p.Count() != 4 {
		t.Error("NewPListFromSlice ctor did not initialize count properly")
	}

	if p.First() != "abc" {
		t.Error("NewPListFromSlice ctor did not initialize first")
	}

	if p.Next().First() != "def" {
		t.Error("NewPListFromSlice ctor did not initialize next to correct tail")
	}
}

func TestNewPlist1N(t *testing.T) {
	tail := NewPListFromSlice([]interface{}{"def", 2, 3})

	p := NewPList1N("abc", tail, 3)

	if p.Meta() != nil {
		t.Error("NewPlist1N ctor should have nil meta")
	}

	if p.First() != "abc" {
		t.Error("NewPlist1N ctor did not initialize first")
	}

	if p.Next() != tail {
		t.Error("NewPlist1N ctor did not initialize next to correct tail")
	}

	if p.Count() != 3 {
		t.Error("NewPlist1N ctor did not initialize count properly")
	}
}

// TODO: add tests for c-tor with meta -- we need a PMap implementation first

func TestPListImplementInterfaces(t *testing.T) {
	var c interface{} = NewPList1("abc")

	if _, ok := c.(iseq.MetaW); !ok {
		t.Error("PList must implement MetaW")
	}

	if _, ok := c.(iseq.Meta); !ok {
		t.Error("PList must implement Meta")
	}

	if _, ok := c.(iseq.PCollection); !ok {
		t.Error("PList must implement PCollection")
	}

	if _, ok := c.(iseq.PStack); !ok {
		t.Error("PList must implement PStack")
	}

	if _, ok := c.(iseq.PList); !ok {
		t.Error("PList must implement PList")
	}

	if _, ok := c.(iseq.Seqable); !ok {
		t.Error("PList must implement Seqable")
	}

	if _, ok := c.(iseq.Counted); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Equivable); !ok {
		t.Error("PList must implement Equatable")
	}

	if _, ok := c.(iseq.Hashable); !ok {
		t.Error("PList must implement Hashable")
	}
}

func createLongerPList() *PList {
	return NewPListFromSlice([]interface{}{"def", createComplicatedCons(), 3})
}

func TestPListCount(t *testing.T) {
	c := createLongerPList()
	if c.Count() != 3 {
		t.Errorf("Count: expected 3, got %v", c.Count())
	}

	if c.Count1() != 3 {
		t.Errorf("Count1: expected 3, got %v", c.Count1())
	}
}

func TestPListSeq(t *testing.T) {

	sl := []interface{}{"abc", "def", 2, 3}
	pl := NewPListFromSlice(sl)

	if pl.Seq() != pl {
		t.Error("Seq should return self")
	}

	i := 0
	for s := pl.Seq(); s != nil; s, i = s.Next(), i+1 {
		if f, e := s.First(), sl[i]; f != e {
			t.Errorf("for Seq, on element %v, expected %v, got %v", i, f, e)
		}
	}
}

func TestPListCons(t *testing.T) {
	c1 := NewPList1("abc")
	c2 := c1.ConsS("def")

	if c2.First() != "def" {
		t.Error("Cons has a bad first element")
	}

	if c2.Next() != c1 {
		t.Error("Cons has a bad rest")
	}

	// TODO: test preservation of meta when we have a PMap implementation
}

// TODO: test Seq has meta

func TestPListEmpty(t *testing.T) {
	c := NewPList1("abc")
	e := c.Empty()
	if e.Count() != 0 {
		t.Error("Empty returns a non-empty list")
	}

	// TODO: test preservation of meta when we have a PMap implementation
}

func TestPListEquiv(t *testing.T) {
	c1 := createLongerPList()
	c2 := createLongerPList()
	if c1 == c2 {
		t.Error("Expect two calls to createLongerPList to return distinct structs")
	}
	if !c1.Equiv(c1) {
		t.Error("Expect cons to be equiv to itself")
	}
	if !c1.Equiv(c2) {
		t.Error("Expect cons to equiv similar cons")
	}

	c3 := NewCons("abc", nil)
	if c1.Equiv(c3) {
		t.Error("cons equiv dissimilar list")
	}
}

func TestPListZeroValue(t *testing.T) {
	c1 := new(PList)
	if !c1.Equiv(CachedEmptyList) {
		t.Error("Zero-value PList should be equiv to an EmptyList")
	}
}
