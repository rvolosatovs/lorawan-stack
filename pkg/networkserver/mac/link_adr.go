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

package mac

import (
	"context"
	"math"

	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	EvtEnqueueLinkADRRequest = defineEnqueueMACRequestEvent(
		"link_adr", "link ADR",
		events.WithDataType(&ttnpb.MACCommand_LinkADRReq{}),
	)()
	EvtReceiveLinkADRAccept = defineReceiveMACAcceptEvent(
		"link_adr", "link ADR",
		events.WithDataType(&ttnpb.MACCommand_LinkADRAns{}),
	)()
	EvtReceiveLinkADRReject = defineReceiveMACRejectEvent(
		"link_adr", "link ADR",
		events.WithDataType(&ttnpb.MACCommand_LinkADRAns{}),
	)()
)

const (
	noChangeDataRateIndex = ttnpb.DATA_RATE_15
	noChangeTXPowerIndex  = 15
)

type linkADRReqParameters struct {
	Masks         []band.ChMaskCntlPair
	DataRateIndex ttnpb.DataRateIndex
	TxPowerIndex  uint32
	NbTrans       uint32
}

func generateLinkADRReq(ctx context.Context, dev *ttnpb.EndDevice, defaults ttnpb.MACSettings, phy *band.Band) (linkADRReqParameters, bool, error) {
	if dev.GetMulticast() || dev.GetMACState() == nil {
		return linkADRReqParameters{}, false, nil
	}
	if len(dev.MACState.DesiredParameters.Channels) > int(phy.MaxUplinkChannels) {
		return linkADRReqParameters{}, false, ErrCorruptedMACState.New()
	}

	currentChs := make([]bool, phy.MaxUplinkChannels)
	for i, ch := range dev.MACState.CurrentParameters.Channels {
		currentChs[i] = ch.GetEnableUplink()
	}
	desiredChs := make([]bool, phy.MaxUplinkChannels)
	for i, ch := range dev.MACState.DesiredParameters.Channels {
		isEnabled := ch.GetEnableUplink()
		if isEnabled && ch.UplinkFrequency == 0 {
			return linkADRReqParameters{}, false, ErrCorruptedMACState.New()
		}
		if DeviceNeedsNewChannelReqAtIndex(dev, i) {
			currentChs[i] = ch != nil && ch.UplinkFrequency != 0
		}
		desiredChs[i] = isEnabled
	}

	useADR := DeviceUseADR(dev, defaults, phy)
	switch {
	case !band.EqualChMasks(currentChs, desiredChs):
		// NOTE: LinkADRReq is scheduled regardless of ADR settings if channel mask is required, which often is the case with ABP devices or when ChMask CFList is not supported/used.
	case !useADR:
		return linkADRReqParameters{}, false, nil
	case dev.MACState.DesiredParameters.ADRNbTrans != dev.MACState.CurrentParameters.ADRNbTrans,
		dev.MACState.DesiredParameters.ADRDataRateIndex != dev.MACState.CurrentParameters.ADRDataRateIndex,
		dev.MACState.DesiredParameters.ADRTxPowerIndex != dev.MACState.CurrentParameters.ADRTxPowerIndex:
	default:
		return linkADRReqParameters{}, false, nil
	}
	desiredMasks, err := phy.GenerateChMasks(currentChs, desiredChs)
	if err != nil {
		return linkADRReqParameters{}, false, err
	}
	if len(desiredMasks) > math.MaxUint16 {
		// Something is really wrong.
		return linkADRReqParameters{}, false, ErrCorruptedMACState.New()
	}

	var (
		drIdx      ttnpb.DataRateIndex
		txPowerIdx uint32
		nbTrans    uint32
	)
	if !useADR {
		drIdx = dev.MACState.CurrentParameters.ADRDataRateIndex
		txPowerIdx = dev.MACState.CurrentParameters.ADRTxPowerIndex
		nbTrans = dev.MACState.CurrentParameters.ADRNbTrans
	} else {
		minDataRateIndex, maxDataRateIndex, ok := channelDataRateRange(dev.MACState.DesiredParameters.Channels...)
		if !ok ||
			dev.MACState.DesiredParameters.ADRTxPowerIndex > uint32(phy.MaxTxPowerIndex()) ||
			dev.MACState.DesiredParameters.ADRDataRateIndex > phy.MaxADRDataRateIndex {
			return linkADRReqParameters{}, false, ErrCorruptedMACState.New()
		}
		if dev.MACState.CurrentParameters.ADRDataRateIndex > minDataRateIndex {
			minDataRateIndex = dev.MACState.CurrentParameters.ADRDataRateIndex
		}
		if dev.MACState.DesiredParameters.ADRDataRateIndex < minDataRateIndex || dev.MACState.DesiredParameters.ADRDataRateIndex > maxDataRateIndex {
			return linkADRReqParameters{}, false, ErrCorruptedMACState.New()
		}

		drIdx = dev.MACState.DesiredParameters.ADRDataRateIndex
		txPowerIdx = dev.MACState.DesiredParameters.ADRTxPowerIndex
		nbTrans = dev.MACState.DesiredParameters.ADRNbTrans
		switch {
		case !deviceRejectedADRDataRateIndex(dev, drIdx) && !deviceRejectedADRTXPowerIndex(dev, txPowerIdx):
			// Only send the desired DataRateIndex and TXPowerIndex if neither of them were rejected.

		case len(desiredMasks) == 0 && dev.MACState.DesiredParameters.ADRNbTrans == dev.MACState.CurrentParameters.ADRNbTrans:
			log.FromContext(ctx).Debug("Either desired data rate index or TX power output index have been rejected and there are no channel mask and NbTrans changes desired, avoid enqueueing LinkADRReq")
			return linkADRReqParameters{}, false, nil

		case dev.MACState.LoRaWANVersion.HasNoChangeDataRateIndex() && !deviceRejectedADRDataRateIndex(dev, noChangeDataRateIndex) &&
			dev.MACState.LoRaWANVersion.HasNoChangeTXPowerIndex() && !deviceRejectedADRTXPowerIndex(dev, noChangeTXPowerIndex):
			drIdx = noChangeDataRateIndex
			txPowerIdx = noChangeTXPowerIndex

		default:
			for deviceRejectedADRDataRateIndex(dev, drIdx) || deviceRejectedADRTXPowerIndex(dev, txPowerIdx) {
				// Since either data rate or TX power index (or both) were rejected by the device, undo the
				// desired ADR adjustments step-by-step until possibly fitting index pair is found.
				if drIdx == minDataRateIndex && (deviceRejectedADRDataRateIndex(dev, drIdx) || txPowerIdx == 0) {
					log.FromContext(ctx).WithFields(log.Fields(
						"current_adr_nb_trans", dev.MACState.CurrentParameters.ADRNbTrans,
						"desired_adr_nb_trans", dev.MACState.DesiredParameters.ADRNbTrans,
						"desired_mask_count", len(desiredMasks),
					)).Warn("Device rejected either all available data rate indexes or all available TX power output indexes and there are channel mask or NbTrans changes desired, avoid enqueueing LinkADRReq")
					return linkADRReqParameters{}, false, nil
				}
				for drIdx > minDataRateIndex && (deviceRejectedADRDataRateIndex(dev, drIdx) || txPowerIdx == 0 && deviceRejectedADRTXPowerIndex(dev, txPowerIdx)) {
					// Increase data rate until a non-rejected index is found.
					// Set TX power to maximum possible value.
					drIdx--
					txPowerIdx = uint32(phy.MaxTxPowerIndex())
				}
				for txPowerIdx > 0 && deviceRejectedADRTXPowerIndex(dev, txPowerIdx) {
					// Increase TX output power until a non-rejected index is found.
					txPowerIdx--
				}
			}
		}
	}
	if drIdx == dev.MACState.CurrentParameters.ADRDataRateIndex && dev.MACState.LoRaWANVersion.HasNoChangeDataRateIndex() && !deviceRejectedADRDataRateIndex(dev, noChangeDataRateIndex) {
		drIdx = noChangeDataRateIndex
	}
	if txPowerIdx == dev.MACState.CurrentParameters.ADRTxPowerIndex && dev.MACState.LoRaWANVersion.HasNoChangeTXPowerIndex() && !deviceRejectedADRTXPowerIndex(dev, noChangeTXPowerIndex) {
		txPowerIdx = noChangeTXPowerIndex
	}
	return linkADRReqParameters{
		Masks:         desiredMasks,
		DataRateIndex: drIdx,
		TxPowerIndex:  txPowerIdx,
		NbTrans:       nbTrans,
	}, true, nil
}

func DeviceNeedsLinkADRReq(ctx context.Context, dev *ttnpb.EndDevice, defaults ttnpb.MACSettings, phy *band.Band) bool {
	_, required, err := generateLinkADRReq(ctx, dev, defaults, phy)
	return err == nil && required
}

func EnqueueLinkADRReq(ctx context.Context, dev *ttnpb.EndDevice, maxDownLen, maxUpLen uint16, defaults ttnpb.MACSettings, phy *band.Band) (EnqueueState, error) {
	params, required, err := generateLinkADRReq(ctx, dev, defaults, phy)
	if err != nil {
		return EnqueueState{
			MaxDownLen: maxDownLen,
			MaxUpLen:   maxUpLen,
		}, err
	}
	if !required {
		return EnqueueState{
			MaxDownLen: maxDownLen,
			MaxUpLen:   maxUpLen,
			Ok:         true,
		}, nil
	}

	var st EnqueueState
	dev.MACState.PendingRequests, st = enqueueMACCommand(ttnpb.CID_LINK_ADR, maxDownLen, maxUpLen, func(nDown, nUp uint16) ([]*ttnpb.MACCommand, uint16, events.Builders, bool) {
		if int(nDown) < len(params.Masks) {
			return nil, 0, nil, false
		}

		uplinksNeeded := uint16(1)
		if dev.MACState.LoRaWANVersion.Compare(ttnpb.MAC_V1_1) < 0 {
			uplinksNeeded = uint16(len(params.Masks))
		}
		if nUp < uplinksNeeded {
			return nil, 0, nil, false
		}
		evs := make(events.Builders, 0, len(params.Masks))
		cmds := make([]*ttnpb.MACCommand, 0, len(params.Masks))
		for i, m := range params.Masks {
			req := &ttnpb.MACCommand_LinkADRReq{
				DataRateIndex:      params.DataRateIndex,
				TxPowerIndex:       params.TxPowerIndex,
				NbTrans:            params.NbTrans,
				ChannelMaskControl: uint32(m.Cntl),
				ChannelMask:        params.Masks[i].Mask[:],
			}
			cmds = append(cmds, req.MACCommand())
			evs = append(evs, EvtEnqueueLinkADRRequest.With(events.WithData(req)))
			log.FromContext(ctx).WithFields(log.Fields(
				"data_rate_index", req.DataRateIndex,
				"nb_trans", req.NbTrans,
				"tx_power_index", req.TxPowerIndex,
				"channel_mask_control", req.ChannelMaskControl,
				"channel_mask", req.ChannelMask,
			)).Debug("Enqueued LinkADRReq")
		}
		return cmds, uplinksNeeded, evs, true
	}, dev.MACState.PendingRequests...)
	return st, nil
}

func HandleLinkADRAns(ctx context.Context, dev *ttnpb.EndDevice, pld *ttnpb.MACCommand_LinkADRAns, dupCount uint, fps *frequencyplans.Store) (events.Builders, error) {
	if pld == nil {
		return nil, ErrNoPayload.New()
	}
	if (dev.MACState.LoRaWANVersion.Compare(ttnpb.MAC_V1_0_2) < 0 || dev.MACState.LoRaWANVersion.Compare(ttnpb.MAC_V1_1) >= 0) && dupCount != 0 {
		return nil, ErrInvalidPayload.New()
	}

	Evt := EvtReceiveLinkADRAccept
	if !pld.ChannelMaskAck || !pld.DataRateIndexAck || !pld.TxPowerIndexAck {
		Evt = EvtReceiveLinkADRReject

		// See "Table 6: LinkADRAns status bits signification" of LoRaWAN 1.1 specification
		if !pld.ChannelMaskAck {
			log.FromContext(ctx).Warn("Either Network Server sent a channel mask, which enables a yet undefined channel or requires all channels to be disabled, or device is malfunctioning.")
		}
	}
	evs := events.Builders{Evt.With(events.WithData(pld))}

	phy, err := DeviceBand(dev, fps)
	if err != nil {
		return evs, err
	}

	handler := handleMACResponseBlock
	if dev.MACState.LoRaWANVersion.Compare(ttnpb.MAC_V1_0_2) < 0 {
		handler = handleMACResponse
	}
	var n uint
	var req *ttnpb.MACCommand_LinkADRReq
	dev.MACState.PendingRequests, err = handler(ttnpb.CID_LINK_ADR, func(cmd *ttnpb.MACCommand) error {
		if dev.MACState.LoRaWANVersion.Compare(ttnpb.MAC_V1_0_2) >= 0 && dev.MACState.LoRaWANVersion.Compare(ttnpb.MAC_V1_1) < 0 && n > dupCount+1 {
			return ErrInvalidPayload.New()
		}
		n++

		req = cmd.GetLinkADRReq()
		if req.NbTrans > 15 || len(req.ChannelMask) != 16 || req.ChannelMaskControl > 7 {
			panic("Network Server scheduled an invalid LinkADR command")
		}
		if !pld.ChannelMaskAck || !pld.DataRateIndexAck || !pld.TxPowerIndexAck {
			return nil
		}
		var mask [16]bool
		for i, v := range req.ChannelMask {
			mask[i] = v
		}
		m, err := phy.ParseChMask(mask, uint8(req.ChannelMaskControl))
		if err != nil {
			return err
		}
		for i, masked := range m {
			if int(i) >= len(dev.MACState.CurrentParameters.Channels) || dev.MACState.CurrentParameters.Channels[i] == nil {
				if !masked {
					continue
				}
				return ErrCorruptedMACState.WithCause(ErrUnknownChannel)
			}
			dev.MACState.CurrentParameters.Channels[i].EnableUplink = masked
		}
		return nil
	}, dev.MACState.PendingRequests...)
	if err != nil || req == nil {
		return evs, err
	}

	if !pld.DataRateIndexAck {
		if i := searchDataRateIndex(req.DataRateIndex, dev.MACState.RejectedADRDataRateIndexes...); i == len(dev.MACState.RejectedADRDataRateIndexes) || dev.MACState.RejectedADRDataRateIndexes[i] != req.DataRateIndex {
			dev.MACState.RejectedADRDataRateIndexes = append(dev.MACState.RejectedADRDataRateIndexes, ttnpb.DATA_RATE_0)
			copy(dev.MACState.RejectedADRDataRateIndexes[i+1:], dev.MACState.RejectedADRDataRateIndexes[i:])
			dev.MACState.RejectedADRDataRateIndexes[i] = req.DataRateIndex
		}
	}
	if !pld.TxPowerIndexAck {
		if i := searchUint32(req.TxPowerIndex, dev.MACState.RejectedADRTxPowerIndexes...); i == len(dev.MACState.RejectedADRTxPowerIndexes) || dev.MACState.RejectedADRTxPowerIndexes[i] != req.TxPowerIndex {
			dev.MACState.RejectedADRTxPowerIndexes = append(dev.MACState.RejectedADRTxPowerIndexes, 0)
			copy(dev.MACState.RejectedADRTxPowerIndexes[i+1:], dev.MACState.RejectedADRTxPowerIndexes[i:])
			dev.MACState.RejectedADRTxPowerIndexes[i] = req.TxPowerIndex
		}
	}
	if !pld.ChannelMaskAck || !pld.DataRateIndexAck || !pld.TxPowerIndexAck {
		return evs, nil
	}
	if !dev.MACState.LoRaWANVersion.HasNoChangeDataRateIndex() || req.DataRateIndex != noChangeDataRateIndex {
		dev.MACState.CurrentParameters.ADRDataRateIndex = req.DataRateIndex
		dev.RecentADRUplinks = nil
	}
	if !dev.MACState.LoRaWANVersion.HasNoChangeTXPowerIndex() || req.TxPowerIndex != noChangeTXPowerIndex {
		dev.MACState.CurrentParameters.ADRTxPowerIndex = req.TxPowerIndex
		dev.RecentADRUplinks = nil
	}
	if req.NbTrans > 0 && dev.MACState.CurrentParameters.ADRNbTrans != req.NbTrans {
		dev.MACState.CurrentParameters.ADRNbTrans = req.NbTrans
		dev.RecentADRUplinks = nil
	}
	return evs, nil
}
