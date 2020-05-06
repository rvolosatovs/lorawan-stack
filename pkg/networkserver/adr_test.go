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
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func newADRUplink(fCnt uint32, maxSNR float32, gtwCount uint, confirmed bool, tx ttnpb.TxSettings) *ttnpb.UplinkMessage {
	if gtwCount == 0 {
		gtwCount = 1 + uint(rand.Int()%100)
	}
	mds := make([]*ttnpb.RxMetadata, 0, gtwCount)
	for i := uint(0); i < gtwCount; i++ {
		mds = append(mds, &ttnpb.RxMetadata{
			SNR: float32(-rand.Int31n(math.MaxInt32+int32(maxSNR)-1)) - rand.Float32() + maxSNR,
		})
	}
	mds[rand.Intn(len(mds))].SNR = maxSNR

	mType := ttnpb.MType_UNCONFIRMED_UP
	if confirmed {
		mType = ttnpb.MType_CONFIRMED_UP
	}

	return &ttnpb.UplinkMessage{
		Payload: &ttnpb.Message{
			MHDR: ttnpb.MHDR{
				MType: mType,
			},
			Payload: &ttnpb.Message_MACPayload{
				MACPayload: &ttnpb.MACPayload{
					FHDR: ttnpb.FHDR{
						FCtrl: ttnpb.FCtrl{
							ADR: true,
						},
						FCnt: fCnt,
					},
				},
			},
		},
		RxMetadata: mds,
		Settings:   tx,
	}
}

type adrMatrixRow struct {
	FCnt         uint32
	MaxSNR       float32
	GtwDiversity uint
	Confirmed    bool
	TxSettings   ttnpb.TxSettings
}

func adrMatrixToUplinks(m []adrMatrixRow) []*ttnpb.UplinkMessage {
	if len(m) > 20 {
		panic("ADR matrix contains more than 20 rows")
	}

	ups := make([]*ttnpb.UplinkMessage, 0, 20)
	for _, r := range m {
		ups = append(ups, newADRUplink(r.FCnt, r.MaxSNR, r.GtwDiversity, r.Confirmed, r.TxSettings))
	}
	return ups
}

func TestLossRate(t *testing.T) {
	for _, tc := range []struct {
		Name    string
		Uplinks []*ttnpb.UplinkMessage
		NbTrans uint32
		Rate    float32
	}{
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 13},
			}),
			NbTrans: 1,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 12},
				{FCnt: 13},
				{FCnt: 13},
			}),
			NbTrans: 2,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 11},
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 12},
				{FCnt: 12},
				{FCnt: 13},
				{FCnt: 13},
				{FCnt: 13},
			}),
			NbTrans: 3,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 12},
				{FCnt: 13},
			}),
			NbTrans: 2,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 12},
				{FCnt: 12},
				{FCnt: 13},
			}),
			NbTrans: 3,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 13},
			}),
			NbTrans: 1,
			Rate:    1. / 3.,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 14},
			}),
			NbTrans: 1,
			Rate:    2. / 4.,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 13},
				{FCnt: 15},
			}),
			NbTrans: 1,
			Rate:    2. / 5.,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 13},
				{FCnt: 13},
			}),
			NbTrans: 2,
			Rate:    1. / 6.,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 13},
			}),
			NbTrans: 2,
			Rate:    1. / 4.,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 12},
				{FCnt: 13},
				{FCnt: 13},
				{FCnt: 13},
			}),
			NbTrans: 3,
			Rate:    1. / 7.,
		},
		{
			Uplinks: adrMatrixToUplinks([]adrMatrixRow{
				{FCnt: 11},
				{FCnt: 12},
				{FCnt: 1},
				{FCnt: 1},
				{FCnt: 3},
				{FCnt: 3},
			}),
			NbTrans: 3,
			Rate:    3. / 7.,
		},
	} {
		t.Run(func() string {
			var ss []string
			for _, up := range tc.Uplinks {
				ss = append(ss, fmt.Sprintf("%d", up.Payload.GetMACPayload().FHDR.FCnt))
			}
			return fmt.Sprintf("NbTrans %d/%s", tc.NbTrans, strings.Join(ss, ","))
		}(), func(t *testing.T) {
			assertions.New(t).So(lossRate(tc.NbTrans, tc.Uplinks...), should.Equal, tc.Rate)
		})
	}
}

func TestAdaptDataRate(t *testing.T) {
	semtechPaperUplinks := adrMatrixToUplinks([]adrMatrixRow{
		{FCnt: 10, MaxSNR: -6, GtwDiversity: 2},
		{FCnt: 11, MaxSNR: -7, GtwDiversity: 2},
		{FCnt: 12, MaxSNR: -25, GtwDiversity: 1},
		{FCnt: 13, MaxSNR: -25, GtwDiversity: 1},
		{FCnt: 14, MaxSNR: -10, GtwDiversity: 2},
		{FCnt: 16, MaxSNR: -25, GtwDiversity: 1},
		{FCnt: 17, MaxSNR: -10, GtwDiversity: 2},
		{FCnt: 19, MaxSNR: -10, GtwDiversity: 3},
		{FCnt: 20, MaxSNR: -6, GtwDiversity: 2},
		{FCnt: 21, MaxSNR: -7, GtwDiversity: 2},
		{FCnt: 22, MaxSNR: -25, GtwDiversity: 0},
		{FCnt: 23, MaxSNR: -25, GtwDiversity: 1},
		{FCnt: 24, MaxSNR: -10, GtwDiversity: 2},
		{FCnt: 25, MaxSNR: -10, GtwDiversity: 2},
		{FCnt: 26, MaxSNR: -25, GtwDiversity: 1},
		{FCnt: 27, MaxSNR: -8, GtwDiversity: 2},
		{FCnt: 28, MaxSNR: -10, GtwDiversity: 2},
		{FCnt: 29, MaxSNR: -10, GtwDiversity: 3},
		{FCnt: 30, MaxSNR: -9, GtwDiversity: 3},
		{
			FCnt: 31, MaxSNR: -7, GtwDiversity: 2,
			TxSettings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							SpreadingFactor: 12,
							Bandwidth:       125000,
						},
					},
				},
				DataRateIndex: 0,
			},
		},
	})
	for _, tc := range []struct {
		Name       string
		Device     *ttnpb.EndDevice
		DeviceDiff func(*ttnpb.EndDevice)
		Error      error
	}{
		{
			Name: "adapted example from Semtech paper/no rejections",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					CurrentParameters: ttnpb.MACParameters{
						ADRNbTrans:      1,
						ADRTxPowerIndex: 1,
						Channels:        MakeDefaultEU868CurrentChannels(),
					},
					DesiredParameters: ttnpb.MACParameters{
						ADRDataRateIndex: ttnpb.DATA_RATE_4,
						ADRNbTrans:       3,
						ADRTxPowerIndex:  2,
						Channels:         MakeDefaultEU868CurrentChannels(),
					},
				},
				MACSettings: &ttnpb.MACSettings{
					ADRMargin: &pbtypes.FloatValue{
						Value: 2,
					},
				},
				RecentADRUplinks: semtechPaperUplinks,
			},
			DeviceDiff: func(dev *ttnpb.EndDevice) {
				dev.MACState.DesiredParameters.ADRDataRateIndex = ttnpb.DATA_RATE_4
				dev.MACState.DesiredParameters.ADRTxPowerIndex = 1
				dev.MACState.DesiredParameters.ADRNbTrans = 1
			},
		},
		{
			Name: "adapted example from Semtech paper/rejected DR:(1,4)",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					CurrentParameters: ttnpb.MACParameters{
						ADRNbTrans:      1,
						ADRTxPowerIndex: 1,
						Channels:        MakeDefaultEU868CurrentChannels(),
					},
					DesiredParameters: ttnpb.MACParameters{
						Channels: MakeDefaultEU868CurrentChannels(),
					},
					RejectedADRDataRateIndexes: []ttnpb.DataRateIndex{
						ttnpb.DATA_RATE_1, ttnpb.DATA_RATE_4,
					},
				},
				MACSettings: &ttnpb.MACSettings{
					ADRMargin: &pbtypes.FloatValue{
						Value: 2,
					},
				},
				RecentADRUplinks: semtechPaperUplinks,
			},
			DeviceDiff: func(dev *ttnpb.EndDevice) {
				dev.MACState.DesiredParameters.ADRDataRateIndex = ttnpb.DATA_RATE_3
				dev.MACState.DesiredParameters.ADRTxPowerIndex = 2
				dev.MACState.DesiredParameters.ADRNbTrans = 1
			},
		},
		{
			Name: "adapted example from Semtech paper/rejected TXPower:(1)",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					CurrentParameters: ttnpb.MACParameters{
						ADRNbTrans:      1,
						ADRTxPowerIndex: 1,
						Channels:        MakeDefaultEU868CurrentChannels(),
					},
					DesiredParameters: ttnpb.MACParameters{
						Channels: MakeDefaultEU868CurrentChannels(),
					},
					RejectedADRTxPowerIndexes: []uint32{
						1,
					},
				},
				MACSettings: &ttnpb.MACSettings{
					ADRMargin: &pbtypes.FloatValue{
						Value: 2,
					},
				},
				RecentADRUplinks: semtechPaperUplinks,
			},
			DeviceDiff: func(dev *ttnpb.EndDevice) {
				dev.MACState.DesiredParameters.ADRDataRateIndex = ttnpb.DATA_RATE_4
				dev.MACState.DesiredParameters.ADRTxPowerIndex = 0
				dev.MACState.DesiredParameters.ADRNbTrans = 1
			},
		},
		{
			Name: "adapted example from Semtech paper/rejected TXPower:(0,1)",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					CurrentParameters: ttnpb.MACParameters{
						ADRNbTrans:      1,
						ADRTxPowerIndex: 1,
						Channels:        MakeDefaultEU868CurrentChannels(),
					},
					DesiredParameters: ttnpb.MACParameters{
						Channels: MakeDefaultEU868CurrentChannels(),
					},
					RejectedADRTxPowerIndexes: []uint32{
						0, 1,
					},
				},
				MACSettings: &ttnpb.MACSettings{
					ADRMargin: &pbtypes.FloatValue{
						Value: 2,
					},
				},
				RecentADRUplinks: semtechPaperUplinks,
			},
			DeviceDiff: func(dev *ttnpb.EndDevice) {
				dev.MACState.DesiredParameters.ADRDataRateIndex = ttnpb.DATA_RATE_3
				dev.MACState.DesiredParameters.ADRTxPowerIndex = 2
				dev.MACState.DesiredParameters.ADRNbTrans = 1
			},
		},
		{
			Name: "adapted example from Semtech paper/rejected DR:(1,4), rejected TXPower:(0,2,3)",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					CurrentParameters: ttnpb.MACParameters{
						ADRNbTrans:      1,
						ADRTxPowerIndex: 1,
						Channels:        MakeDefaultEU868CurrentChannels(),
					},
					DesiredParameters: ttnpb.MACParameters{
						Channels: MakeDefaultEU868CurrentChannels(),
					},
					RejectedADRTxPowerIndexes: []uint32{
						0, 2, 3,
					},
					RejectedADRDataRateIndexes: []ttnpb.DataRateIndex{
						ttnpb.DATA_RATE_1, ttnpb.DATA_RATE_4,
					},
				},
				MACSettings: &ttnpb.MACSettings{
					ADRMargin: &pbtypes.FloatValue{
						Value: 2,
					},
				},
				RecentADRUplinks: semtechPaperUplinks,
			},
			DeviceDiff: func(dev *ttnpb.EndDevice) {
				dev.MACState.DesiredParameters.ADRDataRateIndex = ttnpb.DATA_RATE_3
				dev.MACState.DesiredParameters.ADRTxPowerIndex = 1
				dev.MACState.DesiredParameters.ADRNbTrans = 1
			},
		},
		{
			Name: "adapted example from Semtech paper/rejected DR:(3), rejected TXPower:(0,1)",
			Device: &ttnpb.EndDevice{
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				MACState: &ttnpb.MACState{
					CurrentParameters: ttnpb.MACParameters{
						ADRNbTrans:      1,
						ADRTxPowerIndex: 1,
						Channels:        MakeDefaultEU868CurrentChannels(),
					},
					DesiredParameters: ttnpb.MACParameters{
						Channels: MakeDefaultEU868CurrentChannels(),
					},
					RejectedADRTxPowerIndexes: []uint32{
						0, 1,
					},
					RejectedADRDataRateIndexes: []ttnpb.DataRateIndex{
						ttnpb.DATA_RATE_3,
					},
				},
				MACSettings: &ttnpb.MACSettings{
					ADRMargin: &pbtypes.FloatValue{
						Value: 2,
					},
				},
				RecentADRUplinks: semtechPaperUplinks,
			},
			DeviceDiff: func(dev *ttnpb.EndDevice) {
				dev.MACState.DesiredParameters.ADRDataRateIndex = ttnpb.DATA_RATE_2
				dev.MACState.DesiredParameters.ADRTxPowerIndex = 3
				dev.MACState.DesiredParameters.ADRNbTrans = 1
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			dev := CopyEndDevice(tc.Device)
			fp := FrequencyPlan(dev.FrequencyPlanID)
			err := adaptDataRate(
				log.NewContext(test.ContextWithT(test.Context(), t), test.GetLogger(t)),
				dev,
				Band(fp.BandID, dev.LoRaWANPHYVersion), ttnpb.MACSettings{},
			)
			if !a.So(err, should.Equal, tc.Error) {
				t.Fatalf("ADR failed with: %s", err)
			}
			expected := CopyEndDevice(tc.Device)
			if tc.DeviceDiff != nil {
				tc.DeviceDiff(expected)
			}
			a.So(dev, should.Resemble, expected)
		})
	}
}
