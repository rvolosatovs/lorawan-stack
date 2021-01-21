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

package oauth

import (
	"context"
	"net/http"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/openshift/osin"
	"go.thethings.network/lorawan-stack/v3/pkg/account/session"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	web_errors "go.thethings.network/lorawan-stack/v3/pkg/errors/web"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/web"
	"go.thethings.network/lorawan-stack/v3/pkg/web/middleware"
	"go.thethings.network/lorawan-stack/v3/pkg/webui"
)

// Server is the interface for the OAuth server.
type Server interface {
	web.Registerer

	Authorize(authorizePage echo.HandlerFunc) echo.HandlerFunc
	Token(c echo.Context) error
}

type server struct {
	c          *component.Component
	config     Config
	osinConfig *osin.ServerConfig
	store      Store
	session    session.Session
}

// Store used by the OAuth server.
type Store interface {
	// UserStore and UserSessionStore are needed for user login/logout.
	store.UserStore
	store.UserSessionStore
	// ClientStore is needed for getting the OAuth client.
	store.ClientStore
	// OAuth is needed for OAuth authorizations.
	store.OAuthStore
}

// NewServer returns a new OAuth server on top of the given store.
func NewServer(c *component.Component, store Store, config Config) (Server, error) {
	s := &server{
		c:       c,
		config:  config,
		store:   store,
		session: session.Session{Store: store},
	}

	if s.config.Mount == "" {
		s.config.Mount = s.config.UI.MountPath()
	}

	s.osinConfig = &osin.ServerConfig{
		AuthorizationExpiration: int32((5 * time.Minute).Seconds()),
		AccessExpiration:        int32(time.Hour.Seconds()),
		TokenType:               "bearer",
		AllowedAuthorizeTypes: osin.AllowedAuthorizeType{
			osin.CODE,
		},
		AllowedAccessTypes: osin.AllowedAccessType{
			osin.AUTHORIZATION_CODE,
			osin.REFRESH_TOKEN,
			osin.PASSWORD,
		},
		ErrorStatusCode:           http.StatusBadRequest,
		AllowClientSecretInParams: true,
		RedirectUriSeparator:      redirectURISeparator,
		RetainTokenAfterRefresh:   false,
	}

	return s, nil
}

type ctxKeyType struct{}

var ctxKey ctxKeyType

func (s *server) configFromContext(ctx context.Context) *Config {
	if config, ok := ctx.Value(ctxKey).(*Config); ok {
		return config
	}
	return &s.config
}

func (s *server) now() time.Time { return time.Now().UTC() }

func (s *server) oauth2(ctx context.Context) *osin.Server {
	oauth2 := osin.NewServer(s.osinConfig, &storage{
		ctx:     ctx,
		clients: s.store,
		oauth:   s.store,
	})
	oauth2.AuthorizeTokenGen = s
	oauth2.AccessTokenGen = s
	oauth2.Now = s.now
	oauth2.Logger = s
	return oauth2
}

func (s *server) Printf(format string, v ...interface{}) {
	log.FromContext(s.c.Context()).Warnf(format, v...)
}

// These errors map to errors in the osin library.
var (
	errInvalidRequest          = errors.DefineInvalidArgument("invalid_request", "invalid or missing request parameter")
	errUnauthorizedClient      = errors.DefinePermissionDenied("unauthorized_client", "client is not authorized to request a token using this method")
	errAccessDenied            = errors.DefinePermissionDenied("access_denied", "access denied")
	errUnsupportedResponseType = errors.DefineUnimplemented("unsupported_response_type", "unsupported response type")
	errInvalidScope            = errors.DefineInvalidArgument("invalid_scope", "invalid scope")
	errUnsupportedGrantType    = errors.DefineUnimplemented("unsupported_grant_type", "unsupported grant type")
	errInvalidGrant            = errors.DefinePermissionDenied("invalid_grant", "invalid, expired or revoked authorization code")
	errInvalidClient           = errors.DefinePermissionDenied("invalid client", "invalid or unauthenticated client")
	errInternal                = errors.Define("internal", "internal error {id}")
	errInvalidRedirectURI      = errors.DefinePermissionDenied("invalid_redirect_uri", "invalid redirect URI")
)

func (s *server) output(c echo.Context, resp *osin.Response) error {
	headers := c.Response().Header()
	for i, k := range resp.Headers {
		for _, v := range k {
			headers.Add(i, v)
		}
	}

	var osinErr error
	if resp.IsError {
		switch resp.ErrorId {
		case osin.E_INVALID_REQUEST:
			osinErr = errInvalidRequest
		case osin.E_UNAUTHORIZED_CLIENT:
			osinErr = errUnauthorizedClient
		case osin.E_ACCESS_DENIED:
			osinErr = errAccessDenied
		case osin.E_UNSUPPORTED_RESPONSE_TYPE:
			osinErr = errUnsupportedResponseType
		case osin.E_INVALID_SCOPE:
			osinErr = errInvalidScope
		case osin.E_UNSUPPORTED_GRANT_TYPE:
			osinErr = errUnsupportedGrantType
		case osin.E_INVALID_GRANT:
			osinErr = errInvalidGrant
		case osin.E_INVALID_CLIENT:
			osinErr = errInvalidClient
		default:
			osinErr = errInternal
		}
		if resp.InternalError != nil {
			if ttnErr, ok := errors.From(resp.InternalError); ok {
				osinErr = ttnErr
			} else if _, isURIValidationError := resp.InternalError.(osin.UriValidationError); isURIValidationError {
				osinErr = errInvalidRedirectURI.WithCause(resp.InternalError)
			} else {
				osinErr = osinErr.(errors.Definition).WithCause(resp.InternalError)
			}
		}
		log.FromContext(c.Request().Context()).WithError(osinErr).Warn("OAuth error")
	}

	if resp.Type == osin.REDIRECT {
		location, err := resp.GetRedirectUrl()
		if err != nil {
			return err
		}
		uiMount := strings.TrimSuffix(s.config.UI.MountPath(), "/")
		if strings.HasPrefix(location, "/code") || strings.HasPrefix(location, "/local-callback") {
			location = uiMount + location
		}
		return c.Redirect(http.StatusFound, location)
	}

	if osinErr != nil {
		return osinErr
	}

	return c.JSON(resp.StatusCode, resp.Output)
}

func (s *server) RegisterRoutes(server *web.Server) {
	root := server.Group(
		s.config.Mount,
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				config := s.configFromContext(c.Request().Context())
				c.Set("template_data", config.UI.TemplateData)
				frontendConfig := config.UI.FrontendConfig
				frontendConfig.Language = config.UI.TemplateData.Language
				c.Set("app_config", struct {
					FrontendConfig
				}{
					FrontendConfig: frontendConfig,
				})
				return next(c)
			}
		},
		web_errors.ErrorMiddleware(map[string]web_errors.ErrorRenderer{
			"text/html": webui.Template,
		}),
	)

	csrfMiddleware := middleware.CSRF("_csrf", "/", s.config.CSRFAuthKey)

	page := root.Group("", csrfMiddleware)

	// The logout route is currently in use by existing OAuth clients. As part of
	// the public API it should not be removed in this major.
	page.GET("/logout", s.ClientLogout)

	page.GET("/authorize", s.Authorize(webui.Template.Handler), s.redirectToLogin)
	page.POST("/authorize", s.Authorize(webui.Template.Handler), s.redirectToLogin)

	root.GET("/local-callback", s.redirectToLocal)

	// No CSRF here:
	root.POST("/token", s.Token)
}
