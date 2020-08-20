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

	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	evtReceiveDeviceModeIndication = defineReceiveMACIndicationEvent(
		"device_mode", "device mode",
		events.WithDataType(&ttnpb.MACCommand_DeviceModeInd{}),
	)()
	evtEnqueueDeviceModeConfirmation = defineEnqueueMACConfirmationEvent(
		"device_mode", "device mode",
		events.WithDataType(&ttnpb.MACCommand_DeviceModeConf{}),
	)()
)

func handleDeviceModeInd(ctx context.Context, dev *ttnpb.EndDevice, pld *ttnpb.MACCommand_DeviceModeInd) (events.Builders, error) {
	if pld == nil {
		return nil, errNoPayload.New()
	}

	evs := events.Builders{
		evtReceiveDeviceModeIndication.With(events.WithData(pld)),
	}
	switch {
	case pld.Class == ttnpb.CLASS_C && dev.SupportsClassC && dev.MACState.DeviceClass != ttnpb.CLASS_C:
		evs = append(evs, evtClassCSwitch.With(events.WithData(dev.MACState.DeviceClass)))
		dev.MACState.DeviceClass = ttnpb.CLASS_C

	case pld.Class == ttnpb.CLASS_A && dev.MACState.DeviceClass != ttnpb.CLASS_A:
		evs = append(evs, evtClassASwitch.With(events.WithData(dev.MACState.DeviceClass)))
		dev.MACState.DeviceClass = ttnpb.CLASS_A
	}
	conf := &ttnpb.MACCommand_DeviceModeConf{
		Class: dev.MACState.DeviceClass,
	}
	dev.MACState.QueuedResponses = append(dev.MACState.QueuedResponses, conf.MACCommand())
	return append(evs,
		evtEnqueueDeviceModeConfirmation.With(events.WithData(conf)),
	), nil
}
