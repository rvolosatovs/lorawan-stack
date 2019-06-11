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

package basicstationlns

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/basicstation"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io/basicstationlns/messages"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"go.thethings.network/lorawan-stack/pkg/web"
	"google.golang.org/grpc/metadata"
)

var (
	errEmptyGatewayEUI           = errors.Define("empty_gateway_eui", "empty gateway EUI")
	errMessageTypeNotImplemented = errors.DefineUnimplemented("message_type_not_implemented", "message of type `{type}` is not implemented")
)

type srv struct {
	ctx      context.Context
	server   io.Server
	upgrader *websocket.Upgrader
	tokens   io.DownlinkTokens
}

func (*srv) Protocol() string { return "basicstation" }

// New returns a new Basic Station frontend that can be registered in the web server.
func New(ctx context.Context, server io.Server) web.Registerer {
	ctx = log.NewContextWithField(ctx, "namespace", "gatewayserver/io/basicstation")
	return &srv{
		ctx:      ctx,
		server:   server,
		upgrader: &websocket.Upgrader{},
	}
}

func (s *srv) RegisterRoutes(server *web.Server) {
	group := server.Group(ttnpb.HTTPAPIPrefix + "/gs/io/basicstation")
	group.GET("/discover", s.handleDiscover)
	group.GET("/traffic/:id", s.handleTraffic)
}

func (s *srv) handleDiscover(c echo.Context) error {
	logger := log.FromContext(s.ctx).WithFields(log.Fields(
		"endpoint", "discover",
		"remote_addr", c.Request().RemoteAddr,
	))
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.WithError(err).Debug("Failed to upgrade request to websocket connection")
		return err
	}
	defer ws.Close()

	_, data, err := ws.ReadMessage()
	if err != nil {
		logger.WithError(err).Debug("Failed to read message")
		return err
	}
	var req messages.DiscoverQuery
	if err := json.Unmarshal(data, &req); err != nil {
		logger.WithError(err).Debug("Failed to parse discover query message")
		return err
	}

	if req.EUI.IsZero() {
		writeDiscoverError(s.ctx, ws, "Invalid request")
		return errEmptyGatewayEUI
	}

	ids := ttnpb.GatewayIdentifiers{
		EUI: &req.EUI.EUI64,
	}
	ctx, ids, err := s.server.FillGatewayContext(s.ctx, ids)
	if err != nil {
		logger.WithError(err).Debug("Failed to fill gateway context")
		writeDiscoverError(ctx, ws, "Router not provisioned")
		return err
	}

	scheme := "ws"
	if c.IsTLS() {
		scheme = "wss"
	}
	res := messages.DiscoverResponse{
		EUI: req.EUI,
		Muxs: basicstation.EUI{
			Prefix: "muxs",
		},
		URI: fmt.Sprintf("%s://%s%s", scheme, c.Request().Host, c.Echo().URI(s.handleTraffic, ids.GatewayID)),
	}
	data, err = json.Marshal(res)
	if err != nil {
		logger.WithError(err).Warn("Failed to marshal response message")
		writeDiscoverError(ctx, ws, "Router not provisioned")
		return err
	}
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		logger.WithError(err).Warn("Failed to write discover response message")
		return err
	}
	logger.Debug("Sent discover response message")
	return nil
}

func (s *srv) handleTraffic(c echo.Context) error {
	id := c.Param("id")
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	ctx := s.ctx
	if auth != "" {
		md := metadata.New(map[string]string{
			"id":            id,
			"authorization": fmt.Sprintf("Bearer %s", auth),
		})
		if ctxMd, ok := metadata.FromIncomingContext(s.ctx); ok {
			md = metadata.Join(ctxMd, md)
		}
		ctx = metadata.NewIncomingContext(s.ctx, md)
	}

	logger := log.FromContext(ctx).WithFields(log.Fields(
		"endpoint", "traffic",
		"remote_addr", c.Request().RemoteAddr,
	))

	ctx, ids, err := s.server.FillGatewayContext(ctx, ttnpb.GatewayIdentifiers{GatewayID: id})
	if err != nil {
		logger.WithError(err).Debug("Failed to fill gateway context")
		return err
	}

	if auth == "" {
		ctx = rights.NewContext(ctx, rights.Rights{
			GatewayRights: map[string]*ttnpb.Rights{
				id: {
					Rights: []ttnpb.Right{ttnpb.RIGHT_GATEWAY_LINK},
				},
			},
		})
	}

	uid := unique.ID(ctx, ids)
	ctx = log.NewContextWithField(ctx, "gateway_uid", uid)

	fp, err := s.server.GetFrequencyPlan(ctx, ids)
	if err != nil {
		logger.WithError(err).Warn("Failed to get frequency plan")
		return err
	}

	conn, err := s.server.Connect(ctx, s, ids)
	if err != nil {
		logger.WithError(err).Warn("Failed to connect")
		return err
	}
	if err := s.server.ClaimDownlink(ctx, ids); err != nil {
		logger.WithError(err).Error("Failed to claim downlink")
		return err
	}
	defer func() {
		if err := s.server.UnclaimDownlink(ctx, ids); err != nil {
			logger.WithError(err).Error("Failed to unclaim downlink")
		}
	}()

	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.WithError(err).Debug("Failed to upgrade request to websocket connection")
		return err
	}
	defer ws.Close()

	// Process downlinks in a separate go routine
	go func() {
		for {
			select {
			case <-conn.Context().Done():
				return
			case down := <-conn.Down():
				dlTime := time.Now()
				xTimeUpper, ok := conn.RetrieveCustomValue().(int64)
				if !ok {
					logger.WithError(err).Error("Failed to retrieve XTime")
					continue
				}
				dnmsg := messages.FromDownlinkMessage(ids, *down, int64(s.tokens.Next(down.CorrelationIDs, dlTime)), dlTime, xTimeUpper)
				msg, err := dnmsg.MarshalJSON()
				if err != nil {
					logger.WithError(err).Error("Failed to marshal downlink message")
					continue
				}

				logger.Info("Send downlink message")
				if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
					logger.WithError(err).Error("Failed to send downlink message")
					conn.Disconnect(err)
				}
			}
		}
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			logger.WithError(err).Debug("Failed to read message")
			conn.Disconnect(err)
			return nil
		}

		typ, err := messages.Type(data)
		if err != nil {
			logger.WithError(err).Debug("Failed to parse message type")
			continue
		}
		logger = logger.WithFields(log.Fields(
			"upstream_type", typ,
		))
		receivedAt := time.Now()

		switch typ {
		case messages.TypeUpstreamVersion:
			var version messages.Version
			if err := json.Unmarshal(data, &version); err != nil {
				logger.WithError(err).Debug("Failed to unmarshal version message")
				return err
			}
			logger = logger.WithFields(log.Fields(
				"station", version.Station,
				"firmware", version.Firmware,
				"model", version.Model,
			))
			cfg, err := messages.GetRouterConfig(*fp, version.IsProduction(), time.Now())
			if err != nil {
				logger.WithError(err).Warn("Failed to generate router configuration")
				return err
			}
			data, err = cfg.MarshalJSON()
			if err != nil {
				logger.WithError(err).Warn("Failed to marshal response message")
				return err
			}
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.WithError(err).Warn("Failed to send router configuration")
				return err
			}

		case messages.TypeUpstreamJoinRequest:
			var jreq messages.JoinRequest
			if err := json.Unmarshal(data, &jreq); err != nil {
				logger.WithError(err).Warn("Failed to unmarshal join-request message")
				return nil
			}
			conn.StoreCustomValue(jreq.UpInfo.XTime)
			up, err := jreq.ToUplinkMessage(ids, fp.BandID, receivedAt)
			if err != nil {
				logger.WithError(err).Debug("Failed to parse join-request message")
				return nil
			}
			if err := conn.HandleUp(up); err != nil {
				logger.WithError(err).Warn("Failed to handle uplink message")
			}
			recordRTT(conn, receivedAt, jreq.RefTime)

		case messages.TypeUpstreamUplinkDataFrame:
			var updf messages.UplinkDataFrame
			if err := json.Unmarshal(data, &updf); err != nil {
				logger.WithError(err).Warn("Failed to unmarshal uplink data frame")
				return nil
			}
			up, err := updf.ToUplinkMessage(ids, fp.BandID, receivedAt)
			conn.StoreCustomValue(updf.UpInfo.XTime)
			if err != nil {
				logger.WithError(err).Debug("Failed to parse uplink data frame")
				return nil
			}
			if err := conn.HandleUp(up); err != nil {
				logger.WithError(err).Warn("Failed to handle uplink message")
			}
			recordRTT(conn, receivedAt, updf.RefTime)

		case messages.TypeUpstreamTxConfirmation:
			var txConf messages.TxConfirmation
			if err := json.Unmarshal(data, &txConf); err != nil {
				logger.WithError(err).Warn("Failed to unmarshal Tx acknowledgement frame")
				return nil
			}
			if cids, _, ok := s.tokens.Get(uint16(txConf.Diid), receivedAt); ok {
				txAck := messages.ToTxAcknowledgment(cids)
				if err := conn.HandleTxAck(&txAck); err != nil {
					logger.WithField("diid", txConf.Diid).Debug("Failed to handle Tx acknowledgement")
				}
			} else {
				logger.WithField("diid", txConf.Diid).Debug("Tx acknowledgement either does not correspond to a downlink message or arrived too late")
			}
			recordRTT(conn, receivedAt, txConf.RefTime)

		case messages.TypeUpstreamProprietaryDataFrame:
			return errMessageTypeNotImplemented.WithAttributes("type", typ)
		case messages.TypeUpstreamRemoteShell:
			return errMessageTypeNotImplemented.WithAttributes("type", typ)
		case messages.TypeUpstreamTimeSync:
			return errMessageTypeNotImplemented.WithAttributes("type", typ)

		default:
			// Unknown message types are ignored by the server
			logger.WithField("message_type", typ).Debug("Unknown message type")
		}
	}
}

// writeDiscoverError sends the error messages during the discovery on the WS connection to the station.
func writeDiscoverError(ctx context.Context, ws *websocket.Conn, msg string) {
	logger := log.FromContext(ctx)
	errMsg, err := json.Marshal(messages.DiscoverResponse{Error: msg})
	if err != nil {
		logger.WithError(err).Debug("Failed to marshal error message")
		return
	}
	if err := ws.WriteMessage(websocket.TextMessage, errMsg); err != nil {
		logger.WithError(err).Debug("Failed to write error response message")
	}
}

func recordRTT(conn *io.Connection, receivedAt time.Time, refTime float64) {
	sec, nsec := math.Modf(refTime)
	if sec != 0 {
		ref := time.Unix(int64(sec), int64(nsec*1e9))
		conn.RecordRTT(receivedAt.Sub(ref))
	}
}
