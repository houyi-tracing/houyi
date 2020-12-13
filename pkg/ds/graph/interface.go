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
	// Service returns the service name of.
	Service() string

	// Operation returns the operation name.
	Operation() string

	// GetOuts returns all out nodes.
	GetOuts() []ExecutionGraphNode

	// GetIns returns all in nodes.
	GetIns() []ExecutionGraphNode

	// OutN returns the number of out nodes.
	OutN() int

	// InN returns the number of in nodes.
	InN() int

	// HasIn returns true if inputted node is in the in nodes.
	HasIn(node ExecutionGraphNode) bool

	// HasOut returns true if inputted node is in the out nodes.
	HasOut(node ExecutionGraphNode) bool

	// AddIn adds the node as an in node of current node. The other node must call AddOut.
	AddIn(node ExecutionGraphNode)

	// AddOut adds the node as an out node of current node. The other node must call AddIn.
	AddOut(node ExecutionGraphNode)

	// RemoveIn removes the in node of current node. The other node must call RemoveOut.
	RemoveIn(node ExecutionGraphNode)

	// RemoveOut removes the out node of current node. The other node must call RemoveIn.
	RemoveOut(node ExecutionGraphNode)

	String() string
}

type ExecutionGraph interface {
	// Add adds a new operation into this graph. If operation is already in this graph, it would returns it's node.
	Add(string, string) ExecutionGraphNode

	// AddEdge adds a edge from a node to another node.
	AddEdge(from, to ExecutionGraphNode)

	// AddRoot makes a node which is already in this graph as a new root.
	// A root node is a node relate to operation directly called by outside users.
	AddRoot(root ExecutionGraphNode)

	// AllServices returns names of all services in this graph.
	AllServices() []string

	// AllOperations returns a map[string][]string containing all operations with their services.
	AllOperations() map[string][]string

	// GetNode returns the graph node via service name and operation name. It would returns error if operation does not
	// exist in this graph.
	GetNode(service, operation string) (ExecutionGraphNode, error)

	// GetRootsOf returns  all roots of inputted node.
	GetRootsOf(node ExecutionGraphNode) []ExecutionGraphNode

	// Has returns true if operation is in this graph, else false.
	Has(string, string) bool

	// HasEdge returns true if there is one and only one edge between inpputed two nodes.
	HasEdge(from, to ExecutionGraphNode) bool

	// IsRoot returns true if this operation of service svc is root node, else false.
	IsRoot(svc, op string) bool

	// RemoveExpired removes all expired nodes and their edges.
	RemoveExpired() []ExecutionGraphNode

	// Remove removes a node from this graph. Nothing would happen when remove a node not in this graph.
	Remove(rmNode ExecutionGraphNode)

	// RemoveEdge removes the edge between inputted two nodes.
	RemoveEdge(from, to ExecutionGraphNode)

	// Size returns the number of operations in this graph.
	Size() int

	// Refresh reset the timer for inputted node.
	Refresh(node ExecutionGraphNode)
}

type nodeExpirationTimer interface {
	// Timing would reset the timer for inputted node.
	// This node will expire at the time = now + duration.
	Timing(ExecutionGraphNode, time.Duration)

	// IsExpired returns true if node was expired, else false.
	// The returned error would not be nil if node was not in this timer.
	IsExpired(node ExecutionGraphNode) (bool, error)

	// Remove removes a node from this timer.
	Remove(node ExecutionGraphNode)
}

type graphNodeMap interface {
	// Add adds a new node into map.
	Add(node ExecutionGraphNode)

	// Has returns true if node is already in this map, else false.
	Has(node ExecutionGraphNode) bool

	// Get returns the node relate to the operation. It would return error if this operation is not in this map.
	Get(service string, operation string) (ExecutionGraphNode, error)

	// Remove removes a node from this map.
	Remove(node ExecutionGraphNode)

	// Clear removes all nodes in this map.
	Clear()

	// Size returns the number of operations in this map.
	Size() int

	// AllServices returns all services in this map.
	AllServices() []string

	// AllOperations returns a map[string][]string (service->operations) containing all operations with their services.
	AllOperations() map[string][]string

	// AllNodes returns all nodes in this map.
	AllNodes() []ExecutionGraphNode
}
