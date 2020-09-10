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

package networkserver_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/mohae/deepcopy"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
)

func TestHandleUplink(t *testing.T) {
	type LegacyTestEnvironment struct {
		TestEnvironment
		ApplicationUplinks ApplicationUplinkQueueEnvironment
		DeviceRegistry     DeviceRegistryEnvironment
		DownlinkTasks      DownlinkTaskQueueEnvironment
		UplinkDeduplicator UplinkDeduplicatorEnvironment
	}

	joinGetByEUIPaths := [...]string{
		"frequency_plan_id",
		"lorawan_phy_version",
		"lorawan_version",
		"mac_settings",
		"session.dev_addr",
		"supports_class_b",
		"supports_class_c",
		"supports_join",
	}
	joinSetByIDGetPaths := [...]string{
		"frequency_plan_id",
		"lorawan_phy_version",
		"pending_session.queued_application_downlinks",
		"queued_application_downlinks",
		"recent_uplinks",
		"session.queued_application_downlinks",
	}
	joinSetByIDSetPaths := [...]string{
		"pending_mac_state",
		"recent_uplinks",
	}

	dataSetByIDGetPaths := [...]string{
		"frequency_plan_id",
		"last_dev_status_received_at",
		"lorawan_phy_version",
		"lorawan_version",
		"mac_settings",
		"mac_state",
		"multicast",
		"pending_mac_state",
		"pending_session",
		"queued_application_downlinks", // deprecated
		"recent_adr_uplinks",
		"recent_uplinks",
		"session",
		"supports_class_b",
		"supports_class_c",
		"supports_join",
	}
	dataRangeByAddrPaths := dataSetByIDGetPaths

	const (
		DeduplicationWindow = 24 * time.Millisecond
		CooldownWindow      = 42 * time.Millisecond
		CollectionWindow    = DeduplicationWindow + CooldownWindow

		Rx1Delay = ttnpb.RX_DELAY_5

		FPort = 0x42
		FCnt  = 42
	)

	makeApplicationDownlink := func(s string) *ttnpb.ApplicationDownlink {
		return &ttnpb.ApplicationDownlink{
			SessionKeyID: []byte("app-down-1-session-key-id"),
			FPort:        FPort,
			FCnt:         0x32,
			FRMPayload:   []byte(s),
			Confirmed:    true,
			Priority:     ttnpb.TxSchedulePriority_HIGH,
			CorrelationIDs: []string{
				"app-down-1-correlation-id-1",
			},
		}
	}

	makeJoinSetDevice := func(getDevice *ttnpb.EndDevice, decodedMsg *ttnpb.UplinkMessage, joinReq *ttnpb.JoinRequest, joinResp *ttnpb.JoinResponse) *ttnpb.EndDevice {
		keys := CopySessionKeys(&joinResp.SessionKeys)
		keys.AppSKey = nil

		macState := test.Must(NewMACState(getDevice, frequencyplans.NewStore(test.FrequencyPlansFetcher), ttnpb.MACSettings{})).(*ttnpb.MACState)
		macState.RxWindowsAvailable = true
		macState.QueuedJoinAccept = &ttnpb.MACState_JoinAccept{
			Keys:    *keys,
			Payload: joinResp.RawPayload,
			Request: *joinReq,
		}
		setDevice := CopyEndDevice(getDevice)
		setDevice.RecentUplinks = AppendRecentUplink(setDevice.RecentUplinks, decodedMsg, RecentUplinkCount)
		setDevice.PendingMACState = macState
		return setDevice
	}

	filterEndDevice := func(dev *ttnpb.EndDevice, paths ...string) *ttnpb.EndDevice {
		dev, err := ttnpb.FilterGetEndDevice(dev, paths...)
		if err != nil {
			panic(fmt.Errorf("failed to filter device: %w", err))
		}
		return dev
	}

	eventRPCErrorEqual := func(a, b events.Event) bool {
		if test.EventEqual(a, b) {
			return true
		}
		aErr, aOk := a.Data().(errors.Interface)
		bErr, bOk := b.Data().(errors.Interface)
		if !aOk || !bOk || !errors.Resemble(aErr, bErr) {
			return false
		}
		ap, err := events.Proto(a)
		if err != nil {
			return false
		}
		bp, err := events.Proto(b)
		if err != nil {
			return false
		}
		ap.UniqueID = ""
		bp.UniqueID = ""
		ap.Time = time.Time{}
		bp.Time = time.Time{}
		ap.Data = nil
		bp.Data = nil
		return reflect.DeepEqual(ap, bp)
	}

	assertHandleUplinkResponse := func(ctx context.Context, handleUplinkErrCh <-chan error, expectedErr error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for NetworkServer.HandleUplink to return")
			return false

		case resErr := <-handleUplinkErrCh:
			return a.So(resErr, should.EqualErrorOrDefinition, expectedErr)
		}
	}
	assertHandleUplink := func(ctx context.Context, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error, up *ttnpb.UplinkMessage, f func() bool, expectedErr error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		errCh := handle(ctx, up)
		return assertions.New(t).So(test.AllTrue(
			f(),
			assertHandleUplinkResponse(ctx, errCh, expectedErr),
		), should.BeTrue)
	}
	assertDeduplicateUplink := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, expectedUp *ttnpb.UplinkMessage, ok bool, err error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		return AssertDeduplicateUplink(ctx, env.UplinkDeduplicator.DeduplicateUplink, func(ctx context.Context, up *ttnpb.UplinkMessage, window time.Duration) bool {
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
				a.So(up, should.Resemble, expectedUp),
				a.So(window, should.Resemble, CollectionWindow),
			)
		},
			UplinkDeduplicatorDeduplicateUplinkResponse{
				Ok:    ok,
				Error: err,
			},
		)
	}
	assertAccumulatedMetadata := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, expectedUp *ttnpb.UplinkMessage, mds []*ttnpb.RxMetadata, err error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		return AssertAccumulatedMetadata(ctx, env.UplinkDeduplicator.AccumulatedMetadata, func(ctx context.Context, up *ttnpb.UplinkMessage) bool {
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
				a.So(up, should.Resemble, expectedUp),
			)
		},
			UplinkDeduplicatorAccumulatedMetadataResponse{
				Metadata: mds,
				Error:    err,
			},
		)
	}
	assertDownlinkTaskAdd := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, expectedIDs ttnpb.EndDeviceIdentifiers, expectedStartAt time.Time, replace bool, err error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		return AssertDownlinkTaskAddRequest(ctx, env.DownlinkTasks.Add, func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, startAt time.Time, replace bool) bool {
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
				a.So(ids, should.Resemble, expectedIDs),
				a.So(startAt, should.Resemble, expectedStartAt),
				a.So(replace, should.Equal, replace),
			)
		},
			err,
		)
	}
	assertInteropJoin := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, joinReq *ttnpb.JoinRequest, joinResp *ttnpb.JoinResponse, err error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		return AssertInteropClientHandleJoinRequestRequest(ctx, env.InteropClient.HandleJoinRequest,
			func(ctx context.Context, id types.NetID, req *ttnpb.JoinRequest) bool {
				joinReq.DevAddr = req.DevAddr
				return test.AllTrue(
					a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
					a.So(id, should.Equal, NetID),
					a.So(req, should.NotBeNil),
					a.So(req.DevAddr, should.NotBeEmpty),
					a.So(req.DevAddr.NwkID(), should.Resemble, NetID.ID()),
					a.So(req.DevAddr.NetIDType(), should.Equal, NetID.Type()),
					a.So(req, should.Resemble, joinReq),
				)
			},
			InteropClientHandleJoinRequestResponse{
				Response: joinResp,
				Error:    err,
			},
		)
	}
	assertClusterLocalJoin := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, expectedIDs ttnpb.EndDeviceIdentifiers, joinReq *ttnpb.JoinRequest, joinResp *ttnpb.JoinResponse, err error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		return AssertNsJsPeerHandleAuthJoinRequest(ctx, env.Cluster.GetPeer, env.Cluster.Auth,
			func(ctx context.Context, ids ttnpb.Identifiers) bool {
				return test.AllTrue(
					a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
					a.So(ids, should.Resemble, expectedIDs),
				)
			},
			func(ctx context.Context, req *ttnpb.JoinRequest) bool {
				joinReq.DevAddr = req.DevAddr
				return test.AllTrue(
					a.So(req, should.NotBeNil),
					a.So(req.DevAddr, should.NotBeEmpty),
					a.So(req.DevAddr.NwkID(), should.Resemble, NetID.ID()),
					a.So(req.DevAddr.NetIDType(), should.Equal, NetID.Type()),
					a.So(req, should.Resemble, joinReq),
				)
			},
			&grpc.EmptyCallOption{},
			NsJsHandleJoinResponse{
				Response: joinResp,
				Error:    err,
			},
		)
	}
	assertJoinGetByEUI := func(ctx context.Context, env LegacyTestEnvironment, upCIDs []string, getDevice *ttnpb.EndDevice, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		var getCtx context.Context
		return getCtx, AssertDeviceRegistryGetByEUI(ctx, env.DeviceRegistry.GetByEUI, func(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) bool {
			getCtx = ctx
			ctxCIDs := events.CorrelationIDsFromContext(ctx)
			for _, id := range upCIDs {
				a.So(ctxCIDs, should.Contain, id)
			}
			return test.AllTrue(
				a.So(ctxCIDs, should.HaveLength, 2+len(upCIDs)),
				a.So(joinEUI, should.Resemble, JoinEUI),
				a.So(devEUI, should.Resemble, DevEUI),
				a.So(paths, should.HaveSameElementsDeep, joinGetByEUIPaths[:]),
			)
		},
			func(ctx context.Context) DeviceRegistryGetByEUIResponse {
				if getDevice != nil {
					getDevice = filterEndDevice(CopyEndDevice(getDevice), joinGetByEUIPaths[:]...)
				}
				getCtx = context.WithValue(getCtx, &struct{}{}, "get")
				return DeviceRegistryGetByEUIResponse{
					Device:  getDevice,
					Context: getCtx,
					Error:   err,
				}
			})
	}
	assertJoinSetByID := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, getDevice, setDevice *ttnpb.EndDevice, expectedErr error, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		var setCtx context.Context
		return setCtx, AssertDeviceRegistrySetByID(ctx, env.DeviceRegistry.SetByID, func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) bool {
			setCtx = ctx
			dev, sets, err := f(ctx, CopyEndDevice(getDevice))
			var expectedSetPaths []string
			if setDevice != nil {
				expectedSetPaths = joinSetByIDSetPaths[:]
			}
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
				a.So(appID, should.Resemble, AppID),
				a.So(devID, should.Resemble, DevID),
				a.So(gets, should.HaveSameElementsDeep, joinSetByIDGetPaths),
				a.So(sets, should.HaveSameElementsDeep, expectedSetPaths),
				a.So(dev, should.ResembleFields, setDevice, sets),
				a.So(err, should.EqualErrorOrDefinition, expectedErr),
			)
		},
			func(ctx context.Context) DeviceRegistrySetByIDResponse {
				if setDevice != nil {
					setDevice = filterEndDevice(CopyEndDevice(setDevice), joinSetByIDGetPaths[:]...)
				}
				setCtx = context.WithValue(setCtx, &struct{}{}, "set")
				return DeviceRegistrySetByIDResponse{
					Device:  setDevice,
					Context: setCtx,
					Error:   err,
				}
			})
	}
	assertJoinApplicationUp := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, setDevice *ttnpb.EndDevice, joinResp *ttnpb.JoinResponse, recvAt time.Time, err error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		return AssertApplicationUplinkQueueAddRequest(ctx, env.ApplicationUplinks.Add, func(ctx context.Context, ups ...*ttnpb.ApplicationUp) bool {
			ids := *deepcopy.Copy(&setDevice.EndDeviceIdentifiers).(*ttnpb.EndDeviceIdentifiers)
			ids.DevAddr = &setDevice.PendingMACState.QueuedJoinAccept.Request.DevAddr

			queue := setDevice.GetPendingSession().GetQueuedApplicationDownlinks()
			if setDevice.Session != nil {
				queue = setDevice.Session.QueuedApplicationDownlinks
			}
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
				a.So(ups, should.Resemble, []*ttnpb.ApplicationUp{
					{
						CorrelationIDs:       events.CorrelationIDsFromContext(expectedCtx),
						EndDeviceIdentifiers: ids,
						Up: &ttnpb.ApplicationUp_JoinAccept{
							JoinAccept: &ttnpb.ApplicationJoinAccept{
								AppSKey:              joinResp.AppSKey,
								InvalidatedDownlinks: queue,
								ReceivedAt:           recvAt,
								SessionKeyID:         joinResp.SessionKeyID,
							},
						},
					},
				}),
			)
		}, err)
	}
	assertJoinDeduplicateSequence := func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, msg *ttnpb.UplinkMessage, chIdx uint8, drIdx ttnpb.DataRateIndex, dev *ttnpb.EndDevice, ok bool, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		getCtx, getOk := assertJoinGetByEUI(ctx, env, msg.CorrelationIDs, dev, nil)
		if !a.So(getOk, should.BeTrue) {
			return nil, false
		}
		*msg = *WithMatchedUplinkSettings(msg, chIdx, drIdx)
		msg.ReceivedAt = clock.Now()
		msg.CorrelationIDs = events.CorrelationIDsFromContext(getCtx)
		clock.Set(msg.ReceivedAt.Add(DeduplicationWindow))
		return getCtx, assertDeduplicateUplink(ctx, env, getCtx, msg, ok, err)
	}
	assertJoinGetPeerSequence := func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, msg *ttnpb.UplinkMessage, chIdx uint8, drIdx ttnpb.DataRateIndex, dev *ttnpb.EndDevice, peer cluster.Peer, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		getCtx, ok := assertJoinDeduplicateSequence(ctx, env, clock, msg, chIdx, drIdx, dev, true, nil)
		if !a.So(ok, should.BeTrue) {
			return nil, false
		}
		return getCtx, test.AssertClusterGetPeerRequest(ctx, env.Cluster.GetPeer, func(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, getCtx),
				a.So(role, should.Equal, ttnpb.ClusterRole_JOIN_SERVER),
				a.So(ids, should.Resemble, dev.EndDeviceIdentifiers),
			)
		},
			test.ClusterGetPeerResponse{
				Peer:  peer,
				Error: err,
			},
		)
	}
	assertJoinClusterLocalSequence := func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, msg *ttnpb.UplinkMessage, chIdx uint8, drIdx ttnpb.DataRateIndex, dev *ttnpb.EndDevice, joinReq *ttnpb.JoinRequest, joinResp *ttnpb.JoinResponse, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		getCtx, ok := assertJoinDeduplicateSequence(ctx, env, clock, msg, chIdx, drIdx, dev, true, nil)
		if !a.So(ok, should.BeTrue) {
			return nil, false
		}
		joinReq.CorrelationIDs = msg.CorrelationIDs
		return getCtx, assertClusterLocalJoin(ctx, env, getCtx, dev.EndDeviceIdentifiers, joinReq, joinResp, err)
	}
	assertJoinInteropSequence := func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, peerNotFound bool, msg *ttnpb.UplinkMessage, chIdx uint8, drIdx ttnpb.DataRateIndex, dev *ttnpb.EndDevice, joinReq *ttnpb.JoinRequest, joinResp *ttnpb.JoinResponse, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		var getCtx context.Context
		var ok bool
		if peerNotFound {
			getCtx, ok = assertJoinGetPeerSequence(ctx, env, clock, msg, chIdx, drIdx, dev, nil, ErrTestNotFound)
			joinReq.CorrelationIDs = msg.CorrelationIDs
		} else {
			getCtx, ok = assertJoinClusterLocalSequence(ctx, env, clock, msg, chIdx, drIdx, dev, joinReq, nil, ErrTestNotFound)
		}
		if !a.So(ok, should.BeTrue) {
			return nil, false
		}
		return getCtx, assertInteropJoin(ctx, env, getCtx, joinReq, joinResp, err)
	}

	assertDataRangeByAddr := func(ctx context.Context, env LegacyTestEnvironment, upCIDs []string, err error, getDevices ...*ttnpb.EndDevice) ([]context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		var rangeCtx context.Context
		var fCtxs []context.Context
		return fCtxs, AssertDeviceRegistryRangeByAddr(ctx, env.DeviceRegistry.RangeByAddr, func(ctx context.Context, devAddr types.DevAddr, paths []string, f func(context.Context, *ttnpb.EndDevice) bool) bool {
			rangeCtx = ctx
			ctxCIDs := events.CorrelationIDsFromContext(ctx)
			for _, id := range upCIDs {
				a.So(ctxCIDs, should.Contain, id)
			}
			if !a.So(test.AllTrue(
				a.So(ctxCIDs, should.HaveLength, 2+len(upCIDs)),
				a.So(devAddr, should.Resemble, DevAddr),
				a.So(paths, should.HaveSameElementsDeep, dataRangeByAddrPaths[:]),
			), should.BeTrue) {
				return false
			}
			for i, getDevice := range getDevices {
				fCtx := context.WithValue(rangeCtx, &struct{}{}, fmt.Sprintf("range:%d", i))
				fCtxs = append(fCtxs, fCtx)
				a.So(f(
					fCtx,
					filterEndDevice(CopyEndDevice(getDevice), dataRangeByAddrPaths[:]...),
				), should.BeTrue)
			}
			return true
		},
			err,
		)
	}
	assertDataDeduplicateSequence := func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, msg *ttnpb.UplinkMessage, chIdx uint8, drIdx ttnpb.DataRateIndex, devs []*ttnpb.EndDevice, idx int, ok bool, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		rangeCtxs, rangeOk := assertDataRangeByAddr(ctx, env, msg.CorrelationIDs, nil, devs...)
		if !a.So(rangeOk, should.BeTrue) || !a.So(len(rangeCtxs), should.BeGreaterThan, idx) {
			return nil, false
		}
		rangeCtx := rangeCtxs[idx]
		*msg = *WithMatchedUplinkSettings(msg, chIdx, drIdx)
		msg.ReceivedAt = clock.Now()
		msg.CorrelationIDs = events.CorrelationIDsFromContext(rangeCtx)
		clock.Set(msg.ReceivedAt.Add(DeduplicationWindow))
		return rangeCtx, assertDeduplicateUplink(ctx, env, rangeCtx, msg, ok, err)
	}
	assertDataSetByID := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, getDevice, setDevice *ttnpb.EndDevice, expectedSets []string, expectedErr error, err error) (context.Context, bool) {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		var setCtx context.Context
		return setCtx, AssertDeviceRegistrySetByID(ctx, env.DeviceRegistry.SetByID, func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) bool {
			setCtx = ctx
			dev, sets, err := f(ctx, CopyEndDevice(getDevice))
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
				a.So(appID, should.Resemble, AppID),
				a.So(devID, should.Resemble, DevID),
				a.So(gets, should.HaveSameElementsDeep, dataSetByIDGetPaths),
				a.So(sets, should.HaveSameElementsDeep, expectedSets),
				a.So(dev, should.ResembleFields, setDevice, sets),
				a.So(err, should.EqualErrorOrDefinition, expectedErr),
			)
		},
			func(ctx context.Context) DeviceRegistrySetByIDResponse {
				if setDevice != nil {
					setDevice = filterEndDevice(CopyEndDevice(setDevice), dataSetByIDGetPaths[:]...)
				}
				setCtx = context.WithValue(setCtx, &struct{}{}, "set")
				return DeviceRegistrySetByIDResponse{
					Device:  setDevice,
					Context: setCtx,
					Error:   err,
				}
			})
	}
	assertDataApplicationUp := func(ctx context.Context, env LegacyTestEnvironment, expectedCtx context.Context, setDevice *ttnpb.EndDevice, msg *ttnpb.UplinkMessage, err error) bool {
		t := test.MustTFromContext(ctx)
		t.Helper()
		a := assertions.New(t)
		macPayload := msg.Payload.GetMACPayload()
		return AssertApplicationUplinkQueueAddRequest(ctx, env.ApplicationUplinks.Add, func(ctx context.Context, ups ...*ttnpb.ApplicationUp) bool {
			return test.AllTrue(
				a.So(ctx, should.HaveParentContextOrEqual, expectedCtx),
				a.So(ups, should.Resemble, []*ttnpb.ApplicationUp{
					{
						CorrelationIDs:       events.CorrelationIDsFromContext(expectedCtx),
						EndDeviceIdentifiers: setDevice.EndDeviceIdentifiers,
						Up: &ttnpb.ApplicationUp_UplinkMessage{
							UplinkMessage: &ttnpb.ApplicationUplink{
								Confirmed:    msg.Payload.MType == ttnpb.MType_CONFIRMED_UP,
								FCnt:         macPayload.FCnt,
								FPort:        macPayload.FPort,
								FRMPayload:   macPayload.FRMPayload,
								ReceivedAt:   msg.ReceivedAt,
								RxMetadata:   msg.RxMetadata,
								SessionKeyID: setDevice.Session.SessionKeyID,
								Settings:     msg.Settings,
							},
						},
					},
				}),
			)
		}, err)
	}

	uplinkValidationErr, ok := errors.From((&ttnpb.UplinkMessage{}).ValidateFields())
	if !ok {
		t.Fatal("Failed to construct uplink validation error")
	}
	invalidUplinkSettingsErr := uplinkValidationErr.WithAttributes("field", "settings")

	makeChDRName := func(chIdx uint8, drIdx ttnpb.DataRateIndex, parts ...string) string {
		return MakeTestCaseName(append(parts, fmt.Sprintf("Channel:%d", chIdx), fmt.Sprintf("DR:%d", drIdx))...)
	}

	type TestCase struct {
		Name    string
		Handler func(context.Context, LegacyTestEnvironment, *test.MockClock, func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool
	}
	tcs := []TestCase{
		{
			Name: "No settings",
			Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)
				return a.So(test.AllTrue(
					assertHandleUplinkResponse(ctx, handle(ctx, &ttnpb.UplinkMessage{}), invalidUplinkSettingsErr),
					assertHandleUplinkResponse(ctx, handle(ctx, &ttnpb.UplinkMessage{
						RawPayload: []byte("testpayload"),
						RxMetadata: RxMetadata[:1],
					}), invalidUplinkSettingsErr),
					assertHandleUplinkResponse(ctx, handle(ctx, &ttnpb.UplinkMessage{
						RawPayload: []byte("testpayload"),
					}), invalidUplinkSettingsErr),
				), should.BeTrue)
			},
		},
	}
	for _, uplinkMDs := range [][]*ttnpb.RxMetadata{
		nil,
		RxMetadata[0:1],
		RxMetadata[1:4],
	} {
		uplinkMDs := uplinkMDs
		makeMDName := func(parts ...string) string {
			return MakeTestCaseName(append(parts, fmt.Sprintf("Metadata length:%d", len(uplinkMDs)))...)
		}
		ForEachBand(t, func(makeLoopName func(...string) string, phy *band.Band, phyVersion ttnpb.PHYVersion) {
			chIdx := uint8(len(phy.UplinkChannels) - 1)
			ch := phy.UplinkChannels[chIdx]
			drIdx := ch.MaxDataRate
			dr := phy.DataRates[drIdx].Rate

			makeName := func(parts ...string) string {
				return makeMDName(makeChDRName(chIdx, drIdx, makeLoopName(parts...)))
			}
			tcs = append(tcs,
				TestCase{
					Name: makeName("No payload"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						return assertions.New(test.MustTFromContext(ctx)).So(assertHandleUplinkResponse(ctx, handle(ctx, &ttnpb.UplinkMessage{
							RxMetadata: uplinkMDs,
							Settings:   MakeUplinkSettings(dr, ch.Frequency),
						}), ErrDecodePayload.WithCause(lorawan.UnmarshalMessage(nil, nil))), should.BeTrue)
					},
				},
				TestCase{
					Name: makeName("Unknown Major"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						return assertions.New(test.MustTFromContext(ctx)).So(assertHandleUplinkResponse(ctx, handle(ctx, &ttnpb.UplinkMessage{
							RawPayload: []byte{
								/* MHDR */
								0b000_000_01,
								/* Join-request */
								/** JoinEUI **/
								JoinEUI[7], JoinEUI[6], JoinEUI[5], JoinEUI[4], JoinEUI[3], JoinEUI[2], JoinEUI[1], JoinEUI[0],
								/** DevEUI **/
								DevEUI[7], DevEUI[6], DevEUI[5], DevEUI[4], DevEUI[3], DevEUI[2], DevEUI[1], DevEUI[0],
								/** DevNonce **/
								0x01, 0x00,
								/* MIC */
								0x03, 0x02, 0x01, 0x00,
							},
							RxMetadata: uplinkMDs,
							Settings:   MakeUplinkSettings(dr, ch.Frequency),
						}), ErrUnsupportedLoRaWANVersion.WithAttributes("version", uint32(1))), should.BeTrue)
					},
				},
				TestCase{
					Name: makeName("Invalid MType"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						return assertions.New(test.MustTFromContext(ctx)).So(assertHandleUplinkResponse(ctx, handle(ctx, &ttnpb.UplinkMessage{
							RawPayload: bytes.Repeat([]byte{0x20}, 33),
							RxMetadata: uplinkMDs,
							Settings:   MakeUplinkSettings(dr, ch.Frequency),
						}), nil), should.BeTrue)
					},
				},
				TestCase{
					Name: makeName("Proprietary MType"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						phyPayload := []byte{
							/* MHDR */
							0b111_000_00,
						}
						return assertions.New(test.MustTFromContext(ctx)).So(assertHandleUplinkResponse(ctx, handle(ctx, &ttnpb.UplinkMessage{
							RawPayload: phyPayload,
							RxMetadata: uplinkMDs,
							Settings:   MakeUplinkSettings(dr, ch.Frequency),
						}), ErrDecodePayload.WithCause(lorawan.UnmarshalMessage(phyPayload, &ttnpb.Message{}))), should.BeTrue)
					},
				},
			)
		})

		ForEachFrequencyPlanLoRaWANVersionPair(t, func(makeLoopName func(...string) string, fpID string, fp *frequencyplans.FrequencyPlan, phy *band.Band, macVersion ttnpb.MACVersion, phyVersion ttnpb.PHYVersion) {
			makeJoinResponse := func(withAppSKey bool) *ttnpb.JoinResponse {
				return &ttnpb.JoinResponse{
					RawPayload:  bytes.Repeat([]byte{0x42}, 17),
					SessionKeys: *MakeSessionKeys(macVersion, withAppSKey),
				}
			}
			makeJoinDevice := func(clock *test.MockClock) *ttnpb.EndDevice {
				return &ttnpb.EndDevice{
					EndDeviceIdentifiers: *MakeOTAAIdentifiers(nil),
					FrequencyPlanID:      fpID,
					LoRaWANPHYVersion:    phyVersion,
					LoRaWANVersion:       macVersion,
					MACSettings: &ttnpb.MACSettings{
						Rx1Delay: &ttnpb.RxDelayValue{
							Value: ttnpb.RX_DELAY_3,
						},
						DesiredRx2DataRateIndex: &ttnpb.DataRateIndexValue{
							Value: ttnpb.DATA_RATE_2,
						},
					},
					SupportsJoin: true,
					CreatedAt:    clock.Now(),
					UpdatedAt:    clock.Now(),
				}
			}

			chIdx := uint8(len(phy.UplinkChannels) - 1)
			ch := phy.UplinkChannels[chIdx]
			drIdx := ch.MaxDataRate
			dr := phy.DataRates[drIdx].Rate

			makeNsJsJoinRequest := func(devAddr *types.DevAddr, correlationIDs ...string) *ttnpb.JoinRequest {
				return MakeNsJsJoinRequest(macVersion, phyVersion, fp, devAddr, ttnpb.RX_DELAY_3, 0, ttnpb.DATA_RATE_2, correlationIDs...)
			}
			makeJoinRequest := func(decodePayload bool) *ttnpb.UplinkMessage {
				return MakeJoinRequest(decodePayload, dr, ch.Frequency, uplinkMDs...)
			}

			makeJoinName := func(parts ...string) string {
				return makeMDName(makeChDRName(chIdx, drIdx, makeLoopName(append([]string{"Join-request"}, parts...)...)))
			}
			tcs = append(tcs,
				TestCase{
					Name: makeJoinName("Get fail"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							_, ok := assertJoinGetByEUI(ctx, env, msg.CorrelationIDs, nil, ErrTestInternal)
							if !a.So(ok, should.BeTrue) {
								return false
							}
							return ok
						}, ErrTestInternal), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get ABP device"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						getDevice.SupportsJoin = false
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinGetByEUI(ctx, env, msg.CorrelationIDs, getDevice, nil)
							decodedMsg.ReceivedAt = clock.Now()
							decodedMsg.CorrelationIDs = events.CorrelationIDsFromContext(getCtx)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrABPJoinRequest),
								),
							), should.BeTrue)
						}, nil), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "Deduplication fail"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, getDevice, false, ErrTestInternal)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
								),
							), should.BeTrue)
						}, ErrTestInternal), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "Duplicate uplink"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, getDevice, false, nil)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrDuplicate),
								),
							), should.BeTrue)
						}, nil), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS fail"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinClusterLocalSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, getDevice, joinReq, nil, ErrTestInternal)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
								),
								a.So(env.Events, should.ReceiveEventsFunc, eventRPCErrorEqual,
									EvtClusterJoinFail.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
								),
							), should.BeTrue)
						}, ErrTestInternal), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS not found", "Interop JS fail"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinInteropSequence(ctx, env, clock, true, decodedMsg, chIdx, drIdx, getDevice, joinReq, nil, ErrTestInternal)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtInteropJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
								),
								a.So(env.Events, should.ReceiveEventsFunc, eventRPCErrorEqual,
									EvtInteropJoinFail.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
								),
							), should.BeTrue)
						}, ErrTestInternal), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS does not contain device", "Interop JS accept", "Metadata merge fail", "Current downlink queue", "Set fail on read"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						getDevice.Session = &ttnpb.Session{
							DevAddr: DevAddr,
							QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
								makeApplicationDownlink("test"),
							},
						}
						getDevice.DevAddr = &getDevice.Session.DevAddr
						joinResp := makeJoinResponse(true)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinInteropSequence(ctx, env, clock, false, decodedMsg, chIdx, drIdx, getDevice, joinReq, joinResp, nil)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
								),
								a.So(env.Events, should.ReceiveEventFunc, eventRPCErrorEqual, EvtClusterJoinFail.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestNotFound)),
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtInteropJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
									EvtInteropJoinSuccess.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, JoinResponseWithoutKeys(joinResp)),
								),
								assertAccumulatedMetadata(ctx, env, getCtx, decodedMsg, nil, ErrTestInternal),
								AssertDeviceRegistrySetByID(ctx, env.DeviceRegistry.SetByID, func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) bool {
									return test.AllTrue(
										a.So(ctx, should.HaveParentContextOrEqual, getCtx),
										a.So(appID, should.Resemble, AppID),
										a.So(devID, should.Resemble, DevID),
										a.So(gets, should.HaveSameElementsDeep, joinSetByIDGetPaths),
									)
								}, func(context.Context) DeviceRegistrySetByIDResponse {
									return DeviceRegistrySetByIDResponse{
										Error: ErrTestInternal,
									}
								}),
								a.So(env.Events, should.ReceiveEventResembling,
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
								),
							), should.BeTrue)
						}, ErrTestInternal), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS does not contain device", "Interop JS accept", "Metadata merge fail", "Current downlink queue", "Device deleted during handling"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						getDevice.Session = &ttnpb.Session{
							DevAddr: DevAddr,
							QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
								makeApplicationDownlink("test"),
							},
						}
						getDevice.DevAddr = &getDevice.Session.DevAddr
						joinResp := makeJoinResponse(true)
						joinReq := makeNsJsJoinRequest(nil)
						innerErr := ErrOutdatedData
						registryErr := ErrTestInternal.WithCause(innerErr)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinInteropSequence(ctx, env, clock, false, decodedMsg, chIdx, drIdx, getDevice, joinReq, joinResp, nil)
							if !a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
								),
								a.So(env.Events, should.ReceiveEventFunc, eventRPCErrorEqual, EvtClusterJoinFail.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestNotFound)),
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtInteropJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
									EvtInteropJoinSuccess.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, JoinResponseWithoutKeys(joinResp)),
								),
								assertAccumulatedMetadata(ctx, env, getCtx, decodedMsg, nil, ErrTestInternal),
							), should.BeTrue) {
								return false
							}
							_, ok = assertJoinSetByID(ctx, env, getCtx, nil, nil, innerErr, registryErr)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventResembling,
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, registryErr),
								),
							), should.BeTrue)
						}, registryErr), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS does not contain device", "Interop JS accept", "Metadata merge fail", "Current downlink queue", "Set fail on write"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						getDevice.Session = &ttnpb.Session{
							DevAddr: DevAddr,
							QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
								makeApplicationDownlink("test"),
							},
						}
						getDevice.DevAddr = &getDevice.Session.DevAddr
						joinResp := makeJoinResponse(true)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinInteropSequence(ctx, env, clock, false, decodedMsg, chIdx, drIdx, getDevice, joinReq, joinResp, nil)
							if !a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
								),
								a.So(env.Events, should.ReceiveEventFunc, eventRPCErrorEqual, EvtClusterJoinFail.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestNotFound)),
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtInteropJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
									EvtInteropJoinSuccess.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, JoinResponseWithoutKeys(joinResp)),
								),
								assertAccumulatedMetadata(ctx, env, getCtx, decodedMsg, nil, ErrTestInternal),
							), should.BeTrue) {
								return false
							}
							_, ok = assertJoinSetByID(ctx, env, getCtx, getDevice, makeJoinSetDevice(getDevice, decodedMsg, joinReq, joinResp), nil, ErrTestInternal)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventResembling,
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
								),
							), should.BeTrue)
						}, ErrTestInternal), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS does not contain device", "Interop JS accept", "Metadata merge fail", "Current downlink queue", "Set success", "Downlink add success", "AppSKey present", "Application uplink add success"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						getDevice.Session = &ttnpb.Session{
							DevAddr: DevAddr,
							QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
								makeApplicationDownlink("test"),
							},
						}
						getDevice.DevAddr = &getDevice.Session.DevAddr
						joinResp := makeJoinResponse(true)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinInteropSequence(ctx, env, clock, false, decodedMsg, chIdx, drIdx, getDevice, joinReq, joinResp, nil)
							if !a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
								),
								a.So(env.Events, should.ReceiveEventFunc, eventRPCErrorEqual, EvtClusterJoinFail.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestNotFound)),
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtInteropJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
									EvtInteropJoinSuccess.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, JoinResponseWithoutKeys(joinResp)),
								),
								assertAccumulatedMetadata(ctx, env, getCtx, decodedMsg, nil, ErrTestInternal),
							), should.BeTrue) {
								return false
							}
							joinRespRecvAt := clock.Now()
							clock.Add(time.Nanosecond)
							setDevice := makeJoinSetDevice(getDevice, decodedMsg, joinReq, joinResp)
							setCtx, ok := assertJoinSetByID(ctx, env, getCtx, getDevice, setDevice, nil, nil)
							return test.AllTrue(
								ok,
								assertDownlinkTaskAdd(ctx, env, setCtx, setDevice.EndDeviceIdentifiers, decodedMsg.ReceivedAt.Add(-InfrastructureDelay/2+phy.JoinAcceptDelay1-joinReq.RxDelay.Duration()/2-NSScheduleWindow()), true, nil),
								assertJoinApplicationUp(ctx, env, setCtx, setDevice, joinResp, joinRespRecvAt, nil),
								a.So(env.Events, should.ReceiveEventResembling,
									EvtProcessJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
								),
							)
						}, nil), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS accept", "Metadata merge success", "Pending downlink queue", "Set success", "Downlink add fail", "No AppSKey", "Application uplink add fail"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						getDevice.PendingSession = &ttnpb.Session{
							DevAddr: DevAddr,
							QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
								makeApplicationDownlink("test"),
							},
						}
						joinResp := makeJoinResponse(false)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinClusterLocalSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, getDevice, joinReq, joinResp, nil)
							if !a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
									EvtClusterJoinSuccess.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, JoinResponseWithoutKeys(joinResp)),
								),
								assertAccumulatedMetadata(ctx, env, getCtx, decodedMsg, RxMetadata[:], nil),
							), should.BeTrue) {
								return false
							}
							joinRespRecvAt := clock.Now()
							clock.Add(time.Nanosecond)
							decodedMsg.RxMetadata = RxMetadata[:]
							setDevice := makeJoinSetDevice(getDevice, decodedMsg, joinReq, joinResp)
							setCtx, ok := assertJoinSetByID(ctx, env, getCtx, getDevice, setDevice, nil, nil)
							return a.So(test.AllTrue(
								ok,
								assertDownlinkTaskAdd(ctx, env, setCtx, setDevice.EndDeviceIdentifiers, decodedMsg.ReceivedAt.Add(-InfrastructureDelay/2+phy.JoinAcceptDelay1-joinReq.RxDelay.Duration()/2-NSScheduleWindow()), true, ErrTestInternal),
								assertJoinApplicationUp(ctx, env, setCtx, setDevice, joinResp, joinRespRecvAt, ErrTestInternal),
								a.So(env.Events, should.ReceiveEventResembling,
									EvtProcessJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
								),
							), should.BeTrue)
						}, nil), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS accept", "Metadata merge success", "Both downlink queues", "Set success", "Downlink add fail", "No AppSKey", "Application uplink add success"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						getDevice.Session = &ttnpb.Session{
							DevAddr: DevAddr,
							QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
								makeApplicationDownlink("test"),
							},
						}
						getDevice.DevAddr = &getDevice.Session.DevAddr
						getDevice.PendingSession = &ttnpb.Session{
							DevAddr: DevAddr,
							QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
								makeApplicationDownlink("pending"),
								makeApplicationDownlink("other"),
								makeApplicationDownlink("foo"),
							},
						}
						joinResp := makeJoinResponse(false)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinClusterLocalSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, getDevice, joinReq, joinResp, nil)
							if !a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
									EvtClusterJoinSuccess.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, JoinResponseWithoutKeys(joinResp)),
								),
								assertAccumulatedMetadata(ctx, env, getCtx, decodedMsg, RxMetadata[:], nil),
							), should.BeTrue) {
								return false
							}
							joinRespRecvAt := clock.Now()
							clock.Add(time.Nanosecond)
							decodedMsg.RxMetadata = RxMetadata[:]
							setDevice := makeJoinSetDevice(getDevice, decodedMsg, joinReq, joinResp)
							setCtx, ok := assertJoinSetByID(ctx, env, getCtx, getDevice, setDevice, nil, nil)
							return a.So(test.AllTrue(
								ok,
								assertDownlinkTaskAdd(ctx, env, setCtx, setDevice.EndDeviceIdentifiers, decodedMsg.ReceivedAt.Add(-InfrastructureDelay/2+phy.JoinAcceptDelay1-joinReq.RxDelay.Duration()/2-NSScheduleWindow()), true, nil),
								assertJoinApplicationUp(ctx, env, setCtx, setDevice, joinResp, joinRespRecvAt, ErrTestInternal),
								a.So(env.Events, should.ReceiveEventResembling,
									EvtProcessJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
								),
							), should.BeTrue)
						}, nil), should.BeTrue)
					},
				},
				TestCase{
					Name: makeJoinName("Get OTAA device", "First uplink", "Cluster-local JS accept", "Metadata merge success", "No downlink queue", "Set fail on write"),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						t := test.MustTFromContext(ctx)
						a := assertions.New(t)
						msg := makeJoinRequest(false)
						decodedMsg := makeJoinRequest(true)
						getDevice := makeJoinDevice(clock)
						joinResp := makeJoinResponse(false)
						joinReq := makeNsJsJoinRequest(nil)
						return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
							getCtx, ok := assertJoinClusterLocalSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, getDevice, joinReq, joinResp, nil)
							if !a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventsResembling,
									EvtReceiveJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, decodedMsg),
									EvtClusterJoinAttempt.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, joinReq),
									EvtClusterJoinSuccess.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, JoinResponseWithoutKeys(joinResp)),
								),
								assertAccumulatedMetadata(ctx, env, getCtx, decodedMsg, RxMetadata[:], nil),
							), should.BeTrue) {
								return false
							}
							decodedMsg.RxMetadata = RxMetadata[:]
							setDevice := makeJoinSetDevice(getDevice, decodedMsg, joinReq, joinResp)
							_, ok = assertJoinSetByID(ctx, env, getCtx, getDevice, setDevice, nil, ErrTestInternal)
							return a.So(test.AllTrue(
								ok,
								a.So(env.Events, should.ReceiveEventResembling,
									EvtDropJoinRequest.NewWithIdentifiersAndData(getCtx, getDevice.EndDeviceIdentifiers, ErrTestInternal),
								),
							), should.BeTrue)
						}, ErrTestInternal), should.BeTrue)
					},
				},
			)

			for _, typ := range []ttnpb.RejoinType{
				ttnpb.RejoinType_CONTEXT,
				ttnpb.RejoinType_KEYS,
				ttnpb.RejoinType_SESSION,
			} {
				typ := typ
				if macVersion.Compare(ttnpb.MAC_V1_1) < 0 {
					continue
				}
				makeRejoinRequest := func(decodePayload bool) *ttnpb.UplinkMessage {
					var phyPayload []byte
					switch typ {
					case ttnpb.RejoinType_CONTEXT, ttnpb.RejoinType_KEYS:
						phyPayload = []byte{
							/* MHDR */
							0b110_000_00,
							byte(typ),
							NetID[2], NetID[1], NetID[0],
							DevEUI[7], DevEUI[6], DevEUI[5], DevEUI[4], DevEUI[3], DevEUI[2], DevEUI[1], DevEUI[0],
							/* RejoinCnt0 */
							0x01, 0x00,
							/* MIC */
							0x03, 0x02, 0x01, 0x00,
						}
					case ttnpb.RejoinType_SESSION:
						phyPayload = []byte{
							/* MHDR */
							0b110_000_00,
							byte(typ),
							JoinEUI[7], JoinEUI[6], JoinEUI[5], JoinEUI[4], JoinEUI[3], JoinEUI[2], JoinEUI[1], JoinEUI[0],
							DevEUI[7], DevEUI[6], DevEUI[5], DevEUI[4], DevEUI[3], DevEUI[2], DevEUI[1], DevEUI[0],
							/* RejoinCnt1 */
							0x01, 0x00,
							/* MIC */
							0x03, 0x02, 0x01, 0x00,
						}
					default:
						panic(fmt.Sprintf("unknown rejoin type `%d`", typ))
					}
					msg := &ttnpb.UplinkMessage{
						CorrelationIDs: []string{
							"rejoin-request-correlation-id-1",
							"rejoin-request-correlation-id-2",
							"rejoin-request-correlation-id-3",
						},
						RawPayload: phyPayload,
						RxMetadata: uplinkMDs,
						Settings:   MakeUplinkSettings(dr, ch.Frequency),
					}
					if decodePayload {
						var pld *ttnpb.RejoinRequestPayload
						switch typ {
						case ttnpb.RejoinType_CONTEXT, ttnpb.RejoinType_KEYS:
							pld = &ttnpb.RejoinRequestPayload{
								DevEUI:     *DevEUI.Copy(&types.EUI64{}),
								NetID:      *NetID.Copy(&types.NetID{}),
								RejoinCnt:  uint32(binary.LittleEndian.Uint16(phyPayload[13:14])),
								RejoinType: typ,
							}
						case ttnpb.RejoinType_SESSION:
							pld = &ttnpb.RejoinRequestPayload{
								JoinEUI:    *JoinEUI.Copy(&types.EUI64{}),
								DevEUI:     *DevEUI.Copy(&types.EUI64{}),
								RejoinCnt:  uint32(binary.LittleEndian.Uint16(phyPayload[18:19])),
								RejoinType: typ,
							}
						}
						msg.Payload = &ttnpb.Message{
							MHDR: ttnpb.MHDR{
								MType: ttnpb.MType_REJOIN_REQUEST,
								Major: ttnpb.Major_LORAWAN_R1,
							},
							MIC: phyPayload[len(phyPayload)-4:],
							Payload: &ttnpb.Message_RejoinRequestPayload{
								RejoinRequestPayload: pld,
							},
						}
					}
					return msg
				}

				tcs = append(tcs, TestCase{
					Name: makeMDName(makeChDRName(chIdx, drIdx, makeLoopName(fmt.Sprintf("Rejoin-request Type %d", typ)))),
					Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
						return assertions.New(test.MustTFromContext(ctx)).So(assertHandleUplinkResponse(ctx, handle(ctx, makeRejoinRequest(false)), ErrRejoinRequest), should.BeTrue)
					},
				})
			}
		})

		ForEachLoRaWANVersionPair(t, func(makeLoopName func(...string) string, macVersion ttnpb.MACVersion, phyVersion ttnpb.PHYVersion) {
			fpID := test.EUFrequencyPlanID
			fp := FrequencyPlan(fpID)
			phy := LoRaWANBands[fp.BandID][phyVersion]
			chIdx := uint8(len(phy.UplinkChannels) - 1)
			ch := phy.UplinkChannels[chIdx]
			drIdx := ch.MaxDataRate
			dr := phy.DataRates[drIdx].Rate

			makeDataRangeDevice := func(clock *test.MockClock, useADR bool) *ttnpb.EndDevice {
				sk := *MakeSessionKeys(macVersion, false)
				dev := &ttnpb.EndDevice{
					EndDeviceIdentifiers: *MakeOTAAIdentifiers(&DevAddr),
					FrequencyPlanID:      fpID,
					LoRaWANPHYVersion:    phyVersion,
					LoRaWANVersion:       macVersion,
					MACSettings: &ttnpb.MACSettings{
						UseADR: &pbtypes.BoolValue{
							Value: useADR,
						},
						Rx1Delay: &ttnpb.RxDelayValue{
							Value: Rx1Delay,
						},
					},
					Session: &ttnpb.Session{
						DevAddr:     *DevAddr.Copy(&types.DevAddr{}),
						SessionKeys: sk,
						LastFCntUp:  FCnt - 1,
						StartedAt:   clock.Now(),
						QueuedApplicationDownlinks: []*ttnpb.ApplicationDownlink{
							{
								SessionKeyID: sk.SessionKeyID,
								FPort:        FPort,
								FCnt:         42,
							},
						},
					},
					CreatedAt: clock.Now(),
					UpdatedAt: clock.Now(),
				}
				dev.MACState = test.Must(NewMACState(dev, frequencyplans.NewStore(test.FrequencyPlansFetcher), ttnpb.MACSettings{})).(*ttnpb.MACState)
				dev.MACState.CurrentParameters.ADRNbTrans = 2
				return dev
			}

			const matchIdx = 2
			makeDataRangeDevices := func(clock *test.MockClock, withMatch bool, useADR bool) []*ttnpb.EndDevice {
				withLastFCntUp := func(dev *ttnpb.EndDevice, fCnt uint32) *ttnpb.EndDevice {
					ret := CopyEndDevice(dev)
					ret.Session.LastFCntUp = fCnt
					return ret
				}
				match := makeDataRangeDevice(clock, useADR)
				rets := []*ttnpb.EndDevice{
					withLastFCntUp(match, FCnt+2),
					withLastFCntUp(match, FCnt+25),
				}
				if withMatch {
					rets = append(rets, match)
				}
				return append(rets,
					withLastFCntUp(match, FCnt+42),
				)
			}

			for _, confirmed := range [2]bool{true, false} {
				confirmed := confirmed
				makeDataUplink := func(decodePayload bool, adr bool, frmPayload []byte) *ttnpb.UplinkMessage {
					return MakeDataUplink(DataUplinkConfig{
						MACVersion:    macVersion,
						DecodePayload: decodePayload,
						Confirmed:     confirmed,
						DevAddr:       DevAddr,
						FCtrl:         ttnpb.FCtrl{ADR: adr},
						FCnt:          FCnt,
						FPort:         FPort,
						FRMPayload:    frmPayload,
						FOpts:         []byte{byte(ttnpb.CID_LINK_CHECK)},
						DataRate:      dr,
						DataRateIndex: drIdx,
						Frequency:     ch.Frequency,
						ChannelIndex:  chIdx,
						RxMetadata:    uplinkMDs,
					})
				}
				dataSetByIDSetPaths := [...]string{
					"mac_state",
					"pending_mac_state",
					"pending_session",
					"recent_adr_uplinks",
					"recent_uplinks",
					"session",
				}
				makeDataSetDevice := func(ctx context.Context, getDevice *ttnpb.EndDevice, decodedMsg *ttnpb.UplinkMessage) (*ttnpb.EndDevice, events.Builders) {
					setDevice := CopyEndDevice(getDevice)
					setDevice.MACState.QueuedResponses = nil
					evs := test.Must(HandleLinkCheckReq(test.Context(), setDevice, decodedMsg)).(events.Builders)
					setDevice.MACState.RecentUplinks = AppendRecentUplink(setDevice.MACState.RecentUplinks, WithMatchedUplinkSettings(decodedMsg, chIdx, drIdx), RecentUplinkCount)
					setDevice.MACState.RxWindowsAvailable = true
					setDevice.RecentUplinks = AppendRecentUplink(setDevice.RecentUplinks, WithMatchedUplinkSettings(decodedMsg, chIdx, drIdx), RecentUplinkCount)
					setDevice.Session.LastFCntUp = FCnt
					if !decodedMsg.Payload.GetMACPayload().ADR {
						setDevice.MACState.CurrentParameters.ADRDataRateIndex = ttnpb.DATA_RATE_0
						setDevice.MACState.CurrentParameters.ADRTxPowerIndex = 0
					}
					setDevice.MACState.DesiredParameters.ADRDataRateIndex = setDevice.MACState.CurrentParameters.ADRDataRateIndex
					setDevice.MACState.DesiredParameters.ADRTxPowerIndex = setDevice.MACState.CurrentParameters.ADRTxPowerIndex
					setDevice.MACState.DesiredParameters.ADRNbTrans = setDevice.MACState.CurrentParameters.ADRNbTrans
					if decodedMsg.Payload.GetMACPayload().ADR {
						setDevice.RecentADRUplinks = AppendRecentUplink(setDevice.RecentADRUplinks, WithMatchedUplinkSettings(decodedMsg, chIdx, drIdx), OptimalADRUplinkCount)
						test.Must(nil, AdaptDataRate(ctx, setDevice, phy, ttnpb.MACSettings{}))
					}
					return setDevice, evs
				}

				makeName := func(parts ...string) string {
					mTypeStr := "Unconfirmed Data"
					if confirmed {
						mTypeStr = "Confirmed Data"
					}
					return makeMDName(makeChDRName(chIdx, drIdx, makeLoopName(append([]string{mTypeStr}, parts...)...)))
				}
				for _, conf := range []struct {
					MakeName       func(...string) string
					MakeDataUplink func(bool) *ttnpb.UplinkMessage
				}{
					{
						MakeName: func(parts ...string) string {
							return makeName(append(parts, "No ADR", "No FRMPayload")...)
						},
						MakeDataUplink: func(decoded bool) *ttnpb.UplinkMessage {
							return makeDataUplink(decoded, false, nil)
						},
					},
					{
						MakeName: func(parts ...string) string {
							return makeName(append(parts, "No ADR", "FRMPayload")...)
						},
						MakeDataUplink: func(decoded bool) *ttnpb.UplinkMessage {
							return makeDataUplink(decoded, false, []byte("test"))
						},
					},
					{
						MakeName: func(parts ...string) string {
							return makeName(append(parts, "ADR", "No FRMPayload")...)
						},
						MakeDataUplink: func(decoded bool) *ttnpb.UplinkMessage {
							return makeDataUplink(decoded, true, nil)
						},
					},
					{
						MakeName: func(parts ...string) string {
							return makeName(append(parts, "ADR", "FRMPayload")...)
						},
						MakeDataUplink: func(decoded bool) *ttnpb.UplinkMessage {
							return makeDataUplink(decoded, true, []byte("test"))
						},
					},
				} {
					conf := conf
					tcs = append(tcs,
						TestCase{
							Name: conf.MakeName("Range fail"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									_, ok := assertDataRangeByAddr(ctx, env, msg.CorrelationIDs, ErrTestInternal)
									if !a.So(ok, should.BeTrue) {
										return false
									}
									return ok
								}, ErrTestInternal), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "No devices"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									_, ok := assertDataRangeByAddr(ctx, env, msg.CorrelationIDs, nil)
									if !a.So(ok, should.BeTrue) {
										return false
									}
									return ok
								}, ErrDeviceNotFound), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "No match"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									_, ok := assertDataRangeByAddr(ctx, env, msg.CorrelationIDs, nil, makeDataRangeDevices(clock, false, true)...)
									if !a.So(ok, should.BeTrue) {
										return false
									}
									return ok
								}, ErrDeviceNotFound), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "Deduplication fail"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, false, ErrTestInternal)
									return a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventsResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
											EvtDropDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, ErrTestInternal),
										),
									), should.BeTrue)
								}, ErrTestInternal), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "Duplicate uplink"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, false, nil)
									return a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventsResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
											EvtDropDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, ErrDuplicate),
										),
									), should.BeTrue)
								}, nil), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "First uplink", "Metadata merge fail", "Set fail on read"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, true, nil)
									return a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
										),
										assertAccumulatedMetadata(ctx, env, rangeCtx, decodedMsg, nil, ErrTestInternal),
										AssertDeviceRegistrySetByID(ctx, env.DeviceRegistry.SetByID, func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) bool {
											return test.AllTrue(
												a.So(ctx, should.HaveParentContextOrEqual, rangeCtx),
												a.So(appID, should.Resemble, AppID),
												a.So(devID, should.Resemble, DevID),
												a.So(gets, should.HaveSameElementsDeep, dataSetByIDGetPaths),
											)
										}, func(context.Context) DeviceRegistrySetByIDResponse {
											return DeviceRegistrySetByIDResponse{
												Error: ErrTestInternal,
											}
										}),
										a.So(env.Events, should.ReceiveEventResembling,
											EvtDropDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, ErrTestInternal),
										),
									), should.BeTrue)
								}, ErrTestInternal), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "First uplink", "Metadata merge success", "Device deleted during handling"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								innerErr := ErrOutdatedData
								registryErr := ErrTestInternal.WithCause(innerErr)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, true, nil)
									if !a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
										),
										assertAccumulatedMetadata(ctx, env, rangeCtx, decodedMsg, RxMetadata[:], nil),
									), should.BeTrue) {
										return false
									}
									decodedMsg.RxMetadata = RxMetadata[:]
									_, ok = assertDataSetByID(ctx, env, rangeCtx, nil, nil, nil, innerErr, registryErr)
									return a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtDropDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, registryErr),
										),
									), should.BeTrue)
								}, registryErr), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "First uplink", "Metadata merge success", "Rematch fail"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								innerErr := ErrOutdatedData.WithCause(ErrDeviceNotFound)
								registryErr := ErrTestInternal.WithCause(innerErr)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, true, nil)
									if !a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
										),
										assertAccumulatedMetadata(ctx, env, rangeCtx, decodedMsg, RxMetadata[:], nil),
									), should.BeTrue) {
										return false
									}
									decodedMsg.RxMetadata = RxMetadata[:]
									getDevice := CopyEndDevice(rangeDevices[matchIdx])
									getDevice.UpdatedAt = clock.Now()
									getDevice.MACState = nil
									getDevice.Session = nil
									_, ok = assertDataSetByID(ctx, env, rangeCtx, getDevice, nil, nil, innerErr, registryErr)
									return a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtDropDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, registryErr),
										),
									), should.BeTrue)
								}, registryErr), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "First uplink", "Metadata merge fail", "No rematch", "Set success", "NbTrans=1", "Downlink add fail", "Application uplink add fail"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, true, nil)
									if !a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
										),
										assertAccumulatedMetadata(ctx, env, rangeCtx, decodedMsg, nil, ErrTestInternal),
									), should.BeTrue) {
										return false
									}
									getDevice := CopyEndDevice(rangeDevices[matchIdx])
									setDevice, macEvs := makeDataSetDevice(ctx, getDevice, decodedMsg)
									setCtx, ok := assertDataSetByID(ctx, env, rangeCtx, getDevice, setDevice, dataSetByIDSetPaths[:], nil, nil)
									return a.So(test.AllTrue(
										ok,
										func() bool {
											if len(uplinkMDs) == 0 {
												// No downlink task should be added if no downlink paths are available.
												return true
											}
											return assertDownlinkTaskAdd(ctx, env, setCtx, setDevice.EndDeviceIdentifiers, clock.Now().Add(NSScheduleWindow()), true, ErrTestInternal)
										}(),
										assertDataApplicationUp(ctx, env, setCtx, setDevice, decodedMsg, ErrTestInternal),
										a.So(env.Events, should.ReceiveEventsResembling,
											macEvs.New(setCtx, events.WithIdentifiers(setDevice.EndDeviceIdentifiers)),
											EvtProcessDataUplink.NewWithIdentifiersAndData(setCtx, setDevice.EndDeviceIdentifiers, decodedMsg),
										),
									), should.BeTrue)
								}, nil), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "First uplink", "Metadata merge success", "Rematch success", "Set success", "NbTrans=1", "Downlink add success", "Application uplink add success"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, true, nil)
									if !a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
										),
										assertAccumulatedMetadata(ctx, env, rangeCtx, decodedMsg, RxMetadata[:], nil),
									), should.BeTrue) {
										return false
									}
									decodedMsg.RxMetadata = RxMetadata[:]
									getDevice := CopyEndDevice(rangeDevices[matchIdx])
									getDevice.UpdatedAt = clock.Now()
									setDevice, macEvs := makeDataSetDevice(ctx, getDevice, decodedMsg)
									setCtx, ok := assertDataSetByID(ctx, env, rangeCtx, getDevice, setDevice, dataSetByIDSetPaths[:], nil, nil)
									return a.So(test.AllTrue(
										ok,
										assertDownlinkTaskAdd(ctx, env, setCtx, setDevice.EndDeviceIdentifiers, clock.Now().Add(NSScheduleWindow()), true, nil),
										assertDataApplicationUp(ctx, env, setCtx, setDevice, WithMatchedUplinkSettings(decodedMsg, chIdx, drIdx), nil),
										a.So(env.Events, should.ReceiveEventsResembling,
											macEvs.New(setCtx, events.WithIdentifiers(setDevice.EndDeviceIdentifiers)),
											EvtProcessDataUplink.NewWithIdentifiersAndData(setCtx, setDevice.EndDeviceIdentifiers, decodedMsg),
										),
									), should.BeTrue)
								}, nil), should.BeTrue)
							},
						},
						TestCase{
							Name: conf.MakeName("Range success", "First uplink", "Metadata merge success", "No rematch", "Set success", "NbTrans=2", "Downlink add success"),
							Handler: func(ctx context.Context, env LegacyTestEnvironment, clock *test.MockClock, handle func(context.Context, *ttnpb.UplinkMessage) <-chan error) bool {
								t := test.MustTFromContext(ctx)
								a := assertions.New(t)
								msg := conf.MakeDataUplink(false)
								decodedMsg := conf.MakeDataUplink(true)
								rangeDevices := makeDataRangeDevices(clock, true, true)
								prevMsg := CopyUplinkMessage(decodedMsg)
								prevMsg.ReceivedAt = clock.Now()
								rangeDevice, _ := makeDataSetDevice(ctx, rangeDevices[matchIdx], prevMsg)
								rangeDevices[matchIdx] = rangeDevice
								return a.So(assertHandleUplink(ctx, handle, msg, func() bool {
									rangeCtx, ok := assertDataDeduplicateSequence(ctx, env, clock, decodedMsg, chIdx, drIdx, rangeDevices, matchIdx, true, nil)
									if !a.So(test.AllTrue(
										ok,
										a.So(env.Events, should.ReceiveEventResembling,
											EvtReceiveDataUplink.NewWithIdentifiersAndData(rangeCtx, rangeDevices[matchIdx].EndDeviceIdentifiers, decodedMsg),
										),
										assertAccumulatedMetadata(ctx, env, rangeCtx, decodedMsg, RxMetadata[:], nil),
									), should.BeTrue) {
										return false
									}
									decodedMsg.RxMetadata = RxMetadata[:]
									setDevice, macEvs := makeDataSetDevice(ctx, rangeDevice, decodedMsg)
									setCtx, ok := assertDataSetByID(ctx, env, rangeCtx, rangeDevice, setDevice, dataSetByIDSetPaths[:], nil, nil)
									return a.So(test.AllTrue(
										ok,
										assertDownlinkTaskAdd(ctx, env, setCtx, setDevice.EndDeviceIdentifiers, clock.Now().Add(NSScheduleWindow()), true, ErrTestInternal),
										a.So(env.Events, should.ReceiveEventsResembling,
											macEvs.New(setCtx, events.WithIdentifiers(setDevice.EndDeviceIdentifiers)),
											EvtProcessDataUplink.NewWithIdentifiersAndData(setCtx, setDevice.EndDeviceIdentifiers, decodedMsg),
										),
									), should.BeTrue)
								}, nil), should.BeTrue)
							},
						},
					)
				}
			}
		})
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			dtq, dtqEnv, dtqClose := newMockDownlinkTaskQueue(t)
			defer dtqClose()

			dr, drEnv, drClose := newMockDeviceRegistry(t)
			defer drClose()

			ud, udEnv, udClose := newMockUplinkDeduplicator(t)
			defer udClose()

			auq, auqEnv, auqClose := newMockApplicationUplinkQueue(t)
			defer auqClose()

			ns, ctx, env, stop := StartTest(
				t,
				TestConfig{
					NetworkServer: Config{
						NetID:               *NetID.Copy(&types.NetID{}),
						DeduplicationWindow: DeduplicationWindow,
						CooldownWindow:      CooldownWindow,

						DownlinkTasks:      dtq,
						Devices:            dr,
						UplinkDeduplicator: ud,
						ApplicationUplinkQueue: ApplicationUplinkQueueConfig{
							Queue: auq,
						},
					},
					Timeout: (1 << 10) * test.Delay,
				},
			)
			defer stop()

			<-dtqEnv.Pop

			clock := test.NewMockClock(time.Now().UTC())
			defer SetMockClock(clock)()

			if !tc.Handler(ctx, LegacyTestEnvironment{
				TestEnvironment: env,

				DownlinkTasks:      dtqEnv,
				DeviceRegistry:     drEnv,
				UplinkDeduplicator: udEnv,
				ApplicationUplinks: auqEnv,
			}, clock, func(ctx context.Context, msg *ttnpb.UplinkMessage) <-chan error {
				ch := make(chan error)
				go func() {
					_, err := ttnpb.NewGsNsClient(ns.LoopbackConn()).HandleUplink(ctx, CopyUplinkMessage(msg))
					ttnErr, ok := errors.From(err)
					if ok {
						ch <- ttnErr
					} else {
						ch <- err
					}
					close(ch)
				}()
				return ch
			}) {
				t.Error("Test handler failed")
			}
			assertions.New(t).So(AssertNetworkServerClose(ctx, ns), should.BeTrue)
		})
	}
}
