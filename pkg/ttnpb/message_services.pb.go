// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/message_services.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

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

type ProcessUplinkMessageRequest struct {
	EndDeviceIdentifiers `protobuf:"bytes,1,opt,name=ids,embedded=ids" json:"ids"`
	EndDeviceVersionIDs  EndDeviceVersionIdentifiers `protobuf:"bytes,2,opt,name=end_device_version_ids,json=endDeviceVersionIds" json:"end_device_version_ids"`
	Message              ApplicationUplink           `protobuf:"bytes,3,opt,name=message" json:"message"`
	Parameter            string                      `protobuf:"bytes,4,opt,name=parameter,proto3" json:"parameter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *ProcessUplinkMessageRequest) Reset()      { *m = ProcessUplinkMessageRequest{} }
func (*ProcessUplinkMessageRequest) ProtoMessage() {}
func (*ProcessUplinkMessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_services_3d37a928d34fca12, []int{0}
}
func (m *ProcessUplinkMessageRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessUplinkMessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessUplinkMessageRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ProcessUplinkMessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessUplinkMessageRequest.Merge(dst, src)
}
func (m *ProcessUplinkMessageRequest) XXX_Size() int {
	return m.Size()
}
func (m *ProcessUplinkMessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessUplinkMessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessUplinkMessageRequest proto.InternalMessageInfo

func (m *ProcessUplinkMessageRequest) GetEndDeviceVersionIDs() EndDeviceVersionIdentifiers {
	if m != nil {
		return m.EndDeviceVersionIDs
	}
	return EndDeviceVersionIdentifiers{}
}

func (m *ProcessUplinkMessageRequest) GetMessage() ApplicationUplink {
	if m != nil {
		return m.Message
	}
	return ApplicationUplink{}
}

func (m *ProcessUplinkMessageRequest) GetParameter() string {
	if m != nil {
		return m.Parameter
	}
	return ""
}

type ProcessDownlinkMessageRequest struct {
	EndDeviceIdentifiers `protobuf:"bytes,1,opt,name=ids,embedded=ids" json:"ids"`
	EndDeviceVersionIDs  EndDeviceVersionIdentifiers `protobuf:"bytes,2,opt,name=end_device_version_ids,json=endDeviceVersionIds" json:"end_device_version_ids"`
	Message              ApplicationDownlink         `protobuf:"bytes,3,opt,name=message" json:"message"`
	Parameter            string                      `protobuf:"bytes,4,opt,name=parameter,proto3" json:"parameter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *ProcessDownlinkMessageRequest) Reset()      { *m = ProcessDownlinkMessageRequest{} }
func (*ProcessDownlinkMessageRequest) ProtoMessage() {}
func (*ProcessDownlinkMessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_services_3d37a928d34fca12, []int{1}
}
func (m *ProcessDownlinkMessageRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessDownlinkMessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessDownlinkMessageRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ProcessDownlinkMessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessDownlinkMessageRequest.Merge(dst, src)
}
func (m *ProcessDownlinkMessageRequest) XXX_Size() int {
	return m.Size()
}
func (m *ProcessDownlinkMessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessDownlinkMessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessDownlinkMessageRequest proto.InternalMessageInfo

func (m *ProcessDownlinkMessageRequest) GetEndDeviceVersionIDs() EndDeviceVersionIdentifiers {
	if m != nil {
		return m.EndDeviceVersionIDs
	}
	return EndDeviceVersionIdentifiers{}
}

func (m *ProcessDownlinkMessageRequest) GetMessage() ApplicationDownlink {
	if m != nil {
		return m.Message
	}
	return ApplicationDownlink{}
}

func (m *ProcessDownlinkMessageRequest) GetParameter() string {
	if m != nil {
		return m.Parameter
	}
	return ""
}

func init() {
	proto.RegisterType((*ProcessUplinkMessageRequest)(nil), "ttn.lorawan.v3.ProcessUplinkMessageRequest")
	golang_proto.RegisterType((*ProcessUplinkMessageRequest)(nil), "ttn.lorawan.v3.ProcessUplinkMessageRequest")
	proto.RegisterType((*ProcessDownlinkMessageRequest)(nil), "ttn.lorawan.v3.ProcessDownlinkMessageRequest")
	golang_proto.RegisterType((*ProcessDownlinkMessageRequest)(nil), "ttn.lorawan.v3.ProcessDownlinkMessageRequest")
}
func (this *ProcessUplinkMessageRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProcessUplinkMessageRequest)
	if !ok {
		that2, ok := that.(ProcessUplinkMessageRequest)
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
	if !this.EndDeviceIdentifiers.Equal(&that1.EndDeviceIdentifiers) {
		return false
	}
	if !this.EndDeviceVersionIDs.Equal(&that1.EndDeviceVersionIDs) {
		return false
	}
	if !this.Message.Equal(&that1.Message) {
		return false
	}
	if this.Parameter != that1.Parameter {
		return false
	}
	return true
}
func (this *ProcessDownlinkMessageRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProcessDownlinkMessageRequest)
	if !ok {
		that2, ok := that.(ProcessDownlinkMessageRequest)
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
	if !this.EndDeviceIdentifiers.Equal(&that1.EndDeviceIdentifiers) {
		return false
	}
	if !this.EndDeviceVersionIDs.Equal(&that1.EndDeviceVersionIDs) {
		return false
	}
	if !this.Message.Equal(&that1.Message) {
		return false
	}
	if this.Parameter != that1.Parameter {
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

// Client API for UplinkMessageProcessor service

type UplinkMessageProcessorClient interface {
	Process(ctx context.Context, in *ProcessUplinkMessageRequest, opts ...grpc.CallOption) (*ApplicationUplink, error)
}

type uplinkMessageProcessorClient struct {
	cc *grpc.ClientConn
}

func NewUplinkMessageProcessorClient(cc *grpc.ClientConn) UplinkMessageProcessorClient {
	return &uplinkMessageProcessorClient{cc}
}

func (c *uplinkMessageProcessorClient) Process(ctx context.Context, in *ProcessUplinkMessageRequest, opts ...grpc.CallOption) (*ApplicationUplink, error) {
	out := new(ApplicationUplink)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.UplinkMessageProcessor/Process", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UplinkMessageProcessor service

type UplinkMessageProcessorServer interface {
	Process(context.Context, *ProcessUplinkMessageRequest) (*ApplicationUplink, error)
}

func RegisterUplinkMessageProcessorServer(s *grpc.Server, srv UplinkMessageProcessorServer) {
	s.RegisterService(&_UplinkMessageProcessor_serviceDesc, srv)
}

func _UplinkMessageProcessor_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessUplinkMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UplinkMessageProcessorServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.UplinkMessageProcessor/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UplinkMessageProcessorServer).Process(ctx, req.(*ProcessUplinkMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UplinkMessageProcessor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.UplinkMessageProcessor",
	HandlerType: (*UplinkMessageProcessorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _UplinkMessageProcessor_Process_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/message_services.proto",
}

// Client API for DownlinkMessageProcessor service

type DownlinkMessageProcessorClient interface {
	Process(ctx context.Context, in *ProcessDownlinkMessageRequest, opts ...grpc.CallOption) (*ApplicationDownlink, error)
}

type downlinkMessageProcessorClient struct {
	cc *grpc.ClientConn
}

func NewDownlinkMessageProcessorClient(cc *grpc.ClientConn) DownlinkMessageProcessorClient {
	return &downlinkMessageProcessorClient{cc}
}

func (c *downlinkMessageProcessorClient) Process(ctx context.Context, in *ProcessDownlinkMessageRequest, opts ...grpc.CallOption) (*ApplicationDownlink, error) {
	out := new(ApplicationDownlink)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.DownlinkMessageProcessor/Process", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DownlinkMessageProcessor service

type DownlinkMessageProcessorServer interface {
	Process(context.Context, *ProcessDownlinkMessageRequest) (*ApplicationDownlink, error)
}

func RegisterDownlinkMessageProcessorServer(s *grpc.Server, srv DownlinkMessageProcessorServer) {
	s.RegisterService(&_DownlinkMessageProcessor_serviceDesc, srv)
}

func _DownlinkMessageProcessor_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessDownlinkMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownlinkMessageProcessorServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.DownlinkMessageProcessor/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownlinkMessageProcessorServer).Process(ctx, req.(*ProcessDownlinkMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DownlinkMessageProcessor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.DownlinkMessageProcessor",
	HandlerType: (*DownlinkMessageProcessorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _DownlinkMessageProcessor_Process_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/message_services.proto",
}

func (m *ProcessUplinkMessageRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessUplinkMessageRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintMessageServices(dAtA, i, uint64(m.EndDeviceIdentifiers.Size()))
	n1, err := m.EndDeviceIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x12
	i++
	i = encodeVarintMessageServices(dAtA, i, uint64(m.EndDeviceVersionIDs.Size()))
	n2, err := m.EndDeviceVersionIDs.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	dAtA[i] = 0x1a
	i++
	i = encodeVarintMessageServices(dAtA, i, uint64(m.Message.Size()))
	n3, err := m.Message.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	if len(m.Parameter) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintMessageServices(dAtA, i, uint64(len(m.Parameter)))
		i += copy(dAtA[i:], m.Parameter)
	}
	return i, nil
}

func (m *ProcessDownlinkMessageRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessDownlinkMessageRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintMessageServices(dAtA, i, uint64(m.EndDeviceIdentifiers.Size()))
	n4, err := m.EndDeviceIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	dAtA[i] = 0x12
	i++
	i = encodeVarintMessageServices(dAtA, i, uint64(m.EndDeviceVersionIDs.Size()))
	n5, err := m.EndDeviceVersionIDs.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	dAtA[i] = 0x1a
	i++
	i = encodeVarintMessageServices(dAtA, i, uint64(m.Message.Size()))
	n6, err := m.Message.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n6
	if len(m.Parameter) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintMessageServices(dAtA, i, uint64(len(m.Parameter)))
		i += copy(dAtA[i:], m.Parameter)
	}
	return i, nil
}

func encodeVarintMessageServices(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedProcessUplinkMessageRequest(r randyMessageServices, easy bool) *ProcessUplinkMessageRequest {
	this := &ProcessUplinkMessageRequest{}
	v1 := NewPopulatedEndDeviceIdentifiers(r, easy)
	this.EndDeviceIdentifiers = *v1
	v2 := NewPopulatedEndDeviceVersionIdentifiers(r, easy)
	this.EndDeviceVersionIDs = *v2
	v3 := NewPopulatedApplicationUplink(r, easy)
	this.Message = *v3
	this.Parameter = randStringMessageServices(r)
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedProcessDownlinkMessageRequest(r randyMessageServices, easy bool) *ProcessDownlinkMessageRequest {
	this := &ProcessDownlinkMessageRequest{}
	v4 := NewPopulatedEndDeviceIdentifiers(r, easy)
	this.EndDeviceIdentifiers = *v4
	v5 := NewPopulatedEndDeviceVersionIdentifiers(r, easy)
	this.EndDeviceVersionIDs = *v5
	v6 := NewPopulatedApplicationDownlink(r, easy)
	this.Message = *v6
	this.Parameter = randStringMessageServices(r)
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyMessageServices interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneMessageServices(r randyMessageServices) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringMessageServices(r randyMessageServices) string {
	v7 := r.Intn(100)
	tmps := make([]rune, v7)
	for i := 0; i < v7; i++ {
		tmps[i] = randUTF8RuneMessageServices(r)
	}
	return string(tmps)
}
func randUnrecognizedMessageServices(r randyMessageServices, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldMessageServices(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldMessageServices(dAtA []byte, r randyMessageServices, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateMessageServices(dAtA, uint64(key))
		v8 := r.Int63()
		if r.Intn(2) == 0 {
			v8 *= -1
		}
		dAtA = encodeVarintPopulateMessageServices(dAtA, uint64(v8))
	case 1:
		dAtA = encodeVarintPopulateMessageServices(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateMessageServices(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateMessageServices(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateMessageServices(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateMessageServices(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *ProcessUplinkMessageRequest) Size() (n int) {
	var l int
	_ = l
	l = m.EndDeviceIdentifiers.Size()
	n += 1 + l + sovMessageServices(uint64(l))
	l = m.EndDeviceVersionIDs.Size()
	n += 1 + l + sovMessageServices(uint64(l))
	l = m.Message.Size()
	n += 1 + l + sovMessageServices(uint64(l))
	l = len(m.Parameter)
	if l > 0 {
		n += 1 + l + sovMessageServices(uint64(l))
	}
	return n
}

func (m *ProcessDownlinkMessageRequest) Size() (n int) {
	var l int
	_ = l
	l = m.EndDeviceIdentifiers.Size()
	n += 1 + l + sovMessageServices(uint64(l))
	l = m.EndDeviceVersionIDs.Size()
	n += 1 + l + sovMessageServices(uint64(l))
	l = m.Message.Size()
	n += 1 + l + sovMessageServices(uint64(l))
	l = len(m.Parameter)
	if l > 0 {
		n += 1 + l + sovMessageServices(uint64(l))
	}
	return n
}

func sovMessageServices(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMessageServices(x uint64) (n int) {
	return sovMessageServices((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *ProcessUplinkMessageRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProcessUplinkMessageRequest{`,
		`EndDeviceIdentifiers:` + strings.Replace(strings.Replace(this.EndDeviceIdentifiers.String(), "EndDeviceIdentifiers", "EndDeviceIdentifiers", 1), `&`, ``, 1) + `,`,
		`EndDeviceVersionIDs:` + strings.Replace(strings.Replace(this.EndDeviceVersionIDs.String(), "EndDeviceVersionIdentifiers", "EndDeviceVersionIdentifiers", 1), `&`, ``, 1) + `,`,
		`Message:` + strings.Replace(strings.Replace(this.Message.String(), "ApplicationUplink", "ApplicationUplink", 1), `&`, ``, 1) + `,`,
		`Parameter:` + fmt.Sprintf("%v", this.Parameter) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProcessDownlinkMessageRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProcessDownlinkMessageRequest{`,
		`EndDeviceIdentifiers:` + strings.Replace(strings.Replace(this.EndDeviceIdentifiers.String(), "EndDeviceIdentifiers", "EndDeviceIdentifiers", 1), `&`, ``, 1) + `,`,
		`EndDeviceVersionIDs:` + strings.Replace(strings.Replace(this.EndDeviceVersionIDs.String(), "EndDeviceVersionIdentifiers", "EndDeviceVersionIdentifiers", 1), `&`, ``, 1) + `,`,
		`Message:` + strings.Replace(strings.Replace(this.Message.String(), "ApplicationDownlink", "ApplicationDownlink", 1), `&`, ``, 1) + `,`,
		`Parameter:` + fmt.Sprintf("%v", this.Parameter) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringMessageServices(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ProcessUplinkMessageRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessageServices
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
			return fmt.Errorf("proto: ProcessUplinkMessageRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessUplinkMessageRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndDeviceIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EndDeviceIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndDeviceVersionIDs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EndDeviceVersionIDs.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Message.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Parameter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Parameter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessageServices(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessageServices
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
func (m *ProcessDownlinkMessageRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessageServices
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
			return fmt.Errorf("proto: ProcessDownlinkMessageRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessDownlinkMessageRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndDeviceIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EndDeviceIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndDeviceVersionIDs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EndDeviceVersionIDs.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Message.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Parameter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessageServices
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
				return ErrInvalidLengthMessageServices
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Parameter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessageServices(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessageServices
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
func skipMessageServices(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessageServices
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
					return 0, ErrIntOverflowMessageServices
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
					return 0, ErrIntOverflowMessageServices
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
				return 0, ErrInvalidLengthMessageServices
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMessageServices
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
				next, err := skipMessageServices(dAtA[start:])
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
	ErrInvalidLengthMessageServices = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessageServices   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("lorawan-stack/api/message_services.proto", fileDescriptor_message_services_3d37a928d34fca12)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/message_services.proto", fileDescriptor_message_services_3d37a928d34fca12)
}

var fileDescriptor_message_services_3d37a928d34fca12 = []byte{
	// 533 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x94, 0xb1, 0x6f, 0xd3, 0x40,
	0x14, 0xc6, 0xef, 0xd2, 0x8a, 0xd2, 0x43, 0x62, 0x70, 0xa5, 0x2a, 0x4a, 0xe1, 0x25, 0xa4, 0x0c,
	0x91, 0x20, 0xb6, 0x94, 0xfe, 0x03, 0xb4, 0x84, 0x81, 0x01, 0x09, 0x45, 0x02, 0x09, 0x24, 0x14,
	0x39, 0xc9, 0xd5, 0x39, 0x25, 0xb9, 0x33, 0x77, 0x97, 0x44, 0x0c, 0x48, 0x1d, 0x3b, 0x32, 0x32,
	0x22, 0xa6, 0x8e, 0x1d, 0xbb, 0x51, 0x89, 0x25, 0x63, 0xc6, 0x4e, 0x51, 0x7d, 0x5e, 0x3a, 0x76,
	0xec, 0x88, 0xe2, 0xb8, 0x32, 0x89, 0x49, 0x81, 0xb5, 0x9b, 0x9f, 0xef, 0x7b, 0xef, 0x7b, 0xef,
	0xf7, 0x4e, 0x47, 0x4a, 0x5d, 0x21, 0xdd, 0xa1, 0xcb, 0xcb, 0x4a, 0xbb, 0xcd, 0x8e, 0xe3, 0xfa,
	0xcc, 0xe9, 0x51, 0xa5, 0x5c, 0x8f, 0xd6, 0x15, 0x95, 0x03, 0xd6, 0xa4, 0xca, 0xf6, 0xa5, 0xd0,
	0xc2, 0xba, 0xaf, 0x35, 0xb7, 0x63, 0xb5, 0x3d, 0xd8, 0xc9, 0x95, 0x3d, 0xa6, 0xdb, 0xfd, 0x86,
	0xdd, 0x14, 0x3d, 0xc7, 0x13, 0x9e, 0x70, 0x22, 0x59, 0xa3, 0xbf, 0x1f, 0x45, 0x51, 0x10, 0x7d,
	0xcd, 0xd2, 0x73, 0xc5, 0xb4, 0x11, 0xe5, 0xad, 0x7a, 0x8b, 0x4e, 0x3d, 0x62, 0xcd, 0x76, 0x5a,
	0xc3, 0x5a, 0x94, 0x6b, 0xb6, 0xcf, 0xa8, 0x8c, 0xfb, 0xc8, 0x15, 0x96, 0x76, 0x1c, 0x2b, 0x8a,
	0x3f, 0x32, 0x64, 0xeb, 0xb5, 0x14, 0x4d, 0xaa, 0xd4, 0x1b, 0xbf, 0xcb, 0x78, 0xe7, 0xd5, 0xec,
	0xbc, 0x46, 0x3f, 0xf6, 0xa9, 0xd2, 0xd6, 0x33, 0xb2, 0xc2, 0x5a, 0x2a, 0x8b, 0x0b, 0xb8, 0x74,
	0xaf, 0xf2, 0xd8, 0x9e, 0x9f, 0xcb, 0x7e, 0xc1, 0x5b, 0xd5, 0xa8, 0xa9, 0x97, 0x89, 0xf5, 0xde,
	0xdd, 0xd1, 0x24, 0x8f, 0xc6, 0x93, 0x3c, 0xae, 0x4d, 0x53, 0xad, 0xcf, 0x64, 0x33, 0x69, 0xbe,
	0x3e, 0xa0, 0x52, 0x31, 0xc1, 0xeb, 0xd3, 0xa2, 0x99, 0xa8, 0xe8, 0x93, 0xa5, 0x45, 0xdf, 0xce,
	0xb4, 0xbf, 0xd7, 0xde, 0x9a, 0xd6, 0x36, 0x93, 0xfc, 0x46, 0x4a, 0x54, 0x55, 0xb5, 0x0d, 0x9a,
	0xca, 0x54, 0xd6, 0x2e, 0x59, 0x8b, 0x47, 0xce, 0xae, 0x44, 0x7e, 0x8f, 0x16, 0xfd, 0x76, 0x7d,
	0xbf, 0xcb, 0x9a, 0xae, 0x66, 0x82, 0xcf, 0x10, 0xec, 0xad, 0x4e, 0x5d, 0x6a, 0xd7, 0x79, 0xd6,
	0x03, 0xb2, 0xee, 0xbb, 0xd2, 0xed, 0x51, 0x4d, 0x65, 0x76, 0xb5, 0x80, 0x4b, 0xeb, 0xb5, 0xe4,
	0x47, 0xf1, 0x67, 0x86, 0x3c, 0x8c, 0x09, 0x56, 0xc5, 0x90, 0xdf, 0x46, 0x86, 0xcf, 0x17, 0x19,
	0x6e, 0xdf, 0xc0, 0xf0, 0x1a, 0xc2, 0x7f, 0x51, 0xac, 0x28, 0xb2, 0x39, 0x77, 0xff, 0x62, 0xa2,
	0x42, 0x5a, 0xef, 0xc8, 0x5a, 0x1c, 0x58, 0xa9, 0x31, 0x6f, 0xb8, 0xb9, 0xb9, 0xbf, 0xef, 0xb9,
	0xf2, 0x89, 0x64, 0x17, 0x56, 0x96, 0xd8, 0x7e, 0x48, 0x6c, 0xcb, 0x4b, 0x6c, 0xff, 0xbc, 0xee,
	0xdc, 0x3f, 0xc1, 0xf9, 0x8e, 0x47, 0x01, 0xe0, 0x71, 0x00, 0xf8, 0x2c, 0x00, 0x74, 0x1e, 0x00,
	0xba, 0x08, 0x00, 0x5d, 0x06, 0x80, 0xae, 0x02, 0xc0, 0x07, 0x06, 0xf0, 0xa1, 0x01, 0x74, 0x64,
	0x00, 0x1f, 0x1b, 0x40, 0x27, 0x06, 0xd0, 0xa9, 0x01, 0x34, 0x32, 0x80, 0xc7, 0x06, 0xf0, 0x99,
	0x01, 0x74, 0x6e, 0x00, 0x5f, 0x18, 0x40, 0x97, 0x06, 0xf0, 0x95, 0x01, 0x74, 0x10, 0x02, 0x3a,
	0x0c, 0x01, 0x7f, 0x09, 0x01, 0x7d, 0x0d, 0x01, 0x7f, 0x0b, 0x01, 0x1d, 0x85, 0x80, 0x8e, 0x43,
	0xc0, 0x27, 0x21, 0xe0, 0xd3, 0x10, 0xf0, 0xfb, 0xa7, 0x9e, 0xb0, 0x75, 0x9b, 0xea, 0x36, 0xe3,
	0x9e, 0xb2, 0x39, 0xd5, 0x43, 0x21, 0x3b, 0xce, 0xfc, 0x2b, 0xe1, 0x77, 0x3c, 0x47, 0x6b, 0xee,
	0x37, 0x1a, 0x77, 0xa2, 0x37, 0x62, 0xe7, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x31, 0x96, 0x69,
	0xc6, 0xf9, 0x04, 0x00, 0x00,
}
