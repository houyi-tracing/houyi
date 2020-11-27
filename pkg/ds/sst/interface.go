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

type SampleStrategyTree interface {
	// Add add a new operation into this SST.
	// It would be insert as a descendant ot root and be close to root as possible.
	Add(string, string)

	// AddService adds a new service into this SST.
	AddService(string)

	// Has returns true if operation exist, else false
	Has(string, string) bool

	// HasService returns true if service exist, else false.
	HasService(string) bool

	// Promote promotes a node in this SST.
	// The promoted node would increase a level and if the grand parent has not room to add, the LRU would be demoted.
	Promote(string, string) error

	// GetSamplingRate returns sample strategy the SST represents.
	// Upper bound is the maximum probability of single label, and the lower bound is minimum one
	// and 0 <= (lower bound) <= (upper bound) <= 1.
	//
	// ATTENTION: Due to the definition of Sample Strategy Tree, the probability of the direct child nodes of root node
	// could not exceed 1/N, where N is the immutable number of maximum Children for whole SST. Thus, the result upper
	// bound would be min[inputted upper bound, 1/N].
	GetSamplingRate() map[string]map[string]float64

	// GetServiceSamplingRate returns the sampling rates of all operation of inputted service.
	// It would returns error if service does not exist in SST.
	GetServiceSamplingRate(string) (map[string]float64, error)

	// GetOperationSamplingRate returns sampling rate of single operation.
	// It would return error if this operation not exist.
	GetOperationSamplingRate(string, string) (float64, error)

	// All returns a slice that contains all services in SST.
	All() map[string][]string

	// Remove removes a operation from SST.
	Remove(string, string)

	// Prune removes leaf nodes that has probability lower than the threshold.
	Prune(float64)
}

type TreeNode interface {
	// AddChild add a node as a child of current node.
	// This method would return a previous least-recently-used child node of current node if there has not room for adding
	// a new child for current node.
	AddChild(TreeNode) TreeNode

	// Children return a slice that contains all child nodes of current node.
	Children() []TreeNode

	// ChildrenN returns the number of child nodes of current node.
	ChildrenN() int

	// HasChild returns true if input node is child of current node.
	HasChild(TreeNode) bool

	// HasRoom returns true if the number of current node's children has not exceed.
	// the maximum, else false.
	HasRoom() bool

	// InsertAsDesc insert a leaf node as a descendant of current node.
	InsertAsDesc(TreeNode)

	// IsLeaf returns true if current treeNode is a leaf treeNode, else false.
	IsLeaf() bool

	// Parent returns the parent node of current node.
	Parent() TreeNode

	// PathDepression removes current node which has only one child node and makes it's only child node to replace it.
	PathDepression()

	// RemoveChild deletes the input child node of current node.
	RemoveChild(TreeNode)

	// Service returns the service this tree node relate to.
	Service() string

	// Operation returns the operation this tree node relate to.
	Operation() string
}

type TreeNodeMap interface {
	// AddSvc adds a new service into map.
	AddSvc(svc string)

	// AddOp adds a new operation into map
	AddOp(svc string, op string, node TreeNode)

	// Delete deletes a operation in this map.
	// Nothing would happen if the operation does not exist.
	Delete(svc string, op string)

	// Has returns true if the operation is in this map, else false.
	HasOp(svc string, op string) bool

	// HasSvc returns true if the service is in this map, else false.
	HasSvc(svc string) bool

	// Update updates the TreeNode of operation
	Update(svc string, op string, node TreeNode)

	// Get returns the TreeNode relate to the operation.
	// The returned error would not be nil if operation does not exist in this map.
	GetNode(svc string, op string) (TreeNode, error)

	// GetOps returns a map[string]TreeNode which contains all operations of inputted service in this map.
	// The returned error would not be nil if the service does not exist in this map.
	GetOps(svc string) (map[string]TreeNode, error)

	// Size returns the number of operations in this map.
	Size() int

	// Size returns all operations in this map with type map[string][]string.
	All() map[string][]string
}
