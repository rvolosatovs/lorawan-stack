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
	evtEnqueueForceRejoinRequest = defineEnqueueMACRequestEvent("force_rejoin", "force rejoin")()
)

func enqueueForceRejoinReq(ctx context.Context, dev *ttnpb.EndDevice, maxDownLen, maxUpLen uint16) (uint16, uint16, bool) {
	// TODO: Generate ForceRejoinReq(https://github.com/TheThingsIndustries/ttn/issues/837)
	var pld *ttnpb.MACCommand_ForceRejoinReq
	if pld != nil {
		events.Publish(evtEnqueueForceRejoinRequest(ctx, dev.EndDeviceIdentifiers, pld))
	}
	return maxDownLen, maxUpLen, true
}
