// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/organization_services.proto

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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for OrganizationRegistry service

type OrganizationRegistryClient interface {
	// Create a new organization. This also sets the current user as first collaborator
	// with all possible rights.
	CreateOrganization(ctx context.Context, in *CreateOrganizationRequest, opts ...grpc.CallOption) (*Organization, error)
	// Get the organization with the given identifiers, selecting the fields given
	// by the field mask. The method may return more or less fields, depending on
	// the rights of the caller.
	GetOrganization(ctx context.Context, in *GetOrganizationRequest, opts ...grpc.CallOption) (*Organization, error)
	// List organizations. See request message for details.
	ListOrganizations(ctx context.Context, in *ListOrganizationsRequest, opts ...grpc.CallOption) (*Organizations, error)
	UpdateOrganization(ctx context.Context, in *UpdateOrganizationRequest, opts ...grpc.CallOption) (*Organization, error)
	DeleteOrganization(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error)
}

type organizationRegistryClient struct {
	cc *grpc.ClientConn
}

func NewOrganizationRegistryClient(cc *grpc.ClientConn) OrganizationRegistryClient {
	return &organizationRegistryClient{cc}
}

func (c *organizationRegistryClient) CreateOrganization(ctx context.Context, in *CreateOrganizationRequest, opts ...grpc.CallOption) (*Organization, error) {
	out := new(Organization)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationRegistry/CreateOrganization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationRegistryClient) GetOrganization(ctx context.Context, in *GetOrganizationRequest, opts ...grpc.CallOption) (*Organization, error) {
	out := new(Organization)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationRegistry/GetOrganization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationRegistryClient) ListOrganizations(ctx context.Context, in *ListOrganizationsRequest, opts ...grpc.CallOption) (*Organizations, error) {
	out := new(Organizations)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationRegistry/ListOrganizations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationRegistryClient) UpdateOrganization(ctx context.Context, in *UpdateOrganizationRequest, opts ...grpc.CallOption) (*Organization, error) {
	out := new(Organization)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationRegistry/UpdateOrganization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationRegistryClient) DeleteOrganization(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationRegistry/DeleteOrganization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OrganizationRegistry service

type OrganizationRegistryServer interface {
	// Create a new organization. This also sets the current user as first collaborator
	// with all possible rights.
	CreateOrganization(context.Context, *CreateOrganizationRequest) (*Organization, error)
	// Get the organization with the given identifiers, selecting the fields given
	// by the field mask. The method may return more or less fields, depending on
	// the rights of the caller.
	GetOrganization(context.Context, *GetOrganizationRequest) (*Organization, error)
	// List organizations. See request message for details.
	ListOrganizations(context.Context, *ListOrganizationsRequest) (*Organizations, error)
	UpdateOrganization(context.Context, *UpdateOrganizationRequest) (*Organization, error)
	DeleteOrganization(context.Context, *OrganizationIdentifiers) (*types.Empty, error)
}

func RegisterOrganizationRegistryServer(s *grpc.Server, srv OrganizationRegistryServer) {
	s.RegisterService(&_OrganizationRegistry_serviceDesc, srv)
}

func _OrganizationRegistry_CreateOrganization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrganizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationRegistryServer).CreateOrganization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationRegistry/CreateOrganization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationRegistryServer).CreateOrganization(ctx, req.(*CreateOrganizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationRegistry_GetOrganization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrganizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationRegistryServer).GetOrganization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationRegistry/GetOrganization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationRegistryServer).GetOrganization(ctx, req.(*GetOrganizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationRegistry_ListOrganizations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOrganizationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationRegistryServer).ListOrganizations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationRegistry/ListOrganizations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationRegistryServer).ListOrganizations(ctx, req.(*ListOrganizationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationRegistry_UpdateOrganization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOrganizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationRegistryServer).UpdateOrganization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationRegistry/UpdateOrganization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationRegistryServer).UpdateOrganization(ctx, req.(*UpdateOrganizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationRegistry_DeleteOrganization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrganizationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationRegistryServer).DeleteOrganization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationRegistry/DeleteOrganization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationRegistryServer).DeleteOrganization(ctx, req.(*OrganizationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrganizationRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.OrganizationRegistry",
	HandlerType: (*OrganizationRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrganization",
			Handler:    _OrganizationRegistry_CreateOrganization_Handler,
		},
		{
			MethodName: "GetOrganization",
			Handler:    _OrganizationRegistry_GetOrganization_Handler,
		},
		{
			MethodName: "ListOrganizations",
			Handler:    _OrganizationRegistry_ListOrganizations_Handler,
		},
		{
			MethodName: "UpdateOrganization",
			Handler:    _OrganizationRegistry_UpdateOrganization_Handler,
		},
		{
			MethodName: "DeleteOrganization",
			Handler:    _OrganizationRegistry_DeleteOrganization_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/organization_services.proto",
}

// Client API for OrganizationAccess service

type OrganizationAccessClient interface {
	ListOrganizationRights(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*Rights, error)
	GenerateOrganizationAPIKey(ctx context.Context, in *SetOrganizationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	ListOrganizationAPIKeys(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*APIKeys, error)
	// Set the rights of an existing organization API key. To generate an API key,
	// the GenerateOrganizationAPIKey should be used. To delete an API key, set it
	// without any rights.
	SetOrganizationAPIKey(ctx context.Context, in *SetOrganizationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	// Set the rights of a collaborator (member) on the organization. Users
	// are considered to be a collaborator if they have at least one right on the
	// organization.
	// Note that only users can collaborate (be member of) an organization.
	SetOrganizationCollaborator(ctx context.Context, in *SetOrganizationCollaboratorRequest, opts ...grpc.CallOption) (*types.Empty, error)
	ListOrganizationCollaborators(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*Collaborators, error)
}

type organizationAccessClient struct {
	cc *grpc.ClientConn
}

func NewOrganizationAccessClient(cc *grpc.ClientConn) OrganizationAccessClient {
	return &organizationAccessClient{cc}
}

func (c *organizationAccessClient) ListOrganizationRights(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*Rights, error) {
	out := new(Rights)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationAccess/ListOrganizationRights", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationAccessClient) GenerateOrganizationAPIKey(ctx context.Context, in *SetOrganizationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationAccess/GenerateOrganizationAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationAccessClient) ListOrganizationAPIKeys(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*APIKeys, error) {
	out := new(APIKeys)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationAccess/ListOrganizationAPIKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationAccessClient) SetOrganizationAPIKey(ctx context.Context, in *SetOrganizationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationAccess/SetOrganizationAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationAccessClient) SetOrganizationCollaborator(ctx context.Context, in *SetOrganizationCollaboratorRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationAccess/SetOrganizationCollaborator", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationAccessClient) ListOrganizationCollaborators(ctx context.Context, in *OrganizationIdentifiers, opts ...grpc.CallOption) (*Collaborators, error) {
	out := new(Collaborators)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.OrganizationAccess/ListOrganizationCollaborators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OrganizationAccess service

type OrganizationAccessServer interface {
	ListOrganizationRights(context.Context, *OrganizationIdentifiers) (*Rights, error)
	GenerateOrganizationAPIKey(context.Context, *SetOrganizationAPIKeyRequest) (*APIKey, error)
	ListOrganizationAPIKeys(context.Context, *OrganizationIdentifiers) (*APIKeys, error)
	// Set the rights of an existing organization API key. To generate an API key,
	// the GenerateOrganizationAPIKey should be used. To delete an API key, set it
	// without any rights.
	SetOrganizationAPIKey(context.Context, *SetOrganizationAPIKeyRequest) (*APIKey, error)
	// Set the rights of a collaborator (member) on the organization. Users
	// are considered to be a collaborator if they have at least one right on the
	// organization.
	// Note that only users can collaborate (be member of) an organization.
	SetOrganizationCollaborator(context.Context, *SetOrganizationCollaboratorRequest) (*types.Empty, error)
	ListOrganizationCollaborators(context.Context, *OrganizationIdentifiers) (*Collaborators, error)
}

func RegisterOrganizationAccessServer(s *grpc.Server, srv OrganizationAccessServer) {
	s.RegisterService(&_OrganizationAccess_serviceDesc, srv)
}

func _OrganizationAccess_ListOrganizationRights_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrganizationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationAccessServer).ListOrganizationRights(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationAccess/ListOrganizationRights",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationAccessServer).ListOrganizationRights(ctx, req.(*OrganizationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationAccess_GenerateOrganizationAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOrganizationAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationAccessServer).GenerateOrganizationAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationAccess/GenerateOrganizationAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationAccessServer).GenerateOrganizationAPIKey(ctx, req.(*SetOrganizationAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationAccess_ListOrganizationAPIKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrganizationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationAccessServer).ListOrganizationAPIKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationAccess/ListOrganizationAPIKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationAccessServer).ListOrganizationAPIKeys(ctx, req.(*OrganizationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationAccess_SetOrganizationAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOrganizationAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationAccessServer).SetOrganizationAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationAccess/SetOrganizationAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationAccessServer).SetOrganizationAPIKey(ctx, req.(*SetOrganizationAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationAccess_SetOrganizationCollaborator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOrganizationCollaboratorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationAccessServer).SetOrganizationCollaborator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationAccess/SetOrganizationCollaborator",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationAccessServer).SetOrganizationCollaborator(ctx, req.(*SetOrganizationCollaboratorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationAccess_ListOrganizationCollaborators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrganizationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationAccessServer).ListOrganizationCollaborators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.OrganizationAccess/ListOrganizationCollaborators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationAccessServer).ListOrganizationCollaborators(ctx, req.(*OrganizationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrganizationAccess_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.OrganizationAccess",
	HandlerType: (*OrganizationAccessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListOrganizationRights",
			Handler:    _OrganizationAccess_ListOrganizationRights_Handler,
		},
		{
			MethodName: "GenerateOrganizationAPIKey",
			Handler:    _OrganizationAccess_GenerateOrganizationAPIKey_Handler,
		},
		{
			MethodName: "ListOrganizationAPIKeys",
			Handler:    _OrganizationAccess_ListOrganizationAPIKeys_Handler,
		},
		{
			MethodName: "SetOrganizationAPIKey",
			Handler:    _OrganizationAccess_SetOrganizationAPIKey_Handler,
		},
		{
			MethodName: "SetOrganizationCollaborator",
			Handler:    _OrganizationAccess_SetOrganizationCollaborator_Handler,
		},
		{
			MethodName: "ListOrganizationCollaborators",
			Handler:    _OrganizationAccess_ListOrganizationCollaborators_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/organization_services.proto",
}

func init() {
	proto.RegisterFile("lorawan-stack/api/organization_services.proto", fileDescriptor_organization_services_7e2f7f02a3b48700)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/organization_services.proto", fileDescriptor_organization_services_7e2f7f02a3b48700)
}

var fileDescriptor_organization_services_7e2f7f02a3b48700 = []byte{
	// 756 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0x4f, 0x4c, 0x14, 0x49,
	0x14, 0xc6, 0xab, 0x36, 0xd9, 0xdd, 0xa4, 0x0f, 0x6c, 0xb6, 0x76, 0x17, 0xb2, 0x0d, 0xbc, 0x6c,
	0x7a, 0x37, 0xcb, 0x1f, 0x99, 0x6a, 0x01, 0x43, 0x22, 0x31, 0x21, 0x08, 0x04, 0x51, 0xa3, 0x06,
	0xe3, 0x85, 0x8b, 0xe9, 0x19, 0x8a, 0x9e, 0xce, 0x0c, 0x5d, 0x63, 0x57, 0x0d, 0x38, 0x12, 0x22,
	0xf1, 0x44, 0x3c, 0x18, 0xa3, 0x17, 0x13, 0x0f, 0x1a, 0x4f, 0x98, 0x78, 0xe0, 0x62, 0xc2, 0x91,
	0x23, 0x47, 0x12, 0x2f, 0x1c, 0x99, 0x6e, 0x0e, 0x1c, 0x49, 0xbc, 0x70, 0x34, 0x53, 0xd3, 0x23,
	0xdd, 0x3d, 0x93, 0x19, 0x46, 0xbd, 0x75, 0x77, 0x7d, 0xf5, 0xde, 0x6f, 0xbe, 0x57, 0x5f, 0x8d,
	0x96, 0xca, 0x73, 0xcf, 0x5a, 0xb5, 0xdc, 0x94, 0x90, 0x56, 0x26, 0x67, 0x5a, 0x05, 0xc7, 0xe4,
	0x9e, 0x6d, 0xb9, 0xce, 0x23, 0x4b, 0x3a, 0xdc, 0xbd, 0x2f, 0x98, 0xb7, 0xe2, 0x64, 0x98, 0xa0,
	0x05, 0x8f, 0x4b, 0x4e, 0x3a, 0xa4, 0x74, 0x69, 0xb8, 0x85, 0xae, 0x8c, 0xea, 0x29, 0xdb, 0x91,
	0xd9, 0x62, 0x9a, 0x66, 0xf8, 0xb2, 0x69, 0x73, 0x9b, 0x9b, 0x4a, 0x96, 0x2e, 0x2e, 0xa9, 0x37,
	0xf5, 0xa2, 0x9e, 0xaa, 0xdb, 0xf5, 0x1e, 0x9b, 0x73, 0x3b, 0xcf, 0x54, 0x1b, 0xcb, 0x75, 0xb9,
	0x54, 0x4d, 0xc2, 0xe2, 0x7a, 0x77, 0xb8, 0xfa, 0xb5, 0x06, 0x5b, 0x2e, 0xc8, 0x52, 0xb8, 0xf8,
	0x6f, 0x3d, 0xa8, 0xb3, 0xc8, 0x5c, 0xe9, 0x2c, 0x39, 0xcc, 0xab, 0x55, 0xf8, 0xaf, 0xf9, 0xaf,
	0x09, 0x55, 0x50, 0xaf, 0xf2, 0x1c, 0x3b, 0x2b, 0xc3, 0x2a, 0x23, 0x9f, 0x7f, 0xd6, 0xfe, 0xbc,
	0x1d, 0xd9, 0x36, 0xcf, 0x6c, 0x47, 0x48, 0xaf, 0x44, 0x1e, 0x6a, 0x64, 0xca, 0x63, 0x96, 0x64,
	0xd1, 0x55, 0x32, 0x40, 0xe3, 0xa6, 0xd0, 0x7a, 0xcd, 0x3c, 0x7b, 0x50, 0x64, 0x42, 0xea, 0x3d,
	0x49, 0x69, 0x54, 0x64, 0xfc, 0xfd, 0xe4, 0xd3, 0xd1, 0xcb, 0x9f, 0xfe, 0x30, 0x3a, 0x62, 0xd0,
	0x62, 0x1c, 0x0f, 0x92, 0x17, 0x58, 0xfb, 0x6d, 0x96, 0xc9, 0x58, 0xdf, 0xff, 0x93, 0xc5, 0x12,
	0x82, 0xf3, 0x35, 0xbd, 0xac, 0x9a, 0x8e, 0x92, 0xe1, 0x78, 0x53, 0x73, 0x2d, 0x76, 0x0c, 0x9c,
	0x45, 0x41, 0x13, 0x1f, 0xd6, 0xc9, 0x07, 0xac, 0xfd, 0x7e, 0xd3, 0x11, 0xb1, 0xa6, 0x82, 0xf4,
	0x27, 0xdb, 0xd5, 0x49, 0x6a, 0x60, 0xbd, 0xcd, 0xc0, 0x84, 0x71, 0x4b, 0x91, 0x5d, 0x23, 0x09,
	0x3b, 0x16, 0xc6, 0xc8, 0x25, 0xb3, 0x28, 0x98, 0x27, 0xcc, 0xb5, 0x0c, 0xcf, 0xe7, 0xad, 0x34,
	0xf7, 0x2c, 0xc9, 0x3d, 0x5a, 0xf9, 0xa6, 0x40, 0xc3, 0x87, 0xf5, 0xf8, 0x3e, 0xf2, 0x1a, 0x6b,
	0xe4, 0x5e, 0x61, 0xb1, 0xe5, 0xf8, 0xea, 0x35, 0xe7, 0x73, 0xf2, 0x8a, 0xe2, 0x1d, 0xd3, 0x9b,
	0x3a, 0x49, 0x1b, 0x39, 0x59, 0x99, 0xf0, 0x63, 0x8d, 0x4c, 0xb3, 0x3c, 0x4b, 0xc0, 0xf5, 0x35,
	0xeb, 0x38, 0x77, 0x76, 0xfe, 0xf5, 0x4e, 0x5a, 0x0d, 0x0f, 0xad, 0x85, 0x87, 0xce, 0x54, 0xc2,
	0x63, 0xf4, 0x2b, 0x28, 0x63, 0xf0, 0x9f, 0x16, 0xe3, 0x5d, 0x1f, 0x39, 0xfa, 0x55, 0x23, 0xd1,
	0xea, 0x93, 0x99, 0x0c, 0x13, 0x82, 0x3c, 0xc5, 0x5a, 0x67, 0x72, 0x82, 0xf3, 0x2a, 0x2d, 0xed,
	0xc0, 0x25, 0x84, 0xd5, 0x02, 0x86, 0xa9, 0xe0, 0x06, 0x48, 0x5f, 0x2b, 0xb8, 0x30, 0x9f, 0xe4,
	0x3d, 0xd6, 0xf4, 0x59, 0xe6, 0x32, 0x2f, 0x31, 0xa0, 0xc9, 0x3b, 0x73, 0x37, 0x58, 0x89, 0x0c,
	0x25, 0xfb, 0xdc, 0x8d, 0x27, 0xa2, 0x2a, 0xab, 0x4d, 0xb3, 0x8e, 0xaa, 0xba, 0x6c, 0xcc, 0x28,
	0xaa, 0x09, 0x63, 0xbc, 0xed, 0x44, 0x54, 0xae, 0x92, 0x54, 0x8e, 0x95, 0x54, 0x64, 0x9f, 0x61,
	0xad, 0x2b, 0x69, 0x5c, 0xb5, 0x43, 0x1b, 0xce, 0x75, 0x35, 0x66, 0x14, 0xc6, 0xb0, 0x82, 0xbc,
	0x40, 0x06, 0x5a, 0x5a, 0x57, 0x63, 0xaa, 0x98, 0xf7, 0x57, 0x43, 0x43, 0x7e, 0x90, 0x6f, 0xd7,
	0x15, 0xd2, 0xb4, 0x3e, 0xf1, 0xed, 0xbe, 0x99, 0x6b, 0x61, 0x1a, 0x3e, 0x62, 0xad, 0x3b, 0x01,
	0x31, 0x15, 0xc9, 0x3a, 0x19, 0x69, 0x41, 0x1c, 0x15, 0x9f, 0x71, 0x37, 0x8e, 0xc8, 0x77, 0x70,
	0x47, 0xef, 0x1f, 0x35, 0xf4, 0x37, 0x58, 0xeb, 0x4d, 0x0e, 0x3d, 0xca, 0xd2, 0xc6, 0xe8, 0xeb,
	0x6e, 0xc7, 0x58, 0x1d, 0x63, 0x4c, 0x51, 0x5f, 0x24, 0xb4, 0xe5, 0x01, 0x88, 0x41, 0x5e, 0x7d,
	0x87, 0xf7, 0xca, 0x80, 0xf7, 0xcb, 0x80, 0x0f, 0xca, 0x80, 0x0e, 0xcb, 0x80, 0x8e, 0xcb, 0x80,
	0x4e, 0xca, 0x80, 0x4e, 0xcb, 0x80, 0x37, 0x7c, 0xc0, 0x9b, 0x3e, 0xa0, 0x2d, 0x1f, 0xf0, 0xb6,
	0x0f, 0x68, 0xc7, 0x07, 0xb4, 0xeb, 0x03, 0xda, 0xf3, 0x01, 0xef, 0xfb, 0x80, 0x0f, 0x7c, 0x40,
	0x87, 0x3e, 0xe0, 0x63, 0x1f, 0xd0, 0x89, 0x0f, 0xf8, 0xd4, 0x07, 0xb4, 0x11, 0x00, 0xda, 0x0c,
	0x00, 0x3f, 0x0f, 0x00, 0xbd, 0x0a, 0x00, 0xbf, 0x0d, 0x00, 0x6d, 0x05, 0x80, 0xb6, 0x03, 0xc0,
	0x3b, 0x01, 0xe0, 0xdd, 0x00, 0xf0, 0xc2, 0x90, 0xcd, 0xa9, 0xcc, 0x32, 0x99, 0x75, 0x5c, 0x5b,
	0x50, 0x97, 0xc9, 0x55, 0xee, 0xe5, 0xcc, 0xf8, 0x1f, 0x71, 0x21, 0x67, 0x9b, 0x52, 0xba, 0x85,
	0x74, 0xfa, 0x17, 0x35, 0xa2, 0xd1, 0x2f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x50, 0xe8, 0x9d, 0x7e,
	0x9e, 0x08, 0x00, 0x00,
}
