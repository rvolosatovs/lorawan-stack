// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package band

import "github.com/TheThingsNetwork/ttn/pkg/types"

var in_865_867 Band

const (
	// IN_865_867 is the ID of the Indian frequency plan
	IN_865_867 ID = "IN_865_867"
)

func init() {
	defaultChannels := []Channel{
		{Frequency: 865062500, DataRateIndexes: []int{0, 1, 2, 3, 4, 5}},
		{Frequency: 865402500, DataRateIndexes: []int{0, 1, 2, 3, 4, 5}},
		{Frequency: 865985000, DataRateIndexes: []int{0, 1, 2, 3, 4, 5}},
	}
	inBeaconChannel := uint32(866500000)
	in_865_867 = Band{
		ID: IN_865_867,

		UplinkChannels:   defaultChannels,
		DownlinkChannels: defaultChannels,

		BandDutyCycles: []DutyCycle{
			{
				MinFrequency: 865000000,
				MaxFrequency: 867000000,
				DutyCycle:    1,
			},
		},

		DataRates: [16]DataRate{
			{Rate: types.DataRate{LoRa: "SF12BW125"}, DefaultMaxSize: maxPayloadSize{59, 51}, NoRepeaterMaxSize: maxPayloadSize{59, 51}},
			{Rate: types.DataRate{LoRa: "SF11BW125"}, DefaultMaxSize: maxPayloadSize{59, 51}, NoRepeaterMaxSize: maxPayloadSize{59, 51}},
			{Rate: types.DataRate{LoRa: "SF10BW125"}, DefaultMaxSize: maxPayloadSize{59, 51}, NoRepeaterMaxSize: maxPayloadSize{59, 51}},
			{Rate: types.DataRate{LoRa: "SF9BW125"}, DefaultMaxSize: maxPayloadSize{123, 115}, NoRepeaterMaxSize: maxPayloadSize{123, 115}},
			{Rate: types.DataRate{LoRa: "SF8BW125"}, DefaultMaxSize: maxPayloadSize{230, 222}, NoRepeaterMaxSize: maxPayloadSize{250, 242}},
			{Rate: types.DataRate{LoRa: "SF7BW125"}, DefaultMaxSize: maxPayloadSize{230, 222}, NoRepeaterMaxSize: maxPayloadSize{250, 242}},
			{}, // RFU
			{Rate: types.DataRate{FSK: 50000}, DefaultMaxSize: maxPayloadSize{230, 222}, NoRepeaterMaxSize: maxPayloadSize{250, 242}},
			{}, {}, {}, {}, {}, {}, {}, // RFU
			{}, // Used by LinkADRReq starting from LoRaWAN 1.1, RFU before
		},

		ReceiveDelay1:    defaultReceiveDelay1,
		ReceiveDelay2:    defaultReceiveDelay2,
		JoinAcceptDelay1: defaultJoinAcceptDelay2,
		JoinAcceptDelay2: defaultJoinAcceptDelay2,
		MaxFCntGap:       defaultMaxFCntGap,
		AdrAckLimit:      defaultAdrAckLimit,
		AdrAckDelay:      defaultAdrAckDelay,
		MinAckTimeout:    defaultAckTimeout - defaultAckTimeoutMargin,
		MaxAckTimeout:    defaultAckTimeout + defaultAckTimeoutMargin,

		DefaultMaxEIRP: 30,
		TxOffset: func() [16]float32 {
			offset := [16]float32{}
			for i := 0; i < 11; i++ {
				offset[i] = float32(0 - 2*i)
			}
			return offset
		}(),

		Rx1Parameters: func(frequency uint64, dataRateIndex, rx1DROffset int, _ bool) (int, uint64) {
			effectiveRx1DROffset := rx1DROffset
			if effectiveRx1DROffset > 5 {
				effectiveRx1DROffset = 5 - rx1DROffset
			}

			outDataRateIndex := dataRateIndex - effectiveRx1DROffset
			if outDataRateIndex < 5 {
				outDataRateIndex = 5
			}
			return outDataRateIndex, frequency
		},

		DefaultRx2Parameters: Rx2Parameters{2, 866550000},

		Beacon: Beacon{
			DataRateIndex:    4,
			CodingRate:       "4/5",
			BroadcastChannel: func(_ float64) uint32 { return inBeaconChannel },
			PingSlotChannels: []uint32{inBeaconChannel},
		},

		// No LoRaWAN 1.0
		// No LoRaWAN 1.0.1
		regionalParameters1_0_2: self,
		regionalParameters1_1A:  self,
	}
	All = append(All, in_865_867)
}
