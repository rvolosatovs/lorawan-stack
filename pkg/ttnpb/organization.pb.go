// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: go.thethings.network/lorawan-stack/api/organization.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"

import time "time"

import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import strings "strings"
import reflect "reflect"

import io "io"

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
	UpdatedAt            time.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,stdtime" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Organization) Reset()      { *m = Organization{} }
func (*Organization) ProtoMessage() {}
func (*Organization) Descriptor() ([]byte, []int) {
	return fileDescriptor_organization_e25c0bf788b04f35, []int{0}
}
func (m *Organization) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Organization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Organization.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Organization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Organization.Merge(dst, src)
}
func (m *Organization) XXX_Size() int {
	return m.Size()
}
func (m *Organization) XXX_DiscardUnknown() {
	xxx_messageInfo_Organization.DiscardUnknown(m)
}

var xxx_messageInfo_Organization proto.InternalMessageInfo

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
	Rights               []Right  `protobuf:"varint,3,rep,packed,name=rights,enum=ttn.lorawan.v3.Right" json:"rights,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrganizationMember) Reset()      { *m = OrganizationMember{} }
func (*OrganizationMember) ProtoMessage() {}
func (*OrganizationMember) Descriptor() ([]byte, []int) {
	return fileDescriptor_organization_e25c0bf788b04f35, []int{1}
}
func (m *OrganizationMember) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrganizationMember) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrganizationMember.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *OrganizationMember) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrganizationMember.Merge(dst, src)
}
func (m *OrganizationMember) XXX_Size() int {
	return m.Size()
}
func (m *OrganizationMember) XXX_DiscardUnknown() {
	xxx_messageInfo_OrganizationMember.DiscardUnknown(m)
}

var xxx_messageInfo_OrganizationMember proto.InternalMessageInfo

func (m *OrganizationMember) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

func init() {
	proto.RegisterType((*Organization)(nil), "ttn.lorawan.v3.Organization")
	golang_proto.RegisterType((*Organization)(nil), "ttn.lorawan.v3.Organization")
	proto.RegisterType((*OrganizationMember)(nil), "ttn.lorawan.v3.OrganizationMember")
	golang_proto.RegisterType((*OrganizationMember)(nil), "ttn.lorawan.v3.OrganizationMember")
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
	i = encodeVarintOrganization(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)))
	n2, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	dAtA[i] = 0x42
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)))
	n3, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i:])
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
	v2 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.CreatedAt = *v2
	v3 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
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
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovOrganization(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)
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
func (this *Organization) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Organization{`,
		`OrganizationIdentifiers:` + strings.Replace(strings.Replace(this.OrganizationIdentifiers.String(), "OrganizationIdentifiers", "OrganizationIdentifiers", 1), `&`, ``, 1) + `,`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Description:` + fmt.Sprintf("%v", this.Description) + `,`,
		`URL:` + fmt.Sprintf("%v", this.URL) + `,`,
		`Location:` + fmt.Sprintf("%v", this.Location) + `,`,
		`Email:` + fmt.Sprintf("%v", this.Email) + `,`,
		`CreatedAt:` + strings.Replace(strings.Replace(this.CreatedAt.String(), "Timestamp", "types.Timestamp", 1), `&`, ``, 1) + `,`,
		`UpdatedAt:` + strings.Replace(strings.Replace(this.UpdatedAt.String(), "Timestamp", "types.Timestamp", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrganizationMember) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrganizationMember{`,
		`OrganizationIdentifiers:` + strings.Replace(strings.Replace(this.OrganizationIdentifiers.String(), "OrganizationIdentifiers", "OrganizationIdentifiers", 1), `&`, ``, 1) + `,`,
		`UserIdentifiers:` + strings.Replace(strings.Replace(this.UserIdentifiers.String(), "UserIdentifiers", "UserIdentifiers", 1), `&`, ``, 1) + `,`,
		`Rights:` + fmt.Sprintf("%v", this.Rights) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringOrganization(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
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
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.CreatedAt, dAtA[iNdEx:postIndex]); err != nil {
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
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
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
	proto.RegisterFile("go.thethings.network/lorawan-stack/api/organization.proto", fileDescriptor_organization_e25c0bf788b04f35)
}
func init() {
	golang_proto.RegisterFile("go.thethings.network/lorawan-stack/api/organization.proto", fileDescriptor_organization_e25c0bf788b04f35)
}

var fileDescriptor_organization_e25c0bf788b04f35 = []byte{
	// 566 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0x3f, 0x4c, 0xdb, 0x4c,
	0x18, 0xc6, 0xef, 0x08, 0x7f, 0xc2, 0xf1, 0x89, 0xaf, 0x3a, 0xb5, 0x92, 0x9b, 0x4a, 0xaf, 0x23,
	0x96, 0x32, 0x14, 0x5b, 0x85, 0xa5, 0x1d, 0x81, 0x2e, 0x95, 0x5a, 0x55, 0xb2, 0xca, 0x52, 0x55,
	0x42, 0x97, 0xe4, 0x70, 0x4e, 0xc4, 0x3e, 0xeb, 0x7c, 0x2e, 0x6a, 0x27, 0x46, 0x46, 0xc6, 0x8e,
	0x55, 0x27, 0x46, 0x46, 0x46, 0x46, 0x46, 0x46, 0x26, 0x8a, 0xcf, 0x43, 0x19, 0x19, 0x19, 0x2b,
	0x9f, 0x9d, 0x12, 0xa2, 0x0e, 0xa9, 0xba, 0xdd, 0x7b, 0xcf, 0xf3, 0xfc, 0xa2, 0xf7, 0xc9, 0x99,
	0xbc, 0x0c, 0xa5, 0xa7, 0xfb, 0x5c, 0xf7, 0x45, 0x1c, 0xa6, 0x5e, 0xcc, 0xf5, 0x9e, 0x54, 0xbb,
	0xfe, 0x40, 0x2a, 0xb6, 0xc7, 0xe2, 0x95, 0x54, 0xb3, 0xee, 0xae, 0xcf, 0x12, 0xe1, 0x4b, 0x15,
	0xb2, 0x58, 0x7c, 0x61, 0x5a, 0xc8, 0xd8, 0x4b, 0x94, 0xd4, 0x92, 0x2e, 0x6a, 0x1d, 0x7b, 0xb5,
	0xd3, 0xfb, 0xb4, 0xd6, 0x7a, 0x3e, 0x21, 0x8a, 0x65, 0xba, 0x5f, 0x21, 0x5a, 0x2f, 0x26, 0x8c,
	0x88, 0x1e, 0x8f, 0xb5, 0xd8, 0x11, 0x5c, 0xa5, 0x75, 0x72, 0x6d, 0xc2, 0xa4, 0x12, 0x61, 0x5f,
	0x0f, 0x43, 0x2b, 0xa1, 0xd0, 0xfd, 0xac, 0xe3, 0x75, 0x65, 0xe4, 0x87, 0x32, 0x94, 0xbe, 0xbd,
	0xee, 0x64, 0x3b, 0x76, 0xb2, 0x83, 0x3d, 0xd5, 0xf6, 0x27, 0xa1, 0x94, 0xe1, 0x80, 0xdf, 0xb9,
	0x78, 0x94, 0xe8, 0xcf, 0xb5, 0xe8, 0x8e, 0x8b, 0x5a, 0x44, 0x3c, 0xd5, 0x2c, 0x4a, 0x2a, 0xc3,
	0xd2, 0xd5, 0x14, 0xf9, 0xef, 0xdd, 0x48, 0x6b, 0x74, 0x93, 0x34, 0x44, 0x2f, 0x75, 0x70, 0x1b,
	0x2f, 0x2f, 0xac, 0x3e, 0xf5, 0xee, 0xb7, 0xe7, 0x8d, 0x5a, 0x5f, 0xdf, 0xad, 0xbb, 0xd1, 0x3c,
	0xbb, 0x74, 0xd1, 0xf9, 0xa5, 0x8b, 0x83, 0x32, 0x4d, 0x29, 0x99, 0x8e, 0x59, 0xc4, 0x9d, 0xa9,
	0x36, 0x5e, 0x9e, 0x0f, 0xec, 0x99, 0xb6, 0xc9, 0x42, 0x8f, 0xa7, 0x5d, 0x25, 0x92, 0x32, 0xec,
	0x34, 0xac, 0x34, 0x7a, 0x45, 0x1f, 0x93, 0x46, 0xa6, 0x06, 0xce, 0x74, 0xa9, 0x6c, 0xcc, 0x99,
	0x4b, 0xb7, 0xb1, 0x15, 0xbc, 0x09, 0xca, 0x3b, 0xda, 0x22, 0xcd, 0x81, 0xec, 0xda, 0x9f, 0x75,
	0x66, 0x6c, 0xf2, 0xf7, 0x4c, 0x1f, 0x92, 0x19, 0x1e, 0x31, 0x31, 0x70, 0x66, 0xad, 0x50, 0x0d,
	0x74, 0x93, 0x90, 0xae, 0xe2, 0x4c, 0xf3, 0xde, 0x36, 0xd3, 0xce, 0x9c, 0x5d, 0xa7, 0xe5, 0x55,
	0x75, 0x78, 0xc3, 0x3a, 0xbc, 0xf7, 0xc3, 0x3a, 0xaa, 0x0d, 0x0e, 0x7f, 0xb8, 0x38, 0x98, 0xaf,
	0x73, 0xeb, 0xba, 0x84, 0x64, 0x49, 0x6f, 0x08, 0x69, 0xfe, 0x0d, 0xa4, 0xce, 0xad, 0xeb, 0xa5,
	0x9f, 0x98, 0xd0, 0xd1, 0xde, 0xde, 0xf2, 0xa8, 0xc3, 0x15, 0xfd, 0x48, 0x1e, 0x8c, 0x3e, 0xd7,
	0xed, 0x7f, 0x6a, 0xfd, 0x7f, 0x79, 0xcf, 0x92, 0xd2, 0x57, 0xa4, 0x99, 0xa5, 0x5c, 0x59, 0xea,
	0x94, 0xa5, 0xba, 0xe3, 0xd4, 0xad, 0x94, 0xab, 0x3f, 0xd3, 0xe6, 0x32, 0x2b, 0xa5, 0x74, 0x85,
	0xcc, 0x56, 0x4f, 0xd3, 0x69, 0xb4, 0x1b, 0xcb, 0x8b, 0xab, 0x8f, 0xc6, 0x19, 0x41, 0xa9, 0x06,
	0xb5, 0x69, 0xe3, 0x3b, 0x3e, 0xcb, 0x01, 0x9f, 0xe7, 0x80, 0x2f, 0x72, 0x40, 0x57, 0x39, 0xe0,
	0xeb, 0x1c, 0xd0, 0x4d, 0x0e, 0xe8, 0x36, 0x07, 0xbc, 0x6f, 0x00, 0x1f, 0x18, 0x40, 0x47, 0x06,
	0xf0, 0xb1, 0x01, 0x74, 0x62, 0x00, 0x9d, 0x1a, 0xc0, 0x67, 0x06, 0xf0, 0xb9, 0x01, 0x7c, 0x61,
	0x00, 0x5d, 0x19, 0xc0, 0xd7, 0x06, 0xd0, 0x8d, 0x01, 0x7c, 0x6b, 0x00, 0xed, 0x17, 0x80, 0x0e,
	0x0a, 0xc0, 0x87, 0x05, 0xa0, 0xaf, 0x05, 0xe0, 0x6f, 0x05, 0xa0, 0xa3, 0x02, 0xd0, 0x71, 0x01,
	0xf8, 0xa4, 0x00, 0x7c, 0x5a, 0x00, 0xfe, 0xf0, 0x6c, 0x82, 0xaf, 0x2c, 0xd9, 0x0d, 0x7d, 0xad,
	0xe3, 0xa4, 0xd3, 0x99, 0xb5, 0xff, 0xdb, 0xda, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc0, 0xee,
	0x0c, 0x32, 0x54, 0x04, 0x00, 0x00,
}
