// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"sync/atomic"
)

type IdGenerator struct {
	id uint64
}

func (g *IdGenerator) Next() uint64 {
	return atomic.AddUint64(&g.id, 1)
}
