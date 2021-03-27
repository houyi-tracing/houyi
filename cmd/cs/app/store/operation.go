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

package store

import (
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
	"sync"
	"time"
)

type OperationStore interface {
	Start()
	Stop()

	UpToDate(op *api_v1.Operation, isIngress bool, qps float64)
	QpsWeight(op *api_v1.Operation) float64
}

type tItem struct {
	upSince   time.Time
	isIngress bool
	qps       float64
	op        *api_v1.Operation
}

type opStore struct {
	sync.RWMutex

	logger          *zap.Logger
	m               map[string]map[string]*tItem
	refreshInterval time.Duration
	sst             sst.SamplingStrategyTree
	tg              tg.TraceGraph
	seed            gossip.Seed
	stopChan        chan *sync.WaitGroup
}

func NewOperationStore(logger *zap.Logger, interval time.Duration, seed gossip.Seed) OperationStore {
	return &opStore{
		logger:          logger,
		m:               make(map[string]map[string]*tItem),
		refreshInterval: interval,
		seed:            seed,
		stopChan:        make(chan *sync.WaitGroup),
	}
}

func (t *opStore) Start() {
	go t.background(t.refreshInterval, t.removeExpired)
}

func (t *opStore) Stop() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	t.stopChan <- wg
	wg.Wait()
}

func (t *opStore) UpToDate(op *api_v1.Operation, isIngress bool, qps float64) {
	t.Lock()
	defer t.Unlock()

	svcName, opName := op.Service, op.Operation
	if _, has := t.m[svcName]; !has {
		t.m[svcName] = make(map[string]*tItem)
	}
	if item, has := t.m[svcName][opName]; has {
		item.upSince = time.Now()
		item.qps = qps
		item.isIngress = isIngress
	} else {
		t.m[svcName][opName] = &tItem{
			upSince:   time.Now(),
			op:        op,
			qps:       qps,
			isIngress: isIngress,
		}
	}
}

func (t *opStore) QpsWeight(op *api_v1.Operation) float64 {
	t.RLock()
	defer t.RUnlock()

	if t.has(op) {
		svcN, opN := op.Service, op.Operation
		curr := t.m[svcN][opN]
		if curr.qps == 0 || !curr.isIngress {
			return 1.0
		}
		total := 0.0
		for _, opMap := range t.m {
			for _, item := range opMap {
				if item.isIngress && item.qps != 0 {
					total += 1 / item.qps
				}
			}
		}
		return 1 / curr.qps / total
	}
	return 1.0
}

func (t *opStore) has(op *api_v1.Operation) bool {
	if _, hasSvc := t.m[op.Service]; hasSvc {
		if _, hasOp := t.m[op.Service][op.Operation]; hasOp {
			return true
		}
	}
	return false
}

func (t *opStore) background(interval time.Duration, f func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f()
		case wg := <-t.stopChan:
			wg.Done()
			return
		}
	}
}

func (t *opStore) removeExpired() {
	t.Lock()
	defer t.Unlock()

	now := time.Now()
	for svc, opMap := range t.m {
		for op, item := range opMap {
			if item.upSince.Add(t.refreshInterval).Before(now) {
				t.logger.Debug("monger expired operation",
					zap.String("service", svc), zap.String("operation", op))

				_ = t.sst.Prune(item.op)
				_ = t.tg.Remove(item.op)
				t.seed.MongerExpiredOperation(item.op)
				delete(t.m[svc], op)
				if len(t.m[svc]) == 0 {
					delete(t.m, svc)
				}
			}
		}
	}
}
