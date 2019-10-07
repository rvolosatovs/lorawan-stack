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

package applicationserver

import (
	"context"
	"net/http"
	"time"

	"go.thethings.network/lorawan-stack/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io/packages"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io/pubsub"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io/web"
	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// LinkMode defines how applications are linked to their Network Server.
type LinkMode int

const (
	// LinkAll links all applications in the link registry to their Network Server automatically.
	LinkAll LinkMode = iota
	// LinkExplicit links applications on request.
	LinkExplicit
)

// InteropClient is a client, which Application Server can use for interoperability.
type InteropClient interface {
	GetAppSKey(ctx context.Context, asID string, req *ttnpb.SessionKeyRequest) (*ttnpb.AppSKeyResponse, error)
}

// InteropConfig represents interoperability client configuration.
type InteropConfig struct {
	config.InteropClient `name:",squash"`
	ID                   string `name:"id" description:"AS-ID used for interoperability"`
}

// Config represents the ApplicationServer configuration.
type Config struct {
	LinkMode            string                    `name:"link-mode" description:"Mode to link applications to their Network Server (all, explicit)"`
	Devices             DeviceRegistry            `name:"-"`
	Links               LinkRegistry              `name:"-"`
	MQTT                config.MQTT               `name:"mqtt" description:"MQTT configuration"`
	Webhooks            WebhooksConfig            `name:"webhooks" description:"Webhooks configuration"`
	PubSub              PubSubConfig              `name:"pubsub" description:"Pub/sub messaging configuration"`
	ApplicationPackages ApplicationPackagesConfig `name:"application-packages" description:"Application packages configuration"`
	Interop             InteropConfig             `name:"interop" description:"Interop client configuration"`
}

var errLinkMode = errors.DefineInvalidArgument("link_mode", "invalid link mode `{value}`")

// GetLinkMode returns the converted configuration's link mode to LinkMode.
func (c Config) GetLinkMode() (LinkMode, error) {
	switch c.LinkMode {
	case "all":
		return LinkAll, nil
	case "explicit":
		return LinkExplicit, nil
	default:
		return LinkMode(0), errLinkMode.WithAttributes("value", c.LinkMode)
	}
}

var (
	errWebhooksRegistry = errors.DefineInvalidArgument("webhooks_registry", "invalid webhooks registry")
	errWebhooksTarget   = errors.DefineInvalidArgument("webhooks_target", "invalid webhooks target `{target}`")
)

// WebhooksConfig defines the configuration of the webhooks integration.
type WebhooksConfig struct {
	Registry  web.WebhookRegistry `name:"-"`
	Target    string              `name:"target" description:"Target of the integration (direct)"`
	Timeout   time.Duration       `name:"timeout" description:"Wait timeout of the target to process the request"`
	QueueSize int                 `name:"queue-size" description:"Number of requests to queue"`
	Workers   int                 `name:"workers" description:"Number of workers to process requests"`
	Templates web.TemplatesConfig `name:"templates" description:"The store of the webhook templates"`
}

// PubSubConfig contains go-cloud PubSub configuration of the Application Server.
type PubSubConfig struct {
	Registry pubsub.Registry `name:"-"`
}

// ApplicationPackagesConfig contains application packages associations configuration.
type ApplicationPackagesConfig struct {
	Registry packages.Registry `name:"-"`
}

// NewWebhooks returns a new web.Webhooks based on the configuration.
// If Target is empty, this method returns nil.
func (c WebhooksConfig) NewWebhooks(ctx context.Context, server io.Server) (web.Webhooks, error) {
	var target web.Sink
	switch c.Target {
	case "":
		return nil, nil
	case "direct":
		target = &web.HTTPClientSink{
			Client: &http.Client{
				Timeout: c.Timeout,
			},
		}
	default:
		return nil, errWebhooksTarget.WithAttributes("target", c.Target)
	}
	if c.Registry == nil {
		return nil, errWebhooksRegistry
	}
	if c.QueueSize > 0 || c.Workers > 0 {
		target = &web.QueuedSink{
			Target:  target,
			Queue:   make(chan *http.Request, c.QueueSize),
			Workers: c.Workers,
		}
	}
	if controllable, ok := target.(web.ControllableSink); ok {
		go func() {
			if err := controllable.Run(ctx); err != nil && !errors.IsCanceled(err) {
				log.FromContext(ctx).WithError(err).Error("Webhooks target sink failed")
			}
		}()
	}
	return web.NewWebhooks(ctx, server, c.Registry, target), nil
}

// NewPubSub returns a new pubsub.PubSub based on the configuration.
// If the registry is nil, it returns nil.
func (c PubSubConfig) NewPubSub(comp *component.Component, server io.Server) (*pubsub.PubSub, error) {
	if c.Registry == nil {
		return nil, nil
	}
	return pubsub.New(comp, server, c.Registry)
}

// NewApplicationPackages returns a new applications packages frontend based on the configuration.
// If the registry is nil, it returns nil.
func (c ApplicationPackagesConfig) NewApplicationPackages(ctx context.Context, server io.Server) (packages.Server, error) {
	if c.Registry == nil {
		return nil, nil
	}
	return packages.New(ctx, server, c.Registry)
}
