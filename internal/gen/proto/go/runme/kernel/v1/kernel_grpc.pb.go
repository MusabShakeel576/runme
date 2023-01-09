// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: runme/kernel/v1/kernel.proto

package kernelv1

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

// KernelServiceClient is the client API for KernelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KernelServiceClient interface {
	PostSession(ctx context.Context, in *PostSessionRequest, opts ...grpc.CallOption) (*PostSessionResponse, error)
	DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*DeleteSessionResponse, error)
	Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (KernelService_ExecuteClient, error)
}

type kernelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKernelServiceClient(cc grpc.ClientConnInterface) KernelServiceClient {
	return &kernelServiceClient{cc}
}

func (c *kernelServiceClient) PostSession(ctx context.Context, in *PostSessionRequest, opts ...grpc.CallOption) (*PostSessionResponse, error) {
	out := new(PostSessionResponse)
	err := c.cc.Invoke(ctx, "/runme.kernel.v1.KernelService/PostSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kernelServiceClient) DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*DeleteSessionResponse, error) {
	out := new(DeleteSessionResponse)
	err := c.cc.Invoke(ctx, "/runme.kernel.v1.KernelService/DeleteSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kernelServiceClient) Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (KernelService_ExecuteClient, error) {
	stream, err := c.cc.NewStream(ctx, &KernelService_ServiceDesc.Streams[0], "/runme.kernel.v1.KernelService/Execute", opts...)
	if err != nil {
		return nil, err
	}
	x := &kernelServiceExecuteClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KernelService_ExecuteClient interface {
	Recv() (*ExecuteResponse, error)
	grpc.ClientStream
}

type kernelServiceExecuteClient struct {
	grpc.ClientStream
}

func (x *kernelServiceExecuteClient) Recv() (*ExecuteResponse, error) {
	m := new(ExecuteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// KernelServiceServer is the server API for KernelService service.
// All implementations must embed UnimplementedKernelServiceServer
// for forward compatibility
type KernelServiceServer interface {
	PostSession(context.Context, *PostSessionRequest) (*PostSessionResponse, error)
	DeleteSession(context.Context, *DeleteSessionRequest) (*DeleteSessionResponse, error)
	Execute(*ExecuteRequest, KernelService_ExecuteServer) error
	mustEmbedUnimplementedKernelServiceServer()
}

// UnimplementedKernelServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKernelServiceServer struct {
}

func (UnimplementedKernelServiceServer) PostSession(context.Context, *PostSessionRequest) (*PostSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostSession not implemented")
}
func (UnimplementedKernelServiceServer) DeleteSession(context.Context, *DeleteSessionRequest) (*DeleteSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSession not implemented")
}
func (UnimplementedKernelServiceServer) Execute(*ExecuteRequest, KernelService_ExecuteServer) error {
	return status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedKernelServiceServer) mustEmbedUnimplementedKernelServiceServer() {}

// UnsafeKernelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KernelServiceServer will
// result in compilation errors.
type UnsafeKernelServiceServer interface {
	mustEmbedUnimplementedKernelServiceServer()
}

func RegisterKernelServiceServer(s grpc.ServiceRegistrar, srv KernelServiceServer) {
	s.RegisterService(&KernelService_ServiceDesc, srv)
}

func _KernelService_PostSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KernelServiceServer).PostSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runme.kernel.v1.KernelService/PostSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KernelServiceServer).PostSession(ctx, req.(*PostSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KernelService_DeleteSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KernelServiceServer).DeleteSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runme.kernel.v1.KernelService/DeleteSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KernelServiceServer).DeleteSession(ctx, req.(*DeleteSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KernelService_Execute_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecuteRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KernelServiceServer).Execute(m, &kernelServiceExecuteServer{stream})
}

type KernelService_ExecuteServer interface {
	Send(*ExecuteResponse) error
	grpc.ServerStream
}

type kernelServiceExecuteServer struct {
	grpc.ServerStream
}

func (x *kernelServiceExecuteServer) Send(m *ExecuteResponse) error {
	return x.ServerStream.SendMsg(m)
}

// KernelService_ServiceDesc is the grpc.ServiceDesc for KernelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KernelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "runme.kernel.v1.KernelService",
	HandlerType: (*KernelServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostSession",
			Handler:    _KernelService_PostSession_Handler,
		},
		{
			MethodName: "DeleteSession",
			Handler:    _KernelService_DeleteSession_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Execute",
			Handler:       _KernelService_Execute_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "runme/kernel/v1/kernel.proto",
}