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

package provisioning_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/joinserver/provisioning"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func TestMicrochip(t *testing.T) {
	a := assertions.New(t)
	data := []byte(`[{
		"payload": "eyJ2ZXJzaW9uIjogMSwgIm1vZGVsIjogIkFURUNDNjA4QSIsICJwYXJ0TnVtYmVyIjogIkFURUNDNjA4QS1NQUhUMiIsICJtYW51ZmFjdHVyZXIiOiB7Im9yZ2FuaXphdGlvbk5hbWUiOiAiTWljcm9jaGlwIFRlY2hub2xvZ3kgSW5jIiwgIm9yZ2FuaXphdGlvbmFsVW5pdE5hbWUiOiAiU2VjdXJlIFByb2R1Y3RzIEdyb3VwIn0sICJwcm92aXNpb25lciI6IHsib3JnYW5pemF0aW9uTmFtZSI6ICJNaWNyb2NoaXAgVGVjaG5vbG9neSBJbmMiLCAib3JnYW5pemF0aW9uYWxVbml0TmFtZSI6ICJTZWN1cmUgUHJvZHVjdHMgR3JvdXAifSwgImRpc3RyaWJ1dG9yIjogeyJvcmdhbml6YXRpb25OYW1lIjogIk1pY3JvY2hpcCBUZWNobm9sb2d5IEluYyIsICJvcmdhbml6YXRpb25hbFVuaXROYW1lIjogIk1pY3JvY2hpcCBEaXJlY3QifSwgImdyb3VwSWQiOiAiSkVFTERBSldTUDJKQUdBMyIsICJwcm92aXNpb25pbmdUaW1lc3RhbXAiOiAiMjAxOS0wMS0xNFQxODozMjoyMS44NjdaIiwgInVuaXF1ZUlkIjogIjAxMjM3YTAwNWIwOGJjYzUyNyIsICJwdWJsaWNLZXlTZXQiOiB7ImtleXMiOiBbeyJraWQiOiAiMSIsICJrdHkiOiAiRUMiLCAiY3J2IjogIlAtMjU2IiwgIngiOiAienNTNjNUMnRUdXJTT2E1dmRVamlkU0U2NVJGc3VYOERjNDYwQkpmVE1nND0iLCAieSI6ICI2MmZIa3E1MzVWWE5Ubnc0ZXUxeDRhYl9fM3daRHkyVUh0Q3I3WkpzNmowPSIsICJ4NWMiOiBbIk1JSUI5VENDQVp1Z0F3SUJBZ0lRZE5OdXJPVi9UNTEycnpIb2t2ZlEyekFLQmdncWhrak9QUVFEQWpCUE1TRXdId1lEVlFRS0RCaE5hV055YjJOb2FYQWdWR1ZqYUc1dmJHOW5lU0JKYm1NeEtqQW9CZ05WQkFNTUlVTnllWEIwYnlCQmRYUm9aVzUwYVdOaGRHbHZiaUJUYVdkdVpYSWdSall3TVRBZ0Z3MHhPVEF4TVRReE9EQXdNREJhR0E4eU1EUTNNREV4TkRFNE1EQXdNRm93UmpFaE1COEdBMVVFQ2d3WVRXbGpjbTlqYUdsd0lGUmxZMmh1YjJ4dloza2dTVzVqTVNFd0h3WURWUVFEREJnd01USXpOMEV3TURWQ01EaENRME0xTWpjZ1FWUkZRME13V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVRPeExyZFBhMU82dEk1cm05MVNPSjFJVHJsRVd5NWZ3TnpqclFFbDlNeUR1dG54NUt1ZCtWVnpVNThPSHJ0Y2VHbS8vOThHUTh0bEI3UXErMlNiT285bzJBd1hqQU1CZ05WSFJNQkFmOEVBakFBTUE0R0ExVWREd0VCL3dRRUF3SURpREFkQmdOVkhRNEVGZ1FVNzBNUy9CdjVMK2gvUWdKbkp6a3A4b0phVy84d0h3WURWUjBqQkJnd0ZvQVU4VE1wdElQSDdkWnVRZGtRbHl3cytCbXR5dWd3Q2dZSUtvWkl6ajBFQXdJRFNBQXdSUUlnTFpML2lJSElTaE1mMXRTSkFkRTRLcXI4Tlp0QmllSVhleG1QS2F1cEtzb0NJUUNYZWtQUTdCQm1BOHpiR0NjVnhDRTM5ZEdnQTNsSGhFMnRDWlp1Y04zS0tBPT0iLCAiTUlJQ0JEQ0NBYXFnQXdJQkFnSVFhc2ExbEttdzR1WG5haEdQNXdCZEFEQUtCZ2dxaGtqT1BRUURBakJQTVNFd0h3WURWUVFLREJoTmFXTnliMk5vYVhBZ1ZHVmphRzV2Ykc5bmVTQkpibU14S2pBb0JnTlZCQU1NSVVOeWVYQjBieUJCZFhSb1pXNTBhV05oZEdsdmJpQlNiMjkwSUVOQklEQXdNakFnRncweE9ERXlNVFF4T1RBd01EQmFHQTh5TURRNU1USXhOREU1TURBd01Gb3dUekVoTUI4R0ExVUVDZ3dZVFdsamNtOWphR2x3SUZSbFkyaHViMnh2WjNrZ1NXNWpNU293S0FZRFZRUUREQ0ZEY25sd2RHOGdRWFYwYUdWdWRHbGpZWFJwYjI0Z1UybG5ibVZ5SUVZMk1ERXdXVEFUQmdjcWhrak9QUUlCQmdncWhrak9QUU1CQndOQ0FBVFc5QzVRTU5PT2VKNUZIMkU1THhqdHpNVjlQbXNLL1R6R1hwQmswVm9NYzhsSVZmWkJTWUhxUkh6QXpoRWF5RTRBem1LZDhnbGJGbFdYYTNXRWhON0NvMll3WkRBT0JnTlZIUThCQWY4RUJBTUNBWVl3RWdZRFZSMFRBUUgvQkFnd0JnRUIvd0lCQURBZEJnTlZIUTRFRmdRVThUTXB0SVBIN2RadVFka1FseXdzK0JtdHl1Z3dId1lEVlIwakJCZ3dGb0FVZXUxOWJjYTNlSjJ5T0FHbDZFcU1zS1FPS293d0NnWUlLb1pJemowRUF3SURTQUF3UlFJaEFPYnk0N2pjejJHT2Q3M0NPQzVjQXJrNnlCNDRyd2hJeWFQbHJEMERTU2FVQWlCakU1RGkzMXgrc0RxY3FvR09jaFNhMmJlcnF4Q0taU1dyRTExUWd4ZjBJdz09Il19XX19",
		"protected": "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiIsImtpZCI6IkIwdWFETExpeUtlLVNvclA3MUYzQk5vVnJmWSIsIng1dCNTMjU2IjoiVmFURzFGWjV4Z25CUXV4U2FkV09iWFRIMmhQUkxqSFQtNER6VmFKTTNvbyJ9",
		"signature": "rOrjAQyI1FQO-RmNjS0ggA7t4U7l9epwCbTQYnyb92j2EOqZ3vNhF5RAMldvG68v6w25mKrGaPG8wdn--TSy4w"
	}]`)

	provisioner := provisioning.Get("microchip")
	entries, err := provisioner.Decode(data)
	a.So(err, should.BeNil)
	if !a.So(entries, should.HaveLength, 1) {
		t.FailNow()
	}
	entry := entries[0]
	a.So(entry.Fields["uniqueId"].GetStringValue(), should.Equal, "01237a005b08bcc527")

	joinEUI, err := provisioner.DefaultJoinEUI(entry)
	a.So(err, should.BeNil)
	a.So(joinEUI, should.Resemble, types.EUI64{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x00, 0x00, 0x00})

	devEUI, err := provisioner.DefaultDevEUI(entry)
	a.So(err, should.BeNil)
	a.So(devEUI, should.Resemble, types.EUI64{0x01, 0x23, 0x7a, 0x00, 0x5b, 0x08, 0xbc, 0xc5})

	deviceID, err := provisioner.DeviceID(joinEUI, devEUI, entry)
	a.So(err, should.BeNil)
	a.So(deviceID, should.Equal, "sn-01237a005b08bcc527")

	uniqueID, err := provisioner.UniqueID(entry)
	a.So(err, should.BeNil)
	a.So(uniqueID, should.Equal, "01237A005B08BCC527")
}
