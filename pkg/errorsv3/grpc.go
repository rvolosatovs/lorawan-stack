// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"context"
	"encoding/hex"

	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorDetails that can be carried over API.
type ErrorDetails interface {
	Namespace() string
	Name() string
	MessageFormat() string
	PublicAttributes() map[string]interface{}
	CorrelationID() string
}

// FromGRPCStatus converts the gRPC status message into an Error.
func FromGRPCStatus(status *status.Status) Error {
	err := build(Definition{
		code:          uint32(status.Code()),
		messageFormat: status.Message(),
	}, 0)
	if ErrorDetailsFromProto == nil {
		return err
	}
	detailMsgs := status.Details()
	detailProtos := make([]proto.Message, 0, len(detailMsgs))
	for _, msg := range detailMsgs { // convert to []proto.Message
		if msg, ok := msg.(proto.Message); ok {
			detailProtos = append(detailProtos, msg)
		}
	}
	details, rest := ErrorDetailsFromProto(detailProtos...)
	if len(rest) != 0 {
		detailIfaces := make([]interface{}, len(rest))
		for i, iface := range rest { // convert to []interface{}
			detailIfaces[i] = iface
		}
		err.details = detailIfaces
	}
	if details != nil {
		if namespace := details.Namespace(); namespace != "" {
			err.namespace = namespace
		}
		if name := details.Name(); name != "" {
			err.name = name
		}
		if messageFormat := details.MessageFormat(); messageFormat != "" {
			err.messageFormat = messageFormat
		}
		if attributes := details.PublicAttributes(); len(attributes) != 0 {
			err.attributes = attributes
		}
		if correlationID := details.CorrelationID(); correlationID != "" {
			err.correlationID = correlationID
		}
	}
	return err
}

// ErrorDetailsToProto converts the given ErrorDetails into a protobuf-encoded message.
//
// This variable is set by pkg/ttnpb.
var ErrorDetailsToProto func(e ErrorDetails) (msg proto.Message)

// ErrorDetailsFromProto ranges over the given protobuf-encoded messages
// to extract the ErrorDetails. It returns details if present, as well as the
// rest of the details.
//
// This variable is set by pkg/ttnpb.
var ErrorDetailsFromProto func(msg ...proto.Message) (details ErrorDetails, rest []proto.Message)

// setGRPCStatus sets a (marshaled) gRPC status in the error definition.
//
// This func should be called when the error definition is created. Doing that
// makes that we have to convert to a gRPC status only once instead of on every call.
func (d *Definition) setGRPCStatus() {
	s := status.New(codes.Code(d.Code()), d.String())
	if ErrorDetailsToProto != nil {
		if proto := ErrorDetailsToProto(d); proto != nil {
			var err error
			s, err = s.WithDetails(proto)
			if err != nil {
				panic(err) // ErrorDetailsToProto generated an invalid proto.
			}
		}
	}
	d.grpcStatus = s
}

// GRPCStatus returns the Definition as a gRPC status message.
func (d Definition) GRPCStatus() *status.Status {
	return d.grpcStatus // initialized when defined (with setGRPCStatus).
}

func (e *Error) clearGRPCStatus() {
	e.grpcStatus.Store((*status.Status)(nil))
}

// GRPCStatus converts the Error into a gRPC status message.
func (e *Error) GRPCStatus() *status.Status {
	if s, ok := e.grpcStatus.Load().(*status.Status); ok && s != nil {
		return s
	}
	s := status.New(codes.Code(e.Code()), e.String())
	protoDetails := make([]proto.Message, 0, len(e.Details())+1)
	for _, details := range e.Details() {
		if details, ok := details.(proto.Message); ok {
			protoDetails = append(protoDetails, details)
		}
	}
	if ErrorDetailsToProto != nil {
		if proto := ErrorDetailsToProto(e); proto != nil {
			protoDetails = append(protoDetails, proto)
		}
	}
	if len(protoDetails) != 0 {
		var err error
		s, err = s.WithDetails(protoDetails...)
		if err != nil {
			panic(err) // invalid details in the error or ErrorDetailsToProto generated an invalid proto.
		}
	}
	e.grpcStatus.Store(s)
	return s
}

// UnaryServerInterceptor makes sure that returned TTN errors contain a CorrelationID.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		res, err := handler(ctx, req)
		if ttnErr, ok := From(err); ok && ttnErr != nil {
			if ttnErr.correlationID == "" {
				ttnErr.correlationID = hex.EncodeToString(uuid.NewV4().Bytes()) // Compliant with Sentry.
			}
			err = ttnErr
		}
		return res, err
	}
}

// StreamServerInterceptor makes sure that returned TTN errors contain a CorrelationID.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, stream)
		if ttnErr, ok := From(err); ok && ttnErr != nil {
			if ttnErr.correlationID == "" {
				ttnErr.correlationID = hex.EncodeToString(uuid.NewV4().Bytes()) // Compliant with Sentry.
			}
			err = ttnErr
		}
		return err
	}
}

// UnaryClientInterceptor converts gRPC errors to regular errors.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err, ok := From(err); ok && err != nil {
			return err
		}
		return err
	}
}

type wrappedStream struct {
	grpc.ClientStream
}

func (w wrappedStream) SendMsg(m interface{}) error {
	err := w.ClientStream.SendMsg(m)
	if err, ok := From(err); ok && err != nil {
		return err
	}
	return err
}
func (w wrappedStream) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)
	if err, ok := From(err); ok && err != nil {
		return err
	}
	return err
}

// StreamClientInterceptor converts gRPC errors to regular errors.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err, ok := From(err); ok && err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		return wrappedStream{s}, nil
	}
}
