// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

import (
	"errors"
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/murmur3"
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
	if err != nil {
		return 0
	}
	return h
}

// HashE returns a hash code for an object.
// Uses iseq.Hashable.Hash if the interface is implemented.
// It special cases Go numbers and strings.
// Returns an error if the object is not covered by these cases.
func HashE(v interface{}) (uint32, error) {
	if h, ok := v.(iseq.Hashable); ok {
		return h.Hash(), nil
	}

	switch v := v.(type) {
	case bool, int, int8, int32, int64:
		return murmur3.HashUInt64(uint64(reflect.ValueOf(v).Int())), nil
	case uint, uint8, uint32, uint64:
		return murmur3.HashUInt64(uint64(reflect.ValueOf(v).Uint())), nil
	case float32, float64:
		return murmur3.HashUInt64(math.Float64bits(reflect.ValueOf(v).Float())), nil
	case nil:
		return murmur3.HashUInt64(0), nil
	case string:
		return murmur3.HashString(v), nil
	case complex64, complex128:
		return HashComplex128(v.(complex128)), nil
	}
	return 0, errors.New("don't know how to hash")
}

// IsHashable returns true if Hash/HashE can compute a hash for this object.
func IsHashable(v interface{}) bool {
	if _, ok := v.(iseq.Hashable); ok {
		return true
	}

	switch v.(type) {
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

// HashSeq computes a hash for an iseq.Seq
func HashSeq(s iseq.Seq) uint32 {
	return HashOrdered(s)
}

// HashMap computes a hash for an iseq.PMap
func HashMap(m iseq.PMap) uint32 {
	return HashUnordered(m.Seq())
}

// HashOrdered computes a hash for an iseq.Seq, where order is important
func HashOrdered(s iseq.Seq) uint32 {
	n := int32(0)
	hash := uint32(1)

	for ; s != nil; s = s.Next() {
		hash = 31*hash + Hash(s.First)
		n++
	}
	return murmur3.FinalizeCollHash(hash, n)
}

// HashUnordered computes a hash for an iseq.Seq, independent of order of elements
func HashUnordered(s iseq.Seq) uint32 {
	n := int32(0)
	hash := uint32(0)

	for ; s != nil; s = s.Next() {
		hash += Hash(s.First)
		n++
	}
	return murmur3.FinalizeCollHash(hash, n)
}

// HashComplex128 computes a hash for a complex128
func HashComplex128(c complex128) uint32 {
	hash := murmur3.MixHash(
		murmur3.HashUInt64(math.Float64bits(real(c))),
		murmur3.HashUInt64(math.Float64bits(imag(c))))
	return murmur3.Finalize(hash, 2)
}
