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

package oauth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/v3/pkg/auth"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/pbkdf2"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/oauth"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/webui"
	"golang.org/x/net/publicsuffix"
)

type loginFormData struct {
	encoding string
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type authorizeFormData struct {
	encoding  string
	Authorize bool `json:"authorize"`
}

var (
	mockSession = &ttnpb.UserSession{
		UserIdentifiers: ttnpb.UserIdentifiers{UserID: "user"},
		SessionID:       "session_id",
		CreatedAt:       time.Now().Truncate(time.Second),
		SessionSecret:   "secret-1234",
	}
	mockUser = &ttnpb.User{
		UserIdentifiers: ttnpb.UserIdentifiers{UserID: "user"},
	}
	mockClient = &ttnpb.Client{
		ClientIdentifiers:  ttnpb.ClientIdentifiers{ClientID: "client"},
		State:              ttnpb.STATE_APPROVED,
		Grants:             []ttnpb.GrantType{ttnpb.GRANT_AUTHORIZATION_CODE, ttnpb.GRANT_REFRESH_TOKEN},
		RedirectURIs:       []string{"https://uri/callback", "http://uri/callback"},
		LogoutRedirectURIs: []string{"https://uri/logout-callback", "http://uri/logout-callback", "http://uri/alternative-logout-callback", "http://other-host/logout-callback"},
		Rights:             []ttnpb.Right{ttnpb.RIGHT_USER_INFO},
	}
	mockAccessToken = &ttnpb.OAuthAccessToken{
		UserIDs:       ttnpb.UserIdentifiers{UserID: "user"},
		UserSessionID: "session_id",
	}
)

var authCookie *http.Cookie

func init() {
	ctx := test.Context()

	hashValidator := pbkdf2.Default()
	hashValidator.Iterations = 10
	ctx = auth.NewContextWithHashValidator(ctx, hashValidator)

	password, err := auth.Hash(ctx, "pass")
	if err != nil {
		panic(err)
	}
	mockUser.Password = password

	secret, err := auth.Hash(ctx, "secret")
	if err != nil {
		panic(err)
	}
	mockClient.Secret = secret

	// Create session cookie
	sc := securecookie.New(
		[]byte("12345678123456781234567812345678"),
		[]byte("12345678123456781234567812345678"),
	)
	authCookieContent := &auth.CookieShape{
		UserID:        "user",
		SessionID:     "session_id",
		SessionSecret: "secret-1234",
	}
	encoded, _ := sc.Encode("_session", authCookieContent)
	authCookie = &http.Cookie{
		Name:     "_session",
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
}

func TestOAuthFlow(t *testing.T) {
	store := &mockStore{}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			HTTP: config.HTTP{
				Cookie: config.Cookie{
					HashKey:  []byte("12345678123456781234567812345678"),
					BlockKey: []byte("12345678123456781234567812345678"),
				},
			},
		},
	})
	s, err := oauth.NewServer(c, store, oauth.Config{
		Mount:       "/oauth",
		CSRFAuthKey: []byte("12345678123456781234567812345678"),
		UI: oauth.UIConfig{
			TemplateData: webui.TemplateData{
				SiteName:     "The Things Network",
				Title:        "OAuth",
				CanonicalURL: "https://example.com/oauth",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	c.RegisterWeb(s)
	componenttest.StartComponent(t, c)

	var csrfToken string
	var r *http.Request

	// Obtain CSRF token.
	r = httptest.NewRequest("GET", "/oauth/login", nil)
	r.URL.Scheme, r.URL.Host = "http", r.Host
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	c.ServeHTTP(rr, r)
	csrfToken = rr.Header().Get("X-CSRF-Token")

	if cookies := rr.Result().Cookies(); len(cookies) > 0 {
		jar.SetCookies(r.URL, cookies)
	}

	for _, tt := range []struct {
		Name             string
		StoreSetup       func(*mockStore)
		StoreCheck       func(*testing.T, *mockStore)
		UseCookie        *http.Cookie
		Method           string
		Path             string
		Body             interface{}
		ExpectedCode     int
		ExpectedRedirect string
		ExpectedBody     string
	}{
		{
			Name: "redirect to root",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
			},
			Method:           "GET",
			Path:             "/oauth/authorize",
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "/oauth",
		},
		{
			Name: "authorization page",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = mockClient
				s.err.getAuthorization = mockErrNotFound
			},
			Method:       "GET",
			Path:         "/oauth/authorize?client_id=client&redirect_uri=http://uri/callback&response_type=code&state=foo",
			UseCookie:    authCookie,
			ExpectedCode: http.StatusOK,
			ExpectedBody: `"client":{"ids":{"client_id":"client"}`,
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.req.clientIDs.GetClientID(), should.Equal, "client")
				a.So(s.calls, should.Contain, "GetClient")
				a.So(s.calls, should.Contain, "GetAuthorization")
			},
		},
		{
			Name: "client not found",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = nil
				s.err.getClient = mockErrNotFound
			},
			Method:       "GET",
			Path:         "/oauth/authorize?client_id=client&redirect_uri=http://uri/callback&response_type=code&state=foo",
			UseCookie:    authCookie,
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: `NotFound`,
		},
		{
			Name: "invalid redirect uri",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = &ttnpb.Client{
					ClientIdentifiers: ttnpb.ClientIdentifiers{ClientID: "client"},
					State:             ttnpb.STATE_REJECTED,
					Grants:            []ttnpb.GrantType{ttnpb.GRANT_AUTHORIZATION_CODE, ttnpb.GRANT_REFRESH_TOKEN},
					RedirectURIs:      []string{"https://uri/callback", "http://uri/callback"},
					Rights:            []ttnpb.Right{ttnpb.RIGHT_USER_INFO},
				}
			},
			Method:       "GET",
			Path:         "/oauth/authorize?client_id=client&redirect_uri=http://other-uri/callback&response_type=code&state=foo",
			UseCookie:    authCookie,
			ExpectedCode: http.StatusForbidden,
			ExpectedBody: `redirect URI`,
		},
		{
			Name: "client not approved",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = &ttnpb.Client{
					ClientIdentifiers: ttnpb.ClientIdentifiers{ClientID: "client"},
					State:             ttnpb.STATE_REQUESTED,
					Grants:            []ttnpb.GrantType{ttnpb.GRANT_AUTHORIZATION_CODE, ttnpb.GRANT_REFRESH_TOKEN},
					RedirectURIs:      []string{"https://uri/callback", "http://uri/callback"},
					Rights:            []ttnpb.Right{ttnpb.RIGHT_USER_INFO},
				}
			},
			Method:           "GET",
			Path:             "/oauth/authorize?client_id=client&redirect_uri=http://uri/callback&response_type=code&state=foo",
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "http://uri/callback?error=invalid_client",
		},
		{
			Name: "client suspended",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = &ttnpb.Client{
					ClientIdentifiers: ttnpb.ClientIdentifiers{ClientID: "client"},
					State:             ttnpb.STATE_SUSPENDED,
					Grants:            []ttnpb.GrantType{ttnpb.GRANT_AUTHORIZATION_CODE, ttnpb.GRANT_REFRESH_TOKEN},
					RedirectURIs:      []string{"https://uri/callback", "http://uri/callback"},
					Rights:            []ttnpb.Right{ttnpb.RIGHT_USER_INFO},
				}
			},
			Method:           "GET",
			Path:             "/oauth/authorize?client_id=client&redirect_uri=http://uri/callback&response_type=code&state=foo",
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "http://uri/callback?error=invalid_client",
		},
		{
			Name: "client rejected",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = &ttnpb.Client{
					ClientIdentifiers: ttnpb.ClientIdentifiers{ClientID: "client"},
					State:             ttnpb.STATE_REJECTED,
					Grants:            []ttnpb.GrantType{ttnpb.GRANT_AUTHORIZATION_CODE, ttnpb.GRANT_REFRESH_TOKEN},
					RedirectURIs:      []string{"https://uri/callback", "http://uri/callback"},
					Rights:            []ttnpb.Right{ttnpb.RIGHT_USER_INFO},
				}
			},
			Method:           "GET",
			Path:             "/oauth/authorize?client_id=client&redirect_uri=http://uri/callback&response_type=code&state=foo",
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "http://uri/callback?error=invalid_client",
		},
		{
			Name: "client without grant",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = &ttnpb.Client{
					ClientIdentifiers: ttnpb.ClientIdentifiers{ClientID: "client"},
					State:             ttnpb.STATE_APPROVED,
					Grants:            []ttnpb.GrantType{},
					RedirectURIs:      []string{"https://uri/callback", "http://uri/callback"},
					Rights:            []ttnpb.Right{ttnpb.RIGHT_USER_INFO},
				}
			},
			Method:           "GET",
			Path:             "/oauth/authorize?client_id=client&redirect_uri=http://uri/callback&response_type=code&state=foo",
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "http://uri/callback?error=invalid_grant",
		},
		{
			Name: "authorize client",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.user = mockUser
				s.res.client = mockClient
				s.err.getAuthorization = mockErrNotFound
			},
			Method:           "POST",
			Path:             "/oauth/authorize?client_id=client&redirect_uri=http://uri/callback&response_type=code&state=foo",
			Body:             authorizeFormData{encoding: "form", Authorize: true},
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "http://uri/callback?code=",
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "Authorize")
				a.So(s.req.authorization.UserIDs, should.Resemble, mockUser.UserIdentifiers)
				a.So(s.req.authorization.ClientIDs, should.Resemble, mockClient.ClientIdentifiers)
				a.So(s.calls, should.Contain, "CreateAuthorizationCode")
				a.So(s.req.authorizationCode.UserIDs, should.Resemble, mockUser.UserIdentifiers)
				a.So(s.req.authorizationCode.ClientIDs, should.Resemble, mockClient.ClientIdentifiers)
				a.So(s.req.authorizationCode.UserSessionID, should.Equal, mockSession.SessionID)
				a.So(s.req.authorizationCode.Rights, should.Resemble, mockClient.Rights)
				a.So(s.req.authorizationCode.Code, should.NotBeEmpty)
				a.So(s.req.authorizationCode.RedirectURI, should.Equal, "http://uri/callback")
				a.So(s.req.authorizationCode.State, should.Equal, "foo")
			},
		},
		{
			Name: "client-initiated logout",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.client = mockClient
				s.res.accessToken = mockAccessToken
			},
			Method:           "GET",
			Path:             "/oauth/logout?access_token_id=access-token-id&post_logout_redirect_uri=http://uri/alternative-logout-callback?foo=bar",
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "/alternative-logout-callback?foo=bar",
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "DeleteSession")
				a.So(s.calls, should.Contain, "GetAccessToken")
				a.So(s.calls, should.Contain, "DeleteAccessToken")
				a.So(s.calls, should.Contain, "GetClient")
			},
		},
		{
			Name: "client-initiated logout with redirect to different host",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.client = mockClient
				s.res.accessToken = mockAccessToken
			},
			Method:           "GET",
			Path:             "/oauth/logout?access_token_id=access-token-id&post_logout_redirect_uri=http://other-host/logout-callback?foo=bar",
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "http://other-host/logout-callback?foo=bar",
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "DeleteSession")
				a.So(s.calls, should.Contain, "GetAccessToken")
				a.So(s.calls, should.Contain, "DeleteAccessToken")
				a.So(s.calls, should.Contain, "GetClient")
			},
		},
		{
			Name: "client-initiated logout without redirect uri",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.client = mockClient
				s.res.accessToken = mockAccessToken
			},
			Method:           "GET",
			Path:             "/oauth/logout?access_token_id=access-token-id",
			UseCookie:        authCookie,
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "/logout-callback",
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "DeleteSession")
				a.So(s.calls, should.Contain, "GetAccessToken")
				a.So(s.calls, should.Contain, "DeleteAccessToken")
				a.So(s.calls, should.Contain, "GetClient")
			},
		},
		{
			Name: "client-initiated logout with invalid redirect uri",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.client = mockClient
				s.res.accessToken = mockAccessToken
			},
			Method:       "GET",
			Path:         "/oauth/logout?access_token_id=access-token-id&post_logout_redirect_uri=http://uri/false-callback",
			UseCookie:    authCookie,
			ExpectedCode: http.StatusBadRequest,
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "DeleteSession")
				a.So(s.calls, should.Contain, "DeleteAccessToken")
				a.So(s.calls, should.Contain, "GetAccessToken")
				a.So(s.calls, should.Contain, "GetClient")
			},
		},
		{
			Name: "client-initiated without logout redirect URI set in the client",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.client = mockClient
				s.res.client.LogoutRedirectURIs = []string{}
				s.res.accessToken = mockAccessToken
			},
			Method:           "GET",
			Path:             "/oauth/logout?access_token_id=access-token-id&post_logout_redirect_uri=http://uri/logout-callback",
			ExpectedCode:     http.StatusFound,
			ExpectedRedirect: "/oauth",
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "DeleteSession")
				a.So(s.calls, should.Contain, "DeleteAccessToken")
				a.So(s.calls, should.Contain, "GetAccessToken")
				a.So(s.calls, should.Contain, "GetClient")
			},
		},
		{
			Name: "client-initiated logout without access token",
			StoreSetup: func(s *mockStore) {
				s.res.session = mockSession
				s.res.client = mockClient
				s.res.accessToken = mockAccessToken
			},
			Method:       "GET",
			Path:         "/oauth/logout",
			ExpectedCode: http.StatusForbidden,
		},
	} {
		name := tt.Name
		if name == "" {
			name = fmt.Sprintf("%s %s", tt.Method, tt.Path)
		}
		t.Run(name, func(t *testing.T) {
			store.reset()
			if tt.StoreSetup != nil {
				tt.StoreSetup(store)
			}

			req := httptest.NewRequest(tt.Method, tt.Path, nil)
			req.URL.Scheme, req.URL.Host = "http", req.Host

			var body *bytes.Buffer

			req.Header.Set("X-CSRF-Token", csrfToken)

			for _, c := range jar.Cookies(req.URL) {
				req.AddCookie(c)
			}

			if tt.UseCookie != nil {
				req.AddCookie(tt.UseCookie)
			}

			var contentType string
			switch b := tt.Body.(type) {
			case authorizeFormData:
				if b.encoding == "json" {
					json, _ := json.Marshal(b)
					body = bytes.NewBuffer(json)
					contentType = "application/json"
				}
				if b.encoding == "form" {
					values := url.Values{
						"authorize": []string{strconv.FormatBool(b.Authorize)},
					}
					if csrfToken != "" {
						values.Set("csrf", csrfToken)
					}
					body = bytes.NewBuffer([]byte(values.Encode()))
					contentType = "application/x-www-form-urlencoded"
				}
			}

			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}
			if body != nil {
				req.Body = ioutil.NopCloser(body)
				req.ContentLength = int64(body.Len())
			}

			res := httptest.NewRecorder()

			c.ServeHTTP(res, req)

			a := assertions.New(t)
			a.So(res.Code, should.Equal, tt.ExpectedCode)
			a.So(res.Header().Get("location"), should.ContainSubstring, tt.ExpectedRedirect)
			if tt.ExpectedBody != "" {
				if a.So(res.Body, should.NotBeNil) {
					a.So(res.Body.String(), should.ContainSubstring, tt.ExpectedBody)
				}
			}
			if cookies := res.Result().Cookies(); len(cookies) > 0 {
				jar.SetCookies(req.URL, cookies)
			}

			if tt.StoreCheck != nil {
				tt.StoreCheck(t, store)
			}
		})
	}
}

func TestTokenExchange(t *testing.T) {
	store := &mockStore{}
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			HTTP: config.HTTP{
				Cookie: config.Cookie{
					HashKey:  []byte("12345678123456781234567812345678"),
					BlockKey: []byte("12345678123456781234567812345678"),
				},
			},
		},
	})
	s, err := oauth.NewServer(c, store, oauth.Config{
		Mount: "/oauth",
		UI: oauth.UIConfig{
			TemplateData: webui.TemplateData{
				SiteName: "The Things Network",
				Title:    "OAuth",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	c.RegisterWeb(s)
	componenttest.StartComponent(t, c)

	for _, tt := range []struct {
		Name         string
		StoreSetup   func(*mockStore)
		StoreCheck   func(*testing.T, *mockStore)
		Method       string
		Path         string
		Body         interface{}
		ExpectedCode int
		ExpectedBody string
	}{
		{
			Name: "Exchange Authorization Code",
			StoreSetup: func(s *mockStore) {
				s.res.client = mockClient
				s.res.authorizationCode = &ttnpb.OAuthAuthorizationCode{
					UserIDs:       mockUser.UserIdentifiers,
					ClientIDs:     mockClient.ClientIdentifiers,
					UserSessionID: mockSession.SessionID,
					Rights:        mockClient.Rights,
					Code:          "the code",
					RedirectURI:   "http://uri/callback",
					State:         "foo",
					CreatedAt:     time.Now().Truncate(time.Second),
					ExpiresAt:     time.Now().Truncate(time.Second).Add(time.Hour),
				}
			},
			Method: "POST",
			Path:   "/oauth/token",
			Body: map[string]string{
				"grant_type":    "authorization_code",
				"code":          "the code",
				"redirect_uri":  "http://uri/callback",
				"client_id":     "client",
				"client_secret": "secret",
			},
			ExpectedCode: http.StatusOK,
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "GetAuthorizationCode")
				a.So(s.calls, should.Contain, "DeleteAuthorizationCode")
				a.So(s.req.code, should.Equal, "the code")
				a.So(s.calls, should.Contain, "CreateAccessToken")
				a.So(s.req.token.UserIDs, should.Resemble, mockUser.UserIdentifiers)
				a.So(s.req.token.ClientIDs, should.Resemble, mockClient.ClientIdentifiers)
				a.So(s.req.token.UserSessionID, should.Equal, mockSession.SessionID)
				a.So(s.req.token.Rights, should.Resemble, mockClient.Rights)
				a.So(s.req.token.AccessToken, should.NotBeEmpty)
				a.So(s.req.token.AccessToken, should.NotBeEmpty)
				a.So(s.req.token.RefreshToken, should.NotBeEmpty)
			},
		},
		{
			Name: "Exchange Refresh Token",
			StoreSetup: func(s *mockStore) {
				s.res.client = mockClient
				s.res.accessToken = &ttnpb.OAuthAccessToken{
					UserIDs:       mockUser.UserIdentifiers,
					ClientIDs:     mockClient.ClientIdentifiers,
					UserSessionID: mockSession.SessionID,
					ID:            "SFUBFRKYTGULGPAXXM4SHIBYMKCPTIMQBM63ZGQ",
					RefreshToken:  "PBKDF2$sha256$20000$IGAiKs46xX_M64E5$4xpyqnQT8SOa_Vf4xhEPk6WOZnhmAjG2mqGQiYBhm2s",
					Rights:        mockClient.Rights,
					CreatedAt:     time.Now().Truncate(time.Second),
					ExpiresAt:     time.Now().Truncate(time.Second).Add(time.Hour),
				}
			},
			Method: "POST",
			Path:   "/oauth/token",
			Body: map[string]string{
				"grant_type":    "refresh_token",
				"refresh_token": "OJSWM.IBTFXELDVVT64Y26IZZFFNSL7GWZY2Y3ALQQI3A.GCPIASDUP7UZJ6YL5OP2ESZB7CKRFV4JJQYTMDOSDIOE7O75IAMQ",
				"client_id":     "client",
				"client_secret": "secret",
			},
			ExpectedCode: http.StatusOK,
			StoreCheck: func(t *testing.T, s *mockStore) {
				a := assertions.New(t)
				a.So(s.calls, should.Contain, "GetAccessToken")
				a.So(s.calls, should.Contain, "DeleteAccessToken")
				a.So(s.req.tokenID, should.Equal, "IBTFXELDVVT64Y26IZZFFNSL7GWZY2Y3ALQQI3A")
				a.So(s.calls, should.Contain, "CreateAccessToken")
				a.So(s.req.token.UserIDs, should.Resemble, mockUser.UserIdentifiers)
				a.So(s.req.token.ClientIDs, should.Resemble, mockClient.ClientIdentifiers)
				a.So(s.req.token.Rights, should.Resemble, mockClient.Rights)
				a.So(s.req.token.AccessToken, should.NotBeEmpty)
				a.So(s.req.token.RefreshToken, should.NotBeEmpty)
				a.So(s.req.token.UserSessionID, should.Equal, mockSession.SessionID)
				a.So(s.req.previousID, should.Equal, "IBTFXELDVVT64Y26IZZFFNSL7GWZY2Y3ALQQI3A")
			},
		},
	} {
		name := tt.Name
		if name == "" {
			name = fmt.Sprintf("%s %s", tt.Method, tt.Path)
		}
		t.Run(name, func(t *testing.T) {
			store.reset()
			if tt.StoreSetup != nil {
				tt.StoreSetup(store)
			}

			var body *bytes.Buffer
			var contentType string

			if tt.Body != nil {
				switch ttBody := tt.Body.(type) {
				case url.Values:
					body = bytes.NewBufferString(ttBody.Encode())
					contentType = "application/x-www-form-urlencoded"
				default:
					json, _ := json.Marshal(tt.Body)
					body = bytes.NewBuffer(json)
					contentType = "application/json"
				}
			}

			req := httptest.NewRequest(tt.Method, tt.Path, body)
			req.URL.Scheme, req.URL.Host = "http", req.Host

			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}

			res := httptest.NewRecorder()

			c.ServeHTTP(res, req)

			a := assertions.New(t)
			a.So(res.Code, should.Equal, tt.ExpectedCode)
			if tt.ExpectedBody != "" {
				if a.So(res.Body, should.NotBeNil) {
					a.So(res.Body.String(), should.ContainSubstring, tt.ExpectedBody)
				}
			}

			if tt.StoreCheck != nil {
				tt.StoreCheck(t, store)
			}
		})
	}
}
