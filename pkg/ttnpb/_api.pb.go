// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: go.thethings.network/lorawan-stack/api/_api.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

/*
The Things Network v3 API
*/

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func init() {
	proto.RegisterFile("go.thethings.network/lorawan-stack/api/_api.proto", fileDescriptor__api_d519cbdbd6083473)
}
func init() {
	golang_proto.RegisterFile("go.thethings.network/lorawan-stack/api/_api.proto", fileDescriptor__api_d519cbdbd6083473)
}

var fileDescriptor__api_d519cbdbd6083473 = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0xce, 0x21, 0x4e, 0x03, 0x41,
	0x14, 0x06, 0xe0, 0xf7, 0x0c, 0x02, 0x81, 0xe0, 0x00, 0xff, 0x09, 0x60, 0x26, 0xa4, 0x37, 0xe0,
	0x18, 0x18, 0x32, 0x25, 0xcd, 0x76, 0xb3, 0x64, 0x66, 0xb2, 0x7d, 0xa1, 0xb6, 0xb2, 0x12, 0x89,
	0x24, 0xa8, 0xca, 0xca, 0xca, 0xca, 0xca, 0x95, 0x2b, 0x77, 0xde, 0x33, 0x2b, 0x57, 0xae, 0x24,
	0x21, 0x1c, 0xa0, 0xfe, 0x13, 0xdf, 0xed, 0x53, 0x95, 0x9c, 0xac, 0x57, 0xb2, 0xae, 0x63, 0xb5,
	0x71, 0x71, 0x25, 0xdb, 0xd4, 0x36, 0xfe, 0x3d, 0xb5, 0x61, 0x1b, 0xe2, 0xe3, 0x46, 0xc2, 0x5b,
	0xe3, 0x43, 0xae, 0xfd, 0x6b, 0xc8, 0xb5, 0xcb, 0x6d, 0x92, 0x74, 0x7f, 0x27, 0x12, 0xdd, 0xbf,
	0x70, 0x1f, 0x8b, 0xe7, 0x1f, 0xbe, 0x14, 0x70, 0x57, 0xc0, 0x7d, 0x01, 0x0d, 0x05, 0x34, 0x16,
	0xd0, 0x54, 0x40, 0x73, 0x01, 0xef, 0x14, 0xbc, 0x57, 0xd0, 0x41, 0xc1, 0x47, 0x05, 0x9d, 0x14,
	0x74, 0x56, 0xd0, 0x45, 0xc1, 0x9d, 0x82, 0x7b, 0x05, 0x0d, 0x0a, 0x1e, 0x15, 0x34, 0x29, 0x78,
	0x56, 0xd0, 0xce, 0x40, 0x7b, 0x03, 0x7f, 0x1a, 0xe8, 0xcb, 0xc0, 0xdf, 0x06, 0x3a, 0x18, 0xe8,
	0x68, 0xe0, 0x93, 0x81, 0xcf, 0x06, 0x7e, 0x79, 0xb8, 0x62, 0x9e, 0x9b, 0xca, 0x8b, 0xc4, 0xbc,
	0x5c, 0xde, 0xfc, 0xdd, 0x17, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x43, 0x8a, 0xb2, 0x01, 0xf0,
	0x00, 0x00, 0x00,
}
