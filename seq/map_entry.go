// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
)

// A MapEntry represents a key/value pair. Implements iseq.MapEntry.
type MapEntry struct {
	key interface{}
	val interface{}
}

// Key returns the key.
func (m MapEntry) Key() interface{} {
	return m.key
}

// Val returns the value.
func (m MapEntry) Val() interface{} {
	return m.val
}

// Equiv returns true if its argument is an iseq.MapEntry with equivalent key and value.
func (m MapEntry) Equiv(o interface{}) bool {
	if you, ok := o.(iseq.MapEntry); ok {
		return sequtil.Equiv(m.key, you.Key()) && sequtil.Equiv(m.val, you.Val())
	}
	return false
}

// Hash computes a hash value (based only on the key)
func (m MapEntry) Hash() uint32 {
	return sequtil.Hash(m.key)
}
