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

// Package cluster defines an interface for clustering network components and provides a simple implementation.
package cluster

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"os"
	"strings"

	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/random"
	"go.thethings.network/lorawan-stack/pkg/rpcclient"
	"go.thethings.network/lorawan-stack/pkg/rpcserver"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

// Cluster interface that is implemented by all different clustering implementations.
type Cluster interface {
	// Join the cluster.
	Join() error
	// Leave the cluster.
	Leave() error

	// GetPeers returns peers with the given role.
	GetPeers(ctx context.Context, role ttnpb.ClusterRole) ([]Peer, error)
	// GetPeer returns a peer with the given role, and a responsibility for the
	// given identifiers. If the identifiers are nil, this function returns a random
	// peer from the list that would be returned by GetPeers.
	GetPeer(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) (Peer, error)
	// GetPeerConn returns the gRPC client connection of a peer, if the peer is available as
	// as per GetPeer.
	GetPeerConn(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) (*grpc.ClientConn, error)

	// ClaimIDs can be used to indicate that the current peer takes
	// responsibility for entities identified by ids.
	// Claiming an already claimed ID will transfer the claim (without notifying
	// the previous holder). Releasing a non-claimed ID is a no-op. An error may
	// only be returned if the claim/unclaim couldn't be communicated with the cluster.
	ClaimIDs(ctx context.Context, ids ttnpb.Identifiers) error
	// UnclaimIDs can be used to indicate that the current peer
	// releases responsibility for entities identified by ids.
	// The specified context ctx may already be done before calling this function.
	UnclaimIDs(ctx context.Context, ids ttnpb.Identifiers) error

	// TLS returns whether the cluster uses TLS for cluster connections.
	TLS() bool
	// Auth returns a gRPC CallOption that can be used to identify the component within the cluster.
	Auth() grpc.CallOption
	// WithVerifiedSource verifies if the caller providing this context is a component from the cluster, and returns a
	// new context with that information.
	WithVerifiedSource(context.Context) context.Context
}

// Option to apply at cluster initialization.
type Option interface {
	apply(*cluster)
}

type optionFunc func(*cluster)

func (f optionFunc) apply(c *cluster) { f(c) }

// WithConn bypasses the standard mechanism for connecting to the "self" peer.
func WithConn(conn *grpc.ClientConn) Option {
	return optionFunc(func(c *cluster) {
		c.self.conn = conn
	})
}

// WithServices registers the given services on the "self" peer.
func WithServices(services ...rpcserver.Registerer) Option {
	return optionFunc(func(c *cluster) {
		for _, service := range services {
			if roles := service.Roles(); len(roles) > 0 {
				c.self.roles = append(c.self.roles, roles...)
			}
		}
	})
}

// WithTLSConfig sets the TLS config to use in cluster connections.
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return optionFunc(func(c *cluster) {
		c.tlsConfig = tlsConfig
	})
}

// Config represents clustering configuration.
type Config struct {
	Join              []string `name:"join" description:"Addresses of cluster peers to join"`
	Name              string   `name:"name" description:"Name of the current cluster peer (default: $HOSTNAME)"`
	Address           string   `name:"address" description:"Address to use for cluster communication"`
	IdentityServer    string   `name:"identity-server" description:"Address for the Identity Server"`
	GatewayServer     string   `name:"gateway-server" description:"Address for the Gateway Server"`
	NetworkServer     string   `name:"network-server" description:"Address for the Network Server"`
	ApplicationServer string   `name:"application-server" description:"Address for the Application Server"`
	JoinServer        string   `name:"join-server" description:"Address for the Join Server"`
	CryptoServer      string   `name:"crypto-server" description:"Address for the Crypto Server"`
	TLS               bool     `name:"tls" description:"Do cluster gRPC over TLS"`
	Keys              []string `name:"keys" description:"Keys used to communicate between components of the cluster. The first one will be used by the cluster to identify itself"`
}

// CustomNew allows you to replace the clustering implementation. New will call CustomNew if not nil.
var CustomNew func(ctx context.Context, config *Config, options ...Option) (Cluster, error)

// New instantiates a new clustering implementation.
// The basic clustering implementation allows for a cluster setup with a single-instance deployment of each component
// (GS/NS/AS/JS).
// Network operators can use their own clustering logic, which can be activated by setting the CustomNew variable.
func New(ctx context.Context, config *Config, options ...Option) (Cluster, error) {
	if CustomNew != nil {
		return CustomNew(ctx, config, options...)
	}
	return defaultNew(ctx, config, options...)
}

func defaultNew(ctx context.Context, config *Config, options ...Option) (Cluster, error) {
	c := &cluster{
		ctx:   ctx,
		tls:   config.TLS,
		peers: make(map[string]*peer),
	}

	if err := c.loadKeys(ctx, config.Keys...); err != nil {
		return nil, err
	}

	c.self = &peer{
		name:   config.Name,
		target: config.Address,
	}
	if c.self.name == "" {
		c.self.name, _ = os.Hostname()
	}
	c.peers[c.self.name] = c.self

	for _, option := range options {
		option.apply(c)
	}

	c.addPeer("is", config.IdentityServer, ttnpb.ClusterRole_ACCESS, ttnpb.ClusterRole_ENTITY_REGISTRY)
	c.addPeer("gs", config.GatewayServer, ttnpb.ClusterRole_GATEWAY_SERVER)
	c.addPeer("ns", config.NetworkServer, ttnpb.ClusterRole_NETWORK_SERVER)
	c.addPeer("as", config.ApplicationServer, ttnpb.ClusterRole_APPLICATION_SERVER)
	c.addPeer("js", config.JoinServer, ttnpb.ClusterRole_JOIN_SERVER)
	c.addPeer("cs", config.CryptoServer, ttnpb.ClusterRole_CRYPTO_SERVER)

	for _, join := range config.Join {
		c.peers[join] = &peer{
			name:   join,
			target: join,
		}
	}

	return c, nil
}

type cluster struct {
	ctx       context.Context
	tls       bool
	tlsConfig *tls.Config
	peers     map[string]*peer
	self      *peer

	keys [][]byte
}

var (
	errPeerConnection    = errors.DefineUnavailable("peer_connection", "connection to peer `{name}` on `{address}` failed")
	errPeerEmptyTarget   = errors.DefineInvalidArgument("peer_empty_target", "peer target address is empty")
	errInvalidClusterKey = errors.DefineInvalidArgument("cluster_key", "invalid cluster key")
	errInvalidKeyLength  = errors.DefineInvalidArgument("key_length", "invalid key length %d, must be 16, 24 or 32 bytes")
)

func (c *cluster) loadKeys(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		decodedKey, err := hex.DecodeString(key)
		if err != nil {
			return errInvalidClusterKey.WithCause(err)
		}
		switch len(decodedKey) {
		case 16, 24, 32:
		default:
			return errInvalidClusterKey.WithCause(errInvalidKeyLength)
		}
		c.keys = append(c.keys, decodedKey)
	}
	if c.keys == nil {
		c.keys = [][]byte{random.Bytes(32)}
		log.FromContext(ctx).WithField("key", hex.EncodeToString(c.keys[0])).Warn("No cluster key configured, generated a random one")
	}
	return nil
}

func (c *cluster) addPeer(name string, target string, roles ...ttnpb.ClusterRole) {
	if target == "" {
		return
	}
	var filteredRoles []ttnpb.ClusterRole
	for _, role := range roles {
		if !c.self.HasRole(role) {
			filteredRoles = append(filteredRoles, role)
		}
	}
	if len(filteredRoles) == 0 {
		return
	}
	c.peers[name] = &peer{
		name:   name,
		target: target,
		roles:  filteredRoles,
	}
}

func (c *cluster) Join() (err error) {
	options := rpcclient.DefaultDialOptions(c.ctx)
	if c.tls {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(c.tlsConfig)))
	} else {
		options = append(options, grpc.WithInsecure())
	}
	for _, peer := range c.peers {
		if peer.conn != nil {
			continue
		}
		peer.ctx, peer.cancel = context.WithCancel(c.ctx)
		logger := log.FromContext(c.ctx).WithFields(log.Fields(
			"target", peer.target,
			"name", peer.Name(),
			"roles", peer.Roles(),
		))
		if peer.target == "" {
			logger.Warn("Not connecting to peer, empty address.")
			peer.connErr = errPeerEmptyTarget
			continue
		}
		logger.Debug("Connecting to peer...")
		peer.conn, peer.connErr = grpc.DialContext(peer.ctx, peer.target, options...)
		if peer.connErr != nil {
			return errPeerConnection.WithCause(peer.connErr).WithAttributes("name", peer.name, "address", peer.target)
		}
		logger.Debug("Connected to peer")
	}
	return nil
}

func (c *cluster) Leave() error {
	for _, peer := range c.peers {
		if peer.conn != nil {
			if err := peer.conn.Close(); err != nil {
				return err
			}
		}
		if peer.cancel != nil {
			peer.cancel()
		}
	}
	return nil
}

func (c *cluster) GetPeers(ctx context.Context, role ttnpb.ClusterRole) ([]Peer, error) {
	var matches []Peer
	for _, peer := range c.peers {
		if !peer.HasRole(role) {
			continue
		}
		conn, err := peer.Conn()
		if err != nil {
			continue
		}
		if conn.GetState() == connectivity.Ready {
			matches = append(matches, peer)
		}
	}
	return matches, nil
}

// overridePeerRole may change the peer role depending on the identifiers.
func overridePeerRole(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) ttnpb.ClusterRole {
	switch role {
	case ttnpb.ClusterRole_GATEWAY_SERVER:
		if ids != nil && ids.EntityType() == "gateway" && ids.IDString() == PacketBrokerGatewayID.GatewayID {
			return ttnpb.ClusterRole_PACKET_BROKER_AGENT
		}
	default:
	}
	return role
}

var errPeerUnavailable = errors.DefineUnavailable("peer_unavailable", "{cluster_role} cluster peer unavailable")

func (c *cluster) GetPeer(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) (Peer, error) {
	role = overridePeerRole(ctx, role, ids)
	matches, err := c.GetPeers(ctx, role)
	if err != nil {
		return nil, err
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	// The reference cluster only has a single instance of each component, so we don't need to filter on IDs.
	return nil, errPeerUnavailable.WithAttributes("cluster_role", strings.Title(strings.Replace(role.String(), "_", " ", -1)))
}

func (c *cluster) GetPeerConn(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) (*grpc.ClientConn, error) {
	peer, err := c.GetPeer(ctx, role, ids)
	if err != nil {
		return nil, err
	}
	return peer.Conn()
}

// ClaimIDs is a no-op in the reference implementation.
// The reference cluster only has a single instance of each component, so we don't need to claim.
func (c *cluster) ClaimIDs(ctx context.Context, ids ttnpb.Identifiers) error {
	return nil
}

// UnclaimIDs is a no-op in the reference implementation.
// The reference cluster only has a single instance of each component, so we don't need to unclaim.
func (c *cluster) UnclaimIDs(ctx context.Context, ids ttnpb.Identifiers) error {
	return nil
}
