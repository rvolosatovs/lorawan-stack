// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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
	evtEnqueueTxParamSetupRequest = defineEnqueueMACRequestEvent("tx_param_setup", "Tx parameter setup")
	evtReceiveTxParamSetupAnswer  = defineReceiveMACAnswerEvent("tx_param_setup", "Tx parameter setup")
)

func enqueueTxParamSetupReq(ctx context.Context, dev *ttnpb.EndDevice, maxDownLen, maxUpLen uint16) (uint16, uint16, bool) {
	if dev.MACState.DesiredParameters.MaxEIRP == dev.MACState.CurrentParameters.MaxEIRP &&
		dev.MACState.DesiredParameters.DownlinkDwellTime == dev.MACState.CurrentParameters.DownlinkDwellTime &&
		dev.MACState.DesiredParameters.UplinkDwellTime == dev.MACState.CurrentParameters.UplinkDwellTime {
		return maxDownLen, maxUpLen, true
	}

	var ok bool
	dev.MACState.PendingRequests, maxDownLen, maxUpLen, ok = enqueueMACCommand(ttnpb.CID_TX_PARAM_SETUP, maxDownLen, maxUpLen, func(nDown, nUp uint16) ([]*ttnpb.MACCommand, uint16, bool) {
		if nDown < 1 || nUp < 1 {
			return nil, 0, false
		}
		pld := &ttnpb.MACCommand_TxParamSetupReq{
			MaxEIRPIndex:      ttnpb.Float32ToDeviceEIRP(dev.MACState.DesiredParameters.MaxEIRP),
			DownlinkDwellTime: dev.MACState.DesiredParameters.DownlinkDwellTime,
			UplinkDwellTime:   dev.MACState.DesiredParameters.UplinkDwellTime,
		}
		events.Publish(evtEnqueueTxParamSetupRequest(ctx, dev.EndDeviceIdentifiers, pld))
		return []*ttnpb.MACCommand{pld.MACCommand()}, 1, true

	}, dev.MACState.PendingRequests...)
	return maxDownLen, maxUpLen, ok
}

func handleTxParamSetupAns(ctx context.Context, dev *ttnpb.EndDevice) (err error) {
	dev.MACState.PendingRequests, err = handleMACResponse(ttnpb.CID_TX_PARAM_SETUP, func(cmd *ttnpb.MACCommand) error {
		events.Publish(evtReceiveTxParamSetupAnswer(ctx, dev.EndDeviceIdentifiers, nil))

		req := cmd.GetTxParamSetupReq()

		dev.MACState.CurrentParameters.DownlinkDwellTime = req.DownlinkDwellTime
		dev.MACState.CurrentParameters.UplinkDwellTime = req.UplinkDwellTime
		dev.MACState.CurrentParameters.MaxEIRP = ttnpb.DeviceEIRPToFloat32(req.MaxEIRPIndex)
		if ttnpb.Float32ToDeviceEIRP(dev.MACState.DesiredParameters.MaxEIRP) == req.MaxEIRPIndex {
			dev.MACState.DesiredParameters.MaxEIRP = dev.MACState.CurrentParameters.MaxEIRP
		}
		return nil

	}, dev.MACState.PendingRequests...)
	return
}
