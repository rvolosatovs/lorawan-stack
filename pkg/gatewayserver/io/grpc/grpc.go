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

package grpc

import (
	"context"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"go.thethings.network/lorawan-stack/pkg/validate"
	"google.golang.org/grpc/peer"
)

type impl struct {
	server io.Server
}

// New returns a new gRPC frontend.
func New(server io.Server) ttnpb.GtwGsServer {
	return &impl{server}
}

var errConnect = errors.Define("connect", "failed to connect gateway `{gateway_uid}`")

// LinkGateway links the gateway to the Gateway Server.
func (s *impl) LinkGateway(link ttnpb.GtwGs_LinkGatewayServer) (err error) {
	ctx := log.NewContextWithField(link.Context(), "namespace", "gatewayserver/io/grpc")

	ids := ttnpb.GatewayIdentifiers{
		GatewayID: rpcmetadata.FromIncomingContext(ctx).ID,
	}
	if err = validate.ID(ids.GatewayID); err != nil {
		return
	}
	if err = rights.RequireGateway(ctx, ids, ttnpb.RIGHT_GATEWAY_LINK); err != nil {
		return
	}
	ctx, ids, err = s.server.FillGatewayContext(ctx, ids)
	if err != nil {
		return
	}

	if peer, ok := peer.FromContext(ctx); ok {
		ctx = log.NewContextWithField(ctx, "remote_addr", peer.Addr.String())
	}
	uid := unique.ID(ctx, ids)
	ctx = log.NewContextWithField(ctx, "gateway_uid", uid)
	logger := log.FromContext(ctx)
	conn, err := s.server.Connect(ctx, "grpc", ids)
	if err != nil {
		logger.WithError(err).Warn("Failed to connect")
		return errConnect.WithCause(err).WithAttributes("gateway_uid", uid)
	}
	logger.Info("Connected")
	defer logger.Info("Disconnected")
	if err = s.server.ClaimDownlink(ctx, ids); err != nil {
		logger.WithError(err).Error("Failed to claim downlink")
		return
	}

	go func() {
		for {
			select {
			case <-conn.Context().Done():
				return
			case down := <-conn.Down():
				msg := &ttnpb.GatewayDown{
					DownlinkMessage: down,
				}
				logger.Info("Sending downlink message")
				if err := link.Send(msg); err != nil {
					logger.WithError(err).Warn("Failed to send message")
					conn.Disconnect(err)
					return
				}
			}
		}
	}()

	for {
		msg, err := link.Recv()
		if err != nil {
			if !errors.IsCanceled(err) {
				logger.WithError(err).Warn("Link failed")
			}
			conn.Disconnect(err)
			return err
		}
		now := time.Now()

		logger.WithFields(log.Fields(
			"has_status", msg.GatewayStatus != nil,
			"uplink_count", len(msg.UplinkMessages),
		)).Debug("Received message")

		for _, up := range msg.UplinkMessages {
			up.ReceivedAt = now
			if err := conn.HandleUp(up); err != nil {
				logger.WithError(err).Warn("Failed to handle uplink message")
			}
		}
		if msg.GatewayStatus != nil {
			if err := conn.HandleStatus(msg.GatewayStatus); err != nil {
				logger.WithError(err).Warn("Failed to handle status message")
			}
		}
		if msg.TxAcknowledgment != nil {
			if err := conn.HandleTxAck(msg.TxAcknowledgment); err != nil {
				logger.WithError(err).Warn("Failed to handle Tx acknowledgement")
			}
		}
	}
}

func (s *impl) GetConcentratorConfig(ctx context.Context, _ *pbtypes.Empty) (*ttnpb.ConcentratorConfig, error) {
	ctx = log.NewContextWithField(ctx, "namespace", "gatewayserver/io/grpc")

	ids := ttnpb.GatewayIdentifiers{
		GatewayID: rpcmetadata.FromIncomingContext(ctx).ID,
	}
	if err := validate.ID(ids.GatewayID); err != nil {
		return nil, err
	}
	if err := rights.RequireGateway(ctx, ids, ttnpb.RIGHT_GATEWAY_LINK); err != nil {
		return nil, err
	}
	fp, err := s.server.GetFrequencyPlan(ctx, ids)
	if err != nil {
		return nil, err
	}
	return fp.ToConcentratorConfig(), nil
}
