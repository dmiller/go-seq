// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sequtil

import (
	"github.com/dmiller/go-seq/iseq"
	"reflect"
)

func DefaultCompareFn(k1 interface{}, k2 interface{}) int {
	if k1 == k2 {
		return 0
	}
	if k1 != nil {
		if k2 == nil {
			return 1
		}
		if c, ok := k1.(iseq.Comparer); ok {
			return c.Compare(k2)
		}
		if c, ok := k2.(iseq.Comparer); ok {
			return -c.Compare(k1)
		}
		if s, ok := k1.(string); ok {
			return CompareString(s, k2)
		}
		if IsComparableNumeric(k1) {
			return CompareComparableNumeric(k1, k2)
		}
		panic("Can't compare")
	}
	return -1
}

func IsComparableNumeric(v interface{}) bool {

	switch v.(type) {
	case bool, int, int8, int32, int64,
		uint, uint8, uint32, uint64,
		float32, float64:

		return true
	}
	return false
}

func CompareString(s string, x interface{}) int {
	if s2, ok := x.(string); ok {
		if s < s2 {
			return -1
		}
		if s == s2 {
			return 0
		}
		return 1
	}

	panic "can't compare string to non-string, non-iseq.Comparer"
}

func CompareComparableNumeric(x1 interface{}, x2 interface{}) int {
	// n1 should be numeric
	switch x1 := x1.(type) {
	case bool:
		b1 := bool(x1)
		if b1 {
			return compareNumericInt(int64(1), x2)
		} else {
			return compareNumericInt(int64(0), x2)
		}
	case int, int8, int32, int64:
		n1 := reflect.ValueOf(x1).Int()
		return compareNumericInt(n1, x2)
	case uint, uint8, uint32, uint64:
		n1 := reflect.ValueOf(x1).Uint()
		return compareNumericUint(n1, x2)
	case float32, float64:
		n1 := reflect.ValueOf(x1).Float()
		return compareNumericFloat(n1, x2)
	}
	panic("Expect first arg to be numeric")
}

func compareNumericInt(n1 int64, x2 interface{}) int {
	switch x2 := x2.(type) {
	case bool:
		b2 := bool(x2)
		var n2 int64
		if b2 {
			n2 = 1
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
		return 0

	case int, int8, int32, int64:
		n2 := reflect.ValueOf(x2).Int()
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
		return 0

	case uint, uint8, uint32, uint64:
		n2 := reflect.ValueOf(x2).Uint()
		if n1 < 0 {
			return -1
		}
		un1 := uint64(n2)
		if un1 < n2 {
			return -1
		}
		if un1 > n2 {
			return 1
		}
		return 0

	case float32, float64:
		n2 := reflect.ValueOf(x2).Float()
		fn1 := float64(n1)
		if fn1 < n2 {
			return -1
		}
		if fn1 > n2 {
			return 1
		}
		return 0
	}
	return -1 // what else, other than panic?
}

func compareNumericUint(n1 uint64, x2 interface{}) int {
	switch x2 := x2.(type) {
	case bool:
		b2 := bool(x2)
		var n2 uint64
		if b2 {
			n2 = 1
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
		return 0

	case int, int8, int32, int64:
		n2 := reflect.ValueOf(x2).Int()
		if n2 < 0 {
			return 1
		}
		un2 := uint64(n2)
		if n1 < un2 {
			return -1
		}
		if n1 > un2 {
			return 1
		}
		return 0

	case uint, uint8, uint32, uint64:
		n2 := reflect.ValueOf(x2).Uint()
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
		return 0

	case float32, float64:
		n2 := reflect.ValueOf(x2).Float()
		fn1 := float64(n1)
		if fn1 < n2 {
			return -1
		}
		if fn1 > n2 {
			return 1
		}
		return 0
	}
	return -1 // what else, other than panic?
}

func compareNumericFloat(n1 float64, x2 interface{}) int {
	var n2 float64
	switch x2 := x2.(type) {
	case bool, int, int8, int32, int64:
		n2 = float64(reflect.ValueOf(x2).Int())
	case uint, uint8, uint32, uint64:
		n2 = float64(reflect.ValueOf(x2).Uint())
	case float32, float64:
		n2 = reflect.ValueOf(x2).Float()
	default:
		return -1 // what else, other than panic?
	}
	if n1 < n2 {
		return -1
	}
	if n1 > n2 {
		return 1
	}
	return 0
}
