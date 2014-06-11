// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

// BitCount32 counts the on bits in an int32.
// No multiplication, for those that care.
// This algo is in a lot of places.
// See, for example, http://aggregate.org/MAGIC/#Population%20Count%20(Ones%20Count)
func BitCount(x int32) int {
	x = x - ((x >> 1) & 0x55555555)
	x = (((x >> 2) & 0x33333333) + (x & 0x33333333))
	x = (((x >> 4) + x) & 0x0f0f0f0f)
	return int(((x * 0x01010101) >> 24))
}

// BitCountU32 counts the on bits in an uint32.
// No multiplication, for those that care.
// This algo is in a lot of places.
// See, for example, http://aggregate.org/MAGIC/#Population%20Count%20(Ones%20Count)
func BitCountU32(x uint32) int {
	x = x - ((x >> 1) & 0x55555555)
	x = (((x >> 2) & 0x33333333) + (x & 0x33333333))
	x = (((x >> 4) + x) & 0x0f0f0f0f)
	x = x + (x >> 8)
	x = x + (x >> 16)
	return int(x & 0x0000003f)
}
