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

package mock

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/scheduling"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"go.thethings.network/lorawan-stack/pkg/util/test"
)

type server struct {
	store          *frequencyplans.Store
	connectionsCh  chan *io.Connection
	downlinkClaims sync.Map
}

// Server represents a testing io.Server.
type Server interface {
	io.Server

	HasDownlinkClaim(context.Context, ttnpb.GatewayIdentifiers) bool
	Connections() <-chan *io.Connection
}

// NewServer instantiates a new Server.
func NewServer() Server {
	return &server{
		store:         frequencyplans.NewStore(test.FrequencyPlansFetcher),
		connectionsCh: make(chan *io.Connection, 10),
	}
}

// FillContext implements io.Server.
func (s *server) FillGatewayContext(ctx context.Context, ids ttnpb.GatewayIdentifiers) (context.Context, ttnpb.GatewayIdentifiers, error) {
	if ids.IsZero() {
		return nil, ttnpb.GatewayIdentifiers{}, errors.New("the identifiers are zero")
	}
	if ids.GatewayID != "" {
		return ctx, ids, nil
	}
	ids.GatewayID = fmt.Sprintf("eui-%v", strings.ToLower(ids.EUI.String()))
	return ctx, ids, nil
}

// Connect implements io.Server.
func (s *server) Connect(ctx context.Context, protocol string, ids ttnpb.GatewayIdentifiers) (*io.Connection, error) {
	if err := rights.RequireGateway(ctx, ids, ttnpb.RIGHT_GATEWAY_LINK); err != nil {
		return nil, err
	}
	gtw := &ttnpb.Gateway{GatewayIdentifiers: ids}
	fp, err := s.GetFrequencyPlan(ctx, ids)
	if err != nil {
		return nil, err
	}
	scheduler, err := scheduling.NewScheduler(ctx, fp, true)
	if err != nil {
		return nil, err
	}
	conn := io.NewConnection(ctx, protocol, gtw, fp, scheduler)
	select {
	case s.connectionsCh <- conn:
	default:
	}
	return conn, nil
}

// GetFrequencyPlan implements io.Server.
func (s *server) GetFrequencyPlan(ctx context.Context, _ ttnpb.GatewayIdentifiers) (*frequencyplans.FrequencyPlan, error) {
	return s.store.GetByID(test.EUFrequencyPlanID)
}

// ClaimDownlink implements io.Server.
func (s *server) ClaimDownlink(ctx context.Context, ids ttnpb.GatewayIdentifiers) error {
	s.downlinkClaims.Store(unique.ID(ctx, ids), true)
	return nil
}

// UnclaimDownlink implements io.Server.
func (s *server) UnclaimDownlink(ctx context.Context, ids ttnpb.GatewayIdentifiers) error {
	s.downlinkClaims.Delete(unique.ID(ctx, ids))
	return nil
}

func (s *server) HasDownlinkClaim(ctx context.Context, ids ttnpb.GatewayIdentifiers) bool {
	_, ok := s.downlinkClaims.Load(unique.ID(ctx, ids))
	return ok
}

func (s *server) Connections() <-chan *io.Connection {
	return s.connectionsCh
}
