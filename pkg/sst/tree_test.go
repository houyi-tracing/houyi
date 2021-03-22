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

package sst

import (
	"fmt"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

const (
	maxN = 4
)

func TestAddDuplicatedOperations(t *testing.T) {
	tree := NewSamplingStrategyTree(maxN)
	opsN := 100

	for i := 0; i < opsN; i++ {
		err := tree.Add(&api_v1.Operation{
			Service:   "",
			Operation: "",
		})
		if i != 0 {
			assert.Error(t, err)
		}
	}
	assert.Equal(t, 1, tree.(*sst).root.childN())
}

func TestAddManyNodes(t *testing.T) {
	tree := NewSamplingStrategyTree(maxN)

	opsN := 1000000
	for i := 0; i < opsN; i++ {
		err := tree.Add(&api_v1.Operation{
			Service:   "",
			Operation: fmt.Sprintf("%d", i),
		})
		assert.Nil(t, err)
	}
}

func TestPrune(t *testing.T) {
	tree := NewSamplingStrategyTree(maxN)

	opsN := 1000000
	for i := 0; i < opsN; i++ {
		err := tree.Add(&api_v1.Operation{
			Service:   "",
			Operation: fmt.Sprintf("%d", i),
		})
		assert.Nil(t, err)

		if rand.Float64() < 0.1 {
			deleteOp := &api_v1.Operation{
				Service:   "",
				Operation: fmt.Sprintf("%d", rand.Intn(i)),
			}
			flag := tree.Has(deleteOp)
			err = tree.Prune(deleteOp)
			if flag {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		}
	}

	root := tree.(*sst).root
	check(root, root, t)
}

func TestSumOfSamplingRatesMustEqualToOne(t *testing.T) {
	tree := NewSamplingStrategyTree(maxN)

	opsN := 1000000
	for i := 0; i < opsN; i++ {
		err := tree.Add(&api_v1.Operation{
			Service:   "",
			Operation: fmt.Sprintf("%d", i),
		})
		assert.Nil(t, err)
	}

	ops := tree.(*sst).nodes.allOperations()
	sum := 0.0
	for _, op := range ops {
		s, err := tree.Generate(op)
		assert.Nil(t, err)
		sum += s.SamplingRate
	}
	absErr := 1e-10
	assert.Less(t, math.Abs(1.0-sum), absErr)
}

func TestSamplingRateMustBeGreaterAfterPromoting(t *testing.T) {
	tree := NewSamplingStrategyTree(maxN)

	opsN := 1000000
	for i := 0; i < opsN; i++ {
		err := tree.Add(&api_v1.Operation{
			Service:   "",
			Operation: fmt.Sprintf("%d", i),
		})
		assert.Nil(t, err)
	}

	root := tree.(*sst).root
	check(root, root, t)

	times := 1000000
	for i := 0; i < times; i++ {
		promoteOp := &api_v1.Operation{
			Service:   "",
			Operation: fmt.Sprintf("%d", rand.Intn(opsN)),
		}
		if tree.Has(promoteOp) {
			s, err := tree.Generate(promoteOp)
			assert.Nil(t, err)
			err = tree.Promote(promoteOp)
			assert.Nil(t, err)
			newS, _ := tree.Generate(promoteOp)
			assert.LessOrEqual(t, s.SamplingRate, newS.SamplingRate)
		}
	}

	check(root, root, t)
}

func Test(t *testing.T) {
	order := 4
	N := 50
	times := 40

	tree := NewSamplingStrategyTree(order)

	ops := make([]*api_v1.Operation, 0, N)

	for i := 0; i < N; i++ {
		op := &api_v1.Operation{
			Service:   "svc",
			Operation: fmt.Sprintf("op_%d", i),
		}
		tree.Add(op)
		ops = append(ops, op)
	}

	for k := 0; k < N; k++ {
		sr, _ := tree.Generate(ops[k])
		fmt.Printf("%f", sr.SamplingRate)
		if k != N-1 {
			fmt.Printf("\t")
		} else {
			fmt.Printf("\n")
		}
	}

	for i := 0; i < N*2; i++ {
		for j := 0; j < times; j++ {
			tree.Promote(ops[i%N])
		}
		for k := 0; k < N; k++ {
			sr, _ := tree.Generate(ops[k])
			fmt.Printf("%.12f", sr.SamplingRate)
			if k != N-1 {
				fmt.Printf("\t")
			} else {
				fmt.Printf("\n")
			}
		}
	}
}

func check(root, currNode *treeNode, t *testing.T) int {
	sum := 0
	for _, n := range currNode.childNodes.all() {
		if root != currNode {
			assert.Greater(t, currNode.childN(), 1)
		}
		if n.isLeaf() {
			assert.Equal(t, 1, n.leafCnt)
			sum++
		} else {
			sum += check(root, n, t)
		}
	}
	assert.Equal(t, sum, currNode.leafCnt)
	return sum
}
