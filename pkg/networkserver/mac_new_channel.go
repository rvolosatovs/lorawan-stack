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
	evtEnqueueNewChannelRequest = defineEnqueueMACRequestEvent("new_channel", "new channel")()
	evtReceiveNewChannelAccept  = defineReceiveMACAcceptEvent("new_channel", "new channel")()
	evtReceiveNewChannelReject  = defineReceiveMACRejectEvent("new_channel", "new channel")()
)

func enqueueNewChannelReq(ctx context.Context, dev *ttnpb.EndDevice, maxDownLen, maxUpLen uint16) (uint16, uint16, bool) {
	var ok bool
	dev.MACState.PendingRequests, maxDownLen, maxUpLen, ok = enqueueMACCommand(ttnpb.CID_NEW_CHANNEL, maxDownLen, maxUpLen, func(nDown, nUp uint16) ([]*ttnpb.MACCommand, uint16, bool) {
		var cmds []*ttnpb.MACCommand
		for i, ch := range dev.MACState.DesiredParameters.Channels {
			if i < len(dev.MACState.CurrentParameters.Channels) &&
				ch.UplinkFrequency == dev.MACState.CurrentParameters.Channels[i].UplinkFrequency &&
				ch.MinDataRateIndex == dev.MACState.CurrentParameters.Channels[i].MinDataRateIndex &&
				ch.MaxDataRateIndex == dev.MACState.CurrentParameters.Channels[i].MaxDataRateIndex {
				continue
			}
			if nDown < 1 || nUp < 1 {
				return cmds, uint16(len(cmds)), false
			}
			nDown--
			nUp--

			pld := &ttnpb.MACCommand_NewChannelReq{
				ChannelIndex:     uint32(i),
				Frequency:        ch.UplinkFrequency,
				MinDataRateIndex: ch.MinDataRateIndex,
				MaxDataRateIndex: ch.MaxDataRateIndex,
			}
			log.FromContext(ctx).WithFields(log.Fields(
				"index", pld.ChannelIndex,
				"frequency", pld.Frequency,
				"min_data_rate_index", pld.MinDataRateIndex,
				"max_data_rate_index", pld.MaxDataRateIndex,
			)).Debug("Enqueued NewChannelReq")
			cmds = append(cmds, pld.MACCommand())

			events.Publish(evtEnqueueNewChannelRequest(ctx, dev.EndDeviceIdentifiers, pld))
		}
		return cmds, uint16(len(cmds)), true
	}, dev.MACState.PendingRequests...)
	return maxDownLen, maxUpLen, ok
}

func handleNewChannelAns(ctx context.Context, dev *ttnpb.EndDevice, pld *ttnpb.MACCommand_NewChannelAns) (err error) {
	if pld == nil {
		return errNoPayload
	}

	if !pld.DataRateAck || !pld.FrequencyAck {
		events.Publish(evtReceiveNewChannelReject(ctx, dev.EndDeviceIdentifiers, pld))
	} else {
		events.Publish(evtReceiveNewChannelAccept(ctx, dev.EndDeviceIdentifiers, pld))
	}

	dev.MACState.PendingRequests, err = handleMACResponse(ttnpb.CID_NEW_CHANNEL, func(cmd *ttnpb.MACCommand) error {
		if !pld.DataRateAck || !pld.FrequencyAck {
			return nil
		}

		req := cmd.GetNewChannelReq()

		if uint(req.ChannelIndex) >= uint(len(dev.MACState.CurrentParameters.Channels)) {
			dev.MACState.CurrentParameters.Channels = append(dev.MACState.CurrentParameters.Channels, make([]*ttnpb.MACParameters_Channel, 1+int(req.ChannelIndex-uint32(len(dev.MACState.CurrentParameters.Channels))))...)
		}

		ch := dev.MACState.CurrentParameters.Channels[req.ChannelIndex]
		if ch == nil {
			ch = &ttnpb.MACParameters_Channel{
				DownlinkFrequency: req.Frequency,
			}
			dev.MACState.CurrentParameters.Channels[req.ChannelIndex] = ch
		}

		ch.UplinkFrequency = req.Frequency
		ch.MinDataRateIndex = req.MinDataRateIndex
		ch.MaxDataRateIndex = req.MaxDataRateIndex
		ch.EnableUplink = true
		return nil
	}, dev.MACState.PendingRequests...)
	return
}
