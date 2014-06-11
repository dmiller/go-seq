// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

import (
	"encoding/binary"
	"errors"
	"github.com/dmiller/go-seq/iseq"
	"hash"
	"hash/fnv"
	//"fmt"
	"math"
	"reflect"
)

// Hash returns a hash code for an object.
// Uses iseq.Hashable.Hash if the interface is implemented.
// Otherwise, special cases Go numbers and strings.
// Returns a default value if not covered by these cases.
// Returning a default value is really not good for hashing performance.
// But it is one way not to have an error code and to avoid a panic.
// Use IsHashable to determine if Hash is supported.
// Or call HashE which has an error return.
func Hash(v interface{}) uint32 {
	h, err := HashE(v)
	if err != nil
		return 0
	return h
}

// Hash returns a hash code for an object.
// Uses iseq.Hashable.Hash if the interface is implemented.
// It special cases Go numbers and strings.
// Returns an error if the object is not covered by these cases.
func HashE(v interface{}) (uint32,error) {
	if h, ok := v.(iseq.Hashable); ok {
		return h.Hash(), nil
	}

	switch v := v.(type) {
	case bool, int, int8, int32, int64:
		return HashUint64(uint64(reflect.ValueOf(v).Int())), nil
	case uint, uint8, uint32, uint64:s
		return HashUint64(uint64(reflect.ValueOf(v).Uint())), nil
	case float32, float64:
		return HashUint64(math.Float64bits(reflect.ValueOf(v).Float())), nil
	case nil:
		return HashUint64(0), nil
	case string:
		return hashString(v), nil
	case complex64, complex128:
		return hashComplex128(v.(complex128)), nil
	}
	return 0, errors.New("don't know how to hash")
}

// IsHashable returns true if Hash/HashE can compute a hash for this object.
func IsHashable(v interface{}) bool {
	if h, ok := v.(iseq.Hashable); ok {
		return true
	}

	switch v := v.(type) {
	case bool, int, int8, int32, int64, 
		 uint, uint8, uint32, uint64,
		float32, float64,
	    nil,
		string,
		complex64, complex128:
		return true
	}
	return false
}


func HashSeq(s iseq.Seq) uint32 {
	return HashOrdered(s)
}

func HashMap(m iseq.PMap) uint32 {
	return HashUnordered(m)
}


func HashOrdered(s iseq.Seq) {
	n int32 := 0
	hash uint32 := 1

	for ; s != nil; s = s.Next() {
		hash = 31 * hash + Hash(s.First)
		n += 1
	}
	return murmur3.FinalizeCollHash(hash,n)
}        

func HashUnordered(s iseq.Seq) {
	n int32 := 0
	hash uint32 := 0

	for ; s != nil; s = s.Next() {
		hash += Hash(s.First)
		n += 1
	}
	return murmur3.FinalizeCollHash(hash,n)
}


func HashUint64(v uint64) uint32 {
	h := fnv.New32()
	AddHashUint64(h, v)
	return h.Sum32()
}



func hashComplex128(v complex128) uint32 {
	h := fnv.New32()
	addHashComplex128(h, v)
	return h.Sum32()
}



func hashString(s string) uint32 {
	h := fnv.New32()
	addHashString(h, s)
	return h.Sum32()
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



