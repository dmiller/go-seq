// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

import (
	"github.com/dmiller/go-seq/iseq"
)

// MapCons adds (conses) a new key/value pair onto an iseq.PMap
// A MapEntry adds its key/value.
// A PVector uses v[2*i] v[2*i+1] as key/value pairs
// Otherwise, we need a sequence of iseq.MapEntry values
// Assumes its argument is one of the above; else panics
func MapCons(m iseq.PMap, o interface{}) iseq.PMap {
	if me, ok := o.(iseq.MapEntry); ok {
		return m.AssocM(me.Key(), me.Val())
	}

	if v, ok := o.(iseq.PVector); ok {
		if v.Count() != 2 {
			panic("Vector arg to map cons must be a pair")
		}
		return m.AssocM(v.Nth(0), v.Nth(1))
	}

	ret := m
	for s := ConvertToSeq(o); s != nil; s = s.Next() {
		me := s.First().(iseq.MapEntry)
		ret = ret.AssocM(me.Key(), me.Val())
	}
	return ret
}

// ConvertToSeq attempt to convert its argument to an iseq.Seq
// internal use: assumes the arg can be converted, else panics
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

	panic("Cannot convert object to sequence")
}
