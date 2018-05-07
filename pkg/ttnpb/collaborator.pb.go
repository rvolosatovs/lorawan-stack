// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/collaborator.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ApplicationCollaborator struct {
	// application_ids are the application identifiers.
	ApplicationIdentifiers `protobuf:"bytes,1,opt,name=application_ids,json=applicationIds,embedded=application_ids" json:"application_ids"`
	// collaborator_ids is either an organization's ID or user's ID.
	OrganizationOrUserIdentifiers `protobuf:"bytes,2,opt,name=collaborator_ids,json=collaboratorIds,embedded=collaborator_ids" json:"collaborator_ids"`
	// rights is the list of rights the user bears to the application.
	Rights []Right `protobuf:"varint,3,rep,packed,name=rights,enum=ttn.v3.Right" json:"rights,omitempty"`
}

func (m *ApplicationCollaborator) Reset()         { *m = ApplicationCollaborator{} }
func (m *ApplicationCollaborator) String() string { return proto.CompactTextString(m) }
func (*ApplicationCollaborator) ProtoMessage()    {}
func (*ApplicationCollaborator) Descriptor() ([]byte, []int) {
	return fileDescriptorCollaborator, []int{0}
}

func (m *ApplicationCollaborator) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

type GatewayCollaborator struct {
	// gateway_ids are the gateway identifiers.
	GatewayIdentifiers `protobuf:"bytes,1,opt,name=gateway_ids,json=gatewayIds,embedded=gateway_ids" json:"gateway_ids"`
	// collaborator_ids is either an organization's ID or user's ID.
	OrganizationOrUserIdentifiers `protobuf:"bytes,2,opt,name=collaborator_ids,json=collaboratorIds,embedded=collaborator_ids" json:"collaborator_ids"`
	// rights is the list of rights the user bears to the application.
	Rights []Right `protobuf:"varint,3,rep,packed,name=rights,enum=ttn.v3.Right" json:"rights,omitempty"`
}

func (m *GatewayCollaborator) Reset()                    { *m = GatewayCollaborator{} }
func (m *GatewayCollaborator) String() string            { return proto.CompactTextString(m) }
func (*GatewayCollaborator) ProtoMessage()               {}
func (*GatewayCollaborator) Descriptor() ([]byte, []int) { return fileDescriptorCollaborator, []int{1} }

func (m *GatewayCollaborator) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

func init() {
	proto.RegisterType((*ApplicationCollaborator)(nil), "ttn.v3.ApplicationCollaborator")
	golang_proto.RegisterType((*ApplicationCollaborator)(nil), "ttn.v3.ApplicationCollaborator")
	proto.RegisterType((*GatewayCollaborator)(nil), "ttn.v3.GatewayCollaborator")
	golang_proto.RegisterType((*GatewayCollaborator)(nil), "ttn.v3.GatewayCollaborator")
}
func (this *ApplicationCollaborator) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*ApplicationCollaborator)
	if !ok {
		that2, ok := that.(ApplicationCollaborator)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *ApplicationCollaborator")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *ApplicationCollaborator but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *ApplicationCollaborator but is not nil && this == nil")
	}
	if !this.ApplicationIdentifiers.Equal(&that1.ApplicationIdentifiers) {
		return fmt.Errorf("ApplicationIdentifiers this(%v) Not Equal that(%v)", this.ApplicationIdentifiers, that1.ApplicationIdentifiers)
	}
	if !this.OrganizationOrUserIdentifiers.Equal(&that1.OrganizationOrUserIdentifiers) {
		return fmt.Errorf("OrganizationOrUserIdentifiers this(%v) Not Equal that(%v)", this.OrganizationOrUserIdentifiers, that1.OrganizationOrUserIdentifiers)
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
func (this *ApplicationCollaborator) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ApplicationCollaborator)
	if !ok {
		that2, ok := that.(ApplicationCollaborator)
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
	if !this.ApplicationIdentifiers.Equal(&that1.ApplicationIdentifiers) {
		return false
	}
	if !this.OrganizationOrUserIdentifiers.Equal(&that1.OrganizationOrUserIdentifiers) {
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
func (this *GatewayCollaborator) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*GatewayCollaborator)
	if !ok {
		that2, ok := that.(GatewayCollaborator)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *GatewayCollaborator")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *GatewayCollaborator but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *GatewayCollaborator but is not nil && this == nil")
	}
	if !this.GatewayIdentifiers.Equal(&that1.GatewayIdentifiers) {
		return fmt.Errorf("GatewayIdentifiers this(%v) Not Equal that(%v)", this.GatewayIdentifiers, that1.GatewayIdentifiers)
	}
	if !this.OrganizationOrUserIdentifiers.Equal(&that1.OrganizationOrUserIdentifiers) {
		return fmt.Errorf("OrganizationOrUserIdentifiers this(%v) Not Equal that(%v)", this.OrganizationOrUserIdentifiers, that1.OrganizationOrUserIdentifiers)
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
func (this *GatewayCollaborator) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GatewayCollaborator)
	if !ok {
		that2, ok := that.(GatewayCollaborator)
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
	if !this.GatewayIdentifiers.Equal(&that1.GatewayIdentifiers) {
		return false
	}
	if !this.OrganizationOrUserIdentifiers.Equal(&that1.OrganizationOrUserIdentifiers) {
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
func (m *ApplicationCollaborator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplicationCollaborator) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintCollaborator(dAtA, i, uint64(m.ApplicationIdentifiers.Size()))
	n1, err := m.ApplicationIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x12
	i++
	i = encodeVarintCollaborator(dAtA, i, uint64(m.OrganizationOrUserIdentifiers.Size()))
	n2, err := m.OrganizationOrUserIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	if len(m.Rights) > 0 {
		dAtA4 := make([]byte, len(m.Rights)*10)
		var j3 int
		for _, num := range m.Rights {
			for num >= 1<<7 {
				dAtA4[j3] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j3++
			}
			dAtA4[j3] = uint8(num)
			j3++
		}
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCollaborator(dAtA, i, uint64(j3))
		i += copy(dAtA[i:], dAtA4[:j3])
	}
	return i, nil
}

func (m *GatewayCollaborator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GatewayCollaborator) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintCollaborator(dAtA, i, uint64(m.GatewayIdentifiers.Size()))
	n5, err := m.GatewayIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	dAtA[i] = 0x12
	i++
	i = encodeVarintCollaborator(dAtA, i, uint64(m.OrganizationOrUserIdentifiers.Size()))
	n6, err := m.OrganizationOrUserIdentifiers.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n6
	if len(m.Rights) > 0 {
		dAtA8 := make([]byte, len(m.Rights)*10)
		var j7 int
		for _, num := range m.Rights {
			for num >= 1<<7 {
				dAtA8[j7] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j7++
			}
			dAtA8[j7] = uint8(num)
			j7++
		}
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCollaborator(dAtA, i, uint64(j7))
		i += copy(dAtA[i:], dAtA8[:j7])
	}
	return i, nil
}

func encodeVarintCollaborator(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedApplicationCollaborator(r randyCollaborator, easy bool) *ApplicationCollaborator {
	this := &ApplicationCollaborator{}
	v1 := NewPopulatedApplicationIdentifiers(r, easy)
	this.ApplicationIdentifiers = *v1
	v2 := NewPopulatedOrganizationOrUserIdentifiers(r, easy)
	this.OrganizationOrUserIdentifiers = *v2
	v3 := r.Intn(10)
	this.Rights = make([]Right, v3)
	for i := 0; i < v3; i++ {
		this.Rights[i] = Right([]int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43}[r.Intn(44)])
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedGatewayCollaborator(r randyCollaborator, easy bool) *GatewayCollaborator {
	this := &GatewayCollaborator{}
	v4 := NewPopulatedGatewayIdentifiers(r, easy)
	this.GatewayIdentifiers = *v4
	v5 := NewPopulatedOrganizationOrUserIdentifiers(r, easy)
	this.OrganizationOrUserIdentifiers = *v5
	v6 := r.Intn(10)
	this.Rights = make([]Right, v6)
	for i := 0; i < v6; i++ {
		this.Rights[i] = Right([]int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43}[r.Intn(44)])
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyCollaborator interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneCollaborator(r randyCollaborator) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringCollaborator(r randyCollaborator) string {
	v7 := r.Intn(100)
	tmps := make([]rune, v7)
	for i := 0; i < v7; i++ {
		tmps[i] = randUTF8RuneCollaborator(r)
	}
	return string(tmps)
}
func randUnrecognizedCollaborator(r randyCollaborator, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldCollaborator(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldCollaborator(dAtA []byte, r randyCollaborator, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateCollaborator(dAtA, uint64(key))
		v8 := r.Int63()
		if r.Intn(2) == 0 {
			v8 *= -1
		}
		dAtA = encodeVarintPopulateCollaborator(dAtA, uint64(v8))
	case 1:
		dAtA = encodeVarintPopulateCollaborator(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateCollaborator(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateCollaborator(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateCollaborator(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateCollaborator(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *ApplicationCollaborator) Size() (n int) {
	var l int
	_ = l
	l = m.ApplicationIdentifiers.Size()
	n += 1 + l + sovCollaborator(uint64(l))
	l = m.OrganizationOrUserIdentifiers.Size()
	n += 1 + l + sovCollaborator(uint64(l))
	if len(m.Rights) > 0 {
		l = 0
		for _, e := range m.Rights {
			l += sovCollaborator(uint64(e))
		}
		n += 1 + sovCollaborator(uint64(l)) + l
	}
	return n
}

func (m *GatewayCollaborator) Size() (n int) {
	var l int
	_ = l
	l = m.GatewayIdentifiers.Size()
	n += 1 + l + sovCollaborator(uint64(l))
	l = m.OrganizationOrUserIdentifiers.Size()
	n += 1 + l + sovCollaborator(uint64(l))
	if len(m.Rights) > 0 {
		l = 0
		for _, e := range m.Rights {
			l += sovCollaborator(uint64(e))
		}
		n += 1 + sovCollaborator(uint64(l)) + l
	}
	return n
}

func sovCollaborator(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCollaborator(x uint64) (n int) {
	return sovCollaborator((x << 1) ^ uint64((int64(x) >> 63)))
}
func (m *ApplicationCollaborator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCollaborator
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
			return fmt.Errorf("proto: ApplicationCollaborator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplicationCollaborator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApplicationIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCollaborator
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
				return ErrInvalidLengthCollaborator
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ApplicationIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrganizationOrUserIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCollaborator
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
				return ErrInvalidLengthCollaborator
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OrganizationOrUserIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v Right
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCollaborator
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
						return ErrIntOverflowCollaborator
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
					return ErrInvalidLengthCollaborator
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v Right
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCollaborator
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
			skippy, err := skipCollaborator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCollaborator
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
func (m *GatewayCollaborator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCollaborator
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
			return fmt.Errorf("proto: GatewayCollaborator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GatewayCollaborator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GatewayIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCollaborator
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
				return ErrInvalidLengthCollaborator
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrganizationOrUserIdentifiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCollaborator
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
				return ErrInvalidLengthCollaborator
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OrganizationOrUserIdentifiers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v Right
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCollaborator
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
						return ErrIntOverflowCollaborator
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
					return ErrInvalidLengthCollaborator
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v Right
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCollaborator
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
			skippy, err := skipCollaborator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCollaborator
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
func skipCollaborator(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCollaborator
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
					return 0, ErrIntOverflowCollaborator
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
					return 0, ErrIntOverflowCollaborator
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
				return 0, ErrInvalidLengthCollaborator
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCollaborator
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
				next, err := skipCollaborator(dAtA[start:])
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
	ErrInvalidLengthCollaborator = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCollaborator   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/collaborator.proto", fileDescriptorCollaborator)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/collaborator.proto", fileDescriptorCollaborator)
}

var fileDescriptorCollaborator = []byte{
	// 425 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x92, 0x31, 0x6c, 0xd3, 0x40,
	0x14, 0x86, 0xef, 0x51, 0x29, 0x42, 0x57, 0xd1, 0xa2, 0x30, 0x50, 0x65, 0x78, 0xad, 0x2a, 0x55,
	0x2a, 0x03, 0x36, 0x6a, 0x05, 0x3b, 0x45, 0x08, 0xb1, 0x50, 0x11, 0x95, 0xa5, 0x0b, 0x3a, 0xa7,
	0xee, 0xf9, 0xd4, 0xe0, 0xb3, 0xce, 0x57, 0x2a, 0x98, 0x32, 0x66, 0x64, 0x64, 0x83, 0x31, 0x63,
	0xc6, 0x8c, 0x19, 0x33, 0x66, 0xcc, 0x14, 0xe2, 0xbb, 0x25, 0x1b, 0x19, 0x33, 0xa2, 0x9c, 0x63,
	0xc5, 0x48, 0x91, 0xc8, 0xc8, 0x64, 0xdf, 0x7b, 0xf7, 0x7d, 0xba, 0xff, 0xe9, 0xd1, 0x17, 0x5c,
	0xe8, 0xe8, 0x36, 0xf0, 0x1a, 0xf2, 0x93, 0x7f, 0x11, 0x85, 0x17, 0x91, 0x88, 0x79, 0xfa, 0x2e,
	0xd4, 0x77, 0x52, 0xdd, 0xf8, 0x5a, 0xc7, 0x3e, 0x4b, 0x84, 0xdf, 0x90, 0xcd, 0x26, 0x0b, 0xa4,
	0x62, 0x5a, 0x2a, 0x2f, 0x51, 0x52, 0xcb, 0x6a, 0x45, 0xeb, 0xd8, 0xfb, 0x7c, 0x5a, 0x7b, 0xbe,
	0x09, 0x2f, 0xae, 0xc2, 0x58, 0x8b, 0x6b, 0x11, 0xaa, 0x34, 0xc7, 0x6b, 0xcf, 0x36, 0xc1, 0x94,
	0xe0, 0x91, 0x2e, 0x88, 0xa7, 0x25, 0x82, 0x4b, 0x2e, 0x7d, 0x57, 0x0e, 0x6e, 0xaf, 0xdd, 0xc9,
	0x1d, 0xdc, 0x5f, 0x7e, 0xfd, 0xf0, 0x37, 0xd0, 0xc7, 0x2f, 0x93, 0xa4, 0x29, 0x1a, 0x4c, 0x0b,
	0x19, 0xbf, 0x2a, 0x25, 0xa8, 0xbe, 0xa7, 0xbb, 0x6c, 0xd5, 0xfa, 0x28, 0xae, 0xd2, 0x3d, 0x38,
	0x80, 0xe3, 0xed, 0x13, 0xf4, 0xf2, 0x54, 0x5e, 0x89, 0x7c, 0xbb, 0x7a, 0xfb, 0xd9, 0xfd, 0xc1,
	0x78, 0x9f, 0x0c, 0xc7, 0xfb, 0x50, 0xdf, 0x61, 0xe5, 0x1b, 0x69, 0xf5, 0x92, 0x3e, 0x2c, 0x0f,
	0xc9, 0x39, 0xef, 0x39, 0xe7, 0x51, 0xe1, 0x3c, 0x57, 0x9c, 0xc5, 0xe2, 0xab, 0x43, 0xce, 0xd5,
	0x87, 0x34, 0x54, 0xeb, 0xd5, 0xbb, 0x65, 0xd1, 0xc2, 0x7d, 0x44, 0x2b, 0xf9, 0x24, 0xf6, 0xb6,
	0x0e, 0xb6, 0x8e, 0x77, 0x4e, 0x1e, 0x14, 0xc6, 0xfa, 0xa2, 0x5a, 0x5f, 0x36, 0x0f, 0x7f, 0x01,
	0x7d, 0xf4, 0x86, 0xe9, 0xf0, 0x8e, 0x7d, 0xf9, 0x2b, 0xed, 0x6b, 0xba, 0xcd, 0xf3, 0x72, 0x29,
	0x69, 0xad, 0x70, 0x2c, 0x89, 0xf5, 0x4f, 0xa1, 0xbc, 0xe8, 0xfe, 0x0f, 0x09, 0xcf, 0x7e, 0xc0,
	0x20, 0x43, 0x18, 0x66, 0x08, 0xa3, 0x0c, 0x61, 0x92, 0x21, 0x4c, 0x33, 0x24, 0xb3, 0x0c, 0xc9,
	0x3c, 0x43, 0x68, 0x19, 0x24, 0x6d, 0x83, 0xa4, 0x63, 0x10, 0xba, 0x06, 0x49, 0xcf, 0x20, 0xf4,
	0x0d, 0xc2, 0xc0, 0x20, 0x0c, 0x0d, 0xc2, 0xc8, 0x20, 0x99, 0x18, 0x84, 0xa9, 0x41, 0x32, 0x33,
	0x08, 0x73, 0x83, 0xa4, 0x65, 0x91, 0xb4, 0x2d, 0xc2, 0x37, 0x8b, 0xe4, 0xbb, 0x45, 0xf8, 0x69,
	0x91, 0x74, 0x2c, 0x92, 0xae, 0x45, 0xe8, 0x59, 0x84, 0xbe, 0x45, 0xb8, 0x7c, 0xf2, 0xaf, 0x5d,
	0x4d, 0x6e, 0xf8, 0xe2, 0x9b, 0x04, 0x41, 0xc5, 0x2d, 0xdf, 0xe9, 0x9f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xf5, 0xcc, 0xa7, 0x96, 0x56, 0x03, 0x00, 0x00,
}
