package seq

import (
	"iseq"
)

type arrayChunk struct {
	array  []interface{}
	offset int
	end    int
}

// c-tors

func newArrayChunk2(a []interface{}, offset int) *arrayChunk {
	return newArrayChunk3(a, offset, len(a))
}

func newArrayChunk3(a []interface{}, offset, end int) *arrayChunk {
	return &arrayChunk{array: a, offset: offset, end: end}
}

// ArrayChunk must implement interfaces
//   Counted Indexed Chunk

// interface Counted

func (a *arrayChunk) Count1() int {
	return a.end - a.offset
}

// interface Indexed

func (a *arrayChunk) Nth(i int) interface{} {
	return a.array[a.offset+i]
}

func (a *arrayChunk) NthD(i int, notFound interface{}) interface{} {
	if i >= 0 && i < a.Count1() {
		return a.Nth(i)
	}
	return notFound
}

// interface Chunk

func (a *arrayChunk) DropFirst() iseq.Chunk {
	if a.offset == a.end {
		panic("dropFirst of empty chunk")
	}

	return newArrayChunk3(a.array, a.offset+1, a.end)
}
