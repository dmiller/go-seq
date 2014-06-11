// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

import (
	"github.com/dmiller/go-seq/iseq"
)


// Equiv returns true if the objects are 'equivalent'.
// Two == objects are equivalent.
// Else, if either is iseq.Equivable, we default to that interface.
// Otherwise, not equivalent.
func Equiv(o1 interface{}, o2 interface{}) bool {
	if o1 == o2 {
		return true
	}

	// TODO: make sure we have handled interface nils
	if o1 == nil || o2 == nil
	{
		return false
	}

	if e1, ok := o1.(iseq.Equivable); ok {
		return e1.Equiv(o2)
	}

	if e2, ok := o2.(iseq.Equivable); ok {
		return e2.Equiv(o1)
	}

	return false
}

// MapEquiv returns true if its arguments are equivalent as maps.
// First argument is an iseq.PMap.
// Second argument must be convertible to a map. (Someeday we may handle go maps.)
// To be equivalent, must contain equivalent keys and values.
func MapEquiv(m1 iseq.PMap, obj interface{}) bool {
	if m1 == obj {
		return true
	}

	if _, ok := obj.(map[interface{}]interface{}); ok {
		// TODO: figure out how to handle go maps
		return false
	}

	if m2, ok := obj.(iseq.PMap); ok {
		if m1.Count() != m2.Count() {
			return false
		}

		for s := m1.Seq(); s != nil; s = s.Next() {
			me := s.First().(iseq.MapEntry)
			found := m2.ContainsKey(me.Key())
			if !found || !Equiv(me.Val(), m2.ValAt(me.Key())) {
				return false
			}
		}
		return true
	}
	return false
}


//  Returns true if the sequences are element-by-element equivalent.
func SeqEquiv(s1 iseq.Seq, s2 iseq.Seq) bool {
	if s1 == s2 {
		return true
	}

	iter2 := s2

	for iter1 := s1; iter1 != nil; iter1 = iter1.Next() {
		if iter2 == nil || !Equiv(iter1.First(), iter2.First()) {
			return false
		}
		iter2 = iter2.Next()
	}

	return iter2 == nil
}

