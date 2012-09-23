// Package sequtils contains some generally useful functions for implementing the Clojure sequences
package sequtils

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"seq"
	//"fmt"
	"math"
	"reflect"
)

func Equals(o1 interface{}, o2 interface{}) bool {
	if o1 == o2 {
		return true
	}

	if e1, ok1 := o1.(seq.Equatable); ok1 {
		return e1.Equals(o2)
	}

	return false
}

func SeqEquals(s1 seq.Seq, s2 seq.Seq) bool {
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

func Equiv(o1 interface{}, o2 interface{}) bool {
	if o1 == o2 {
		return true
	}
	if o1 != nil {
		// TODO: Determine how to handle numbers. Do we want Clojure's semantics?
		// Go's semantics says the o1 == o2 case is enough
		pc1, ok1 := o1.(seq.PersistentCollection)
		if ok1 {
			return pc1.Equiv(o2)
		}

		pc2, ok2 := o2.(seq.PersistentCollection)
		if ok2 {
			return pc2.Equiv(o1)
		}

		return Equals(o1, o2)
	}

	return false
}

func SeqEquiv(s1 seq.Seq, s2 seq.Seq) bool {
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

	if cnt, ok := o.(seq.Counted); ok {
		return cnt.Count1()
	}

	if pc, ok := o.(seq.PersistentCollection); ok {
		s := pc.Seq()
		i := 0
		for ; s != nil; s = s.Next() {
			if c, ok := s.(seq.Counted); ok {
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

var (
	zeroBytes = make([]byte, 4)
)

func HashSeq(seq seq.Seq) uint32 {
	h := fnv.New32()
	AddHashSeq(h, seq)
	return h.Sum32()
}

func AddHashSeq(h hash.Hash, seq seq.Seq) {
	for s := seq; s != nil; s = s.Next() {
		if f := s.First(); f == nil {
			h.Write(zeroBytes)
		} else {
			AddHash(h, f)
		}
	}
}

func HashUint64(v uint64) uint32 {
	h := fnv.New32()
	AddHashUint64(h, v)
	return h.Sum32()
}

func AddHashUint64(h hash.Hash, v uint64) {
	data := make([]byte, 4)
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
	if sh, ok := v.(seq.Hashable); ok {
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
	if h, ok := v.(seq.Hashable); ok {
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
