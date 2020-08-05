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
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/scheduling"
	"go.thethings.network/lorawan-stack/v3/pkg/toa"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestScheduleAtWithBandDutyCycle(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fps := map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}}
	timeSource := &mockTimeSource{
		Time: time.Unix(0, 0),
	}
	scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
	a.So(err, should.BeNil)

	for i, tc := range []struct {
		SyncWithGatewayAbsolute bool
		PayloadSize             int
		Settings                ttnpb.TxSettings
		Priority                ttnpb.TxSchedulePriority
		NPercentileRTT          *time.Duration
		MedianRTT               *time.Duration
		ExpectedToa             time.Duration
		ExpectedStarts          scheduling.ConcentratorTime
		ExpectedError           *errors.Definition
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
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa: 41216 * time.Microsecond,
			// Too late for transmission with ScheduleTimeShort.
			ExpectedError: &scheduling.ErrTooLate,
		},
		{
			SyncWithGatewayAbsolute: true,
			PayloadSize:             51,
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
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa: 2465792 * time.Microsecond,
			// Too late for transmission with ScheduleTimeShort.
			ExpectedError: &scheduling.ErrTooLate,
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
				Timestamp:  300000,
			},
			Priority:       ttnpb.TxSchedulePriority_NORMAL,
			NPercentileRTT: durationPtr(550 * time.Millisecond),
			ExpectedToa:    41216 * time.Microsecond,
			// Too late for transmission with RTT.
			ExpectedError: &scheduling.ErrTooLate,
		},
		{
			SyncWithGatewayAbsolute: true,
			PayloadSize:             51,
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
				Time:       timePtr(time.Unix(0, int64(300*time.Millisecond))),
			},
			Priority:       ttnpb.TxSchedulePriority_NORMAL,
			NPercentileRTT: durationPtr(500 * time.Millisecond),
			ExpectedToa:    2465792 * time.Microsecond,
			// Too late for transmission with RTT.
			ExpectedError: &scheduling.ErrTooLate,
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
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa: 2465792 * time.Microsecond,
			// Exceeding dwell time of 2 seconds.
			ExpectedError: &scheduling.ErrDwellTime,
		},
		{
			PayloadSize: 16,
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
				Frequency:  868100000,
				Time:       timePtr(time.Unix(0, int64(1*time.Second))),
			},
			Priority:       ttnpb.TxSchedulePriority_HIGHEST,
			MedianRTT:      durationPtr(200 * time.Millisecond),
			ExpectedToa:    51456 * time.Microsecond,
			ExpectedStarts: 1000000000 - 200000000/2 - 51456000,
		},
		{
			SyncWithGatewayAbsolute: true,
			PayloadSize:             16,
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
				Frequency:  868100000,
				Time:       timePtr(time.Unix(0, int64(1*time.Minute))),
			},
			Priority:       ttnpb.TxSchedulePriority_HIGHEST,
			ExpectedToa:    51456 * time.Microsecond,
			ExpectedStarts: 60000000000 - 51456000,
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
			Priority:       ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa:    41216 * time.Microsecond,
			ExpectedStarts: 20000000000,
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
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa: 41216 * time.Microsecond,
			// Overlapping with previous transmission.
			ExpectedError: &scheduling.ErrConflict,
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
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa: 41216 * time.Microsecond,
			// Right after previous transmission; not respecting time-off-air.
			ExpectedError: &scheduling.ErrConflict,
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
			Priority:       ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa:    41216 * time.Microsecond,
			ExpectedStarts: 20000000000 + 41216000 + 1000000000,
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
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa: 1318912 * time.Microsecond,
			// Exceeds duty-cycle limitation of 1% in 868.0 - 868.6.
			ExpectedError: &scheduling.ErrDutyCycle,
		},
	} {
		tcok := t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assertions.New(t)
			if tc.SyncWithGatewayAbsolute {
				scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
			} else {
				scheduler.Sync(0, timeSource.Time)
			}
			d, err := toa.Compute(tc.PayloadSize, tc.Settings)
			a.So(err, should.BeNil)
			a.So(d, should.Equal, tc.ExpectedToa)
			rtts := &mockRTTs{}
			if tc.NPercentileRTT != nil {
				rtts.NPercentile = *tc.NPercentileRTT
				rtts.Count = 10
			}
			if tc.MedianRTT != nil {
				rtts.Median = *tc.MedianRTT
				rtts.Count = 10
			}
			em, err := scheduler.ScheduleAt(ctx, scheduling.Options{
				PayloadSize: tc.PayloadSize,
				TxSettings:  tc.Settings,
				RTTs:        rtts,
				Priority:    tc.Priority,
			})
			if tc.ExpectedError != nil {
				if !a.So(err, should.HaveSameErrorDefinitionAs, *tc.ExpectedError) {
					t.Fatalf("Unexpected error: %v", err)
				}
				return
			}
			if !a.So(err, should.BeNil) {
				t.FailNow()
			}
			a.So(em.Starts(), should.Equal, tc.ExpectedStarts)
		})
		if !tcok {
			t.FailNow()
		}
	}
}

func TestScheduleAtWithFrequencyPlanDutyCycle(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fps := map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
		BandID: band.EU_863_870,
		SubBands: []frequencyplans.SubBandParameters{
			{
				MinFrequency: 0,
				MaxFrequency: math.MaxUint64,
				DutyCycle:    0,
			},
		},
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
	}}
	timeSource := &mockTimeSource{
		Time: time.Unix(0, 0),
	}
	scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
	a.So(err, should.BeNil)

	for i, tc := range []struct {
		SyncWithGatewayAbsolute bool
		PayloadSize             int
		Settings                ttnpb.TxSettings
		Priority                ttnpb.TxSchedulePriority
		MaxRTT                  *time.Duration
		MedianRTT               *time.Duration
		ExpectedToa             time.Duration
		ExpectedStarts          scheduling.ConcentratorTime
	}{
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
				Timestamp:  20000000,
			},
			Priority:       ttnpb.TxSchedulePriority_NORMAL,
			ExpectedToa:    1318912 * time.Microsecond, // Normally violating duty-cycle limitation of 1% in 868.0 - 868.6.
			ExpectedStarts: 20000000000,
		},
	} {
		tcok := t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assertions.New(t)
			if tc.SyncWithGatewayAbsolute {
				scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
			} else {
				scheduler.Sync(0, timeSource.Time)
			}
			d, err := toa.Compute(tc.PayloadSize, tc.Settings)
			a.So(err, should.BeNil)
			a.So(d, should.Equal, tc.ExpectedToa)
			rtts := &mockRTTs{}
			if tc.MaxRTT != nil {
				rtts.Max = *tc.MaxRTT
				rtts.Count = 1
			}
			if tc.MedianRTT != nil {
				rtts.Median = *tc.MedianRTT
				rtts.Count = 1
			}
			em, err := scheduler.ScheduleAt(ctx, scheduling.Options{
				PayloadSize: tc.PayloadSize,
				TxSettings:  tc.Settings,
				RTTs:        rtts,
				Priority:    tc.Priority,
			})
			if !a.So(err, should.BeNil) {
				t.FailNow()
			}
			a.So(em.Starts(), should.Equal, tc.ExpectedStarts)
		})
		if !tcok {
			t.FailNow()
		}
	}
}

func TestScheduleAnytime(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fps := map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}}
	scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, nil)
	a.So(err, should.BeNil)
	scheduler.SyncWithGatewayAbsolute(0, time.Now(), time.Unix(0, 0))

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

	// Scheduling two items, considering time-on-air and time-off-air.
	// Time-on-air is 41216 us, time-off-air is 1000000 us.
	// 1: [1000000, 2041216]
	// 2: [4000000, 5041216]
	_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings:  settingsAt(869525000, 7, 1000000),
		Priority:    ttnpb.TxSchedulePriority_NORMAL,
	})
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings:  settingsAt(869525000, 7, 4000000),
		Priority:    ttnpb.TxSchedulePriority_NORMAL,
	})
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	// Try schedule a transmission from 1000000 us.
	// Time-on-air is 41216 us, time-off-air is 1000000 us.
	// It fits between 1 and 2, so it should be right after 1.
	// 1: [1000000, 2041216]
	// 3: [2041216, 3082432]
	// 2: [4000000, 5041216]
	em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings:  settingsAt(869525000, 7, 1000000),
		Priority:    ttnpb.TxSchedulePriority_NORMAL,
	})
	if !a.So(err, should.BeNil) || !a.So(em.Starts(), should.Equal, 2041216*time.Microsecond) {
		t.FailNow()
	}

	// Try schedule another transmission from 1000000 us.
	// Time-on-air is 41216 us, time-off-air is 1000000 us.
	// It does not fit between 1, 3 and 2, so it should be right after 2.
	// 1: [1000000, 2041216]
	// 3: [2041216, 3082432]
	// 2: [4000000, 5041216]
	// 4: [5041216, 5082432]
	em, err = scheduler.ScheduleAnytime(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings:  settingsAt(869525000, 7, 1000000),
		Priority:    ttnpb.TxSchedulePriority_NORMAL,
	})
	if !a.So(err, should.BeNil) || !a.So(em.Starts(), should.Equal, 5041216*time.Microsecond) {
		t.FailNow()
	}

	// Try schedule another transmission from 1000000 us.
	// Time-on-air is 991232 us, time-off-air is 1000000 us.
	// It's 9.91% in a 10% duty-cycle sub-band, almost hitting the limit, so it should be pushed to right after transmission 4.
	// Transmission starts then at 5041216 (start of 4) + 41216 (time-on-air of 4) + 10000000 (duty-cycle window) - 991232 (this time-on-air).
	// 1: [1000000, 2041216]
	// 3: [2041216, 3082432]
	// 2: [4000000, 5041216]
	// 4: [5041216, 5082432]
	// 5: [14091200, 15082432]
	em, err = scheduler.ScheduleAnytime(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings:  settingsAt(869525000, 12, 1000000),
		Priority:    ttnpb.TxSchedulePriority_HIGHEST,
	})
	if !a.So(err, should.BeNil) || !a.So(em.Starts(), should.Equal, 14091200*time.Microsecond) {
		t.FailNow()
	}

	// Try schedule another transmission from 1000000 us.
	// Time-on-air is 991232 us, time-off-air is 1000000 us.
	// It's 9.91% in a 1% duty-cycle sub-band, so it hits the duty-cycle limitation.
	_, err = scheduler.ScheduleAnytime(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings:  settingsAt(868100000, 12, 1000000),
		Priority:    ttnpb.TxSchedulePriority_HIGHEST,
	})
	a.So(err, should.HaveSameErrorDefinitionAs, scheduling.ErrDutyCycle)
}

func TestScheduleAnytimeShort(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fps := map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}}

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

	// Gateway time; too late (100 ms) without RTT.
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 0),
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, scheduling.ScheduleTimeLong)
	}

	// Gateway time; too late (100 ms) without RTT, with high schedule_anytime_delay
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduleAnytimeDelay := time.Second
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, &scheduleAnytimeDelay, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 0),
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, time.Second)
	}

	// Gateway time; too late (100 ms) without RTT, with low schedule_anytime_delay
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduleAnytimeDelay := time.Millisecond
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, &scheduleAnytimeDelay, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 0),
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, scheduling.ScheduleTimeShort)
	}

	// Gateway time; too late (100 ms) without RTT, with zero schedule_anytime_delay
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduleAnytimeDelay := time.Duration(0)
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, &scheduleAnytimeDelay, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 0),
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, scheduling.ScheduleTimeLong)
	}

	// Gateway time; too late (10 ms) with RTT.
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		rtts := &mockRTTs{
			Max:   40 * time.Millisecond,
			Count: 1,
		}
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 0),
			RTTs:        rtts,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, scheduling.ScheduleTimeLong)
	}

	// Timestamp; too late (100 ms) without RTT.
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 100*1000),
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, scheduling.ScheduleTimeShort)
	}

	// Timestamp; too late (10 ms) with too few RTTs (5).
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		rtts := &mockRTTs{
			NPercentile: 40 * time.Millisecond,
			Count:       3,
		}
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 10*1000),
			RTTs:        rtts,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, scheduling.ScheduleTimeShort)
	}

	// Timestamp; too late (10 ms) with enough RTTs.
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		scheduler.SyncWithGatewayAbsolute(0, timeSource.Time, time.Unix(0, 0))
		rtts := &mockRTTs{
			NPercentile: 40 * time.Millisecond,
			Count:       20,
		}
		em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings:  settingsAt(869525000, 7, nil, 10*1000),
			RTTs:        rtts,
			Priority:    ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
		a.So(time.Duration(em.Starts()), should.Equal, 40*time.Millisecond+scheduling.QueueDelay)
	}
}

func TestScheduleAnytimeClassC(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fps := map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}}

	timeSource := &mockTimeSource{
		Time: time.Unix(0, 0),
	}
	scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
	a.So(err, should.BeNil)
	scheduler.Sync(0, timeSource.Time)

	// Schedule a join-accept.
	_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings: ttnpb.TxSettings{
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
		},
		Priority: ttnpb.TxSchedulePriority_HIGHEST,
	})
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	// Schedule another class A transmission after.
	_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings: ttnpb.TxSettings{
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
			Timestamp:  7000000,
		},
		Priority: ttnpb.TxSchedulePriority_HIGHEST,
	})
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	// Fast forward 9 seconds.
	timeSource.Time = time.Unix(9, 0)
	scheduler.Sync(9000000, timeSource.Time)

	// Schedule any time.
	em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings: ttnpb.TxSettings{
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
		},
		Priority: ttnpb.TxSchedulePriority_HIGHEST,
	})
	a.So(err, should.BeNil)
	a.So(time.Duration(em.Starts()), should.Equal, 9*time.Second+scheduling.ScheduleTimeLong)
}

func TestSchedulerWithMultipleFrequencyPlans(t *testing.T) {
	ctx := test.Context()
	for _, tc := range []struct {
		Name                 string
		FrequencyPlans       map[string]*frequencyplans.FrequencyPlan
		ExpectedSubBandCount int
		ErrorAssertion       func(error) bool
	}{
		{
			Name: "RepeatedNoOverlap",
			FrequencyPlans: map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
				BandID: band.EU_863_870,
				TimeOffAir: frequencyplans.TimeOffAir{
					Duration: time.Second,
				},
				DwellTime: frequencyplans.DwellTime{
					Downlinks: boolPtr(true),
					Duration:  durationPtr(2 * time.Second),
				},
			},
				"EU_863_870_Custom": {
					BandID: band.EU_863_870,
					TimeOffAir: frequencyplans.TimeOffAir{
						Duration: time.Second,
					},
					DwellTime: frequencyplans.DwellTime{
						Downlinks: boolPtr(true),
						Duration:  durationPtr(2 * time.Second),
					},
				},
			},
			ExpectedSubBandCount: 6,
		},
		{
			Name: "UnionOfNonOverlapping",
			FrequencyPlans: map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
				BandID: band.EU_863_870,
				TimeOffAir: frequencyplans.TimeOffAir{
					Duration: time.Second,
				},
				DwellTime: frequencyplans.DwellTime{
					Downlinks: boolPtr(true),
					Duration:  durationPtr(2 * time.Second),
				},
				SubBands: []frequencyplans.SubBandParameters{
					// Fictional Band S
					{
						MinFrequency: 870000000,
						MaxFrequency: 875000000,
						DutyCycle:    0.01,
						MaxEIRP:      float32Ptr(16.25),
					},
				},
			},
				"EU_863_870_Custom": {
					BandID: band.EU_863_870,
					TimeOffAir: frequencyplans.TimeOffAir{
						Duration: time.Second,
					},
					DwellTime: frequencyplans.DwellTime{
						Downlinks: boolPtr(true),
						Duration:  durationPtr(2 * time.Second),
					},
				},
			},
			ExpectedSubBandCount: 7,
		},
		{
			Name: "MismatchedTimeOffAir",
			FrequencyPlans: map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
				BandID: band.EU_863_870,
				TimeOffAir: frequencyplans.TimeOffAir{
					Duration: 2 * time.Second,
				},
				DwellTime: frequencyplans.DwellTime{
					Downlinks: boolPtr(true),
					Duration:  durationPtr(2 * time.Second),
				},
			},
				"EU_863_870_Custom": {
					BandID: band.EU_863_870,
					TimeOffAir: frequencyplans.TimeOffAir{
						Duration: time.Second,
					},
					DwellTime: frequencyplans.DwellTime{
						Downlinks: boolPtr(true),
						Duration:  durationPtr(2 * time.Second),
					},
				},
			},
			ErrorAssertion: func(err error) bool {
				return errors.IsInvalidArgument(err)
			},
		},
		{
			Name: "OverlappingSubBands",
			FrequencyPlans: map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
				BandID: band.EU_863_870,
				TimeOffAir: frequencyplans.TimeOffAir{
					Duration: time.Second,
				},
				DwellTime: frequencyplans.DwellTime{
					Downlinks: boolPtr(true),
					Duration:  durationPtr(2 * time.Second),
				},
				SubBands: []frequencyplans.SubBandParameters{
					// Fictional Band S
					{
						MinFrequency: 869000000,
						MaxFrequency: 873000000,
						DutyCycle:    0.01,
						MaxEIRP:      float32Ptr(16.25),
					},
				},
			},
				"EU_863_870_Custom": {
					BandID: band.EU_863_870,
					TimeOffAir: frequencyplans.TimeOffAir{
						Duration: time.Second,
					},
					DwellTime: frequencyplans.DwellTime{
						Downlinks: boolPtr(true),
						Duration:  durationPtr(2 * time.Second),
					},
				},
			},
			ErrorAssertion: func(err error) bool {
				return errors.IsInvalidArgument(err)
			},
		},
		{
			Name: "OverlappingSubBandsFromBand",
			FrequencyPlans: map[string]*frequencyplans.FrequencyPlan{"AS_923": {
				// This is a fictional test case since currently we don't support mix-band frequency plans (https://github.com/TheThingsNetwork/lorawan-stack/issues/1394).
				BandID: band.AS_923,
				TimeOffAir: frequencyplans.TimeOffAir{
					Duration: time.Second,
				},
				DwellTime: frequencyplans.DwellTime{
					Downlinks: boolPtr(true),
					Duration:  durationPtr(2 * time.Second),
				},
			},
				"AU_915_928": {
					BandID: band.AU_915_928,
					TimeOffAir: frequencyplans.TimeOffAir{
						Duration: time.Second,
					},
					DwellTime: frequencyplans.DwellTime{
						Downlinks: boolPtr(true),
						Duration:  durationPtr(2 * time.Second),
					},
				},
			},
			ErrorAssertion: func(err error) bool {
				return errors.IsInvalidArgument(err)
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			timeSource := &mockTimeSource{
				Time: time.Unix(0, 0),
			}
			scheduler, err := scheduling.NewScheduler(ctx, tc.FrequencyPlans, true, nil, timeSource)
			if err != nil {
				if tc.ErrorAssertion == nil || !a.So(tc.ErrorAssertion(err), should.BeTrue) {
					t.Fatalf("Unexpected error: %v", err)
				}
			} else if tc.ErrorAssertion != nil {
				t.Fatalf("Expected error")
			} else {
				if !a.So(scheduler.SubBandCount(), should.Equal, tc.ExpectedSubBandCount) {
					t.Fatalf("Invalid number of sub bands: %v", scheduler.SubBandCount())
				}
			}
		})
	}
}

func TestSchedulingWithMultipleFrequencyPlans(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fps := map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
		SubBands: []frequencyplans.SubBandParameters{
			// Fictional Band S
			{
				MinFrequency: 870000000,
				MaxFrequency: 875000000,
				DutyCycle:    0.01,
				MaxEIRP:      float32Ptr(16.25),
			},
		},
	},
		"EU_863_870_Custom": {
			BandID: band.EU_863_870,
			TimeOffAir: frequencyplans.TimeOffAir{
				Duration: time.Second,
			},
			DwellTime: frequencyplans.DwellTime{
				Downlinks: boolPtr(true),
				Duration:  durationPtr(2 * time.Second),
			},
		},
	}

	timeSource := &mockTimeSource{
		Time: time.Unix(0, 0),
	}
	scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
	a.So(err, should.BeNil)
	scheduler.Sync(0, timeSource.Time)

	// Schedule a join-accept.
	_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings: ttnpb.TxSettings{
			DataRate: ttnpb.DataRate{
				Modulation: &ttnpb.DataRate_LoRa{
					LoRa: &ttnpb.LoRaDataRate{
						Bandwidth:       125000,
						SpreadingFactor: 7,
					},
				},
			},
			CodingRate: "4/5",
			Frequency:  871100000,
			Timestamp:  5000000,
		},
		Priority: ttnpb.TxSchedulePriority_HIGHEST,
	})
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	// Schedule another class A transmission after.
	_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings: ttnpb.TxSettings{
			DataRate: ttnpb.DataRate{
				Modulation: &ttnpb.DataRate_LoRa{
					LoRa: &ttnpb.LoRaDataRate{
						Bandwidth:       125000,
						SpreadingFactor: 7,
					},
				},
			},
			CodingRate: "4/5",
			Frequency:  872100000,
			Timestamp:  7000000,
		},
		Priority: ttnpb.TxSchedulePriority_HIGHEST,
	})
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	// Fast forward 9 seconds.
	timeSource.Time = time.Unix(9, 0)
	scheduler.Sync(9000000, timeSource.Time)

	// Schedule any time.
	em, err := scheduler.ScheduleAnytime(ctx, scheduling.Options{
		PayloadSize: 10,
		TxSettings: ttnpb.TxSettings{
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
		},
		Priority: ttnpb.TxSchedulePriority_HIGHEST,
	})
	a.So(err, should.BeNil)
	a.So(time.Duration(em.Starts()), should.Equal, 9*time.Second+scheduling.ScheduleTimeLong)
}

func TestScheduleSyncViaUplinkToken(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()
	fps := map[string]*frequencyplans.FrequencyPlan{test.EUFrequencyPlanID: {
		BandID: band.EU_863_870,
		TimeOffAir: frequencyplans.TimeOffAir{
			Duration: time.Second,
		},
		DwellTime: frequencyplans.DwellTime{
			Downlinks: boolPtr(true),
			Duration:  durationPtr(2 * time.Second),
		},
	}}

	// No clock sync and no uplink token.
	{
		timeSource := &mockTimeSource{
			Time: time.Now(),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings: ttnpb.TxSettings{
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
			},
			Priority: ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(errors.IsUnavailable(err), should.BeTrue)
	}

	// No clock sync, clock sync via uplink token, scheduling takes place too late.
	{
		timeSource := &mockTimeSource{
			Time: time.Unix(11, 0),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings: ttnpb.TxSettings{
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
				Timestamp:  6000000,
			},
			UplinkToken: &ttnpb.UplinkToken{
				ServerTime:       time.Unix(10, 0),
				Timestamp:        5000000,
				ConcentratorTime: 5000000000,
			},
			Priority: ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(errors.IsFailedPrecondition(err), should.BeTrue)
	}

	// No clock sync, clock sync via uplink token, scheduling works.
	{
		timeSource := &mockTimeSource{
			Time: time.Unix(11, 0),
		}
		scheduler, err := scheduling.NewScheduler(ctx, fps, true, nil, timeSource)
		a.So(err, should.BeNil)
		_, err = scheduler.ScheduleAt(ctx, scheduling.Options{
			PayloadSize: 10,
			TxSettings: ttnpb.TxSettings{
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
				Timestamp:  7000000,
			},
			UplinkToken: &ttnpb.UplinkToken{
				ServerTime:       time.Unix(10, 0),
				Timestamp:        5000000,
				ConcentratorTime: 5000000000,
			},
			Priority: ttnpb.TxSchedulePriority_NORMAL,
		})
		a.So(err, should.BeNil)
	}
}
