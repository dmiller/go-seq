package seqimpl

import "seq"

// AMeta provides a slot to hold a 'meta' value
type AMeta struct {
	meta seq.PersistentMap
}

func (o *AMeta) Meta() seq.PersistentMap {
	return o.meta
}
