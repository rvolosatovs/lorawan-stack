// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/gatewayserver.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import types "github.com/gogo/protobuf/types"

import (
	context "context"

	grpc "google.golang.org/grpc"
)

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

// GatewayUp may contain zero or more uplink messages and/or a status message for the gateway.
type GatewayUp struct {
	// UplinkMessages received by the gateway.
	UplinkMessages       []*UplinkMessage  `protobuf:"bytes,1,rep,name=uplink_messages,json=uplinkMessages,proto3" json:"uplink_messages,omitempty"`
	GatewayStatus        *GatewayStatus    `protobuf:"bytes,2,opt,name=gateway_status,json=gatewayStatus,proto3" json:"gateway_status,omitempty"`
	TxAcknowledgment     *TxAcknowledgment `protobuf:"bytes,3,opt,name=tx_acknowledgment,json=txAcknowledgment,proto3" json:"tx_acknowledgment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *GatewayUp) Reset()      { *m = GatewayUp{} }
func (*GatewayUp) ProtoMessage() {}
func (*GatewayUp) Descriptor() ([]byte, []int) {
	return fileDescriptor_gatewayserver_14f5294a2b645388, []int{0}
}
func (m *GatewayUp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GatewayUp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GatewayUp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *GatewayUp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GatewayUp.Merge(dst, src)
}
func (m *GatewayUp) XXX_Size() int {
	return m.Size()
}
func (m *GatewayUp) XXX_DiscardUnknown() {
	xxx_messageInfo_GatewayUp.DiscardUnknown(m)
}

var xxx_messageInfo_GatewayUp proto.InternalMessageInfo

func (m *GatewayUp) GetUplinkMessages() []*UplinkMessage {
	if m != nil {
		return m.UplinkMessages
	}
	return nil
}

func (m *GatewayUp) GetGatewayStatus() *GatewayStatus {
	if m != nil {
		return m.GatewayStatus
	}
	return nil
}

func (m *GatewayUp) GetTxAcknowledgment() *TxAcknowledgment {
	if m != nil {
		return m.TxAcknowledgment
	}
	return nil
}

// GatewayDown contains downlink messages for the gateway.
type GatewayDown struct {
	// DownlinkMessage for the gateway.
	DownlinkMessage      *DownlinkMessage `protobuf:"bytes,1,opt,name=downlink_message,json=downlinkMessage,proto3" json:"downlink_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GatewayDown) Reset()      { *m = GatewayDown{} }
func (*GatewayDown) ProtoMessage() {}
func (*GatewayDown) Descriptor() ([]byte, []int) {
	return fileDescriptor_gatewayserver_14f5294a2b645388, []int{1}
}
func (m *GatewayDown) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GatewayDown) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GatewayDown.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *GatewayDown) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GatewayDown.Merge(dst, src)
}
func (m *GatewayDown) XXX_Size() int {
	return m.Size()
}
func (m *GatewayDown) XXX_DiscardUnknown() {
	xxx_messageInfo_GatewayDown.DiscardUnknown(m)
}

var xxx_messageInfo_GatewayDown proto.InternalMessageInfo

func (m *GatewayDown) GetDownlinkMessage() *DownlinkMessage {
	if m != nil {
		return m.DownlinkMessage
	}
	return nil
}

func init() {
	proto.RegisterType((*GatewayUp)(nil), "ttn.lorawan.v3.GatewayUp")
	golang_proto.RegisterType((*GatewayUp)(nil), "ttn.lorawan.v3.GatewayUp")
	proto.RegisterType((*GatewayDown)(nil), "ttn.lorawan.v3.GatewayDown")
	golang_proto.RegisterType((*GatewayDown)(nil), "ttn.lorawan.v3.GatewayDown")
}
func (this *GatewayUp) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GatewayUp)
	if !ok {
		that2, ok := that.(GatewayUp)
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
	if len(this.UplinkMessages) != len(that1.UplinkMessages) {
		return false
	}
	for i := range this.UplinkMessages {
		if !this.UplinkMessages[i].Equal(that1.UplinkMessages[i]) {
			return false
		}
	}
	if !this.GatewayStatus.Equal(that1.GatewayStatus) {
		return false
	}
	if !this.TxAcknowledgment.Equal(that1.TxAcknowledgment) {
		return false
	}
	return true
}
func (this *GatewayDown) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GatewayDown)
	if !ok {
		that2, ok := that.(GatewayDown)
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
	if !this.DownlinkMessage.Equal(that1.DownlinkMessage) {
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

// GtwGsClient is the client API for GtwGs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GtwGsClient interface {
	// Link the gateway to the Gateway Server.
	LinkGateway(ctx context.Context, opts ...grpc.CallOption) (GtwGs_LinkGatewayClient, error)
	// GetConcentratorConfig associated to the gateway.
	GetConcentratorConfig(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*ConcentratorConfig, error)
}

type gtwGsClient struct {
	cc *grpc.ClientConn
}

func NewGtwGsClient(cc *grpc.ClientConn) GtwGsClient {
	return &gtwGsClient{cc}
}

func (c *gtwGsClient) LinkGateway(ctx context.Context, opts ...grpc.CallOption) (GtwGs_LinkGatewayClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GtwGs_serviceDesc.Streams[0], "/ttn.lorawan.v3.GtwGs/LinkGateway", opts...)
	if err != nil {
		return nil, err
	}
	x := &gtwGsLinkGatewayClient{stream}
	return x, nil
}

type GtwGs_LinkGatewayClient interface {
	Send(*GatewayUp) error
	Recv() (*GatewayDown, error)
	grpc.ClientStream
}

type gtwGsLinkGatewayClient struct {
	grpc.ClientStream
}

func (x *gtwGsLinkGatewayClient) Send(m *GatewayUp) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gtwGsLinkGatewayClient) Recv() (*GatewayDown, error) {
	m := new(GatewayDown)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gtwGsClient) GetConcentratorConfig(ctx context.Context, in *types.Empty, opts ...grpc.CallOption) (*ConcentratorConfig, error) {
	out := new(ConcentratorConfig)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.GtwGs/GetConcentratorConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GtwGsServer is the server API for GtwGs service.
type GtwGsServer interface {
	// Link the gateway to the Gateway Server.
	LinkGateway(GtwGs_LinkGatewayServer) error
	// GetConcentratorConfig associated to the gateway.
	GetConcentratorConfig(context.Context, *types.Empty) (*ConcentratorConfig, error)
}

func RegisterGtwGsServer(s *grpc.Server, srv GtwGsServer) {
	s.RegisterService(&_GtwGs_serviceDesc, srv)
}

func _GtwGs_LinkGateway_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GtwGsServer).LinkGateway(&gtwGsLinkGatewayServer{stream})
}

type GtwGs_LinkGatewayServer interface {
	Send(*GatewayDown) error
	Recv() (*GatewayUp, error)
	grpc.ServerStream
}

type gtwGsLinkGatewayServer struct {
	grpc.ServerStream
}

func (x *gtwGsLinkGatewayServer) Send(m *GatewayDown) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gtwGsLinkGatewayServer) Recv() (*GatewayUp, error) {
	m := new(GatewayUp)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GtwGs_GetConcentratorConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtwGsServer).GetConcentratorConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.GtwGs/GetConcentratorConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtwGsServer).GetConcentratorConfig(ctx, req.(*types.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _GtwGs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.GtwGs",
	HandlerType: (*GtwGsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConcentratorConfig",
			Handler:    _GtwGs_GetConcentratorConfig_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "LinkGateway",
			Handler:       _GtwGs_LinkGateway_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "lorawan-stack/api/gatewayserver.proto",
}

// NsGsClient is the client API for NsGs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NsGsClient interface {
	// ScheduleDownlink instructs the Gateway Server to schedule a downlink message.
	// The Gateway Server may refuse if there are any conflicts in the schedule or
	// if a duty cycle prevents the gateway from transmitting.
	ScheduleDownlink(ctx context.Context, in *DownlinkMessage, opts ...grpc.CallOption) (*types.Empty, error)
}

type nsGsClient struct {
	cc *grpc.ClientConn
}

func NewNsGsClient(cc *grpc.ClientConn) NsGsClient {
	return &nsGsClient{cc}
}

func (c *nsGsClient) ScheduleDownlink(ctx context.Context, in *DownlinkMessage, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.NsGs/ScheduleDownlink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NsGsServer is the server API for NsGs service.
type NsGsServer interface {
	// ScheduleDownlink instructs the Gateway Server to schedule a downlink message.
	// The Gateway Server may refuse if there are any conflicts in the schedule or
	// if a duty cycle prevents the gateway from transmitting.
	ScheduleDownlink(context.Context, *DownlinkMessage) (*types.Empty, error)
}

func RegisterNsGsServer(s *grpc.Server, srv NsGsServer) {
	s.RegisterService(&_NsGs_serviceDesc, srv)
}

func _NsGs_ScheduleDownlink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownlinkMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NsGsServer).ScheduleDownlink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.NsGs/ScheduleDownlink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NsGsServer).ScheduleDownlink(ctx, req.(*DownlinkMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _NsGs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.NsGs",
	HandlerType: (*NsGsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ScheduleDownlink",
			Handler:    _NsGs_ScheduleDownlink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/gatewayserver.proto",
}

// GsClient is the client API for Gs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GsClient interface {
	// Get statistics about the current gateway connection to the Gateway Server.
	// This is not persisted between reconnects.
	GetGatewayConnectionStats(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*GatewayConnectionStats, error)
}

type gsClient struct {
	cc *grpc.ClientConn
}

func NewGsClient(cc *grpc.ClientConn) GsClient {
	return &gsClient{cc}
}

func (c *gsClient) GetGatewayConnectionStats(ctx context.Context, in *GatewayIdentifiers, opts ...grpc.CallOption) (*GatewayConnectionStats, error) {
	out := new(GatewayConnectionStats)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.Gs/GetGatewayConnectionStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GsServer is the server API for Gs service.
type GsServer interface {
	// Get statistics about the current gateway connection to the Gateway Server.
	// This is not persisted between reconnects.
	GetGatewayConnectionStats(context.Context, *GatewayIdentifiers) (*GatewayConnectionStats, error)
}

func RegisterGsServer(s *grpc.Server, srv GsServer) {
	s.RegisterService(&_Gs_serviceDesc, srv)
}

func _Gs_GetGatewayConnectionStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GsServer).GetGatewayConnectionStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.Gs/GetGatewayConnectionStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GsServer).GetGatewayConnectionStats(ctx, req.(*GatewayIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _Gs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.Gs",
	HandlerType: (*GsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGatewayConnectionStats",
			Handler:    _Gs_GetGatewayConnectionStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/gatewayserver.proto",
}

func (m *GatewayUp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GatewayUp) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.UplinkMessages) > 0 {
		for _, msg := range m.UplinkMessages {
			dAtA[i] = 0xa
			i++
			i = encodeVarintGatewayserver(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.GatewayStatus != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintGatewayserver(dAtA, i, uint64(m.GatewayStatus.Size()))
		n1, err := m.GatewayStatus.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.TxAcknowledgment != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintGatewayserver(dAtA, i, uint64(m.TxAcknowledgment.Size()))
		n2, err := m.TxAcknowledgment.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *GatewayDown) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GatewayDown) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.DownlinkMessage != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintGatewayserver(dAtA, i, uint64(m.DownlinkMessage.Size()))
		n3, err := m.DownlinkMessage.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func encodeVarintGatewayserver(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedGatewayUp(r randyGatewayserver, easy bool) *GatewayUp {
	this := &GatewayUp{}
	if r.Intn(10) != 0 {
		v1 := r.Intn(5)
		this.UplinkMessages = make([]*UplinkMessage, v1)
		for i := 0; i < v1; i++ {
			this.UplinkMessages[i] = NewPopulatedUplinkMessage(r, easy)
		}
	}
	if r.Intn(10) != 0 {
		this.GatewayStatus = NewPopulatedGatewayStatus(r, easy)
	}
	if r.Intn(10) != 0 {
		this.TxAcknowledgment = NewPopulatedTxAcknowledgment(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedGatewayDown(r randyGatewayserver, easy bool) *GatewayDown {
	this := &GatewayDown{}
	if r.Intn(10) != 0 {
		this.DownlinkMessage = NewPopulatedDownlinkMessage(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyGatewayserver interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneGatewayserver(r randyGatewayserver) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringGatewayserver(r randyGatewayserver) string {
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
		tmps[i] = randUTF8RuneGatewayserver(r)
	}
	return string(tmps)
}
func randUnrecognizedGatewayserver(r randyGatewayserver, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldGatewayserver(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldGatewayserver(dAtA []byte, r randyGatewayserver, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateGatewayserver(dAtA, uint64(key))
		v3 := r.Int63()
		if r.Intn(2) == 0 {
			v3 *= -1
		}
		dAtA = encodeVarintPopulateGatewayserver(dAtA, uint64(v3))
	case 1:
		dAtA = encodeVarintPopulateGatewayserver(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateGatewayserver(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateGatewayserver(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateGatewayserver(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateGatewayserver(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *GatewayUp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.UplinkMessages) > 0 {
		for _, e := range m.UplinkMessages {
			l = e.Size()
			n += 1 + l + sovGatewayserver(uint64(l))
		}
	}
	if m.GatewayStatus != nil {
		l = m.GatewayStatus.Size()
		n += 1 + l + sovGatewayserver(uint64(l))
	}
	if m.TxAcknowledgment != nil {
		l = m.TxAcknowledgment.Size()
		n += 1 + l + sovGatewayserver(uint64(l))
	}
	return n
}

func (m *GatewayDown) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DownlinkMessage != nil {
		l = m.DownlinkMessage.Size()
		n += 1 + l + sovGatewayserver(uint64(l))
	}
	return n
}

func sovGatewayserver(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozGatewayserver(x uint64) (n int) {
	return sovGatewayserver((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *GatewayUp) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&GatewayUp{`,
		`UplinkMessages:` + strings.Replace(fmt.Sprintf("%v", this.UplinkMessages), "UplinkMessage", "UplinkMessage", 1) + `,`,
		`GatewayStatus:` + strings.Replace(fmt.Sprintf("%v", this.GatewayStatus), "GatewayStatus", "GatewayStatus", 1) + `,`,
		`TxAcknowledgment:` + strings.Replace(fmt.Sprintf("%v", this.TxAcknowledgment), "TxAcknowledgment", "TxAcknowledgment", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *GatewayDown) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&GatewayDown{`,
		`DownlinkMessage:` + strings.Replace(fmt.Sprintf("%v", this.DownlinkMessage), "DownlinkMessage", "DownlinkMessage", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGatewayserver(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *GatewayUp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGatewayserver
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
			return fmt.Errorf("proto: GatewayUp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GatewayUp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UplinkMessages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayserver
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
				return ErrInvalidLengthGatewayserver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UplinkMessages = append(m.UplinkMessages, &UplinkMessage{})
			if err := m.UplinkMessages[len(m.UplinkMessages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GatewayStatus", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayserver
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
				return ErrInvalidLengthGatewayserver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GatewayStatus == nil {
				m.GatewayStatus = &GatewayStatus{}
			}
			if err := m.GatewayStatus.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxAcknowledgment", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayserver
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
				return ErrInvalidLengthGatewayserver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TxAcknowledgment == nil {
				m.TxAcknowledgment = &TxAcknowledgment{}
			}
			if err := m.TxAcknowledgment.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGatewayserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGatewayserver
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
func (m *GatewayDown) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGatewayserver
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
			return fmt.Errorf("proto: GatewayDown: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GatewayDown: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DownlinkMessage", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGatewayserver
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
				return ErrInvalidLengthGatewayserver
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DownlinkMessage == nil {
				m.DownlinkMessage = &DownlinkMessage{}
			}
			if err := m.DownlinkMessage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGatewayserver(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGatewayserver
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
func skipGatewayserver(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGatewayserver
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
					return 0, ErrIntOverflowGatewayserver
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
					return 0, ErrIntOverflowGatewayserver
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
				return 0, ErrInvalidLengthGatewayserver
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGatewayserver
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
				next, err := skipGatewayserver(dAtA[start:])
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
	ErrInvalidLengthGatewayserver = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGatewayserver   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("lorawan-stack/api/gatewayserver.proto", fileDescriptor_gatewayserver_14f5294a2b645388)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/gatewayserver.proto", fileDescriptor_gatewayserver_14f5294a2b645388)
}

var fileDescriptor_gatewayserver_14f5294a2b645388 = []byte{
	// 586 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xb1, 0x53, 0x13, 0x41,
	0x18, 0xc5, 0x77, 0x41, 0x9d, 0x71, 0x19, 0x01, 0x77, 0x46, 0x27, 0x84, 0xf1, 0x23, 0x13, 0x47,
	0x87, 0x42, 0x2e, 0x4e, 0xf8, 0x0b, 0x14, 0xf4, 0x06, 0x15, 0x0b, 0x22, 0x85, 0x36, 0xcc, 0x25,
	0x59, 0x36, 0x37, 0x49, 0x76, 0x6f, 0x6e, 0xbf, 0x70, 0xd0, 0x51, 0x52, 0x5a, 0x5a, 0x3a, 0x36,
	0x52, 0x52, 0x52, 0x52, 0x52, 0x52, 0x52, 0x72, 0x7b, 0x0d, 0x25, 0x25, 0xa5, 0xc3, 0xe5, 0x32,
	0x90, 0xc4, 0xd3, 0xee, 0xbe, 0xdd, 0xdf, 0xbe, 0x77, 0x6f, 0xef, 0x1d, 0x7b, 0xd1, 0xd1, 0xa1,
	0x17, 0x79, 0x6a, 0xc9, 0xa0, 0xd7, 0x68, 0x57, 0xbc, 0xc0, 0xaf, 0x48, 0x0f, 0x45, 0xe4, 0xed,
	0x19, 0x11, 0xee, 0x88, 0xd0, 0x09, 0x42, 0x8d, 0x9a, 0x4f, 0x23, 0x2a, 0x27, 0x43, 0x9d, 0x9d,
	0xe5, 0xe2, 0x92, 0xf4, 0xb1, 0xd5, 0xab, 0x3b, 0x0d, 0xdd, 0xad, 0x48, 0x2d, 0x75, 0x25, 0xc5,
	0xea, 0xbd, 0xed, 0x74, 0x4a, 0x87, 0xf4, 0xa9, 0x7f, 0xbc, 0x38, 0x2f, 0xb5, 0x96, 0x1d, 0x71,
	0x4b, 0x89, 0x6e, 0x80, 0x7b, 0xd9, 0xe6, 0x42, 0xee, 0x2b, 0x64, 0xc0, 0xf3, 0x71, 0xc0, 0x6f,
	0x0a, 0x85, 0xfe, 0xb6, 0x2f, 0x42, 0x93, 0x41, 0xa5, 0x71, 0xa8, 0x2b, 0x8c, 0xf1, 0xa4, 0xf8,
	0x07, 0x11, 0x0a, 0xe9, 0x6b, 0xe5, 0x75, 0xfa, 0x44, 0xf9, 0x92, 0xb2, 0x87, 0x6e, 0xdf, 0x7a,
	0x33, 0xe0, 0xef, 0xd9, 0x4c, 0x2f, 0xe8, 0xf8, 0xaa, 0xbd, 0x35, 0x10, 0x2a, 0xd0, 0xd2, 0xe4,
	0xe2, 0x54, 0xf5, 0x99, 0x33, 0x7c, 0x1b, 0xce, 0x66, 0x8a, 0xad, 0xf7, 0xa9, 0x8d, 0xe9, 0xde,
	0xdd, 0xd1, 0xf0, 0x55, 0x36, 0x9d, 0xe5, 0xd9, 0x32, 0xe8, 0x61, 0xcf, 0x14, 0x26, 0x4a, 0xf4,
	0x6f, 0x32, 0x99, 0x75, 0x2d, 0x85, 0x36, 0x1e, 0xc9, 0xbb, 0x23, 0x5f, 0x67, 0x8f, 0x71, 0x77,
	0xcb, 0x6b, 0xb4, 0x95, 0x8e, 0x3a, 0xa2, 0x29, 0xbb, 0x42, 0x61, 0x61, 0x32, 0x15, 0x2a, 0x8d,
	0x0a, 0x7d, 0xd9, 0x7d, 0x33, 0xc4, 0x6d, 0xcc, 0xe2, 0xc8, 0x4a, 0xf9, 0x2b, 0x9b, 0xca, 0xec,
	0x56, 0x75, 0xa4, 0xf8, 0x07, 0x36, 0xdb, 0xd4, 0x91, 0xba, 0x9b, 0xb6, 0x40, 0x53, 0xf1, 0x85,
	0x51, 0xf1, 0xd5, 0x8c, 0x1b, 0xc4, 0x9d, 0x69, 0x0e, 0x2f, 0x54, 0x7f, 0x53, 0x76, 0xdf, 0xc5,
	0xc8, 0x35, 0x7c, 0x8d, 0x4d, 0x7d, 0xf2, 0x55, 0x3b, 0x33, 0xe2, 0x73, 0x39, 0x81, 0x37, 0x83,
	0xe2, 0x7c, 0xce, 0xd6, 0x8d, 0xd9, 0x22, 0x7d, 0x4d, 0x79, 0x8d, 0x3d, 0x71, 0x05, 0xae, 0x68,
	0xd5, 0x10, 0x0a, 0x43, 0x0f, 0x75, 0xb8, 0xa2, 0xd5, 0xb6, 0x2f, 0xf9, 0x53, 0xa7, 0xdf, 0x2d,
	0x67, 0xd0, 0x2d, 0xe7, 0xdd, 0x4d, 0xb7, 0x8a, 0xe5, 0x51, 0xc5, 0xf1, 0xb3, 0xd5, 0x1a, 0xbb,
	0xf7, 0xd9, 0xb8, 0x86, 0x7f, 0x64, 0xb3, 0xb5, 0x46, 0x4b, 0x34, 0x7b, 0x1d, 0x31, 0x48, 0xc7,
	0xff, 0x97, 0xbb, 0x98, 0x63, 0x5c, 0xed, 0xb2, 0x09, 0xd7, 0x70, 0xc9, 0xe6, 0x5c, 0x81, 0x59,
	0x8a, 0x15, 0xad, 0x94, 0x68, 0xa0, 0xaf, 0xd5, 0xcd, 0xc7, 0x34, 0xbc, 0x9c, 0x93, 0x76, 0xed,
	0xb6, 0xd5, 0xc5, 0x97, 0x39, 0xcc, 0x88, 0xd6, 0xdb, 0x5f, 0xf4, 0x34, 0x06, 0x7a, 0x16, 0x03,
	0x3d, 0x8f, 0x81, 0x5c, 0xc4, 0x40, 0x2e, 0x63, 0x20, 0x57, 0x31, 0x90, 0xeb, 0x18, 0xe8, 0xbe,
	0x05, 0x7a, 0x60, 0x81, 0x1c, 0x5a, 0xa0, 0x47, 0x16, 0xc8, 0xb1, 0x05, 0x72, 0x62, 0x81, 0x9c,
	0x5a, 0xa0, 0x67, 0x16, 0xe8, 0xb9, 0x05, 0x72, 0x61, 0x81, 0x5e, 0x5a, 0x20, 0x57, 0x16, 0xe8,
	0xb5, 0x05, 0xb2, 0x9f, 0x00, 0x39, 0x48, 0x80, 0x7e, 0x4f, 0x80, 0xfc, 0x48, 0x80, 0xfe, 0x4c,
	0x80, 0x1c, 0x26, 0x40, 0x8e, 0x12, 0xa0, 0xc7, 0x09, 0xd0, 0x93, 0x04, 0xe8, 0xb7, 0x57, 0x52,
	0x3b, 0xd8, 0x12, 0xd8, 0xf2, 0x95, 0x34, 0x8e, 0x12, 0x18, 0xe9, 0xb0, 0x5d, 0x19, 0xfe, 0xc3,
	0x82, 0xb6, 0xac, 0x20, 0xaa, 0xa0, 0x5e, 0x7f, 0x90, 0xde, 0xd1, 0xf2, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x60, 0x0f, 0xc8, 0xf4, 0x6e, 0x04, 0x00, 0x00,
}
