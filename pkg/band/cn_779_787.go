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

package band

import (
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

//revive:disable:var-naming

var cn_779_787 Band

// CN_779_787 is the ID of the Chinese 779-787Mhz band
const CN_779_787 = "CN_779_787"

//revive:enable:var-naming

func init() {
	defaultChannels := []Channel{
		{Frequency: 779500000, MinDataRate: 0, MaxDataRate: 5},
		{Frequency: 779500000, MinDataRate: 0, MaxDataRate: 5},
		{Frequency: 779900000, MinDataRate: 0, MaxDataRate: 5},
		{Frequency: 780500000, MinDataRate: 0, MaxDataRate: 5},
		{Frequency: 780700000, MinDataRate: 0, MaxDataRate: 5},
		{Frequency: 780900000, MinDataRate: 0, MaxDataRate: 5},
	}
	const cnBeaconFrequency = 785000000

	downlinkDRTable := [8][6]ttnpb.DataRateIndex{
		{0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0},
		{2, 1, 0, 0, 0, 0},
		{3, 2, 1, 0, 0, 0},
		{4, 3, 2, 1, 0, 0},
		{5, 4, 3, 2, 1, 0},
		{6, 5, 4, 3, 2, 1},
		{7, 6, 5, 4, 3, 2},
	}

	cn_779_787 = Band{
		ID: CN_779_787,

		MaxUplinkChannels: 16,
		UplinkChannels:    defaultChannels,

		MaxDownlinkChannels: 16,
		DownlinkChannels:    defaultChannels,

		// See IEEE 11-11/0972r0
		SubBands: []SubBandParameters{
			{
				MinFrequency: 779000000,
				MaxFrequency: 787000000,
				DutyCycle:    0.01,
				MaxEIRP:      10.0 + eirpDelta,
			},
		},

		DataRates: map[ttnpb.DataRateIndex]DataRate{
			0: makeLoRaDataRate(12, 125000, makeConstMaxMACPayloadSizeFunc(59)),
			1: makeLoRaDataRate(11, 125000, makeConstMaxMACPayloadSizeFunc(59)),
			2: makeLoRaDataRate(10, 125000, makeConstMaxMACPayloadSizeFunc(59)),
			3: makeLoRaDataRate(9, 125000, makeConstMaxMACPayloadSizeFunc(123)),
			4: makeLoRaDataRate(8, 125000, makeConstMaxMACPayloadSizeFunc(230)),
			5: makeLoRaDataRate(7, 125000, makeConstMaxMACPayloadSizeFunc(230)),
			6: makeLoRaDataRate(7, 250000, makeConstMaxMACPayloadSizeFunc(250)),
			7: makeFSKDataRate(50000, makeConstMaxMACPayloadSizeFunc(230)),
		},
		MaxADRDataRateIndex: 5,

		ReceiveDelay1:    defaultReceiveDelay1,
		ReceiveDelay2:    defaultReceiveDelay2,
		JoinAcceptDelay1: defaultJoinAcceptDelay1,
		JoinAcceptDelay2: defaultJoinAcceptDelay2,
		MaxFCntGap:       defaultMaxFCntGap,
		ADRAckLimit:      defaultADRAckLimit,
		ADRAckDelay:      defaultADRAckDelay,
		MinAckTimeout:    defaultAckTimeout - defaultAckTimeoutMargin,
		MaxAckTimeout:    defaultAckTimeout + defaultAckTimeoutMargin,

		DefaultMaxEIRP: 12.15,
		TxOffset: [16]float32{
			0, -2, -4, -6, -8, -10,
			0, 0, 0, 0, 0, 0, 0, 0, 0, // RFU
			0, // Used by LinkADRReq starting from LoRaWAN Regional Parameters 1.1, RFU before
		},
		MaxTxPowerIndex: 5,

		Rx1Channel: channelIndexIdentity,
		Rx1DataRate: func(idx ttnpb.DataRateIndex, offset uint32, _ bool) (ttnpb.DataRateIndex, error) {
			if idx > 7 {
				return 0, errDataRateIndexTooHigh.WithAttributes("max", 7)
			}
			if offset > 5 {
				return 0, errDataRateOffsetTooHigh.WithAttributes("max", 5)
			}
			return downlinkDRTable[idx][offset], nil
		},

		GenerateChMasks: generateChMask16,
		ParseChMask:     parseChMask16,

		LoRaCodingRate: "4/5",

		FreqMultiplier:   100,
		ImplementsCFList: true,
		CFListType:       ttnpb.CFListType_FREQUENCIES,

		DefaultRx2Parameters: Rx2Parameters{0, 786000000},

		Beacon: Beacon{
			DataRateIndex:    3,
			CodingRate:       "4/5",
			ComputeFrequency: func(_ float64) uint64 { return cnBeaconFrequency },
		},
		PingSlotFrequency: uint64Ptr(cnBeaconFrequency),

		regionalParameters1_0:       bandIdentity,
		regionalParameters1_0_1:     bandIdentity,
		regionalParameters1_0_2RevA: bandIdentity,
		regionalParameters1_0_2RevB: bandIdentity,
		regionalParameters1_0_3RevA: bandIdentity,
		regionalParameters1_1RevA:   bandIdentity,
	}
	All[CN_779_787] = cn_779_787
}
