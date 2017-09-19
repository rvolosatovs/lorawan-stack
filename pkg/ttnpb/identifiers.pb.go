// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/identifiers.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import github_com_TheThingsNetwork_ttn_pkg_types "github.com/TheThingsNetwork/ttn/pkg/types"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type GatewayIdentifiers struct {
	// TTN Gateway ID
	GatewayID string `protobuf:"bytes,1,opt,name=gateway_id,json=gatewayId,proto3" json:"gateway_id,omitempty"`
	// TTN Tenant ID (in case of multi-tenant network stack)
	TenantID string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
}

func (m *GatewayIdentifiers) Reset()                    { *m = GatewayIdentifiers{} }
func (*GatewayIdentifiers) ProtoMessage()               {}
func (*GatewayIdentifiers) Descriptor() ([]byte, []int) { return fileDescriptorIdentifiers, []int{0} }

func (m *GatewayIdentifiers) GetGatewayID() string {
	if m != nil {
		return m.GatewayID
	}
	return ""
}

func (m *GatewayIdentifiers) GetTenantID() string {
	if m != nil {
		return m.TenantID
	}
	return ""
}

// End device identifiers are carried with uplink and downlink messages
// Unknown fields are left empty
type EndDeviceIdentifiers struct {
	// TTN Device ID
	DeviceID string `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	// TTN Application ID
	ApplicationID string `protobuf:"bytes,2,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	// TTN Tenant ID (in case of multi-tenant network server)
	TenantID string `protobuf:"bytes,3,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	// LoRaWAN DevEUI
	DevEUI *github_com_TheThingsNetwork_ttn_pkg_types.EUI64 `protobuf:"bytes,4,opt,name=dev_eui,json=devEui,proto3,customtype=github.com/TheThingsNetwork/ttn/pkg/types.EUI64" json:"dev_eui,omitempty"`
	// LoRaWAN JoinEUI (or AppEUI for LoRaWAN 1.0 end devices)
	JoinEUI *github_com_TheThingsNetwork_ttn_pkg_types.EUI64 `protobuf:"bytes,5,opt,name=join_eui,json=joinEui,proto3,customtype=github.com/TheThingsNetwork/ttn/pkg/types.EUI64" json:"join_eui,omitempty"`
	// LoRaWAN DevAddr
	DevAddr *github_com_TheThingsNetwork_ttn_pkg_types.DevAddr `protobuf:"bytes,6,opt,name=dev_addr,json=devAddr,proto3,customtype=github.com/TheThingsNetwork/ttn/pkg/types.DevAddr" json:"dev_addr,omitempty"`
}

func (m *EndDeviceIdentifiers) Reset()                    { *m = EndDeviceIdentifiers{} }
func (*EndDeviceIdentifiers) ProtoMessage()               {}
func (*EndDeviceIdentifiers) Descriptor() ([]byte, []int) { return fileDescriptorIdentifiers, []int{1} }

func (m *EndDeviceIdentifiers) GetDeviceID() string {
	if m != nil {
		return m.DeviceID
	}
	return ""
}

func (m *EndDeviceIdentifiers) GetApplicationID() string {
	if m != nil {
		return m.ApplicationID
	}
	return ""
}

func (m *EndDeviceIdentifiers) GetTenantID() string {
	if m != nil {
		return m.TenantID
	}
	return ""
}

type ApplicationIdentifiers struct {
	// TTN Application ID
	ApplicationID string `protobuf:"bytes,1,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	// TTN Tenant ID (in case of multi-tenant network stack)
	TenantID string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
}

func (m *ApplicationIdentifiers) Reset()      { *m = ApplicationIdentifiers{} }
func (*ApplicationIdentifiers) ProtoMessage() {}
func (*ApplicationIdentifiers) Descriptor() ([]byte, []int) {
	return fileDescriptorIdentifiers, []int{2}
}

func (m *ApplicationIdentifiers) GetApplicationID() string {
	if m != nil {
		return m.ApplicationID
	}
	return ""
}

func (m *ApplicationIdentifiers) GetTenantID() string {
	if m != nil {
		return m.TenantID
	}
	return ""
}

type UserIdentifiers struct {
	// username of the user
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// email address of the user
	Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	// TTN Tenant ID (in case of multi-tenant network stack)
	TenantID string `protobuf:"bytes,3,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
}

func (m *UserIdentifiers) Reset()                    { *m = UserIdentifiers{} }
func (*UserIdentifiers) ProtoMessage()               {}
func (*UserIdentifiers) Descriptor() ([]byte, []int) { return fileDescriptorIdentifiers, []int{3} }

func (m *UserIdentifiers) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserIdentifiers) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserIdentifiers) GetTenantID() string {
	if m != nil {
		return m.TenantID
	}
	return ""
}

type ClientIdentifiers struct {
	// TTN Client ID
	ClientID string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// TTN Tenant ID (in case of multi-tenant network stack)
	TenantID string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
}

func (m *ClientIdentifiers) Reset()                    { *m = ClientIdentifiers{} }
func (*ClientIdentifiers) ProtoMessage()               {}
func (*ClientIdentifiers) Descriptor() ([]byte, []int) { return fileDescriptorIdentifiers, []int{4} }

func (m *ClientIdentifiers) GetClientID() string {
	if m != nil {
		return m.ClientID
	}
	return ""
}

func (m *ClientIdentifiers) GetTenantID() string {
	if m != nil {
		return m.TenantID
	}
	return ""
}

func init() {
	proto.RegisterType((*GatewayIdentifiers)(nil), "ttn.v3.GatewayIdentifiers")
	proto.RegisterType((*EndDeviceIdentifiers)(nil), "ttn.v3.EndDeviceIdentifiers")
	proto.RegisterType((*ApplicationIdentifiers)(nil), "ttn.v3.ApplicationIdentifiers")
	proto.RegisterType((*UserIdentifiers)(nil), "ttn.v3.UserIdentifiers")
	proto.RegisterType((*ClientIdentifiers)(nil), "ttn.v3.ClientIdentifiers")
}
func (m *GatewayIdentifiers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GatewayIdentifiers) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.GatewayID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.GatewayID)))
		i += copy(dAtA[i:], m.GatewayID)
	}
	if len(m.TenantID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.TenantID)))
		i += copy(dAtA[i:], m.TenantID)
	}
	return i, nil
}

func (m *EndDeviceIdentifiers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EndDeviceIdentifiers) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.DeviceID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.DeviceID)))
		i += copy(dAtA[i:], m.DeviceID)
	}
	if len(m.ApplicationID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.ApplicationID)))
		i += copy(dAtA[i:], m.ApplicationID)
	}
	if len(m.TenantID) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.TenantID)))
		i += copy(dAtA[i:], m.TenantID)
	}
	if m.DevEUI != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(m.DevEUI.Size()))
		n1, err := m.DevEUI.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.JoinEUI != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(m.JoinEUI.Size()))
		n2, err := m.JoinEUI.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.DevAddr != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(m.DevAddr.Size()))
		n3, err := m.DevAddr.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func (m *ApplicationIdentifiers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplicationIdentifiers) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ApplicationID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.ApplicationID)))
		i += copy(dAtA[i:], m.ApplicationID)
	}
	if len(m.TenantID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.TenantID)))
		i += copy(dAtA[i:], m.TenantID)
	}
	return i, nil
}

func (m *UserIdentifiers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UserIdentifiers) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Username) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.Username)))
		i += copy(dAtA[i:], m.Username)
	}
	if len(m.Email) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.Email)))
		i += copy(dAtA[i:], m.Email)
	}
	if len(m.TenantID) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.TenantID)))
		i += copy(dAtA[i:], m.TenantID)
	}
	return i, nil
}

func (m *ClientIdentifiers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClientIdentifiers) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ClientID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.ClientID)))
		i += copy(dAtA[i:], m.ClientID)
	}
	if len(m.TenantID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIdentifiers(dAtA, i, uint64(len(m.TenantID)))
		i += copy(dAtA[i:], m.TenantID)
	}
	return i, nil
}

func encodeFixed64Identifiers(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32Identifiers(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintIdentifiers(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *GatewayIdentifiers) Size() (n int) {
	var l int
	_ = l
	l = len(m.GatewayID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	l = len(m.TenantID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	return n
}

func (m *EndDeviceIdentifiers) Size() (n int) {
	var l int
	_ = l
	l = len(m.DeviceID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	l = len(m.ApplicationID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	l = len(m.TenantID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	if m.DevEUI != nil {
		l = m.DevEUI.Size()
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	if m.JoinEUI != nil {
		l = m.JoinEUI.Size()
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	if m.DevAddr != nil {
		l = m.DevAddr.Size()
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	return n
}

func (m *ApplicationIdentifiers) Size() (n int) {
	var l int
	_ = l
	l = len(m.ApplicationID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	l = len(m.TenantID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	return n
}

func (m *UserIdentifiers) Size() (n int) {
	var l int
	_ = l
	l = len(m.Username)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	l = len(m.TenantID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	return n
}

func (m *ClientIdentifiers) Size() (n int) {
	var l int
	_ = l
	l = len(m.ClientID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	l = len(m.TenantID)
	if l > 0 {
		n += 1 + l + sovIdentifiers(uint64(l))
	}
	return n
}

func sovIdentifiers(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozIdentifiers(x uint64) (n int) {
	return sovIdentifiers(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *GatewayIdentifiers) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&GatewayIdentifiers{`,
		`GatewayID:` + fmt.Sprintf("%v", this.GatewayID) + `,`,
		`TenantID:` + fmt.Sprintf("%v", this.TenantID) + `,`,
		`}`,
	}, "")
	return s
}
func (this *EndDeviceIdentifiers) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&EndDeviceIdentifiers{`,
		`DeviceID:` + fmt.Sprintf("%v", this.DeviceID) + `,`,
		`ApplicationID:` + fmt.Sprintf("%v", this.ApplicationID) + `,`,
		`TenantID:` + fmt.Sprintf("%v", this.TenantID) + `,`,
		`DevEUI:` + fmt.Sprintf("%v", this.DevEUI) + `,`,
		`JoinEUI:` + fmt.Sprintf("%v", this.JoinEUI) + `,`,
		`DevAddr:` + fmt.Sprintf("%v", this.DevAddr) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ApplicationIdentifiers) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplicationIdentifiers{`,
		`ApplicationID:` + fmt.Sprintf("%v", this.ApplicationID) + `,`,
		`TenantID:` + fmt.Sprintf("%v", this.TenantID) + `,`,
		`}`,
	}, "")
	return s
}
func (this *UserIdentifiers) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&UserIdentifiers{`,
		`Username:` + fmt.Sprintf("%v", this.Username) + `,`,
		`Email:` + fmt.Sprintf("%v", this.Email) + `,`,
		`TenantID:` + fmt.Sprintf("%v", this.TenantID) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ClientIdentifiers) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ClientIdentifiers{`,
		`ClientID:` + fmt.Sprintf("%v", this.ClientID) + `,`,
		`TenantID:` + fmt.Sprintf("%v", this.TenantID) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringIdentifiers(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *GatewayIdentifiers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentifiers
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
			return fmt.Errorf("proto: GatewayIdentifiers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GatewayIdentifiers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GatewayID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GatewayID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TenantID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TenantID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentifiers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIdentifiers
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
func (m *EndDeviceIdentifiers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentifiers
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
			return fmt.Errorf("proto: EndDeviceIdentifiers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EndDeviceIdentifiers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeviceID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeviceID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApplicationID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ApplicationID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TenantID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TenantID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DevEUI", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_TheThingsNetwork_ttn_pkg_types.EUI64
			m.DevEUI = &v
			if err := m.DevEUI.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JoinEUI", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_TheThingsNetwork_ttn_pkg_types.EUI64
			m.JoinEUI = &v
			if err := m.JoinEUI.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DevAddr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_TheThingsNetwork_ttn_pkg_types.DevAddr
			m.DevAddr = &v
			if err := m.DevAddr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentifiers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIdentifiers
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
func (m *ApplicationIdentifiers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentifiers
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
			return fmt.Errorf("proto: ApplicationIdentifiers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplicationIdentifiers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApplicationID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ApplicationID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TenantID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TenantID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentifiers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIdentifiers
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
func (m *UserIdentifiers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentifiers
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
			return fmt.Errorf("proto: UserIdentifiers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UserIdentifiers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Username", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Username = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TenantID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TenantID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentifiers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIdentifiers
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
func (m *ClientIdentifiers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIdentifiers
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
			return fmt.Errorf("proto: ClientIdentifiers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClientIdentifiers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TenantID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIdentifiers
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
				return ErrInvalidLengthIdentifiers
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TenantID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIdentifiers(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIdentifiers
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
func skipIdentifiers(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIdentifiers
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
					return 0, ErrIntOverflowIdentifiers
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
					return 0, ErrIntOverflowIdentifiers
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
				return 0, ErrInvalidLengthIdentifiers
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowIdentifiers
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
				next, err := skipIdentifiers(dAtA[start:])
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
	ErrInvalidLengthIdentifiers = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIdentifiers   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/identifiers.proto", fileDescriptorIdentifiers)
}

var fileDescriptorIdentifiers = []byte{
	// 493 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x41, 0x6b, 0xd4, 0x40,
	0x14, 0xc7, 0x1b, 0x6b, 0xb3, 0xc9, 0xd0, 0x55, 0x1a, 0x8a, 0x2c, 0x3d, 0x24, 0x65, 0x4f, 0x5d,
	0xd0, 0x0d, 0x52, 0x2b, 0x5e, 0x14, 0xba, 0x66, 0x29, 0xf1, 0x20, 0x32, 0xec, 0x0a, 0x7a, 0x59,
	0x66, 0x33, 0xaf, 0xd9, 0xb1, 0xbb, 0x33, 0x21, 0x99, 0xa4, 0xf4, 0x20, 0xf8, 0x7d, 0xfc, 0x22,
	0x1e, 0x3d, 0x4a, 0x0f, 0x41, 0xe7, 0xe4, 0xd1, 0x8f, 0x20, 0x99, 0x6c, 0xdd, 0x58, 0xa1, 0xb6,
	0xbd, 0xcd, 0x7b, 0xef, 0xff, 0x7e, 0x2f, 0x33, 0xff, 0x17, 0x74, 0x10, 0x33, 0x39, 0xcb, 0xa7,
	0xfd, 0x48, 0x2c, 0xfc, 0xd1, 0x0c, 0x46, 0x33, 0xc6, 0xe3, 0xec, 0x35, 0xc8, 0x53, 0x91, 0x9e,
	0xf8, 0x52, 0x72, 0x9f, 0x24, 0xcc, 0x67, 0x14, 0xb8, 0x64, 0xc7, 0x0c, 0xd2, 0xac, 0x9f, 0xa4,
	0x42, 0x0a, 0xc7, 0x94, 0x92, 0xf7, 0x8b, 0xfd, 0x9d, 0x47, 0x8d, 0xf6, 0x58, 0xc4, 0xc2, 0xd7,
	0xe5, 0x69, 0x7e, 0xac, 0x23, 0x1d, 0xe8, 0x53, 0xdd, 0xd6, 0x5d, 0x20, 0xe7, 0x88, 0x48, 0x38,
	0x25, 0x67, 0xe1, 0x0a, 0xe9, 0x3c, 0x44, 0x28, 0xae, 0xb3, 0x13, 0x46, 0x3b, 0xc6, 0xae, 0xb1,
	0x67, 0x0f, 0xda, 0xaa, 0xf4, 0xec, 0x0b, 0x6d, 0x80, 0xed, 0xf8, 0xa2, 0xcd, 0xe9, 0x21, 0x5b,
	0x02, 0x27, 0x5c, 0x56, 0xe2, 0x3b, 0x5a, 0xbc, 0xa9, 0x4a, 0xcf, 0x1a, 0xe9, 0x64, 0x18, 0x60,
	0xab, 0x2e, 0x87, 0xb4, 0xfb, 0x79, 0x1d, 0x6d, 0x0f, 0x39, 0x0d, 0xa0, 0x60, 0x11, 0x34, 0x27,
	0xf6, 0x90, 0x4d, 0x75, 0x72, 0x35, 0x50, 0x33, 0x96, 0xca, 0x00, 0x5b, 0x74, 0xd9, 0xe3, 0x3c,
	0x43, 0xf7, 0x48, 0x92, 0xcc, 0x59, 0x44, 0x24, 0x13, 0x7c, 0x35, 0x73, 0x4b, 0x95, 0x5e, 0xfb,
	0x70, 0x55, 0x09, 0x03, 0xdc, 0x6e, 0x08, 0x2f, 0x7f, 0xe8, 0xfa, 0x55, 0x1f, 0xea, 0xbc, 0x45,
	0x2d, 0x0a, 0xc5, 0x04, 0x72, 0xd6, 0xb9, 0xbb, 0x6b, 0xec, 0x6d, 0x0e, 0x9e, 0x9f, 0x97, 0x9e,
	0xff, 0x3f, 0x6b, 0x92, 0x93, 0xd8, 0x97, 0x67, 0x09, 0x64, 0xfd, 0xe1, 0x38, 0x7c, 0xfa, 0x44,
	0x95, 0x9e, 0x19, 0x40, 0x31, 0x1c, 0x87, 0xd8, 0xa4, 0x50, 0x0c, 0x73, 0xe6, 0xbc, 0x43, 0xd6,
	0x07, 0xc1, 0xb8, 0x06, 0x6f, 0x68, 0xf0, 0x8b, 0xdb, 0x81, 0x5b, 0xaf, 0x04, 0xe3, 0x15, 0xb9,
	0x55, 0xf1, 0x2a, 0xf4, 0x1b, 0x54, 0xbd, 0xd1, 0x84, 0x50, 0x9a, 0x76, 0x4c, 0x8d, 0x3e, 0x38,
	0x2f, 0xbd, 0xc7, 0xd7, 0x47, 0x07, 0x50, 0x1c, 0x52, 0x9a, 0xe2, 0xea, 0xe6, 0xd5, 0xa1, 0xfb,
	0x11, 0x3d, 0x68, 0xbe, 0x67, 0xc3, 0xae, 0x7f, 0x3d, 0x30, 0x6e, 0xe3, 0xc1, 0xd5, 0xcb, 0xc2,
	0xd1, 0xfd, 0x71, 0x06, 0x69, 0x73, 0xee, 0x0e, 0xb2, 0xf2, 0x0c, 0x52, 0x4e, 0x16, 0x50, 0x4f,
	0xc4, 0x7f, 0x62, 0x67, 0x1b, 0x6d, 0xc0, 0x82, 0xb0, 0x79, 0x4d, 0xc5, 0x75, 0x70, 0x03, 0xcf,
	0xbb, 0x0c, 0x6d, 0xbd, 0x9c, 0x33, 0xa8, 0xce, 0x7f, 0x2d, 0x66, 0xa4, 0x93, 0x97, 0x16, 0x73,
	0xa9, 0x0c, 0xb0, 0x15, 0x2d, 0x7b, 0x6e, 0x70, 0xb5, 0xc1, 0xd1, 0xb7, 0x1f, 0xee, 0xda, 0x27,
	0xe5, 0x1a, 0x5f, 0x94, 0x6b, 0x7c, 0x55, 0xae, 0xf1, 0x5d, 0xb9, 0xc6, 0x4f, 0xe5, 0xae, 0xfd,
	0x52, 0xae, 0xf1, 0xbe, 0x77, 0x2d, 0xcf, 0x24, 0x4f, 0xa6, 0x53, 0x53, 0xff, 0xc6, 0xfb, 0xbf,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x55, 0x1a, 0x65, 0x68, 0x36, 0x04, 0x00, 0x00,
}
