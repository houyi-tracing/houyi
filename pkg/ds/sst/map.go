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

import "fmt"

// treeNodeMap contains pointer of tree nodes of all operations in services
// This map has not lock to avoid concurrent read & write to map, so that the outside must hold a RWMutex to operate this map.
type treeNodeMap map[string]map[string]TreeNode

func NewTreeNodeMap() TreeNodeMap {
	retMe := make(treeNodeMap)
	return &retMe
}

func (m *treeNodeMap) HasOp(svc, op string) bool {
	return m.has(svc, op)
}

func (m *treeNodeMap) HasSvc(svc string) bool {
	return m.hasService(svc)
}

func (m *treeNodeMap) AddOp(svc, op string, node TreeNode) {
	if !m.hasService(svc) {
		(*m)[svc] = make(map[string]TreeNode)
	}
	(*m)[svc][op] = node
}

func (m *treeNodeMap) AddSvc(svc string) {
	if !m.hasService(svc) {
		(*m)[svc] = make(map[string]TreeNode)
	}
}

func (m *treeNodeMap) Delete(svc, op string) {
	if m.has(svc, op) {
		delete((*m)[svc], op)
	}
}

func (m *treeNodeMap) Update(svc, op string, node TreeNode) {
	if m.has(svc, op) {
		(*m)[svc][op] = node
	}
}

func (m *treeNodeMap) GetNode(svc, op string) (TreeNode, error) {
	if m.has(svc, op) {
		return (*m)[svc][op], nil
	} else {
		return nil, fmt.Errorf("%s, %s not exist in svc map", svc, op)
	}
}

func (m *treeNodeMap) Size() int {
	sum := 0
	for _, innerMap := range *m {
		sum += len(innerMap)
	}
	return sum
}

func (m *treeNodeMap) All() map[string][]string {
	all := make(map[string][]string)
	for service, innerMap := range *m {
		for op := range innerMap {
			all[service] = append(all[service], op)
		}
	}
	return all
}

func (m *treeNodeMap) GetOps(svc string) (map[string]TreeNode, error) {
	if m.hasService(svc) {
		return (*m)[svc], nil
	} else {
		return nil, fmt.Errorf("service not exist: %v", svc)
	}
}

func (m *treeNodeMap) has(s, o string) bool {
	if _, has := (*m)[s]; has {
		_, ok := (*m)[s][o]
		return ok
	}
	return false
}

func (m *treeNodeMap) hasService(s string) bool {
	_, has := (*m)[s]
	return has
}
