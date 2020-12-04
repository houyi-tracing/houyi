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
	"github.com/magiconair/properties/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	httpPort = 8088
)

func TestHC(t *testing.T) {
	hc := NewHC(zap.NewNop())
	assert.Equal(t, hc.getState().status, Unavailable)
	hc.Set(Ready)
	assert.Equal(t, hc.getState().status, Ready)
	hc.Set(Broken)
	assert.Equal(t, hc.getState().status, Broken)

	hc.Ready()

	go func() {
		http.Handle("/health", hc.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(2 * time.Second)
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", httpPort))
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	hcResp := hcResponse{}
	_ = json.Unmarshal(bytes, &hcResp)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, hcResp.StatusMsg, "Service Available")
	_ = resp.Body.Close()
}
