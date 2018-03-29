// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/client.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"

import time "time"

import strconv "strconv"

import types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// ReviewingState enum defines all the possible admin reviewing states that a
// request can be at. For example a third-party client creation request.
type ReviewingState int32

const (
	// Denotes that the request is pending to review.
	STATE_PENDING ReviewingState = 0
	// Denotes that the request has been reviewed and approved by an admin.
	STATE_APPROVED ReviewingState = 1
	// Denotes that the request has been reviewed and rejected by an admin.
	STATE_REJECTED ReviewingState = 2
)

var ReviewingState_name = map[int32]string{
	0: "STATE_PENDING",
	1: "STATE_APPROVED",
	2: "STATE_REJECTED",
}
var ReviewingState_value = map[string]int32{
	"STATE_PENDING":  0,
	"STATE_APPROVED": 1,
	"STATE_REJECTED": 2,
}

func (ReviewingState) EnumDescriptor() ([]byte, []int) { return fileDescriptorClient, []int{0} }

// Client is the message that defines a third-party client on the network.
type Client struct {
	// Client identifiers.
	ClientIdentifiers `protobuf:"bytes,1,opt,name=ids,embedded=ids" json:""`
	// description is the description of the client.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// secret is the secret used to prove the client identity.
	// This is a read-only field.
	Secret string `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
	// redirect_uri is the OAuth redirect URI of the client.
	RedirectURI string `protobuf:"bytes,4,opt,name=redirect_uri,json=redirectUri,proto3" json:"redirect_uri,omitempty"`
	// state denotes the reviewing state of the client by admin.
	// This field can only be modified by admins.
	State ReviewingState `protobuf:"varint,5,opt,name=state,proto3,enum=ttn.v3.ReviewingState" json:"state,omitempty"`
	// official_labeled denotes if a client has been labeled as an official
	// third-party client by the tenant admin.
	// This field can only be modified by admins.
	OfficialLabeled bool `protobuf:"varint,6,opt,name=official_labeled,json=officialLabeled,proto3" json:"official_labeled,omitempty"`
	// grants denotes which OAuth2 flows can the client use to get a token.
	// This field can only be modified by admins.
	Grants []GrantType `protobuf:"varint,7,rep,packed,name=grants,enum=ttn.v3.GrantType" json:"grants,omitempty"`
	// Rights denotes what rights the client will have access to.
	Rights []Right `protobuf:"varint,8,rep,packed,name=rights,enum=ttn.v3.Right" json:"rights,omitempty"`
	// creator_ids are the identifiers of the user that created the client.
	CreatorIDs UserIdentifiers `protobuf:"bytes,9,opt,name=creator_ids,json=creatorIds" json:"creator_ids"`
	// created_at denotes when the client was created.
	// This a read-only field.
	CreatedAt time.Time `protobuf:"bytes,10,opt,name=created_at,json=createdAt,stdtime" json:"created_at"`
	// updated_at is the last time the client was updated.
	// This is a read-only field.
	UpdatedAt time.Time `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt,stdtime" json:"updated_at"`
}

func (m *Client) Reset()                    { *m = Client{} }
func (m *Client) String() string            { return proto.CompactTextString(m) }
func (*Client) ProtoMessage()               {}
func (*Client) Descriptor() ([]byte, []int) { return fileDescriptorClient, []int{0} }

func (m *Client) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Client) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *Client) GetRedirectURI() string {
	if m != nil {
		return m.RedirectURI
	}
	return ""
}

func (m *Client) GetState() ReviewingState {
	if m != nil {
		return m.State
	}
	return STATE_PENDING
}

func (m *Client) GetOfficialLabeled() bool {
	if m != nil {
		return m.OfficialLabeled
	}
	return false
}

func (m *Client) GetGrants() []GrantType {
	if m != nil {
		return m.Grants
	}
	return nil
}

func (m *Client) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

func (m *Client) GetCreatorIDs() UserIdentifiers {
	if m != nil {
		return m.CreatorIDs
	}
	return UserIdentifiers{}
}

func (m *Client) GetCreatedAt() time.Time {
	if m != nil {
		return m.CreatedAt
	}
	return time.Time{}
}

func (m *Client) GetUpdatedAt() time.Time {
	if m != nil {
		return m.UpdatedAt
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*Client)(nil), "ttn.v3.Client")
	golang_proto.RegisterType((*Client)(nil), "ttn.v3.Client")
	proto.RegisterEnum("ttn.v3.ReviewingState", ReviewingState_name, ReviewingState_value)
	golang_proto.RegisterEnum("ttn.v3.ReviewingState", ReviewingState_name, ReviewingState_value)
}
func (x ReviewingState) String() string {
	s, ok := ReviewingState_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *Client) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*Client)
	if !ok {
		that2, ok := that.(Client)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *Client")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *Client but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *Client but is not nil && this == nil")
	}
	if !this.ClientIdentifiers.Equal(&that1.ClientIdentifiers) {
		return fmt.Errorf("ClientIdentifiers this(%v) Not Equal that(%v)", this.ClientIdentifiers, that1.ClientIdentifiers)
	}
	if this.Description != that1.Description {
		return fmt.Errorf("Description this(%v) Not Equal that(%v)", this.Description, that1.Description)
	}
	if this.Secret != that1.Secret {
		return fmt.Errorf("Secret this(%v) Not Equal that(%v)", this.Secret, that1.Secret)
	}
	if this.RedirectURI != that1.RedirectURI {
		return fmt.Errorf("RedirectURI this(%v) Not Equal that(%v)", this.RedirectURI, that1.RedirectURI)
	}
	if this.State != that1.State {
		return fmt.Errorf("State this(%v) Not Equal that(%v)", this.State, that1.State)
	}
	if this.OfficialLabeled != that1.OfficialLabeled {
		return fmt.Errorf("OfficialLabeled this(%v) Not Equal that(%v)", this.OfficialLabeled, that1.OfficialLabeled)
	}
	if len(this.Grants) != len(that1.Grants) {
		return fmt.Errorf("Grants this(%v) Not Equal that(%v)", len(this.Grants), len(that1.Grants))
	}
	for i := range this.Grants {
		if this.Grants[i] != that1.Grants[i] {
			return fmt.Errorf("Grants this[%v](%v) Not Equal that[%v](%v)", i, this.Grants[i], i, that1.Grants[i])
		}
	}
	if len(this.Rights) != len(that1.Rights) {
		return fmt.Errorf("Rights this(%v) Not Equal that(%v)", len(this.Rights), len(that1.Rights))
	}
	for i := range this.Rights {
		if this.Rights[i] != that1.Rights[i] {
			return fmt.Errorf("Rights this[%v](%v) Not Equal that[%v](%v)", i, this.Rights[i], i, that1.Rights[i])
		}
	}
	if !this.CreatorIDs.Equal(&that1.CreatorIDs) {
		return fmt.Errorf("CreatorIDs this(%v) Not Equal that(%v)", this.CreatorIDs, that1.CreatorIDs)
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return fmt.Errorf("CreatedAt this(%v) Not Equal that(%v)", this.CreatedAt, that1.CreatedAt)
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return fmt.Errorf("UpdatedAt this(%v) Not Equal that(%v)", this.UpdatedAt, that1.UpdatedAt)
	}
	return nil
}
func (this *Client) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Client)
	if !ok {
		that2, ok := that.(Client)
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
	if !this.ClientIdentifiers.Equal(&that1.ClientIdentifiers) {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if this.Secret != that1.Secret {
		return false
	}
	if this.RedirectURI != that1.RedirectURI {
		return false
	}
	if this.State != that1.State {
		return false
	}
	if this.OfficialLabeled != that1.OfficialLabeled {
		return false
	}
	if len(this.Grants) != len(that1.Grants) {
		return false
	}
	for i := range this.Grants {
		if this.Grants[i] != that1.Grants[i] {
			return false
		}
	}
	if len(this.Rights) != len(that1.Rights) {
		return false
	}
	for i := range this.Rights {
		if this.Rights[i] != that1.Rights[i] {
			return false
		}
	}
	if !this.CreatorIDs.Equal(&that1.CreatorIDs) {
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
func (m *Client) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Client) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintClient(dAtA, i, uint64(m.ClientIdentifiers.Size()))
	n1, err := m.ClientIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Description) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintClient(dAtA, i, uint64(len(m.Description)))
		i += copy(dAtA[i:], m.Description)
	}
	if len(m.Secret) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintClient(dAtA, i, uint64(len(m.Secret)))
		i += copy(dAtA[i:], m.Secret)
	}
	if len(m.RedirectURI) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintClient(dAtA, i, uint64(len(m.RedirectURI)))
		i += copy(dAtA[i:], m.RedirectURI)
	}
	if m.State != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintClient(dAtA, i, uint64(m.State))
	}
	if m.OfficialLabeled {
		dAtA[i] = 0x30
		i++
		if m.OfficialLabeled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.Grants) > 0 {
		dAtA3 := make([]byte, len(m.Grants)*10)
		var j2 int
		for _, num := range m.Grants {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		dAtA[i] = 0x3a
		i++
		i = encodeVarintClient(dAtA, i, uint64(j2))
		i += copy(dAtA[i:], dAtA3[:j2])
	}
	if len(m.Rights) > 0 {
		dAtA5 := make([]byte, len(m.Rights)*10)
		var j4 int
		for _, num := range m.Rights {
			for num >= 1<<7 {
				dAtA5[j4] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j4++
			}
			dAtA5[j4] = uint8(num)
			j4++
		}
		dAtA[i] = 0x42
		i++
		i = encodeVarintClient(dAtA, i, uint64(j4))
		i += copy(dAtA[i:], dAtA5[:j4])
	}
	dAtA[i] = 0x4a
	i++
	i = encodeVarintClient(dAtA, i, uint64(m.CreatorIDs.Size()))
	n6, err := m.CreatorIDs.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n6
	dAtA[i] = 0x52
	i++
	i = encodeVarintClient(dAtA, i, uint64(types.SizeOfStdTime(m.CreatedAt)))
	n7, err := types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n7
	dAtA[i] = 0x5a
	i++
	i = encodeVarintClient(dAtA, i, uint64(types.SizeOfStdTime(m.UpdatedAt)))
	n8, err := types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n8
	return i, nil
}

func encodeVarintClient(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedClient(r randyClient, easy bool) *Client {
	this := &Client{}
	v1 := NewPopulatedClientIdentifiers(r, easy)
	this.ClientIdentifiers = *v1
	this.Description = randStringClient(r)
	this.Secret = randStringClient(r)
	this.RedirectURI = randStringClient(r)
	this.State = ReviewingState([]int32{0, 1, 2}[r.Intn(3)])
	this.OfficialLabeled = bool(r.Intn(2) == 0)
	v2 := r.Intn(10)
	this.Grants = make([]GrantType, v2)
	for i := 0; i < v2; i++ {
		this.Grants[i] = GrantType([]int32{0, 1, 2}[r.Intn(3)])
	}
	v3 := r.Intn(10)
	this.Rights = make([]Right, v3)
	for i := 0; i < v3; i++ {
		this.Rights[i] = Right([]int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 17, 18, 31, 32, 33, 34, 35, 36, 37, 38, 39, 51, 52, 53, 54, 55, 56, 57, 58, 71, 72, 73, 74, 75, 76, 77, 78, 79}[r.Intn(40)])
	}
	v4 := NewPopulatedUserIdentifiers(r, easy)
	this.CreatorIDs = *v4
	v5 := types.NewPopulatedStdTime(r, easy)
	this.CreatedAt = *v5
	v6 := types.NewPopulatedStdTime(r, easy)
	this.UpdatedAt = *v6
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyClient interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneClient(r randyClient) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringClient(r randyClient) string {
	v7 := r.Intn(100)
	tmps := make([]rune, v7)
	for i := 0; i < v7; i++ {
		tmps[i] = randUTF8RuneClient(r)
	}
	return string(tmps)
}
func randUnrecognizedClient(r randyClient, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldClient(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldClient(dAtA []byte, r randyClient, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateClient(dAtA, uint64(key))
		v8 := r.Int63()
		if r.Intn(2) == 0 {
			v8 *= -1
		}
		dAtA = encodeVarintPopulateClient(dAtA, uint64(v8))
	case 1:
		dAtA = encodeVarintPopulateClient(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateClient(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateClient(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateClient(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateClient(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *Client) Size() (n int) {
	var l int
	_ = l
	l = m.ClientIdentifiers.Size()
	n += 1 + l + sovClient(uint64(l))
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovClient(uint64(l))
	}
	l = len(m.Secret)
	if l > 0 {
		n += 1 + l + sovClient(uint64(l))
	}
	l = len(m.RedirectURI)
	if l > 0 {
		n += 1 + l + sovClient(uint64(l))
	}
	if m.State != 0 {
		n += 1 + sovClient(uint64(m.State))
	}
	if m.OfficialLabeled {
		n += 2
	}
	if len(m.Grants) > 0 {
		l = 0
		for _, e := range m.Grants {
			l += sovClient(uint64(e))
		}
		n += 1 + sovClient(uint64(l)) + l
	}
	if len(m.Rights) > 0 {
		l = 0
		for _, e := range m.Rights {
			l += sovClient(uint64(e))
		}
		n += 1 + sovClient(uint64(l)) + l
	}
	l = m.CreatorIDs.Size()
	n += 1 + l + sovClient(uint64(l))
	l = types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovClient(uint64(l))
	l = types.SizeOfStdTime(m.UpdatedAt)
	n += 1 + l + sovClient(uint64(l))
	return n
}

func sovClient(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozClient(x uint64) (n int) {
	return sovClient((x << 1) ^ uint64((int64(x) >> 63)))
}
func (m *Client) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClient
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
			return fmt.Errorf("proto: Client: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Client: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
				return ErrInvalidLengthClient
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ClientIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
				return ErrInvalidLengthClient
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Secret", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
				return ErrInvalidLengthClient
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Secret = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RedirectURI", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
				return ErrInvalidLengthClient
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RedirectURI = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OfficialLabeled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
			m.OfficialLabeled = bool(v != 0)
		case 7:
			if wireType == 0 {
				var v GrantType
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClient
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (GrantType(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Grants = append(m.Grants, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClient
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
					return ErrInvalidLengthClient
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v GrantType
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowClient
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (GrantType(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Grants = append(m.Grants, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Grants", wireType)
			}
		case 8:
			if wireType == 0 {
				var v Right
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClient
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
						return ErrIntOverflowClient
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
					return ErrInvalidLengthClient
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v Right
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowClient
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
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatorIDs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
				return ErrInvalidLengthClient
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CreatorIDs.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
				return ErrInvalidLengthClient
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := types.StdTimeUnmarshal(&m.CreatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClient
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
				return ErrInvalidLengthClient
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
			skippy, err := skipClient(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthClient
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
func skipClient(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowClient
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
					return 0, ErrIntOverflowClient
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
					return 0, ErrIntOverflowClient
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
				return 0, ErrInvalidLengthClient
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowClient
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
				next, err := skipClient(dAtA[start:])
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
	ErrInvalidLengthClient = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowClient   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/client.proto", fileDescriptorClient)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/client.proto", fileDescriptorClient)
}

var fileDescriptorClient = []byte{
	// 650 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x3b, 0x4c, 0xdb, 0x4e,
	0x18, 0xbf, 0xe3, 0x91, 0x3f, 0x5c, 0x20, 0x84, 0x1b, 0xf8, 0xbb, 0x19, 0xbe, 0x44, 0x95, 0x2a,
	0x85, 0xaa, 0x75, 0xaa, 0xa0, 0x0e, 0x1d, 0xf3, 0x12, 0x4a, 0x85, 0x28, 0x35, 0xa1, 0x43, 0x97,
	0xc8, 0x89, 0x2f, 0xce, 0x89, 0x60, 0x5b, 0xf6, 0x05, 0xd4, 0x8d, 0x91, 0x91, 0xb1, 0x5b, 0x2b,
	0x75, 0x61, 0x64, 0xe8, 0xc0, 0xc8, 0xc8, 0xc8, 0xc8, 0x94, 0xe2, 0xf3, 0x82, 0x3a, 0x31, 0x32,
	0x56, 0x7e, 0x15, 0x3a, 0x95, 0x4e, 0xf6, 0xf7, 0x7b, 0xe9, 0xd3, 0xef, 0xd3, 0x91, 0x57, 0x26,
	0x17, 0xc3, 0x71, 0x4f, 0xed, 0xdb, 0x7b, 0x95, 0xce, 0x90, 0x75, 0x86, 0xdc, 0x32, 0xbd, 0x4d,
	0x26, 0x0e, 0x6c, 0x77, 0xb7, 0x22, 0x84, 0x55, 0xd1, 0x1d, 0x5e, 0xe9, 0x8f, 0x38, 0xb3, 0x84,
	0xea, 0xb8, 0xb6, 0xb0, 0x69, 0x46, 0x08, 0x4b, 0xdd, 0x5f, 0x2b, 0xa8, 0x8f, 0x71, 0xea, 0x63,
	0x31, 0x8c, 0x7d, 0x85, 0xd7, 0x8f, 0xd1, 0x73, 0x83, 0x59, 0x82, 0x0f, 0x38, 0x73, 0xbd, 0xc4,
	0xf6, 0xa8, 0x05, 0x5d, 0x6e, 0x0e, 0x45, 0xea, 0x78, 0xf9, 0xc0, 0x61, 0xda, 0xa6, 0x5d, 0x89,
	0xe0, 0xde, 0x78, 0x10, 0x4d, 0xd1, 0x10, 0xfd, 0x25, 0xf2, 0xa2, 0x69, 0xdb, 0xe6, 0x88, 0xdd,
	0xab, 0x04, 0xdf, 0x63, 0x9e, 0xd0, 0xf7, 0x9c, 0x58, 0xf0, 0xf4, 0xfb, 0x0c, 0xc9, 0x34, 0xa2,
	0x06, 0xe8, 0x1b, 0x32, 0xcd, 0x0d, 0x4f, 0xc1, 0x25, 0x5c, 0xce, 0x56, 0x9f, 0xa8, 0x71, 0x13,
	0x6a, 0x4c, 0xb6, 0xef, 0x57, 0xaf, 0x2f, 0x5c, 0x4c, 0x8a, 0xe8, 0x72, 0x52, 0xc4, 0x3f, 0x27,
	0x45, 0xa4, 0x85, 0x1e, 0x5a, 0x22, 0x59, 0x83, 0x79, 0x7d, 0x97, 0x3b, 0x82, 0xdb, 0x96, 0x32,
	0x55, 0xc2, 0xe5, 0x79, 0xed, 0x21, 0x44, 0x57, 0x48, 0xc6, 0x63, 0x7d, 0x97, 0x09, 0x65, 0x3a,
	0x22, 0x93, 0x89, 0x56, 0xc9, 0x82, 0xcb, 0x0c, 0xee, 0xb2, 0xbe, 0xe8, 0x8e, 0x5d, 0xae, 0xcc,
	0x84, 0x6c, 0x7d, 0x49, 0x4e, 0x8a, 0x59, 0x2d, 0xc1, 0x77, 0xb4, 0xb6, 0x96, 0x4d, 0x45, 0x3b,
	0x2e, 0xa7, 0x2f, 0xc8, 0xac, 0x27, 0x74, 0xc1, 0x94, 0xd9, 0x12, 0x2e, 0xe7, 0xaa, 0x2b, 0xe9,
	0xaa, 0x1a, 0xdb, 0xe7, 0xec, 0x80, 0x5b, 0xe6, 0x76, 0xc8, 0x6a, 0xb1, 0x88, 0xae, 0x92, 0xbc,
	0x3d, 0x18, 0xf0, 0x3e, 0xd7, 0x47, 0xdd, 0x91, 0xde, 0x63, 0x23, 0x66, 0x28, 0x99, 0x12, 0x2e,
	0xcf, 0x69, 0x4b, 0x29, 0xbe, 0x11, 0xc3, 0x74, 0x95, 0x64, 0x4c, 0x57, 0xb7, 0x84, 0xa7, 0xfc,
	0x57, 0x9a, 0x2e, 0xe7, 0xaa, 0xcb, 0x69, 0xf2, 0x7a, 0x88, 0x76, 0x3e, 0x39, 0x4c, 0x4b, 0x04,
	0xf4, 0x19, 0xc9, 0xc4, 0x77, 0x51, 0xe6, 0x22, 0xe9, 0xe2, 0xef, 0x25, 0x42, 0x54, 0x4b, 0x48,
	0xba, 0x41, 0xb2, 0x7d, 0x97, 0xe9, 0xc2, 0x76, 0xbb, 0x61, 0xb7, 0xf3, 0x51, 0xb7, 0xff, 0xa7,
	0xda, 0x1d, 0x8f, 0xb9, 0x0f, 0x9b, 0xa5, 0x61, 0xb3, 0x72, 0x52, 0x24, 0x8d, 0xd8, 0xd3, 0x6e,
	0x7a, 0x1a, 0x49, 0xfc, 0x6d, 0xc3, 0xa3, 0x0d, 0x12, 0x4f, 0xcc, 0xe8, 0xea, 0x42, 0x21, 0x51,
	0x58, 0x41, 0x8d, 0x4f, 0xac, 0xa6, 0x27, 0x56, 0x3b, 0xe9, 0x89, 0xeb, 0x73, 0x61, 0xde, 0xf1,
	0x8f, 0x22, 0xd6, 0xe6, 0x13, 0x5f, 0x4d, 0x84, 0x21, 0x63, 0xc7, 0x48, 0x43, 0xb2, 0xff, 0x12,
	0x92, 0xf8, 0x6a, 0xe2, 0xf9, 0x7b, 0x92, 0xfb, 0xb3, 0x6d, 0xba, 0x4c, 0x16, 0xb7, 0x3b, 0xb5,
	0x4e, 0xab, 0xbb, 0xd5, 0xda, 0x6c, 0xb6, 0x37, 0xd7, 0xf3, 0x88, 0x52, 0x92, 0x8b, 0xa1, 0xda,
	0xd6, 0x96, 0xf6, 0xee, 0x43, 0xab, 0x99, 0xc7, 0xf7, 0x98, 0xd6, 0x7a, 0xdb, 0x6a, 0x74, 0x5a,
	0xcd, 0xfc, 0x54, 0x61, 0xe6, 0xe8, 0x1b, 0xa0, 0xfa, 0x17, 0x7c, 0xe1, 0x03, 0xbe, 0xf4, 0x01,
	0x5f, 0xf9, 0x80, 0xaf, 0x7d, 0xc0, 0x37, 0x3e, 0xa0, 0x5b, 0x1f, 0xd0, 0x9d, 0x0f, 0xf8, 0x50,
	0x02, 0x3a, 0x92, 0x80, 0x4e, 0x24, 0xe0, 0x53, 0x09, 0xe8, 0x4c, 0x02, 0x3e, 0x97, 0x80, 0x2f,
	0x24, 0xe0, 0x4b, 0x09, 0xf8, 0x4a, 0x02, 0xba, 0x96, 0x80, 0x6f, 0x24, 0xa0, 0x5b, 0x09, 0xf8,
	0x4e, 0x02, 0x3a, 0x0c, 0x00, 0x1d, 0x05, 0x80, 0x8f, 0x03, 0x40, 0x9f, 0x03, 0xc0, 0x5f, 0x03,
	0x40, 0x27, 0x01, 0xa0, 0xd3, 0x00, 0xf0, 0x59, 0x00, 0xf8, 0x3c, 0x00, 0xfc, 0x71, 0xf5, 0x6f,
	0x4f, 0xd0, 0xd9, 0x35, 0xc3, 0xaf, 0xd3, 0xeb, 0x65, 0xa2, 0x76, 0xd6, 0x7e, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x6b, 0xe6, 0x26, 0xf5, 0x57, 0x04, 0x00, 0x00,
}
