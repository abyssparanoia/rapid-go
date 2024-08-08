// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: rapid/debug_api/v1/api.proto

/*
Package debug_apiv1 is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package debug_apiv1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_DebugV1Service_CreateStaffIDToken_0(ctx context.Context, marshaler runtime.Marshaler, client DebugV1ServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CreateStaffIDTokenRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.CreateStaffIDToken(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_DebugV1Service_CreateStaffIDToken_0(ctx context.Context, marshaler runtime.Marshaler, server DebugV1ServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq CreateStaffIDTokenRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.CreateStaffIDToken(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterDebugV1ServiceHandlerServer registers the http handlers for service DebugV1Service to "mux".
// UnaryRPC     :call DebugV1ServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterDebugV1ServiceHandlerFromEndpoint instead.
// GRPC interceptors will not work for this type of registration. To use interceptors, you must use the "runtime.WithMiddlewares" option in the "runtime.NewServeMux" call.
func RegisterDebugV1ServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server DebugV1ServiceServer) error {

	mux.Handle("POST", pattern_DebugV1Service_CreateStaffIDToken_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/rapid.debug_api.v1.DebugV1Service/CreateStaffIDToken", runtime.WithHTTPPathPattern("/debug/v1/staffs/-/id_token"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_DebugV1Service_CreateStaffIDToken_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_DebugV1Service_CreateStaffIDToken_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterDebugV1ServiceHandlerFromEndpoint is same as RegisterDebugV1ServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterDebugV1ServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterDebugV1ServiceHandler(ctx, mux, conn)
}

// RegisterDebugV1ServiceHandler registers the http handlers for service DebugV1Service to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterDebugV1ServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterDebugV1ServiceHandlerClient(ctx, mux, NewDebugV1ServiceClient(conn))
}

// RegisterDebugV1ServiceHandlerClient registers the http handlers for service DebugV1Service
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "DebugV1ServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "DebugV1ServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "DebugV1ServiceClient" to call the correct interceptors. This client ignores the HTTP middlewares.
func RegisterDebugV1ServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client DebugV1ServiceClient) error {

	mux.Handle("POST", pattern_DebugV1Service_CreateStaffIDToken_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/rapid.debug_api.v1.DebugV1Service/CreateStaffIDToken", runtime.WithHTTPPathPattern("/debug/v1/staffs/-/id_token"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_DebugV1Service_CreateStaffIDToken_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_DebugV1Service_CreateStaffIDToken_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_DebugV1Service_CreateStaffIDToken_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 2, 4}, []string{"debug", "v1", "staffs", "-", "id_token"}, ""))
)

var (
	forward_DebugV1Service_CreateStaffIDToken_0 = runtime.ForwardResponseMessage
)
