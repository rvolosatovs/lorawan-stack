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
	"crypto/subtle"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/heptiolabs/healthcheck"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/metrics"
	"go.thethings.network/lorawan-stack/pkg/web"
)

const (
	metricsUsername = "metrics"
	pprofUsername   = "pprof"
	healthUsername  = "health"
)

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

func (c *Component) serveHTTP(lis net.Listener) error {
	return http.Serve(lis, c)
}

func (c *Component) httpListenerConfigs() []endpoint {
	return []endpoint{
		{toNativeListener: Listener.TCP, address: c.config.HTTP.Listen, protocol: "HTTP"},
		{toNativeListener: Listener.TLS, address: c.config.HTTP.ListenTLS, protocol: "HTTPS"},
	}
}

func (c *Component) listenWeb() (err error) {
	err = c.serveOnListeners(c.httpListenerConfigs(), (*Component).serveHTTP, "web")
	if err != nil {
		return
	}

	if c.config.HTTP.PProf.Enable {
		var middleware []echo.MiddlewareFunc
		if c.config.HTTP.PProf.Password != "" {
			middleware = append(middleware, c.basicAuth(pprofUsername, c.config.HTTP.PProf.Password))
		}
		g := c.web.RootGroup("/debug/pprof", middleware...)
		g.GET("", func(c echo.Context) error { return c.Redirect(http.StatusFound, c.Path()+"/") })
		g.GET("/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
		g.GET("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
		g.GET("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	}

	if c.config.HTTP.Metrics.Enable {
		var middleware []echo.MiddlewareFunc
		if c.config.HTTP.Metrics.Password != "" {
			middleware = append(middleware, c.basicAuth(metricsUsername, c.config.HTTP.Metrics.Password))
		}
		g := c.web.RootGroup("/metrics", middleware...)
		g.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusFound, strings.TrimSuffix(c.Path(), "/")) })
		g.GET("", echo.WrapHandler(metrics.Exporter), func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Request().Header.Del("Accept-Encoding")
				return next(c)
			}
		})
	}

	if c.config.HTTP.Health.Enable {
		var middleware []echo.MiddlewareFunc
		if c.config.HTTP.Health.Password != "" {
			middleware = append(middleware, c.basicAuth(healthUsername, c.config.HTTP.Health.Password))
		}
		g := c.web.RootGroup("/healthz", middleware...)
		g.GET("/live", echo.WrapHandler(http.HandlerFunc(c.healthHandler.LiveEndpoint)))
		g.GET("/ready", echo.WrapHandler(http.HandlerFunc(c.healthHandler.ReadyEndpoint)))
	}

	return nil
}

func (c *Component) basicAuth(username, password string) echo.MiddlewareFunc {
	usernameBytes, passwordBytes := []byte(username), []byte(password)
	return middleware.BasicAuth(func(username string, password string, ctx echo.Context) (bool, error) {
		usernameCompare := subtle.ConstantTimeCompare([]byte(username), usernameBytes)
		passwordCompare := subtle.ConstantTimeCompare([]byte(password), passwordBytes)
		if usernameCompare != 1 || passwordCompare != 1 {
			c.Logger().WithFields(log.Fields(
				"namespace", "web",
				"url", ctx.Path(),
				"remote_addr", ctx.RealIP(),
			)).Warn("Basic auth failed")
			return false, nil
		}
		return true, nil
	})
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		c.grpc.Server.ServeHTTP(w, r)
	} else {
		c.web.ServeHTTP(w, r)
	}
}
