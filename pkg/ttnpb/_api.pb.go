// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/_api.proto

/*
Package ttnpb is a generated protocol buffer package.

The Things Network v3 API

It is generated from these files:
	github.com/TheThingsNetwork/ttn/api/_api.proto
	github.com/TheThingsNetwork/ttn/api/application.proto
	github.com/TheThingsNetwork/ttn/api/applicationserver.proto
	github.com/TheThingsNetwork/ttn/api/auth.proto
	github.com/TheThingsNetwork/ttn/api/client.proto
	github.com/TheThingsNetwork/ttn/api/cluster.proto
	github.com/TheThingsNetwork/ttn/api/collaborator.proto
	github.com/TheThingsNetwork/ttn/api/end_device.proto
	github.com/TheThingsNetwork/ttn/api/gateway.proto
	github.com/TheThingsNetwork/ttn/api/gatewayserver.proto
	github.com/TheThingsNetwork/ttn/api/identifiers.proto
	github.com/TheThingsNetwork/ttn/api/identityserver.proto
	github.com/TheThingsNetwork/ttn/api/join.proto
	github.com/TheThingsNetwork/ttn/api/joinserver.proto
	github.com/TheThingsNetwork/ttn/api/lorawan.proto
	github.com/TheThingsNetwork/ttn/api/messages.proto
	github.com/TheThingsNetwork/ttn/api/metadata.proto
	github.com/TheThingsNetwork/ttn/api/networkserver.proto
	github.com/TheThingsNetwork/ttn/api/rights.proto
	github.com/TheThingsNetwork/ttn/api/tokenkey.proto
	github.com/TheThingsNetwork/ttn/api/user.proto

It has these top-level messages:
	Application
	ApplicationUp
	ApplicationUplink
	ApplicationDownlink
	ApplicationDownlinks
	DownlinkQueueRequest
	APIKey
	Client
	PeerInfo
	ApplicationCollaborator
	GatewayCollaborator
	KeyEnvelope
	RootKeys
	SessionKeys
	Session
	EndDevice
	EndDevices
	MACSettings
	MACState
	MACInfo
	Gateway
	GatewayPrivacySettings
	GatewayAntenna
	GatewayRadio
	GatewayStatus
	GatewayObservations
	GatewayUp
	GatewayDown
	FrequencyPlan
	FrequencyPlanRequest
	UserIdentifier
	ApplicationIdentifier
	GatewayIdentifier
	EndDeviceIdentifiers
	ClientIdentifier
	PullConfigurationRequest
	IdentityServerSettings
	UpdateSettingsRequest
	CreateUserRequest
	UpdateUserRequest
	UpdateUserPasswordRequest
	GenerateUserAPIKeyRequest
	ListUserAPIKeysResponse
	UpdateUserAPIKeyRequest
	RemoveUserAPIKeyRequest
	ValidateUserEmailRequest
	ListAuthorizedClientsResponse
	CreateApplicationRequest
	ListApplicationsResponse
	UpdateApplicationRequest
	GenerateApplicationAPIKeyRequest
	ListApplicationAPIKeysResponse
	UpdateApplicationAPIKeyRequest
	RemoveApplicationAPIKeyRequest
	ListApplicationCollaboratorsResponse
	ListApplicationRightsResponse
	CreateGatewayRequest
	ListGatewaysResponse
	UpdateGatewayRequest
	GenerateGatewayAPIKeyRequest
	ListGatewayAPIKeysResponse
	UpdateGatewayAPIKeyRequest
	RemoveGatewayAPIKeyRequest
	ListGatewayCollaboratorsResponse
	ListGatewayRightsResponse
	CreateClientRequest
	ListClientsResponse
	UpdateClientRequest
	JoinRequest
	JoinResponse
	SessionKeyRequest
	NwkSKeysResponse
	AppSKeyResponse
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
	MACCommand
	UplinkMessage
	DownlinkMessage
	RxMetadata
	TxMetadata
	Location
	GetTokenInfoRequest
	GetTokenInfoResponse
	GetKeyInfoRequest
	GetKeyInfoResponse
	User
*/
package ttnpb

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

func init() { proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/_api.proto", fileDescriptorXApi) }
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/_api.proto", fileDescriptorXApi)
}

var fileDescriptorXApi = []byte{
	// 203 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xce, 0x21, 0x52, 0xc3, 0x40,
	0x14, 0x06, 0xe0, 0xf7, 0x9b, 0x0a, 0x24, 0x07, 0xf8, 0x3d, 0x66, 0x57, 0xf4, 0x06, 0x1c, 0x00,
	0x55, 0x85, 0x61, 0x1a, 0xa6, 0x93, 0xec, 0x74, 0xc8, 0xee, 0x94, 0x05, 0x6c, 0x65, 0x25, 0x12,
	0x07, 0xb2, 0xb2, 0xb2, 0xb2, 0xb2, 0xb2, 0x32, 0x32, 0xfb, 0x9e, 0x89, 0x8c, 0x8c, 0x64, 0xe0,
	0x02, 0x55, 0x9f, 0xfd, 0x6e, 0x5c, 0x1d, 0x72, 0xf3, 0x56, 0xb9, 0xe7, 0xf8, 0xe2, 0x17, 0xcd,
	0x6a, 0xd1, 0x84, 0xb6, 0x7e, 0x7d, 0x58, 0xe5, 0x8f, 0xb8, 0x59, 0xfb, 0x9c, 0x5b, 0xbf, 0x4c,
	0xc1, 0x3f, 0x2d, 0x53, 0x70, 0x69, 0x13, 0x73, 0xbc, 0x9d, 0xe5, 0xdc, 0xba, 0xf7, 0xf9, 0xfd,
	0x37, 0xce, 0x85, 0xb8, 0x14, 0xa2, 0x2b, 0x44, 0x5f, 0x88, 0xa1, 0x50, 0xc6, 0x42, 0x99, 0x0a,
	0xb1, 0x55, 0xca, 0x4e, 0x29, 0x7b, 0x25, 0x0e, 0x4a, 0x39, 0x2a, 0x71, 0x52, 0xe2, 0xac, 0xc4,
	0x45, 0x89, 0x4e, 0x29, 0xbd, 0x12, 0x83, 0x52, 0x46, 0x25, 0x26, 0xa5, 0x6c, 0x8d, 0xb2, 0x33,
	0xe2, 0xd3, 0x28, 0x5f, 0x46, 0xfc, 0x18, 0x65, 0x6f, 0x94, 0x83, 0x11, 0x47, 0x23, 0x4e, 0x46,
	0x3c, 0xde, 0x5d, 0xeb, 0xa6, 0x75, 0xfd, 0x67, 0xaa, 0xaa, 0xd9, 0x7f, 0x78, 0xfe, 0x1b, 0x00,
	0x00, 0xff, 0xff, 0xe5, 0x44, 0x0c, 0x8a, 0xe2, 0x00, 0x00, 0x00,
}
