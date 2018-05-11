// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/applicationserver.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf3 "github.com/gogo/protobuf/types"
import google_protobuf4 "github.com/gogo/protobuf/types"

import context "context"
import grpc "google.golang.org/grpc"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type SetApplicationRequest struct {
	Application Application                 `protobuf:"bytes,1,opt,name=application" json:"application"`
	FieldMask   *google_protobuf4.FieldMask `protobuf:"bytes,2,opt,name=field_mask,json=fieldMask" json:"field_mask,omitempty"`
}

func (m *SetApplicationRequest) Reset()         { *m = SetApplicationRequest{} }
func (m *SetApplicationRequest) String() string { return proto.CompactTextString(m) }
func (*SetApplicationRequest) ProtoMessage()    {}
func (*SetApplicationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptorApplicationserver, []int{0}
}

func (m *SetApplicationRequest) GetApplication() Application {
	if m != nil {
		return m.Application
	}
	return Application{}
}

func (m *SetApplicationRequest) GetFieldMask() *google_protobuf4.FieldMask {
	if m != nil {
		return m.FieldMask
	}
	return nil
}

func init() {
	proto.RegisterType((*SetApplicationRequest)(nil), "ttn.lorawan.v3.SetApplicationRequest")
	golang_proto.RegisterType((*SetApplicationRequest)(nil), "ttn.lorawan.v3.SetApplicationRequest")
}
func (this *SetApplicationRequest) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*SetApplicationRequest)
	if !ok {
		that2, ok := that.(SetApplicationRequest)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *SetApplicationRequest")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *SetApplicationRequest but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *SetApplicationRequest but is not nil && this == nil")
	}
	if !this.Application.Equal(&that1.Application) {
		return fmt.Errorf("Application this(%v) Not Equal that(%v)", this.Application, that1.Application)
	}
	if !this.FieldMask.Equal(that1.FieldMask) {
		return fmt.Errorf("FieldMask this(%v) Not Equal that(%v)", this.FieldMask, that1.FieldMask)
	}
	return nil
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
	SetApplication(ctx context.Context, in *SetApplicationRequest, opts ...grpc.CallOption) (*google_protobuf3.Empty, error)
	// DeleteApplication deletes the application that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteApplication(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*google_protobuf3.Empty, error)
}

type asApplicationRegistryClient struct {
	cc *grpc.ClientConn
}

func NewAsApplicationRegistryClient(cc *grpc.ClientConn) AsApplicationRegistryClient {
	return &asApplicationRegistryClient{cc}
}

func (c *asApplicationRegistryClient) GetApplication(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*Application, error) {
	out := new(Application)
	err := grpc.Invoke(ctx, "/ttn.lorawan.v3.AsApplicationRegistry/GetApplication", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asApplicationRegistryClient) SetApplication(ctx context.Context, in *SetApplicationRequest, opts ...grpc.CallOption) (*google_protobuf3.Empty, error) {
	out := new(google_protobuf3.Empty)
	err := grpc.Invoke(ctx, "/ttn.lorawan.v3.AsApplicationRegistry/SetApplication", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asApplicationRegistryClient) DeleteApplication(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*google_protobuf3.Empty, error) {
	out := new(google_protobuf3.Empty)
	err := grpc.Invoke(ctx, "/ttn.lorawan.v3.AsApplicationRegistry/DeleteApplication", in, out, c.cc, opts...)
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
	SetApplication(context.Context, *SetApplicationRequest) (*google_protobuf3.Empty, error)
	// DeleteApplication deletes the application that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteApplication(context.Context, *ApplicationIdentifiers) (*google_protobuf3.Empty, error)
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
	Metadata: "github.com/TheThingsNetwork/ttn/api/applicationserver.proto",
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
	stream, err := grpc.NewClientStream(ctx, &_As_serviceDesc.Streams[0], c.cc, "/ttn.lorawan.v3.As/Subscribe", opts...)
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
	Metadata: "github.com/TheThingsNetwork/ttn/api/applicationserver.proto",
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
		this.FieldMask = google_protobuf4.NewPopulatedFieldMask(r, easy)
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
				m.FieldMask = &google_protobuf4.FieldMask{}
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
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/applicationserver.proto", fileDescriptorApplicationserver)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/applicationserver.proto", fileDescriptorApplicationserver)
}

var fileDescriptorApplicationserver = []byte{
	// 552 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x3d, 0x6c, 0xd3, 0x40,
	0x14, 0xc7, 0xef, 0x05, 0x84, 0x54, 0x57, 0x8a, 0xc0, 0x52, 0x50, 0x95, 0xc2, 0xa3, 0x04, 0x81,
	0x68, 0x24, 0x6c, 0x48, 0xc5, 0x40, 0x99, 0x52, 0xbe, 0xc4, 0x50, 0x86, 0x24, 0x2c, 0x59, 0x2a,
	0x3b, 0xb9, 0x38, 0xa7, 0x24, 0x3e, 0xe3, 0xbb, 0xb4, 0xaa, 0x10, 0x28, 0x62, 0xea, 0x88, 0xc4,
	0x00, 0x1b, 0x0c, 0x0c, 0x1d, 0x3b, 0x76, 0x2c, 0x5b, 0xc6, 0x4a, 0x2c, 0x9d, 0x50, 0x7d, 0x66,
	0xe8, 0xd8, 0xb1, 0x23, 0x8a, 0xe3, 0x92, 0x0f, 0x20, 0x2d, 0x9d, 0x7c, 0xcf, 0xef, 0xfd, 0xdf,
	0xfb, 0xbd, 0xff, 0xd9, 0xda, 0x43, 0x87, 0xc9, 0x7a, 0xdb, 0x36, 0x2a, 0xbc, 0x65, 0x96, 0xea,
	0xb4, 0x54, 0x67, 0xae, 0x23, 0x5e, 0x50, 0xb9, 0xc6, 0xfd, 0x86, 0x29, 0xa5, 0x6b, 0x5a, 0x1e,
	0x33, 0x2d, 0xcf, 0x6b, 0xb2, 0x8a, 0x25, 0x19, 0x77, 0x05, 0xf5, 0x57, 0xa9, 0x6f, 0x78, 0x3e,
	0x97, 0x5c, 0x4f, 0x4a, 0xe9, 0x1a, 0x4d, 0xee, 0x5b, 0x6b, 0x96, 0x6b, 0xac, 0x2e, 0xa4, 0xef,
	0xff, 0x67, 0xb3, 0x7e, 0x9b, 0xd3, 0xc9, 0x58, 0x95, 0xba, 0x92, 0xd5, 0x18, 0xf5, 0x45, 0x2c,
	0xbb, 0x33, 0x24, 0x73, 0xb8, 0xc3, 0xcd, 0xe8, 0xb5, 0xdd, 0xae, 0x45, 0x51, 0x14, 0x44, 0xa7,
	0xb8, 0xfc, 0x8a, 0xc3, 0xb9, 0xd3, 0xa4, 0x7d, 0x06, 0xd7, 0xe5, 0xb2, 0xbf, 0x4f, 0x9c, 0x9d,
	0x8d, 0xb3, 0xbf, 0x7b, 0xd0, 0x96, 0x27, 0xd7, 0xe3, 0xe4, 0xdc, 0x78, 0xb2, 0xc6, 0x68, 0xb3,
	0xba, 0xd2, 0xb2, 0x44, 0xa3, 0x5f, 0x91, 0xf9, 0x08, 0x5a, 0xaa, 0x48, 0x65, 0x7e, 0xb0, 0x5b,
	0x81, 0xbe, 0x6a, 0x53, 0x21, 0xf5, 0x47, 0xda, 0xf4, 0xd0, 0xc6, 0x33, 0x30, 0x07, 0xb7, 0xa7,
	0x73, 0xb3, 0xc6, 0xa8, 0x73, 0xc6, 0x90, 0x70, 0xe9, 0x7c, 0xf7, 0xc7, 0x35, 0x52, 0x18, 0x56,
	0xe9, 0x0f, 0x34, 0x6d, 0x30, 0x72, 0x26, 0x11, 0xf5, 0x48, 0x1b, 0x7d, 0x2a, 0xe3, 0x98, 0xca,
	0x78, 0xda, 0x2b, 0x59, 0xb6, 0x44, 0xa3, 0x30, 0x55, 0x3b, 0x3e, 0xe6, 0xbe, 0x9d, 0xd3, 0x52,
	0x79, 0x31, 0x02, 0xe6, 0x30, 0x21, 0xfd, 0x75, 0xbd, 0x03, 0x5a, 0xf2, 0xd9, 0x08, 0xb3, 0x7e,
	0x6b, 0x02, 0xd7, 0xf3, 0xc1, 0x05, 0xa4, 0x27, 0xf1, 0x67, 0xe6, 0xdf, 0x7d, 0xff, 0xf9, 0x21,
	0x71, 0x43, 0xbf, 0x6e, 0x5a, 0x62, 0xe4, 0xd3, 0x31, 0x5f, 0x0f, 0x45, 0x2b, 0xac, 0xfa, 0x46,
	0xff, 0x0a, 0x5a, 0x72, 0xd4, 0x36, 0xfd, 0xe6, 0x78, 0xeb, 0xbf, 0xda, 0x9a, 0xbe, 0xfc, 0xc7,
	0xf6, 0x4f, 0x7a, 0x17, 0x96, 0x29, 0x46, 0xc3, 0x97, 0x33, 0xf7, 0x26, 0x0e, 0x37, 0x58, 0x55,
	0x18, 0x63, 0x30, 0x8b, 0x90, 0x2d, 0xa7, 0x32, 0x17, 0xc7, 0x75, 0x8b, 0x90, 0xd5, 0xdf, 0x6a,
	0x97, 0x1e, 0xd3, 0x26, 0x95, 0xf4, 0x2c, 0x5e, 0xfd, 0x8b, 0x34, 0xb6, 0x29, 0x7b, 0xb2, 0x4d,
	0xb9, 0xb2, 0x96, 0xc8, 0x0b, 0xbd, 0xa4, 0x4d, 0x15, 0xdb, 0xb6, 0xa8, 0xf8, 0xcc, 0xa6, 0xa7,
	0x9e, 0x7e, 0x75, 0x42, 0xdd, 0x4b, 0xef, 0x2e, 0x2c, 0x7d, 0x86, 0x6e, 0x80, 0xb0, 0x1b, 0x20,
	0xec, 0x05, 0x08, 0xfb, 0x01, 0xc2, 0x41, 0x80, 0xe4, 0x30, 0x40, 0x72, 0x14, 0x20, 0x74, 0x14,
	0x92, 0x0d, 0x85, 0x64, 0x53, 0x21, 0x6c, 0x29, 0x24, 0xdb, 0x0a, 0x61, 0x47, 0x21, 0x74, 0x15,
	0xc2, 0xae, 0x42, 0xd8, 0x53, 0x48, 0xf6, 0x15, 0xc2, 0x81, 0x42, 0x72, 0xa8, 0x10, 0x8e, 0x14,
	0x92, 0x4e, 0x88, 0x64, 0x23, 0x44, 0x78, 0x1f, 0x22, 0xf9, 0x14, 0x22, 0x7c, 0x09, 0x91, 0x6c,
	0x86, 0x48, 0xb6, 0x42, 0x84, 0xed, 0x10, 0x61, 0x27, 0x44, 0x28, 0xcf, 0x9f, 0xf4, 0xcf, 0x7b,
	0x0d, 0xa7, 0xf7, 0xf4, 0x6c, 0xfb, 0x42, 0x64, 0xdc, 0xc2, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xc6, 0xf0, 0xb6, 0xc2, 0xab, 0x04, 0x00, 0x00,
}
