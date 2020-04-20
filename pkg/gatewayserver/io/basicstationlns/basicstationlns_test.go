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

package basicstationlns_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/gorilla/websocket"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/basicstation"
	"go.thethings.network/lorawan-stack/pkg/cluster"
	"go.thethings.network/lorawan-stack/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/pkg/component/test"
	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	. "go.thethings.network/lorawan-stack/pkg/gatewayserver/io/basicstationlns"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io/basicstationlns/messages"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io/mock"
	"go.thethings.network/lorawan-stack/pkg/log"
	pfconfig "go.thethings.network/lorawan-stack/pkg/pfconfig/basicstationlns"
	"go.thethings.network/lorawan-stack/pkg/pfconfig/shared"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

var (
	serverAddress          = "127.0.0.1:0"
	registeredGatewayUID   = "eui-0101010101010101"
	registeredGatewayID    = ttnpb.GatewayIdentifiers{GatewayID: registeredGatewayUID}
	registeredGateway      = ttnpb.Gateway{GatewayIdentifiers: registeredGatewayID, FrequencyPlanID: "EU_863_870"}
	registeredGatewayToken = "secrettoken"

	discoveryEndPoint      = "/router-info"
	connectionRootEndPoint = "/traffic/"

	testTrafficEndPoint = "/traffic/eui-0101010101010101"

	timeout               = (1 << 5) * test.Delay
	defaultWSPingInterval = (1 << 3) * test.Delay
)

func eui64Ptr(eui types.EUI64) *types.EUI64 { return &eui }

func TestClientTokenAuth(t *testing.T) {
	a := assertions.New(t)
	ctx := log.NewContext(test.Context(), test.GetLogger(t))
	ctx = newContextWithRightsFetcher(ctx)
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	is, isAddr := mock.NewIS(ctx)
	is.Add(ctx, registeredGatewayID, registeredGatewayToken)
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)
	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)
	gs := mock.NewServer(c)

	bsWebServer := New(ctx, gs, false, defaultWSPingInterval)
	lis, err := net.Listen("tcp", serverAddress)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	defer lis.Close()
	go func() error {
		return http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bsWebServer.ServeHTTP(w, r)
		}))
	}()
	servAddr := fmt.Sprintf("ws://%s", lis.Addr().String())

	for _, tc := range []struct {
		Name           string
		GatewayID      string
		AuthToken      string
		TokenPrefix    string
		ErrorAssertion func(err error) bool
	}{
		{
			Name:           "RegisteredGatewayAndValidKey",
			GatewayID:      registeredGatewayID.GatewayID,
			AuthToken:      registeredGatewayToken,
			ErrorAssertion: nil,
		},
		{
			Name:           "RegisteredGatewayAndValidKey",
			GatewayID:      registeredGatewayID.GatewayID,
			AuthToken:      registeredGatewayToken,
			TokenPrefix:    "Bearer ",
			ErrorAssertion: nil,
		},
		{
			Name:      "RegisteredGatewayAndInValidKey",
			GatewayID: registeredGatewayID.GatewayID,
			AuthToken: "invalidToken",
			ErrorAssertion: func(err error) bool {
				return err == websocket.ErrBadHandshake
			},
		},
		{
			Name:           "RegisteredGatewayAndNoKey",
			GatewayID:      registeredGatewayID.GatewayID,
			ErrorAssertion: nil,
		},
		{
			Name:      "UnregisteredGateway",
			GatewayID: "eui-1122334455667788",
			AuthToken: registeredGatewayToken,
			ErrorAssertion: func(err error) bool {
				return err == websocket.ErrBadHandshake
			},
		},
	} {
		t.Run(fmt.Sprintf("%s", tc.Name), func(t *testing.T) {
			a := assertions.New(t)
			h := http.Header{}
			h.Set("Authorization", fmt.Sprintf("%s%s", tc.TokenPrefix, tc.AuthToken))
			conn, _, err := websocket.DefaultDialer.Dial(servAddr+connectionRootEndPoint+tc.GatewayID, h)
			if err != nil {
				if tc.ErrorAssertion == nil || !a.So(tc.ErrorAssertion(err), should.BeTrue) {
					t.Fatalf("Unexpected error: %v", err)
				}
			} else if tc.ErrorAssertion != nil {
				t.Fatalf("Expected error")
			}
			if conn != nil {
				conn.Close()
			}
		})
	}
}

func TestClientSideTLS(t *testing.T) {
	// TODO: https://github.com/TheThingsNetwork/lorawan-stack/issues/558
}

func TestDiscover(t *testing.T) {
	a := assertions.New(t)
	ctx := log.NewContext(test.Context(), test.GetLogger(t))
	ctx = newContextWithRightsFetcher(ctx)
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	is, isAddr := mock.NewIS(ctx)
	is.Add(ctx, registeredGatewayID, registeredGatewayToken)
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)
	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)
	gs := mock.NewServer(c)

	bsWebServer := New(ctx, gs, false, defaultWSPingInterval)
	lis, err := net.Listen("tcp", serverAddress)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	defer lis.Close()
	go func() error {
		return http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bsWebServer.ServeHTTP(w, r)
		}))
	}()
	servAddr := fmt.Sprintf("ws://%s", lis.Addr().String())

	// Invalid Endpoints
	for i, tc := range []struct {
		URL string
	}{
		{
			URL: servAddr + "/api/v3/gs/io/basicstation/discover",
		},
		{
			URL: servAddr + discoveryEndPoint + "/eui-0101010101010101",
		},
	} {
		t.Run(fmt.Sprintf("InvalidDiscoveryEndPoint/%d", i), func(t *testing.T) {
			a := assertions.New(t)
			_, res, err := websocket.DefaultDialer.Dial(tc.URL, nil)
			if res.StatusCode != http.StatusNotFound {
				t.Fatalf("Unexpected response received: %v", res.Status)
			}
			if !a.So(err, should.Equal, websocket.ErrBadHandshake) {
				t.Fatalf("Connection failed: %v", err)
			}
		})
	}

	// Test Queries
	for _, tc := range []struct {
		Name     string
		Query    interface{}
		Response messages.DiscoverResponse
	}{
		{
			Name:     "EmptyEUI",
			Query:    messages.DiscoverQuery{},
			Response: messages.DiscoverResponse{Error: "Empty router EUI provided"},
		},
		{
			Name:     "EmptyStruct",
			Query:    struct{}{},
			Response: messages.DiscoverResponse{Error: "Empty router EUI provided"},
		},
		{
			Name: "InvalidJSONKey",
			Query: struct {
				EUI string `json:"route"`
			}{EUI: `"01-02-03-04-05-06-07-08"`},
			Response: messages.DiscoverResponse{Error: "Empty router EUI provided"},
		},
	} {
		t.Run(fmt.Sprintf("InvalidQuery/%s", tc.Name), func(t *testing.T) {
			a := assertions.New(t)
			conn, _, err := websocket.DefaultDialer.Dial(servAddr+discoveryEndPoint, nil)
			if !a.So(err, should.BeNil) {
				t.Fatalf("Connection failed: %v", err)
			}
			defer conn.Close()
			req, err := json.Marshal(tc.Query)
			if err != nil {
				panic(err)
			}
			if err := conn.WriteMessage(websocket.TextMessage, req); err != nil {
				t.Fatalf("Failed to write message: %v", err)
			}

			resCh := make(chan []byte)
			go func() {
				_, data, err := conn.ReadMessage()
				if err != nil {
					close(resCh)
					if err == websocket.ErrBadHandshake {
						return
					}
					t.Fatalf("Failed to read message: %v", err)
				}
				resCh <- data
			}()
			select {
			case res := <-resCh:
				var response messages.DiscoverResponse
				if err := json.Unmarshal(res, &response); err != nil {
					t.Fatalf("Failed to unmarshal response `%s`: %v", string(res), err)
				}
				a.So(response, should.Resemble, tc.Response)
			case <-time.After(timeout):
				t.Fatal("Read message timeout")
			}
		})
	}

	for _, tc := range []struct {
		Name  string
		Query interface{}
	}{
		{
			Name: "InvalidLength",
			Query: struct {
				EUI string `json:"router"`
			}{EUI: `"01-02-03-04-05-06-07-08-09"`},
		},
		{
			Name: "InvalidLength",
			Query: struct {
				EUI string `json:"router"`
			}{EUI: `"01:02:03:04:05:06:07:08:09"`},
		},
		{
			Name: "InvalidEUIFormat",
			Query: struct {
				EUI string `json:"router"`
			}{EUI: `"01:02:03:04:05:06:07-08"`},
		},
	} {
		t.Run(fmt.Sprintf("InvalidQuery/%s", tc.Name), func(t *testing.T) {
			a := assertions.New(t)
			conn, _, err := websocket.DefaultDialer.Dial(servAddr+discoveryEndPoint, nil)
			if !a.So(err, should.BeNil) {
				t.Fatalf("Connection failed: %v", err)
			}
			defer conn.Close()
			req, err := json.Marshal(tc.Query)
			if err != nil {
				panic(err)
			}
			if err := conn.WriteMessage(websocket.TextMessage, req); err != nil {
				t.Fatalf("Failed to write message: %v", err)
			}

			go func() {
				_, _, err := conn.ReadMessage()
				if err == nil {
					t.Fatalf("Expected connection closure with error but received none")
				}
			}()
		})
	}

	// Valid
	for i, tc := range []struct {
		EndPointEUI string
		EUI         types.EUI64
		Query       interface{}
	}{
		{
			EndPointEUI: "1111111111111111",
			EUI:         types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
			Query: messages.DiscoverQuery{
				EUI: basicstation.EUI{
					Prefix: "router",
					EUI64:  types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("ValidQuery/%d", i), func(t *testing.T) {
			a := assertions.New(t)
			conn, _, err := websocket.DefaultDialer.Dial(servAddr+discoveryEndPoint, nil)
			if !a.So(err, should.BeNil) {
				t.Fatalf("Connection failed: %v", err)
			}
			defer conn.Close()
			req, err := json.Marshal(tc.Query)
			if err != nil {
				panic(err)
			}
			if err := conn.WriteMessage(websocket.TextMessage, req); err != nil {
				t.Fatalf("Failed to write message: %v", err)
			}
			resCh := make(chan []byte)
			go func() {
				_, data, err := conn.ReadMessage()
				if err != nil {
					close(resCh)
					if err == websocket.ErrBadHandshake {
						return
					} else {
						t.Fatalf("Failed to read message: %v", err)
					}
				}
				resCh <- data
			}()
			select {
			case res := <-resCh:
				var response messages.DiscoverResponse
				if err := json.Unmarshal(res, &response); err != nil {
					t.Fatalf("Failed to unmarshal response `%s`: %v", string(res), err)
				}
				a.So(response, should.Resemble, messages.DiscoverResponse{
					EUI: basicstation.EUI{Prefix: "router", EUI64: tc.EUI},
					Muxs: basicstation.EUI{
						Prefix: "muxs",
					},
					URI: servAddr + connectionRootEndPoint + "eui-" + tc.EndPointEUI,
				})
			case <-time.After(timeout):
				t.Fatalf("Read message timeout")
			}
			conn.Close()
		})
	}
}

func TestVersion(t *testing.T) {
	a := assertions.New(t)
	ctx := log.NewContext(test.Context(), test.GetLogger(t))
	ctx = newContextWithRightsFetcher(ctx)
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	is, isAddr := mock.NewIS(ctx)
	is.Add(ctx, registeredGatewayID, registeredGatewayToken)
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)
	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)
	gs := mock.NewServer(c)

	bsWebServer := New(ctx, gs, false, defaultWSPingInterval)
	lis, err := net.Listen("tcp", serverAddress)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	defer lis.Close()
	go func() error {
		return http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bsWebServer.ServeHTTP(w, r)
		}))
	}()
	servAddr := fmt.Sprintf("ws://%s", lis.Addr().String())

	conn, _, err := websocket.DefaultDialer.Dial(servAddr+testTrafficEndPoint, nil)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	var gsConn *io.Connection
	select {
	case gsConn = <-gs.Connections():
	case <-time.After(timeout):
		t.Fatal("Connection timeout")
	}

	for _, tc := range []struct {
		Name                  string
		VersionQuery          interface{}
		ExpectedRouterConfig  interface{}
		ExpectedStatusMessage ttnpb.GatewayStatus
	}{
		{
			Name: "VersionProd",
			VersionQuery: messages.Version{
				Station:  "test-station",
				Firmware: "1.0.0",
				Package:  "test-package",
				Model:    "test-model",
				Protocol: 2,
				Features: "prod gps",
			},
			ExpectedRouterConfig: pfconfig.RouterConfig{
				Region:         "EU863",
				HardwareSpec:   "sx1301/1",
				FrequencyRange: []int{863000000, 870000000},
				DataRates: [16][3]int{
					{12, 125, 0},
					{11, 125, 0},
					{10, 125, 0},
					{9, 125, 0},
					{8, 125, 0},
					{7, 125, 0},
					{7, 250, 0},
					{0, 0, 0},
				},
				SX1301Config: []shared.SX1301Config{
					{
						LoRaWANPublic: true,
						ClockSource:   1,
						AntennaGain:   0,
						Radios: []shared.RFConfig{
							{
								Enable:     true,
								Frequency:  867500000,
								TxEnable:   true,
								RSSIOffset: -166,
							},
							{
								Enable:     true,
								Frequency:  868500000,
								TxEnable:   false,
								TxFreqMin:  0,
								TxFreqMax:  0,
								RSSIOffset: -166,
							},
						},
						Channels: []shared.IFConfig{
							{Enable: true, Radio: 0, IFValue: 600000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 800000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 1000000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: -400000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: -200000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 200000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 400000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
						},
						LoRaStandardChannel: &shared.IFConfig{Enable: true, Radio: 0, IFValue: 800000, Bandwidth: 250000, SpreadFactor: 7, Datarate: 0},
						FSKChannel:          &shared.IFConfig{Enable: true, Radio: 0, IFValue: 1300000, Bandwidth: 125000, SpreadFactor: 0, Datarate: 50000},
						TxLUTConfigs:        []shared.TxLUTConfig{},
					},
				},
			},
			ExpectedStatusMessage: ttnpb.GatewayStatus{
				Versions: map[string]string{
					"station":  "test-station",
					"firmware": "1.0.0",
					"package":  "test-package",
				},
				Advanced: &pbtypes.Struct{
					Fields: map[string]*pbtypes.Value{
						"model": {
							Kind: &pbtypes.Value_StringValue{StringValue: "test-model"},
						},
						"features": {
							Kind: &pbtypes.Value_StringValue{StringValue: "prod gps"},
						},
					},
				},
			},
		},
		{
			Name: "VersionDebug",
			VersionQuery: messages.Version{
				Station:  "test-station-rc1",
				Firmware: "1.0.0",
				Package:  "test-package",
				Model:    "test-model",
				Protocol: 2,
				Features: "rmtsh gps",
			},
			ExpectedRouterConfig: pfconfig.RouterConfig{
				Region:         "EU863",
				HardwareSpec:   "sx1301/1",
				FrequencyRange: []int{863000000, 870000000},
				DataRates: [16][3]int{
					{12, 125, 0},
					{11, 125, 0},
					{10, 125, 0},
					{9, 125, 0},
					{8, 125, 0},
					{7, 125, 0},
					{7, 250, 0},
					{0, 0, 0},
				},
				NoCCA:       true,
				NoDwellTime: true,
				NoDutyCycle: true,
				SX1301Config: []shared.SX1301Config{
					{
						LoRaWANPublic: true,
						ClockSource:   1,
						AntennaGain:   0,
						Radios: []shared.RFConfig{
							{
								Enable:     true,
								Frequency:  867500000,
								TxEnable:   true,
								RSSIOffset: -166,
							},
							{
								Enable:     true,
								Frequency:  868500000,
								TxEnable:   false,
								TxFreqMin:  0,
								TxFreqMax:  0,
								RSSIOffset: -166,
							},
						},
						Channels: []shared.IFConfig{
							{Enable: true, Radio: 0, IFValue: 600000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 800000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 1000000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: -400000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: -200000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 200000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: true, Radio: 0, IFValue: 400000, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
						},
						LoRaStandardChannel: &shared.IFConfig{Enable: true, Radio: 0, IFValue: 800000, Bandwidth: 250000, SpreadFactor: 7, Datarate: 0},
						FSKChannel:          &shared.IFConfig{Enable: true, Radio: 0, IFValue: 1300000, Bandwidth: 125000, SpreadFactor: 0, Datarate: 50000},
						TxLUTConfigs:        []shared.TxLUTConfig{},
					},
				},
			},
			ExpectedStatusMessage: ttnpb.GatewayStatus{
				Versions: map[string]string{
					"station":  "test-station-rc1",
					"firmware": "1.0.0",
					"package":  "test-package",
				},
				Advanced: &pbtypes.Struct{
					Fields: map[string]*pbtypes.Value{
						"model": {
							Kind: &pbtypes.Value_StringValue{StringValue: "test-model"},
						},
						"features": {
							Kind: &pbtypes.Value_StringValue{StringValue: "rmtsh gps"},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			reqVersion, err := json.Marshal(tc.VersionQuery)
			if err != nil {
				panic(err)
			}
			if err := conn.WriteMessage(websocket.TextMessage, reqVersion); err != nil {
				t.Fatalf("Failed to write message: %v", err)
			}

			resCh := make(chan []byte)
			go func() {
				_, data, err := conn.ReadMessage()
				if err != nil {
					t.Fatalf("Failed to read message: %v", err)
				}
				resCh <- data
			}()
			select {
			case res := <-resCh:
				var response pfconfig.RouterConfig
				if err := json.Unmarshal(res, &response); err != nil {
					t.Fatalf("Failed to unmarshal response `%s`: %v", string(res), err)
				}
				response.MuxTime = tc.ExpectedRouterConfig.(pfconfig.RouterConfig).MuxTime
				a.So(response, should.Resemble, tc.ExpectedRouterConfig)
			case <-time.After(timeout):
				t.Fatalf("Read message timeout")
			}
			select {
			case stat := <-gsConn.Status():
				a.So(time.Since(stat.Time), should.BeLessThan, timeout)
				stat.Time = time.Time{}
				a.So(stat, should.Resemble, &tc.ExpectedStatusMessage)
			case <-time.After(timeout):
				t.Fatalf("Read message timeout")
			}
		})
	}
}

func TestTraffic(t *testing.T) {
	a := assertions.New(t)
	ctx := log.NewContext(test.Context(), test.GetLogger(t))
	ctx = newContextWithRightsFetcher(ctx)
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	is, isAddr := mock.NewIS(ctx)
	is.Add(ctx, registeredGatewayID, registeredGatewayToken)
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)
	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)
	gs := mock.NewServer(c)

	bsWebServer := New(ctx, gs, false, defaultWSPingInterval)
	lis, err := net.Listen("tcp", serverAddress)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	defer lis.Close()
	go func() error {
		return http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bsWebServer.ServeHTTP(w, r)
		}))
	}()
	servAddr := fmt.Sprintf("ws://%s", lis.Addr().String())

	wsConn, _, err := websocket.DefaultDialer.Dial(servAddr+testTrafficEndPoint, nil)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Connection failed: %v", err)
	}
	defer wsConn.Close()

	var gsConn *io.Connection
	select {
	case gsConn = <-gs.Connections():
	case <-time.After(timeout):
		t.Fatal("Connection timeout")
	}

	for _, tc := range []struct {
		Name                    string
		InputBSUpstream         interface{}
		InputNetworkDownstream  *ttnpb.DownlinkMessage
		InputDownlinkPath       *ttnpb.DownlinkPath
		ExpectedBSDownstream    interface{}
		ExpectedNetworkUpstream interface{}
	}{
		{
			Name: "JoinRequest",
			InputBSUpstream: messages.JoinRequest{
				MHdr:     0,
				DevEUI:   basicstation.EUI{Prefix: "DevEui", EUI64: types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
				JoinEUI:  basicstation.EUI{Prefix: "JoinEui", EUI64: types.EUI64{0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22}},
				DevNonce: 18000,
				MIC:      12345678,
				RadioMetaData: messages.RadioMetaData{
					DataRate:  1,
					Frequency: 868300000,
					UpInfo: messages.UpInfo{
						RxTime: 1548059982,
						XTime:  12666373963464220,
						RSSI:   89,
						SNR:    9.25,
					},
				},
			},
			ExpectedNetworkUpstream: ttnpb.UplinkMessage{
				Payload: &ttnpb.Message{
					MHDR: ttnpb.MHDR{MType: ttnpb.MType_JOIN_REQUEST, Major: ttnpb.Major_LORAWAN_R1},
					MIC:  []byte{0x4E, 0x61, 0xBC, 0x00},
					Payload: &ttnpb.Message_JoinRequestPayload{JoinRequestPayload: &ttnpb.JoinRequestPayload{
						JoinEUI:  types.EUI64{0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22},
						DevEUI:   types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
						DevNonce: [2]byte{0x46, 0x50},
					}},
				},
				RxMetadata: []*ttnpb.RxMetadata{{
					GatewayIdentifiers: ttnpb.GatewayIdentifiers{
						GatewayID: "eui-0101010101010101",
						EUI:       &types.EUI64{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
					},
					Time:        &[]time.Time{time.Unix(1548059982, 0)}[0],
					Timestamp:   (uint32)(12666373963464220 & 0xFFFFFFFF),
					RSSI:        89,
					ChannelRSSI: 89,
					SNR:         9.25,
				}},
				Settings: ttnpb.TxSettings{
					Frequency:  868300000,
					CodingRate: "4/5",
					Time:       &[]time.Time{time.Unix(1548059982, 0)}[0],
					Timestamp:  (uint32)(12666373963464220 & 0xFFFFFFFF),
					DataRate: ttnpb.DataRate{Modulation: &ttnpb.DataRate_LoRa{LoRa: &ttnpb.LoRaDataRate{
						SpreadingFactor: 11,
						Bandwidth:       125000,
					}}},
				},
			},
		},
		{
			Name: "UplinkFrame",
			InputBSUpstream: messages.UplinkDataFrame{
				MHdr:       0x40,
				DevAddr:    0x11223344,
				FCtrl:      0x30,
				FPort:      0x00,
				FCnt:       25,
				FOpts:      "FD",
				FRMPayload: "5fcc",
				MIC:        12345678,
				RadioMetaData: messages.RadioMetaData{
					DataRate:  1,
					Frequency: 868300000,
					UpInfo: messages.UpInfo{
						RxTime: 1548059982,
						XTime:  12666373963464220,
						RSSI:   89,
						SNR:    9.25,
					},
				},
			},
			ExpectedNetworkUpstream: ttnpb.UplinkMessage{
				Payload: &ttnpb.Message{
					MHDR: ttnpb.MHDR{MType: ttnpb.MType_UNCONFIRMED_UP, Major: ttnpb.Major_LORAWAN_R1},
					MIC:  []byte{0x4E, 0x61, 0xBC, 0x00},
					Payload: &ttnpb.Message_MACPayload{MACPayload: &ttnpb.MACPayload{
						FPort:      0,
						FRMPayload: []byte{0x5F, 0xCC},
						FHDR: ttnpb.FHDR{
							DevAddr: [4]byte{0x11, 0x22, 0x33, 0x44},
							FCtrl: ttnpb.FCtrl{
								Ack:    true,
								ClassB: true,
							},
							FCnt:  25,
							FOpts: []byte{0xFD},
						},
					}},
				},
				RxMetadata: []*ttnpb.RxMetadata{
					{
						GatewayIdentifiers: ttnpb.GatewayIdentifiers{
							GatewayID: "eui-0101010101010101",
							EUI:       &types.EUI64{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
						},
						Time:        &[]time.Time{time.Unix(1548059982, 0)}[0],
						Timestamp:   (uint32)(12666373963464220 & 0xFFFFFFFF),
						RSSI:        89,
						ChannelRSSI: 89,
						SNR:         9.25,
					},
				},
				Settings: ttnpb.TxSettings{
					Frequency:  868300000,
					Time:       &[]time.Time{time.Unix(1548059982, 0)}[0],
					Timestamp:  (uint32)(12666373963464220 & 0xFFFFFFFF),
					CodingRate: "4/5",
					DataRate: ttnpb.DataRate{Modulation: &ttnpb.DataRate_LoRa{LoRa: &ttnpb.LoRaDataRate{
						SpreadingFactor: 11,
						Bandwidth:       125000,
					}}},
				},
			},
		},
		{
			Name: "Downlink",
			InputNetworkDownstream: &ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				EndDeviceIDs: &ttnpb.EndDeviceIdentifiers{
					DeviceID: "testdevice",
					DevEUI:   eui64Ptr(types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}),
				},
				Settings: &ttnpb.DownlinkMessage_Request{
					Request: &ttnpb.TxRequest{
						Class:            ttnpb.CLASS_A,
						Priority:         ttnpb.TxSchedulePriority_NORMAL,
						Rx1Delay:         ttnpb.RX_DELAY_1,
						Rx1DataRateIndex: 5,
						Rx1Frequency:     868100000,
						FrequencyPlanID:  test.EUFrequencyPlanID,
					},
				},
				CorrelationIDs: []string{"correlation1", "correlation2"},
			},

			InputDownlinkPath: &ttnpb.DownlinkPath{
				Path: &ttnpb.DownlinkPath_UplinkToken{
					UplinkToken: io.MustUplinkToken(ttnpb.GatewayAntennaIdentifiers{GatewayIdentifiers: registeredGatewayID}, 1553759666),
				},
			},
			ExpectedBSDownstream: messages.DownlinkMessage{
				DevEUI:      "00-00-00-00-00-00-00-00",
				DeviceClass: 0,
				Pdu:         "596d7868616d74686332356b4a334d3d3d",
				Diid:        1,
				RxDelay:     1,
				Rx1Freq:     868100000,
				Rx1DR:       5,
				XTime:       1553759666,
				Priority:    25,
				MuxTime:     1554300787.123456,
			},
		},
		{
			Name: "FollowUpTxAck",
			InputBSUpstream: messages.TxConfirmation{
				Diid:  1,
				XTime: 1548059982,
			},
			ExpectedNetworkUpstream: ttnpb.TxAcknowledgment{
				CorrelationIDs: []string{"correlation1", "correlation2"},
				Result:         ttnpb.TxAcknowledgment_SUCCESS,
			},
		},
		{
			Name: "RepeatedTxAck",
			InputBSUpstream: messages.TxConfirmation{
				Diid:  1,
				XTime: 1548059982,
			},
			ExpectedNetworkUpstream: ttnpb.TxAcknowledgment{
				CorrelationIDs: []string{"correlation1", "correlation2"},
				Result:         ttnpb.TxAcknowledgment_SUCCESS,
			},
		},
		{
			Name: "RandomTxAck",
			InputBSUpstream: messages.TxConfirmation{
				Diid:  2,
				XTime: 1548059982,
			},
			ExpectedNetworkUpstream: nil,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			if tc.InputBSUpstream != nil {
				switch v := tc.InputBSUpstream.(type) {
				case messages.TxConfirmation:
					req, err := json.Marshal(v)
					if err != nil {
						panic(err)
					}
					if err := wsConn.WriteMessage(websocket.TextMessage, req); err != nil {
						t.Fatalf("Failed to write message: %v", err)
					}
					select {
					case ack := <-gsConn.TxAck():
						if !a.So(*ack, should.Resemble, tc.ExpectedNetworkUpstream) {
							t.Fatalf("Invalid TxAck: %v", ack)
						}
					case <-time.After(timeout):
						if tc.ExpectedNetworkUpstream != nil {
							t.Fatalf("Read message timeout")
						}
					}

				case messages.UplinkDataFrame, messages.JoinRequest:
					req, err := json.Marshal(v)
					if err != nil {
						panic(err)
					}
					if err := wsConn.WriteMessage(websocket.TextMessage, req); err != nil {
						t.Fatalf("Failed to write message: %v", err)
					}
					select {
					case up := <-gsConn.Up():
						a.So(time.Since(up.ReceivedAt), should.BeLessThan, timeout)
						up.ReceivedAt = time.Time{}
						var payload ttnpb.Message
						a.So(lorawan.UnmarshalMessage(up.RawPayload, &payload), should.BeNil)
						if !a.So(&payload, should.Resemble, up.Payload) {
							t.Fatalf("Invalid RawPayload: %v", up.RawPayload)
						}
						up.RawPayload = nil
						up.RxMetadata[0].UplinkToken = nil
						expectedUp := tc.ExpectedNetworkUpstream.(ttnpb.UplinkMessage)
						a.So(up.UplinkMessage, should.Resemble, &expectedUp)
					case <-time.After(timeout):
						t.Fatalf("Read message timeout")
					}
				}
			}

			if tc.InputNetworkDownstream != nil {
				if _, err := gsConn.ScheduleDown(tc.InputDownlinkPath, tc.InputNetworkDownstream); err != nil {
					t.Fatalf("Failed to send downlink: %v", err)
				}

				resCh := make(chan []byte)
				go func() {
					_, data, err := wsConn.ReadMessage()
					if err != nil {
						t.Fatalf("Failed to read message: %v", err)
					}
					resCh <- data
				}()
				select {
				case res := <-resCh:
					switch tc.ExpectedBSDownstream.(type) {
					case messages.DownlinkMessage:
						var msg messages.DownlinkMessage
						if err := json.Unmarshal(res, &msg); err != nil {
							t.Fatalf("Failed to unmarshal response `%s`: %v", string(res), err)
						}
						msg.XTime = tc.ExpectedBSDownstream.(messages.DownlinkMessage).XTime
						msg.MuxTime = tc.ExpectedBSDownstream.(messages.DownlinkMessage).MuxTime
						if !a.So(msg, should.Resemble, tc.ExpectedBSDownstream.(messages.DownlinkMessage)) {
							t.Fatalf("Incorrect Downlink received: %s", string(res))
						}
					}
				case <-time.After(timeout):
					t.Fatalf("Read message timeout")
				}
			}
		})
	}
}

func TestRTT(t *testing.T) {
	a := assertions.New(t)
	ctx := log.NewContext(test.Context(), test.GetLogger(t))
	ctx = newContextWithRightsFetcher(ctx)
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	is, isAddr := mock.NewIS(ctx)
	is.Add(ctx, registeredGatewayID, registeredGatewayToken)
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)
	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)
	gs := mock.NewServer(c)

	bsWebServer := New(ctx, gs, false, defaultWSPingInterval)
	lis, err := net.Listen("tcp", serverAddress)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	defer lis.Close()
	go func() error {
		return http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bsWebServer.ServeHTTP(w, r)
		}))
	}()
	servAddr := fmt.Sprintf("ws://%s", lis.Addr().String())

	wsConn, _, err := websocket.DefaultDialer.Dial(servAddr+testTrafficEndPoint, nil)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Connection failed: %v", err)
	}
	defer wsConn.Close()

	var gsConn *io.Connection
	select {
	case gsConn = <-gs.Connections():
	case <-time.After(timeout):
		t.Fatal("Connection timeout")
	}

	var MuxTime, RxTime float64
	for _, tc := range []struct {
		Name                   string
		InputBSUpstream        interface{}
		InputNetworkDownstream *ttnpb.DownlinkMessage
		InputDownlinkPath      *ttnpb.DownlinkPath
		WaitTime               time.Duration
		ExpectedRTTStatsCount  int
	}{
		{
			Name: "JoinRequest",
			InputBSUpstream: messages.JoinRequest{
				MHdr:     0,
				DevEUI:   basicstation.EUI{Prefix: "DevEui", EUI64: types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
				JoinEUI:  basicstation.EUI{Prefix: "JoinEui", EUI64: types.EUI64{0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22}},
				DevNonce: 18000,
				MIC:      12345678,
				RadioMetaData: messages.RadioMetaData{
					DataRate:  1,
					Frequency: 868300000,
					UpInfo: messages.UpInfo{
						RxTime: 1548059982,
						XTime:  12666373963464220,
						RSSI:   89,
						SNR:    9.25,
					},
				},
			},
			ExpectedRTTStatsCount: 0,
		},
		{
			Name: "Downlink",
			InputNetworkDownstream: &ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				EndDeviceIDs: &ttnpb.EndDeviceIdentifiers{
					DeviceID: "testdevice",
					DevEUI:   eui64Ptr(types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}),
				},
				Settings: &ttnpb.DownlinkMessage_Request{
					Request: &ttnpb.TxRequest{
						Class:            ttnpb.CLASS_A,
						Priority:         ttnpb.TxSchedulePriority_NORMAL,
						Rx1Delay:         ttnpb.RX_DELAY_1,
						Rx1DataRateIndex: 5,
						Rx1Frequency:     868100000,
						FrequencyPlanID:  test.EUFrequencyPlanID,
					},
				},
				CorrelationIDs: []string{"correlation1", "correlation2"},
			},

			InputDownlinkPath: &ttnpb.DownlinkPath{
				Path: &ttnpb.DownlinkPath_UplinkToken{
					UplinkToken: io.MustUplinkToken(ttnpb.GatewayAntennaIdentifiers{GatewayIdentifiers: registeredGatewayID}, 1553759666),
				},
			},
		},
		{
			Name: "FollowUpTxAck",
			InputBSUpstream: messages.TxConfirmation{
				Diid:  1,
				XTime: 1548059982,
			},
			ExpectedRTTStatsCount: 1,
			WaitTime:              1 << 4 * test.Delay,
		},
		{
			Name: "RepeatedTxAck",
			InputBSUpstream: messages.TxConfirmation{
				Diid:  1,
				XTime: 1548059982,
			},
			ExpectedRTTStatsCount: 2,
			WaitTime:              0,
		},
		{
			Name: "UplinkFrame",
			InputBSUpstream: messages.UplinkDataFrame{
				MHdr:       0x40,
				DevAddr:    0x11223344,
				FCtrl:      0x30,
				FPort:      0x00,
				FCnt:       25,
				FOpts:      "FD",
				FRMPayload: "5fcc",
				MIC:        12345678,
				RadioMetaData: messages.RadioMetaData{
					DataRate:  1,
					Frequency: 868300000,
					UpInfo: messages.UpInfo{
						RxTime: 1548059982,
						XTime:  12666373963464220,
						RSSI:   89,
						SNR:    9.25,
					},
				},
			},
			ExpectedRTTStatsCount: 3,
			WaitTime:              1 << 3 * test.Delay,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			if tc.InputBSUpstream != nil {
				switch v := tc.InputBSUpstream.(type) {
				case messages.TxConfirmation:
					if MuxTime != 0 {
						time.Sleep(tc.WaitTime)
						now := float64(time.Now().UnixNano()) / float64(time.Second)
						v.RefTime = now - RxTime + MuxTime
					}
					req, err := json.Marshal(v)
					if err != nil {
						panic(err)
					}
					if err := wsConn.WriteMessage(websocket.TextMessage, req); err != nil {
						t.Fatalf("Failed to write message: %v", err)
					}
					select {
					case ack := <-gsConn.TxAck():
						if ack.Result != ttnpb.TxAcknowledgment_SUCCESS {
							t.Fatalf("Tx Acknowledgement failed")
						}
					case <-time.After(timeout):
						t.Fatalf("Read message timeout")
					}

				case messages.UplinkDataFrame:
					if MuxTime != 0 {
						time.Sleep(tc.WaitTime)
						now := float64(time.Now().UnixNano()) / float64(time.Second)
						v.RefTime = now - RxTime + MuxTime
					}
					req, err := json.Marshal(v)
					if err != nil {
						panic(err)
					}
					if err := wsConn.WriteMessage(websocket.TextMessage, req); err != nil {
						t.Fatalf("Failed to write message: %v", err)
					}
					select {
					case up := <-gsConn.Up():
						var payload ttnpb.Message
						a.So(lorawan.UnmarshalMessage(up.RawPayload, &payload), should.BeNil)
						if !a.So(&payload, should.Resemble, up.Payload) {
							t.Fatalf("Invalid RawPayload: %v", up.RawPayload)
						}
					case <-time.After(timeout):
						t.Fatalf("Read message timeout")
					}

				case messages.JoinRequest:
					if MuxTime != 0 {
						time.Sleep(tc.WaitTime)
						now := float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/(1e9)
						v.RefTime = now - RxTime + MuxTime
					}
					req, err := json.Marshal(v)
					if err != nil {
						panic(err)
					}
					if err := wsConn.WriteMessage(websocket.TextMessage, req); err != nil {
						t.Fatalf("Failed to write message: %v", err)
					}
					select {
					case up := <-gsConn.Up():
						var payload ttnpb.Message
						a.So(lorawan.UnmarshalMessage(up.RawPayload, &payload), should.BeNil)
						if !a.So(&payload, should.Resemble, up.Payload) {
							t.Fatalf("Invalid RawPayload: %v", up.RawPayload)
						}
					case <-time.After(timeout):
						t.Fatalf("Read message timeout")
					}
				}

				if MuxTime > 0 {
					// Atleast one downlink is needed for the first muxtime.
					min, max, median, count := gsConn.RTTStats()
					if !a.So(count, should.Equal, tc.ExpectedRTTStatsCount) {
						t.Fatalf("Incorrect Stats entries recorded: %d", count)
					}
					if !a.So(min, should.BeGreaterThan, 0) {
						t.Fatalf("Incorrect min: %s", min)
					}
					if tc.ExpectedRTTStatsCount > 1 {
						if !a.So(max, should.BeGreaterThan, min) {
							t.Fatalf("Incorrect max: %s", max)
						}
						if !a.So(median, should.BeBetween, min, max) {
							t.Fatalf("Incorrect median: %s", median)
						}
					}
				}
			}

			if tc.InputNetworkDownstream != nil {
				if _, err := gsConn.ScheduleDown(tc.InputDownlinkPath, tc.InputNetworkDownstream); err != nil {
					t.Fatalf("Failed to send downlink: %v", err)
				}

				resCh := make(chan []byte)
				go func() {
					_, data, err := wsConn.ReadMessage()
					if err != nil {
						t.Fatalf("Failed to read message: %v", err)
					}
					resCh <- data
				}()
				select {
				case res := <-resCh:
					var msg messages.DownlinkMessage
					if err := json.Unmarshal(res, &msg); err != nil {
						t.Fatalf("Failed to unmarshal response `%s`: %v", string(res), err)
					}
					MuxTime = msg.MuxTime
					RxTime = float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/(1e9)
				case <-time.After(timeout):
					t.Fatalf("Read message timeout")
				}
			}
		})
	}
}

func TestPingPong(t *testing.T) {
	a := assertions.New(t)
	ctx := log.NewContext(test.Context(), test.GetLogger(t))
	ctx = newContextWithRightsFetcher(ctx)
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	is, isAddr := mock.NewIS(ctx)
	is.Add(ctx, registeredGatewayID, registeredGatewayToken)
	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)
	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)
	gs := mock.NewServer(c)

	bsWebServer := New(ctx, gs, false, defaultWSPingInterval)
	lis, err := net.Listen("tcp", serverAddress)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	defer lis.Close()
	go func() error {
		return http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bsWebServer.ServeHTTP(w, r)
		}))
	}()
	servAddr := fmt.Sprintf("ws://%s", lis.Addr().String())

	conn, _, err := websocket.DefaultDialer.Dial(servAddr+testTrafficEndPoint, nil)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	pingCh := make(chan []byte)
	pongCh := make(chan []byte)

	// Read server ping
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			//  The ping/pong handlers are called only after ws.ReadMessage() receives a ping/pong message. The data read here is irrelevant.
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	conn.SetPingHandler(func(data string) error {
		pingCh <- []byte{}
		return nil
	})

	conn.SetPongHandler(func(data string) error {
		pongCh <- []byte{}
		return nil
	})

	select {
	case <-pingCh:
		t.Log("Received server ping")
		break
	case <-time.After(timeout):
		t.Fatalf("Server ping timeout")
	}

	// Client Ping, Server Pong
	if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		t.Fatalf("Failed to ping server: %v", err)
	}
	select {
	case <-pongCh:
		t.Log("Received server pong")
		break
	case <-time.After(timeout):
		t.Fatalf("Server pong timeout")
	}
}
