// Copyright (c) 2020 The Houyi Authors.
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

import (
	"github.com/yan-fuhai/go-ds/cache"
	"math"
)

type treeNode struct {
	// service name
	service string

	// operation name
	operation string

	// maximum number of children
	maxChildrenN int

	// number of leaf nodes in this tree, if this node is a leaf node,
	// the value of this filed is 0.
	leafCnt int

	// children nodes
	childNodes cache.Set

	// parent treeNode of current node
	parent TreeNode
}

// newTreeNode returns a new tree node pointer.
func newTreeNode(maxChildrenN int, parent TreeNode) TreeNode {
	if maxChildrenN < 1 {
		panic("number of maximum children must be greater than or equal to 1")
	}

	return &treeNode{
		maxChildrenN: maxChildrenN,
		leafCnt:      0,
		parent:       parent,
		childNodes:   cache.NewSet(maxChildrenN),
	}
}

// newLeafNode returns a new leaf node pointer.
func newLeafNode(service string, operation string, maxChildrenN int, parent TreeNode) TreeNode {
	return &treeNode{
		service:      service,
		operation:    operation,
		maxChildrenN: maxChildrenN,
		leafCnt:      0,
		childNodes:   nil,
		parent:       parent,
	}
}

func (n *treeNode) AddChild(child TreeNode) TreeNode {
	if childNode, ok := child.(*treeNode); ok {
		childNode.parent = n
		lruNode := n.childNodes.Add(child)

		if childNode.IsLeaf() {
			n.leafCnt++
		} else {
			n.leafCnt += childNode.leafCnt
		}

		if lruNode != nil {
			return lruNode.(TreeNode)
		}
	}
	return nil
}

func (n *treeNode) RemoveChild(child TreeNode) {
	if childNode, ok := child.(*treeNode); ok && n.childNodes.Has(childNode) {
		n.childNodes.Delete(childNode)
		childNode.parent = nil
		if childNode.IsLeaf() {
			n.leafCnt--
		}
	}
}

func (n *treeNode) Children() []TreeNode {
	ret := make([]TreeNode, 0, n.childNodes.Size())
	for _, c := range n.childNodes.Keys() {
		ret = append(ret, c.(*treeNode))
	}
	return ret
}

func (n *treeNode) ChildrenN() int {
	return n.childNodes.Size()
}

func (n *treeNode) HasChild(child TreeNode) bool {
	return n.childNodes.Has(child)
}

func (n *treeNode) Parent() TreeNode {
	return n.parent
}

func (n *treeNode) PathCompression() {
	if !n.IsLeaf() && n.childNodes.Size() == 1 {
		onlyChild := n.childNodes.Keys()[0].(*treeNode)
		parent := n.parent.(*treeNode)

		parent.childNodes.Delete(n)
		parent.childNodes.Add(onlyChild)
		onlyChild.parent = parent
	}
}

func (n *treeNode) IsLeaf() bool {
	return n.leafCnt == 0 && n.childNodes == nil
}

func (n *treeNode) HasRoom() bool {
	return n.childNodes.Size() < n.maxChildrenN
}

func (n *treeNode) InsertAsDesc(leaf TreeNode) {
	n.insertLeaf(leaf)
}

func (n *treeNode) Service() string {
	return n.service
}

func (n *treeNode) Operation() string {
	return n.operation
}

func (n *treeNode) insertLeaf(leaf TreeNode) {
	if n.IsLeaf() {
		if newParent, ok := n.split().(*treeNode); ok && newParent != nil {
			newParent.childNodes.Add(leaf)
			leaf.(*treeNode).parent = newParent
			newParent.leafCnt = 2
		}
	} else {
		n.leafCnt++ // increase leaf node count of current node
		if n.HasRoom() {
			// insert the leaf as a direct child node of current node.
			n.childNodes.Add(leaf)
			if node, ok := leaf.(*treeNode); ok {
				node.parent = n
			}
		} else {
			// insert the leaf as a descendant of current node.
			if next, ok := n.leastLeafChild().(*treeNode); ok {
				next.insertLeaf(leaf)
			}
		}
	}
}

// split generate a new branch node to be the parent of current node
// and the current node would be the child node of itself.
// Current node must be leaf.
func (n *treeNode) split() TreeNode {
	if n.IsLeaf() {
		oldParent := n.parent.(*treeNode)
		newParent := newTreeNode(n.maxChildrenN, oldParent).(*treeNode)

		oldParent.childNodes.Add(newParent)
		oldParent.childNodes.Delete(n)
		newParent.childNodes.Add(n)

		newParent.parent = oldParent
		n.parent = newParent
		return newParent
	}
	return nil
}

// leastLeafChild returns the child node which has least leaf nodes among all child nodes.
// The LEAF node has the HIGHEST priority to be returned.
func (n *treeNode) leastLeafChild() TreeNode {
	if n.childNodes.Size() != 0 {
		ret := (*treeNode)(nil)
		min := math.MaxInt64
		children := n.childNodes.Keys()
		for i := len(children) - 1; i >= 0; i-- {
			// Literate from tail to head because our purpose is found out the LRU leaf node
			// when current node has direct leaf nodes.
			child := children[i].(*treeNode)
			if child.IsLeaf() {
				return child
			}
			if child.leafCnt < min {
				min, ret = child.leafCnt, child
			}
		}
		return ret
	}
	return nil
}
