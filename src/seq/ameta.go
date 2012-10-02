package seq

import (
	"iseq"
)

// AMeta provides a slot to hold a 'meta' value
type AMeta struct {
	meta iseq.PersistentMap
}

func (o *AMeta) Meta() iseq.PersistentMap {
	return o.meta
}
