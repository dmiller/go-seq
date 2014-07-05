// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"testing"
)

// Here we test only those parts of Ref that do not involve transactions

func TestNewRef(t *testing.T) {

	o := new(interface{})

	r := NewRef(o)

	if h := r.MinHistory(); h != 0 {
		t.Errorf("For a new Ref, minHistory should be 0, found %d", h)
	}

	if h := r.MaxHistory(); h != DefaultMaxHistory {
		t.Errorf("For a new Ref, maxHistory should be %d, found %d", DefaultMaxHistory, h)
	}

	if v := r.Deref(nil); v != o {
		t.Errorf("For a new Ref, value should be supplied object. Found %v, expected %v", v, o)
	}

}
