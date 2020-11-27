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

import "time"

type ExecutionGraphNode interface {
	Service() string
	Operation() string
	GetOuts() []ExecutionGraphNode
	GetIns() []ExecutionGraphNode
	OutN() int
	InN() int
	HasIn(node ExecutionGraphNode) bool
	HasOut(node ExecutionGraphNode) bool
	AddIn(node ExecutionGraphNode)
	AddOut(node ExecutionGraphNode)
	RemoveIn(node ExecutionGraphNode)
	RemoveOut(node ExecutionGraphNode)
}

type ExecutionGraph interface {
	Add(string, string) ExecutionGraphNode
	AddEdge(from, to ExecutionGraphNode)
	AddRoot(root ExecutionGraphNode)
	AllServices() []string
	AllOperations() map[string][]string
	GetNode(service, operation string) (ExecutionGraphNode, error)
	GetRootsOf(node ExecutionGraphNode) []ExecutionGraphNode
	Has(string, string) bool
	HasEdge(from, to ExecutionGraphNode) bool
	RemoveExpired()
	Remove(rmNode ExecutionGraphNode)
	RemoveEdge(from, to ExecutionGraphNode)
	Size() int
}

type ExpirationTimer interface {
	Timing(ExecutionGraphNode, time.Duration)
	IsExpired(node ExecutionGraphNode) (bool, error)
	Remove(node ExecutionGraphNode)
}

type NodeMap interface {
	Add(node ExecutionGraphNode)
	Has(node ExecutionGraphNode) bool
	Get(service string, operation string) (ExecutionGraphNode, error)
	Remove(node ExecutionGraphNode)
	Clear()
	Size() int
	AllServices() []string
	AllOperations() map[string][]string
	AllNodes() []ExecutionGraphNode
}
