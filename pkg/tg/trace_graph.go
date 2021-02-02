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
	"fmt"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/yan-fuhai/go-ds/set"
	"go.uber.org/zap"
	"strings"
	"sync"
)

const (
	fakeRootService   = "fake-root-service"
	fakeRootOperation = "fake-root-operation"
)

const (
	OperationDoesNotExistErr = "operation does not exist in this trace graph"
	OperationAlreadyExistErr = "operation already exist in this trace graph"
)

type traceGraph struct {
	sync.RWMutex

	logger *zap.Logger
	nodes  nodeMap

	// globalRoot is used to mark some trace graph node as entries.
	// Entries have below features:
	//  1. The operation of entries would not be called by other operations.
	//  2. Node of every entry has a relation from globalRoot to itself.
	// ATTENTION: globalRoot is not in traceGraph.nodes.
	globalRoot *node
}

func NewTraceGraph(logger *zap.Logger) TraceGraph {
	fakeRootOp := &api_v1.Operation{
		Service:   fakeRootService,
		Operation: fakeRootOperation,
	}
	return &traceGraph{
		RWMutex:    sync.RWMutex{},
		logger:     logger,
		nodes:      newNodeMap(),
		globalRoot: newNode(fakeRootOp),
	}
}

func (t *traceGraph) Add(op *api_v1.Operation) error {
	t.Lock()
	defer t.Unlock()

	if !t.has(op) {
		n := newNode(op)
		t.nodes.Add(op.GetService(), op.GetOperation(), n)

		// mark a new operation as an ingress operation because there are not other operations calling it.
		addRelation(t.globalRoot, n)
		t.logger.Debug("added operation",
			zap.String("service", op.Service), zap.String("operation", op.Operation))
		return nil
	} else {
		return fmt.Errorf(OperationAlreadyExistErr)
	}
}

func (t *traceGraph) Remove(op *api_v1.Operation) error {
	t.Lock()
	defer t.Unlock()

	if t.has(op) {
		rmNode := t.get(op)

		// remove all relations related to this node for garbage collection
		for _, in := range rmNode.in.All() {
			in.RemoveOut(rmNode)
		}
		for _, out := range rmNode.out.All() {
			out.RemoveIn(rmNode)
		}

		t.nodes.Remove(op.Service, op.Operation)
		t.logger.Debug("removed operation",
			zap.String("service", op.Service), zap.String("operation", op.Operation))
		return nil
	} else {
		return fmt.Errorf(OperationDoesNotExistErr)
	}
}

func (t *traceGraph) Has(op *api_v1.Operation) bool {
	t.RLock()
	defer t.RUnlock()

	return t.has(op)
}

func (t *traceGraph) AddRelation(rel *api_v1.Relation) error {
	t.Lock()
	defer t.Unlock()

	from, to := rel.GetFrom(), rel.GetTo()
	if t.has(from) && t.has(to) {
		fromNode, toNode := t.get(from), t.get(to)
		addRelation(fromNode, toNode)

		if toNode.HasIn(t.globalRoot) {
			removeRelation(t.globalRoot, toNode)
		}

		t.logger.Debug("added relation", zap.String("relation", rel.String()))
		return nil
	} else {
		return fmt.Errorf(OperationDoesNotExistErr)
	}
}

func (t *traceGraph) RemoveRelation(rel *api_v1.Relation) error {
	t.Lock()
	defer t.Unlock()

	from, to := rel.GetFrom(), rel.GetTo()
	if t.has(from) && t.has(to) {
		// from -> to
		fromNode, toNode := t.get(from), t.get(to)
		removeRelation(fromNode, toNode)

		if toNode.InCnt() == 0 && toNode.OutCnt() != 0 {
			addRelation(t.globalRoot, toNode)
		}

		t.logger.Debug("removed relation", zap.String("relation", rel.String()))
		return nil
	} else {
		return fmt.Errorf(OperationDoesNotExistErr)
	}
}

func (t *traceGraph) HasRelation(rel *api_v1.Relation) bool {
	t.RLock()
	defer t.RUnlock()

	from, to := rel.GetFrom(), rel.GetTo()
	if t.has(from) && t.has(to) {
		fromNode, toNode := t.get(from), t.get(to)
		return fromNode.HasOut(toNode) && toNode.HasIn(fromNode)
	} else {
		return false
	}
}

func (t *traceGraph) IsIngress(op *api_v1.Operation) bool {
	t.RLock()
	defer t.RUnlock()

	if t.has(op) {
		n := t.get(op)
		return n.HasIn(t.globalRoot) && t.globalRoot.HasOut(n)
	} else {
		return false
	}
}

func (t *traceGraph) GetIngresses(op *api_v1.Operation) ([]*api_v1.Operation, error) {
	t.RLock()
	defer t.RUnlock()

	if t.has(op) {
		entries := make([]*api_v1.Operation, 0)
		entries = t.searchIngresses(t.get(op), entries, set.NewSet())
		return entries, nil
	} else {
		return nil, fmt.Errorf(OperationDoesNotExistErr)
	}
}

func (t *traceGraph) Traces(op *api_v1.Operation) ([]*api_v1.TraceNode, error) {
	t.RLock()
	defer t.RUnlock()

	if t.has(op) {
		traces := make([]*api_v1.TraceNode, 0)

		ingresses := make([]*api_v1.Operation, 0)
		ingresses = t.searchIngresses(t.get(op), ingresses, set.NewSet())
		for _, e := range ingresses {
			traces = append(traces, generateTrace(t.get(e)))
		}

		return traces, nil
	} else {
		return nil, fmt.Errorf(OperationDoesNotExistErr)
	}
}

func (t *traceGraph) Size() int {
	t.RLock()
	defer t.RUnlock()

	return t.nodes.Size()
}

func (t *traceGraph) has(op *api_v1.Operation) bool {
	return t.nodes.Has(op.GetService(), op.GetOperation())
}

func (t *traceGraph) get(op *api_v1.Operation) *node {
	return t.nodes.Get(op.GetService(), op.GetOperation())
}

func (t *traceGraph) searchIngresses(n *node, result []*api_v1.Operation, hasSearched set.Set) []*api_v1.Operation {
	if hasSearched.Has(n) {
		splitStr := make([]string, len(result))
		for i, rn := range result {
			splitStr[len(result)-1-i] = rn.String()
		}
		t.logger.Fatal("cycled call found",
			zap.String("call path", strings.Join(splitStr, "->")+n.String()))
		return result
	}

	hasSearched.Add(n)
	if n.HasIn(t.globalRoot) && t.globalRoot.HasOut(n) {
		result = append(result, n.operation)
	} else {
		for _, next := range n.in.All() {
			result = t.searchIngresses(next, result, hasSearched)
		}
	}
	return result
}

func generateTrace(root *node) *api_v1.TraceNode {
	if root != nil {
		tn := &api_v1.TraceNode{
			Name:     fmt.Sprintf("%s:%s", root.operation.Service, root.operation.Operation),
			Children: make([]*api_v1.TraceNode, 0),
		}
		for _, outNode := range root.out.All() {
			if ret := generateTrace(outNode); ret != nil {
				tn.Children = append(tn.Children, ret)
			}
		}
		return tn
	}
	return nil
}
