// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: dropper_service/dropper.proto

package dropper_service

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

// DropperServiceClient is the client API for DropperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DropperServiceClient interface {
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (DropperService_WatchClient, error)
}

type dropperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDropperServiceClient(cc grpc.ClientConnInterface) DropperServiceClient {
	return &dropperServiceClient{cc}
}

func (c *dropperServiceClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (DropperService_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &DropperService_ServiceDesc.Streams[0], "/dropperv1.DropperService/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &dropperServiceWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DropperService_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type dropperServiceWatchClient struct {
	grpc.ClientStream
}

func (x *dropperServiceWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DropperServiceServer is the server API for DropperService service.
// All implementations must embed UnimplementedDropperServiceServer
// for forward compatibility
type DropperServiceServer interface {
	Watch(*WatchRequest, DropperService_WatchServer) error
	mustEmbedUnimplementedDropperServiceServer()
}

// UnimplementedDropperServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDropperServiceServer struct {
}

func (UnimplementedDropperServiceServer) Watch(*WatchRequest, DropperService_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedDropperServiceServer) mustEmbedUnimplementedDropperServiceServer() {}

// UnsafeDropperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DropperServiceServer will
// result in compilation errors.
type UnsafeDropperServiceServer interface {
	mustEmbedUnimplementedDropperServiceServer()
}

func RegisterDropperServiceServer(s grpc.ServiceRegistrar, srv DropperServiceServer) {
	s.RegisterService(&DropperService_ServiceDesc, srv)
}

func _DropperService_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DropperServiceServer).Watch(m, &dropperServiceWatchServer{stream})
}

type DropperService_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type dropperServiceWatchServer struct {
	grpc.ServerStream
}

func (x *dropperServiceWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

// DropperService_ServiceDesc is the grpc.ServiceDesc for DropperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DropperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dropperv1.DropperService",
	HandlerType: (*DropperServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _DropperService_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "dropper_service/dropper.proto",
}
