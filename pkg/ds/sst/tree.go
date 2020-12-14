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

import (
	"fmt"
	"go.uber.org/zap"
)

type sampleStrategyTree struct {
	logger *zap.Logger

	maxChildrenN int
	root         TreeNode
	leafNodeMap  TreeNodeMap
}

// NewSampleStrategyTree returns a new SST pointer.
func NewSampleStrategyTree(maxChildrenN int, logger *zap.Logger) SampleStrategyTree {
	sst := &sampleStrategyTree{
		maxChildrenN: maxChildrenN,
		root:         newTreeNode(maxChildrenN, nil),
		leafNodeMap:  NewTreeNodeMap(),
		logger:       logger,
	}
	return sst
}

func (t *sampleStrategyTree) AddService(service string) {
	if !t.leafNodeMap.HasSvc(service) {
		t.leafNodeMap.AddSvc(service)
	}
}

func (t *sampleStrategyTree) Add(service string, operation string) {
	if !t.leafNodeMap.HasOp(service, operation) {
		leafNode := newLeafNode(service, operation, t.maxChildrenN, nil)
		t.root.InsertAsDesc(leafNode)
		t.leafNodeMap.AddOp(service, operation, leafNode)
	}
}

func (t *sampleStrategyTree) All() map[string][]string {
	return t.leafNodeMap.All()
}

func (t *sampleStrategyTree) Promote(service, operation string) error {
	if t.leafNodeMap.HasOp(service, operation) {
		promNode, err := t.leafNodeMap.GetNode(service, operation)
		if err != nil {
			return err
		}

		if promNode.Parent() == t.root {
			// the node locating at the top level of SST should not be promoted.
			return nil
		}

		grandParent, parent := promNode.Parent().Parent(), promNode.Parent()
		t.promoteChild(grandParent, parent, promNode)
		return nil
	} else {
		return fmt.Errorf("sample strategy tree has not serivce: %v", operation)
	}
}

func (t *sampleStrategyTree) GetSamplingRate() map[string]map[string]float64 {
	sampleStrategy := make(map[string]map[string]float64)
	t.getSampleStrategy(t.root, 1, &sampleStrategy)
	return sampleStrategy
}

func (t *sampleStrategyTree) GetServiceSamplingRate(service string) (map[string]float64, error) {
	if ops, err := t.leafNodeMap.GetOps(service); err == nil {
		retMe := make(map[string]float64)
		for op, tNode := range ops {
			retMe[op] = getSamplingRateFromBottomToTop(tNode)
		}
		return retMe, nil
	} else {
		return nil, err
	}
}

func (t *sampleStrategyTree) GetOperationSamplingRate(service, operation string) (float64, error) {
	if node, err := t.leafNodeMap.GetNode(service, operation); err == nil {
		return getSamplingRateFromBottomToTop(node), nil
	} else {
		return 0, err
	}
}

func (t *sampleStrategyTree) Remove(service, operation string) {
	if toRmNode, err := t.leafNodeMap.GetNode(service, operation); err == nil {
		if !toRmNode.IsLeaf() {
			t.logger.Error("for removing operation, the tree node relate to it to remove must be leaf node",
				zap.String("service", service),
				zap.String("operation", operation))
			return
		}
		parent := toRmNode.Parent()
		parent.RemoveChild(toRmNode)
		if parent != t.root {
			parent.PathDepression()
		}
		t.leafNodeMap.Delete(service, operation)
	}
}

func (t *sampleStrategyTree) Prune(threshold float64) {
	removeLeafNodeWithThreshold(nil, t.root, 1, threshold, t.leafNodeMap)
	updateLeafCnt(t.root)
}

func (t *sampleStrategyTree) Has(service, operation string) bool {
	return t.leafNodeMap.HasOp(service, operation)
}

func (t *sampleStrategyTree) HasService(service string) bool {
	return t.leafNodeMap.HasSvc(service)
}

// getSampleStrategy generates sample strategy from this SST and returns a map.
// The keys of map are label names and the value would be the sampling probability of operation relative with it.
func (t *sampleStrategyTree) getSampleStrategy(
	root TreeNode,
	curProbability float64,
	sampleStrategy *map[string]map[string]float64) {

	svc, op := root.Service(), root.Operation()
	if root.IsLeaf() {
		if _, has := (*sampleStrategy)[svc]; !has {
			(*sampleStrategy)[svc] = make(map[string]float64)
		}

		if _, has := (*sampleStrategy)[svc][op]; has {
			t.logger.Info("found cycled calls",
				zap.String("service", svc),
				zap.String("operation", op))
		} else {
			(*sampleStrategy)[svc][op] = curProbability
		}
	} else {
		n := float64(root.ChildrenN())
		for _, c := range root.Children() {
			t.getSampleStrategy(c, curProbability/n, sampleStrategy)
		}
	}
}

// promoteChild promotes a child node of current node as a direct child of current node's parent node.
func (t *sampleStrategyTree) promoteChild(grandParent, parent, child TreeNode) {
	if grandParent.HasChild(parent) && parent.HasChild(child) && child.IsLeaf() {
		parent.RemoveChild(child)
		demoteNode := grandParent.(*treeNode).childNodes.Add(child).(TreeNode)
		if demoteNode != nil {
			// demoteNode != nil means that grandParent node has not room to add child.
			// In this case, it would remove the LRU node among its child nodes and returns it.
			// To receive the returned node, current node parent would add it as a child.
			parent.AddChild(demoteNode)
		} else {
			parent.PathDepression()
		}
	}
}

// removeLeafNodeWithThreshold removes all leaf nodes whose probability is below the threshold.
func removeLeafNodeWithThreshold(parent, cur TreeNode, curProb, threshold float64, sMap TreeNodeMap) {
	if cur.IsLeaf() {
		if curProb < threshold {
			parent.RemoveChild(cur)
			sMap.Delete(cur.Service(), cur.Operation())
		}
	} else {
		n := float64(cur.ChildrenN())
		for _, c := range cur.Children() {
			removeLeafNodeWithThreshold(cur, c, curProb/n, threshold, sMap)
		}
		if parent != nil && cur.ChildrenN() == 0 {
			parent.RemoveChild(cur) // remove nodes that have not child nodes after pruning.
		}
	}
}

// updateLeafCnt updates the value of filed, leaf count, of every node in this tree by DFS algorithm.
func updateLeafCnt(root TreeNode) int {
	sum := 0
	if !root.IsLeaf() {
		for _, c := range root.Children() {
			if c.IsLeaf() {
				sum += 1
			} else {
				sum += updateLeafCnt(c)
			}
		}
	}
	if tNode, ok := root.(*treeNode); ok {
		tNode.leafCnt = sum
	}
	return sum
}

func getSamplingRateFromBottomToTop(leafNode TreeNode) float64 {
	if leafNode.IsLeaf() {
		ret := 1.0
		p := leafNode.Parent()
		for p != nil {
			ret /= float64(p.ChildrenN())
			p = p.Parent()
		}
		return ret
	}
	return 0
}
