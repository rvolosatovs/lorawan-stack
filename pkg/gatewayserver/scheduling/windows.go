// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package scheduling

import (
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
)

// Span is a time window requested for usage by the consumer entity.
type Span struct {
	Start    time.Time
	Duration time.Duration
}

// End returns the timestamp at which the timespan ends
func (s Span) End() time.Time {
	return s.Start.Add(s.Duration)
}

// Precedes returns true if a portion of the timespan is located before the span passed as a parameter
func (s Span) Precedes(newS Span) bool {
	return s.Start.Before(newS.Start)
}

// Contains returns true if the given time is contained between the beginning and the end of this timespan
func (s Span) Contains(ts time.Time) bool {
	if ts.Before(s.Start) {
		return false
	}
	if ts.After(s.End()) {
		return false
	}
	return true
}

// IsProlongedBy returns true if after the span ends, there is still a portion of the span passed as parameter
func (s Span) IsProlongedBy(newS Span) bool {
	return s.End().Before(newS.End())
}

// Overlaps returns true if the two timespans overlap
func (s Span) Overlaps(newS Span) bool {
	if newS.End().Before(s.Start) || newS.End() == s.Start {
		return false
	}

	if s.End().Before(newS.Start) || s.End() == newS.Start {
		return false
	}

	return true
}

func filterWithinInterval(spans []Span, start, end time.Time) []Span {
	filteredSpans := []Span{}

	for _, span := range spans {
		if span.End().Before(start) || span.Start.After(end) {
			continue
		}

		filteredSpans = append(filteredSpans, span)
	}

	return filteredSpans
}

// spanDurationSum takes an array of non-overlapping spans, a start and an end time, and determines the sum of the duration of these spans within the interval determined by start and end.
func spanDurationSum(spans []Span, start, end time.Time) time.Duration {
	var duration time.Duration

	spans = filterWithinInterval(spans, start, end)
	for _, span := range spans {
		spanInterval := struct {
			start, end time.Time
		}{start: span.Start, end: span.End()}

		if spanInterval.start.Before(start) {
			spanInterval.start = start
		}
		if spanInterval.end.After(end) {
			spanInterval.end = end
		}

		duration = duration + spanInterval.end.Sub(spanInterval.start)
	}

	return duration
}

func (s Span) timeOffAir(timeOffAir *ttnpb.FrequencyPlan_TimeOffAir) (timeOffAirSpan Span) {
	timeOffAirSpan = Span{Start: s.End(), Duration: 0}

	if timeOffAir == nil {
		return
	}

	timeOffAirSpan.Duration = time.Duration(timeOffAir.Fraction * float32(s.Duration))
	if timeOffAir.Duration != nil {
		if *timeOffAir.Duration > timeOffAirSpan.Duration {
			timeOffAirSpan.Duration = *timeOffAir.Duration
		}
	}

	return
}
