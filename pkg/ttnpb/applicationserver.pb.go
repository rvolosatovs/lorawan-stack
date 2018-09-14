// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/applicationserver.proto

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

type ApplicationLink struct {
	NetworkServerAddress string                    `protobuf:"bytes,1,opt,name=network_server_address,json=networkServerAddress,proto3" json:"network_server_address,omitempty"`
	APIKey               string                    `protobuf:"bytes,2,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	DefaultFormatters    *MessagePayloadFormatters `protobuf:"bytes,3,opt,name=default_formatters,json=defaultFormatters" json:"default_formatters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *ApplicationLink) Reset()      { *m = ApplicationLink{} }
func (*ApplicationLink) ProtoMessage() {}
func (*ApplicationLink) Descriptor() ([]byte, []int) {
	return fileDescriptor_applicationserver_fa9ff3978a1857d5, []int{0}
}
func (m *ApplicationLink) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ApplicationLink) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ApplicationLink.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ApplicationLink) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplicationLink.Merge(dst, src)
}
func (m *ApplicationLink) XXX_Size() int {
	return m.Size()
}
func (m *ApplicationLink) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplicationLink.DiscardUnknown(m)
}

var xxx_messageInfo_ApplicationLink proto.InternalMessageInfo

func (m *ApplicationLink) GetNetworkServerAddress() string {
	if m != nil {
		return m.NetworkServerAddress
	}
	return ""
}

func (m *ApplicationLink) GetAPIKey() string {
	if m != nil {
		return m.APIKey
	}
	return ""
}

func (m *ApplicationLink) GetDefaultFormatters() *MessagePayloadFormatters {
	if m != nil {
		return m.DefaultFormatters
	}
	return nil
}

type SetApplicationLinkRequest struct {
	ApplicationIdentifiers `protobuf:"bytes,1,opt,name=application_ids,json=applicationIds,embedded=application_ids" json:"application_ids"`
	ApplicationLink        `protobuf:"bytes,2,opt,name=link,embedded=link" json:"link"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_sizecache          int32    `json:"-"`
}

func (m *SetApplicationLinkRequest) Reset()      { *m = SetApplicationLinkRequest{} }
func (*SetApplicationLinkRequest) ProtoMessage() {}
func (*SetApplicationLinkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_applicationserver_fa9ff3978a1857d5, []int{1}
}
func (m *SetApplicationLinkRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetApplicationLinkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetApplicationLinkRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *SetApplicationLinkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetApplicationLinkRequest.Merge(dst, src)
}
func (m *SetApplicationLinkRequest) XXX_Size() int {
	return m.Size()
}
func (m *SetApplicationLinkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetApplicationLinkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetApplicationLinkRequest proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ApplicationLink)(nil), "ttn.lorawan.v3.ApplicationLink")
	golang_proto.RegisterType((*ApplicationLink)(nil), "ttn.lorawan.v3.ApplicationLink")
	proto.RegisterType((*SetApplicationLinkRequest)(nil), "ttn.lorawan.v3.SetApplicationLinkRequest")
	golang_proto.RegisterType((*SetApplicationLinkRequest)(nil), "ttn.lorawan.v3.SetApplicationLinkRequest")
}
func (this *ApplicationLink) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ApplicationLink)
	if !ok {
		that2, ok := that.(ApplicationLink)
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
	if this.NetworkServerAddress != that1.NetworkServerAddress {
		return false
	}
	if this.APIKey != that1.APIKey {
		return false
	}
	if !this.DefaultFormatters.Equal(that1.DefaultFormatters) {
		return false
	}
	return true
}
func (this *SetApplicationLinkRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SetApplicationLinkRequest)
	if !ok {
		that2, ok := that.(SetApplicationLinkRequest)
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
	if !this.ApplicationIdentifiers.Equal(&that1.ApplicationIdentifiers) {
		return false
	}
	if !this.ApplicationLink.Equal(&that1.ApplicationLink) {
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

// Client API for As service

type AsClient interface {
	GetLink(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*ApplicationLink, error)
	SetLink(ctx context.Context, in *SetApplicationLinkRequest, opts ...grpc.CallOption) (*types.Empty, error)
	DeleteLink(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error)
	Subscribe(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (As_SubscribeClient, error)
}

type asClient struct {
	cc *grpc.ClientConn
}

func NewAsClient(cc *grpc.ClientConn) AsClient {
	return &asClient{cc}
}

func (c *asClient) GetLink(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*ApplicationLink, error) {
	out := new(ApplicationLink)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.As/GetLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asClient) SetLink(ctx context.Context, in *SetApplicationLinkRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.As/SetLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asClient) DeleteLink(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.As/DeleteLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asClient) Subscribe(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (As_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_As_serviceDesc.Streams[0], "/ttn.lorawan.v3.As/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &asSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type As_SubscribeClient interface {
	Recv() (*ApplicationUp, error)
	grpc.ClientStream
}

type asSubscribeClient struct {
	grpc.ClientStream
}

func (x *asSubscribeClient) Recv() (*ApplicationUp, error) {
	m := new(ApplicationUp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for As service

type AsServer interface {
	GetLink(context.Context, *ApplicationIdentifiers) (*ApplicationLink, error)
	SetLink(context.Context, *SetApplicationLinkRequest) (*types.Empty, error)
	DeleteLink(context.Context, *ApplicationIdentifiers) (*types.Empty, error)
	Subscribe(*ApplicationIdentifiers, As_SubscribeServer) error
}

func RegisterAsServer(s *grpc.Server, srv AsServer) {
	s.RegisterService(&_As_serviceDesc, srv)
}

func _As_GetLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplicationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsServer).GetLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.As/GetLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsServer).GetLink(ctx, req.(*ApplicationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _As_SetLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetApplicationLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsServer).SetLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.As/SetLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsServer).SetLink(ctx, req.(*SetApplicationLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _As_DeleteLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplicationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsServer).DeleteLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.As/DeleteLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsServer).DeleteLink(ctx, req.(*ApplicationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _As_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ApplicationIdentifiers)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AsServer).Subscribe(m, &asSubscribeServer{stream})
}

type As_SubscribeServer interface {
	Send(*ApplicationUp) error
	grpc.ServerStream
}

type asSubscribeServer struct {
	grpc.ServerStream
}

func (x *asSubscribeServer) Send(m *ApplicationUp) error {
	return x.ServerStream.SendMsg(m)
}

var _As_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.As",
	HandlerType: (*AsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLink",
			Handler:    _As_GetLink_Handler,
		},
		{
			MethodName: "SetLink",
			Handler:    _As_SetLink_Handler,
		},
		{
			MethodName: "DeleteLink",
			Handler:    _As_DeleteLink_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _As_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "lorawan-stack/api/applicationserver.proto",
}

func (m *ApplicationLink) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplicationLink) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.NetworkServerAddress) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintApplicationserver(dAtA, i, uint64(len(m.NetworkServerAddress)))
		i += copy(dAtA[i:], m.NetworkServerAddress)
	}
	if len(m.APIKey) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintApplicationserver(dAtA, i, uint64(len(m.APIKey)))
		i += copy(dAtA[i:], m.APIKey)
	}
	if m.DefaultFormatters != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintApplicationserver(dAtA, i, uint64(m.DefaultFormatters.Size()))
		n1, err := m.DefaultFormatters.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *SetApplicationLinkRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetApplicationLinkRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintApplicationserver(dAtA, i, uint64(m.ApplicationIdentifiers.Size()))
	n2, err := m.ApplicationIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	dAtA[i] = 0x12
	i++
	i = encodeVarintApplicationserver(dAtA, i, uint64(m.ApplicationLink.Size()))
	n3, err := m.ApplicationLink.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func encodeVarintApplicationserver(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedApplicationLink(r randyApplicationserver, easy bool) *ApplicationLink {
	this := &ApplicationLink{}
	this.NetworkServerAddress = randStringApplicationserver(r)
	this.APIKey = randStringApplicationserver(r)
	if r.Intn(10) != 0 {
		this.DefaultFormatters = NewPopulatedMessagePayloadFormatters(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedSetApplicationLinkRequest(r randyApplicationserver, easy bool) *SetApplicationLinkRequest {
	this := &SetApplicationLinkRequest{}
	v1 := NewPopulatedApplicationIdentifiers(r, easy)
	this.ApplicationIdentifiers = *v1
	v2 := NewPopulatedApplicationLink(r, easy)
	this.ApplicationLink = *v2
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyApplicationserver interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneApplicationserver(r randyApplicationserver) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringApplicationserver(r randyApplicationserver) string {
	v3 := r.Intn(100)
	tmps := make([]rune, v3)
	for i := 0; i < v3; i++ {
		tmps[i] = randUTF8RuneApplicationserver(r)
	}
	return string(tmps)
}
func randUnrecognizedApplicationserver(r randyApplicationserver, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldApplicationserver(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldApplicationserver(dAtA []byte, r randyApplicationserver, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateApplicationserver(dAtA, uint64(key))
		v4 := r.Int63()
		if r.Intn(2) == 0 {
			v4 *= -1
		}
		dAtA = encodeVarintPopulateApplicationserver(dAtA, uint64(v4))
	case 1:
		dAtA = encodeVarintPopulateApplicationserver(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateApplicationserver(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateApplicationserver(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateApplicationserver(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateApplicationserver(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *ApplicationLink) Size() (n int) {
	var l int
	_ = l
	l = len(m.NetworkServerAddress)
	if l > 0 {
		n += 1 + l + sovApplicationserver(uint64(l))
	}
	l = len(m.APIKey)
	if l > 0 {
		n += 1 + l + sovApplicationserver(uint64(l))
	}
	if m.DefaultFormatters != nil {
		l = m.DefaultFormatters.Size()
		n += 1 + l + sovApplicationserver(uint64(l))
	}
	return n
}

func (m *SetApplicationLinkRequest) Size() (n int) {
	var l int
	_ = l
	l = m.ApplicationIdentifiers.Size()
	n += 1 + l + sovApplicationserver(uint64(l))
	l = m.ApplicationLink.Size()
	n += 1 + l + sovApplicationserver(uint64(l))
	return n
}

func sovApplicationserver(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozApplicationserver(x uint64) (n int) {
	return sovApplicationserver((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *ApplicationLink) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplicationLink{`,
		`NetworkServerAddress:` + fmt.Sprintf("%v", this.NetworkServerAddress) + `,`,
		`APIKey:` + fmt.Sprintf("%v", this.APIKey) + `,`,
		`DefaultFormatters:` + strings.Replace(fmt.Sprintf("%v", this.DefaultFormatters), "MessagePayloadFormatters", "MessagePayloadFormatters", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *SetApplicationLinkRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SetApplicationLinkRequest{`,
		`ApplicationIdentifiers:` + strings.Replace(strings.Replace(this.ApplicationIdentifiers.String(), "ApplicationIdentifiers", "ApplicationIdentifiers", 1), `&`, ``, 1) + `,`,
		`ApplicationLink:` + strings.Replace(strings.Replace(this.ApplicationLink.String(), "ApplicationLink", "ApplicationLink", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringApplicationserver(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ApplicationLink) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplicationserver
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
			return fmt.Errorf("proto: ApplicationLink: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplicationLink: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkServerAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplicationserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApplicationserver
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetworkServerAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field APIKey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplicationserver
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApplicationserver
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.APIKey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultFormatters", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplicationserver
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
				return ErrInvalidLengthApplicationserver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DefaultFormatters == nil {
				m.DefaultFormatters = &MessagePayloadFormatters{}
			}
			if err := m.DefaultFormatters.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplicationserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApplicationserver
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
func (m *SetApplicationLinkRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplicationserver
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
			return fmt.Errorf("proto: SetApplicationLinkRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetApplicationLinkRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApplicationIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplicationserver
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
				return ErrInvalidLengthApplicationserver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ApplicationIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApplicationLink", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplicationserver
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
				return ErrInvalidLengthApplicationserver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ApplicationLink.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplicationserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApplicationserver
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
func skipApplicationserver(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowApplicationserver
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
					return 0, ErrIntOverflowApplicationserver
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
					return 0, ErrIntOverflowApplicationserver
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
				return 0, ErrInvalidLengthApplicationserver
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowApplicationserver
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
				next, err := skipApplicationserver(dAtA[start:])
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
	ErrInvalidLengthApplicationserver = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowApplicationserver   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("lorawan-stack/api/applicationserver.proto", fileDescriptor_applicationserver_fa9ff3978a1857d5)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/applicationserver.proto", fileDescriptor_applicationserver_fa9ff3978a1857d5)
}

var fileDescriptor_applicationserver_fa9ff3978a1857d5 = []byte{
	// 646 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4d, 0x4c, 0x13, 0x41,
	0x18, 0x9d, 0x01, 0x03, 0x32, 0x24, 0x10, 0x37, 0x86, 0x60, 0xd5, 0xaf, 0xa4, 0x24, 0xa4, 0x10,
	0xd9, 0x35, 0xc5, 0x78, 0xd0, 0x78, 0x28, 0xf1, 0x27, 0x04, 0x4d, 0x48, 0xab, 0x31, 0x7a, 0x69,
	0xa6, 0xed, 0x74, 0x99, 0xb4, 0xdd, 0x59, 0x77, 0xa6, 0x90, 0x46, 0x4d, 0x08, 0x27, 0x6e, 0x9a,
	0x78, 0xf1, 0x68, 0x3c, 0x71, 0xe4, 0x62, 0xc2, 0x91, 0x8b, 0x09, 0x47, 0x12, 0x2f, 0x9c, 0x08,
	0x3b, 0xeb, 0x81, 0x9b, 0x1c, 0x39, 0x9a, 0xee, 0x2e, 0xd0, 0x82, 0xfc, 0x78, 0xeb, 0xf4, 0xbd,
	0x79, 0xef, 0x7d, 0x6f, 0xbf, 0x5d, 0x32, 0x5e, 0x13, 0x1e, 0x5d, 0xa4, 0xce, 0xa4, 0x54, 0xb4,
	0x54, 0xb5, 0xa8, 0xcb, 0x2d, 0xea, 0xba, 0x35, 0x5e, 0xa2, 0x8a, 0x0b, 0x47, 0x32, 0x6f, 0x81,
	0x79, 0xa6, 0xeb, 0x09, 0x25, 0x8c, 0x01, 0xa5, 0x1c, 0x33, 0xa6, 0x9b, 0x0b, 0x53, 0x89, 0x49,
	0x9b, 0xab, 0xf9, 0x46, 0xd1, 0x2c, 0x89, 0xba, 0x65, 0x0b, 0x5b, 0x58, 0x21, 0xad, 0xd8, 0xa8,
	0x84, 0xa7, 0xf0, 0x10, 0xfe, 0x8a, 0xae, 0x27, 0x6e, 0xd9, 0x42, 0xd8, 0x35, 0x16, 0x59, 0x38,
	0x8e, 0x50, 0x91, 0x43, 0x8c, 0xde, 0x8c, 0xd1, 0x23, 0x0d, 0x56, 0x77, 0x55, 0x33, 0x06, 0x47,
	0x4f, 0x87, 0xe4, 0x65, 0xe6, 0x28, 0x5e, 0xe1, 0xcc, 0x3b, 0x54, 0x18, 0x39, 0x4d, 0xaa, 0x33,
	0x29, 0xa9, 0xcd, 0x62, 0x46, 0xea, 0x27, 0x26, 0x83, 0xd9, 0xe3, 0xe1, 0x9e, 0x73, 0xa7, 0x6a,
	0xdc, 0x23, 0x43, 0x0e, 0x53, 0x8b, 0xc2, 0xab, 0x16, 0xa2, 0x61, 0x0b, 0xb4, 0x5c, 0xf6, 0x98,
	0x94, 0xc3, 0x78, 0x04, 0xa7, 0xfb, 0x72, 0xd7, 0x63, 0x34, 0x1f, 0x82, 0xd9, 0x08, 0x33, 0x46,
	0x49, 0x2f, 0x75, 0x79, 0xa1, 0xca, 0x9a, 0xc3, 0x5d, 0x2d, 0xda, 0x34, 0xd1, 0x3b, 0xc9, 0x9e,
	0xec, 0xdc, 0xcc, 0x2c, 0x6b, 0xe6, 0x7a, 0xa8, 0xcb, 0x67, 0x59, 0xd3, 0x78, 0x4d, 0x8c, 0x32,
	0xab, 0xd0, 0x46, 0x4d, 0x15, 0x2a, 0xc2, 0xab, 0x53, 0xa5, 0x98, 0x27, 0x87, 0xbb, 0x47, 0x70,
	0xba, 0x3f, 0x93, 0x36, 0x3b, 0xcb, 0x34, 0x5f, 0x44, 0x51, 0xe7, 0x68, 0xb3, 0x26, 0x68, 0xf9,
	0xe9, 0x11, 0x3f, 0x77, 0x2d, 0xd6, 0x38, 0xfe, 0x2b, 0xf5, 0x03, 0x93, 0x1b, 0x79, 0xa6, 0x4e,
	0x8c, 0x92, 0x63, 0xef, 0x1a, 0x4c, 0x2a, 0xe3, 0x0d, 0x19, 0x6c, 0x7b, 0x82, 0x05, 0x5e, 0x8e,
	0x46, 0xe9, 0xcf, 0x8c, 0x9d, 0xf4, 0x6c, 0x13, 0x98, 0x39, 0xae, 0x73, 0xfa, 0xea, 0xe6, 0x4e,
	0x12, 0x6d, 0xed, 0x24, 0x71, 0x6e, 0x80, 0xb6, 0x33, 0xa4, 0xf1, 0x88, 0x5c, 0xa9, 0x71, 0xa7,
	0x1a, 0xce, 0xdc, 0x9f, 0x49, 0x9e, 0xa3, 0xd7, 0x0a, 0xd4, 0x26, 0x14, 0x5e, 0xcb, 0xfc, 0xe9,
	0x26, 0x5d, 0x59, 0x69, 0x2c, 0x63, 0xd2, 0xfb, 0x8c, 0xa9, 0xb0, 0xfe, 0x4b, 0x66, 0x4a, 0x5c,
	0xe4, 0x95, 0x32, 0x97, 0x7f, 0xfd, 0xfe, 0xd2, 0x95, 0x36, 0xc6, 0x2c, 0x2a, 0x3b, 0x36, 0xd8,
	0x7a, 0xdf, 0xd9, 0xc6, 0x47, 0xab, 0x95, 0xc5, 0xf8, 0x84, 0x49, 0x6f, 0x3e, 0x0e, 0x31, 0x7e,
	0x52, 0xfc, 0xcc, 0x72, 0x13, 0x43, 0x66, 0xb4, 0xa7, 0xe6, 0xe1, 0x9e, 0x9a, 0x4f, 0x5a, 0x7b,
	0x9a, 0xca, 0x86, 0xf6, 0x0f, 0x13, 0xf7, 0x2f, 0xb2, 0x97, 0xe6, 0xbf, 0xe2, 0x3c, 0xc0, 0x13,
	0xc6, 0x07, 0x42, 0x1e, 0xb3, 0x1a, 0x53, 0xec, 0xbf, 0x8a, 0x39, 0x2b, 0x50, 0xdc, 0xc7, 0xc4,
	0x65, 0xfb, 0x78, 0x49, 0xfa, 0xf2, 0x8d, 0xa2, 0x2c, 0x79, 0xbc, 0xc8, 0x2e, 0x6d, 0x7e, 0xfb,
	0x1c, 0xde, 0x2b, 0xf7, 0x2e, 0x9e, 0xfe, 0x8e, 0x37, 0x7d, 0xc0, 0x5b, 0x3e, 0xe0, 0x6d, 0x1f,
	0xd0, 0xae, 0x0f, 0x68, 0xcf, 0x07, 0xb4, 0xef, 0x03, 0x3a, 0xf0, 0x01, 0x2f, 0x69, 0xc0, 0x2b,
	0x1a, 0xd0, 0xaa, 0x06, 0xbc, 0xa6, 0x01, 0xad, 0x6b, 0x40, 0x1b, 0x1a, 0xd0, 0xa6, 0x06, 0xbc,
	0xa5, 0x01, 0x6f, 0x6b, 0x40, 0xbb, 0x1a, 0xf0, 0x9e, 0x06, 0xb4, 0xaf, 0x01, 0x1f, 0x68, 0x40,
	0x4b, 0x01, 0xa0, 0x95, 0x00, 0xf0, 0xe7, 0x00, 0xd0, 0xd7, 0x00, 0xf0, 0xb7, 0x00, 0xd0, 0x6a,
	0x00, 0x68, 0x2d, 0x00, 0xbc, 0x1e, 0x00, 0xde, 0x08, 0x00, 0xbf, 0xbd, 0x63, 0x0b, 0x53, 0xcd,
	0x33, 0x35, 0xcf, 0x1d, 0x5b, 0x9a, 0xf1, 0x1b, 0x6c, 0x75, 0x7e, 0x1f, 0xdc, 0xaa, 0x6d, 0x29,
	0xe5, 0xb8, 0xc5, 0x62, 0x4f, 0x58, 0xdd, 0xd4, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x42, 0x20,
	0x3d, 0xf7, 0x0b, 0x05, 0x00, 0x00,
}
