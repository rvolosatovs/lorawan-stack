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

package component

import (
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"go.thethings.network/lorawan-stack/v3/pkg/metrics"
	"go.thethings.network/lorawan-stack/v3/pkg/web"
	"go.thethings.network/lorawan-stack/v3/pkg/webmiddleware"
)

const (
	metricsUsername = "metrics"
	pprofUsername   = "pprof"
	healthUsername  = "health"
)

func (c *Component) initWeb() error {
	webOptions := []web.Option{
		web.WithContextFiller(c.FillContext),
		web.WithTrustedProxies(c.config.HTTP.TrustedProxies...),
		web.WithCookieKeys(c.config.HTTP.Cookie.HashKey, c.config.HTTP.Cookie.BlockKey),
		web.WithStatic(c.config.HTTP.Static.Mount, c.config.HTTP.Static.SearchPath...),
	}
	if c.config.HTTP.RedirectToHost != "" {
		webOptions = append(webOptions, web.WithRedirectToHost(c.config.HTTP.RedirectToHost))
	}
	if c.config.HTTP.RedirectToHTTPS {
		httpAddr, err := net.ResolveTCPAddr("tcp", c.config.HTTP.Listen)
		if err != nil {
			return err
		}
		httpsAddr, err := net.ResolveTCPAddr("tcp", c.config.HTTP.ListenTLS)
		if err != nil {
			return err
		}
		if httpsAddr.Port == 0 {
			httpsAddr.Port = 443
		}
		webOptions = append(webOptions, web.WithRedirectToHTTPS(httpAddr.Port, httpsAddr.Port))
		if httpAddr.Port != 80 && httpsAddr.Port != 443 {
			webOptions = append(webOptions, web.WithRedirectToHTTPS(80, 443))
		}
	}
	web, err := web.New(c.ctx, webOptions...)
	if err != nil {
		return err
	}
	c.web = web
	return nil
}

// RegisterWeb registers a web subsystem to the component.
func (c *Component) RegisterWeb(s web.Registerer) {
	c.webSubsystems = append(c.webSubsystems, s)
}

// RegisterLivenessCheck registers a liveness check for the component.
func (c *Component) RegisterLivenessCheck(name string, check healthcheck.Check) {
	c.healthHandler.AddLivenessCheck(name, check)
}

// RegisterReadinessCheck registers a readiness check for the component.
func (c *Component) RegisterReadinessCheck(name string, check healthcheck.Check) {
	c.healthHandler.AddReadinessCheck(name, check)
}

func (c *Component) serveWeb(lis net.Listener) error {
	srv := http.Server{
		Handler:           c,
		ReadTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-c.Context().Done()
		srv.Close()
	}()
	return srv.Serve(lis)
}

func (c *Component) webEndpoints() []Endpoint {
	return []Endpoint{
		NewTCPEndpoint(c.config.HTTP.Listen, "Web"),
		NewTLSEndpoint(c.config.HTTP.ListenTLS, "Web", WithNextProtos("h2", "http/1.1")),
	}
}

// listenWeb starts the web listeners on the addresses and endpoints configured in the HTTP section.
func (c *Component) listenWeb() (err error) {
	err = c.serveOnEndpoints(c.webEndpoints(), (*Component).serveWeb, "web")
	if err != nil {
		return
	}

	if c.config.HTTP.PProf.Enable {
		g := c.web.RootRouter().NewRoute().Subrouter()
		if c.config.HTTP.PProf.Password != "" {
			g.Use(mux.MiddlewareFunc(webmiddleware.BasicAuth(
				"pprof",
				webmiddleware.AuthUser(pprofUsername, c.config.HTTP.PProf.Password),
			)))
		}
		g.HandleFunc("/debug/pprof/profile", pprof.Profile)
		g.HandleFunc("/debug/pprof/trace", pprof.Trace)
		g.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
		g.Handle("/debug/pprof", http.RedirectHandler("/debug/pprof/", http.StatusMovedPermanently))
	}

	if c.config.HTTP.Metrics.Enable {
		g := c.web.RootRouter().NewRoute().Subrouter()
		if c.config.HTTP.Metrics.Password != "" {
			g.Use(mux.MiddlewareFunc(webmiddleware.BasicAuth(
				"metrics",
				webmiddleware.AuthUser(metricsUsername, c.config.HTTP.Metrics.Password),
			)))
		}
		g.Handle("/metrics", metrics.Exporter)
	}

	if c.config.HTTP.Health.Enable {
		g := c.web.RootRouter().NewRoute().Subrouter()
		if c.config.HTTP.Health.Password != "" {
			g.Use(mux.MiddlewareFunc(webmiddleware.BasicAuth(
				"health",
				webmiddleware.AuthUser(healthUsername, c.config.HTTP.Health.Password),
			)))
		}
		g.HandleFunc("/healthz/live", c.healthHandler.LiveEndpoint)
		g.HandleFunc("/healthz/ready", c.healthHandler.ReadyEndpoint)
	}

	return nil
}
