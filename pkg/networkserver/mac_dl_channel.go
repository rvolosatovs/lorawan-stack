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
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	evtEnqueueDLChannelRequest = defineEnqueueMACRequestEvent("dl_channel", "downlink Rx1 channel frequency modification")()
	evtReceiveDLChannelAccept  = defineReceiveMACAcceptEvent("dl_channel", "downlink Rx1 channel frequency modification")()
	evtReceiveDLChannelReject  = defineReceiveMACRejectEvent("dl_channel", "downlink Rx1 channel frequency modification")()
)

func enqueueDLChannelReq(ctx context.Context, dev *ttnpb.EndDevice, maxDownLen, maxUpLen uint16) (uint16, uint16, bool) {
	var ok bool
	dev.MACState.PendingRequests, maxDownLen, maxUpLen, ok = enqueueMACCommand(ttnpb.CID_DL_CHANNEL, maxDownLen, maxUpLen, func(nDown, nUp uint16) ([]*ttnpb.MACCommand, uint16, bool) {
		var cmds []*ttnpb.MACCommand
		for i := 0; i < len(dev.MACState.DesiredParameters.Channels) && i < len(dev.MACState.CurrentParameters.Channels); i++ {
			if dev.MACState.DesiredParameters.Channels[i].DownlinkFrequency == dev.MACState.CurrentParameters.Channels[i].DownlinkFrequency {
				continue
			}
			if nDown < 1 || nUp < 1 {
				return cmds, uint16(len(cmds)), false
			}
			nDown--
			nUp--

			pld := &ttnpb.MACCommand_DLChannelReq{
				ChannelIndex: uint32(i),
				Frequency:    dev.MACState.DesiredParameters.Channels[i].DownlinkFrequency,
			}
			cmds = append(cmds, pld.MACCommand())

			events.Publish(evtEnqueueDLChannelRequest(ctx, dev.EndDeviceIdentifiers, pld))
		}
		return cmds, uint16(len(cmds)), true

	}, dev.MACState.PendingRequests...)
	return maxDownLen, maxUpLen, ok
}

func handleDLChannelAns(ctx context.Context, dev *ttnpb.EndDevice, pld *ttnpb.MACCommand_DLChannelAns) (err error) {
	if pld == nil {
		return errNoPayload
	}

	if !pld.ChannelIndexAck || !pld.FrequencyAck {
		events.Publish(evtReceiveDLChannelReject(ctx, dev.EndDeviceIdentifiers, pld))
	} else {
		events.Publish(evtReceiveDLChannelAccept(ctx, dev.EndDeviceIdentifiers, pld))
	}

	dev.MACState.PendingRequests, err = handleMACResponse(ttnpb.CID_DL_CHANNEL, func(cmd *ttnpb.MACCommand) error {
		if !pld.ChannelIndexAck || !pld.FrequencyAck {
			return nil
		}

		req := cmd.GetDLChannelReq()
		if uint(req.ChannelIndex) >= uint(len(dev.MACState.CurrentParameters.Channels)) || dev.MACState.CurrentParameters.Channels[req.ChannelIndex] == nil {
			return errCorruptedMACState.WithCause(errUnknownChannel)
		}
		dev.MACState.CurrentParameters.Channels[req.ChannelIndex].DownlinkFrequency = req.Frequency
		return nil

	}, dev.MACState.PendingRequests...)
	return
}
