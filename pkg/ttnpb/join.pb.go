// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/join.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import google_protobuf4 "github.com/gogo/protobuf/types"

import github_com_TheThingsNetwork_ttn_pkg_types "github.com/TheThingsNetwork/ttn/pkg/types"

import bytes "bytes"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Join Request
type JoinRequest struct {
	// Raw uplink bytes (PHYPayload)
	RawPayload []byte `protobuf:"bytes,1,opt,name=raw_payload,json=rawPayload,proto3" json:"raw_payload,omitempty"`
	// LoRaWAN Payload
	Payload Message `protobuf:"bytes,2,opt,name=payload" json:"payload"`
	// End device identifiers
	// - this includes the DevAddr assigned by the network server
	EndDeviceIdentifiers `protobuf:"bytes,3,opt,name=end_device,json=endDevice,embedded=end_device" json:"end_device"`
	// MAC version selected by the network server
	SelectedMacVersion string `protobuf:"bytes,4,opt,name=selected_mac_version,json=selectedMacVersion,proto3" json:"selected_mac_version,omitempty"`
	// NetID of the accepting network server
	NetID github_com_TheThingsNetwork_ttn_pkg_types.NetID `protobuf:"bytes,5,opt,name=net_id,json=netId,proto3,customtype=github.com/TheThingsNetwork/ttn/pkg/types.NetID" json:"net_id"`
	// Downlink Settings
	DownlinkSettings DLSettings `protobuf:"bytes,6,opt,name=downlink_settings,json=downlinkSettings" json:"downlink_settings"`
	// RX Delay in seconds
	RxDelay uint32 `protobuf:"varint,7,opt,name=rx_delay,json=rxDelay,proto3" json:"rx_delay,omitempty"`
	// Optional CFList
	CFList *CFList `protobuf:"bytes,8,opt,name=cf_list,json=cfList" json:"cf_list,omitempty"`
}

func (m *JoinRequest) Reset()                    { *m = JoinRequest{} }
func (m *JoinRequest) String() string            { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()               {}
func (*JoinRequest) Descriptor() ([]byte, []int) { return fileDescriptorJoin, []int{0} }

func (m *JoinRequest) GetRawPayload() []byte {
	if m != nil {
		return m.RawPayload
	}
	return nil
}

func (m *JoinRequest) GetPayload() Message {
	if m != nil {
		return m.Payload
	}
	return Message{}
}

func (m *JoinRequest) GetSelectedMacVersion() string {
	if m != nil {
		return m.SelectedMacVersion
	}
	return ""
}

func (m *JoinRequest) GetDownlinkSettings() DLSettings {
	if m != nil {
		return m.DownlinkSettings
	}
	return DLSettings{}
}

func (m *JoinRequest) GetRxDelay() uint32 {
	if m != nil {
		return m.RxDelay
	}
	return 0
}

func (m *JoinRequest) GetCFList() *CFList {
	if m != nil {
		return m.CFList
	}
	return nil
}

// Answer to the Join Request
type JoinResponse struct {
	// Raw uplink bytes (PHYPayload)
	RawPayload []byte `protobuf:"bytes,1,opt,name=raw_payload,json=rawPayload,proto3" json:"raw_payload,omitempty"`
	// The session keys
	SessionKeys `protobuf:"bytes,2,opt,name=session_keys,json=sessionKeys,embedded=session_keys" json:"session_keys"`
	// Lifetime of the session
	Lifetime *google_protobuf4.Duration `protobuf:"bytes,3,opt,name=lifetime" json:"lifetime,omitempty"`
}

func (m *JoinResponse) Reset()                    { *m = JoinResponse{} }
func (m *JoinResponse) String() string            { return proto.CompactTextString(m) }
func (*JoinResponse) ProtoMessage()               {}
func (*JoinResponse) Descriptor() ([]byte, []int) { return fileDescriptorJoin, []int{1} }

func (m *JoinResponse) GetRawPayload() []byte {
	if m != nil {
		return m.RawPayload
	}
	return nil
}

func (m *JoinResponse) GetLifetime() *google_protobuf4.Duration {
	if m != nil {
		return m.Lifetime
	}
	return nil
}

func init() {
	proto.RegisterType((*JoinRequest)(nil), "ttn.v3.JoinRequest")
	golang_proto.RegisterType((*JoinRequest)(nil), "ttn.v3.JoinRequest")
	proto.RegisterType((*JoinResponse)(nil), "ttn.v3.JoinResponse")
	golang_proto.RegisterType((*JoinResponse)(nil), "ttn.v3.JoinResponse")
}
func (this *JoinRequest) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*JoinRequest)
	if !ok {
		that2, ok := that.(JoinRequest)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *JoinRequest")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *JoinRequest but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *JoinRequest but is not nil && this == nil")
	}
	if !bytes.Equal(this.RawPayload, that1.RawPayload) {
		return fmt.Errorf("RawPayload this(%v) Not Equal that(%v)", this.RawPayload, that1.RawPayload)
	}
	if !this.Payload.Equal(&that1.Payload) {
		return fmt.Errorf("Payload this(%v) Not Equal that(%v)", this.Payload, that1.Payload)
	}
	if !this.EndDeviceIdentifiers.Equal(&that1.EndDeviceIdentifiers) {
		return fmt.Errorf("EndDeviceIdentifiers this(%v) Not Equal that(%v)", this.EndDeviceIdentifiers, that1.EndDeviceIdentifiers)
	}
	if this.SelectedMacVersion != that1.SelectedMacVersion {
		return fmt.Errorf("SelectedMacVersion this(%v) Not Equal that(%v)", this.SelectedMacVersion, that1.SelectedMacVersion)
	}
	if !this.NetID.Equal(that1.NetID) {
		return fmt.Errorf("NetID this(%v) Not Equal that(%v)", this.NetID, that1.NetID)
	}
	if !this.DownlinkSettings.Equal(&that1.DownlinkSettings) {
		return fmt.Errorf("DownlinkSettings this(%v) Not Equal that(%v)", this.DownlinkSettings, that1.DownlinkSettings)
	}
	if this.RxDelay != that1.RxDelay {
		return fmt.Errorf("RxDelay this(%v) Not Equal that(%v)", this.RxDelay, that1.RxDelay)
	}
	if !this.CFList.Equal(that1.CFList) {
		return fmt.Errorf("CFList this(%v) Not Equal that(%v)", this.CFList, that1.CFList)
	}
	return nil
}
func (this *JoinRequest) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*JoinRequest)
	if !ok {
		that2, ok := that.(JoinRequest)
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
	if !bytes.Equal(this.RawPayload, that1.RawPayload) {
		return false
	}
	if !this.Payload.Equal(&that1.Payload) {
		return false
	}
	if !this.EndDeviceIdentifiers.Equal(&that1.EndDeviceIdentifiers) {
		return false
	}
	if this.SelectedMacVersion != that1.SelectedMacVersion {
		return false
	}
	if !this.NetID.Equal(that1.NetID) {
		return false
	}
	if !this.DownlinkSettings.Equal(&that1.DownlinkSettings) {
		return false
	}
	if this.RxDelay != that1.RxDelay {
		return false
	}
	if !this.CFList.Equal(that1.CFList) {
		return false
	}
	return true
}
func (this *JoinResponse) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*JoinResponse)
	if !ok {
		that2, ok := that.(JoinResponse)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *JoinResponse")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *JoinResponse but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *JoinResponse but is not nil && this == nil")
	}
	if !bytes.Equal(this.RawPayload, that1.RawPayload) {
		return fmt.Errorf("RawPayload this(%v) Not Equal that(%v)", this.RawPayload, that1.RawPayload)
	}
	if !this.SessionKeys.Equal(&that1.SessionKeys) {
		return fmt.Errorf("SessionKeys this(%v) Not Equal that(%v)", this.SessionKeys, that1.SessionKeys)
	}
	if !this.Lifetime.Equal(that1.Lifetime) {
		return fmt.Errorf("Lifetime this(%v) Not Equal that(%v)", this.Lifetime, that1.Lifetime)
	}
	return nil
}
func (this *JoinResponse) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*JoinResponse)
	if !ok {
		that2, ok := that.(JoinResponse)
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
	if !bytes.Equal(this.RawPayload, that1.RawPayload) {
		return false
	}
	if !this.SessionKeys.Equal(&that1.SessionKeys) {
		return false
	}
	if !this.Lifetime.Equal(that1.Lifetime) {
		return false
	}
	return true
}
func (m *JoinRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *JoinRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.RawPayload) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintJoin(dAtA, i, uint64(len(m.RawPayload)))
		i += copy(dAtA[i:], m.RawPayload)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintJoin(dAtA, i, uint64(m.Payload.Size()))
	n1, err := m.Payload.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x1a
	i++
	i = encodeVarintJoin(dAtA, i, uint64(m.EndDeviceIdentifiers.Size()))
	n2, err := m.EndDeviceIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	if len(m.SelectedMacVersion) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintJoin(dAtA, i, uint64(len(m.SelectedMacVersion)))
		i += copy(dAtA[i:], m.SelectedMacVersion)
	}
	dAtA[i] = 0x2a
	i++
	i = encodeVarintJoin(dAtA, i, uint64(m.NetID.Size()))
	n3, err := m.NetID.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	dAtA[i] = 0x32
	i++
	i = encodeVarintJoin(dAtA, i, uint64(m.DownlinkSettings.Size()))
	n4, err := m.DownlinkSettings.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	if m.RxDelay != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintJoin(dAtA, i, uint64(m.RxDelay))
	}
	if m.CFList != nil {
		dAtA[i] = 0x42
		i++
		i = encodeVarintJoin(dAtA, i, uint64(m.CFList.Size()))
		n5, err := m.CFList.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	return i, nil
}

func (m *JoinResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *JoinResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.RawPayload) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintJoin(dAtA, i, uint64(len(m.RawPayload)))
		i += copy(dAtA[i:], m.RawPayload)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintJoin(dAtA, i, uint64(m.SessionKeys.Size()))
	n6, err := m.SessionKeys.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n6
	if m.Lifetime != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintJoin(dAtA, i, uint64(m.Lifetime.Size()))
		n7, err := m.Lifetime.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	return i, nil
}

func encodeVarintJoin(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedJoinRequest(r randyJoin, easy bool) *JoinRequest {
	this := &JoinRequest{}
	v1 := r.Intn(100)
	this.RawPayload = make([]byte, v1)
	for i := 0; i < v1; i++ {
		this.RawPayload[i] = byte(r.Intn(256))
	}
	v2 := NewPopulatedMessage(r, easy)
	this.Payload = *v2
	v3 := NewPopulatedEndDeviceIdentifiers(r, easy)
	this.EndDeviceIdentifiers = *v3
	this.SelectedMacVersion = string(randStringJoin(r))
	v4 := github_com_TheThingsNetwork_ttn_pkg_types.NewPopulatedNetID(r)
	this.NetID = *v4
	v5 := NewPopulatedDLSettings(r, easy)
	this.DownlinkSettings = *v5
	this.RxDelay = uint32(r.Uint32())
	if r.Intn(10) != 0 {
		this.CFList = NewPopulatedCFList(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedJoinResponse(r randyJoin, easy bool) *JoinResponse {
	this := &JoinResponse{}
	v6 := r.Intn(100)
	this.RawPayload = make([]byte, v6)
	for i := 0; i < v6; i++ {
		this.RawPayload[i] = byte(r.Intn(256))
	}
	v7 := NewPopulatedSessionKeys(r, easy)
	this.SessionKeys = *v7
	if r.Intn(10) != 0 {
		this.Lifetime = google_protobuf4.NewPopulatedDuration(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyJoin interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneJoin(r randyJoin) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringJoin(r randyJoin) string {
	v8 := r.Intn(100)
	tmps := make([]rune, v8)
	for i := 0; i < v8; i++ {
		tmps[i] = randUTF8RuneJoin(r)
	}
	return string(tmps)
}
func randUnrecognizedJoin(r randyJoin, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldJoin(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldJoin(dAtA []byte, r randyJoin, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateJoin(dAtA, uint64(key))
		v9 := r.Int63()
		if r.Intn(2) == 0 {
			v9 *= -1
		}
		dAtA = encodeVarintPopulateJoin(dAtA, uint64(v9))
	case 1:
		dAtA = encodeVarintPopulateJoin(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateJoin(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateJoin(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateJoin(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateJoin(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *JoinRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.RawPayload)
	if l > 0 {
		n += 1 + l + sovJoin(uint64(l))
	}
	l = m.Payload.Size()
	n += 1 + l + sovJoin(uint64(l))
	l = m.EndDeviceIdentifiers.Size()
	n += 1 + l + sovJoin(uint64(l))
	l = len(m.SelectedMacVersion)
	if l > 0 {
		n += 1 + l + sovJoin(uint64(l))
	}
	l = m.NetID.Size()
	n += 1 + l + sovJoin(uint64(l))
	l = m.DownlinkSettings.Size()
	n += 1 + l + sovJoin(uint64(l))
	if m.RxDelay != 0 {
		n += 1 + sovJoin(uint64(m.RxDelay))
	}
	if m.CFList != nil {
		l = m.CFList.Size()
		n += 1 + l + sovJoin(uint64(l))
	}
	return n
}

func (m *JoinResponse) Size() (n int) {
	var l int
	_ = l
	l = len(m.RawPayload)
	if l > 0 {
		n += 1 + l + sovJoin(uint64(l))
	}
	l = m.SessionKeys.Size()
	n += 1 + l + sovJoin(uint64(l))
	if m.Lifetime != nil {
		l = m.Lifetime.Size()
		n += 1 + l + sovJoin(uint64(l))
	}
	return n
}

func sovJoin(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozJoin(x uint64) (n int) {
	return sovJoin(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *JoinRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowJoin
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
			return fmt.Errorf("proto: JoinRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: JoinRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RawPayload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RawPayload = append(m.RawPayload[:0], dAtA[iNdEx:postIndex]...)
			if m.RawPayload == nil {
				m.RawPayload = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Payload.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndDeviceIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EndDeviceIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SelectedMacVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SelectedMacVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.NetID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DownlinkSettings", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DownlinkSettings.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RxDelay", wireType)
			}
			m.RxDelay = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RxDelay |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CFList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CFList == nil {
				m.CFList = &CFList{}
			}
			if err := m.CFList.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipJoin(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthJoin
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
func (m *JoinResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowJoin
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
			return fmt.Errorf("proto: JoinResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: JoinResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RawPayload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RawPayload = append(m.RawPayload[:0], dAtA[iNdEx:postIndex]...)
			if m.RawPayload == nil {
				m.RawPayload = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SessionKeys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SessionKeys.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Lifetime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJoin
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
				return ErrInvalidLengthJoin
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Lifetime == nil {
				m.Lifetime = &google_protobuf4.Duration{}
			}
			if err := m.Lifetime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipJoin(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthJoin
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
func skipJoin(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowJoin
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
					return 0, ErrIntOverflowJoin
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
					return 0, ErrIntOverflowJoin
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
				return 0, ErrInvalidLengthJoin
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowJoin
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
				next, err := skipJoin(dAtA[start:])
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
	ErrInvalidLengthJoin = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowJoin   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/join.proto", fileDescriptorJoin) }
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/join.proto", fileDescriptorJoin)
}

var fileDescriptorJoin = []byte{
	// 633 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x3f, 0x4c, 0xdb, 0x4e,
	0x14, 0xf6, 0xfd, 0x80, 0x24, 0x5c, 0xf8, 0xb5, 0xd4, 0xed, 0x10, 0x50, 0xf5, 0x12, 0x31, 0xa5,
	0x43, 0xed, 0xb6, 0x94, 0x1d, 0xa5, 0xa1, 0x12, 0x14, 0x50, 0x65, 0x50, 0x2b, 0x75, 0xb1, 0x9c,
	0xf8, 0x62, 0xae, 0x31, 0x77, 0xa9, 0xef, 0x42, 0xc8, 0xc6, 0xc8, 0xd8, 0xb1, 0x5b, 0x3b, 0x32,
	0x32, 0x32, 0x32, 0x32, 0x55, 0x8c, 0xa8, 0x43, 0x84, 0xcf, 0x0b, 0x23, 0x23, 0x63, 0xe5, 0x7f,
	0x0d, 0x1b, 0x99, 0xfc, 0xde, 0xfb, 0xee, 0x7b, 0xcf, 0xdf, 0xfb, 0x1e, 0x36, 0x3c, 0x2a, 0xf7,
	0xfa, 0x2d, 0xa3, 0xcd, 0xf7, 0xcd, 0xdd, 0x3d, 0xb2, 0xbb, 0x47, 0x99, 0x27, 0xb6, 0x89, 0x1c,
	0xf0, 0xa0, 0x6b, 0x4a, 0xc9, 0x4c, 0xa7, 0x47, 0xcd, 0xaf, 0x9c, 0x32, 0xa3, 0x17, 0x70, 0xc9,
	0xf5, 0x82, 0x94, 0xcc, 0x38, 0x58, 0x5e, 0x7c, 0x79, 0x8f, 0xe7, 0x71, 0x8f, 0x9b, 0x09, 0xdc,
	0xea, 0x77, 0x92, 0x2c, 0x49, 0x92, 0x28, 0xa5, 0x2d, 0xbe, 0x9d, 0x64, 0x0c, 0x61, 0xae, 0xed,
	0x92, 0x03, 0xda, 0x26, 0x19, 0x6b, 0x65, 0x12, 0x16, 0x75, 0x09, 0x93, 0xb4, 0x43, 0x49, 0x20,
	0x32, 0xda, 0xeb, 0x49, 0x68, 0x3e, 0x0f, 0x9c, 0x81, 0x93, 0xc9, 0x5a, 0x04, 0x8f, 0x73, 0xcf,
	0x27, 0x63, 0x15, 0x6e, 0x3f, 0x70, 0x24, 0xe5, 0x19, 0xbe, 0xf4, 0x7b, 0x0a, 0x97, 0x37, 0x38,
	0x65, 0x16, 0xf9, 0xd6, 0x27, 0x42, 0xea, 0x55, 0x5c, 0x0e, 0x9c, 0x81, 0xdd, 0x73, 0x86, 0x3e,
	0x77, 0xdc, 0x0a, 0xaa, 0xa1, 0xfa, 0x9c, 0x85, 0x03, 0x67, 0xf0, 0x31, 0xad, 0xe8, 0x26, 0x2e,
	0xe6, 0xe0, 0x7f, 0x35, 0x54, 0x2f, 0xbf, 0x79, 0x6c, 0xa4, 0x9b, 0x33, 0xb6, 0x88, 0x10, 0x8e,
	0x47, 0x1a, 0xd3, 0x17, 0xa3, 0xaa, 0x66, 0xe5, 0xaf, 0xf4, 0x35, 0x8c, 0xc7, 0xfa, 0x2b, 0x53,
	0x09, 0xe7, 0x79, 0xce, 0x59, 0x63, 0x6e, 0x33, 0x01, 0xd6, 0xc7, 0x62, 0x1b, 0xa5, 0xb8, 0xc1,
	0xe5, 0xa8, 0x8a, 0xac, 0x59, 0x92, 0xe3, 0xfa, 0x2b, 0xfc, 0x4c, 0x10, 0x9f, 0xb4, 0x25, 0x71,
	0xed, 0x7d, 0xa7, 0x6d, 0x1f, 0x90, 0x40, 0x50, 0xce, 0x2a, 0xd3, 0x35, 0x54, 0x9f, 0xb5, 0xf4,
	0x1c, 0xdb, 0x72, 0xda, 0x9f, 0x52, 0x44, 0xff, 0x8c, 0x0b, 0x8c, 0x48, 0x9b, 0xba, 0x95, 0x99,
	0x58, 0x45, 0x63, 0x35, 0x6e, 0xfb, 0x67, 0x54, 0x35, 0x1f, 0xda, 0x62, 0xaf, 0xeb, 0x99, 0x72,
	0xd8, 0x23, 0xc2, 0xd8, 0x26, 0x72, 0xbd, 0xa9, 0x46, 0xd5, 0x99, 0x24, 0xb0, 0x66, 0x18, 0x91,
	0xeb, 0xb1, 0xa2, 0x27, 0x2e, 0x1f, 0x30, 0x9f, 0xb2, 0xae, 0x2d, 0x88, 0x94, 0x31, 0xbf, 0x52,
	0x48, 0x84, 0xe9, 0xb9, 0xb0, 0xe6, 0xe6, 0x4e, 0x86, 0x64, 0xfb, 0x98, 0xcf, 0x29, 0x79, 0x5d,
	0x5f, 0xc0, 0xa5, 0xe0, 0xd0, 0x76, 0x89, 0xef, 0x0c, 0x2b, 0xc5, 0x1a, 0xaa, 0xff, 0x6f, 0x15,
	0x83, 0xc3, 0x66, 0x9c, 0xea, 0xcb, 0xb8, 0xd8, 0xee, 0xd8, 0x3e, 0x15, 0xb2, 0x52, 0x4a, 0xfa,
	0x3e, 0xca, 0xfb, 0xbe, 0x7b, 0xbf, 0x49, 0x85, 0x6c, 0x60, 0x35, 0xaa, 0x16, 0xd2, 0xd8, 0x2a,
	0xb4, 0x3b, 0xf1, 0x77, 0x63, 0xba, 0x34, 0x3b, 0x8f, 0x97, 0x4e, 0x10, 0x9e, 0x4b, 0x0d, 0x15,
	0x3d, 0xce, 0x04, 0x79, 0xd8, 0xd1, 0x55, 0x3c, 0x27, 0x88, 0x88, 0x57, 0x66, 0x77, 0xc9, 0x50,
	0x64, 0xb6, 0x3e, 0xcd, 0x27, 0xee, 0xa4, 0xd8, 0x07, 0x32, 0xbc, 0xef, 0x4c, 0x59, 0x8c, 0xcb,
	0xfa, 0x0a, 0x2e, 0xf9, 0xb4, 0x43, 0x24, 0xdd, 0xcf, 0x0d, 0x5e, 0x30, 0xd2, 0xbb, 0x33, 0xf2,
	0xbb, 0x33, 0x9a, 0xd9, 0xdd, 0x59, 0xff, 0x9e, 0x36, 0x7e, 0xa2, 0x8b, 0x10, 0xd0, 0x65, 0x08,
	0xe8, 0x2a, 0x04, 0x74, 0x1d, 0x02, 0xba, 0x09, 0x41, 0xbb, 0x0d, 0x41, 0xbb, 0x0b, 0x01, 0x1d,
	0x29, 0xd0, 0x8e, 0x15, 0x68, 0x27, 0x0a, 0xd0, 0xa9, 0x02, 0xed, 0x4c, 0x01, 0x3a, 0x57, 0x80,
	0x2e, 0x14, 0xa0, 0x4b, 0x05, 0xe8, 0x4a, 0x81, 0x76, 0xad, 0x00, 0xdd, 0x28, 0xd0, 0x6e, 0x15,
	0xa0, 0x3b, 0x05, 0xda, 0x51, 0x04, 0xda, 0x71, 0x04, 0xe8, 0x7b, 0x04, 0xda, 0x8f, 0x08, 0xd0,
	0xaf, 0x08, 0xb4, 0x93, 0x08, 0xb4, 0xd3, 0x08, 0xd0, 0x59, 0x04, 0xe8, 0x3c, 0x02, 0xf4, 0xe5,
	0xc5, 0x44, 0xfe, 0x4b, 0xd6, 0x6b, 0xb5, 0x0a, 0xc9, 0xef, 0x2f, 0xff, 0x0d, 0x00, 0x00, 0xff,
	0xff, 0xda, 0x6b, 0x7a, 0xd1, 0x4d, 0x04, 0x00, 0x00,
}
