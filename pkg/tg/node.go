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

package tg

import (
	"github.com/houyi-tracing/houyi/idl/api_v1"
)

type node struct {
	operation *api_v1.Operation
	in        nodeMap
	out       nodeMap
}

func newNode(op *api_v1.Operation) *node {
	return &node{
		operation: op,
		in:        newNodeMap(),
		out:       newNodeMap(),
	}
}

func (n *node) AddIn(node *node) {
	n.in.Add(node.operation.Service, node.operation.Operation, node)
}

func (n *node) AddOut(node *node) {
	n.out.Add(node.operation.Service, node.operation.Operation, node)
}

func (n *node) HasIn(node *node) bool {
	return n.in.Has(node.operation.Service, node.operation.Operation)
}

func (n *node) HasOut(node *node) bool {
	return n.out.Has(node.operation.Service, node.operation.Operation)
}

func (n *node) RemoveIn(node *node) {
	n.in.Remove(node.operation.Service, node.operation.Operation)
}

func (n *node) RemoveOut(node *node) {
	n.out.Remove(node.operation.Service, node.operation.Operation)
}

func (n *node) UpdateIn(newNode *node) {
	n.in.Update(newNode.operation.Service, newNode.operation.Operation, newNode)
}

func (n *node) UpdateOut(newNode *node) {
	n.out.Update(newNode.operation.Service, newNode.operation.Operation, newNode)
}

func (n *node) InCnt() int {
	return n.in.Size()
}

func (n *node) OutCnt() int {
	return n.out.Size()
}

func (n *node) AllIn() []*api_v1.Operation {
	inNodes := n.in.All()
	ret := make([]*api_v1.Operation, 0, len(inNodes))
	for _, in := range inNodes {
		ret = append(ret, in.operation)
	}
	return ret
}

func (n *node) AllOut() []*api_v1.Operation {
	outNodes := n.out.All()
	ret := make([]*api_v1.Operation, 0, len(outNodes))
	for _, out := range outNodes {
		ret = append(ret, out.operation)
	}
	return ret
}

func (n *node) String() string {
	return n.operation.String()
}

func addRelation(from, to *node) {
	if from != nil && to != nil {
		from.AddOut(to)
		to.AddIn(from)
	}
}

func removeRelation(from, to *node) {
	if from != nil && to != nil {
		from.RemoveOut(to)
		to.RemoveIn(from)
	}
}

// ===================== Node Map =====================
type nodeMap map[string]map[string]*node

func newNodeMap() nodeMap {
	return make(nodeMap)
}

func (m nodeMap) Add(svc, op string, n *node) {
	if _, hasSvc := m[svc]; !hasSvc {
		m[svc] = make(map[string]*node)
	}
	m[svc][op] = n
}

func (m nodeMap) Has(svc, op string) bool {
	if _, hasSvc := m[svc]; hasSvc {
		if _, hasOp := m[svc][op]; hasOp {
			return true
		}
	}
	return false
}

func (m nodeMap) Remove(svc, op string) {
	if m.Has(svc, op) {
		delete(m[svc], op)
		if len(m[svc]) == 0 {
			delete(m, svc)
		}
	}
}

func (m nodeMap) Get(svc, op string) *node {
	if m.Has(svc, op) {
		return m[svc][op]
	} else {
		return nil
	}
}

func (m nodeMap) Update(svc, op string, node *node) {
	if m.Has(svc, op) {
		m[svc][op] = node
	}
}

func (m nodeMap) All() []*node {
	ret := make([]*node, 0)
	for _, ops := range m {
		for _, n := range ops {
			ret = append(ret, n)
		}
	}
	return ret
}

func (m nodeMap) Services() []string {
	ret := make([]string, 0, len(m))
	for svc := range m {
		ret = append(ret, svc)
	}
	return ret
}

func (m nodeMap) Operations(svc string) []string {
	if svcMap, ok := m[svc]; ok {
		ret := make([]string, 0, len(svcMap))
		for op := range svcMap {
			ret = append(ret, op)
		}
		return ret
	} else {
		return []string{}
	}
}

func (m nodeMap) Size() int {
	return len(m)
}
