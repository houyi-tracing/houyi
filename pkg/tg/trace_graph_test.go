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

package tg

import (
	"encoding/json"
	"fmt"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestAutomaticallyFindEntries(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tg := NewTraceGraph(logger)

	op1 := &api_v1.Operation{
		Service:   "1",
		Operation: "1",
	}
	op2 := &api_v1.Operation{
		Service:   "2",
		Operation: "2",
	}
	op3 := &api_v1.Operation{
		Service:   "3",
		Operation: "3",
	}

	assert.Nil(t, tg.Add(op1))
	assert.NotNil(t, tg.Add(op1))
	assert.Nil(t, tg.Add(op2))
	assert.NotNil(t, tg.Add(op2))
	assert.Nil(t, tg.Add(op3))
	assert.NotNil(t, tg.Add(op3))
	assert.False(t, tg.IsIngress(op1))
	assert.False(t, tg.IsIngress(op2))
	assert.False(t, tg.IsIngress(op3))

	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: op1,
		To:   op2,
	}))
	assert.True(t, tg.IsIngress(op1))

	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: op2,
		To:   op3,
	}))
	assert.Nil(t, tg.RemoveRelation(&api_v1.Relation{
		From: op1,
		To:   op2,
	}))
	assert.False(t, tg.IsIngress(op1))
	assert.True(t, tg.IsIngress(op2))

	assert.Nil(t, tg.RemoveRelation(&api_v1.Relation{
		From: op2,
		To:   op3,
	}))
	assert.False(t, tg.IsIngress(op2))
}

func TestMultipleEntries(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tg := NewTraceGraph(logger)

	op1 := &api_v1.Operation{
		Service:   "1",
		Operation: "1",
	}
	op2 := &api_v1.Operation{
		Service:   "2",
		Operation: "2",
	}
	op3 := &api_v1.Operation{
		Service:   "3",
		Operation: "3",
	}
	op4 := &api_v1.Operation{
		Service:   "4",
		Operation: "4",
	}
	op5 := &api_v1.Operation{
		Service:   "5",
		Operation: "5",
	}

	assert.Nil(t, tg.Add(op1))
	assert.Nil(t, tg.Add(op2))
	assert.Nil(t, tg.Add(op3))
	assert.Nil(t, tg.Add(op4))
	assert.Nil(t, tg.Add(op5))

	// 1->3
	// 2->3
	// 3->4
	// 4->5
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: op1,
		To:   op3,
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: op2,
		To:   op3,
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: op3,
		To:   op4,
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: op4,
		To:   op5,
	}))
	entries, err := tg.GetIngresses(op5)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(entries))
}

func TestAddDuplicatedOperation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tg := NewTraceGraph(logger)

	op1 := &api_v1.Operation{
		Service:   "1",
		Operation: "1",
	}
	op2 := &api_v1.Operation{
		Service:   "1",
		Operation: "1",
	}

	assert.Nil(t, tg.Add(op1))
	assert.NotNil(t, tg.Add(op2))
	assert.Equal(t, 1, tg.Size())
}

func TestGenerateTraces(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tg := NewTraceGraph(logger)

	opN := 7
	ops := make([]*api_v1.Operation, 7)
	for i := 0; i < opN; i++ {
		ops[i] = &api_v1.Operation{
			Service:   fmt.Sprintf("%d", i),
			Operation: fmt.Sprintf("%d", i),
		}
		assert.Nil(t, tg.Add(ops[i]))
	}

	// 0
	// 1 	2
	// 3 4 	5 6
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: ops[0],
		To:   ops[1],
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: ops[0],
		To:   ops[2],
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: ops[1],
		To:   ops[3],
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: ops[1],
		To:   ops[4],
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: ops[2],
		To:   ops[5],
	}))
	assert.Nil(t, tg.AddRelation(&api_v1.Relation{
		From: ops[2],
		To:   ops[6],
	}))

	entries, err := tg.GetIngresses(ops[6])
	assert.Nil(t, err)
	assert.Equal(t, 1, len(entries))
	traces, err := tg.Traces(ops[5])
	assert.Nil(t, err)

	for _, trace := range traces {
		if bytes, err := json.Marshal(trace); err != nil {
			logger.Error("", zap.Error(err))
		} else {
			fmt.Println(string(bytes))
		}
	}
}
