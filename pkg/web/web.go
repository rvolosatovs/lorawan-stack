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

package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	echo "github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.thethings.network/lorawan-stack/v3/pkg/auth"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/fillcontext"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/random"
	"go.thethings.network/lorawan-stack/v3/pkg/web/cookie"
	"go.thethings.network/lorawan-stack/v3/pkg/webhandlers"
	"go.thethings.network/lorawan-stack/v3/pkg/webmiddleware"
	"go.thethings.network/lorawan-stack/v3/pkg/webui"
	"gopkg.in/yaml.v2"
)

// Registerer allows components to register their services to the web server.
type Registerer interface {
	RegisterRoutes(s *Server)
}

// Server is the server.
type Server struct {
	// The root HTTP router.
	root *mux.Router

	// The main HTTP router.
	router *mux.Router

	// The HTTP router for API.
	apiRouter *mux.Router

	// The legacy HTTP framework.
	echo *echo.Echo
}

type options struct {
	cookieHashKey  []byte
	cookieBlockKey []byte

	staticMount       string
	staticSearchPaths []string

	trustedProxies []string

	contextFillers []fillcontext.Filler

	redirectToHost  string
	redirectToHTTPS map[int]int

	logIgnorePaths []string
}

// Option for the web server
type Option func(*options)

// WithContextFiller sets context fillers that are executed on every request context.
func WithContextFiller(contextFillers ...fillcontext.Filler) Option {
	return func(o *options) {
		o.contextFillers = append(o.contextFillers, contextFillers...)
	}
}

// WithTrustedProxies adds trusted proxies from which proxy headers are trusted.
func WithTrustedProxies(cidrs ...string) Option {
	return func(o *options) {
		o.trustedProxies = append(o.trustedProxies, cidrs...)
	}
}

// WithCookieKeys sets the cookie hash key and block key.
func WithCookieKeys(hashKey, blockKey []byte) Option {
	return func(o *options) {
		o.cookieHashKey, o.cookieBlockKey = hashKey, blockKey
	}
}

// WithStatic sets the mount and search paths for static assets.
func WithStatic(mount string, searchPaths ...string) Option {
	return func(o *options) {
		o.staticMount, o.staticSearchPaths = mount, searchPaths
	}
}

// WithRedirectToHost redirects all requests to this host.
func WithRedirectToHost(target string) Option {
	return func(o *options) {
		o.redirectToHost = target
	}
}

// WithRedirectToHTTPS redirects HTTP requests to HTTPS.
func WithRedirectToHTTPS(from, to int) Option {
	return func(o *options) {
		if o.redirectToHTTPS == nil {
			o.redirectToHTTPS = make(map[int]int)
		}
		o.redirectToHTTPS[from] = to
	}
}

// WithLogIgnorePaths silences log messages for a list of URLs.
func WithLogIgnorePaths(paths []string) Option {
	return func(o *options) {
		o.logIgnorePaths = paths
	}
}

// New builds a new server.
func New(ctx context.Context, opts ...Option) (*Server, error) {
	logger := log.FromContext(ctx).WithField("namespace", "web")

	options := new(options)
	for _, opt := range opts {
		opt(options)
	}

	hashKey, blockKey := options.cookieHashKey, options.cookieBlockKey

	if len(hashKey) == 0 || isZeros(hashKey) {
		hashKey = random.Bytes(64)
		logger.Warn("No cookie hash key configured, generated a random one")
	}

	if len(hashKey) != 32 && len(hashKey) != 64 {
		return nil, errors.New("Expected cookie hash key to be 32 or 64 bytes long")
	}

	if len(blockKey) == 0 || isZeros(blockKey) {
		blockKey = random.Bytes(32)
		logger.Warn("No cookie block key configured, generated a random one")
	}

	if len(blockKey) != 32 {
		return nil, errors.New("Expected cookie block key to be 32 bytes long")
	}

	var proxyConfiguration webmiddleware.ProxyConfiguration
	proxyConfiguration.ParseAndAddTrusted(options.trustedProxies...)
	root := mux.NewRouter()
	root.NotFoundHandler = http.HandlerFunc(webhandlers.NotFound)
	root.Use(
		mux.MiddlewareFunc(webmiddleware.Recover()),
		mux.MiddlewareFunc(webmiddleware.FillContext(options.contextFillers...)),
		mux.MiddlewareFunc(webmiddleware.RequestURL()),
		mux.MiddlewareFunc(webmiddleware.RequestID()),
		mux.MiddlewareFunc(webmiddleware.ProxyHeaders(proxyConfiguration)),
		mux.MiddlewareFunc(webmiddleware.MaxBody(1<<24)), // 16 MB.
		mux.MiddlewareFunc(webmiddleware.SecurityHeaders()),
		mux.MiddlewareFunc(webmiddleware.Log(logger, options.logIgnorePaths)),
		mux.MiddlewareFunc(webmiddleware.Cookies(hashKey, blockKey)),
	)

	var redirectConfig webmiddleware.RedirectConfiguration
	if options.redirectToHost != "" {
		if host, portStr, err := net.SplitHostPort(options.redirectToHost); err == nil {
			redirectConfig.HostName = func(string) string { return host }
			port, err := strconv.ParseUint(portStr, 10, 0)
			if err != nil {
				return nil, err
			}
			redirectConfig.Port = func(uint) uint { return uint(port) }
		} else {
			redirectConfig.HostName = func(string) string { return options.redirectToHost }
		}
	}
	if options.redirectToHTTPS != nil {
		redirectConfig.Scheme = func(string) string { return "https" }
		// Only redirect to HTTPS port if no port redirection has been configured
		if redirectConfig.Port == nil {
			redirectConfig.Port = func(current uint) uint {
				return uint(options.redirectToHTTPS[int(current)])
			}
		}
	}

	router := root.NewRoute().Subrouter()
	router.Use(
		mux.MiddlewareFunc(webmiddleware.Redirect(redirectConfig)),
	)

	corsSkip := func(r *http.Request) bool {
		authVal := r.Header.Get("Authorization")
		if authVal != "" || !(strings.HasPrefix(authVal, "Bearer ")) {
			token := strings.TrimPrefix(authVal, "Bearer ")
			tokenType, _, _, err := auth.SplitToken(token)
			if err == nil {
				if tokenType == auth.SessionToken {
					return false
				}
			}
		}
		return true
	}
	csrfSkip := func(r *http.Request) bool {
		return !corsSkip(r)
	}

	apiRouter := mux.NewRouter()
	apiRouter.NotFoundHandler = http.HandlerFunc(webhandlers.NotFound)
	apiRouter.Use(
		mux.MiddlewareFunc(webmiddleware.CookieAuth("_session")),
		mux.MiddlewareFunc(webmiddleware.Conditional(
			webmiddleware.CSRF(hashKey, csrf.CookieName("_csrf"), csrf.Path("/"), csrf.SameSite(csrf.SameSiteStrictMode)),
			csrfSkip,
		)),
		mux.MiddlewareFunc(
			webmiddleware.Conditional(
				webmiddleware.CORS(webmiddleware.CORSConfig{
					AllowedHeaders:   []string{"Authorization", "Content-Type", "X-CSRF-Token"},
					AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
					AllowedOrigins:   []string{"*"},
					ExposedHeaders:   []string{"Date", "Content-Length", "X-Request-Id", "X-Total-Count", "X-Warning"},
					MaxAge:           600,
					AllowCredentials: true,
				}),
				corsSkip,
			),
		),
	)
	root.PathPrefix("/api/").Handler(apiRouter)

	server := echo.New()

	server.Logger = &noopLogger{}
	server.HTTPErrorHandler = errorHandler

	server.Use(
		echomiddleware.Gzip(),
		cookie.Cookies(hashKey, blockKey),
	)

	s := &Server{
		root:      root,
		router:    router,
		apiRouter: apiRouter,
		echo:      server,
	}

	var staticPath string
	for _, path := range options.staticSearchPaths {
		if s, err := os.Stat(path); err == nil && s.IsDir() {
			staticPath = path
			break
		}
	}
	if staticPath != "" {
		staticDir := http.Dir(staticPath)
		logger := logger.WithFields(log.Fields("path", staticDir, "mount", options.staticMount))
		s.Static(options.staticMount, staticDir)

		// register hashed filenames
		manifest, err := ioutil.ReadFile(filepath.Join(staticPath, "manifest.yaml"))
		if err != nil {
			logger.WithError(err).Warn("Failed to load manifest.yaml")
			return s, nil
		}
		hashedFiles := make(map[string]string)
		err = yaml.Unmarshal(manifest, &hashedFiles)
		if err != nil {
			return nil, errors.New("Corrupted manifest.yaml").WithCause(err)
		}
		for original, hashed := range hashedFiles {
			webui.RegisterHashedFile(original, hashed)
		}
		logger.Debug("Loaded manifest.yaml")
		logger.Debug("Serving static assets")
	} else {
		logger.WithField("search_paths", options.staticSearchPaths).Warn("No static assets found in any search path")
	}

	return s, nil
}

func isZeros(buf []byte) bool {
	for _, b := range buf {
		if b != 0x00 {
			return false
		}
	}

	return true
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.root.ServeHTTP(w, r)
}

// RootRouter returns the root router.
// In most cases the Router() should be used instead of the root router.
func (s *Server) RootRouter() *mux.Router {
	return s.root
}

// Router returns the main router.
func (s *Server) Router() *mux.Router {
	return s.router
}

// APIRouter returns the API router.
func (s *Server) APIRouter() *mux.Router {
	return s.apiRouter
}

func (s *Server) getRouter(path string) *mux.Router {
	if strings.HasPrefix(path, "/api/") {
		return s.apiRouter
	}
	return s.router
}

var echoVar = regexp.MustCompile(":[^/]+")

// replaceEchoVars replaces Echo path variables with Mux path variables.
// "/users/:user_id" will become "/users/{user_id}"
func replaceEchoVars(path string) string {
	return echoVar.ReplaceAllStringFunc(path, func(s string) string {
		return fmt.Sprintf("{%s}", strings.TrimPrefix(s, ":"))
	})
}

// Group creates a sub group.
func (s *Server) Group(prefix string, middleware ...echo.MiddlewareFunc) *echo.Group {
	path := "/" + strings.Trim(prefix, "/")
	pathWithSlash := path + "/"
	router := s.getRouter(pathWithSlash)
	router.PathPrefix(replaceEchoVars(pathWithSlash)).Handler(s.echo)
	if !strings.HasSuffix(path, "/") {
		router.Handle(replaceEchoVars(path), http.RedirectHandler(pathWithSlash, http.StatusPermanentRedirect))
	}
	return s.echo.Group(path, middleware...)
}

// GET registers a GET handler at path.
func (s *Server) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	router := s.getRouter(path)
	router.Path(path).Methods(http.MethodGet).Handler(s.echo)
	return s.echo.GET(replaceEchoVars(path), h, m...)
}

// HEAD registers a HEAD handler at path.
func (s *Server) HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	router := s.getRouter(path)
	router.Path(path).Methods(http.MethodHead).Handler(s.echo)
	return s.echo.HEAD(replaceEchoVars(path), h, m...)
}

// POST registers a POST handler at path.
func (s *Server) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	router := s.getRouter(path)
	router.Path(path).Methods(http.MethodPost).Handler(s.echo)
	return s.echo.POST(replaceEchoVars(path), h, m...)
}

// DELETE registers a DELETE handler at path.
func (s *Server) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	router := s.getRouter(path)
	router.Path(path).Methods(http.MethodDelete).Handler(s.echo)
	return s.echo.DELETE(replaceEchoVars(path), h, m...)
}

// Static adds the http.FileSystem under the defined prefix.
func (s *Server) Static(prefix string, fs http.FileSystem) {
	prefix = "/" + strings.Trim(prefix, "/") + "/"
	s.router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(fs)))
}

// Prefix returns a route for the given path prefix.
func (s *Server) Prefix(prefix string) *mux.Route {
	return s.getRouter(prefix).PathPrefix(prefix)
}
