// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"sync/atomic"
)

// An IDGenerator provides a sequential stream of uint64 values (safe for concurrent access)
type IDGenerator struct {
	id uint64
}

// Next returns the next value in sequence
func (g *IDGenerator) Next() uint64 {
	return atomic.AddUint64(&g.id, 1)
}
