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
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

//const (
//	defaultExpireDuration = time.Minute * 1
//)

//
//func TestAddService(t *testing.T) {
//	g := NewExecutionGraph(defaultExpireDuration)
//	g.AddOp("hello")
//	all := g.All()
//	assert.Equal(t, 1, len(all))
//	fmt.Println(all)
//
//	g.AddOp("world")
//	all = g.All()
//	assert.Equal(t, 2, len(all))
//	fmt.Println(all)
//
//	g.AddEdge("hello", "world")
//	assert.True(t, g.HasEdge("hello", "world"))
//	assert.True(t, g.HasEdge("world", "hello"))
//	assert.Equal(t, "world", g.GetRootsOf("hello")[0])
//
//	g.AddOp("register")
//	g.AddEdge("world", "register")
//	assert.True(t, g.HasEdge("register", "world"))
//	assert.True(t, g.HasEdge("world", "register"))
//	assert.True(t, g.HasEdgeFromTo("world", "register"))
//	assert.False(t, g.HasEdgeFromTo("register", "world"))
//	assert.Equal(t, "register", g.GetRootsOf("hello")[0])
//
//	g.Remove("world")
//	assert.False(t, g.HasEdgeFromTo("world", "register"))
//	assert.False(t, g.HasEdgeFromTo("register", "world"))
//	fmt.Println(g.All())
//}
//
//func TestAddRoot(t *testing.T) {
//	tg := NewExecutionGraph(defaultExpireDuration)
//	rootOp := "root"
//	child1 := "child1"
//	child2 := "child2"
//
//	tg.AddRoot(rootOp)
//	tg.AddEdge(child1, rootOp)
//	assert.Equal(t, rootOp, tg.GetRootsOf(child1)[0])
//
//	tg.AddRoot(child1)
//	tg.AddEdge(child2, child1)
//	rootsOfChild2 := tg.GetRootsOf(child2)
//	assert.Equal(t, child1, rootsOfChild2[0])
//	assert.Equal(t, 2, len(rootsOfChild2))
//}
//
//func TestRemoveExpired(t *testing.T) {
//	tg := NewExecutionGraph(time.Second * 1)
//	tg.AddOp("test1")
//	tg.AddOp("test2")
//	time.Sleep(time.Second * 5)
//	tg.RemoveExpired()
//	assert.Equal(t, 0, tg.Size())
//}

func TestGetRoots(t *testing.T) {
	logger, _ := zap.NewProduction()
	tg := NewExecutionGraph(logger, time.Minute)

	tg.Add("svc1", "op1")
	tg.Add("svc2", "op2")
	tg.Add("svc3", "op3")
	node1, _ := tg.GetNode("svc1", "op1")
	node2, _ := tg.GetNode("svc2", "op2")
	node3, _ := tg.GetNode("svc3", "op3")
	tg.AddEdge(node1, node2)
	tg.AddEdge(node2, node3)
	tg.AddRoot(node1)
	tg.AddRoot(node2)

	roots := tg.GetRootsOf(node3)
	for _, r := range roots {
		fmt.Println(r)
	}
}
