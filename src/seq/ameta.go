package seq

import (
	"iseq"
)

// AMeta provides a slot to hold a 'meta' value
type AMeta struct {
	meta iseq.PMap
}

func (o *AMeta) Meta() iseq.PMap {
	return o.meta
}
