// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/organization.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import _ "github.com/gogo/protobuf/types"

import time "time"

import types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

type Organization struct {
	// Identifiers of the organization.
	OrganizationIdentifiers `protobuf:"bytes,1,opt,name=ids,embedded=ids" json:"ids"`
	// name is the organization's name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// description is an organization's description.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// url is the URL of the organization website.
	URL string `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	// location is the organization's location, e.g. "Amsterdam, Europe".
	Location string `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	// email is a generic contact email address that is shown as contact email.
	Email string `protobuf:"bytes,6,opt,name=email,proto3" json:"email,omitempty"`
	// created_at denotes when the user was created.
	// This is a read-only field.
	CreatedAt time.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt,stdtime" json:"created_at"`
	// updated_at is the last time the user was updated.
	// This is a read-only field.
	UpdatedAt time.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,stdtime" json:"updated_at"`
}

func (m *Organization) Reset()                    { *m = Organization{} }
func (m *Organization) String() string            { return proto.CompactTextString(m) }
func (*Organization) ProtoMessage()               {}
func (*Organization) Descriptor() ([]byte, []int) { return fileDescriptorOrganization, []int{0} }

func (m *Organization) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Organization) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Organization) GetURL() string {
	if m != nil {
		return m.URL
	}
	return ""
}

func (m *Organization) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Organization) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Organization) GetCreatedAt() time.Time {
	if m != nil {
		return m.CreatedAt
	}
	return time.Time{}
}

func (m *Organization) GetUpdatedAt() time.Time {
	if m != nil {
		return m.UpdatedAt
	}
	return time.Time{}
}

type OrganizationMember struct {
	// organization_ids are the organization's identifiers.
	OrganizationIdentifiers `protobuf:"bytes,1,opt,name=organization_ids,json=organizationIds,embedded=organization_ids" json:"organization_ids"`
	// user_ids are the user's identifiers.
	UserIdentifiers `protobuf:"bytes,2,opt,name=user_ids,json=userIds,embedded=user_ids" json:"user_ids"`
	// rights is the list of rights the user bears to the organization.
	Rights []Right `protobuf:"varint,3,rep,packed,name=rights,enum=ttn.v3.Right" json:"rights,omitempty"`
}

func (m *OrganizationMember) Reset()                    { *m = OrganizationMember{} }
func (m *OrganizationMember) String() string            { return proto.CompactTextString(m) }
func (*OrganizationMember) ProtoMessage()               {}
func (*OrganizationMember) Descriptor() ([]byte, []int) { return fileDescriptorOrganization, []int{1} }

func (m *OrganizationMember) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

func init() {
	proto.RegisterType((*Organization)(nil), "ttn.v3.Organization")
	golang_proto.RegisterType((*Organization)(nil), "ttn.v3.Organization")
	proto.RegisterType((*OrganizationMember)(nil), "ttn.v3.OrganizationMember")
	golang_proto.RegisterType((*OrganizationMember)(nil), "ttn.v3.OrganizationMember")
}
func (this *Organization) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*Organization)
	if !ok {
		that2, ok := that.(Organization)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *Organization")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *Organization but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *Organization but is not nil && this == nil")
	}
	if !this.OrganizationIdentifiers.Equal(&that1.OrganizationIdentifiers) {
		return fmt.Errorf("OrganizationIdentifiers this(%v) Not Equal that(%v)", this.OrganizationIdentifiers, that1.OrganizationIdentifiers)
	}
	if this.Name != that1.Name {
		return fmt.Errorf("Name this(%v) Not Equal that(%v)", this.Name, that1.Name)
	}
	if this.Description != that1.Description {
		return fmt.Errorf("Description this(%v) Not Equal that(%v)", this.Description, that1.Description)
	}
	if this.URL != that1.URL {
		return fmt.Errorf("URL this(%v) Not Equal that(%v)", this.URL, that1.URL)
	}
	if this.Location != that1.Location {
		return fmt.Errorf("Location this(%v) Not Equal that(%v)", this.Location, that1.Location)
	}
	if this.Email != that1.Email {
		return fmt.Errorf("Email this(%v) Not Equal that(%v)", this.Email, that1.Email)
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return fmt.Errorf("CreatedAt this(%v) Not Equal that(%v)", this.CreatedAt, that1.CreatedAt)
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return fmt.Errorf("UpdatedAt this(%v) Not Equal that(%v)", this.UpdatedAt, that1.UpdatedAt)
	}
	return nil
}
func (this *Organization) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Organization)
	if !ok {
		that2, ok := that.(Organization)
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
	if !this.OrganizationIdentifiers.Equal(&that1.OrganizationIdentifiers) {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if this.URL != that1.URL {
		return false
	}
	if this.Location != that1.Location {
		return false
	}
	if this.Email != that1.Email {
		return false
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return false
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return false
	}
	return true
}
func (this *OrganizationMember) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*OrganizationMember)
	if !ok {
		that2, ok := that.(OrganizationMember)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *OrganizationMember")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *OrganizationMember but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *OrganizationMember but is not nil && this == nil")
	}
	if !this.OrganizationIdentifiers.Equal(&that1.OrganizationIdentifiers) {
		return fmt.Errorf("OrganizationIdentifiers this(%v) Not Equal that(%v)", this.OrganizationIdentifiers, that1.OrganizationIdentifiers)
	}
	if !this.UserIdentifiers.Equal(&that1.UserIdentifiers) {
		return fmt.Errorf("UserIdentifiers this(%v) Not Equal that(%v)", this.UserIdentifiers, that1.UserIdentifiers)
	}
	if len(this.Rights) != len(that1.Rights) {
		return fmt.Errorf("Rights this(%v) Not Equal that(%v)", len(this.Rights), len(that1.Rights))
	}
	for i := range this.Rights {
		if this.Rights[i] != that1.Rights[i] {
			return fmt.Errorf("Rights this[%v](%v) Not Equal that[%v](%v)", i, this.Rights[i], i, that1.Rights[i])
		}
	}
	return nil
}
func (this *OrganizationMember) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrganizationMember)
	if !ok {
		that2, ok := that.(OrganizationMember)
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
	if !this.OrganizationIdentifiers.Equal(&that1.OrganizationIdentifiers) {
		return false
	}
	if !this.UserIdentifiers.Equal(&that1.UserIdentifiers) {
		return false
	}
	if len(this.Rights) != len(that1.Rights) {
		return false
	}
	for i := range this.Rights {
		if this.Rights[i] != that1.Rights[i] {
			return false
		}
	}
	return true
}
func (m *Organization) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Organization) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(m.OrganizationIdentifiers.Size()))
	n1, err := m.OrganizationIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Description) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Description)))
		i += copy(dAtA[i:], m.Description)
	}
	if len(m.URL) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.URL)))
		i += copy(dAtA[i:], m.URL)
	}
	if len(m.Location) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Location)))
		i += copy(dAtA[i:], m.Location)
	}
	if len(m.Email) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Email)))
		i += copy(dAtA[i:], m.Email)
	}
	dAtA[i] = 0x3a
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(types.SizeOfStdTime(m.CreatedAt)))
	n2, err := types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	dAtA[i] = 0x42
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(types.SizeOfStdTime(m.UpdatedAt)))
	n3, err := types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func (m *OrganizationMember) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrganizationMember) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(m.OrganizationIdentifiers.Size()))
	n4, err := m.OrganizationIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	dAtA[i] = 0x12
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(m.UserIdentifiers.Size()))
	n5, err := m.UserIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	if len(m.Rights) > 0 {
		dAtA7 := make([]byte, len(m.Rights)*10)
		var j6 int
		for _, num := range m.Rights {
			for num >= 1<<7 {
				dAtA7[j6] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j6++
			}
			dAtA7[j6] = uint8(num)
			j6++
		}
		dAtA[i] = 0x1a
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(j6))
		i += copy(dAtA[i:], dAtA7[:j6])
	}
	return i, nil
}

func encodeVarintOrganization(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedOrganization(r randyOrganization, easy bool) *Organization {
	this := &Organization{}
	v1 := NewPopulatedOrganizationIdentifiers(r, easy)
	this.OrganizationIdentifiers = *v1
	this.Name = randStringOrganization(r)
	this.Description = randStringOrganization(r)
	this.URL = randStringOrganization(r)
	this.Location = randStringOrganization(r)
	this.Email = randStringOrganization(r)
	v2 := types.NewPopulatedStdTime(r, easy)
	this.CreatedAt = *v2
	v3 := types.NewPopulatedStdTime(r, easy)
	this.UpdatedAt = *v3
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedOrganizationMember(r randyOrganization, easy bool) *OrganizationMember {
	this := &OrganizationMember{}
	v4 := NewPopulatedOrganizationIdentifiers(r, easy)
	this.OrganizationIdentifiers = *v4
	v5 := NewPopulatedUserIdentifiers(r, easy)
	this.UserIdentifiers = *v5
	v6 := r.Intn(10)
	this.Rights = make([]Right, v6)
	for i := 0; i < v6; i++ {
		this.Rights[i] = Right([]int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43}[r.Intn(44)])
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyOrganization interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneOrganization(r randyOrganization) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringOrganization(r randyOrganization) string {
	v7 := r.Intn(100)
	tmps := make([]rune, v7)
	for i := 0; i < v7; i++ {
		tmps[i] = randUTF8RuneOrganization(r)
	}
	return string(tmps)
}
func randUnrecognizedOrganization(r randyOrganization, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldOrganization(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldOrganization(dAtA []byte, r randyOrganization, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		v8 := r.Int63()
		if r.Intn(2) == 0 {
			v8 *= -1
		}
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(v8))
	case 1:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateOrganization(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *Organization) Size() (n int) {
	var l int
	_ = l
	l = m.OrganizationIdentifiers.Size()
	n += 1 + l + sovOrganization(uint64(l))
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.URL)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.Location)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovOrganization(uint64(l))
	l = types.SizeOfStdTime(m.UpdatedAt)
	n += 1 + l + sovOrganization(uint64(l))
	return n
}

func (m *OrganizationMember) Size() (n int) {
	var l int
	_ = l
	l = m.OrganizationIdentifiers.Size()
	n += 1 + l + sovOrganization(uint64(l))
	l = m.UserIdentifiers.Size()
	n += 1 + l + sovOrganization(uint64(l))
	if len(m.Rights) > 0 {
		l = 0
		for _, e := range m.Rights {
			l += sovOrganization(uint64(e))
		}
		n += 1 + sovOrganization(uint64(l)) + l
	}
	return n
}

func sovOrganization(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozOrganization(x uint64) (n int) {
	return sovOrganization((x << 1) ^ uint64((int64(x) >> 63)))
}
func (m *Organization) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrganization
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
			return fmt.Errorf("proto: Organization: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Organization: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrganizationIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OrganizationIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field URL", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.URL = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Location = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := types.StdTimeUnmarshal(&m.CreatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := types.StdTimeUnmarshal(&m.UpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOrganization(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOrganization
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
func (m *OrganizationMember) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrganization
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
			return fmt.Errorf("proto: OrganizationMember: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrganizationMember: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrganizationIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OrganizationIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.UserIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v Right
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowOrganization
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (Right(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Rights = append(m.Rights, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowOrganization
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthOrganization
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v Right
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowOrganization
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (Right(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Rights = append(m.Rights, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Rights", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipOrganization(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOrganization
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
func skipOrganization(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOrganization
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
					return 0, ErrIntOverflowOrganization
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
					return 0, ErrIntOverflowOrganization
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
				return 0, ErrInvalidLengthOrganization
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowOrganization
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
				next, err := skipOrganization(dAtA[start:])
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
	ErrInvalidLengthOrganization = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOrganization   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/organization.proto", fileDescriptorOrganization)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/organization.proto", fileDescriptorOrganization)
}

var fileDescriptorOrganization = []byte{
	// 553 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x3f, 0x6c, 0xd3, 0x4e,
	0x14, 0xc7, 0xef, 0x35, 0x6d, 0x9a, 0x5e, 0x7f, 0x3f, 0x40, 0x27, 0x24, 0x4c, 0x90, 0x5e, 0xa2,
	0x4a, 0x48, 0x65, 0xc0, 0x41, 0xad, 0x60, 0x81, 0xa5, 0x65, 0x42, 0xe2, 0x8f, 0x64, 0xa5, 0x0b,
	0x4b, 0xe5, 0x24, 0x57, 0xe7, 0xd4, 0xd8, 0x67, 0xd9, 0x67, 0x10, 0x4c, 0x1d, 0x3b, 0x76, 0x64,
	0x83, 0xb1, 0x63, 0xc7, 0x8e, 0x1d, 0xbb, 0xd1, 0xb1, 0x53, 0xa9, 0xcf, 0x4b, 0xc7, 0x8e, 0x1d,
	0x91, 0xcf, 0x36, 0x35, 0x08, 0x89, 0xc0, 0xe4, 0x7b, 0xef, 0xfb, 0x3e, 0x5f, 0xeb, 0x7d, 0x7d,
	0xa6, 0x4f, 0x3c, 0xa1, 0xc6, 0xc9, 0xc0, 0x1e, 0x4a, 0xbf, 0xd7, 0x1f, 0xf3, 0xfe, 0x58, 0x04,
	0x5e, 0xfc, 0x9a, 0xab, 0xf7, 0x32, 0xda, 0xee, 0x29, 0x15, 0xf4, 0xdc, 0x50, 0xf4, 0x64, 0xe4,
	0xb9, 0x81, 0xf8, 0xe8, 0x2a, 0x21, 0x03, 0x3b, 0x8c, 0xa4, 0x92, 0xac, 0xa9, 0x54, 0x60, 0xbf,
	0x5b, 0x6d, 0xdb, 0xd3, 0xf0, 0x6e, 0xa2, 0xc6, 0x05, 0xd7, 0x7e, 0x3c, 0xcd, 0xbc, 0x18, 0xf1,
	0x40, 0x89, 0x2d, 0xc1, 0xa3, 0xb8, 0xc4, 0x1e, 0x4d, 0x83, 0x45, 0xc2, 0x1b, 0xab, 0x8a, 0x78,
	0x58, 0x23, 0x3c, 0xe9, 0xc9, 0x9e, 0x69, 0x0f, 0x92, 0x2d, 0x53, 0x99, 0xc2, 0x9c, 0xca, 0xf1,
	0x7b, 0x9e, 0x94, 0xde, 0x84, 0x5f, 0x4f, 0x71, 0x3f, 0x54, 0x1f, 0x4a, 0xb1, 0xf3, 0xab, 0xa8,
	0x84, 0xcf, 0x63, 0xe5, 0xfa, 0x61, 0x31, 0xb0, 0x74, 0x3a, 0x43, 0xff, 0x7b, 0x53, 0x0b, 0x89,
	0x3d, 0xa5, 0x0d, 0x31, 0x8a, 0x2d, 0xe8, 0xc2, 0xf2, 0xe2, 0x4a, 0xc7, 0x2e, 0xc2, 0xb2, 0xeb,
	0x23, 0x2f, 0xae, 0x77, 0x5c, 0x6f, 0x1d, 0x9f, 0x75, 0xc8, 0xc9, 0x59, 0x07, 0x9c, 0x9c, 0x62,
	0x8c, 0xce, 0x06, 0xae, 0xcf, 0xad, 0x99, 0x2e, 0x2c, 0x2f, 0x38, 0xe6, 0xcc, 0xba, 0x74, 0x71,
	0xc4, 0xe3, 0x61, 0x24, 0xc2, 0x1c, 0xb6, 0x1a, 0x46, 0xaa, 0xb7, 0xd8, 0x5d, 0xda, 0x48, 0xa2,
	0x89, 0x35, 0x9b, 0x2b, 0xeb, 0xf3, 0xfa, 0xac, 0xd3, 0xd8, 0x70, 0x5e, 0x3a, 0x79, 0x8f, 0xb5,
	0x69, 0x6b, 0x22, 0x87, 0xe6, 0xb5, 0xd6, 0x9c, 0x21, 0x7f, 0xd4, 0xec, 0x36, 0x9d, 0xe3, 0xbe,
	0x2b, 0x26, 0x56, 0xd3, 0x08, 0x45, 0xc1, 0x9e, 0x53, 0x3a, 0x8c, 0xb8, 0xab, 0xf8, 0x68, 0xd3,
	0x55, 0xd6, 0xbc, 0x59, 0xa3, 0x6d, 0x17, 0x31, 0xd8, 0x55, 0x0c, 0x76, 0xbf, 0x8a, 0xa1, 0xd8,
	0x60, 0xef, 0x5b, 0x07, 0x9c, 0x85, 0x92, 0x5b, 0x53, 0xb9, 0x49, 0x12, 0x8e, 0x2a, 0x93, 0xd6,
	0xdf, 0x98, 0x94, 0xdc, 0x9a, 0x5a, 0xfa, 0x0a, 0x94, 0xd5, 0x73, 0x7b, 0xc5, 0xfd, 0x01, 0x8f,
	0x58, 0x9f, 0xde, 0xaa, 0xdf, 0xca, 0xcd, 0x7f, 0x4a, 0xfb, 0xa6, 0xfc, 0x69, 0x24, 0x66, 0xcf,
	0x68, 0x2b, 0x89, 0x79, 0x64, 0xdc, 0x66, 0x8c, 0xdb, 0x9d, 0xca, 0x6d, 0x23, 0xe6, 0xd1, 0xef,
	0x5d, 0xe6, 0x13, 0x23, 0xc5, 0xec, 0x3e, 0x6d, 0x16, 0x57, 0xd0, 0x6a, 0x74, 0x1b, 0xcb, 0x37,
	0x56, 0xfe, 0xaf, 0x58, 0x27, 0xef, 0x3a, 0xa5, 0xb8, 0xfe, 0x19, 0x8e, 0x53, 0x84, 0x93, 0x14,
	0xe1, 0x34, 0x45, 0x38, 0x4f, 0x11, 0x2e, 0x52, 0x24, 0x97, 0x29, 0x92, 0xab, 0x14, 0x61, 0x47,
	0x23, 0xd9, 0xd5, 0x48, 0xf6, 0x35, 0xc2, 0x81, 0x46, 0x72, 0xa8, 0x11, 0x8e, 0x34, 0xc2, 0xb1,
	0x46, 0x38, 0xd1, 0x08, 0xa7, 0x1a, 0xc9, 0xb9, 0x46, 0xb8, 0xd0, 0x48, 0x2e, 0x35, 0xc2, 0x95,
	0x46, 0xb2, 0x93, 0x21, 0xd9, 0xcd, 0x10, 0xf6, 0x32, 0x24, 0x9f, 0x32, 0x84, 0x2f, 0x19, 0x92,
	0xfd, 0x0c, 0xc9, 0x41, 0x86, 0x70, 0x98, 0x21, 0x1c, 0x65, 0x08, 0x6f, 0x1f, 0xfc, 0xe9, 0x17,
	0x0a, 0xb7, 0xbd, 0xfc, 0x19, 0x0e, 0x06, 0x4d, 0xf3, 0x71, 0x56, 0xbf, 0x07, 0x00, 0x00, 0xff,
	0xff, 0x33, 0x55, 0xff, 0x78, 0x1d, 0x04, 0x00, 0x00,
}
