// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: lorawan-stack/api/identityserver.proto

/*
Package ttnpb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package ttnpb

import (
	"io"
	"net/http"

	"context"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

func request_EntityAccess_AuthInfo_0(ctx context.Context, marshaler runtime.Marshaler, client EntityAccessClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq types.Empty
	var metadata runtime.ServerMetadata

	msg, err := client.AuthInfo(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterEntityAccessHandlerFromEndpoint is same as RegisterEntityAccessHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterEntityAccessHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterEntityAccessHandler(ctx, mux, conn)
}

// RegisterEntityAccessHandler registers the http handlers for service EntityAccess to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterEntityAccessHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterEntityAccessHandlerClient(ctx, mux, NewEntityAccessClient(conn))
}

// RegisterEntityAccessHandlerClient registers the http handlers for service EntityAccess
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "EntityAccessClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "EntityAccessClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "EntityAccessClient" to call the correct interceptors.
func RegisterEntityAccessHandlerClient(ctx context.Context, mux *runtime.ServeMux, client EntityAccessClient) error {

	mux.Handle("GET", pattern_EntityAccess_AuthInfo_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_EntityAccess_AuthInfo_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_EntityAccess_AuthInfo_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_EntityAccess_AuthInfo_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"auth_info"}, ""))
)

var (
	forward_EntityAccess_AuthInfo_0 = runtime.ForwardResponseMessage
)
