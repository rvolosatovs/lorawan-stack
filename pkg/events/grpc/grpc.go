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

// Package grpc contains an implementation of the EventsServer, which is used to
// stream all events published for a set of identifiers.
package grpc

import (
	"context"
	"runtime"

	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights/rightsutil"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware/warning"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const workersPerCPU = 2

// NewEventsServer returns a new EventsServer on the given PubSub.
func NewEventsServer(ctx context.Context, pubsub events.PubSub) *EventsServer {
	srv := &EventsServer{
		ctx:    ctx,
		pubsub: pubsub,
		events: make(events.Channel, 256),
		filter: events.NewIdentifierFilter(),
	}

	hander := events.ContextHandler(ctx, srv.events)
	pubsub.Subscribe("**", hander)
	go func() {
		<-ctx.Done()
		pubsub.Unsubscribe("**", hander)
		close(srv.events)
	}()

	for i := 0; i < runtime.NumCPU()*workersPerCPU; i++ {
		go func() {
			for evt := range srv.events {
				proto, err := events.Proto(evt)
				if err != nil {
					return
				}
				srv.filter.Notify(marshaledEvent{
					Event: evt,
					proto: proto,
				})
			}
		}()
	}

	return srv
}

type marshaledEvent struct {
	events.Event
	proto *ttnpb.Event
}

// EventsServer streams events from a PubSub over gRPC.
type EventsServer struct {
	ctx    context.Context
	pubsub events.PubSub
	events events.Channel
	filter events.IdentifierFilter
}

var (
	evtStreamStart = events.Define("events.stream.start", "start event stream")
	evtStreamStop  = events.Define("events.stream.stop", "stop event stream")
)

// Stream implements the EventsServer interface.
func (srv *EventsServer) Stream(req *ttnpb.StreamEventsRequest, stream ttnpb.Events_StreamServer) error {
	ctx := stream.Context()

	if len(req.Identifiers) == 0 {
		return nil
	}

	if err := rights.RequireAny(ctx, req.Identifiers...); err != nil {
		return err
	}

	ch := make(events.Channel, 8)
	handler := events.ContextHandler(ctx, ch)
	srv.filter.Subscribe(ctx, req, handler)
	defer srv.filter.Unsubscribe(ctx, req, handler)

	if req.Tail > 0 || req.After != nil {
		warning.Add(ctx, "Historical events not implemented")
	}

	if err := stream.SendHeader(metadata.MD{}); err != nil {
		return err
	}

	evtStreamStart := evtStreamStart(ctx, req, req)
	srv.pubsub.Publish(evtStreamStart)

	evtStreamStartVisible, err := rightsutil.EventIsVisible(ctx, evtStreamStart)
	if err != nil {
		return err
	}
	if !evtStreamStartVisible {
		evt, err := events.Proto(evtStreamStart)
		if err != nil {
			return err
		}
		if err := stream.Send(evt); err != nil {
			return err
		}
	}
	defer srv.pubsub.Publish(evtStreamStop(ctx, req, req))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-ch:
			isVisible, err := rightsutil.EventIsVisible(ctx, evt)
			if err != nil {
				return err
			}
			if !isVisible {
				continue
			}
			marshaled := evt.(marshaledEvent)
			if err := stream.Send(marshaled.proto); err != nil {
				return err
			}
		}
	}
}

// Roles implements rpcserver.Registerer.
func (srv *EventsServer) Roles() []ttnpb.ClusterRole {
	return nil
}

// RegisterServices implements rpcserver.Registerer.
func (srv *EventsServer) RegisterServices(s *grpc.Server) {
	ttnpb.RegisterEventsServer(s, srv)
}

// RegisterHandlers implements rpcserver.Registerer.
func (srv *EventsServer) RegisterHandlers(s *grpc_runtime.ServeMux, conn *grpc.ClientConn) {
	ttnpb.RegisterEventsHandler(srv.ctx, s, conn)
}
