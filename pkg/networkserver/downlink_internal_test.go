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

package networkserver

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/mohae/deepcopy"
	ulid "github.com/oklog/ulid/v2"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/crypto"
	"go.thethings.network/lorawan-stack/v3/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/gpstime"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
)

func TestAppendRecentDownlink(t *testing.T) {
	downs := [...]*ttnpb.DownlinkMessage{
		{
			RawPayload: []byte("test1"),
		},
		{
			RawPayload: []byte("test2"),
		},
		{
			RawPayload: []byte("test3"),
		},
	}
	for _, tc := range []struct {
		Recent   []*ttnpb.DownlinkMessage
		Down     *ttnpb.DownlinkMessage
		Window   int
		Expected []*ttnpb.DownlinkMessage
	}{
		{
			Down:     downs[0],
			Window:   1,
			Expected: downs[:1],
		},
		{
			Recent:   downs[:1],
			Down:     downs[1],
			Window:   1,
			Expected: downs[1:2],
		},
		{
			Recent:   downs[:2],
			Down:     downs[2],
			Window:   1,
			Expected: downs[2:3],
		},
		{
			Recent:   downs[:1],
			Down:     downs[1],
			Window:   2,
			Expected: downs[:2],
		},
		{
			Recent:   downs[:2],
			Down:     downs[2],
			Window:   2,
			Expected: downs[1:3],
		},
	} {
		t.Run(fmt.Sprintf("recent_length:%d,window:%v", len(tc.Recent), tc.Window), func(t *testing.T) {
			a := assertions.New(t)
			recent := CopyDownlinkMessages(tc.Recent...)
			down := CopyDownlinkMessage(tc.Down)
			ret := appendRecentDownlink(recent, down, tc.Window)
			a.So(recent, should.Resemble, tc.Recent)
			a.So(down, should.Resemble, tc.Down)
			a.So(ret, should.Resemble, tc.Expected)
		})
	}
}

func TestProcessDownlinkTask(t *testing.T) {
	// TODO: Refactor. (https://github.com/TheThingsNetwork/lorawan-stack/issues/2475)

	getPaths := []string{
		"frequency_plan_id",
		"last_dev_status_received_at",
		"lorawan_phy_version",
		"mac_settings",
		"mac_state",
		"multicast",
		"pending_mac_state",
		"queued_application_downlinks",
		"recent_downlinks",
		"recent_uplinks",
		"session",
	}

	const appIDString = "process-downlink-test-app-id"
	appID := ttnpb.ApplicationIdentifiers{ApplicationID: appIDString}
	const devID = "process-downlink-test-dev-id"

	devAddr := types.DevAddr{0x42, 0xff, 0xff, 0xff}

	fNwkSIntKey := types.AES128Key{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nwkSEncKey := types.AES128Key{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	sNwkSIntKey := types.AES128Key{0x42, 0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	sessionKeys := &ttnpb.SessionKeys{
		FNwkSIntKey: &ttnpb.KeyEnvelope{
			Key: &fNwkSIntKey,
		},
		NwkSEncKey: &ttnpb.KeyEnvelope{
			Key: &nwkSEncKey,
		},
		SNwkSIntKey: &ttnpb.KeyEnvelope{
			Key: &sNwkSIntKey,
		},
		SessionKeyID: []byte{0x11, 0x22, 0x33, 0x44},
	}

	const rx1Delay = ttnpb.RX_DELAY_1
	makeEU868macParameters := func(ver ttnpb.PHYVersion) ttnpb.MACParameters {
		params := MakeDefaultEU868CurrentMACParameters(ver)
		params.Rx1Delay = rx1Delay
		params.Channels = append(params.Channels, &ttnpb.MACParameters_Channel{
			UplinkFrequency:   430000000,
			DownlinkFrequency: 431000000,
			MinDataRateIndex:  ttnpb.DATA_RATE_0,
			MaxDataRateIndex:  ttnpb.DATA_RATE_3,
		})
		return params
	}

	assertGetGatewayPeers := func(ctx context.Context, getPeerCh <-chan test.ClusterGetPeerRequest, peer124, peer3 cluster.Peer) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()

		a := assertions.New(t)
		return test.AssertClusterGetPeerRequestSequence(ctx, getPeerCh,
			[]test.ClusterGetPeerResponse{
				{Error: errors.New("peer not found")},
				{Peer: peer124},
				{Peer: peer124},
				{Peer: peer3},
				{Peer: peer124},
			},
			func(reqCtx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.ClusterRole_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-0",
					})
			},
			func(reqCtx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.ClusterRole_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-1",
					})
			},
			func(reqCtx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.ClusterRole_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-2",
					})
			},
			func(reqCtx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.ClusterRole_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-3",
					})
			},
			func(reqCtx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
				return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
					a.So(role, should.Equal, ttnpb.ClusterRole_GATEWAY_SERVER) &&
					a.So(ids, should.Resemble, ttnpb.GatewayIdentifiers{
						GatewayID: "gateway-test-4",
					})
			},
		)
	}

	assertScheduleGateways := func(ctx context.Context, authCh <-chan test.ClusterAuthRequest, scheduleDownlink124Ch, scheduleDownlink3Ch <-chan NsGsScheduleDownlinkRequest, payload []byte, makeTxRequest func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest, fixedPaths bool, resps ...NsGsScheduleDownlinkResponse) (*ttnpb.DownlinkMessage, bool) {
		if len(resps) < 1 || len(resps) > 3 {
			panic("invalid response count specified")
		}

		t := test.MustTFromContext(ctx)
		t.Helper()

		a := assertions.New(t)

		makePath := func(i int) *ttnpb.DownlinkPath {
			if fixedPaths {
				return &ttnpb.DownlinkPath{
					Path: &ttnpb.DownlinkPath_Fixed{
						Fixed: &GatewayAntennaIdentifiers[i],
					},
				}
			}
			return &ttnpb.DownlinkPath{
				Path: &ttnpb.DownlinkPath_UplinkToken{
					UplinkToken: func() *ttnpb.RxMetadata {
						switch i {
						case 1:
							return RxMetadata[0]
						case 2:
							return RxMetadata[4]
						case 3:
							return RxMetadata[1]
						case 4:
							return RxMetadata[5]
						default:
							panic(fmt.Sprintf("Invalid index requested: %d", i))
						}
					}().UplinkToken,
				},
			}
		}

		var lastDown *ttnpb.DownlinkMessage
		var correlationIDs []string
		if !a.So(AssertAuthNsGsScheduleDownlinkRequest(ctx, authCh, scheduleDownlink124Ch,
			func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool {
				correlationIDs = msg.CorrelationIDs
				lastDown = &ttnpb.DownlinkMessage{
					CorrelationIDs: correlationIDs,
					RawPayload:     payload,
					Settings: &ttnpb.DownlinkMessage_Request{
						Request: makeTxRequest(
							makePath(1),
							makePath(2),
						),
					},
				}
				return a.So(msg, should.Resemble, lastDown)
			},
			grpc.EmptyCallOption{},
			resps[0],
		), should.BeTrue) {
			t.Error("Downlink assertion failed for gateways 1 and 2")
			return nil, false
		}
		t.Logf("Downlink correlation IDs: %v", correlationIDs)
		if len(resps) == 1 {
			return lastDown, true
		}

		lastDown = &ttnpb.DownlinkMessage{
			CorrelationIDs: correlationIDs,
			RawPayload:     payload,
			Settings: &ttnpb.DownlinkMessage_Request{
				Request: makeTxRequest(
					makePath(3),
				),
			},
		}
		if !a.So(AssertAuthNsGsScheduleDownlinkRequest(ctx, authCh, scheduleDownlink3Ch,
			func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool {
				return a.So(msg, should.Resemble, lastDown)
			},
			grpc.EmptyCallOption{},
			resps[1],
		), should.BeTrue) {
			t.Error("Downlink assertion failed for gateway 3")
			return nil, false
		}
		if len(resps) == 2 {
			return lastDown, true
		}

		lastDown = &ttnpb.DownlinkMessage{
			CorrelationIDs: correlationIDs,
			RawPayload:     payload,
			Settings: &ttnpb.DownlinkMessage_Request{
				Request: makeTxRequest(
					makePath(4),
				),
			},
		}
		if !a.So(AssertAuthNsGsScheduleDownlinkRequest(ctx, authCh, scheduleDownlink124Ch,
			func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool {
				return a.So(msg, should.Resemble, lastDown)
			},
			grpc.EmptyCallOption{},
			resps[2],
		), should.BeTrue) {
			t.Error("Downlink assertion failed for gateway 4")
			return nil, false
		}
		return lastDown, true
	}

	assertScheduleRxMetadataGateways := func(ctx context.Context, authCh <-chan test.ClusterAuthRequest, scheduleDownlink124Ch, scheduleDownlink3Ch <-chan NsGsScheduleDownlinkRequest, payload []byte, makeTxRequest func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest, resps ...NsGsScheduleDownlinkResponse) (*ttnpb.DownlinkMessage, bool) {
		return assertScheduleGateways(ctx, authCh, scheduleDownlink124Ch, scheduleDownlink3Ch, payload, makeTxRequest, false, resps...)
	}

	assertScheduleClassBCGateways := func(ctx context.Context, authCh <-chan test.ClusterAuthRequest, scheduleDownlink124Ch, scheduleDownlink3Ch <-chan NsGsScheduleDownlinkRequest, payload []byte, makeTxRequest func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest, resps ...NsGsScheduleDownlinkResponse) (*ttnpb.DownlinkMessage, bool) {
		return assertScheduleGateways(ctx, authCh, scheduleDownlink124Ch, scheduleDownlink3Ch, payload, makeTxRequest, true, resps...)
	}

	for _, tc := range []struct {
		Name               string
		DownlinkPriorities DownlinkPriorities
		Handler            func(context.Context, TestEnvironment) bool
		ErrorAssertion     func(*testing.T, error) bool
	}{
		{
			Name: "no device",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now().UTC())
					}()
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, nil)
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeEmpty)
					a.So(resp.Device, should.ResembleFields, &ttnpb.EndDevice{}, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "no MAC state",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeEmpty)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows closed",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-time.Second),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeEmpty)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.1/no uplink",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, time.Now().UTC())
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters:  makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters:  makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:        ttnpb.CLASS_A,
						LoRaWANVersion:     ttnpb.MAC_V1_1,
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeNil)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.1/no session",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-time.Second),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeEmpty)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.1/RX1,RX2 expired",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters:  makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters:  makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:        ttnpb.CLASS_A,
						LoRaWANVersion:     ttnpb.MAC_V1_1,
						RxWindowsAvailable: true,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-2*time.Second - time.Nanosecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
					},
					MACSettings: &ttnpb.MACSettings{
						StatusTimePeriodicity:  DurationPtr(0),
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeNil)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.1/RX1,RX2 available/no MAC/no application downlink",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-100 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					MACSettings: &ttnpb.MACSettings{
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
						StatusTimePeriodicity:  DurationPtr(0),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeEmpty)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.0.3/RX1,RX2 available/no MAC/generic application downlink/FCnt too low",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_3_REV_A,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_0_3_REV_A),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_0_3_REV_A),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_0_3,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-500 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					MACSettings: &ttnpb.MACSettings{
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
						StatusTimePeriodicity:  DurationPtr(0),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x22,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
							{
								CorrelationIDs: []string{"correlation-app-down-3", "correlation-app-down-4"},
								FCnt:           0x23,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"session.queued_application_downlinks",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				lastUp := lastUplink(getDevice.MACState.RecentUplinks...)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for ApplicationUplinks.Add to be called")

				case req := <-env.ApplicationUplinks.Add:
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.Uplinks, should.Resemble, []*ttnpb.ApplicationUp{
						{
							EndDeviceIdentifiers: getDevice.EndDeviceIdentifiers,
							CorrelationIDs:       lastUp.CorrelationIDs,
							Up: &ttnpb.ApplicationUp_DownlinkQueueInvalidated{
								DownlinkQueueInvalidated: &ttnpb.ApplicationInvalidatedDownlinks{
									Downlinks:    getDevice.Session.QueuedApplicationDownlinks,
									LastFCntDown: getDevice.Session.LastNFCntDown,
								},
							},
						},
					})
					close(req.Response)
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.0.3/RX1,RX2 available/no MAC/generic application downlink/application downlink exceeds length limit",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_3_REV_A,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_0_3_REV_A),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_0_3_REV_A),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_0_3,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-500 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					MACSettings: &ttnpb.MACSettings{
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
						StatusTimePeriodicity:  DurationPtr(0),
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     bytes.Repeat([]byte("x"), 250),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"session.queued_application_downlinks",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				lastUp := lastUplink(getDevice.MACState.RecentUplinks...)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for ApplicationUplinks.Add to be called")

				case req := <-env.ApplicationUplinks.Add:
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.Uplinks, should.Resemble, []*ttnpb.ApplicationUp{
						{
							EndDeviceIdentifiers: getDevice.EndDeviceIdentifiers,
							CorrelationIDs:       append(lastUp.CorrelationIDs, getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs...),
							Up: &ttnpb.ApplicationUp_DownlinkFailed{
								DownlinkFailed: &ttnpb.ApplicationDownlinkFailed{
									ApplicationDownlink: *getDevice.Session.QueuedApplicationDownlinks[0],
									Error:               *ttnpb.ErrorDetailsToProto(errApplicationDownlinkTooLong.WithAttributes("length", 250, "max", uint16(51))),
								},
							},
						},
					})
					close(req.Response)
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.1/RX1,RX2 available/MAC answers/MAC requests/generic application downlink/data+MAC/RX1,RX2/EU868/scheduling fail",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0110,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0b0000_0001,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Delay:         getDevice.MACState.CurrentParameters.Rx1Delay,
							Rx1Frequency:     431000000,
							Rx2DataRateIndex: getDevice.MACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.MACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				lastUp := lastUplink(getDevice.MACState.RecentUplinks...)
				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(lastUp.CorrelationIDs)+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range lastUp.CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.RxWindowsAvailable = false
				setDevice.MACState.QueuedResponses = nil

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.queued_responses",
						"mac_state.rx_windows_available",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.1/RX1,RX2 available/MAC answers/MAC requests/generic application downlink/data+MAC/RX1,RX2/EU868",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-500 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				downAt := clock.Add(time.Millisecond).Add(time.Second)

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0110,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0b0000_0001,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Delay:         getDevice.MACState.CurrentParameters.Rx1Delay,
							Rx1Frequency:     431000000,
							Rx2DataRateIndex: getDevice.MACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.MACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				lastUp := lastUplink(getDevice.MACState.RecentUplinks...)
				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(lastUp.CorrelationIDs)+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range lastUp.CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.LastConfirmedDownlinkAt = &downAt
				setDevice.MACState.LastDownlinkAt = &downAt
				setDevice.MACState.PendingRequests = []*ttnpb.MACCommand{
					{
						CID: ttnpb.CID_DEV_STATUS,
					},
				}
				setDevice.MACState.QueuedResponses = nil
				setDevice.MACState.RecentDownlinks = []*ttnpb.DownlinkMessage{
					lastDown,
				}
				setDevice.MACState.RxWindowsAvailable = false
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}
				setDevice.RecentDownlinks = []*ttnpb.DownlinkMessage{
					lastDown,
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.last_confirmed_downlink_at",
						"mac_state.last_downlink_at",
						"mac_state.pending_application_downlink",
						"mac_state.pending_requests",
						"mac_state.queued_responses",
						"mac_state.recent_downlinks",
						"mac_state.rx_windows_available",
						"recent_downlinks",
						"session",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class A/windows open/1.1/RX1,RX2 available/MAC answers/MAC requests/generic application downlink/application downlink does not fit due to FOpts/MAC/RX1,RX2/EU868",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				// NOTE: Maximum MACPayload length in both Rx1(DR0) and RX2(DR1) is 59. There are 6 bytes of FOpts, hence maximum fitting application downlink length is 59-8-6 == 45.
				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-500 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x15,
								FRMPayload:     bytes.Repeat([]byte{0x42}, 46),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				downAt := clock.Add(time.Millisecond).Add(time.Second)

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_1_0110,
							/*** FCnt ***/
							0x25, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x25,
							[]byte{
								/* ResetConf */
								0x01, 0b0000_0001,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x25,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Delay:         getDevice.MACState.CurrentParameters.Rx1Delay,
							Rx1Frequency:     431000000,
							Rx2DataRateIndex: getDevice.MACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.MACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				lastUp := lastUplink(getDevice.MACState.RecentUplinks...)
				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(lastUp.CorrelationIDs)) {
					for _, cid := range lastUp.CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.LastDownlinkAt = &downAt
				setDevice.MACState.LastConfirmedDownlinkAt = &downAt
				setDevice.MACState.QueuedResponses = nil
				setDevice.MACState.PendingRequests = []*ttnpb.MACCommand{
					{
						CID: ttnpb.CID_DEV_STATUS,
					},
				}
				setDevice.MACState.RecentDownlinks = []*ttnpb.DownlinkMessage{
					lastDown,
				}
				setDevice.MACState.RxWindowsAvailable = false
				setDevice.RecentDownlinks = []*ttnpb.DownlinkMessage{
					lastDown,
				}
				setDevice.Session.LastNFCntDown = 0x25

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.last_confirmed_downlink_at",
						"mac_state.last_downlink_at",
						"mac_state.pending_application_downlink",
						"mac_state.pending_requests",
						"mac_state.queued_responses",
						"mac_state.recent_downlinks",
						"mac_state.rx_windows_available",
						"recent_downlinks",
						"session",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		// Adapted from https://github.com/TheThingsNetwork/lorawan-stack/issues/866#issue-461484955.
		{
			Name: "Class A/windows open/1.0.2/RX1,RX2 available/MAC answers/MAC requests/generic application downlink/data+MAC/RX2 does not fit/RX1/EU868",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-500 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_6,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x15,
								FRMPayload:     []byte("AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUU="),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				downAt := clock.Add(time.Millisecond).Add(time.Second)

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0110,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0b0000_0001,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x15)

						/** FRMPayload **/
						b = append(b, []byte("AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUU=")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1DataRateIndex: ttnpb.DATA_RATE_6,
							Rx1Delay:         getDevice.MACState.CurrentParameters.Rx1Delay,
							Rx1Frequency:     431000000,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				lastUp := lastUplink(getDevice.MACState.RecentUplinks...)
				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(lastUp.CorrelationIDs)+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range lastUp.CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.LastConfirmedDownlinkAt = &downAt
				setDevice.MACState.LastDownlinkAt = &downAt
				setDevice.MACState.PendingRequests = []*ttnpb.MACCommand{
					{
						CID: ttnpb.CID_DEV_STATUS,
					},
				}
				setDevice.MACState.QueuedResponses = nil
				setDevice.MACState.RecentDownlinks = append(setDevice.MACState.RecentDownlinks, lastDown)
				setDevice.MACState.RxWindowsAvailable = false
				setDevice.RecentDownlinks = append(setDevice.RecentDownlinks, lastDown)
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.last_confirmed_downlink_at",
						"mac_state.last_downlink_at",
						"mac_state.pending_application_downlink",
						"mac_state.pending_requests",
						"mac_state.queued_responses",
						"mac_state.recent_downlinks",
						"mac_state.rx_windows_available",
						"recent_downlinks",
						"session",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class B/windows closed/ping slot",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := gpstime.Parse(10000 * beaconPeriod).Add(time.Second + 200*time.Millisecond)
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				nextDown := &ttnpb.ApplicationDownlink{
					CorrelationIDs: []string{"correlation-app-down-3", "correlation-app-down-4"},
					FCnt:           0x43,
					FPort:          0x2,
					FRMPayload:     []byte("nextTestPayload"),
					Priority:       ttnpb.TxSchedulePriority_HIGH,
					SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
				}
				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassBTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_B,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{
							Value: ttnpb.PING_EVERY_8S,
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{
										FHDR: ttnpb.FHDR{
											FCtrl: ttnpb.FCtrl{
												ClassB: true,
											},
										},
									}},
								},
								ReceivedAt: start.Add(-500 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								Confirmed:      true,
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
							nextDown,
						},
					},
				}

				delay := getDevice.MACState.DesiredParameters.Rx1Delay.Duration() / 2

				pingAt, ok := nextPingSlotAt(ctx, getDevice, start.Add(delay))
				if !ok {
					t.Errorf("Failed to determine ping slot time")
					return false
				}
				clock.Set(pingAt.Add(-nsScheduleWindow() - delay))

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				clock.Add(time.Second)

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b101_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_1_0001,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_B,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx2DataRateIndex: ttnpb.DATA_RATE_3,
							Rx2Frequency:     869525000,
							AbsoluteTime:     TimePtr(pingAt.UTC()),
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				downAt := clock.Now().Add(time.Second)

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.LastConfirmedDownlinkAt = &downAt
				setDevice.MACState.LastDownlinkAt = &downAt
				setDevice.MACState.LastNetworkInitiatedDownlinkAt = &downAt
				setDevice.MACState.PendingApplicationDownlink = setDevice.Session.QueuedApplicationDownlinks[0]
				setDevice.MACState.PendingRequests = []*ttnpb.MACCommand{
					{
						CID: ttnpb.CID_DEV_STATUS,
					},
				}
				setDevice.MACState.RecentDownlinks = append(setDevice.MACState.RecentDownlinks, lastDown)
				setDevice.RecentDownlinks = append(setDevice.RecentDownlinks, lastDown)
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{
					nextDown,
				}
				setDevice.Session.LastConfFCntDown = 0x42

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.last_confirmed_downlink_at",
						"mac_state.last_downlink_at",
						"mac_state.last_network_initiated_downlink_at",
						"mac_state.pending_application_downlink",
						"mac_state.pending_requests",
						"mac_state.queued_responses",
						"mac_state.recent_downlinks",
						"mac_state.rx_windows_available",
						"recent_downlinks",
						"session",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				clock.Add(time.Second)

				if timeout, interval := *setDevice.MACSettings.ClassBTimeout, networkInitiatedDownlinkInterval; timeout < interval {
					panic(fmt.Sprintf("class B timeout less than networkInitiatedDownlinkInterval (%v < %v)", timeout, interval))
				}
				nextPingSlot, ok := nextPingSlotAt(ctx, setDevice, downAt.Add(*setDevice.MACSettings.ClassBTimeout))
				if !a.So(ok, should.BeTrue) {
					t.Error("Failed to compute next ping slot")
				}
				if !AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(reqCtx context.Context, ids ttnpb.EndDeviceIdentifiers, startAt time.Time, replace bool) bool {
					return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
						a.So(ids, should.Resemble, setDevice.EndDeviceIdentifiers) &&
						a.So(replace, should.BeTrue) &&
						a.So(startAt, should.Resemble, nextPingSlot.Add(-getDevice.MACState.CurrentParameters.Rx1Delay.Duration()-NSScheduleWindow()))
				},
					nil,
				) {
					t.Error("Downlink task add assertion failed")
					return false
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/1.1/RX1,RX2 available/MAC answers/MAC requests/generic application downlink/data+MAC/RX1,RX2/EU868",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-500 * time.Millisecond),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				downAt := clock.Add(time.Millisecond).Add(time.Second)

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0110,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* ResetConf */
								0x01, 0b0000_0001,
								/* LinkCheckAns */
								0x02, 0x02, 0x05,
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Delay:         getDevice.MACState.CurrentParameters.Rx1Delay,
							Rx1Frequency:     431000000,
							Rx2DataRateIndex: getDevice.MACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.MACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				lastUp := lastUplink(getDevice.MACState.RecentUplinks...)
				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(lastUp.CorrelationIDs)+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range lastUp.CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters:       makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters:       makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:             ttnpb.CLASS_C,
						LastConfirmedDownlinkAt: &downAt,
						LastDownlinkAt:          &downAt,
						LoRaWANVersion:          ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							{
								CID: ttnpb.CID_DEV_STATUS,
							},
						},
						RecentDownlinks: []*ttnpb.DownlinkMessage{
							lastDown,
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							CopyUplinkMessage(lastUp),
						},
					},
					RecentDownlinks: []*ttnpb.DownlinkMessage{
						lastDown,
					},
					Session: &ttnpb.Session{
						DevAddr:                    devAddr,
						LastNFCntDown:              0x24,
						SessionKeys:                *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					},
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.last_confirmed_downlink_at",
						"mac_state.last_downlink_at",
						"mac_state.pending_application_downlink",
						"mac_state.pending_requests",
						"mac_state.queued_responses",
						"mac_state.recent_downlinks",
						"mac_state.rx_windows_available",
						"recent_downlinks",
						"session",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/1.1/RX1,RX2 expired/MAC answers/MAC requests/generic application downlink/data+MAC/RXC/EU868",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				clock.Add(time.Millisecond)

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedResponses: []*ttnpb.MACCommand{
							(&ttnpb.MACCommand_ResetConf{
								MinorVersion: 1,
							}).MACCommand(),
							(&ttnpb.MACCommand_LinkCheckAns{
								Margin:       2,
								GatewayCount: 5,
							}).MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-rx1Delay.Duration() - time.Second),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				downAt := clock.Add(time.Millisecond).Add(time.Second)

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0001,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_C,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx2DataRateIndex: getDevice.MACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.MACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.LastConfirmedDownlinkAt = &downAt
				setDevice.MACState.LastDownlinkAt = &downAt
				setDevice.MACState.LastNetworkInitiatedDownlinkAt = &downAt
				setDevice.MACState.PendingRequests = []*ttnpb.MACCommand{
					{
						CID: ttnpb.CID_DEV_STATUS,
					},
				}
				setDevice.MACState.QueuedResponses = nil
				setDevice.MACState.RecentDownlinks = append(setDevice.MACState.RecentDownlinks, lastDown)
				setDevice.MACState.RxWindowsAvailable = false
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}
				setDevice.RecentDownlinks = append(setDevice.RecentDownlinks, lastDown)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.last_confirmed_downlink_at",
						"mac_state.last_downlink_at",
						"mac_state.last_network_initiated_downlink_at",
						"mac_state.pending_application_downlink",
						"mac_state.pending_requests",
						"mac_state.queued_responses",
						"mac_state.recent_downlinks",
						"mac_state.rx_windows_available",
						"recent_downlinks",
						"session",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/1.1/RX1,RX2 expired/no MAC answers/MAC requests/classBC application downlink/absolute time within window/no forced gateways/data+MAC/RXC/EU868",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				absTime := start.Add(infrastructureDelay)

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:         test.EUFrequencyPlanID,
					LastDevStatusReceivedAt: TimePtr(start),
					LoRaWANPHYVersion:       ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout:         DurationPtr(42 * time.Second),
						StatusTimePeriodicity: DurationPtr(time.Hour),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-2 * time.Second),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
								ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
									AbsoluteTime: &absTime,
								},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				downAt := clock.Add(time.Millisecond).Add(time.Second)

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0000,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_C,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_NORMAL,
							Rx2DataRateIndex: ttnpb.DATA_RATE_0,
							Rx2Frequency:     869525000,
							AbsoluteTime:     &absTime,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.LastDownlinkAt = &downAt
				setDevice.MACState.LastNetworkInitiatedDownlinkAt = &downAt
				setDevice.MACState.RxWindowsAvailable = false
				setDevice.MACState.RecentDownlinks = append(setDevice.MACState.RecentDownlinks, lastDown)
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}
				setDevice.RecentDownlinks = append(setDevice.RecentDownlinks, lastDown)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"mac_state.last_confirmed_downlink_at",
						"mac_state.last_downlink_at",
						"mac_state.last_network_initiated_downlink_at",
						"mac_state.pending_application_downlink",
						"mac_state.pending_requests",
						"mac_state.queued_responses",
						"mac_state.recent_downlinks",
						"mac_state.rx_windows_available",
						"recent_downlinks",
						"session",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows closed/1.1/no MAC answers/MAC requests/classBC application downlink with absolute time/no forced gateways/MAC/RXC/EU868/non-retryable errors",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				absTime := start.Add(rx1Delay.Duration()).UTC()

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-rx1Delay.Duration()),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
								ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
									AbsoluteTime: &absTime,
								},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0001,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_C,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							AbsoluteTime:     &absTime,
							Rx2DataRateIndex: getDevice.MACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.MACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test").WithDetails(&ttnpb.ScheduleDownlinkErrorDetails{
							PathErrors: []*ttnpb.ErrorDetails{
								ttnpb.ErrorDetailsToProto(errors.DefineAborted(ulid.MustNew(0, test.Randy).String(), "aborted")),
								ttnpb.ErrorDetailsToProto(errors.DefineResourceExhausted(ulid.MustNew(0, test.Randy).String(), "resource exhausted")),
							},
						}),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test").WithDetails(&ttnpb.ScheduleDownlinkErrorDetails{
							PathErrors: []*ttnpb.ErrorDetails{
								ttnpb.ErrorDetailsToProto(errors.DefineFailedPrecondition(ulid.MustNew(0, test.Randy).String(), "failed precondition")),
							},
						}),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test").WithDetails(&ttnpb.ScheduleDownlinkErrorDetails{
							PathErrors: []*ttnpb.ErrorDetails{
								ttnpb.ErrorDetailsToProto(errors.DefineResourceExhausted(ulid.MustNew(0, test.Randy).String(), "resource exhausted")),
							},
						}),
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.MACState.RxWindowsAvailable = false
				setDevice.MACState.QueuedResponses = nil
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"session.queued_application_downlinks",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
					close(setFuncRespCh)
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for ApplicationUplinks.Add to be called")

				case req := <-env.ApplicationUplinks.Add:
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.Uplinks, should.Resemble, []*ttnpb.ApplicationUp{
						{
							EndDeviceIdentifiers: getDevice.EndDeviceIdentifiers,
							CorrelationIDs:       getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs,
							Up: &ttnpb.ApplicationUp_DownlinkFailed{
								DownlinkFailed: &ttnpb.ApplicationDownlinkFailed{
									ApplicationDownlink: *getDevice.Session.QueuedApplicationDownlinks[0],
									Error:               *ttnpb.ErrorDetailsToProto(errInvalidAbsoluteTime),
								},
							},
						},
					})
					close(req.Response)
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
					close(popFuncRespCh)
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows closed/1.1/no MAC answers/MAC requests/classBC application downlink/forced gateways/MAC/RXC/EU868/retryable error",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				now := clock.Add(time.Millisecond)

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout: DurationPtr(42 * time.Second),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-2 * rx1Delay.Duration()),
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
								ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
									Gateways: GatewayAntennaIdentifiers[:],
								},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				now = clock.Add(time.Millisecond)

				lastDown, ok := assertScheduleClassBCGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					func() []byte {
						b := []byte{
							/* MHDR */
							0b011_000_00,
							/* MACPayload */
							/** FHDR **/
							/*** DevAddr ***/
							devAddr[3], devAddr[2], devAddr[1], devAddr[0],
							/*** FCtrl ***/
							0b1_0_0_0_0001,
							/*** FCnt ***/
							0x42, 0x00,
						}

						/** FOpts **/
						b = append(b, test.Must(crypto.EncryptDownlink(
							nwkSEncKey,
							devAddr,
							0x24,
							[]byte{
								/* DevStatusReq */
								0x06,
							},
							true,
						)).([]byte)...)

						/** FPort **/
						b = append(b, 0x1)

						/** FRMPayload **/
						b = append(b, []byte("testPayload")...)

						/* MIC */
						mic := test.Must(crypto.ComputeDownlinkMIC(
							sNwkSIntKey,
							devAddr,
							0,
							0x42,
							b,
						)).([4]byte)
						return append(b, mic[:]...)
					}(),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_C,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGH,
							Rx2DataRateIndex: getDevice.MACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.MACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test").WithDetails(&ttnpb.ScheduleDownlinkErrorDetails{
							PathErrors: []*ttnpb.ErrorDetails{
								ttnpb.ErrorDetailsToProto(errors.DefineAborted(ulid.MustNew(0, test.Randy).String(), "aborted")),
								ttnpb.ErrorDetailsToProto(errors.DefineResourceExhausted(ulid.MustNew(0, test.Randy).String(), "resource exhausted")),
							},
						}),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test").WithDetails(&ttnpb.ScheduleDownlinkErrorDetails{
							PathErrors: []*ttnpb.ErrorDetails{
								ttnpb.ErrorDetailsToProto(errors.DefineCorruption(ulid.MustNew(0, test.Randy).String(), "corruption")), // retryable
							},
						}),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test").WithDetails(&ttnpb.ScheduleDownlinkErrorDetails{
							PathErrors: []*ttnpb.ErrorDetails{
								ttnpb.ErrorDetailsToProto(errors.DefineResourceExhausted(ulid.MustNew(0, test.Randy).String(), "resource exhausted")),
							},
						}),
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs)) {
					for _, cid := range getDevice.Session.QueuedApplicationDownlinks[0].CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeNil)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				if !AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(reqCtx context.Context, ids ttnpb.EndDeviceIdentifiers, startAt time.Time, replace bool) bool {
					return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
						a.So(ids, should.Resemble, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
							DevAddr:                &devAddr,
						}) &&
						a.So(replace, should.BeTrue) &&
						a.So(startAt, should.Resemble, now.Add(downlinkRetryInterval+NSScheduleWindow()))
				},
					nil,
				) {
					t.Error("Downlink task add assertion failed")
					return false
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/1.1/RX1,RX2 available/no MAC/classBC application downlink/absolute time outside window",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				absTime := clock.Now().Add(42 * time.Hour).UTC()

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout:          DurationPtr(42 * time.Second),
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
						StatusTimePeriodicity:  DurationPtr(0),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: clock.Now().Add(-time.Second),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
								ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
									AbsoluteTime: &absTime,
								},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeEmpty)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				if !AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(reqCtx context.Context, ids ttnpb.EndDeviceIdentifiers, startAt time.Time, replace bool) bool {
					return a.So(reqCtx, should.HaveParentContextOrEqual, ctx) &&
						a.So(ids, should.Resemble, setDevice.EndDeviceIdentifiers) &&
						a.So(replace, should.BeTrue) &&
						a.So(startAt, should.Resemble, absTime.Add(-rx1Delay.Duration()-NSScheduleWindow()))
				},
					nil,
				) {
					t.Error("Downlink task add assertion failed")
					return false
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "Class C/windows open/1.1/RX1,RX2 available/no MAC/expired application downlinks",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACSettings: &ttnpb.MACSettings{
						ClassCTimeout:          DurationPtr(42 * time.Second),
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
						StatusTimePeriodicity:  DurationPtr(0),
					},
					MACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_C,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{
							{
								CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
								DeviceChannelIndex: 3,
								Payload: &ttnpb.Message{
									MHDR: ttnpb.MHDR{
										MType: ttnpb.MType_UNCONFIRMED_UP,
									},
									Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
								},
								ReceivedAt: start.Add(-time.Second),
								RxMetadata: RxMetadata[:],
								Settings: ttnpb.TxSettings{
									DataRateIndex: ttnpb.DATA_RATE_0,
									Frequency:     430000000,
								},
							},
						},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
								ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
									AbsoluteTime: TimePtr(start.Add(-2).UTC()),
								},
							},
							{
								CorrelationIDs: []string{"correlation-app-down-3", "correlation-app-down-4"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
								ClassBC: &ttnpb.ApplicationDownlink_ClassBC{
									AbsoluteTime: TimePtr(start.Add(-1).UTC()),
								},
							},
						},
					},
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.Session.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.BeNil)
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},

		{
			Name: "join-accept/windows open/RX1,RX2 available/no active MAC state/EU868",
			DownlinkPriorities: DownlinkPriorities{
				JoinAccept:             ttnpb.TxSchedulePriority_HIGHEST,
				MACCommands:            ttnpb.TxSchedulePriority_HIGH,
				MaxApplicationDownlink: ttnpb.TxSchedulePriority_NORMAL,
			},
			Handler: func(ctx context.Context, env TestEnvironment) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)

				start := time.Now().UTC()
				clock := test.NewMockClock(start)
				defer SetMockClock(clock)()

				var popRespCh chan<- error
				popFuncRespCh := make(chan error)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop to be called")
					return false

				case req := <-env.DownlinkTasks.Pop:
					popRespCh = req.Response
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					go func() {
						popFuncRespCh <- req.Func(req.Context, ttnpb.EndDeviceIdentifiers{
							ApplicationIdentifiers: appID,
							DeviceID:               devID,
						}, start)
					}()
				}

				getDevice := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &types.DevAddr{0x42, 0xff, 0xff, 0xff},
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					PendingMACState: &ttnpb.MACState{
						CurrentParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DesiredParameters: makeEU868macParameters(ttnpb.PHY_V1_1_REV_B),
						DeviceClass:       ttnpb.CLASS_A,
						LoRaWANVersion:    ttnpb.MAC_V1_1,
						QueuedJoinAccept: &ttnpb.MACState_JoinAccept{
							Keys:    *sessionKeys,
							Payload: bytes.Repeat([]byte{0x42}, 33),
							Request: ttnpb.JoinRequest{
								DevAddr: devAddr,
							},
						},
						RxWindowsAvailable: true,
					},
					RecentUplinks: []*ttnpb.UplinkMessage{
						{
							CorrelationIDs:     []string{"correlation-up-1", "correlation-up-2"},
							DeviceChannelIndex: 3,
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_JOIN_REQUEST,
								},
								Payload: &ttnpb.Message_JoinRequestPayload{JoinRequestPayload: &ttnpb.JoinRequestPayload{
									JoinEUI:  types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
									DevEUI:   types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
									DevNonce: types.DevNonce{0x00, 0x42},
								}},
							},
							ReceivedAt: start.Add(-time.Second),
							RxMetadata: RxMetadata[:],
							Settings: ttnpb.TxSettings{
								DataRateIndex: ttnpb.DATA_RATE_0,
								Frequency:     430000000,
							},
						},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 0x24,
						SessionKeys:   *sessionKeys,
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								CorrelationIDs: []string{"correlation-app-down-1", "correlation-app-down-2"},
								FCnt:           0x42,
								FPort:          0x1,
								FRMPayload:     []byte("testPayload"),
								Priority:       ttnpb.TxSchedulePriority_HIGHEST,
								SessionKeyID:   []byte{0x11, 0x22, 0x33, 0x44},
							},
						},
					},
					SupportsJoin: true,
				}

				var setRespCh chan<- DeviceRegistrySetByIDResponse
				var setCtx context.Context
				setFuncRespCh := make(chan DeviceRegistrySetByIDRequestFuncResponse)
				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID to be called")
					return false

				case req := <-env.DeviceRegistry.SetByID:
					setRespCh = req.Response
					setCtx = req.Context
					a.So(req.Context, should.HaveParentContextOrEqual, ctx)
					a.So(req.ApplicationIdentifiers, should.Resemble, appID)
					a.So(req.DeviceID, should.Resemble, devID)
					a.So(req.Paths, should.HaveSameElementsDeep, getPaths)

					go func() {
						dev, sets, err := req.Func(req.Context, CopyEndDevice(getDevice))
						setFuncRespCh <- DeviceRegistrySetByIDRequestFuncResponse{
							Device: dev,
							Paths:  sets,
							Error:  err,
						}
					}()
				}

				scheduleDownlink124Ch := make(chan NsGsScheduleDownlinkRequest)
				peer124 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink124Ch),
				})

				scheduleDownlink3Ch := make(chan NsGsScheduleDownlinkRequest)
				peer3 := NewGSPeer(ctx, &MockNsGsServer{
					ScheduleDownlinkFunc: MakeNsGsScheduleDownlinkChFunc(scheduleDownlink3Ch),
				})

				if !a.So(assertGetGatewayPeers(ctx, env.Cluster.GetPeer, peer124, peer3), should.BeTrue) {
					return false
				}

				lastDown, ok := assertScheduleRxMetadataGateways(
					ctx,
					env.Cluster.Auth,
					scheduleDownlink124Ch,
					scheduleDownlink3Ch,
					bytes.Repeat([]byte{0x42}, 33),
					func(paths ...*ttnpb.DownlinkPath) *ttnpb.TxRequest {
						return &ttnpb.TxRequest{
							Class:            ttnpb.CLASS_A,
							DownlinkPaths:    paths,
							Priority:         ttnpb.TxSchedulePriority_HIGHEST,
							Rx1DataRateIndex: ttnpb.DATA_RATE_0,
							Rx1Delay:         ttnpb.RX_DELAY_5,
							Rx1Frequency:     431000000,
							Rx2DataRateIndex: getDevice.PendingMACState.CurrentParameters.Rx2DataRateIndex,
							Rx2Frequency:     getDevice.PendingMACState.CurrentParameters.Rx2Frequency,
							FrequencyPlanID:  test.EUFrequencyPlanID,
						}
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Error: errors.New("test"),
					},
					NsGsScheduleDownlinkResponse{
						Response: &ttnpb.ScheduleDownlinkResponse{
							Delay: time.Second,
						},
					},
				)
				if !a.So(ok, should.BeTrue) {
					t.Error("Scheduling assertion failed")
					return false
				}

				lastUp := lastUplink(getDevice.RecentUplinks...)
				if a.So(lastDown.CorrelationIDs, should.HaveLength, 1+len(lastUp.CorrelationIDs)) {
					for _, cid := range lastUp.CorrelationIDs {
						a.So(lastDown.CorrelationIDs, should.Contain, cid)
					}
				}

				setDevice := CopyEndDevice(getDevice)
				setDevice.PendingMACState.QueuedJoinAccept = nil
				setDevice.PendingMACState.PendingJoinRequest = &ttnpb.JoinRequest{
					DevAddr: devAddr,
				}
				setDevice.PendingMACState.RxWindowsAvailable = false
				setDevice.PendingSession = &ttnpb.Session{
					DevAddr:     devAddr,
					SessionKeys: *sessionKeys,
				}
				setDevice.RecentDownlinks = append(setDevice.RecentDownlinks, lastDown)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID callback to return")

				case resp := <-setFuncRespCh:
					a.So(resp.Error, should.BeNil)
					a.So(resp.Paths, should.HaveSameElementsDeep, []string{
						"pending_mac_state.pending_join_request",
						"pending_mac_state.queued_join_accept",
						"pending_mac_state.rx_windows_available",
						"pending_session.dev_addr",
						"pending_session.keys",
						"recent_downlinks",
					})
					a.So(resp.Device, should.ResembleFields, setDevice, resp.Paths)
				}
				close(setFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DeviceRegistry.SetByID response to be processed")

				case setRespCh <- DeviceRegistrySetByIDResponse{
					Device:  CopyEndDevice(setDevice),
					Context: setCtx,
				}:
				}

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop callback to return")

				case resp := <-popFuncRespCh:
					a.So(resp, should.BeNil)
				}
				close(popFuncRespCh)

				select {
				case <-ctx.Done():
					t.Error("Timed out while waiting for DownlinkTasks.Pop response to be processed")

				case popRespCh <- nil:
				}

				return true
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			ns, ctx, env, stop := StartTest(t, component.Config{}, Config{}, (1<<10)*test.Delay)
			defer stop()
			go LogEvents(t, env.Events)

			ns.downlinkPriorities = tc.DownlinkPriorities

			<-env.DownlinkTasks.Pop

			processDownlinkTaskErrCh := make(chan error)
			go func() {
				err := ns.processDownlinkTask(ctx)
				select {
				case <-ctx.Done():
					t.Log("NetworkServer.processDownlinkTask took too long to return")
					return

				default:
					processDownlinkTaskErrCh <- err
				}
			}()

			res := tc.Handler(ctx, env)
			if !a.So(res, should.BeTrue) {
				t.Error("Test handler failed")
				return
			}
			select {
			case <-ctx.Done():
				t.Error("Timed out while waiting for NetworkServer.processDownlinkTask to return")
				return

			case err := <-processDownlinkTaskErrCh:
				if tc.ErrorAssertion != nil {
					a.So(tc.ErrorAssertion(t, err), should.BeTrue)
				} else {
					a.So(err, should.BeNil)
				}
			}
			close(processDownlinkTaskErrCh)
		})
	}
}

func TestGenerateDownlink(t *testing.T) {
	const appIDString = "process-downlink-test-app-id"
	appID := ttnpb.ApplicationIdentifiers{ApplicationID: appIDString}
	const devID = "process-downlink-test-dev-id"

	devAddr := types.DevAddr{0x42, 0xff, 0xff, 0xff}

	fNwkSIntKey := types.AES128Key{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nwkSEncKey := types.AES128Key{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	sNwkSIntKey := types.AES128Key{0x42, 0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	encodeMessage := func(msg *ttnpb.Message, ver ttnpb.MACVersion, confFCnt uint32) []byte {
		msg = deepcopy.Copy(msg).(*ttnpb.Message)
		mac := msg.GetMACPayload()

		if len(mac.FRMPayload) > 0 && mac.FPort == 0 || len(mac.FOpts) > 0 {
			var key types.AES128Key
			switch ver {
			case ttnpb.MAC_V1_0, ttnpb.MAC_V1_0_1, ttnpb.MAC_V1_0_2:
				key = fNwkSIntKey
			case ttnpb.MAC_V1_1:
				key = nwkSEncKey
			default:
				panic(fmt.Errorf("unknown version %s", ver))
			}

			var err error
			if len(mac.FOpts) > 0 {
				mac.FOpts, err = crypto.EncryptDownlink(key, mac.DevAddr, mac.FCnt, mac.FOpts, true)
			} else {
				mac.FRMPayload, err = crypto.EncryptDownlink(key, mac.DevAddr, mac.FCnt, mac.FRMPayload, false)
			}
			if err != nil {
				t.Fatal("Failed to encrypt downlink MAC buffer")
			}
		}

		b, err := lorawan.MarshalMessage(*msg)
		if err != nil {
			t.Fatal("Failed to marshal downlink")
		}

		var key types.AES128Key
		switch ver {
		case ttnpb.MAC_V1_0, ttnpb.MAC_V1_0_1, ttnpb.MAC_V1_0_2:
			key = fNwkSIntKey
		case ttnpb.MAC_V1_1:
			key = sNwkSIntKey
		default:
			panic(fmt.Errorf("unknown version %s", ver))
		}

		mic, err := crypto.ComputeDownlinkMIC(key, mac.DevAddr, confFCnt, mac.FCnt, b)
		if err != nil {
			t.Fatal("Failed to compute MIC")
		}
		return append(b, mic[:]...)
	}

	encodeMAC := func(phy *band.Band, cmds ...*ttnpb.MACCommand) (b []byte) {
		for _, cmd := range cmds {
			b = test.Must(lorawan.DefaultMACCommands.AppendDownlink(*phy, b, *cmd)).([]byte)
		}
		return
	}

	for _, tc := range []struct {
		Name                         string
		Device                       *ttnpb.EndDevice
		Bytes                        []byte
		ApplicationDownlinkAssertion func(t *testing.T, down *ttnpb.ApplicationDownlink) bool
		DeviceAssertion              func(*testing.T, *ttnpb.EndDevice) bool
		Error                        error
	}{
		{
			Name: "1.1/no app downlink/no MAC/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				},
				Session:           ttnpb.NewPopulatedSession(test.Randy, false),
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Error: errNoDownlink,
		},
		{
			Name: "1.1/no app downlink/status after 1 downlink/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 3},
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:      ttnpb.MAC_V1_1,
					LastDevStatusFCntUp: 2,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				},
				Session: &ttnpb.Session{
					LastFCntUp: 4,
				},
				LoRaWANPHYVersion:       ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:         band.EU_863_870,
				LastDevStatusReceivedAt: TimePtr(time.Unix(42, 0)),
			},
			Error: errNoDownlink,
		},
		{
			Name: "1.1/no app downlink/status after an hour/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusTimePeriodicity: DurationPtr(24 * time.Hour),
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				},
				LoRaWANPHYVersion:       ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:         band.EU_863_870,
				LastDevStatusReceivedAt: TimePtr(time.Now()),
				Session:                 ttnpb.NewPopulatedSession(test.Randy, false),
			},
			Error: errNoDownlink,
		},
		{
			Name: "1.1/no app downlink/no MAC/ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_CONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{
								MACPayload: &ttnpb.MACPayload{
									FHDR: ttnpb.FHDR{
										FCnt: 24,
									},
								},
							},
						},
					}},
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr:       devAddr,
					LastNFCntDown: 41,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: true,
								ADR: true,
							},
							FCnt: 42,
						},
					},
				},
			}, ttnpb.MAC_V1_1, 24),
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				return assertions.New(t).So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{{
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_CONFIRMED_UP,
								},
								Payload: &ttnpb.Message_MACPayload{
									MACPayload: &ttnpb.MACPayload{
										FHDR: ttnpb.FHDR{
											FCnt: 24,
										},
									},
								},
							},
						}},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 41,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
				})
			},
		},
		{
			Name: "1.1/unconfirmed app downlink/no MAC/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							Confirmed:  false,
							FCnt:       42,
							FPort:      1,
							FRMPayload: []byte("test"),
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  false,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				return assertions.New(t).So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{{
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_UNCONFIRMED_UP,
								},
								Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
							},
						}},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr: devAddr,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
				})
			},
		},
		{
			Name: "1.1/unconfirmed app downlink/no MAC/ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_CONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{
								MACPayload: &ttnpb.MACPayload{
									FHDR: ttnpb.FHDR{
										FCnt: 24,
									},
								},
							},
						},
					}},
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							Confirmed:  false,
							FCnt:       42,
							FPort:      1,
							FRMPayload: []byte("test"),
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: true,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 24),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  false,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				return assertions.New(t).So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{{
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_CONFIRMED_UP,
								},
								Payload: &ttnpb.Message_MACPayload{
									MACPayload: &ttnpb.MACPayload{
										FHDR: ttnpb.FHDR{
											FCnt: 24,
										},
									},
								},
							},
						}},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr: devAddr,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
				})
			},
		},
		{
			Name: "1.1/confirmed app downlink/no MAC/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							Confirmed:  true,
							FCnt:       42,
							FPort:      1,
							FRMPayload: []byte("test"),
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_CONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  true,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{{
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_UNCONFIRMED_UP,
								},
								Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
							},
						}},
					},
					Session: &ttnpb.Session{
						DevAddr: devAddr,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
				})
			},
		},
		{
			Name: "1.1/confirmed app downlink/no MAC/ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_CONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{
								MACPayload: &ttnpb.MACPayload{
									FHDR: ttnpb.FHDR{
										FCnt: 24,
									},
								},
							},
						},
					}},
					RxWindowsAvailable: true,
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
					QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
						{
							Confirmed:  true,
							FCnt:       42,
							FPort:      1,
							FRMPayload: []byte("test"),
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_CONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: true,
								ADR: true,
							},
							FCnt: 42,
						},
						FPort:      1,
						FRMPayload: []byte("test"),
					},
				},
			}, ttnpb.MAC_V1_1, 24),
			ApplicationDownlinkAssertion: func(t *testing.T, down *ttnpb.ApplicationDownlink) bool {
				return assertions.New(t).So(down, should.Resemble, &ttnpb.ApplicationDownlink{
					Confirmed:  true,
					FCnt:       42,
					FPort:      1,
					FRMPayload: []byte("test"),
				})
			},
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						RecentUplinks: []*ttnpb.UplinkMessage{{
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_CONFIRMED_UP,
								},
								Payload: &ttnpb.Message_MACPayload{
									MACPayload: &ttnpb.MACPayload{
										FHDR: ttnpb.FHDR{
											FCnt: 24,
										},
									},
								},
							},
						}},
						RxWindowsAvailable: true,
					},
					Session: &ttnpb.Session{
						DevAddr: devAddr,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
				})
			},
		},
		{
			Name: "1.1/no app downlink/status(count)/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 3},
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion:      ttnpb.MAC_V1_1,
					LastDevStatusFCntUp: 4,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				},
				Session: &ttnpb.Session{
					DevAddr:       devAddr,
					LastFCntUp:    99,
					LastNFCntDown: 41,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
							FOpts: encodeMAC(
								LoRaWANBands[band.EU_863_870][ttnpb.PHY_V1_1_REV_B],
								ttnpb.CID_DEV_STATUS.MACCommand(),
							),
						},
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACSettings: &ttnpb.MACSettings{
						StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 3},
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion:      ttnpb.MAC_V1_1,
						LastDevStatusFCntUp: 4,
						PendingRequests: []*ttnpb.MACCommand{
							ttnpb.CID_DEV_STATUS.MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{{
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_UNCONFIRMED_UP,
								},
								Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
							},
						}},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastFCntUp:    99,
						LastNFCntDown: 41,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
				})
			},
		},
		{
			Name: "1.1/no app downlink/status(time/zero time)/no ack",
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: appID,
					DeviceID:               devID,
					DevAddr:                &devAddr,
				},
				MACSettings: &ttnpb.MACSettings{
					StatusTimePeriodicity: DurationPtr(time.Nanosecond),
				},
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					RecentUplinks: []*ttnpb.UplinkMessage{{
						Payload: &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_UNCONFIRMED_UP,
							},
							Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
						},
					}},
				},
				Session: &ttnpb.Session{
					DevAddr:       devAddr,
					LastNFCntDown: 41,
					SessionKeys: ttnpb.SessionKeys{
						NwkSEncKey: &ttnpb.KeyEnvelope{
							Key: &nwkSEncKey,
						},
						SNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &sNwkSIntKey,
						},
					},
				},
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				FrequencyPlanID:   band.EU_863_870,
			},
			Bytes: encodeMessage(&ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
					Major: ttnpb.Major_LORAWAN_R1,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							DevAddr: devAddr,
							FCtrl: ttnpb.FCtrl{
								Ack: false,
								ADR: true,
							},
							FCnt: 42,
							FOpts: encodeMAC(
								LoRaWANBands[band.EU_863_870][ttnpb.PHY_V1_1_REV_B],
								ttnpb.CID_DEV_STATUS.MACCommand(),
							),
						},
					},
				},
			}, ttnpb.MAC_V1_1, 0),
			DeviceAssertion: func(t *testing.T, dev *ttnpb.EndDevice) bool {
				a := assertions.New(t)
				if !a.So(dev.MACState, should.NotBeNil) {
					t.FailNow()
				}
				return a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						ApplicationIdentifiers: appID,
						DeviceID:               devID,
						DevAddr:                &devAddr,
					},
					MACSettings: &ttnpb.MACSettings{
						StatusTimePeriodicity: DurationPtr(time.Nanosecond),
					},
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						PendingRequests: []*ttnpb.MACCommand{
							ttnpb.CID_DEV_STATUS.MACCommand(),
						},
						RecentUplinks: []*ttnpb.UplinkMessage{{
							Payload: &ttnpb.Message{
								MHDR: ttnpb.MHDR{
									MType: ttnpb.MType_UNCONFIRMED_UP,
								},
								Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{}},
							},
						}},
					},
					Session: &ttnpb.Session{
						DevAddr:       devAddr,
						LastNFCntDown: 41,
						SessionKeys: ttnpb.SessionKeys{
							NwkSEncKey: &ttnpb.KeyEnvelope{
								Key: &nwkSEncKey,
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &sNwkSIntKey,
							},
						},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					FrequencyPlanID:   band.EU_863_870,
				})
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			logger := test.GetLogger(t)

			ctx := test.ContextWithTB(test.Context(), t)
			ctx = log.NewContext(ctx, logger)
			ctx, cancel := context.WithTimeout(ctx, (1<<7)*test.Delay)
			defer cancel()

			c := component.MustNew(
				log.Noop,
				&component.Config{},
				component.WithClusterNew(func(context.Context, *cluster.Config, ...cluster.Option) (cluster.Cluster, error) {
					return &test.MockCluster{
						JoinFunc: test.ClusterJoinNilFunc,
					}, nil
				}),
			)
			c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)

			componenttest.StartComponent(t, c)

			ns := &NetworkServer{
				Component: c,
				ctx:       ctx,
				defaultMACSettings: ttnpb.MACSettings{
					StatusTimePeriodicity:  DurationPtr(0),
					StatusCountPeriodicity: &pbtypes.UInt32Value{Value: 0},
				},
			}

			dev := CopyEndDevice(tc.Device)
			phy, err := deviceBand(dev, ns.FrequencyPlans)
			if !a.So(err, should.BeNil) {
				t.Fail()
				return
			}

			genDown, genState, err := ns.generateDataDownlink(ctx, dev, phy, dev.MACState.DeviceClass, time.Now(), math.MaxUint16, math.MaxUint16)
			if tc.Error != nil {
				a.So(err, should.EqualErrorOrDefinition, tc.Error)
				a.So(genDown, should.BeNil)
				return
			}
			// TODO: Assert AS uplinks generated(https://github.com/TheThingsNetwork/lorawan-stack/issues/631).

			if !a.So(err, should.BeNil) || !a.So(genDown, should.NotBeNil) {
				t.Fail()
				return
			}

			a.So(genDown.Payload, should.Resemble, tc.Bytes)
			if tc.ApplicationDownlinkAssertion != nil {
				a.So(tc.ApplicationDownlinkAssertion(t, genState.ApplicationDownlink), should.BeTrue)
			} else {
				a.So(genState.ApplicationDownlink, should.BeNil)
			}

			if tc.DeviceAssertion != nil {
				a.So(tc.DeviceAssertion(t, dev), should.BeTrue)
			} else {
				a.So(dev, should.Resemble, tc.Device)
			}
		})
	}
}
