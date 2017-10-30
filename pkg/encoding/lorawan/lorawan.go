// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

// Package lorawan provides LoRaWAN decoding/encoding interfaces.
package lorawan

// Marshaler is the interface implemented by an object that can
// marshal itself into LoRaWAN form.
//
// MarshalLoRaWAN encodes the receiver into LoRaWAN form and returns the result.
type Marshaler interface {
	MarshalLoRaWAN() (data []byte, err error)
}

// Appender is the interface implemented by an object that can append LoRaWAN representation of itself to a byte slice.
//
// AppendLoRaWAN encodes the receiver into LoRaWAN form, appends it to supplied byte slice and returns the result.
type Appender interface {
	AppendLoRaWAN(dst []byte) (data []byte, err error)
}

// Unmarshaler is the interface implemented by an object that can
// unmarshal a LoRaWAN representation of itself.
//
// UnmarshalLoRaWAN must be able to decode the form generated by MarshalLoRaWAN.
// UnmarshalLoRaWAN must copy the data if it wishes to retain the data
// after returning.
type Unmarshaler interface {
	UnmarshalLoRaWAN(data []byte) error
}
