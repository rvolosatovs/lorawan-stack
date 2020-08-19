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
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	evtEnqueuePingSlotInfoAnswer = defineEnqueueMACAnswerEvent(
		"ping_slot_info", "ping slot info",
	)()
	evtReceivePingSlotInfoRequest = defineReceiveMACRequestEvent(
		"ping_slot_info", "ping slot info",
		events.WithDataType(&ttnpb.MACCommand_PingSlotInfoReq{}),
	)()
)

func handlePingSlotInfoReq(ctx context.Context, dev *ttnpb.EndDevice, pld *ttnpb.MACCommand_PingSlotInfoReq) (events.Builders, error) {
	if pld == nil {
		return nil, errNoPayload.New()
	}

	evs := events.Builders{
		evtReceivePingSlotInfoRequest.With(events.WithData(pld)),
	}
	if dev.MACState.DeviceClass != ttnpb.CLASS_A {
		log.FromContext(ctx).Debug("Ignore PingSlotInfoReq from device not in class A mode")
		return evs, nil
	}

	dev.MACState.PingSlotPeriodicity = &ttnpb.PingSlotPeriodValue{Value: pld.Period}
	dev.MACState.QueuedResponses = append(dev.MACState.QueuedResponses, ttnpb.CID_PING_SLOT_INFO.MACCommand())
	return append(evs,
		evtEnqueuePingSlotInfoAnswer,
	), nil
}
