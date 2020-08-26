// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package javascript

import (
	"testing"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/gogoproto"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestLegacyEncode(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()
	host := New()

	eui := types.EUI64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
			ApplicationID: "foo-app",
		},
		DeviceID: "foo-device",
		DevEUI:   &eui,
	}

	message := &ttnpb.ApplicationDownlink{
		DecodedPayload: &pbtypes.Struct{
			Fields: map[string]*pbtypes.Value{
				"temperature": {
					Kind: &pbtypes.Value_NumberValue{
						NumberValue: -21.3,
					},
				},
			},
		},
	}

	// Return constant byte array.
	{
		script := `
		function Encoder(payload, f_port) {
			return [1, 2, 3]
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{1, 2, 3})
	}

	// Encode temperature.
	{
		script := `
		function Encoder(payload, f_port) {
			var val = payload.temperature * 100
			return [
				(val >> 8) & 0xff,
				val & 0xff
			]
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{247, 174})
	}

	// Return nothing.
	{
		script := `
		function Encoder(payload, f_port) {
			return null
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Return an object.
	{
		script := `
		function Encoder(payload, f_port) {
			return {
				value: 42
			}
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}
}

func TestEncode(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()
	host := New()

	eui := types.EUI64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
			ApplicationID: "foo-app",
		},
		DeviceID: "foo-device",
		DevEUI:   &eui,
	}

	message := &ttnpb.ApplicationDownlink{
		DecodedPayload: &pbtypes.Struct{
			Fields: map[string]*pbtypes.Value{
				"temperature": {
					Kind: &pbtypes.Value_NumberValue{
						NumberValue: -21.3,
					},
				},
			},
		},
		FPort: 2,
	}

	// Return constant byte array and FPort.
	{
		script := `
		function encode(input) {
			return {
				bytes: [1, 2, 3],
				fPort: 42
			}
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{1, 2, 3})
		a.So(message.FPort, should.Equal, 42)
	}

	// Encode temperature.
	{
		script := `
		function encode(input) {
			var val = input.data.temperature * 100
			return {
				bytes: [
					(val >> 8) & 0xff,
					val & 0xff
				],
				fPort: input.fPort
			}
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{247, 174})
	}

	// Return nothing.
	{
		script := `
		function encode(input) {
			return null
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Return undefined.
	{
		script := `
		function encode(input) {
			return undefined
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Return errors.
	{
		script := `
		function encode(input) {
			return {
				bytes: [1, 2, 3],
				errors: ["error 1", "error 2"]
			}
		}
		`
		err := host.Encode(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutputErrors.WithAttributes("errors", "error 1, error 2"))
	}
}

func TestLegacyDecode(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()
	host := New()

	eui := types.EUI64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
			ApplicationID: "foo-app",
		},
		DeviceID: "foo-device",
		DevEUI:   &eui,
	}

	message := &ttnpb.ApplicationUplink{
		FRMPayload: []byte{0xF7, 0xAE},
	}

	// Return constant object.
	{
		script := `
		function Decoder(payload, f_port) {
			return {
				temperature: -21.3
			}
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		m, err := gogoproto.Map(message.DecodedPayload)
		a.So(err, should.BeNil)
		a.So(m, should.Resemble, map[string]interface{}{
			"temperature": -21.3,
		})
	}

	// Decode bytes.
	{
		script := `
		function Decoder(payload, f_port) {
			return {
				temperature: (((payload[0] & 0x80 ? payload[0] - 0x100 : payload[0]) << 8) | payload[1]) / 100
			}
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		m, err := gogoproto.Map(message.DecodedPayload)
		a.So(err, should.BeNil)
		a.So(m, should.Resemble, map[string]interface{}{
			"temperature": -21.3,
		})
	}

	// Return invalid type.
	{
		script := `
		function Decoder(payload, f_port) {
			return 42
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.NotBeNil)
	}

	// Catch error.
	{
		script := `
		function Decoder(payload, f_port) {
			throw Error('unknown error')
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.NotBeNil)
	}
}

func TestDecode(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()
	host := New()

	eui := types.EUI64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
			ApplicationID: "foo-app",
		},
		DeviceID: "foo-device",
		DevEUI:   &eui,
	}

	message := &ttnpb.ApplicationUplink{
		FRMPayload: []byte{247, 174},
	}

	// Return constant object.
	{
		script := `
		function decode(input) {
			return {
				data: {
					temperature: (((input.bytes[0] & 0x80 ? input.bytes[0] - 0x100 : input.bytes[0]) << 8) | input.bytes[1]) / 100
				}
			}
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		m, err := gogoproto.Map(message.DecodedPayload)
		a.So(err, should.BeNil)
		a.So(m, should.Resemble, map[string]interface{}{
			"temperature": -21.3,
		})
	}

	// Return invalid type.
	{
		script := `
		function decode(input) {
			return {
				data: 42
			}
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Catch error.
	{
		script := `
		function decode(input) {
			throw Error('unknown error')
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.NotBeNil)
	}

	// Return errors.
	{
		script := `
		function decode(input) {
			return {
				errors: ["error 1", "error 2"]
			}
		}
		`
		err := host.Decode(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutputErrors.WithAttributes("errors", "error 1, error 2"))
	}
}
