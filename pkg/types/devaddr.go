// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package types

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// DevAddr is a 32-bit LoRaWAN device address.
type DevAddr [4]byte

// MinDevAddr is the lowest value possible for a DevAddr.
var MinDevAddr = DevAddr{0x00, 0x00, 0x00, 0x00}

// MaxDevAddr is the highest value possible for a DevAddr.
var MaxDevAddr = DevAddr{0xFF, 0xFF, 0xFF, 0xFF}

// IsZero returns true iff the type is zero.
func (addr DevAddr) IsZero() bool { return addr == [4]byte{} }

// String implements the Stringer interface.
func (addr DevAddr) String() string { return strings.ToUpper(hex.EncodeToString(addr[:])) }

// GoString implements the GoStringer interface.
func (addr DevAddr) GoString() string { return addr.String() }

// Size implements the Sizer interface.
func (addr DevAddr) Size() int { return 4 }

// Equal returns true iff addresses are equal.
func (addr DevAddr) Equal(other DevAddr) bool { return addr == other }

// Marshal implements the proto.Marshaler interface.
func (addr DevAddr) Marshal() ([]byte, error) { return addr.MarshalBinary() }

// MarshalTo implements the MarshalerTo function required by generated protobuf.
func (addr DevAddr) MarshalTo(data []byte) (int, error) { return marshalBinaryBytesTo(data, addr[:]) }

// Unmarshal implements the proto.Unmarshaler interface.
func (addr *DevAddr) Unmarshal(data []byte) error { return addr.UnmarshalBinary(data) }

// MarshalJSON implements the json.Marshaler interface.
func (addr DevAddr) MarshalJSON() ([]byte, error) { return marshalJSONBytes(addr[:]) }

// UnmarshalJSON implements the json.Unmarshaler interface.
func (addr *DevAddr) UnmarshalJSON(data []byte) error {
	*addr = [4]byte{}
	return unmarshalJSONBytes(addr[:], data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (addr DevAddr) MarshalBinary() ([]byte, error) { return marshalBinaryBytes(addr[:]) }

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (addr *DevAddr) UnmarshalBinary(data []byte) error {
	*addr = [4]byte{}
	return unmarshalBinaryBytes(addr[:], data)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (addr DevAddr) MarshalText() ([]byte, error) { return marshalTextBytes(addr[:]) }

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (addr *DevAddr) UnmarshalText(data []byte) error {
	*addr = [4]byte{}
	return unmarshalTextBytes(addr[:], data)
}

// Value implements sql.Valuer interface.
func (addr DevAddr) Value() (driver.Value, error) {
	return addr.MarshalText()
}

// Scan implements sql.Scanner interface.
func (addr *DevAddr) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return ErrTypeAssertion
	}
	return addr.UnmarshalText(data)
}

// MarshalNumber returns the DevAddr in a decimal form.
func (addr DevAddr) MarshalNumber() uint32 {
	return binary.BigEndian.Uint32(addr[:])
}

// UnmarshalNumber retrieves a DevAddr from a decimal form.
func (addr *DevAddr) UnmarshalNumber(n uint32) {
	*addr = [4]byte{}
	binary.BigEndian.PutUint32(addr[:], n)
}

// Before returns true if the DevAddr is strictly inferior to the DevAddr passed as an argument.
func (addr DevAddr) Before(a DevAddr) bool {
	if addr.MarshalNumber() < a.MarshalNumber() {
		return true
	}
	return false
}

// After returns true if the DevAddr is strictly superior to the DevAddr passed as an argument.
func (addr DevAddr) After(a DevAddr) bool {
	if addr.MarshalNumber() > a.MarshalNumber() {
		return true
	}
	return false
}

// BeforeOrEqual returns true if the DevAddr is inferior or equal to the DevAddr passed as an argument.
func (addr DevAddr) BeforeOrEqual(a DevAddr) bool {
	return addr.Before(a) || addr == a
}

// AfterOrEqual returns true if the DevAddr is superior or equal to the DevAddr passed as an argument.
func (addr DevAddr) AfterOrEqual(a DevAddr) bool {
	return addr.After(a) || addr == a
}

// ErrInvalidDevAddrPrefix can be returned when unmarshaling an invalid slice into a prefix
var ErrInvalidDevAddrPrefix = errors.New("invalid device address prefix")

// DevAddrPrefix is a DevAddr with a prefix length.
type DevAddrPrefix struct {
	DevAddr DevAddr
	Length  uint8
}

// NbItems returns the number of items that this prefix encapsulates.
func (prefix DevAddrPrefix) NbItems() uint64 {
	return uint64(math.Pow(2, float64(32-prefix.Length)))
}

// FirstDevAddrCovered returns the first DevAddr covered, in the numeric order.
func (prefix DevAddrPrefix) FirstDevAddrCovered() DevAddr {
	return prefix.DevAddr.Mask(prefix.Length)
}

func (prefix DevAddrPrefix) firstNumericDevAddrCovered() uint32 {
	return prefix.FirstDevAddrCovered().MarshalNumber()
}

func (prefix DevAddrPrefix) lastNumericDevAddrCovered() uint32 {
	return prefix.firstNumericDevAddrCovered() + uint32(prefix.NbItems()-1)
}

// LastDevAddrCovered returns the last DevAddr covered, in the numeric order.
func (prefix DevAddrPrefix) LastDevAddrCovered() DevAddr {
	result := DevAddr{}

	lastDevAddrNumeric := prefix.lastNumericDevAddrCovered()
	hex := fmt.Sprintf("%08X", lastDevAddrNumeric)
	result.UnmarshalText([]byte(hex))

	return result
}

// IsZero returns true iff the type is zero.
func (prefix DevAddrPrefix) IsZero() bool { return prefix.Length == 0 }

// String implements the Stringer interface.
func (prefix DevAddrPrefix) String() string {
	return fmt.Sprintf("%s/%d", prefix.DevAddr, prefix.Length)
}

// GoString implements the GoStringer interface.
func (prefix DevAddrPrefix) GoString() string { return prefix.String() }

// Size implements the Sizer interface.
func (prefix DevAddrPrefix) Size() int { return 5 }

// Equal returns true iff prefixes are equal.
func (prefix DevAddrPrefix) Equal(other DevAddrPrefix) bool {
	return prefix.Length == other.Length && prefix.DevAddr.Equal(other.DevAddr)
}

// MarshalTo implements the MarshalerTo function required by generated protobuf.
func (prefix DevAddrPrefix) MarshalTo(data []byte) (int, error) {
	return marshalBinaryBytesTo(data, append(prefix.DevAddr[:], prefix.Length))
}

// Marshal implements the proto.Marshaler interface.
func (prefix DevAddrPrefix) Marshal() ([]byte, error) { return prefix.MarshalBinary() }

// Unmarshal implements the proto.Unmarshaler interface.
func (prefix *DevAddrPrefix) Unmarshal(data []byte) error { return prefix.UnmarshalBinary(data) }

// MarshalJSON implements the json.Marshaler interface.
func (prefix DevAddrPrefix) MarshalJSON() ([]byte, error) {
	return append([]byte(`"`+base64Encoding.EncodeToString(prefix.DevAddr[:])), '/', prefix.Length, '"'), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (prefix *DevAddrPrefix) UnmarshalJSON(data []byte) error {
	if string(data) == `""` {
		*prefix = DevAddrPrefix{}
		return nil
	}
	if len(data) != 12 {
		return ErrInvalidLength
	}
	if data[9] != '/' {
		return ErrInvalidDevAddrPrefix
	}
	b := make([]byte, base64Encoding.DecodedLen(8))
	n, err := base64Encoding.Decode(b, data[1:9])
	if err != nil {
		return err
	}
	if n != 4 || copy(prefix.DevAddr[:], b) != 4 {
		return ErrInvalidDevAddrPrefix
	}
	prefix.Length = data[10]
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (prefix DevAddrPrefix) MarshalBinary() ([]byte, error) {
	return marshalBinaryBytes(append(prefix.DevAddr[:], prefix.Length))
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (prefix *DevAddrPrefix) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		*prefix = DevAddrPrefix{}
		return nil
	}
	if len(data) != 5 {
		return ErrInvalidLength
	}
	if err := prefix.DevAddr.Unmarshal(data[:4]); err != nil {
		return err
	}
	prefix.Length = data[4]
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (prefix DevAddrPrefix) MarshalText() ([]byte, error) {
	b, err := prefix.DevAddr.MarshalText()
	if err != nil {
		return nil, err
	}
	return append(append(b, '/'), []byte(strconv.Itoa(int(prefix.Length)))...), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (prefix *DevAddrPrefix) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*prefix = DevAddrPrefix{}
		return nil
	}
	if len(data) != 10 && len(data) != 11 {
		return ErrInvalidLength
	}
	if data[8] != '/' {
		return ErrInvalidDevAddrPrefix
	}
	if err := prefix.DevAddr.UnmarshalText(data[:8]); err != nil {
		return err
	}
	// transform length from number character range
	if len(data) == 10 {
		prefix.Length = data[9] - '0'
	} else {
		prefix.Length = (data[9]-'0')*10 + (data[10] - '0')
	}
	return nil
}

// Value implements driver.Valuer interface.
func (prefix DevAddrPrefix) Value() (driver.Value, error) {
	return prefix.MarshalText()
}

// Scan implements sql.Scanner interface.
func (prefix *DevAddrPrefix) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return ErrTypeAssertion
	}
	return prefix.UnmarshalText(data)
}

// NwkID of the DevAddr.
func (addr DevAddr) NwkID() byte {
	return addr[0] >> 1
}

// WithPrefix returns the DevAddr, but with the first length bits replaced by the Prefix
func (addr DevAddr) WithPrefix(prefix DevAddrPrefix) (prefixed DevAddr) {
	k := uint(prefix.Length)
	for i := 0; i < 4; i++ {
		if k >= 8 {
			prefixed[i] = prefix.DevAddr[i] & 0xff
			k -= 8
			continue
		}
		prefixed[i] = (prefix.DevAddr[i] & ^byte(0xff>>k)) | (addr[i] & byte(0xff>>k))
		k = 0
	}
	return
}

// Mask returns a copy of the DevAddr with only the first "bits" bits.
func (addr DevAddr) Mask(bits uint8) (masked DevAddr) {
	return (DevAddr{}).WithPrefix(DevAddrPrefix{addr, bits})
}

// HasPrefix returns true iff the DevAddr has a prefix of given length.
func (addr DevAddr) HasPrefix(prefix DevAddrPrefix) bool { return prefix.Matches(addr) }

// Matches returns true iff the DevAddr matches the prefix.
func (prefix DevAddrPrefix) Matches(addr DevAddr) bool {
	return addr.Mask(prefix.Length) == prefix.DevAddr.Mask(prefix.Length)
}
