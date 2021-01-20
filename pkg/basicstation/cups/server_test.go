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

package cups

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestGetTrust(t *testing.T) {
	a := assertions.New(t)

	s := NewServer(nil)

	for _, addr := range []string{
		"thethingsnetwork.org:443",
		"https://thethingsnetwork.org:443",
		"https://thethingsnetwork.org",
	} {
		cert, err := s.getTrust(addr)
		a.So(err, should.BeNil)
		a.So(cert, should.NotBeNil)
	}
}

type mockGatewayClientData struct {
	ctx struct {
		GetIdentifiersForEUI context.Context
		Create               context.Context
		Get                  context.Context
		Update               context.Context
		CreateAPIKey         context.Context
		UpdateAPIKey         context.Context
	}
	req struct {
		GetIdentifiersForEUI *ttnpb.GetGatewayIdentifiersForEUIRequest
		Create               *ttnpb.CreateGatewayRequest
		Get                  *ttnpb.GetGatewayRequest
		Update               *ttnpb.UpdateGatewayRequest
		CreateAPIKey         *ttnpb.CreateGatewayAPIKeyRequest
		UpdateAPIKey         *ttnpb.UpdateGatewayAPIKeyRequest
	}
	opts struct {
		GetIdentifiersForEUI []grpc.CallOption
		Create               []grpc.CallOption
		Get                  []grpc.CallOption
		Update               []grpc.CallOption
		CreateAPIKey         []grpc.CallOption
		UpdateAPIKey         []grpc.CallOption
	}
	res struct {
		GetIdentifiersForEUI *ttnpb.GatewayIdentifiers
		Create               *ttnpb.Gateway
		Get                  *ttnpb.Gateway
		Update               *ttnpb.Gateway
		CreateAPIKey         *ttnpb.APIKey
		UpdateAPIKey         *ttnpb.APIKey
	}
	err struct {
		GetIdentifiersForEUI error
		Create               error
		Get                  error
		Update               error
		CreateAPIKey         error
		UpdateAPIKey         error
	}
}

type mockGatewayClient struct {
	mockGatewayClientData
	ttnpb.GatewayRegistryClient
	ttnpb.GatewayAccessClient
}

func (m *mockGatewayClient) reset() {
	m.mockGatewayClientData = mockGatewayClientData{}
}

func (m *mockGatewayClient) GetIdentifiersForEUI(ctx context.Context, in *ttnpb.GetGatewayIdentifiersForEUIRequest, opts ...grpc.CallOption) (*ttnpb.GatewayIdentifiers, error) {
	m.ctx.GetIdentifiersForEUI, m.req.GetIdentifiersForEUI, m.opts.GetIdentifiersForEUI = ctx, in, opts
	return m.res.GetIdentifiersForEUI, m.err.GetIdentifiersForEUI
}

func (m *mockGatewayClient) Create(ctx context.Context, in *ttnpb.CreateGatewayRequest, opts ...grpc.CallOption) (*ttnpb.Gateway, error) {
	m.ctx.Create, m.req.Create, m.opts.Create = ctx, in, opts
	return m.res.Create, m.err.Create
}

func (m *mockGatewayClient) Get(ctx context.Context, in *ttnpb.GetGatewayRequest, opts ...grpc.CallOption) (*ttnpb.Gateway, error) {
	m.ctx.Get, m.req.Get, m.opts.Get = ctx, in, opts
	return m.res.Get, m.err.Get
}

func (m *mockGatewayClient) Update(ctx context.Context, in *ttnpb.UpdateGatewayRequest, opts ...grpc.CallOption) (*ttnpb.Gateway, error) {
	m.ctx.Update, m.req.Update, m.opts.Update = ctx, in, opts
	return m.res.Update, m.err.Update
}

func (m *mockGatewayClient) CreateAPIKey(ctx context.Context, in *ttnpb.CreateGatewayAPIKeyRequest, opts ...grpc.CallOption) (*ttnpb.APIKey, error) {
	m.ctx.CreateAPIKey, m.req.CreateAPIKey, m.opts.CreateAPIKey = ctx, in, opts
	return m.res.CreateAPIKey, m.err.CreateAPIKey
}

func (m *mockGatewayClient) UpdateAPIKey(ctx context.Context, in *ttnpb.UpdateGatewayAPIKeyRequest, opts ...grpc.CallOption) (*ttnpb.APIKey, error) {
	m.ctx.UpdateAPIKey, m.req.UpdateAPIKey, m.opts.UpdateAPIKey = ctx, in, opts
	return m.res.UpdateAPIKey, m.err.UpdateAPIKey
}

const updateInfoRequest = `{
  "router": "58a0:cbff:fe80:19",
  "cupsUri": "https://thethingsnetwork.org:443",
  "tcUri": "",
  "cupsCredCrc": 1398343300,
  "tcCredCrc": 0,
  "station": "2.0.0(minihub/debug) 2018-12-06 09:30:35",
  "model": "minihub",
  "package": "2.0.0",
  "keys": [
    392840017
  ]
}`

var (
	mockFallbackAuth = grpc.PerRPCCredentials(nil)
	mockAuthFunc     = func(ctx context.Context) grpc.CallOption {
		return mockFallbackAuth
	}
	mockGatewayEUI    = types.EUI64{0x58, 0xA0, 0xCB, 0xFF, 0xFE, 0x80, 0x00, 0x19}
	mockErrNotFound   = grpc.Errorf(codes.NotFound, "not found")
	mockRightsFetcher = rights.FetcherFunc(func(ctx context.Context, ids ttnpb.Identifiers) (*ttnpb.Rights, error) {
		md := rpcmetadata.FromIncomingContext(ctx)
		if md.AuthType == "Bearer" {
			return ttnpb.RightsFrom(ttnpb.RIGHT_GATEWAY_INFO, ttnpb.RIGHT_GATEWAY_SETTINGS_BASIC, ttnpb.RIGHT_GATEWAY_READ_SECRETS), nil
		}
		return nil, rights.ErrNoGatewayRights
	})
)

func TestServer(t *testing.T) {
	tlsServer := httptest.NewTLSServer(http.HandlerFunc(http.NotFound))
	defer tlsServer.Close()
	tlsServerURL, _ := url.Parse(tlsServer.URL)

	cupsURI := (&url.URL{Scheme: "https", Host: tlsServerURL.Host}).String()
	lnsURI := (&url.URL{Scheme: "wss", Host: tlsServerURL.Host}).String()

	mockGateway := func(hasLNSSecret, redirectCUPS, updateCUPSCreds bool) *ttnpb.Gateway {
		secret := &ttnpb.Secret{
			KeyID: "test-key",
			Value: []byte("KEYCONTENTS"),
		}
		gtw := ttnpb.Gateway{
			GatewayIdentifiers: ttnpb.GatewayIdentifiers{
				GatewayID: "test-gateway",
				EUI:       &mockGatewayEUI,
			},
			Attributes: map[string]string{
				cupsStationAttribute: "2.0.0(minihub/debug) 2018-12-06 09:30:35",
				cupsModelAttribute:   "minihub",
				cupsPackageAttribute: "2.0.0",
			},
			GatewayServerAddress: lnsURI,
		}
		if hasLNSSecret {
			gtw.LBSLNSSecret = secret
		}
		if redirectCUPS {
			gtw.TargetCUPSURI = cupsURI
			gtw.TargetCUPSKey = secret
		}
		if updateCUPSCreds {
			gtw.TargetCUPSKey = secret
		}
		return &gtw
	}

	e := echo.New()

	for _, tt := range []struct {
		Name           string
		StoreSetup     func(*mockGatewayClient)
		Options        []Option
		RequestSetup   func(*http.Request)
		AssertError    func(error) bool
		AssertStore    func(*assertions.Assertion, *mockGatewayClient)
		AssertResponse func(*assertions.Assertion, *httptest.ResponseRecorder)
	}{
		{
			Name: "No Auth",
			StoreSetup: func(c *mockGatewayClient) {
				c.res.Get = mockGateway(false, false, false)
				c.res.GetIdentifiersForEUI = &c.res.Get.GatewayIdentifiers
			},
			RequestSetup: func(req *http.Request) {
				req.Header.Del(echo.HeaderAuthorization)
			},
			AssertError: func(err error) bool {
				return errors.IsUnauthenticated(err)
			},
		},
		{
			Name: "Not Found",
			StoreSetup: func(c *mockGatewayClient) {
				c.err.GetIdentifiersForEUI = mockErrNotFound
			},
			AssertError: func(err error) bool {
				return errors.IsNotFound(err)
			},
			AssertStore: func(a *assertions.Assertion, c *mockGatewayClient) {
				a.So(c.req.GetIdentifiersForEUI.EUI, should.Equal, mockGatewayEUI)
			},
		},
		{
			Name: "Register New Gateway",
			StoreSetup: func(c *mockGatewayClient) {
				c.err.GetIdentifiersForEUI = mockErrNotFound
				c.res.Create = &ttnpb.Gateway{
					GatewayIdentifiers: ttnpb.GatewayIdentifiers{
						GatewayID: "eui-58a0cbfffe800019",
						EUI:       &mockGatewayEUI,
					},
					LBSLNSSecret: &ttnpb.Secret{
						KeyID: "some-key-id",
						Value: []byte("KEYCONTENTS"),
					},
					TargetCUPSKey: &ttnpb.Secret{
						KeyID: "test-key",
						Value: []byte("KEYCONTENTS"),
					},
					TargetCUPSURI: cupsURI,
				}
				c.res.CreateAPIKey = &ttnpb.APIKey{
					ID:  "KEYID",
					Key: "KEYCONTENTS",
				}
				c.res.Get = c.res.Create
			},
			Options: []Option{
				WithRegisterUnknown(&ttnpb.OrganizationOrUserIdentifiers{}, mockAuthFunc),
				WithDefaultLNSURI(lnsURI),
				WithAllowCUPSURIUpdate(true),
			},
			AssertError: func(err error) bool {
				return err == nil
			},
			AssertResponse: func(a *assertions.Assertion, rec *httptest.ResponseRecorder) {
				var res UpdateInfoResponse
				err := res.UnmarshalBinary(rec.Body.Bytes())
				a.So(err, should.BeNil)
				a.So(res.CUPSURI, should.Equal, cupsURI)
				a.So(res.LNSURI, should.BeEmpty)
				a.So(res.CUPSCredentials, should.NotBeEmpty)
				a.So(res.LNSCredentials, should.BeEmpty)
				a.So(res.SignatureKeyCRC, should.BeZeroValue)
				a.So(res.Signature, should.BeEmpty)
				a.So(res.UpdateData, should.BeEmpty)
			},
			AssertStore: func(a *assertions.Assertion, s *mockGatewayClient) {
				if a.So(s.req.Create, should.NotBeNil) {
					a.So(s.req.Create.GatewayIdentifiers.GatewayID, should.Equal, "eui-58a0cbfffe800019")
					a.So(s.req.Create.GatewayIdentifiers.EUI, should.Resemble, &mockGatewayEUI)
				}
				if a.So(s.req.Update, should.NotBeNil) {
					a.So(s.req.Update.GatewayIdentifiers.GatewayID, should.Equal, "eui-58a0cbfffe800019")
					a.So(s.req.Update.GatewayIdentifiers.EUI, should.Resemble, &mockGatewayEUI)
					expectedAttributes := mockGateway(false, false, false).Attributes
					for _, attr := range []string{
						cupsStationAttribute,
						cupsModelAttribute,
						cupsPackageAttribute,
					} {
						a.So(s.req.Update.Attributes[attr], should.Equal, expectedAttributes[attr])
					}
				}
			},
		},
		{
			Name: "CUPS Not Enabled For Gateway",
			StoreSetup: func(c *mockGatewayClient) {
				c.res.Get = mockGateway(false, false, false)
				c.res.GetIdentifiersForEUI = &c.res.Get.GatewayIdentifiers
			},
			Options: []Option{
				WithExplicitEnable(true),
			},
			AssertError: func(err error) bool {
				return errors.IsUnauthenticated(err)
			},
		},
		{
			Name: "Existing Gateway",
			StoreSetup: func(c *mockGatewayClient) {
				c.res.Get = mockGateway(true, false, false)
				c.res.GetIdentifiersForEUI = &c.res.Get.GatewayIdentifiers
			},
			Options: []Option{
				WithAllowCUPSURIUpdate(true),
			},
			RequestSetup: func(req *http.Request) {
				req.Header.Set(echo.HeaderAuthorization, "Bearer KEYCONTENTS")
			},
			AssertError: func(err error) bool {
				return err == nil
			},
			AssertResponse: func(a *assertions.Assertion, rec *httptest.ResponseRecorder) {
				var res UpdateInfoResponse
				err := res.UnmarshalBinary(rec.Body.Bytes())
				a.So(err, should.BeNil)
				a.So(res.CUPSURI, should.BeEmpty)
				a.So(res.LNSURI, should.Equal, lnsURI)
				a.So(res.CUPSCredentials, should.BeEmpty)
				a.So(res.LNSCredentials, should.NotBeEmpty)
				a.So(res.SignatureKeyCRC, should.BeZeroValue)
				a.So(res.Signature, should.BeEmpty)
				a.So(res.UpdateData, should.BeEmpty)
			},
			AssertStore: func(a *assertions.Assertion, s *mockGatewayClient) {
				if a.So(s.req.Update, should.NotBeNil) {
					a.So(s.req.Update.GatewayIdentifiers.GatewayID, should.Equal, "test-gateway")
					a.So(s.req.Update.GatewayIdentifiers.EUI, should.Resemble, &mockGatewayEUI)
					expectedAttributes := mockGateway(true, false, false).Attributes
					for _, attr := range []string{
						cupsStationAttribute,
						cupsModelAttribute,
						cupsPackageAttribute,
					} {
						a.So(s.req.Update.Attributes[attr], should.Equal, expectedAttributes[attr])
					}
				}
			},
		},
		{
			Name: "Existing Gateway Without LNS Secret",
			StoreSetup: func(c *mockGatewayClient) {
				c.res.Get = mockGateway(false, false, false)
				c.res.GetIdentifiersForEUI = &c.res.Get.GatewayIdentifiers
			},
			Options: []Option{
				WithAllowCUPSURIUpdate(true),
			},
			RequestSetup: func(req *http.Request) {
				req.Header.Set(echo.HeaderAuthorization, "Bearer KEYCONTENTS")
			},
			AssertError: func(err error) bool {
				return errors.IsNotFound(err)
			},
		},
		{
			Name: "CUPS redirection",
			StoreSetup: func(c *mockGatewayClient) {
				c.res.Get = mockGateway(false, true, false)
				c.res.GetIdentifiersForEUI = &c.res.Get.GatewayIdentifiers
			},
			Options: []Option{
				WithAllowCUPSURIUpdate(true),
			},
			RequestSetup: func(req *http.Request) {
				req.Header.Set(echo.HeaderAuthorization, "Bearer KEYCONTENTS")
			},
			AssertError: func(err error) bool {
				return err == nil
			},
			AssertResponse: func(a *assertions.Assertion, rec *httptest.ResponseRecorder) {
				var res UpdateInfoResponse
				err := res.UnmarshalBinary(rec.Body.Bytes())
				a.So(err, should.BeNil)
				a.So(res.CUPSURI, should.Equal, cupsURI)
				a.So(res.LNSURI, should.BeEmpty)
				a.So(res.CUPSCredentials, should.NotBeEmpty)
				a.So(res.LNSCredentials, should.BeEmpty)
				a.So(res.SignatureKeyCRC, should.BeZeroValue)
				a.So(res.Signature, should.BeEmpty)
				a.So(res.UpdateData, should.BeEmpty)
			},
			AssertStore: func(a *assertions.Assertion, s *mockGatewayClient) {
				if a.So(s.req.Update, should.NotBeNil) {
					a.So(s.req.Update.GatewayIdentifiers.GatewayID, should.Equal, "test-gateway")
					a.So(s.req.Update.GatewayIdentifiers.EUI, should.Resemble, &mockGatewayEUI)
					expectedAttributes := mockGateway(false, true, false).Attributes
					for _, attr := range []string{
						cupsStationAttribute,
						cupsModelAttribute,
						cupsPackageAttribute,
					} {
						a.So(s.req.Update.Attributes[attr], should.Equal, expectedAttributes[attr])
					}
				}
			},
		},
		{
			Name: "CUPS Credentials Update",
			StoreSetup: func(c *mockGatewayClient) {
				c.res.Get = mockGateway(false, false, true)
				c.res.GetIdentifiersForEUI = &c.res.Get.GatewayIdentifiers
			},
			Options: []Option{
				WithAllowCUPSURIUpdate(true),
			},
			RequestSetup: func(req *http.Request) {
				req.Header.Set(echo.HeaderAuthorization, "Bearer KEYCONTENTS")
			},
			AssertError: func(err error) bool {
				return err == nil
			},
			AssertResponse: func(a *assertions.Assertion, rec *httptest.ResponseRecorder) {
				var res UpdateInfoResponse
				err := res.UnmarshalBinary(rec.Body.Bytes())
				a.So(err, should.BeNil)
				a.So(res.CUPSURI, should.Equal, "https://thethingsnetwork.org:443") //Same as the update-info request.
				a.So(res.LNSURI, should.BeEmpty)
				a.So(res.CUPSCredentials, should.NotBeEmpty)
				a.So(res.LNSCredentials, should.BeEmpty)
				a.So(res.SignatureKeyCRC, should.BeZeroValue)
				a.So(res.Signature, should.BeEmpty)
				a.So(res.UpdateData, should.BeEmpty)
			},
			AssertStore: func(a *assertions.Assertion, s *mockGatewayClient) {
				if a.So(s.req.Update, should.NotBeNil) {
					a.So(s.req.Update.GatewayIdentifiers.GatewayID, should.Equal, "test-gateway")
					a.So(s.req.Update.GatewayIdentifiers.EUI, should.Resemble, &mockGatewayEUI)
					expectedAttributes := mockGateway(false, false, true).Attributes
					for _, attr := range []string{
						cupsStationAttribute,
						cupsModelAttribute,
						cupsPackageAttribute,
					} {
						a.So(s.req.Update.Attributes[attr], should.Equal, expectedAttributes[attr])
					}
				}
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			a := assertions.New(t)
			store := &mockGatewayClient{}
			if tt.StoreSetup != nil {
				tt.StoreSetup(store)
			}

			s := NewServer(componenttest.NewComponent(t, &component.Config{}), append([]Option{
				WithTLSConfig(&tls.Config{
					InsecureSkipVerify: true,
				}),
				WithAuth(mockAuthFunc),
				WithRegistries(store, store),
			}, tt.Options...)...)
			req := httptest.NewRequest(http.MethodPost, "/update-info", strings.NewReader(updateInfoRequest))
			ctx := test.Context()
			ctx = log.NewContext(ctx, test.GetLogger(t))
			ctx = rights.NewContextWithFetcher(ctx, mockRightsFetcher)
			req = req.WithContext(ctx)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, "random string")
			if tt.RequestSetup != nil {
				tt.RequestSetup(req)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := s.UpdateInfo(c)
			if tt.AssertError != nil && !a.So(tt.AssertError(err), should.BeTrue) {
				t.Fatalf("Unexpected error :%v", err)
			}
			if tt.AssertResponse != nil {
				tt.AssertResponse(a, rec)
			}
			if tt.AssertStore != nil {
				tt.AssertStore(a, store)
			}
		})
	}
}
