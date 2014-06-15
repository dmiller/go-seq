// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

import (
	"math/rand"
	"testing"
	"time"
)

// Several test tables for obvious values, just to get warmed up

var signedTests = []struct {
	in  int32
	out int
}{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 1},
	{5, 2},
	{6, 2},
	{7, 3},
	{-1, 32},
	{-2, 31},
	{-3, 31},
	{-4, 30},
	{-5, 31},
	{-6, 30},
	{-7, 30},
}

var unsignedTests = []struct {
	in  uint32
	out int
}{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 1},
	{5, 2},
	{6, 2},
	{7, 3},
	{0xFFFFFFFF, 32},
	{0xF0F0F0F0, 16},
	{0x87878787, 16},
}

// Constructing (u)int32s from bytes and counting the bits

var byteCounts = make([]int, 256)

func init() {
	for i := 0; i < 256; i++ {
		byteCounts[i] = byteBitCount(byte(i))
	}
}

// Attributed to Brian Kernighan
func byteBitCount(b byte) int {
	c := 0
	for ; b != 0; c++ {
		b &= b - 1 // clear the least significan bit set
	}
	return c
}

func bytesToUInt32(b1, b2, b3, b4 byte) uint32 {
	return uint32(b1) | uint32(b2)<<8 | uint32(b3)<<16 | uint32(b4)<<24
}

func generateRandomUInt32AndCount(r *rand.Rand) (uint32, int) {
	b1 := byte(r.Intn(256))
	b2 := byte(r.Intn(256))
	b3 := byte(r.Intn(256))
	b4 := byte(r.Intn(256))
	return bytesToUInt32(b1, b2, b3, b4),
		byteBitCount(b1) + byteBitCount(b2) + byteBitCount(b3) + byteBitCount(b4)
}

func TestBitCountRandom(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10000; i++ {
		v, expect := generateRandomUInt32AndCount(r)
		cnt := BitCount(int32(v))
		if cnt != expect {
			t.Errorf("%d. BitCount(%d) => %d, want %d", i, int32(v), cnt, expect)
		}
	}
}

func TestBitCountU32Random(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10000; i++ {
		v, expect := generateRandomUInt32AndCount(r)
		cnt := BitCountU32(v)
		if cnt != expect {
			t.Errorf("%d. BitCount(%d) => %d, want %d", i, v, cnt, expect)
		}
	}
}

func TestBitCount(t *testing.T) {

	for i, tt := range signedTests {
		c := BitCount(tt.in)
		if c != tt.out {
			t.Errorf("%d. BitCount(%d) => %d, want %d", i, tt.in, c, tt.out)
		}
	}
}

func TestBitCountU32(t *testing.T) {

	for i, tt := range unsignedTests {
		c := BitCountU32(tt.in)
		if c != tt.out {
			t.Errorf("%d. BitCountU32(%d) => %d, want %d", i, tt.in, c, tt.out)
		}
	}
}
