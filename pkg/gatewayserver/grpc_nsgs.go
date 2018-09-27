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

package gatewayserver

import (
	"context"

	"github.com/gogo/protobuf/types"
	clusterauth "go.thethings.network/lorawan-stack/pkg/auth/cluster"
	errors "go.thethings.network/lorawan-stack/pkg/errorsv3"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"go.thethings.network/lorawan-stack/pkg/validate"
)

var errInvalidAntennaIndex = errors.DefineNotFound("antenna_not_found", "antenna `{index}` not found")

// ScheduleDownlink instructs the Gateway Server to schedule a downlink message.
// The Gateway Server may refuse if there are any conflicts in the schedule or
// if a duty cycle prevents the gateway from transmitting.
func (gs *GatewayServer) ScheduleDownlink(ctx context.Context, down *ttnpb.DownlinkMessage) (*types.Empty, error) {
	if err := clusterauth.Authorized(ctx); err != nil {
		return nil, err
	}

	ids := down.TxMetadata.GatewayIdentifiers
	// TODO: Remove validation (https://github.com/TheThingsIndustries/lorawan-stack/issues/1058)
	if err := validate.ID(ids.GatewayID); err != nil {
		return nil, err
	}

	uid := unique.ID(ctx, ids)
	val, ok := gs.connections.Load(uid)
	if !ok {
		return nil, errNotConnected.WithAttributes("gateway_uid", uid)
	}
	conn := val.(*io.Connection)
	gtw := conn.Gateway()
	if len(gtw.Antennas) <= int(down.TxMetadata.AntennaIndex) {
		return nil, errInvalidAntennaIndex.WithAttributes("index", down.TxMetadata.AntennaIndex)
	}
	down.Settings.TxPower -= int32(gtw.Antennas[down.TxMetadata.AntennaIndex].Gain)

	if err := conn.SendDown(down); err != nil {
		return nil, err
	}

	ctx = events.ContextWithCorrelationID(ctx, events.CorrelationIDsFromContext(conn.Context())...)
	registerSendDownlink(ctx, conn.Gateway(), down)
	return &types.Empty{}, nil
}
