// Copyright (c) 2021 The Houyi Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sst

import "github.com/houyi-tracing/houyi/idl/api_v1"

type treeNode struct {
	op         *api_v1.Operation
	maxN       int // maximum of child nodes
	leafCnt    int
	parent     *treeNode
	childNodes *nodeSet
}

func newRoot(maxN int) *treeNode {
	return &treeNode{
		op:         &api_v1.Operation{},
		maxN:       maxN,
		leafCnt:    0,
		parent:     nil,
		childNodes: newNodeSet(maxN),
	}
}

func newBranchNode(maxN int, parent *treeNode) *treeNode {
	return &treeNode{
		op:         &api_v1.Operation{},
		maxN:       maxN,
		leafCnt:    0,
		parent:     parent,
		childNodes: newNodeSet(maxN),
	}
}

func newLeafNode(maxN int, parent *treeNode, op *api_v1.Operation) *treeNode {
	return &treeNode{
		op:         op,
		maxN:       maxN,
		leafCnt:    1,
		parent:     parent,
		childNodes: nil,
	}
}

func (n *treeNode) addChild(child *treeNode) {
	if n.isLeaf() {
		n.splitSelfAndMerge(child)
	} else {
		if n.hasRoom() {
			n.childNodes.add(child)
			child.parent = n
		} else {
			next := findNext(n.childNodes.all())
			next.addChild(child)
		}
		n.leafCnt += child.leafCnt
	}
}

func (n *treeNode) rmChild(child *treeNode) {
	if !n.isLeaf() {
		n.childNodes.remove(child)
		n.leafCnt -= child.leafCnt
	}
}

func (n *treeNode) hasChild(child *treeNode) bool {
	return !n.isLeaf() && n.childNodes.has(child)
}

func (n *treeNode) hasRoom() bool {
	if !n.isLeaf() {
		return n.childNodes.size() < n.maxN
	} else {
		return false
	}
}

func (n *treeNode) childN() int {
	if !n.isLeaf() {
		return n.childNodes.size()
	} else {
		return 0
	}
}

func (n *treeNode) isLeaf() bool {
	return n.childNodes == nil
}

func (n *treeNode) lruChild(exclude *treeNode) *treeNode {
	if !n.isLeaf() && n.childNodes.size() >= 1 {
		return n.childNodes.lruNode(exclude)
	} else {
		return nil
	}
}

func (n *treeNode) shrink() {
	if !n.isLeaf() && n.childNodes.size() == 1 {
		onlyChild := n.childNodes.all()[0]
		parent := n.parent
		parent.childNodes.remove(n)
		parent.childNodes.add(onlyChild)
		onlyChild.parent = parent
	}
}

func (n *treeNode) splitSelfAndMerge(other *treeNode) {
	grandParent := n.parent
	parent := newBranchNode(n.maxN, grandParent)

	n.parent = parent
	other.parent = parent
	parent.childNodes.add(n)
	parent.childNodes.add(other) // to make OTHER newer, we should add it after adding N.

	grandParent.childNodes.remove(n)
	grandParent.childNodes.add(parent)

	// Move parent to the end of LRU double linked-list to make "n" keep it's LRU property because
	// add function of tree node would add a node as newest node by default.
	grandParent.childNodes.outOfDate(parent)

	parent.leafCnt = n.leafCnt + other.leafCnt
}

func findNext(nodes []*treeNode) *treeNode {
	cnt := len(nodes)
	if cnt == 0 {
		return nil
	}

	nextNode := nodes[cnt-1]
	minLeafCnt := nextNode.leafCnt

	for i := cnt - 1; i >= 0; i-- {
		n := nodes[i]
		if n.isLeaf() {
			return n
		} else {
			if n.leafCnt < minLeafCnt {
				minLeafCnt, nextNode = n.leafCnt, n
			}
		}
	}

	return nextNode
}

// ============== Node Set ============== //

// doubleListNode
type doubleListNode struct {
	val   *treeNode
	left  *doubleListNode
	right *doubleListNode
}

// nodeSet is LRU set for store child nodes
type nodeSet struct {
	capacity int
	head     *doubleListNode
	tail     *doubleListNode
	itemMap  map[*treeNode]*doubleListNode
}

// newNodeSet returns a new nodeSet pointer.
func newNodeSet(capacity int) *nodeSet {
	head, tail := &doubleListNode{}, &doubleListNode{}
	head.right, tail.left = tail, head
	return &nodeSet{
		capacity: capacity,
		head:     head,
		tail:     tail,
		itemMap:  make(map[*treeNode]*doubleListNode),
	}
}

// size returns the size of this cache. The size would not exceed the capacity of this cache.
func (s *nodeSet) size() int {
	return len(s.itemMap)
}

// cap returns the capacity of this cache.
func (s *nodeSet) cap() int {
	return s.capacity
}

// resize sets a new capacity for this cache.
func (s *nodeSet) resize(capacity int) {
	for len(s.itemMap) > capacity {
		delete(s.itemMap, removeTail(s.head, s.tail).val)
	}
	s.capacity = capacity
}

// clear clears this cache.
func (s *nodeSet) clear() {
	s.head.right = s.tail
	s.tail.left = s.head
	s.itemMap = make(map[*treeNode]*doubleListNode)
}

// has returns true if k already exist in this cache, else false.
func (s *nodeSet) has(tn *treeNode) bool {
	nPtr, has := s.itemMap[tn]
	if has {
		moveToHead(s.head, nPtr)
	}
	return has
}

// add adds a new val into this nodeSet.
// If the size of this nodeSet exceed capacity after this operation, the Cache val would be removed and be returned.
func (s *nodeSet) add(tn *treeNode) interface{} {
	if nPtr, has := s.itemMap[tn]; has {
		moveToHead(s.head, nPtr)
	} else {
		newNode := &doubleListNode{
			val:   tn,
			left:  nil,
			right: nil,
		}
		s.itemMap[tn] = newNode
		addToHead(s.head, newNode)
	}

	if len(s.itemMap) > s.capacity {
		tail := removeTail(s.head, s.tail)
		delete(s.itemMap, tail.val)
		return tail.val
	}
	return nil
}

// remove deletes an val from this nodeSet and returns nothing, no matter whether or not it already exist in this nodeSet.
func (s *nodeSet) remove(tn *treeNode) {
	if nPtr, has := s.itemMap[tn]; has {
		removeNode(nPtr)
		delete(s.itemMap, tn)
	}
}

// all return a slice which contains all unique items in this nodeSet.
func (s *nodeSet) all() []*treeNode {
	all := make([]*treeNode, 0, len(s.itemMap))
	p := s.head.right
	for p != s.tail {
		all = append(all, p.val)
		p = p.right
	}
	return all
}

func (s *nodeSet) upToDate(n *treeNode) {
	if node, has := s.itemMap[n]; has {
		moveToHead(s.head, node)
	}
}

func (s *nodeSet) outOfDate(n *treeNode) {
	if node, has := s.itemMap[n]; has {
		moveToTail(s.tail, node)
	}
}

func (s *nodeSet) lruNode(exclude *treeNode) *treeNode {
	if s.size() != 0 {
		if o := s.tail.left; o.val == exclude {
			return o.left.val
		} else {
			return o.val
		}
	} else {
		return nil
	}
}

// moveToHead moves a node to the head of double linked-list.
func moveToHead(head, node *doubleListNode) {
	if node != nil {
		removeNode(node)
		addToHead(head, node)
	}
}

// removeNode removes a node from double linked-list.
func removeNode(node *doubleListNode) {
	node.right.left, node.left.right = node.left, node.right
}

// addToHead adds a node to the head of double linked-list.
func addToHead(head, node *doubleListNode) {
	node.left, node.right = head, head.right
	head.right.left, head.right = node, node
}

// addToTail add a node to the tail of double linked-list.
func addToTail(tail, node *doubleListNode) {
	node.right, node.left = tail, tail.left
	tail.left.right, tail.left = node, node
}

// removeTail removes the tail node of double linked-list.
func removeTail(head, tail *doubleListNode) *doubleListNode {
	if head.right != tail {
		ret := tail.left
		removeNode(tail.left)
		return ret
	}
	return nil
}

// moveToTail moves a node to the tail of double linked-list.
func moveToTail(tail, node *doubleListNode) {
	if node != nil {
		removeNode(node)
		addToTail(tail, node)
	}
}

// ============== Node Map ============== //
type nodeMap map[string]map[string]*treeNode

func newNodeMap() nodeMap {
	return make(nodeMap)
}

func (n nodeMap) add(svc, op string, node *treeNode) {
	if _, has := n[svc]; !has {
		n[svc] = make(map[string]*treeNode)
	}
	n[svc][op] = node
}

func (n nodeMap) has(svc, op string) bool {
	if _, hasSvc := n[svc]; hasSvc {
		if _, hasOp := n[svc][op]; hasOp {
			return true
		}
	}
	return false
}

func (n nodeMap) remove(svc, op string) {
	delete(n[svc], op)
	if len(n[svc]) == 0 {
		delete(n, svc)
	}
}

func (n nodeMap) get(svc, op string) *treeNode {
	if n.has(svc, op) {
		return n[svc][op]
	} else {
		return nil
	}
}

func (n nodeMap) allOperations() []*api_v1.Operation {
	ret := make([]*api_v1.Operation, 0, len(n))
	for svc := range n {
		for _, tn := range n[svc] {
			ret = append(ret, &api_v1.Operation{
				Service:   tn.op.GetService(),
				Operation: tn.op.GetOperation(),
			})
		}
	}
	return ret
}

func (n nodeMap) allNodes() []*treeNode {
	ret := make([]*treeNode, 0, len(n))
	for svc := range n {
		for _, tn := range n[svc] {
			ret = append(ret, tn)
		}
	}
	return ret
}

func (n nodeMap) size() int {
	return len(n)
}
