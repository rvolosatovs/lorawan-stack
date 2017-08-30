// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/application.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import google_protobuf1 "github.com/gogo/protobuf/types"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// ApplicationUp wraps multiple application-layer uplink types
type ApplicationUp struct {
	// Types that are valid to be assigned to Up:
	//	*ApplicationUp_UplinkMessage
	Up isApplicationUp_Up `protobuf_oneof:"up"`
}

func (m *ApplicationUp) Reset()                    { *m = ApplicationUp{} }
func (*ApplicationUp) ProtoMessage()               {}
func (*ApplicationUp) Descriptor() ([]byte, []int) { return fileDescriptorApplication, []int{0} }

type isApplicationUp_Up interface {
	isApplicationUp_Up()
	MarshalTo([]byte) (int, error)
	Size() int
}

type ApplicationUp_UplinkMessage struct {
	UplinkMessage *ApplicationUplink `protobuf:"bytes,1,opt,name=uplink_message,json=uplinkMessage,oneof"`
}

func (*ApplicationUp_UplinkMessage) isApplicationUp_Up() {}

func (m *ApplicationUp) GetUp() isApplicationUp_Up {
	if m != nil {
		return m.Up
	}
	return nil
}

func (m *ApplicationUp) GetUplinkMessage() *ApplicationUplink {
	if x, ok := m.GetUp().(*ApplicationUp_UplinkMessage); ok {
		return x.UplinkMessage
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ApplicationUp) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ApplicationUp_OneofMarshaler, _ApplicationUp_OneofUnmarshaler, _ApplicationUp_OneofSizer, []interface{}{
		(*ApplicationUp_UplinkMessage)(nil),
	}
}

func _ApplicationUp_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ApplicationUp)
	// up
	switch x := m.Up.(type) {
	case *ApplicationUp_UplinkMessage:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.UplinkMessage); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ApplicationUp.Up has unexpected type %T", x)
	}
	return nil
}

func _ApplicationUp_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ApplicationUp)
	switch tag {
	case 1: // up.uplink_message
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ApplicationUplink)
		err := b.DecodeMessage(msg)
		m.Up = &ApplicationUp_UplinkMessage{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ApplicationUp_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ApplicationUp)
	// up
	switch x := m.Up.(type) {
	case *ApplicationUp_UplinkMessage:
		s := proto.Size(x.UplinkMessage)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ApplicationUplink struct {
	FPort      uint32 `protobuf:"varint,1,opt,name=f_port,json=fPort,proto3" json:"f_port,omitempty"`
	FCnt       uint32 `protobuf:"varint,2,opt,name=f_cnt,json=fCnt,proto3" json:"f_cnt,omitempty"`
	FrmPayload []byte `protobuf:"bytes,3,opt,name=frm_payload,json=frmPayload,proto3" json:"frm_payload,omitempty"`
}

func (m *ApplicationUplink) Reset()                    { *m = ApplicationUplink{} }
func (*ApplicationUplink) ProtoMessage()               {}
func (*ApplicationUplink) Descriptor() ([]byte, []int) { return fileDescriptorApplication, []int{1} }

func (m *ApplicationUplink) GetFPort() uint32 {
	if m != nil {
		return m.FPort
	}
	return 0
}

func (m *ApplicationUplink) GetFCnt() uint32 {
	if m != nil {
		return m.FCnt
	}
	return 0
}

func (m *ApplicationUplink) GetFrmPayload() []byte {
	if m != nil {
		return m.FrmPayload
	}
	return nil
}

type ApplicationDownlink struct {
	FPort      uint32 `protobuf:"varint,1,opt,name=f_port,json=fPort,proto3" json:"f_port,omitempty"`
	FCnt       uint32 `protobuf:"varint,2,opt,name=f_cnt,json=fCnt,proto3" json:"f_cnt,omitempty"`
	FrmPayload []byte `protobuf:"bytes,3,opt,name=frm_payload,json=frmPayload,proto3" json:"frm_payload,omitempty"`
}

func (m *ApplicationDownlink) Reset()                    { *m = ApplicationDownlink{} }
func (*ApplicationDownlink) ProtoMessage()               {}
func (*ApplicationDownlink) Descriptor() ([]byte, []int) { return fileDescriptorApplication, []int{2} }

func (m *ApplicationDownlink) GetFPort() uint32 {
	if m != nil {
		return m.FPort
	}
	return 0
}

func (m *ApplicationDownlink) GetFCnt() uint32 {
	if m != nil {
		return m.FCnt
	}
	return 0
}

func (m *ApplicationDownlink) GetFrmPayload() []byte {
	if m != nil {
		return m.FrmPayload
	}
	return nil
}

type ApplicationDownlinks struct {
	Downlinks []*ApplicationDownlink `protobuf:"bytes,1,rep,name=downlinks" json:"downlinks,omitempty"`
}

func (m *ApplicationDownlinks) Reset()                    { *m = ApplicationDownlinks{} }
func (*ApplicationDownlinks) ProtoMessage()               {}
func (*ApplicationDownlinks) Descriptor() ([]byte, []int) { return fileDescriptorApplication, []int{3} }

func (m *ApplicationDownlinks) GetDownlinks() []*ApplicationDownlink {
	if m != nil {
		return m.Downlinks
	}
	return nil
}

type DownlinkQueueRequest struct {
	Downlinks            []*ApplicationDownlink `protobuf:"bytes,1,rep,name=downlinks" json:"downlinks,omitempty"`
	EndDeviceIdentifiers `protobuf:"bytes,2,opt,name=end_device,json=endDevice,embedded=end_device" json:"end_device"`
}

func (m *DownlinkQueueRequest) Reset()                    { *m = DownlinkQueueRequest{} }
func (*DownlinkQueueRequest) ProtoMessage()               {}
func (*DownlinkQueueRequest) Descriptor() ([]byte, []int) { return fileDescriptorApplication, []int{4} }

func (m *DownlinkQueueRequest) GetDownlinks() []*ApplicationDownlink {
	if m != nil {
		return m.Downlinks
	}
	return nil
}

func init() {
	proto.RegisterType((*ApplicationUp)(nil), "ttn.v3.ApplicationUp")
	proto.RegisterType((*ApplicationUplink)(nil), "ttn.v3.ApplicationUplink")
	proto.RegisterType((*ApplicationDownlink)(nil), "ttn.v3.ApplicationDownlink")
	proto.RegisterType((*ApplicationDownlinks)(nil), "ttn.v3.ApplicationDownlinks")
	proto.RegisterType((*DownlinkQueueRequest)(nil), "ttn.v3.DownlinkQueueRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ApplicationDownlinkQueue service

type ApplicationDownlinkQueueClient interface {
	DownlinkQueueReplace(ctx context.Context, in *DownlinkQueueRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	DownlinkQueuePush(ctx context.Context, in *DownlinkQueueRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	DownlinkQueueList(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*ApplicationDownlinks, error)
	DownlinkQueueClear(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
}

type applicationDownlinkQueueClient struct {
	cc *grpc.ClientConn
}

func NewApplicationDownlinkQueueClient(cc *grpc.ClientConn) ApplicationDownlinkQueueClient {
	return &applicationDownlinkQueueClient{cc}
}

func (c *applicationDownlinkQueueClient) DownlinkQueueReplace(ctx context.Context, in *DownlinkQueueRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueueReplace", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationDownlinkQueueClient) DownlinkQueuePush(ctx context.Context, in *DownlinkQueueRequest, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueuePush", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationDownlinkQueueClient) DownlinkQueueList(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*ApplicationDownlinks, error) {
	out := new(ApplicationDownlinks)
	err := grpc.Invoke(ctx, "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueueList", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationDownlinkQueueClient) DownlinkQueueClear(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueueClear", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ApplicationDownlinkQueue service

type ApplicationDownlinkQueueServer interface {
	DownlinkQueueReplace(context.Context, *DownlinkQueueRequest) (*google_protobuf1.Empty, error)
	DownlinkQueuePush(context.Context, *DownlinkQueueRequest) (*google_protobuf1.Empty, error)
	DownlinkQueueList(context.Context, *EndDeviceIdentifiers) (*ApplicationDownlinks, error)
	DownlinkQueueClear(context.Context, *EndDeviceIdentifiers) (*google_protobuf1.Empty, error)
}

func RegisterApplicationDownlinkQueueServer(s *grpc.Server, srv ApplicationDownlinkQueueServer) {
	s.RegisterService(&_ApplicationDownlinkQueue_serviceDesc, srv)
}

func _ApplicationDownlinkQueue_DownlinkQueueReplace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownlinkQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueueReplace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueueReplace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueueReplace(ctx, req.(*DownlinkQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationDownlinkQueue_DownlinkQueuePush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownlinkQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueuePush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueuePush",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueuePush(ctx, req.(*DownlinkQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationDownlinkQueue_DownlinkQueueList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueueList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueueList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueueList(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationDownlinkQueue_DownlinkQueueClear_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueueClear(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.ApplicationDownlinkQueue/DownlinkQueueClear",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationDownlinkQueueServer).DownlinkQueueClear(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApplicationDownlinkQueue_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.v3.ApplicationDownlinkQueue",
	HandlerType: (*ApplicationDownlinkQueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DownlinkQueueReplace",
			Handler:    _ApplicationDownlinkQueue_DownlinkQueueReplace_Handler,
		},
		{
			MethodName: "DownlinkQueuePush",
			Handler:    _ApplicationDownlinkQueue_DownlinkQueuePush_Handler,
		},
		{
			MethodName: "DownlinkQueueList",
			Handler:    _ApplicationDownlinkQueue_DownlinkQueueList_Handler,
		},
		{
			MethodName: "DownlinkQueueClear",
			Handler:    _ApplicationDownlinkQueue_DownlinkQueueClear_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/TheThingsNetwork/ttn/api/application.proto",
}

func (m *ApplicationUp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplicationUp) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Up != nil {
		nn1, err := m.Up.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	return i, nil
}

func (m *ApplicationUp_UplinkMessage) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.UplinkMessage != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintApplication(dAtA, i, uint64(m.UplinkMessage.Size()))
		n2, err := m.UplinkMessage.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}
func (m *ApplicationUplink) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplicationUplink) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.FPort != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintApplication(dAtA, i, uint64(m.FPort))
	}
	if m.FCnt != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintApplication(dAtA, i, uint64(m.FCnt))
	}
	if len(m.FrmPayload) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintApplication(dAtA, i, uint64(len(m.FrmPayload)))
		i += copy(dAtA[i:], m.FrmPayload)
	}
	return i, nil
}

func (m *ApplicationDownlink) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplicationDownlink) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.FPort != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintApplication(dAtA, i, uint64(m.FPort))
	}
	if m.FCnt != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintApplication(dAtA, i, uint64(m.FCnt))
	}
	if len(m.FrmPayload) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintApplication(dAtA, i, uint64(len(m.FrmPayload)))
		i += copy(dAtA[i:], m.FrmPayload)
	}
	return i, nil
}

func (m *ApplicationDownlinks) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplicationDownlinks) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Downlinks) > 0 {
		for _, msg := range m.Downlinks {
			dAtA[i] = 0xa
			i++
			i = encodeVarintApplication(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *DownlinkQueueRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DownlinkQueueRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Downlinks) > 0 {
		for _, msg := range m.Downlinks {
			dAtA[i] = 0xa
			i++
			i = encodeVarintApplication(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintApplication(dAtA, i, uint64(m.EndDeviceIdentifiers.Size()))
	n3, err := m.EndDeviceIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func encodeFixed64Application(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Application(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintApplication(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ApplicationUp) Size() (n int) {
	var l int
	_ = l
	if m.Up != nil {
		n += m.Up.Size()
	}
	return n
}

func (m *ApplicationUp_UplinkMessage) Size() (n int) {
	var l int
	_ = l
	if m.UplinkMessage != nil {
		l = m.UplinkMessage.Size()
		n += 1 + l + sovApplication(uint64(l))
	}
	return n
}
func (m *ApplicationUplink) Size() (n int) {
	var l int
	_ = l
	if m.FPort != 0 {
		n += 1 + sovApplication(uint64(m.FPort))
	}
	if m.FCnt != 0 {
		n += 1 + sovApplication(uint64(m.FCnt))
	}
	l = len(m.FrmPayload)
	if l > 0 {
		n += 1 + l + sovApplication(uint64(l))
	}
	return n
}

func (m *ApplicationDownlink) Size() (n int) {
	var l int
	_ = l
	if m.FPort != 0 {
		n += 1 + sovApplication(uint64(m.FPort))
	}
	if m.FCnt != 0 {
		n += 1 + sovApplication(uint64(m.FCnt))
	}
	l = len(m.FrmPayload)
	if l > 0 {
		n += 1 + l + sovApplication(uint64(l))
	}
	return n
}

func (m *ApplicationDownlinks) Size() (n int) {
	var l int
	_ = l
	if len(m.Downlinks) > 0 {
		for _, e := range m.Downlinks {
			l = e.Size()
			n += 1 + l + sovApplication(uint64(l))
		}
	}
	return n
}

func (m *DownlinkQueueRequest) Size() (n int) {
	var l int
	_ = l
	if len(m.Downlinks) > 0 {
		for _, e := range m.Downlinks {
			l = e.Size()
			n += 1 + l + sovApplication(uint64(l))
		}
	}
	l = m.EndDeviceIdentifiers.Size()
	n += 1 + l + sovApplication(uint64(l))
	return n
}

func sovApplication(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozApplication(x uint64) (n int) {
	return sovApplication(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *ApplicationUp) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplicationUp{`,
		`Up:` + fmt.Sprintf("%v", this.Up) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ApplicationUp_UplinkMessage) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplicationUp_UplinkMessage{`,
		`UplinkMessage:` + strings.Replace(fmt.Sprintf("%v", this.UplinkMessage), "ApplicationUplink", "ApplicationUplink", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ApplicationUplink) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplicationUplink{`,
		`FPort:` + fmt.Sprintf("%v", this.FPort) + `,`,
		`FCnt:` + fmt.Sprintf("%v", this.FCnt) + `,`,
		`FrmPayload:` + fmt.Sprintf("%v", this.FrmPayload) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ApplicationDownlink) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplicationDownlink{`,
		`FPort:` + fmt.Sprintf("%v", this.FPort) + `,`,
		`FCnt:` + fmt.Sprintf("%v", this.FCnt) + `,`,
		`FrmPayload:` + fmt.Sprintf("%v", this.FrmPayload) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ApplicationDownlinks) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplicationDownlinks{`,
		`Downlinks:` + strings.Replace(fmt.Sprintf("%v", this.Downlinks), "ApplicationDownlink", "ApplicationDownlink", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *DownlinkQueueRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&DownlinkQueueRequest{`,
		`Downlinks:` + strings.Replace(fmt.Sprintf("%v", this.Downlinks), "ApplicationDownlink", "ApplicationDownlink", 1) + `,`,
		`EndDeviceIdentifiers:` + strings.Replace(strings.Replace(this.EndDeviceIdentifiers.String(), "EndDeviceIdentifiers", "EndDeviceIdentifiers", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringApplication(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ApplicationUp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplication
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
			return fmt.Errorf("proto: ApplicationUp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplicationUp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UplinkMessage", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
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
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &ApplicationUplink{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Up = &ApplicationUp_UplinkMessage{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplication(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApplication
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
func (m *ApplicationUplink) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplication
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
			return fmt.Errorf("proto: ApplicationUplink: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplicationUplink: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FPort", wireType)
			}
			m.FPort = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FPort |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FCnt", wireType)
			}
			m.FCnt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FCnt |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FrmPayload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FrmPayload = append(m.FrmPayload[:0], dAtA[iNdEx:postIndex]...)
			if m.FrmPayload == nil {
				m.FrmPayload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplication(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApplication
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
func (m *ApplicationDownlink) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplication
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
			return fmt.Errorf("proto: ApplicationDownlink: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplicationDownlink: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FPort", wireType)
			}
			m.FPort = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FPort |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FCnt", wireType)
			}
			m.FCnt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FCnt |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FrmPayload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FrmPayload = append(m.FrmPayload[:0], dAtA[iNdEx:postIndex]...)
			if m.FrmPayload == nil {
				m.FrmPayload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplication(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApplication
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
func (m *ApplicationDownlinks) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplication
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
			return fmt.Errorf("proto: ApplicationDownlinks: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplicationDownlinks: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Downlinks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
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
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Downlinks = append(m.Downlinks, &ApplicationDownlink{})
			if err := m.Downlinks[len(m.Downlinks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplication(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApplication
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
func (m *DownlinkQueueRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplication
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
			return fmt.Errorf("proto: DownlinkQueueRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DownlinkQueueRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Downlinks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
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
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Downlinks = append(m.Downlinks, &ApplicationDownlink{})
			if err := m.Downlinks[len(m.Downlinks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndDeviceIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
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
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EndDeviceIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplication(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApplication
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
func skipApplication(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowApplication
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
					return 0, ErrIntOverflowApplication
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
					return 0, ErrIntOverflowApplication
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
				return 0, ErrInvalidLengthApplication
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowApplication
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
				next, err := skipApplication(dAtA[start:])
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
	ErrInvalidLengthApplication = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowApplication   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/application.proto", fileDescriptorApplication)
}

var fileDescriptorApplication = []byte{
	// 493 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0x3f, 0x6f, 0xd3, 0x40,
	0x14, 0xcf, 0xa5, 0x6d, 0x44, 0x2f, 0x04, 0xa9, 0xd7, 0x82, 0x42, 0x5a, 0x39, 0x51, 0xa6, 0x30,
	0x60, 0x4b, 0xa9, 0x18, 0x18, 0x49, 0x1b, 0x41, 0x51, 0x81, 0xd4, 0x2a, 0x03, 0x2c, 0xc6, 0x7f,
	0x9e, 0x9d, 0x53, 0xec, 0xbb, 0xc3, 0x77, 0x6e, 0xd5, 0x8d, 0x8f, 0xc1, 0x47, 0xea, 0x98, 0x91,
	0xa9, 0x02, 0x4f, 0x8c, 0x7c, 0x04, 0x94, 0x33, 0x26, 0x0d, 0x0d, 0x2d, 0x02, 0xa6, 0xe4, 0xbd,
	0xf7, 0xfb, 0xf7, 0x4e, 0xcf, 0xf8, 0x51, 0x44, 0xd5, 0x38, 0xf3, 0x4c, 0x9f, 0x27, 0xd6, 0xf1,
	0x18, 0x8e, 0xc7, 0x94, 0x45, 0xf2, 0x25, 0xa8, 0x53, 0x9e, 0x4e, 0x2c, 0xa5, 0x98, 0xe5, 0x0a,
	0x6a, 0xb9, 0x42, 0xc4, 0xd4, 0x77, 0x15, 0xe5, 0xcc, 0x14, 0x29, 0x57, 0x9c, 0xd4, 0x94, 0x62,
	0xe6, 0xc9, 0x6e, 0xeb, 0xe1, 0x25, 0x7a, 0xc4, 0x23, 0x6e, 0xe9, 0xb1, 0x97, 0x85, 0xba, 0xd2,
	0x85, 0xfe, 0x57, 0xd0, 0x5a, 0x7f, 0xe4, 0x46, 0x03, 0x60, 0x8a, 0x86, 0x14, 0x52, 0xf9, 0x83,
	0xb6, 0x1d, 0x71, 0x1e, 0xc5, 0x30, 0x17, 0x87, 0x44, 0xa8, 0xb3, 0x62, 0xd8, 0x7d, 0x83, 0x1b,
	0x4f, 0xe6, 0xf9, 0x5e, 0x0b, 0x32, 0xc0, 0x77, 0x32, 0x11, 0x53, 0x36, 0x71, 0x12, 0x90, 0xd2,
	0x8d, 0xa0, 0x89, 0x3a, 0xa8, 0x57, 0xef, 0xdf, 0x37, 0x8b, 0xd0, 0xe6, 0x02, 0x7c, 0x06, 0x7c,
	0x56, 0xb1, 0x1b, 0x05, 0xe5, 0x45, 0xc1, 0x18, 0xac, 0xe2, 0x6a, 0x26, 0xba, 0xef, 0xf0, 0xc6,
	0x15, 0x2c, 0xb9, 0x8b, 0x6b, 0xa1, 0x23, 0x78, 0xaa, 0xb4, 0x6c, 0xc3, 0x5e, 0x0b, 0x47, 0x3c,
	0x55, 0x64, 0x13, 0xaf, 0x85, 0x8e, 0xcf, 0x54, 0xb3, 0xaa, 0xbb, 0xab, 0xe1, 0x1e, 0x53, 0xa4,
	0x8d, 0xeb, 0x61, 0x9a, 0x38, 0xc2, 0x3d, 0x8b, 0xb9, 0x1b, 0x34, 0x57, 0x3a, 0xa8, 0x77, 0xdb,
	0xc6, 0x61, 0x9a, 0x8c, 0x8a, 0x4e, 0xd7, 0xc3, 0x9b, 0x97, 0x1c, 0xf6, 0xf9, 0x29, 0xfb, 0xff,
	0x1e, 0x47, 0x78, 0x6b, 0x89, 0x87, 0x24, 0x8f, 0xf1, 0x7a, 0x50, 0x16, 0x4d, 0xd4, 0x59, 0xe9,
	0xd5, 0xfb, 0xdb, 0x4b, 0x9e, 0xa8, 0x24, 0xd8, 0x73, 0x74, 0xf7, 0x23, 0xc2, 0x5b, 0x65, 0xff,
	0x28, 0x83, 0x0c, 0x6c, 0x78, 0x9f, 0x81, 0x54, 0xff, 0xa0, 0x49, 0x86, 0x18, 0x03, 0x0b, 0x9c,
	0x00, 0x4e, 0xa8, 0x0f, 0x7a, 0xc3, 0x7a, 0x7f, 0xa7, 0xe4, 0x0e, 0x59, 0xb0, 0xaf, 0x07, 0x07,
	0xf3, 0xe3, 0x18, 0xdc, 0x3a, 0xbf, 0x68, 0x57, 0xa6, 0x17, 0x6d, 0x64, 0xaf, 0x43, 0x39, 0xef,
	0x4f, 0xab, 0xb8, 0xb9, 0xc4, 0x49, 0xa7, 0x24, 0x87, 0x57, 0x62, 0x8b, 0xd8, 0xf5, 0x81, 0xfc,
	0xf4, 0x59, 0xb6, 0x54, 0xeb, 0x9e, 0x59, 0xdc, 0x9f, 0x59, 0xde, 0x9f, 0x39, 0x9c, 0xdd, 0x1f,
	0x39, 0xc0, 0x1b, 0x0b, 0xf8, 0x51, 0x26, 0xc7, 0x7f, 0x29, 0xf5, 0xea, 0x17, 0xa9, 0x43, 0x2a,
	0x15, 0xb9, 0x76, 0xfb, 0xd6, 0xce, 0x35, 0xef, 0x2a, 0xc9, 0x73, 0x4c, 0x16, 0x04, 0xf7, 0x62,
	0x70, 0xd3, 0x1b, 0x14, 0x7f, 0x13, 0x6e, 0xf0, 0xf4, 0xd3, 0x17, 0xa3, 0xf2, 0x21, 0x37, 0xd0,
	0x79, 0x6e, 0xa0, 0x69, 0x6e, 0xa0, 0xcf, 0xb9, 0x81, 0xbe, 0xe6, 0x46, 0xe5, 0x5b, 0x6e, 0xa0,
	0xb7, 0x0f, 0x6e, 0xfa, 0xa6, 0xc5, 0x24, 0x9a, 0xfd, 0x0a, 0xcf, 0xab, 0x69, 0xe1, 0xdd, 0xef,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x22, 0xba, 0xb8, 0x75, 0x04, 0x00, 0x00,
}
