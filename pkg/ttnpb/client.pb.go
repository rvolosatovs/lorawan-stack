// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/client.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
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

// ClientState enum defines all the possible tenant admin reviewing states
// that a third-party client request can be at.
type ClientState int32

const (
	// State that denotes that the client request is pending to review by the
	// tenant admin.
	STATE_PENDING ClientState = 0
	// Denotes that the client request has been approved by the tenant admin
	// and therefore the client can be used.
	STATE_APPROVED ClientState = 1
	// Denotes that the client request has beenr rejected by the tenant admin
	// and therefore it cannot be used.
	STATE_REJECTED ClientState = 2
)

var ClientState_name = map[int32]string{
	0: "STATE_PENDING",
	1: "STATE_APPROVED",
	2: "STATE_REJECTED",
}
var ClientState_value = map[string]int32{
	"STATE_PENDING":  0,
	"STATE_APPROVED": 1,
	"STATE_REJECTED": 2,
}

func (ClientState) EnumDescriptor() ([]byte, []int) { return fileDescriptorClient, []int{0} }

// Scope enum defines the different scopes a third-party client can have access to.
type ClientScope int32

const (
	// Denotes whether if the client has access to manage user's applications.
	SCOPE_APPLICATION ClientScope = 0
	// Denotes wheter if the client has r-w access to user's profile.
	SCOPE_PROFILE ClientScope = 1
)

var ClientScope_name = map[int32]string{
	0: "SCOPE_APPLICATION",
	1: "SCOPE_PROFILE",
}
var ClientScope_value = map[string]int32{
	"SCOPE_APPLICATION": 0,
	"SCOPE_PROFILE":     1,
}

func (ClientScope) EnumDescriptor() ([]byte, []int) { return fileDescriptorClient, []int{1} }

// Grant enum defines the OAuth2 flows a third-party client can use to get
// access to a token.
type ClientGrant int32

const (
	// Grant type used to exchange an authorization code for an access token.
	GRANT_AUTHORIZATION_CODE ClientGrant = 0
	// Grant type used to exchange an user ID and password for an access token.
	GRANT_PASSWORD ClientGrant = 1
	// Grant type used to exchange a refresh token for an access token.
	GRANT_REFRESH_TOKEN ClientGrant = 2
)

var ClientGrant_name = map[int32]string{
	0: "GRANT_AUTHORIZATION_CODE",
	1: "GRANT_PASSWORD",
	2: "GRANT_REFRESH_TOKEN",
}
var ClientGrant_value = map[string]int32{
	"GRANT_AUTHORIZATION_CODE": 0,
	"GRANT_PASSWORD":           1,
	"GRANT_REFRESH_TOKEN":      2,
}

func (ClientGrant) EnumDescriptor() ([]byte, []int) { return fileDescriptorClient, []int{2} }

// Client is the message that defines a third-party client on the network.
type Client struct {
	ClientIdentifier `protobuf:"bytes,1,opt,name=client,embedded=client" json:"client"`
	// description is the description of the client.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// secret is the secret used to prove the client identity.
	Secret string `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
	// callback_uri is the callback URI of the client.
	CallbackURI string `protobuf:"bytes,4,opt,name=callback_uri,json=callbackUri,proto3" json:"callback_uri,omitempty" db:"callback_uri"`
	// state denotes the reviewing state of the client by the tenant admin.
	State ClientState `protobuf:"varint,5,opt,name=state,proto3,enum=ttn.v3.ClientState" json:"state,omitempty"`
	// official denotes that a client has been labelled as an official third-party
	// client by the tenant admins and therefore this fact will be shown to
	Official bool `protobuf:"varint,6,opt,name=official,proto3" json:"official,omitempty"`
	// grants denotes which OAuth2 flows can the client use to get a token.
	Grants []ClientGrant `protobuf:"varint,7,rep,packed,name=grants,enum=ttn.v3.ClientGrant" json:"grants,omitempty"`
	// scope denotes what scopes the client will have access to.
	Scope []ClientScope `protobuf:"varint,8,rep,packed,name=scope,enum=ttn.v3.ClientScope" json:"scope,omitempty"`
	// created_at denotes when the client was created.
	CreatedAt time.Time `protobuf:"bytes,9,opt,name=created_at,json=createdAt,stdtime" json:"created_at"`
	// deleted_at denotes when the client was deleted.
	DeletedAt *time.Time `protobuf:"bytes,10,opt,name=deleted_at,json=deletedAt,stdtime" json:"deleted_at,omitempty"`
}

func (m *Client) Reset()                    { *m = Client{} }
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

func (m *Client) GetCallbackURI() string {
	if m != nil {
		return m.CallbackURI
	}
	return ""
}

func (m *Client) GetState() ClientState {
	if m != nil {
		return m.State
	}
	return STATE_PENDING
}

func (m *Client) GetOfficial() bool {
	if m != nil {
		return m.Official
	}
	return false
}

func (m *Client) GetGrants() []ClientGrant {
	if m != nil {
		return m.Grants
	}
	return nil
}

func (m *Client) GetScope() []ClientScope {
	if m != nil {
		return m.Scope
	}
	return nil
}

func (m *Client) GetCreatedAt() time.Time {
	if m != nil {
		return m.CreatedAt
	}
	return time.Time{}
}

func (m *Client) GetDeletedAt() *time.Time {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

func init() {
	proto.RegisterType((*Client)(nil), "ttn.v3.Client")
	proto.RegisterEnum("ttn.v3.ClientState", ClientState_name, ClientState_value)
	proto.RegisterEnum("ttn.v3.ClientScope", ClientScope_name, ClientScope_value)
	proto.RegisterEnum("ttn.v3.ClientGrant", ClientGrant_name, ClientGrant_value)
}
func (x ClientState) String() string {
	s, ok := ClientState_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x ClientScope) String() string {
	s, ok := ClientScope_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x ClientGrant) String() string {
	s, ok := ClientGrant_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
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
	i = encodeVarintClient(dAtA, i, uint64(m.ClientIdentifier.Size()))
	n1, err := m.ClientIdentifier.MarshalTo(dAtA[i:])
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
	if len(m.CallbackURI) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintClient(dAtA, i, uint64(len(m.CallbackURI)))
		i += copy(dAtA[i:], m.CallbackURI)
	}
	if m.State != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintClient(dAtA, i, uint64(m.State))
	}
	if m.Official {
		dAtA[i] = 0x30
		i++
		if m.Official {
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
	if len(m.Scope) > 0 {
		dAtA5 := make([]byte, len(m.Scope)*10)
		var j4 int
		for _, num := range m.Scope {
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
	i = encodeVarintClient(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)))
	n6, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n6
	if m.DeletedAt != nil {
		dAtA[i] = 0x52
		i++
		i = encodeVarintClient(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.DeletedAt)))
		n7, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.DeletedAt, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	return i, nil
}

func encodeFixed64Client(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32Client(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
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
func (m *Client) Size() (n int) {
	var l int
	_ = l
	l = m.ClientIdentifier.Size()
	n += 1 + l + sovClient(uint64(l))
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovClient(uint64(l))
	}
	l = len(m.Secret)
	if l > 0 {
		n += 1 + l + sovClient(uint64(l))
	}
	l = len(m.CallbackURI)
	if l > 0 {
		n += 1 + l + sovClient(uint64(l))
	}
	if m.State != 0 {
		n += 1 + sovClient(uint64(m.State))
	}
	if m.Official {
		n += 2
	}
	if len(m.Grants) > 0 {
		l = 0
		for _, e := range m.Grants {
			l += sovClient(uint64(e))
		}
		n += 1 + sovClient(uint64(l)) + l
	}
	if len(m.Scope) > 0 {
		l = 0
		for _, e := range m.Scope {
			l += sovClient(uint64(e))
		}
		n += 1 + sovClient(uint64(l)) + l
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovClient(uint64(l))
	if m.DeletedAt != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.DeletedAt)
		n += 1 + l + sovClient(uint64(l))
	}
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
	return sovClient(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Client) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Client{`,
		`ClientIdentifier:` + strings.Replace(strings.Replace(this.ClientIdentifier.String(), "ClientIdentifier", "ClientIdentifier", 1), `&`, ``, 1) + `,`,
		`Description:` + fmt.Sprintf("%v", this.Description) + `,`,
		`Secret:` + fmt.Sprintf("%v", this.Secret) + `,`,
		`CallbackURI:` + fmt.Sprintf("%v", this.CallbackURI) + `,`,
		`State:` + fmt.Sprintf("%v", this.State) + `,`,
		`Official:` + fmt.Sprintf("%v", this.Official) + `,`,
		`Grants:` + fmt.Sprintf("%v", this.Grants) + `,`,
		`Scope:` + fmt.Sprintf("%v", this.Scope) + `,`,
		`CreatedAt:` + strings.Replace(strings.Replace(this.CreatedAt.String(), "Timestamp", "google_protobuf1.Timestamp", 1), `&`, ``, 1) + `,`,
		`DeletedAt:` + strings.Replace(fmt.Sprintf("%v", this.DeletedAt), "Timestamp", "google_protobuf1.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringClient(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
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
				return fmt.Errorf("proto: wrong wireType = %d for field ClientIdentifier", wireType)
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
			if err := m.ClientIdentifier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
				return fmt.Errorf("proto: wrong wireType = %d for field CallbackURI", wireType)
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
			m.CallbackURI = string(dAtA[iNdEx:postIndex])
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
				m.State |= (ClientState(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Official", wireType)
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
			m.Official = bool(v != 0)
		case 7:
			if wireType == 0 {
				var v ClientGrant
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClient
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ClientGrant(b) & 0x7F) << shift
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
					var v ClientGrant
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowClient
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ClientGrant(b) & 0x7F) << shift
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
				var v ClientScope
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClient
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ClientScope(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Scope = append(m.Scope, v)
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
					var v ClientScope
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowClient
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ClientScope(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Scope = append(m.Scope, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Scope", wireType)
			}
		case 9:
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
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.CreatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeletedAt", wireType)
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
			if m.DeletedAt == nil {
				m.DeletedAt = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.DeletedAt, dAtA[iNdEx:postIndex]); err != nil {
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

var fileDescriptorClient = []byte{
	// 598 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x3d, 0x6f, 0xd3, 0x40,
	0x18, 0xce, 0xf5, 0xc3, 0xa4, 0x17, 0xa8, 0xd2, 0xab, 0x00, 0x2b, 0x42, 0x4e, 0x54, 0x31, 0xa4,
	0x45, 0x38, 0xa8, 0x15, 0x4b, 0x25, 0x84, 0x5c, 0xc7, 0x4d, 0x0d, 0x95, 0x1d, 0x9d, 0x5d, 0x90,
	0x3a, 0x10, 0xd9, 0xce, 0xc5, 0x3d, 0xd5, 0xb5, 0x2d, 0xfb, 0x0a, 0x2b, 0x23, 0x23, 0xff, 0x81,
	0x3f, 0xd3, 0xb1, 0x23, 0x53, 0x00, 0x4f, 0x0c, 0x0c, 0x88, 0x5f, 0x80, 0x7c, 0xe7, 0xb4, 0x41,
	0x20, 0x95, 0xc9, 0x7e, 0x9f, 0xaf, 0x7b, 0xfc, 0x9e, 0x0c, 0x9f, 0x84, 0x94, 0x9d, 0x9c, 0xfb,
	0x6a, 0x90, 0x9c, 0xf5, 0xdc, 0x13, 0xe2, 0x9e, 0xd0, 0x38, 0xcc, 0x2d, 0xc2, 0xde, 0x25, 0xd9,
	0x69, 0x8f, 0xb1, 0xb8, 0xe7, 0xa5, 0xb4, 0x17, 0x44, 0x94, 0xc4, 0x4c, 0x4d, 0xb3, 0x84, 0x25,
	0x48, 0x62, 0x2c, 0x56, 0xdf, 0xee, 0xb4, 0x1e, 0xcf, 0x39, 0xc3, 0x24, 0x4c, 0x7a, 0x9c, 0xf6,
	0xcf, 0x27, 0x7c, 0xe2, 0x03, 0x7f, 0x13, 0xb6, 0xd6, 0xd3, 0xff, 0x39, 0x88, 0x8e, 0x49, 0xcc,
	0xe8, 0x84, 0x92, 0x2c, 0xaf, 0x6c, 0xed, 0x30, 0x49, 0xc2, 0x88, 0x5c, 0x87, 0x33, 0x7a, 0x46,
	0x72, 0xe6, 0x9d, 0xa5, 0x42, 0xb0, 0xf1, 0x63, 0x11, 0x4a, 0x3a, 0xef, 0x87, 0x76, 0xa1, 0x24,
	0x9a, 0xca, 0xa0, 0x03, 0xba, 0x8d, 0x6d, 0x59, 0x15, 0x55, 0x55, 0xc1, 0x9b, 0x57, 0xe1, 0x7b,
	0xf5, 0x8b, 0x69, 0xbb, 0x76, 0x39, 0x6d, 0x03, 0x5c, 0x39, 0x50, 0x07, 0x36, 0xc6, 0x24, 0x0f,
	0x32, 0x9a, 0x32, 0x9a, 0xc4, 0xf2, 0x42, 0x07, 0x74, 0x57, 0xf0, 0x3c, 0x84, 0xee, 0x41, 0x29,
	0x27, 0x41, 0x46, 0x98, 0xbc, 0xc8, 0xc9, 0x6a, 0x42, 0x03, 0x78, 0x3b, 0xf0, 0xa2, 0xc8, 0xf7,
	0x82, 0xd3, 0xd1, 0x79, 0x46, 0xe5, 0xa5, 0x92, 0xdd, 0x7b, 0x58, 0x4c, 0xdb, 0x0d, 0xbd, 0xc2,
	0x8f, 0xb0, 0xf9, 0x6b, 0xda, 0x5e, 0x1b, 0xfb, 0xbb, 0x1b, 0xf3, 0xd2, 0x0d, 0xdc, 0x98, 0x8d,
	0x47, 0x19, 0x45, 0x9b, 0x70, 0x39, 0x67, 0x1e, 0x23, 0xf2, 0x72, 0x07, 0x74, 0x57, 0xb7, 0xd7,
	0xff, 0x6c, 0xef, 0x94, 0x14, 0x16, 0x0a, 0xd4, 0x82, 0xf5, 0x64, 0x32, 0xa1, 0x01, 0xf5, 0x22,
	0x59, 0xea, 0x80, 0x6e, 0x1d, 0x5f, 0xcd, 0xe8, 0x11, 0x94, 0xc2, 0xcc, 0x8b, 0x59, 0x2e, 0xdf,
	0xea, 0x2c, 0xfe, 0x9d, 0x33, 0x28, 0x39, 0x5c, 0x49, 0xf8, 0x99, 0x41, 0x92, 0x12, 0xb9, 0xfe,
	0x2f, 0xad, 0x53, 0x52, 0x58, 0x28, 0x90, 0x0e, 0x61, 0x90, 0x11, 0x8f, 0x91, 0xf1, 0xc8, 0x63,
	0xf2, 0x0a, 0xdf, 0x70, 0x4b, 0x15, 0xd7, 0xa3, 0xce, 0xae, 0x47, 0x75, 0x67, 0xd7, 0x23, 0x76,
	0xfc, 0xf1, 0x4b, 0x1b, 0xe0, 0x95, 0xca, 0xa7, 0x31, 0xf4, 0x1c, 0xc2, 0x31, 0x89, 0x48, 0x15,
	0x02, 0x6f, 0x0c, 0x59, 0x12, 0x01, 0x95, 0x47, 0x63, 0x5b, 0x16, 0x6c, 0xcc, 0xed, 0x03, 0xad,
	0xc1, 0x3b, 0x8e, 0xab, 0xb9, 0xc6, 0x68, 0x68, 0x58, 0x7d, 0xd3, 0x1a, 0x34, 0x6b, 0x08, 0xc1,
	0x55, 0x01, 0x69, 0xc3, 0x21, 0xb6, 0x5f, 0x19, 0xfd, 0x26, 0xb8, 0xc6, 0xb0, 0xf1, 0xc2, 0xd0,
	0x5d, 0xa3, 0xdf, 0x5c, 0x68, 0x2d, 0x7d, 0xf8, 0xa4, 0xd4, 0xb6, 0x9e, 0x5d, 0xe5, 0xf1, 0x8f,
	0xbc, 0x0b, 0xd7, 0x1c, 0xdd, 0x1e, 0x72, 0xf3, 0xa1, 0xa9, 0x6b, 0xae, 0x69, 0x5b, 0xcd, 0x1a,
	0x3f, 0x86, 0xc3, 0x43, 0x6c, 0xef, 0x9b, 0x87, 0x46, 0x13, 0x54, 0xf6, 0x37, 0x33, 0x3b, 0x5f,
	0x2b, 0x7a, 0x00, 0xe5, 0x01, 0xd6, 0x2c, 0x77, 0xa4, 0x1d, 0xb9, 0x07, 0x36, 0x36, 0x8f, 0x79,
	0xc0, 0x48, 0xb7, 0xfb, 0x86, 0x68, 0x26, 0xd8, 0xa1, 0xe6, 0x38, 0xaf, 0x6d, 0x5c, 0x36, 0xbb,
	0x0f, 0xd7, 0x05, 0x86, 0x8d, 0x7d, 0x6c, 0x38, 0x07, 0x23, 0xd7, 0x7e, 0x69, 0x58, 0xb3, 0x7a,
	0x7b, 0x83, 0xcf, 0xdf, 0x94, 0xda, 0xfb, 0x42, 0x01, 0x17, 0x85, 0x02, 0x2e, 0x0b, 0x05, 0x7c,
	0x2d, 0x14, 0xf0, 0xbd, 0x50, 0x6a, 0x3f, 0x0b, 0x05, 0x1c, 0x6f, 0xde, 0xf4, 0x4f, 0xa5, 0xa7,
	0x61, 0xf9, 0x4c, 0x7d, 0x5f, 0xe2, 0xcb, 0xdd, 0xf9, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x6b, 0x77,
	0x9a, 0xb5, 0xf0, 0x03, 0x00, 0x00,
}
