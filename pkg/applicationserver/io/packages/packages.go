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

package packages

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcserver"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"google.golang.org/grpc"
)

type server struct {
	ctx context.Context

	io       io.Server
	registry Registry

	handlers map[string]ApplicationPackageHandler
}

// Server is an application packages frontend.
type Server interface {
	rpcserver.Registerer
}

// New returns an application packages server wrapping the given registries and handlers.
func New(ctx context.Context, io io.Server, registry Registry, handlers map[string]ApplicationPackageHandler) (Server, error) {
	ctx = log.NewContextWithField(ctx, "namespace", "applicationserver/io/packages")
	s := &server{
		ctx:      ctx,
		io:       io,
		registry: registry,
		handlers: handlers,
	}
	sub, err := io.Subscribe(ctx, "applicationpackages", nil, false)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sub.Context().Done():
				return
			case up := <-sub.Up():
				if err := s.handleUp(up.Context, up.ApplicationUp); err != nil {
					log.FromContext(s.ctx).WithError(err).Warn("Failed to handle message")
				}
			}
		}
	}()
	return s, nil
}

type associationsPair struct {
	defaultAssociation *ttnpb.ApplicationPackageDefaultAssociation
	association        *ttnpb.ApplicationPackageAssociation
}

type associationsMap map[string]*associationsPair

func (s *server) findAssociations(ctx context.Context, ids ttnpb.EndDeviceIdentifiers) (associationsMap, error) {
	paths := []string{
		"data",
		"ids",
		"package_name",
	}
	associations, err := s.registry.ListAssociations(ctx, ids, paths)
	if err != nil {
		return nil, err
	}
	defaults, err := s.registry.ListDefaultAssociations(ctx, ids.ApplicationIdentifiers, paths)
	if err != nil {
		return nil, err
	}
	m := make(associationsMap)
	for _, association := range associations {
		m[association.PackageName] = &associationsPair{
			association: association,
		}
	}
	for _, defaultAssociation := range defaults {
		if pair, ok := m[defaultAssociation.PackageName]; ok {
			pair.defaultAssociation = defaultAssociation
		} else {
			m[defaultAssociation.PackageName] = &associationsPair{
				defaultAssociation: defaultAssociation,
			}
		}
	}
	return m, nil
}

func (s *server) handleUp(ctx context.Context, msg *ttnpb.ApplicationUp) error {
	ctx = log.NewContextWithField(ctx, "device_uid", unique.ID(ctx, msg.EndDeviceIdentifiers))
	associations, err := s.findAssociations(ctx, msg.EndDeviceIdentifiers)
	if err != nil {
		return err
	}
	for name, pair := range associations {
		if handler, ok := s.handlers[name]; ok {
			ctx := log.NewContextWithField(ctx, "package", name)
			err := handler.HandleUp(ctx, pair.defaultAssociation, pair.association, msg)
			if err != nil {
				return err
			}
		} else {
			return errNotImplemented.WithAttributes("name", name)
		}
	}
	return nil
}

// Roles implements the rpcserver.Registerer interface.
func (s *server) Roles() []ttnpb.ClusterRole {
	return nil
}

// RegisterServices registers the services of the registered application packages.
func (s *server) RegisterServices(gs *grpc.Server) {
	ttnpb.RegisterApplicationPackageRegistryServer(gs, s)
	for _, subsystem := range s.handlers {
		subsystem.RegisterServices(gs)
	}
}

// RegisterHandlers registers the handlers of the registered application packages.
func (s *server) RegisterHandlers(rs *runtime.ServeMux, conn *grpc.ClientConn) {
	ttnpb.RegisterApplicationPackageRegistryHandler(s.ctx, rs, conn)
	for _, subsystem := range s.handlers {
		subsystem.RegisterHandlers(rs, conn)
	}
}
