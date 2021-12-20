// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// DisplayServiceClient is the client API for DisplayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DisplayServiceClient interface {
	Show(ctx context.Context, in *ShowRequest, opts ...grpc.CallOption) (*ShowResponse, error)
}

type displayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDisplayServiceClient(cc grpc.ClientConnInterface) DisplayServiceClient {
	return &displayServiceClient{cc}
}

func (c *displayServiceClient) Show(ctx context.Context, in *ShowRequest, opts ...grpc.CallOption) (*ShowResponse, error) {
	out := new(ShowResponse)
	err := c.cc.Invoke(ctx, "/display.DisplayService/Show", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DisplayServiceServer is the server API for DisplayService service.
// All implementations should embed UnimplementedDisplayServiceServer
// for forward compatibility
type DisplayServiceServer interface {
	Show(context.Context, *ShowRequest) (*ShowResponse, error)
}

// UnimplementedDisplayServiceServer should be embedded to have forward compatible implementations.
type UnimplementedDisplayServiceServer struct {
}

func (UnimplementedDisplayServiceServer) Show(context.Context, *ShowRequest) (*ShowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Show not implemented")
}

// UnsafeDisplayServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DisplayServiceServer will
// result in compilation errors.
type UnsafeDisplayServiceServer interface {
	mustEmbedUnimplementedDisplayServiceServer()
}

func RegisterDisplayServiceServer(s grpc.ServiceRegistrar, srv DisplayServiceServer) {
	s.RegisterService(&_DisplayService_serviceDesc, srv)
}

func _DisplayService_Show_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisplayServiceServer).Show(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/display.DisplayService/Show",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisplayServiceServer).Show(ctx, req.(*ShowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DisplayService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "display.DisplayService",
	HandlerType: (*DisplayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Show",
			Handler:    _DisplayService_Show_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "display.proto",
}
