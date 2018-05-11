// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/rights.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strconv "strconv"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Right is the enum that defines all the different rights to do something in the network.
type Right int32

const (
	// This is the invalid right and should not be used. It can denote a parsing error.
	RIGHT_INVALID Right = 0
	// The right to view information about the user.
	RIGHT_USER_INFO Right = 1
	// The right to write the basic settings of the user.
	RIGHT_USER_SETTINGS_BASIC Right = 2
	// The right to delete the account of the user.
	RIGHT_USER_DELETE Right = 3
	// The right to view authorized, authorize and de-authorize OAuth clients of the user.
	RIGHT_USER_AUTHORIZED_CLIENTS Right = 4
	// The right to list all of the applications the user is a collaborator of.
	RIGHT_USER_APPLICATIONS_LIST Right = 5
	// The right to create an application under the account of the user.
	RIGHT_USER_APPLICATIONS_CREATE Right = 6
	// The right to list gateways the user is collaborator of.
	RIGHT_USER_GATEWAYS_LIST Right = 7
	// The right to register a gateway under the account of the user.
	RIGHT_USER_GATEWAYS_CREATE Right = 8
	// The right to create, list, view, update and delete OAuth clients under the account of the user.
	RIGHT_USER_CLIENTS Right = 9
	// The right to performs actions that requires the user to be an admin.
	RIGHT_USER_ADMIN Right = 10
	// The right to view and write the API keys of the user.
	RIGHT_USER_SETTINGS_API_KEYS Right = 11
	// The right to create a new organization under the account of the user.
	RIGHT_USER_ORGANIZATIONS_CREATE Right = 12
	// The right to list all of the organizations the user is a member of.
	RIGHT_USER_ORGANIZATIONS_LIST Right = 13
	// The right to view information about the application.
	RIGHT_APPLICATION_INFO Right = 14
	// The right to write the basic settings of the application.
	RIGHT_APPLICATION_SETTINGS_BASIC Right = 15
	// The right to view and write the API keys of the application.
	RIGHT_APPLICATION_SETTINGS_API_KEYS Right = 16
	// The right to view and edit the collaborators of the application.
	RIGHT_APPLICATION_SETTINGS_COLLABORATORS Right = 17
	// The right to delete the application.
	RIGHT_APPLICATION_DELETE Right = 18
	// The right to view the devices of an application.
	RIGHT_APPLICATION_DEVICES_READ Right = 19
	// The right to register devices to an application.
	RIGHT_APPLICATION_DEVICES_WRITE Right = 20
	// The right to read traffic of the application (uplink and downlink).
	RIGHT_APPLICATION_TRAFFIC_READ Right = 21
	// The right to write uplink traffic of the application.
	RIGHT_APPLICATION_TRAFFIC_UP_WRITE Right = 22
	// The right to write downlink traffic of the application.
	RIGHT_APPLICATION_TRAFFIC_DOWN_WRITE Right = 23
	// The right to link to a Network Server for traffic exchange, i.e. write uplink
	// and read downlink (API keys only).
	RIGHT_APPLICATION_LINK Right = 24
	// The right to view information about the gateway.
	RIGHT_GATEWAY_INFO Right = 25
	// The right to write the basic settings of the gateway.
	RIGHT_GATEWAY_SETTINGS_BASIC Right = 26
	// The right to view and write the API keys of the gateway.
	RIGHT_GATEWAY_SETTINGS_API_KEYS Right = 27
	// The right to view and edit collaborators of the gateway.
	RIGHT_GATEWAY_SETTINGS_COLLABORATORS Right = 28
	// The right to delete the gateway.
	RIGHT_GATEWAY_DELETE Right = 29
	// The right to read traffic coming from the gateway.
	RIGHT_GATEWAY_TRAFFIC_READ Right = 30
	// The right to write downlink traffic to a gateway.
	RIGHT_GATEWAY_TRAFFIC_DOWN_WRITE Right = 31
	// The right to link to a Gateway Server for traffic exchange, i.e. write uplink
	// and read downlink (API keys only).
	RIGHT_GATEWAY_LINK Right = 32
	// The right to view the status of the gateway.
	RIGHT_GATEWAY_STATUS_READ Right = 33
	// The right to view the location of the gateway.
	RIGHT_GATEWAY_LOCATION_READ Right = 34
	// The right to view information about the organization.
	RIGHT_ORGANIZATION_INFO Right = 35
	// The right to write the basic settings of the organization.
	RIGHT_ORGANIZATION_SETTINGS_BASIC Right = 36
	// The right to view and write the API keys of the organization.
	RIGHT_ORGANIZATION_SETTINGS_API_KEYS Right = 37
	// The right to view and edit members of the organization.
	RIGHT_ORGANIZATION_SETTINGS_MEMBERS Right = 38
	// The right to delete the organization.
	RIGHT_ORGANIZATION_DELETE Right = 39
	// The right to create an application under an organization.
	RIGHT_ORGANIZATION_APPLICATIONS_CREATE Right = 40
	// The right to list the applications the organization is collaborator of.
	RIGHT_ORGANIZATION_APPLICATIONS_LIST Right = 41
	// The right to create a gateway under an organization.
	RIGHT_ORGANIZATION_GATEWAYS_CREATE Right = 42
	// The right to list the gateways the organization is collaborator of.
	RIGHT_ORGANIZATION_GATEWAYS_LIST Right = 43
)

var Right_name = map[int32]string{
	0:  "RIGHT_INVALID",
	1:  "RIGHT_USER_INFO",
	2:  "RIGHT_USER_SETTINGS_BASIC",
	3:  "RIGHT_USER_DELETE",
	4:  "RIGHT_USER_AUTHORIZED_CLIENTS",
	5:  "RIGHT_USER_APPLICATIONS_LIST",
	6:  "RIGHT_USER_APPLICATIONS_CREATE",
	7:  "RIGHT_USER_GATEWAYS_LIST",
	8:  "RIGHT_USER_GATEWAYS_CREATE",
	9:  "RIGHT_USER_CLIENTS",
	10: "RIGHT_USER_ADMIN",
	11: "RIGHT_USER_SETTINGS_API_KEYS",
	12: "RIGHT_USER_ORGANIZATIONS_CREATE",
	13: "RIGHT_USER_ORGANIZATIONS_LIST",
	14: "RIGHT_APPLICATION_INFO",
	15: "RIGHT_APPLICATION_SETTINGS_BASIC",
	16: "RIGHT_APPLICATION_SETTINGS_API_KEYS",
	17: "RIGHT_APPLICATION_SETTINGS_COLLABORATORS",
	18: "RIGHT_APPLICATION_DELETE",
	19: "RIGHT_APPLICATION_DEVICES_READ",
	20: "RIGHT_APPLICATION_DEVICES_WRITE",
	21: "RIGHT_APPLICATION_TRAFFIC_READ",
	22: "RIGHT_APPLICATION_TRAFFIC_UP_WRITE",
	23: "RIGHT_APPLICATION_TRAFFIC_DOWN_WRITE",
	24: "RIGHT_APPLICATION_LINK",
	25: "RIGHT_GATEWAY_INFO",
	26: "RIGHT_GATEWAY_SETTINGS_BASIC",
	27: "RIGHT_GATEWAY_SETTINGS_API_KEYS",
	28: "RIGHT_GATEWAY_SETTINGS_COLLABORATORS",
	29: "RIGHT_GATEWAY_DELETE",
	30: "RIGHT_GATEWAY_TRAFFIC_READ",
	31: "RIGHT_GATEWAY_TRAFFIC_DOWN_WRITE",
	32: "RIGHT_GATEWAY_LINK",
	33: "RIGHT_GATEWAY_STATUS_READ",
	34: "RIGHT_GATEWAY_LOCATION_READ",
	35: "RIGHT_ORGANIZATION_INFO",
	36: "RIGHT_ORGANIZATION_SETTINGS_BASIC",
	37: "RIGHT_ORGANIZATION_SETTINGS_API_KEYS",
	38: "RIGHT_ORGANIZATION_SETTINGS_MEMBERS",
	39: "RIGHT_ORGANIZATION_DELETE",
	40: "RIGHT_ORGANIZATION_APPLICATIONS_CREATE",
	41: "RIGHT_ORGANIZATION_APPLICATIONS_LIST",
	42: "RIGHT_ORGANIZATION_GATEWAYS_CREATE",
	43: "RIGHT_ORGANIZATION_GATEWAYS_LIST",
}
var Right_value = map[string]int32{
	"RIGHT_INVALID":                            0,
	"RIGHT_USER_INFO":                          1,
	"RIGHT_USER_SETTINGS_BASIC":                2,
	"RIGHT_USER_DELETE":                        3,
	"RIGHT_USER_AUTHORIZED_CLIENTS":            4,
	"RIGHT_USER_APPLICATIONS_LIST":             5,
	"RIGHT_USER_APPLICATIONS_CREATE":           6,
	"RIGHT_USER_GATEWAYS_LIST":                 7,
	"RIGHT_USER_GATEWAYS_CREATE":               8,
	"RIGHT_USER_CLIENTS":                       9,
	"RIGHT_USER_ADMIN":                         10,
	"RIGHT_USER_SETTINGS_API_KEYS":             11,
	"RIGHT_USER_ORGANIZATIONS_CREATE":          12,
	"RIGHT_USER_ORGANIZATIONS_LIST":            13,
	"RIGHT_APPLICATION_INFO":                   14,
	"RIGHT_APPLICATION_SETTINGS_BASIC":         15,
	"RIGHT_APPLICATION_SETTINGS_API_KEYS":      16,
	"RIGHT_APPLICATION_SETTINGS_COLLABORATORS": 17,
	"RIGHT_APPLICATION_DELETE":                 18,
	"RIGHT_APPLICATION_DEVICES_READ":           19,
	"RIGHT_APPLICATION_DEVICES_WRITE":          20,
	"RIGHT_APPLICATION_TRAFFIC_READ":           21,
	"RIGHT_APPLICATION_TRAFFIC_UP_WRITE":       22,
	"RIGHT_APPLICATION_TRAFFIC_DOWN_WRITE":     23,
	"RIGHT_APPLICATION_LINK":                   24,
	"RIGHT_GATEWAY_INFO":                       25,
	"RIGHT_GATEWAY_SETTINGS_BASIC":             26,
	"RIGHT_GATEWAY_SETTINGS_API_KEYS":          27,
	"RIGHT_GATEWAY_SETTINGS_COLLABORATORS":     28,
	"RIGHT_GATEWAY_DELETE":                     29,
	"RIGHT_GATEWAY_TRAFFIC_READ":               30,
	"RIGHT_GATEWAY_TRAFFIC_DOWN_WRITE":         31,
	"RIGHT_GATEWAY_LINK":                       32,
	"RIGHT_GATEWAY_STATUS_READ":                33,
	"RIGHT_GATEWAY_LOCATION_READ":              34,
	"RIGHT_ORGANIZATION_INFO":                  35,
	"RIGHT_ORGANIZATION_SETTINGS_BASIC":        36,
	"RIGHT_ORGANIZATION_SETTINGS_API_KEYS":     37,
	"RIGHT_ORGANIZATION_SETTINGS_MEMBERS":      38,
	"RIGHT_ORGANIZATION_DELETE":                39,
	"RIGHT_ORGANIZATION_APPLICATIONS_CREATE":   40,
	"RIGHT_ORGANIZATION_APPLICATIONS_LIST":     41,
	"RIGHT_ORGANIZATION_GATEWAYS_CREATE":       42,
	"RIGHT_ORGANIZATION_GATEWAYS_LIST":         43,
}

func (Right) EnumDescriptor() ([]byte, []int) { return fileDescriptorRights, []int{0} }

func init() {
	proto.RegisterEnum("ttn.lorawan.v3.Right", Right_name, Right_value)
	golang_proto.RegisterEnum("ttn.lorawan.v3.Right", Right_name, Right_value)
}
func (x Right) String() string {
	s, ok := Right_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/rights.proto", fileDescriptorRights)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/rights.proto", fileDescriptorRights)
}

var fileDescriptorRights = []byte{
	// 740 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x95, 0xbb, 0x53, 0xeb, 0x46,
	0x14, 0xc6, 0xb5, 0x09, 0x90, 0x64, 0x09, 0xb0, 0x2c, 0x6f, 0x03, 0x87, 0xf7, 0x33, 0x89, 0x9d,
	0x19, 0xfe, 0x02, 0xd9, 0x5e, 0xcc, 0x0e, 0x42, 0xf2, 0x48, 0x6b, 0x18, 0x68, 0x3c, 0x38, 0x43,
	0x6c, 0x0f, 0x89, 0xed, 0x31, 0x22, 0xb4, 0x94, 0x94, 0x29, 0xd3, 0xe5, 0xce, 0xdc, 0x86, 0x92,
	0x92, 0x92, 0x92, 0x92, 0xee, 0x52, 0x62, 0xa9, 0xa1, 0xa4, 0xa4, 0xbc, 0x83, 0x24, 0xac, 0x87,
	0x65, 0x6e, 0x65, 0x6b, 0xbf, 0xdf, 0x1e, 0x9d, 0xfd, 0xbe, 0x33, 0x2b, 0xfc, 0x7b, 0xb9, 0x6a,
	0x56, 0xce, 0x4b, 0xc9, 0x3f, 0xea, 0x7f, 0xa7, 0x44, 0xe5, 0x44, 0x54, 0xaa, 0xb5, 0xf2, 0x99,
	0x7a, 0x62, 0x5e, 0xd4, 0x9b, 0xa7, 0x29, 0xd3, 0xac, 0xa5, 0x8e, 0x1b, 0xd5, 0x54, 0xb3, 0x5a,
	0xae, 0x98, 0x67, 0xc9, 0x46, 0xb3, 0x6e, 0xd6, 0xe9, 0xa0, 0x69, 0xd6, 0x92, 0x7f, 0xd5, 0x9b,
	0xc7, 0x17, 0xc7, 0xb5, 0xe4, 0x3f, 0x5b, 0x89, 0xdf, 0x02, 0x15, 0xca, 0xf5, 0x72, 0x3d, 0xe5,
	0x60, 0xa5, 0xf3, 0x3f, 0x9d, 0x27, 0xe7, 0xc1, 0xf9, 0xe7, 0x6e, 0xdf, 0xfc, 0xd2, 0x8f, 0x7b,
	0xf5, 0xb7, 0x7a, 0x74, 0x18, 0x0f, 0xe8, 0x3c, 0xb7, 0x23, 0x8a, 0x5c, 0xdd, 0x97, 0x15, 0x9e,
	0x25, 0x12, 0x1d, 0xc1, 0x43, 0xee, 0x52, 0xc1, 0x60, 0x7a, 0x91, 0xab, 0xdb, 0x1a, 0x41, 0x74,
	0x16, 0x4f, 0x05, 0x16, 0x0d, 0x26, 0x04, 0x57, 0x73, 0x46, 0x31, 0x2d, 0x1b, 0x3c, 0x43, 0xbe,
	0xa3, 0x63, 0x78, 0x38, 0x20, 0x67, 0x99, 0xc2, 0x04, 0x23, 0xdf, 0xd3, 0x05, 0x3c, 0x1b, 0x58,
	0x96, 0x0b, 0x62, 0x47, 0xd3, 0xf9, 0x11, 0xcb, 0x16, 0x33, 0x0a, 0x67, 0xaa, 0x30, 0x48, 0x0f,
	0x9d, 0xc7, 0x33, 0x41, 0x24, 0x9f, 0x57, 0x78, 0x46, 0x16, 0x5c, 0x53, 0x8d, 0xa2, 0xc2, 0x0d,
	0x41, 0x7a, 0xe9, 0x22, 0x86, 0x6e, 0x44, 0x46, 0x67, 0xb2, 0x60, 0xa4, 0x8f, 0xce, 0xe0, 0xc9,
	0x00, 0x93, 0x93, 0x05, 0x3b, 0x90, 0x0f, 0xbd, 0x0a, 0x3f, 0x50, 0xc0, 0x89, 0x38, 0xd5, 0xdb,
	0xfd, 0x23, 0x1d, 0xc7, 0x34, 0xa0, 0xbf, 0xf7, 0xf6, 0x13, 0x1d, 0xc5, 0x24, 0xf8, 0xe6, 0xec,
	0x1e, 0x57, 0x09, 0x8e, 0x74, 0xdc, 0xb6, 0x42, 0xce, 0xf3, 0xe2, 0x2e, 0x3b, 0x34, 0x48, 0x3f,
	0x5d, 0xc2, 0x73, 0x01, 0x42, 0xd3, 0x73, 0xb2, 0xca, 0x8f, 0xc2, 0x2d, 0xff, 0x1c, 0xf1, 0x26,
	0x0c, 0x39, 0x7d, 0x0f, 0xd0, 0x04, 0x1e, 0x77, 0x91, 0xc0, 0xa1, 0xdd, 0x40, 0x06, 0xe9, 0x32,
	0x9e, 0xef, 0xd4, 0x22, 0xb9, 0x0c, 0xd1, 0x35, 0xbc, 0xf4, 0x01, 0xd5, 0x6e, 0x99, 0xd0, 0x5f,
	0xf1, 0xfa, 0x07, 0x60, 0x46, 0x53, 0x14, 0x39, 0xad, 0xe9, 0xb2, 0xd0, 0x74, 0x83, 0x0c, 0xfb,
	0x76, 0x07, 0x69, 0x2f, 0x75, 0xea, 0x07, 0x16, 0x56, 0xf7, 0x79, 0x86, 0x19, 0x45, 0x9d, 0xc9,
	0x59, 0x32, 0xe2, 0x5b, 0x14, 0xc7, 0x1c, 0xe8, 0x5c, 0x30, 0x32, 0x1a, 0x5f, 0x48, 0xe8, 0xf2,
	0xf6, 0x36, 0xcf, 0xb8, 0x85, 0xc6, 0xe8, 0x2a, 0x5e, 0xec, 0xce, 0x14, 0xf2, 0x5e, 0xad, 0x71,
	0xba, 0x8e, 0x97, 0xbb, 0x73, 0x59, 0xed, 0x40, 0xf5, 0xc8, 0x89, 0x78, 0xd7, 0x15, 0xae, 0xee,
	0x92, 0x49, 0x7f, 0x52, 0xbc, 0x21, 0x72, 0xd3, 0x98, 0xf2, 0x67, 0xe2, 0x7d, 0x3d, 0x92, 0x44,
	0xc2, 0x3f, 0x70, 0x07, 0xd1, 0x4e, 0x61, 0xda, 0x6f, 0xb2, 0x03, 0x0a, 0x27, 0x30, 0x43, 0x27,
	0xf1, 0x68, 0x98, 0xf4, 0xdc, 0x9f, 0xf5, 0x87, 0xfd, 0x5d, 0x09, 0x19, 0x06, 0xfe, 0xe0, 0x44,
	0xf5, 0x80, 0x09, 0x73, 0x9d, 0x07, 0x75, 0x0c, 0x98, 0xf7, 0xef, 0x81, 0x76, 0x87, 0x42, 0x16,
	0x05, 0x2f, 0xd6, 0x05, 0x3a, 0x87, 0xa7, 0x23, 0xdb, 0x34, 0xcf, 0x40, 0x07, 0x58, 0xa4, 0xd3,
	0x78, 0xc2, 0x05, 0x82, 0x03, 0xef, 0xba, 0xb8, 0x44, 0x57, 0xf0, 0x42, 0x8c, 0x18, 0xb1, 0x72,
	0xd9, 0x77, 0x29, 0x1e, 0x6b, 0xfb, 0xb9, 0xe2, 0x8f, 0x7f, 0x3c, 0xb9, 0xc7, 0xf6, 0xd2, 0x4c,
	0x37, 0xc8, 0xaa, 0x7f, 0xac, 0x10, 0xe8, 0x79, 0xba, 0x46, 0x37, 0xf1, 0x6a, 0x8c, 0x1c, 0x77,
	0x15, 0xad, 0x77, 0xe9, 0xae, 0xf3, 0x62, 0xdb, 0xf0, 0x47, 0x37, 0x44, 0x46, 0xaf, 0xa7, 0x4d,
	0x3f, 0xb1, 0x78, 0xce, 0xa9, 0xf6, 0x4b, 0xa2, 0xe7, 0xea, 0x33, 0x48, 0xe9, 0xff, 0xd1, 0x7d,
	0x0b, 0xd0, 0x43, 0x0b, 0xd0, 0x63, 0x0b, 0xd0, 0x53, 0x0b, 0xd0, 0x73, 0x0b, 0xa4, 0x97, 0x16,
	0x48, 0xaf, 0x2d, 0x40, 0x97, 0x16, 0x48, 0x57, 0x16, 0x48, 0xd7, 0x16, 0xa0, 0x1b, 0x0b, 0xa4,
	0x5b, 0x0b, 0xd0, 0x9d, 0x05, 0xe8, 0xde, 0x02, 0xf4, 0x60, 0x01, 0x7a, 0xb4, 0x40, 0x7a, 0xb2,
	0x00, 0x3d, 0x5b, 0x20, 0xbd, 0x58, 0x80, 0x5e, 0x2d, 0x90, 0x2e, 0x6d, 0x90, 0xae, 0x6c, 0x40,
	0xff, 0xda, 0x20, 0xfd, 0x67, 0x03, 0xfa, 0x64, 0x83, 0x74, 0x6d, 0x83, 0x74, 0x63, 0x03, 0xba,
	0xb5, 0x01, 0xdd, 0xd9, 0x80, 0x8e, 0x36, 0xbe, 0xf5, 0x05, 0x6b, 0x9c, 0x96, 0xdf, 0x7e, 0x1b,
	0xa5, 0x52, 0x9f, 0xf3, 0x09, 0xda, 0xfa, 0x1a, 0x00, 0x00, 0xff, 0xff, 0xfd, 0x6a, 0x9a, 0x63,
	0xf5, 0x06, 0x00, 0x00,
}
