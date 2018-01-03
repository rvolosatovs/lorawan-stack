// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

// Package scheduling offer convenience methods to manage RF packets that must respect scheduling constraints
package scheduling

import (
	"context"
	"sync"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/band"
	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
)

var (
	// ErrDutyCycleFull is returned is the duty cycle prevents scheduling of a downlink
	ErrDutyCycleFull = &errors.ErrDescriptor{
		Code:           1,
		MessageFormat:  "Duty cycle between { min_frequency } and { max_frequency } full, exceeded quota of { quota }",
		SafeAttributes: []string{"min_frequency", "max_frequency", "quota"},
	}
	// ErrOverlap is returned if there is an already existing scheduling overlapping
	ErrOverlap = &errors.ErrDescriptor{
		Code:          2,
		MessageFormat: "Window overlap",
	}
	// ErrTimeOffAir is returned if time-off-air constraints prevent scheduling of the new downlink
	ErrTimeOffAir = &errors.ErrDescriptor{
		Code:          3,
		MessageFormat: "Time-off-air constraints prevent scheduling",
	}
	// ErrNoSubBandFound is returned when an operation fails because there is no sub band for the given channel
	ErrNoSubBandFound = &errors.ErrDescriptor{
		Code:           4,
		MessageFormat:  "No sub band found for the given channel { channel }",
		SafeAttributes: []string{"channel"},
	}
	// ErrDwellTime is returned when an operation fails because the packet does not respect the dwell time
	ErrDwellTime = &errors.ErrDescriptor{
		Code:           5,
		MessageFormat:  "Packet time-on-air duration is greater than this band's dwell time ({ packet_duration } > { dwell_time })",
		SafeAttributes: []string{"packet_duration", "dwell_time"},
	}
)

func init() {
	ErrDutyCycleFull.Register()
	ErrOverlap.Register()
	ErrTimeOffAir.Register()
	ErrNoSubBandFound.Register()
	ErrDwellTime.Register()
}

// Scheduler is an abstraction for an entity that manages the packet's timespans.
type Scheduler interface {
	// ScheduleAt adds the requested timespan to its internal schedule. If, because of its internal constraints (e.g. for duty cycles, not respecting the duty cycle), it returns ErrScheduleFull. If another error prevents scheduling, it is returned.
	ScheduleAt(s Span, channel uint64) error
	// ScheduleAnytime requires a scheduling window if there is no time.Time constraint
	ScheduleAnytime(minimum time.Time, d time.Duration, channel uint64) (Span, error)
	// RegisterEmission that has happened during that timespan, on that specific channel
	RegisterEmission(s Span, channel uint64) error
}

// FrequencyPlanScheduler returns a scheduler based on the frequency plan, and starts a goroutine for cleanup. The scheduler is based on the dwell time, time off air, and the frequency plan's band. Assumption is made that no two duty cycles on a given band overlap.
func FrequencyPlanScheduler(ctx context.Context, fp ttnpb.FrequencyPlan) (Scheduler, error) {
	scheduler := &frequencyPlanScheduling{
		dwellTime:  fp.DwellTime,
		timeOffAir: fp.TimeOffAir,
		subBands:   []*subBandScheduling{},
	}

	band, err := band.GetByID(fp.BandID)
	if err != nil {
		return nil, errors.NewWithCause("Could not find band associated to that frequency plan", err)
	}

	for _, subBand := range band.BandDutyCycles {
		scheduling := &subBandScheduling{
			dutyCycle:         subBand,
			schedulingWindows: []packetWindow{},

			mu: sync.Mutex{},
		}
		scheduler.subBands = append(scheduler.subBands, scheduling)
		go scheduling.bgCleanup(ctx)
	}

	return scheduler, nil
}

type frequencyPlanScheduling struct {
	dwellTime  *time.Duration
	timeOffAir *ttnpb.FrequencyPlan_TimeOffAir

	subBands []*subBandScheduling
}

func (f frequencyPlanScheduling) findSubBand(channel uint64) (*subBandScheduling, error) {
	for _, subBand := range f.subBands {
		if subBand.dutyCycle.Comprises(channel) {
			return subBand, nil
		}
	}

	return nil, ErrNoSubBandFound.New(errors.Attributes{"channel": channel})
}

func (f frequencyPlanScheduling) ScheduleAt(s Span, channel uint64) error {
	if f.dwellTime != nil && s.Duration > *f.dwellTime {
		return ErrDwellTime.New(errors.Attributes{"packet_duration": s.Duration, "dwell_time": *f.dwellTime})
	}

	subBand, err := f.findSubBand(channel)
	if err != nil {
		return err
	}

	return subBand.ScheduleAt(s, f.timeOffAir)
}

func (f frequencyPlanScheduling) ScheduleAnytime(minimum time.Time, d time.Duration, channel uint64) (Span, error) {
	subBand, err := f.findSubBand(channel)
	if err != nil {
		return Span{}, err
	}

	return subBand.ScheduleAnytime(minimum, d, f.timeOffAir)
}

func (f frequencyPlanScheduling) RegisterEmission(s Span, channel uint64) error {
	subBand, err := f.findSubBand(channel)
	if err != nil {
		return err
	}

	subBand.RegisterEmission(packetWindow{window: s, timeOffAir: s.timeOffAir(f.timeOffAir)})
	return nil
}
