// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.1
// source: proto/amqp_proxy_go.proto

package AmqpProxy

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConsumerClient is the client API for Consumer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsumerClient interface {
	Do(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Reply, error)
}

type consumerClient struct {
	cc grpc.ClientConnInterface
}

func NewConsumerClient(cc grpc.ClientConnInterface) ConsumerClient {
	return &consumerClient{cc}
}

func (c *consumerClient) Do(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/AmqpProxy.consumer/do", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsumerServer is the server API for Consumer service.
// All implementations should embed UnimplementedConsumerServer
// for forward compatibility
type ConsumerServer interface {
	Do(context.Context, *Params) (*Reply, error)
}

// UnimplementedConsumerServer should be embedded to have forward compatible implementations.
type UnimplementedConsumerServer struct {
}

func (UnimplementedConsumerServer) Do(context.Context, *Params) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Do not implemented")
}

// UnsafeConsumerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsumerServer will
// result in compilation errors.
type UnsafeConsumerServer interface {
	mustEmbedUnimplementedConsumerServer()
}

func RegisterConsumerServer(s grpc.ServiceRegistrar, srv ConsumerServer) {
	s.RegisterService(&Consumer_ServiceDesc, srv)
}

func _Consumer_Do_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Params)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsumerServer).Do(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AmqpProxy.consumer/do",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsumerServer).Do(ctx, req.(*Params))
	}
	return interceptor(ctx, in, info, handler)
}

// Consumer_ServiceDesc is the grpc.ServiceDesc for Consumer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Consumer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AmqpProxy.consumer",
	HandlerType: (*ConsumerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "do",
			Handler:    _Consumer_Do_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/amqp_proxy_go.proto",
}