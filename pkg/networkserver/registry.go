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
	"context"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/internal/registry"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
)

type UplinkMatch interface {
	ApplicationIdentifiers() ttnpb.ApplicationIdentifiers
	DeviceID() string
	LoRaWANVersion() ttnpb.MACVersion
	FNwkSIntKey() *ttnpb.KeyEnvelope
	FCnt() uint32
	LastFCnt() uint32
	IsPending() bool
	ResetsFCnt() *pbtypes.BoolValue
}

// DeviceRegistry is a registry, containing devices.
type DeviceRegistry interface {
	GetByEUI(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, context.Context, error)
	GetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string) (*ttnpb.EndDevice, context.Context, error)
	RangeByUplinkMatches(ctx context.Context, up *ttnpb.UplinkMessage, cacheTTL time.Duration, f func(context.Context, UplinkMatch) (bool, error)) error
	SetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error)
}

var errDeviceExists = errors.DefineAlreadyExists("device_exists", "device already exists")

// CreateDevice creates device dev in r.
func CreateDevice(ctx context.Context, r DeviceRegistry, dev *ttnpb.EndDevice, paths ...string) (*ttnpb.EndDevice, context.Context, error) {
	return r.SetByID(ctx, dev.ApplicationIdentifiers, dev.DeviceID, ttnpb.EndDeviceFieldPathsTopLevel, func(_ context.Context, stored *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
		if stored != nil {
			return nil, nil, errDeviceExists
		}
		return dev, paths, nil
	})
}

// DeleteDevice deletes device identified by appID, devID from r.
func DeleteDevice(ctx context.Context, r DeviceRegistry, appID ttnpb.ApplicationIdentifiers, devID string) error {
	_, _, err := r.SetByID(ctx, appID, devID, nil, func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) { return nil, nil, nil })
	return err
}

func logRegistryRPCError(ctx context.Context, err error, msg string) {
	logger := log.FromContext(ctx).WithError(err)
	var printLog func(string)
	switch {
	case errors.IsNotFound(err), errors.IsInvalidArgument(err):
		printLog = logger.Debug
	case errors.IsFailedPrecondition(err):
		printLog = logger.Warn
	default:
		printLog = logger.Error
	}
	printLog(msg)
}

type replacedEndDeviceFieldRegistryWrapper struct {
	DeviceRegistry
	fields []registry.ReplacedEndDeviceField
}

func (w replacedEndDeviceFieldRegistryWrapper) GetByEUI(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, context.Context, error) {
	paths, replaced := registry.MatchReplacedEndDeviceFields(paths, w.fields)
	dev, ctx, err := w.DeviceRegistry.GetByEUI(ctx, joinEUI, devEUI, paths)
	if err != nil || dev == nil {
		return dev, ctx, err
	}
	for _, d := range replaced {
		d.GetTransform(dev)
	}
	return dev, ctx, nil
}

func (w replacedEndDeviceFieldRegistryWrapper) GetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string) (*ttnpb.EndDevice, context.Context, error) {
	paths, replaced := registry.MatchReplacedEndDeviceFields(paths, w.fields)
	dev, ctx, err := w.DeviceRegistry.GetByID(ctx, appID, devID, paths)
	if err != nil || dev == nil {
		return dev, ctx, err
	}
	for _, d := range replaced {
		d.GetTransform(dev)
	}
	return dev, ctx, nil
}

func (w replacedEndDeviceFieldRegistryWrapper) SetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
	paths, replaced := registry.MatchReplacedEndDeviceFields(paths, w.fields)
	dev, ctx, err := w.DeviceRegistry.SetByID(ctx, appID, devID, paths, func(ctx context.Context, dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
		if dev != nil {
			for _, d := range replaced {
				d.GetTransform(dev)
			}
		}
		dev, paths, err := f(ctx, dev)
		if err != nil || dev == nil {
			return dev, paths, err
		}
		for _, d := range replaced {
			if ttnpb.HasAnyField(paths, d.Old) {
				paths = ttnpb.AddFields(paths, d.New)
			}
			d.SetTransform(dev, d.MatchedOld, d.MatchedNew)
		}
		return dev, paths, nil
	})
	if err != nil || dev == nil {
		return dev, ctx, err
	}
	for _, d := range replaced {
		d.GetTransform(dev)
	}
	return dev, ctx, nil
}

func wrapEndDeviceRegistryWithReplacedFields(r DeviceRegistry, fields ...registry.ReplacedEndDeviceField) DeviceRegistry {
	return replacedEndDeviceFieldRegistryWrapper{
		DeviceRegistry: r,
		fields:         fields,
	}
}

var replacedEndDeviceFields = []registry.ReplacedEndDeviceField{
	{
		Old: "mac_state.current_parameters.adr_ack_delay",
		New: "mac_state.current_parameters.adr_ack_delay_exponent",
		GetTransform: func(dev *ttnpb.EndDevice) {
			if dev.MACState == nil {
				return
			}
			dev.MACState.CurrentParameters.ADRAckDelay = uint32(dev.MACState.CurrentParameters.ADRAckDelayExponent.GetValue())
		},
		SetTransform: func(dev *ttnpb.EndDevice, _, _ bool) error {
			if dev.MACState == nil {
				return nil
			}
			// Replicate old behavior for backwards-compatibility.
			dev.MACState.CurrentParameters.ADRAckDelay = 0
			return nil
		},
	},
	{
		Old: "mac_state.current_parameters.adr_ack_limit",
		New: "mac_state.current_parameters.adr_ack_limit_exponent",
		GetTransform: func(dev *ttnpb.EndDevice) {
			if dev.MACState == nil {
				return
			}
			dev.MACState.CurrentParameters.ADRAckLimit = uint32(dev.MACState.CurrentParameters.ADRAckLimitExponent.GetValue())
		},
		SetTransform: func(dev *ttnpb.EndDevice, _, _ bool) error {
			if dev.MACState == nil {
				return nil
			}
			// Replicate old behavior for backwards-compatibility.
			dev.MACState.CurrentParameters.ADRAckLimit = 0
			return nil
		},
	},
	{
		Old: "mac_state.current_parameters.ping_slot_data_rate_index",
		New: "mac_state.current_parameters.ping_slot_data_rate_index_value",
		GetTransform: func(dev *ttnpb.EndDevice) {
			if dev.MACState == nil {
				return
			}
			dev.MACState.CurrentParameters.PingSlotDataRateIndex = dev.MACState.CurrentParameters.PingSlotDataRateIndexValue.GetValue()
		},
		SetTransform: func(dev *ttnpb.EndDevice, _, _ bool) error {
			if dev.MACState == nil {
				return nil
			}
			// Replicate old behavior for backwards-compatibility.
			dev.MACState.CurrentParameters.PingSlotDataRateIndex = 0
			return nil
		},
	},
	{
		Old: "mac_state.desired_parameters.adr_ack_delay",
		New: "mac_state.desired_parameters.adr_ack_delay_exponent",
		GetTransform: func(dev *ttnpb.EndDevice) {
			if dev.MACState == nil {
				return
			}
			dev.MACState.DesiredParameters.ADRAckDelay = uint32(dev.MACState.DesiredParameters.ADRAckDelayExponent.GetValue())
		},
		SetTransform: func(dev *ttnpb.EndDevice, _, _ bool) error {
			if dev.MACState == nil {
				return nil
			}
			// Replicate old behavior for backwards-compatibility.
			dev.MACState.DesiredParameters.ADRAckDelay = 0
			return nil
		},
	},
	{
		Old: "mac_state.desired_parameters.adr_ack_limit",
		New: "mac_state.desired_parameters.adr_ack_limit_exponent",
		GetTransform: func(dev *ttnpb.EndDevice) {
			if dev.MACState == nil {
				return
			}
			dev.MACState.DesiredParameters.ADRAckLimit = uint32(dev.MACState.DesiredParameters.ADRAckLimitExponent.GetValue())
		},
		SetTransform: func(dev *ttnpb.EndDevice, _, _ bool) error {
			if dev.MACState == nil {
				return nil
			}
			// Replicate old behavior for backwards-compatibility.
			dev.MACState.DesiredParameters.ADRAckLimit = 0
			return nil
		},
	},
	{
		Old: "mac_state.desired_parameters.ping_slot_data_rate_index",
		New: "mac_state.desired_parameters.ping_slot_data_rate_index_value",
		GetTransform: func(dev *ttnpb.EndDevice) {
			if dev.MACState == nil {
				return
			}
			dev.MACState.DesiredParameters.PingSlotDataRateIndex = dev.MACState.DesiredParameters.PingSlotDataRateIndexValue.GetValue()
		},
		SetTransform: func(dev *ttnpb.EndDevice, _, _ bool) error {
			if dev.MACState == nil {
				return nil
			}
			// Replicate old behavior for backwards-compatibility.
			dev.MACState.DesiredParameters.PingSlotDataRateIndex = 0
			return nil
		},
	},
	{
		Old: "queued_application_downlinks",
		New: "session.queued_application_downlinks",
		GetTransform: func(dev *ttnpb.EndDevice) {
			switch {
			case dev.QueuedApplicationDownlinks == nil && dev.GetSession().GetQueuedApplicationDownlinks() == nil:
				return

			case dev.QueuedApplicationDownlinks != nil:
				if dev.Session == nil {
					dev.Session = &ttnpb.Session{}
				}
				dev.Session.QueuedApplicationDownlinks = dev.QueuedApplicationDownlinks

			default:
				dev.QueuedApplicationDownlinks = dev.Session.QueuedApplicationDownlinks
			}
		},
		SetTransform: func(dev *ttnpb.EndDevice, useOld, useNew bool) error {
			switch {
			case useOld && useNew:
				oldValue := dev.QueuedApplicationDownlinks
				newValue := dev.GetSession().GetQueuedApplicationDownlinks()
				n := len(oldValue)
				if n != len(newValue) {
					return errInvalidFieldValue.WithAttributes("field", "queued_application_downlinks")
				}
				for i := 0; i < n; i++ {
					if !oldValue[i].Equal(newValue[i]) {
						return errInvalidFieldValue.WithAttributes("field", "queued_application_downlinks")
					}
				}

			case useNew:
				dev.QueuedApplicationDownlinks = nil

			case dev.QueuedApplicationDownlinks == nil:
				if dev.Session != nil {
					dev.Session.QueuedApplicationDownlinks = nil
				}

			default:
				if dev.Session == nil {
					dev.Session = &ttnpb.Session{}
				}
				dev.Session.QueuedApplicationDownlinks = dev.QueuedApplicationDownlinks
			}
			dev.QueuedApplicationDownlinks = nil
			return nil
		},
	},
}
