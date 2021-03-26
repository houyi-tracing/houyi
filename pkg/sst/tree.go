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

import (
	"fmt"
	"github.com/houyi-tracing/houyi/idl/api_v1"
)

const (
	notExistErr     = "operation doest not exist"
	alreadyExistErr = "operation already exist"
)

type sst struct {
	root  *treeNode
	nodes nodeMap
	maxN  int
}

func NewSamplingStrategyTree(maxN int) SamplingStrategyTree {
	return &sst{
		root:  newRoot(maxN),
		nodes: newNodeMap(),
		maxN:  maxN,
	}
}

func (t *sst) Add(op *api_v1.Operation) error {
	if !t.hasOp(op) {
		svcName, opName := op.GetService(), op.GetOperation()
		newNode := newLeafNode(t.maxN, nil, op)
		t.nodes.add(svcName, opName, newNode)
		t.root.addChild(newNode)
		return nil
	} else {
		return fmt.Errorf(alreadyExistErr)
	}
}

func (t *sst) Has(op *api_v1.Operation) bool {
	return t.hasOp(op)
}

func (t *sst) Promote(op *api_v1.Operation) error {
	if t.hasOp(op) {
		node := t.getNode(op)
		if node.parent == t.root {
			t.root.childNodes.upToDate(node)
			return nil
		}
		grandParent, parent := node.parent.parent, node.parent
		t.promote(grandParent, parent, node)
		return nil
	} else {
		return fmt.Errorf(notExistErr)
	}
}

func (t *sst) Generate(op *api_v1.Operation) (float64, error) {
	if t.hasOp(op) {
		node := t.getNode(op)
		sr, parent := 1.0, node.parent
		for parent != nil {
			sr *= 1.0 / float64(parent.childNodes.size())
			parent = parent.parent
		}
		return sr, nil
	} else {
		return 0.0, fmt.Errorf(notExistErr)
	}
}

func (t *sst) Prune(op *api_v1.Operation) error {
	if t.hasOp(op) {
		node := t.getNode(op)
		currP := node.parent
		currP.childNodes.remove(node)
		p := currP

		for p != nil {
			p.leafCnt -= node.leafCnt
			p = p.parent
		}

		if currP != t.root {
			currP.shrink()
		}
		t.nodes.remove(op.GetService(), op.GetOperation())
		return nil
	} else {
		return fmt.Errorf(notExistErr)
	}
}

func (t *sst) hasOp(op *api_v1.Operation) bool {
	return t.nodes.has(op.GetService(), op.GetOperation())
}

func (t *sst) getNode(op *api_v1.Operation) *treeNode {
	return t.nodes.get(op.GetService(), op.GetOperation())
}

func (t *sst) promote(gp, p, node *treeNode) {
	p.childNodes.remove(node)

	if gp.hasRoom() {
		gp.childNodes.add(node)

		node.parent = gp
		p.leafCnt -= node.leafCnt

		p.shrink()
	} else {
		lruNode := gp.lruChild(p)
		if p.childN() > 2 {
			lruNode.splitSelfAndMerge(node)
			p.leafCnt -= node.leafCnt
		} else {
			gp.childNodes.remove(lruNode)
			lruNode.parent = p

			p.childNodes.add(lruNode)
			gp.childNodes.add(node)
			node.parent = gp

			p.leafCnt = p.leafCnt - node.leafCnt + lruNode.leafCnt
		}
	}
}
