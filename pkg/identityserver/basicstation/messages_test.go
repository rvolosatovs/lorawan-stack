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

package basicstation_test

import (
	"bytes"
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/identityserver/basicstation"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func TestUpdateInfoResponse(t *testing.T) {
	for _, tt := range []struct {
		Name string
		basicstation.UpdateInfoResponse
	}{
		{Name: "Empty"},
		{Name: "Full", UpdateInfoResponse: basicstation.UpdateInfoResponse{
			CUPSURI:         "https://cups.example.com",
			LNSURI:          "https://lns.example.com",
			CUPSCredentials: bytes.Repeat([]byte("CUPS CREDENTIALS"), 1000),
			LNSCredentials:  bytes.Repeat([]byte("LNS CREDENTIALS"), 1000),
			SignatureKeyCRC: 12345678,
			Signature:       bytes.Repeat([]byte("THIS IS THE SIGNATURE"), 100),
			UpdateData:      bytes.Repeat([]byte("THIS IS THE UPDATE DATA"), 1000),
		}},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			a := assertions.New(t)

			data, err := tt.UpdateInfoResponse.MarshalBinary()
			a.So(err, should.BeNil)

			var dec basicstation.UpdateInfoResponse
			err = dec.UnmarshalBinary(data)
			a.So(err, should.BeNil)
			a.So(dec, should.Resemble, tt.UpdateInfoResponse)
		})
	}
}
