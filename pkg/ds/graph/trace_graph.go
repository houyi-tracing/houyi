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
	nodes    NodeMap
	timer    ExpirationTimer
	duration time.Duration
	logger   *zap.Logger

	// fakeRoot is used to mark some trace graph node as true roots.
	// True roots have below features:
	//  1. The operation of those roots would be called by other operations,
	//  2. Or, the operation of those roots would be called directly by user requests.
	// ATTENTION: fakeRoot is not in traceGraph.nodes.
	fakeRoot ExecutionGraphNode
}

func NewExecutionGraph(logger *zap.Logger, duration time.Duration) *traceGraph {
	return &traceGraph{
		nodes:    NewNodeMap(),
		timer:    NewExpirationTimer(),
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

	// from -> to
	from.AddOut(to)
	to.AddIn(from)
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

func (t *traceGraph) RemoveExpired() {
	t.lock.Lock()
	defer t.lock.Unlock()

	for _, node := range t.nodes.AllNodes() {
		if expired, err := t.timer.IsExpired(node); expired && err == nil {
			t.removeNode(node)
			t.logger.Debug("remove expired node in trace graph",
				zap.String("service", node.Service()),
				zap.String("operation", node.Operation()))
		} else {
			if err != nil {
				t.logger.Error(
					"graph contains node but the node does not exist in timer",
					zap.Error(err))
			}
		}
	}
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

	// remove node
	t.nodes.Remove(node)

	// remove timer
	t.timer.Remove(node)
}

func (t *traceGraph) searchRoots(node ExecutionGraphNode, roots []ExecutionGraphNode, hasSearched set.Set) []ExecutionGraphNode {
	// fakeRoot is used to mark some trace graph node as roots.
	hasSearched.Add(node)

	if node.OutN() == 0 || node.HasOut(t.fakeRoot) {
		roots = append(roots, node)
	}

	if node.OutN() != 0 {
		for _, next := range node.GetOuts() {
			if next != t.fakeRoot && !hasSearched.Has(next) {
				roots = t.searchRoots(next, roots, hasSearched)
			}
		}
	}
	return roots
}
