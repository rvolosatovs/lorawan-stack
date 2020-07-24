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

package gatewayserver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/smartystreets/assertions"
	clusterauth "go.thethings.network/lorawan-stack/v3/pkg/auth/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/v3/pkg/errorcontext"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/basicstationlns"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/basicstationlns/messages"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/udp"
	gsredis "go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/upstream/mock"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcclient"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	encoding "go.thethings.network/lorawan-stack/v3/pkg/ttnpb/udp"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	registeredGatewayID  = "eui-aaee000000000000"
	registeredGatewayKey = "secret"
	registeredGatewayEUI = types.EUI64{0xAA, 0xEE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	unregisteredGatewayID  = "eui-bbff000000000000"
	unregisteredGatewayEUI = types.EUI64{0xBB, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	timeout        = (1 << 5) * test.Delay
	wsPingInterval = (1 << 3) * test.Delay
)

func TestGatewayServer(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()
	is, isAddr := startMockIS(ctx)
	ns, nsAddr := mock.StartNS(ctx)

	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":9187",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
				NetworkServer:  nsAddr,
			},
		},
	})
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)

	config := &gatewayserver.Config{
		RequireRegisteredGateways:         false,
		UpdateGatewayLocationDebounceTime: 0,
		MQTT: config.MQTT{
			Listen: ":1882",
		},
		UDP: gatewayserver.UDPConfig{
			Config: udp.Config{
				PacketHandlers:      2,
				PacketBuffer:        10,
				DownlinkPathExpires: 100 * time.Millisecond,
				ConnectionExpires:   250 * time.Millisecond,
				ScheduleLateTime:    0,
				AddrChangeBlock:     250 * time.Millisecond,
			},
			Listeners: map[string]string{
				":1700": test.EUFrequencyPlanID,
			},
		},
		BasicStation: gatewayserver.BasicStationConfig{
			Listen: ":1887",
			Config: basicstationlns.Config{
				WSPingInterval:       wsPingInterval,
				AllowUnauthenticated: true,
			},
		},
		UpdateConnectionStatsDebounceTime: 0,
	}
	if os.Getenv("TEST_REDIS") == "1" {
		statsRedisClient, statsFlush := test.NewRedis(t, "gatewayserver_test")
		defer statsFlush()
		defer statsRedisClient.Close()
		statsRegistry := &gsredis.GatewayConnectionStatsRegistry{
			Redis: statsRedisClient,
		}
		config.Stats = statsRegistry
	}

	gs, err := gatewayserver.New(c, config)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	roles := gs.Roles()
	a.So(len(roles), should.Equal, 1)
	a.So(roles[0], should.Equal, ttnpb.ClusterRole_GATEWAY_SERVER)

	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_NETWORK_SERVER)
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)

	time.Sleep(timeout) // Wait for tasks to be started.

	for _, publicLocation := range []bool{true, false} {
		is.add(ctx, ttnpb.GatewayIdentifiers{
			GatewayID: registeredGatewayID,
			EUI:       &registeredGatewayEUI,
		}, registeredGatewayKey, publicLocation, true)

		for _, ptc := range []struct {
			Protocol               string
			SupportsStatus         bool
			DetectsInvalidMessages bool
			DetectsDisconnect      bool
			TimeoutOnInvalidAuth   bool
			ValidAuth              func(ctx context.Context, ids ttnpb.GatewayIdentifiers, key string) bool
			Link                   func(ctx context.Context, t *testing.T, ids ttnpb.GatewayIdentifiers, key string, upCh <-chan *ttnpb.GatewayUp, downCh chan<- *ttnpb.GatewayDown) error
		}{
			{
				Protocol:          "grpc",
				SupportsStatus:    true,
				DetectsDisconnect: true,
				ValidAuth: func(ctx context.Context, ids ttnpb.GatewayIdentifiers, key string) bool {
					return ids.GatewayID == registeredGatewayID && key == registeredGatewayKey
				},
				Link: func(ctx context.Context, t *testing.T, ids ttnpb.GatewayIdentifiers, key string, upCh <-chan *ttnpb.GatewayUp, downCh chan<- *ttnpb.GatewayDown) error {
					conn, err := grpc.Dial(":9187", append(rpcclient.DefaultDialOptions(ctx), grpc.WithInsecure(), grpc.WithBlock())...)
					if err != nil {
						return err
					}
					defer conn.Close()
					md := rpcmetadata.MD{
						ID:            ids.GatewayID,
						AuthType:      "Bearer",
						AuthValue:     key,
						AllowInsecure: true,
					}
					client := ttnpb.NewGtwGsClient(conn)
					_, err = client.GetConcentratorConfig(ctx, ttnpb.Empty, grpc.PerRPCCredentials(md))
					if err != nil {
						return err
					}
					link, err := client.LinkGateway(ctx, grpc.PerRPCCredentials(md))
					if err != nil {
						return err
					}
					ctx, cancel := errorcontext.New(ctx)
					// Write upstream.
					go func() {
						for {
							select {
							case <-ctx.Done():
								return
							case msg := <-upCh:
								if err := link.Send(msg); err != nil {
									cancel(err)
									return
								}
							}
						}
					}()
					// Read downstream.
					go func() {
						for {
							msg, err := link.Recv()
							if err != nil {
								cancel(err)
								return
							}
							downCh <- msg
						}
					}()
					<-ctx.Done()
					return ctx.Err()
				},
			},
			{
				Protocol:             "mqtt",
				SupportsStatus:       true,
				DetectsDisconnect:    true,
				TimeoutOnInvalidAuth: true, // The MQTT client keeps reconnecting on invalid auth.
				ValidAuth: func(ctx context.Context, ids ttnpb.GatewayIdentifiers, key string) bool {
					return ids.GatewayID == registeredGatewayID && key == registeredGatewayKey
				},
				Link: func(ctx context.Context, t *testing.T, ids ttnpb.GatewayIdentifiers, key string, upCh <-chan *ttnpb.GatewayUp, downCh chan<- *ttnpb.GatewayDown) error {
					if ids.GatewayID == "" {
						t.SkipNow()
					}
					ctx, cancel := errorcontext.New(ctx)
					clientOpts := mqtt.NewClientOptions()
					clientOpts.AddBroker("tcp://0.0.0.0:1882")
					clientOpts.SetUsername(unique.ID(ctx, ids))
					clientOpts.SetPassword(key)
					clientOpts.SetAutoReconnect(false)
					clientOpts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
						cancel(err)
					})
					client := mqtt.NewClient(clientOpts)
					if token := client.Connect(); !token.WaitTimeout(timeout) {
						return context.DeadlineExceeded
					} else if err := token.Error(); err != nil {
						return err
					}
					defer client.Disconnect(uint(timeout / time.Millisecond))
					// Write upstream.
					go func() {
						for {
							select {
							case <-ctx.Done():
								return
							case up := <-upCh:
								for _, msg := range up.UplinkMessages {
									buf, err := msg.Marshal()
									if err != nil {
										cancel(err)
										return
									}
									if token := client.Publish(fmt.Sprintf("v3/%v/up", unique.ID(ctx, ids)), 1, false, buf); token.Wait() && token.Error() != nil {
										cancel(token.Error())
										return
									}
								}
								if up.GatewayStatus != nil {
									buf, err := up.GatewayStatus.Marshal()
									if err != nil {
										cancel(err)
										return
									}
									if token := client.Publish(fmt.Sprintf("v3/%v/status", unique.ID(ctx, ids)), 1, false, buf); token.Wait() && token.Error() != nil {
										cancel(token.Error())
										return
									}
								}
								if up.TxAcknowledgment != nil {
									buf, err := up.TxAcknowledgment.Marshal()
									if err != nil {
										cancel(err)
										return
									}
									if token := client.Publish(fmt.Sprintf("v3/%v/down/ack", unique.ID(ctx, ids)), 1, false, buf); token.Wait() && token.Error() != nil {
										cancel(token.Error())
										return
									}
								}
							}
						}
					}()
					// Read downstream.
					token := client.Subscribe(fmt.Sprintf("v3/%v/down", unique.ID(ctx, ids)), 1, func(_ mqtt.Client, raw mqtt.Message) {
						var msg ttnpb.GatewayDown
						if err := msg.Unmarshal(raw.Payload()); err != nil {
							cancel(err)
							return
						}
						downCh <- &msg
					})
					if token.Wait() && token.Error() != nil {
						return token.Error()
					}
					<-ctx.Done()
					return ctx.Err()
				},
			},
			{
				Protocol:       "udp",
				SupportsStatus: true,
				ValidAuth: func(ctx context.Context, ids ttnpb.GatewayIdentifiers, key string) bool {
					return ids.EUI != nil
				},
				Link: func(ctx context.Context, t *testing.T, ids ttnpb.GatewayIdentifiers, key string, upCh <-chan *ttnpb.GatewayUp, downCh chan<- *ttnpb.GatewayDown) error {
					if ids.EUI == nil {
						t.SkipNow()
					}
					upConn, err := net.Dial("udp", ":1700")
					if err != nil {
						return err
					}
					downConn, err := net.Dial("udp", ":1700")
					if err != nil {
						return err
					}
					ctx, cancel := errorcontext.New(ctx)
					// Write upstream.
					go func() {
						var token byte
						var readBuf [65507]byte
						for {
							select {
							case <-ctx.Done():
								return
							case up := <-upCh:
								token++
								packet := encoding.Packet{
									GatewayEUI:      ids.EUI,
									ProtocolVersion: encoding.Version1,
									Token:           [2]byte{0x00, token},
									PacketType:      encoding.PushData,
									Data:            &encoding.Data{},
								}
								packet.Data.RxPacket, packet.Data.Stat, packet.Data.TxPacketAck = encoding.FromGatewayUp(up)
								if packet.Data.TxPacketAck != nil {
									packet.PacketType = encoding.TxAck
								}
								writeBuf, err := packet.MarshalBinary()
								if err != nil {
									cancel(err)
									return
								}
								switch packet.PacketType {
								case encoding.PushData:
									if _, err := upConn.Write(writeBuf); err != nil {
										cancel(err)
										return
									}
									if _, err := upConn.Read(readBuf[:]); err != nil {
										cancel(err)
										return
									}
								case encoding.TxAck:
									if _, err := downConn.Write(writeBuf); err != nil {
										cancel(err)
										return
									}
								}
							}
						}
					}()
					// Engage downstream by sending PULL_DATA every 10ms.
					go func() {
						var token byte
						ticker := time.NewTicker(10 * time.Millisecond)
						for {
							select {
							case <-ctx.Done():
								ticker.Stop()
								return
							case <-ticker.C:
								token++
								pull := encoding.Packet{
									GatewayEUI:      ids.EUI,
									ProtocolVersion: encoding.Version1,
									Token:           [2]byte{0x01, token},
									PacketType:      encoding.PullData,
								}
								buf, err := pull.MarshalBinary()
								if err != nil {
									cancel(err)
									return
								}
								if _, err := downConn.Write(buf); err != nil {
									cancel(err)
									return
								}
							}
						}
					}()
					// Read downstream; PULL_RESP and PULL_ACK.
					go func() {
						var buf [65507]byte
						for {
							n, err := downConn.Read(buf[:])
							if err != nil {
								cancel(err)
								return
							}
							packetBuf := make([]byte, n)
							copy(packetBuf, buf[:])
							var packet encoding.Packet
							if err := packet.UnmarshalBinary(packetBuf); err != nil {
								cancel(err)
								return
							}
							switch packet.PacketType {
							case encoding.PullResp:
								msg, err := encoding.ToDownlinkMessage(packet.Data.TxPacket)
								if err != nil {
									cancel(err)
									return
								}
								downCh <- &ttnpb.GatewayDown{
									DownlinkMessage: msg,
								}
							}
						}
					}()
					<-ctx.Done()
					time.Sleep(config.UDP.ConnectionExpires * 150 / 100) // Ensure that connection expires.
					return ctx.Err()
				},
			},
			{
				Protocol:               "basicstation",
				SupportsStatus:         false,
				DetectsDisconnect:      true,
				DetectsInvalidMessages: true,
				ValidAuth: func(ctx context.Context, ids ttnpb.GatewayIdentifiers, key string) bool {
					return ids.EUI != nil
				},
				Link: func(ctx context.Context, t *testing.T, ids ttnpb.GatewayIdentifiers, key string, upCh <-chan *ttnpb.GatewayUp, downCh chan<- *ttnpb.GatewayDown) error {
					if ids.EUI == nil {
						t.SkipNow()
					}
					wsConn, _, err := websocket.DefaultDialer.Dial("ws://0.0.0.0:1887/traffic/"+registeredGatewayID, nil)
					if err != nil {
						return err
					}
					defer wsConn.Close()
					ctx, cancel := errorcontext.New(ctx)
					// Write upstream.
					go func() {
						for {
							select {
							case <-ctx.Done():
								return
							case msg := <-upCh:
								for _, uplink := range msg.UplinkMessages {
									var payload ttnpb.Message
									if err := lorawan.UnmarshalMessage(uplink.RawPayload, &payload); err != nil {
										// Ignore invalid uplinks
										continue
									}
									var bsUpstream []byte
									if payload.GetMType() == ttnpb.MType_JOIN_REQUEST {
										var jreq messages.JoinRequest
										err := jreq.FromUplinkMessage(uplink, test.EUFrequencyPlanID)
										if err != nil {
											cancel(err)
											return
										}
										bsUpstream, err = jreq.MarshalJSON()
										if err != nil {
											cancel(err)
											return
										}
									}
									if payload.GetMType() == ttnpb.MType_UNCONFIRMED_UP || payload.GetMType() == ttnpb.MType_CONFIRMED_UP {
										var updf messages.UplinkDataFrame
										err := updf.FromUplinkMessage(uplink, test.EUFrequencyPlanID)
										if err != nil {
											cancel(err)
											return
										}
										bsUpstream, err = updf.MarshalJSON()
										if err != nil {
											cancel(err)
											return
										}
									}
									if err := wsConn.WriteMessage(websocket.TextMessage, bsUpstream); err != nil {
										cancel(err)
										return
									}
								}
								if msg.TxAcknowledgment != nil {
									txConf := messages.TxConfirmation{
										Diid:  0,
										XTime: time.Now().Unix(),
									}
									bsUpstream, err := txConf.MarshalJSON()
									if err != nil {
										cancel(err)
										return
									}
									if err := wsConn.WriteMessage(websocket.TextMessage, bsUpstream); err != nil {
										cancel(err)
										return
									}
								}
							}
						}
					}()
					// Read downstream.
					go func() {
						for {
							_, data, err := wsConn.ReadMessage()
							if err != nil {
								cancel(err)
								return
							}
							var msg messages.DownlinkMessage
							if err := json.Unmarshal(data, &msg); err != nil {
								cancel(err)
								return
							}
							dlmesg := msg.ToDownlinkMessage()
							downCh <- &ttnpb.GatewayDown{
								DownlinkMessage: &dlmesg,
							}
						}
					}()
					<-ctx.Done()
					return ctx.Err()
				},
			},
		} {
			t.Run(fmt.Sprintf("Authenticate/%v", ptc.Protocol), func(t *testing.T) {
				for _, ctc := range []struct {
					Name string
					ID   ttnpb.GatewayIdentifiers
					Key  string
				}{
					{
						Name: "ValidIDAndKey",
						ID:   ttnpb.GatewayIdentifiers{GatewayID: registeredGatewayID},
						Key:  registeredGatewayKey,
					},
					{
						Name: "InvalidKey",
						ID:   ttnpb.GatewayIdentifiers{GatewayID: registeredGatewayID},
						Key:  "invalid-key",
					},
					{
						Name: "InvalidIDAndKey",
						ID:   ttnpb.GatewayIdentifiers{GatewayID: "invalid-gateway"},
						Key:  "invalid-key",
					},
					{
						Name: "RegisteredEUI",
						ID:   ttnpb.GatewayIdentifiers{EUI: &registeredGatewayEUI},
					},
					{
						Name: "UnregisteredEUI",
						ID:   ttnpb.GatewayIdentifiers{EUI: &unregisteredGatewayEUI},
					},
				} {
					t.Run(ctc.Name, func(t *testing.T) {
						ctx, cancel := context.WithCancel(ctx)
						upCh := make(chan *ttnpb.GatewayUp)
						downCh := make(chan *ttnpb.GatewayDown)

						upEvents := map[string]events.Channel{}
						for _, event := range []string{"gs.gateway.connect"} {
							upEvents[event] = make(events.Channel, 5)
						}
						defer test.SetDefaultEventsPubSub(&test.MockEventPubSub{
							PublishFunc: func(ev events.Event) {
								switch name := ev.Name(); name {
								case "gs.gateway.connect":
									go func() {
										upEvents[name] <- ev
									}()
								default:
									t.Logf("%s event published", name)
								}
							},
						})()

						connectedWithInvalidAuth := make(chan struct{}, 1)
						expectedProperLink := make(chan struct{}, 1)
						go func() {
							select {
							case <-upEvents["gs.gateway.connect"]:
								if !ptc.ValidAuth(ctx, ctc.ID, ctc.Key) {
									connectedWithInvalidAuth <- struct{}{}
								}
							case <-time.After(timeout):
								if ptc.ValidAuth(ctx, ctc.ID, ctc.Key) {
									expectedProperLink <- struct{}{}
								}
							}
							time.Sleep(test.Delay)
							cancel()
						}()
						err := ptc.Link(ctx, t, ctc.ID, ctc.Key, upCh, downCh)
						if !errors.IsCanceled(err) && ptc.ValidAuth(ctx, ctc.ID, ctc.Key) {
							t.Fatalf("Expect canceled context but have %v", err)
						}
						select {
						case <-connectedWithInvalidAuth:
							t.Fatal("Expected link error due to invalid auth")
						case <-expectedProperLink:
							t.Fatal("Expected proper link")
						default:
						}
					})
				}
			})

			// Wait for gateway disconnection to be processed.
			time.Sleep(timeout)

			t.Run(fmt.Sprintf("DetectDisconnect/%v", ptc.Protocol), func(t *testing.T) {
				if !ptc.DetectsDisconnect {
					t.SkipNow()
				}

				id := ttnpb.GatewayIdentifiers{
					GatewayID: registeredGatewayID,
					EUI:       &registeredGatewayEUI,
				}

				ctx1, fail1 := errorcontext.New(ctx)
				defer fail1(context.Canceled)
				go func() {
					upCh := make(chan *ttnpb.GatewayUp)
					downCh := make(chan *ttnpb.GatewayDown)
					err := ptc.Link(ctx1, t, id, registeredGatewayKey, upCh, downCh)
					fail1(err)
				}()
				select {
				case <-ctx1.Done():
					t.Fatalf("Expected no link error on first connection but have %v", ctx1.Err())
				case <-time.After(timeout):
				}

				ctx2, cancel2 := context.WithDeadline(ctx, time.Now().Add(4*timeout))
				upCh := make(chan *ttnpb.GatewayUp)
				downCh := make(chan *ttnpb.GatewayDown)
				err := ptc.Link(ctx2, t, id, registeredGatewayKey, upCh, downCh)
				cancel2()
				if !errors.IsDeadlineExceeded(err) {
					t.Fatalf("Expected deadline exceeded on second connection but have %v", err)
				}
				select {
				case <-ctx1.Done():
					t.Logf("First connection failed when second connected with %v", ctx1.Err())
				case <-time.After(4 * timeout):
					t.Fatalf("Expected link failure on first connection when second connected")
				}
			})

			// Wait for gateway disconnection to be processed.
			time.Sleep(3 * timeout)

			t.Run(fmt.Sprintf("Traffic/%v", ptc.Protocol), func(t *testing.T) {
				a := assertions.New(t)

				ctx, cancel := context.WithCancel(ctx)
				upCh := make(chan *ttnpb.GatewayUp)
				downCh := make(chan *ttnpb.GatewayDown)
				ids := ttnpb.GatewayIdentifiers{
					GatewayID: registeredGatewayID,
					EUI:       &registeredGatewayEUI,
				}
				// Setup a stats client with independent context to query whether the gateway is connected and statistics on
				// upstream and downstream.
				statsConn, err := grpc.Dial(":9187", append(rpcclient.DefaultDialOptions(test.Context()), grpc.WithInsecure(), grpc.WithBlock())...)
				if !a.So(err, should.BeNil) {
					t.FailNow()
				}
				defer statsConn.Close()
				statsCtx := metadata.AppendToOutgoingContext(test.Context(),
					"id", ids.GatewayID,
					"authorization", fmt.Sprintf("Bearer %v", registeredGatewayKey),
				)
				statsClient := ttnpb.NewGsClient(statsConn)

				// The gateway should not be connected before testing traffic.
				t.Run("NotConnected", func(t *testing.T) {
					_, err := statsClient.GetGatewayConnectionStats(statsCtx, &ids)
					if !a.So(errors.IsNotFound(err), should.BeTrue) {
						t.Fatal("Expected gateway not to be connected yet, but it is")
					}
				})

				if ptc.SupportsStatus {
					t.Run("UpdateLocation", func(t *testing.T) {
						for _, tc := range []struct {
							Name           string
							UpdateLocation bool
							Up             *ttnpb.GatewayUp
							ExpectLocation ttnpb.Location
						}{
							{
								Name:           "NoUpdate",
								UpdateLocation: false,
								Up: &ttnpb.GatewayUp{
									GatewayStatus: &ttnpb.GatewayStatus{
										Time: time.Unix(424242, 0),
										AntennaLocations: []*ttnpb.Location{
											{
												Source:    ttnpb.SOURCE_GPS,
												Altitude:  10,
												Latitude:  12,
												Longitude: 14,
											},
										},
									},
								},
								ExpectLocation: ttnpb.Location{
									Source: ttnpb.SOURCE_GPS,
								},
							},
							{
								Name:           "NoLocation",
								UpdateLocation: true,
								Up: &ttnpb.GatewayUp{
									GatewayStatus: &ttnpb.GatewayStatus{
										Time: time.Unix(424242, 0),
									},
								},
								ExpectLocation: ttnpb.Location{
									Source: ttnpb.SOURCE_GPS,
								},
							},
							{
								Name:           "Update",
								UpdateLocation: true,
								Up: &ttnpb.GatewayUp{
									GatewayStatus: &ttnpb.GatewayStatus{
										Time: time.Unix(42424242, 0),
										AntennaLocations: []*ttnpb.Location{
											{
												Source:    ttnpb.SOURCE_GPS,
												Altitude:  10,
												Latitude:  12,
												Longitude: 14,
											},
										},
									},
								},
								ExpectLocation: ttnpb.Location{
									Source:    ttnpb.SOURCE_GPS,
									Altitude:  10,
									Latitude:  12,
									Longitude: 14,
								},
							},
						} {
							t.Run(tc.Name, func(t *testing.T) {
								a := assertions.New(t)

								gtw, err := is.Get(ctx, &ttnpb.GetGatewayRequest{
									GatewayIdentifiers: ids,
								})
								a.So(err, should.BeNil)
								gtw.Antennas[0].Location = ttnpb.Location{
									Source: ttnpb.SOURCE_GPS,
								}
								gtw.UpdateLocationFromStatus = tc.UpdateLocation
								gtw, err = is.Get(ctx, &ttnpb.GetGatewayRequest{
									GatewayIdentifiers: ids,
								})
								a.So(err, should.BeNil)
								a.So(gtw.UpdateLocationFromStatus, should.Equal, tc.UpdateLocation)

								ctx, cancel := context.WithCancel(ctx)
								upCh := make(chan *ttnpb.GatewayUp)
								downCh := make(chan *ttnpb.GatewayDown)

								wg := &sync.WaitGroup{}
								wg.Add(1)
								go func() {
									defer wg.Done()
									err := ptc.Link(ctx, t, ids, registeredGatewayKey, upCh, downCh)
									if !errors.IsCanceled(err) {
										t.Fatalf("Expected context canceled, but have %v", err)
									}
								}()

								select {
								case upCh <- tc.Up:
								case <-time.After(timeout):
									t.Fatalf("Failed to send message to upstream channel")
								}

								time.Sleep(timeout)
								gtw, err = is.Get(ctx, &ttnpb.GetGatewayRequest{
									GatewayIdentifiers: ids,
								})
								a.So(err, should.BeNil)
								a.So(gtw.Antennas[0].Location, should.Resemble, tc.ExpectLocation)

								cancel()
								wg.Wait()
							})
						}
					})
				}
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := ptc.Link(ctx, t, ids, registeredGatewayKey, upCh, downCh)
					if !errors.IsCanceled(err) {
						t.Fatalf("Expected context canceled, but have %v", err)
					}
				}()

				// Expected location for RxMetadata
				gtw, err := is.Get(ctx, &ttnpb.GetGatewayRequest{
					GatewayIdentifiers: ids,
				})
				location := &gtw.Antennas[0].Location
				if !publicLocation {
					location = nil
				}

				t.Run("Upstream", func(t *testing.T) {
					uplinkCount := 0
					for _, tc := range []struct {
						Name     string
						Up       *ttnpb.GatewayUp
						Forwards []uint32 // Timestamps of uplink messages in Up that are being forwarded.
					}{
						{
							Name: "GatewayStatus",
							Up: &ttnpb.GatewayUp{
								GatewayStatus: &ttnpb.GatewayStatus{
									Time: time.Unix(424242, 0),
								},
							},
						},
						{
							Name: "TxAck",
							Up: &ttnpb.GatewayUp{
								TxAcknowledgment: &ttnpb.TxAcknowledgment{
									Result: ttnpb.TxAcknowledgment_SUCCESS,
								},
							},
						},
						{
							Name: "OneValidLoRa",
							Up: &ttnpb.GatewayUp{
								UplinkMessages: []*ttnpb.UplinkMessage{
									{
										Settings: ttnpb.TxSettings{
											DataRate: ttnpb.DataRate{
												Modulation: &ttnpb.DataRate_LoRa{
													LoRa: &ttnpb.LoRaDataRate{
														SpreadingFactor: 7,
														Bandwidth:       250000,
													},
												},
											},
											CodingRate: "4/5",
											Frequency:  867900000,
											Timestamp:  100,
										},
										RxMetadata: []*ttnpb.RxMetadata{
											{
												GatewayIdentifiers: ids,
												Timestamp:          100,
												RSSI:               -69,
												ChannelRSSI:        -69,
												SNR:                11,
												Location:           location,
											},
										},
										RawPayload: randomUpDataPayload(types.DevAddr{0x26, 0x01, 0xff, 0xff}, 1, 6),
									},
								},
							},
							Forwards: []uint32{100},
						},
						{
							Name: "OneValidFSK",
							Up: &ttnpb.GatewayUp{
								UplinkMessages: []*ttnpb.UplinkMessage{
									{
										Settings: ttnpb.TxSettings{
											DataRate: ttnpb.DataRate{
												Modulation: &ttnpb.DataRate_FSK{
													FSK: &ttnpb.FSKDataRate{
														BitRate: 50000,
													},
												},
											},
											Frequency: 867900000,
											Timestamp: 100,
										},
										RxMetadata: []*ttnpb.RxMetadata{
											{
												GatewayIdentifiers: ids,
												Timestamp:          100,
												RSSI:               -69,
												ChannelRSSI:        -69,
												SNR:                11,
												Location:           location,
											},
										},
										RawPayload: randomUpDataPayload(types.DevAddr{0x26, 0x01, 0xff, 0xff}, 1, 6),
									},
								},
							},
							Forwards: []uint32{100},
						},
						{
							Name: "OneGarbageWithStatus",
							Up: &ttnpb.GatewayUp{
								UplinkMessages: []*ttnpb.UplinkMessage{
									{
										Settings: ttnpb.TxSettings{
											DataRate: ttnpb.DataRate{
												Modulation: &ttnpb.DataRate_LoRa{
													LoRa: &ttnpb.LoRaDataRate{
														SpreadingFactor: 9,
														Bandwidth:       125000,
													},
												},
											},
											CodingRate: "4/5",
											Frequency:  868500000,
											Timestamp:  100,
										},
										RxMetadata: []*ttnpb.RxMetadata{
											{
												GatewayIdentifiers: ids,
												Timestamp:          100,
												RSSI:               -112,
												ChannelRSSI:        -112,
												SNR:                2,
												Location:           location,
											},
										},
										RawPayload: []byte{0xff, 0x02, 0x03}, // Garbage; doesn't get forwarded.
									},
									{
										Settings: ttnpb.TxSettings{
											DataRate: ttnpb.DataRate{
												Modulation: &ttnpb.DataRate_LoRa{
													LoRa: &ttnpb.LoRaDataRate{
														SpreadingFactor: 7,
														Bandwidth:       125000,
													},
												},
											},
											CodingRate: "4/5",
											Frequency:  868100000,
											Timestamp:  200,
										},
										RxMetadata: []*ttnpb.RxMetadata{
											{
												GatewayIdentifiers: ids,
												Timestamp:          200,
												RSSI:               -69,
												ChannelRSSI:        -69,
												SNR:                11,
												Location:           location,
											},
										},
										RawPayload: randomUpDataPayload(types.DevAddr{0x26, 0x01, 0xff, 0xff}, 1, 6),
									},
									{
										Settings: ttnpb.TxSettings{
											DataRate: ttnpb.DataRate{
												Modulation: &ttnpb.DataRate_LoRa{
													LoRa: &ttnpb.LoRaDataRate{
														SpreadingFactor: 12,
														Bandwidth:       125000,
													},
												},
											},
											CodingRate: "4/5",
											Frequency:  867700000,
											Timestamp:  300,
										},
										RxMetadata: []*ttnpb.RxMetadata{
											{
												GatewayIdentifiers: ids,
												Timestamp:          300,
												RSSI:               -36,
												ChannelRSSI:        -36,
												SNR:                5,
												Location:           location,
											},
										},
										RawPayload: randomJoinRequestPayload(
											types.EUI64{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
											types.EUI64{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
										),
									},
								},
								GatewayStatus: &ttnpb.GatewayStatus{
									Time: time.Unix(4242424, 0),
								},
							},
							Forwards: []uint32{200, 300},
						},
					} {
						t.Run(tc.Name, func(t *testing.T) {
							a := assertions.New(t)

							upEvents := map[string]events.Channel{}
							for _, event := range []string{"gs.up.receive", "gs.down.tx.success", "gs.down.tx.fail", "gs.status.receive", "gs.status.forward"} {
								upEvents[event] = make(events.Channel, 5)
							}
							defer test.SetDefaultEventsPubSub(&test.MockEventPubSub{
								PublishFunc: func(ev events.Event) {
									switch name := ev.Name(); name {
									case "gs.up.receive", "gs.down.tx.success", "gs.down.tx.fail", "gs.status.receive", "gs.status.forward":
										go func() {
											upEvents[name] <- ev
										}()
									default:
										t.Logf("%s event published", name)
									}
								},
							})()

							select {
							case upCh <- tc.Up:
							case <-time.After(timeout):
								t.Fatalf("Failed to send message to upstream channel")
							}
							if ptc.DetectsInvalidMessages {
								uplinkCount += len(tc.Forwards)
							} else {
								uplinkCount += len(tc.Up.UplinkMessages)
							}

							notSeen := make(map[uint32]struct{})
							for _, t := range tc.Forwards {
								notSeen[t] = struct{}{}
							}
							for len(notSeen) > 0 {
								select {
								case msg := <-ns.Up():
									var expected *ttnpb.UplinkMessage
									for _, up := range tc.Up.UplinkMessages {
										if ts := up.Settings.Timestamp; ts == msg.Settings.Timestamp {
											if _, ok := notSeen[ts]; !ok {
												t.Fatalf("Not expecting message %v", msg)
											}
											expected = up
											delete(notSeen, ts)
											break
										}
									}
									if expected == nil {
										t.Fatalf("Received unexpected message")
									}
									a.So(time.Since(msg.ReceivedAt), should.BeLessThan, timeout)
									a.So(msg.Settings, should.Resemble, expected.Settings)
									for _, md := range msg.RxMetadata {
										a.So(md.UplinkToken, should.NotBeEmpty)
										md.UplinkToken = nil
									}
									a.So(msg.RxMetadata, should.Resemble, expected.RxMetadata)
									a.So(msg.RawPayload, should.Resemble, expected.RawPayload)
								case <-time.After(timeout):
									t.Fatal("Expected uplink timeout")
								}
								select {
								case evt := <-upEvents["gs.up.receive"]:
									a.So(evt.Name(), should.Equal, "gs.up.receive")
								case <-time.After(timeout):
									t.Fatal("Expected uplink event timeout")
								}
							}
							if expected := tc.Up.TxAcknowledgment; expected != nil {
								select {
								case <-upEvents["gs.down.tx.success"]:
								case evt := <-upEvents["gs.down.tx.fail"]:
									received, ok := evt.Data().(ttnpb.TxAcknowledgment_Result)
									if !ok {
										t.Fatal("No acknowledgment attached to the downlink emission fail event")
									}
									a.So(received, should.Resemble, expected.Result)
								case <-time.After(timeout):
									t.Fatal("Expected Tx acknowledgment event timeout")
								}
							}
							if tc.Up.GatewayStatus != nil && ptc.SupportsStatus {
								select {
								case <-upEvents["gs.status.receive"]:
								case <-time.After(timeout):
									t.Fatal("Expected gateway status event timeout")
								}

								select {
								case <-upEvents["gs.status.forward"]:
								case <-time.After(timeout):
									t.Fatal("Expected gateway status forward event timeout")
								}
							}

							time.Sleep(2 * timeout)

							conn, ok := gs.GetConnection(ctx, ids)
							a.So(ok, should.BeTrue)
							a.So(conn.Stats(), should.NotBeNil)
							if config.Stats != nil {
								a.So(gs.UpdateConnectionStats(conn), should.BeNil)
							}

							stats, err := statsClient.GetGatewayConnectionStats(statsCtx, &ids)
							if !a.So(err, should.BeNil) {
								t.FailNow()
							}
							a.So(stats.UplinkCount, should.Equal, uplinkCount)

							if tc.Up.GatewayStatus != nil && ptc.SupportsStatus {
								if !a.So(stats.LastStatus, should.NotBeNil) {
									t.FailNow()
								}
								a.So(stats.LastStatus.Time, should.Equal, tc.Up.GatewayStatus.Time)
							}
						})
					}
				})

				t.Run("Downstream", func(t *testing.T) {
					ctx := clusterauth.NewContext(test.Context(), nil)
					downlinkCount := 0
					for _, tc := range []struct {
						Name                     string
						Message                  *ttnpb.DownlinkMessage
						ErrorAssertion           func(error) bool
						RxWindowDetailsAssertion []func(error) bool
					}{
						{
							Name: "InvalidSettingsType",
							Message: &ttnpb.DownlinkMessage{
								RawPayload: randomDownDataPayload(types.DevAddr{0x26, 0x01, 0xff, 0xff}, 1, 6),
								Settings: &ttnpb.DownlinkMessage_Scheduled{
									Scheduled: &ttnpb.TxSettings{
										DataRate: ttnpb.DataRate{
											Modulation: &ttnpb.DataRate_LoRa{
												LoRa: &ttnpb.LoRaDataRate{
													SpreadingFactor: 12,
													Bandwidth:       125000,
												},
											},
										},
										CodingRate: "4/5",
										Frequency:  869525000,
										Downlink: &ttnpb.TxSettings_Downlink{
											TxPower: 10,
										},
										Timestamp: 100,
									},
								},
							},
							ErrorAssertion: errors.IsInvalidArgument, // Network Server may send Tx request only.
						},
						{
							Name: "NotConnected",
							Message: &ttnpb.DownlinkMessage{
								Settings: &ttnpb.DownlinkMessage_Request{
									Request: &ttnpb.TxRequest{
										Class: ttnpb.CLASS_C,
										DownlinkPaths: []*ttnpb.DownlinkPath{
											{
												Path: &ttnpb.DownlinkPath_Fixed{
													Fixed: &ttnpb.GatewayAntennaIdentifiers{
														GatewayIdentifiers: ttnpb.GatewayIdentifiers{
															GatewayID: "not-connected",
														},
													},
												},
											},
										},
										FrequencyPlanID: test.EUFrequencyPlanID,
									},
								},
							},
							ErrorAssertion: errors.IsAborted, // The gateway is not connected.
						},
						{
							Name: "ValidClassA",
							Message: &ttnpb.DownlinkMessage{
								RawPayload: randomDownDataPayload(types.DevAddr{0x26, 0x01, 0xff, 0xff}, 1, 6),
								Settings: &ttnpb.DownlinkMessage_Request{
									Request: &ttnpb.TxRequest{
										Class: ttnpb.CLASS_A,
										DownlinkPaths: []*ttnpb.DownlinkPath{
											{
												Path: &ttnpb.DownlinkPath_UplinkToken{
													UplinkToken: io.MustUplinkToken(
														ttnpb.GatewayAntennaIdentifiers{
															GatewayIdentifiers: ttnpb.GatewayIdentifiers{
																GatewayID: registeredGatewayID,
															},
														},
														10000000,
														time.Unix(0, 10000000*1000),
													),
												},
											},
										},
										Priority:         ttnpb.TxSchedulePriority_NORMAL,
										Rx1Delay:         ttnpb.RX_DELAY_1,
										Rx1DataRateIndex: 5,
										Rx1Frequency:     868100000,
										FrequencyPlanID:  test.EUFrequencyPlanID,
									},
								},
							},
						},
						{
							Name: "ValidClassAWithoutFrequencyPlanInTxRequest",
							Message: &ttnpb.DownlinkMessage{
								RawPayload: randomDownDataPayload(types.DevAddr{0x26, 0x01, 0xff, 0xff}, 1, 6),
								Settings: &ttnpb.DownlinkMessage_Request{
									Request: &ttnpb.TxRequest{
										Class: ttnpb.CLASS_A,
										DownlinkPaths: []*ttnpb.DownlinkPath{
											{
												Path: &ttnpb.DownlinkPath_UplinkToken{
													UplinkToken: io.MustUplinkToken(
														ttnpb.GatewayAntennaIdentifiers{
															GatewayIdentifiers: ttnpb.GatewayIdentifiers{
																GatewayID: registeredGatewayID,
															},
														},
														20000000,
														time.Unix(0, 20000000*1000),
													),
												},
											},
										},
										Priority:         ttnpb.TxSchedulePriority_NORMAL,
										Rx1Delay:         ttnpb.RX_DELAY_1,
										Rx1DataRateIndex: 5,
										Rx1Frequency:     868100000,
									},
								},
							},
						},
						{
							Name: "ConflictClassA",
							Message: &ttnpb.DownlinkMessage{
								RawPayload: randomDownDataPayload(types.DevAddr{0x26, 0x02, 0xff, 0xff}, 1, 6),
								Settings: &ttnpb.DownlinkMessage_Request{
									Request: &ttnpb.TxRequest{
										Class: ttnpb.CLASS_A,
										DownlinkPaths: []*ttnpb.DownlinkPath{
											{
												Path: &ttnpb.DownlinkPath_UplinkToken{
													UplinkToken: io.MustUplinkToken(
														ttnpb.GatewayAntennaIdentifiers{
															GatewayIdentifiers: ttnpb.GatewayIdentifiers{
																GatewayID: registeredGatewayID,
															},
														},
														10000000,
														time.Unix(0, 10000000*1000),
													),
												},
											},
										},
										Priority:         ttnpb.TxSchedulePriority_NORMAL,
										Rx1Delay:         ttnpb.RX_DELAY_1,
										Rx1DataRateIndex: 5,
										Rx1Frequency:     868100000,
										FrequencyPlanID:  test.EUFrequencyPlanID,
									},
								},
							},
							ErrorAssertion: errors.IsAborted,
							RxWindowDetailsAssertion: []func(error) bool{
								errors.IsResourceExhausted,  // Rx1 conflicts with previous.
								errors.IsFailedPrecondition, // Rx2 not provided.
							},
						},
						{
							Name: "ValidClassC",
							Message: &ttnpb.DownlinkMessage{
								RawPayload: randomDownDataPayload(types.DevAddr{0x26, 0x02, 0xff, 0xff}, 1, 6),
								Settings: &ttnpb.DownlinkMessage_Request{
									Request: &ttnpb.TxRequest{
										Class: ttnpb.CLASS_C,
										DownlinkPaths: []*ttnpb.DownlinkPath{
											{
												Path: &ttnpb.DownlinkPath_Fixed{
													Fixed: &ttnpb.GatewayAntennaIdentifiers{
														GatewayIdentifiers: ttnpb.GatewayIdentifiers{
															GatewayID: registeredGatewayID,
														},
													},
												},
											},
										},
										Priority:         ttnpb.TxSchedulePriority_NORMAL,
										Rx1Delay:         ttnpb.RX_DELAY_1,
										Rx1DataRateIndex: 5,
										Rx1Frequency:     868100000,
										FrequencyPlanID:  test.EUFrequencyPlanID,
									},
								},
							},
						},
						{
							Name: "ValidClassCWithoutFrequencyPlanInTxRequest",
							Message: &ttnpb.DownlinkMessage{
								RawPayload: randomDownDataPayload(types.DevAddr{0x26, 0x02, 0xff, 0xff}, 1, 6),
								Settings: &ttnpb.DownlinkMessage_Request{
									Request: &ttnpb.TxRequest{
										Class: ttnpb.CLASS_C,
										DownlinkPaths: []*ttnpb.DownlinkPath{
											{
												Path: &ttnpb.DownlinkPath_Fixed{
													Fixed: &ttnpb.GatewayAntennaIdentifiers{
														GatewayIdentifiers: ttnpb.GatewayIdentifiers{
															GatewayID: registeredGatewayID,
														},
													},
												},
											},
										},
										Priority:         ttnpb.TxSchedulePriority_NORMAL,
										Rx1Delay:         ttnpb.RX_DELAY_1,
										Rx1DataRateIndex: 5,
										Rx1Frequency:     868100000,
									},
								},
							},
						},
					} {
						t.Run(tc.Name, func(t *testing.T) {
							a := assertions.New(t)

							_, err := gs.ScheduleDownlink(ctx, tc.Message)
							if err != nil {
								if tc.ErrorAssertion == nil || !a.So(tc.ErrorAssertion(err), should.BeTrue) {
									t.Fatalf("Unexpected error: %v", err)
								}
								if tc.RxWindowDetailsAssertion != nil {
									a.So(err, should.HaveSameErrorDefinitionAs, gatewayserver.ErrSchedule)
									if !a.So(errors.Details(err), should.HaveLength, 1) {
										t.FailNow()
									}
									details := errors.Details(err)[0].(*ttnpb.ScheduleDownlinkErrorDetails)
									if !a.So(details, should.NotBeNil) || !a.So(details.PathErrors, should.HaveLength, 1) {
										t.FailNow()
									}
									errSchedulePathCause := errors.Cause(ttnpb.ErrorDetailsFromProto(details.PathErrors[0]))
									a.So(errors.IsAborted(errSchedulePathCause), should.BeTrue)
									for i, assert := range tc.RxWindowDetailsAssertion {
										if !a.So(errors.Details(errSchedulePathCause), should.HaveLength, 1) {
											t.FailNow()
										}
										errSchedulePathCauseDetails := errors.Details(errSchedulePathCause)[0].(*ttnpb.ScheduleDownlinkErrorDetails)
										if !a.So(errSchedulePathCauseDetails, should.NotBeNil) {
											t.FailNow()
										}
										if i >= len(errSchedulePathCauseDetails.PathErrors) {
											t.Fatalf("Expected error in Rx window %d", i+1)
										}
										errRxWindow := ttnpb.ErrorDetailsFromProto(errSchedulePathCauseDetails.PathErrors[i])
										if !a.So(assert(errRxWindow), should.BeTrue) {
											t.Fatalf("Unexpected Rx window %d error: %v", i+1, errRxWindow)
										}
									}
								}
								return
							} else if tc.ErrorAssertion != nil {
								t.Fatalf("Expected error")
							}
							downlinkCount++

							select {
							case msg := <-downCh:
								settings := msg.DownlinkMessage.GetScheduled()
								a.So(settings, should.NotBeNil)
							case <-time.After(timeout):
								t.Fatal("Expected downlink timeout")
							}

							time.Sleep(2 * timeout)

							conn, ok := gs.GetConnection(ctx, ids)
							a.So(ok, should.BeTrue)
							a.So(conn.Stats(), should.NotBeNil)
							if config.Stats != nil {
								a.So(gs.UpdateConnectionStats(conn), should.BeNil)
							}

							stats, err := statsClient.GetGatewayConnectionStats(statsCtx, &ids)
							if !a.So(err, should.BeNil) {
								t.FailNow()
							}
							a.So(stats.DownlinkCount, should.Equal, downlinkCount)
						})
					}
				})

				cancel()
				wg.Wait()

				// Wait for disconnection to be processed.
				time.Sleep(timeout)

				// After canceling the context and awaiting the link, the connection should be gone.
				t.Run("Disconnected", func(t *testing.T) {
					_, err := statsClient.GetGatewayConnectionStats(statsCtx, &ids)
					if !a.So(errors.IsNotFound(err), should.BeTrue) {
						t.Fatalf("Expected gateway to be disconnected, but it's not")
					}
				})
			})
		}
	}
}
