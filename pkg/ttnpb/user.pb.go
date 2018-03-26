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

import types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

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
	State ReviewingState `protobuf:"varint,6,opt,name=state,proto3,enum=ttn.v3.ReviewingState" json:"state,omitempty"`
	// created_at denotes when the user was created.
	// This is a read-only field.
	CreatedAt time.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt,stdtime" json:"created_at"`
	// updated_at is the last time the user was updated.
	// This is a read-only field.
	UpdatedAt time.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,stdtime" json:"updated_at"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptorUser, []int{0} }

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
		i = encodeVarintUser(dAtA, i, uint64(types.SizeOfStdTime(*m.ValidatedAt)))
		n2, err := types.StdTimeMarshalTo(*m.ValidatedAt, dAtA[i:])
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
	i = encodeVarintUser(dAtA, i, uint64(types.SizeOfStdTime(m.CreatedAt)))
	n3, err := types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	dAtA[i] = 0x42
	i++
	i = encodeVarintUser(dAtA, i, uint64(types.SizeOfStdTime(m.UpdatedAt)))
	n4, err := types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i:])
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
	v1 := NewPopulatedUserIdentifiers(r, easy)
	this.UserIdentifiers = *v1
	this.Password = randStringUser(r)
	this.Name = randStringUser(r)
	if r.Intn(10) != 0 {
		this.ValidatedAt = types.NewPopulatedStdTime(r, easy)
	}
	this.Admin = bool(r.Intn(2) == 0)
	this.State = ReviewingState([]int32{0, 1, 2}[r.Intn(3)])
	v2 := types.NewPopulatedStdTime(r, easy)
	this.CreatedAt = *v2
	v3 := types.NewPopulatedStdTime(r, easy)
	this.UpdatedAt = *v3
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
	v4 := r.Intn(100)
	tmps := make([]rune, v4)
	for i := 0; i < v4; i++ {
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
		v5 := r.Int63()
		if r.Intn(2) == 0 {
			v5 *= -1
		}
		dAtA = encodeVarintPopulateUser(dAtA, uint64(v5))
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
		l = types.SizeOfStdTime(*m.ValidatedAt)
		n += 1 + l + sovUser(uint64(l))
	}
	if m.Admin {
		n += 2
	}
	if m.State != 0 {
		n += 1 + sovUser(uint64(m.State))
	}
	l = types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovUser(uint64(l))
	l = types.SizeOfStdTime(m.UpdatedAt)
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
	return sovUser((x << 1) ^ uint64((int64(x) >> 63)))
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
			if err := types.StdTimeUnmarshal(m.ValidatedAt, dAtA[iNdEx:postIndex]); err != nil {
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
			if err := types.StdTimeUnmarshal(&m.UpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
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
	// 490 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x3f, 0x6f, 0xd3, 0x40,
	0x18, 0xc6, 0xef, 0x6d, 0x93, 0x90, 0x5e, 0x11, 0x83, 0x85, 0xc0, 0x0a, 0xd2, 0x9b, 0x88, 0x29,
	0x48, 0xe0, 0xa0, 0x46, 0x7c, 0x80, 0xa4, 0x13, 0x0b, 0x83, 0x09, 0x0b, 0x0b, 0x72, 0xe2, 0xab,
	0x73, 0x6a, 0xfc, 0x47, 0xf6, 0x9b, 0x44, 0x6c, 0x1d, 0x3b, 0x76, 0x64, 0x83, 0xb1, 0x63, 0xc7,
	0x8e, 0x1d, 0x33, 0x66, 0xec, 0x54, 0xea, 0xf3, 0xd2, 0xb1, 0x63, 0x27, 0x84, 0x72, 0xb6, 0x53,
	0xc4, 0xd2, 0x74, 0xf2, 0xbd, 0xf7, 0x3c, 0xbf, 0xbb, 0xe7, 0xb1, 0x8e, 0x5b, 0x9e, 0xa4, 0xf1,
	0x74, 0x68, 0x8d, 0x42, 0xbf, 0x33, 0x18, 0x8b, 0xc1, 0x58, 0x06, 0x5e, 0xf2, 0x49, 0xd0, 0x3c,
	0x8c, 0x0f, 0x3b, 0x44, 0x41, 0xc7, 0x89, 0x64, 0x67, 0x9a, 0x88, 0xd8, 0x8a, 0xe2, 0x90, 0x42,
	0xa3, 0x46, 0x14, 0x58, 0xb3, 0x6e, 0xe3, 0xfd, 0x26, 0xdc, 0x68, 0x22, 0x45, 0x40, 0x39, 0xd9,
	0xf8, 0xb0, 0x09, 0x21, 0x5d, 0x11, 0x90, 0x3c, 0x90, 0x22, 0x4e, 0x0a, 0x6c, 0xa3, 0x8b, 0x62,
	0xe9, 0x8d, 0xa9, 0x24, 0xde, 0xfd, 0x43, 0x78, 0xa1, 0x17, 0x76, 0xf4, 0xf6, 0x70, 0x7a, 0xa0,
	0x27, 0x3d, 0xe8, 0x55, 0x61, 0x7f, 0xe5, 0x85, 0xa1, 0x37, 0x11, 0xf7, 0x2e, 0xe1, 0x47, 0xf4,
	0xbd, 0x10, 0x9b, 0xff, 0x8b, 0x24, 0x7d, 0x91, 0x90, 0xe3, 0x47, 0xb9, 0xe1, 0xf5, 0x9f, 0x2d,
	0x5e, 0xf9, 0x92, 0x88, 0xd8, 0xe8, 0xf2, 0x6d, 0xe9, 0x26, 0x26, 0xb4, 0xa0, 0xbd, 0xbb, 0xf7,
	0xd2, 0xca, 0x7f, 0x93, 0xb5, 0x92, 0x3e, 0xde, 0x77, 0xea, 0xd7, 0x17, 0x57, 0x4d, 0xb6, 0xbc,
	0x6a, 0x82, 0xbd, 0x72, 0x1b, 0x0d, 0x5e, 0x8f, 0x9c, 0x24, 0x99, 0x87, 0xb1, 0x6b, 0x6e, 0xb5,
	0xa0, 0xbd, 0x63, 0xaf, 0x67, 0xc3, 0xe0, 0x95, 0xc0, 0xf1, 0x85, 0xb9, 0xad, 0xf7, 0xf5, 0xda,
	0xd8, 0xe7, 0x4f, 0x67, 0xce, 0x44, 0xba, 0x0e, 0x09, 0xf7, 0x9b, 0x43, 0x66, 0x45, 0xdf, 0xd6,
	0xb0, 0xf2, 0x94, 0x56, 0x99, 0xd2, 0x1a, 0x94, 0x29, 0xfb, 0x95, 0x93, 0xdf, 0x4d, 0xb0, 0x77,
	0xd7, 0x54, 0x8f, 0x8c, 0xe7, 0xbc, 0xea, 0xb8, 0xbe, 0x0c, 0xcc, 0x6a, 0x0b, 0xda, 0x75, 0x3b,
	0x1f, 0x8c, 0xb7, 0xbc, 0x9a, 0x90, 0x43, 0xc2, 0xac, 0xb5, 0xa0, 0xfd, 0x6c, 0xef, 0x45, 0xd9,
	0xc0, 0x16, 0x33, 0x29, 0xe6, 0x32, 0xf0, 0x3e, 0xaf, 0x54, 0x3b, 0x37, 0x19, 0xfb, 0x9c, 0x8f,
	0x62, 0x51, 0xc6, 0x78, 0xf2, 0x60, 0x0c, 0xdd, 0x5b, 0x47, 0xd9, 0x29, 0xb8, 0x1e, 0xad, 0x0e,
	0x99, 0x46, 0xeb, 0x2e, 0xf5, 0xc7, 0x1c, 0x52, 0x70, 0x3d, 0xea, 0xff, 0x84, 0x45, 0x8a, 0xb0,
	0x4c, 0x11, 0x2e, 0x53, 0x84, 0xeb, 0x14, 0xe1, 0x26, 0x45, 0x76, 0x9b, 0x22, 0xbb, 0x4b, 0x11,
	0x8e, 0x14, 0xb2, 0x63, 0x85, 0xec, 0x54, 0x21, 0x9c, 0x29, 0x64, 0xe7, 0x0a, 0xe1, 0x42, 0x21,
	0x2c, 0x14, 0xc2, 0x52, 0x21, 0x5c, 0x2a, 0x64, 0xd7, 0x0a, 0xe1, 0x46, 0x21, 0xbb, 0x55, 0x08,
	0x77, 0x0a, 0xd9, 0x51, 0x86, 0xec, 0x38, 0x43, 0x38, 0xc9, 0x90, 0xfd, 0xc8, 0x10, 0x7e, 0x65,
	0xc8, 0x4e, 0x33, 0x64, 0x67, 0x19, 0xc2, 0x79, 0x86, 0x70, 0x91, 0x21, 0x7c, 0x7d, 0xf3, 0xd0,
	0xb3, 0x8c, 0x0e, 0xbd, 0xd5, 0x37, 0x1a, 0x0e, 0x6b, 0xba, 0x4a, 0xf7, 0x6f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xcf, 0x28, 0xa8, 0xd9, 0x6b, 0x03, 0x00, 0x00,
}
