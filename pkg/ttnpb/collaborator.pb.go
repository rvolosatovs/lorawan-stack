// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/collaborator.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Collaborator is the message to define a collaborator of an entity on the network.
type Collaborator struct {
	UserIDIdentifier `protobuf:"bytes,1,opt,name=user,embedded=user" json:"user"`
	// rights is the list of rights the user bears.
	Rights []Right `protobuf:"varint,2,rep,packed,name=rights,enum=ttn.v3.Right" json:"rights,omitempty"`
}

func (m *Collaborator) Reset()                    { *m = Collaborator{} }
func (*Collaborator) ProtoMessage()               {}
func (*Collaborator) Descriptor() ([]byte, []int) { return fileDescriptorCollaborator, []int{0} }

func (m *Collaborator) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

func init() {
	proto.RegisterType((*Collaborator)(nil), "ttn.v3.Collaborator")
}
func (m *Collaborator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Collaborator) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintCollaborator(dAtA, i, uint64(m.UserIDIdentifier.Size()))
	n1, err := m.UserIDIdentifier.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Rights) > 0 {
		dAtA3 := make([]byte, len(m.Rights)*10)
		var j2 int
		for _, num := range m.Rights {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		dAtA[i] = 0x12
		i++
		i = encodeVarintCollaborator(dAtA, i, uint64(j2))
		i += copy(dAtA[i:], dAtA3[:j2])
	}
	return i, nil
}

func encodeFixed64Collaborator(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32Collaborator(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
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
func (m *Collaborator) Size() (n int) {
	var l int
	_ = l
	l = m.UserIDIdentifier.Size()
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
	return sovCollaborator(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Collaborator) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Collaborator{`,
		`UserIDIdentifier:` + strings.Replace(strings.Replace(this.UserIDIdentifier.String(), "UserIDIdentifier", "UserIDIdentifier", 1), `&`, ``, 1) + `,`,
		`Rights:` + fmt.Sprintf("%v", this.Rights) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringCollaborator(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Collaborator) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Collaborator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Collaborator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserIDIdentifier", wireType)
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
			if err := m.UserIDIdentifier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
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

var fileDescriptorCollaborator = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x4b, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x0f, 0xc9, 0x48, 0x0d, 0xc9, 0xc8, 0xcc, 0x4b, 0x2f,
	0xf6, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x2f, 0x29, 0xc9, 0xd3, 0x4f, 0x2c, 0xc8, 0xd4,
	0x4f, 0xce, 0xcf, 0xc9, 0x49, 0x4c, 0xca, 0x2f, 0x4a, 0x2c, 0xc9, 0x2f, 0xd2, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x29, 0xc9, 0xd3, 0x2b, 0x33, 0x96, 0xd2, 0x45, 0xd2, 0x9f, 0x9e,
	0x9f, 0x9e, 0xaf, 0x0f, 0x96, 0x4e, 0x2a, 0x4d, 0x03, 0xf3, 0xc0, 0x1c, 0x30, 0x0b, 0xa2, 0x4d,
	0xca, 0x94, 0x18, 0xeb, 0x32, 0x53, 0x52, 0xf3, 0x4a, 0x32, 0xd3, 0x32, 0x53, 0x8b, 0x8a, 0xa1,
	0xda, 0x0c, 0x88, 0xd1, 0x56, 0x94, 0x99, 0x9e, 0x51, 0x02, 0xd5, 0xa1, 0x94, 0xcb, 0xc5, 0xe3,
	0x8c, 0xe4, 0x6a, 0x21, 0x33, 0x2e, 0x96, 0xd2, 0xe2, 0xd4, 0x22, 0x09, 0x46, 0x05, 0x46, 0x0d,
	0x6e, 0x23, 0x09, 0x3d, 0x88, 0xf3, 0xf5, 0x42, 0x8b, 0x53, 0x8b, 0x3c, 0x5d, 0x3c, 0xe1, 0x16,
	0x3a, 0x71, 0x9c, 0xb8, 0x27, 0xcf, 0x70, 0xe1, 0x9e, 0x3c, 0x63, 0x10, 0x58, 0xbd, 0x90, 0x2a,
	0x17, 0x1b, 0xc4, 0x5c, 0x09, 0x26, 0x05, 0x66, 0x0d, 0x3e, 0x23, 0x5e, 0x98, 0xce, 0x20, 0x90,
	0x68, 0x10, 0x54, 0xd2, 0xc9, 0xfd, 0xc6, 0x43, 0x39, 0x86, 0x86, 0x47, 0x72, 0x8c, 0x27, 0x1e,
	0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x8b, 0x47, 0x72, 0x0c, 0x1f,
	0x1e, 0xc9, 0x31, 0x46, 0x69, 0x12, 0x72, 0x7e, 0x41, 0x76, 0x3a, 0x88, 0x2e, 0x48, 0x4a, 0x62,
	0x03, 0x3b, 0xdf, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x21, 0xbe, 0x24, 0x99, 0x98, 0x01, 0x00,
	0x00,
}
