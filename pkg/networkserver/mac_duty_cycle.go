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
	evtMACDutyCycleRequest = events.Define("ns.mac.duty_cycle.request", "request duty cycle") // TODO(#988): publish when requesting
	evtMACDutyCycle        = events.Define("ns.mac.duty_cycle.accept", "device accepted duty cycle request")
)

func handleDutyCycleAns(ctx context.Context, dev *ttnpb.EndDevice) (err error) {
	dev.MACState.PendingRequests, err = handleMACResponse(ttnpb.CID_DUTY_CYCLE, func(cmd *ttnpb.MACCommand) {
		req := cmd.GetDutyCycleReq()

		dev.MACState.DutyCycle = req.MaxDutyCycle

		events.Publish(evtMACDutyCycle(ctx, dev.EndDeviceIdentifiers, req))
	}, dev.MACState.PendingRequests...)
	return
}
