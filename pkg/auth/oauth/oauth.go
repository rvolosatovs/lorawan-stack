// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package oauth

import (
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/TheThingsNetwork/ttn/pkg/auth"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/types"
	"github.com/TheThingsNetwork/ttn/pkg/web"
	"github.com/TheThingsNetwork/ttn/pkg/web/middleware"
	"github.com/labstack/echo"
)

type Server struct {
	keys       *auth.Keys
	iss        string
	oauth      *osin.Server
	authorizer Authorizer
}

func New(iss string, keys *auth.Keys, store store.OAuthStore, authorizer Authorizer) *Server {
	config := &osin.ServerConfig{
		AuthorizationExpiration:     60 * 5,  // 5 minutes
		AccessExpiration:            60 * 60, // 1 hour
		ErrorStatusCode:             http.StatusUnauthorized,
		RequirePKCEForPublicClients: false,
		RedirectUriSeparator:        "",
		RetainTokenAfterRefresh:     false,
		AllowClientSecretInParams:   false,
		TokenType:                   "bearer",
		AllowedAuthorizeTypes: osin.AllowedAuthorizeType{
			osin.CODE,
		},
		AllowedAccessTypes: osin.AllowedAccessType{
			osin.AUTHORIZATION_CODE,
			osin.REFRESH_TOKEN,
			osin.CLIENT_CREDENTIALS,
		},
	}

	storage := &storage{
		keys:  keys,
		store: store,
	}

	s := &Server{
		keys:       keys,
		iss:        iss,
		oauth:      osin.NewServer(config, storage),
		authorizer: authorizer,
	}

	s.oauth.AccessTokenGen = s

	return s
}

// Register registers the server to the web server.
func (s *Server) Register(server *web.Server) {
	group := server.Group.Group("/oauth", middleware.JSONForm)
	group.Any("/token", s.tokenHandler)
	group.Any("/authorize", s.authorizationHandler)
}

func (s *Server) tokenHandler(c echo.Context) error {
	req := c.Request()
	resp := s.oauth.NewResponse()
	defer resp.Close()

	ar := s.oauth.HandleAccessRequest(resp, req)
	if ar == nil {
		return output(c, resp)
	}

	client := ar.Client.(types.Client).GetClient()

	switch ar.Type {
	case osin.AUTHORIZATION_CODE:
		ar.Authorized = client != nil && client.Grants.AuthorizationCode
	case osin.REFRESH_TOKEN:
		ar.Authorized = client != nil && client.Grants.RefreshToken
	case osin.PASSWORD:
		ar.Authorized = client != nil && client.Grants.Password
	case osin.CLIENT_CREDENTIALS, osin.ASSERTION, osin.IMPLICIT:
		// not supported
		ar.Authorized = false
	}

	s.oauth.FinishAccessRequest(resp, req, ar)

	return output(c, resp)
}

func (s *Server) authorizationHandler(c echo.Context) error {
	req := c.Request()
	resp := s.oauth.NewResponse()
	defer resp.Close()

	ar := s.oauth.HandleAuthorizeRequest(resp, req)
	if ar == nil {
		return output(c, resp)
	}
	client := ar.Client.(types.Client)

	// make sure client supports authorization code
	if !client.GetClient().Grants.AuthorizationCode {
		resp.SetError(osin.E_INVALID_CLIENT, "")
		s.oauth.FinishAuthorizeRequest(resp, req, ar)
		return output(c, resp)
	}

	// TODO: match the rights of the client to the scope of the request

	// make sure the user is logged in or redirect
	err := s.authorizer.CheckLogin(c)
	if err != nil || c.Response().Committed {
		return err
	}

	// check if the user authorized, or redner the form
	authorized, err := s.authorizer.Authorize(c, client)
	if err != nil || c.Response().Committed {
		return err
	}

	ar.Authorized = authorized

	s.oauth.FinishAuthorizeRequest(resp, req, ar)
	return output(c, resp)
}

func output(c echo.Context, resp *osin.Response) error {
	if resp.IsError && resp.InternalError != nil {
		// TODO: log internal error
	}

	headers := c.Response().Header()

	// Add headers
	for i, k := range resp.Headers {
		for _, v := range k {
			headers.Add(i, v)
		}
	}

	if resp.Type == osin.REDIRECT {
		// output redirect with parameters
		location, err := resp.GetRedirectUrl()
		if err != nil {
			return err
		}
		headers.Add("Location", location)

		return c.NoContent(http.StatusFound)
	}

	return c.JSON(resp.StatusCode, resp.Output)
}
