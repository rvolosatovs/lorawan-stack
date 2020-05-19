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

package gatewayconfigurationserver

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/web"
	"go.thethings.network/lorawan-stack/v3/pkg/webmiddleware"
	"google.golang.org/grpc"
)

// Server implements the CUPS endpoints used by The Things Gateway.
type Server struct {
	component *component.Component

	ttgConfig TheThingsGatewayConfig

	registry ttnpb.GatewayRegistryClient
	auth     func(context.Context) grpc.CallOption
}

func (s *Server) getRegistry(ctx context.Context, ids *ttnpb.GatewayIdentifiers) (ttnpb.GatewayRegistryClient, error) {
	if s.registry != nil {
		return s.registry, nil
	}
	cc, err := s.component.GetPeerConn(ctx, ttnpb.ClusterRole_ENTITY_REGISTRY, ids)
	if err != nil {
		return nil, err
	}
	return ttnpb.NewGatewayRegistryClient(cc), nil
}

func (s *Server) getAuth(ctx context.Context) grpc.CallOption {
	if s.auth != nil {
		return s.auth(ctx)
	}
	return s.component.WithClusterAuth()
}

// Option configures the Server.
type Option func(s *Server)

// WithRegistry overrides the Server's gateway registry.
func WithRegistry(registry ttnpb.GatewayRegistryClient) Option {
	return func(s *Server) {
		s.registry = registry
	}
}

// WithAuth overrides the Server's auth func.
func WithAuth(auth func(ctx context.Context) grpc.CallOption) Option {
	return func(s *Server) {
		s.auth = auth
	}
}

// WithTheThingsGatewayConfig overrides the Server's configuration for The Things Gateway.
func WithTheThingsGatewayConfig(config TheThingsGatewayConfig) Option {
	return func(s *Server) {
		s.ttgConfig = config
	}
}

// RegisterRoutes implements the web.Registerer interface.
func (s *Server) RegisterRoutes(server *web.Server) {
	router := server.Prefix("/api/v2/").Subrouter()
	router.Use(
		mux.MiddlewareFunc(webmiddleware.Namespace("gatewayconfigurationserver/v2")),
		rewriteAuthorization,
		mux.MiddlewareFunc(webmiddleware.Metadata("Authorization")),
		validateAndFillIDs,
	)

	router.HandleFunc("/gateways/{gateway_id}", s.handleGetGateway).Methods(http.MethodGet)
	router.HandleFunc("/frequency-plans/{frequency_plan_id}", s.handleGetFrequencyPlan).Methods(http.MethodGet)
}

// New returns a new v2 GCS on top of the given gateway registry.
func New(c *component.Component, options ...Option) *Server {
	s := &Server{
		component: c,
	}
	for _, opt := range options {
		opt(s)
	}
	c.RegisterWeb(s)
	return s
}
