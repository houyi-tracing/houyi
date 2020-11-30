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

import (
	"github.com/yan-fuhai/go-ds/set"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	defaultRootParentService   = "root-service"
	defaultRootParentOperation = "root-operation"
)

type traceGraph struct {
	lock     sync.RWMutex // make trace graph to be thread-safe
	nodes    graphNodeMap
	timer    nodeExpirationTimer
	duration time.Duration
	logger   *zap.Logger

	// fakeRoot is used to mark some trace graph nodes as true roots.
	// True roots have below features:
	//  1. The operation of those roots would be called by other operations,
	//  2. Or, the operation of those roots would be called directly by user requests.
	// ATTENTION: fakeRoot is not in traceGraph.nodes.
	fakeRoot ExecutionGraphNode
}

func NewExecutionGraph(logger *zap.Logger, duration time.Duration) *traceGraph {
	return &traceGraph{
		nodes:    newNodeMap(),
		timer:    newExpirationTimer(),
		duration: duration,
		logger:   logger,
		fakeRoot: NewExecutionGraphNode(defaultRootParentService, defaultRootParentOperation),
	}
}

func (t *traceGraph) Add(service, operation string) ExecutionGraphNode {
	t.lock.Lock()
	defer t.lock.Unlock()

	if node, err := t.nodes.Get(service, operation); err != nil {
		newNode := NewExecutionGraphNode(service, operation)
		t.nodes.Add(newNode)
		return newNode
	} else {
		return node
	}
}

func (t *traceGraph) AddEdge(from, to ExecutionGraphNode) {
	t.lock.Lock()
	defer t.lock.Unlock()
	defer t.timer.Timing(from, t.duration)
	defer t.timer.Timing(to, t.duration)

	if from != to {
		// from -> to
		from.AddOut(to)
		to.AddIn(from)
	}
}

func (t *traceGraph) AddRoot(root ExecutionGraphNode) {
	t.lock.Lock()
	defer t.lock.Unlock()
	defer t.timer.Timing(root, t.duration)

	// root -> t.fakeRoot
	t.fakeRoot.AddIn(root)
	root.AddOut(t.fakeRoot)
}

func (t *traceGraph) AllServices() []string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.nodes.AllServices()
}

func (t *traceGraph) AllOperations() map[string][]string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.nodes.AllOperations()
}

func (t *traceGraph) GetNode(service string, operation string) (ExecutionGraphNode, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if ret, err := t.nodes.Get(service, operation); err == nil {
		t.timer.Timing(ret, t.duration)
		return ret, nil
	} else {
		return nil, err
	}
}

func (t *traceGraph) GetRootsOf(node ExecutionGraphNode) []ExecutionGraphNode {
	t.lock.RLock()
	defer t.lock.RUnlock()
	defer t.timer.Timing(node, t.duration)

	retMe := make([]ExecutionGraphNode, 0)
	if t.nodes.Has(node) {
		// TODO
	}
	return retMe
}

func (t *traceGraph) Has(service string, operation string) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.has(service, operation)
}

func (t *traceGraph) HasEdge(from, to ExecutionGraphNode) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	defer t.timer.Timing(from, t.duration)
	defer t.timer.Timing(to, t.duration)

	return from.HasOut(to) && to.HasIn(from)
}

func (t *traceGraph) RemoveExpired() []ExecutionGraphNode {
	t.lock.Lock()
	defer t.lock.Unlock()

	removedNodes := make([]ExecutionGraphNode, 0)

	for _, node := range t.nodes.AllNodes() {
		if expired, err := t.timer.IsExpired(node); expired && err == nil {
			t.removeNode(node)
			t.logger.Debug("remove expired node in trace graph",
				zap.String("service", node.Service()),
				zap.String("operation", node.Operation()))
			removedNodes = append(removedNodes, node)
		} else {
			if err != nil {
				t.logger.Error(
					"graph contains node but the node does not exist in timer",
					zap.Error(err),
					zap.Stringer("node", node))
			}
		}
	}

	return removedNodes
}

func (t *traceGraph) Remove(rmNode ExecutionGraphNode) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.removeNode(rmNode)
}

func (t *traceGraph) RemoveEdge(from, to ExecutionGraphNode) {
	t.lock.Lock()
	defer t.lock.Unlock()

	from.RemoveOut(to)
	to.RemoveIn(from)
}

func (t *traceGraph) Size() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.nodes.Size()
}

func (t *traceGraph) Refresh(node ExecutionGraphNode) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.nodes.Has(node) {
		t.timer.Timing(node, t.duration)
	}
}

func (t *traceGraph) has(svc, op string) bool {
	if node, err := t.nodes.Get(svc, op); node != nil && err == nil {
		t.timer.Timing(node, t.duration)
		return true
	} else {
		return false
	}
}

func (t *traceGraph) removeNode(node ExecutionGraphNode) {
	// remove edges
	for _, from := range node.GetIns() {
		from.RemoveOut(node)
	}
	for _, to := range node.GetOuts() {
		to.RemoveIn(node)
	}

	// remove node from node map
	t.nodes.Remove(node)

	// remove node from timer
	t.timer.Remove(node)
}

func (t *traceGraph) searchRoots(node ExecutionGraphNode, roots []ExecutionGraphNode, hasSearched set.Set) []ExecutionGraphNode {
	if hasSearched.Has(node) {
		t.logger.Error("cycled call found",
			zap.String("current service", node.Service()),
			zap.String("current operation", node.Operation()))
		return roots
	}

	hasSearched.Add(node)

	// fakeRoot is used to mark some trace graph node as roots.
	if node.OutN() == 0 || node.HasOut(t.fakeRoot) {
		roots = append(roots, node)
	}

	if node.OutN() != 0 {
		for _, next := range node.GetOuts() {
			if next != t.fakeRoot {
				roots = t.searchRoots(next, roots, hasSearched)
			}
		}
	}
	return roots
}
