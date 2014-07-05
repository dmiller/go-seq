// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"testing"
	"time"
)

func TestNewLatch(t *testing.T) {
	c := NewCountDownLatch(10)

	if c.Completed() {
		t.Error("New latch should not be completed")
	}

	if c.TimedOut() {
		t.Error("New latch should not be timed out")
	}
}

func TestAwaitWithTimeout(t *testing.T) {
	c := NewCountDownLatch(1)

	v := c.Await(time.Second)

	if !v {
		t.Error("Await that timed out should return true")
	}

	if !c.Completed() {
		t.Error("After timed-out await, latch should be completed")
	}

	if !c.TimedOut() {
		t.Error("After timed-out await, latch should be timed out")
	}

}

func TestAwaitAlreadyCountedDown(t *testing.T) {
	c := NewCountDownLatch(1)
	c.CountDown()

	v := c.Await(time.Second)

	if v {
		t.Error("Await that succeeded on entry should return false")
	}

	if !c.Completed() {
		t.Error("After Await that succeeded on entry, latch should be completed")
	}

	if c.TimedOut() {
		t.Error("After Await that succeeded on entry, latch should not be timed out")
	}
}

func TestAwaitCountDown(t *testing.T) {
	c := NewCountDownLatch(1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		c.CountDown()
	}()

	t0 := time.Now()
	v := c.Await(100 * time.Second)
	t1 := time.Now()

	if v {
		t.Error("Await that succeeded should return false")
	}

	dur := t1.Sub(t0)
	if dur/time.Millisecond > 1000 {
		t.Errorf("Time should be under 1000 msecs, was %v", dur)
	}
}

func TestAwaitCountDown2(t *testing.T) {
	c := NewCountDownLatch(2)

	go func() {
		time.Sleep(100 * time.Millisecond)
		c.CountDown()
		time.Sleep(100 * time.Millisecond)
		c.CountDown()
	}()

	t0 := time.Now()
	v := c.Await(100 * time.Second)
	t1 := time.Now()

	if v {
		t.Error("Await that succeeded should return false")
	}

	dur := t1.Sub(t0)
	if dur/time.Millisecond > 1000 {
		t.Errorf("Time should be under 1000 msecs, was %v", dur)
	}
}

func TestCallAwaitTwice(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Await second call should panic")
		}
		if r != nil && r != "Latch already completed" {
			t.Errorf("Await second call should panic with 'completed' message, got %v", r)
		}
	}()

	c := NewCountDownLatch(2)
	c.Await(1 * time.Millisecond)
	c.Await(1 * time.Millisecond)
}

func TestCountDownTooFar(t *testing.T) {

	n := 0

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("CountDown past zero should panic")
		}
		if r != nil && r != "Latch count already zero" {
			t.Errorf("CountDown past zero should panic with 'zero count' message, got %v", r)
		}
		if n != 2 {
			t.Errorf("Exactly two CountDowns should have occurred, actual count = %d", n)
		}
	}()

	c := NewCountDownLatch(2)
	c.CountDown()
	n++
	c.CountDown()
	n++
	c.CountDown()
	n++
}
