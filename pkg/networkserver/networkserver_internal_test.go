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

package networkserver

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/kr/pretty"
	"github.com/mohae/deepcopy"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/crypto"
	"go.thethings.network/lorawan-stack/pkg/deviceregistry"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/store"
	"go.thethings.network/lorawan-stack/pkg/store/mapstore"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
)

const (
	RecentUplinkCount = recentUplinkCount
)

var (
	NewMACState = newMACState
)

func TestAccumulator(t *testing.T) {
	a := assertions.New(t)

	acc := &metadataAccumulator{newAccumulator()}
	a.So(func() { acc.Add() }, should.NotPanic)

	vals := []*ttnpb.RxMetadata{
		ttnpb.NewPopulatedRxMetadata(test.Randy, false),
		ttnpb.NewPopulatedRxMetadata(test.Randy, false),
		nil,
		ttnpb.NewPopulatedRxMetadata(test.Randy, false),
		ttnpb.NewPopulatedRxMetadata(test.Randy, false),
		ttnpb.NewPopulatedRxMetadata(test.Randy, false),
		ttnpb.NewPopulatedRxMetadata(test.Randy, false),
	}

	a.So(acc.Accumulated(), should.BeEmpty)

	acc.Add(vals[0], vals[1], vals[2])
	for _, v := range vals[:3] {
		a.So(acc.Accumulated(), should.Contain, v)
	}

	for i := 2; i < len(vals); i++ {
		acc.Add(vals[i])
		for _, v := range vals[:i] {
			a.So(acc.Accumulated(), should.Contain, v)
		}
	}

	acc.Reset()
	a.So(acc.Accumulated(), should.BeEmpty)
}

var _ ttnpb.NsGsClient = &MockNsGsClient{}

type MockNsGsClient struct {
	*test.MockClientStream
	ScheduleDownlinkFunc func(ctx context.Context, in *ttnpb.DownlinkMessage, opts ...grpc.CallOption) (*pbtypes.Empty, error)
}

func (cl *MockNsGsClient) ScheduleDownlink(ctx context.Context, in *ttnpb.DownlinkMessage, opts ...grpc.CallOption) (*pbtypes.Empty, error) {
	if cl.ScheduleDownlinkFunc == nil {
		return nil, nil
	}
	return cl.ScheduleDownlinkFunc(ctx, in, opts...)
}

func TestScheduleDownlink(t *testing.T) {
	newRX2 := func(down *ttnpb.ApplicationDownlink, fp ttnpb.FrequencyPlan, dev *ttnpb.EndDevice) *ttnpb.DownlinkMessage {
		band := test.Must(band.GetByID(fp.BandID)).(band.Band)
		st := dev.MACState
		drIdx := st.Rx2DataRateIndex

		msg := &ttnpb.DownlinkMessage{
			EndDeviceIdentifiers: dev.EndDeviceIdentifiers,
			Settings: ttnpb.TxSettings{
				DataRateIndex:         drIdx,
				CodingRate:            "4/5",
				PolarizationInversion: true,
				Frequency:             st.Rx2Frequency,
				TxPower:               int32(band.DefaultMaxEIRP),
			},
			RawPayload: test.Must((ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_DOWN,
				},
				Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{
					FHDR: ttnpb.FHDR{
						DevAddr: *dev.EndDeviceIdentifiers.DevAddr,
						FCtrl: ttnpb.FCtrl{
							ADR:      false,
							Ack:      false,
							FPending: false,
						},
						FCnt:  down.FCnt,
						FOpts: nil,
					},
					FPort:      down.FPort,
					FRMPayload: down.FRMPayload,
				}},
			}).MarshalLoRaWAN()).([]byte),
		}

		mic := test.Must(crypto.ComputeDownlinkMIC(
			*dev.Session.SessionKeys.SNwkSIntKey.Key,
			*dev.EndDeviceIdentifiers.DevAddr,
			down.FCnt,
			msg.RawPayload,
		)).([4]byte)
		msg.RawPayload = append(msg.RawPayload, mic[:]...)

		test.Must(nil, setDownlinkModulation(&msg.Settings, band.DataRates[drIdx]))
		return msg
	}

	newRX1 := func(down *ttnpb.ApplicationDownlink, fp ttnpb.FrequencyPlan, dev *ttnpb.EndDevice, up *ttnpb.UplinkMessage) *ttnpb.DownlinkMessage {
		msg := newRX2(down, fp, dev)

		sets := up.Settings
		st := dev.MACState
		band := test.Must(band.GetByID(fp.BandID)).(band.Band)
		drIdx := test.Must(band.Rx1DataRate(sets.DataRateIndex, st.Rx1DataRateOffset, st.DownlinkDwellTime)).(uint32)

		msg.Settings.ChannelIndex = test.Must(band.Rx1Channel(sets.ChannelIndex)).(uint32)
		msg.Settings.Frequency = uint64(fp.Channels[msg.Settings.ChannelIndex].Frequency)
		msg.Settings.DataRateIndex = drIdx

		test.Must(nil, setDownlinkModulation(&msg.Settings, band.DataRates[drIdx]))
		return msg
	}

	t.Run("Empty queue", func(t *testing.T) {
		a := assertions.New(t)

		reg := deviceregistry.New(store.NewTypedMapStoreClient(mapstore.New()))
		ns := test.Must(New(
			component.MustNew(test.GetLogger(t), &component.Config{}),
			&Config{
				Registry:            reg,
				JoinServers:         nil,
				DeduplicationWindow: 42,
				CooldownWindow:      42,
			},
		)).(*NetworkServer)
		test.Must(nil, ns.Start())

		ed := ttnpb.NewPopulatedEndDevice(test.Randy, false)
		for ed.Session == nil || len(ed.RecentUplinks) == 0 {
			ed = ttnpb.NewPopulatedEndDevice(test.Randy, false)
		}
		ed.QueuedApplicationDownlinks = nil
		dev := test.Must(reg.Create(ed)).(*deviceregistry.Device)

		err := ns.scheduleApplicationDownlink(context.Background(), dev, nil, nil)
		a.So(err, should.BeNil)
	})

	t.Run("No recent uplinks", func(t *testing.T) {
		a := assertions.New(t)

		reg := deviceregistry.New(store.NewTypedMapStoreClient(mapstore.New()))
		ns := test.Must(New(
			component.MustNew(test.GetLogger(t), &component.Config{}),
			&Config{
				Registry:            reg,
				JoinServers:         nil,
				DeduplicationWindow: 42,
				CooldownWindow:      42,
			},
		)).(*NetworkServer)
		test.Must(nil, ns.Start())

		ed := ttnpb.NewPopulatedEndDevice(test.Randy, false)
		for ed.Session == nil || len(ed.QueuedApplicationDownlinks) == 0 {
			ed = ttnpb.NewPopulatedEndDevice(test.Randy, false)
		}
		ed.RecentUplinks = nil
		dev := test.Must(reg.Create(ed)).(*deviceregistry.Device)

		err := ns.scheduleApplicationDownlink(context.Background(), dev, nil, nil)
		a.So(err, should.BeError)
	})

	t.Run("No GS matching", func(t *testing.T) {
		a := assertions.New(t)

		reg := deviceregistry.New(store.NewTypedMapStoreClient(mapstore.New()))
		scheduleCtx := context.Background()

		ns := test.Must(New(
			component.MustNew(test.GetLogger(t), &component.Config{}),
			&Config{
				Registry:            reg,
				JoinServers:         nil,
				DeduplicationWindow: 42,
				CooldownWindow:      42,
			},
			WithNsGsClientFunc(func(ctx context.Context, id ttnpb.GatewayIdentifiers) (ttnpb.NsGsClient, error) {
				a.So(ctx, should.Resemble, scheduleCtx)
				return nil, errors.New("Test")
			}),
		)).(*NetworkServer)
		test.Must(nil, ns.Start())

		ed := ttnpb.NewPopulatedEndDevice(test.Randy, false)
		for ed.Session == nil || len(ed.RecentUplinks) == 0 {
			ed = ttnpb.NewPopulatedEndDevice(test.Randy, false)
		}
		ed.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{
			ttnpb.NewPopulatedApplicationDownlink(test.Randy, false),
		}
		dev := test.Must(reg.Create(ed)).(*deviceregistry.Device)

		err := ns.scheduleApplicationDownlink(context.Background(), dev, nil, nil)
		a.So(err, should.BeError)
	})

	t.Run("Rx1", func(t *testing.T) {
		a := assertions.New(t)

		reg := deviceregistry.New(store.NewTypedMapStoreClient(mapstore.New()))
		gateways := make(map[string]ttnpb.NsGsClient)
		scheduleCtx := context.Background()
		wg := &sync.WaitGroup{}

		fpStore, err := test.NewFrequencyPlansStore()
		if !a.So(err, should.BeNil) {
			t.FailNow()
		}
		defer fpStore.Destroy()

		ns := test.Must(New(
			component.MustNew(test.GetLogger(t), &component.Config{ServiceBase: config.ServiceBase{
				FrequencyPlans: config.FrequencyPlans{
					StoreDirectory: fpStore.Directory(),
				}}}),
			&Config{
				Registry:            reg,
				JoinServers:         nil,
				DeduplicationWindow: 42,
				CooldownWindow:      42,
			},
			WithNsGsClientFunc(func(ctx context.Context, id ttnpb.GatewayIdentifiers) (ttnpb.NsGsClient, error) {
				defer wg.Done()
				a.So(ctx, should.Resemble, scheduleCtx)

				cl, ok := gateways[id.UniqueID(ctx)]
				if !ok {
					t.Error("Non-existing gateway lookup")
					return nil, errors.New("Not found")
				}
				return cl, nil
			}),
		)).(*NetworkServer)
		test.Must(nil, ns.Start())

		down := ttnpb.NewPopulatedApplicationDownlink(test.Randy, false)

		ed := ttnpb.NewPopulatedEndDevice(test.Randy, false)
		for ed.Session == nil {
			ed = ttnpb.NewPopulatedEndDevice(test.Randy, false)
		}
		ed.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{
			down,
		}
		dev := test.Must(reg.Create(ed)).(*deviceregistry.Device)

		up := ttnpb.NewPopulatedUplinkMessageUplink(
			test.Randy,
			*types.NewPopulatedAES128Key(test.Randy),
			*types.NewPopulatedAES128Key(test.Randy),
			false,
		)

		mds := append(make([]*ttnpb.RxMetadata, 0), up.RxMetadata...)
		sort.SliceStable(mds, func(i, j int) bool {
			return mds[i].SNR > mds[j].SNR
		})

		wg.Add(len(mds))

		slots := []*ttnpb.DownlinkMessage{
			newRX1(
				down,
				test.Must(ns.Component.FrequencyPlans.GetByID(dev.FrequencyPlanID)).(ttnpb.FrequencyPlan),
				dev.EndDevice,
				up,
			),
		}

		var cnt uint32
		for i, md := range mds {
			i, md := i, deepcopy.Copy(md).(*ttnpb.RxMetadata)

			slots := deepcopy.Copy(slots).([]*ttnpb.DownlinkMessage)
			slots[0].TxMetadata = ttnpb.TxMetadata{
				GatewayIdentifiers: md.GatewayIdentifiers,
				Timestamp:          md.Timestamp + uint64(time.Duration(dev.MACState.RxDelay)*time.Second),
			}

			var n uint32
			gateways[md.GatewayIdentifiers.UniqueID(scheduleCtx)] = &MockNsGsClient{
				MockClientStream: &test.MockClientStream{},
				ScheduleDownlinkFunc: func(ctx context.Context, in *ttnpb.DownlinkMessage, opts ...grpc.CallOption) (*pbtypes.Empty, error) {
					defer atomic.AddUint32(&cnt, 1)
					defer atomic.AddUint32(&n, 1)

					a.So(ctx, should.Resemble, scheduleCtx)
					a.So(cnt, should.Equal, i+int(n)*len(mds))
					if !a.So(n, should.BeLessThan, len(slots)) {
						return nil, errors.New("Too many slots scheduled")
					}
					if !a.So(in, should.Resemble, slots[n]) {
						pretty.Ldiff(t, in, slots[n])
					}

					if i == len(mds)-1 && n == uint32(len(slots)-1) {
						return nil, nil
					}
					return nil, errors.New("Test")
				},
			}
		}

		err = ns.scheduleApplicationDownlink(context.Background(), dev, up, nil)
		a.So(err, should.BeNil)
		a.So(cnt, should.Equal, len(mds)*len(slots))
		a.So(test.WaitTimeout(20*test.Delay, wg.Wait), should.BeTrue)
	})

	t.Run("Rx2", func(t *testing.T) {
		a := assertions.New(t)

		reg := deviceregistry.New(store.NewTypedMapStoreClient(mapstore.New()))
		gateways := make(map[string]ttnpb.NsGsClient)
		scheduleCtx := context.Background()
		wg := &sync.WaitGroup{}

		fpStore, err := test.NewFrequencyPlansStore()
		if !a.So(err, should.BeNil) {
			t.FailNow()
		}
		defer fpStore.Destroy()

		ns := test.Must(New(
			component.MustNew(test.GetLogger(t), &component.Config{ServiceBase: config.ServiceBase{
				FrequencyPlans: config.FrequencyPlans{
					StoreDirectory: fpStore.Directory(),
				}}}),
			&Config{
				Registry:            reg,
				JoinServers:         nil,
				DeduplicationWindow: 42,
				CooldownWindow:      42,
			},
			WithNsGsClientFunc(func(ctx context.Context, id ttnpb.GatewayIdentifiers) (ttnpb.NsGsClient, error) {
				defer wg.Done()
				a.So(ctx, should.Resemble, scheduleCtx)

				cl, ok := gateways[id.UniqueID(ctx)]
				if !ok {
					t.Error("Non-existing gateway lookup")
					return nil, errors.New("Not found")
				}
				return cl, nil
			}),
		)).(*NetworkServer)
		test.Must(nil, ns.Start())

		down := ttnpb.NewPopulatedApplicationDownlink(test.Randy, false)

		ed := ttnpb.NewPopulatedEndDevice(test.Randy, false)
		for ed.Session == nil {
			ed = ttnpb.NewPopulatedEndDevice(test.Randy, false)
		}
		ed.QueuedApplicationDownlinks = []*ttnpb.ApplicationDownlink{
			down,
		}
		dev := test.Must(reg.Create(ed)).(*deviceregistry.Device)

		up := ttnpb.NewPopulatedUplinkMessageUplink(
			test.Randy,
			*types.NewPopulatedAES128Key(test.Randy),
			*types.NewPopulatedAES128Key(test.Randy),
			false,
		)

		mds := append(make([]*ttnpb.RxMetadata, 0), up.RxMetadata...)
		sort.SliceStable(mds, func(i, j int) bool {
			return mds[i].SNR > mds[j].SNR
		})

		wg.Add(len(mds) * 2)

		slots := []*ttnpb.DownlinkMessage{
			newRX1(
				down,
				test.Must(ns.Component.FrequencyPlans.GetByID(dev.FrequencyPlanID)).(ttnpb.FrequencyPlan),
				dev.EndDevice,
				up,
			),
			newRX2(
				down,
				test.Must(ns.Component.FrequencyPlans.GetByID(dev.FrequencyPlanID)).(ttnpb.FrequencyPlan),
				dev.EndDevice,
			),
		}

		var cnt uint32
		for i, md := range mds {
			i, md := i, deepcopy.Copy(md).(*ttnpb.RxMetadata)

			slots := deepcopy.Copy(slots).([]*ttnpb.DownlinkMessage)
			slots[0].TxMetadata = ttnpb.TxMetadata{
				GatewayIdentifiers: md.GatewayIdentifiers,
				Timestamp:          md.Timestamp + uint64(time.Duration(dev.MACState.RxDelay)*time.Second),
			}

			slots[1].TxMetadata = slots[0].TxMetadata
			slots[1].TxMetadata.Timestamp += uint64(time.Second.Nanoseconds())

			var n uint32
			gateways[md.GatewayIdentifiers.UniqueID(scheduleCtx)] = &MockNsGsClient{
				MockClientStream: &test.MockClientStream{},
				ScheduleDownlinkFunc: func(ctx context.Context, in *ttnpb.DownlinkMessage, opts ...grpc.CallOption) (*pbtypes.Empty, error) {
					defer atomic.AddUint32(&cnt, 1)
					defer atomic.AddUint32(&n, 1)

					a.So(ctx, should.Resemble, scheduleCtx)
					a.So(cnt, should.Equal, i+int(n)*len(mds))
					if !a.So(n, should.BeLessThan, len(slots)) {
						return nil, errors.New("Too many slots scheduled")
					}
					if !a.So(in, should.Resemble, slots[n]) {
						pretty.Ldiff(t, in, slots[n])
					}

					if i == len(mds)-1 && n == uint32(len(slots)-1) {
						return nil, nil
					}
					return nil, errors.New("Test")
				},
			}
		}

		err = ns.scheduleApplicationDownlink(context.Background(), dev, up, nil)
		a.So(err, should.BeNil)
		a.So(cnt, should.Equal, len(mds)*len(slots))
		a.So(test.WaitTimeout(20*test.Delay, wg.Wait), should.BeTrue)
	})
}
