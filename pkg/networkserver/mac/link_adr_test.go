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

package mac_test

import (
	"context"
	"fmt"
	"testing"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal/test"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal/test/shared"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/mac"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestLinkADRReq(t *testing.T) {
	for _, tc := range []struct {
		Name                                             string
		BandID                                           string
		LoRaWANVersion                                   ttnpb.MACVersion
		LoRaWANPHYVersion                                ttnpb.PHYVersion
		CurrentChannels, DesiredChannels                 []*ttnpb.MACParameters_Channel
		CurrentADRDataRateIndex, DesiredADRDataRateIndex ttnpb.DataRateIndex
		CurrentADRTxPowerIndex, DesiredADRTxPowerIndex   uint32
		CurrentADRNbTrans, DesiredADRNbTrans             uint32
		RejectedADRDataRateIndexes                       []ttnpb.DataRateIndex
		RejectedADRTxPowerIndexes                        []uint32
		NoADRCommands, ADRCommands                       []*ttnpb.MACCommand_LinkADRReq
		NoADRErrorAssertion, ADRErrorAssertion           func(*testing.T, error) bool
	}{
		{
			Name:              "no channels",
			BandID:            band.US_902_928,
			LoRaWANVersion:    ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion: ttnpb.PHY_V1_0_3_REV_A,
			CurrentADRNbTrans: 1,
			DesiredADRNbTrans: 1,
		},
		{
			Name:              "invalid channel",
			BandID:            band.US_902_928,
			LoRaWANVersion:    ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion: ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:   MakeDefaultUS915FSB2DesiredChannels(),
			CurrentADRNbTrans: 1,
			DesiredChannels: []*ttnpb.MACParameters_Channel{
				{EnableUplink: true},
			},
			DesiredADRNbTrans: 1,
			NoADRErrorAssertion: func(t *testing.T, err error) bool {
				a, _ := test.New(t)
				return a.So(err, should.BeError)
			},
		},
		{
			Name:              "invalid channel count",
			BandID:            band.EU_863_870,
			LoRaWANVersion:    ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion: ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:   MakeDefaultEU868CurrentChannels(),
			CurrentADRNbTrans: 1,
			DesiredChannels:   MakeDefaultUS915FSB2DesiredChannels(),
			DesiredADRNbTrans: 1,
			NoADRErrorAssertion: func(t *testing.T, err error) bool {
				a, _ := test.New(t)
				return a.So(err, should.BeError)
			},
		},
		{
			Name:              "invalid band channels",
			BandID:            band.EU_863_870,
			LoRaWANVersion:    ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion: ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:   MakeDefaultUS915CurrentChannels(),
			CurrentADRNbTrans: 1,
			DesiredChannels:   MakeDefaultUS915FSB2DesiredChannels(),
			DesiredADRNbTrans: 1,
			NoADRErrorAssertion: func(t *testing.T, err error) bool {
				a, _ := test.New(t)
				return a.So(err, should.BeError)
			},
		},
		{
			Name:                    "non-existent data rate",
			BandID:                  band.EU_863_870,
			LoRaWANVersion:          ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion:       ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:         MakeDefaultEU868DesiredChannels(),
			CurrentADRNbTrans:       1,
			DesiredChannels:         MakeDefaultEU868DesiredChannels(),
			DesiredADRDataRateIndex: ttnpb.DATA_RATE_15,
			DesiredADRNbTrans:       1,
			ADRErrorAssertion: func(t *testing.T, err error) bool {
				a, _ := test.New(t)
				return a.So(err, should.BeError)
			},
		},
		{
			Name:                    "data rate too low",
			BandID:                  band.EU_863_870,
			LoRaWANVersion:          ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion:       ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:         MakeDefaultEU868DesiredChannels(),
			CurrentADRNbTrans:       1,
			CurrentADRDataRateIndex: ttnpb.DATA_RATE_2,
			DesiredChannels:         MakeDefaultEU868DesiredChannels(),
			DesiredADRDataRateIndex: ttnpb.DATA_RATE_1,
			DesiredADRNbTrans:       1,
			ADRErrorAssertion: func(t *testing.T, err error) bool {
				a, _ := test.New(t)
				return a.So(err, should.BeError)
			},
		},
		{
			Name:                   "TX power too high",
			BandID:                 band.EU_863_870,
			LoRaWANVersion:         ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion:      ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:        MakeDefaultEU868DesiredChannels(),
			CurrentADRNbTrans:      1,
			DesiredChannels:        MakeDefaultEU868DesiredChannels(),
			DesiredADRTxPowerIndex: 14,
			DesiredADRNbTrans:      1,
			ADRErrorAssertion: func(t *testing.T, err error) bool {
				a, _ := test.New(t)
				return a.So(err, should.BeError)
			},
		},
		{
			Name:              "ABP channel setup",
			BandID:            band.US_902_928,
			LoRaWANVersion:    ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion: ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:   MakeDefaultUS915CurrentChannels(),
			CurrentADRNbTrans: 1,
			DesiredChannels:   MakeDefaultUS915FSB2DesiredChannels(),
			DesiredADRNbTrans: 1,
			NoADRCommands: []*ttnpb.MACCommand_LinkADRReq{
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						false, false, false, false, false, false, false, false,
					},
					ChannelMaskControl: 7,
					NbTrans:            1,
				},
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						true, true, true, true, true, true, true, true,
					},
					NbTrans: 1,
				},
			},
		},
		{
			Name:              "ABP channel setup",
			BandID:            band.US_902_928,
			LoRaWANVersion:    ttnpb.MAC_V1_1,
			LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
			CurrentChannels:   MakeDefaultUS915CurrentChannels(),
			CurrentADRNbTrans: 1,
			DesiredChannels:   MakeDefaultUS915FSB2DesiredChannels(),
			DesiredADRNbTrans: 1,
			NoADRCommands: []*ttnpb.MACCommand_LinkADRReq{
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						false, false, false, false, false, false, false, false,
					},
					ChannelMaskControl: 7,
					DataRateIndex:      ttnpb.DATA_RATE_15,
					TxPowerIndex:       15,
					NbTrans:            1,
				},
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						true, true, true, true, true, true, true, true,
					},
					DataRateIndex: ttnpb.DATA_RATE_15,
					TxPowerIndex:  15,
					NbTrans:       1,
				},
			},
		},
		{
			Name:                    "ADR",
			BandID:                  band.US_902_928,
			LoRaWANVersion:          ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion:       ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:         MakeDefaultUS915FSB2DesiredChannels(),
			CurrentADRNbTrans:       1,
			DesiredChannels:         MakeDefaultUS915FSB2DesiredChannels(),
			DesiredADRDataRateIndex: ttnpb.DATA_RATE_3,
			DesiredADRTxPowerIndex:  1,
			DesiredADRNbTrans:       2,
			RejectedADRDataRateIndexes: []ttnpb.DataRateIndex{
				ttnpb.DATA_RATE_2,
			},
			RejectedADRTxPowerIndexes: []uint32{
				0,
				1,
			},
			ADRCommands: []*ttnpb.MACCommand_LinkADRReq{
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						true, true, true, true, true, true, true, true,
					},
					DataRateIndex: ttnpb.DATA_RATE_1,
					TxPowerIndex:  15,
					NbTrans:       2,
				},
			},
		},
		{
			Name:                    "ADR",
			BandID:                  band.EU_863_870,
			LoRaWANVersion:          ttnpb.MAC_V1_0_1,
			LoRaWANPHYVersion:       ttnpb.PHY_V1_0_1,
			CurrentADRDataRateIndex: ttnpb.DATA_RATE_1,
			CurrentADRNbTrans:       1,
			CurrentChannels:         MakeDefaultEU868DesiredChannels(),
			DesiredADRDataRateIndex: ttnpb.DATA_RATE_5,
			DesiredADRNbTrans:       2,
			DesiredADRTxPowerIndex:  3,
			DesiredChannels:         MakeDefaultEU868DesiredChannels(),
			RejectedADRDataRateIndexes: []ttnpb.DataRateIndex{
				ttnpb.DATA_RATE_1,
				ttnpb.DATA_RATE_2,
				ttnpb.DATA_RATE_3,
				ttnpb.DATA_RATE_4,
			},
			ADRCommands: []*ttnpb.MACCommand_LinkADRReq{
				{
					ChannelMask: []bool{
						true, true, true, true, true, true, true, true,
						false, false, false, false, false, false, false, false,
					},
					DataRateIndex: ttnpb.DATA_RATE_5,
					TxPowerIndex:  3,
					NbTrans:       2,
				},
			},
		},
		{
			Name:                    "ADR",
			BandID:                  band.EU_863_870,
			LoRaWANVersion:          ttnpb.MAC_V1_0_1,
			LoRaWANPHYVersion:       ttnpb.PHY_V1_0_1,
			CurrentADRDataRateIndex: ttnpb.DATA_RATE_1,
			CurrentADRNbTrans:       1,
			CurrentChannels:         MakeDefaultEU868DesiredChannels(),
			DesiredADRDataRateIndex: ttnpb.DATA_RATE_5,
			DesiredADRNbTrans:       1,
			DesiredADRTxPowerIndex:  3,
			DesiredChannels:         MakeDefaultEU868DesiredChannels(),
			RejectedADRDataRateIndexes: []ttnpb.DataRateIndex{
				ttnpb.DATA_RATE_1,
				ttnpb.DATA_RATE_2,
				ttnpb.DATA_RATE_3,
				ttnpb.DATA_RATE_4,
				ttnpb.DATA_RATE_5,
			},
		},
		{
			Name:                    "ABP channel setup + ADR",
			BandID:                  band.US_902_928,
			LoRaWANVersion:          ttnpb.MAC_V1_0_3,
			LoRaWANPHYVersion:       ttnpb.PHY_V1_0_3_REV_A,
			CurrentChannels:         MakeDefaultUS915CurrentChannels(),
			CurrentADRNbTrans:       1,
			CurrentADRDataRateIndex: ttnpb.DATA_RATE_1,
			DesiredChannels:         MakeDefaultUS915FSB2DesiredChannels(),
			DesiredADRNbTrans:       2,
			DesiredADRDataRateIndex: ttnpb.DATA_RATE_2,
			DesiredADRTxPowerIndex:  3,
			NoADRCommands: []*ttnpb.MACCommand_LinkADRReq{
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						false, false, false, false, false, false, false, false,
					},
					ChannelMaskControl: 7,
					DataRateIndex:      ttnpb.DATA_RATE_1,
					NbTrans:            1,
				},
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						true, true, true, true, true, true, true, true,
					},
					DataRateIndex: ttnpb.DATA_RATE_1,
					NbTrans:       1,
				},
			},
			ADRCommands: []*ttnpb.MACCommand_LinkADRReq{
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						false, false, false, false, false, false, false, false,
					},
					ChannelMaskControl: 7,
					DataRateIndex:      ttnpb.DATA_RATE_2,
					TxPowerIndex:       3,
					NbTrans:            2,
				},
				{
					ChannelMask: []bool{
						false, false, false, false, false, false, false, false,
						true, true, true, true, true, true, true, true,
					},
					DataRateIndex: ttnpb.DATA_RATE_2,
					TxPowerIndex:  3,
					NbTrans:       2,
				},
			},
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name: fmt.Sprintf("%s/band:%s,MAC:%s,PHY:%s,DR:%d->%d,TX:%d->%d,NB:%d->%d,rejected_DR:%v,rejected_TX:%v",
				tc.Name,
				tc.BandID,
				tc.LoRaWANVersion,
				tc.LoRaWANPHYVersion,
				tc.CurrentADRDataRateIndex, tc.DesiredADRDataRateIndex,
				tc.CurrentADRTxPowerIndex, tc.DesiredADRTxPowerIndex,
				tc.CurrentADRNbTrans, tc.DesiredADRNbTrans,
				fmt.Sprintf("[%s]", test.JoinStringsf("%d", ",", false, tc.RejectedADRDataRateIndexes)),
				fmt.Sprintf("[%s]", test.JoinStringsf("%d", ",", false, tc.RejectedADRTxPowerIndexes)),
			),
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				makeDevice := func(adr bool) *ttnpb.EndDevice {
					return CopyEndDevice(&ttnpb.EndDevice{
						MACState: &ttnpb.MACState{
							LoRaWANVersion: tc.LoRaWANVersion,
							CurrentParameters: ttnpb.MACParameters{
								Channels:         tc.CurrentChannels,
								ADRDataRateIndex: tc.CurrentADRDataRateIndex,
								ADRTxPowerIndex:  tc.CurrentADRTxPowerIndex,
								ADRNbTrans:       tc.CurrentADRNbTrans,
							},
							DesiredParameters: ttnpb.MACParameters{
								Channels:         tc.DesiredChannels,
								ADRDataRateIndex: tc.DesiredADRDataRateIndex,
								ADRTxPowerIndex:  tc.DesiredADRTxPowerIndex,
								ADRNbTrans:       tc.DesiredADRNbTrans,
							},
							RejectedADRDataRateIndexes: tc.RejectedADRDataRateIndexes,
							RejectedADRTxPowerIndexes:  tc.RejectedADRTxPowerIndexes,
						},
						MACSettings: func() *ttnpb.MACSettings {
							if DefaultMACSettings.UseADR != nil && DefaultMACSettings.UseADR.Value == adr || DefaultMACSettings.UseADR == nil && adr {
								return nil
							}
							return &ttnpb.MACSettings{
								UseADR: &pbtypes.BoolValue{Value: adr},
							}
						}(),
					})
				}
				phy := LoRaWANBands[tc.BandID][tc.LoRaWANPHYVersion]

				test.RunSubtestFromContext(ctx, test.SubtestConfig{
					Name:     "DeviceNeedsLinkADRReq",
					Parallel: true,
					Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
						dev := makeDevice(true)
						a.So(DeviceNeedsLinkADRReq(ctx, dev, DefaultMACSettings, phy), func() func(interface{}, ...interface{}) string {
							if len(tc.ADRCommands) > 0 || len(tc.NoADRCommands) > 0 {
								return should.BeTrue
							}
							return should.BeFalse
						}())
						a.So(dev, should.Resemble, makeDevice(true))

						dev = makeDevice(false)
						a.So(DeviceNeedsLinkADRReq(ctx, dev, DefaultMACSettings, phy), func() func(interface{}, ...interface{}) string {
							if len(tc.NoADRCommands) > 0 {
								return should.BeTrue
							}
							return should.BeFalse
						}())
						a.So(dev, should.Resemble, makeDevice(false))
					},
				})
				for adr, cmds := range map[bool][]*ttnpb.MACCommand_LinkADRReq{
					true: func() []*ttnpb.MACCommand_LinkADRReq {
						if len(tc.ADRCommands) == 0 {
							return tc.NoADRCommands
						}
						return tc.ADRCommands
					}(),
					false: tc.NoADRCommands,
				} {
					for _, n := range func() []int {
						switch len(cmds) {
						case 0:
							return []int{0}
						default:
							return []int{0, len(cmds)}
						}
					}() {
						adr := adr
						cmdsFit := n >= len(cmds)
						cmdLen := (1 + lorawan.DefaultMACCommands[ttnpb.CID_LINK_ADR].DownlinkLength) * uint16(n)
						cmds := cmds[:n]
						answerLen := (1 + lorawan.DefaultMACCommands[ttnpb.CID_LINK_ADR].UplinkLength) * func() uint16 {
							switch {
							case n == 0:
								return 0
							case tc.LoRaWANVersion.Compare(ttnpb.MAC_V1_1) >= 0:
								return 1
							default:
								return uint16(n)
							}
						}()
						test.RunSubtestFromContext(ctx, test.SubtestConfig{
							Name:     fmt.Sprintf("EnqueueLinkADRReq/adr:%v,max_down_len:%d", adr, cmdLen),
							Parallel: true,
							Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
								dev := makeDevice(adr)
								st, err := EnqueueLinkADRReq(ctx, dev, cmdLen, answerLen, DefaultMACSettings, phy)
								if errorAssertion := func() func(*testing.T, error) bool {
									switch {
									case !adr:
										return tc.NoADRErrorAssertion
									case tc.ADRErrorAssertion != nil:
										return tc.ADRErrorAssertion
									default:
										return tc.NoADRErrorAssertion
									}
								}(); errorAssertion != nil {
									if !a.So(errorAssertion(t, err), should.BeTrue) {
										t.FailNow()
									}
									a.So(st, should.Resemble, EnqueueState{
										MaxDownLen: cmdLen,
										MaxUpLen:   answerLen,
										Ok:         false,
									})
									return
								}
								if !a.So(err, should.BeNil) {
									t.Fatalf("Failed to enqueue LinkADRReq: %s", err)
								}
								expectedDevice := makeDevice(adr)
								var expectedEventBuilders []events.Builder
								for _, cmd := range cmds {
									expectedDevice.MACState.PendingRequests = append(expectedDevice.MACState.PendingRequests, cmd.MACCommand())
									expectedEventBuilders = append(expectedEventBuilders, EvtEnqueueLinkADRRequest.BindData(cmd))
								}
								a.So(st.QueuedEvents, should.ResembleEventBuilders, events.Builders(expectedEventBuilders))
								if a.So(st, should.Resemble, EnqueueState{
									QueuedEvents: st.QueuedEvents,
									Ok:           cmdsFit,
								}) {
									a.So(dev, should.Resemble, expectedDevice)
								}
							},
						})
					}
				}
			},
		})
	}
}

func TestHandleLinkADRAns(t *testing.T) {
	recentADRUplinks := []*ttnpb.UplinkMessage{
		{
			Payload: &ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_UP,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							FCtrl: ttnpb.FCtrl{
								ADR: true,
							},
							FCnt: 42,
						},
					},
				},
			},
		},
		{
			Payload: &ttnpb.Message{
				MHDR: ttnpb.MHDR{
					MType: ttnpb.MType_UNCONFIRMED_UP,
				},
				Payload: &ttnpb.Message_MACPayload{
					MACPayload: &ttnpb.MACPayload{
						FHDR: ttnpb.FHDR{
							FCtrl: ttnpb.FCtrl{
								ADR: true,
							},
							FCnt: 43,
						},
					},
				},
			},
		},
	}

	for _, tc := range []struct {
		Name             string
		Device, Expected *ttnpb.EndDevice
		Payload          *ttnpb.MACCommand_LinkADRAns
		DupCount         uint
		Events           events.Builders
		Error            error
	}{
		{
			Name: "nil payload",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Expected: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Error: ErrNoPayload,
		},
		{
			Name: "no request",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Expected: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Payload: &ttnpb.MACCommand_LinkADRAns{
				ChannelMaskAck:   true,
				DataRateIndexAck: true,
				TxPowerIndexAck:  true,
			},
			Events: events.Builders{
				EvtReceiveLinkADRAccept.With(events.WithData(&ttnpb.MACCommand_LinkADRAns{
					ChannelMaskAck:   true,
					DataRateIndexAck: true,
					TxPowerIndexAck:  true,
				})),
			},
			Error: ErrRequestNotFound,
		},
		{
			Name: "1 request/all ack",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					CurrentParameters: ttnpb.MACParameters{
						Channels: []*ttnpb.MACParameters_Channel{
							nil,
							{UplinkFrequency: 42},
							{DownlinkFrequency: 23},
							nil,
						},
					},
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_4,
							TxPowerIndex:  42,
							ChannelMask: []bool{
								false, true, false, false,
								false, false, false, false,
								false, false, false, false,
								false, false, false, false,
							},
						}).MACCommand(),
					},
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Expected: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					CurrentParameters: ttnpb.MACParameters{
						ADRDataRateIndex: ttnpb.DATA_RATE_4,
						ADRTxPowerIndex:  42,
						Channels: []*ttnpb.MACParameters_Channel{
							nil,
							{
								EnableUplink:    true,
								UplinkFrequency: 42,
							},
							{
								EnableUplink:      false,
								DownlinkFrequency: 23,
							},
							nil,
						},
					},
					PendingRequests: []*ttnpb.MACCommand{},
				},
			},
			Payload: &ttnpb.MACCommand_LinkADRAns{
				ChannelMaskAck:   true,
				DataRateIndexAck: true,
				TxPowerIndexAck:  true,
			},
			Events: events.Builders{
				EvtReceiveLinkADRAccept.With(events.WithData(&ttnpb.MACCommand_LinkADRAns{
					ChannelMaskAck:   true,
					DataRateIndexAck: true,
					TxPowerIndexAck:  true,
				})),
			},
		},
		{
			Name: "1.1/2 requests/all ack",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					CurrentParameters: ttnpb.MACParameters{
						Channels: []*ttnpb.MACParameters_Channel{
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							nil,
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
						},
					},
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_5,
							TxPowerIndex:  42,
							ChannelMask: []bool{
								true, true, true, false,
								true, true, true, true,
								true, true, true, true,
								true, true, false, false,
							},
						}).MACCommand(),
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_10,
							TxPowerIndex:  43,
							ChannelMask: []bool{
								false, true, true, false,
								true, true, true, true,
								true, true, true, true,
								true, true, false, false,
							},
						}).MACCommand(),
					},
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Expected: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_1,
					CurrentParameters: ttnpb.MACParameters{
						ADRDataRateIndex: ttnpb.DATA_RATE_10,
						ADRTxPowerIndex:  43,
						Channels: []*ttnpb.MACParameters_Channel{
							{EnableUplink: false},
							{EnableUplink: true},
							{EnableUplink: true},
							nil,
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: false},
						},
					},
					PendingRequests: []*ttnpb.MACCommand{},
				},
			},
			Payload: &ttnpb.MACCommand_LinkADRAns{
				ChannelMaskAck:   true,
				DataRateIndexAck: true,
				TxPowerIndexAck:  true,
			},
			Events: events.Builders{
				EvtReceiveLinkADRAccept.With(events.WithData(&ttnpb.MACCommand_LinkADRAns{
					ChannelMaskAck:   true,
					DataRateIndexAck: true,
					TxPowerIndexAck:  true,
				})),
			},
		},
		{
			Name:     "1.0.2/2 requests/all ack",
			DupCount: 1,
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_0_2,
					CurrentParameters: ttnpb.MACParameters{
						Channels: []*ttnpb.MACParameters_Channel{
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							nil,
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
						},
					},
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_5,
							TxPowerIndex:  42,
							ChannelMask: []bool{
								true, true, true, false,
								true, true, true, true,
								true, true, true, true,
								true, true, false, false,
							},
						}).MACCommand(),
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_10,
							TxPowerIndex:  43,
							ChannelMask: []bool{
								false, true, true, false,
								true, true, true, true,
								true, true, true, true,
								true, true, false, false,
							},
						}).MACCommand(),
					},
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Expected: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_0_2,
					CurrentParameters: ttnpb.MACParameters{
						ADRDataRateIndex: ttnpb.DATA_RATE_10,
						ADRTxPowerIndex:  43,
						Channels: []*ttnpb.MACParameters_Channel{
							{EnableUplink: false},
							{EnableUplink: true},
							{EnableUplink: true},
							nil,
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: false},
						},
					},
					PendingRequests: []*ttnpb.MACCommand{},
				},
			},
			Payload: &ttnpb.MACCommand_LinkADRAns{
				ChannelMaskAck:   true,
				DataRateIndexAck: true,
				TxPowerIndexAck:  true,
			},
			Events: events.Builders{
				EvtReceiveLinkADRAccept.With(events.WithData(&ttnpb.MACCommand_LinkADRAns{
					ChannelMaskAck:   true,
					DataRateIndexAck: true,
					TxPowerIndexAck:  true,
				})),
			},
		},
		{
			Name: "1.0/2 requests/all ack",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_0,
					CurrentParameters: ttnpb.MACParameters{
						Channels: []*ttnpb.MACParameters_Channel{
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							nil,
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
						},
					},
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_5,
							TxPowerIndex:  42,
							ChannelMask: []bool{
								true, true, true, false,
								true, true, true, true,
								true, true, true, true,
								true, true, false, false,
							},
						}).MACCommand(),
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_10,
							TxPowerIndex:  43,
							ChannelMask: []bool{
								false, true, true, false,
								true, true, true, true,
								true, true, true, true,
								true, true, false, false,
							},
						}).MACCommand(),
					},
				},
				RecentADRUplinks: recentADRUplinks,
			},
			Expected: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_0,
					CurrentParameters: ttnpb.MACParameters{
						ADRDataRateIndex: ttnpb.DATA_RATE_5,
						ADRTxPowerIndex:  42,
						Channels: []*ttnpb.MACParameters_Channel{
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							nil,
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: true},
							{EnableUplink: false},
						},
					},
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_10,
							TxPowerIndex:  43,
							ChannelMask: []bool{
								false, true, true, false,
								true, true, true, true,
								true, true, true, true,
								true, true, false, false,
							},
						}).MACCommand(),
					},
				},
			},
			Payload: &ttnpb.MACCommand_LinkADRAns{
				ChannelMaskAck:   true,
				DataRateIndexAck: true,
				TxPowerIndexAck:  true,
			},
			Events: events.Builders{
				EvtReceiveLinkADRAccept.With(events.WithData(&ttnpb.MACCommand_LinkADRAns{
					ChannelMaskAck:   true,
					DataRateIndexAck: true,
					TxPowerIndexAck:  true,
				})),
			},
		},
		{
			Name: "1.0.2/2 requests/US915 FSB2",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.USFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion:    ttnpb.MAC_V1_0_2,
					CurrentParameters: MakeDefaultUS915CurrentMACParameters(ttnpb.PHY_V1_0_2_REV_B),
					DesiredParameters: MakeDefaultUS915FSB2DesiredMACParameters(ttnpb.PHY_V1_0_2_REV_B),
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex:      ttnpb.DATA_RATE_3,
							TxPowerIndex:       1,
							ChannelMaskControl: 7,
							NbTrans:            3,
							ChannelMask: []bool{
								false, false, false, false,
								false, false, false, false,
								false, false, false, false,
								false, false, false, false,
							},
						}).MACCommand(),
						(&ttnpb.MACCommand_LinkADRReq{
							DataRateIndex: ttnpb.DATA_RATE_3,
							TxPowerIndex:  1,
							NbTrans:       3,
							ChannelMask: []bool{
								false, false, false, false,
								false, false, false, false,
								true, true, true, true,
								true, true, true, true,
							},
						}).MACCommand(),
					},
				},
			},
			Expected: &ttnpb.EndDevice{
				FrequencyPlanID:   test.USFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					LoRaWANVersion: ttnpb.MAC_V1_0_2,
					CurrentParameters: func() ttnpb.MACParameters {
						params := MakeDefaultUS915FSB2DesiredMACParameters(ttnpb.PHY_V1_0_2_REV_B)
						params.ADRDataRateIndex = ttnpb.DATA_RATE_3
						params.ADRTxPowerIndex = 1
						params.ADRNbTrans = 3
						return params
					}(),
					DesiredParameters: MakeDefaultUS915FSB2DesiredMACParameters(ttnpb.PHY_V1_0_2_REV_B),
					PendingRequests:   []*ttnpb.MACCommand{},
				},
			},
			Payload: &ttnpb.MACCommand_LinkADRAns{
				ChannelMaskAck:   true,
				DataRateIndexAck: true,
				TxPowerIndexAck:  true,
			},
			Events: events.Builders{
				EvtReceiveLinkADRAccept.With(events.WithData(&ttnpb.MACCommand_LinkADRAns{
					ChannelMaskAck:   true,
					DataRateIndexAck: true,
					TxPowerIndexAck:  true,
				})),
			},
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				dev := CopyEndDevice(tc.Device)

				evs, err := HandleLinkADRAns(ctx, dev, tc.Payload, tc.DupCount, frequencyplans.NewStore(test.FrequencyPlansFetcher))
				if tc.Error != nil && !a.So(err, should.EqualErrorOrDefinition, tc.Error) ||
					tc.Error == nil && !a.So(err, should.BeNil) {
					t.FailNow()
				}
				a.So(dev, should.Resemble, tc.Expected)
				a.So(evs, should.ResembleEventBuilders, tc.Events)
			},
		})
	}
}
