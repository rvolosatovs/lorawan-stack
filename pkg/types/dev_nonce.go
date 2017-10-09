// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package types

import (
	"database/sql/driver"
	"encoding/hex"
	"strings"
)

// DevNonce is randomly generated in the join procedure
type DevNonce [2]byte

// IsZero returns true iff the type is zero
func (dn DevNonce) IsZero() bool { return dn == [2]byte{} }

// String implements the Stringer interface
func (dn DevNonce) String() string { return strings.ToUpper(hex.EncodeToString(dn[:])) }

// GoString implements the GoStringer interface
func (dn DevNonce) GoString() string { return dn.String() }

// Size implements the Sizer interface
func (dn DevNonce) Size() int { return 2 }

// Equal returns true iff nonces are equal
func (dn DevNonce) Equal(other DevNonce) bool { return dn == other }

// Marshal implements the proto.Marshaler interface
func (dn DevNonce) Marshal() ([]byte, error) { return dn.MarshalBinary() }

// MarshalTo implements the MarshalerTo function required by generated protobuf
func (dn DevNonce) MarshalTo(data []byte) (int, error) { return marshalBinaryBytesTo(data, dn[:]) }

// Unmarshal implements the proto.Unmarshaler interface
func (dn *DevNonce) Unmarshal(data []byte) error { return dn.UnmarshalBinary(data) }

// MarshalJSON implements the json.Marshaler interface
func (dn DevNonce) MarshalJSON() ([]byte, error) { return marshalJSONBytes(dn[:]) }

// UnmarshalJSON implements the json.Unmarshaler interface
func (dn *DevNonce) UnmarshalJSON(data []byte) error {
	*dn = [2]byte{}
	return unmarshalJSONBytes(dn[:], data)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (dn DevNonce) MarshalBinary() ([]byte, error) { return marshalBinaryBytes(dn[:]) }

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (dn *DevNonce) UnmarshalBinary(data []byte) error {
	*dn = [2]byte{}
	return unmarshalBinaryBytes(dn[:], data)
}

// MarshalText implements the encoding.TextMarshaler interface
func (dn DevNonce) MarshalText() ([]byte, error) { return marshalTextBytes(dn[:]) }

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (dn *DevNonce) UnmarshalText(data []byte) error {
	*dn = [2]byte{}
	return unmarshalTextBytes(dn[:], data)
}

// Value implements driver.Valuer interface.
func (dn DevNonce) Value() (driver.Value, error) {
	return dn.MarshalText()
}

// Scan implements sql.Scanner interface.
func (dn *DevNonce) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return ErrTypeAssertion
	}
	return dn.UnmarshalText(data)
}
