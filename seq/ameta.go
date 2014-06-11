// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"github.com/dmiller/go-seq/iseq"
)

// AMeta provides a slot to hold a 'meta' value
type AMeta struct {
	meta iseq.PMap
}

func (o *AMeta) Meta() iseq.PMap {
	return o.meta
}
