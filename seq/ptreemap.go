// Copyright 2012 David Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package seq

import (
	"fmt"
	"github.com/dmiller/go-seq/iseq"
	"github.com/dmiller/go-seq/sequtil"
	"hash"
)

// PTreeMap implements a persistent Red-Black tree.
// Instances of this struct are constant values.
// See Okasaki, Kahrs, Larsen, et al

type PTreeMap struct {
	comp iseq.CompareFn
	tree tmNode
	count int
	AMeta
	hash uint32
}

var (
	EmptyPTreeMap = &PTreeMap{comp: sequtil.DefaultCompareFn}
	)


// factories

// TODO: need a factory for creating from an arbitrary Go map

func CreateEmptyPTreeMap(comp iseq.CompareFn) *PTreeMap {
	return &PTreeMap{comp: comp}
}

func NewPTreeMapFromSeq(items iseq.Seq) *PTreeMap {
	return NewPTreeMapFromSeqC(items,sequtil.DefaultCompareFn)
}

func NewPTreeMapFromSeqC(items iseq.Seq, comp iseq.CompareFn) *PTreeMap {

	ret := CreateEmptyPTreeMap(comp)

	for i := 0; items != nil; items,i  = items.Next().Next(), i+1 {
		if items.Next() == nil {
			panic(fmt.Sprintf("No value supplied for key: %v",items.First()))
		}
		ret = ret.AssocM(items.First(),items.Next().First()).(*PTreeMap)
	}
	return ret
}

func NewPTreeMapFromSlice(s []interface{}) *PTreeMap {
	return NewPTreeMapFromSliceC(s,sequtil.DefaultCompareFn)
}

func NewPTreeMapFromSliceC(s []interface{}, comp iseq.CompareFn) *PTreeMap {
	ret := CreateEmptyPTreeMap(comp)
	for i:=0; i<len(s); i = i+2 {
		ret = ret.AssocM(s[i],s[i+1]).(*PTreeMap)
	}
	return ret
}


func NewPTreeMapFromItems(items ...interface{}) *PTreeMap {
	return NewPTreeMapFromSliceC(items,sequtil.DefaultCompareFn)
}

func NewPTreeMapFromItemsC(comp iseq.CompareFn, items ...interface{}) *PTreeMap {
	return NewPTreeMapFromSliceC(items,comp)
}

// PTreeMap needs to implement the following iseq interfaces:
//		Obj Meta Seqable PCollection Lookup Associative Counted PMap Reversible Sorted?
//  Are we going to do EditableCollection?
//  Also, Equatable and Hashable
//
// interface Meta is covered by the AMeta embedding
// TODO: IEditableCollection

// interface iseq.Obj

func (m *PTreeMap) WithMeta(meta iseq.PMap) iseq.Obj {
	return &PTreeMap{comp: m.comp, tree: m.tree, count: m.count, AMeta:AMeta{m.meta}}
}


// interface iseq.Associative, iseq.Lookup


func (m *PTreeMap) ContainsKey(key interface{}) bool {
	return m.tmNodeAt(key) != nil
}


func (m *PTreeMap) EntryAt(key interface{}) iseq.MapEntry {
	return m.tmNodeAt(key)
}

func (m *PTreeMap) ValAt(key interface{}) interface{} {
	return m.ValAtD(key,nil)
}

func (m *PTreeMap) ValAtD(key interface{}, notFound interface{}) interface{} {
	if n := m.tmNodeAt(key); n != nil {
		return n.val()
	}
	return notFound 
}

func (m *PTreeMap) Assoc(key interface{}, val interface{}) iseq.Associative {
	return m.AssocM(key,val)
}


// interface iseq.PMap

func (m *PTreeMap) AssocM(key interface{}, val interface{}) iseq.PMap {
	tree ,foundNode := m.addNode(m.tree,key,val)
	if tree == nil {
		if foundNode.val() == val {
			return m
		}
		return &PTreeMap{comp: m.comp, tree: m.replace(m.tree, key, val), count: m.count, AMeta: AMeta{m.meta}}
	}
	return &PTreeMap{comp: m.comp, tree: tree.blacken(), count: m.count+1, AMeta: AMeta{m.meta}}
}


func (m *PTreeMap) Without(key interface{}) iseq.PMap {
	tree, foundNode := m.remove(m.tree,key)
	if tree == nil {
		if foundNode == nil {
			return m
		}
		return &PTreeMap{comp: m.comp, AMeta: AMeta{m.meta}}
	}
	return &PTreeMap{comp: m.comp, tree: tree.blacken(), count: m.count-1, AMeta: AMeta{m.meta}}
}

func (m *PTreeMap) ConsM(e iseq.MapEntry) iseq.PMap {
	return sequtil.MapCons(m,e)
}

// interface iseq.PCollection, iseq.Seqable, iseq.Counted

func (m *PTreeMap) Count() int {
	return m.count
}

func (m *PTreeMap) Cons(o interface{}) iseq.PCollection {
	return sequtil.MapCons(m,o)
}

func (m *PTreeMap) Empty() iseq.PCollection {
	return &PTreeMap{comp: m.comp, AMeta: AMeta{m.meta}}
}

func (m *PTreeMap) Equiv(o interface{}) bool {
	return sequtil.Equiv(m,o)
}

func (m *PTreeMap) Count1() int {
	return m.count
}

func (m *PTreeMap) Seq() iseq.Seq {
	if m.count > 0 {
		return createTmnodeSeq(m.tree,true,m.count)
	}
	return nil
}

// interface Reversible

func (m *PTreeMap) Rseq() iseq.Seq {
	if m.count > 0 {
		return createTmnodeSeq(m.tree,false,m.count)
	}
	return nil
}


// interface Sorted

func (m *PTreeMap) Comparator() iseq.CompareFn {
	return m.comp
}

func (m *PTreeMap) EntryKey(entry interface{}) interface{} {
	if me, ok := entry.(iseq.MapEntry); ok {
		return me.Key()
	}
	panic("Expected an iseq.MapEntry")
}

func (m *PTreeMap) SeqA(ascending bool) iseq.Seq {
	if m.count > 0 {
		return createTmnodeSeq(m.tree,ascending,m.count)
	}
	return nil
}

func (m *PTreeMap) SeqFrom(key interface{}, ascending bool) iseq.Seq {
	if m.count > 0 {
		var stack iseq.Seq
		t := m.tree
		for t != nil {
			c := m.doCompare(key,t.key()) 
			if c == 0 {
				stack = smartCons(t,stack)
				return createTmnodeSeqFromStack(stack,ascending)
			} else if ascending {
				if c < 0 {
					stack = smartCons(t,stack)
					t = t.left()
				} else {
					t = t.right()
				}
			} else {
				if c > 0 {
					stack = smartCons(t,stack)
					t = t.right()
				} else {
					t = t.left()
				}
			}
		}
		if stack != nil {
			return createTmnodeSeqFromStack(stack,ascending)
		}
	}
	return nil
}

// interfaces Equatable, Hashable

func (p *PTreeMap) Equals(o interface{}) bool {
	return sequtil.MapEquals(p,o)
}

func (p *PTreeMap) Hash() uint32 {
	if p.hash == 0 {
		p.hash = sequtil.HashMap(p)
	}
	return p.hash
}

func (p *PTreeMap) AddHash(h hash.Hash) {
	sequtil.AddHashMap(h,p)
}

// tree operations

func (m *PTreeMap) tmNodeAt(key interface{}) tmNode {
	t := m.tree
	for t != nil {
		c :=m.doCompare(key,t.key())
		switch  {
			case c == 0: return t
			case c < 0: t = t.left()
			default: t = t.right()
		}
	}
	return t
}

func (m *PTreeMap) doCompare(k1, k2 interface{}) int {
	return m.comp(k1,k2)
}

func (m *PTreeMap) addNode(t tmNode, key, val interface{}) (newRoot tmNode, addNode tmNode) {
	if t == nil {
		return makeRed(key,val,nil,nil),nil
	}	
	c := m.doCompare(key,t.key())
	if c == 0 {
		return nil,t
	}
	var ins, addedNode tmNode
	if c < 0 {
		ins, addedNode = m.addNode(t.left(),key,val)
	} else {
		ins, addedNode = m.addNode(t.right(),key,val)
	}
	if ins == nil {
		return nil,addedNode
	}
	if c < 0 {
		return t.addLeft(ins),addedNode
	}
	return t.addRight(ins),addedNode
}

func (m *PTreeMap) remove(t tmNode, key interface{}) (newRoot tmNode, remdNode tmNode){
	if t == nil {
		return nil,nil
	}
	c := m.doCompare(key,t.key())
	if c == 0 {
		return appendTmnode(t.left(),t.right()),t
	}
	var del,remNode tmNode
	if c < 0 {
		del,remNode = m.remove(t.left(),key)
	} else {
		del,remNode = m.remove(t.right(),key)
	}
	if  del == nil && remNode == nil {
		return nil, nil
	}
	if c < 0 {
		if isBlack(t.left()) {
			return balanceLeftDel(t.key(),t.val(),del,t.right()),remNode
		} else {
			return makeRed(t.key(),t.val(),del,t.right()),remNode
		}
	}
	if isBlack(t.right()) {
		return balanceRightDel(t.key(),t.val(),t.left(),del),remNode
	}
	return makeRed(t.key(),t.val(),t.left(),del),remNode
}

func appendTmnode(left, right tmNode) tmNode {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if isRed(left) {
		if isRed(right) {
			app := appendTmnode(left.right(), right.left())
			if isRed(app) {
				makeRed(app.key(), app.val(),
					makeRed(left.key(),left.val(),left.left(),app.left()),
					makeRed(right.key(),right.val(),app.right(),right.right()))
			}
		} else {
			return makeRed(left.key(),left.val(),left.left(),appendTmnode(left.right(),right))
		}
	}

	if isRed(right) {
		return makeRed(right.key(),right.val(),appendTmnode(left,right.left()),right.right())
	}

	app := appendTmnode(left.right(),right.left())
	if isRed(app) {
		return makeRed(app.key(),app.val(),
				makeBlack(left.key(),left.val(),left.left(),app.left()),
				makeBlack(right.key(),right.val(),app.right(),right.right()))
	}
	return balanceLeftDel(left.key(),left.val(),left.left(),makeBlack(right.key(),right.val(),app,right.right()))
}

func balanceLeftDel(key, val interface{}, del, right tmNode) tmNode {
	if isRed(del) {
		return makeRed(key,val,del.blacken(),right)
	}
	if isBlack(right) {
		return rightBalance(key,val,del,right.redden())
	}
	if isRed(right) && isBlack(right.left()) {
		return makeRed(right.left().key(), right.left().val(),
				makeBlack(key,val,del,right.left().left()),
				rightBalance(right.key(),right.val(),right.left().right(),right.right().redden()))
	}
	panic("Invariant violation!")
}

func balanceRightDel(key, val interface{}, left, del tmNode) tmNode {
	if isRed(del) {
		return makeRed(key,val,left,del.blacken())
	}
	if isBlack(left) {
		return leftBalance(key,val,left.redden(),del)
	}
	if isRed(left) && isBlack(left.right()) {
		return makeRed(left.right().key(), left.right().val(),
				leftBalance(left.key(),left.val(),left.left().redden(),left.right().left()),
				makeBlack(key,val,left.right().right(),del))
	}
	panic("Invariant violation!")	
}

func leftBalance(key, val interface{}, ins, right tmNode) tmNode {
	insIsRed := isRed(ins)
	if insIsRed && isRed(ins.left()) {
		return makeRed(ins.key(),ins.val(),ins.left().blacken(), makeBlack(key,val,ins.right(),right))
	}
	if insIsRed && isRed(ins.right()) {
		return makeRed(ins.right().key(),ins.right().val(),
				makeBlack(ins.key(),ins.val(),ins.left(),ins.right().left()),
				makeBlack(key,val,ins.right().right(),right))
	}
	return makeBlack(key,val,ins,right)
}

func rightBalance(key, val interface{}, left, ins tmNode) tmNode {
	insIsRed := isRed(ins)
	if insIsRed && isRed(ins.right()) {
		return makeRed(ins.key(),ins.val(),makeBlack(key,val,left,ins.left()),ins.right().blacken())
	}
	if insIsRed && isRed(ins.left()) {
		return makeRed(ins.left().key(),ins.left().val(),
				makeBlack(key,val,left,ins.left().left()),
				makeBlack(ins.key(),ins.val(),ins.left().right(),ins.right()))
	}
	return makeBlack(key,val,left,ins)
}

func (m *PTreeMap) replace(t tmNode, key, val interface{}) tmNode {
	c := m.doCompare(key,t.key())
	v := t.val()
	l := t.left()
	r := t.right()

	if c == 0 {
		v = val
	} 
	if c < 0 {
		l = m.replace(t.left(),key,val)
	} 
	if c > 0 {
		r = m.replace(t.right(),key,val)
	}

	return t.replace(t.key(),v,l,r)
}

func makeRed(key, val interface{}, left, right tmNode) tmNode {
	return &redTmnode{baseTmnode:baseTmnode{key,val,left,right}}
}

func makeBlack(key, val interface{}, left, right tmNode) tmNode {
	return &blackTmnode{baseTmnode:baseTmnode{key,val,left,right}}
}


type tmNode interface {
	iseq.MapEntry
	addLeft(ins tmNode) tmNode
	addRight(ins tmNode) tmNode
	removeLeft(ins tmNode) tmNode
	removeRight(ins tmNode) tmNode
	blacken() tmNode
	redden() tmNode
	replace(key, val interface{}, left, right tmNode) tmNode
	balanceLeft(parent tmNode) tmNode
	balanceRight(parent tmNode) tmNode
	key() interface{}
	val() interface{}
	left() tmNode
	right() tmNode
}

type baseTmnode struct {
	_key interface{}
	_val interface{}
	_left tmNode
	_right tmNode
}

func isBlack(x tmNode) bool {
	_,ok := x.(*blackTmnode)
	return ok
}

func isRed(x tmNode) bool {
	_,ok := x.(*redTmnode)
	return ok
}

// interface iseq.MapEntry

func (n *baseTmnode) Key() interface{} {
	return n._key
}

func (n *baseTmnode) Val()  interface{} {
	return n._val
}

// basic tmNode methods

func (n *baseTmnode) key()  interface{} {
	return n._key
}

func (n *baseTmnode) val()  interface{} {
	return n._val
}

func (n *baseTmnode) left() tmNode {
	return n._left
}

func (n *baseTmnode) right() tmNode {
	return n._right
}


type blackTmnode struct {
	baseTmnode
}

func (n *blackTmnode) balanceLeft(parent tmNode) tmNode { 
	return makeBlack(parent.key(),parent.val(),n,parent.right())
}
func (n *blackTmnode) balanceRight(parent tmNode) tmNode {
	return makeBlack(parent.key(), parent.val(), parent.left(),n)
}

func (n *blackTmnode) addLeft(ins tmNode) tmNode {
	return ins.balanceLeft(n)
}

func (n *blackTmnode) addRight(ins tmNode) tmNode {
	return ins.balanceRight(n)
}

func (n *blackTmnode) removeLeft(del tmNode) tmNode { 
	return balanceLeftDel(n._key,n._val,del,n._right)
}

func (n *blackTmnode) removeRight(del tmNode) tmNode {
	return balanceRightDel(n._key,n._val,n._left,del)
}

func (n *blackTmnode) blacken() tmNode { 
	return n
}

func (n *blackTmnode) redden() tmNode { 
	return makeRed(n._key,n._val,n._left,n._right)
}

func (n *blackTmnode) replace(key, val interface{}, left, right tmNode) tmNode { 
	return makeBlack(key,val,left,right)
}

type redTmnode struct {
	baseTmnode
}

func (n *redTmnode) addLeft(ins tmNode) tmNode { 
	return makeRed(n._key,n._val,ins,n._right)
}

func (n *redTmnode) addRight(ins tmNode) tmNode { 
	return makeRed(n._key,n._val,n._left,ins)
}

func (n *redTmnode) removeLeft(del tmNode) tmNode {
	return makeRed(n._key,n._val,del,n._right)
}

func (n *redTmnode) removeRight(del tmNode) tmNode {
	return makeRed(n._key,n._val,n._left,del)
}

func (n *redTmnode) blacken() tmNode {
	return makeBlack(n._key,n._val,n._left,n._right)
}

func (n *redTmnode) replace(key, val interface{}, left, right tmNode) tmNode {
	return makeRed(key,val,left,right)
}

func (n *redTmnode) balanceLeft(parent tmNode) tmNode { 
	if n._left != nil && isRed(n._left) {
		return makeRed(n._key,n._val,n._left.blacken(),makeBlack(parent.key(),parent.val(),n._right,parent.right()))
	}
	if n._right != nil && isRed(n._right) {
		return makeRed(n._right.key(),n._right.val(),
				makeBlack(n._key,n._val,n._left,n._right.left()),
				makeBlack(parent.key(),parent.val(),n._right.right(),parent.right()))
	}
	return makeBlack(parent.key(),parent.val(),n,parent.right())
}

func (n *redTmnode) balanceRight(parent tmNode) tmNode {
	if n._right!= nil && isRed(n._right) {
		return makeRed(n._key,n._val,makeBlack(parent.key(),parent.val(),parent.left(),n._left), n._right.blacken())
	}
	if n._left != nil && isRed(n._left) {
		return makeRed(n._left.key(), n._left.val(),
				makeBlack(parent.key(),parent.val(),parent.left(),n._left.left()),
				makeBlack(n._key,n._val,n._left.right(),n._right))
	}
	return makeBlack(parent.key(), parent.val(), parent.left(),n)
}

func (n *redTmnode) redden() tmNode {
	panic("Invariant violation!")
}

func createTmnodeSeq(t tmNode, ascending bool, count int) iseq.Seq {
	return &tmNodeSeq{pushTmnodeSeq(t,nil,ascending),ascending,count,AMeta{nil}}
}

func createTmnodeSeqFromStack(s iseq.Seq, asc bool) iseq.Seq {
	return &tmNodeSeq{s,asc,-1,AMeta{nil}}
}

func pushTmnodeSeq(t tmNode, stack iseq.Seq, asc bool) iseq.Seq {
	for t != nil {
		stack = smartCons(t,stack)
		if asc {
			t = t.left()
		} else {
			t = t.right()
		}
	}
	return stack
}


type tmNodeSeq struct {
	stack iseq.Seq
	asc bool
	cnt int
	AMeta
}

// interface iseq.Obj

func (t *tmNodeSeq) WithMeta(meta iseq.PMap) iseq.Obj {
	return &tmNodeSeq{AMeta:AMeta{meta},stack:t.stack,asc:t.asc,cnt:t.cnt}
}

// interface iseq.Seq

func (t *tmNodeSeq) First() interface{} {
	return t.stack.First()
}

func (t *tmNodeSeq) Next() iseq.Seq {
	node, ok := t.stack.First().(tmNode)
	if !ok { panic("Unexpected node type")}
	var n2 tmNode
	if t.asc {
		n2 = node.right()
	} else {
		n2 = node.left()
	}
	nextStack := pushTmnodeSeq(n2,t.stack.Next(),t.asc)
	if nextStack != nil {
		return &tmNodeSeq{nextStack,t.asc,t.cnt-1,AMeta{nil}}
	}
	return nil
}

func (t *tmNodeSeq) More() iseq.Seq {
	return moreFromSeq(t)
}

func (t *tmNodeSeq) SCons(o interface{}) iseq.Seq {
	return NewCons(o,t)
}


// interface iseq.Seqable

func (t *tmNodeSeq) Seq() iseq.Seq {
	return t
}

// interface iseq.PCollection

func (t *tmNodeSeq) Count() int {
	if t.cnt < 0 {
		return sequtil.SeqCount(t)
	}
	return t.cnt
}

func (t *tmNodeSeq) Cons(o interface{}) iseq.PCollection {
	return NewCons(o,t)
}

func (t *tmNodeSeq) Empty() iseq.PCollection {
	return CachedEmptyList;
}

func (t *tmNodeSeq) Equiv(o interface{}) bool {
	// TODO: revisit Equiv
	return sequtil.Equals(t,o)
}

// TODO: Test that keys are ordered when seq'd

// TODO: Finish tests
