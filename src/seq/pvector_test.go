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
