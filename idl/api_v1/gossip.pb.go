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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: gossip.proto

package api_v1

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Message_MessageType int32

const (
	Message_NEW_RELATION      Message_MessageType = 0
	Message_NEW_OPERATION     Message_MessageType = 1
	Message_EXPIRED_OPERATION Message_MessageType = 2
	Message_EVALUATING_TAGS   Message_MessageType = 3
)

// Enum value maps for Message_MessageType.
var (
	Message_MessageType_name = map[int32]string{
		0: "NEW_RELATION",
		1: "NEW_OPERATION",
		2: "EXPIRED_OPERATION",
		3: "EVALUATING_TAGS",
	}
	Message_MessageType_value = map[string]int32{
		"NEW_RELATION":      0,
		"NEW_OPERATION":     1,
		"EXPIRED_OPERATION": 2,
		"EVALUATING_TAGS":   3,
	}
)

func (x Message_MessageType) Enum() *Message_MessageType {
	p := new(Message_MessageType)
	*p = x
	return p
}

func (x Message_MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Message_MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_gossip_proto_enumTypes[0].Descriptor()
}

func (Message_MessageType) Type() protoreflect.EnumType {
	return &file_gossip_proto_enumTypes[0]
}

func (x Message_MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Message_MessageType.Descriptor instead.
func (Message_MessageType) EnumDescriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{1, 0}
}

type EvaluatingTags struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tags []*EvaluatingTag `protobuf:"bytes,1,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *EvaluatingTags) Reset() {
	*x = EvaluatingTags{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EvaluatingTags) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EvaluatingTags) ProtoMessage() {}

func (x *EvaluatingTags) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EvaluatingTags.ProtoReflect.Descriptor instead.
func (*EvaluatingTags) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{0}
}

func (x *EvaluatingTags) GetTags() []*EvaluatingTag {
	if x != nil {
		return x.Tags
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgId   int64               `protobuf:"varint,1,opt,name=msgId,proto3" json:"msgId,omitempty"`
	MsgType Message_MessageType `protobuf:"varint,2,opt,name=msgType,proto3,enum=gossip.Message_MessageType" json:"msgType,omitempty"`
	// Types that are assignable to Msg:
	//	*Message_Operation
	//	*Message_Relation
	//	*Message_EvaluateTags
	Msg isMessage_Msg `protobuf_oneof:"msg"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetMsgId() int64 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

func (x *Message) GetMsgType() Message_MessageType {
	if x != nil {
		return x.MsgType
	}
	return Message_NEW_RELATION
}

func (m *Message) GetMsg() isMessage_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *Message) GetOperation() *Operation {
	if x, ok := x.GetMsg().(*Message_Operation); ok {
		return x.Operation
	}
	return nil
}

func (x *Message) GetRelation() *Relation {
	if x, ok := x.GetMsg().(*Message_Relation); ok {
		return x.Relation
	}
	return nil
}

func (x *Message) GetEvaluateTags() *EvaluatingTags {
	if x, ok := x.GetMsg().(*Message_EvaluateTags); ok {
		return x.EvaluateTags
	}
	return nil
}

type isMessage_Msg interface {
	isMessage_Msg()
}

type Message_Operation struct {
	Operation *Operation `protobuf:"bytes,3,opt,name=operation,proto3,oneof"`
}

type Message_Relation struct {
	Relation *Relation `protobuf:"bytes,4,opt,name=relation,proto3,oneof"`
}

type Message_EvaluateTags struct {
	EvaluateTags *EvaluatingTags `protobuf:"bytes,5,opt,name=evaluateTags,proto3,oneof"`
}

func (*Message_Operation) isMessage_Msg() {}

func (*Message_Relation) isMessage_Msg() {}

func (*Message_EvaluateTags) isMessage_Msg() {}

type NullReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NullReply) Reset() {
	*x = NullReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NullReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NullReply) ProtoMessage() {}

func (x *NullReply) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NullReply.ProtoReflect.Descriptor instead.
func (*NullReply) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{2}
}

type Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port int64  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Peer) Reset() {
	*x = Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{3}
}

func (x *Peer) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Peer) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port int64 `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterRequest) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

type RegisterRely struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId     int64   `protobuf:"varint,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Interval   int64   `protobuf:"varint,2,opt,name=interval,proto3" json:"interval,omitempty"`
	RandomPick int64   `protobuf:"varint,3,opt,name=randomPick,proto3" json:"randomPick,omitempty"`
	ProbToR    float64 `protobuf:"fixed64,4,opt,name=probToR,proto3" json:"probToR,omitempty"`
}

func (x *RegisterRely) Reset() {
	*x = RegisterRely{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRely) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRely) ProtoMessage() {}

func (x *RegisterRely) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRely.ProtoReflect.Descriptor instead.
func (*RegisterRely) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{5}
}

func (x *RegisterRely) GetNodeId() int64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *RegisterRely) GetInterval() int64 {
	if x != nil {
		return x.Interval
	}
	return 0
}

func (x *RegisterRely) GetRandomPick() int64 {
	if x != nil {
		return x.RandomPick
	}
	return 0
}

func (x *RegisterRely) GetProbToR() float64 {
	if x != nil {
		return x.ProbToR
	}
	return 0
}

type HeartbeatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId int64 `protobuf:"varint,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Port   int64 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *HeartbeatRequest) Reset() {
	*x = HeartbeatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatRequest) ProtoMessage() {}

func (x *HeartbeatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatRequest.ProtoReflect.Descriptor instead.
func (*HeartbeatRequest) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{6}
}

func (x *HeartbeatRequest) GetNodeId() int64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *HeartbeatRequest) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

type HeartbeatReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId int64   `protobuf:"varint,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Peers  []*Peer `protobuf:"bytes,3,rep,name=peers,proto3" json:"peers,omitempty"`
}

func (x *HeartbeatReply) Reset() {
	*x = HeartbeatReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gossip_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatReply) ProtoMessage() {}

func (x *HeartbeatReply) ProtoReflect() protoreflect.Message {
	mi := &file_gossip_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatReply.ProtoReflect.Descriptor instead.
func (*HeartbeatReply) Descriptor() ([]byte, []int) {
	return file_gossip_proto_rawDescGZIP(), []int{7}
}

func (x *HeartbeatReply) GetNodeId() int64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *HeartbeatReply) GetPeers() []*Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

var File_gossip_proto protoreflect.FileDescriptor

var file_gossip_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x1a, 0x0b, 0x68, 0x6f, 0x75, 0x79, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x0e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6e,
	0x67, 0x54, 0x61, 0x67, 0x73, 0x12, 0x28, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x68, 0x6f, 0x75, 0x79, 0x69, 0x2e, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x22,
	0xdc, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d,
	0x73, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49,
	0x64, 0x12, 0x35, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x30, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x68, 0x6f,
	0x75, 0x79, 0x69, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52,
	0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x08, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x68,
	0x6f, 0x75, 0x79, 0x69, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52,
	0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3c, 0x0a, 0x0c, 0x65, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x65, 0x54, 0x61, 0x67, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x54, 0x61, 0x67, 0x73, 0x48, 0x00, 0x52, 0x0c, 0x65, 0x76, 0x61, 0x6c, 0x75,
	0x61, 0x74, 0x65, 0x54, 0x61, 0x67, 0x73, 0x22, 0x5e, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x45, 0x57, 0x5f, 0x52, 0x45,
	0x4c, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x45, 0x57, 0x5f,
	0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x45,
	0x58, 0x50, 0x49, 0x52, 0x45, 0x44, 0x5f, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x56, 0x41, 0x4c, 0x55, 0x41, 0x54, 0x49, 0x4e, 0x47,
	0x5f, 0x54, 0x41, 0x47, 0x53, 0x10, 0x03, 0x42, 0x05, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x0b,
	0x0a, 0x09, 0x4e, 0x75, 0x6c, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x2a, 0x0a, 0x04, 0x50,
	0x65, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x25, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x7c,
	0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x6c, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x50, 0x69, 0x63, 0x6b,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x50, 0x69,
	0x63, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x62, 0x54, 0x6f, 0x52, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x62, 0x54, 0x6f, 0x52, 0x22, 0x3e, 0x0a, 0x10,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x4c, 0x0a, 0x0e,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x50,
	0x65, 0x65, 0x72, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x32, 0x34, 0x0a, 0x04, 0x53, 0x65,
	0x65, 0x64, 0x12, 0x2c, 0x0a, 0x04, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x0f, 0x2e, 0x67, 0x6f, 0x73,
	0x73, 0x69, 0x70, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x11, 0x2e, 0x67, 0x6f,
	0x73, 0x73, 0x69, 0x70, 0x2e, 0x4e, 0x75, 0x6c, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x32, 0x88, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12, 0x3b, 0x0a,
	0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x17, 0x2e, 0x67, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x09, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x18, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70,
	0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74,
	0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x2b, 0x5a, 0x29, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x75, 0x79, 0x69, 0x2d,
	0x74, 0x72, 0x61, 0x63, 0x69, 0x6e, 0x67, 0x2f, 0x68, 0x6f, 0x75, 0x79, 0x69, 0x2f, 0x69, 0x64,
	0x6c, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gossip_proto_rawDescOnce sync.Once
	file_gossip_proto_rawDescData = file_gossip_proto_rawDesc
)

func file_gossip_proto_rawDescGZIP() []byte {
	file_gossip_proto_rawDescOnce.Do(func() {
		file_gossip_proto_rawDescData = protoimpl.X.CompressGZIP(file_gossip_proto_rawDescData)
	})
	return file_gossip_proto_rawDescData
}

var file_gossip_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gossip_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_gossip_proto_goTypes = []interface{}{
	(Message_MessageType)(0), // 0: gossip.Message.MessageType
	(*EvaluatingTags)(nil),   // 1: gossip.EvaluatingTags
	(*Message)(nil),          // 2: gossip.Message
	(*NullReply)(nil),        // 3: gossip.NullReply
	(*Peer)(nil),             // 4: gossip.Peer
	(*RegisterRequest)(nil),  // 5: gossip.RegisterRequest
	(*RegisterRely)(nil),     // 6: gossip.RegisterRely
	(*HeartbeatRequest)(nil), // 7: gossip.HeartbeatRequest
	(*HeartbeatReply)(nil),   // 8: gossip.HeartbeatReply
	(*EvaluatingTag)(nil),    // 9: houyi.EvaluatingTag
	(*Operation)(nil),        // 10: houyi.Operation
	(*Relation)(nil),         // 11: houyi.Relation
}
var file_gossip_proto_depIdxs = []int32{
	9,  // 0: gossip.EvaluatingTags.tags:type_name -> houyi.EvaluatingTag
	0,  // 1: gossip.Message.msgType:type_name -> gossip.Message.MessageType
	10, // 2: gossip.Message.operation:type_name -> houyi.Operation
	11, // 3: gossip.Message.relation:type_name -> houyi.Relation
	1,  // 4: gossip.Message.evaluateTags:type_name -> gossip.EvaluatingTags
	4,  // 5: gossip.HeartbeatReply.peers:type_name -> gossip.Peer
	2,  // 6: gossip.Seed.Sync:input_type -> gossip.Message
	5,  // 7: gossip.Registry.Register:input_type -> gossip.RegisterRequest
	7,  // 8: gossip.Registry.Heartbeat:input_type -> gossip.HeartbeatRequest
	3,  // 9: gossip.Seed.Sync:output_type -> gossip.NullReply
	6,  // 10: gossip.Registry.Register:output_type -> gossip.RegisterRely
	8,  // 11: gossip.Registry.Heartbeat:output_type -> gossip.HeartbeatReply
	9,  // [9:12] is the sub-list for method output_type
	6,  // [6:9] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_gossip_proto_init() }
func file_gossip_proto_init() {
	if File_gossip_proto != nil {
		return
	}
	file_houyi_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_gossip_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EvaluatingTags); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NullReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Peer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRely); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gossip_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_gossip_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Message_Operation)(nil),
		(*Message_Relation)(nil),
		(*Message_EvaluateTags)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gossip_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_gossip_proto_goTypes,
		DependencyIndexes: file_gossip_proto_depIdxs,
		EnumInfos:         file_gossip_proto_enumTypes,
		MessageInfos:      file_gossip_proto_msgTypes,
	}.Build()
	File_gossip_proto = out.File
	file_gossip_proto_rawDesc = nil
	file_gossip_proto_goTypes = nil
	file_gossip_proto_depIdxs = nil
}