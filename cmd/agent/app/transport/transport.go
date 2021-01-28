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

package transport

import (
	"context"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/routing"
	jaeger "github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// CollectorTransport reuses gRPC connection from agent to collector.
type CollectorTransport struct {
	logger *zap.Logger

	c    jaeger.CollectorServiceClient
	conn *grpc.ClientConn
	ep   *routing.Endpoint
}

func NewCollectorTransport(logger *zap.Logger, ep *routing.Endpoint) *CollectorTransport {
	ct := &CollectorTransport{
		logger: logger,
		ep:     ep,
	}
	conn, err := grpc.Dial(ep.String(), grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Failed to dial to remote collector", zap.Error(err))
		return nil
	}
	ct.conn = conn
	ct.c = jaeger.NewCollectorServiceClient(conn)
	return ct
}

func (t *CollectorTransport) PostSpans(ctx context.Context, req *jaeger.PostSpansRequest) (*jaeger.PostSpansResponse, error) {
	if t.conn.GetState() == connectivity.Shutdown {
		conn, err := grpc.Dial(t.ep.String(), grpc.WithInsecure())
		if err != nil {
			t.logger.Fatal("Connection is closed and client failed to dial to remote collector", zap.Error(err))
			return &jaeger.PostSpansResponse{}, err
		}
		t.conn = conn
		t.c = jaeger.NewCollectorServiceClient(conn)
	}
	return t.c.PostSpans(ctx, req)
}

// CollectorTransport reuses gRPC connection from agent to strategy manager.
type StrategyManagerTransport struct {
	logger *zap.Logger

	c    api_v1.StrategyManagerClient
	conn *grpc.ClientConn
	ep   *routing.Endpoint
}

func NewStrategyManagerTransport(logger *zap.Logger, ep *routing.Endpoint) *StrategyManagerTransport {
	ct := &StrategyManagerTransport{
		logger: logger,
		ep:     ep,
	}
	conn, err := grpc.Dial(ep.String(), grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Failed to dial to remote strategy manager", zap.Error(err))
		return nil
	}
	ct.conn = conn
	ct.c = api_v1.NewStrategyManagerClient(conn)
	return ct
}

func (t *StrategyManagerTransport) GetStrategy(ctx context.Context, req *api_v1.StrategyRequest) (*api_v1.StrategyResponse, error) {
	if t.conn.GetState() == connectivity.Shutdown {
		conn, err := grpc.Dial(t.ep.String(), grpc.WithInsecure())
		if err != nil {
			t.logger.Fatal("Connection is closed and client failed to dial to remote strategy manager", zap.Error(err))
			return &api_v1.StrategyResponse{}, err
		}
		t.conn = conn
		t.c = api_v1.NewStrategyManagerClient(conn)
	}
	return t.c.GetStrategy(ctx, req)
}
