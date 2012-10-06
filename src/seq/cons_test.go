package seq

import (
	"iseq"
	"testing"
)

func TestConsCtors(t *testing.T) {
	c := NewCons("abc", nil)

	if c.Meta() != nil {
		t.Error("NewCons ctor should have nil meta")
	}

	if c.First() != "abc" {
		t.Error("NewCons ctor did not initialize first")
	}

	c1 := NewCons("def", c)
	if c1.First() != "def" {
		t.Error("NewCons ctor did not initialize first")
	}

	if c1.Next() != c {
		t.Error("NewCons ctor did nto initialize more/next")
	}

	// TODO: add tests for c-tor with meta -- we need a PMap implementation first
}

func TestConsImplementInterfaces(t *testing.T) {
	var c interface{} = NewCons("abc", nil)

	if _, ok := c.(iseq.Obj); !ok {
		t.Error("Cons must implement Obj")
	}

	if _, ok := c.(iseq.Meta); !ok {
		t.Error("Cons must implement Meta")
	}

	if _, ok := c.(iseq.PCollection); !ok {
		t.Error("Cons must implement PCollection")
	}

	if _, ok := c.(iseq.Seqable); !ok {
		t.Error("Cons must implement Seqable")
	}

	if _, ok := c.(iseq.Equatable); !ok {
		t.Error("Cons must implement Equatable")
	}

	if _, ok := c.(iseq.Hashable); !ok {
		t.Error("Cons must implement Hashable")
	}
}

func createComplicatedCons() *Cons {
	c1 := NewCons(1, nil)
	c2 := NewCons(2, c1)
	c3 := NewCons("abc", nil)
	c4 := NewCons(c3, c2)
	c5 := NewCons("def", c4)
	return c5
}

func TestConsCount(t *testing.T) {
	c := createComplicatedCons()
	if c.Count() != 4 {
		t.Errorf("Count: expected 4, got %v", c.Count())
	}
}

func TestConsSeq(t *testing.T) {
	c1 := NewCons("abc", nil)
	c2 := createComplicatedCons()
	if c1.Seq() != c1 {
		t.Error("Seq should return self")
	}
	if c2.Seq() != c2 {
		t.Error("Seq should return self")
	}
}

func TestConsEmpty(t *testing.T) {
	c := NewCons("abc", nil)
	e := c.Empty()
	if e != CachedEmptyList {
		t.Error("Empty should be  CachedEmptyList")
	}
}

func TestConsEquiv(t *testing.T) {
	c1 := createComplicatedCons()
	c2 := createComplicatedCons()
	if c1 == c2 {
		t.Error("Expect two calls to createComplicatedCons to return distinct structs")
	}
	if !c1.Equiv(c1) {
		t.Error("Expect cons to be equiv to itself")
	}
	if !c1.Equiv(c2) {
		t.Error("Expect cons to equiv similar cons")
	}

	c3 := NewCons("abc", nil)
	if c1.Equiv(c3) {
		t.Error("cons equiv dissimilar cons")
	}
}
