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

package scheduling_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/scheduling"
	"go.thethings.network/lorawan-stack/pkg/toa"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func TestScheduleAt(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fp := &frequencyplans.FrequencyPlan{
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}
	scheduler, err := scheduling.NewScheduler(ctx, fp, true)
	a.So(err, should.BeNil)
	scheduler.SyncWithGateway(0, time.Now(), time.Unix(0, 0))

	for i, tc := range []struct {
		PayloadSize    int
		Settings       ttnpb.TxSettings
		Priority       ttnpb.TxSchedulePriority
		ExpectedToa    time.Duration
		ErrorAssertion *errors.Definition
	}{
		{
			PayloadSize: 10,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 7,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Timestamp:  100,
			},
			ExpectedToa: 41216 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			// Too late for transmission.
			ErrorAssertion: &scheduling.ErrTooLate,
		},
		{
			PayloadSize: 51,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 12,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Time:       timePtr(time.Unix(0, int64(100*time.Millisecond))),
			},
			ExpectedToa: 2465792 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			// Too late for transmission.
			ErrorAssertion: &scheduling.ErrTooLate,
		},
		{
			PayloadSize: 51,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 12,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Timestamp:  20000000,
			},
			ExpectedToa: 2465792 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			// Exceeding dwell time of 2 seconds.
			ErrorAssertion: &scheduling.ErrDwellTime,
		},
		{
			PayloadSize: 10,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 7,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Timestamp:  20000000,
			},
			ExpectedToa: 41216 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		},
		{
			PayloadSize: 10,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 7,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Timestamp:  20000000,
			},
			ExpectedToa: 41216 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			// Overlapping with previous transmission.
			ErrorAssertion: &scheduling.ErrConflict,
		},
		{
			PayloadSize: 10,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 7,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Timestamp:  20000000,
			},
			ExpectedToa: 41216 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			// Right after previous transmission; not respecting time-off-air and queue delay.
			ErrorAssertion: &scheduling.ErrConflict,
		},
		{
			PayloadSize: 10,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 7,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Timestamp:  20000000 + 41216 + 1000000, // time-on-air + time-off-air.
			},
			ExpectedToa: 41216 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			// Right after previous transmission; not respecting queue delay.
			ErrorAssertion: &scheduling.ErrConflict,
		},
		{
			PayloadSize: 10,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 7,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  869525000,
				Timestamp:  20000000 + 41216 + 1000000 + 30000, // time-on-air + time-off-air + queue delay
			},
			ExpectedToa: 41216 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		},
		{
			PayloadSize: 20,
			Settings: ttnpb.TxSettings{
				DataRate: ttnpb.DataRate{
					Modulation: &ttnpb.DataRate_LoRa{
						LoRa: &ttnpb.LoRaDataRate{
							Bandwidth:       125000,
							SpreadingFactor: 12,
						},
					},
				},
				CodingRate: "4/5",
				Frequency:  868100000,
				Timestamp:  30000000, // In next duty-cycle window; discard previous.
			},
			ExpectedToa: 1318912 * time.Microsecond,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			// Exceeds duty-cycle limitation of 1% in 868.0 - 868.6.
			ErrorAssertion: &scheduling.ErrDutyCycle,
		},
	} {
		tcok := t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assertions.New(t)
			d, err := toa.Compute(tc.PayloadSize, tc.Settings)
			a.So(err, should.BeNil)
			a.So(d, should.Equal, tc.ExpectedToa)
			_, err = scheduler.ScheduleAt(ctx, tc.PayloadSize, tc.Settings, tc.Priority)
			if err != nil && tc.ErrorAssertion == nil {
				a.So(err, should.BeNil)
			} else if err != nil {
				if !a.So(err, should.HaveSameErrorDefinitionAs, *tc.ErrorAssertion) {
					t.Fatalf("Unexpected error: %v", err)
				}
			} else {
				a.So(err, should.BeNil)
			}
		})
		if !tcok {
			t.FailNow()
		}
	}
}

func TestScheduleAnytime(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fp := &frequencyplans.FrequencyPlan{
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}
	scheduler, err := scheduling.NewScheduler(ctx, fp, true)
	a.So(err, should.BeNil)
	scheduler.SyncWithGateway(0, time.Now(), time.Unix(0, 0))

	settingsAt := func(frequency uint64, sf, t uint32) ttnpb.TxSettings {
		return ttnpb.TxSettings{
			DataRate: ttnpb.DataRate{
				Modulation: &ttnpb.DataRate_LoRa{
					LoRa: &ttnpb.LoRaDataRate{
						Bandwidth:       125000,
						SpreadingFactor: sf,
					},
				},
			},
			CodingRate: "4/5",
			Frequency:  frequency,
			Timestamp:  t,
		}
	}

	// Scheduling two items, occupying considering time-on-air, time-off-air and queue delay.
	// Time-on-air is 41216 us, time-off-air is 1000000 us, queue delay is 30000 us.
	// 1: [1000000, 2071216]
	// 2: [4000000, 5071216]
	_, err = scheduler.ScheduleAt(ctx, 10, settingsAt(869525000, 7, 1000000), ttnpb.TxSchedulePriority_NORMAL)
	a.So(err, should.BeNil)
	_, err = scheduler.ScheduleAt(ctx, 10, settingsAt(869525000, 7, 4000000), ttnpb.TxSchedulePriority_NORMAL)
	a.So(err, should.BeNil)

	// Try schedule a transmission from 1000000 us.
	// Time-on-air is 41216 us, time-off-air is 1000000 us, queue delay is 30000 us.
	// It fits between 1 and 2, so it should be right after 1.
	// 1: [1000000, 2071216]
	// 3: [2071216, 3142432]
	// 2: [4000000, 5071216]
	em, err := scheduler.ScheduleAnytime(ctx, 10, settingsAt(869525000, 7, 1000000), ttnpb.TxSchedulePriority_NORMAL)
	a.So(err, should.BeNil)
	a.So(em.Starts(), should.Equal, 2071216*time.Microsecond)

	// Try schedule another transmission from 1000000 us.
	// Time-on-air is 41216 us, time-off-air is 1000000 us, queue delay is 30000 us.
	// It does not fit between 1, 3, and 2, so it should be right after 2.
	// 1: [1000000, 2071216]
	// 3: [2071216, 3142432]
	// 2: [4000000, 5071216]
	// 4: [5071216, 6142432]
	em, err = scheduler.ScheduleAnytime(ctx, 10, settingsAt(869525000, 7, 1000000), ttnpb.TxSchedulePriority_NORMAL)
	a.So(err, should.BeNil)
	a.So(em.Starts(), should.Equal, 5071216*time.Microsecond)

	// Try schedule another transmission from 1000000 us.
	// Time-on-air is 991232 us, time-off-air is 1000000 us, queue delay is 30000 us.
	// It's 9.91% in a 10% duty-cycle sub-band, almost hitting the limit, so it should be pushed to right after transmission 4.
	// Transmission starts then at 5071216 (start of 4) + 41216 (time-on-air of 4) + 10000000 (duty-cycle window) - 991232 (this time-on-air).
	// 1: [1000000, 2071216]
	// 3: [2071216, 3142432]
	// 2: [4000000, 5071216]
	// 4: [5071216, 6142432]
	// 5: [14121200, 15112432]
	em, err = scheduler.ScheduleAnytime(ctx, 10, settingsAt(869525000, 12, 1000000), ttnpb.TxSchedulePriority_HIGHEST)
	a.So(err, should.BeNil)
	a.So(em.Starts(), should.Equal, 14121200*time.Microsecond)

	// Try schedule another transmission from 1000000 us.
	// Time-on-air is 991232 us, time-off-air is 1000000 us, queue delay is 30000 us.
	// It's 9.91% in a 1% duty-cycle sub-band, so it hits the duty-cycle limitation.
	_, err = scheduler.ScheduleAnytime(ctx, 10, settingsAt(868100000, 12, 1000000), ttnpb.TxSchedulePriority_HIGHEST)
	a.So(err, should.HaveSameErrorDefinitionAs, scheduling.ErrDutyCycle)
}

func TestScheduleAnytimeShort(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fp := &frequencyplans.FrequencyPlan{
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}

	settingsAt := func(frequency uint64, sf uint32, time *time.Time, timestamp uint32) ttnpb.TxSettings {
		return ttnpb.TxSettings{
			DataRate: ttnpb.DataRate{
				Modulation: &ttnpb.DataRate_LoRa{
					LoRa: &ttnpb.LoRaDataRate{
						Bandwidth:       125000,
						SpreadingFactor: sf,
					},
				},
			},
			CodingRate: "4/5",
			Frequency:  frequency,
			Time:       time,
			Timestamp:  timestamp,
		}
	}

	// Gateway time; too late (100 ms).
	{
		scheduler, err := scheduling.NewScheduler(ctx, fp, true)
		a.So(err, should.BeNil)
		scheduler.SyncWithGateway(0, time.Now(), time.Unix(0, 0))
		em, err := scheduler.ScheduleAnytime(ctx, 10, settingsAt(869525000, 7, timePtr(time.Unix(0, int64(100*time.Millisecond))), 0), ttnpb.TxSchedulePriority_NORMAL)
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.AlmostEqual, scheduling.ScheduleTimeShort, test.Delay>>6)
	}

	// Timestamp; too late (100 ms).
	{
		scheduler, err := scheduling.NewScheduler(ctx, fp, true)
		a.So(err, should.BeNil)
		scheduler.SyncWithGateway(0, time.Now(), time.Unix(0, 0))
		em, err := scheduler.ScheduleAnytime(ctx, 10, settingsAt(869525000, 7, nil, 100*1000), ttnpb.TxSchedulePriority_NORMAL)
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.AlmostEqual, scheduling.ScheduleTimeShort, test.Delay/1000)
	}
}

func TestScheduleAnytimeClassC(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fp := &frequencyplans.FrequencyPlan{
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}

	scheduler, err := scheduling.NewScheduler(ctx, fp, true)
	a.So(err, should.BeNil)

	// Schedule a join-accept. The clock is not synced; this is scheduled in the past at 5 seconds.
	_, err = scheduler.ScheduleAt(ctx, 10, ttnpb.TxSettings{
		DataRate: ttnpb.DataRate{
			Modulation: &ttnpb.DataRate_LoRa{
				LoRa: &ttnpb.LoRaDataRate{
					Bandwidth:       125000,
					SpreadingFactor: 7,
				},
			},
		},
		CodingRate: "4/5",
		Frequency:  868100000,
		Timestamp:  5000000,
	}, ttnpb.TxSchedulePriority_NORMAL)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	// Fast forward 9 seconds.
	scheduler.Sync(9000000, time.Now().Add(-3*time.Second))

	em, err := scheduler.ScheduleAnytime(ctx, 10, ttnpb.TxSettings{
		DataRate: ttnpb.DataRate{
			Modulation: &ttnpb.DataRate_LoRa{
				LoRa: &ttnpb.LoRaDataRate{
					Bandwidth:       125000,
					SpreadingFactor: 7,
				},
			},
		},
		CodingRate: "4/5",
		Frequency:  869525000,
	}, ttnpb.TxSchedulePriority_NORMAL)
	a.So(err, should.BeNil)
	a.So(time.Duration(em.Starts()), should.AlmostEqual, 12*time.Second+scheduling.ScheduleTimeLong, test.Delay>>6)
}
