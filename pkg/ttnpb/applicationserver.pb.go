// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: go.thethings.network/lorawan-stack/api/applicationserver.proto

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

type SetApplicationRequest struct {
	Application          Application      `protobuf:"bytes,1,opt,name=application" json:"application"`
	FieldMask            *types.FieldMask `protobuf:"bytes,2,opt,name=field_mask,json=fieldMask" json:"field_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SetApplicationRequest) Reset()      { *m = SetApplicationRequest{} }
func (*SetApplicationRequest) ProtoMessage() {}
func (*SetApplicationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_applicationserver_cf2f36e2a5290989, []int{0}
}
func (m *SetApplicationRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetApplicationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetApplicationRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *SetApplicationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetApplicationRequest.Merge(dst, src)
}
func (m *SetApplicationRequest) XXX_Size() int {
	return m.Size()
}
func (m *SetApplicationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetApplicationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetApplicationRequest proto.InternalMessageInfo

func (m *SetApplicationRequest) GetApplication() Application {
	if m != nil {
		return m.Application
	}
	return Application{}
}

func (m *SetApplicationRequest) GetFieldMask() *types.FieldMask {
	if m != nil {
		return m.FieldMask
	}
	return nil
}

func init() {
	proto.RegisterType((*SetApplicationRequest)(nil), "ttn.lorawan.v3.SetApplicationRequest")
	golang_proto.RegisterType((*SetApplicationRequest)(nil), "ttn.lorawan.v3.SetApplicationRequest")
}
func (this *SetApplicationRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SetApplicationRequest)
	if !ok {
		that2, ok := that.(SetApplicationRequest)
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
	if !this.Application.Equal(&that1.Application) {
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

// Client API for AsApplicationRegistry service

type AsApplicationRegistryClient interface {
	// GetApplication returns the application that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetApplication(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*Application, error)
	// SetApplication creates or updates the application
	SetApplication(ctx context.Context, in *SetApplicationRequest, opts ...grpc.CallOption) (*Application, error)
	// DeleteApplication deletes the application that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteApplication(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error)
}

type asApplicationRegistryClient struct {
	cc *grpc.ClientConn
}

func NewAsApplicationRegistryClient(cc *grpc.ClientConn) AsApplicationRegistryClient {
	return &asApplicationRegistryClient{cc}
}

func (c *asApplicationRegistryClient) GetApplication(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*Application, error) {
	out := new(Application)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.AsApplicationRegistry/GetApplication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asApplicationRegistryClient) SetApplication(ctx context.Context, in *SetApplicationRequest, opts ...grpc.CallOption) (*Application, error) {
	out := new(Application)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.AsApplicationRegistry/SetApplication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asApplicationRegistryClient) DeleteApplication(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.AsApplicationRegistry/DeleteApplication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AsApplicationRegistry service

type AsApplicationRegistryServer interface {
	// GetApplication returns the application that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetApplication(context.Context, *ApplicationIdentifiers) (*Application, error)
	// SetApplication creates or updates the application
	SetApplication(context.Context, *SetApplicationRequest) (*Application, error)
	// DeleteApplication deletes the application that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteApplication(context.Context, *ApplicationIdentifiers) (*types.Empty, error)
}

func RegisterAsApplicationRegistryServer(s *grpc.Server, srv AsApplicationRegistryServer) {
	s.RegisterService(&_AsApplicationRegistry_serviceDesc, srv)
}

func _AsApplicationRegistry_GetApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplicationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsApplicationRegistryServer).GetApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.AsApplicationRegistry/GetApplication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsApplicationRegistryServer).GetApplication(ctx, req.(*ApplicationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _AsApplicationRegistry_SetApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsApplicationRegistryServer).SetApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.AsApplicationRegistry/SetApplication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsApplicationRegistryServer).SetApplication(ctx, req.(*SetApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AsApplicationRegistry_DeleteApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplicationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsApplicationRegistryServer).DeleteApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.AsApplicationRegistry/DeleteApplication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsApplicationRegistryServer).DeleteApplication(ctx, req.(*ApplicationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _AsApplicationRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.AsApplicationRegistry",
	HandlerType: (*AsApplicationRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetApplication",
			Handler:    _AsApplicationRegistry_GetApplication_Handler,
		},
		{
			MethodName: "SetApplication",
			Handler:    _AsApplicationRegistry_SetApplication_Handler,
		},
		{
			MethodName: "DeleteApplication",
			Handler:    _AsApplicationRegistry_DeleteApplication_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "go.thethings.network/lorawan-stack/api/applicationserver.proto",
}

// Client API for As service

type AsClient interface {
	Subscribe(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (As_SubscribeClient, error)
}

type asClient struct {
	cc *grpc.ClientConn
}

func NewAsClient(cc *grpc.ClientConn) AsClient {
	return &asClient{cc}
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
	Subscribe(*ApplicationIdentifiers, As_SubscribeServer) error
}

func RegisterAsServer(s *grpc.Server, srv AsServer) {
	s.RegisterService(&_As_serviceDesc, srv)
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
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _As_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "go.thethings.network/lorawan-stack/api/applicationserver.proto",
}

func (m *SetApplicationRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetApplicationRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintApplicationserver(dAtA, i, uint64(m.Application.Size()))
	n1, err := m.Application.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if m.FieldMask != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintApplicationserver(dAtA, i, uint64(m.FieldMask.Size()))
		n2, err := m.FieldMask.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
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
func NewPopulatedSetApplicationRequest(r randyApplicationserver, easy bool) *SetApplicationRequest {
	this := &SetApplicationRequest{}
	v1 := NewPopulatedApplication(r, easy)
	this.Application = *v1
	if r.Intn(10) != 0 {
		this.FieldMask = types.NewPopulatedFieldMask(r, easy)
	}
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
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
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
		v3 := r.Int63()
		if r.Intn(2) == 0 {
			v3 *= -1
		}
		dAtA = encodeVarintPopulateApplicationserver(dAtA, uint64(v3))
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
func (m *SetApplicationRequest) Size() (n int) {
	var l int
	_ = l
	l = m.Application.Size()
	n += 1 + l + sovApplicationserver(uint64(l))
	if m.FieldMask != nil {
		l = m.FieldMask.Size()
		n += 1 + l + sovApplicationserver(uint64(l))
	}
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
func (this *SetApplicationRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SetApplicationRequest{`,
		`Application:` + strings.Replace(strings.Replace(this.Application.String(), "Application", "Application", 1), `&`, ``, 1) + `,`,
		`FieldMask:` + strings.Replace(fmt.Sprintf("%v", this.FieldMask), "FieldMask", "types.FieldMask", 1) + `,`,
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
func (m *SetApplicationRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SetApplicationRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetApplicationRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Application", wireType)
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
			if err := m.Application.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
			if m.FieldMask == nil {
				m.FieldMask = &types.FieldMask{}
			}
			if err := m.FieldMask.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
	proto.RegisterFile("go.thethings.network/lorawan-stack/api/applicationserver.proto", fileDescriptor_applicationserver_cf2f36e2a5290989)
}
func init() {
	golang_proto.RegisterFile("go.thethings.network/lorawan-stack/api/applicationserver.proto", fileDescriptor_applicationserver_cf2f36e2a5290989)
}

var fileDescriptor_applicationserver_cf2f36e2a5290989 = []byte{
	// 555 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x3f, 0x6c, 0xd3, 0x40,
	0x14, 0xc6, 0xef, 0x15, 0x84, 0x54, 0x57, 0xaa, 0xc0, 0x52, 0x50, 0x95, 0xc2, 0xa3, 0x04, 0x81,
	0x20, 0xa2, 0x67, 0x48, 0x17, 0xe8, 0x80, 0x94, 0xf2, 0x4f, 0x0c, 0x5d, 0x12, 0x58, 0xb2, 0x54,
	0x76, 0x72, 0x71, 0x4e, 0x49, 0x7c, 0xc6, 0x77, 0x69, 0x55, 0x21, 0x50, 0xc4, 0xd4, 0x11, 0x89,
	0x01, 0x46, 0xc4, 0x54, 0xb6, 0x8e, 0x1d, 0x3b, 0xa1, 0x8c, 0x95, 0x58, 0x3a, 0xa1, 0xfa, 0xcc,
	0xd0, 0xb1, 0x63, 0x47, 0x54, 0xc7, 0x21, 0x7f, 0x40, 0x21, 0xb0, 0xdd, 0xf9, 0xbd, 0xef, 0xbb,
	0xdf, 0xfb, 0xee, 0x6c, 0x3c, 0x70, 0x05, 0x55, 0x35, 0xa6, 0x6a, 0xdc, 0x73, 0x25, 0xf5, 0x98,
	0xda, 0x10, 0x41, 0xdd, 0x6a, 0x88, 0xc0, 0xde, 0xb0, 0xbd, 0x45, 0xa9, 0xec, 0x72, 0xdd, 0xb2,
	0x7d, 0x6e, 0xd9, 0xbe, 0xdf, 0xe0, 0x65, 0x5b, 0x71, 0xe1, 0x49, 0x16, 0xac, 0xb3, 0x80, 0xfa,
	0x81, 0x50, 0xc2, 0x9c, 0x55, 0xca, 0xa3, 0x49, 0x3b, 0x5d, 0x5f, 0x4a, 0xdf, 0xfb, 0x77, 0xbf,
	0xae, 0xd3, 0xc4, 0x4a, 0x5e, 0x61, 0x9e, 0xe2, 0x55, 0xce, 0x02, 0x99, 0x28, 0x17, 0x5d, 0xae,
	0x6a, 0x2d, 0x87, 0x96, 0x45, 0xd3, 0x72, 0x85, 0x2b, 0xac, 0xf8, 0xb3, 0xd3, 0xaa, 0xc6, 0xbb,
	0x78, 0x13, 0xaf, 0x92, 0xf6, 0x4b, 0xae, 0x10, 0x6e, 0x83, 0x75, 0x31, 0x3c, 0x4f, 0xa8, 0xee,
	0x54, 0x49, 0x75, 0x3e, 0xa9, 0xfe, 0xf2, 0x60, 0x4d, 0x5f, 0x6d, 0x26, 0xc5, 0x85, 0xd1, 0x62,
	0x95, 0xb3, 0x46, 0x65, 0xad, 0x69, 0xcb, 0x7a, 0xb7, 0x23, 0xf3, 0x01, 0x8c, 0x54, 0x91, 0xa9,
	0x7c, 0x7f, 0xbc, 0x02, 0x7b, 0xd9, 0x62, 0x52, 0x99, 0x0f, 0x8d, 0x99, 0x81, 0xa1, 0xe7, 0x60,
	0x01, 0x6e, 0xce, 0xe4, 0xe6, 0xe9, 0x70, 0x7e, 0x74, 0x40, 0xb8, 0x72, 0xb6, 0xf3, 0xfd, 0x0a,
	0x29, 0x0c, 0xaa, 0xcc, 0xfb, 0x86, 0xd1, 0x3f, 0x72, 0x6e, 0x2a, 0xf6, 0x48, 0xd3, 0x2e, 0x15,
	0xed, 0x51, 0xd1, 0x27, 0xa7, 0x2d, 0xab, 0xb6, 0xac, 0x17, 0xa6, 0xab, 0xbd, 0x65, 0xee, 0xeb,
	0x19, 0x23, 0x95, 0x97, 0x43, 0x60, 0x2e, 0x97, 0x2a, 0xd8, 0x34, 0xdb, 0x60, 0xcc, 0x3e, 0x1d,
	0x62, 0x36, 0x6f, 0x8c, 0xe1, 0x7a, 0xd6, 0xbf, 0x80, 0xf4, 0x38, 0xfe, 0xcc, 0xad, 0xb7, 0xdf,
	0x7e, 0xbc, 0x9f, 0xba, 0x66, 0x5e, 0xb5, 0x6c, 0x39, 0xf4, 0x80, 0xac, 0x57, 0x03, 0xbb, 0x35,
	0x5e, 0x79, 0x6d, 0x7e, 0x01, 0x63, 0x76, 0x38, 0x36, 0xf3, 0xfa, 0xa8, 0xf5, 0x1f, 0x63, 0x1d,
	0x4f, 0x50, 0x8c, 0x09, 0x56, 0x33, 0x77, 0xc7, 0x12, 0x50, 0x5e, 0x91, 0x74, 0x84, 0x68, 0x19,
	0xb2, 0xa5, 0x54, 0xe6, 0xfc, 0xa8, 0x6e, 0x19, 0xb2, 0xe6, 0x1b, 0xe3, 0xc2, 0x23, 0xd6, 0x60,
	0x8a, 0xfd, 0x4f, 0x60, 0x17, 0x7f, 0xbb, 0xac, 0xc7, 0xa7, 0xef, 0xab, 0x97, 0x55, 0xf6, 0xef,
	0x59, 0xe5, 0x4a, 0xc6, 0x54, 0x5e, 0x9a, 0xcf, 0x8d, 0xe9, 0x62, 0xcb, 0x91, 0xe5, 0x80, 0x3b,
	0x6c, 0xe2, 0xd3, 0x2f, 0x8f, 0xe9, 0x7b, 0xe1, 0xdf, 0x81, 0x95, 0xcf, 0xd0, 0x09, 0x11, 0xf6,
	0x43, 0x84, 0x83, 0x10, 0xc9, 0x61, 0x88, 0xe4, 0x28, 0x44, 0x72, 0x1c, 0x22, 0x39, 0x09, 0x11,
	0xda, 0x1a, 0x61, 0x4b, 0x23, 0xd9, 0xd6, 0x08, 0x3b, 0x1a, 0xc9, 0xae, 0x46, 0xb2, 0xa7, 0x91,
	0x74, 0x34, 0xc2, 0xbe, 0x46, 0x38, 0xd0, 0x48, 0x0e, 0x35, 0xc2, 0x91, 0x46, 0x72, 0xac, 0x11,
	0x4e, 0x34, 0x92, 0x76, 0x84, 0x64, 0x2b, 0x42, 0x78, 0x17, 0x21, 0xf9, 0x18, 0x21, 0x7c, 0x8a,
	0x90, 0x6c, 0x47, 0x48, 0x76, 0x22, 0x84, 0xdd, 0x08, 0x61, 0x2f, 0x42, 0x28, 0xdd, 0x9e, 0xe0,
	0xdf, 0xf7, 0xeb, 0xae, 0xa5, 0x94, 0xe7, 0x3b, 0xce, 0xb9, 0x38, 0xbb, 0xa5, 0x9f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xe0, 0xf9, 0x7f, 0x03, 0xbc, 0x04, 0x00, 0x00,
}
