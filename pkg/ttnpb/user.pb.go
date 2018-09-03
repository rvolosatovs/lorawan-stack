// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: go.thethings.network/lorawan-stack/api/user.proto

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

// User is the message that defines an user on the network.
type User struct {
	// User identifiers.
	UserIdentifiers `protobuf:"bytes,1,opt,name=ids,embedded=ids" json:"ids"`
	// password is the user's password.
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// name is the user's full name.
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// validated_at denotes the time the email was validated at.
	// This is a read-only field.
	ValidatedAt *time.Time `protobuf:"bytes,4,opt,name=validated_at,json=validatedAt,stdtime" json:"validated_at,omitempty"`
	// admin denotes whether or not the user has administrative rights within the tenancy.
	// This field can only be modified by admins.
	Admin bool `protobuf:"varint,5,opt,name=admin,proto3" json:"admin,omitempty"`
	// state denotes the reviewing state of the user by admin.
	// This field can only be modified by admins.
	State ReviewingState `protobuf:"varint,6,opt,name=state,proto3,enum=ttn.lorawan.v3.ReviewingState" json:"state,omitempty"`
	// created_at denotes when the user was created.
	// This is a read-only field.
	CreatedAt time.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt,stdtime" json:"created_at"`
	// updated_at is the last time the user was updated.
	// This is a read-only field.
	UpdatedAt time.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,stdtime" json:"updated_at"`
	// The last time the password was updated (read-only field).
	PasswordUpdatedAt time.Time `protobuf:"bytes,9,opt,name=password_updated_at,json=passwordUpdatedAt,stdtime" json:"password_updated_at"`
	// Require the user to update its password (modifiable only by admins).
	RequirePasswordUpdate bool     `protobuf:"varint,10,opt,name=require_password_update,json=requirePasswordUpdate,proto3" json:"require_password_update,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *User) Reset()      { *m = User{} }
func (*User) ProtoMessage() {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_328bde6787001cf6, []int{0}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_User.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return m.Size()
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetValidatedAt() *time.Time {
	if m != nil {
		return m.ValidatedAt
	}
	return nil
}

func (m *User) GetAdmin() bool {
	if m != nil {
		return m.Admin
	}
	return false
}

func (m *User) GetState() ReviewingState {
	if m != nil {
		return m.State
	}
	return STATE_PENDING
}

func (m *User) GetCreatedAt() time.Time {
	if m != nil {
		return m.CreatedAt
	}
	return time.Time{}
}

func (m *User) GetUpdatedAt() time.Time {
	if m != nil {
		return m.UpdatedAt
	}
	return time.Time{}
}

func (m *User) GetPasswordUpdatedAt() time.Time {
	if m != nil {
		return m.PasswordUpdatedAt
	}
	return time.Time{}
}

func (m *User) GetRequirePasswordUpdate() bool {
	if m != nil {
		return m.RequirePasswordUpdate
	}
	return false
}

func init() {
	proto.RegisterType((*User)(nil), "ttn.lorawan.v3.User")
	golang_proto.RegisterType((*User)(nil), "ttn.lorawan.v3.User")
}
func (this *User) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*User)
	if !ok {
		that2, ok := that.(User)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *User")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *User but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *User but is not nil && this == nil")
	}
	if !this.UserIdentifiers.Equal(&that1.UserIdentifiers) {
		return fmt.Errorf("UserIdentifiers this(%v) Not Equal that(%v)", this.UserIdentifiers, that1.UserIdentifiers)
	}
	if this.Password != that1.Password {
		return fmt.Errorf("Password this(%v) Not Equal that(%v)", this.Password, that1.Password)
	}
	if this.Name != that1.Name {
		return fmt.Errorf("Name this(%v) Not Equal that(%v)", this.Name, that1.Name)
	}
	if that1.ValidatedAt == nil {
		if this.ValidatedAt != nil {
			return fmt.Errorf("this.ValidatedAt != nil && that1.ValidatedAt == nil")
		}
	} else if !this.ValidatedAt.Equal(*that1.ValidatedAt) {
		return fmt.Errorf("ValidatedAt this(%v) Not Equal that(%v)", this.ValidatedAt, that1.ValidatedAt)
	}
	if this.Admin != that1.Admin {
		return fmt.Errorf("Admin this(%v) Not Equal that(%v)", this.Admin, that1.Admin)
	}
	if this.State != that1.State {
		return fmt.Errorf("State this(%v) Not Equal that(%v)", this.State, that1.State)
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return fmt.Errorf("CreatedAt this(%v) Not Equal that(%v)", this.CreatedAt, that1.CreatedAt)
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return fmt.Errorf("UpdatedAt this(%v) Not Equal that(%v)", this.UpdatedAt, that1.UpdatedAt)
	}
	if !this.PasswordUpdatedAt.Equal(that1.PasswordUpdatedAt) {
		return fmt.Errorf("PasswordUpdatedAt this(%v) Not Equal that(%v)", this.PasswordUpdatedAt, that1.PasswordUpdatedAt)
	}
	if this.RequirePasswordUpdate != that1.RequirePasswordUpdate {
		return fmt.Errorf("RequirePasswordUpdate this(%v) Not Equal that(%v)", this.RequirePasswordUpdate, that1.RequirePasswordUpdate)
	}
	return nil
}
func (this *User) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*User)
	if !ok {
		that2, ok := that.(User)
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
	if !this.UserIdentifiers.Equal(&that1.UserIdentifiers) {
		return false
	}
	if this.Password != that1.Password {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if that1.ValidatedAt == nil {
		if this.ValidatedAt != nil {
			return false
		}
	} else if !this.ValidatedAt.Equal(*that1.ValidatedAt) {
		return false
	}
	if this.Admin != that1.Admin {
		return false
	}
	if this.State != that1.State {
		return false
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return false
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return false
	}
	if !this.PasswordUpdatedAt.Equal(that1.PasswordUpdatedAt) {
		return false
	}
	if this.RequirePasswordUpdate != that1.RequirePasswordUpdate {
		return false
	}
	return true
}
func (m *User) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *User) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintUser(dAtA, i, uint64(m.UserIdentifiers.Size()))
	n1, err := m.UserIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Password) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Password)))
		i += copy(dAtA[i:], m.Password)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.ValidatedAt != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.ValidatedAt)))
		n2, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.ValidatedAt, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.Admin {
		dAtA[i] = 0x28
		i++
		if m.Admin {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.State != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintUser(dAtA, i, uint64(m.State))
	}
	dAtA[i] = 0x3a
	i++
	i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)))
	n3, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	dAtA[i] = 0x42
	i++
	i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)))
	n4, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	dAtA[i] = 0x4a
	i++
	i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.PasswordUpdatedAt)))
	n5, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.PasswordUpdatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	if m.RequirePasswordUpdate {
		dAtA[i] = 0x50
		i++
		if m.RequirePasswordUpdate {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func encodeVarintUser(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedUser(r randyUser, easy bool) *User {
	this := &User{}
	v1 := NewPopulatedUserIdentifiers(r, easy)
	this.UserIdentifiers = *v1
	this.Password = randStringUser(r)
	this.Name = randStringUser(r)
	if r.Intn(10) != 0 {
		this.ValidatedAt = github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	}
	this.Admin = bool(r.Intn(2) == 0)
	this.State = ReviewingState([]int32{0, 1, 2}[r.Intn(3)])
	v2 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.CreatedAt = *v2
	v3 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.UpdatedAt = *v3
	v4 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.PasswordUpdatedAt = *v4
	this.RequirePasswordUpdate = bool(r.Intn(2) == 0)
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyUser interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneUser(r randyUser) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringUser(r randyUser) string {
	v5 := r.Intn(100)
	tmps := make([]rune, v5)
	for i := 0; i < v5; i++ {
		tmps[i] = randUTF8RuneUser(r)
	}
	return string(tmps)
}
func randUnrecognizedUser(r randyUser, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldUser(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldUser(dAtA []byte, r randyUser, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateUser(dAtA, uint64(key))
		v6 := r.Int63()
		if r.Intn(2) == 0 {
			v6 *= -1
		}
		dAtA = encodeVarintPopulateUser(dAtA, uint64(v6))
	case 1:
		dAtA = encodeVarintPopulateUser(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateUser(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateUser(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateUser(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateUser(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *User) Size() (n int) {
	var l int
	_ = l
	l = m.UserIdentifiers.Size()
	n += 1 + l + sovUser(uint64(l))
	l = len(m.Password)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.ValidatedAt != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.ValidatedAt)
		n += 1 + l + sovUser(uint64(l))
	}
	if m.Admin {
		n += 2
	}
	if m.State != 0 {
		n += 1 + sovUser(uint64(m.State))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovUser(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)
	n += 1 + l + sovUser(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.PasswordUpdatedAt)
	n += 1 + l + sovUser(uint64(l))
	if m.RequirePasswordUpdate {
		n += 2
	}
	return n
}

func sovUser(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozUser(x uint64) (n int) {
	return sovUser((x << 1) ^ uint64((int64(x) >> 63)))
}
func (this *User) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&User{`,
		`UserIdentifiers:` + strings.Replace(strings.Replace(this.UserIdentifiers.String(), "UserIdentifiers", "UserIdentifiers", 1), `&`, ``, 1) + `,`,
		`Password:` + fmt.Sprintf("%v", this.Password) + `,`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`ValidatedAt:` + strings.Replace(fmt.Sprintf("%v", this.ValidatedAt), "Timestamp", "types.Timestamp", 1) + `,`,
		`Admin:` + fmt.Sprintf("%v", this.Admin) + `,`,
		`State:` + fmt.Sprintf("%v", this.State) + `,`,
		`CreatedAt:` + strings.Replace(strings.Replace(this.CreatedAt.String(), "Timestamp", "types.Timestamp", 1), `&`, ``, 1) + `,`,
		`UpdatedAt:` + strings.Replace(strings.Replace(this.UpdatedAt.String(), "Timestamp", "types.Timestamp", 1), `&`, ``, 1) + `,`,
		`PasswordUpdatedAt:` + strings.Replace(strings.Replace(this.PasswordUpdatedAt.String(), "Timestamp", "types.Timestamp", 1), `&`, ``, 1) + `,`,
		`RequirePasswordUpdate:` + fmt.Sprintf("%v", this.RequirePasswordUpdate) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringUser(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *User) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUser
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
			return fmt.Errorf("proto: User: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: User: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.UserIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Password", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Password = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ValidatedAt == nil {
				m.ValidatedAt = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.ValidatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Admin", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Admin = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= (ReviewingState(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
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
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PasswordUpdatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.PasswordUpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequirePasswordUpdate", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.RequirePasswordUpdate = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipUser(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthUser
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
func skipUser(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
				return 0, ErrInvalidLengthUser
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowUser
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
				next, err := skipUser(dAtA[start:])
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
	ErrInvalidLengthUser = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUser   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("go.thethings.network/lorawan-stack/api/user.proto", fileDescriptor_user_328bde6787001cf6)
}
func init() {
	golang_proto.RegisterFile("go.thethings.network/lorawan-stack/api/user.proto", fileDescriptor_user_328bde6787001cf6)
}

var fileDescriptor_user_328bde6787001cf6 = []byte{
	// 543 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x3f, 0x4c, 0xdb, 0x40,
	0x14, 0xc6, 0xef, 0x95, 0x40, 0xc3, 0x51, 0x21, 0xd5, 0x6d, 0x55, 0x2b, 0x95, 0x5e, 0xa2, 0x4e,
	0x19, 0x8a, 0xad, 0x42, 0x55, 0x55, 0xea, 0x04, 0x4c, 0xdd, 0x2a, 0x17, 0x96, 0x2e, 0xe8, 0x12,
	0x1f, 0xce, 0x89, 0xf8, 0x4f, 0xed, 0x97, 0x44, 0xdd, 0x18, 0x19, 0x19, 0x3b, 0x74, 0xa8, 0x3a,
	0x31, 0x32, 0x32, 0x32, 0x66, 0xcc, 0xc8, 0x44, 0xf1, 0x79, 0x61, 0x64, 0x64, 0xac, 0xe2, 0xd8,
	0xa1, 0x64, 0x69, 0xd8, 0xee, 0xdd, 0xfb, 0xbe, 0xdf, 0x7d, 0xef, 0x49, 0xc7, 0xdf, 0x7a, 0xa1,
	0x45, 0x1d, 0x49, 0x1d, 0x15, 0x78, 0x89, 0x15, 0x48, 0x1a, 0x84, 0xf1, 0x81, 0xdd, 0x0d, 0x63,
	0x31, 0x10, 0xc1, 0x5a, 0x42, 0xa2, 0x7d, 0x60, 0x8b, 0x48, 0xd9, 0xbd, 0x44, 0xc6, 0x56, 0x14,
	0x87, 0x14, 0x1a, 0xab, 0x44, 0x81, 0x55, 0x28, 0xac, 0xfe, 0x46, 0x6d, 0x63, 0x4e, 0x44, 0xbb,
	0xab, 0x64, 0x40, 0x13, 0x48, 0xed, 0xc3, 0x9c, 0x26, 0xe5, 0xca, 0x80, 0xd4, 0xbe, 0x92, 0x71,
	0x52, 0x38, 0xe7, 0x7d, 0x2e, 0x56, 0x5e, 0x87, 0x4a, 0xd3, 0x9a, 0xa7, 0xa8, 0xd3, 0x6b, 0x59,
	0xed, 0xd0, 0xb7, 0xbd, 0xd0, 0x0b, 0xed, 0xfc, 0xba, 0xd5, 0xdb, 0xcf, 0xab, 0xbc, 0xc8, 0x4f,
	0x85, 0xfc, 0x95, 0x17, 0x86, 0x5e, 0x57, 0xde, 0xa9, 0xa4, 0x1f, 0xd1, 0xf7, 0xa2, 0x59, 0x9f,
	0x6d, 0x92, 0xf2, 0x65, 0x42, 0xc2, 0x8f, 0x26, 0x82, 0xd7, 0x3f, 0x2b, 0xbc, 0xb2, 0x9b, 0xc8,
	0xd8, 0xf8, 0xc8, 0x17, 0x94, 0x9b, 0x98, 0xd0, 0x80, 0xe6, 0xca, 0x7a, 0xdd, 0xba, 0xbf, 0x37,
	0x6b, 0x2c, 0xf9, 0x74, 0x37, 0xde, 0x56, 0x75, 0x78, 0x59, 0x67, 0xa3, 0xcb, 0x3a, 0x38, 0x63,
	0x97, 0x51, 0xe3, 0xd5, 0x48, 0x24, 0xc9, 0x20, 0x8c, 0x5d, 0xf3, 0x51, 0x03, 0x9a, 0xcb, 0xce,
	0xb4, 0x36, 0x0c, 0x5e, 0x09, 0x84, 0x2f, 0xcd, 0x85, 0xfc, 0x3e, 0x3f, 0x1b, 0xdb, 0xfc, 0x49,
	0x5f, 0x74, 0x95, 0x2b, 0x48, 0xba, 0x7b, 0x82, 0xcc, 0x4a, 0xfe, 0x6a, 0xcd, 0x9a, 0xa4, 0xb5,
	0xca, 0xb4, 0xd6, 0x4e, 0x99, 0x76, 0xab, 0x72, 0xfc, 0xa7, 0x0e, 0xce, 0xca, 0xd4, 0xb5, 0x49,
	0xc6, 0x73, 0xbe, 0x28, 0x5c, 0x5f, 0x05, 0xe6, 0x62, 0x03, 0x9a, 0x55, 0x67, 0x52, 0x18, 0xef,
	0xf8, 0x62, 0x42, 0x82, 0xa4, 0xb9, 0xd4, 0x80, 0xe6, 0xea, 0x3a, 0xce, 0x4e, 0xe2, 0xc8, 0xbe,
	0x92, 0x03, 0x15, 0x78, 0x5f, 0xc6, 0x2a, 0x67, 0x22, 0x36, 0xb6, 0x39, 0x6f, 0xc7, 0xb2, 0x8c,
	0xf3, 0xf8, 0xbf, 0x71, 0xf2, 0xf9, 0xf3, 0x48, 0xcb, 0x85, 0x6f, 0x93, 0xc6, 0x90, 0x5e, 0x34,
	0x9d, 0xa9, 0xfa, 0x10, 0x48, 0xe1, 0xdb, 0x24, 0x63, 0x87, 0x3f, 0x2b, 0x57, 0xb7, 0xf7, 0x0f,
	0x6d, 0xf9, 0x01, 0xb4, 0xa7, 0x25, 0x60, 0x77, 0x4a, 0x7d, 0xcf, 0x5f, 0xc6, 0xf2, 0x5b, 0x4f,
	0xc5, 0x72, 0x6f, 0x86, 0x6e, 0xf2, 0x7c, 0x7b, 0x2f, 0x8a, 0xf6, 0xe7, 0x7b, 0xd6, 0xad, 0xdf,
	0x30, 0x4c, 0x11, 0x46, 0x29, 0xc2, 0x45, 0x8a, 0xec, 0x2a, 0x45, 0xb8, 0x4e, 0x91, 0xdd, 0xa4,
	0xc8, 0x6e, 0x53, 0x84, 0x43, 0x8d, 0x70, 0xa4, 0x91, 0x9d, 0x68, 0x84, 0x53, 0x8d, 0xec, 0x4c,
	0x23, 0x3b, 0xd7, 0x08, 0x43, 0x8d, 0x30, 0xd2, 0x08, 0x17, 0x1a, 0xd9, 0x95, 0x46, 0xb8, 0xd6,
	0xc8, 0x6e, 0x34, 0xc2, 0xad, 0x46, 0x76, 0x98, 0x21, 0x3b, 0xca, 0x10, 0x8e, 0x33, 0x64, 0x3f,
	0x32, 0x84, 0x5f, 0x19, 0xb2, 0x93, 0x0c, 0xd9, 0x69, 0x86, 0x70, 0x96, 0x21, 0x9c, 0x67, 0x08,
	0x5f, 0xdf, 0xcc, 0xf1, 0x6f, 0xa2, 0x03, 0xcf, 0x26, 0x0a, 0xa2, 0x56, 0x6b, 0x29, 0xdf, 0xc6,
	0xc6, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcf, 0x79, 0xe5, 0x33, 0x20, 0x04, 0x00, 0x00,
}
