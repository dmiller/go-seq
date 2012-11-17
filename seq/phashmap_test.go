// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"fmt"
	"github.com/dmiller/go-seq/iseq"
	"math/rand"
	"testing"
)

func TestPHashMapImplementInterfaces(t *testing.T) {
	var c interface{} = NewPHashMapFromItems("abc", "def")

	if _, ok := c.(iseq.Obj); !ok {
		t.Error("PHashMap must implement Obj")
	}

	if _, ok := c.(iseq.Meta); !ok {
		t.Error("PHashMap must implement Meta")
	}

	if _, ok := c.(iseq.PCollection); !ok {
		t.Error("PHashMap must implement PCollection")
	}

	if _, ok := c.(iseq.PMap); !ok {
		t.Error("PHashMap must implement PMap")
	}

	if _, ok := c.(iseq.Lookup); !ok {
		t.Error("PHashMap must implement Counted")
	}

	if _, ok := c.(iseq.Associative); !ok {
		t.Error("PHashMap must implement Counted")
	}

	if _, ok := c.(iseq.Seqable); !ok {
		t.Error("PHashMap must implement Seqable")
	}

	if _, ok := c.(iseq.Counted); !ok {
		t.Error("PHashMap must implement Counted")
	}

	if _, ok := c.(iseq.Equatable); !ok {
		t.Error("PHashMap must implement Equatable")
	}

	if _, ok := c.(iseq.Hashable); !ok {
		t.Error("PHashMap must implement Hashable")
	}
}


// Factory tests

func TestPHashMapISeqFactoryWorks(t *testing.T) {
	var seq iseq.Seq = NewPListFromSlice([]interface{}{"def", 2, "abc", 3, "pqr", 7})
	m := NewPHashMapFromSeq(seq)
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSeq has wrong count, expected %v, got %v", 3, m.Count())
	}

	for s := seq; s != nil; s = s.Next().Next() {
		if m.ValAt(s.First()) != s.Next().First() {
			t.Errorf("NewPHashMapFromSeq: expected key %v => %v, found %v instead", s.First(), s.Next().First(), m.ValAt(s.First()))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSeq: found key that should not be there")
	}
}

func TestPHashMapISeqFactoryWorksWithDupKey(t *testing.T) {
	var seq iseq.Seq = NewPListFromSlice([]interface{}{"def", 2, "abc", 3, "pqr", 7})
	m := NewPHashMapFromSeq(NewCons("def",NewCons("99",seq)))
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSeq has wrong count, expected %v, got %v", 3, m.Count())
	}

	for s := seq; s != nil; s = s.Next().Next() {
		if m.ValAt(s.First()) != s.Next().First() {
			t.Errorf("NewPHashMapFromSeq: expected key %v => %v, found %v instead", s.First(), s.Next().First(), m.ValAt(s.First()))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSeq: found key that should not be there")
	}
}

func TestPHashMapISeqFactoryOnEmpty(t *testing.T) {
	m := NewPHashMapFromSeq(nil)
	if m.Count() != 0  {
		t.Errorf("NewPHashMapFromSeq: on nil, should have count 0, got %v",m.Count())
	}
}


func TestPHashMapSliceFactoryWorks(t *testing.T) {
	s := []interface{}{"def",2,"abc",3,"pqr",7}
	m := NewPHashMapFromSlice(s)
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSlice has wrong count, expected %v, got %v", 3, m.Count())
	}

	for i := 0; i < len(s); i += 2 {
		if m.ValAt(s[i]) != s[i+1] {
			t.Errorf("NewPHashMapFromSlice: expected key %v => %v, found %v instead", s[i], s[i+1], m.ValAt(s[i]))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSlice: found key that should not be there")
	}
}

func TestPHashMapSliceFactoryWorksWithDupKey(t *testing.T) {
	s4 := []interface{}{"def",2,"abc",3,"pqr",7, "def", 99} 
	s3 := s4[2:]
	m := NewPHashMapFromSlice(s4)
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSlice has wrong count, expected %v, got %v", 3, m.Count())
	}

	for i := 0; i < len(s3); i += 2 {
		if m.ValAt(s3[i]) != s3[i+1] {
			t.Errorf("NewPHashMapFromSlice: expected key %v => %v, found %v instead", s3[i], s3[i+1], m.ValAt(s3[i]))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSlice: found key that should not be there")
	}
}

func TestPHashMapSliceFactoryOnEmpty(t *testing.T) {
	m := NewPHashMapFromSlice([]interface{}{})
	if m.Count() != 0  {
		t.Errorf("NewPHashMapFromSlice: on nil, should have count 0, got %v",m.Count())
	}
}

// TODO: Eventually, move this to a benchmark
func createBigSliceForPHashMapTest(n int) []interface{} {
	s := make([]interface{},2*n)
	for i := 0; i<2*n ; i+=2 {
		r := rand.Int63()
		s[i] = fmt.Sprintf("%d",r)
		s[i+1] = r
	}
	return s
}

func TestPHashMapGoesBig(t *testing.T) {
	sizes := []int{10,100,1000,10000,100000}
	for _,n := range sizes {
		s := createBigSliceForPHashMapTest(n)
		//fmt.Printf("Testing big PHashMap creation: %v items\n",n)
		m := NewPHashMapFromSlice(s)
		if m.Count() != n {
			t.Errorf("NewPHashMapFromSlice has wrong count, expected %v, got %v", n, m.Count())
		}

		for i := 0; i < len(s); i += 2 {
			if m.ValAt(s[i]) != s[i+1] {
				t.Errorf("NewPHashMapFromSlice: expected key %v => %v, found %v instead", s[i], s[i+1], m.ValAt(s[i]))
				break
			}
		}
	}
}


// interface iseq.Obj

// TODO: test WithMeta once we have PMap

// interface iseq.Associative, iseq.Lookup

const bigTestPHashMapSize = 100000
var bigTestPHashMapSlice = createBigSliceForPHashMapTest(bigTestPHashMapSize)
var bigTestPHashMap = NewPHashMapFromSlice(bigTestPHashMapSlice)

func TestPHashMapContainsAndValAtAndEntryAt(t *testing.T) {
	m1 := NewPHashMapFromItems(1,2,3,4,nil,6)
	m2 := NewPHashMapFromItems(1,2,3,4,5,6)

	if !m1.ContainsKey(nil) {
		t.Error("Expected to find nil key")
	}
	if m2.ContainsKey(nil) {
		t.Error("Did not expect to find nil key")
	}
	for _,k := range []int{1,3,5} {
		if ! m2.ContainsKey(k) {
			t.Errorf("Expected to find key %v",k)
		}
		if m2.ValAt(k) != k+1 {
			t.Errorf("Expected value %v for key %v, got %v",k+1,k,m2.ValAt(k))
		}
		if m2.ValAtD(k,12) != k+1 {
			t.Errorf("Expected value %v for key %v, got %v",k+1,k,m2.ValAt(k))
		}
		if me := m2.EntryAt(k); me.Key() != k || me.Val() != k+1 {
			t.Errorf("Expected map entry (%v, %v), got (%v, %v)",k,k+1,me.Key(),me.Val())
		}
	}
	for _,k := range []int{2,4,6} {
		if m2.ContainsKey(k) {
			t.Errorf("Did not expect to find key %v",k)
		}
		if v := m2.ValAtD(k,12); v != 12 {
			t.Errorf("Expected default value key %v, got %v",k,v)
		}
		if m2.EntryAt(k) != nil {
			t.Error("Expected nil for MapEntry")
		}
	}

	for i := 0; i < len(bigTestPHashMapSlice); i += 2 {
		k := bigTestPHashMapSlice[i]
		v := bigTestPHashMapSlice[i+1]
		if  !bigTestPHashMap.ContainsKey(k) {
				t.Errorf("Expected to find key %v (item %v)",k,i)
				break
		}
		if bigTestPHashMap.ValAt(k) != v {
			t.Errorf("Expected value %v for key %v, got %v",v,k, bigTestPHashMap.ValAt(k))
		}
		if bigTestPHashMap.ValAtD(k,12) != v {
			t.Errorf("Expected value %v for key %v, got %v",v,k, bigTestPHashMap.ValAt(k))
		}

		if me := bigTestPHashMap.EntryAt(k); me.Key() != k || me.Val() != v {
			t.Errorf("Expected map entry (%v, %v), got (%v, %v)",k,v,me.Key(),me.Val())
		}
	}

	for _,k := range []interface{}{-1,nil} {
		if bigTestPHashMap.ContainsKey(k) {
			t.Errorf("Did not expect to find key %v", k)
		}

		if bigTestPHashMap.EntryAt(k) != nil {
			t.Error("Expected nil MapEntry for key %v", k)
		}

		if v := bigTestPHashMap.ValAtD(k,12); v != 12 {
			t.Errorf("Expected default value for key %v, got %v",k,v)
		}
	} 
}

// interface iseq.PMap
// AssocM, Without, ConsM

func TestPHashMapAssoc(t *testing.T) {
	m1 := NewPHashMapFromItems(1,2,3,4,nil,6)
	m2 := NewPHashMapFromItems(1,2,3,4,5,6)

	// assoc'ing nil
	if m1.AssocM(nil,6) != m1 {
		t.Error("Assoc'ing nil with same value should return same PMap")
	}

	m1a := m1.AssocM(nil,12)
	if m1a.Count() != m1.Count() {
		t.Error("Assoc'ing existing key (nil) should not change count")
	}

	if m1a.ValAt(nil) != 12 {
		t.Error("Assoc'ing nil -- wrong value found")
	}
	for _,k := range []interface{}{1,3} {
		if m1a.ValAt(k) != m1.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m1.ValAt(k),m1a.ValAt(k))
		}
	}

	m2a := m2.AssocM(nil,12)
	if m2a.Count() != m1.Count() + 1 {
		t.Error("Assoc'ing new key (nil) should increase count")
	}

	if m2a.ValAt(nil) != 12 {
		t.Error("Assoc'ing nil -- wrong value found")
	}
	for _,k := range []interface{}{1,3,5} {
		if m2a.ValAt(k) != m2.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m2.ValAt(k),m2a.ValAt(k))
		}
	}

	// assoc'ing a non-nil new key

	m2b := m2.AssocM(7,8)
	if m2b.Count() != m2.Count() + 1 {
		t.Error("Assoc'ing new key should increase count")
	}
	for _,k := range []interface{}{1,3,5} {
		if m2b.ValAt(k) != m2.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m2.ValAt(k),m2a.ValAt(k))
		}
	}
	if m2b.ValAt(7) != 8 {
			t.Error("On key 7, expected 8, got %v",m2a.ValAt(7))

	}

	// assoc'ing a non-nil existing key

	m2b = m2.AssocM(3,8)
	if m2b.Count() != m2.Count() {
		t.Error("Assoc'ing existing key should leave count unchanged")
	}
	for _,k := range []interface{}{1,5} {
		if m2b.ValAt(k) != m2.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m2.ValAt(k),m2a.ValAt(k))
		}
	}
	if m2b.ValAt(3) != 8 {
			t.Error("On key 3, expected 8, got %v",m2a.ValAt(7))

	}
}

func TestPHashMapWithout(t *testing.T) {
	m1 := NewPHashMapFromItems(1,2,3,4,nil,6)
	m2 := NewPHashMapFromItems(1,2,3,4,5,6)

	// without'ing nil

	m1a := m1.Without(nil)
	if m1a.Count() != m1.Count() - 1 {
		t.Error("Without'ing existing key (nil) should decrease count")
	}

	if m1a.ContainsKey(nil)  {
		t.Error("Without'ing nil -- nil key still present")
	}

	for _,k := range []interface{}{1,3} {
		if m1a.ValAt(k) != m1.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m1.ValAt(k),m1a.ValAt(k))
		}
	}

	m2a := m2.Without(nil)
	if m2a.Count() != m1.Count() {
		t.Error("Without'ing a non-present key (nil) should not change count")
	}

	if m2a.ContainsKey(nil) {
		t.Error("Without'ing nil -- key suddenly appeared")
	}
	for _,k := range []interface{}{1,3,5} {
		if m2a.ValAt(k) != m2.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m2.ValAt(k),m2a.ValAt(k))
		}
	}

	// without'ing a non-nil non-present key

	m2b := m2.Without(7)
	if m2b.Count() != m2.Count() {
		t.Error("Without'ing a non-present key should not change count")
	}
	for _,k := range []interface{}{1,3,5} {
		if m2b.ValAt(k) != m2.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m2.ValAt(k),m2a.ValAt(k))
		}
	}
	if m2b.ContainsKey(7) {
			t.Error("Without'ing a non-present key, it suddenly appears")

	}

	// assoc'ing a non-nil present key

	m2b = m2.Without(3)
	if m2b.Count() != m2.Count() -1 {
		t.Error("Without'ing a present key should decrement count")
	}
	for _,k := range []interface{}{1,5} {
		if m2b.ValAt(k) != m2.ValAt(k) {
			t.Error("On key %v, expected %v, got %v",k,m2.ValAt(k),m2a.ValAt(k))
		}
	}
	if m2b.ContainsKey(3) {
			t.Error("Without'ing present key, but it's still there")

	}
}

// TODO: Finish tests
