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

package applicationserver

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io"
	iogrpc "go.thethings.network/lorawan-stack/pkg/applicationserver/io/grpc"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io/mqtt"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io/web"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/crypto"
	"go.thethings.network/lorawan-stack/pkg/crypto/cryptoutil"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/messageprocessors"
	"go.thethings.network/lorawan-stack/pkg/messageprocessors/cayennelpp"
	"go.thethings.network/lorawan-stack/pkg/messageprocessors/javascript"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"google.golang.org/grpc"
)

// ApplicationServer implements the Application Server component.
//
// The Application Server exposes the As, AppAs and AsEndDeviceRegistry services.
type ApplicationServer struct {
	*component.Component

	linkMode       LinkMode
	linkRegistry   LinkRegistry
	deviceRegistry DeviceRegistry
	formatter      payloadFormatter
	webhooks       web.Webhooks

	links              sync.Map
	defaultSubscribers []*io.Subscription
}

var (
	errListenFrontend = errors.DefineFailedPrecondition(
		"listen_frontend",
		"failed to start frontend listener `{protocol}` on address `{address}`",
	)
)

// New returns new *ApplicationServer.
func New(c *component.Component, conf *Config) (as *ApplicationServer, err error) {
	linkMode, err := conf.GetLinkMode()
	if err != nil {
		return nil, err
	}
	as = &ApplicationServer{
		Component:      c,
		linkMode:       linkMode,
		linkRegistry:   conf.Links,
		deviceRegistry: conf.Devices,
		formatter: payloadFormatter{
			repository: c.GetBaseConfig(c.Context()).DeviceRepository.Client(),
			upFormatters: map[ttnpb.PayloadFormatter]messageprocessors.PayloadDecoder{
				ttnpb.PayloadFormatter_FORMATTER_JAVASCRIPT: javascript.New(),
				ttnpb.PayloadFormatter_FORMATTER_CAYENNELPP: cayennelpp.New(),
			},
			downFormatters: map[ttnpb.PayloadFormatter]messageprocessors.PayloadEncoder{
				ttnpb.PayloadFormatter_FORMATTER_JAVASCRIPT: javascript.New(),
				ttnpb.PayloadFormatter_FORMATTER_CAYENNELPP: cayennelpp.New(),
			},
		},
	}

	ctx, cancel := context.WithCancel(c.Context())
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	for _, version := range []struct {
		Format mqtt.Format
		Config MQTTConfig
	}{
		{
			Format: mqtt.JSON,
			Config: conf.MQTT,
		},
	} {
		for _, lis := range []struct {
			Listen   string
			Protocol string
			Net      func(component.Listener) (net.Listener, error)
		}{
			{
				Listen:   version.Config.Listen,
				Protocol: "tcp",
				Net:      component.Listener.TCP,
			},
			{
				Listen:   version.Config.ListenTLS,
				Protocol: "tls",
				Net:      component.Listener.TLS,
			},
		} {
			if lis.Listen == "" {
				continue
			}
			var componentLis component.Listener
			var netLis net.Listener
			componentLis, err = as.ListenTCP(lis.Listen)
			if err == nil {
				netLis, err = lis.Net(componentLis)
			}
			if err != nil {
				return nil, errListenFrontend.WithCause(err).WithAttributes(
					"protocol", lis.Protocol,
					"address", lis.Listen,
				)
			}
			mqtt.Start(ctx, as, netLis, version.Format, lis.Protocol)
		}
	}

	if webhooks, err := conf.Webhooks.NewWebhooks(as.FillContext(as.Context()), as); err != nil {
		return nil, err
	} else if webhooks != nil {
		as.webhooks = webhooks
		as.defaultSubscribers = append(as.defaultSubscribers, webhooks.NewSubscription())
		c.RegisterWeb(webhooks)
	}

	c.RegisterGRPC(as)
	if as.linkMode == LinkAll {
		c.RegisterTask("link_all", as.linkAll, component.TaskRestartOnFailure)
	}
	return as, nil
}

// RegisterServices registers services provided by as at s.
func (as *ApplicationServer) RegisterServices(s *grpc.Server) {
	ttnpb.RegisterAsServer(s, as)
	ttnpb.RegisterAsEndDeviceRegistryServer(s, &deviceRegistryRPC{
		registry: as.deviceRegistry,
	})
	ttnpb.RegisterAppAsServer(s, iogrpc.New(as))
	if as.webhooks != nil {
		ttnpb.RegisterApplicationWebhookRegistryServer(s, web.NewWebhookRegistryRPC(as.webhooks.Registry()))
	}
}

// RegisterHandlers registers gRPC handlers.
func (as *ApplicationServer) RegisterHandlers(s *runtime.ServeMux, conn *grpc.ClientConn) {
	ttnpb.RegisterAsHandler(as.Context(), s, conn)
	ttnpb.RegisterAsEndDeviceRegistryHandler(as.Context(), s, conn)
	if as.webhooks != nil {
		ttnpb.RegisterApplicationWebhookRegistryHandler(as.Context(), s, conn)
	}
}

// Roles returns the roles that the Application Server fulfills.
func (as *ApplicationServer) Roles() []ttnpb.PeerInfo_Role {
	return []ttnpb.PeerInfo_Role{ttnpb.PeerInfo_APPLICATION_SERVER}
}

// FillApplicationContext fills the given context and identifiers.
func (as *ApplicationServer) FillApplicationContext(ctx context.Context, ids ttnpb.ApplicationIdentifiers) (context.Context, ttnpb.ApplicationIdentifiers, error) {
	return as.FillContext(ctx), ids, nil
}

// Subscribe subscribes an application or integration by its identifiers to the Application Server, and returns a
// io.Subscription for traffic and control.
func (as *ApplicationServer) Subscribe(ctx context.Context, protocol string, ids ttnpb.ApplicationIdentifiers) (*io.Subscription, error) {
	if err := rights.RequireApplication(ctx, ids, ttnpb.RIGHT_APPLICATION_TRAFFIC_READ); err != nil {
		return nil, err
	}

	uid := unique.ID(ctx, ids)
	ctx = log.NewContextWithField(
		events.ContextWithCorrelationID(ctx, fmt.Sprintf("as:conn:%s", events.NewCorrelationID())),
		"application_uid", uid,
	)

	l, err := as.getLink(ctx, ids)
	if err != nil {
		return nil, err
	}
	sub := io.NewSubscription(ctx, protocol, &ids)
	l.subscribeCh <- sub
	go func() {
		<-sub.Context().Done()
		l.unsubscribeCh <- sub
	}()
	return sub, nil
}

var (
	errDeviceNotFound  = errors.DefineNotFound("device_not_found", "device `{device_uid}` not found")
	errNoDeviceSession = errors.DefineFailedPrecondition("no_device_session", "no device session; check device activation")
)

func (as *ApplicationServer) downlinkQueueOp(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, items []*ttnpb.ApplicationDownlink, op func(ttnpb.AsNsClient, context.Context, *ttnpb.DownlinkQueueRequest, ...grpc.CallOption) (*pbtypes.Empty, error)) error {
	ctx = events.ContextWithCorrelationID(ctx, fmt.Sprintf("as:downlink:%s", events.NewCorrelationID()))
	for _, item := range items {
		item.CorrelationIDs = append(item.CorrelationIDs, events.CorrelationIDsFromContext(ctx)...)
	}
	logger := log.FromContext(ctx)
	link, err := as.getLink(ctx, ids.ApplicationIdentifiers)
	if err != nil {
		return err
	}
	<-link.connReady
	var encryptedItems []*ttnpb.ApplicationDownlink
	_, err = as.deviceRegistry.Set(ctx, ids,
		[]string{
			"session",
			"formatters",
			"version_ids",
		},
		func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
			if dev == nil {
				return nil, nil, errDeviceNotFound.WithAttributes("device_uid", unique.ID(ctx, ids))
			}
			if dev.Session == nil {
				return nil, nil, errNoDeviceSession
			}
			for _, item := range items {
				registerReceiveDownlink(ctx, ids, item)
				item.SessionKeyID = dev.Session.SessionKeyID
				item.FCnt = dev.Session.LastAFCntDown + 1
				if err := as.encodeAndEncrypt(ctx, dev, item, link.DefaultFormatters); err != nil {
					logger.WithError(err).Warn("Dropping downlink message; encoding and encryption failed")
					registerDropDownlink(ctx, ids, item, err)
					continue
				}
				item.DecodedPayload = nil
				item.CorrelationIDs = events.CorrelationIDsFromContext(ctx)
				dev.Session.LastAFCntDown = item.FCnt
				encryptedItems = append(encryptedItems, item)
			}
			client := ttnpb.NewAsNsClient(link.conn)
			req := &ttnpb.DownlinkQueueRequest{
				EndDeviceIdentifiers: ids,
				Downlinks:            encryptedItems,
			}
			_, err = op(client, ctx, req, link.callOpts...)
			if err != nil {
				return nil, nil, err
			}
			return dev, []string{"session.last_a_f_cnt_down"}, nil
		},
	)
	if err != nil {
		for _, item := range encryptedItems {
			registerDropDownlink(ctx, ids, item, err)
		}
		return err
	}
	atomic.AddUint64(&link.downlinks, uint64(len(encryptedItems)))
	atomic.StoreInt64(&link.lastDownlinkTime, time.Now().UnixNano())
	for _, item := range encryptedItems {
		registerForwardDownlink(ctx, ids, item, link.connName)
	}
	return nil
}

// DownlinkQueuePush pushes the given downlink messages to the end device's application downlink queue.
// This operation changes FRMPayload in the given items.
func (as *ApplicationServer) DownlinkQueuePush(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, items []*ttnpb.ApplicationDownlink) error {
	return as.downlinkQueueOp(ctx, ids, items, ttnpb.AsNsClient.DownlinkQueuePush)
}

// DownlinkQueueReplace replaces the end device's application downlink queue with the given downlink messages.
// This operation changes FRMPayload in the given items.
func (as *ApplicationServer) DownlinkQueueReplace(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, items []*ttnpb.ApplicationDownlink) error {
	return as.downlinkQueueOp(ctx, ids, items, ttnpb.AsNsClient.DownlinkQueueReplace)
}

var errNoAppSKey = errors.DefineCorruption("no_app_s_key", "no AppSKey")

// DownlinkQueueList lists the application downlink queue of the given end device.
func (as *ApplicationServer) DownlinkQueueList(ctx context.Context, ids ttnpb.EndDeviceIdentifiers) ([]*ttnpb.ApplicationDownlink, error) {
	dev, err := as.deviceRegistry.Get(ctx, ids, []string{"session"})
	if err != nil {
		return nil, err
	}
	if dev.Session == nil {
		return nil, errNoDeviceSession
	}
	if dev.Session.AppSKey == nil {
		return nil, errNoAppSKey
	}
	appSKey, err := cryptoutil.UnwrapAES128Key(*dev.Session.AppSKey, as.KeyVault)
	if err != nil {
		return nil, err
	}
	link, err := as.getLink(ctx, ids.ApplicationIdentifiers)
	if err != nil {
		return nil, err
	}
	<-link.connReady
	client := ttnpb.NewAsNsClient(link.conn)
	res, err := client.DownlinkQueueList(ctx, &ids, link.callOpts...)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Downlinks {
		item.FRMPayload, err = crypto.DecryptDownlink(appSKey, dev.Session.DevAddr, item.FCnt, item.FRMPayload)
		if err != nil {
			return nil, err
		}
	}
	return res.Downlinks, nil
}

var errJSUnavailable = errors.DefineUnavailable("join_server_unavailable", "Join Server unavailable for JoinEUI `{join_eui}`")

func (as *ApplicationServer) fetchAppSKey(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, sessionKeyID []byte) (ttnpb.KeyEnvelope, error) {
	// TODO: Lookup Join Server (https://github.com/TheThingsNetwork/lorawan-stack/issues/4)
	js := as.GetPeer(ctx, ttnpb.PeerInfo_JOIN_SERVER, ids)
	if js == nil {
		return ttnpb.KeyEnvelope{}, errJSUnavailable.WithAttributes("join_eui", *ids.JoinEUI)
	}
	client := ttnpb.NewAsJsClient(js.Conn())
	req := &ttnpb.SessionKeyRequest{
		SessionKeyID: sessionKeyID,
		DevEUI:       *ids.DevEUI,
	}
	res, err := client.GetAppSKey(ctx, req, as.WithClusterAuth())
	if err != nil {
		return ttnpb.KeyEnvelope{}, err
	}
	return res.AppSKey, nil
}

func (as *ApplicationServer) handleUp(ctx context.Context, up *ttnpb.ApplicationUp, link *link) error {
	ctx = log.NewContextWithField(ctx, "device_uid", unique.ID(ctx, up.EndDeviceIdentifiers))
	switch p := up.Up.(type) {
	case *ttnpb.ApplicationUp_JoinAccept:
		return as.handleJoinAccept(ctx, up.EndDeviceIdentifiers, p.JoinAccept, link)
	case *ttnpb.ApplicationUp_UplinkMessage:
		return as.handleUplink(ctx, up.EndDeviceIdentifiers, p.UplinkMessage, link)
	case *ttnpb.ApplicationUp_DownlinkQueueInvalidated:
		return as.handleDownlinkQueueInvalidated(ctx, up.EndDeviceIdentifiers, p.DownlinkQueueInvalidated, link)
	case *ttnpb.ApplicationUp_DownlinkQueued:
		return as.decryptDownlinkMessage(ctx, up.EndDeviceIdentifiers, p.DownlinkQueued)
	case *ttnpb.ApplicationUp_DownlinkSent:
		return as.decryptDownlinkMessage(ctx, up.EndDeviceIdentifiers, p.DownlinkSent)
	case *ttnpb.ApplicationUp_DownlinkFailed:
		return as.decryptDownlinkMessage(ctx, up.EndDeviceIdentifiers, &p.DownlinkFailed.ApplicationDownlink)
	case *ttnpb.ApplicationUp_DownlinkAck:
		return as.decryptDownlinkMessage(ctx, up.EndDeviceIdentifiers, p.DownlinkAck)
	case *ttnpb.ApplicationUp_DownlinkNack:
		return as.handleDownlinkNack(ctx, up.EndDeviceIdentifiers, p.DownlinkNack, link)
	default:
		return nil
	}
}

var errFetchAppSKey = errors.Define("app_s_key", "failed to get AppSKey")

func (as *ApplicationServer) handleJoinAccept(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, joinAccept *ttnpb.ApplicationJoinAccept, link *link) error {
	logger := log.FromContext(ctx).WithFields(log.Fields(
		"join_eui", ids.JoinEUI,
		"dev_eui", ids.DevEUI,
		"session_key_id", joinAccept.SessionKeyID,
	))
	_, err := as.deviceRegistry.Set(ctx, ids,
		[]string{
			"session",
			"pending_session",
		},
		func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
			var mask []string
			if dev == nil {
				return nil, nil, errDeviceNotFound.WithAttributes("device_uid", unique.ID(ctx, ids))
			}
			var appSKey ttnpb.KeyEnvelope
			if joinAccept.AppSKey != nil {
				logger.Debug("Received AppSKey from Network Server")
				appSKey = *joinAccept.AppSKey
			} else {
				logger.Debug("Fetching AppSKey from Join Server...")
				key, err := as.fetchAppSKey(ctx, ids, joinAccept.SessionKeyID)
				if err != nil {
					return nil, nil, errFetchAppSKey.WithCause(err)
				}
				appSKey = key
				logger.Debug("Fetched AppSKey from Join Server")
			}
			session := &ttnpb.Session{
				DevAddr: *ids.DevAddr,
				SessionKeys: ttnpb.SessionKeys{
					SessionKeyID: joinAccept.SessionKeyID,
					AppSKey:      &appSKey,
				},
				StartedAt: time.Now().UTC(),
			}
			if joinAccept.PendingSession {
				dev.PendingSession = session
				mask = append(mask, "pending_session")
			} else {
				previousSession := dev.Session
				dev.Session = session
				dev.PendingSession = nil
				mask = append(mask, "session", "pending_session")
				if len(joinAccept.InvalidatedDownlinks) > 0 {
					// The Network Server reset the downlink queue as the new security session invalidated it. The invalidated
					// downlink queue is passed as part of the join-accept and the Application Server should recalculate it. This
					// changes the LastAFCntDown in the session, so it should be run as part of the transaction.
					if err := as.recalculateDownlinkQueue(ctx, dev, previousSession, joinAccept.InvalidatedDownlinks, 1, link); err != nil {
						logger.WithError(err).WithField("count", len(joinAccept.InvalidatedDownlinks)).Warn("Failed to recalculate downlink queue, items lost")
					}
				}
			}
			return dev, mask, nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (as *ApplicationServer) handleUplink(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, uplink *ttnpb.ApplicationUplink, link *link) error {
	ctx = log.NewContextWithField(ctx, "session_key_id", uplink.SessionKeyID)
	logger := log.FromContext(ctx)
	dev, err := as.deviceRegistry.Set(ctx, ids,
		[]string{
			"session",
			"pending_session",
			"formatters",
			"version_ids",
		},
		func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
			var mask []string
			if dev == nil {
				return nil, nil, errDeviceNotFound.WithAttributes("device_uid", unique.ID(ctx, ids))
			}
			if dev.Session == nil || !bytes.Equal(dev.Session.SessionKeyID, uplink.SessionKeyID) {
				logger := logger.WithField("session_key_id", uplink.SessionKeyID)
				previousSession := dev.Session
				if dev.PendingSession != nil && bytes.Equal(dev.PendingSession.SessionKeyID, uplink.SessionKeyID) {
					logger.Debug("Switching to pending session")
					dev.Session = dev.PendingSession
				} else {
					appSKey, err := as.fetchAppSKey(ctx, ids, uplink.SessionKeyID)
					if err != nil {
						return nil, nil, errFetchAppSKey.WithCause(err)
					}
					dev.Session = &ttnpb.Session{
						DevAddr: *ids.DevAddr,
						SessionKeys: ttnpb.SessionKeys{
							SessionKeyID: uplink.SessionKeyID,
							AppSKey:      &appSKey,
						},
						StartedAt: time.Now().UTC(),
					}
					logger.Debug("Restored session")
				}
				dev.PendingSession = nil
				// At this point, the application downlink queue in the Network Server is invalid; recalculation is necessary.
				// Next AFCntDown 1 is assumed. If this is a LoRaWAN 1.0.x end device and the Network Server sent MAC layer
				// downlink already, the Network Server will trigger the DownlinkQueueInvalidated event. Therefore, this
				// recalculation may result in another recalculation.
				client := ttnpb.NewAsNsClient(link.conn)
				res, err := client.DownlinkQueueList(ctx, &ids, link.callOpts...)
				if err != nil {
					log.WithError(err).Warn("Failed to list downlink queue for recalculation; clearing the downlink queue")
					req := &ttnpb.DownlinkQueueRequest{
						EndDeviceIdentifiers: ids,
					}
					_, err = client.DownlinkQueueReplace(ctx, req, link.callOpts...)
					if err != nil {
						log.WithError(err).Warn("Failed to clear the downlink queue; any queued items in the Network Server are invalid")
						events.Publish(evtInvalidQueueDataDown(ctx, ids, err))
					} else {
						events.Publish(evtLostQueueDataDown(ctx, ids, err))
					}
				} else if err := as.recalculateDownlinkQueue(ctx, dev, previousSession, res.Downlinks, 1, link); err != nil {
					log.WithError(err).Warn("Failed to recalculate downlink queue")
				}
				mask = append(mask, "session", "pending_session")
			} else if dev.Session.AppSKey == nil {
				return nil, nil, errNoAppSKey
			}
			return dev, mask, nil
		},
	)
	if err != nil {
		return err
	}
	if err := as.decryptAndDecode(ctx, dev, uplink, link.DefaultFormatters); err != nil {
		return err
	}
	// TODO: Run uplink messages through location solvers async (https://github.com/TheThingsNetwork/lorawan-stack/issues/37)
	return nil
}

func (as *ApplicationServer) handleDownlinkQueueInvalidated(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, invalid *ttnpb.ApplicationInvalidatedDownlinks, link *link) error {
	_, err := as.deviceRegistry.Set(ctx, ids,
		[]string{
			"session",
		},
		func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
			if dev == nil {
				return nil, nil, errDeviceNotFound.WithAttributes("device_uid", unique.ID(ctx, ids))
			}
			if err := as.recalculateDownlinkQueue(ctx, dev, nil, invalid.Downlinks, invalid.LastFCntDown+1, link); err != nil {
				return nil, nil, err
			}
			return dev, []string{"session"}, nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (as *ApplicationServer) handleDownlinkNack(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, msg *ttnpb.ApplicationDownlink, link *link) error {
	client := ttnpb.NewAsNsClient(link.conn)
	res, err := client.DownlinkQueueList(ctx, &ids, link.callOpts...)
	if err != nil {
		log.WithError(err).Warn("Failed to list the downlink queue for inserting nacked downlink message")
		registerDropDownlink(ctx, ids, msg, err)
	} else {
		queue := append([]*ttnpb.ApplicationDownlink{msg}, res.Downlinks...)
		_, err := as.deviceRegistry.Set(ctx, ids,
			[]string{
				"session",
			},
			func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
				if err := as.recalculateDownlinkQueue(ctx, dev, nil, queue, msg.FCnt+1, link); err != nil {
					return nil, nil, err
				}
				return dev, []string{"session"}, nil
			},
		)
		if err != nil {
			log.WithError(err).Warn("Failed to recalculate downlink queue with inserted nacked downlink message")
			registerDropDownlink(ctx, ids, msg, err)
		}
	}
	// Decrypt the message as it will be sent to upstream after handling it.
	if err := as.decryptDownlinkMessage(ctx, ids, msg); err != nil {
		return err
	}
	return nil
}

func (as *ApplicationServer) decryptDownlinkMessage(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, msg *ttnpb.ApplicationDownlink) error {
	dev, err := as.deviceRegistry.Get(ctx, ids, []string{"session"})
	if err != nil {
		return err
	}
	if dev.Session == nil || !bytes.Equal(dev.Session.SessionKeyID, msg.SessionKeyID) || dev.Session.AppSKey == nil {
		return errNoAppSKey
	}
	appSKey, err := cryptoutil.UnwrapAES128Key(*dev.Session.AppSKey, as.KeyVault)
	if err != nil {
		return err
	}
	msg.FRMPayload, err = crypto.DecryptDownlink(appSKey, dev.Session.DevAddr, msg.FCnt, msg.FRMPayload)
	if err != nil {
		return err
	}
	return nil
}

// recalculateDownlinkQueue decrypts items in the given invalid downlink queue, encrypts the items with frame counters
// starting from the given frame counter, and replaces the downlink queue in the Network Server.
// If re-encrypting a message fails, the message is skipped.
// This method requires the given end device's session to be set. This method mutates the end device's session LastAFCntDown.
// This method does not change the contents of the given invalid downlink queue.
func (as *ApplicationServer) recalculateDownlinkQueue(ctx context.Context, dev *ttnpb.EndDevice, previousSession *ttnpb.Session, invalid []*ttnpb.ApplicationDownlink, nextAFCntDown uint32, link *link) (err error) {
	logger := log.FromContext(ctx)
	logger.WithFields(log.Fields(
		"count", len(invalid),
		"next_a_f_cnt_down", nextAFCntDown,
	)).Debug("Recalculating downlink queue")
	defer func() {
		// If something fails, clear the downlink queue as an empty downlink queue is better than a downlink queue
		// with items that are encrypted with the wrong AppSKey.
		if err != nil {
			logger.WithError(err).Warn("Recalculate downlink queue failed; clearing the downlink queue")
			dev.Session.LastAFCntDown = nextAFCntDown - 1
			client := ttnpb.NewAsNsClient(link.conn)
			req := &ttnpb.DownlinkQueueRequest{
				EndDeviceIdentifiers: dev.EndDeviceIdentifiers,
			}
			if _, err := client.DownlinkQueueReplace(ctx, req, link.callOpts...); err != nil {
				log.WithError(err).Warn("Failed to clear the downlink queue; any queued items in the Network Server are invalid")
				events.Publish(evtInvalidQueueDataDown(ctx, dev.EndDeviceIdentifiers, err))
			} else {
				events.Publish(evtLostQueueDataDown(ctx, dev.EndDeviceIdentifiers, err))
			}
		}
	}()
	if dev.Session == nil || dev.Session.AppSKey == nil {
		return errNoAppSKey
	}
	newAppSKey, err := cryptoutil.UnwrapAES128Key(*dev.Session.AppSKey, as.KeyVault)
	if err != nil {
		return err
	}
	dev.Session.LastAFCntDown = nextAFCntDown - 1
	valid := make([]*ttnpb.ApplicationDownlink, 0, len(invalid))
	for _, oldItem := range invalid {
		logger := logger.WithFields(log.Fields(
			"f_port", oldItem.FPort,
			"f_cnt", oldItem.FCnt,
			"session_key_id", oldItem.SessionKeyID,
		))
		var oldSession *ttnpb.Session
		for _, s := range []*ttnpb.Session{previousSession, dev.Session} {
			if s != nil && bytes.Equal(s.SessionKeyID, oldItem.SessionKeyID) {
				oldSession = s
				break
			}
		}
		if oldSession == nil || oldSession.AppSKey == nil {
			logger.Warn("Dropping downlink message; session not found or AppSKey not available")
			registerDropDownlink(ctx, dev.EndDeviceIdentifiers, oldItem, err)
			continue
		}
		// TODO: Cache unwrapped keys (https://github.com/TheThingsNetwork/lorawan-stack/issues/36)
		oldAppSKey, err := cryptoutil.UnwrapAES128Key(*oldSession.AppSKey, as.KeyVault)
		if err != nil {
			logger.WithError(err).Warn("Dropping downlink message; failed to unwrap AppSKey for decryption")
			registerDropDownlink(ctx, dev.EndDeviceIdentifiers, oldItem, err)
			continue
		}
		frmPayload, err := crypto.DecryptDownlink(oldAppSKey, oldSession.DevAddr, oldItem.FCnt, oldItem.FRMPayload)
		if err != nil {
			logger.WithError(err).Warn("Dropping downlink message; failed to decrypt")
			registerDropDownlink(ctx, dev.EndDeviceIdentifiers, oldItem, err)
			continue
		}
		newItem := &ttnpb.ApplicationDownlink{
			SessionKeyID:   dev.Session.SessionKeyID,
			FPort:          oldItem.FPort,
			FCnt:           dev.Session.LastAFCntDown + 1,
			Confirmed:      oldItem.Confirmed,
			ClassBC:        oldItem.ClassBC,
			CorrelationIDs: oldItem.CorrelationIDs,
		}
		newItem.FRMPayload, err = crypto.EncryptDownlink(newAppSKey, dev.Session.DevAddr, newItem.FCnt, frmPayload)
		if err != nil {
			logger.WithError(err).Warn("Dropping downlink message; failed to encrypt")
			registerDropDownlink(ctx, dev.EndDeviceIdentifiers, oldItem, err)
			continue
		}
		valid = append(valid, newItem)
		dev.Session.LastAFCntDown = newItem.FCnt
	}
	client := ttnpb.NewAsNsClient(link.conn)
	req := &ttnpb.DownlinkQueueRequest{
		EndDeviceIdentifiers: dev.EndDeviceIdentifiers,
		Downlinks:            valid,
	}
	_, err = client.DownlinkQueueReplace(ctx, req, link.callOpts...)
	return err
}
