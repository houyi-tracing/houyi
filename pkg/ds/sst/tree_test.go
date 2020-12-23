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
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"math"
	"os"
	"testing"
)

func checkSampleStrategy(t *testing.T, sampleStrategy map[string]map[string]float64) {
	sum := 0.0
	for _, innerMap := range sampleStrategy {
		for _, v := range innerMap {
			sum += v
		}
	}
	assert.Less(t, math.Abs(sum-1.0), 0.00000001)
}

func dfs(root TreeNode, t *testing.T) int {
	sum := 0
	if !root.IsLeaf() {
		for _, c := range root.Children() {
			if c.IsLeaf() {
				sum += 1
			} else {
				sum += dfs(c, t)
			}
		}
	}
	assert.Equal(t, sum, root.(*treeNode).leafCnt)
	return sum
}

func checkLeafCnt(sst SampleStrategyTree, t *testing.T) {
	dfs(sst.(*sampleStrategyTree).root, t)
}

func TestAdd(t *testing.T) {
	logger, _ := zap.NewProduction()
	sst := NewSampleStrategyTree(4, logger)

	newSvc := "svc-1"
	sst.AddService(newSvc)
	assert.True(t, sst.HasService(newSvc))

	newOp := "op-1"
	sst.Add(newSvc, newOp)
	assert.True(t, sst.Has(newSvc, newOp))

	sr := sst.GetSamplingRate()
	assert.Equal(t, 1.0, sr[newSvc][newOp])

	newOp = "op-2"
	sst.Add(newSvc, newOp)
	assert.True(t, sst.Has(newSvc, newOp))

	sr = sst.GetSamplingRate()
	assert.Equal(t, 0.5, sr[newSvc][newOp])
}

func TestPrune(t *testing.T) {
	logger, _ := zap.NewProduction()
	sst := NewSampleStrategyTree(4, logger)

	ops := []string{"op1", "op2", "op3", "op4", "op5"}

	svc := "svc"
	sst.AddService(svc)
	assert.True(t, sst.HasService(svc))
	for _, op := range ops {
		sst.Add(svc, op)
		assert.True(t, sst.Has(svc, op))
	}

	svcSr, _ := sst.GetServiceSamplingRate(svc)
	fmt.Printf("%v\n", svcSr)

	sr := sst.GetSamplingRate()
	assert.Equal(t, 0.125, sr[svc]["op1"])
	assert.Equal(t, 0.25, sr[svc]["op2"])
	assert.Equal(t, 0.25, sr[svc]["op3"])
	assert.Equal(t, 0.25, sr[svc]["op4"])
	assert.Equal(t, 0.125, sr[svc]["op5"])

	sst.Prune(0.25)
	assert.False(t, sst.Has(svc, "op1"))
	assert.True(t, sst.Has(svc, "op2"))
	assert.True(t, sst.Has(svc, "op3"))
	assert.True(t, sst.Has(svc, "op4"))
	assert.False(t, sst.Has(svc, "op5"))

	sst.Prune(1)
	assert.Equal(t, 0, len(sst.All()))

	for _, s := range []string{"svc-1", "svc-2", "svc-3", "svc-4", "svc-5", "svc-6"} {
		for _, op := range ops {
			sst.Add(s, op)
		}
	}

	checkLeafCnt(sst, t)

	allSr := sst.GetSamplingRate()
	fmt.Printf("%v\n", allSr)

	sst.Prune(0.15)
	checkLeafCnt(sst, t)
	checkSampleStrategy(t, allSr)
}

func TestPromote(t *testing.T) {
	logger, _ := zap.NewProduction()
	sst := NewSampleStrategyTree(4, logger)

	svc := "svc"
	sst.AddService(svc)
	assert.True(t, sst.HasService(svc))
	for _, op := range []string{"op1", "op2", "op3", "op4", "op5"} {
		sst.Add(svc, op)
		assert.True(t, sst.Has(svc, op))
	}

	sr := sst.GetSamplingRate()
	fmt.Printf("sampling rate:%v\n", sr)

	err := sst.Promote(svc, "op5")
	assert.NoError(t, err)
	sr = sst.GetSamplingRate()
	fmt.Printf("sampling rate:%v\n", sr)
	assert.Greater(t, sr[svc]["op5"], sr[svc]["op1"])

	checkLeafCnt(sst, t)
	checkSampleStrategy(t, sst.GetSamplingRate())
}

func TestSST(t *testing.T) {
	logger, _ := zap.NewProduction()
	svc := "svc"
	ops := []string{"A", "B", "C", "D", "E", "F"}
	sst := NewSampleStrategyTree(2, logger)
	for _, op := range ops {
		sst.Add(svc, op)
	}

	sr := sst.GetSamplingRate()
	for _, op := range ops {
		fmt.Printf("%f\t", sr[svc][op])
	}
	fmt.Println()

	for _, op := range ops {
		_ = sst.Promote(svc, op)
		sr = sst.GetSamplingRate()
		for _, op := range ops {
			fmt.Printf("%f\t", sr[svc][op])
		}
		fmt.Println()
	}
}

func TestSSTLargeScale(t *testing.T) {
	logger, _ := zap.NewProduction()
	svc := "svc"
	n := 10
	maxCN := 3

	opts := make([]string, 0)
	for i := 0; i < n; i++ {
		opts = append(opts, fmt.Sprintf("op%d", i))
	}

	sst := NewSampleStrategyTree(maxCN, logger)
	for _, op := range opts {
		sst.Add(svc, op)
	}

	f, err := os.OpenFile("sst_test.csv", os.O_TRUNC|os.O_CREATE, 0655)

	writeToFile := func(sr map[string]map[string]float64) {
		if err != nil {
			logger.Error("", zap.Error(err))
		}
		for i, op := range opts {
			f.WriteString(fmt.Sprintf("%f", sr[svc][op]))
			if i != len(opts)-1 {
				f.WriteString(",")
			} else {
				f.WriteString("\n")
			}
		}
	}

	check := func(sr map[string]map[string]float64) bool {
		sum := 0.0
		for _, opMap := range sr {
			for _, srVal := range opMap {
				sum += srVal
			}
		}
		return (sum - 1.0) < 0.00001
	}

	writeToFile(sst.GetSamplingRate())

	toPromote := map[int]int{
		0: 4,
		1: 4,
		2: 4,
		3: 4,
		4: 4,
		5: 4,
		6: 4,
		7: 4,
		8: 4,
		9: 4,
	}

	for n, opN := range toPromote {
		f.WriteString(fmt.Sprintf("PROMOTE: %d\n", n))
		for i := 0; i < opN; i++ {
			sst.Promote(svc, fmt.Sprintf("op%d", n))
			sr := sst.GetSamplingRate()
			writeToFile(sr)
			assert.True(t, check(sr))
		}
		f.WriteString("\n")
	}

	f.Close()
}
