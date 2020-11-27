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

package grpc_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	. "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/grpc"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/mock"
	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
)

var (
	registeredApplicationID  = ttnpb.ApplicationIdentifiers{ApplicationID: "test-app"}
	registeredApplicationUID = unique.ID(test.Context(), registeredApplicationID)
	registeredApplicationKey = "test-key"

	timeout = (1 << 6) * test.Delay
)

func TestAuthentication(t *testing.T) {
	ctx := log.NewContext(test.Context(), test.GetLogger(t))

	is, isAddr := startMockIS(ctx)
	is.add(ctx, registeredApplicationID, registeredApplicationKey)

	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	as := mock.NewServer(c)
	srv := New(as)
	c.RegisterGRPC(&mockRegisterer{ctx, srv})
	componenttest.StartComponent(t, c)
	defer c.Close()

	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)

	client := ttnpb.NewAppAsClient(c.LoopbackConn())

	for _, tc := range []struct {
		ID  ttnpb.ApplicationIdentifiers
		Key string
		OK  bool
	}{
		{
			ID:  registeredApplicationID,
			Key: registeredApplicationKey,
			OK:  true,
		},
		{
			ID:  registeredApplicationID,
			Key: "invalid-key",
			OK:  false,
		},
		{
			ID:  ttnpb.ApplicationIdentifiers{ApplicationID: "invalid-application"},
			Key: "invalid-key",
			OK:  false,
		},
	} {
		t.Run(fmt.Sprintf("%v:%v", tc.ID.ApplicationID, tc.Key), func(t *testing.T) {
			a := assertions.New(t)

			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			creds := grpc.PerRPCCredentials(rpcmetadata.MD{
				AuthType:      "Bearer",
				AuthValue:     tc.Key,
				AllowInsecure: true,
			})

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := client.Subscribe(ctx, &tc.ID, creds)
				if tc.OK && err != nil && !a.So(errors.IsCanceled(err), should.BeTrue) {
					t.Fatalf("Unexpected link error: %v", err)
				}
				if !tc.OK && !a.So(errors.IsCanceled(err), should.BeFalse) {
					t.FailNow()
				}
			}()

			wg.Wait()
		})
	}
}

type erroredApplicationUp struct {
	*ttnpb.ApplicationUp
	error
}

func TestTraffic(t *testing.T) {
	ctx := log.NewContext(test.Context(), test.GetLogger(t))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	is, isAddr := startMockIS(ctx)
	is.add(ctx, registeredApplicationID, registeredApplicationKey)

	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	as := mock.NewServer(c)
	srv := New(as)
	c.RegisterGRPC(&mockRegisterer{ctx, srv})
	componenttest.StartComponent(t, c)
	defer c.Close()

	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)

	client := ttnpb.NewAppAsClient(c.LoopbackConn())

	creds := grpc.PerRPCCredentials(rpcmetadata.MD{
		AuthType:      "Bearer",
		AuthValue:     registeredApplicationKey,
		AllowInsecure: true,
	})
	badCreds := grpc.PerRPCCredentials(rpcmetadata.MD{
		AuthType:      "Bearer",
		AuthValue:     "barfoo",
		AllowInsecure: true,
	})

	upCh := make(chan erroredApplicationUp, 10)
	stream, err := client.Subscribe(ctx, &registeredApplicationID, creds)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}
	go func() {
		for ctx.Err() == nil {
			up, err := stream.Recv()
			upCh <- erroredApplicationUp{up, err}
		}
	}()

	var sub *io.Subscription
	select {
	case sub = <-as.Subscriptions():
	case <-time.After(timeout):
		t.Fatal("Subscription timeout")
	}

	t.Run("Upstream", func(t *testing.T) {
		a := assertions.New(t)

		up := &ttnpb.ApplicationUp{
			EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
				ApplicationIdentifiers: registeredApplicationID,
				DeviceID:               "foo-device",
			},
			Up: &ttnpb.ApplicationUp_UplinkMessage{
				UplinkMessage: &ttnpb.ApplicationUplink{
					FRMPayload: []byte{0x01, 0x02, 0x03},
				},
			},
		}
		if err := sub.Publish(ctx, up); !a.So(err, should.BeNil) {
			t.FailNow()
		}

		select {
		case actual := <-upCh:
			a.So(actual.ApplicationUp, should.Resemble, up)
			a.So(actual.error, should.BeNil)
		case <-time.After(timeout):
			t.Fatal("Receive expected upstream message timeout")
		}
	})

	t.Run("Downstream", func(t *testing.T) {
		a := assertions.New(t)
		ids := ttnpb.EndDeviceIdentifiers{
			ApplicationIdentifiers: registeredApplicationID,
			DeviceID:               "foo-device",
		}

		// List: unauthorized.
		{
			_, err := client.DownlinkQueueList(ctx, &ids, badCreds)
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		// List: happy flow; no items.
		{
			res, err := client.DownlinkQueueList(ctx, &ids, creds)
			a.So(err, should.BeNil)
			a.So(res.Downlinks, should.HaveLength, 0)
		}

		// Push: unauthorized.
		{
			_, err := client.DownlinkQueuePush(ctx, &ttnpb.DownlinkQueueRequest{
				EndDeviceIdentifiers: ids,
				Downlinks: []*ttnpb.ApplicationDownlink{
					{
						FPort:      1,
						FRMPayload: []byte{0x01, 0x01, 0x01},
					},
				},
			}, badCreds)
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		// Push and assert content: happy flow.
		{
			_, err := client.DownlinkQueuePush(ctx, &ttnpb.DownlinkQueueRequest{
				EndDeviceIdentifiers: ids,
				Downlinks: []*ttnpb.ApplicationDownlink{
					{
						SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44}, // This gets discarded.
						FPort:          1,
						FRMPayload:     []byte{0x01, 0x01, 0x01},
						Confirmed:      true,
						CorrelationIDs: []string{"test"},
					},
					{
						FPort:      2,
						FRMPayload: []byte{0x02, 0x02, 0x02},
					},
				},
			}, creds)
			a.So(err, should.BeNil)
		}
		{
			_, err := client.DownlinkQueuePush(ctx, &ttnpb.DownlinkQueueRequest{
				EndDeviceIdentifiers: ids,
				Downlinks: []*ttnpb.ApplicationDownlink{
					{
						FPort:      3,
						FRMPayload: []byte{0x03, 0x03, 0x03},
					},
				},
			}, creds)
			a.So(err, should.BeNil)
		}
		{
			res, err := client.DownlinkQueueList(ctx, &ids, creds)
			a.So(err, should.BeNil)
			a.So(res.Downlinks, should.HaveLength, 3)
			a.So(res.Downlinks, should.Resemble, []*ttnpb.ApplicationDownlink{
				{
					FPort:          1,
					Confirmed:      true,
					FRMPayload:     []byte{0x01, 0x01, 0x01},
					CorrelationIDs: []string{"test"},
				},
				{
					FPort:      2,
					FRMPayload: []byte{0x02, 0x02, 0x02},
				},
				{
					FPort:      3,
					FRMPayload: []byte{0x03, 0x03, 0x03},
				},
			})
		}

		// Replace: unauthorized.
		{
			_, err := client.DownlinkQueueReplace(ctx, &ttnpb.DownlinkQueueRequest{
				EndDeviceIdentifiers: ids,
				Downlinks: []*ttnpb.ApplicationDownlink{
					{
						FPort:      4,
						FRMPayload: []byte{0x04, 0x04, 0x04},
					},
				},
			}, badCreds)
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		// Replace and assert content: happy flow.
		{
			_, err := client.DownlinkQueueReplace(ctx, &ttnpb.DownlinkQueueRequest{
				EndDeviceIdentifiers: ids,
				Downlinks: []*ttnpb.ApplicationDownlink{
					{
						FPort:      4,
						FRMPayload: []byte{0x04, 0x04, 0x04},
						Confirmed:  true,
					},
				},
			}, creds)
			a.So(err, should.BeNil)
		}
		{
			res, err := client.DownlinkQueueList(ctx, &ids, creds)
			a.So(err, should.BeNil)
			a.So(res.Downlinks, should.HaveLength, 1)
			a.So(res.Downlinks, should.Resemble, []*ttnpb.ApplicationDownlink{
				{
					FPort:      4,
					FRMPayload: []byte{0x04, 0x04, 0x04},
					Confirmed:  true,
				},
			})
		}
	})
}

type mockMQTTConfigProvider struct {
	config.MQTT
}

func (p mockMQTTConfigProvider) GetMQTTConfig(context.Context) (*config.MQTT, error) {
	return &p.MQTT, nil
}

func TestMQTTConfig(t *testing.T) {
	a := assertions.New(t)
	ctx := log.NewContext(test.Context(), test.GetLogger(t))

	is, isAddr := startMockIS(ctx)
	is.add(ctx, registeredApplicationID, registeredApplicationKey)

	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	as := mock.NewServer(c)
	srv := New(as, WithMQTTConfigProvider(&mockMQTTConfigProvider{
		MQTT: config.MQTT{
			PublicAddress:    "example.com:1883",
			PublicTLSAddress: "example.com:8883",
		},
	}))
	c.RegisterGRPC(&mockRegisterer{ctx, srv})
	componenttest.StartComponent(t, c)
	defer c.Close()

	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)

	client := ttnpb.NewAppAsClient(c.LoopbackConn())

	creds := grpc.PerRPCCredentials(rpcmetadata.MD{
		AuthType:      "Bearer",
		AuthValue:     registeredApplicationKey,
		AllowInsecure: true,
	})

	info, err := client.GetMQTTConnectionInfo(ctx, &registeredApplicationID, creds)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(info, should.Resemble, &ttnpb.MQTTConnectionInfo{
		Username:         registeredApplicationUID,
		PublicAddress:    "example.com:1883",
		PublicTLSAddress: "example.com:8883",
	})
}
