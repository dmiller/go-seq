// Copyright 2014 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"fmt"

	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
)

// PHashMap is a persistent rendition of Phil Bagwell's Hash Array Mapped Trie.
//
// The modifications to make this persistent are by Rich Hickey.
// The errors are mine.
//
// Rich's notes in the Clojure/Java implementation:
//
// Uses path copying for persistence.
// HashCollision leaves vs extended hashing.
// Node polymorphism vs conditionals
// No sub-tree pools or root-resizing.
//
type PHashMap struct {
	count    int
	root     hmnode
	hasNil   bool
	nilValue interface{}
	AMeta
	hash uint32
}

var (
	EmptyPHashMap    = &PHashMap{}
	phmNotFoundValue = new(interface{})
)

// factories

// TODO: need a factory for creating from an arbitrary Go map

func NewPHashMapFromSeq(items iseq.Seq) *PHashMap {
	// TODO: transients
	ret := EmptyPHashMap

	for i := 0; items != nil; items, i = items.Next().Next(), i+1 {
		if items.Next() == nil {
			panic(fmt.Sprintf("No value supplied for key: %v", items.First()))
		}
		ret = ret.AssocM(items.First(), items.Next().First()).(*PHashMap)
		// if checkDup && ret.Count1() != i+1 {
		// 	panic(fmt.Sprintf("Duplicate key: %v",items.First()))
		// }
	}
	return ret
}

func NewPHashMapFromSlice(s []interface{}) *PHashMap {
	// TODO: transients
	ret := EmptyPHashMap
	for i := 0; i < len(s); i = i + 2 {
		ret = ret.AssocM(s[i], s[i+1]).(*PHashMap)
		// if checkDup && ret.Count1() != i+1 {
		// 	panic(fmt.Sprintf("Duplicate key: %v",s[i]))
		// }
	}
	return ret
}

func NewPHashMapFromItems(items ...interface{}) *PHashMap {
	return NewPHashMapFromSlice(items)
}

// PHashMap needs to implement the following iseq interfaces:
//	Meta MetaW Seqable PCollection Lookup Associative Counted PMap
//  Are we going to do EditableCollection?
//  Also, Equivable and Hashable
//
// interface Meta is covered by the AMeta embedding
// TODO: IEditableCollection

// interface iseq.MetaW

func (m *PHashMap) WithMeta(meta iseq.PMap) iseq.MetaW {
	return &PHashMap{AMeta: AMeta{meta},
		count:    m.count,
		root:     m.root,
		hasNil:   m.hasNil,
		nilValue: m.nilValue}
}

// interface iseq.Associative, iseq.Lookup

func (m *PHashMap) ContainsKey(key interface{}) bool {
	if key == nil {
		return m.hasNil
	}

	if m.root == nil {
		return false
	}

	return m.root.findD(0, Hash(key), key, phmNotFoundValue) != phmNotFoundValue
}

func (m *PHashMap) EntryAt(key interface{}) iseq.MapEntry {
	if key == nil {
		if m.hasNil {
			return &MapEntry{nil, m.nilValue}
		}
		return nil
	}
	if m.root == nil {
		return nil
	}

	return m.root.find(0, Hash(key), key)
}

func (m *PHashMap) Assoc(key interface{}, val interface{}) iseq.Associative {
	return m.AssocM(key, val)
}

func (m *PHashMap) ValAt(key interface{}) interface{} {
	return m.ValAtD(key, nil)
}

func (m *PHashMap) ValAtD(key interface{}, notFound interface{}) interface{} {
	if key == nil {
		if m.hasNil {
			return m.nilValue
		}
		return notFound
	}

	if m.root == nil {
		return notFound
	}

	return m.root.findD(0, Hash(key), key, notFound)
}

// interface iseq.PMap

func (m *PHashMap) AssocM(key interface{}, val interface{}) iseq.PMap {
	if key == nil {
		if m.hasNil && val == m.nilValue {
			return m
		}
		newCount := m.count
		if !m.hasNil {
			newCount++
		}
		return &PHashMap{AMeta: AMeta{m.Meta()},
			count:    newCount,
			root:     m.root,
			hasNil:   true,
			nilValue: val}
	}
	newRoot := m.root
	if newRoot == nil {
		newRoot = emptyBitmapIndexedHmnode
	}
	newRoot, addedLeaf := newRoot.assoc2(0, Hash(key), key, val)
	if newRoot == m.root {
		return m
	}
	newCount := m.count
	if addedLeaf {
		newCount++
	}
	return &PHashMap{AMeta: AMeta{m.Meta()},
		count:    newCount,
		root:     newRoot,
		hasNil:   m.hasNil,
		nilValue: m.nilValue}
}

func (m *PHashMap) Without(key interface{}) iseq.PMap {
	if key == nil {
		if m.hasNil {
			return &PHashMap{AMeta: AMeta{m.Meta()},
				count:    m.count - 1,
				root:     m.root,
				hasNil:   false,
				nilValue: nil}
		}
		return m
	}
	newRoot := m.root.without(0, Hash(key), key)
	if newRoot == m.root {
		return m
	}
	return &PHashMap{AMeta: AMeta{m.Meta()},
		count:    m.count - 1,
		root:     newRoot,
		hasNil:   m.hasNil,
		nilValue: m.nilValue}
}

func (m *PHashMap) ConsM(e iseq.MapEntry) iseq.PMap {
	return sequtil.MapCons(m, e)
}

// interface iseq.PCollection, iseq.Seqable, iseq.Counted

func (m *PHashMap) Count() int {
	return m.count
}

func (m *PHashMap) Count1() int {
	return m.count
}

func (m *PHashMap) Cons(o interface{}) iseq.PCollection {
	return sequtil.MapCons(m, o)
}

func (m *PHashMap) Empty() iseq.PCollection {
	return EmptyPHashMap.WithMeta(m.Meta()).(iseq.PCollection)
}

func (m *PHashMap) Seq() iseq.Seq {
	var s iseq.Seq
	if m.root != nil {
		s = m.root.getNodeSeq()
	}
	if m.hasNil {
		return NewCons(MapEntry{nil, m.nilValue}, s)
	}
	return s
}

// interfaces Equivable, Hashable

func (m *PHashMap) Equiv(o interface{}) bool {
	return sequtil.MapEquiv(m, o)
}

func (p *PHashMap) Hash() uint32 {
	if p.hash == 0 {
		p.hash = sequtil.HashMap(p)
	}
	return p.hash
}

// Nodes in the trie

type hmnode interface {
	assoc(shift uint32, hash uint32, key interface{}, val interface{}) hmnode
	assoc2(shift uint32, hash uint32, key interface{}, val interface{}) (hmnode, bool)
	without(shift uint32, hash uint32, key interface{}) hmnode
	find(shift uint32, hash uint32, key interface{}) iseq.MapEntry
	findD(shift uint32, hash uint32, key interface{}, notFound interface{}) interface{}
	getNodeSeq() iseq.Seq
	//getHash() uint32 -- in the Java code, but does not appear to be used
}

// Slice manipulation

func cloneAndSetNodeSlice(src []hmnode, i int, a hmnode) []hmnode {
	clone := make([]hmnode, len(src))
	copy(clone, src)
	clone[i] = a
	return clone
}

func cloneAndSetObjectSlice(src []interface{}, i int, a interface{}) []interface{} {
	clone := make([]interface{}, len(src))
	copy(clone, src)
	clone[i] = a
	return clone
}

func cloneAndSetObjectSlice2(src []interface{}, i int, a interface{}, j int, b interface{}) []interface{} {
	clone := make([]interface{}, len(src))
	copy(clone, src)
	clone[i] = a
	clone[j] = b
	return clone
}

func removePair(src []interface{}, i int) []interface{} {
	dest := make([]interface{}, len(src)-2)
	copy(dest, src[:2*i])
	copy(dest[2*i:], src[2*(i+1):])
	return dest
}

// Node factories

func createNode(shift uint32, key1 interface{}, val1 interface{}, key2hash uint32, key2 interface{}, val2 interface{}) hmnode {
	key1hash := Hash(key1)
	if key1hash == key2hash {
		return &hashCollisionHmnode{key1hash, 2, []interface{}{key1, val1, key2, val2}}
	}
	node := emptyBitmapIndexedHmnode.assoc(shift, key1hash, key1, val1).assoc(shift, key2hash, key2, val2)
	return node
}

func imask(hash uint32, shift uint32) int {
	return int(mask(hash, shift))
}

func mask(hash uint32, shift uint32) uint32 {
	return (hash >> shift) & 0x01f
}

func bitpos(hash uint32, shift uint32) uint32 {
	return 1 << mask(hash, shift)
}

// arrayHmnode implements node of key-value pairs in the trie
//
// Must implement interface hmnode
type arrayHmnode struct {
	count int
	array []hmnode
}

func (a *arrayHmnode) assoc(shift uint32, hash uint32, key interface{}, val interface{}) hmnode {
	node, _ := a.assoc2(shift, hash, key, val)
	return node
}

func (a *arrayHmnode) assoc2(shift uint32, hash uint32, key interface{}, val interface{}) (hmnode, bool) {
	idx := imask(hash, shift)
	node := a.array[idx]
	if node == nil {
		newNode, addedLeaf := emptyBitmapIndexedHmnode.assoc2(shift+5, hash, key, val)
		return &arrayHmnode{a.count + 1, cloneAndSetNodeSlice(a.array, idx, newNode)}, addedLeaf
	}
	anode, addedLeaf := node.assoc2(shift+5, hash, key, val)
	if anode == node {
		return a, addedLeaf
	}
	return &arrayHmnode{a.count, cloneAndSetNodeSlice(a.array, idx, anode)}, addedLeaf
}

func (a *arrayHmnode) without(shift uint32, hash uint32, key interface{}) hmnode {
	idx := imask(hash, shift)
	node := a.array[idx]
	if node == nil {
		return a
	}
	n := node.without(shift+5, hash, key)
	if n == node {
		return a
	}
	if n == nil {
		if a.count <= 8 { // shrink
			return a.pack(idx)
		}
		return &arrayHmnode{a.count - 1, cloneAndSetNodeSlice(a.array, idx, n)}
	}
	return &arrayHmnode{a.count, cloneAndSetNodeSlice(a.array, idx, n)}
}

func (a *arrayHmnode) find(shift uint32, hash uint32, key interface{}) iseq.MapEntry {
	idx := imask(hash, shift)
	node := a.array[idx]
	if node == nil {
		return nil
	}
	return node.find(shift+5, hash, key)
}

func (a *arrayHmnode) findD(shift uint32, hash uint32, key interface{}, notFound interface{}) interface{} {
	idx := imask(hash, shift)
	node := a.array[idx]
	if node == nil {
		return notFound
	}
	return node.findD(shift+5, hash, key, notFound)

}

func (a *arrayHmnode) getNodeSeq() iseq.Seq {
	return createArrayHmnodeSeq(nil, a.array, 0, nil)
}

// func (a *arrayHmnode) getHash() uint32 {
// }

func (a *arrayHmnode) pack(idx int) hmnode {
	newArray := make([]interface{}, 2*(a.count-1))
	j := 1
	var bitmap uint32 = 0
	//TODO: change these to range iterations over subslices
	for i := 0; i < idx; i++ {
		if a.array[i] != nil {
			newArray[j] = a.array[i]
			bitmap = bitmap | 1<<uint32(i)
			j = j + 2
		}
	}
	for i := idx + 1; i < len(a.array); i++ {
		if a.array[i] != nil {
			newArray[j] = a.array[i]
			bitmap = bitmap | 1<<uint32(i)
			j = j + 2
		}
	}
	return &bitmapIndexedHmnode{bitmap, newArray}
}

//  Seq implementation for arrayHmnode

type arrayHmnodeSeq struct {
	nodes []hmnode
	i     int
	s     iseq.Seq
	AMeta
}

func createArrayHmnodeSeq(meta iseq.PMap, nodes []hmnode, i int, s iseq.Seq) *arrayHmnodeSeq {
	if s != nil {
		return &arrayHmnodeSeq{AMeta: AMeta{meta}, nodes: nodes, i: i, s: s}
	}
	for j := i; j < len(nodes); j++ {
		if nodes[j] != nil {
			ns := nodes[j].getNodeSeq()
			if ns != nil {
				return &arrayHmnodeSeq{AMeta: AMeta{meta}, nodes: nodes, i: j + 1, s: ns}
			}
		}
	}
	return nil
}

func (a *arrayHmnodeSeq) WithMeta(meta iseq.PMap) iseq.MetaW {
	return createArrayHmnodeSeq(meta, a.nodes, a.i, a.s)
}

func (a *arrayHmnodeSeq) First() interface{} {
	return a.s.First()
}

func (a *arrayHmnodeSeq) Next() iseq.Seq {
	return createArrayHmnodeSeq(nil, a.nodes, a.i, a.s.Next())
}

func (a *arrayHmnodeSeq) Cons(o interface{}) iseq.PCollection {
	return NewCons(o, a)
}

func (a *arrayHmnodeSeq) ConsS(o interface{}) iseq.Seq {
	return NewCons(o, a)
}

func (a *arrayHmnodeSeq) More() iseq.Seq {
	return moreFromSeq(a)
}

func (a *arrayHmnodeSeq) Count() int {
	return sequtil.SeqCount(a)
}

func (a *arrayHmnodeSeq) Empty() iseq.PCollection {
	return CachedEmptyList
}

func (a *arrayHmnodeSeq) Equiv(o interface{}) bool {
	return sequtil.Equiv(a, o)
}

func (a *arrayHmnodeSeq) Seq() iseq.Seq {
	return a
}

// bitmapIndexedHmnode represents an internal node in the trie, not full.
type bitmapIndexedHmnode struct {
	bitmap uint32
	array  []interface{}
}

var emptyBitmapIndexedHmnode = &bitmapIndexedHmnode{}

func (b *bitmapIndexedHmnode) index(bit uint32) int {
	return sequtil.BitCountU32(b.bitmap & (bit - 1))
}

func (b *bitmapIndexedHmnode) assoc(shift uint32, hash uint32, key interface{}, val interface{}) hmnode {
	node, _ := b.assoc2(shift, hash, key, val)
	return node
}

func (b *bitmapIndexedHmnode) assoc2(shift uint32, hash uint32, key interface{}, val interface{}) (hmnode, bool) {
	bit := bitpos(hash, shift)
	idx := b.index(bit)
	if (b.bitmap & bit) != 0 {
		keyOrNil := b.array[2*idx]
		valOrNode := b.array[2*idx+1]
		if keyOrNil == nil {
			n, ok := valOrNode.(hmnode)
			if !ok {
				panic("Unexpected node type")
			}
			n, addedLeaf := n.assoc2(shift+5, hash, key, val)
			if n == valOrNode {
				return b, false
			}
			return &bitmapIndexedHmnode{b.bitmap, cloneAndSetObjectSlice(b.array, 2*idx+1, n)}, addedLeaf
		}
		if sequtil.Equiv(key, keyOrNil) {
			if val == valOrNode {
				return b, false
			}
			return &bitmapIndexedHmnode{b.bitmap, cloneAndSetObjectSlice(b.array, 2*idx+1, val)}, false
		}
		return &bitmapIndexedHmnode{b.bitmap, cloneAndSetObjectSlice2(b.array, 2*idx, nil, 2*idx+1, createNode(shift+5, keyOrNil, valOrNode, hash, key, val))}, true
	}

	n := sequtil.BitCountU32(b.bitmap)
	if n >= 16 {
		nodes := make([]hmnode, 32)
		jdx := imask(hash, shift)
		nodes[jdx] = emptyBitmapIndexedHmnode.assoc(shift+5, hash, key, val)
		for i, j := 0, 0; i < 32; i++ {
			if ((b.bitmap >> uint(i)) & 1) != 0 {
				if b.array[j] == nil {
					if nn, ok := b.array[j+1].(hmnode); ok {
						nodes[i] = nn
					} else {
						panic("Unexpected node type")
					}
				} else {
					nodes[i] = emptyBitmapIndexedHmnode.assoc(shift+5, Hash(b.array[j]), b.array[j], b.array[j+1])
				}
				j += 2
			}
		}
		return &arrayHmnode{n + 1, nodes}, true
	}

	newArray := make([]interface{}, 2*(n+1))
	copy(newArray, b.array[0:2*idx])
	newArray[2*idx] = key
	newArray[2*idx+1] = val
	copy(newArray[2*(idx+1):], b.array[2*idx:])
	return &bitmapIndexedHmnode{b.bitmap | bit, newArray}, true
}

func (b *bitmapIndexedHmnode) without(shift uint32, hash uint32, key interface{}) hmnode {
	bit := bitpos(hash, shift)
	if (b.bitmap & bit) == 0 {
		return b
	}

	idx := b.index(bit)
	keyOrNil := b.array[2*idx]
	valOrNode := b.array[2*idx+1]
	if keyOrNil == nil {
		n := valOrNode.(hmnode).without(shift+5, hash, key)
		// TOOD: use switch
		if n == valOrNode {
			return b
		}
		if n != nil {
			return &bitmapIndexedHmnode{b.bitmap, cloneAndSetObjectSlice(b.array, 2*idx+1, n)}
		}
		if b.bitmap == bit {
			return nil
		}
		return &bitmapIndexedHmnode{b.bitmap ^ bit, removePair(b.array, idx)}
	}
	if sequtil.Equiv(key, keyOrNil) {
		// TODO: Collapse  (TODO in Java code)
		return &bitmapIndexedHmnode{b.bitmap ^ bit, removePair(b.array, idx)}
	}
	return b
}

func (b *bitmapIndexedHmnode) find(shift uint32, hash uint32, key interface{}) iseq.MapEntry {
	bit := bitpos(hash, shift)
	if (b.bitmap & bit) == 0 {
		return nil
	}

	// TODO: Factor out the following three lines -- repeated
	idx := b.index(bit)
	keyOrNil := b.array[2*idx]
	valOrNode := b.array[2*idx+1]
	if keyOrNil == nil {
		return valOrNode.(hmnode).find(shift+5, hash, key)
	}
	if sequtil.Equiv(key, keyOrNil) {
		return MapEntry{keyOrNil, valOrNode}
	}
	return nil
}

func (b *bitmapIndexedHmnode) findD(shift uint32, hash uint32, key interface{}, notFound interface{}) interface{} {
	bit := bitpos(hash, shift)
	if (b.bitmap & bit) == 0 {
		return notFound
	}

	// TODO: Factor out the following three lines -- repeated
	idx := b.index(bit)
	keyOrNil := b.array[2*idx]
	valOrNode := b.array[2*idx+1]
	if keyOrNil == nil {
		return valOrNode.(hmnode).findD(shift+5, hash, key, notFound)
	}
	if sequtil.Equiv(key, keyOrNil) {
		return valOrNode
	}
	return notFound
}

func (b *bitmapIndexedHmnode) getNodeSeq() iseq.Seq {
	return createHmnodeSeq(b.array)
}

// func (b *bitmapIndexedHmnode) getHash() uint32 {

// }

// hashCollisionHmnode represents a leaf node corresponding to multiple map entries, all with keys that have the same hash value.
type hashCollisionHmnode struct {
	hash  uint32
	count int
	array []interface{}
}

func (h *hashCollisionHmnode) findIndex(key interface{}) int {
	for i := 0; i < 2*h.count; i = i + 2 {
		if sequtil.Equiv(key, h.array[i]) {
			return i
		}
	}
	return -1
}

func (h *hashCollisionHmnode) assoc(shift uint32, hash uint32, key interface{}, val interface{}) hmnode {
	node, _ := h.assoc2(shift, hash, key, val)
	return node
}

func (h *hashCollisionHmnode) assoc2(shift uint32, hash uint32, key interface{}, val interface{}) (hmnode, bool) {
	if h.hash == hash {
		idx := h.findIndex(key)
		if idx != -1 {
			if h.array[idx+1] == val {
				return h, false
			}
			return &hashCollisionHmnode{hash, h.count, cloneAndSetObjectSlice(h.array, idx+1, val)}, false
		}
		newArray := make([]interface{}, len(h.array)+2)
		copy(newArray, h.array)
		newArray[len(h.array)] = key
		newArray[len(h.array)+1] = val
		return &hashCollisionHmnode{hash, h.count + 1, newArray}, true
	}
	// nest it in a bitmap node
	ret, addedLeaf := (&bitmapIndexedHmnode{bitpos(h.hash, shift), []interface{}{nil, h}}).assoc2(shift, hash, key, val)
	return ret, addedLeaf
}

func (h *hashCollisionHmnode) without(shift uint32, hash uint32, key interface{}) hmnode {
	idx := h.findIndex(key)
	// TOOD: use switch
	if idx == -1 {
		return h
	}
	if h.count == 1 {
		return nil
	}
	return &hashCollisionHmnode{hash, h.count - 1, removePair(h.array, idx/2)}

}

func (h *hashCollisionHmnode) find(shift uint32, hash uint32, key interface{}) iseq.MapEntry {
	idx := h.findIndex(key)
	if idx < 0 {
		return nil
	}
	if sequtil.Equiv(key, h.array[idx]) {
		return &MapEntry{h.array[idx], h.array[idx+1]}
	}
	return nil
}

func (h *hashCollisionHmnode) findD(shift uint32, hash uint32, key interface{}, notFound interface{}) interface{} {
	idx := h.findIndex(key)
	if idx < 0 {
		return notFound
	}
	if sequtil.Equiv(key, h.array[idx]) {
		return h.array[idx+1]
	}
	return notFound
}

func (h *hashCollisionHmnode) getNodeSeq() iseq.Seq {
	return createHmnodeSeq(h.array)
}

// func (h *hashCollisionHmnode) getHash() uint32 {

// }

// hmnodeSeq represents an iseq.Seq across an hmnode
type hmnodeSeq struct {
	array []interface{}
	i     int
	s     iseq.Seq
	AMeta
}

func createHmnodeSeq(array []interface{}) *hmnodeSeq {
	return createHmnodeSeq3(array, 0, nil)
}

func createHmnodeSeq3(array []interface{}, i int, s iseq.Seq) *hmnodeSeq {
	if s != nil {
		return &hmnodeSeq{array: array, i: i, s: s}
	}
	for j := i; i < len(array); j = j + 2 {
		if array[j] != nil {
			return &hmnodeSeq{array: array, i: j, s: nil}
		}
		// TODO: nil comparison on interface: fix
		node, ok := array[j+1].(hmnode)
		if !ok {
			panic("Bad node type")
		}
		if node != nil {
			if nodeSeq := node.getNodeSeq(); nodeSeq != nil {
				return &hmnodeSeq{array: array, i: j + 2, s: nodeSeq}
			}
		}
	}
	return nil

}

// hmnodeSeq must implement the following iseq interfaces:
//  Meta, MetaW, Seqable, PCollection, Seq

// interface iseq.MetaW
func (h *hmnodeSeq) WithMeta(meta iseq.PMap) iseq.MetaW {
	return &hmnodeSeq{AMeta: AMeta{meta}, array: h.array, i: h.i, s: h.s}
}

// interface iseq.Seqable

func (h *hmnodeSeq) Seq() iseq.Seq {
	return h
}

// interface iseq.PCollection

func (h *hmnodeSeq) Count() int {
	return sequtil.SeqCount(h)
}

func (h *hmnodeSeq) Cons(o interface{}) iseq.PCollection {
	return NewCons(o, h)
}

func (h *hmnodeSeq) Empty() iseq.PCollection {
	return CachedEmptyList
}

// TODO: Check to make sure not a loop
func (h *hmnodeSeq) Equiv(o interface{}) bool {
	return sequtil.Equiv(h, o)
}

// interface iseq.Seq

func (h *hmnodeSeq) First() interface{} {
	if h.s != nil {
		return h.s.First()
	}
	return MapEntry{h.array[h.i], h.array[h.i+1]}
}

func (h *hmnodeSeq) Next() iseq.Seq {
	if h.s != nil {
		return createHmnodeSeq3(h.array, h.i, h.s.Next())
	}
	return createHmnodeSeq3(h.array, h.i+2, nil)
}

// TODO: pick these up from ASeq
func (h *hmnodeSeq) More() iseq.Seq {
	return moreFromSeq(h)

}

func (h *hmnodeSeq) ConsS(o interface{}) iseq.Seq {
	return NewCons(o, h)
}

func Hash(k interface{}) uint32 {
	return sequtil.Hash(k)
}
