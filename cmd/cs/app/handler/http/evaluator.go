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

package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/houyi-tracing/houyi/cmd/cs/app/handler/http/model"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/houyi-tracing/houyi/route"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

const (
	EqualTo              = "=="
	NotEqualTo           = "!="
	GreaterThan          = ">"
	GreaterThanOrEqualTo = ">="
	LessThan             = "<"
	LessThanOrEqualTo    = "<="
)

type EvaluatorHttpHandlerParams struct {
	Logger         *zap.Logger
	Evaluator      evaluator.Evaluator
	GossipRegistry gossip.Registry
}

type EvaluatorHttpHandler struct {
	logger   *zap.Logger
	eval     evaluator.Evaluator
	registry gossip.Registry
}

func NewEvaluatorHttpHandler(params *EvaluatorHttpHandlerParams) *EvaluatorHttpHandler {
	return &EvaluatorHttpHandler{
		logger:   params.Logger,
		eval:     params.Evaluator,
		registry: params.GossipRegistry,
	}
}

func (h *EvaluatorHttpHandler) RegisterRoutes(e *gin.Engine) {
	e.GET(route.GetEvaluatorTagsRoute, h.getEvaluatorTags)
	e.POST(route.UpdateEvaluatorTagsRoute, h.updateEvaluatorTags)
}

func (h *EvaluatorHttpHandler) getEvaluatorTags(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	tags := h.eval.Get()
	c.JSON(http.StatusOK, gin.H{
		"result": convertToJsonTags(tags.GetTags()),
	})
}

func (h *EvaluatorHttpHandler) updateEvaluatorTags(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	tags := make([]model.Tag, 0)
	err := c.BindJSON(&tags)
	if err == nil {
		convertedTags := convertToTags(tags)
		peers := h.registry.AllSeeds()
		h.eval.Update(&api_v1.EvaluatingTags{
			Tags: convertedTags,
		})
		for _, p := range peers {
			h.doUpdate(p.GetIp(), convertedTags)
		}
		c.JSON(http.StatusOK, gin.H{
			"result": "OK",
		})
	} else {
		h.logger.Error("failed to parse JSON from request's body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
		})
	}
}

func (h *EvaluatorHttpHandler) doUpdate(ip string, tags []*api_v1.EvaluatingTag) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, ports.CollectorGrpcListenPort), grpc.WithInsecure())
	if err != nil {
		h.logger.Debug("failed to dail collector", zap.String("ip", ip))
		return
	}
	c := api_v1.NewEvaluatorManagerClient(conn)
	_, err = c.UpdateTags(context.TODO(), &api_v1.UpdateTagsRequest{Tags: tags})
	if err != nil {
		h.logger.Error("failed to send request for updating tags", zap.Error(err))
	}
}

func convertToJsonTags(tags []*api_v1.EvaluatingTag) []model.Tag {
	ret := make([]model.Tag, 0)
	for _, t := range tags {
		newTag := model.Tag{}
		newTag.Name = t.TagName

		switch t.OperationType {
		case api_v1.EvaluatingTag_EQUAL_TO:
			newTag.Operator = EqualTo
		case api_v1.EvaluatingTag_NOT_EQUAL_TO:
			newTag.Operator = NotEqualTo
		case api_v1.EvaluatingTag_GREATER_THAN:
			newTag.Operator = GreaterThan
		case api_v1.EvaluatingTag_GREATER_THAN_OR_EQUAL_TO:
			newTag.Operator = GreaterThanOrEqualTo
		case api_v1.EvaluatingTag_LESS_THAN:
			newTag.Operator = LessThan
		case api_v1.EvaluatingTag_LESS_THAN_OR_EQUAL_TO:
			newTag.Operator = LessThanOrEqualTo
		}

		switch t.ValueType {
		case api_v1.EvaluatingTag_INTEGER:
			newTag.Value = t.GetIntegerVal()
		case api_v1.EvaluatingTag_STRING:
			newTag.Value = t.GetStringVal()
		case api_v1.EvaluatingTag_FLOAT:
			newTag.Value = t.GetFloatVal()
		case api_v1.EvaluatingTag_BOOLEAN:
			newTag.Value = t.GetBooleanVal()
		}

		ret = append(ret, newTag)
	}
	return ret
}

func convertToTags(tags []model.Tag) []*api_v1.EvaluatingTag {
	ret := make([]*api_v1.EvaluatingTag, 0)
	for _, t := range tags {
		newTag := &api_v1.EvaluatingTag{}
		newTag.TagName = t.Name

		switch t.Operator {
		case EqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_EQUAL_TO
		case NotEqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_NOT_EQUAL_TO
		case GreaterThan:
			newTag.OperationType = api_v1.EvaluatingTag_GREATER_THAN
		case GreaterThanOrEqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_GREATER_THAN_OR_EQUAL_TO
		case LessThan:
			newTag.OperationType = api_v1.EvaluatingTag_LESS_THAN
		case LessThanOrEqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_LESS_THAN_OR_EQUAL_TO
		}

		switch val := t.Value.(type) {
		case int64:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: val}
		case int32:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case int16:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case int8:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case int:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint64:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint32:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint16:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint8:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case float64:
			if IsInteger(val) {
				newTag.ValueType = api_v1.EvaluatingTag_INTEGER
				newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
			} else {
				newTag.ValueType = api_v1.EvaluatingTag_FLOAT
				newTag.Value = &api_v1.EvaluatingTag_FloatVal{FloatVal: val}
			}
		case float32:
			if Is32Integer(val) {
				newTag.ValueType = api_v1.EvaluatingTag_INTEGER
				newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
			} else {
				newTag.ValueType = api_v1.EvaluatingTag_FLOAT
				newTag.Value = &api_v1.EvaluatingTag_FloatVal{FloatVal: float64(val)}
			}
		case string:
			newTag.ValueType = api_v1.EvaluatingTag_STRING
			newTag.Value = &api_v1.EvaluatingTag_StringVal{StringVal: val}
		case bool:
			newTag.ValueType = api_v1.EvaluatingTag_BOOLEAN
			newTag.Value = &api_v1.EvaluatingTag_BooleanVal{BooleanVal: val}
		}

		ret = append(ret, newTag)
	}
	return ret
}

func IsInteger(n float64) bool {
	return n-float64(int(n)) == 0
}

func Is32Integer(n float32) bool {
	return n-float32(int(n)) == 0
}
