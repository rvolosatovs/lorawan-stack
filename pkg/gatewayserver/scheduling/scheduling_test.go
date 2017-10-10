// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package scheduling_test

import (
	"testing"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/band"
	"github.com/TheThingsNetwork/ttn/pkg/gatewayserver/scheduling"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func emptyEUFrequencyPlan() ttnpb.FrequencyPlan {
	return ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
	}
}

func TestEmptyScheduler(t *testing.T) {
	a := assertions.New(t)

	s, err := scheduling.FrequencyPlanScheduler(emptyEUFrequencyPlan())
	a.So(err, should.BeNil)

	askingTime := time.Now().Add(time.Minute)
	askingDuration := time.Second
	_, err = s.AskScheduling(askingTime, askingDuration, 0)
	a.So(err, should.NotBeNil)

	w := scheduling.Window{Start: askingTime, Duration: askingDuration}
	err = s.RegisterEmission(w, 0)
	a.So(err, should.NotBeNil)

	err = s.Schedule(w, 0)
	a.So(err, should.NotBeNil)
}

func TestInexistantBand(t *testing.T) {
	a := assertions.New(t)

	_, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{})
	a.So(err, should.NotBeNil)
}

func TestDwellTimeBlocking(t *testing.T) {
	a := assertions.New(t)

	dwellTime := time.Microsecond
	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID:    string(band.EU_863_870),
		DwellTime: &dwellTime,
	})
	a.So(err, should.BeNil)

	err = s.Schedule(scheduling.Window{Start: time.Now(), Duration: time.Minute}, 0)
	a.So(err, should.NotBeNil)
}

func TestAskScheduling(t *testing.T) {
	a := assertions.New(t)

	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
	})
	a.So(err, should.BeNil)

	askingTime := time.Now()
	err = s.Schedule(scheduling.Window{Start: askingTime, Duration: time.Microsecond}, 863000000)
	a.So(err, should.BeNil)

	schedule, err := s.AskScheduling(askingTime.Add(time.Hour), time.Microsecond, 863000000)
	a.So(err, should.BeNil)
	a.So(schedule.Start, should.Equal, askingTime.Add(time.Hour))
	a.So(schedule.Duration, should.Equal, time.Microsecond)

	_, err = s.AskScheduling(askingTime.Add(time.Hour).Add(-1*time.Microsecond), time.Minute, 863000000)
	a.So(err, should.NotBeNil)
}

func TestAskScheduling2(t *testing.T) {
	a := assertions.New(t)

	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
	})
	a.So(err, should.BeNil)

	askingTime := time.Now()
	err = s.Schedule(scheduling.Window{Start: askingTime, Duration: time.Microsecond}, 863000000)
	a.So(err, should.BeNil)

	schedule, err := s.AskScheduling(askingTime.Add(time.Hour), time.Microsecond, 863000000)
	a.So(err, should.BeNil)
	a.So(schedule.Start, should.Equal, askingTime.Add(time.Hour))
	a.So(schedule.Duration, should.Equal, time.Microsecond)

	schedule2, err := s.AskScheduling(askingTime.Add(time.Hour), time.Microsecond, 863000000)
	a.So(err, should.BeNil)
	a.So(schedule2.Start, should.Equal, askingTime.Add(time.Hour).Add(time.Microsecond))
	a.So(schedule2.Duration, should.Equal, time.Microsecond)
}

func TestAskSchedulingFullDutyCycle(t *testing.T) {
	a := assertions.New(t)

	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
	})
	a.So(err, should.BeNil)

	askingTime := time.Now()
	scheduleDuration := time.Duration(180 * time.Millisecond)
	err = s.Schedule(scheduling.Window{Start: askingTime, Duration: scheduleDuration}, 863000000)
	a.So(err, should.BeNil)

	schedule, err := s.AskScheduling(askingTime, scheduleDuration, 863000000)
	a.So(err, should.BeNil)
	expectedSchedule2Time := askingTime.Add(5 * time.Minute).Add(-120 * time.Millisecond)
	a.So(schedule.Start, should.Equal, expectedSchedule2Time)
	a.So(schedule.Duration, should.Equal, scheduleDuration)

	schedule, err = s.AskScheduling(askingTime, scheduleDuration, 863000000)
	a.So(err, should.BeNil)
	expectedSchedule3Time := expectedSchedule2Time.Add(5 * time.Minute).Add(-120 * time.Millisecond)
	a.So(schedule.Start, should.Equal, expectedSchedule3Time)
	a.So(schedule.Duration, should.Equal, scheduleDuration)
}

func TestAskSchedulingFullDutyCycleAfterRegisteredEmission(t *testing.T) {
	a := assertions.New(t)

	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
	})
	a.So(err, should.BeNil)

	askingTime := time.Now()
	scheduleDuration := time.Duration(180 * time.Millisecond)
	err = s.RegisterEmission(scheduling.Window{Start: askingTime, Duration: scheduleDuration}, 863000000)
	a.So(err, should.BeNil)

	schedule, err := s.AskScheduling(askingTime, scheduleDuration, 863000000)
	a.So(err, should.BeNil)
	expectedSchedule2Time := askingTime.Add(5 * time.Minute).Add(-120 * time.Millisecond)
	a.So(schedule.Start, should.Equal, expectedSchedule2Time)
	a.So(schedule.Duration, should.Equal, scheduleDuration)
}

func TestScheduleFullDutyCycle(t *testing.T) {
	a := assertions.New(t)

	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
	})
	a.So(err, should.BeNil)

	askingTime := time.Now()
	scheduleDuration := time.Duration(180 * time.Millisecond)
	err = s.Schedule(scheduling.Window{Start: askingTime, Duration: scheduleDuration}, 863000000)
	a.So(err, should.BeNil)

	err = s.Schedule(scheduling.Window{Start: askingTime, Duration: scheduleDuration}, 863000000)
	a.So(err, should.NotBeNil)

	err = s.Schedule(scheduling.Window{Start: askingTime.Add(200 * time.Millisecond), Duration: scheduleDuration}, 863000000)
	a.So(err, should.NotBeNil)

	err = s.Schedule(scheduling.Window{Start: askingTime.Add(5 * time.Minute).Add(-120 * time.Millisecond), Duration: scheduleDuration}, 863000000)
	a.So(err, should.BeNil)
}

func TestScheduleOrdering(t *testing.T) {
	a := assertions.New(t)

	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
	})
	a.So(err, should.BeNil)

	askingTime := time.Now().Add(time.Minute)
	scheduleDuration := time.Duration(time.Millisecond)

	w, err := s.AskScheduling(askingTime, scheduleDuration, 863000000)
	a.So(err, should.BeNil)
	a.So(w.Start, should.Equal, askingTime)

	err = s.Schedule(scheduling.Window{Start: askingTime.Add(time.Second), Duration: scheduleDuration}, 863000000)
	a.So(err, should.BeNil)

	err = s.Schedule(scheduling.Window{Start: askingTime.Add(50 * scheduleDuration), Duration: scheduleDuration}, 863000000)
	a.So(err, should.BeNil)
}

func TestTimeOffAirError(t *testing.T) {
	a := assertions.New(t)

	toa := time.Minute
	s, err := scheduling.FrequencyPlanScheduler(ttnpb.FrequencyPlan{
		BandID: string(band.EU_863_870),
		TimeOffAir: &ttnpb.FrequencyPlan_TimeOffAir{
			Duration: &toa,
		},
	})
	a.So(err, should.BeNil)

	askingTime := time.Now()
	scheduleDuration := time.Duration(60 * time.Millisecond)
	err = s.Schedule(scheduling.Window{Start: askingTime, Duration: scheduleDuration}, 863000000)
	a.So(err, should.BeNil)

	err = s.Schedule(scheduling.Window{Start: askingTime.Add(90 * time.Millisecond), Duration: scheduleDuration}, 863000000)
	a.So(err, should.NotBeNil)
}
