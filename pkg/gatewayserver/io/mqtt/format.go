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

package mqtt

import (
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io/mqtt/topics"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// Format formats topics, downlink, uplink and status messages.
type Format interface {
	topics.Layout

	FromDownlink(down *ttnpb.DownlinkMessage, ids ttnpb.GatewayIdentifiers) ([]byte, error)
	ToUplink(message []byte, ids ttnpb.GatewayIdentifiers) (*ttnpb.UplinkMessage, error)
	ToStatus(message []byte, ids ttnpb.GatewayIdentifiers) (*ttnpb.GatewayStatus, error)
	ToTxAck(message []byte, ids ttnpb.GatewayIdentifiers) (*ttnpb.TxAcknowledgment, error)
}

var errNotSupported = errors.DefineFailedPrecondition("not_supported", "not supported")
