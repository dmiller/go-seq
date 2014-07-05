// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"testing"
)

func TestIdGenSequentialValues(t *testing.T) {
	c := new(IDGenerator)

	v1 := c.Next()

	if v1 != 1 {
		t.Errorf("First value should be 1, got %d", v1)
	}

	v2 := c.Next()

	if v2 != 2 {
		t.Errorf("Second value should be 2, got %d", v2)
	}
}
