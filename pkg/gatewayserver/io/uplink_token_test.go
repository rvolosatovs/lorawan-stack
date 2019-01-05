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

package io_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func TestUplinkToken(t *testing.T) {
	a := assertions.New(t)

	inIDs := ttnpb.GatewayAntennaIdentifiers{
		GatewayIdentifiers: ttnpb.GatewayIdentifiers{
			GatewayID: "foo-gateway",
		},
		AntennaIndex: 0,
	}
	inTimestamp := uint32(12345678)

	uplinkToken, err := io.UplinkToken(inIDs, inTimestamp)
	a.So(err, should.BeNil)

	outIDs, outTimestamp, err := io.ParseUplinkToken(uplinkToken)
	a.So(err, should.BeNil)
	a.So(outIDs, should.Resemble, inIDs)
	a.So(outTimestamp, should.Equal, inTimestamp)
}
