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

import "fmt"

// service -> operation -> trace graph node
type nodeMap map[string]map[string]ExecutionGraphNode

func newNodeMap() graphNodeMap {
	ret := make(nodeMap)
	return &ret
}

func (ns *nodeMap) Add(node ExecutionGraphNode) {
	svc, op := node.Service(), node.Operation()
	if _, has := (*ns)[svc]; !has {
		(*ns)[svc] = make(map[string]ExecutionGraphNode)
	}
	(*ns)[svc][op] = node
}

func (ns *nodeMap) Has(node ExecutionGraphNode) bool {
	return ns.has(node.Service(), node.Operation())
}

func (ns *nodeMap) Get(service, operation string) (ExecutionGraphNode, error) {
	if ns.has(service, operation) {
		return (*ns)[service][operation], nil
	} else {
		return nil, fmt.Errorf("operation not exist: service=%v, operation=%v", service, operation)
	}
}

func (ns *nodeMap) Remove(node ExecutionGraphNode) {
	svc, op := node.Service(), node.Operation()
	if ns.has(svc, op) {
		delete((*ns)[svc], op)
	}
}

func (ns *nodeMap) Clear() {
	*ns = make(nodeMap)
}

func (ns *nodeMap) Size() int {
	size := 0
	for _, innerMap := range *ns {
		size += len(innerMap)
	}
	return size
}

func (ns *nodeMap) AllServices() []string {
	retMe := make([]string, 0, len(*ns))
	for svc := range *ns {
		retMe = append(retMe, svc)
	}
	return retMe
}

func (ns *nodeMap) AllOperations() map[string][]string {
	retMe := make(map[string][]string)
	for svc, opMap := range *ns {
		retMe[svc] = make([]string, 0, len(opMap))
		for op := range opMap {
			retMe[svc] = append(retMe[svc], op)
		}
	}
	return retMe
}

func (ns *nodeMap) AllNodes() []ExecutionGraphNode {
	retMe := make([]ExecutionGraphNode, 0)
	for _, opMap := range *ns {
		for _, node := range opMap {
			retMe = append(retMe, node)
		}
	}
	return retMe
}

func (ns *nodeMap) has(service, operation string) bool {
	if _, has := (*ns)[service]; has {
		if _, ok := (*ns)[service][operation]; ok {
			return true
		}
	}
	return false
}
