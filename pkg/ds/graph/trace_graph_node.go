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

package graph

type traceGraphNode struct {
	service   string
	operation string
	inEdges   NodeMap
	outEdges  NodeMap
}

func NewExecutionGraphNode(service, operation string) ExecutionGraphNode {
	return &traceGraphNode{
		service:   service,
		operation: operation,
		inEdges:   NewNodeMap(),
		outEdges:  NewNodeMap(),
	}
}

func (tgn *traceGraphNode) Service() string {
	return tgn.service
}

func (tgn *traceGraphNode) Operation() string {
	return tgn.operation
}

func (tgn *traceGraphNode) GetOuts() []ExecutionGraphNode {
	return tgn.outEdges.AllNodes()
}

func (tgn *traceGraphNode) GetIns() []ExecutionGraphNode {
	return tgn.inEdges.AllNodes()
}

func (tgn *traceGraphNode) HasIn(node ExecutionGraphNode) bool {
	return tgn.inEdges.Has(node)
}

func (tgn *traceGraphNode) HasOut(node ExecutionGraphNode) bool {
	return tgn.outEdges.Has(node)
}

func (tgn *traceGraphNode) AddIn(node ExecutionGraphNode) {
	if !tgn.HasIn(node) {
		tgn.inEdges.Add(node)
	}
}

func (tgn *traceGraphNode) AddOut(node ExecutionGraphNode) {
	if !tgn.HasOut(node) {
		tgn.outEdges.Add(node)
	}
}

func (tgn *traceGraphNode) RemoveIn(node ExecutionGraphNode) {
	if tgn.HasIn(node) {
		tgn.inEdges.Remove(node)
	}
}

func (tgn *traceGraphNode) RemoveOut(node ExecutionGraphNode) {
	if tgn.HasOut(node) {
		tgn.outEdges.Remove(node)
	}
}

func (tgn *traceGraphNode) InN() int {
	return tgn.inEdges.Size()
}

func (tgn *traceGraphNode) OutN() int {
	return tgn.outEdges.Size()
}
