// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
	"hash"
)

type MapEntry struct {
	key interface{}
	val interface{}
}

func (me MapEntry) Key() interface{} {
	return me.key
}

func (me MapEntry) Val() interface{} {
	return me.val
}

func (me MapEntry) Equals(o interface{}) bool {
	if you, ok := o.(iseq.MapEntry); ok {
		return sequtil.Equals(me.key, you.Key()) && sequtil.Equals(me.val, you.Val())
	}
	return false
}

func (me MapEntry) Hash() uint32 {
	return sequtil.Hash(me.key)
}

func (me MapEntry) AddHash(h hash.Hash) {
	sequtil.AddHash(h, me.key)
}
