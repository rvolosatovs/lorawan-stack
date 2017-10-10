// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/user.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import _ "github.com/gogo/protobuf/types"

import time "time"

import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// User is the message that defines an user on the network.
type User struct {
	// id is the user's unique and immutable ID.
	UserIdentifier `protobuf:"bytes,1,opt,name=id,embedded=id" json:"id"`
	// email is the user's email address.
	Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	// password is the user's password.
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	// name is the user's full name.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// validated denotes if the email address has been validated.
	Validated bool `protobuf:"varint,5,opt,name=validated,proto3" json:"validated,omitempty"`
	// admin denotes whether or not the user has administrative rights within the
	// tenancy.
	Admin bool `protobuf:"varint,6,opt,name=admin,proto3" json:"admin,omitempty"`
	// created_at denotes when the user was created.
	CreatedAt time.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt,stdtime" json:"created_at"`
	// updated_at is the last time the user was updated.
	UpdatedAt time.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,stdtime" json:"updated_at"`
	// archived_at is the time when the user was archived and therefore
	// permantly disabled.
	ArchivedAt time.Time `protobuf:"bytes,9,opt,name=archived_at,json=archivedAt,stdtime" json:"archived_at"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptorUser, []int{0} }

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

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

func (m *User) GetValidated() bool {
	if m != nil {
		return m.Validated
	}
	return false
}

func (m *User) GetAdmin() bool {
	if m != nil {
		return m.Admin
	}
	return false
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

func (m *User) GetArchivedAt() time.Time {
	if m != nil {
		return m.ArchivedAt
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*User)(nil), "ttn.v3.User")
	golang_proto.RegisterType((*User)(nil), "ttn.v3.User")
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
	if !this.UserIdentifier.Equal(&that1.UserIdentifier) {
		return fmt.Errorf("UserIdentifier this(%v) Not Equal that(%v)", this.UserIdentifier, that1.UserIdentifier)
	}
	if this.Email != that1.Email {
		return fmt.Errorf("Email this(%v) Not Equal that(%v)", this.Email, that1.Email)
	}
	if this.Password != that1.Password {
		return fmt.Errorf("Password this(%v) Not Equal that(%v)", this.Password, that1.Password)
	}
	if this.Name != that1.Name {
		return fmt.Errorf("Name this(%v) Not Equal that(%v)", this.Name, that1.Name)
	}
	if this.Validated != that1.Validated {
		return fmt.Errorf("Validated this(%v) Not Equal that(%v)", this.Validated, that1.Validated)
	}
	if this.Admin != that1.Admin {
		return fmt.Errorf("Admin this(%v) Not Equal that(%v)", this.Admin, that1.Admin)
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return fmt.Errorf("CreatedAt this(%v) Not Equal that(%v)", this.CreatedAt, that1.CreatedAt)
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return fmt.Errorf("UpdatedAt this(%v) Not Equal that(%v)", this.UpdatedAt, that1.UpdatedAt)
	}
	if !this.ArchivedAt.Equal(that1.ArchivedAt) {
		return fmt.Errorf("ArchivedAt this(%v) Not Equal that(%v)", this.ArchivedAt, that1.ArchivedAt)
	}
	return nil
}
func (this *User) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
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
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if !this.UserIdentifier.Equal(&that1.UserIdentifier) {
		return false
	}
	if this.Email != that1.Email {
		return false
	}
	if this.Password != that1.Password {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Validated != that1.Validated {
		return false
	}
	if this.Admin != that1.Admin {
		return false
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return false
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return false
	}
	if !this.ArchivedAt.Equal(that1.ArchivedAt) {
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
	i = encodeVarintUser(dAtA, i, uint64(m.UserIdentifier.Size()))
	n1, err := m.UserIdentifier.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Email) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Email)))
		i += copy(dAtA[i:], m.Email)
	}
	if len(m.Password) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Password)))
		i += copy(dAtA[i:], m.Password)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintUser(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.Validated {
		dAtA[i] = 0x28
		i++
		if m.Validated {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.Admin {
		dAtA[i] = 0x30
		i++
		if m.Admin {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	dAtA[i] = 0x3a
	i++
	i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)))
	n2, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	dAtA[i] = 0x42
	i++
	i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)))
	n3, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	dAtA[i] = 0x4a
	i++
	i = encodeVarintUser(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.ArchivedAt)))
	n4, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.ArchivedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
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
	v1 := NewPopulatedUserIdentifier(r, easy)
	this.UserIdentifier = *v1
	this.Email = string(randStringUser(r))
	this.Password = string(randStringUser(r))
	this.Name = string(randStringUser(r))
	this.Validated = bool(bool(r.Intn(2) == 0))
	this.Admin = bool(bool(r.Intn(2) == 0))
	v2 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.CreatedAt = *v2
	v3 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.UpdatedAt = *v3
	v4 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.ArchivedAt = *v4
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
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *User) Size() (n int) {
	var l int
	_ = l
	l = m.UserIdentifier.Size()
	n += 1 + l + sovUser(uint64(l))
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Password)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.Validated {
		n += 2
	}
	if m.Admin {
		n += 2
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovUser(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)
	n += 1 + l + sovUser(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.ArchivedAt)
	n += 1 + l + sovUser(uint64(l))
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
	return sovUser(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
				return fmt.Errorf("proto: wrong wireType = %d for field UserIdentifier", wireType)
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
			if err := m.UserIdentifier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
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
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
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
		case 4:
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
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validated", wireType)
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
			m.Validated = bool(v != 0)
		case 6:
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
				return fmt.Errorf("proto: wrong wireType = %d for field ArchivedAt", wireType)
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
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.ArchivedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
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

func init() { proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/user.proto", fileDescriptorUser) }
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/user.proto", fileDescriptorUser)
}

var fileDescriptorUser = []byte{
	// 465 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xbe, 0x97, 0xa6, 0x21, 0xb9, 0x6e, 0x16, 0x42, 0x56, 0x40, 0x2f, 0x11, 0x53, 0x18, 0x38,
	0x23, 0x2a, 0x7e, 0x40, 0x8a, 0x18, 0x58, 0x18, 0xa2, 0xb0, 0xb0, 0xa0, 0x4b, 0x7c, 0x75, 0x4e,
	0x8d, 0x7d, 0x96, 0x7d, 0x49, 0xc5, 0x56, 0x89, 0xa5, 0x63, 0x47, 0x36, 0x18, 0x3b, 0x76, 0xec,
	0xd8, 0x31, 0x63, 0xc6, 0x4e, 0xa5, 0xbe, 0x5b, 0x3a, 0x76, 0xec, 0x88, 0x7c, 0x4e, 0x5a, 0xc4,
	0x42, 0x33, 0xf9, 0x7d, 0xef, 0x7b, 0xdf, 0xf7, 0xde, 0x67, 0x1d, 0x65, 0x91, 0xd4, 0x93, 0xd9,
	0x88, 0x8d, 0x55, 0x1c, 0x0c, 0x27, 0x62, 0x38, 0x91, 0x49, 0x94, 0x7f, 0x12, 0xfa, 0x50, 0x65,
	0x07, 0x81, 0xd6, 0x49, 0xc0, 0x53, 0x19, 0xcc, 0x72, 0x91, 0xb1, 0x34, 0x53, 0x5a, 0x79, 0x0d,
	0xad, 0x13, 0x36, 0xdf, 0x6d, 0xbf, 0xfe, 0x4b, 0x17, 0xa9, 0x48, 0x05, 0x8e, 0x1e, 0xcd, 0xf6,
	0x1d, 0x72, 0xc0, 0x55, 0x95, 0xac, 0xfd, 0xee, 0x31, 0x6b, 0x64, 0x28, 0x12, 0x2d, 0xf7, 0xa5,
	0xc8, 0xf2, 0x95, 0xac, 0x13, 0x29, 0x15, 0x4d, 0xc5, 0x83, 0xb9, 0x96, 0xb1, 0xc8, 0x35, 0x8f,
	0xd3, 0xd5, 0xc0, 0xf3, 0x7f, 0x07, 0x44, 0x9c, 0xea, 0x6f, 0x15, 0xf9, 0xf2, 0xfb, 0x16, 0xad,
	0x7f, 0xce, 0x45, 0xe6, 0xbd, 0xa1, 0x35, 0x19, 0xfa, 0xd0, 0x85, 0xde, 0xce, 0xdb, 0x67, 0xac,
	0x4a, 0xc0, 0x4a, 0xe6, 0xe3, 0xfd, 0xc6, 0xbd, 0xe6, 0xe2, 0xaa, 0x43, 0x96, 0x57, 0x1d, 0x18,
	0xd4, 0x64, 0xe8, 0x3d, 0xa5, 0xdb, 0x22, 0xe6, 0x72, 0xea, 0xd7, 0xba, 0xd0, 0x6b, 0x0d, 0x2a,
	0xe0, 0xb5, 0x69, 0x33, 0xe5, 0x79, 0x7e, 0xa8, 0xb2, 0xd0, 0xdf, 0x72, 0xc4, 0x3d, 0xf6, 0x3c,
	0x5a, 0x4f, 0x78, 0x2c, 0xfc, 0xba, 0xeb, 0xbb, 0xda, 0x7b, 0x41, 0x5b, 0x73, 0x3e, 0x95, 0x21,
	0xd7, 0x22, 0xf4, 0xb7, 0xbb, 0xd0, 0x6b, 0x0e, 0x1e, 0x1a, 0xe5, 0x0e, 0x1e, 0xc6, 0x32, 0xf1,
	0x1b, 0x8e, 0xa9, 0x80, 0xf7, 0x9e, 0xd2, 0x71, 0x26, 0xca, 0x81, 0xaf, 0x5c, 0xfb, 0x4f, 0xdc,
	0xcd, 0x6d, 0x56, 0xc5, 0x64, 0xeb, 0x98, 0x6c, 0xb8, 0xfe, 0x0f, 0xd5, 0xdd, 0x27, 0xbf, 0x3b,
	0x30, 0x68, 0xad, 0x74, 0x7d, 0x5d, 0x9a, 0xcc, 0xd2, 0x70, 0x6d, 0xd2, 0xdc, 0xc4, 0x64, 0xa5,
	0xeb, 0x6b, 0xef, 0x03, 0xdd, 0xe1, 0xd9, 0x78, 0x22, 0xe7, 0x95, 0x4b, 0x6b, 0x03, 0x17, 0xba,
	0x16, 0xf6, 0xf5, 0xde, 0x4f, 0x58, 0x14, 0x08, 0xcb, 0x02, 0xe1, 0xb2, 0x40, 0xb8, 0x2e, 0x10,
	0x6e, 0x0a, 0x24, 0xb7, 0x05, 0x92, 0xbb, 0x02, 0xe1, 0xc8, 0x20, 0x39, 0x36, 0x48, 0x4e, 0x0d,
	0xc2, 0x99, 0x41, 0x72, 0x6e, 0x10, 0x2e, 0x0c, 0xc2, 0xc2, 0x20, 0x2c, 0x0d, 0xc2, 0xa5, 0x41,
	0x72, 0x6d, 0x10, 0x6e, 0x0c, 0x92, 0x5b, 0x83, 0x70, 0x67, 0x90, 0x1c, 0x59, 0x24, 0xc7, 0x16,
	0xe1, 0xc4, 0x22, 0xf9, 0x61, 0x11, 0x7e, 0x59, 0x24, 0xa7, 0x16, 0xc9, 0x99, 0x45, 0x38, 0xb7,
	0x08, 0x17, 0x16, 0xe1, 0xcb, 0xab, 0xff, 0xbd, 0xb8, 0xf4, 0x20, 0x2a, 0xbf, 0xe9, 0x68, 0xd4,
	0x70, 0x59, 0x76, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x46, 0x0a, 0xe1, 0x34, 0x0c, 0x03, 0x00,
	0x00,
}
