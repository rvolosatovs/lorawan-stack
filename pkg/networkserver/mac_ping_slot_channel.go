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
	evtEnqueuePingSlotChannelRequest = defineEnqueueMACRequestEvent("ping_slot_channel", "ping slot channel")()
	evtReceivePingSlotChannelAnswer  = defineReceiveMACAcceptEvent("ping_slot_channel", "ping slot channel")()
)

func deviceNeedsPingSlotChannelReq(dev *ttnpb.EndDevice) bool {
	switch {
	case dev.GetMulticast(),
		dev.GetMACState() == nil:
		return false
	case dev.MACState.DesiredParameters.PingSlotFrequency != dev.MACState.CurrentParameters.PingSlotFrequency:
		return true
	case dev.MACState.DesiredParameters.PingSlotDataRateIndexValue == nil:
		return false
	case dev.MACState.CurrentParameters.PingSlotDataRateIndexValue == nil,
		dev.MACState.DesiredParameters.PingSlotDataRateIndexValue.Value != dev.MACState.CurrentParameters.PingSlotDataRateIndexValue.Value:
		return true
	}
	return false
}

func enqueuePingSlotChannelReq(ctx context.Context, dev *ttnpb.EndDevice, maxDownLen, maxUpLen uint16) macCommandEnqueueState {
	if !deviceNeedsPingSlotChannelReq(dev) {
		return macCommandEnqueueState{
			MaxDownLen: maxDownLen,
			MaxUpLen:   maxUpLen,
			Ok:         true,
		}
	}

	var st macCommandEnqueueState
	dev.MACState.PendingRequests, st = enqueueMACCommand(ttnpb.CID_PING_SLOT_CHANNEL, maxDownLen, maxUpLen, func(nDown, nUp uint16) ([]*ttnpb.MACCommand, uint16, events.Builders, bool) {
		if nDown < 1 || nUp < 1 {
			return nil, 0, nil, false
		}
		req := &ttnpb.MACCommand_PingSlotChannelReq{
			Frequency:     dev.MACState.DesiredParameters.PingSlotFrequency,
			DataRateIndex: dev.MACState.DesiredParameters.PingSlotDataRateIndexValue.Value,
		}
		log.FromContext(ctx).WithFields(log.Fields(
			"frequency", req.Frequency,
			"data_rate_index", req.DataRateIndex,
		)).Debug("Enqueued PingSlotChannelReq")
		return []*ttnpb.MACCommand{
				req.MACCommand(),
			},
			1,
			events.Builders{
				evtEnqueuePingSlotChannelRequest.With(events.WithData(req)),
			},
			true
	}, dev.MACState.PendingRequests...)
	return st
}

func handlePingSlotChannelAns(ctx context.Context, dev *ttnpb.EndDevice, pld *ttnpb.MACCommand_PingSlotChannelAns) (events.Builders, error) {
	if pld == nil {
		return nil, errNoPayload.New()
	}

	var err error
	dev.MACState.PendingRequests, err = handleMACResponse(ttnpb.CID_PING_SLOT_CHANNEL, func(cmd *ttnpb.MACCommand) error {
		req := cmd.GetPingSlotChannelReq()

		dev.MACState.CurrentParameters.PingSlotFrequency = req.Frequency
		dev.MACState.CurrentParameters.PingSlotDataRateIndexValue = &ttnpb.DataRateIndexValue{Value: req.DataRateIndex}
		return nil
	}, dev.MACState.PendingRequests...)
	return events.Builders{
		evtReceivePingSlotChannelAnswer.With(events.WithData(pld)),
	}, err
}
