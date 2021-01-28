// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SeedClient is the client API for Seed service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeedClient interface {
	Sync(ctx context.Context, in *Message, opts ...grpc.CallOption) (*NullReply, error)
}

type seedClient struct {
	cc grpc.ClientConnInterface
}

func NewSeedClient(cc grpc.ClientConnInterface) SeedClient {
	return &seedClient{cc}
}

func (c *seedClient) Sync(ctx context.Context, in *Message, opts ...grpc.CallOption) (*NullReply, error) {
	out := new(NullReply)
	err := c.cc.Invoke(ctx, "/gossip.Seed/Sync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SeedServer is the server API for Seed service.
// All implementations must embed UnimplementedSeedServer
// for forward compatibility
type SeedServer interface {
	Sync(context.Context, *Message) (*NullReply, error)
	mustEmbedUnimplementedSeedServer()
}

// UnimplementedSeedServer must be embedded to have forward compatible implementations.
type UnimplementedSeedServer struct {
}

func (UnimplementedSeedServer) Sync(context.Context, *Message) (*NullReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (UnimplementedSeedServer) mustEmbedUnimplementedSeedServer() {}

// UnsafeSeedServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SeedServer will
// result in compilation errors.
type UnsafeSeedServer interface {
	mustEmbedUnimplementedSeedServer()
}

func RegisterSeedServer(s grpc.ServiceRegistrar, srv SeedServer) {
	s.RegisterService(&_Seed_serviceDesc, srv)
}

func _Seed_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeedServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gossip.Seed/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeedServer).Sync(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _Seed_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gossip.Seed",
	HandlerType: (*SeedServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sync",
			Handler:    _Seed_Sync_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gossip.proto",
}

// RegistryClient is the client API for GossipRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistryClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterRely, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatReply, error)
}

type registryClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistryClient(cc grpc.ClientConnInterface) RegistryClient {
	return &registryClient{cc}
}

func (c *registryClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterRely, error) {
	out := new(RegisterRely)
	err := c.cc.Invoke(ctx, "/gossip.GossipRegistry/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatReply, error) {
	out := new(HeartbeatReply)
	err := c.cc.Invoke(ctx, "/gossip.GossipRegistry/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistryServer is the server API for GossipRegistry service.
// All implementations must embed UnimplementedRegistryServer
// for forward compatibility
type RegistryServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterRely, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatReply, error)
	mustEmbedUnimplementedRegistryServer()
}

// UnimplementedRegistryServer must be embedded to have forward compatible implementations.
type UnimplementedRegistryServer struct {
}

func (UnimplementedRegistryServer) Register(context.Context, *RegisterRequest) (*RegisterRely, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRegistryServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedRegistryServer) mustEmbedUnimplementedRegistryServer() {}

// UnsafeRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistryServer will
// result in compilation errors.
type UnsafeRegistryServer interface {
	mustEmbedUnimplementedRegistryServer()
}

func RegisterRegistryServer(s grpc.ServiceRegistrar, srv RegistryServer) {
	s.RegisterService(&_Registry_serviceDesc, srv)
}

func _Registry_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gossip.GossipRegistry/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gossip.GossipRegistry/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Registry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gossip.GossipRegistry",
	HandlerType: (*RegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Registry_Register_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _Registry_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gossip.proto",
}
