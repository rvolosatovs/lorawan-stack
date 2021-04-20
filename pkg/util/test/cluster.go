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

package test

import (
	"context"
	"reflect"

	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcclient"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcserver"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"google.golang.org/grpc"
)

// MockPeer is a mock cluster.Peer used for testing.
type MockPeer struct {
	NameFunc    func() string
	ConnFunc    func() (*grpc.ClientConn, error)
	HasRoleFunc func(ttnpb.ClusterRole) bool
	RolesFunc   func() []ttnpb.ClusterRole
	TagsFunc    func() map[string]string
}

// Name calls NameFunc if set and panics otherwise.
func (m MockPeer) Name() string {
	if m.NameFunc == nil {
		panic("Name called, but not set")
	}
	return m.NameFunc()
}

// Conn calls ConnFunc if set and panics otherwise.
func (m MockPeer) Conn() (*grpc.ClientConn, error) {
	if m.ConnFunc == nil {
		panic("Conn called, but not set")
	}
	return m.ConnFunc()
}

// HasRole calls HasRoleFunc if set and panics otherwise.
func (m MockPeer) HasRole(r ttnpb.ClusterRole) bool {
	if m.HasRoleFunc == nil {
		panic("HasRole called, but not set")
	}
	return m.HasRoleFunc(r)
}

// Roles calls RolesFunc if set and panics otherwise.
func (m MockPeer) Roles() []ttnpb.ClusterRole {
	if m.RolesFunc == nil {
		panic("Roles called, but not set")
	}
	return m.RolesFunc()
}

// Tags calls TagsFunc if set and panics otherwise.
func (m MockPeer) Tags() map[string]string {
	if m.TagsFunc == nil {
		panic("Tags called, but not set")
	}
	return m.TagsFunc()
}

// NewGRPCServerPeer creates a new MockPeer with ConnFunc, which always returns the same loopback connection to the server itself.
// srv is the implementation of the gRPC interface.
// registrators represents a slice of functions, which register the gRPC interface implementation at a gRPC server.
func NewGRPCServerPeer(ctx context.Context, srv interface{}, registrators ...interface{}) (*MockPeer, error) {
	grpcSrv := rpcserver.New(ctx).Server
	for _, r := range registrators {
		reflect.ValueOf(r).Call([]reflect.Value{
			reflect.ValueOf(grpcSrv),
			reflect.ValueOf(srv),
		})
	}
	conn, err := rpcserver.StartLoopback(ctx, grpcSrv, rpcclient.DefaultDialOptions(ctx)...)
	if err != nil {
		return nil, err
	}
	return &MockPeer{
		ConnFunc: func() (*grpc.ClientConn, error) { return conn, nil },
	}, nil
}

// MockCluster is a mock cluster.Cluster used for testing.
type MockCluster struct {
	JoinFunc               func() error
	LeaveFunc              func() error
	GetPeersFunc           func(ctx context.Context, role ttnpb.ClusterRole) ([]cluster.Peer, error)
	GetPeerFunc            func(ctx context.Context, role ttnpb.ClusterRole, ids cluster.EntityIdentifiers) (cluster.Peer, error)
	ClaimIDsFunc           func(ctx context.Context, ids cluster.EntityIdentifiers) error
	UnclaimIDsFunc         func(ctx context.Context, ids cluster.EntityIdentifiers) error
	TLSFunc                func() bool
	AuthFunc               func() grpc.CallOption
	WithVerifiedSourceFunc func(ctx context.Context) context.Context
}

// Join calls JoinFunc if set and panics otherwise.
func (m MockCluster) Join() error {
	if m.JoinFunc == nil {
		panic("Join called, but not set")
	}
	return m.JoinFunc()
}

// Leave calls LeaveFunc if set and panics otherwise.
func (m MockCluster) Leave() error {
	if m.LeaveFunc == nil {
		panic("Leave called, but not set")
	}
	return m.LeaveFunc()
}

// GetPeers calls GetPeersFunc if set and panics otherwise.
func (m MockCluster) GetPeers(ctx context.Context, role ttnpb.ClusterRole) ([]cluster.Peer, error) {
	if m.GetPeersFunc == nil {
		panic("GetPeers called, but not set")
	}
	return m.GetPeersFunc(ctx, role)
}

// GetPeer calls GetPeerFunc if set and panics otherwise.
func (m MockCluster) GetPeer(ctx context.Context, role ttnpb.ClusterRole, ids cluster.EntityIdentifiers) (cluster.Peer, error) {
	if m.GetPeerFunc == nil {
		panic("GetPeer called, but not set")
	}
	return m.GetPeerFunc(ctx, role, ids)
}

// GetPeerConn calls GetPeer and then Conn.
func (m MockCluster) GetPeerConn(ctx context.Context, role ttnpb.ClusterRole, ids cluster.EntityIdentifiers) (*grpc.ClientConn, error) {
	peer, err := m.GetPeer(ctx, role, ids)
	if err != nil {
		return nil, err
	}
	return peer.Conn()
}

// ClaimIDs calls ClaimIDsFunc if set and panics otherwise.
func (m MockCluster) ClaimIDs(ctx context.Context, ids cluster.EntityIdentifiers) error {
	if m.ClaimIDsFunc == nil {
		panic("ClaimIDs called, but not set")
	}
	return m.ClaimIDsFunc(ctx, ids)
}

// UnclaimIDs calls UnclaimIDsFunc if set and panics otherwise.
func (m MockCluster) UnclaimIDs(ctx context.Context, ids cluster.EntityIdentifiers) error {
	if m.UnclaimIDsFunc == nil {
		panic("UnclaimIDs called, but not set")
	}
	return m.UnclaimIDsFunc(ctx, ids)
}

// TLS calls TLSFunc if set and panics otherwise.
func (m MockCluster) TLS() bool {
	if m.TLSFunc == nil {
		panic("TLS called, but not set")
	}
	return m.TLSFunc()
}

// Auth calls AuthFunc if set and panics otherwise.
func (m MockCluster) Auth() grpc.CallOption {
	if m.AuthFunc == nil {
		panic("Auth called, but not set")
	}
	return m.AuthFunc()
}

// WithVerifiedSource calls WithVerifiedSourceFunc if set and panics otherwise.
func (m MockCluster) WithVerifiedSource(ctx context.Context) context.Context {
	if m.WithVerifiedSourceFunc == nil {
		panic("WithVerifiedSource called, but not set")
	}
	return m.WithVerifiedSourceFunc(ctx)
}

type ClusterAuthRequest struct {
	Response chan<- grpc.CallOption
}

func MakeClusterAuthChFunc(reqCh chan<- ClusterAuthRequest) func() grpc.CallOption {
	return func() grpc.CallOption {
		respCh := make(chan grpc.CallOption)
		reqCh <- ClusterAuthRequest{
			Response: respCh,
		}
		return <-respCh
	}
}

type ClusterGetPeerResponse struct {
	Peer  cluster.Peer
	Error error
}

type ClusterGetPeerRequest struct {
	Context     context.Context
	Role        ttnpb.ClusterRole
	Identifiers cluster.EntityIdentifiers
	Response    chan<- ClusterGetPeerResponse
}

func MakeClusterGetPeerChFunc(reqCh chan<- ClusterGetPeerRequest) func(context.Context, ttnpb.ClusterRole, cluster.EntityIdentifiers) (cluster.Peer, error) {
	return func(ctx context.Context, role ttnpb.ClusterRole, ids cluster.EntityIdentifiers) (cluster.Peer, error) {
		respCh := make(chan ClusterGetPeerResponse)
		reqCh <- ClusterGetPeerRequest{
			Context:     ctx,
			Role:        role,
			Identifiers: ids,
			Response:    respCh,
		}
		resp := <-respCh
		return resp.Peer, resp.Error
	}
}

type ClusterJoinRequest struct {
	Response chan<- error
}

func MakeClusterJoinChFunc(reqCh chan<- ClusterJoinRequest) func() error {
	return func() error {
		respCh := make(chan error)
		reqCh <- ClusterJoinRequest{
			Response: respCh,
		}
		return <-respCh
	}
}

func ClusterJoinNilFunc() error { return nil }

func AssertClusterAuthRequest(ctx context.Context, reqCh <-chan ClusterAuthRequest, resp grpc.CallOption) bool {
	t := MustTFromContext(ctx)
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("Timed out while waiting for Cluster.Auth to be called")
		return false

	case req := <-reqCh:
		t.Log("Cluster.Auth called")
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for Cluster.Auth response to be processed")
			return false

		case req.Response <- resp:
			return true
		}
	}
}

func AssertClusterGetPeerRequest(ctx context.Context, reqCh <-chan ClusterGetPeerRequest, assert func(ctx, reqCtx context.Context, role ttnpb.ClusterRole, ids cluster.EntityIdentifiers) (ClusterGetPeerResponse, bool)) bool {
	t := MustTFromContext(ctx)
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("Timed out while waiting for Cluster.GetPeer to be called")
		return false

	case req := <-reqCh:
		t.Log("Cluster.GetPeer called")
		resp, ok := assert(ctx, req.Context, req.Role, req.Identifiers)
		if !ok {
			t.Error("Cluster.GetPeer request assertion failed")
			return false
		}
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for Cluster.GetPeer response to be processed")
			return false

		case req.Response <- resp:
			return true
		}
	}
}
