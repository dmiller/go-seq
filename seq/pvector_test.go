// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	//"fmt"
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
	"testing"
)

//  PVector needs to implement the following seq interfaces:
//        MetaW Meta Seqable PCollection Lookup Associative PStack PVector Counted Reversible Indexed
//  Are we going to do EditableCollection?
//  Also, Equatable and Hashable
func TestPVectorImplementInterfaces(t *testing.T) {
	var c interface{} = NewPVectorFromItems("abc", "def")

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

	if _, ok := c.(iseq.Lookup); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Associative); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.PVector); !ok {
		t.Error("PList must implement PList")
	}

	if _, ok := c.(iseq.Seqable); !ok {
		t.Error("PList must implement Seqable")
	}

	if _, ok := c.(iseq.Counted); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Indexed); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Reversible); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Equivable); !ok {
		t.Error("PList must implement Equatable")
	}

	if _, ok := c.(iseq.Hashable); !ok {
		t.Error("PList must implement Hashable")
	}
}

// Factory tests

func TestPVectorISeqFactoryWorks(t *testing.T) {
	var seq iseq.Seq = NewPListFromSlice([]interface{}{"def", 2, 3})
	v := NewPVectorFromISeq(seq)
	if v.Count() != 3 {
		t.Errorf("NewPVectorFromISeq has wrong count, expected %v, got %v", 3, v.Count())
	}
	for i, s := 0, seq; i < v.Count(); i, s = i+1, s.Next() {
		if v.Nth(i) != s.First() {
			t.Errorf("NewPVectorFromISeq: expected element %v = %v, found %v instead", i, v.Nth(i), s.First())
		}
	}
}

func makeRangePlist(size int) (*PList, []interface{}) {
	sl := make([]interface{}, size)
	for i := 0; i < len(sl); i++ {
		sl[i] = i + 10
	}

	return NewPListFromSlice(sl), sl
}

func makeRangePVector(size int) (*PVector, []interface{}) {
	sl := make([]interface{}, size)
	for i := 0; i < len(sl); i++ {
		sl[i] = i + 10
	}

	return NewPVectorFromSlice(sl), sl
}

func TestPVectorISeqFactoryWorksOnLargeSeq(t *testing.T) {
	sizes := []int{
		1000,   // this should get us out of the first node
		100000} // this should get us out of the second level

	for is := 0; is < len(sizes); is++ {
		size := sizes[is]
		pl, _ := makeRangePlist(size)

		v := NewPVectorFromISeq(pl)
		if v.Count() != size {
			t.Errorf("NewPVectorFromISeq has wrong count, expected %v, got %v", size, v.Count())
		}
		for i, s := 0, pl.Seq(); i < v.Count(); i, s = i+1, s.Next() {
			if v.Nth(i) != s.First() {
				t.Errorf("NewPVectorFromISeq: expected element %v = %v, found %v instead", i, v.Nth(i), s.First())
			}
		}
	}
}

func TestPVectorSliceFactoryWorks(t *testing.T) {
	sl := []interface{}{"def", 2, 3}
	v := NewPVectorFromSlice(sl)
	if v.Count() != 3 {
		t.Errorf("NewPVectorFromSlice has wrong count, expected %v, got %v", 3, v.Count())
	}
	for i := 0; i < v.Count(); i = i + 1 {
		if v.Nth(i) != sl[i] {
			t.Errorf("NewPVectorFromSlice: expected element %v = %v, found %v instead", i, v.Nth(i), sl[i])
		}
	}
}

func TestPVectorFromItemsWorks(t *testing.T) {
	sl := []interface{}{"def", 2, 3}
	v := NewPVectorFromItems("def", 2, 3)
	if v.Count() != 3 {
		t.Errorf("NewPVectorFromSlice has wrong count, expected %v, got %v", 3, v.Count())
	}
	for i := 0; i < v.Count(); i = i + 1 {
		if v.Nth(i) != sl[i] {
			t.Errorf("NewPVectorFromSlice: expected element %v = %v, found %v instead", i, v.Nth(i), sl[i])
		}
	}
}

// Equality tests

func TestPVectorEquals(t *testing.T) {
	sl := []interface{}{"def", 2, 3}
	v := NewPVectorFromSlice(sl)
	if v.Equals(sl) {
		t.Error("PVector should not Equals a non-PVector")
	}

	v1 := NewPVectorFromSlice(sl)
	sl2 := []interface{}{"def", 2, 4}
	v2 := NewPVectorFromSlice(sl2)

	if !v.Equals(v1) {
		t.Error("PVector should equal equal PVector")
	}

	if v.Equals(v2) {
		t.Error("PVector should not equal non-equal PVector")
	}

	var seq iseq.Seq = NewPListFromSlice(sl)
	if !v.Equals(seq) {
		t.Error("PVector should equal equivalent iseq.Seq")
	}

	v0 := NewPVectorFromItems()
	if v0.Count() != 0 {
		t.Error("PVector of no items should have zero count")
	}
}

// iseq.Meta, iseq.MetaW tests

// TODO: test WithMeta once we have PMap

// iseq.Sequable tests

func TestPVectorSeqable(t *testing.T) {
	s0 := NewPVectorFromItems().Seq()
	if s0 != nil {
		t.Error("PVectorSeq on zero count should be nil")
	}

	s3 := NewPVectorFromItems("abc", 4, 5).Seq()
	if s3.Count() != 3 {
		t.Errorf("PVectorSeq on non-empty items should have count %v, got %v", s3.Count())
	}
	if s3.First() != "abc" {
		t.Error("PVectorSeq first item wrong")
	}
	if s3.Next().First() != 4 {
		t.Error("PVectorSeq second item wrong")
	}
	if s3.Next().Next().First() != 5 {
		t.Error("PVectorSeq third item wrong")
	}
	if s3.Next().Next().Next() != nil {
		t.Error("PVectorSeq should have nil after last item")
	}
}

// iseq.PCollection tests

func TestPVectorCons(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)
	v4 := v3.ConsV(12)

	c4 := v3.Cons(12)

	if c4.Count() != 4 {
		t.Error("Cons of PVector has wrong count")
	}
	vc4, vok := c4.(*PVector)
	if !vok {
		t.Errorf("Cons of PVector should be a PVector, got %T", vc4)
	}

	for i := 0; i < c4.Count(); i++ {
		if vc4.Nth(i) != v4.Nth(i) {
			t.Error("Cons of PVector has wrong item at %v", i)
		}
	}

	if !sequtil.SeqEquiv(c4.Seq(), v4.Seq()) {
		t.Error("Something wrong with Seq()ing on PVector")
	}
}

func TestPVectorEmpty(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)
	e := v3.Empty()

	if e.Count() != 0 {
		t.Error("Expected PVector.Empty() to have 0 count")
	}

	// TODO: test Empty copies Meta once we have PMap
}

// PVector.Seq() tests

func TestPVectorConses(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)

	s4 := v3.Seq().ConsS(12)

	if s4.Count() != 4 {
		t.Error("SCons of PVector.Seq has wrong count")
	}

	if s4.First() != 12 {
		t.Error("SCons of PVector.Seq has wrong first item")
	}
	if !sequtil.SeqEquiv(s4.Next(), v3.Seq()) {
		t.Error("SCons of PVector.Seq has wrong next seq")
	}
}

// iseq.PVector tests

func TestPVectorConsV(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)

	v4 := v3.ConsV(12)
	if v4.Count() != 4 {
		t.Error("ConsV of PVector has wrong count")
	}
	slice := []interface{}{"abc", 4, 5, 12}
	for i := 0; i < len(slice); i++ {
		if v4.Nth(i) != slice[i] {
			t.Error("ConsV of PVector has wrong item at %v", i)
		}
	}
}

func TestPVectorAssocN(t *testing.T) {
	v3o := NewPVectorFromItems("abc", 4, 5)
	v3 := NewPVectorFromItems("abc", 4, 5)
	v3a := v3.AssocN(1, "def")
	if !sequtil.SeqEquiv(v3o.Seq(), v3.Seq()) {
		t.Error("PVector.AssocN: appears to have mutated original")
	}
	if v3a.Count() != v3.Count() {
		t.Errorf("PVector.AssocN: count should be %v, got %v", v3.Count(), v3a.Count())
	}
	v3t := NewPVectorFromItems("abc", "def", 5)
	if !sequtil.SeqEquiv(v3a.Seq(), v3t.Seq()) {
		t.Errorf("PVector.AssocN: wrong items")
	}

	v4a := v3.AssocN(3, "pqr")
	if !sequtil.SeqEquiv(v3o.Seq(), v3.Seq()) {
		t.Error("PVector.AssocN: appears to have mutated original")
	}
	if v4a.Count() != v3.Count()+1 {
		t.Errorf("PVector.AssocN: appending count should be %v, got %v", v3.Count()+1, v3a.Count())
	}
	v4t := NewPVectorFromItems("abc", 4, 5, "pqr")
	if !sequtil.SeqEquiv(v4a.Seq(), v4t.Seq()) {
		t.Errorf("PVector.AssocN: wrong items on append")
	}

	v1a := NewPVectorFromItems().AssocN(0, "abc")
	if c := v1a.Count(); c != 1 {
		t.Errorf("PVector.AssocN: appending to empty should have one element, found %v", c)
	}

	if e := v1a.Nth(0); e != "abc" {
		t.Errorf("PVector.AssocN: appending to empty, first (only) element should be abc, got %v", e)
	}
}

func TestPVectorAssocNBadIndexLow(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)

	defer func() {
		if r := recover(); r == nil {
			t.Error("PVector.AssocN: expected panic on out-of-bounds index, but it executed normally")
		}
	}()

	v3.Assoc(-10, "def")
}

func TestPVectorAssocNBadIndexHigh(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)

	defer func() {
		if r := recover(); r == nil {
			t.Error("PVector.AssocN: expected panic on out-of-bounds index, but it executed normally")
		}
	}()

	v3.Assoc(10, "def")
}

func TestPVectorAssocNBig(t *testing.T) {
	v, _ := makeRangePVector(100000)
	var nv iseq.PVector = v
	for i := 0; i < 110000; i++ {
		nv = nv.AssocN(i, i+20)
	}

	if nv.Count() != 110000 {
		t.Error("PVector.AssocN: lots of appends at end of big vector: wrong number of items")
	}
	for i := 0; i < nv.Count(); i++ {
		if nv.Nth(i) != i+20 {
			t.Errorf("PVector.AssocN: lots of appends at end of big vector: bad value at %v", i)
		}
	}
}

// iseq.Counted tests

func TestPVectorCount1(t *testing.T) {
	v, _ := makeRangePVector(1000)
	if v.Count1() != 1000 {
		t.Errorf("PVector.Count1: expected %v, found %v", 1000, v.Count1())
	}
}

// iseq.Lookup, iseq.Associative tests

func TestPVectorValAtEtc(t *testing.T) {
	v, sl := makeRangePVector(1000)

	for i := 0; i < len(sl); i++ {
		if sl[i] != v.ValAt(i) {
			t.Errorf("PVector.ValAt(%v) = %v, expected %v", i, v.Nth(i), sl[i])
		}
		if sl[i] != v.ValAtD(i, "a") {
			t.Errorf("PVector.ValAt(%v) = %v, expected %v", i, v.Nth(i), sl[i])
		}
		if !v.ContainsKey(i) {
			t.Errorf("PVector.ContainsKey(%v) returned false, expected true", i)
		}
		if me := v.EntryAt(i); me.Key() != i || me.Val() != sl[i] {
			t.Errorf("PVector.EntryAt(%v) = <%v, %v>, expected <%v, %v>", i, me.Key(), me.Val(), i, sl[i])
		}
	}

	if v.ValAt(-100) != nil {
		t.Error("PVector.ValAt on out-of-range index should return nil")
	}

	if v.ValAtD(-100, "a") != "a" {
		t.Error("PVector.ValAtD on out-of-range index should return default")
	}

	if v.ContainsKey(-100) {
		t.Error("PVector.ContainsKey on out-of-range index should be false")
	}

	if v.EntryAt(-100) != nil {
		t.Error("PVector.EntryAt on out-of-range index should return nil")
	}

	if v.ValAt(100000) != nil {
		t.Error("PVector.ValAt on out-of-range index should return nil")
	}

	if v.ValAtD(100000, "a") != "a" {
		t.Error("PVector.ValAtD on out-of-range index should return default")
	}

	if v.ContainsKey(100000) {
		t.Error("PVector.ContainsKey on out-of-range index should be false")
	}

	if v.EntryAt(100000) != nil {
		t.Error("PVector.EntryAt on out-of-range index should return nil")
	}

	if v.ValAt("b") != nil {
		t.Error("PVector.ValAt on non-numeric index should return nil")
	}

	if v.ValAtD("b", "a") != "a" {
		t.Error("PVector.ValAtD on non-numeric index should return default")
	}

	if v.ContainsKey("b") {
		t.Error("PVector.ContainsKey on out-of-range index should be false")
	}

	if v.EntryAt("b") != nil {
		t.Error("PVector.EntryAt on non-numeric index should return nil")
	}
}

func TestPVectorAssoc(t *testing.T) {
	// Not black box here.  Assuming Assoc delegates to AssocN on numeric key
	// Only testing non-numeric key handling

	v3 := NewPVectorFromItems("abc", 4, 5)
	defer func() {
		if r := recover(); r == nil {
			t.Error("PVector.Assoc: expected panic on non-numeric key, but it executed normally")
		}
	}()

	v3.Assoc("a", 12)
}

// iseq.Indexed tests

func TestPVectorIndexed(t *testing.T) {
	v, sl := makeRangePVector(1000)

	for i := 0; i < len(sl); i++ {
		if sl[i] != v.Nth(i) {
			t.Errorf("PVector.Nth(%v) = %v, expected %v", i, v.Nth(i), sl[i])
		}
		if sl[i] != v.NthD(i, "a") {
			t.Errorf("PVector.NthD(%v,'a') = %v, expected %v", i, v.NthD(i, "a"), sl[i])
		}
	}

	if v.NthD(-100, "a") != "a" {
		t.Error("PVector.NthD: should return default on out-of-bounds index")
	}

	if v.NthD(100000, "a") != "a" {
		t.Error("PVector.NthD: should return default on out-of-bounds index")
	}
}

func TestPVectorIndexedBadIndexLow(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)

	defer func() {
		if r := recover(); r == nil {
			t.Error("PVector.Nth: expected panic on out-of-bounds index, but it executed normally")
		}
	}()

	v3.Nth(-10)
}

func TestPVectorIndexedBadIndexHigh(t *testing.T) {
	v3 := NewPVectorFromItems("abc", 4, 5)

	defer func() {
		if r := recover(); r == nil {
			t.Error("PVector.Nth: expected panic on out-of-bounds index, but it executed normally")
		}
	}()

	v3.Nth(10)
}

// iseq.PStack tests

func TestPVectorPeek(t *testing.T) {
	if EmptyPVector.Peek() != nil {
		t.Error("PVector.Peek: expected Peek of empty to be nil")
	}

	v3 := NewPVectorFromItems("abc", 4, 5)
	if v3.Peek() != 5 {
		t.Error("PVector.Peek: expected Peek to return last item")
	}
}

func TestPVectorPop(t *testing.T) {

	v1 := NewPVectorFromItems("a")
	if v1.Pop().Count() != 0 {
		t.Error("PVector.Pop: expected Pop to return collection of count 0")
	}

	v, _ := makeRangePVector(100000)
	var s iseq.PStack = v
	for i := 16; i < 100000; i++ {
		s = s.Pop()
	}

	vs, ok := s.(iseq.PVector)
	if !ok {
		t.Errorf("PVector.Pop: expected pop result to be PVector, got a %T", s)
	}

	if cnt := vs.Count(); cnt != 16 {
		t.Errorf("PVector.Pop: after many pops, expected 16 count, got %v", cnt)
	}
	for i := 0; i < vs.Count1(); i++ {
		if vs.Nth(i) != v.Nth(i) {
			t.Errorf("PVector.Pop: expected vs[%v] = %v, found %v", i, v.Nth(i), vs.Nth(i))
		}
	}
}

func TestPVectorPopOnEmpty(t *testing.T) {
	v0 := NewPVectorFromItems()

	defer func() {
		if r := recover(); r == nil {
			t.Error("PVector.Pop: expected panic on Pop of empty vector, but it was okay")
		}
	}()

	v0.Pop()
}

// TODO: Add Rseq tests
