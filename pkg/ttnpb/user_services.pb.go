// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/user_services.proto

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

// Client API for UserRegistry service

type UserRegistryClient interface {
	// Register a new user. This method may be restricted by network settings.
	Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error)
	// Get the user with the given identifiers, selecting the fields given by the
	// field mask. The method may return more or less fields, depending on the rights
	// of the caller.
	Get(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error)
	Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error)
	// Create a temporary password that can be used for updating a forgotten password.
	// The generated password is sent to the user's email address.
	CreateTemporaryPassword(ctx context.Context, in *CreateTemporaryPasswordRequest, opts ...grpc.CallOption) (*types.Empty, error)
	UpdatePassword(ctx context.Context, in *UpdateUserPasswordRequest, opts ...grpc.CallOption) (*types.Empty, error)
	Delete(ctx context.Context, in *UserIdentifiers, opts ...grpc.CallOption) (*types.Empty, error)
}

type userRegistryClient struct {
	cc *grpc.ClientConn
}

func NewUserRegistryClient(cc *grpc.ClientConn) UserRegistryClient {
	return &userRegistryClient{cc}
}

func (c *userRegistryClient) Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserRegistry/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRegistryClient) Get(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserRegistry/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRegistryClient) Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserRegistry/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRegistryClient) CreateTemporaryPassword(ctx context.Context, in *CreateTemporaryPasswordRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserRegistry/CreateTemporaryPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRegistryClient) UpdatePassword(ctx context.Context, in *UpdateUserPasswordRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserRegistry/UpdatePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRegistryClient) Delete(ctx context.Context, in *UserIdentifiers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserRegistry/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserRegistry service

type UserRegistryServer interface {
	// Register a new user. This method may be restricted by network settings.
	Create(context.Context, *CreateUserRequest) (*User, error)
	// Get the user with the given identifiers, selecting the fields given by the
	// field mask. The method may return more or less fields, depending on the rights
	// of the caller.
	Get(context.Context, *GetUserRequest) (*User, error)
	Update(context.Context, *UpdateUserRequest) (*User, error)
	// Create a temporary password that can be used for updating a forgotten password.
	// The generated password is sent to the user's email address.
	CreateTemporaryPassword(context.Context, *CreateTemporaryPasswordRequest) (*types.Empty, error)
	UpdatePassword(context.Context, *UpdateUserPasswordRequest) (*types.Empty, error)
	Delete(context.Context, *UserIdentifiers) (*types.Empty, error)
}

func RegisterUserRegistryServer(s *grpc.Server, srv UserRegistryServer) {
	s.RegisterService(&_UserRegistry_serviceDesc, srv)
}

func _UserRegistry_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRegistryServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserRegistry/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRegistryServer).Create(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRegistry_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRegistryServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserRegistry/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRegistryServer).Get(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRegistry_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRegistryServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserRegistry/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRegistryServer).Update(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRegistry_CreateTemporaryPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTemporaryPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRegistryServer).CreateTemporaryPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserRegistry/CreateTemporaryPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRegistryServer).CreateTemporaryPassword(ctx, req.(*CreateTemporaryPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRegistry_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRegistryServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserRegistry/UpdatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRegistryServer).UpdatePassword(ctx, req.(*UpdateUserPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRegistry_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRegistryServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserRegistry/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRegistryServer).Delete(ctx, req.(*UserIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.UserRegistry",
	HandlerType: (*UserRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _UserRegistry_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _UserRegistry_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _UserRegistry_Update_Handler,
		},
		{
			MethodName: "CreateTemporaryPassword",
			Handler:    _UserRegistry_CreateTemporaryPassword_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _UserRegistry_UpdatePassword_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _UserRegistry_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/user_services.proto",
}

// Client API for UserAccess service

type UserAccessClient interface {
	ListRights(ctx context.Context, in *UserIdentifiers, opts ...grpc.CallOption) (*Rights, error)
	CreateAPIKey(ctx context.Context, in *CreateUserAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	ListAPIKeys(ctx context.Context, in *UserIdentifiers, opts ...grpc.CallOption) (*APIKeys, error)
	// Update the rights of an existing user API key. To generate an API key,
	// the CreateAPIKey should be used. To delete an API key, update it
	// with zero rights.
	UpdateAPIKey(ctx context.Context, in *UpdateUserAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
}

type userAccessClient struct {
	cc *grpc.ClientConn
}

func NewUserAccessClient(cc *grpc.ClientConn) UserAccessClient {
	return &userAccessClient{cc}
}

func (c *userAccessClient) ListRights(ctx context.Context, in *UserIdentifiers, opts ...grpc.CallOption) (*Rights, error) {
	out := new(Rights)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserAccess/ListRights", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAccessClient) CreateAPIKey(ctx context.Context, in *CreateUserAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserAccess/CreateAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAccessClient) ListAPIKeys(ctx context.Context, in *UserIdentifiers, opts ...grpc.CallOption) (*APIKeys, error) {
	out := new(APIKeys)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserAccess/ListAPIKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAccessClient) UpdateAPIKey(ctx context.Context, in *UpdateUserAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserAccess/UpdateAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserAccess service

type UserAccessServer interface {
	ListRights(context.Context, *UserIdentifiers) (*Rights, error)
	CreateAPIKey(context.Context, *CreateUserAPIKeyRequest) (*APIKey, error)
	ListAPIKeys(context.Context, *UserIdentifiers) (*APIKeys, error)
	// Update the rights of an existing user API key. To generate an API key,
	// the CreateAPIKey should be used. To delete an API key, update it
	// with zero rights.
	UpdateAPIKey(context.Context, *UpdateUserAPIKeyRequest) (*APIKey, error)
}

func RegisterUserAccessServer(s *grpc.Server, srv UserAccessServer) {
	s.RegisterService(&_UserAccess_serviceDesc, srv)
}

func _UserAccess_ListRights_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAccessServer).ListRights(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserAccess/ListRights",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAccessServer).ListRights(ctx, req.(*UserIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAccess_CreateAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAccessServer).CreateAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserAccess/CreateAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAccessServer).CreateAPIKey(ctx, req.(*CreateUserAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAccess_ListAPIKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAccessServer).ListAPIKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserAccess/ListAPIKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAccessServer).ListAPIKeys(ctx, req.(*UserIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAccess_UpdateAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAccessServer).UpdateAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserAccess/UpdateAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAccessServer).UpdateAPIKey(ctx, req.(*UpdateUserAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserAccess_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.UserAccess",
	HandlerType: (*UserAccessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRights",
			Handler:    _UserAccess_ListRights_Handler,
		},
		{
			MethodName: "CreateAPIKey",
			Handler:    _UserAccess_CreateAPIKey_Handler,
		},
		{
			MethodName: "ListAPIKeys",
			Handler:    _UserAccess_ListAPIKeys_Handler,
		},
		{
			MethodName: "UpdateAPIKey",
			Handler:    _UserAccess_UpdateAPIKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/user_services.proto",
}

// Client API for UserInvitationRegistry service

type UserInvitationRegistryClient interface {
	Send(ctx context.Context, in *SendInvitationRequest, opts ...grpc.CallOption) (*Invitation, error)
	List(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*Invitations, error)
	Delete(ctx context.Context, in *DeleteInvitationRequest, opts ...grpc.CallOption) (*types.Empty, error)
}

type userInvitationRegistryClient struct {
	cc *grpc.ClientConn
}

func NewUserInvitationRegistryClient(cc *grpc.ClientConn) UserInvitationRegistryClient {
	return &userInvitationRegistryClient{cc}
}

func (c *userInvitationRegistryClient) Send(ctx context.Context, in *SendInvitationRequest, opts ...grpc.CallOption) (*Invitation, error) {
	out := new(Invitation)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserInvitationRegistry/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userInvitationRegistryClient) List(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*Invitations, error) {
	out := new(Invitations)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserInvitationRegistry/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userInvitationRegistryClient) Delete(ctx context.Context, in *DeleteInvitationRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UserInvitationRegistry/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserInvitationRegistry service

type UserInvitationRegistryServer interface {
	Send(context.Context, *SendInvitationRequest) (*Invitation, error)
	List(context.Context, *types.Empty) (*Invitations, error)
	Delete(context.Context, *DeleteInvitationRequest) (*types.Empty, error)
}

func RegisterUserInvitationRegistryServer(s *grpc.Server, srv UserInvitationRegistryServer) {
	s.RegisterService(&_UserInvitationRegistry_serviceDesc, srv)
}

func _UserInvitationRegistry_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendInvitationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserInvitationRegistryServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserInvitationRegistry/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserInvitationRegistryServer).Send(ctx, req.(*SendInvitationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserInvitationRegistry_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserInvitationRegistryServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserInvitationRegistry/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserInvitationRegistryServer).List(ctx, req.(*types.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserInvitationRegistry_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteInvitationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserInvitationRegistryServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UserInvitationRegistry/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserInvitationRegistryServer).Delete(ctx, req.(*DeleteInvitationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserInvitationRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.UserInvitationRegistry",
	HandlerType: (*UserInvitationRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _UserInvitationRegistry_Send_Handler,
		},
		{
			MethodName: "List",
			Handler:    _UserInvitationRegistry_List_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _UserInvitationRegistry_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/user_services.proto",
}

func init() {
	proto.RegisterFile("lorawan-stack/api/user_services.proto", fileDescriptor_user_services_3c18ee3daf352a55)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/user_services.proto", fileDescriptor_user_services_3c18ee3daf352a55)
}

var fileDescriptor_user_services_3c18ee3daf352a55 = []byte{
	// 762 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x55, 0x3f, 0x4c, 0xeb, 0x46,
	0x18, 0xf7, 0xbd, 0xf7, 0x9a, 0xe1, 0x5e, 0x84, 0xfa, 0x4e, 0x4f, 0x2f, 0xc5, 0xa0, 0x0f, 0xe1,
	0x16, 0xa1, 0x22, 0x72, 0x96, 0x80, 0x89, 0x8d, 0xfe, 0x11, 0x42, 0xed, 0x40, 0x29, 0x2c, 0xad,
	0xd4, 0xc8, 0x49, 0x0e, 0xe7, 0x1a, 0xb0, 0x5d, 0xdf, 0x05, 0x14, 0x21, 0x2a, 0xc4, 0x84, 0xd4,
	0xa5, 0x52, 0x87, 0x76, 0xac, 0x3a, 0xd1, 0x8d, 0x91, 0x91, 0x91, 0x11, 0xa9, 0x0b, 0x23, 0xb1,
	0x3b, 0x30, 0x55, 0x8c, 0x8c, 0x95, 0xef, 0xe2, 0x84, 0x38, 0xf1, 0x23, 0x9b, 0x7d, 0xdf, 0xef,
	0x7e, 0x7f, 0xbe, 0xbb, 0xcf, 0xc6, 0x73, 0x7b, 0x7e, 0xe8, 0x1c, 0x3a, 0x5e, 0x59, 0x48, 0xa7,
	0xd6, 0xb4, 0x9d, 0x80, 0xdb, 0x2d, 0xc1, 0xc2, 0x8a, 0x60, 0xe1, 0x01, 0xaf, 0x31, 0x41, 0x83,
	0xd0, 0x97, 0x3e, 0x99, 0x90, 0xd2, 0xa3, 0x5d, 0x28, 0x3d, 0x58, 0x36, 0xcb, 0x2e, 0x97, 0x8d,
	0x56, 0x95, 0xd6, 0xfc, 0x7d, 0xdb, 0xf5, 0x5d, 0xdf, 0x56, 0xb0, 0x6a, 0x6b, 0x57, 0xbd, 0xa9,
	0x17, 0xf5, 0xa4, 0xb7, 0x9b, 0xd3, 0xae, 0xef, 0xbb, 0x7b, 0x4c, 0xd1, 0x3b, 0x9e, 0xe7, 0x4b,
	0x47, 0x72, 0xdf, 0xeb, 0x92, 0x9b, 0x53, 0xdd, 0x6a, 0x8f, 0x83, 0xed, 0x07, 0xb2, 0xdd, 0x2d,
	0x7e, 0x3c, 0x6c, 0x90, 0xd7, 0x99, 0x27, 0xf9, 0x2e, 0x67, 0x61, 0xca, 0x00, 0xc3, 0xa0, 0x90,
	0xbb, 0x0d, 0x99, 0xd6, 0xa7, 0x47, 0xa7, 0xd4, 0xd5, 0xa5, 0xbf, 0x3f, 0xc0, 0xc5, 0x1d, 0xc1,
	0xc2, 0x2d, 0xe6, 0x72, 0x21, 0xc3, 0x36, 0xd9, 0xc6, 0x85, 0xcf, 0x43, 0xe6, 0x48, 0x46, 0x66,
	0xe9, 0x60, 0x70, 0xaa, 0xd7, 0x35, 0xfa, 0xa7, 0x16, 0x13, 0xd2, 0x7c, 0x9b, 0x85, 0x24, 0x45,
	0xeb, 0xcd, 0xe9, 0x3f, 0xff, 0xfe, 0xf6, 0xe2, 0xb5, 0x55, 0x50, 0x42, 0x62, 0x15, 0x2d, 0x90,
	0x1f, 0xf0, 0xcb, 0x75, 0x26, 0x09, 0x64, 0xf1, 0xeb, 0x4c, 0x3e, 0xcf, 0x37, 0xab, 0xf8, 0xa6,
	0xc8, 0xa4, 0xe6, 0xb3, 0x8f, 0xd4, 0x29, 0xf1, 0xba, 0xa0, 0xdd, 0x87, 0x63, 0xe2, 0xe2, 0xc2,
	0x4e, 0x50, 0x1f, 0xe9, 0x5a, 0xaf, 0x3f, 0xaf, 0xf2, 0x89, 0x52, 0x01, 0x73, 0x40, 0x85, 0x3e,
	0x55, 0x49, 0x82, 0xfc, 0x8e, 0x70, 0x49, 0xf7, 0x61, 0x9b, 0xed, 0x07, 0x7e, 0xe8, 0x84, 0xed,
	0x4d, 0x47, 0x88, 0x43, 0x3f, 0xac, 0x13, 0x3a, 0xba, 0x61, 0x43, 0xc0, 0xd4, 0xc7, 0x3b, 0xaa,
	0x0f, 0x9f, 0xa6, 0x87, 0x4f, 0xbf, 0x4c, 0x0e, 0xdf, 0x5a, 0x51, 0x4e, 0xa8, 0xb5, 0x98, 0x9b,
	0xd7, 0x96, 0x29, 0x67, 0x25, 0x48, 0xd5, 0x4f, 0x11, 0x9e, 0xd0, 0x59, 0x7b, 0x86, 0x3e, 0xcd,
	0xef, 0xc5, 0xb8, 0x5e, 0xca, 0xca, 0xcb, 0xbc, 0x69, 0xe5, 0x7b, 0x49, 0x1d, 0x24, 0xed, 0xf9,
	0x1e, 0x17, 0xbe, 0x60, 0x7b, 0x4c, 0x32, 0x32, 0x33, 0xaa, 0xc9, 0x1b, 0xfd, 0xdb, 0x9b, 0xab,
	0xf8, 0x91, 0x52, 0x24, 0x0b, 0x1f, 0x66, 0x14, 0x8f, 0x97, 0xfe, 0x7b, 0x89, 0x71, 0xc2, 0xb2,
	0x56, 0xab, 0x31, 0x21, 0xc8, 0x2e, 0xc6, 0x5f, 0x73, 0x21, 0xb7, 0xd4, 0x65, 0x1f, 0x47, 0x2f,
	0x03, 0xd0, 0x1b, 0xad, 0x19, 0xa5, 0x37, 0x49, 0x4a, 0x59, 0xbd, 0xee, 0x18, 0x91, 0x9f, 0x71,
	0x51, 0x1f, 0xe4, 0xda, 0xe6, 0xc6, 0x57, 0xac, 0x4d, 0xe6, 0xf3, 0xe7, 0x42, 0x23, 0xfa, 0x3d,
	0xcd, 0x00, 0x75, 0x39, 0xed, 0xa9, 0xf5, 0x9e, 0x9e, 0x3a, 0x01, 0x2f, 0x37, 0x59, 0x5b, 0xcd,
	0xce, 0x8f, 0xf8, 0x75, 0x92, 0x53, 0x6f, 0x1e, 0x23, 0x68, 0x69, 0xb4, 0xac, 0xc8, 0x9d, 0xa3,
	0xbe, 0x1c, 0xf9, 0x05, 0xe1, 0xa2, 0xbe, 0x24, 0x79, 0x61, 0xfb, 0x57, 0x68, 0xbc, 0xb0, 0xab,
	0x4a, 0x74, 0xc5, 0xb4, 0x9f, 0x0f, 0x6b, 0x1f, 0x39, 0x01, 0xaf, 0x34, 0x59, 0x9b, 0xea, 0x61,
	0x5b, 0xba, 0x78, 0x81, 0xdf, 0xa9, 0x74, 0xde, 0x01, 0xd7, 0x9f, 0xcd, 0xde, 0x67, 0xaa, 0x8a,
	0x5f, 0x7d, 0xcb, 0xbc, 0x3a, 0x99, 0xcb, 0xca, 0x26, 0xab, 0x4f, 0xf1, 0xda, 0x9d, 0x99, 0x85,
	0xf5, 0x21, 0x56, 0x49, 0x39, 0x7c, 0x63, 0x15, 0x6d, 0xde, 0x5b, 0x54, 0x8d, 0xff, 0x06, 0xbf,
	0x4a, 0x1a, 0x4f, 0x72, 0x6e, 0xaa, 0x39, 0x95, 0x4f, 0x2a, 0xac, 0xb7, 0x8a, 0x75, 0x82, 0x0c,
	0xb0, 0x92, 0x4a, 0x6f, 0x3e, 0x86, 0x1a, 0xab, 0xd7, 0x87, 0xad, 0xe7, 0xcd, 0x49, 0x57, 0x60,
	0x61, 0x40, 0xe0, 0xb3, 0xbf, 0xd0, 0x75, 0x07, 0xd0, 0x4d, 0x07, 0xd0, 0x6d, 0x07, 0x8c, 0xbb,
	0x0e, 0x18, 0xf7, 0x1d, 0x30, 0x1e, 0x3a, 0x60, 0x3c, 0x76, 0x00, 0x9d, 0x44, 0x80, 0xce, 0x22,
	0x30, 0xce, 0x23, 0x40, 0x17, 0x11, 0x18, 0x97, 0x11, 0x18, 0x57, 0x11, 0x18, 0xd7, 0x11, 0xa0,
	0x9b, 0x08, 0xd0, 0x6d, 0x04, 0xc6, 0x5d, 0x04, 0xe8, 0x3e, 0x02, 0xe3, 0x21, 0x02, 0xf4, 0x18,
	0x81, 0x71, 0x12, 0x83, 0x71, 0x16, 0x03, 0xfa, 0x35, 0x06, 0xe3, 0x8f, 0x18, 0xd0, 0x9f, 0x31,
	0x18, 0xe7, 0x31, 0x18, 0x17, 0x31, 0xa0, 0xcb, 0x18, 0xd0, 0x55, 0x0c, 0xe8, 0xbb, 0x45, 0xd7,
	0xa7, 0xb2, 0xc1, 0x64, 0x83, 0x7b, 0xae, 0xa0, 0x1e, 0x93, 0x87, 0x7e, 0xd8, 0xb4, 0x07, 0xff,
	0x3c, 0x41, 0xd3, 0xb5, 0xa5, 0xf4, 0x82, 0x6a, 0xb5, 0xa0, 0xa2, 0x2c, 0xff, 0x1f, 0x00, 0x00,
	0xff, 0xff, 0x70, 0xdc, 0xf6, 0x62, 0x81, 0x07, 0x00, 0x00,
}
