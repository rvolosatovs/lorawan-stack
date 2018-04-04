// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
