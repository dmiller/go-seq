package seqimpl

import "seq"


type AMeta struct {
	meta seq.PersistentMap
}

func (o *AMeta) Meta() seq.PersistentMap {
	return o.meta
}

