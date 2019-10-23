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

	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	evtEnqueuePingSlotChannelRequest = defineEnqueueMACRequestEvent("ping_slot_channel", "ping slot channel")()
	evtReceivePingSlotChannelAnswer  = defineReceiveMACAcceptEvent("ping_slot_channel", "ping slot channel")()
)

func deviceNeedsPingSlotChannelReq(dev *ttnpb.EndDevice) bool {
	return dev.MACState != nil &&
		(dev.MACState.DesiredParameters.PingSlotDataRateIndex != dev.MACState.CurrentParameters.PingSlotDataRateIndex ||
			dev.MACState.DesiredParameters.PingSlotFrequency != dev.MACState.CurrentParameters.PingSlotFrequency)
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
	dev.MACState.PendingRequests, st = enqueueMACCommand(ttnpb.CID_PING_SLOT_CHANNEL, maxDownLen, maxUpLen, func(nDown, nUp uint16) ([]*ttnpb.MACCommand, uint16, []events.DefinitionDataClosure, bool) {
		if nDown < 1 || nUp < 1 {
			return nil, 0, nil, false
		}
		req := &ttnpb.MACCommand_PingSlotChannelReq{
			Frequency:     dev.MACState.DesiredParameters.PingSlotFrequency,
			DataRateIndex: dev.MACState.DesiredParameters.PingSlotDataRateIndex,
		}
		log.FromContext(ctx).WithFields(log.Fields(
			"frequency", req.Frequency,
			"data_rate_index", req.DataRateIndex,
		)).Debug("Enqueued PingSlotChannelReq")
		return []*ttnpb.MACCommand{
				req.MACCommand(),
			},
			1,
			[]events.DefinitionDataClosure{
				evtEnqueuePingSlotChannelRequest.BindData(req),
			},
			true
	}, dev.MACState.PendingRequests...)
	return st
}

func handlePingSlotChannelAns(ctx context.Context, dev *ttnpb.EndDevice, pld *ttnpb.MACCommand_PingSlotChannelAns) ([]events.DefinitionDataClosure, error) {
	if pld == nil {
		return nil, errNoPayload
	}

	var err error
	dev.MACState.PendingRequests, err = handleMACResponse(ttnpb.CID_PING_SLOT_CHANNEL, func(cmd *ttnpb.MACCommand) error {
		req := cmd.GetPingSlotChannelReq()

		dev.MACState.CurrentParameters.PingSlotDataRateIndex = req.DataRateIndex
		dev.MACState.CurrentParameters.PingSlotFrequency = req.Frequency
		return nil
	}, dev.MACState.PendingRequests...)
	return []events.DefinitionDataClosure{
		evtReceivePingSlotChannelAnswer.BindData(pld),
	}, err
}
