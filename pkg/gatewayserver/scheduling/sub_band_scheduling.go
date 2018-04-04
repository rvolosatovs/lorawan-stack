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

package scheduling

import (
	"context"
	"sync"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/band"
	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
)

const (
	dutyCycleWindow = 5 * time.Minute
)

type packetWindow struct {
	window     Span
	timeOffAir Span
}

func (w packetWindow) withTimeOffAir() Span {
	return Span{
		Start:    w.window.Start,
		Duration: w.window.Duration + w.timeOffAir.Duration,
	}
}

type schedulingWindow struct {
	packetWindow

	added time.Time
}

type subBandScheduling struct {
	dutyCycle         band.DutyCycle
	schedulingWindows []schedulingWindow

	mu sync.Mutex
}

func (w schedulingWindow) expired() bool {
	return time.Now().Sub(w.added) > 2*dutyCycleWindow
}

func (s *subBandScheduling) bgCleanup(ctx context.Context) {
	cleanupTicker := time.NewTicker(dutyCycleWindow)
	for {
		select {
		case <-ctx.Done():
			cleanupTicker.Stop()
			return
		case <-cleanupTicker.C:
			s.mu.Lock()
			for i, w := range s.schedulingWindows {
				if w.expired() {
					s.schedulingWindows = append(s.schedulingWindows[:i], s.schedulingWindows[i+1:]...)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *subBandScheduling) addScheduling(w packetWindow) {
	window := schedulingWindow{packetWindow: w, added: time.Now()}
	for i, window := range s.schedulingWindows {
		if w.window.StartsBefore(window.window) {
			s.schedulingWindows = append(s.schedulingWindows[:i], append([]schedulingWindow{window}, s.schedulingWindows[i:]...)...)
			return
		}
	}
	s.schedulingWindows = append(s.schedulingWindows, window)
}

func (s *subBandScheduling) RegisterEmission(w packetWindow) {
	s.mu.Lock()
	s.addScheduling(w)
	s.mu.Unlock()
}

// Schedule adds the requested time window to its internal schedule. If, because of its internal constraints (e.g. for duty cycles, not respecting the duty cycle), it returns ErrScheduleFull. If another error prevents scheduling, it is returned.
func (s *subBandScheduling) ScheduleAt(w Span, timeOffAir *ttnpb.FrequencyPlan_TimeOffAir) error {
	s.mu.Lock()
	err := s.schedule(w, timeOffAir)
	s.mu.Unlock()
	return err
}

func (s *subBandScheduling) schedule(w Span, timeOffAir *ttnpb.FrequencyPlan_TimeOffAir) error {
	windowWithTimeOffAir := packetWindow{window: w, timeOffAir: w.timeOffAir(timeOffAir)}

	emissionWindows := []Span{w}

	for _, scheduledWindow := range s.schedulingWindows {
		emissionWindows = append(emissionWindows, scheduledWindow.window)

		if scheduledWindow.window.Overlaps(w) {
			return ErrOverlap.New(errors.Attributes{})
		}
		if scheduledWindow.withTimeOffAir().Overlaps(windowWithTimeOffAir.withTimeOffAir()) {
			return ErrTimeOffAir.New(errors.Attributes{})
		}
	}

	precedingWindowsAirtime := sumWithinInterval(emissionWindows, w.End().Add(-1*dutyCycleWindow), w.End())
	prolongingWindowsAirtime := sumWithinInterval(emissionWindows, w.Start, w.Start.Add(dutyCycleWindow))

	if prolongingWindowsAirtime > s.dutyCycle.MaxEmissionDuring(dutyCycleWindow) ||
		precedingWindowsAirtime > s.dutyCycle.MaxEmissionDuring(dutyCycleWindow) {
		return ErrDutyCycleFull.New(errors.Attributes{
			"min_frequency": s.dutyCycle.MinFrequency,
			"max_frequency": s.dutyCycle.MaxFrequency,
			"quota":         s.dutyCycle.DutyCycle,
		})
	}

	s.addScheduling(windowWithTimeOffAir)

	return nil
}

// ScheduleAnytime requires a scheduling window if there is no time.Time constraint.
func (s *subBandScheduling) ScheduleAnytime(minimum Timestamp, d time.Duration, timeOffAir *ttnpb.FrequencyPlan_TimeOffAir) (Span, error) {
	minimumSpan := Span{Start: minimum, Duration: d}
	if err := s.ScheduleAt(minimumSpan, timeOffAir); err == nil {
		return minimumSpan, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	potentialTimings := []Timestamp{}
	emissionWindows := []Span{}

	for _, window := range s.schedulingWindows {
		emissionWindows = append(emissionWindows, window.window)
		windowWithTimeOffAir := window.withTimeOffAir()

		windowEnd := windowWithTimeOffAir.End()
		if windowEnd.After(minimum) {
			potentialTimings = append(potentialTimings, windowEnd)
		}
	}

	var potentialTiming Timestamp
	for _, potentialTiming = range potentialTimings {
		w := Span{Start: potentialTiming, Duration: d}
		err := s.schedule(w, timeOffAir)
		if err != nil {
			continue
		}

		return w, nil
	}

	start := firstMomentConsideringDutyCycle(emissionWindows, s.dutyCycle.DutyCycle, potentialTiming, d)
	w := Span{Start: start, Duration: d}
	err := s.schedule(w, timeOffAir)
	return w, err
}

func firstMomentConsideringDutyCycle(spans []Span, dutyCycle float32, minimum Timestamp, duration time.Duration) Timestamp {
	if len(spans) == 0 {
		return minimum
	}

	maxAirtime := time.Duration(dutyCycle * float32(dutyCycleWindow))
	lastWindow := spans[len(spans)-1]

	precedingWindowsAirtime := sumWithinInterval(spans, minimum.Add(-1*dutyCycleWindow).Add(duration), minimum.Add(duration)) + duration

	margin := maxAirtime - (precedingWindowsAirtime - duration)
	return lastWindow.Start.Add(-1 * margin).Add(dutyCycleWindow)
}
