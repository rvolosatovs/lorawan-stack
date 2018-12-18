// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: lorawan-stack/api/oauth_services.proto

/*
Package ttnpb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package ttnpb

import (
	"io"
	"net/http"

	"context"
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

var (
	filter_OAuthAuthorizationRegistry_List_0 = &utilities.DoubleArray{Encoding: map[string]int{"user_ids": 0, "user_id": 1}, Base: []int{1, 1, 1, 0}, Check: []int{0, 1, 2, 3}}
)

func request_OAuthAuthorizationRegistry_List_0(ctx context.Context, marshaler runtime.Marshaler, client OAuthAuthorizationRegistryClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListOAuthClientAuthorizationsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["user_ids.user_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "user_ids.user_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "user_ids.user_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "user_ids.user_id", err)
	}

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_OAuthAuthorizationRegistry_List_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.List(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_OAuthAuthorizationRegistry_ListTokens_0 = &utilities.DoubleArray{Encoding: map[string]int{"user_ids": 0, "user_id": 1, "client_ids": 2, "client_id": 3}, Base: []int{1, 1, 1, 1, 2, 0, 0}, Check: []int{0, 1, 2, 1, 4, 3, 5}}
)

func request_OAuthAuthorizationRegistry_ListTokens_0(ctx context.Context, marshaler runtime.Marshaler, client OAuthAuthorizationRegistryClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListOAuthAccessTokensRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["user_ids.user_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "user_ids.user_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "user_ids.user_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "user_ids.user_id", err)
	}

	val, ok = pathParams["client_ids.client_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "client_ids.client_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "client_ids.client_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "client_ids.client_id", err)
	}

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_OAuthAuthorizationRegistry_ListTokens_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.ListTokens(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_OAuthAuthorizationRegistry_Delete_0 = &utilities.DoubleArray{Encoding: map[string]int{"user_ids": 0, "user_id": 1, "client_ids": 2, "client_id": 3}, Base: []int{1, 1, 1, 1, 2, 0, 0}, Check: []int{0, 1, 2, 1, 4, 3, 5}}
)

func request_OAuthAuthorizationRegistry_Delete_0(ctx context.Context, marshaler runtime.Marshaler, client OAuthAuthorizationRegistryClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq OAuthClientAuthorizationIdentifiers
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["user_ids.user_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "user_ids.user_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "user_ids.user_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "user_ids.user_id", err)
	}

	val, ok = pathParams["client_ids.client_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "client_ids.client_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "client_ids.client_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "client_ids.client_id", err)
	}

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_OAuthAuthorizationRegistry_Delete_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.Delete(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_OAuthAuthorizationRegistry_DeleteToken_0 = &utilities.DoubleArray{Encoding: map[string]int{"user_ids": 0, "user_id": 1, "client_ids": 2, "client_id": 3, "id": 4}, Base: []int{1, 1, 1, 1, 2, 3, 0, 0, 0}, Check: []int{0, 1, 2, 1, 4, 1, 3, 5, 6}}
)

func request_OAuthAuthorizationRegistry_DeleteToken_0(ctx context.Context, marshaler runtime.Marshaler, client OAuthAuthorizationRegistryClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq OAuthAccessTokenIdentifiers
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["user_ids.user_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "user_ids.user_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "user_ids.user_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "user_ids.user_id", err)
	}

	val, ok = pathParams["client_ids.client_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "client_ids.client_id")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "client_ids.client_id", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "client_ids.client_id", err)
	}

	val, ok = pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.ID, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_OAuthAuthorizationRegistry_DeleteToken_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.DeleteToken(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterOAuthAuthorizationRegistryHandlerFromEndpoint is same as RegisterOAuthAuthorizationRegistryHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterOAuthAuthorizationRegistryHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
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

	return RegisterOAuthAuthorizationRegistryHandler(ctx, mux, conn)
}

// RegisterOAuthAuthorizationRegistryHandler registers the http handlers for service OAuthAuthorizationRegistry to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterOAuthAuthorizationRegistryHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterOAuthAuthorizationRegistryHandlerClient(ctx, mux, NewOAuthAuthorizationRegistryClient(conn))
}

// RegisterOAuthAuthorizationRegistryHandlerClient registers the http handlers for service OAuthAuthorizationRegistry
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "OAuthAuthorizationRegistryClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "OAuthAuthorizationRegistryClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "OAuthAuthorizationRegistryClient" to call the correct interceptors.
func RegisterOAuthAuthorizationRegistryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client OAuthAuthorizationRegistryClient) error {

	mux.Handle("GET", pattern_OAuthAuthorizationRegistry_List_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_OAuthAuthorizationRegistry_List_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_OAuthAuthorizationRegistry_List_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_OAuthAuthorizationRegistry_ListTokens_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_OAuthAuthorizationRegistry_ListTokens_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_OAuthAuthorizationRegistry_ListTokens_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("DELETE", pattern_OAuthAuthorizationRegistry_Delete_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_OAuthAuthorizationRegistry_Delete_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_OAuthAuthorizationRegistry_Delete_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("DELETE", pattern_OAuthAuthorizationRegistry_DeleteToken_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_OAuthAuthorizationRegistry_DeleteToken_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_OAuthAuthorizationRegistry_DeleteToken_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_OAuthAuthorizationRegistry_List_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1, 2, 2}, []string{"users", "user_ids.user_id", "authorizations"}, ""))

	pattern_OAuthAuthorizationRegistry_ListTokens_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1, 2, 2, 1, 0, 4, 1, 5, 3, 2, 4}, []string{"users", "user_ids.user_id", "authorizations", "client_ids.client_id", "tokens"}, ""))

	pattern_OAuthAuthorizationRegistry_Delete_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"users", "user_ids.user_id", "authorizations", "client_ids.client_id"}, ""))

	pattern_OAuthAuthorizationRegistry_DeleteToken_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1, 2, 2, 1, 0, 4, 1, 5, 3, 2, 4, 1, 0, 4, 1, 5, 5}, []string{"users", "user_ids.user_id", "authorizations", "client_ids.client_id", "tokens", "id"}, ""))
)

var (
	forward_OAuthAuthorizationRegistry_List_0 = runtime.ForwardResponseMessage

	forward_OAuthAuthorizationRegistry_ListTokens_0 = runtime.ForwardResponseMessage

	forward_OAuthAuthorizationRegistry_Delete_0 = runtime.ForwardResponseMessage

	forward_OAuthAuthorizationRegistry_DeleteToken_0 = runtime.ForwardResponseMessage
)
