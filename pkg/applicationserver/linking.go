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
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/errorcontext"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcclient"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware/discover"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (as *ApplicationServer) linkAll(ctx context.Context) error {
	return as.linkRegistry.Range(ctx, nil,
		func(ctx context.Context, ids ttnpb.ApplicationIdentifiers, _ *ttnpb.ApplicationLink) bool {
			as.startLinkTask(ctx, ids)
			return true
		},
	)
}

func (as *ApplicationServer) startLinkTask(ctx context.Context, ids ttnpb.ApplicationIdentifiers) {
	ctx = log.NewContextWithField(ctx, "application_uid", unique.ID(ctx, ids))
	as.StartTask(&component.TaskConfig{
		Context: ctx,
		ID:      "link",
		Func: func(ctx context.Context) error {
			target, err := as.linkRegistry.Get(ctx, ids, []string{
				"network_server_address",
				"api_key",
				"default_formatters",
				"skip_payload_crypto",
			})
			if err != nil && !errors.IsNotFound(err) {
				return err
			} else if err != nil {
				log.FromContext(ctx).WithError(err).Warn("Link not found")
				return nil
			}

			return as.link(ctx, ids, target)
		},
		Restart: component.TaskRestartOnFailure,
		Backoff: io.DialTaskBackoffConfig,
	})
}

type upstreamTrafficHandler func(context.Context, *ttnpb.ApplicationUp, *link) (pass bool, err error)

type link struct {
	// Align for sync/atomic.
	ups,
	downlinks uint64
	linkTime,
	lastUpTime,
	lastDownlinkTime int64

	ttnpb.ApplicationIdentifiers
	ttnpb.ApplicationLink
	ctx    context.Context
	cancel errorcontext.CancelFunc
	closed chan struct{}

	conn      *grpc.ClientConn
	connName  string
	connReady chan struct{}
	callOpts  []grpc.CallOption

	handleUp upstreamTrafficHandler

	subscribeCh   chan *io.Subscription
	unsubscribeCh chan *io.Subscription
	upCh          chan *io.ContextualApplicationUp
}

const linkBufferSize = 10

var errNSPeerNotFound = errors.DefineNotFound("network_server_not_found", "Network Server not found for `{application_uid}`")

func (as *ApplicationServer) connectLink(ctx context.Context, link *link) error {
	var allowInsecure bool
	if link.NetworkServerAddress != "" {
		allowInsecure = as.AllowInsecureForCredentials()
		var (
			dialOpt   grpc.DialOption
			target    string
			targetErr error
		)
		if link.TLS {
			tlsConfig, err := as.GetTLSClientConfig(ctx)
			if err != nil {
				return err
			}
			dialOpt = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
			target, targetErr = discover.Address(ttnpb.ClusterRole_NETWORK_SERVER, link.NetworkServerAddress)
		} else {
			dialOpt = grpc.WithInsecure()
			target, targetErr = discover.DefaultPort(link.NetworkServerAddress, discover.DefaultPorts[false])
		}
		if targetErr != nil {
			return targetErr
		}
		conn, err := grpc.DialContext(ctx, target, append(rpcclient.DefaultDialOptions(ctx), dialOpt)...)
		if err != nil {
			return err
		}
		link.conn = conn
		link.connName = link.NetworkServerAddress
		go func() {
			<-ctx.Done()
			conn.Close()
		}()
	} else {
		allowInsecure = !as.ClusterTLS()
		ns, err := as.GetPeer(ctx, ttnpb.ClusterRole_NETWORK_SERVER, link.ApplicationIdentifiers)
		if err != nil {
			return errNSPeerNotFound.WithCause(err).WithAttributes("application_uid", unique.ID(ctx, link.ApplicationIdentifiers))
		}
		cc, err := ns.Conn()
		if err != nil {
			return errNSPeerNotFound.WithCause(err).WithAttributes("application_uid", unique.ID(ctx, link.ApplicationIdentifiers))
		}
		link.conn = cc
		link.connName = ns.Name()
	}
	link.callOpts = []grpc.CallOption{
		grpc.PerRPCCredentials(rpcmetadata.MD{
			ID:            link.ApplicationID,
			AuthType:      "Bearer",
			AuthValue:     link.APIKey,
			AllowInsecure: allowInsecure,
		}),
	}
	link.linkTime = time.Now().UnixNano()
	close(link.connReady)
	return nil
}

func (as *ApplicationServer) link(ctx context.Context, ids ttnpb.ApplicationIdentifiers, target *ttnpb.ApplicationLink) (err error) {
	uid := unique.ID(ctx, ids)
	ctx = log.NewContextWithField(ctx, "application_uid", uid)
	ctx, cancel := errorcontext.New(ctx)
	defer func() {
		cancel(err)
	}()
	l := &link{
		ApplicationIdentifiers: ids,
		ApplicationLink:        *target,
		ctx:                    ctx,
		cancel:                 cancel,
		closed:                 make(chan struct{}),
		connReady:              make(chan struct{}),
		handleUp:               as.handleUp,
		subscribeCh:            make(chan *io.Subscription, 1),
		unsubscribeCh:          make(chan *io.Subscription, 1),
		upCh:                   make(chan *io.ContextualApplicationUp, linkBufferSize),
	}
	if _, loaded := as.links.LoadOrStore(uid, l); loaded {
		log.FromContext(ctx).Warn("Link already started")
		return nil
	}
	go func() {
		<-ctx.Done()
		as.linkErrors.Store(uid, ctx.Err())
		as.links.Delete(uid)
		if err := ctx.Err(); err != nil {
			log.FromContext(ctx).WithError(err).Warn("Link failed")
			registerLinkFail(ctx, l, err)
		}
		close(l.closed)
	}()
	if err := as.connectLink(ctx, l); err != nil {
		return err
	}
	client := ttnpb.NewAsNsClient(l.conn)
	ctx = log.NewContextWithField(ctx, "network_server", l.connName)
	logger := log.FromContext(ctx)
	logger.Debug("Link")
	stream, err := client.LinkApplication(ctx, l.callOpts...)
	if err != nil {
		logger.WithError(err).Warn("Link setup failed")
		return err
	}
	logger.Info("Linked")
	registerLinkStart(ctx, l)
	go func() {
		<-ctx.Done()
		logger.WithError(ctx.Err()).Info("Unlinked")
		registerLinkStop(ctx, l)
	}()

	go l.run()
	for _, sub := range as.defaultSubscribers {
		sub := sub
		select {
		case <-l.ctx.Done():
			return
		case l.subscribeCh <- sub:
		}
		go func() {
			select {
			case <-l.ctx.Done():
				// Default subscriptions should not be canceled on link failures,
				// and they should skip the subscribe channel since it will get
				// closed.
				return
			case <-sub.Context().Done():
			}
			select {
			case <-l.ctx.Done():
				return
			case l.unsubscribeCh <- sub:
			}
		}()
	}
	for {
		up, err := stream.Recv()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return err
			}
		}
		atomic.AddUint64(&l.ups, 1)
		atomic.StoreInt64(&l.lastUpTime, time.Now().UnixNano())

		err = l.sendUp(ctx, up, func() error { return stream.Send(ttnpb.Empty) })
		if err != nil {
			return err
		}
	}
}

var (
	errNotLinked  = errors.DefineNotFound("not_linked", "not linked to `{application_uid}`")
	errLinkFailed = errors.DefineAborted("link", "link failed")
)

func (as *ApplicationServer) cancelLink(ctx context.Context, ids ttnpb.ApplicationIdentifiers) error {
	uid := unique.ID(ctx, ids)
	if val, ok := as.links.Load(uid); ok {
		l := val.(*link)
		log.FromContext(ctx).WithField("application_uid", uid).Debug("Unlink")
		l.cancel(context.Canceled)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-l.closed:
		}
	} else {
		as.linkErrors.Delete(uid)
	}
	return nil
}

func (as *ApplicationServer) getLink(ctx context.Context, ids ttnpb.ApplicationIdentifiers) (*link, error) {
	uid := unique.ID(ctx, ids)
	val, ok := as.links.Load(uid)
	if !ok {
		if val, ok := as.linkErrors.Load(uid); ok {
			err := val.(error)
			if err != nil && !errors.IsCanceled(err) {
				return nil, errLinkFailed.WithCause(err)
			}
		}
		return nil, errNotLinked.WithAttributes("application_uid", uid)
	}
	link := val.(*link)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-link.connReady:
		return link, nil
	}
}

func (l *link) observeSubscribe(correlationID string, sub *io.Subscription) {
	registerSubscribe(events.ContextWithCorrelationID(l.ctx, correlationID), sub)
	log.FromContext(sub.Context()).Debug("Subscribed")
}

func (l *link) observeUnsubscribe(correlationID string, sub *io.Subscription) {
	registerUnsubscribe(events.ContextWithCorrelationID(l.ctx, correlationID), sub)
	log.FromContext(sub.Context()).Debug("Unsubscribed")
}

func (l *link) run() {
	subscribers := make(map[*io.Subscription]string)
	defer func() {
		for sub, correlationID := range subscribers {
			l.observeUnsubscribe(correlationID, sub)
		}
	}()
	for {
		select {
		case <-l.ctx.Done():
			return
		case sub := <-l.subscribeCh:
			correlationID := fmt.Sprintf("as:subscriber:%s", events.NewCorrelationID())
			subscribers[sub] = correlationID
			l.observeSubscribe(correlationID, sub)
		case sub := <-l.unsubscribeCh:
			if correlationID, ok := subscribers[sub]; ok {
				delete(subscribers, sub)
				l.observeUnsubscribe(correlationID, sub)
			}
		case up := <-l.upCh:
			for sub := range subscribers {
				if err := sub.SendUp(up.Context, up.ApplicationUp); err != nil {
					log.FromContext(sub.Context()).WithError(err).Warn("Send upstream message failed")
				}
			}
		}
	}
}

// SendUp processes the given uplink and then sends it to the application frontends.
func (as *ApplicationServer) SendUp(ctx context.Context, up *ttnpb.ApplicationUp) error {
	link, err := as.getLink(ctx, up.ApplicationIdentifiers)
	if err != nil {
		return err
	}
	return link.sendUp(ctx, up, func() error { return nil })
}

func (l *link) sendUp(ctx context.Context, up *ttnpb.ApplicationUp, ack func() error) error {
	ctx = events.ContextWithCorrelationID(ctx, append(up.CorrelationIDs, fmt.Sprintf("as:up:%s", events.NewCorrelationID()))...)
	up.CorrelationIDs = events.CorrelationIDsFromContext(ctx)
	registerReceiveUp(ctx, up, l.connName)

	now := time.Now().UTC()
	up.ReceivedAt = &now

	pass, handleUpErr := l.handleUp(ctx, up, l)
	if err := ack(); err != nil {
		return err
	}
	if !pass {
		return nil
	}

	if handleUpErr != nil {
		log.FromContext(ctx).WithError(handleUpErr).Warn("Failed to process upstream message")
		registerDropUp(ctx, up, handleUpErr)
		return nil
	}
	ctxUp := &io.ContextualApplicationUp{
		Context:       ctx,
		ApplicationUp: up,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case l.upCh <- ctxUp:
	}
	registerForwardUp(ctx, up)
	return nil
}

// GetLinkTime returns the timestamp when the link got established.
func (l *link) GetLinkTime() time.Time { return time.Unix(0, l.linkTime) }

// GetUpStats returns the upstream statistics.
func (l *link) GetUpStats() (total uint64, t time.Time, ok bool) {
	total = atomic.LoadUint64(&l.ups)
	if ok = total > 0; ok {
		t = time.Unix(0, atomic.LoadInt64(&l.lastUpTime))
	}
	return
}

// GetDownlinkStats returns the downlink statistics.
func (l *link) GetDownlinkStats() (total uint64, t time.Time, ok bool) {
	total = atomic.LoadUint64(&l.downlinks)
	if ok = total > 0; ok {
		t = time.Unix(0, atomic.LoadInt64(&l.lastDownlinkTime))
	}
	return
}
