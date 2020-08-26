// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

func TestLegacyEncodeDownlink(t *testing.T) {
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
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
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
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{247, 174})
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
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}
}

func TestEncodeDownlink(t *testing.T) {
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
		function encodeDownlink(input) {
			return {
				bytes: [1, 2, 3],
				fPort: 42
			}
		}
		`
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{1, 2, 3})
		a.So(message.FPort, should.Equal, 42)
	}

	// Encode temperature.
	{
		script := `
		function encodeDownlink(input) {
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
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{247, 174})
	}

	// The Things Node example.
	{
		message := &ttnpb.ApplicationDownlink{
			DecodedPayload: &pbtypes.Struct{
				Fields: map[string]*pbtypes.Value{
					"color": {
						Kind: &pbtypes.Value_StringValue{
							StringValue: "blue",
						},
					},
				},
			},
		}
		script := `
		function encodeDownlink(input) {
			var colors = ["red", "green", "blue"];
			return {
				bytes: [colors.indexOf(input.data.color)],
				fPort: 4,
			}
		}
		`
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		a.So(message.FRMPayload, should.Resemble, []byte{0x2})
		a.So(message.FPort, should.Equal, 4)
	}

	// Return nothing.
	{
		script := `
		function encodeDownlink(input) {
			return null
		}
		`
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Return undefined.
	{
		script := `
		function encodeDownlink(input) {
			return undefined
		}
		`
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Return errors.
	{
		script := `
		function encodeDownlink(input) {
			return {
				bytes: [1, 2, 3],
				errors: ["error 1", "error 2"]
			}
		}
		`
		err := host.EncodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutputErrors.WithAttributes("errors", "error 1, error 2"))
	}
}

func TestLegacyDecodeUplink(t *testing.T) {
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
		err := host.DecodeUplink(ctx, ids, nil, message, script)
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
		err := host.DecodeUplink(ctx, ids, nil, message, script)
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
		err := host.DecodeUplink(ctx, ids, nil, message, script)
		a.So(err, should.NotBeNil)
	}

	// Catch error.
	{
		script := `
		function Decoder(payload, f_port) {
			throw Error('unknown error')
		}
		`
		err := host.DecodeUplink(ctx, ids, nil, message, script)
		a.So(err, should.NotBeNil)
	}
}

func TestDecodeUplink(t *testing.T) {
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

	// Decode bytes.
	{
		script := `
		function decodeUplink(input) {
			return {
				data: {
					temperature: (((input.bytes[0] & 0x80 ? input.bytes[0] - 0x100 : input.bytes[0]) << 8) | input.bytes[1]) / 100
				}
			}
		}
		`
		err := host.DecodeUplink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		m, err := gogoproto.Map(message.DecodedPayload)
		a.So(err, should.BeNil)
		a.So(m, should.Resemble, map[string]interface{}{
			"temperature": -21.3,
		})
	}

	// The Things Node example.
	{
		message := &ttnpb.ApplicationUplink{
			FPort:      4,
			FRMPayload: []byte{0x0C, 0xB2, 0x04, 0x80, 0xF7, 0xAE},
		}
		script := `
		function decodeUplink(input) {
			var data = {};
			var events = {
				1: 'setup',
				2: 'interval',
				3: 'motion',
				4: 'button'
			};
			data.event = events[input.fPort];
			data.battery = (input.bytes[0] << 8) + input.bytes[1];
			data.light = (input.bytes[2] << 8) + input.bytes[3];
			data.temperature = (((input.bytes[4] & 0x80 ? input.bytes[4] - 0x100 : input.bytes[4]) << 8) + input.bytes[5]) / 100;
			return {
				data: data
			};
		}
		`
		err := host.DecodeUplink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		m, err := gogoproto.Map(message.DecodedPayload)
		a.So(err, should.BeNil)
		a.So(m, should.Resemble, map[string]interface{}{
			"event":       "button",
			"battery":     3250.0,
			"light":       1152.0,
			"temperature": -21.3,
		})
	}

	// Return invalid type.
	{
		script := `
		function decodeUplink(input) {
			return {
				data: 42
			}
		}
		`
		err := host.DecodeUplink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Catch error.
	{
		script := `
		function decodeUplink(input) {
			throw Error('unknown error')
		}
		`
		err := host.DecodeUplink(ctx, ids, nil, message, script)
		a.So(err, should.NotBeNil)
	}

	// Return errors.
	{
		script := `
		function decodeUplink(input) {
			return {
				errors: ["error 1", "error 2"]
			}
		}
		`
		err := host.DecodeUplink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutputErrors.WithAttributes("errors", "error 1, error 2"))
	}
}

func TestDecodeDownlink(t *testing.T) {
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
		FRMPayload: []byte{0xF7, 0xAE},
		FPort:      4,
	}

	// Decode bytes.
	{
		script := `
		function decodeDownlink(input) {
			return {
				data: {
					value: (((input.bytes[0] & 0x80 ? input.bytes[0] - 0x100 : input.bytes[0]) << 8) | input.bytes[1]) / 100
				}
			}
		}
		`
		err := host.DecodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		m, err := gogoproto.Map(message.DecodedPayload)
		a.So(err, should.BeNil)
		a.So(m, should.Resemble, map[string]interface{}{
			"value": -21.3,
		})
	}

	// The Things Node example.
	{
		message := &ttnpb.ApplicationDownlink{
			FPort:      4,
			FRMPayload: []byte{0x02},
		}
		script := `
		function decodeDownlink(input) {
			switch (input.fPort) {
			case 4:
				return {
					data: {
						color: ["red", "green", "blue"][input.bytes[0]]
					}
				}
			default:
				return {
					errors: ["unknown FPort"]
				}
			}
		}
		`
		err := host.DecodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.BeNil)
		m, err := gogoproto.Map(message.DecodedPayload)
		a.So(err, should.BeNil)
		a.So(m, should.Resemble, map[string]interface{}{
			"color": "blue",
		})
	}

	// Return invalid type.
	{
		script := `
		function decodeDownlink(input) {
			return {
				data: 42
			}
		}
		`
		err := host.DecodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutput)
	}

	// Catch error.
	{
		script := `
		function decodeDownlink(input) {
			throw Error('unknown error')
		}
		`
		err := host.DecodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.NotBeNil)
	}

	// Return errors.
	{
		script := `
		function decodeDownlink(input) {
			return {
				errors: ["error 1", "error 2"]
			}
		}
		`
		err := host.DecodeDownlink(ctx, ids, nil, message, script)
		a.So(err, should.HaveSameErrorDefinitionAs, errOutputErrors.WithAttributes("errors", "error 1, error 2"))
	}
}
