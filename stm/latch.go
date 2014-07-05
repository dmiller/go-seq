// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stm

import (
	"sync"
	"sync/atomic"
	"time"
)

// Implements a count down latch with timeouts.
// For a more feature-full implementation, check out https://github.com/R358/brace, the inspiration for this implementation.
type CountDownLatch struct {
	count       int32
	zeroReached chan bool
	timedOut    int32 // really a bool, but need atomic access
	completed   int32 // really a bool, but need atomic access
	awaitMutex  *sync.Mutex
	countMutex  *sync.Mutex
}

// NewCountDownLatch returns a CountDownLatch with an initial count
func NewCountDownLatch(i int) *CountDownLatch {
	return &CountDownLatch{count: int32(i),
		zeroReached: make(chan bool, 1),
		timedOut:    0,
		completed:   0,
		awaitMutex:  new(sync.Mutex),
		countMutex:  new(sync.Mutex),
	}
}

// Completed returns true if the latch is completed (has already been awaited on)
func (c *CountDownLatch) Completed() bool {
	return atomic.LoadInt32(&c.completed) == 1
}

// TimedOut returns true if a latch await has timed out
func (c *CountDownLatch) TimedOut() bool {
	return atomic.LoadInt32(&c.timedOut) == 1
}

// CountDown decrements the latch count by 1
func (c *CountDownLatch) CountDown() {
	c.countMutex.Lock()
	defer c.countMutex.Unlock()

	if c.Completed() {
		panic("Latch already completed")
	}

	if c.count <= 0 {
		panic("Latch count already zero")
	}

	c.count--
	if c.count <= 0 {
		c.zeroReached <- true
	}
}

// Await waits up to the given amount of time for the latch to count down to zero.
// Returns true if timed-out, false if zero reached
func (c *CountDownLatch) Await(dur time.Duration) bool {
	c.awaitMutex.Lock()
	defer func() {
		atomic.StoreInt32(&c.completed, 1)
		c.awaitMutex.Unlock()
	}()

	if c.Completed() {
		panic("Latch already completed")
	}

	to := time.After(dur)

	select {
	case <-to:
		atomic.StoreInt32(&c.timedOut, 1)
		return true

	case <-c.zeroReached:
		atomic.StoreInt32(&c.completed, 1)
		return false
	}
}
