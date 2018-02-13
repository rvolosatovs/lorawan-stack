// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/messageprocessors.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UplinkMessageProcessor service

type UplinkMessageProcessorClient interface {
	Process(ctx context.Context, in *UplinkMessage, opts ...grpc.CallOption) (*UplinkMessage, error)
}

type uplinkMessageProcessorClient struct {
	cc *grpc.ClientConn
}

func NewUplinkMessageProcessorClient(cc *grpc.ClientConn) UplinkMessageProcessorClient {
	return &uplinkMessageProcessorClient{cc}
}

func (c *uplinkMessageProcessorClient) Process(ctx context.Context, in *UplinkMessage, opts ...grpc.CallOption) (*UplinkMessage, error) {
	out := new(UplinkMessage)
	err := grpc.Invoke(ctx, "/ttn.v3.UplinkMessageProcessor/Process", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UplinkMessageProcessor service

type UplinkMessageProcessorServer interface {
	Process(context.Context, *UplinkMessage) (*UplinkMessage, error)
}

func RegisterUplinkMessageProcessorServer(s *grpc.Server, srv UplinkMessageProcessorServer) {
	s.RegisterService(&_UplinkMessageProcessor_serviceDesc, srv)
}

func _UplinkMessageProcessor_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UplinkMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UplinkMessageProcessorServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.UplinkMessageProcessor/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UplinkMessageProcessorServer).Process(ctx, req.(*UplinkMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _UplinkMessageProcessor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.v3.UplinkMessageProcessor",
	HandlerType: (*UplinkMessageProcessorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _UplinkMessageProcessor_Process_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/TheThingsNetwork/ttn/api/messageprocessors.proto",
}

// Client API for DownlinkMessageProcessor service

type DownlinkMessageProcessorClient interface {
	Process(ctx context.Context, in *DownlinkMessage, opts ...grpc.CallOption) (*DownlinkMessage, error)
}

type downlinkMessageProcessorClient struct {
	cc *grpc.ClientConn
}

func NewDownlinkMessageProcessorClient(cc *grpc.ClientConn) DownlinkMessageProcessorClient {
	return &downlinkMessageProcessorClient{cc}
}

func (c *downlinkMessageProcessorClient) Process(ctx context.Context, in *DownlinkMessage, opts ...grpc.CallOption) (*DownlinkMessage, error) {
	out := new(DownlinkMessage)
	err := grpc.Invoke(ctx, "/ttn.v3.DownlinkMessageProcessor/Process", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DownlinkMessageProcessor service

type DownlinkMessageProcessorServer interface {
	Process(context.Context, *DownlinkMessage) (*DownlinkMessage, error)
}

func RegisterDownlinkMessageProcessorServer(s *grpc.Server, srv DownlinkMessageProcessorServer) {
	s.RegisterService(&_DownlinkMessageProcessor_serviceDesc, srv)
}

func _DownlinkMessageProcessor_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownlinkMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownlinkMessageProcessorServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.DownlinkMessageProcessor/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownlinkMessageProcessorServer).Process(ctx, req.(*DownlinkMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _DownlinkMessageProcessor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.v3.DownlinkMessageProcessor",
	HandlerType: (*DownlinkMessageProcessorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _DownlinkMessageProcessor_Process_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/TheThingsNetwork/ttn/api/messageprocessors.proto",
}

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/messageprocessors.proto", fileDescriptorMessageprocessors)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/messageprocessors.proto", fileDescriptorMessageprocessors)
}

var fileDescriptorMessageprocessors = []byte{
	// 285 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0x3d, 0x4c, 0x02, 0x41,
	0x10, 0x85, 0x67, 0x1a, 0x4c, 0xae, 0x24, 0xf1, 0x27, 0x14, 0xaf, 0xb7, 0xb9, 0x4b, 0xa0, 0xb0,
	0xa0, 0x33, 0xb6, 0x1a, 0x4d, 0x30, 0x26, 0x76, 0x40, 0x2e, 0xc7, 0xe5, 0xe4, 0x76, 0x73, 0xbb,
	0x4a, 0x4b, 0x49, 0x69, 0x69, 0xa7, 0x25, 0x25, 0x25, 0x25, 0x25, 0x25, 0x25, 0x25, 0x3b, 0xdb,
	0x50, 0x52, 0x52, 0x1a, 0x05, 0x62, 0x34, 0x18, 0x63, 0x35, 0x99, 0x37, 0xef, 0x7d, 0x93, 0xbc,
	0xa0, 0x9e, 0xa4, 0xb6, 0xf3, 0xd8, 0x0a, 0xdb, 0xaa, 0x1b, 0x35, 0x3a, 0x71, 0xa3, 0x93, 0xe6,
	0x89, 0xb9, 0x8a, 0x6d, 0x4f, 0x15, 0x59, 0x64, 0x6d, 0x1e, 0x35, 0x75, 0x1a, 0x75, 0x63, 0x63,
	0x9a, 0x49, 0xac, 0x0b, 0xd5, 0x8e, 0x8d, 0x51, 0x85, 0x09, 0x75, 0xa1, 0xac, 0x2a, 0x97, 0xac,
	0xcd, 0xc3, 0xa7, 0x5a, 0xa5, 0xfa, 0x0f, 0xc8, 0x36, 0x5b, 0xbd, 0x09, 0x8e, 0x6e, 0xf5, 0x43,
	0x9a, 0x67, 0x97, 0x1b, 0xfd, 0x7a, 0x07, 0x2f, 0x9f, 0x05, 0x07, 0xdb, 0xa5, 0x7c, 0x18, 0x6e,
	0x3e, 0x84, 0xdf, 0xac, 0x95, 0xfd, 0x72, 0xf5, 0x2e, 0x38, 0xb9, 0x50, 0xbd, 0x7c, 0x2f, 0xb4,
	0xfe, 0x05, 0x3d, 0xde, 0xa5, 0x7f, 0x98, 0x2b, 0xbf, 0x1d, 0xce, 0x5f, 0x79, 0xea, 0xc0, 0x33,
	0x07, 0x9e, 0x3b, 0xf0, 0xc2, 0x81, 0x97, 0x0e, 0xb4, 0x72, 0xa0, 0xb5, 0x03, 0xf7, 0x05, 0x34,
	0x10, 0xd0, 0x50, 0xc0, 0x23, 0x01, 0x8d, 0x05, 0x3c, 0x11, 0xf0, 0x54, 0xc0, 0x33, 0x01, 0xcf,
	0x05, 0xb4, 0x10, 0xf0, 0x52, 0x40, 0x2b, 0x01, 0xaf, 0x05, 0xd4, 0xf7, 0xa0, 0x81, 0x07, 0x3f,
	0x7b, 0xd0, 0x8b, 0x07, 0xbf, 0x79, 0xd0, 0xd0, 0x83, 0x46, 0x1e, 0x3c, 0xf6, 0xe0, 0x89, 0x07,
	0xdf, 0x9f, 0xfe, 0x55, 0xab, 0xce, 0x92, 0x8f, 0xa9, 0x5b, 0xad, 0xd2, 0x67, 0xa9, 0xb5, 0xf7,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x7a, 0x07, 0xd3, 0x2b, 0xcf, 0x01, 0x00, 0x00,
}
