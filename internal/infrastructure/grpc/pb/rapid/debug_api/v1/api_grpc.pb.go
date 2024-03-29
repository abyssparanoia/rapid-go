// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: rapid/debug_api/v1/api.proto

package debug_apiv1

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

const (
	DebugV1Service_CreateStaffIDToken_FullMethodName = "/rapid.debug_api.v1.DebugV1Service/CreateStaffIDToken"
)

// DebugV1ServiceClient is the client API for DebugV1Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DebugV1ServiceClient interface {
	CreateStaffIDToken(ctx context.Context, in *CreateStaffIDTokenRequest, opts ...grpc.CallOption) (*CreateStaffIDTokenResponse, error)
}

type debugV1ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDebugV1ServiceClient(cc grpc.ClientConnInterface) DebugV1ServiceClient {
	return &debugV1ServiceClient{cc}
}

func (c *debugV1ServiceClient) CreateStaffIDToken(ctx context.Context, in *CreateStaffIDTokenRequest, opts ...grpc.CallOption) (*CreateStaffIDTokenResponse, error) {
	out := new(CreateStaffIDTokenResponse)
	err := c.cc.Invoke(ctx, DebugV1Service_CreateStaffIDToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DebugV1ServiceServer is the server API for DebugV1Service service.
// All implementations should embed UnimplementedDebugV1ServiceServer
// for forward compatibility
type DebugV1ServiceServer interface {
	CreateStaffIDToken(context.Context, *CreateStaffIDTokenRequest) (*CreateStaffIDTokenResponse, error)
}

// UnimplementedDebugV1ServiceServer should be embedded to have forward compatible implementations.
type UnimplementedDebugV1ServiceServer struct {
}

func (UnimplementedDebugV1ServiceServer) CreateStaffIDToken(context.Context, *CreateStaffIDTokenRequest) (*CreateStaffIDTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStaffIDToken not implemented")
}

// UnsafeDebugV1ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DebugV1ServiceServer will
// result in compilation errors.
type UnsafeDebugV1ServiceServer interface {
	mustEmbedUnimplementedDebugV1ServiceServer()
}

func RegisterDebugV1ServiceServer(s grpc.ServiceRegistrar, srv DebugV1ServiceServer) {
	s.RegisterService(&DebugV1Service_ServiceDesc, srv)
}

func _DebugV1Service_CreateStaffIDToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStaffIDTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugV1ServiceServer).CreateStaffIDToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebugV1Service_CreateStaffIDToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugV1ServiceServer).CreateStaffIDToken(ctx, req.(*CreateStaffIDTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DebugV1Service_ServiceDesc is the grpc.ServiceDesc for DebugV1Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DebugV1Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rapid.debug_api.v1.DebugV1Service",
	HandlerType: (*DebugV1ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateStaffIDToken",
			Handler:    _DebugV1Service_CreateStaffIDToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rapid/debug_api/v1/api.proto",
}
