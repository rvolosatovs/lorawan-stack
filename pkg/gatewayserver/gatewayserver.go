// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package gatewayserver

import (
	"context"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/component"
	"github.com/TheThingsNetwork/ttn/pkg/frequencyplans"
	"github.com/TheThingsNetwork/ttn/pkg/gatewayserver/pool"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const sendUplinkTimeout = 5 * time.Minute

// GatewayServer implements the gateway server component.
//
// The gateway server exposes the Gs, GtwGs and NsGs services.
type GatewayServer struct {
	*component.Component

	gateways       pool.Pool
	frequencyPlans frequencyplans.Store

	nsTags []string
}

// New returns new *GatewayServer.
func New(c *component.Component, conf *Config) *GatewayServer {
	gs := &GatewayServer{
		Component: c,

		gateways:       pool.NewPool(c.Logger(), sendUplinkTimeout),
		frequencyPlans: conf.store(),

		nsTags: conf.NSTags,
	}
	c.RegisterGRPC(gs)
	return gs
}

// RegisterServices registers services provided by gs at s.
func (gs *GatewayServer) RegisterServices(s *grpc.Server) {
	ttnpb.RegisterGsServer(s, gs)
	ttnpb.RegisterGtwGsServer(s, gs)
	ttnpb.RegisterNsGsServer(s, gs)
}

// RegisterHandlers registers gRPC handlers.
func (gs *GatewayServer) RegisterHandlers(s *runtime.ServeMux, conn *grpc.ClientConn) {}

// Roles returns the roles that the gateway server fulfils
func (gs *GatewayServer) Roles() []ttnpb.PeerInfo_Role {
	return []ttnpb.PeerInfo_Role{ttnpb.PeerInfo_GATEWAY_SERVER}
}

// GetGatewayObservations returns gateway information as observed by the gateway server.
func (gs *GatewayServer) GetGatewayObservations(ctx context.Context, id *ttnpb.GatewayIdentifiers) (*ttnpb.GatewayObservations, error) {
	return gs.gateways.GetGatewayObservations(id)
}
