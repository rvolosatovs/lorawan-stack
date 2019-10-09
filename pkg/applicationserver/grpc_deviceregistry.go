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

package applicationserver

import (
	"context"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/crypto/cryptoutil"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	evtCreateEndDevice = events.Define(
		"as.end_device.create", "create end device",
		ttnpb.RIGHT_APPLICATION_DEVICES_READ,
	)
	evtUpdateEndDevice = events.Define(
		"as.end_device.update", "update end device",
		ttnpb.RIGHT_APPLICATION_DEVICES_READ,
	)
	evtDeleteEndDevice = events.Define(
		"as.end_device.delete", "delete end device",
		ttnpb.RIGHT_APPLICATION_DEVICES_READ,
	)
)

type asEndDeviceRegistryServer struct {
	AS *ApplicationServer
}

// Get implements ttnpb.AsEndDeviceRegistryServer.
func (r asEndDeviceRegistryServer) Get(ctx context.Context, req *ttnpb.GetEndDeviceRequest) (*ttnpb.EndDevice, error) {
	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_READ); err != nil {
		return nil, err
	}
	dev, err := r.AS.deviceRegistry.Get(ctx, req.EndDeviceIdentifiers, req.FieldMask.Paths)
	if err != nil {
		return nil, err
	}
	if ttnpb.HasAnyField(req.FieldMask.Paths, "session.keys.app_s_key") && dev.Session != nil && dev.Session.AppSKey != nil {
		key, err := cryptoutil.UnwrapAES128Key(ctx, *dev.Session.AppSKey, r.AS.KeyVault)
		if err != nil {
			return nil, err
		}
		dev.Session.AppSKey = &ttnpb.KeyEnvelope{Key: &key}
	}
	if ttnpb.HasAnyField(req.FieldMask.Paths, "pending_session.keys.app_s_key") && dev.PendingSession != nil && dev.PendingSession.AppSKey != nil {
		key, err := cryptoutil.UnwrapAES128Key(ctx, *dev.PendingSession.AppSKey, r.AS.KeyVault)
		if err != nil {
			return nil, err
		}
		dev.PendingSession.AppSKey = &ttnpb.KeyEnvelope{Key: &key}
	}
	return dev, nil
}

var (
	errInvalidFieldMask  = errors.DefineInvalidArgument("field_mask", "invalid field mask")
	errInvalidFieldValue = errors.DefineInvalidArgument("field_value", "invalid value of field `{field}`")
)

// Set implements ttnpb.AsEndDeviceRegistryServer.
func (r asEndDeviceRegistryServer) Set(ctx context.Context, req *ttnpb.SetEndDeviceRequest) (*ttnpb.EndDevice, error) {
	if ttnpb.HasAnyField(req.FieldMask.Paths, "session.dev_addr") && (req.EndDevice.Session == nil || req.EndDevice.Session.DevAddr.IsZero()) {
		return nil, errInvalidFieldValue.WithAttributes("field", "session.dev_addr")
	}

	if err := rights.RequireApplication(ctx, req.EndDevice.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE); err != nil {
		return nil, err
	}

	var evt events.Event
	dev, err := r.AS.deviceRegistry.Set(ctx, req.EndDevice.EndDeviceIdentifiers, req.FieldMask.Paths, func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
		sets := req.FieldMask.Paths
		if dev == nil {
			evt = evtCreateEndDevice(ctx, req.EndDevice.EndDeviceIdentifiers, nil)
		} else {
			evt = evtUpdateEndDevice(ctx, req.EndDevice.EndDeviceIdentifiers, req.FieldMask.Paths)
			if err := ttnpb.ProhibitFields(req.FieldMask.Paths, "ids.dev_addr"); err != nil {
				return nil, nil, errInvalidFieldMask.WithCause(err)
			}
			if ttnpb.HasAnyField(req.FieldMask.Paths, "session.dev_addr") {
				req.EndDevice.DevAddr = &req.EndDevice.Session.DevAddr
				sets = append(sets, "ids.dev_addr")
			}
			return &req.EndDevice, sets, nil
		}

		if req.EndDevice.DevAddr != nil && !req.EndDevice.DevAddr.IsZero() {
			return nil, nil, errInvalidFieldValue.WithAttributes("field", "ids.dev_addr")
		}
		sets = append(sets,
			"ids.application_ids",
			"ids.device_id",
		)
		if req.EndDevice.JoinEUI != nil && !req.EndDevice.JoinEUI.IsZero() {
			sets = append(sets,
				"ids.join_eui",
			)
		}
		if req.EndDevice.DevEUI != nil && !req.EndDevice.DevEUI.IsZero() {
			sets = append(sets,
				"ids.dev_eui",
			)
		}
		return &req.EndDevice, sets, nil
	})
	if err != nil {
		return nil, err
	}
	if evt != nil {
		events.Publish(evt)
	}
	return dev, nil
}

// Delete implements ttnpb.AsEndDeviceRegistryServer.
func (r asEndDeviceRegistryServer) Delete(ctx context.Context, ids *ttnpb.EndDeviceIdentifiers) (*pbtypes.Empty, error) {
	if err := rights.RequireApplication(ctx, ids.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE); err != nil {
		return nil, err
	}
	var evt events.Event
	_, err := r.AS.deviceRegistry.Set(ctx, *ids, nil, func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
		if dev != nil {
			evt = evtDeleteEndDevice(ctx, ids, nil)
		}
		return nil, nil, nil
	})
	if err != nil {
		return nil, err
	}
	if evt != nil {
		events.Publish(evt)
	}
	return ttnpb.Empty, nil
}
