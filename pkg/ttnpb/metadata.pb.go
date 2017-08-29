// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/metadata.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import google_protobuf2 "github.com/gogo/protobuf/types"
import _ "github.com/gogo/protobuf/types"

import time "time"

import strconv "strconv"

import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// LocationSource indicates the source of a Location
type LocationSource int32

const (
	// The source of the location is not known or not set
	SOURCE_UNKNOWN LocationSource = 0
	// The location is determined by GPS
	SOURCE_GPS LocationSource = 1
	// The location is fixed by configuration
	SOURCE_CONFIG LocationSource = 2
	// The location is set in and updated from a registry
	SOURCE_REGISTRY LocationSource = 3
	// The location is estimated with IP Geolocation
	SOURCE_IP_GEOLOCATION LocationSource = 4
)

var LocationSource_name = map[int32]string{
	0: "SOURCE_UNKNOWN",
	1: "SOURCE_GPS",
	2: "SOURCE_CONFIG",
	3: "SOURCE_REGISTRY",
	4: "SOURCE_IP_GEOLOCATION",
}
var LocationSource_value = map[string]int32{
	"SOURCE_UNKNOWN":        0,
	"SOURCE_GPS":            1,
	"SOURCE_CONFIG":         2,
	"SOURCE_REGISTRY":       3,
	"SOURCE_IP_GEOLOCATION": 4,
}

func (LocationSource) EnumDescriptor() ([]byte, []int) { return fileDescriptorMetadata, []int{0} }

// RxMetadata contains metadata for a received message. Each antenna that receives
// a message corresponds to one RxMetadata.
type RxMetadata struct {
	// ID of the gateway that received the message
	GatewayIdentifiers `protobuf:"bytes,1,opt,name=gateway_id,json=gatewayId,embedded=gateway_id" json:"gateway_id"`
	// Index of the antenna that received the message
	AntennaIndex uint32 `protobuf:"varint,2,opt,name=antenna_index,json=antennaIndex,proto3" json:"antenna_index,omitempty"`
	// Index of the channel that received the message
	ChannelIndex uint32 `protobuf:"varint,3,opt,name=channel_index,json=channelIndex,proto3" json:"channel_index,omitempty"`
	// Gateway's internal timestamp when the Rx finished (nanoseconds)
	// NOTE: most gateways use microsecond timestamps, so conversion may be needed
	Timestamp uint64 `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Real time
	Time *time.Time `protobuf:"bytes,5,opt,name=time,stdtime" json:"time,omitempty"`
	// Received signal strength in dBm
	RSSI float32 `protobuf:"fixed32,6,opt,name=rssi,proto3" json:"rssi,omitempty"`
	// Signal-to-noise-ratio in dB
	SNR float32 `protobuf:"fixed32,7,opt,name=snr,proto3" json:"snr,omitempty"`
	// Location of the antenna
	Location *Location `protobuf:"bytes,8,opt,name=location" json:"location,omitempty"`
	// Advanced metadata fields
	// - can be used for advanced information or experimental features that are not yet formally defined in the API
	// - field names are written in snake_case
	Advanced *google_protobuf2.Struct `protobuf:"bytes,99,opt,name=advanced" json:"advanced,omitempty"`
}

func (m *RxMetadata) Reset()                    { *m = RxMetadata{} }
func (*RxMetadata) ProtoMessage()               {}
func (*RxMetadata) Descriptor() ([]byte, []int) { return fileDescriptorMetadata, []int{0} }

func (m *RxMetadata) GetAntennaIndex() uint32 {
	if m != nil {
		return m.AntennaIndex
	}
	return 0
}

func (m *RxMetadata) GetChannelIndex() uint32 {
	if m != nil {
		return m.ChannelIndex
	}
	return 0
}

func (m *RxMetadata) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *RxMetadata) GetTime() *time.Time {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *RxMetadata) GetRSSI() float32 {
	if m != nil {
		return m.RSSI
	}
	return 0
}

func (m *RxMetadata) GetSNR() float32 {
	if m != nil {
		return m.SNR
	}
	return 0
}

func (m *RxMetadata) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *RxMetadata) GetAdvanced() *google_protobuf2.Struct {
	if m != nil {
		return m.Advanced
	}
	return nil
}

// TxMetadata contains metadata for a to-be-transmitted message.
// It contains gateway-specific fields that are not in the TxSettings message
type TxMetadata struct {
	// ID of the gateway that received the message
	GatewayIdentifiers `protobuf:"bytes,1,opt,name=gateway_id,json=gatewayId,embedded=gateway_id" json:"gateway_id"`
	// Gateway's internal timestamp when the Tx should start (nanoseconds)
	// NOTE: most gateways use microsecond timestamps, so conversion may be needed
	Timestamp uint64 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Real time
	Time *time.Time `protobuf:"bytes,3,opt,name=time,stdtime" json:"time,omitempty"`
	// Advanced metadata fields
	// - can be used for advanced information or experimental features that are not yet formally defined in the API
	// - field names are written in snake_case
	Advanced *google_protobuf2.Struct `protobuf:"bytes,99,opt,name=advanced" json:"advanced,omitempty"`
}

func (m *TxMetadata) Reset()                    { *m = TxMetadata{} }
func (*TxMetadata) ProtoMessage()               {}
func (*TxMetadata) Descriptor() ([]byte, []int) { return fileDescriptorMetadata, []int{1} }

func (m *TxMetadata) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *TxMetadata) GetTime() *time.Time {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *TxMetadata) GetAdvanced() *google_protobuf2.Struct {
	if m != nil {
		return m.Advanced
	}
	return nil
}

// Location of an object
type Location struct {
	// The North–South position (degrees; -90 to +90), where 0 is the equator, North pole is positive, South pole is negative
	Latitude float32 `protobuf:"fixed32,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	// The East-West position (degrees; -180 to +180), where 0 is the Prime Meridian (Greenwich), East is positive , West is negative
	Longitude float32 `protobuf:"fixed32,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	// The altitude (meters), where 0 is the mean sea level
	Altitude int32 `protobuf:"varint,3,opt,name=altitude,proto3" json:"altitude,omitempty"`
	// The accuracy of the location (meters)
	Accuracy int32 `protobuf:"varint,4,opt,name=accuracy,proto3" json:"accuracy,omitempty"`
	// Source of the location information
	Source LocationSource `protobuf:"varint,5,opt,name=source,proto3,enum=ttn.v3.LocationSource" json:"source,omitempty"`
}

func (m *Location) Reset()                    { *m = Location{} }
func (*Location) ProtoMessage()               {}
func (*Location) Descriptor() ([]byte, []int) { return fileDescriptorMetadata, []int{2} }

func (m *Location) GetLatitude() float32 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *Location) GetLongitude() float32 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

func (m *Location) GetAltitude() int32 {
	if m != nil {
		return m.Altitude
	}
	return 0
}

func (m *Location) GetAccuracy() int32 {
	if m != nil {
		return m.Accuracy
	}
	return 0
}

func (m *Location) GetSource() LocationSource {
	if m != nil {
		return m.Source
	}
	return SOURCE_UNKNOWN
}

func init() {
	proto.RegisterType((*RxMetadata)(nil), "ttn.v3.RxMetadata")
	proto.RegisterType((*TxMetadata)(nil), "ttn.v3.TxMetadata")
	proto.RegisterType((*Location)(nil), "ttn.v3.Location")
	proto.RegisterEnum("ttn.v3.LocationSource", LocationSource_name, LocationSource_value)
}
func (x LocationSource) String() string {
	s, ok := LocationSource_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (m *RxMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RxMetadata) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintMetadata(dAtA, i, uint64(m.GatewayIdentifiers.Size()))
	n1, err := m.GatewayIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if m.AntennaIndex != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.AntennaIndex))
	}
	if m.ChannelIndex != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.ChannelIndex))
	}
	if m.Timestamp != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Timestamp))
	}
	if m.Time != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.Time)))
		n2, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.Time, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.RSSI != 0 {
		dAtA[i] = 0x35
		i++
		i = encodeFixed32Metadata(dAtA, i, uint32(math.Float32bits(float32(m.RSSI))))
	}
	if m.SNR != 0 {
		dAtA[i] = 0x3d
		i++
		i = encodeFixed32Metadata(dAtA, i, uint32(math.Float32bits(float32(m.SNR))))
	}
	if m.Location != nil {
		dAtA[i] = 0x42
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Location.Size()))
		n3, err := m.Location.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.Advanced != nil {
		dAtA[i] = 0x9a
		i++
		dAtA[i] = 0x6
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Advanced.Size()))
		n4, err := m.Advanced.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	return i, nil
}

func (m *TxMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxMetadata) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintMetadata(dAtA, i, uint64(m.GatewayIdentifiers.Size()))
	n5, err := m.GatewayIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	if m.Timestamp != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Timestamp))
	}
	if m.Time != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.Time)))
		n6, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.Time, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	if m.Advanced != nil {
		dAtA[i] = 0x9a
		i++
		dAtA[i] = 0x6
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Advanced.Size()))
		n7, err := m.Advanced.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	return i, nil
}

func (m *Location) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Location) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Latitude != 0 {
		dAtA[i] = 0xd
		i++
		i = encodeFixed32Metadata(dAtA, i, uint32(math.Float32bits(float32(m.Latitude))))
	}
	if m.Longitude != 0 {
		dAtA[i] = 0x15
		i++
		i = encodeFixed32Metadata(dAtA, i, uint32(math.Float32bits(float32(m.Longitude))))
	}
	if m.Altitude != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Altitude))
	}
	if m.Accuracy != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Accuracy))
	}
	if m.Source != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintMetadata(dAtA, i, uint64(m.Source))
	}
	return i, nil
}

func encodeFixed64Metadata(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32Metadata(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintMetadata(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RxMetadata) Size() (n int) {
	var l int
	_ = l
	l = m.GatewayIdentifiers.Size()
	n += 1 + l + sovMetadata(uint64(l))
	if m.AntennaIndex != 0 {
		n += 1 + sovMetadata(uint64(m.AntennaIndex))
	}
	if m.ChannelIndex != 0 {
		n += 1 + sovMetadata(uint64(m.ChannelIndex))
	}
	if m.Timestamp != 0 {
		n += 1 + sovMetadata(uint64(m.Timestamp))
	}
	if m.Time != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.Time)
		n += 1 + l + sovMetadata(uint64(l))
	}
	if m.RSSI != 0 {
		n += 5
	}
	if m.SNR != 0 {
		n += 5
	}
	if m.Location != nil {
		l = m.Location.Size()
		n += 1 + l + sovMetadata(uint64(l))
	}
	if m.Advanced != nil {
		l = m.Advanced.Size()
		n += 2 + l + sovMetadata(uint64(l))
	}
	return n
}

func (m *TxMetadata) Size() (n int) {
	var l int
	_ = l
	l = m.GatewayIdentifiers.Size()
	n += 1 + l + sovMetadata(uint64(l))
	if m.Timestamp != 0 {
		n += 1 + sovMetadata(uint64(m.Timestamp))
	}
	if m.Time != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.Time)
		n += 1 + l + sovMetadata(uint64(l))
	}
	if m.Advanced != nil {
		l = m.Advanced.Size()
		n += 2 + l + sovMetadata(uint64(l))
	}
	return n
}

func (m *Location) Size() (n int) {
	var l int
	_ = l
	if m.Latitude != 0 {
		n += 5
	}
	if m.Longitude != 0 {
		n += 5
	}
	if m.Altitude != 0 {
		n += 1 + sovMetadata(uint64(m.Altitude))
	}
	if m.Accuracy != 0 {
		n += 1 + sovMetadata(uint64(m.Accuracy))
	}
	if m.Source != 0 {
		n += 1 + sovMetadata(uint64(m.Source))
	}
	return n
}

func sovMetadata(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMetadata(x uint64) (n int) {
	return sovMetadata(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *RxMetadata) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RxMetadata{`,
		`GatewayIdentifiers:` + strings.Replace(strings.Replace(this.GatewayIdentifiers.String(), "GatewayIdentifiers", "GatewayIdentifiers", 1), `&`, ``, 1) + `,`,
		`AntennaIndex:` + fmt.Sprintf("%v", this.AntennaIndex) + `,`,
		`ChannelIndex:` + fmt.Sprintf("%v", this.ChannelIndex) + `,`,
		`Timestamp:` + fmt.Sprintf("%v", this.Timestamp) + `,`,
		`Time:` + strings.Replace(fmt.Sprintf("%v", this.Time), "Timestamp", "google_protobuf3.Timestamp", 1) + `,`,
		`RSSI:` + fmt.Sprintf("%v", this.RSSI) + `,`,
		`SNR:` + fmt.Sprintf("%v", this.SNR) + `,`,
		`Location:` + strings.Replace(fmt.Sprintf("%v", this.Location), "Location", "Location", 1) + `,`,
		`Advanced:` + strings.Replace(fmt.Sprintf("%v", this.Advanced), "Struct", "google_protobuf2.Struct", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *TxMetadata) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&TxMetadata{`,
		`GatewayIdentifiers:` + strings.Replace(strings.Replace(this.GatewayIdentifiers.String(), "GatewayIdentifiers", "GatewayIdentifiers", 1), `&`, ``, 1) + `,`,
		`Timestamp:` + fmt.Sprintf("%v", this.Timestamp) + `,`,
		`Time:` + strings.Replace(fmt.Sprintf("%v", this.Time), "Timestamp", "google_protobuf3.Timestamp", 1) + `,`,
		`Advanced:` + strings.Replace(fmt.Sprintf("%v", this.Advanced), "Struct", "google_protobuf2.Struct", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Location) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Location{`,
		`Latitude:` + fmt.Sprintf("%v", this.Latitude) + `,`,
		`Longitude:` + fmt.Sprintf("%v", this.Longitude) + `,`,
		`Altitude:` + fmt.Sprintf("%v", this.Altitude) + `,`,
		`Accuracy:` + fmt.Sprintf("%v", this.Accuracy) + `,`,
		`Source:` + fmt.Sprintf("%v", this.Source) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringMetadata(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *RxMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMetadata
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
			return fmt.Errorf("proto: RxMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RxMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GatewayIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
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
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GatewayIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AntennaIndex", wireType)
			}
			m.AntennaIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AntennaIndex |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChannelIndex", wireType)
			}
			m.ChannelIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChannelIndex |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
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
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Time == nil {
				m.Time = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field RSSI", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += 4
			v = uint32(dAtA[iNdEx-4])
			v |= uint32(dAtA[iNdEx-3]) << 8
			v |= uint32(dAtA[iNdEx-2]) << 16
			v |= uint32(dAtA[iNdEx-1]) << 24
			m.RSSI = float32(math.Float32frombits(v))
		case 7:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field SNR", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += 4
			v = uint32(dAtA[iNdEx-4])
			v |= uint32(dAtA[iNdEx-3]) << 8
			v |= uint32(dAtA[iNdEx-2]) << 16
			v |= uint32(dAtA[iNdEx-1]) << 24
			m.SNR = float32(math.Float32frombits(v))
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
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
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Location == nil {
				m.Location = &Location{}
			}
			if err := m.Location.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 99:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Advanced", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
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
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Advanced == nil {
				m.Advanced = &google_protobuf2.Struct{}
			}
			if err := m.Advanced.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMetadata(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMetadata
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
func (m *TxMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMetadata
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
			return fmt.Errorf("proto: TxMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GatewayIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
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
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GatewayIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
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
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Time == nil {
				m.Time = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 99:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Advanced", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
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
				return ErrInvalidLengthMetadata
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Advanced == nil {
				m.Advanced = &google_protobuf2.Struct{}
			}
			if err := m.Advanced.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMetadata(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMetadata
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
func (m *Location) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMetadata
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
			return fmt.Errorf("proto: Location: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Location: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Latitude", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += 4
			v = uint32(dAtA[iNdEx-4])
			v |= uint32(dAtA[iNdEx-3]) << 8
			v |= uint32(dAtA[iNdEx-2]) << 16
			v |= uint32(dAtA[iNdEx-1]) << 24
			m.Latitude = float32(math.Float32frombits(v))
		case 2:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Longitude", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += 4
			v = uint32(dAtA[iNdEx-4])
			v |= uint32(dAtA[iNdEx-3]) << 8
			v |= uint32(dAtA[iNdEx-2]) << 16
			v |= uint32(dAtA[iNdEx-1]) << 24
			m.Longitude = float32(math.Float32frombits(v))
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Altitude", wireType)
			}
			m.Altitude = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Altitude |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Accuracy", wireType)
			}
			m.Accuracy = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Accuracy |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
			}
			m.Source = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetadata
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Source |= (LocationSource(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMetadata(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMetadata
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
func skipMetadata(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMetadata
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
					return 0, ErrIntOverflowMetadata
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
					return 0, ErrIntOverflowMetadata
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
				return 0, ErrInvalidLengthMetadata
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMetadata
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
				next, err := skipMetadata(dAtA[start:])
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
	ErrInvalidLengthMetadata = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMetadata   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/metadata.proto", fileDescriptorMetadata)
}

var fileDescriptorMetadata = []byte{
	// 617 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xcb, 0x6e, 0xd3, 0x4c,
	0x14, 0xce, 0x24, 0x6e, 0xea, 0xce, 0xff, 0x37, 0x84, 0x41, 0x80, 0x1b, 0x55, 0x4e, 0x54, 0x36,
	0x01, 0x81, 0x23, 0xb5, 0xf0, 0x00, 0x24, 0x0a, 0x91, 0x45, 0x71, 0xaa, 0x71, 0x2a, 0x04, 0x9b,
	0x68, 0x62, 0x4f, 0x1d, 0xab, 0xc9, 0x4c, 0x64, 0x4f, 0x7a, 0x59, 0x20, 0xb1, 0x44, 0x62, 0xc3,
	0x3b, 0xb0, 0xe5, 0x41, 0xba, 0xec, 0x12, 0xb1, 0x08, 0xe0, 0x15, 0x4b, 0x1e, 0x01, 0x79, 0x3c,
	0x49, 0x68, 0x11, 0xe2, 0x22, 0x56, 0xf6, 0x77, 0x39, 0x9e, 0xf3, 0x9d, 0x33, 0x32, 0xdc, 0x0e,
	0x42, 0x31, 0x9c, 0x0e, 0x2c, 0x8f, 0x8f, 0x1b, 0xbd, 0x21, 0xed, 0x0d, 0x43, 0x16, 0xc4, 0x0e,
	0x15, 0xc7, 0x3c, 0x3a, 0x6c, 0x08, 0xc1, 0x1a, 0x64, 0x12, 0x36, 0xc6, 0x54, 0x10, 0x9f, 0x08,
	0x62, 0x4d, 0x22, 0x2e, 0x38, 0x2a, 0x0a, 0xc1, 0xac, 0xa3, 0x9d, 0xca, 0xbd, 0xef, 0x6a, 0x03,
	0x1e, 0xf0, 0x86, 0x94, 0x07, 0xd3, 0x03, 0x89, 0x24, 0x90, 0x6f, 0x59, 0x59, 0xe5, 0xc1, 0xef,
	0x1c, 0x15, 0xfa, 0x94, 0x89, 0xf0, 0x20, 0xa4, 0x51, 0xac, 0xca, 0x36, 0x03, 0xce, 0x83, 0x11,
	0x5d, 0x7e, 0x3c, 0x16, 0xd1, 0xd4, 0x13, 0x4a, 0xad, 0x5e, 0x56, 0x45, 0x38, 0xa6, 0xb1, 0x20,
	0xe3, 0x49, 0x66, 0xd8, 0x7a, 0x5d, 0x80, 0x10, 0x9f, 0x3c, 0x51, 0x09, 0x50, 0x0b, 0xc2, 0x80,
	0x08, 0x7a, 0x4c, 0x4e, 0xfb, 0xa1, 0x6f, 0x80, 0x1a, 0xa8, 0xff, 0xb7, 0x5d, 0xb1, 0xb2, 0x40,
	0x56, 0x27, 0x53, 0xec, 0x65, 0x0f, 0x4d, 0xfd, 0x6c, 0x56, 0xcd, 0x9d, 0xcf, 0xaa, 0x00, 0xaf,
	0x05, 0x73, 0x15, 0xdd, 0x82, 0xeb, 0x84, 0x09, 0xca, 0x18, 0xe9, 0x87, 0xcc, 0xa7, 0x27, 0x46,
	0xbe, 0x06, 0xea, 0xeb, 0xf8, 0x7f, 0x45, 0xda, 0x29, 0x97, 0x9a, 0xbc, 0x21, 0x61, 0x8c, 0x8e,
	0x94, 0xa9, 0x90, 0x99, 0x14, 0x99, 0x99, 0x36, 0xe1, 0xda, 0xa2, 0x61, 0x43, 0xab, 0x81, 0xba,
	0x86, 0x97, 0x04, 0xba, 0x0f, 0xb5, 0x14, 0x18, 0x2b, 0xaa, 0xcd, 0x2c, 0xab, 0x35, 0xcf, 0x6a,
	0xf5, 0xe6, 0xce, 0xa6, 0xf6, 0xe6, 0x63, 0x15, 0x60, 0xe9, 0x46, 0x9b, 0x50, 0x8b, 0xe2, 0x38,
	0x34, 0x8a, 0x35, 0x50, 0xcf, 0x37, 0xf5, 0x64, 0x56, 0xd5, 0xb0, 0xeb, 0xda, 0x58, 0xb2, 0x68,
	0x03, 0x16, 0x62, 0x16, 0x19, 0xab, 0x52, 0x5c, 0x4d, 0x66, 0xd5, 0x82, 0xeb, 0x60, 0x9c, 0x72,
	0xe8, 0x2e, 0xd4, 0x47, 0xdc, 0x23, 0x22, 0xe4, 0xcc, 0xd0, 0xe5, 0x91, 0xe5, 0xf9, 0x64, 0x76,
	0x15, 0x8f, 0x17, 0x0e, 0xb4, 0x03, 0x75, 0xe2, 0x1f, 0x11, 0xe6, 0x51, 0xdf, 0xf0, 0xa4, 0xfb,
	0xe6, 0x0f, 0x0d, 0xba, 0x72, 0x55, 0x78, 0x61, 0xdc, 0xfa, 0x00, 0x20, 0xec, 0xfd, 0xe3, 0x6d,
	0x5c, 0x98, 0x61, 0xfe, 0x67, 0x33, 0x2c, 0xfc, 0xd1, 0x0c, 0xff, 0x2a, 0xdc, 0x3b, 0x00, 0xf5,
	0xf9, 0xa0, 0x50, 0x05, 0xea, 0x23, 0x22, 0x42, 0x31, 0xf5, 0xa9, 0x0c, 0x96, 0xc7, 0x0b, 0x9c,
	0x76, 0x3c, 0xe2, 0x2c, 0xc8, 0xc4, 0xbc, 0x14, 0x97, 0x44, 0x5a, 0x49, 0x46, 0xaa, 0x32, 0xed,
	0x7a, 0x05, 0x2f, 0xb0, 0xd4, 0x3c, 0x6f, 0x1a, 0x11, 0xef, 0x54, 0x5e, 0x97, 0x54, 0x53, 0x18,
	0x59, 0xb0, 0x18, 0xf3, 0x69, 0xe4, 0x65, 0xf7, 0xa5, 0xb4, 0x7d, 0xe3, 0xf2, 0xf2, 0x5c, 0xa9,
	0x62, 0xe5, 0xba, 0xf3, 0x02, 0x96, 0x2e, 0x2a, 0x08, 0xc1, 0x92, 0xdb, 0xdd, 0xc7, 0xad, 0x76,
	0x7f, 0xdf, 0x79, 0xec, 0x74, 0x9f, 0x3a, 0xe5, 0x1c, 0x2a, 0x41, 0xa8, 0xb8, 0xce, 0x9e, 0x5b,
	0x06, 0xe8, 0x2a, 0x5c, 0x57, 0xb8, 0xd5, 0x75, 0x1e, 0xd9, 0x9d, 0x72, 0x1e, 0x5d, 0x83, 0x57,
	0x14, 0x85, 0xdb, 0x1d, 0xdb, 0xed, 0xe1, 0x67, 0xe5, 0x02, 0xda, 0x80, 0xd7, 0x15, 0x69, 0xef,
	0xf5, 0x3b, 0xed, 0xee, 0x6e, 0xb7, 0xf5, 0xb0, 0x67, 0x77, 0x9d, 0xb2, 0x56, 0xd1, 0x5e, 0xbd,
	0x35, 0x73, 0xcd, 0xce, 0xfb, 0xcf, 0x66, 0xee, 0x65, 0x62, 0x82, 0xb3, 0xc4, 0x04, 0xe7, 0x89,
	0x09, 0x3e, 0x25, 0x26, 0xf8, 0x92, 0x98, 0xb9, 0xaf, 0x89, 0x09, 0x9e, 0xdf, 0xfe, 0xd5, 0xcf,
	0x62, 0x72, 0x18, 0xa4, 0xcf, 0xc9, 0x60, 0x50, 0x94, 0x1b, 0xd9, 0xf9, 0x16, 0x00, 0x00, 0xff,
	0xff, 0x43, 0xc8, 0x9f, 0xa7, 0xcb, 0x04, 0x00, 0x00,
}
