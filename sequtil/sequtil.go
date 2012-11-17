// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sequtil contains some generally useful functions for implementing the Clojure sequences
package sequtil

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"github.com/dmiller/go-seq/iseq"
	//"fmt"
	"math"
	"reflect"
)

func Equals(o1 interface{}, o2 interface{}) bool {
	if o1 == o2 {
		return true
	}

	if e1, ok1 := o1.(iseq.Equatable); ok1 {
		return e1.Equals(o2)
	}

	return false
}

func SeqEquals(s1 iseq.Seq, s2 iseq.Seq) bool {
	if s1 == s2 {
		return true
	}

	iter2 := s2

	for iter1 := s1; iter1 != nil; iter1 = iter1.Next() {
		if iter2 == nil || !Equals(iter1.First(), iter2.First()) {
			return false
		}
		iter2 = iter2.Next()
	}

	return iter2 == nil
}

func MapEquals(m1 iseq.PMap, obj interface{}) bool {
	if m1 == obj {
		return true
	}

	if _,ok := obj.(map[interface{}]interface{}); ok {
		// TODO: figure out how to handle go maps
		return false
	}

	if m2,ok := obj.(iseq.PMap); ok {
		if m1.Count() != m2.Count() {
			return false
		}

		for s := m1.Seq(); s != nil; s = s.Next() {
			me := s.First().(iseq.MapEntry)
			found := m2.ContainsKey(me.Key())
			if ! found || !Equals(me.Val(),m2.ValAt(me.Key())) {
				return false
			}
		}
		return true
	}
	return false
}

func Equiv(o1 interface{}, o2 interface{}) bool {
	if o1 == o2 {
		return true
	}
	if o1 != nil {
		// TODO: Determine how to handle numbers. Do we want Clojure's semantics?
		// Go's semantics says the o1 == o2 case is enough
		pc1, ok1 := o1.(iseq.PCollection)
		if ok1 {
			return pc1.Equiv(o2)
		}

		pc2, ok2 := o2.(iseq.PCollection)
		if ok2 {
			return pc2.Equiv(o1)
		}

		return Equals(o1, o2)
	}

	return false
}

func SeqEquiv(s1 iseq.Seq, s2 iseq.Seq) bool {
	if s1 == s2 {
		return true
	}

	iter2 := s2

	for iter1 := s1; iter1 != nil; iter1 = iter1.Next() {
		if iter2 == nil || !Equiv(iter1.First(), iter2.First()) {
			return false
		}
		iter2 = iter2.Next()
	}

	return iter2 == nil
}

func Count(o interface{}) int {
	if o == nil {
		return 0
	}

	if cnt, ok := o.(iseq.Counted); ok {
		return cnt.Count1()
	}

	if pc, ok := o.(iseq.PCollection); ok {
		s := pc.Seq()
		i := 0
		for ; s != nil; s = s.Next() {
			if c, ok := s.(iseq.Counted); ok {
				return i + c.Count1()
			}
			i++
		}
		return i
	}

	if s, ok := o.(string); ok {
		return len(s)
	}
	// TODO: Figure out how to  handle arrays, slices, maps in a typeswitch/generic way
	panic("Count not supported on this type")
}

// Only call this on known non-empty
func SeqCount(s0 iseq.Seq) int {
	i := 1 // if we are here, it is non-empty
	for s := s0.Next(); s != nil; s, i = s.Next(), i+1 {
		if cnt, ok := s.(iseq.Counted); ok {
			return i + cnt.Count1()
		}
	}
	return i
}


var (
	zeroBytes = make([]byte, 4)
)

func HashSeq(seq iseq.Seq) uint32 {
	h := fnv.New32()
	AddHashSeq(h, seq)
	return h.Sum32()
}

func AddHashSeq(h hash.Hash, seq iseq.Seq) {
	for s := seq; s != nil; s = s.Next() {
		if f := s.First(); f == nil {
			h.Write(zeroBytes)
		} else {
			AddHash(h, f)
		}
	}
}

func HashMap(m iseq.PMap) uint32 {
	h := fnv.New32()
	AddHashMap(h,m)
	return h.Sum32()
}

func AddHashMap(h hash.Hash, m iseq.PMap) {
	for s := m.Seq(); s != nil; s = s.Next() {
		me := s.First().(iseq.MapEntry)
		AddHash(h,me.Key())
		AddHash(h,me.Val())
	}
}


func HashUint64(v uint64) uint32 {
	h := fnv.New32()
	AddHashUint64(h, v)
	return h.Sum32()
}

func AddHashUint64(h hash.Hash, v uint64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, v)
	h.Write(data)
}

func hashComplex128(v complex128) uint32 {
	h := fnv.New32()
	addHashComplex128(h, v)
	return h.Sum32()
}

func addHashComplex128(h hash.Hash, v complex128) {
	AddHashUint64(h, math.Float64bits(real(v)))
	AddHashUint64(h, math.Float64bits(imag(v)))
}

func hashString(s string) uint32 {
	h := fnv.New32()
	addHashString(h, s)
	return h.Sum32()
}

func addHashString(h hash.Hash, s string) {
	h.Write([]byte(s))
}

func AddHash(h hash.Hash, v interface{}) {
	if sh, ok := v.(iseq.Hashable); ok {
		sh.AddHash(h)
		return
	}

	switch v := v.(type) {
	default:
		panic("Cannot hash element")
	case bool, int, int8, int32, int64:
		AddHashUint64(h, uint64(reflect.ValueOf(v).Int()))
	case uint, uint8, uint32, uint64:
		AddHashUint64(h, uint64(reflect.ValueOf(v).Uint()))
	case float32, float64:
		AddHashUint64(h, math.Float64bits(reflect.ValueOf(v).Float()))
	case nil:
		AddHashUint64(h, 0)
	case string:
		addHashString(h, v)
	case complex64, complex128:
		addHashComplex128(h, v.(complex128))
	}
}

func Hash(v interface{}) uint32 {
	if h, ok := v.(iseq.Hashable); ok {
		return h.Hash()
	}

	switch v := v.(type) {
	case bool, int, int8, int32, int64:
		return HashUint64(uint64(reflect.ValueOf(v).Int()))
	case uint, uint8, uint32, uint64:
		return HashUint64(uint64(reflect.ValueOf(v).Uint()))
	case float32, float64:
		return HashUint64(math.Float64bits(reflect.ValueOf(v).Float()))
	case nil:
		return HashUint64(0)
	case string:
		return hashString(v)
	case complex64, complex128:
		return hashComplex128(v.(complex128))
	}
	panic("Cannot hash element")
}

// TODO: investigate use of IHashEq
func Hasheq(o interface{}) uint32 {
	if o == nil {
		return 1
	}

	// if ihe,ok := o.(iseq.HashEq); ok {
	// 	return ih3.hasheq()
	// }

	return Hash(o)
}

func MapCons(m iseq.PMap, o interface{}) iseq.PMap {
	if me, ok := o.(iseq.MapEntry); ok {
		return m.AssocM(me.Key(),me.Val())
	}

	if v, ok := o.(iseq.PVector); ok {
		if v.Count() != 2 {
			panic("Vector arg to map cons must be a pair")
		}
		return m.AssocM(v.Nth(0),v.Nth(1))
	}
	
	var ret iseq.PMap
	for s := ConvertToSeq(o); s != nil; s = s.Next() {
		me := s.First().(iseq.MapEntry)
		ret = ret.AssocM(me.Key(),me.Val())
	}
	return ret
}

func ConvertToSeq(o interface{}) iseq.Seq {
	// TODO: handle more general cases of maps, slices, arrays
	if o == nil {
		return nil
	}
	if s, ok := o.(iseq.Seq); ok {
		return s
	}

	if s, ok := o.(iseq.Seqable); ok {
		return s.Seq()
	}
	return nil
	// or maybe panic?
}

func BitCount(x int32) int {
	x = x - ((x >> 1) & 0x55555555)
	x = (((x >> 2) & 0x33333333) + (x & 0x33333333))
	x = (((x >> 4) + x) & 0x0f0f0f0f)
	return int(((x * 0x01010101) >> 24))
}


// A variant of the above that avoids multiplying
// This algo is in a lot of places.
// See, for example, http://aggregate.org/MAGIC/#Population%20Count%20(Ones%20Count)
func BitCountU(x uint32) int {
	x = x- ((x >> 1) & 0x55555555);
	x = (((x >> 2) & 0x33333333) + (x & 0x33333333));
	x = (((x >> 4) + x) & 0x0f0f0f0f);
	x = x + (x >> 8);
	x = x +(x >> 16);
	return int(x & 0x0000003f);
}

func DefaultCompareFn(k1 interface{}, k2 interface{}) int {
	if k1 == k2 {
		return 0
	}
	if k1 != nil {
		if k2 == nil {
			return 1
		}
		if c,ok := k1.(iseq.Comparer); ok {
			return c.Compare(k2)
		}
		if s,ok := k1.(string); ok {
			return CompareString(s,k2)
		}
		if IsComparableNumeric(k1) {
			return CompareComparableNumeric(k1,k2)
		}
		panic("Can't compare")
	}
	return -1  
}

func IsComparableNumeric(v interface{}) bool {

	switch v.(type) {
	case bool, int, int8, int32, int64, 
		uint, uint8, uint32, uint64, 
		float32, float64:

		return true
	}
	return false
}

func CompareString(s string, x interface{}) int {
	if s2, ok := x.(string); ok {
		if s < s2 {
			return -1
		}
		if s == s2 {
			return 0
		}
		return 1
	}

	return -1 // don't feel like panicking.
}

func CompareComparableNumeric(x1 interface{}, x2 interface{} ) int {
	// n1 should be numeric
	switch x1 := x1.(type) {
	case bool, int, int8, int32, int64:
		n1 := reflect.ValueOf(x1).Int()
		return compareNumericInt(n1,x2)
	case uint, uint8, uint32, uint64:
		n1 := reflect.ValueOf(x1).Uint()
		return compareNumericUint(n1,x2)
	case float32, float64:
		n1 := reflect.ValueOf(x1).Float()
		return compareNumericFloat(n1,x2)
	}
	panic("Expect first arg to be numeric")
}

func compareNumericInt(n1 int64, x2 interface{}) int {
	switch x2 := x2.(type) {
	case bool, int, int8, int32, int64:
		n2 := reflect.ValueOf(x2).Int()
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		} 
		return 0

	case uint, uint8, uint32, uint64:
		n2 := reflect.ValueOf(x2).Uint()
		if n1 < 0 {
			return -1
		}
		un1 := uint64(n2)
		if un1 < n2 {
			return -1
		}
		if un1 > n2 {
			return 1
		}
		return 0

	case float32, float64:
		n2 := reflect.ValueOf(x2).Float()
		fn1 := float64(n1)
		if fn1 < n2 {
			return -1
		}
		if fn1 > n2 {
			return 1
		}
		return 0
	}
	return -1  // what else, other than panic?
}

func compareNumericUint(n1 uint64, x2 interface{}) int {
	switch x2 := x2.(type) {
	case bool, int, int8, int32, int64:
		n2 := reflect.ValueOf(x2).Int()
		if n2 < 0 {
			return 1
		}
		un2 := uint64(n2)
		if n1 < un2 {
			return -1
		}
		if n1 > un2 {
			return 1
		} 
		return 0

	case uint, uint8, uint32, uint64:
		n2 := reflect.ValueOf(x2).Uint()
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
		return 0

	case float32, float64:
		n2 := reflect.ValueOf(x2).Float()
		fn1 := float64(n1)
		if fn1 < n2 {
			return -1
		}
		if fn1 > n2 {
			return 1
		}
		return 0
	}
	return -1  // what else, other than panic?
}

func compareNumericFloat(n1 float64, x2 interface{}) int {
	var n2 float64
	switch x2 := x2.(type) {
	case bool, int, int8, int32, int64:
		n2 = float64(reflect.ValueOf(x2).Int())
	case uint, uint8, uint32, uint64:
		n2 = float64(reflect.ValueOf(x2).Uint())
	case float32, float64:
		n2 = reflect.ValueOf(x2).Float()
	default:
		return -1  // what else, other than panic?
	}
	if n1 < n2 {
		return -1
	}
	if n1 > n2 {
		return 1
	}
	return 0
}
