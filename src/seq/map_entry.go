package seq

type MapEntry struct {
	key interface{}
	val interface{}
}

func (me MapEntry) Key() interface{} {
	return me.key
}

func (me MapEntry) Val() interface{} {
	return me.val
}
