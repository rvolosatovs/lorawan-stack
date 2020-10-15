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

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/mac"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestNewChannelReq(t *testing.T) {
	for _, tc := range []struct {
		CurrentChannels, DesiredChannels []*ttnpb.MACParameters_Channel
		RejectedFrequencies              []uint64
		RejectedDataRateRanges           map[uint64]*ttnpb.MACState_DataRateRanges
		Commands                         []*ttnpb.MACCommand_NewChannelReq
	}{
		{},
		{
			CurrentChannels: []*ttnpb.MACParameters_Channel{
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				nil,
				{
					UplinkFrequency:  128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
			},
			DesiredChannels: []*ttnpb.MACParameters_Channel{
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				nil,
				{
					UplinkFrequency:  128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
			},
		},
		{
			CurrentChannels: []*ttnpb.MACParameters_Channel{
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  124,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_3,
				},
				{
					UplinkFrequency:  128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
			},
			DesiredChannels: []*ttnpb.MACParameters_Channel{
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  124,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_3,
				},
				{
					UplinkFrequency:  128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
			},
		},
		{
			CurrentChannels: []*ttnpb.MACParameters_Channel{
				{
					UplinkFrequency:  124,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_3,
				},
				nil,
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  129,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					UplinkFrequency:  150,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
			},
			DesiredChannels: []*ttnpb.MACParameters_Channel{
				nil,
				{
					UplinkFrequency:  128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  130,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
			},
			Commands: []*ttnpb.MACCommand_NewChannelReq{
				{},
				{
					ChannelIndex:     1,
					Frequency:        128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					ChannelIndex:     3,
					Frequency:        130,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					ChannelIndex: 4,
				},
			},
		},
		{
			CurrentChannels: []*ttnpb.MACParameters_Channel{
				{
					UplinkFrequency:  124,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_3,
				},
				nil,
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  129,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					UplinkFrequency:  150,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
			},
			DesiredChannels: []*ttnpb.MACParameters_Channel{
				nil,
				{
					UplinkFrequency:  128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  130,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
			},
			Commands: []*ttnpb.MACCommand_NewChannelReq{
				{},
				{
					ChannelIndex:     1,
					Frequency:        128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					ChannelIndex: 4,
				},
			},
			RejectedFrequencies: []uint64{130},
		},
		{
			CurrentChannels: []*ttnpb.MACParameters_Channel{
				{
					UplinkFrequency:  124,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_3,
				},
				nil,
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  129,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					UplinkFrequency:  150,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
			},
			DesiredChannels: []*ttnpb.MACParameters_Channel{
				nil,
				{
					UplinkFrequency:  128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
				{
					UplinkFrequency:  123,
					MinDataRateIndex: ttnpb.DATA_RATE_1,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
				{
					UplinkFrequency:  130,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_5,
				},
			},
			Commands: []*ttnpb.MACCommand_NewChannelReq{
				{
					ChannelIndex:     1,
					Frequency:        128,
					MinDataRateIndex: ttnpb.DATA_RATE_2,
					MaxDataRateIndex: ttnpb.DATA_RATE_4,
				},
			},
			RejectedFrequencies: []uint64{130},
			RejectedDataRateRanges: map[uint64]*ttnpb.MACState_DataRateRanges{
				0: {
					Ranges: []*ttnpb.MACState_DataRateRange{
						{},
					},
				},
			},
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name: func() string {
				formatChannels := func(chs ...*ttnpb.MACParameters_Channel) string {
					return fmt.Sprintf("[%s]", test.JoinStringsMap(func(_, v interface{}) string {
						ch := v.(*ttnpb.MACParameters_Channel)
						if ch == nil {
							return "nil"
						}
						return fmt.Sprintf("(%d,%d-%d)", ch.UplinkFrequency, ch.MinDataRateIndex, ch.MaxDataRateIndex)
					}, ",", chs))
				}
				return fmt.Sprintf("channels:%s->%s,rejected_freqs:[%s],rejected_drs:[%s]",
					formatChannels(tc.CurrentChannels...),
					formatChannels(tc.DesiredChannels...),
					test.JoinStringsf("%d", ",", false, tc.RejectedFrequencies),
					test.JoinStringsMap(func(freq, rs interface{}) string {
						return fmt.Sprintf("%d:[%s]",
							freq,
							test.JoinStringsMap(func(_, v interface{}) string {
								r := v.(*ttnpb.MACState_DataRateRange)
								return fmt.Sprintf("%d-%d", r.MinDataRateIndex, r.MaxDataRateIndex)
							}, "", rs.(*ttnpb.MACState_DataRateRanges).Ranges),
						)
					}, ",", tc.RejectedDataRateRanges),
				)
			}(),
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				makeDevice := func() *ttnpb.EndDevice {
					return CopyEndDevice(&ttnpb.EndDevice{
						MACState: &ttnpb.MACState{
							CurrentParameters: ttnpb.MACParameters{
								Channels: tc.CurrentChannels,
							},
							DesiredParameters: ttnpb.MACParameters{
								Channels: tc.DesiredChannels,
							},
							RejectedFrequencies:    tc.RejectedFrequencies,
							RejectedDataRateRanges: tc.RejectedDataRateRanges,
						},
					})
				}

				test.RunSubtestFromContext(ctx, test.SubtestConfig{
					Name:     "DeviceNeedsNewChannelReqAtIndex",
					Parallel: true,
					Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
						dev := makeDevice()
						max := len(dev.MACState.CurrentParameters.Channels)
						if n := len(dev.MACState.DesiredParameters.Channels); n > max {
							max = n
						}
						needs := make(map[int]struct{}, max)
						for _, cmd := range tc.Commands {
							needs[int(cmd.ChannelIndex)] = struct{}{}
						}
						for i := 0; i <= max+1; i++ {
							a.So(DeviceNeedsNewChannelReqAtIndex(dev, i), func() func(interface{}, ...interface{}) string {
								if _, ok := needs[i]; ok {
									return should.BeTrue
								}
								return should.BeFalse
							}())
						}
					},
				})

				test.RunSubtestFromContext(ctx, test.SubtestConfig{
					Name:     "DeviceNeedsNewChannelReq",
					Parallel: true,
					Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
						dev := makeDevice()
						a.So(DeviceNeedsNewChannelReq(dev), func() func(interface{}, ...interface{}) string {
							if len(tc.Commands) > 0 {
								return should.BeTrue
							}
							return should.BeFalse
						}())
						a.So(dev, should.Resemble, makeDevice())
					},
				})

				for _, n := range func() []int {
					switch len(tc.Commands) {
					case 0:
						return []int{0}
					case 1:
						return []int{0, 1}
					default:
						return []int{0, len(tc.Commands) / 2, len(tc.Commands)}
					}
				}() {
					cmdsFit := n >= len(tc.Commands)
					cmdLen := (1 + lorawan.DefaultMACCommands[ttnpb.CID_NEW_CHANNEL].DownlinkLength) * uint16(n)
					cmds := tc.Commands[:n]
					answerLen := (1 + lorawan.DefaultMACCommands[ttnpb.CID_NEW_CHANNEL].UplinkLength) * uint16(n)
					test.RunSubtestFromContext(ctx, test.SubtestConfig{
						Name:     fmt.Sprintf("EnqueueNewChannelReq/max_down_len:%d", cmdLen),
						Parallel: true,
						Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
							dev := makeDevice()
							st := EnqueueNewChannelReq(ctx, dev, cmdLen, answerLen)
							expectedDevice := makeDevice()
							var expectedEventBuilders []events.Builder
							for _, cmd := range cmds {
								expectedDevice.MACState.PendingRequests = append(expectedDevice.MACState.PendingRequests, cmd.MACCommand())
								expectedEventBuilders = append(expectedEventBuilders, EvtEnqueueNewChannelRequest.BindData(cmd))
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
			},
		})
	}
}

func TestHandleNewChannelAns(t *testing.T) {
	for _, tc := range []struct {
		Name             string
		Device, Expected *ttnpb.EndDevice
		Payload          *ttnpb.MACCommand_NewChannelAns
		Events           events.Builders
		Error            error
	}{
		{
			Name: "nil payload",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{},
			},
			Payload: nil,
			Error:   ErrNoPayload,
		},
		{
			Name: "no request",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{},
			},
			Payload: &ttnpb.MACCommand_NewChannelAns{
				FrequencyAck: true,
				DataRateAck:  true,
			},
			Events: events.Builders{
				EvtReceiveNewChannelAccept.With(events.WithData(&ttnpb.MACCommand_NewChannelAns{
					FrequencyAck: true,
					DataRateAck:  true,
				})),
			},
			Error: ErrRequestNotFound,
		},
		{
			Name: "frequency nack/data rate ack/no rejections",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_NewChannelReq{
							ChannelIndex:     4,
							Frequency:        42,
							MinDataRateIndex: ttnpb.DATA_RATE_2,
							MaxDataRateIndex: ttnpb.DATA_RATE_3,
						}).MACCommand(),
					},
				},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					PendingRequests:     []*ttnpb.MACCommand{},
					RejectedFrequencies: []uint64{42},
				},
			},
			Payload: &ttnpb.MACCommand_NewChannelAns{
				DataRateAck: true,
			},
			Events: events.Builders{
				EvtReceiveNewChannelReject.With(events.WithData(&ttnpb.MACCommand_NewChannelAns{
					DataRateAck: true,
				})),
			},
		},
		{
			Name: "frequency nack/data rate nack/rejected frequencies:(1,2,100)",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_NewChannelReq{
							ChannelIndex:     4,
							Frequency:        42,
							MinDataRateIndex: ttnpb.DATA_RATE_2,
							MaxDataRateIndex: ttnpb.DATA_RATE_3,
						}).MACCommand(),
					},
					RejectedFrequencies: []uint64{1, 2, 100},
				},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					PendingRequests:     []*ttnpb.MACCommand{},
					RejectedFrequencies: []uint64{1, 2, 42, 100},
					RejectedDataRateRanges: map[uint64]*ttnpb.MACState_DataRateRanges{
						42: {
							Ranges: []*ttnpb.MACState_DataRateRange{
								{
									MinDataRateIndex: ttnpb.DATA_RATE_2,
									MaxDataRateIndex: ttnpb.DATA_RATE_3,
								},
							},
						},
					},
				},
			},
			Payload: &ttnpb.MACCommand_NewChannelAns{},
			Events: events.Builders{
				EvtReceiveNewChannelReject.With(events.WithData(&ttnpb.MACCommand_NewChannelAns{})),
			},
		},
		{
			Name: "both ack",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					PendingRequests: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_NewChannelReq{
							ChannelIndex:     4,
							Frequency:        42,
							MinDataRateIndex: ttnpb.DATA_RATE_2,
							MaxDataRateIndex: ttnpb.DATA_RATE_3,
						}).MACCommand(),
					},
				},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					CurrentParameters: ttnpb.MACParameters{
						Channels: []*ttnpb.MACParameters_Channel{
							nil,
							nil,
							nil,
							nil,
							{
								DownlinkFrequency: 42,
								UplinkFrequency:   42,
								MinDataRateIndex:  2,
								MaxDataRateIndex:  3,
								EnableUplink:      true,
							},
						},
					},
					PendingRequests: []*ttnpb.MACCommand{},
				},
			},
			Payload: &ttnpb.MACCommand_NewChannelAns{
				FrequencyAck: true,
				DataRateAck:  true,
			},
			Events: events.Builders{
				EvtReceiveNewChannelAccept.With(events.WithData(&ttnpb.MACCommand_NewChannelAns{
					FrequencyAck: true,
					DataRateAck:  true,
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

				evs, err := HandleNewChannelAns(ctx, dev, tc.Payload)
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
