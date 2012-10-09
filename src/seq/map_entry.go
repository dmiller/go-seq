// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

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
