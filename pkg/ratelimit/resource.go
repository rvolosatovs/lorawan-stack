// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package ratelimit

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
)

// Resource represents an entity on which rate limits apply.
type Resource interface {
	// Key represents the unique identifier for the resource.
	Key() string
	// Classes represents the rate limiting classes for the resource. A resource may
	// belong in multiple classes. The limiter instance will limit access to the resource
	// based on the limits of the first class that matches the configuration. If no
	// class is matched, the limiter instance can fallback to a default rate limit.
	Classes() []string
}

type resource struct {
	key     string
	classes []string
}

func (r *resource) Key() string       { return r.key }
func (r *resource) Classes() []string { return r.classes }

// httpRequestResource represents an HTTP request. Avoid using directly, use HTTPMiddleware instead.
func httpRequestResource(r *http.Request, class string) Resource {
	return &resource{
		key:     fmt.Sprintf("%s:ip:%s:url:%s", class, httpRemoteIP(r), r.URL.Path),
		classes: []string{class, "http"},
	}
}

// echoRequestResource represents an HTTP request. Avoid using directly, use EchoMiddleware instead.
func echoRequestResource(c echo.Context, class string) Resource {
	return &resource{
		key:     fmt.Sprintf("%s:ip:%s:url:%s", class, c.RealIP(), c.Request().URL.Path),
		classes: []string{class, "http"},
	}
}

// grpcMethodResource represents a gRPC request.
func grpcMethodResource(ctx context.Context, fullMethod string, req interface{}) Resource {
	key := fmt.Sprintf("grpc:method:%s:%s", fullMethod, grpcEntityFromRequest(ctx, req))
	if authTokenID := grpcAuthTokenID(ctx); authTokenID != "" {
		key = fmt.Sprintf("%s:token:%s", key, authTokenID)
	}
	return &resource{
		key:     key,
		classes: []string{fmt.Sprintf("grpc:method:%s", fullMethod), "grpc:method"},
	}
}

// grpcStreamAcceptResource represents a new gRPC server stream.
func grpcStreamAcceptResource(ctx context.Context, fullMethod string) Resource {
	key := fmt.Sprintf("grpc:stream:accept:%s", fullMethod)
	if authTokenID := grpcAuthTokenID(ctx); authTokenID != "" {
		key = fmt.Sprintf("%s:token:%s", key, authTokenID)
	}
	return &resource{
		key:     key,
		classes: []string{fmt.Sprintf("grpc:stream:accept:%s", fullMethod), "grpc:stream:accept"},
	}
}

// grpcStreamUpResource represents client messages for a gRPC server stream.
func grpcStreamUpResource(ctx context.Context, fullMethod string) Resource {
	return &resource{
		key:     fmt.Sprintf("grpc:stream:up:%s:streamID:%s", fullMethod, events.NewCorrelationID()),
		classes: []string{fmt.Sprintf("grpc:stream:up:%s", fullMethod), "grpc:stream:up"},
	}
}

// GatewayUpResource represents uplink traffic from a gateway.
func GatewayUpResource(ctx context.Context, ids ttnpb.GatewayIdentifiers) Resource {
	return &resource{
		key:     fmt.Sprintf("gs:up:gtw:%s", unique.ID(ctx, ids)),
		classes: []string{"gs:up"},
	}
}

// GatewayAcceptMQTTConnectionResource represents a new MQTT gateway connection from a remote address.
func GatewayAcceptMQTTConnectionResource(remoteAddr string) Resource {
	remoteIP := remoteAddr
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		remoteIP = host
	}
	return &resource{
		key:     fmt.Sprintf("gs:accept:mqtt:ip:%s", remoteIP),
		classes: []string{"gs:accept:mqtt"},
	}
}

// GatewayUDPTrafficResource represents UDP gateway traffic from a remote IP address.
func GatewayUDPTrafficResource(addr *net.UDPAddr) Resource {
	return &resource{
		key:     fmt.Sprintf("gs:up:udp:ip:%s", addr.IP.String()),
		classes: []string{"gs:up:udp", "gs:up"},
	}
}

// ApplicationAcceptMQTTConnectionResource represents a new MQTT client connection from a remote IP address.
func ApplicationAcceptMQTTConnectionResource(remoteAddr string) Resource {
	remoteIP := remoteAddr
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		remoteIP = host
	}
	return &resource{
		key:     fmt.Sprintf("as:accept:mqtt:ip:%s", remoteIP),
		classes: []string{"as:accept:mqtt"},
	}
}

// ApplicationMQTTDownResource represents downlink traffic for an application from an MQTT client.
func ApplicationMQTTDownResource(ctx context.Context, ids ttnpb.ApplicationIdentifiers, authTokenID string) Resource {
	key := fmt.Sprintf("as:down:mqtt:app:%s", unique.ID(ctx, ids))
	if authTokenID != "" {
		key = fmt.Sprintf("%s:token:%s", key, authTokenID)
	}
	return &resource{
		key:     key,
		classes: []string{"as:down:mqtt"},
	}
}

// ApplicationWebhooksDownResource represents downlink traffic for an application from a webhook.
func ApplicationWebhooksDownResource(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, authTokenID string) Resource {
	key := fmt.Sprintf("as:down:web:dev:%s", unique.ID(ctx, ids))
	if authTokenID != "" {
		key = fmt.Sprintf("%s:token:%s", key, authTokenID)
	}
	return &resource{
		key:     key,
		classes: []string{"as:down:web"},
	}
}

// NewCustomResource returns a new resource. It is used internally by other components.
func NewCustomResource(key string, classes ...string) Resource {
	return &resource{key, classes}
}
