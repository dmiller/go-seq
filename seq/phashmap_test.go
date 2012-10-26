package seq

import (
	"fmt"
	"github.com/dmiller/go-seq/iseq"
	"math/rand"
	"testing"
)

func TestPHashMapImplementInterfaces(t *testing.T) {
	var c interface{} = NewPHashMapFromItems("abc", "def")

	if _, ok := c.(iseq.Obj); !ok {
		t.Error("PHashMap must implement Obj")
	}

	if _, ok := c.(iseq.Meta); !ok {
		t.Error("PHashMap must implement Meta")
	}

	if _, ok := c.(iseq.PCollection); !ok {
		t.Error("PHashMap must implement PCollection")
	}

	if _, ok := c.(iseq.PMap); !ok {
		t.Error("PHashMap must implement PMap")
	}

	if _, ok := c.(iseq.Lookup); !ok {
		t.Error("PHashMap must implement Counted")
	}

	if _, ok := c.(iseq.Associative); !ok {
		t.Error("PHashMap must implement Counted")
	}

	if _, ok := c.(iseq.Seqable); !ok {
		t.Error("PHashMap must implement Seqable")
	}

	if _, ok := c.(iseq.Counted); !ok {
		t.Error("PHashMap must implement Counted")
	}

	if _, ok := c.(iseq.Equatable); !ok {
		t.Error("PHashMap must implement Equatable")
	}

	if _, ok := c.(iseq.Hashable); !ok {
		t.Error("PHashMap must implement Hashable")
	}
}


// Factory tests

func TestPHashMapISeqFactoryWorks(t *testing.T) {
	var seq iseq.Seq = NewPListFromSlice([]interface{}{"def", 2, "abc", 3, "pqr", 7})
	m := NewPHashMapFromSeq(seq)
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSeq has wrong count, expected %v, got %v", 3, m.Count())
	}

	for s := seq; s != nil; s = s.Next().Next() {
		if m.ValAt(s.First()) != s.Next().First() {
			t.Errorf("NewPHashMapFromSeq: expected key %v => %v, found %v instead", s.First(), s.Next().First(), m.ValAt(s.First()))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSeq: found key that should not be there")
	}
}

func TestPHashMapISeqFactoryWorksWithDupKey(t *testing.T) {
	var seq iseq.Seq = NewPListFromSlice([]interface{}{"def", 2, "abc", 3, "pqr", 7})
	m := NewPHashMapFromSeq(NewCons("def",NewCons("99",seq)))
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSeq has wrong count, expected %v, got %v", 3, m.Count())
	}

	for s := seq; s != nil; s = s.Next().Next() {
		if m.ValAt(s.First()) != s.Next().First() {
			t.Errorf("NewPHashMapFromSeq: expected key %v => %v, found %v instead", s.First(), s.Next().First(), m.ValAt(s.First()))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSeq: found key that should not be there")
	}
}

func TestPHashMapISeqFactoryOnEmpty(t *testing.T) {
	m := NewPHashMapFromSeq(nil)
	if m.Count() != 0  {
		t.Errorf("NewPHashMapFromSeq: on nil, should have count 0, got %v",m.Count())
	}
}


func TestPHashMapSliceFactoryWorks(t *testing.T) {
	s := []interface{}{"def",2,"abc",3,"pqr",7}
	m := NewPHashMapFromSlice(s)
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSlice has wrong count, expected %v, got %v", 3, m.Count())
	}

	for i := 0; i < len(s); i += 2 {
		if m.ValAt(s[i]) != s[i+1] {
			t.Errorf("NewPHashMapFromSlice: expected key %v => %v, found %v instead", s[i], s[i+1], m.ValAt(s[i]))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSlice: found key that should not be there")
	}
}

func TestPHashMapSliceFactoryWorksWithDupKey(t *testing.T) {
	s4 := []interface{}{"def",2,"abc",3,"pqr",7, "def", 99} 
	s3 := s4[2:]
	m := NewPHashMapFromSlice(s4)
	if m.Count() != 3 {
		t.Errorf("NewPHashMapFromSlice has wrong count, expected %v, got %v", 3, m.Count())
	}

	for i := 0; i < len(s3); i += 2 {
		if m.ValAt(s3[i]) != s3[i+1] {
			t.Errorf("NewPHashMapFromSlice: expected key %v => %v, found %v instead", s3[i], s3[i+1], m.ValAt(s3[i]))
		}
	}

	if m.ContainsKey("xyz") {
		t.Errorf("NewPHashMapFromSlice: found key that should not be there")
	}
}

func TestPHashMapSliceFactoryOnEmpty(t *testing.T) {
	m := NewPHashMapFromSlice([]interface{}{})
	if m.Count() != 0  {
		t.Errorf("NewPHashMapFromSlice: on nil, should have count 0, got %v",m.Count())
	}
}

// TODO: Eventually, move this to a benchmark
func createBigSliceForPHashMapTest(n int) []interface{} {
	s := make([]interface{},2*n)
	for i := 0; i<2*n ; i+=2 {
		r := rand.Int63()
		s[i] = fmt.Sprintf("%d",r)
		s[i+1] = r
	}
	return s
}

func TestPHashMapGoesBig(t *testing.T) {
	sizes := []int{10,100,1000,10000,100000}
	for _,n := range sizes {
		s := createBigSliceForPHashMapTest(n)
		//fmt.Printf("Testing big PHashMap creation: %v items\n",n)
		m := NewPHashMapFromSlice(s)
		if m.Count() != n {
			t.Errorf("NewPHashMapFromSlice has wrong count, expected %v, got %v", n, m.Count())
		}

		for i := 0; i < len(s); i += 2 {
			if m.ValAt(s[i]) != s[i+1] {
				t.Errorf("NewPHashMapFromSlice: expected key %v => %v, found %v instead", s[i], s[i+1], m.ValAt(s[i]))
				break
			}
		}
	}
}

// func TestPHashMapHasProblem(t *testing.T) {
// 	var m iseq.PMap = EmptyPHashMap
// 	for i := 0; i<100; i++ {
// 		m = m.AssocM(i,i*10)
// 		fmt.Printf("%v\n",i)
// 		if m.Count() != i+1 {
// 			t.Errorf("Count is %v, should be %v\n",m.Count(),i+1)
// 		}
// 	}
// }
