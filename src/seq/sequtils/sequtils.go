// Package sequtils contains some generally useful functions for implementing the Clojure sequences
package sequtils

import (
	"seq"
	//"fmt"
)

func Equals(o1 interface{}, o2 interface{}) bool {
	if o1 == o2 {
		return true
	}

	if e1, ok1 := o1.(seq.Equatable); ok1 {
		return e1.Equals(o2)
	}

	return false
}

func SeqEquals(s1 seq.Seq, s2 seq.Seq) bool {
	if ( s1 == s2 ) {
		return true
	}

	iter2 := s2

	for iter1 := s1; iter1 != nil; iter1 = iter1.Next() {
		if iter2 == nil || ! Equals(iter1.First(),iter2.First()) {
			return false;
		}
		iter2 = iter2.Next() 
	}

	return iter2 == nil
}

func Equiv(o1 interface{}, o2 interface{} ) bool {
	if (o1 == o2) {
		return true
	}
	if ( o1 != nil ) {
		// TODO: Determine how to handle numbers. Do we want Clojure's semantics?
		// Go's semantics says the o1 == o2 case is enough
			pc1, ok1 := o1.(seq.PersistentCollection)
			if ok1 {
				return pc1.Equiv(o2)
			}

			pc2, ok2 := o2.(seq.PersistentCollection)
			if ok2 {
				return pc2.Equiv(o1)
			}

			return Equals(o1,o2)
	}

	return false
}

func SeqEquiv(s1 seq.Seq, s2 seq.Seq) bool {
	if ( s1 == s2 ) {
		return true
	}

	iter2 := s2

	for iter1 := s1; iter1 != nil; iter1 = iter1.Next() {
		if iter2 == nil || ! Equiv(iter1.First(),iter2.First()) {
			return false;
		}
		iter2 = iter2.Next() 
	}

	return iter2 == nil
}

func Count(o interface{}) int {
	if ( o == nil ) {
		return 0
	}

	if cnt, ok := o.(seq.Counted); ok {
		return cnt.CountFast()
	}

	if pc, ok := o.(seq.PersistentCollection); ok {
		s := pc.Seq()
		i := 0
		for ; s!=nil; s = s.Next() {
			if c,ok := s.(seq.Counted); ok {
				return i + c.Count()
			}
			i++
		}
		return i
	}

	if s, ok := o.(string); ok {
		return len(s)
	}
	// TODO: Figure out how to  handle arrays, slices, maps in a typeswitch/generic way
	panic("Count not supported on this type")
}
