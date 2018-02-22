// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/end_device_services.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf2 "github.com/gogo/protobuf/types"
import google_protobuf4 "github.com/gogo/protobuf/types"

import (
	context "context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type SetDeviceRequest struct {
	Device    EndDevice                   `protobuf:"bytes,1,opt,name=device" json:"device"`
	FieldMask *google_protobuf4.FieldMask `protobuf:"bytes,2,opt,name=field_mask,json=fieldMask" json:"field_mask,omitempty"`
}

func (m *SetDeviceRequest) Reset()         { *m = SetDeviceRequest{} }
func (m *SetDeviceRequest) String() string { return proto.CompactTextString(m) }
func (*SetDeviceRequest) ProtoMessage()    {}
func (*SetDeviceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptorEndDeviceServices, []int{0}
}

func (m *SetDeviceRequest) GetDevice() EndDevice {
	if m != nil {
		return m.Device
	}
	return EndDevice{}
}

func (m *SetDeviceRequest) GetFieldMask() *google_protobuf4.FieldMask {
	if m != nil {
		return m.FieldMask
	}
	return nil
}

func init() {
	proto.RegisterType((*SetDeviceRequest)(nil), "ttn.v3.SetDeviceRequest")
	golang_proto.RegisterType((*SetDeviceRequest)(nil), "ttn.v3.SetDeviceRequest")
}
func (this *SetDeviceRequest) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*SetDeviceRequest)
	if !ok {
		that2, ok := that.(SetDeviceRequest)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *SetDeviceRequest")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *SetDeviceRequest but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *SetDeviceRequest but is not nil && this == nil")
	}
	if !this.Device.Equal(&that1.Device) {
		return fmt.Errorf("Device this(%v) Not Equal that(%v)", this.Device, that1.Device)
	}
	if !this.FieldMask.Equal(that1.FieldMask) {
		return fmt.Errorf("FieldMask this(%v) Not Equal that(%v)", this.FieldMask, that1.FieldMask)
	}
	return nil
}
func (this *SetDeviceRequest) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*SetDeviceRequest)
	if !ok {
		that2, ok := that.(SetDeviceRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if !this.Device.Equal(&that1.Device) {
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

// Client API for NsDeviceRegistry service

type NsDeviceRegistryClient interface {
	// ListDevices returns the devices that match the given identifiers
	ListDevices(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevices, error)
	// GetDevice returns the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevice, error)
	// SetDevice creates or updates the device
	SetDevice(ctx context.Context, in *SetDeviceRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	// DeleteDevice deletes the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}

type nsDeviceRegistryClient struct {
	cc *grpc.ClientConn
}

func NewNsDeviceRegistryClient(cc *grpc.ClientConn) NsDeviceRegistryClient {
	return &nsDeviceRegistryClient{cc}
}

func (c *nsDeviceRegistryClient) ListDevices(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevices, error) {
	out := new(EndDevices)
	err := grpc.Invoke(ctx, "/ttn.v3.NsDeviceRegistry/ListDevices", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nsDeviceRegistryClient) GetDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevice, error) {
	out := new(EndDevice)
	err := grpc.Invoke(ctx, "/ttn.v3.NsDeviceRegistry/GetDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nsDeviceRegistryClient) SetDevice(ctx context.Context, in *SetDeviceRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.NsDeviceRegistry/SetDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nsDeviceRegistryClient) DeleteDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.NsDeviceRegistry/DeleteDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for NsDeviceRegistry service

type NsDeviceRegistryServer interface {
	// ListDevices returns the devices that match the given identifiers
	ListDevices(context.Context, *EndDeviceIdentifiers) (*EndDevices, error)
	// GetDevice returns the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetDevice(context.Context, *EndDeviceIdentifiers) (*EndDevice, error)
	// SetDevice creates or updates the device
	SetDevice(context.Context, *SetDeviceRequest) (*google_protobuf2.Empty, error)
	// DeleteDevice deletes the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteDevice(context.Context, *EndDeviceIdentifiers) (*google_protobuf2.Empty, error)
}

func RegisterNsDeviceRegistryServer(s *grpc.Server, srv NsDeviceRegistryServer) {
	s.RegisterService(&_NsDeviceRegistry_serviceDesc, srv)
}

func _NsDeviceRegistry_ListDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NsDeviceRegistryServer).ListDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.NsDeviceRegistry/ListDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NsDeviceRegistryServer).ListDevices(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _NsDeviceRegistry_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NsDeviceRegistryServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.NsDeviceRegistry/GetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NsDeviceRegistryServer).GetDevice(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _NsDeviceRegistry_SetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NsDeviceRegistryServer).SetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.NsDeviceRegistry/SetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NsDeviceRegistryServer).SetDevice(ctx, req.(*SetDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NsDeviceRegistry_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NsDeviceRegistryServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.NsDeviceRegistry/DeleteDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NsDeviceRegistryServer).DeleteDevice(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _NsDeviceRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.v3.NsDeviceRegistry",
	HandlerType: (*NsDeviceRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDevices",
			Handler:    _NsDeviceRegistry_ListDevices_Handler,
		},
		{
			MethodName: "GetDevice",
			Handler:    _NsDeviceRegistry_GetDevice_Handler,
		},
		{
			MethodName: "SetDevice",
			Handler:    _NsDeviceRegistry_SetDevice_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _NsDeviceRegistry_DeleteDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/TheThingsNetwork/ttn/api/end_device_services.proto",
}

// Client API for AsDeviceRegistry service

type AsDeviceRegistryClient interface {
	// ListDevices returns the devices that match the given identifiers
	ListDevices(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevices, error)
	// GetDevice returns the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevice, error)
	// SetDevice creates or updates the device
	SetDevice(ctx context.Context, in *SetDeviceRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	// DeleteDevice deletes the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}

type asDeviceRegistryClient struct {
	cc *grpc.ClientConn
}

func NewAsDeviceRegistryClient(cc *grpc.ClientConn) AsDeviceRegistryClient {
	return &asDeviceRegistryClient{cc}
}

func (c *asDeviceRegistryClient) ListDevices(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevices, error) {
	out := new(EndDevices)
	err := grpc.Invoke(ctx, "/ttn.v3.AsDeviceRegistry/ListDevices", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asDeviceRegistryClient) GetDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevice, error) {
	out := new(EndDevice)
	err := grpc.Invoke(ctx, "/ttn.v3.AsDeviceRegistry/GetDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asDeviceRegistryClient) SetDevice(ctx context.Context, in *SetDeviceRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.AsDeviceRegistry/SetDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *asDeviceRegistryClient) DeleteDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.AsDeviceRegistry/DeleteDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AsDeviceRegistry service

type AsDeviceRegistryServer interface {
	// ListDevices returns the devices that match the given identifiers
	ListDevices(context.Context, *EndDeviceIdentifiers) (*EndDevices, error)
	// GetDevice returns the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetDevice(context.Context, *EndDeviceIdentifiers) (*EndDevice, error)
	// SetDevice creates or updates the device
	SetDevice(context.Context, *SetDeviceRequest) (*google_protobuf2.Empty, error)
	// DeleteDevice deletes the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteDevice(context.Context, *EndDeviceIdentifiers) (*google_protobuf2.Empty, error)
}

func RegisterAsDeviceRegistryServer(s *grpc.Server, srv AsDeviceRegistryServer) {
	s.RegisterService(&_AsDeviceRegistry_serviceDesc, srv)
}

func _AsDeviceRegistry_ListDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsDeviceRegistryServer).ListDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.AsDeviceRegistry/ListDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsDeviceRegistryServer).ListDevices(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _AsDeviceRegistry_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsDeviceRegistryServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.AsDeviceRegistry/GetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsDeviceRegistryServer).GetDevice(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _AsDeviceRegistry_SetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsDeviceRegistryServer).SetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.AsDeviceRegistry/SetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsDeviceRegistryServer).SetDevice(ctx, req.(*SetDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AsDeviceRegistry_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AsDeviceRegistryServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.AsDeviceRegistry/DeleteDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AsDeviceRegistryServer).DeleteDevice(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _AsDeviceRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.v3.AsDeviceRegistry",
	HandlerType: (*AsDeviceRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDevices",
			Handler:    _AsDeviceRegistry_ListDevices_Handler,
		},
		{
			MethodName: "GetDevice",
			Handler:    _AsDeviceRegistry_GetDevice_Handler,
		},
		{
			MethodName: "SetDevice",
			Handler:    _AsDeviceRegistry_SetDevice_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _AsDeviceRegistry_DeleteDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/TheThingsNetwork/ttn/api/end_device_services.proto",
}

// Client API for JsDeviceRegistry service

type JsDeviceRegistryClient interface {
	// ListDevices returns the devices that match the given identifiers
	ListDevices(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevices, error)
	// GetDevice returns the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevice, error)
	// SetDevice creates or updates the device
	SetDevice(ctx context.Context, in *SetDeviceRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	// DeleteDevice deletes the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}

type jsDeviceRegistryClient struct {
	cc *grpc.ClientConn
}

func NewJsDeviceRegistryClient(cc *grpc.ClientConn) JsDeviceRegistryClient {
	return &jsDeviceRegistryClient{cc}
}

func (c *jsDeviceRegistryClient) ListDevices(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevices, error) {
	out := new(EndDevices)
	err := grpc.Invoke(ctx, "/ttn.v3.JsDeviceRegistry/ListDevices", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jsDeviceRegistryClient) GetDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*EndDevice, error) {
	out := new(EndDevice)
	err := grpc.Invoke(ctx, "/ttn.v3.JsDeviceRegistry/GetDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jsDeviceRegistryClient) SetDevice(ctx context.Context, in *SetDeviceRequest, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.JsDeviceRegistry/SetDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jsDeviceRegistryClient) DeleteDevice(ctx context.Context, in *EndDeviceIdentifiers, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/ttn.v3.JsDeviceRegistry/DeleteDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for JsDeviceRegistry service

type JsDeviceRegistryServer interface {
	// ListDevices returns the devices that match the given identifiers
	ListDevices(context.Context, *EndDeviceIdentifiers) (*EndDevices, error)
	// GetDevice returns the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	GetDevice(context.Context, *EndDeviceIdentifiers) (*EndDevice, error)
	// SetDevice creates or updates the device
	SetDevice(context.Context, *SetDeviceRequest) (*google_protobuf2.Empty, error)
	// DeleteDevice deletes the device that matches the given identifiers.
	// If there are multiple matches, an error will be returned.
	DeleteDevice(context.Context, *EndDeviceIdentifiers) (*google_protobuf2.Empty, error)
}

func RegisterJsDeviceRegistryServer(s *grpc.Server, srv JsDeviceRegistryServer) {
	s.RegisterService(&_JsDeviceRegistry_serviceDesc, srv)
}

func _JsDeviceRegistry_ListDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JsDeviceRegistryServer).ListDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.JsDeviceRegistry/ListDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JsDeviceRegistryServer).ListDevices(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _JsDeviceRegistry_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JsDeviceRegistryServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.JsDeviceRegistry/GetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JsDeviceRegistryServer).GetDevice(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _JsDeviceRegistry_SetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JsDeviceRegistryServer).SetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.JsDeviceRegistry/SetDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JsDeviceRegistryServer).SetDevice(ctx, req.(*SetDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JsDeviceRegistry_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndDeviceIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JsDeviceRegistryServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.v3.JsDeviceRegistry/DeleteDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JsDeviceRegistryServer).DeleteDevice(ctx, req.(*EndDeviceIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _JsDeviceRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.v3.JsDeviceRegistry",
	HandlerType: (*JsDeviceRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDevices",
			Handler:    _JsDeviceRegistry_ListDevices_Handler,
		},
		{
			MethodName: "GetDevice",
			Handler:    _JsDeviceRegistry_GetDevice_Handler,
		},
		{
			MethodName: "SetDevice",
			Handler:    _JsDeviceRegistry_SetDevice_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _JsDeviceRegistry_DeleteDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/TheThingsNetwork/ttn/api/end_device_services.proto",
}

func (m *SetDeviceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetDeviceRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintEndDeviceServices(dAtA, i, uint64(m.Device.Size()))
	n1, err := m.Device.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if m.FieldMask != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintEndDeviceServices(dAtA, i, uint64(m.FieldMask.Size()))
		n2, err := m.FieldMask.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func encodeVarintEndDeviceServices(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedSetDeviceRequest(r randyEndDeviceServices, easy bool) *SetDeviceRequest {
	this := &SetDeviceRequest{}
	v1 := NewPopulatedEndDevice(r, easy)
	this.Device = *v1
	if r.Intn(10) != 0 {
		this.FieldMask = google_protobuf4.NewPopulatedFieldMask(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyEndDeviceServices interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneEndDeviceServices(r randyEndDeviceServices) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringEndDeviceServices(r randyEndDeviceServices) string {
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
		tmps[i] = randUTF8RuneEndDeviceServices(r)
	}
	return string(tmps)
}
func randUnrecognizedEndDeviceServices(r randyEndDeviceServices, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldEndDeviceServices(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldEndDeviceServices(dAtA []byte, r randyEndDeviceServices, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateEndDeviceServices(dAtA, uint64(key))
		v3 := r.Int63()
		if r.Intn(2) == 0 {
			v3 *= -1
		}
		dAtA = encodeVarintPopulateEndDeviceServices(dAtA, uint64(v3))
	case 1:
		dAtA = encodeVarintPopulateEndDeviceServices(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateEndDeviceServices(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateEndDeviceServices(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateEndDeviceServices(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateEndDeviceServices(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *SetDeviceRequest) Size() (n int) {
	var l int
	_ = l
	l = m.Device.Size()
	n += 1 + l + sovEndDeviceServices(uint64(l))
	if m.FieldMask != nil {
		l = m.FieldMask.Size()
		n += 1 + l + sovEndDeviceServices(uint64(l))
	}
	return n
}

func sovEndDeviceServices(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozEndDeviceServices(x uint64) (n int) {
	return sovEndDeviceServices((x << 1) ^ uint64((int64(x) >> 63)))
}
func (m *SetDeviceRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEndDeviceServices
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
			return fmt.Errorf("proto: SetDeviceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetDeviceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Device", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEndDeviceServices
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
				return ErrInvalidLengthEndDeviceServices
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Device.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
					return ErrIntOverflowEndDeviceServices
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
				return ErrInvalidLengthEndDeviceServices
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
			skippy, err := skipEndDeviceServices(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEndDeviceServices
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
func skipEndDeviceServices(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEndDeviceServices
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
					return 0, ErrIntOverflowEndDeviceServices
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
					return 0, ErrIntOverflowEndDeviceServices
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
				return 0, ErrInvalidLengthEndDeviceServices
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowEndDeviceServices
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
				next, err := skipEndDeviceServices(dAtA[start:])
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
	ErrInvalidLengthEndDeviceServices = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEndDeviceServices   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/end_device_services.proto", fileDescriptorEndDeviceServices)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/end_device_services.proto", fileDescriptorEndDeviceServices)
}

var fileDescriptorEndDeviceServices = []byte{
	// 617 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x96, 0x3b, 0x4c, 0xdc, 0x4c,
	0x10, 0xc7, 0x77, 0x10, 0x42, 0xc2, 0x7c, 0x05, 0x9f, 0x8b, 0x08, 0x39, 0x68, 0x82, 0x5c, 0x05,
	0xa2, 0xd8, 0x0a, 0x0f, 0x45, 0x20, 0x51, 0x04, 0x41, 0xa2, 0xbc, 0x28, 0x08, 0x15, 0xcd, 0xc9,
	0x87, 0xf7, 0xcc, 0xde, 0x1d, 0xb6, 0xe3, 0x5d, 0x88, 0x10, 0x22, 0x42, 0x49, 0x83, 0x94, 0x26,
	0x4a, 0x9a, 0x74, 0x49, 0x49, 0x49, 0x49, 0x49, 0x49, 0x89, 0x92, 0x86, 0x0a, 0xe1, 0x75, 0x0a,
	0x4a, 0x4a, 0xca, 0xe8, 0xd6, 0x36, 0x20, 0x4e, 0xd1, 0x5d, 0x4e, 0xbc, 0x2a, 0xef, 0xde, 0x3c,
	0xfe, 0x33, 0xff, 0xfd, 0x15, 0xa7, 0x8d, 0x7b, 0x4c, 0x2c, 0x2c, 0x15, 0xad, 0xf9, 0x60, 0xd1,
	0x9e, 0x5d, 0xa0, 0xb3, 0x0b, 0xcc, 0xf7, 0xf8, 0x34, 0x15, 0xef, 0x82, 0xa8, 0x62, 0x0b, 0xe1,
	0xdb, 0x4e, 0xc8, 0x6c, 0xea, 0xbb, 0x05, 0x97, 0x2e, 0xb3, 0x79, 0x5a, 0xe0, 0x34, 0xaa, 0x7d,
	0xb9, 0x15, 0x46, 0x81, 0x08, 0xf4, 0x0e, 0x21, 0x7c, 0x6b, 0x79, 0xc8, 0x18, 0xfe, 0xb7, 0x36,
	0x69, 0xb5, 0x31, 0xd2, 0x4c, 0x15, 0x73, 0xa9, 0x2f, 0x58, 0x89, 0xd1, 0x28, 0x13, 0x35, 0x1e,
	0x9e, 0x2b, 0xf3, 0x02, 0x2f, 0xb0, 0xd5, 0xcf, 0xc5, 0xa5, 0x92, 0xba, 0xa9, 0x8b, 0x3a, 0x65,
	0xe9, 0xbd, 0x5e, 0x10, 0x78, 0x55, 0xaa, 0x9a, 0x39, 0xbe, 0x1f, 0x08, 0x47, 0xb0, 0xc0, 0xcf,
	0x9b, 0xdd, 0xcd, 0xa2, 0xa7, 0x3d, 0xe8, 0x62, 0x28, 0x56, 0xb2, 0x60, 0xdf, 0xc5, 0x60, 0x89,
	0xd1, 0xaa, 0x5b, 0x58, 0x74, 0x78, 0x25, 0xcd, 0x30, 0xdf, 0x6b, 0xdd, 0x6f, 0xa8, 0x98, 0x54,
	0x5b, 0xcd, 0xd0, 0xb7, 0x4b, 0x94, 0x0b, 0xdd, 0xd6, 0x3a, 0xd2, 0x35, 0x7b, 0xa0, 0x0f, 0xee,
	0x77, 0x0d, 0xfe, 0x6f, 0xa5, 0x2e, 0x59, 0x53, 0xbe, 0x9b, 0x66, 0x4e, 0xb4, 0xef, 0x1e, 0xdc,
	0x23, 0x33, 0x59, 0x9a, 0x3e, 0xaa, 0x69, 0x67, 0x8d, 0x7b, 0xda, 0x54, 0x91, 0x61, 0xa5, 0xda,
	0x56, 0xae, 0x6d, 0x3d, 0xad, 0xa5, 0xbc, 0x76, 0x78, 0x65, 0xa6, 0xb3, 0x94, 0x1f, 0x07, 0x0f,
	0xda, 0xb5, 0xee, 0x69, 0x9e, 0xeb, 0x7b, 0x8c, 0x8b, 0x68, 0x45, 0x8f, 0xb4, 0xae, 0x57, 0x8c,
	0x67, 0x53, 0x71, 0xbd, 0xb7, 0x4e, 0xff, 0xf9, 0x99, 0xa7, 0x86, 0x5e, 0x17, 0xe5, 0xe6, 0xa3,
	0x0f, 0xbf, 0x7e, 0x7f, 0x6d, 0x7b, 0xa0, 0xf7, 0xdb, 0x3e, 0xb7, 0x9d, 0x30, 0xac, 0xb2, 0xf9,
	0xd4, 0x34, 0x7b, 0xf5, 0xdc, 0xad, 0xc0, 0xdc, 0x35, 0xdb, 0xcd, 0x44, 0x56, 0xb5, 0xce, 0x67,
	0xb9, 0x11, 0x0d, 0x14, 0xeb, 0xfd, 0x30, 0xc7, 0x95, 0xe0, 0x63, 0x7d, 0xa4, 0x69, 0x41, 0x7b,
	0x35, 0x23, 0x92, 0xb9, 0x6b, 0xfa, 0x4f, 0xd0, 0x3a, 0x4f, 0x9f, 0x41, 0xef, 0xc9, 0xfb, 0x5f,
	0x7c, 0x19, 0xe3, 0x4e, 0x9d, 0xa9, 0x53, 0xb5, 0xd7, 0x36, 0xbf, 0x80, 0xd2, 0xff, 0x04, 0xe6,
	0xcb, 0xfa, 0x01, 0x32, 0x64, 0x99, 0xcb, 0xad, 0x06, 0xb3, 0xa8, 0x9c, 0xb3, 0xb1, 0xc6, 0x60,
	0x60, 0x6e, 0xd4, 0x1c, 0x6e, 0xa5, 0xe3, 0x18, 0x0c, 0xe8, 0x1f, 0x41, 0xfb, 0x6f, 0x92, 0x56,
	0xa9, 0xa0, 0x4d, 0xb9, 0xfa, 0xb7, 0xdd, 0x32, 0x6b, 0x07, 0x5a, 0xb3, 0x56, 0x01, 0xf6, 0xe4,
	0x0a, 0x01, 0x73, 0xae, 0x1b, 0x30, 0xe7, 0xa6, 0x01, 0x73, 0x2e, 0x1d, 0x30, 0xe7, 0x96, 0x00,
	0xe6, 0xb4, 0x0a, 0xd8, 0x8b, 0x2b, 0x04, 0xac, 0x7c, 0xdd, 0x80, 0x95, 0x6f, 0x1a, 0xb0, 0xf2,
	0xa5, 0x03, 0x56, 0xbe, 0x25, 0x80, 0xb5, 0x62, 0xed, 0xc4, 0x77, 0xd8, 0x8d, 0x11, 0xf6, 0x62,
	0x84, 0xfd, 0x18, 0xe1, 0x30, 0x46, 0x38, 0x8a, 0x91, 0x1c, 0xc7, 0x48, 0x4e, 0x62, 0x84, 0x75,
	0x89, 0x64, 0x43, 0x22, 0xd9, 0x94, 0x08, 0x5b, 0x12, 0xc9, 0xb6, 0x44, 0xd8, 0x91, 0x08, 0xbb,
	0x12, 0x61, 0x4f, 0x22, 0xec, 0x4b, 0x24, 0x87, 0x12, 0xe1, 0x48, 0x22, 0x39, 0x96, 0x08, 0x27,
	0x12, 0xc9, 0x7a, 0x82, 0x64, 0x23, 0x41, 0xf8, 0x9c, 0x20, 0xf9, 0x96, 0x20, 0xfc, 0x48, 0x90,
	0x6c, 0x26, 0x48, 0xb6, 0x12, 0x84, 0xed, 0x04, 0x61, 0x27, 0x41, 0x98, 0xeb, 0x6f, 0xf4, 0xe7,
	0x26, 0xac, 0x78, 0xb5, 0x6f, 0x58, 0x2c, 0x76, 0xa8, 0x85, 0x87, 0xfe, 0x04, 0x00, 0x00, 0xff,
	0xff, 0x58, 0x89, 0x44, 0xfb, 0x8d, 0x09, 0x00, 0x00,
}
