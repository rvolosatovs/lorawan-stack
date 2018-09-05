// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/gateway_services.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import types "github.com/gogo/protobuf/types"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import context "context"
import grpc "google.golang.org/grpc"

import strings "strings"
import reflect "reflect"

import io "io"

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

type PullGatewayConfigurationRequest struct {
	GatewayIdentifiers   `protobuf:"bytes,1,opt,name=gateway_ids,json=gatewayIds,embedded=gateway_ids" json:"gateway_ids"`
	FieldMask            *types.FieldMask `protobuf:"bytes,2,opt,name=field_mask,json=fieldMask" json:"field_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *PullGatewayConfigurationRequest) Reset()      { *m = PullGatewayConfigurationRequest{} }
func (*PullGatewayConfigurationRequest) ProtoMessage() {}
func (*PullGatewayConfigurationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_gateway_services_44cd116423204041, []int{0}
}
func (m *PullGatewayConfigurationRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PullGatewayConfigurationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PullGatewayConfigurationRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *PullGatewayConfigurationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PullGatewayConfigurationRequest.Merge(dst, src)
}
func (m *PullGatewayConfigurationRequest) XXX_Size() int {
	return m.Size()
}
func (m *PullGatewayConfigurationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PullGatewayConfigurationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PullGatewayConfigurationRequest proto.InternalMessageInfo

func (m *PullGatewayConfigurationRequest) GetFieldMask() *types.FieldMask {
	if m != nil {
		return m.FieldMask
	}
	return nil
}

func init() {
	proto.RegisterType((*PullGatewayConfigurationRequest)(nil), "ttn.lorawan.v3.PullGatewayConfigurationRequest")
	golang_proto.RegisterType((*PullGatewayConfigurationRequest)(nil), "ttn.lorawan.v3.PullGatewayConfigurationRequest")
}
func (this *PullGatewayConfigurationRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PullGatewayConfigurationRequest)
	if !ok {
		that2, ok := that.(PullGatewayConfigurationRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.GatewayIdentifiers.Equal(&that1.GatewayIdentifiers) {
		return false
	}
	if !this.FieldMask.Equal(that1.FieldMask) {
		return false
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GatewayRegistry service

type GatewayRegistryClient interface {
	// Create a new gateway. This also sets the current user as first collaborator
	// with all possible rights. When including organization identifiers, this instead
	// grants those rights to the given organization.
	CreateGateway(ctx context.Context, in *CreateGatewayRequest, opts ...grpc.CallOption) (*Gateway, error)
	// Get the gateway with the given identifiers, selecting the fields given
	// by the field mask. The method may return more or less fields, depending on
	// the rights of the caller.
	GetGateway(ctx context.Context, in *GetGatewayRequest, opts ...grpc.CallOption) (*Gateway, error)
	// List gateways. See request message for details.
	ListGateways(ctx context.Context, in *ListGatewaysRequest, opts ...grpc.CallOption) (*Gateways, error)
	UpdateGateway(ctx context.Context, in *UpdateGatewayRequest, opts ...grpc.CallOption) (*Gateway, error)
	DeleteGateway(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*types.Empty, error)
}

type gatewayRegistryClient struct {
	cc *grpc.ClientConn
}

func NewGatewayRegistryClient(cc *grpc.ClientConn) GatewayRegistryClient {
	return &gatewayRegistryClient{cc}
}

func (c *gatewayRegistryClient) CreateGateway(ctx context.Context, in *CreateGatewayRequest, opts ...grpc.CallOption) (*Gateway, error) {
	out := new(Gateway)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayRegistry/CreateGateway", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayRegistryClient) GetGateway(ctx context.Context, in *GetGatewayRequest, opts ...grpc.CallOption) (*Gateway, error) {
	out := new(Gateway)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayRegistry/GetGateway", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayRegistryClient) ListGateways(ctx context.Context, in *ListGatewaysRequest, opts ...grpc.CallOption) (*Gateways, error) {
	out := new(Gateways)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayRegistry/ListGateways", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayRegistryClient) UpdateGateway(ctx context.Context, in *UpdateGatewayRequest, opts ...grpc.CallOption) (*Gateway, error) {
	out := new(Gateway)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayRegistry/UpdateGateway", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayRegistryClient) DeleteGateway(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayRegistry/DeleteGateway", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GatewayRegistry service

type GatewayRegistryServer interface {
	// Create a new gateway. This also sets the current user as first collaborator
	// with all possible rights. When including organization identifiers, this instead
	// grants those rights to the given organization.
	CreateGateway(context.Context, *CreateGatewayRequest) (*Gateway, error)
	// Get the gateway with the given identifiers, selecting the fields given
	// by the field mask. The method may return more or less fields, depending on
	// the rights of the caller.
	GetGateway(context.Context, *GetGatewayRequest) (*Gateway, error)
	// List gateways. See request message for details.
	ListGateways(context.Context, *ListGatewaysRequest) (*Gateways, error)
	UpdateGateway(context.Context, *UpdateGatewayRequest) (*Gateway, error)
	DeleteGateway(context.Context, *GatewayIdentifiers) (*types.Empty, error)
}

func RegisterGatewayRegistryServer(s *grpc.Server, srv GatewayRegistryServer) {
	s.RegisterService(&_GatewayRegistry_serviceDesc, srv)
}

func _GatewayRegistry_CreateGateway_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGatewayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayRegistryServer).CreateGateway(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayRegistry/CreateGateway",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayRegistryServer).CreateGateway(ctx, req.(*CreateGatewayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayRegistry_GetGateway_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGatewayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayRegistryServer).GetGateway(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayRegistry/GetGateway",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayRegistryServer).GetGateway(ctx, req.(*GetGatewayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayRegistry_ListGateways_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGatewaysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayRegistryServer).ListGateways(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayRegistry/ListGateways",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayRegistryServer).ListGateways(ctx, req.(*ListGatewaysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayRegistry_UpdateGateway_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGatewayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayRegistryServer).UpdateGateway(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayRegistry/UpdateGateway",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayRegistryServer).UpdateGateway(ctx, req.(*UpdateGatewayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayRegistry_DeleteGateway_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayRegistryServer).DeleteGateway(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayRegistry/DeleteGateway",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayRegistryServer).DeleteGateway(ctx, req.(*GatewayIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _GatewayRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.GatewayRegistry",
	HandlerType: (*GatewayRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGateway",
			Handler:    _GatewayRegistry_CreateGateway_Handler,
		},
		{
			MethodName: "GetGateway",
			Handler:    _GatewayRegistry_GetGateway_Handler,
		},
		{
			MethodName: "ListGateways",
			Handler:    _GatewayRegistry_ListGateways_Handler,
		},
		{
			MethodName: "UpdateGateway",
			Handler:    _GatewayRegistry_UpdateGateway_Handler,
		},
		{
			MethodName: "DeleteGateway",
			Handler:    _GatewayRegistry_DeleteGateway_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/gateway_services.proto",
}

// Client API for GatewayAccess service

type GatewayAccessClient interface {
	ListGatewayRights(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*Rights, error)
	GenerateGatewayAPIKey(ctx context.Context, in *SetGatewayAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	ListGatewayAPIKeys(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*APIKeys, error)
	// Set the rights of an existing gateway API key. To generate an API key,
	// the GenerateGatewayAPIKey should be used. To delete an API key, set it
	// without any rights.
	SetGatewayAPIKey(ctx context.Context, in *SetGatewayAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	// Set the rights of a collaborator on the gateway. Users or organizations
	// are considered to be a collaborator if they have at least one right on the
	// gateway.
	SetGatewayCollaborator(ctx context.Context, in *SetGatewayCollaboratorRequest, opts ...grpc.CallOption) (*types.Empty, error)
	ListGatewayCollaborators(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*Collaborators, error)
}

type gatewayAccessClient struct {
	cc *grpc.ClientConn
}

func NewGatewayAccessClient(cc *grpc.ClientConn) GatewayAccessClient {
	return &gatewayAccessClient{cc}
}

func (c *gatewayAccessClient) ListGatewayRights(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*Rights, error) {
	out := new(Rights)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayAccess/ListGatewayRights", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayAccessClient) GenerateGatewayAPIKey(ctx context.Context, in *SetGatewayAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayAccess/GenerateGatewayAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayAccessClient) ListGatewayAPIKeys(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*APIKeys, error) {
	out := new(APIKeys)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayAccess/ListGatewayAPIKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayAccessClient) SetGatewayAPIKey(ctx context.Context, in *SetGatewayAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayAccess/SetGatewayAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayAccessClient) SetGatewayCollaborator(ctx context.Context, in *SetGatewayCollaboratorRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayAccess/SetGatewayCollaborator", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayAccessClient) ListGatewayCollaborators(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*Collaborators, error) {
	out := new(Collaborators)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GatewayAccess/ListGatewayCollaborators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GatewayAccess service

type GatewayAccessServer interface {
	ListGatewayRights(context.Context, *GatewayIdentifiers) (*Rights, error)
	GenerateGatewayAPIKey(context.Context, *SetGatewayAPIKeyRequest) (*APIKey, error)
	ListGatewayAPIKeys(context.Context, *GatewayIdentifiers) (*APIKeys, error)
	// Set the rights of an existing gateway API key. To generate an API key,
	// the GenerateGatewayAPIKey should be used. To delete an API key, set it
	// without any rights.
	SetGatewayAPIKey(context.Context, *SetGatewayAPIKeyRequest) (*APIKey, error)
	// Set the rights of a collaborator on the gateway. Users or organizations
	// are considered to be a collaborator if they have at least one right on the
	// gateway.
	SetGatewayCollaborator(context.Context, *SetGatewayCollaboratorRequest) (*types.Empty, error)
	ListGatewayCollaborators(context.Context, *GatewayIdentifiers) (*Collaborators, error)
}

func RegisterGatewayAccessServer(s *grpc.Server, srv GatewayAccessServer) {
	s.RegisterService(&_GatewayAccess_serviceDesc, srv)
}

func _GatewayAccess_ListGatewayRights_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayAccessServer).ListGatewayRights(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayAccess/ListGatewayRights",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayAccessServer).ListGatewayRights(ctx, req.(*GatewayIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayAccess_GenerateGatewayAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGatewayAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayAccessServer).GenerateGatewayAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayAccess/GenerateGatewayAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayAccessServer).GenerateGatewayAPIKey(ctx, req.(*SetGatewayAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayAccess_ListGatewayAPIKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayAccessServer).ListGatewayAPIKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayAccess/ListGatewayAPIKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayAccessServer).ListGatewayAPIKeys(ctx, req.(*GatewayIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayAccess_SetGatewayAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGatewayAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayAccessServer).SetGatewayAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayAccess/SetGatewayAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayAccessServer).SetGatewayAPIKey(ctx, req.(*SetGatewayAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayAccess_SetGatewayCollaborator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetGatewayCollaboratorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayAccessServer).SetGatewayCollaborator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayAccess/SetGatewayCollaborator",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayAccessServer).SetGatewayCollaborator(ctx, req.(*SetGatewayCollaboratorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayAccess_ListGatewayCollaborators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayAccessServer).ListGatewayCollaborators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GatewayAccess/ListGatewayCollaborators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayAccessServer).ListGatewayCollaborators(ctx, req.(*GatewayIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _GatewayAccess_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.GatewayAccess",
	HandlerType: (*GatewayAccessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListGatewayRights",
			Handler:    _GatewayAccess_ListGatewayRights_Handler,
		},
		{
			MethodName: "GenerateGatewayAPIKey",
			Handler:    _GatewayAccess_GenerateGatewayAPIKey_Handler,
		},
		{
			MethodName: "ListGatewayAPIKeys",
			Handler:    _GatewayAccess_ListGatewayAPIKeys_Handler,
		},
		{
			MethodName: "SetGatewayAPIKey",
			Handler:    _GatewayAccess_SetGatewayAPIKey_Handler,
		},
		{
			MethodName: "SetGatewayCollaborator",
			Handler:    _GatewayAccess_SetGatewayCollaborator_Handler,
		},
		{
			MethodName: "ListGatewayCollaborators",
			Handler:    _GatewayAccess_ListGatewayCollaborators_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/gateway_services.proto",
}

// Client API for GatewayConfigurator service

type GatewayConfiguratorClient interface {
	PullConfiguration(ctx context.Context, in *PullGatewayConfigurationRequest, opts ...grpc.CallOption) (GatewayConfigurator_PullConfigurationClient, error)
}

type gatewayConfiguratorClient struct {
	cc *grpc.ClientConn
}

func NewGatewayConfiguratorClient(cc *grpc.ClientConn) GatewayConfiguratorClient {
	return &gatewayConfiguratorClient{cc}
}

func (c *gatewayConfiguratorClient) PullConfiguration(ctx context.Context, in *PullGatewayConfigurationRequest, opts ...grpc.CallOption) (GatewayConfigurator_PullConfigurationClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GatewayConfigurator_serviceDesc.Streams[0], "/ttn.lorawan.v3.GatewayConfigurator/PullConfiguration", opts...)
	if err != nil {
		return nil, err
	}
	x := &gatewayConfiguratorPullConfigurationClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GatewayConfigurator_PullConfigurationClient interface {
	Recv() (*Gateway, error)
	grpc.ClientStream
}

type gatewayConfiguratorPullConfigurationClient struct {
	grpc.ClientStream
}

func (x *gatewayConfiguratorPullConfigurationClient) Recv() (*Gateway, error) {
	m := new(Gateway)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for GatewayConfigurator service

type GatewayConfiguratorServer interface {
	PullConfiguration(*PullGatewayConfigurationRequest, GatewayConfigurator_PullConfigurationServer) error
}

func RegisterGatewayConfiguratorServer(s *grpc.Server, srv GatewayConfiguratorServer) {
	s.RegisterService(&_GatewayConfigurator_serviceDesc, srv)
}

func _GatewayConfigurator_PullConfiguration_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullGatewayConfigurationRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GatewayConfiguratorServer).PullConfiguration(m, &gatewayConfiguratorPullConfigurationServer{stream})
}

type GatewayConfigurator_PullConfigurationServer interface {
	Send(*Gateway) error
	grpc.ServerStream
}

type gatewayConfiguratorPullConfigurationServer struct {
	grpc.ServerStream
}

func (x *gatewayConfiguratorPullConfigurationServer) Send(m *Gateway) error {
	return x.ServerStream.SendMsg(m)
}

var _GatewayConfigurator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.GatewayConfigurator",
	HandlerType: (*GatewayConfiguratorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullConfiguration",
			Handler:       _GatewayConfigurator_PullConfiguration_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "lorawan-stack/api/gateway_services.proto",
}

func (m *PullGatewayConfigurationRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PullGatewayConfigurationRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintGatewayServices(dAtA, i, uint64(m.GatewayIdentifiers.Size()))
	n1, err := m.GatewayIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if m.FieldMask != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintGatewayServices(dAtA, i, uint64(m.FieldMask.Size()))
		n2, err := m.FieldMask.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func encodeVarintGatewayServices(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedPullGatewayConfigurationRequest(r randyGatewayServices, easy bool) *PullGatewayConfigurationRequest {
	this := &PullGatewayConfigurationRequest{}
	v1 := NewPopulatedGatewayIdentifiers(r, easy)
	this.GatewayIdentifiers = *v1
	if r.Intn(10) != 0 {
		this.FieldMask = types.NewPopulatedFieldMask(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyGatewayServices interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneGatewayServices(r randyGatewayServices) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringGatewayServices(r randyGatewayServices) string {
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
		tmps[i] = randUTF8RuneGatewayServices(r)
	}
	return string(tmps)
}
func randUnrecognizedGatewayServices(r randyGatewayServices, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldGatewayServices(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldGatewayServices(dAtA []byte, r randyGatewayServices, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateGatewayServices(dAtA, uint64(key))
		v3 := r.Int63()
		if r.Intn(2) == 0 {
			v3 *= -1
		}
		dAtA = encodeVarintPopulateGatewayServices(dAtA, uint64(v3))
	case 1:
		dAtA = encodeVarintPopulateGatewayServices(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateGatewayServices(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateGatewayServices(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateGatewayServices(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateGatewayServices(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *PullGatewayConfigurationRequest) Size() (n int) {
	var l int
	_ = l
	l = m.GatewayIdentifiers.Size()
	n += 1 + l + sovGatewayServices(uint64(l))
	if m.FieldMask != nil {
		l = m.FieldMask.Size()
		n += 1 + l + sovGatewayServices(uint64(l))
	}
	return n
}

func sovGatewayServices(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozGatewayServices(x uint64) (n int) {
	return sovGatewayServices((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *PullGatewayConfigurationRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PullGatewayConfigurationRequest{`,
		`GatewayIdentifiers:` + strings.Replace(strings.Replace(this.GatewayIdentifiers.String(), "GatewayIdentifiers", "GatewayIdentifiers", 1), `&`, ``, 1) + `,`,
		`FieldMask:` + strings.Replace(fmt.Sprintf("%v", this.FieldMask), "FieldMask", "types.FieldMask", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGatewayServices(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *PullGatewayConfigurationRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGatewayServices
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PullGatewayConfigurationRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PullGatewayConfigurationRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GatewayIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayServices
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGatewayServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GatewayIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FieldMask", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayServices
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGatewayServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.FieldMask == nil {
				m.FieldMask = &types.FieldMask{}
			}
			if err := m.FieldMask.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGatewayServices(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGatewayServices
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGatewayServices(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGatewayServices
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGatewayServices
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGatewayServices
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthGatewayServices
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGatewayServices
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGatewayServices(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGatewayServices = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGatewayServices   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("lorawan-stack/api/gateway_services.proto", fileDescriptor_gateway_services_44cd116423204041)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/gateway_services.proto", fileDescriptor_gateway_services_44cd116423204041)
}

var fileDescriptor_gateway_services_44cd116423204041 = []byte{
	// 898 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x4d, 0x8c, 0xdb, 0x44,
	0x14, 0x9e, 0xa9, 0x44, 0x45, 0xa7, 0x04, 0xe8, 0x20, 0xb6, 0x91, 0xa1, 0x93, 0xd6, 0x5d, 0xd8,
	0x6a, 0x69, 0xec, 0x25, 0x15, 0x48, 0x94, 0x03, 0x6a, 0x17, 0x88, 0x56, 0x50, 0x51, 0x05, 0xf5,
	0x92, 0xcb, 0xca, 0x49, 0x26, 0x8e, 0x15, 0xaf, 0x27, 0x78, 0x26, 0x59, 0xa5, 0x51, 0xa4, 0x8a,
	0x53, 0x11, 0x97, 0x95, 0xb8, 0x70, 0xe0, 0x80, 0x90, 0x90, 0x2a, 0x71, 0xe9, 0xb1, 0xc7, 0x1e,
	0xf7, 0xb8, 0x12, 0x12, 0xda, 0xd3, 0x6a, 0x63, 0x73, 0xd8, 0xe3, 0x1e, 0xf7, 0x84, 0x90, 0xc7,
	0x76, 0xec, 0x38, 0x09, 0xc9, 0x8a, 0xdb, 0x8c, 0xdf, 0x9b, 0xef, 0xfb, 0xde, 0xf7, 0xe6, 0xc7,
	0xe8, 0x96, 0xcd, 0x5c, 0x63, 0xd7, 0x70, 0x8a, 0x5c, 0x18, 0xf5, 0xb6, 0x6e, 0x74, 0x2c, 0xdd,
	0x34, 0x04, 0xdd, 0x35, 0xfa, 0xdb, 0x9c, 0xba, 0x3d, 0xab, 0x4e, 0xb9, 0xd6, 0x71, 0x99, 0x60,
	0xf8, 0x75, 0x21, 0x1c, 0x2d, 0xca, 0xd6, 0x7a, 0x77, 0x94, 0xa2, 0x69, 0x89, 0x56, 0xb7, 0xa6,
	0xd5, 0xd9, 0x8e, 0x6e, 0x32, 0x93, 0xe9, 0x32, 0xad, 0xd6, 0x6d, 0xca, 0x99, 0x9c, 0xc8, 0x51,
	0xb8, 0x5c, 0x79, 0xd7, 0x64, 0xcc, 0xb4, 0xa9, 0x64, 0x30, 0x1c, 0x87, 0x09, 0x43, 0x58, 0xcc,
	0x89, 0xc0, 0x95, 0x77, 0xa2, 0xe8, 0x18, 0x83, 0xee, 0x74, 0x44, 0x3f, 0x0a, 0x5e, 0xcf, 0x06,
	0x9b, 0x16, 0xb5, 0x1b, 0xdb, 0x3b, 0x06, 0x6f, 0x47, 0x19, 0x85, 0xb9, 0x55, 0x44, 0x09, 0x37,
	0xa7, 0x13, 0xac, 0x06, 0x75, 0x84, 0xd5, 0xb4, 0xa8, 0x1b, 0x8b, 0x20, 0xd3, 0x49, 0xae, 0x65,
	0xb6, 0x44, 0x14, 0x57, 0xff, 0x80, 0xa8, 0xf0, 0xb0, 0x6b, 0xdb, 0xe5, 0x10, 0x7a, 0x93, 0x39,
	0x4d, 0xcb, 0xec, 0xba, 0xb2, 0x90, 0x0a, 0xfd, 0xae, 0x4b, 0xb9, 0xc0, 0x0f, 0xd0, 0xe5, 0xd8,
	0x3f, 0xab, 0xc1, 0xf3, 0xf0, 0x3a, 0xbc, 0x75, 0xb9, 0xa4, 0x6a, 0x93, 0xde, 0x69, 0x11, 0xc2,
	0x56, 0x22, 0xe1, 0xfe, 0xab, 0xfb, 0x47, 0x05, 0x70, 0x70, 0x54, 0x80, 0x15, 0x64, 0xc6, 0x51,
	0x8e, 0x3f, 0x41, 0x28, 0x29, 0x36, 0x7f, 0x41, 0xa2, 0x29, 0x5a, 0xe8, 0x87, 0x16, 0xfb, 0xa1,
	0x7d, 0x19, 0xa4, 0x3c, 0x30, 0x78, 0xbb, 0x72, 0xa9, 0x19, 0x0f, 0x4b, 0xff, 0xbc, 0x82, 0xde,
	0x88, 0x78, 0x2a, 0xd4, 0xb4, 0xb8, 0x70, 0xfb, 0xf8, 0x77, 0x88, 0x72, 0x9b, 0x2e, 0x35, 0x04,
	0x8d, 0x22, 0x78, 0x35, 0x2b, 0x6d, 0x22, 0x1c, 0x55, 0xa5, 0x5c, 0x9d, 0x53, 0x80, 0xfa, 0xe8,
	0xfb, 0x3f, 0xff, 0xfe, 0xe9, 0xc2, 0x37, 0xea, 0xa5, 0xd8, 0x6f, 0x7e, 0x17, 0xae, 0x57, 0x3f,
	0x53, 0xef, 0xea, 0xcc, 0x35, 0x0d, 0xc7, 0x7a, 0x1c, 0x76, 0x59, 0x1f, 0xa4, 0xa7, 0x81, 0x33,
	0x5a, 0xe6, 0xc3, 0x30, 0x0d, 0x80, 0x05, 0x42, 0x65, 0x2a, 0x62, 0x8d, 0x37, 0xa6, 0xd8, 0xc7,
	0xb1, 0x85, 0x02, 0xd7, 0xa5, 0xc0, 0x55, 0xac, 0x8e, 0xf1, 0xf5, 0x41, 0xaa, 0x41, 0x5a, 0x32,
	0x1e, 0xe2, 0x63, 0x88, 0x5e, 0xfb, 0xda, 0xe2, 0x31, 0x36, 0xc7, 0x37, 0xb3, 0xa8, 0xe9, 0x68,
	0x4c, 0x9d, 0x9f, 0x43, 0xcd, 0xd5, 0x3d, 0x28, 0xc9, 0x7f, 0x80, 0x38, 0xb1, 0xa7, 0xfa, 0x21,
	0xd6, 0xf5, 0x2e, 0xa7, 0x2e, 0xd7, 0x07, 0x75, 0x66, 0xdb, 0x46, 0x8d, 0xb9, 0x86, 0x60, 0xae,
	0x16, 0x7c, 0x93, 0x8a, 0xa2, 0x41, 0x62, 0x48, 0x75, 0x0b, 0x97, 0xb3, 0x76, 0x4e, 0x2c, 0x5d,
	0xde, 0x5b, 0x3c, 0x44, 0xb9, 0x47, 0x9d, 0xc6, 0x7f, 0xf5, 0x7f, 0x22, 0xbc, 0xd0, 0xde, 0xa2,
	0xac, 0x70, 0x4d, 0x99, 0x61, 0xaf, 0x96, 0xb1, 0x37, 0xe8, 0x6b, 0x1b, 0xe5, 0x3e, 0xa7, 0x36,
	0x4d, 0xe8, 0x97, 0x38, 0x19, 0xca, 0xca, 0xd4, 0x7e, 0xff, 0x22, 0xb8, 0x1c, 0x54, 0x22, 0xb9,
	0xf3, 0xeb, 0x2b, 0x33, 0x5b, 0x3b, 0x2c, 0xfd, 0x75, 0x11, 0xe5, 0x22, 0xb8, 0x7b, 0xf5, 0x3a,
	0xe5, 0x1c, 0xf7, 0xd0, 0x95, 0x54, 0x07, 0x2b, 0xf2, 0x6c, 0x2f, 0x29, 0x21, 0x93, 0x13, 0xae,
	0x55, 0xdf, 0x93, 0x12, 0x0a, 0xf8, 0xda, 0x6c, 0x09, 0xd1, 0xf5, 0x81, 0xf7, 0x20, 0x7a, 0xbb,
	0x4c, 0x1d, 0xea, 0x26, 0xce, 0xde, 0x7b, 0xb8, 0xf5, 0x15, 0xed, 0xe3, 0xb5, 0x2c, 0xf0, 0xb7,
	0xe3, 0xad, 0x1d, 0x66, 0xc4, 0x1d, 0x98, 0x52, 0x10, 0x86, 0xd5, 0x8f, 0xa5, 0x82, 0x0d, 0xf5,
	0x83, 0xc5, 0xfb, 0x3b, 0xb8, 0xd0, 0x8a, 0x6d, 0x1a, 0x9e, 0xb0, 0xc7, 0x08, 0xa7, 0xac, 0x08,
	0xc1, 0x96, 0xf3, 0xe2, 0xea, 0x6c, 0x25, 0x5c, 0x5d, 0x93, 0x52, 0x6e, 0xe0, 0xc2, 0x1c, 0x33,
	0x62, 0xfa, 0xc0, 0x8e, 0x37, 0xb3, 0x75, 0xfe, 0x7f, 0x27, 0x3e, 0x95, 0xf4, 0x1f, 0x29, 0x1b,
	0xe7, 0x70, 0x42, 0x1f, 0x44, 0x1b, 0xf3, 0x17, 0x88, 0x56, 0x12, 0xc2, 0xcd, 0xd4, 0xe9, 0xc2,
	0xc5, 0xf9, 0xc2, 0xd2, 0x79, 0x89, 0xbc, 0xd9, 0xbb, 0xf5, 0x3c, 0xf2, 0xd2, 0xa7, 0x5b, 0x76,
	0xeb, 0x47, 0x88, 0xf2, 0xa9, 0x76, 0xa5, 0x79, 0x97, 0x6b, 0xda, 0xb5, 0xa9, 0x6b, 0x3e, 0x0d,
	0xa1, 0xde, 0x96, 0xe2, 0xde, 0xc7, 0xab, 0x73, 0x5a, 0x37, 0x21, 0xa8, 0xd4, 0x43, 0x6f, 0x4d,
	0x3d, 0x81, 0xcc, 0xc5, 0xdb, 0xe8, 0x4a, 0xf0, 0x3a, 0x4e, 0x3c, 0x8b, 0x58, 0xcf, 0x12, 0x2f,
	0x78, 0x40, 0xe7, 0x5e, 0x35, 0x1b, 0xf0, 0xfe, 0x6f, 0x70, 0x7f, 0x44, 0xe0, 0xc1, 0x88, 0xc0,
	0xc3, 0x11, 0x01, 0xc7, 0x23, 0x02, 0x4e, 0x46, 0x04, 0x9c, 0x8e, 0x08, 0x38, 0x1b, 0x11, 0xf8,
	0xc4, 0x23, 0xf0, 0xa9, 0x47, 0xc0, 0x33, 0x8f, 0xc0, 0xe7, 0x1e, 0x01, 0x2f, 0x3c, 0x02, 0x5e,
	0x7a, 0x04, 0xec, 0x7b, 0x04, 0x1e, 0x78, 0x04, 0x1e, 0x7a, 0x04, 0x1c, 0x7b, 0x04, 0x9e, 0x78,
	0x04, 0x9c, 0x7a, 0x04, 0x9e, 0x79, 0x04, 0x3c, 0xf1, 0x09, 0x78, 0xea, 0x13, 0xb8, 0xe7, 0x13,
	0xf0, 0xb3, 0x4f, 0xe0, 0xaf, 0x3e, 0x01, 0xcf, 0x7c, 0x02, 0x9e, 0xfb, 0x04, 0xbe, 0xf0, 0x09,
	0x7c, 0xe9, 0x13, 0x58, 0xbd, 0x6d, 0x32, 0x4d, 0xb4, 0xa8, 0x68, 0x59, 0x8e, 0xc9, 0x35, 0x87,
	0x8a, 0x5d, 0xe6, 0xb6, 0xf5, 0xc9, 0x7f, 0x85, 0x4e, 0xdb, 0xd4, 0x85, 0x70, 0x3a, 0xb5, 0xda,
	0x45, 0xd9, 0xf7, 0x3b, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x72, 0xc1, 0xb2, 0x48, 0x59, 0x09,
	0x00, 0x00,
}
