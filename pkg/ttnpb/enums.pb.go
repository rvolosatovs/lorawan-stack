// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/enums.proto

package ttnpb

import (
	fmt "fmt"
	math "math"
	strconv "strconv"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	golang_proto "github.com/golang/protobuf/proto"
)

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

type DownlinkPathConstraint int32

const (
	// Indicates that the gateway can be selected for downlink without constraints by the Network Server.
	DOWNLINK_PATH_CONSTRAINT_NONE DownlinkPathConstraint = 0
	// Indicates that the gateway can be selected for downlink only if no other or better gateway can be selected.
	DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER DownlinkPathConstraint = 1
	// Indicates that this gateway will never be selected for downlink, even if that results in no available downlink path.
	DOWNLINK_PATH_CONSTRAINT_NEVER DownlinkPathConstraint = 2
)

var DownlinkPathConstraint_name = map[int32]string{
	0: "DOWNLINK_PATH_CONSTRAINT_NONE",
	1: "DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER",
	2: "DOWNLINK_PATH_CONSTRAINT_NEVER",
}

var DownlinkPathConstraint_value = map[string]int32{
	"DOWNLINK_PATH_CONSTRAINT_NONE":         0,
	"DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER": 1,
	"DOWNLINK_PATH_CONSTRAINT_NEVER":        2,
}

func (DownlinkPathConstraint) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e36318a1e2f407cb, []int{0}
}

// State enum defines states that an entity can be in.
type State int32

const (
	// Denotes that the entity has been requested and is pending review by an admin.
	STATE_REQUESTED State = 0
	// Denotes that the entity has been reviewed and approved by an admin.
	STATE_APPROVED State = 1
	// Denotes that the entity has been reviewed and rejected by an admin.
	STATE_REJECTED State = 2
	// Denotes that the entity has been flagged and is pending review by an admin.
	STATE_FLAGGED State = 3
	// Denotes that the entity has been reviewed and suspended by an admin.
	STATE_SUSPENDED State = 4
)

var State_name = map[int32]string{
	0: "STATE_REQUESTED",
	1: "STATE_APPROVED",
	2: "STATE_REJECTED",
	3: "STATE_FLAGGED",
	4: "STATE_SUSPENDED",
}

var State_value = map[string]int32{
	"STATE_REQUESTED": 0,
	"STATE_APPROVED":  1,
	"STATE_REJECTED":  2,
	"STATE_FLAGGED":   3,
	"STATE_SUSPENDED": 4,
}

func (State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e36318a1e2f407cb, []int{1}
}

type ClusterRole int32

const (
	ClusterRole_NONE                      ClusterRole = 0
	ClusterRole_ENTITY_REGISTRY           ClusterRole = 1
	ClusterRole_ACCESS                    ClusterRole = 2
	ClusterRole_GATEWAY_SERVER            ClusterRole = 3
	ClusterRole_NETWORK_SERVER            ClusterRole = 4
	ClusterRole_APPLICATION_SERVER        ClusterRole = 5
	ClusterRole_JOIN_SERVER               ClusterRole = 6
	ClusterRole_CRYPTO_SERVER             ClusterRole = 7
	ClusterRole_DEVICE_TEMPLATE_CONVERTER ClusterRole = 8
	ClusterRole_DEVICE_CLAIMING_SERVER    ClusterRole = 9
)

var ClusterRole_name = map[int32]string{
	0: "NONE",
	1: "ENTITY_REGISTRY",
	2: "ACCESS",
	3: "GATEWAY_SERVER",
	4: "NETWORK_SERVER",
	5: "APPLICATION_SERVER",
	6: "JOIN_SERVER",
	7: "CRYPTO_SERVER",
	8: "DEVICE_TEMPLATE_CONVERTER",
	9: "DEVICE_CLAIMING_SERVER",
}

var ClusterRole_value = map[string]int32{
	"NONE":                      0,
	"ENTITY_REGISTRY":           1,
	"ACCESS":                    2,
	"GATEWAY_SERVER":            3,
	"NETWORK_SERVER":            4,
	"APPLICATION_SERVER":        5,
	"JOIN_SERVER":               6,
	"CRYPTO_SERVER":             7,
	"DEVICE_TEMPLATE_CONVERTER": 8,
	"DEVICE_CLAIMING_SERVER":    9,
}

func (ClusterRole) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e36318a1e2f407cb, []int{2}
}

func init() {
	proto.RegisterEnum("ttn.lorawan.v3.DownlinkPathConstraint", DownlinkPathConstraint_name, DownlinkPathConstraint_value)
	golang_proto.RegisterEnum("ttn.lorawan.v3.DownlinkPathConstraint", DownlinkPathConstraint_name, DownlinkPathConstraint_value)
	proto.RegisterEnum("ttn.lorawan.v3.State", State_name, State_value)
	golang_proto.RegisterEnum("ttn.lorawan.v3.State", State_name, State_value)
	proto.RegisterEnum("ttn.lorawan.v3.ClusterRole", ClusterRole_name, ClusterRole_value)
	golang_proto.RegisterEnum("ttn.lorawan.v3.ClusterRole", ClusterRole_name, ClusterRole_value)
}

func init() { proto.RegisterFile("lorawan-stack/api/enums.proto", fileDescriptor_e36318a1e2f407cb) }
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/enums.proto", fileDescriptor_e36318a1e2f407cb)
}

var fileDescriptor_e36318a1e2f407cb = []byte{
	// 561 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xbf, 0x4f, 0xdb, 0x4e,
	0x18, 0xc6, 0xef, 0x20, 0xf0, 0xe5, 0x7b, 0xa8, 0xe0, 0x1e, 0x12, 0x52, 0x91, 0x78, 0xa5, 0x56,
	0xea, 0x50, 0x54, 0x92, 0x81, 0xbf, 0xc0, 0xb5, 0x5f, 0x82, 0x21, 0xd8, 0xee, 0xf9, 0x08, 0x4a,
	0x17, 0xcb, 0x41, 0x69, 0x12, 0x25, 0xd8, 0x51, 0x72, 0x29, 0x2b, 0x23, 0x23, 0x63, 0xc7, 0xaa,
	0x5d, 0x18, 0x19, 0x19, 0x19, 0x19, 0xd9, 0xca, 0x88, 0xcf, 0x0b, 0x23, 0x23, 0x63, 0xe5, 0xfc,
	0x28, 0xea, 0xc0, 0x76, 0xef, 0xe7, 0x9e, 0xf7, 0xc7, 0x23, 0x3d, 0x6c, 0xbd, 0x9b, 0xf4, 0xa3,
	0x93, 0x28, 0xde, 0x1c, 0xa8, 0xe8, 0xa8, 0x53, 0x8a, 0x7a, 0xed, 0x52, 0x23, 0x1e, 0x1e, 0x0f,
	0x8a, 0xbd, 0x7e, 0xa2, 0x12, 0xbe, 0xa4, 0x54, 0x5c, 0x9c, 0x48, 0x8a, 0xdf, 0xb6, 0xd6, 0x36,
	0x9b, 0x6d, 0xd5, 0x1a, 0xd6, 0x8b, 0x47, 0xc9, 0x71, 0xa9, 0x99, 0x34, 0x93, 0xd2, 0x48, 0x56,
	0x1f, 0x7e, 0x1d, 0x55, 0xa3, 0x62, 0xf4, 0x1a, 0xb7, 0x6f, 0x9c, 0x53, 0xb6, 0x6a, 0x27, 0x27,
	0x71, 0xb7, 0x1d, 0x77, 0xfc, 0x48, 0xb5, 0xac, 0x24, 0x1e, 0xa8, 0x7e, 0xd4, 0x8e, 0x15, 0x7f,
	0xcb, 0xd6, 0x6d, 0xef, 0xd0, 0xad, 0x38, 0xee, 0x5e, 0xe8, 0x9b, 0x72, 0x27, 0xb4, 0x3c, 0x37,
	0x90, 0xc2, 0x74, 0x5c, 0x19, 0xba, 0x9e, 0x8b, 0x06, 0xe1, 0x1f, 0xd8, 0xfb, 0x17, 0x25, 0xbe,
	0xc0, 0x6d, 0x14, 0xa1, 0x27, 0x77, 0x50, 0x18, 0x94, 0xbf, 0x63, 0xf0, 0xf2, 0x34, 0xac, 0xa2,
	0x30, 0x66, 0xd6, 0x0a, 0x67, 0xbf, 0x80, 0x6c, 0xf4, 0xd9, 0x5c, 0xa0, 0x22, 0xd5, 0xe0, 0x2b,
	0x6c, 0x39, 0x90, 0xa6, 0xc4, 0x50, 0xe0, 0xe7, 0x03, 0x0c, 0x24, 0xda, 0x06, 0xe1, 0x9c, 0x2d,
	0x8d, 0xa1, 0xe9, 0xfb, 0xc2, 0xab, 0xa2, 0x6d, 0xd0, 0x67, 0x26, 0x70, 0x17, 0xad, 0x5c, 0x37,
	0xc3, 0x5f, 0xb3, 0x57, 0x63, 0xb6, 0x5d, 0x31, 0xcb, 0x65, 0xb4, 0x8d, 0xd9, 0xe7, 0x79, 0xc1,
	0x41, 0xe0, 0xa3, 0x6b, 0xa3, 0x6d, 0x14, 0x26, 0x3b, 0x7f, 0x53, 0xb6, 0x68, 0x75, 0x87, 0x03,
	0xd5, 0xe8, 0x8b, 0xa4, 0xdb, 0xe0, 0x0b, 0xac, 0x30, 0xb1, 0xb8, 0xc2, 0x96, 0xd1, 0x95, 0x8e,
	0xac, 0x85, 0x02, 0xcb, 0x4e, 0x20, 0x45, 0xcd, 0xa0, 0x9c, 0xb1, 0x79, 0xd3, 0xb2, 0x30, 0x08,
	0x8c, 0x99, 0x7c, 0x79, 0xd9, 0x94, 0x78, 0x68, 0xd6, 0xc2, 0x00, 0x45, 0x6e, 0x64, 0x36, 0x67,
	0x2e, 0xca, 0x43, 0x4f, 0xec, 0x4d, 0x59, 0x81, 0xaf, 0x32, 0x6e, 0xfa, 0x7e, 0xc5, 0xb1, 0x4c,
	0xe9, 0x78, 0xee, 0x94, 0xcf, 0xf1, 0x65, 0xb6, 0xb8, 0xeb, 0x39, 0x7f, 0xc1, 0x7c, 0x7e, 0xb9,
	0x25, 0x6a, 0xbe, 0xf4, 0xa6, 0xe8, 0x3f, 0xbe, 0xce, 0xde, 0xd8, 0x58, 0x75, 0x2c, 0x0c, 0x25,
	0xee, 0xfb, 0x95, 0xdc, 0x83, 0xe5, 0xb9, 0x55, 0x14, 0x12, 0x85, 0xb1, 0xc0, 0xd7, 0xd8, 0xea,
	0xe4, 0xdb, 0xaa, 0x98, 0xce, 0xbe, 0xe3, 0x96, 0xa7, 0xad, 0xff, 0x7f, 0xfa, 0x49, 0x6f, 0x52,
	0xa0, 0xb7, 0x29, 0xd0, 0xbb, 0x14, 0xc8, 0x7d, 0x0a, 0xe4, 0x21, 0x05, 0xf2, 0x98, 0x02, 0x79,
	0x4a, 0x81, 0x9e, 0x6a, 0xa0, 0x67, 0x1a, 0xc8, 0x85, 0x06, 0x7a, 0xa9, 0x81, 0x5c, 0x69, 0x20,
	0xd7, 0x1a, 0xc8, 0x8d, 0x06, 0x7a, 0xab, 0x81, 0xde, 0x69, 0x20, 0xf7, 0x1a, 0xe8, 0x83, 0x06,
	0xf2, 0xa8, 0x81, 0x3e, 0x69, 0x20, 0xa7, 0x19, 0x90, 0xb3, 0x0c, 0xe8, 0x79, 0x06, 0xe4, 0x7b,
	0x06, 0xf4, 0x47, 0x06, 0xe4, 0x22, 0x03, 0x72, 0x99, 0x01, 0xbd, 0xca, 0x80, 0x5e, 0x67, 0x40,
	0xbf, 0x7c, 0x6c, 0x26, 0x45, 0xd5, 0x6a, 0xa8, 0x56, 0x3b, 0x6e, 0x0e, 0x8a, 0x71, 0x43, 0x9d,
	0x24, 0xfd, 0x4e, 0xe9, 0xdf, 0x28, 0xf7, 0x3a, 0xcd, 0x92, 0x52, 0x71, 0xaf, 0x5e, 0x9f, 0x1f,
	0x85, 0x71, 0xeb, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xba, 0xf8, 0x09, 0x5e, 0xec, 0x02, 0x00,
	0x00,
}

func (x DownlinkPathConstraint) String() string {
	s, ok := DownlinkPathConstraint_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x State) String() string {
	s, ok := State_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x ClusterRole) String() string {
	s, ok := ClusterRole_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
