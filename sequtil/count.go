// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

import (
	"github.com/dmiller/go-seq/iseq"
)

// Count computes the length of a sequence.
// Handles nil, iseq.Counted, iseq.PCollection, and strings
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

// SeqCount computes the length of an iseq.Seq
// Only call this on known non-empty iseq.Seq
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
