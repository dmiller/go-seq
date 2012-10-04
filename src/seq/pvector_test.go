package seq

import (
	"iseq"
	"testing"
)

//  PVector needs to implement the following seq interfaces:
//        Obj Meta Seqable PersistentCollection Lookup Associative PersistentStack PersistentVector Counted Reversible Indexed
//  Are we going to do EditableCollection?
//  Also, Equatable and Hashable
func TestPVectorImplementInterfaces(t *testing.T) {
	var c interface{} = NewPVectorFromItems("abc", "def")

	if _, ok := c.(iseq.Obj); !ok {
		t.Error("PList must implement Obj")
	}

	if _, ok := c.(iseq.Meta); !ok {
		t.Error("PList must implement Meta")
	}

	if _, ok := c.(iseq.PersistentCollection); !ok {
		t.Error("PList must implement PersistentCollection")
	}

	if _, ok := c.(iseq.PersistentStack); !ok {
		t.Error("PList must implement PersistentStack")
	}

	if _, ok := c.(iseq.Lookup); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Associative); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.PersistentVector); !ok {
		t.Error("PList must implement PersistentList")
	}

	if _, ok := c.(iseq.Seqable); !ok {
		t.Error("PList must implement Seqable")
	}

	if _, ok := c.(iseq.Counted); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Indexed); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Reversible); !ok {
		t.Error("PList must implement Counted")
	}

	if _, ok := c.(iseq.Equatable); !ok {
		t.Error("PList must implement Equatable")
	}

	if _, ok := c.(iseq.Hashable); !ok {
		t.Error("PList must implement Hashable")
	}
}

func TestPVectorISeqCtorWorks(t *testing.T) {
	var seq iseq.Seq = NewPListFromSlice([]interface{}{"def", 2, 3})
	v := NewPVectorFromISeq(seq)
	if v.Count() != 3 {
		t.Errorf("NewPVectorFromISeq has wrong count, expected %v, got %v", 3, v.Count())
	}
	for i, s := 0, seq; i < v.Count(); i, s = i+1, s.Next() {
		if v.Nth(i) != s.First() {
			t.Errorf("NewPVectorFromISeq: expected element %v = %v, found %v instead", i, v.Nth(i), s.First())
		}
	}
}

func makeRangePlist(size int) *PList {
	sl := make([]interface{}, size)
	for i := 0; i < len(sl); i++ {
		sl[i] = i + 10
	}

	return NewPListFromSlice(sl)
}

func TestPVectorISeqCtorWorksOnLargeSeq(t *testing.T) {
	sizes := []int{
		1000,   // this should get us out of the first node
		100000} // this should get us out of the second level

	for is := 0; is < len(sizes); is++ {
		size := sizes[is]
		pl := makeRangePlist(size)

		v := NewPVectorFromISeq(pl)
		if v.Count() != size {
			t.Errorf("NewPVectorFromISeq has wrong count, expected %v, got %v", size, v.Count())
		}
		for i, s := 0, pl.Seq(); i < v.Count(); i, s = i+1, s.Next() {
			if v.Nth(i) != s.First() {
				t.Errorf("NewPVectorFromISeq: expected element %v = %v, found %v instead", i, v.Nth(i), s.First())
			}
		}
	}
}

func TestPVectorSliceCtorWorks(t *testing.T) {
	sl := []interface{}{"def", 2, 3}
	v := NewPVectorFromSlice(sl)
	if v.Count() != 3 {
		t.Errorf("NewPVectorFromSlice has wrong count, expected %v, got %v", 3, v.Count())
	}
	for i := 0; i < v.Count(); i = i + 1 {
		if v.Nth(i) != sl[i] {
			t.Errorf("NewPVectorFromSlice: expected element %v = %v, found %v instead", i, v.Nth(i), sl[i])
		}
	}
}

func TestPVectorFromItemsWorks(t *testing.T) {
	sl := []interface{}{"def", 2, 3}
	v := NewPVectorFromItems("def", 2, 3)
	if v.Count() != 3 {
		t.Errorf("NewPVectorFromSlice has wrong count, expected %v, got %v", 3, v.Count())
	}
	for i := 0; i < v.Count(); i = i + 1 {
		if v.Nth(i) != sl[i] {
			t.Errorf("NewPVectorFromSlice: expected element %v = %v, found %v instead", i, v.Nth(i), sl[i])
		}
	}
}
