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

package grpc

import (
	"context"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"google.golang.org/grpc/peer"
)

type impl struct {
	server io.Server
}

// New returns a new gRPC frontend.
func New(server io.Server) ttnpb.AppAsServer {
	return &impl{server}
}

var errConnect = errors.Define("connect", "failed to connect application `{application_uid}`")

func (s *impl) Subscribe(ids *ttnpb.ApplicationIdentifiers, stream ttnpb.AppAs_SubscribeServer) error {
	ctx := log.NewContextWithField(stream.Context(), "namespace", "applicationserver/io/grpc")

	if err := rights.RequireApplication(ctx, *ids, ttnpb.RIGHT_APPLICATION_TRAFFIC_READ); err != nil {
		return err
	}

	if peer, ok := peer.FromContext(ctx); ok {
		ctx = log.NewContextWithField(ctx, "remote_addr", peer.Addr.String())
	}
	uid := unique.ID(ctx, ids)
	ctx = log.NewContextWithField(ctx, "application_uid", uid)
	logger := log.FromContext(ctx)

	sub, err := s.server.Subscribe(ctx, "grpc", *ids)
	if err != nil {
		logger.WithError(err).Warn("Failed to connect")
		return errConnect.WithCause(err).WithAttributes("application_uid", uid)
	}
	logger.Info("Subscribed")
	defer logger.Info("Unsubscribed")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case up := <-sub.Up():
			if err := stream.Send(up); err != nil {
				logger.WithError(err).Warn("Failed to send message")
				sub.Disconnect(err)
				return err
			}
		}
	}
}

func clean(items []*ttnpb.ApplicationDownlink) []*ttnpb.ApplicationDownlink {
	res := make([]*ttnpb.ApplicationDownlink, 0, len(items))
	for _, item := range items {
		res = append(res, &ttnpb.ApplicationDownlink{
			FPort:          item.FPort,
			FRMPayload:     item.FRMPayload,
			DecodedPayload: item.DecodedPayload,
			Confirmed:      item.Confirmed,
		})
	}
	return res
}

func (s *impl) DownlinkQueuePush(ctx context.Context, req *ttnpb.DownlinkQueueRequest) (*pbtypes.Empty, error) {
	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_TRAFFIC_DOWN_WRITE); err != nil {
		return nil, err
	}
	if err := s.server.DownlinkQueuePush(ctx, req.EndDeviceIdentifiers, clean(req.Downlinks)); err != nil {
		return nil, err
	}
	return ttnpb.Empty, nil
}

func (s *impl) DownlinkQueueReplace(ctx context.Context, req *ttnpb.DownlinkQueueRequest) (*pbtypes.Empty, error) {
	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_TRAFFIC_DOWN_WRITE); err != nil {
		return nil, err
	}
	if err := s.server.DownlinkQueueReplace(ctx, req.EndDeviceIdentifiers, clean(req.Downlinks)); err != nil {
		return nil, err
	}
	return ttnpb.Empty, nil
}

func (s *impl) DownlinkQueueList(ctx context.Context, ids *ttnpb.EndDeviceIdentifiers) (*ttnpb.ApplicationDownlinks, error) {
	if err := rights.RequireApplication(ctx, ids.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_TRAFFIC_READ); err != nil {
		return nil, err
	}
	items, err := s.server.DownlinkQueueList(ctx, *ids)
	if err != nil {
		return nil, err
	}
	return &ttnpb.ApplicationDownlinks{
		Downlinks: items,
	}, nil
}
