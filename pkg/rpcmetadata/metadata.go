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

// Package rpcmetadata contains utilities for transporting common request metadata over gRPC.
package rpcmetadata

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

// MD contains the TTN metadata fields
type MD struct {
	ID             string
	AuthType       string
	AuthValue      string
	ServiceType    string
	ServiceVersion string
	NetAddress     string
	AllowInsecure  bool

	// Limit is the limit of elements to display per-page.
	Limit uint64

	// Page is the page of elements to display.
	Page uint64

	// Host is the hostname the request is directed to.
	Host string

	// URI is the URI the request is directed to.
	URI string
}

// RequireTransportSecurity returns true if authentication is configured
func (m MD) RequireTransportSecurity() bool {
	if m.AuthType == "" || m.AuthValue == "" {
		return false
	}
	return !m.AllowInsecure
}

// ToMetadata puts the TTN metadata fields in a metadata.MD
func (m MD) ToMetadata() metadata.MD {
	var pairs []string
	if m.ID != "" {
		pairs = append(pairs, "id", m.ID)
	}
	if m.ServiceType != "" {
		pairs = append(pairs, "service-type", m.ServiceType)
	}
	if m.ServiceVersion != "" {
		pairs = append(pairs, "service-version", m.ServiceVersion)
	}
	if m.NetAddress != "" {
		pairs = append(pairs, "net-address", m.NetAddress)
	}
	if m.Host != "" {
		pairs = append(pairs, "host", m.Host)
	}
	if m.URI != "" {
		pairs = append(pairs, "uri", m.URI)
	}
	if m.Limit != 0 {
		pairs = append(pairs, "limit", strconv.FormatUint(m.Limit, 10))
	}
	if m.Page != 0 {
		pairs = append(pairs, "page", strconv.FormatUint(m.Page, 10))
	}
	return metadata.Pairs(pairs...)
}

// ToOutgoingContext puts the TTN metadata fields in an outgoing context.Context
func (m MD) ToOutgoingContext(ctx context.Context) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = metadata.Join(m.ToMetadata(), md)
	return metadata.NewOutgoingContext(ctx, md)
}

// ToIncomingContext puts the TTN metadata fields in an incoming context.Context
func (m MD) ToIncomingContext(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	md = metadata.Join(m.ToMetadata(), md)
	return metadata.NewIncomingContext(ctx, md)
}

// FromMetadata returns the TTN metadata from metadata.MD
func FromMetadata(md metadata.MD) (m MD) {
	if id, ok := md["id"]; ok && len(id) > 0 {
		m.ID = id[0]
	}
	if authorization, ok := md["authorization"]; ok && len(authorization) > 0 {
		if parts := strings.SplitN(authorization[0], " ", 2); len(parts) == 2 {
			m.AuthType, m.AuthValue = parts[0], parts[1]
		}
	}
	if serviceType, ok := md["service-type"]; ok && len(serviceType) > 0 {
		m.ServiceType = serviceType[0]
	}
	if serviceVersion, ok := md["service-version"]; ok && len(serviceVersion) > 0 {
		m.ServiceVersion = serviceVersion[0]
	}
	if netAddress, ok := md["net-address"]; ok && len(netAddress) > 0 {
		m.NetAddress = netAddress[0]
	}
	if host, ok := md["host"]; ok && len(host) > 0 {
		m.Host = host[0]
	}
	if uri, ok := md["uri"]; ok && len(uri) > 0 {
		m.URI = uri[0]
	}
	if limit, ok := md["limit"]; ok && len(limit) > 0 {
		m.Limit, _ = strconv.ParseUint(limit[0], 10, 64)
	}
	if page, ok := md["page"]; ok && len(page) > 0 {
		m.Page, _ = strconv.ParseUint(page[0], 10, 64)
	}
	return
}

// FromOutgoingContext returns the TTN metadata from the outgoing context ctx.
func FromOutgoingContext(ctx context.Context) (m MD) {
	md, _ := metadata.FromOutgoingContext(ctx)
	return FromMetadata(md)
}

// FromIncomingContext returns the TTN metadata from the incoming context ctx.
func FromIncomingContext(ctx context.Context) (m MD) {
	md, _ := metadata.FromIncomingContext(ctx)
	return FromMetadata(md)
}
