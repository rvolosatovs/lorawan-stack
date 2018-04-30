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

package gatewayserver_test

import (
	"context"
	"testing"

	"github.com/TheThingsNetwork/ttn/pkg/component"
	"github.com/TheThingsNetwork/ttn/pkg/config"
	"github.com/TheThingsNetwork/ttn/pkg/gatewayserver"
	"github.com/TheThingsNetwork/ttn/pkg/gatewayserver/pool"
	"github.com/TheThingsNetwork/ttn/pkg/log"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/TheThingsNetwork/ttn/pkg/util/test"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestScheduleDownlinkUnregisteredGateway(t *testing.T) {
	a := assertions.New(t)

	store, err := test.NewFrequencyPlansStore()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	defer store.Destroy()

	logger := test.GetLogger(t)
	c := component.MustNew(test.GetLogger(t), &component.Config{ServiceBase: config.ServiceBase{
		FrequencyPlans: config.FrequencyPlans{StoreDirectory: store.Directory()},
	}})
	gs, err := gatewayserver.New(c, gatewayserver.Config{})
	if !a.So(err, should.BeNil) {
		logger.Fatal("Gateway server could not start")
	}

	_, err = gs.ScheduleDownlink(log.NewContext(context.Background(), logger), &ttnpb.DownlinkMessage{
		TxMetadata: ttnpb.TxMetadata{
			GatewayIdentifiers: ttnpb.GatewayIdentifiers{
				GatewayID: "unknown-downlink",
			},
		},
	})
	a.So(err, should.NotBeNil)
	a.So(pool.ErrGatewayNotConnected.Caused(err), should.BeTrue)

	defer gs.Close()
}
