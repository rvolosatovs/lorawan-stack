// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: lorawan-stack/api/applicationserver_integrations_storage.proto

/*
Package ttnpb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package ttnpb

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = descriptor.ForMessage
var _ = metadata.Join

var (
	filter_ApplicationUpStorage_GetStoredApplicationUp_0 = &utilities.DoubleArray{Encoding: map[string]int{"end_device_ids": 0, "application_ids": 1, "application_id": 2, "device_id": 3, "type": 4}, Base: []int{1, 1, 1, 1, 2, 3, 0, 0, 0}, Check: []int{0, 1, 2, 3, 2, 1, 4, 5, 6}}
)

func request_ApplicationUpStorage_GetStoredApplicationUp_0(ctx context.Context, marshaler runtime.Marshaler, client ApplicationUpStorageClient, req *http.Request, pathParams map[string]string) (ApplicationUpStorage_GetStoredApplicationUpClient, runtime.ServerMetadata, error) {
	var protoReq GetStoredApplicationUpRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["end_device_ids.application_ids.application_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "end_device_ids.application_ids.application_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "end_device_ids.application_ids.application_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "end_device_ids.application_ids.application_id", err)
	}

	val, ok = pathParams["end_device_ids.device_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "end_device_ids.device_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "end_device_ids.device_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "end_device_ids.device_id", err)
	}

	val, ok = pathParams["type"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "type")
	}

	protoReq.Type, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "type", err)
	}

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_ApplicationUpStorage_GetStoredApplicationUp_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	stream, err := client.GetStoredApplicationUp(ctx, &protoReq)
	if err != nil {
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil

}

var (
	filter_ApplicationUpStorage_GetStoredApplicationUp_1 = &utilities.DoubleArray{Encoding: map[string]int{"application_ids": 0, "application_id": 1, "type": 2}, Base: []int{1, 1, 1, 2, 0, 0}, Check: []int{0, 1, 2, 1, 3, 4}}
)

func request_ApplicationUpStorage_GetStoredApplicationUp_1(ctx context.Context, marshaler runtime.Marshaler, client ApplicationUpStorageClient, req *http.Request, pathParams map[string]string) (ApplicationUpStorage_GetStoredApplicationUpClient, runtime.ServerMetadata, error) {
	var protoReq GetStoredApplicationUpRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["application_ids.application_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "application_ids.application_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "application_ids.application_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "application_ids.application_id", err)
	}

	val, ok = pathParams["type"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "type")
	}

	protoReq.Type, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "type", err)
	}

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_ApplicationUpStorage_GetStoredApplicationUp_1); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	stream, err := client.GetStoredApplicationUp(ctx, &protoReq)
	if err != nil {
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil

}

// RegisterApplicationUpStorageHandlerServer registers the http handlers for service ApplicationUpStorage to "mux".
// UnaryRPC     :call ApplicationUpStorageServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterApplicationUpStorageHandlerFromEndpoint instead.
func RegisterApplicationUpStorageHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ApplicationUpStorageServer) error {

	mux.Handle("GET", pattern_ApplicationUpStorage_GetStoredApplicationUp_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	})

	mux.Handle("GET", pattern_ApplicationUpStorage_GetStoredApplicationUp_1, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	})

	return nil
}

// RegisterApplicationUpStorageHandlerFromEndpoint is same as RegisterApplicationUpStorageHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterApplicationUpStorageHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
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

	return RegisterApplicationUpStorageHandler(ctx, mux, conn)
}

// RegisterApplicationUpStorageHandler registers the http handlers for service ApplicationUpStorage to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterApplicationUpStorageHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterApplicationUpStorageHandlerClient(ctx, mux, NewApplicationUpStorageClient(conn))
}

// RegisterApplicationUpStorageHandlerClient registers the http handlers for service ApplicationUpStorage
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ApplicationUpStorageClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ApplicationUpStorageClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ApplicationUpStorageClient" to call the correct interceptors.
func RegisterApplicationUpStorageHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ApplicationUpStorageClient) error {

	mux.Handle("GET", pattern_ApplicationUpStorage_GetStoredApplicationUp_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ApplicationUpStorage_GetStoredApplicationUp_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ApplicationUpStorage_GetStoredApplicationUp_0(ctx, mux, outboundMarshaler, w, req, func() (proto.Message, error) { return resp.Recv() }, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ApplicationUpStorage_GetStoredApplicationUp_1, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ApplicationUpStorage_GetStoredApplicationUp_1(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ApplicationUpStorage_GetStoredApplicationUp_1(ctx, mux, outboundMarshaler, w, req, func() (proto.Message, error) { return resp.Recv() }, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ApplicationUpStorage_GetStoredApplicationUp_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2, 2, 3, 1, 0, 4, 1, 5, 4, 2, 5, 2, 6, 1, 0, 4, 1, 5, 7}, []string{"as", "applications", "end_device_ids.application_ids.application_id", "devices", "end_device_ids.device_id", "packages", "storage", "type"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ApplicationUpStorage_GetStoredApplicationUp_1 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2, 2, 3, 2, 4, 1, 0, 4, 1, 5, 5}, []string{"as", "applications", "application_ids.application_id", "packages", "storage", "type"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_ApplicationUpStorage_GetStoredApplicationUp_0 = runtime.ForwardResponseStream

	forward_ApplicationUpStorage_GetStoredApplicationUp_1 = runtime.ForwardResponseStream
)
