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

package app

import (
	"context"
	"fmt"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestGetTagS(t *testing.T) {
	conn, err := grpc.Dial("192.168.31.77:18760", grpc.WithInsecure(), grpc.WithBlock())
	if conn == nil || err != nil {
		log.Fatal(err)
		return
	} else {
		defer conn.Close()
	}

	client := api_v1.NewEvaluatorManagerClient(conn)
	resp, err := client.GetTags(context.TODO(), &api_v1.GetTagsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

func TestUpdateTags(t *testing.T) {
	conn, err := grpc.Dial("192.168.31.77:18760", grpc.WithInsecure(), grpc.WithBlock())
	if conn == nil || err != nil {
		log.Fatal(err)
		return
	} else {
		defer conn.Close()
	}

	client := api_v1.NewEvaluatorManagerClient(conn)
	resp, err := client.UpdateTags(context.TODO(), &api_v1.UpdateTagsRequest{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       "http_status_code",
				OperationType: api_v1.EvaluatingTag_NOT_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_STRING,
				Value: &api_v1.EvaluatingTag_StringVal{
					StringVal: "500",
				},
			},
		}},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
