// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: rapid/public_api/v1/api.proto

package public_apiv1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PublicV1Service_DeepHealthCheck_FullMethodName = "/rapid.public_api.v1.PublicV1Service/DeepHealthCheck"
)

// PublicV1ServiceClient is the client API for PublicV1Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PublicV1ServiceClient interface {
	DeepHealthCheck(ctx context.Context, in *DeepHealthCheckRequest, opts ...grpc.CallOption) (*DeepHealthCheckResponse, error)
}

type publicV1ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPublicV1ServiceClient(cc grpc.ClientConnInterface) PublicV1ServiceClient {
	return &publicV1ServiceClient{cc}
}

func (c *publicV1ServiceClient) DeepHealthCheck(ctx context.Context, in *DeepHealthCheckRequest, opts ...grpc.CallOption) (*DeepHealthCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeepHealthCheckResponse)
	err := c.cc.Invoke(ctx, PublicV1Service_DeepHealthCheck_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PublicV1ServiceServer is the server API for PublicV1Service service.
// All implementations should embed UnimplementedPublicV1ServiceServer
// for forward compatibility.
type PublicV1ServiceServer interface {
	DeepHealthCheck(context.Context, *DeepHealthCheckRequest) (*DeepHealthCheckResponse, error)
}

// UnimplementedPublicV1ServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPublicV1ServiceServer struct{}

func (UnimplementedPublicV1ServiceServer) DeepHealthCheck(context.Context, *DeepHealthCheckRequest) (*DeepHealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeepHealthCheck not implemented")
}
func (UnimplementedPublicV1ServiceServer) testEmbeddedByValue() {}

// UnsafePublicV1ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PublicV1ServiceServer will
// result in compilation errors.
type UnsafePublicV1ServiceServer interface {
	mustEmbedUnimplementedPublicV1ServiceServer()
}

func RegisterPublicV1ServiceServer(s grpc.ServiceRegistrar, srv PublicV1ServiceServer) {
	// If the following call pancis, it indicates UnimplementedPublicV1ServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PublicV1Service_ServiceDesc, srv)
}

func _PublicV1Service_DeepHealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeepHealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicV1ServiceServer).DeepHealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PublicV1Service_DeepHealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicV1ServiceServer).DeepHealthCheck(ctx, req.(*DeepHealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PublicV1Service_ServiceDesc is the grpc.ServiceDesc for PublicV1Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PublicV1Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rapid.public_api.v1.PublicV1Service",
	HandlerType: (*PublicV1ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeepHealthCheck",
			Handler:    _PublicV1Service_DeepHealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rapid/public_api/v1/api.proto",
}
