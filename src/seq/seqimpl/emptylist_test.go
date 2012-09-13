package seqimpl
import (
	"testing"
	"seq"
)


func TestEmptyListImplementInterfaces(t *testing.T) {
	var c interface{} = CachedEmptyList

	if _,ok := c.(seq.Obj); !ok {
		t.Error("EmptyList must implement Obj")
	}

	if _,ok := c.(seq.Meta); !ok {
		t.Error("EmptyList must implement Meta")
	}

	if _,ok := c.(seq.PersistentCollection); !ok {
		t.Error("EmptyList must implement PersistentCollection")
	}

	if _,ok := c.(seq.PersistentStack); !ok {
		t.Error("EmptyList must implement PersistentStack")
	}

	if _,ok := c.(seq.PersistentList); !ok {
		t.Error("EmptyList must implement PersistentList")
	}

	if _,ok := c.(seq.Seqable); !ok {
		t.Error("EmptyList must implement Seqable")
	}

	if _,ok := c.(seq.Counted); !ok {
		t.Error("EmptyList must implement Counted")
	}

	if _,ok := c.(seq.Equatable); !ok {
		t.Error("EmptyList must implement Equatable")
	}

	if _,ok := c.(seq.Hashable); !ok {
		t.Error("EmptyList must implement Hashable")
	}
}


func TestEmptyListCount(t *testing.T) {
	c := CachedEmptyList
	if c.Count() != 0 {
		t.Errorf("Count: expected 0, got %v",c.Count())
	}

	if c.Count1() != 0 {
		t.Errorf("Count1: expected 0, got %v",c.Count1())
	}
}

func TestEmptyListSeq(t *testing.T) {
	c := CachedEmptyList
	if ( c.Seq() != nil ) {
		t.Error("Seq of EmptyList should be nil")
	}
}

func TestEmptyListEmpty(t *testing.T) {
	c := CachedEmptyList
	e := c.Empty()
	if e !=c {
		t.Error("Empty should be self")
	}

	c1 := &EmptyList{}
	e1 := c1.Empty()
	if e1 != c1 {
		t.Error("Empty should be self")
	}		
}

func TestEmptyListEquiv(t *testing.T) {
	c1 := CachedEmptyList
	c2 := &EmptyList{}
	if ! c1.Equiv(c1) {
		t.Error("Expect empty list to be equiv to itself")
	}
	 if ! c1.Equiv(c2) {
	 	t.Error("Expect empty list to equiv another empty list")
	 }

	// c3 := NewCons("abc",nil)
	// if c1.Equiv(c3) {
	// 	t.Error("Expect empty list not equiv to a non-empty list")
	// }
}

