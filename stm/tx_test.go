// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"runtime"
	"sync"
	"testing"
	"time"
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
	t.Errorf("Abort should prevented us from reaching this point")

}

func TestNoSetAfterCommute(t *testing.T) {
	init1 := 12
	add1 := 10
	add2 := 20

	var point1, point2 bool

	r1 := NewRef(init1)
	ret := 300

	fc := func(old interface{}, args ...interface{}) interface{} {
		return old.(int) + args[0].(int)
	}

	f := func(tx *Tx) interface{} {

		r1.Commute(tx, fc, add1)
		point1 = true

		r1.Set(tx, add2)
		point2 = true

		return ret
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected a panic due to transaction abort, didn't get one")
		} else if r.(error).Error() != "Can't set after commute" {
			t.Errorf("Expected 'can't set' error, got %v", r)

		}

		if v := r1.Deref(nil); v != init1 {
			t.Errorf("Expected r1 to have value %v, got %v", init1, v)
		}

		if !point1 {
			t.Errorf("Expected to do the commute, but failed earlier")
		}

		if point2 {
			t.Errorf("Expected to panic during the set, but made it through")
		}

	}()

	RunInTransaction(f)

	t.Errorf("Should not have reached this point, should have panicked")
}

func TestSimpleCommute(t *testing.T) {
	init1 := 12
	add1 := 10
	add2 := 20

	var save1, save2 interface{}

	r1 := NewRef(init1)
	ret := 300

	fc := func(old interface{}, args ...interface{}) interface{} {
		return old.(int) + args[0].(int)
	}

	f := func(tx *Tx) interface{} {
		r1.Commute(tx, fc, add1)
		save1 = r1.Deref(tx)
		r1.Commute(tx, fc, add2)
		save2 = r1.Deref(tx)
		return ret
	}

	v, e := RunInTransaction(f)

	if e != nil {
		t.Errorf("Expected no error, got %v", e)
	}

	if v != ret {
		t.Errorf("Expected return value of %v, got %v", ret, v)
	}

	if v := r1.Deref(nil); v != init1+add1+add2 {
		t.Errorf("Expected r1 to have value %v, got %v", init1+add1+add2, v)
	}

	if save1 != init1+add1 {
		t.Errorf("Expected save to have value %v, got %v", init1+add1, save1)
	}

	if save2 != init1+add1+add2 {
		t.Errorf("Expected save to have value %v, got %v", init1+add1+add2, save2)
	}
}

func TestSimpleInterference(t *testing.T) {

	runtime.GOMAXPROCS(4)

	init1 := 12
	val1 := 100
	val2 := 200
	add1 := 1

	r1 := NewRef(init1)
	ret := 300

	fc := func(old interface{}, args ...interface{}) interface{} {
		return old.(int) + args[0].(int)
	}

	var done sync.WaitGroup
	done.Add(2)

	var once sync.Once

	nfEnter, nfExit := 0, 0
	ngEnter, ngExit := 0, 0

	ch := make(chan bool, 1)

	signalG := func() {
		ch <- true
	}

	f := func(tx *Tx) interface{} {
		t.Logf("F: Entering")
		nfEnter++
		r1.Set(tx, val1)
		t.Logf("F: between sets")
		once.Do(signalG)
		time.Sleep(20 * time.Nanosecond)
		r1.Set(tx, val2)
		t.Logf("F: after sets")
		nfExit++
		return ret
	}

	go func() {
		defer done.Done()
		t.Logf("F: before TX")
		RunInTransaction(f)
		t.Logf("F: after TX")
	}()

	g := func(tx *Tx) interface{} {
		t.Logf("G: entering")
		ngEnter++
		r1.Commute(tx, fc, add1)
		t.Logf("G: after commute")
		ngExit++
		return ret
	}

	go func() {
		defer done.Done()
		<-ch
		t.Logf("G: before TX")
		RunInTransaction(g)
		t.Logf("G: after TX")
	}()

	time.Sleep(time.Nanosecond)
	done.Wait()

	t.Errorf("%v %v %v %v %v", ngEnter, ngExit, r1.Deref(nil), nfEnter, nfExit)

}
