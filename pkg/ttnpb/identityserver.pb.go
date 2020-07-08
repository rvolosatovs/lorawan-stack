// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/identityserver.proto

package ttnpb

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type AuthInfoResponse struct {
	// Types that are valid to be assigned to AccessMethod:
	//	*AuthInfoResponse_APIKey
	//	*AuthInfoResponse_OAuthAccessToken
	//	*AuthInfoResponse_UserSession
	AccessMethod         isAuthInfoResponse_AccessMethod `protobuf_oneof:"access_method"`
	UniversalRights      *Rights                         `protobuf:"bytes,3,opt,name=universal_rights,json=universalRights,proto3" json:"universal_rights,omitempty"`
	IsAdmin              bool                            `protobuf:"varint,4,opt,name=is_admin,json=isAdmin,proto3" json:"is_admin,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *AuthInfoResponse) Reset()      { *m = AuthInfoResponse{} }
func (*AuthInfoResponse) ProtoMessage() {}
func (*AuthInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1c7e02f6181562c, []int{0}
}
func (m *AuthInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthInfoResponse.Merge(m, src)
}
func (m *AuthInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *AuthInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AuthInfoResponse proto.InternalMessageInfo

type isAuthInfoResponse_AccessMethod interface {
	isAuthInfoResponse_AccessMethod()
	Equal(interface{}) bool
	MarshalTo([]byte) (int, error)
	Size() int
}

type AuthInfoResponse_APIKey struct {
	APIKey *AuthInfoResponse_APIKeyAccess `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3,oneof" json:"api_key,omitempty"`
}
type AuthInfoResponse_OAuthAccessToken struct {
	OAuthAccessToken *OAuthAccessToken `protobuf:"bytes,2,opt,name=oauth_access_token,json=oauthAccessToken,proto3,oneof" json:"oauth_access_token,omitempty"`
}
type AuthInfoResponse_UserSession struct {
	UserSession *UserSession `protobuf:"bytes,5,opt,name=user_session,json=userSession,proto3,oneof" json:"user_session,omitempty"`
}

func (*AuthInfoResponse_APIKey) isAuthInfoResponse_AccessMethod()           {}
func (*AuthInfoResponse_OAuthAccessToken) isAuthInfoResponse_AccessMethod() {}
func (*AuthInfoResponse_UserSession) isAuthInfoResponse_AccessMethod()      {}

func (m *AuthInfoResponse) GetAccessMethod() isAuthInfoResponse_AccessMethod {
	if m != nil {
		return m.AccessMethod
	}
	return nil
}

func (m *AuthInfoResponse) GetAPIKey() *AuthInfoResponse_APIKeyAccess {
	if x, ok := m.GetAccessMethod().(*AuthInfoResponse_APIKey); ok {
		return x.APIKey
	}
	return nil
}

func (m *AuthInfoResponse) GetOAuthAccessToken() *OAuthAccessToken {
	if x, ok := m.GetAccessMethod().(*AuthInfoResponse_OAuthAccessToken); ok {
		return x.OAuthAccessToken
	}
	return nil
}

func (m *AuthInfoResponse) GetUserSession() *UserSession {
	if x, ok := m.GetAccessMethod().(*AuthInfoResponse_UserSession); ok {
		return x.UserSession
	}
	return nil
}

func (m *AuthInfoResponse) GetUniversalRights() *Rights {
	if m != nil {
		return m.UniversalRights
	}
	return nil
}

func (m *AuthInfoResponse) GetIsAdmin() bool {
	if m != nil {
		return m.IsAdmin
	}
	return false
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*AuthInfoResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*AuthInfoResponse_APIKey)(nil),
		(*AuthInfoResponse_OAuthAccessToken)(nil),
		(*AuthInfoResponse_UserSession)(nil),
	}
}

type AuthInfoResponse_APIKeyAccess struct {
	APIKey               `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3,embedded=api_key" json:"api_key"`
	EntityIDs            EntityIdentifiers `protobuf:"bytes,2,opt,name=entity_ids,json=entityIds,proto3" json:"entity_ids"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AuthInfoResponse_APIKeyAccess) Reset()      { *m = AuthInfoResponse_APIKeyAccess{} }
func (*AuthInfoResponse_APIKeyAccess) ProtoMessage() {}
func (*AuthInfoResponse_APIKeyAccess) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1c7e02f6181562c, []int{0, 0}
}
func (m *AuthInfoResponse_APIKeyAccess) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthInfoResponse_APIKeyAccess) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthInfoResponse_APIKeyAccess.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthInfoResponse_APIKeyAccess) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthInfoResponse_APIKeyAccess.Merge(m, src)
}
func (m *AuthInfoResponse_APIKeyAccess) XXX_Size() int {
	return m.Size()
}
func (m *AuthInfoResponse_APIKeyAccess) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthInfoResponse_APIKeyAccess.DiscardUnknown(m)
}

var xxx_messageInfo_AuthInfoResponse_APIKeyAccess proto.InternalMessageInfo

func (m *AuthInfoResponse_APIKeyAccess) GetEntityIDs() EntityIdentifiers {
	if m != nil {
		return m.EntityIDs
	}
	return EntityIdentifiers{}
}

func init() {
	proto.RegisterType((*AuthInfoResponse)(nil), "ttn.lorawan.v3.AuthInfoResponse")
	golang_proto.RegisterType((*AuthInfoResponse)(nil), "ttn.lorawan.v3.AuthInfoResponse")
	proto.RegisterType((*AuthInfoResponse_APIKeyAccess)(nil), "ttn.lorawan.v3.AuthInfoResponse.APIKeyAccess")
	golang_proto.RegisterType((*AuthInfoResponse_APIKeyAccess)(nil), "ttn.lorawan.v3.AuthInfoResponse.APIKeyAccess")
}

func init() {
	proto.RegisterFile("lorawan-stack/api/identityserver.proto", fileDescriptor_a1c7e02f6181562c)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/identityserver.proto", fileDescriptor_a1c7e02f6181562c)
}

var fileDescriptor_a1c7e02f6181562c = []byte{
	// 683 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0x3d, 0x6b, 0x1b, 0x4b,
	0x14, 0x9d, 0xf1, 0xf3, 0xe7, 0x5a, 0xef, 0x59, 0x2c, 0x0f, 0x23, 0xcb, 0x7e, 0x57, 0x7a, 0x0e,
	0x04, 0x37, 0xda, 0x05, 0xfb, 0x0f, 0x44, 0x22, 0x06, 0x1b, 0x17, 0x09, 0x1b, 0x07, 0x42, 0x52,
	0x88, 0x95, 0x34, 0xda, 0x1d, 0x24, 0xcd, 0x2c, 0x3b, 0x23, 0x39, 0xea, 0x4c, 0x2a, 0x93, 0x2a,
	0x90, 0x26, 0x65, 0x9a, 0x80, 0xab, 0x60, 0x52, 0xb9, 0x74, 0xe9, 0xd2, 0x90, 0xc6, 0x95, 0xb0,
	0x66, 0x53, 0xb8, 0x74, 0x69, 0x52, 0x05, 0xcd, 0xae, 0x6c, 0x59, 0x22, 0xa4, 0xdb, 0x7b, 0xcf,
	0x99, 0x73, 0xe7, 0x9e, 0x3d, 0x63, 0x3c, 0x6e, 0xf2, 0xd0, 0x3d, 0x70, 0x59, 0x41, 0x48, 0xb7,
	0xda, 0xb0, 0xdd, 0x80, 0xda, 0xb4, 0x46, 0x98, 0xa4, 0xb2, 0x2b, 0x48, 0xd8, 0x21, 0xa1, 0x15,
	0x84, 0x5c, 0x72, 0xf3, 0x1f, 0x29, 0x99, 0x95, 0x70, 0xad, 0xce, 0x56, 0xb6, 0xe8, 0x51, 0xe9,
	0xb7, 0x2b, 0x56, 0x95, 0xb7, 0x6c, 0xc2, 0x3a, 0xbc, 0x1b, 0x84, 0xfc, 0x6d, 0xd7, 0xd6, 0xe4,
	0x6a, 0xc1, 0x23, 0xac, 0xd0, 0x71, 0x9b, 0xb4, 0xe6, 0x4a, 0x62, 0x4f, 0x7c, 0xc4, 0x92, 0xd9,
	0xc2, 0x88, 0x84, 0xc7, 0x3d, 0x1e, 0x1f, 0xae, 0xb4, 0xeb, 0xba, 0xd2, 0x85, 0xfe, 0x4a, 0xe8,
	0x6b, 0x1e, 0xe7, 0x5e, 0x93, 0xe8, 0x2b, 0xba, 0x8c, 0x71, 0xe9, 0x4a, 0xca, 0x99, 0x48, 0xd0,
	0xd5, 0x04, 0xbd, 0xd3, 0x20, 0xad, 0x40, 0x76, 0x13, 0xf0, 0xd1, 0xef, 0x96, 0xac, 0x53, 0x12,
	0x0e, 0x15, 0xd6, 0x26, 0x49, 0x6d, 0x31, 0xdc, 0x3f, 0xfb, 0xdf, 0x24, 0xca, 0xdd, 0xb6, 0xf4,
	0x13, 0x18, 0x26, 0xe1, 0x90, 0x7a, 0xbe, 0x4c, 0xc4, 0xd7, 0xbf, 0x4d, 0x1b, 0xe9, 0x62, 0x5b,
	0xfa, 0xbb, 0xac, 0xce, 0x1d, 0x22, 0x02, 0xce, 0x04, 0x31, 0xf7, 0x8d, 0x39, 0x37, 0xa0, 0xe5,
	0x06, 0xe9, 0x66, 0x70, 0x1e, 0x6f, 0x2c, 0x6e, 0x16, 0xac, 0x87, 0x2e, 0x5b, 0xe3, 0x47, 0xac,
	0xe2, 0xf3, 0xdd, 0x3d, 0xd2, 0x2d, 0x56, 0xab, 0x44, 0x88, 0x92, 0xa1, 0x7a, 0xb9, 0xd9, 0xb8,
	0xb3, 0x83, 0x9c, 0x59, 0x37, 0xa0, 0x7b, 0xa4, 0x6b, 0xd6, 0x0d, 0x53, 0xdf, 0xac, 0xec, 0x6a,
	0x56, 0x59, 0xf2, 0x06, 0x61, 0x99, 0x29, 0x3d, 0x20, 0x3f, 0x3e, 0xe0, 0xd9, 0x60, 0x42, 0x2c,
	0xb7, 0x3f, 0xe0, 0x95, 0xfe, 0x55, 0xbd, 0x5c, 0x7a, 0xbc, 0xbb, 0x83, 0x9c, 0xb4, 0xd6, 0x1c,
	0xe9, 0x99, 0x4f, 0x8c, 0xd4, 0xc0, 0x9f, 0xb2, 0x20, 0x42, 0x50, 0xce, 0x32, 0x33, 0x7a, 0xc2,
	0xea, 0xf8, 0x84, 0x97, 0x82, 0x84, 0x2f, 0x62, 0xca, 0x0e, 0x72, 0x16, 0xdb, 0xf7, 0xa5, 0x59,
	0x34, 0xd2, 0x6d, 0x46, 0x3b, 0x24, 0x14, 0x6e, 0xb3, 0x1c, 0xdb, 0x95, 0xf9, 0x4b, 0xab, 0x2c,
	0x8f, 0xab, 0x38, 0x1a, 0x75, 0x96, 0xee, 0xf8, 0x71, 0xc3, 0x5c, 0x31, 0xe6, 0xa9, 0x28, 0xbb,
	0xb5, 0x16, 0x65, 0x99, 0xe9, 0x3c, 0xde, 0x98, 0x77, 0xe6, 0xa8, 0x28, 0x0e, 0xca, 0xec, 0x57,
	0x6c, 0xa4, 0x46, 0xed, 0x32, 0x8b, 0xe3, 0x76, 0x4f, 0x4c, 0x89, 0xe9, 0xa5, 0xf4, 0xcf, 0xd2,
	0xcc, 0x7b, 0x3c, 0x95, 0xc6, 0xe7, 0xbd, 0x1c, 0xba, 0xe8, 0xe5, 0xf0, 0x9d, 0xb7, 0x6f, 0x0c,
	0x23, 0x7e, 0x1b, 0x65, 0x5a, 0x13, 0x89, 0xa7, 0xff, 0x8f, 0xab, 0x6c, 0x6b, 0xc6, 0xee, 0x7d,
	0xc0, 0x4a, 0x2b, 0xa3, 0x82, 0xaa, 0x97, 0x5b, 0x48, 0x28, 0x4f, 0x85, 0xb3, 0x40, 0x12, 0xb6,
	0x28, 0x2d, 0x19, 0x7f, 0x27, 0xbf, 0xac, 0x45, 0xa4, 0xcf, 0x6b, 0x9b, 0xbe, 0x91, 0x8a, 0x89,
	0xc9, 0x02, 0xaf, 0x8c, 0xf9, 0x61, 0x20, 0xcc, 0x65, 0x2b, 0x0e, 0xbc, 0x35, 0x0c, 0xbc, 0xb5,
	0x3d, 0x08, 0x7c, 0x36, 0xff, 0xa7, 0x08, 0xad, 0x9b, 0xef, 0xbe, 0xff, 0xf8, 0x38, 0x95, 0x32,
	0x0d, 0x5b, 0xa7, 0x84, 0xb2, 0x3a, 0x2f, 0x7d, 0xc1, 0xe7, 0x7d, 0xc0, 0x17, 0x7d, 0xc0, 0x97,
	0x7d, 0x40, 0x57, 0x7d, 0x40, 0xd7, 0x7d, 0x40, 0x37, 0x7d, 0x40, 0xb7, 0x7d, 0xc0, 0x87, 0x0a,
	0xf0, 0x91, 0x02, 0x74, 0xac, 0x00, 0x9f, 0x28, 0x40, 0xa7, 0x0a, 0xd0, 0x99, 0x02, 0x74, 0xae,
	0x00, 0x5f, 0x28, 0xc0, 0x97, 0x0a, 0xd0, 0x95, 0x02, 0x7c, 0xad, 0x00, 0xdd, 0x28, 0xc0, 0xb7,
	0x0a, 0xd0, 0x61, 0x04, 0xe8, 0x28, 0x02, 0xfc, 0x21, 0x02, 0xf4, 0x29, 0x02, 0xfc, 0x39, 0x02,
	0x74, 0x1c, 0x01, 0x3a, 0x89, 0x00, 0x9f, 0x46, 0x80, 0xcf, 0x22, 0xc0, 0xaf, 0x6d, 0x8f, 0x5b,
	0xd2, 0x27, 0xd2, 0xa7, 0xcc, 0x13, 0x16, 0x23, 0xf2, 0x80, 0x87, 0x0d, 0xfb, 0xe1, 0x53, 0xea,
	0x6c, 0xd9, 0x41, 0xc3, 0xb3, 0xa5, 0x64, 0x41, 0xa5, 0x32, 0xab, 0xb7, 0xdd, 0xfa, 0x15, 0x00,
	0x00, 0xff, 0xff, 0x32, 0xb2, 0x84, 0xf9, 0xb6, 0x04, 0x00, 0x00,
}

func (this *AuthInfoResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AuthInfoResponse)
	if !ok {
		that2, ok := that.(AuthInfoResponse)
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
	if that1.AccessMethod == nil {
		if this.AccessMethod != nil {
			return false
		}
	} else if this.AccessMethod == nil {
		return false
	} else if !this.AccessMethod.Equal(that1.AccessMethod) {
		return false
	}
	if !this.UniversalRights.Equal(that1.UniversalRights) {
		return false
	}
	if this.IsAdmin != that1.IsAdmin {
		return false
	}
	return true
}
func (this *AuthInfoResponse_APIKey) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AuthInfoResponse_APIKey)
	if !ok {
		that2, ok := that.(AuthInfoResponse_APIKey)
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
	if !this.APIKey.Equal(that1.APIKey) {
		return false
	}
	return true
}
func (this *AuthInfoResponse_OAuthAccessToken) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AuthInfoResponse_OAuthAccessToken)
	if !ok {
		that2, ok := that.(AuthInfoResponse_OAuthAccessToken)
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
	if !this.OAuthAccessToken.Equal(that1.OAuthAccessToken) {
		return false
	}
	return true
}
func (this *AuthInfoResponse_UserSession) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AuthInfoResponse_UserSession)
	if !ok {
		that2, ok := that.(AuthInfoResponse_UserSession)
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
	if !this.UserSession.Equal(that1.UserSession) {
		return false
	}
	return true
}
func (this *AuthInfoResponse_APIKeyAccess) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AuthInfoResponse_APIKeyAccess)
	if !ok {
		that2, ok := that.(AuthInfoResponse_APIKeyAccess)
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
	if !this.APIKey.Equal(&that1.APIKey) {
		return false
	}
	if !this.EntityIDs.Equal(&that1.EntityIDs) {
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

// EntityAccessClient is the client API for EntityAccess service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EntityAccessClient interface {
	// AuthInfo returns information about the authentication that is used on the request.
	AuthInfo(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*AuthInfoResponse, error)
}

type entityAccessClient struct {
	cc *grpc.ClientConn
}

func NewEntityAccessClient(cc *grpc.ClientConn) EntityAccessClient {
	return &entityAccessClient{cc}
}

func (c *entityAccessClient) AuthInfo(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*AuthInfoResponse, error) {
	out := new(AuthInfoResponse)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.EntityAccess/AuthInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EntityAccessServer is the server API for EntityAccess service.
type EntityAccessServer interface {
	// AuthInfo returns information about the authentication that is used on the request.
	AuthInfo(context.Context, *types.Empty) (*AuthInfoResponse, error)
}

// UnimplementedEntityAccessServer can be embedded to have forward compatible implementations.
type UnimplementedEntityAccessServer struct {
}

func (*UnimplementedEntityAccessServer) AuthInfo(ctx context.Context, req *types.Empty) (*AuthInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthInfo not implemented")
}

func RegisterEntityAccessServer(s *grpc.Server, srv EntityAccessServer) {
	s.RegisterService(&_EntityAccess_serviceDesc, srv)
}

func _EntityAccess_AuthInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntityAccessServer).AuthInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.EntityAccess/AuthInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntityAccessServer).AuthInfo(ctx, req.(*types.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _EntityAccess_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.EntityAccess",
	HandlerType: (*EntityAccessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthInfo",
			Handler:    _EntityAccess_AuthInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/identityserver.proto",
}

func (m *AuthInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AccessMethod != nil {
		{
			size := m.AccessMethod.Size()
			i -= size
			if _, err := m.AccessMethod.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	if m.IsAdmin {
		i--
		if m.IsAdmin {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.UniversalRights != nil {
		{
			size, err := m.UniversalRights.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIdentityserver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	return len(dAtA) - i, nil
}

func (m *AuthInfoResponse_APIKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthInfoResponse_APIKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.APIKey != nil {
		{
			size, err := m.APIKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIdentityserver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *AuthInfoResponse_OAuthAccessToken) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthInfoResponse_OAuthAccessToken) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OAuthAccessToken != nil {
		{
			size, err := m.OAuthAccessToken.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIdentityserver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *AuthInfoResponse_UserSession) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthInfoResponse_UserSession) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.UserSession != nil {
		{
			size, err := m.UserSession.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIdentityserver(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	return len(dAtA) - i, nil
}
func (m *AuthInfoResponse_APIKeyAccess) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthInfoResponse_APIKeyAccess) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthInfoResponse_APIKeyAccess) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.EntityIDs.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintIdentityserver(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.APIKey.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintIdentityserver(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintIdentityserver(dAtA []byte, offset int, v uint64) int {
	offset -= sovIdentityserver(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func NewPopulatedAuthInfoResponse(r randyIdentityserver, easy bool) *AuthInfoResponse {
	this := &AuthInfoResponse{}
	oneofNumber_AccessMethod := []int32{1, 2, 5}[r.Intn(3)]
	switch oneofNumber_AccessMethod {
	case 1:
		this.AccessMethod = NewPopulatedAuthInfoResponse_APIKey(r, easy)
	case 2:
		this.AccessMethod = NewPopulatedAuthInfoResponse_OAuthAccessToken(r, easy)
	case 5:
		this.AccessMethod = NewPopulatedAuthInfoResponse_UserSession(r, easy)
	}
	if r.Intn(5) != 0 {
		this.UniversalRights = NewPopulatedRights(r, easy)
	}
	this.IsAdmin = bool(r.Intn(2) == 0)
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedAuthInfoResponse_APIKey(r randyIdentityserver, easy bool) *AuthInfoResponse_APIKey {
	this := &AuthInfoResponse_APIKey{}
	this.APIKey = NewPopulatedAuthInfoResponse_APIKeyAccess(r, easy)
	return this
}
func NewPopulatedAuthInfoResponse_OAuthAccessToken(r randyIdentityserver, easy bool) *AuthInfoResponse_OAuthAccessToken {
	this := &AuthInfoResponse_OAuthAccessToken{}
	this.OAuthAccessToken = NewPopulatedOAuthAccessToken(r, easy)
	return this
}
func NewPopulatedAuthInfoResponse_UserSession(r randyIdentityserver, easy bool) *AuthInfoResponse_UserSession {
	this := &AuthInfoResponse_UserSession{}
	this.UserSession = NewPopulatedUserSession(r, easy)
	return this
}
func NewPopulatedAuthInfoResponse_APIKeyAccess(r randyIdentityserver, easy bool) *AuthInfoResponse_APIKeyAccess {
	this := &AuthInfoResponse_APIKeyAccess{}
	v1 := NewPopulatedAPIKey(r, easy)
	this.APIKey = *v1
	v2 := NewPopulatedEntityIdentifiers(r, easy)
	this.EntityIDs = *v2
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyIdentityserver interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneIdentityserver(r randyIdentityserver) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringIdentityserver(r randyIdentityserver) string {
	v3 := r.Intn(100)
	tmps := make([]rune, v3)
	for i := 0; i < v3; i++ {
		tmps[i] = randUTF8RuneIdentityserver(r)
	}
	return string(tmps)
}
func randUnrecognizedIdentityserver(r randyIdentityserver, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldIdentityserver(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldIdentityserver(dAtA []byte, r randyIdentityserver, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateIdentityserver(dAtA, uint64(key))
		v4 := r.Int63()
		if r.Intn(2) == 0 {
			v4 *= -1
		}
		dAtA = encodeVarintPopulateIdentityserver(dAtA, uint64(v4))
	case 1:
		dAtA = encodeVarintPopulateIdentityserver(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateIdentityserver(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateIdentityserver(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateIdentityserver(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateIdentityserver(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *AuthInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AccessMethod != nil {
		n += m.AccessMethod.Size()
	}
	if m.UniversalRights != nil {
		l = m.UniversalRights.Size()
		n += 1 + l + sovIdentityserver(uint64(l))
	}
	if m.IsAdmin {
		n += 2
	}
	return n
}

func (m *AuthInfoResponse_APIKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.APIKey != nil {
		l = m.APIKey.Size()
		n += 1 + l + sovIdentityserver(uint64(l))
	}
	return n
}
func (m *AuthInfoResponse_OAuthAccessToken) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.OAuthAccessToken != nil {
		l = m.OAuthAccessToken.Size()
		n += 1 + l + sovIdentityserver(uint64(l))
	}
	return n
}
func (m *AuthInfoResponse_UserSession) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.UserSession != nil {
		l = m.UserSession.Size()
		n += 1 + l + sovIdentityserver(uint64(l))
	}
	return n
}
func (m *AuthInfoResponse_APIKeyAccess) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.APIKey.Size()
	n += 1 + l + sovIdentityserver(uint64(l))
	l = m.EntityIDs.Size()
	n += 1 + l + sovIdentityserver(uint64(l))
	return n
}

func sovIdentityserver(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIdentityserver(x uint64) (n int) {
	return sovIdentityserver((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *AuthInfoResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AuthInfoResponse{`,
		`AccessMethod:` + fmt.Sprintf("%v", this.AccessMethod) + `,`,
		`UniversalRights:` + strings.Replace(fmt.Sprintf("%v", this.UniversalRights), "Rights", "Rights", 1) + `,`,
		`IsAdmin:` + fmt.Sprintf("%v", this.IsAdmin) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AuthInfoResponse_APIKey) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AuthInfoResponse_APIKey{`,
		`APIKey:` + strings.Replace(fmt.Sprintf("%v", this.APIKey), "AuthInfoResponse_APIKeyAccess", "AuthInfoResponse_APIKeyAccess", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AuthInfoResponse_OAuthAccessToken) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AuthInfoResponse_OAuthAccessToken{`,
		`OAuthAccessToken:` + strings.Replace(fmt.Sprintf("%v", this.OAuthAccessToken), "OAuthAccessToken", "OAuthAccessToken", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AuthInfoResponse_UserSession) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AuthInfoResponse_UserSession{`,
		`UserSession:` + strings.Replace(fmt.Sprintf("%v", this.UserSession), "UserSession", "UserSession", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AuthInfoResponse_APIKeyAccess) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AuthInfoResponse_APIKeyAccess{`,
		`APIKey:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.APIKey), "APIKey", "APIKey", 1), `&`, ``, 1) + `,`,
		`EntityIDs:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.EntityIDs), "EntityIdentifiers", "EntityIdentifiers", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringIdentityserver(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *AuthInfoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentityserver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AuthInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AuthInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field APIKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIdentityserver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &AuthInfoResponse_APIKeyAccess{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.AccessMethod = &AuthInfoResponse_APIKey{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OAuthAccessToken", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIdentityserver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &OAuthAccessToken{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.AccessMethod = &AuthInfoResponse_OAuthAccessToken{v}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UniversalRights", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIdentityserver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UniversalRights == nil {
				m.UniversalRights = &Rights{}
			}
			if err := m.UniversalRights.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsAdmin", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsAdmin = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserSession", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIdentityserver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &UserSession{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.AccessMethod = &AuthInfoResponse_UserSession{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentityserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIdentityserver
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
func (m *AuthInfoResponse_APIKeyAccess) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentityserver
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: APIKeyAccess: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: APIKeyAccess: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field APIKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIdentityserver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.APIKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EntityIDs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIdentityserver
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EntityIDs.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentityserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIdentityserver
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIdentityserver
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
func skipIdentityserver(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIdentityserver
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
					return 0, ErrIntOverflowIdentityserver
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIdentityserver
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
			if length < 0 {
				return 0, ErrInvalidLengthIdentityserver
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupIdentityserver
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthIdentityserver
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthIdentityserver        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIdentityserver          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupIdentityserver = fmt.Errorf("proto: unexpected end of group")
)
