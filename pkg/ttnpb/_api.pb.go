// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/_api.proto

/*
Package ttnpb is a generated protocol buffer package.

The Things Network v3 API

It is generated from these files:
	github.com/TheThingsNetwork/ttn/api/_api.proto
	github.com/TheThingsNetwork/ttn/api/application.proto
	github.com/TheThingsNetwork/ttn/api/end_device.proto
	github.com/TheThingsNetwork/ttn/api/gateway.proto
	github.com/TheThingsNetwork/ttn/api/gatewayserver.proto
	github.com/TheThingsNetwork/ttn/api/identifiers.proto
	github.com/TheThingsNetwork/ttn/api/join.proto
	github.com/TheThingsNetwork/ttn/api/joinserver.proto
	github.com/TheThingsNetwork/ttn/api/lorawan.proto
	github.com/TheThingsNetwork/ttn/api/messages.proto
	github.com/TheThingsNetwork/ttn/api/metadata.proto
	github.com/TheThingsNetwork/ttn/api/networkserver.proto

It has these top-level messages:
	ApplicationUplink
	ApplicationDownlink
	ApplicationDownlinks
	KeyEnvelope
	RootKeys
	SessionKeys
	Session
	EndDevice
	EndDevices
	MACSettings
	MACState
	MACInfo
	GatewayStatus
	GatewayObservations
	GatewayUp
	GatewayDown
	GatewayIdentifiers
	EndDeviceIdentifiers
	ApplicationIdentifiers
	JoinRequest
	JoinResponse
	Message
	MHDR
	MACPayload
	FHDR
	FCtrl
	JoinRequestPayload
	RejoinRequestPayload
	JoinAcceptPayload
	DLSettings
	CFList
	TxSettings
	UplinkMessage
	DownlinkMessage
	RxMetadata
	TxMetadata
	Location
	ApplicationUp
*/
package ttnpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func init() { proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/_api.proto", fileDescriptorXApi) }

var fileDescriptorXApi = []byte{
	// 128 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4b, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x0f, 0xc9, 0x48, 0x0d, 0xc9, 0xc8, 0xcc, 0x4b, 0x2f,
	0xf6, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x2f, 0x29, 0xc9, 0xd3, 0x4f, 0x2c, 0xc8, 0xd4,
	0x8f, 0x4f, 0x2c, 0xc8, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x29, 0xc9, 0xd3,
	0x2b, 0x33, 0x76, 0x72, 0xbf, 0xf1, 0x50, 0x8e, 0xa1, 0xe1, 0x91, 0x1c, 0xe3, 0x89, 0x47, 0x72,
	0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0xf8, 0xe2, 0x91, 0x1c, 0xc3, 0x87, 0x47,
	0x72, 0x8c, 0x51, 0x9a, 0x84, 0x4c, 0x2d, 0xc8, 0x4e, 0x07, 0xd1, 0x05, 0x49, 0x49, 0x6c, 0x60,
	0x73, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x36, 0xfc, 0x8c, 0x9c, 0x89, 0x00, 0x00, 0x00,
}
