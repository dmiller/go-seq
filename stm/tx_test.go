// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"testing"
)

func TestSimpleTxCall(t *testing.T) {
	o := new(interface{})
	f := func(tx *Tx) interface{} {
		return o
	}

	v, e := RunInTransaction(f)

	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if v != o {
		t.Errorf("Expected return value %v, got %v", o, v)
	}
}

func TestBasicRefSetting(t *testing.T) {
	init1 := 12
	init2 := 100
	new1 := 1000
	new2 := 2000

	r1 := NewRef(init1)
	r2 := NewRef(init2)

	r3 := NewRef(-1)
	r4 := NewRef(-1)
	r5 := NewRef(-1)

	ret := 300

	f := func(tx *Tx) interface{} {
		r3.Set(tx, r1.Deref(tx).(int)+r2.Deref(tx).(int))
		r1.Set(tx, new1)
		r4.Set(tx, r1.Deref(tx).(int)+r2.Deref(tx).(int))
		r2.Set(tx, new2)
		r5.Set(tx, r1.Deref(tx).(int)+r2.Deref(tx).(int))
		return ret
	}

	v, e := RunInTransaction(f)

	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if v != ret {
		t.Errorf("Expected return value of %v, got %v", ret, v)
	}

	if v := r1.Deref(nil); v != new1 {
		t.Errorf("Expected r1 to have value %v, got %v", new1, v)
	}

	if v := r2.Deref(nil); v != new2 {
		t.Errorf("Expected r2 to have value %v, got %v", new2, v)
	}

	if v := r3.Deref(nil); v != init1+init2 {
		t.Errorf("Expected r3 to have value %v, got %v", init1+init2, v)
	}

	if v := r4.Deref(nil); v != new1+init2 {
		t.Errorf("Expected r4 to have value %v, got %v", new1+init2, v)
	}

	if v := r5.Deref(nil); v != new1+new2 {
		t.Errorf("Expected r5 to have value %v, got %v", new1+new2, v)
	}
}

func TestNoCommitOnAbort(t *testing.T) {
	init1 := 12
	new1 := 1000

	r1 := NewRef(init1)

	var save1, save2 interface{}

	ret := 300

	f := func(tx *Tx) interface{} {
		save1 = r1.Deref(tx)
		r1.Set(tx, new1)
		save2 = r1.Deref(tx)
		tx.Abort()
		return ret
	}

	defer func() {

		r := recover()

		if r == nil {
			t.Errorf("Expected a panic due to transaction abort, didn't get one")
		} else if r.(error).Error() != "Transaction aborted" {
			t.Errorf("Expected transaction abort error, got %v", r)

		}

		if v := r1.Deref(nil); v != init1 {
			t.Errorf("Expected r1 to have value %v, got %v", init1, v)
		}

		if save1 != init1 {
			t.Errorf("Expected save1 to have value %v, got %v", init1, save1)
		}

		if save2 != new1 {
			t.Errorf("Expected save2 to have value %v, got %v", new1, save2)
		}
	}()

	RunInTransaction(f)

	// shouldn't get here

	if true {
		t.Errorf("Abort should prevented us from reaching this point")
	}

}
