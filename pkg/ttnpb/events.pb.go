// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/events.proto

package ttnpb

import (
	bytes "bytes"
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"
	time "time"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	types "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	_ "github.com/lyft/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Event struct {
	Name                 string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Time                 time.Time            `protobuf:"bytes,2,opt,name=time,proto3,stdtime" json:"time"`
	Identifiers          []*EntityIdentifiers `protobuf:"bytes,3,rep,name=identifiers,proto3" json:"identifiers,omitempty"`
	Data                 *types.Any           `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	CorrelationIDs       []string             `protobuf:"bytes,5,rep,name=correlation_ids,json=correlationIds,proto3" json:"correlation_ids,omitempty"`
	Origin               string               `protobuf:"bytes,6,opt,name=origin,proto3" json:"origin,omitempty"`
	Context              map[string][]byte    `protobuf:"bytes,7,rep,name=context,proto3" json:"context,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Event) Reset()      { *m = Event{} }
func (*Event) ProtoMessage() {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fd8551d68f51e44, []int{0}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Event.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return m.Size()
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Event) GetTime() time.Time {
	if m != nil {
		return m.Time
	}
	return time.Time{}
}

func (m *Event) GetIdentifiers() []*EntityIdentifiers {
	if m != nil {
		return m.Identifiers
	}
	return nil
}

func (m *Event) GetData() *types.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Event) GetCorrelationIDs() []string {
	if m != nil {
		return m.CorrelationIDs
	}
	return nil
}

func (m *Event) GetOrigin() string {
	if m != nil {
		return m.Origin
	}
	return ""
}

func (m *Event) GetContext() map[string][]byte {
	if m != nil {
		return m.Context
	}
	return nil
}

type StreamEventsRequest struct {
	Identifiers []*EntityIdentifiers `protobuf:"bytes,1,rep,name=identifiers,proto3" json:"identifiers,omitempty"`
	// If greater than zero, this will return historical events, up to this maximum when the stream starts.
	// If used in combination with "after", the limit that is reached first, is used.
	// The availability of historical events depends on server support and retention policy.
	Tail uint32 `protobuf:"varint,2,opt,name=tail,proto3" json:"tail,omitempty"`
	// If not empty, this will return historical events after the given time when the stream starts.
	// If used in combination with "tail", the limit that is reached first, is used.
	// The availability of historical events depends on server support and retention policy.
	After                *time.Time `protobuf:"bytes,3,opt,name=after,proto3,stdtime" json:"after,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *StreamEventsRequest) Reset()      { *m = StreamEventsRequest{} }
func (*StreamEventsRequest) ProtoMessage() {}
func (*StreamEventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fd8551d68f51e44, []int{1}
}
func (m *StreamEventsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StreamEventsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StreamEventsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StreamEventsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamEventsRequest.Merge(m, src)
}
func (m *StreamEventsRequest) XXX_Size() int {
	return m.Size()
}
func (m *StreamEventsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamEventsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StreamEventsRequest proto.InternalMessageInfo

func (m *StreamEventsRequest) GetIdentifiers() []*EntityIdentifiers {
	if m != nil {
		return m.Identifiers
	}
	return nil
}

func (m *StreamEventsRequest) GetTail() uint32 {
	if m != nil {
		return m.Tail
	}
	return 0
}

func (m *StreamEventsRequest) GetAfter() *time.Time {
	if m != nil {
		return m.After
	}
	return nil
}

func init() {
	proto.RegisterType((*Event)(nil), "ttn.lorawan.v3.Event")
	golang_proto.RegisterType((*Event)(nil), "ttn.lorawan.v3.Event")
	proto.RegisterMapType((map[string][]byte)(nil), "ttn.lorawan.v3.Event.ContextEntry")
	golang_proto.RegisterMapType((map[string][]byte)(nil), "ttn.lorawan.v3.Event.ContextEntry")
	proto.RegisterType((*StreamEventsRequest)(nil), "ttn.lorawan.v3.StreamEventsRequest")
	golang_proto.RegisterType((*StreamEventsRequest)(nil), "ttn.lorawan.v3.StreamEventsRequest")
}

func init() { proto.RegisterFile("lorawan-stack/api/events.proto", fileDescriptor_4fd8551d68f51e44) }
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/events.proto", fileDescriptor_4fd8551d68f51e44)
}

var fileDescriptor_4fd8551d68f51e44 = []byte{
	// 666 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x3f, 0x6c, 0xd3, 0x40,
	0x14, 0xc6, 0xef, 0xf2, 0xaf, 0xea, 0xb5, 0x94, 0xea, 0x28, 0xc8, 0x44, 0xe8, 0x25, 0xa4, 0x4b,
	0x84, 0x88, 0x83, 0x5a, 0x09, 0xa1, 0xc2, 0x42, 0x4a, 0x87, 0xac, 0x86, 0xa9, 0x0b, 0xba, 0x24,
	0x17, 0xd7, 0x4a, 0x72, 0x17, 0xec, 0x4b, 0x4a, 0xb6, 0x8a, 0xa9, 0x62, 0xaa, 0x60, 0x61, 0x44,
	0x0c, 0xa8, 0x12, 0x4b, 0xc7, 0x8a, 0xa9, 0x63, 0xc7, 0x4a, 0x2c, 0x9d, 0x4a, 0x6d, 0x33, 0x94,
	0xad, 0x63, 0x47, 0x94, 0xb3, 0x4b, 0xd3, 0xa6, 0x12, 0x12, 0xdb, 0x7b, 0x7a, 0xdf, 0xcb, 0xf7,
	0xbd, 0x5f, 0x7c, 0x04, 0xda, 0xd2, 0x65, 0xeb, 0x4c, 0x94, 0x3c, 0xc5, 0xea, 0xad, 0x32, 0xeb,
	0x3a, 0x65, 0xde, 0xe7, 0x42, 0x79, 0x66, 0xd7, 0x95, 0x4a, 0xd2, 0x19, 0xa5, 0x84, 0x19, 0x6b,
	0xcc, 0xfe, 0x62, 0xb6, 0x64, 0x3b, 0x6a, 0xad, 0x57, 0x33, 0xeb, 0xb2, 0x53, 0xb6, 0xa5, 0x2d,
	0xcb, 0x5a, 0x56, 0xeb, 0x35, 0x75, 0xa7, 0x1b, 0x5d, 0x45, 0xeb, 0xd9, 0xa7, 0x23, 0xf2, 0xf6,
	0xa0, 0xa9, 0x22, 0x79, 0xbd, 0x64, 0x73, 0x51, 0xea, 0xb3, 0xb6, 0xd3, 0x60, 0x8a, 0x97, 0xc7,
	0x8a, 0x78, 0xf9, 0x9e, 0x2d, 0xa5, 0xdd, 0xe6, 0x3a, 0x14, 0x13, 0x42, 0x2a, 0xa6, 0x1c, 0x29,
	0xe2, 0x64, 0xd9, 0xbb, 0xf1, 0xf4, 0x6f, 0x00, 0x26, 0x06, 0xf1, 0x28, 0x77, 0x75, 0xa4, 0x9c,
	0x0e, 0xf7, 0x14, 0xeb, 0x74, 0x63, 0xc1, 0xfc, 0xf8, 0xd5, 0x4e, 0x83, 0x0b, 0xe5, 0x34, 0x1d,
	0xee, 0xc6, 0x06, 0x85, 0x6f, 0x49, 0x92, 0x5e, 0x19, 0xb2, 0xa0, 0x94, 0xa4, 0x04, 0xeb, 0x70,
	0x03, 0xe7, 0x71, 0x71, 0xd2, 0xd2, 0x35, 0xad, 0x90, 0xd4, 0xf0, 0x57, 0x8d, 0x44, 0x1e, 0x17,
	0xa7, 0x16, 0xb2, 0x66, 0x64, 0x69, 0x9e, 0x5b, 0x9a, 0xaf, 0xce, 0x2d, 0x2b, 0x74, 0xff, 0x28,
	0x87, 0xb6, 0x7e, 0xe6, 0xf0, 0xf7, 0xdf, 0x7b, 0xc9, 0xf4, 0x7b, 0x9c, 0x98, 0xc5, 0x96, 0xde,
	0xa5, 0xcb, 0x64, 0x6a, 0xc4, 0xd6, 0x48, 0xe6, 0x93, 0xc5, 0xa9, 0x85, 0xfb, 0xe6, 0x65, 0xe4,
	0xe6, 0x8a, 0x50, 0x8e, 0x1a, 0x54, 0x2f, 0x84, 0xd6, 0xe8, 0x16, 0x2d, 0x92, 0x54, 0x83, 0x29,
	0x66, 0xa4, 0x74, 0x90, 0xb9, 0xb1, 0x20, 0xcf, 0xc5, 0xc0, 0xd2, 0x0a, 0x5a, 0x25, 0x37, 0xeb,
	0xd2, 0x75, 0x79, 0x5b, 0x73, 0x7c, 0xed, 0x34, 0x3c, 0x23, 0x9d, 0x4f, 0x16, 0x27, 0x2b, 0xf9,
	0xe0, 0x28, 0x37, 0xb3, 0x7c, 0x31, 0xaa, 0xbe, 0xf0, 0x86, 0x59, 0x27, 0x3f, 0xe0, 0x4c, 0x21,
	0xe5, 0x26, 0x8c, 0x86, 0x35, 0x33, 0xb2, 0x58, 0x6d, 0x78, 0xf4, 0x0e, 0xc9, 0x48, 0xd7, 0xb1,
	0x1d, 0x61, 0x64, 0x34, 0x93, 0xb8, 0xa3, 0xcf, 0xc8, 0x44, 0x5d, 0x0a, 0xc5, 0xdf, 0x2a, 0x63,
	0x42, 0x5f, 0x53, 0x18, 0xbb, 0x66, 0x48, 0xd4, 0x5c, 0x8e, 0x44, 0x2b, 0x42, 0xb9, 0x03, 0xeb,
	0x7c, 0x25, 0xbb, 0x44, 0xa6, 0x47, 0x07, 0x74, 0x96, 0x24, 0x5b, 0x7c, 0x10, 0x63, 0x1f, 0x96,
	0x74, 0x8e, 0xa4, 0xfb, 0xac, 0xdd, 0x8b, 0xb0, 0x4f, 0x5b, 0x51, 0xb3, 0x94, 0x78, 0x82, 0x0b,
	0x5f, 0x31, 0xb9, 0xf5, 0x52, 0xb9, 0x9c, 0x75, 0xb4, 0x83, 0x67, 0xf1, 0x37, 0x3d, 0xee, 0xa9,
	0xab, 0x8c, 0xf1, 0x7f, 0x31, 0xa6, 0x24, 0xa5, 0x98, 0xd3, 0xd6, 0xae, 0x37, 0x2c, 0x5d, 0xd3,
	0xc7, 0x24, 0xcd, 0x9a, 0x8a, 0xbb, 0x46, 0xf2, 0x9f, 0x5f, 0x40, 0x6a, 0xf8, 0xef, 0x5b, 0x91,
	0x7c, 0xa1, 0x41, 0x32, 0x51, 0x42, 0xba, 0x4a, 0x32, 0x51, 0x62, 0x3a, 0x7f, 0x35, 0xcf, 0x35,
	0x97, 0x64, 0x6f, 0x5f, 0x8b, 0xb2, 0x40, 0xdf, 0xfd, 0xf8, 0xf5, 0x31, 0x31, 0x5d, 0x98, 0x88,
	0x1f, 0xee, 0x12, 0x7e, 0xf0, 0x08, 0x57, 0xbe, 0xe0, 0x7d, 0x1f, 0xf0, 0x81, 0x0f, 0xf8, 0xd0,
	0x07, 0x74, 0xec, 0x03, 0x3a, 0xf1, 0x01, 0x9d, 0xfa, 0x80, 0xce, 0x7c, 0xc0, 0x1b, 0x01, 0xe0,
	0xcd, 0x00, 0xd0, 0x76, 0x00, 0x78, 0x27, 0x00, 0xb4, 0x1b, 0x00, 0xda, 0x0b, 0x00, 0xed, 0x07,
	0x80, 0x0f, 0x02, 0xc0, 0x87, 0x01, 0xa0, 0xe3, 0x00, 0xf0, 0x49, 0x00, 0xe8, 0x34, 0x00, 0x7c,
	0x16, 0x00, 0xda, 0x08, 0x01, 0x6d, 0x86, 0x80, 0xb7, 0x42, 0x40, 0x9f, 0x42, 0xc0, 0x9f, 0x43,
	0x40, 0xdb, 0x21, 0xa0, 0x9d, 0x10, 0xf0, 0x6e, 0x08, 0x78, 0x2f, 0x04, 0xbc, 0xfa, 0xd0, 0x96,
	0xa6, 0x5a, 0xe3, 0x6a, 0xcd, 0x11, 0xb6, 0x67, 0x0a, 0xae, 0xd6, 0xa5, 0xdb, 0x2a, 0x5f, 0x7e,
	0x6d, 0xdd, 0x96, 0x5d, 0x56, 0x4a, 0x74, 0x6b, 0xb5, 0x8c, 0x66, 0xb5, 0xf8, 0x27, 0x00, 0x00,
	0xff, 0xff, 0x72, 0xe1, 0x3c, 0x1e, 0x85, 0x04, 0x00, 0x00,
}

func (this *Event) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Event)
	if !ok {
		that2, ok := that.(Event)
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
	if this.Name != that1.Name {
		return false
	}
	if !this.Time.Equal(that1.Time) {
		return false
	}
	if len(this.Identifiers) != len(that1.Identifiers) {
		return false
	}
	for i := range this.Identifiers {
		if !this.Identifiers[i].Equal(that1.Identifiers[i]) {
			return false
		}
	}
	if !this.Data.Equal(that1.Data) {
		return false
	}
	if len(this.CorrelationIDs) != len(that1.CorrelationIDs) {
		return false
	}
	for i := range this.CorrelationIDs {
		if this.CorrelationIDs[i] != that1.CorrelationIDs[i] {
			return false
		}
	}
	if this.Origin != that1.Origin {
		return false
	}
	if len(this.Context) != len(that1.Context) {
		return false
	}
	for i := range this.Context {
		if !bytes.Equal(this.Context[i], that1.Context[i]) {
			return false
		}
	}
	return true
}
func (this *StreamEventsRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*StreamEventsRequest)
	if !ok {
		that2, ok := that.(StreamEventsRequest)
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
	if len(this.Identifiers) != len(that1.Identifiers) {
		return false
	}
	for i := range this.Identifiers {
		if !this.Identifiers[i].Equal(that1.Identifiers[i]) {
			return false
		}
	}
	if this.Tail != that1.Tail {
		return false
	}
	if that1.After == nil {
		if this.After != nil {
			return false
		}
	} else if !this.After.Equal(*that1.After) {
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

// EventsClient is the client API for Events service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventsClient interface {
	// Stream live events, optionally with a tail of historical events (depending on server support and retention policy).
	// Events may arrive out-of-order.
	Stream(ctx context.Context, in *StreamEventsRequest, opts ...grpc.CallOption) (Events_StreamClient, error)
}

type eventsClient struct {
	cc *grpc.ClientConn
}

func NewEventsClient(cc *grpc.ClientConn) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) Stream(ctx context.Context, in *StreamEventsRequest, opts ...grpc.CallOption) (Events_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Events_serviceDesc.Streams[0], "/ttn.lorawan.v3.Events/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Events_StreamClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventsStreamClient struct {
	grpc.ClientStream
}

func (x *eventsStreamClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventsServer is the server API for Events service.
type EventsServer interface {
	// Stream live events, optionally with a tail of historical events (depending on server support and retention policy).
	// Events may arrive out-of-order.
	Stream(*StreamEventsRequest, Events_StreamServer) error
}

func RegisterEventsServer(s *grpc.Server, srv EventsServer) {
	s.RegisterService(&_Events_serviceDesc, srv)
}

func _Events_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamEventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventsServer).Stream(m, &eventsStreamServer{stream})
}

type Events_StreamServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type eventsStreamServer struct {
	grpc.ServerStream
}

func (x *eventsStreamServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

var _Events_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.Events",
	HandlerType: (*EventsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Events_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "lorawan-stack/api/events.proto",
}

func (m *Event) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Event) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintEvents(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.Time)))
	n1, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.Time, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Identifiers) > 0 {
		for _, msg := range m.Identifiers {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintEvents(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Data != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintEvents(dAtA, i, uint64(m.Data.Size()))
		n2, err := m.Data.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if len(m.CorrelationIDs) > 0 {
		for _, s := range m.CorrelationIDs {
			dAtA[i] = 0x2a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if len(m.Origin) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Origin)))
		i += copy(dAtA[i:], m.Origin)
	}
	if len(m.Context) > 0 {
		for k := range m.Context {
			dAtA[i] = 0x3a
			i++
			v := m.Context[k]
			byteSize := 0
			if len(v) > 0 {
				byteSize = 1 + len(v) + sovEvents(uint64(len(v)))
			}
			mapSize := 1 + len(k) + sovEvents(uint64(len(k))) + byteSize
			i = encodeVarintEvents(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintEvents(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			if len(v) > 0 {
				dAtA[i] = 0x12
				i++
				i = encodeVarintEvents(dAtA, i, uint64(len(v)))
				i += copy(dAtA[i:], v)
			}
		}
	}
	return i, nil
}

func (m *StreamEventsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StreamEventsRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Identifiers) > 0 {
		for _, msg := range m.Identifiers {
			dAtA[i] = 0xa
			i++
			i = encodeVarintEvents(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Tail != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintEvents(dAtA, i, uint64(m.Tail))
	}
	if m.After != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintEvents(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.After)))
		n3, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.After, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedEvent(r randyEvents, easy bool) *Event {
	this := &Event{}
	this.Name = randStringEvents(r)
	v1 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.Time = *v1
	if r.Intn(10) != 0 {
		v2 := r.Intn(5)
		this.Identifiers = make([]*EntityIdentifiers, v2)
		for i := 0; i < v2; i++ {
			this.Identifiers[i] = NewPopulatedEntityIdentifiers(r, easy)
		}
	}
	if r.Intn(10) != 0 {
		this.Data = types.NewPopulatedAny(r, easy)
	}
	v3 := r.Intn(10)
	this.CorrelationIDs = make([]string, v3)
	for i := 0; i < v3; i++ {
		this.CorrelationIDs[i] = randStringEvents(r)
	}
	this.Origin = randStringEvents(r)
	if r.Intn(10) != 0 {
		v4 := r.Intn(10)
		this.Context = make(map[string][]byte)
		for i := 0; i < v4; i++ {
			v5 := r.Intn(100)
			v6 := randStringEvents(r)
			this.Context[v6] = make([]byte, v5)
			for i := 0; i < v5; i++ {
				this.Context[v6][i] = byte(r.Intn(256))
			}
		}
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedStreamEventsRequest(r randyEvents, easy bool) *StreamEventsRequest {
	this := &StreamEventsRequest{}
	if r.Intn(10) != 0 {
		v7 := r.Intn(5)
		this.Identifiers = make([]*EntityIdentifiers, v7)
		for i := 0; i < v7; i++ {
			this.Identifiers[i] = NewPopulatedEntityIdentifiers(r, easy)
		}
	}
	this.Tail = r.Uint32()
	if r.Intn(10) != 0 {
		this.After = github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyEvents interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneEvents(r randyEvents) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringEvents(r randyEvents) string {
	v8 := r.Intn(100)
	tmps := make([]rune, v8)
	for i := 0; i < v8; i++ {
		tmps[i] = randUTF8RuneEvents(r)
	}
	return string(tmps)
}
func randUnrecognizedEvents(r randyEvents, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldEvents(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldEvents(dAtA []byte, r randyEvents, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateEvents(dAtA, uint64(key))
		v9 := r.Int63()
		if r.Intn(2) == 0 {
			v9 *= -1
		}
		dAtA = encodeVarintPopulateEvents(dAtA, uint64(v9))
	case 1:
		dAtA = encodeVarintPopulateEvents(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateEvents(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateEvents(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateEvents(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateEvents(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *Event) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovEvents(uint64(l))
	if len(m.Identifiers) > 0 {
		for _, e := range m.Identifiers {
			l = e.Size()
			n += 1 + l + sovEvents(uint64(l))
		}
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovEvents(uint64(l))
	}
	if len(m.CorrelationIDs) > 0 {
		for _, s := range m.CorrelationIDs {
			l = len(s)
			n += 1 + l + sovEvents(uint64(l))
		}
	}
	l = len(m.Origin)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if len(m.Context) > 0 {
		for k, v := range m.Context {
			_ = k
			_ = v
			l = 0
			if len(v) > 0 {
				l = 1 + len(v) + sovEvents(uint64(len(v)))
			}
			mapEntrySize := 1 + len(k) + sovEvents(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovEvents(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *StreamEventsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Identifiers) > 0 {
		for _, e := range m.Identifiers {
			l = e.Size()
			n += 1 + l + sovEvents(uint64(l))
		}
	}
	if m.Tail != 0 {
		n += 1 + sovEvents(uint64(m.Tail))
	}
	if m.After != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.After)
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozEvents(x uint64) (n int) {
	return sovEvents((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *Event) String() string {
	if this == nil {
		return "nil"
	}
	keysForContext := make([]string, 0, len(this.Context))
	for k := range this.Context {
		keysForContext = append(keysForContext, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForContext)
	mapStringForContext := "map[string][]byte{"
	for _, k := range keysForContext {
		mapStringForContext += fmt.Sprintf("%v: %v,", k, this.Context[k])
	}
	mapStringForContext += "}"
	s := strings.Join([]string{`&Event{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Time:` + strings.Replace(strings.Replace(this.Time.String(), "Timestamp", "types.Timestamp", 1), `&`, ``, 1) + `,`,
		`Identifiers:` + strings.Replace(fmt.Sprintf("%v", this.Identifiers), "EntityIdentifiers", "EntityIdentifiers", 1) + `,`,
		`Data:` + strings.Replace(fmt.Sprintf("%v", this.Data), "Any", "types.Any", 1) + `,`,
		`CorrelationIDs:` + fmt.Sprintf("%v", this.CorrelationIDs) + `,`,
		`Origin:` + fmt.Sprintf("%v", this.Origin) + `,`,
		`Context:` + mapStringForContext + `,`,
		`}`,
	}, "")
	return s
}
func (this *StreamEventsRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&StreamEventsRequest{`,
		`Identifiers:` + strings.Replace(fmt.Sprintf("%v", this.Identifiers), "EntityIdentifiers", "EntityIdentifiers", 1) + `,`,
		`Tail:` + fmt.Sprintf("%v", this.Tail) + `,`,
		`After:` + strings.Replace(fmt.Sprintf("%v", this.After), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringEvents(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Event) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: Event: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Event: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Identifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Identifiers = append(m.Identifiers, &EntityIdentifiers{})
			if err := m.Identifiers[len(m.Identifiers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &types.Any{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CorrelationIDs", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CorrelationIDs = append(m.CorrelationIDs, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Origin", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Origin = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Context", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Context == nil {
				m.Context = make(map[string][]byte)
			}
			var mapkey string
			mapvalue := []byte{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowEvents
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowEvents
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthEvents
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthEvents
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapbyteLen uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowEvents
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapbyteLen |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intMapbyteLen := int(mapbyteLen)
					if intMapbyteLen < 0 {
						return ErrInvalidLengthEvents
					}
					postbytesIndex := iNdEx + intMapbyteLen
					if postbytesIndex < 0 {
						return ErrInvalidLengthEvents
					}
					if postbytesIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = make([]byte, mapbyteLen)
					copy(mapvalue, dAtA[iNdEx:postbytesIndex])
					iNdEx = postbytesIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipEvents(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthEvents
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Context[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEvents
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *StreamEventsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: StreamEventsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamEventsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Identifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Identifiers = append(m.Identifiers, &EntityIdentifiers{})
			if err := m.Identifiers[len(m.Identifiers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tail", wireType)
			}
			m.Tail = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Tail |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field After", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.After == nil {
				m.After = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.After, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEvents
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthEvents
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
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthEvents
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowEvents
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
				next, err := skipEvents(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthEvents
				}
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
	ErrInvalidLengthEvents = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents   = fmt.Errorf("proto: integer overflow")
)
