package seq

import (
	"iseq"
	"sequtil"
)

// chunkedSeq provides a chunked seq over a PVector
type chunkedSeq struct {
	vec    *PVector
	node   []interface{}
	idx    int
	offset int
	AMeta
}

//  chunkedSeq needs to implement the following iseq interfaces:
//        Obj Meta Seq IPersistentCollection Sequential ChunkedSeq
//  Also, Equatable and Hashable

// c-tors

func newChunkedSeq(v *PVector, i int, offset int) *chunkedSeq {
	return newChunkedSeqM(nil, v, v.arrayFor(i), i, offset)
}

func newChunkedSeqM(meta iseq.PersistentMap, v *PVector, node []interface{}, i int, offset int) *chunkedSeq {
	return &chunkedSeq{AMeta: AMeta{meta}, vec: v, node: node, idx: i, offset: offset}
}

func newChunkedSeqRaw(v *PVector, node []interface{}, i int, offset int) *chunkedSeq {
	return newChunkedSeqM(nil, v, node, i, offset)
}

// interface Obj

func (c *chunkedSeq) WithMeta(meta iseq.PersistentMap) iseq.Obj {
	if meta == c.meta {
		return c
	}
	return newChunkedSeqM(meta, c.vec, c.node, c.idx, c.offset)
}

// interface ChunkedSeq

func (c *chunkedSeq) ChunkedFirst() iseq.Chunk {
	return newArrayChunk2(c.node, c.offset)
}

func (c *chunkedSeq) ChunkedNext() iseq.Seq {
	if c.idx+len(c.node) < c.vec.cnt {
		return newChunkedSeq(c.vec, c.idx+len(c.node), 0)
	}
	return nil
}

func (c *chunkedSeq) ChunkedMore() iseq.Seq {
	s := c.ChunkedNext()
	if s == nil {
		return CachedEmptyList
	}
	return s
}

// interface PersistentCollection

func (c *chunkedSeq) Seq() iseq.Seq {
	return c
}

func (c *chunkedSeq) Count() int {
	// TODO: find out where this is duplicated; originally from ASeq.cs
	i := 1 // if we are here, it is non-empty
	for s := c.Next(); s != nil; s, i = s.Next(), i+1 {
		if cnt, ok := s.(iseq.Counted); ok {
			return i + cnt.Count1()
		}
	}
	return i
}

func (c *chunkedSeq) Cons(o interface{}) iseq.PersistentCollection {
	return NewCons(o, c)
}

func (c *chunkedSeq) Empty() iseq.PersistentCollection {
	return CachedEmptyList
}

func (c *chunkedSeq) Equiv(o interface{}) bool {
	// TODO: revisit Equiv
	return sequtil.Equals(c, o)
}

func (c *chunkedSeq) First() interface{} {
	return c.node[c.offset]
}

func (c *chunkedSeq) Next() iseq.Seq {
	if c.offset+1 < len(c.node) {
		return newChunkedSeqRaw(c.vec, c.node, c.idx, c.offset+1)
	}
	return c.ChunkedNext()
}

func (c *chunkedSeq) More() iseq.Seq {
	s := c.Next()
	if c == nil {
		return CachedEmptyList
	}
	return s
}

func (c *chunkedSeq) SCons(o interface{}) iseq.Seq {
	return NewCons(o, c)
}
