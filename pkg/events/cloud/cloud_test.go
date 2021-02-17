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

package cloud_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/events/cloud"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	_ "gocloud.dev/pubsub/mempubsub"
)

func Example() {
	// The task starter is used for automatic re-subscription on failure.
	taskStarter := component.StartTaskFunc(component.DefaultStartTask)

	// Import the desired cloud pub-sub drivers (see godoc.org/gocloud.dev).
	// In this example we use "gocloud.dev/pubsub/mempubsub".

	cloudPubSub, err := cloud.NewPubSub(context.TODO(), taskStarter, "mem://events", "mem://events")
	if err != nil {
		// Handle error.
	}

	// Replace the default pubsub so that we will now publish to a Go Cloud pub sub.
	events.SetDefaultPubSub(cloudPubSub)
}

var timeout = (1 << 10) * test.Delay

func TestCloudPubSub(t *testing.T) {
	test.RunTest(t, test.TestConfig{
		Timeout: 4 * timeout,
		Func: func(ctx context.Context, a *assertions.Assertion) {
			events.IncludeCaller = true

			eventCh := make(chan events.Event)
			handler := events.HandlerFunc(func(e events.Event) {
				t.Logf("Received event %v", e)
				a.So(e.Time().IsZero(), should.BeFalse)
				a.So(e.Context(), should.NotBeNil)
				eventCh <- e
			})

			taskStarter := component.StartTaskFunc(component.DefaultStartTask)
			pubsub, err := cloud.NewPubSub(ctx, taskStarter, "mem://events_test", "mem://events_test")
			a.So(err, should.BeNil)

			defer pubsub.Close(ctx)

			subCtx, unsubscribe := context.WithCancel(ctx)
			defer unsubscribe()
			pubsub.Subscribe(subCtx, "cloud.**", nil, handler)

			time.Sleep(timeout)

			ctx = events.ContextWithCorrelationID(ctx, t.Name())

			eui := types.EUI64{1, 2, 3, 4, 5, 6, 7, 8}
			devAddr := types.DevAddr{1, 2, 3, 4}
			appID := ttnpb.ApplicationIdentifiers{
				ApplicationID: "test-app",
			}
			devID := ttnpb.EndDeviceIdentifiers{
				ApplicationIdentifiers: appID,
				DeviceID:               "test-dev",
				DevEUI:                 &eui,
				JoinEUI:                &eui,
				DevAddr:                &devAddr,
			}
			gtwID := ttnpb.GatewayIdentifiers{
				GatewayID: "test-gtw",
				EUI:       &eui,
			}

			test.RunSubtestFromContext(ctx, test.SubtestConfig{
				Name:    "publish_json",
				Timeout: timeout,
				Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
					cloud.SetContentType(pubsub, "application/json")

					pubsub.Publish(events.New(ctx, "cloud.test.evt0", "cloud test event 0", events.WithIdentifiers(appID)))
					select {
					case e := <-eventCh:
						a.So(e.Name(), should.Equal, "cloud.test.evt0")
						if a.So(e.Identifiers(), should.NotBeNil) && a.So(e.Identifiers(), should.HaveLength, 1) {
							a.So(e.Identifiers()[0].GetApplicationIDs(), should.Resemble, &appID)
						}
					case <-ctx.Done():
						t.Error("Did not receive expected event")
						t.FailNow()
					}
				},
			})

			test.RunSubtestFromContext(ctx, test.SubtestConfig{
				Name:    "publish_pb",
				Timeout: timeout,
				Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
					cloud.SetContentType(pubsub, "application/protobuf")

					pubsub.Publish(events.New(ctx, "cloud.test.evt1", "cloud test event 1", events.WithIdentifiers(&devID, &gtwID)))
					select {
					case e := <-eventCh:
						a.So(e.Name(), should.Equal, "cloud.test.evt1")
						if a.So(e.Identifiers(), should.NotBeNil) && a.So(e.Identifiers(), should.HaveLength, 2) {
							a.So(e.Identifiers()[0].GetDeviceIDs(), should.Resemble, &devID)
							a.So(e.Identifiers()[1].GetGatewayIDs(), should.Resemble, &gtwID)
						}
					case <-ctx.Done():
						t.Error("Did not receive expected event")
						t.FailNow()
					}
				},
			})
		},
	})
}
