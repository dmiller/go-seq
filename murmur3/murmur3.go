// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package murmer3 provides functions implementing the Murmur3 hashing algorithm.
// The ClojureJVM version imported the Guava Murmur3 implementation
// and made some changes.
// For ClojureCLR and here, I copied the API stubs, then implemented the API
// based on the algorithm description at
//     http://en.wikipedia.org/wiki/MurmurHash.
//     See also: https://code.google.com/p/smhasher/wiki/MurmurHash3. </p
// The implementations of HashUnordered and HashOrdered taken from ClojureJVM.
package murmur3

const seed uint32 = 0
const c1 uint32 = 0xcc9e2d51
const c2 uint32 = 0x1b873593
const r1 uint32 = 15
const r2 uint32 = 13
const m uint32 = 5
const n uint32 = 0xe6546b64

// The public interface

func HashInt32(input int32) uint32 {
	return HashUInt32(uint32(input))
}

func HashInt64(input int64) uint32 {
	return HashUInt64(uint64(input))
}

func HashUInt32(input uint32) uint32 {
	if input == 0 {
		return 0
	}

	key := MixKey(input)
	hash := MixHash(seed, key)
	return Finalize(hash, 4)
}

func HashUInt64(input uint64) uint32 {
	if input == 0 {
		return 0
	}

	low := uint32(input)
	high := uint32(input >> 32)

	key := MixKey(low)
	hash := MixHash(seed, key)

	key = MixKey(high)
	hash = MixHash(hash, key)

	return Finalize(hash, 8)
}

func HashString(input string) uint32 {

	hash := seed
	len := len(input)

	// step through the string 4 bytes at a time
	for i := 3; i < len; i += 4 {
		key := uint32(input[i-3]) | uint32(input[i-2]<<8) | uint32(input[i-1]<<16) | uint32(input[i]<<24)
		key = MixKey(key)
		hash = MixHash(hash, key)
	}

	// deal with remaining characters

	if len != 0 {
		var key uint32
		switch len % 4 {
		case 1:
			key = uint32(input[len-1])
		case 2:
			key = uint32(input[len-2]) | uint32(input[len-1]<<8)
		case 3:
			key = uint32(input[len-3]) | uint32(input[len-2]<<8) | uint32(input[len-1]<<16)
		}
		key = MixKey(key)
		hash = MixHash(hash, key)
	}

	return Finalize(hash, int32(len))
}

func MixKey(key uint32) uint32 {
	key *= c1
	key = rotateLeft(key, r1)
	key *= c2
	return key
}

func MixHash(hash uint32, key uint32) uint32 {
	hash ^= key
	hash = rotateLeft(hash, r2)
	hash = hash*m + n
	return hash

}

// finalize forces all bits of a hash block to avalanche
func Finalize(hash uint32, length int32) uint32 {
	hash ^= uint32(length)
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16
	return hash
}

func FinalizeCollHash(hash uint32, count int32) uint32 {
	h1 := seed
	k1 := MixKey(hash)
	h1 = MixHash(h1, k1)
	return Finalize(h1, count)
}

// implementation details

func rotateLeft(x uint32, n uint32) uint32 {
	return (x << n) | (x >> (32 - n))
}
