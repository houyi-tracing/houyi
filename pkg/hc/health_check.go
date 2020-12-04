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

package hc

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sync/atomic"
	"time"
)

type Status int

const (
	Unavailable Status = iota
	Ready
	Broken
)

func (s Status) String() string {
	switch s {
	case Unavailable:
		return "unavailable"
	case Ready:
		return "ready"
	case Broken:
		return "broken"
	default:
		return "unknown"
	}
}

type hcResponse struct {
	statusCode int
	StatusMsg  string    `json:"status"`
	UpSince    time.Time `json:"upSince"`
	Uptime     string    `json:"uptime"`
}

type state struct {
	status  Status
	upSince time.Time
}

type HealthCheck struct {
	logger  *zap.Logger
	state   atomic.Value
	respMap map[Status]hcResponse
}

func NewHC(logger *zap.Logger) *HealthCheck {
	hc := &HealthCheck{
		logger: logger,
		respMap: map[Status]hcResponse{
			Unavailable: {
				statusCode: http.StatusServiceUnavailable,
				StatusMsg:  "Service Unavailable",
			},
			Ready: {
				statusCode: http.StatusOK,
				StatusMsg:  "Service Available",
			},
		},
	}
	hc.state.Store(state{
		status: Unavailable,
	})
	return hc
}

func (hc *HealthCheck) Get() Status {
	return hc.getState().status
}

func (hc *HealthCheck) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := hc.getState()
		resp := hc.respMap[state.status]
		w.WriteHeader(resp.statusCode)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(hc.createRespBody(state, resp)); err != nil {
			hc.logger.Error("error for write http response body", zap.Error(err))
		}
	})
}

func (hc *HealthCheck) Ready() {
	hc.Set(Ready)
}

func (hc *HealthCheck) Set(status Status) {
	old := hc.getState()
	newState := state{
		status: status,
	}
	if status == Ready && old.status != Ready {
		newState.upSince = time.Now()
	}
	hc.state.Store(newState)
	hc.logger.Info("Health Check state changed", zap.Stringer("status", status))
}

func (hc *HealthCheck) SetLogger(logger *zap.Logger) {
	hc.logger = logger
}

func (hc *HealthCheck) getState() state {
	return hc.state.Load().(state)
}

func (hc *HealthCheck) createRespBody(state state, resp hcResponse) []byte {
	copyResp := resp
	if state.status == Ready {
		copyResp.UpSince = state.upSince
		copyResp.Uptime = fmt.Sprintf("%s", time.Since(state.upSince))
	}
	bytes, _ := json.Marshal(copyResp)
	return bytes
}
