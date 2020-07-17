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

package networkserver_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/mohae/deepcopy"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestDeviceRegistryGet(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		ContextFunc    func(context.Context) context.Context
		GetByIDFunc    func(context.Context, ttnpb.ApplicationIdentifiers, string, []string) (*ttnpb.EndDevice, context.Context, error)
		KeyVault       map[string][]byte
		Request        *ttnpb.GetEndDeviceRequest
		Device         *ttnpb.EndDevice
		ErrorAssertion func(*testing.T, error) bool
		GetByIDCalls   uint64
	}{
		{
			Name: "No device read rights",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_GATEWAY_SETTINGS_BASIC,
							},
						},
					},
				})
			},
			GetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				err := errors.New("GetByIDFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, ctx, err
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
					},
				},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				if !assertions.New(t).So(errors.IsPermissionDenied(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
		},

		{
			Name: "no keys",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_READ,
							},
						},
					},
				})
			},
			GetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
				})
				return &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					FrequencyPlanID: test.EUFrequencyPlanID,
				}, ctx, nil
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				FrequencyPlanID: test.EUFrequencyPlanID,
			},
			GetByIDCalls: 1,
		},

		{
			Name: "with keys",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_READ,
								ttnpb.RIGHT_APPLICATION_DEVICES_READ_KEYS,
								ttnpb.RIGHT_APPLICATION_LINK,
							},
						},
					},
				})
			},
			GetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"session",
					"queued_application_downlinks",
				})
				return &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					FrequencyPlanID: test.EUFrequencyPlanID,
					Session: &ttnpb.Session{
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								KEKLabel:     "test",
								EncryptedKey: []byte{0x96, 0x77, 0x8b, 0x25, 0xae, 0x6c, 0xa4, 0x35, 0xf9, 0x2b, 0x5b, 0x97, 0xc0, 0x50, 0xae, 0xd2, 0x46, 0x8a, 0xb8, 0xa1, 0x7a, 0xd8, 0x4e, 0x5d},
							},
						},
					},
				}, ctx, nil
			},
			KeyVault: map[string][]byte{
				"test": {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"session",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				FrequencyPlanID: test.EUFrequencyPlanID,
				Session: &ttnpb.Session{
					SessionKeys: ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: AES128KeyPtr(types.AES128Key{0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
						},
					},
				},
			},
			GetByIDCalls: 1,
		},

		{
			Name: "with specific key envelope",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_READ,
								ttnpb.RIGHT_APPLICATION_DEVICES_READ_KEYS,
							},
						},
					},
				})
			},
			GetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"pending_session.keys.f_nwk_s_int_key",
				})
				return &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					PendingSession: &ttnpb.Session{
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								KEKLabel:     "test",
								EncryptedKey: []byte{0x96, 0x77, 0x8b, 0x25, 0xae, 0x6c, 0xa4, 0x35, 0xf9, 0x2b, 0x5b, 0x97, 0xc0, 0x50, 0xae, 0xd2, 0x46, 0x8a, 0xb8, 0xa1, 0x7a, 0xd8, 0x4e, 0x5d},
							},
						},
					},
				}, ctx, nil
			},
			KeyVault: map[string][]byte{
				"test": {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"pending_session.keys.f_nwk_s_int_key",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				PendingSession: &ttnpb.Session{
					SessionKeys: ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: AES128KeyPtr(types.AES128Key{0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
						},
					},
				},
			},
			GetByIDCalls: 1,
		},

		{
			Name: "with specific key",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_READ,
								ttnpb.RIGHT_APPLICATION_DEVICES_READ_KEYS,
							},
						},
					},
				})
			},
			GetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"pending_session.keys.f_nwk_s_int_key.encrypted_key",
					"pending_session.keys.f_nwk_s_int_key.kek_label",
					"pending_session.keys.f_nwk_s_int_key.key",
				})
				return &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					PendingSession: &ttnpb.Session{
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								KEKLabel:     "test",
								EncryptedKey: []byte{0x96, 0x77, 0x8b, 0x25, 0xae, 0x6c, 0xa4, 0x35, 0xf9, 0x2b, 0x5b, 0x97, 0xc0, 0x50, 0xae, 0xd2, 0x46, 0x8a, 0xb8, 0xa1, 0x7a, 0xd8, 0x4e, 0x5d},
							},
						},
					},
				}, ctx, nil
			},
			KeyVault: map[string][]byte{
				"test": {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			},
			Request: &ttnpb.GetEndDeviceRequest{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"pending_session.keys.f_nwk_s_int_key.key",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				PendingSession: &ttnpb.Session{
					SessionKeys: ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: AES128KeyPtr(types.AES128Key{0x0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
						},
					},
				},
			},
			GetByIDCalls: 1,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			var getByIDCalls uint64

			ns, ctx, env, stop := StartTest(
				t,
				component.Config{
					ServiceBase: config.ServiceBase{
						KeyVault: config.KeyVault{
							Provider: "static",
							Static:   tc.KeyVault,
						},
					},
				},
				Config{
					Devices: &MockDeviceRegistry{
						GetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string) (*ttnpb.EndDevice, context.Context, error) {
							atomic.AddUint64(&getByIDCalls, 1)
							return tc.GetByIDFunc(ctx, appID, devID, gets)
						},
					},
				},
				(1<<9)*test.Delay,
			)
			defer stop()

			<-env.DownlinkTasks.Pop

			ns.AddContextFiller(tc.ContextFunc)
			ns.AddContextFiller(func(ctx context.Context) context.Context {
				ctx, cancel := context.WithDeadline(ctx, time.Now().Add(Timeout))
				_ = cancel
				return ctx
			})
			ns.AddContextFiller(func(ctx context.Context) context.Context {
				return test.ContextWithTB(ctx, t)
			})

			req := deepcopy.Copy(tc.Request).(*ttnpb.GetEndDeviceRequest)
			dev, err := ttnpb.NewNsEndDeviceRegistryClient(ns.LoopbackConn()).Get(ctx, req)
			if tc.ErrorAssertion != nil && a.So(tc.ErrorAssertion(t, err), should.BeTrue) {
				a.So(dev, should.BeNil)
			} else if a.So(err, should.BeNil) {
				a.So(dev, should.Resemble, tc.Device)
			}
			a.So(req, should.Resemble, tc.Request)
			a.So(getByIDCalls, should.Equal, tc.GetByIDCalls)
		})
	}
}

func TestDeviceRegistrySet(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		ContextFunc    func(context.Context) context.Context
		AddFunc        func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time, bool) error
		SetByIDFunc    func(context.Context, ttnpb.ApplicationIdentifiers, string, []string, func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error)
		Request        *ttnpb.SetEndDeviceRequest
		Device         *ttnpb.EndDevice
		ErrorAssertion func(*testing.T, error) bool
		AddCalls       uint64
		SetByIDCalls   uint64
	}{
		{
			Name: "No device write rights",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_GATEWAY_SETTINGS_BASIC,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				err := errors.New("SetByIDFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, ctx, err
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0,
					LoRaWANVersion:    ttnpb.MAC_V1_0,
					SupportsJoin:      true,
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"lorawan_phy_version",
						"lorawan_version",
						"supports_join",
					},
				},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				if !assertions.New(t).So(errors.IsPermissionDenied(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
		},

		{
			Name: "Create invalid device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings.adr_margin",
					"supports_class_b",
					"supports_class_c",
					"supports_join",
				})

				dev, sets, err := f(ctx, nil)
				if !a.So(err, should.NotBeNil) {
					return nil, ctx, errors.New("test failed")
				}
				a.So(dev, should.BeNil)
				a.So(sets, should.BeNil)
				a.So(errors.IsInvalidArgument(err), should.BeTrue)
				return nil, ctx, err
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					LoRaWANPHYVersion: ttnpb.PHY_V1_0,
					LoRaWANVersion:    ttnpb.MAC_V1_0,
					SupportsJoin:      true,
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"lorawan_phy_version",
						"lorawan_version",
						"mac_settings.adr_margin",
						"supports_class_b",
						"supports_class_c",
						"supports_join",
					},
				},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				return assertions.New(t).So(errors.IsInvalidArgument(err), should.BeTrue)
			},
		},

		{
			Name: "Create OTAA device with invalid frequency plan",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"last_dev_status_received_at",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings",
					"mac_state",
					"multicast",
					"queued_application_downlinks",
					"recent_uplinks",
					"session.dev_addr",
					"session.last_conf_f_cnt_down",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.queued_application_downlinks",
					"supports_class_b",
					"supports_class_c",
					"supports_join",
				})

				dev, sets, err := f(ctx, nil)
				a.So(dev, should.BeNil)
				a.So(sets, should.BeNil)
				if !a.So(err, should.NotBeNil) {
					return nil, ctx, errors.New("test")
				}
				return nil, ctx, err
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   "invalid-frequency-plan",
					LoRaWANPHYVersion: ttnpb.PHY_V1_0,
					LoRaWANVersion:    ttnpb.MAC_V1_0,
					SupportsJoin:      true,
					MACSettings: &ttnpb.MACSettings{
						ADRMargin: &pbtypes.FloatValue{Value: 4},
					},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"lorawan_phy_version",
						"lorawan_version",
						"mac_settings.adr_margin",
						"supports_class_b",
						"supports_class_c",
						"supports_join",
					},
				},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				return assertions.New(t).So(errors.IsNotFound(err), should.BeTrue)
			},
			SetByIDCalls: 1,
		},

		{
			Name: "Create OTAA device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"last_dev_status_received_at",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings",
					"mac_state",
					"multicast",
					"queued_application_downlinks",
					"recent_uplinks",
					"session.dev_addr",
					"session.last_conf_f_cnt_down",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.queued_application_downlinks",
					"supports_class_b",
					"supports_class_c",
					"supports_join",
				})

				dev, sets, err := f(ctx, nil)
				if !a.So(err, should.BeNil) {
					return nil, ctx, err
				}
				a.So(sets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"ids.application_ids",
					"ids.dev_eui",
					"ids.device_id",
					"ids.join_eui",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings.adr_margin",
					"supports_class_b",
					"supports_class_c",
					"supports_join",
				})
				a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0,
					LoRaWANVersion:    ttnpb.MAC_V1_0,
					SupportsJoin:      true,
					MACSettings: &ttnpb.MACSettings{
						ADRMargin: &pbtypes.FloatValue{Value: 4},
					},
				})
				return dev, ctx, nil
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0,
					LoRaWANVersion:    ttnpb.MAC_V1_0,
					SupportsJoin:      true,
					MACSettings: &ttnpb.MACSettings{
						ADRMargin: &pbtypes.FloatValue{Value: 4},
					},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"lorawan_phy_version",
						"lorawan_version",
						"mac_settings.adr_margin",
						"supports_class_b",
						"supports_class_c",
						"supports_join",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					JoinEUI:                &types.EUI64{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					DevEUI:                 &types.EUI64{0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				},
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0,
				LoRaWANVersion:    ttnpb.MAC_V1_0,
				SupportsJoin:      true,
				MACSettings: &ttnpb.MACSettings{
					ADRMargin: &pbtypes.FloatValue{Value: 4},
				},
			},
			SetByIDCalls: 1,
		},

		{
			// https://github.com/TheThingsNetwork/lorawan-stack/issues/104#issuecomment-465074076
			Name: "Create OTAA device with existing session",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"last_dev_status_received_at",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings",
					"mac_state",
					"multicast",
					"queued_application_downlinks",
					"recent_uplinks",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.key",
					"session.last_conf_f_cnt_down",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.queued_application_downlinks",
					"session.started_at",
					"supports_join",
				})

				dev, sets, err := f(ctx, nil)
				if !a.So(err, should.BeNil) {
					return nil, ctx, err
				}
				a.So(sets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"ids.application_ids",
					"ids.dev_addr",
					"ids.dev_eui",
					"ids.device_id",
					"ids.join_eui",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings.supports_32_bit_f_cnt",
					"mac_settings.use_adr",
					"mac_state",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.encrypted_key",
					"session.keys.f_nwk_s_int_key.kek_label",
					"session.keys.f_nwk_s_int_key.key",
					"session.keys.nwk_s_enc_key.encrypted_key",
					"session.keys.nwk_s_enc_key.kek_label",
					"session.keys.s_nwk_s_int_key.encrypted_key",
					"session.keys.s_nwk_s_int_key.kek_label",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.started_at",
					"supports_join",
				})

				expected := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x70, 0xB3, 0xD5, 0x95, 0x20, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0xA8, 0x17, 0x58, 0xFF, 0xFE, 0x03, 0x22, 0x77},
						DevAddr:                &types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
					LoRaWANVersion:    ttnpb.MAC_V1_0_2,
					SupportsJoin:      true,
					MACSettings: &ttnpb.MACSettings{
						Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
						UseADR:            &pbtypes.BoolValue{Value: true},
					},
					Session: &ttnpb.Session{
						StartedAt:     time.Unix(0, 42).UTC(),
						DevAddr:       types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
						LastFCntUp:    45872,
						LastNFCntDown: 1880,
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
							NwkSEncKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
						},
					},
				}
				macState, err := NewMACState(expected, frequencyplans.NewStore(test.FrequencyPlansFetcher), ttnpb.MACSettings{})
				if !a.So(err, should.BeNil) {
					panic(fmt.Sprintf("failed to reset MAC state: %s", err))
				}
				expected.MACState = macState
				a.So(dev, should.Resemble, expected)
				return dev, ctx, nil
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x70, 0xB3, 0xD5, 0x95, 0x20, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0xA8, 0x17, 0x58, 0xFF, 0xFE, 0x03, 0x22, 0x77},
						DevAddr:                &types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
					LoRaWANVersion:    ttnpb.MAC_V1_0_2,
					SupportsJoin:      true,
					MACSettings: &ttnpb.MACSettings{
						Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
						UseADR:            &pbtypes.BoolValue{Value: true},
					},
					Session: &ttnpb.Session{
						StartedAt:     time.Unix(0, 42).UTC(),
						DevAddr:       types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
						LastFCntUp:    45872,
						LastNFCntDown: 1880,
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &types.AES128Key{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
						},
					},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"lorawan_phy_version",
						"lorawan_version",
						"mac_settings.supports_32_bit_f_cnt",
						"mac_settings.use_adr",
						"session.dev_addr",
						"session.keys.f_nwk_s_int_key.key",
						"session.last_f_cnt_up",
						"session.last_n_f_cnt_down",
						"session.started_at",
						"supports_join",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					JoinEUI:                &types.EUI64{0x70, 0xB3, 0xD5, 0x95, 0x20, 0x00, 0x00, 0x00},
					DevEUI:                 &types.EUI64{0xA8, 0x17, 0x58, 0xFF, 0xFE, 0x03, 0x22, 0x77},
					DevAddr:                &types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
				},
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				LoRaWANVersion:    ttnpb.MAC_V1_0_2,
				SupportsJoin:      true,
				MACSettings: &ttnpb.MACSettings{
					Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
					UseADR:            &pbtypes.BoolValue{Value: true},
				},
				Session: &ttnpb.Session{
					StartedAt:     time.Unix(0, 42).UTC(),
					DevAddr:       types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
					LastFCntUp:    45872,
					LastNFCntDown: 1880,
					SessionKeys: ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &types.AES128Key{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
						},
					},
				},
			},
			SetByIDCalls: 1,
		},

		{
			Name: "Create ABP device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"last_dev_status_received_at",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings",
					"mac_state",
					"multicast",
					"queued_application_downlinks",
					"recent_uplinks",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.key",
					"session.last_conf_f_cnt_down",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.queued_application_downlinks",
					"session.started_at",
					"supports_join",
				})

				dev, sets, err := f(ctx, nil)
				if !a.So(err, should.BeNil) {
					return nil, ctx, err
				}
				a.So(sets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"ids.application_ids",
					"ids.dev_addr",
					"ids.device_id",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings.supports_32_bit_f_cnt",
					"mac_settings.use_adr",
					"mac_state",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.encrypted_key",
					"session.keys.f_nwk_s_int_key.kek_label",
					"session.keys.f_nwk_s_int_key.key",
					"session.keys.nwk_s_enc_key.encrypted_key",
					"session.keys.nwk_s_enc_key.kek_label",
					"session.keys.s_nwk_s_int_key.encrypted_key",
					"session.keys.s_nwk_s_int_key.kek_label",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.started_at",
					"supports_join",
				})

				expected := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						DevAddr:                &types.DevAddr{0x42, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
					LoRaWANVersion:    ttnpb.MAC_V1_0_2,
					MACSettings: &ttnpb.MACSettings{
						Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
						UseADR:            &pbtypes.BoolValue{Value: true},
					},
					Session: &ttnpb.Session{
						StartedAt:     time.Unix(0, 42).UTC(),
						DevAddr:       types.DevAddr{0x42, 0x00, 0x00, 0x00},
						LastFCntUp:    42,
						LastNFCntDown: 4242,
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
							},
							NwkSEncKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
							},
						},
					},
				}
				macState, err := NewMACState(expected, frequencyplans.NewStore(test.FrequencyPlansFetcher), ttnpb.MACSettings{})
				if !a.So(err, should.BeNil) {
					panic(fmt.Sprintf("failed to reset MAC state: %s", err))
				}
				expected.MACState = macState
				a.So(dev, should.Resemble, expected)
				return dev, ctx, nil
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						DevAddr:                &types.DevAddr{0x42, 0x00, 0x00, 0x00},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
					LoRaWANVersion:    ttnpb.MAC_V1_0_2,
					MACSettings: &ttnpb.MACSettings{
						Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
						UseADR:            &pbtypes.BoolValue{Value: true},
					},
					Session: &ttnpb.Session{
						StartedAt:     time.Unix(0, 42).UTC(),
						DevAddr:       types.DevAddr{0x42, 0x00, 0x00, 0x00},
						LastFCntUp:    42,
						LastNFCntDown: 4242,
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &types.AES128Key{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
							},
						},
					},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"lorawan_phy_version",
						"lorawan_version",
						"mac_settings.supports_32_bit_f_cnt",
						"mac_settings.use_adr",
						"session.dev_addr",
						"session.keys.f_nwk_s_int_key.key",
						"session.last_f_cnt_up",
						"session.last_n_f_cnt_down",
						"session.started_at",
						"supports_join",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					DevAddr:                &types.DevAddr{0x42, 0x00, 0x00, 0x00},
				},
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				LoRaWANVersion:    ttnpb.MAC_V1_0_2,
				MACSettings: &ttnpb.MACSettings{
					Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
					UseADR:            &pbtypes.BoolValue{Value: true},
				},
				Session: &ttnpb.Session{
					StartedAt:     time.Unix(0, 42).UTC(),
					DevAddr:       types.DevAddr{0x42, 0x00, 0x00, 0x00},
					LastFCntUp:    42,
					LastNFCntDown: 4242,
					SessionKeys: ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &types.AES128Key{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
						},
					},
				},
			},
			SetByIDCalls: 1,
		},

		{
			// https://github.com/TheThingsNetwork/lorawan-stack/issues/159#issue-411803325
			Name: "Create ABP device with existing session",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE_KEYS,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"last_dev_status_received_at",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings",
					"mac_state",
					"multicast",
					"queued_application_downlinks",
					"recent_uplinks",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.key",
					"session.last_conf_f_cnt_down",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.queued_application_downlinks",
					"session.started_at",
					"supports_join",
				})

				dev, sets, err := f(ctx, nil)
				if !a.So(err, should.BeNil) {
					return nil, ctx, err
				}
				a.So(sets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"ids.application_ids",
					"ids.dev_addr",
					"ids.dev_eui",
					"ids.device_id",
					"ids.join_eui",
					"lorawan_phy_version",
					"lorawan_version",
					"mac_settings.supports_32_bit_f_cnt",
					"mac_settings.use_adr",
					"mac_state",
					"session.dev_addr",
					"session.keys.f_nwk_s_int_key.encrypted_key",
					"session.keys.f_nwk_s_int_key.kek_label",
					"session.keys.f_nwk_s_int_key.key",
					"session.keys.nwk_s_enc_key.encrypted_key",
					"session.keys.nwk_s_enc_key.kek_label",
					"session.keys.s_nwk_s_int_key.encrypted_key",
					"session.keys.s_nwk_s_int_key.kek_label",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.started_at",
					"supports_join",
				})

				expected := &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x70, 0xB3, 0xD5, 0x95, 0x20, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0xA8, 0x17, 0x58, 0xFF, 0xFE, 0x03, 0x22, 0x77},
						DevAddr:                &types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
					LoRaWANVersion:    ttnpb.MAC_V1_0_2,
					MACSettings: &ttnpb.MACSettings{
						Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
						UseADR:            &pbtypes.BoolValue{Value: true},
					},
					Session: &ttnpb.Session{
						StartedAt:     time.Unix(0, 42).UTC(),
						DevAddr:       types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
						LastFCntUp:    45872,
						LastNFCntDown: 1880,
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
							NwkSEncKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
							SNwkSIntKey: &ttnpb.KeyEnvelope{
								EncryptedKey: []byte{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
						},
					},
				}
				macState, err := NewMACState(expected, frequencyplans.NewStore(test.FrequencyPlansFetcher), ttnpb.MACSettings{})
				if !a.So(err, should.BeNil) {
					panic(fmt.Sprintf("failed to reset MAC state: %s", err))
				}
				expected.MACState = macState
				a.So(dev, should.Resemble, expected)
				return dev, ctx, nil
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
						JoinEUI:                &types.EUI64{0x70, 0xB3, 0xD5, 0x95, 0x20, 0x00, 0x00, 0x00},
						DevEUI:                 &types.EUI64{0xA8, 0x17, 0x58, 0xFF, 0xFE, 0x03, 0x22, 0x77},
						DevAddr:                &types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
					LoRaWANVersion:    ttnpb.MAC_V1_0_2,
					MACSettings: &ttnpb.MACSettings{
						Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
						UseADR:            &pbtypes.BoolValue{Value: true},
					},
					Session: &ttnpb.Session{
						StartedAt:     time.Unix(0, 42).UTC(),
						DevAddr:       types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
						LastFCntUp:    45872,
						LastNFCntDown: 1880,
						SessionKeys: ttnpb.SessionKeys{
							FNwkSIntKey: &ttnpb.KeyEnvelope{
								Key: &types.AES128Key{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
							},
						},
					},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
						"lorawan_phy_version",
						"lorawan_version",
						"mac_settings.supports_32_bit_f_cnt",
						"mac_settings.use_adr",
						"session.dev_addr",
						"session.keys.f_nwk_s_int_key.key",
						"session.last_f_cnt_up",
						"session.last_n_f_cnt_down",
						"session.started_at",
						"supports_join",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					JoinEUI:                &types.EUI64{0x70, 0xB3, 0xD5, 0x95, 0x20, 0x00, 0x00, 0x00},
					DevEUI:                 &types.EUI64{0xA8, 0x17, 0x58, 0xFF, 0xFE, 0x03, 0x22, 0x77},
					DevAddr:                &types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
				},
				FrequencyPlanID:   test.EUFrequencyPlanID,
				LoRaWANPHYVersion: ttnpb.PHY_V1_0_2_REV_B,
				LoRaWANVersion:    ttnpb.MAC_V1_0_2,
				MACSettings: &ttnpb.MACSettings{
					Supports32BitFCnt: &pbtypes.BoolValue{Value: true},
					UseADR:            &pbtypes.BoolValue{Value: true},
				},
				Session: &ttnpb.Session{
					StartedAt:     time.Unix(0, 42).UTC(),
					DevAddr:       types.DevAddr{0x01, 0x0b, 0x60, 0x0c},
					LastFCntUp:    45872,
					LastNFCntDown: 1880,
					SessionKeys: ttnpb.SessionKeys{
						FNwkSIntKey: &ttnpb.KeyEnvelope{
							Key: &types.AES128Key{0x9e, 0x2f, 0xb6, 0x1d, 0x73, 0x10, 0xc9, 0x27, 0x98, 0x86, 0xdb, 0x79, 0xfa, 0x52, 0xf9, 0xf4},
						},
					},
				},
			},
			SetByIDCalls: 1,
		},

		{
			Name: "Update with invalid frequency plan",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"last_dev_status_received_at",
					"lorawan_phy_version",
					"mac_settings",
					"mac_state",
					"multicast",
					"queued_application_downlinks",
					"recent_uplinks",
					"session.dev_addr",
					"session.last_conf_f_cnt_down",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.queued_application_downlinks",
				})

				dev, sets, err := f(ctx, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
				})
				a.So(dev, should.BeNil)
				a.So(sets, should.BeNil)
				if !a.So(err, should.NotBeNil) {
					return nil, ctx, errors.New("test")
				}
				return nil, ctx, err
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					FrequencyPlanID: "invalid-frequency-plan",
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"frequency_plan_id",
					},
				},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				return assertions.New(t).So(errors.IsNotFound(err), should.BeTrue)
			},
			SetByIDCalls: 1,
		},

		{
			Name: "Update device desired MAC parameters",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					},
				})
			},
			AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
				err := errors.New("AddFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return err
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.HaveSameElementsDeep, []string{
					"frequency_plan_id",
					"last_dev_status_received_at",
					"lorawan_phy_version",
					"mac_settings",
					"mac_state",
					"multicast",
					"queued_application_downlinks",
					"recent_uplinks",
					"session.dev_addr",
					"session.last_conf_f_cnt_down",
					"session.last_f_cnt_up",
					"session.last_n_f_cnt_down",
					"session.queued_application_downlinks",
				})

				dev, sets, err := f(ctx, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						CurrentParameters: ttnpb.MACParameters{
							Rx2Frequency: 868000000,
						},
						DesiredParameters: ttnpb.MACParameters{
							Rx2Frequency: 868000000,
						},
					},
				})
				if !a.So(err, should.BeNil) {
					return nil, ctx, err
				}
				a.So(sets, should.HaveSameElementsDeep, []string{
					"mac_state.desired_parameters.rx2_frequency",
				})
				a.So(dev, should.Resemble, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					MACState: &ttnpb.MACState{
						DesiredParameters: ttnpb.MACParameters{
							Rx2Frequency: 123456789,
						},
					},
				})
				return &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					FrequencyPlanID:   test.EUFrequencyPlanID,
					LoRaWANPHYVersion: ttnpb.PHY_V1_1_REV_B,
					MACState: &ttnpb.MACState{
						LoRaWANVersion: ttnpb.MAC_V1_1,
						CurrentParameters: ttnpb.MACParameters{
							Rx2Frequency: 868000000,
						},
						DesiredParameters: ttnpb.MACParameters{
							Rx2Frequency: 123456789,
						},
					},
				}, ctx, nil
			},
			Request: &ttnpb.SetEndDeviceRequest{
				EndDevice: ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
					MACState: &ttnpb.MACState{
						DesiredParameters: ttnpb.MACParameters{
							Rx2Frequency: 123456789,
						},
					},
				},
				FieldMask: pbtypes.FieldMask{
					Paths: []string{
						"mac_state.desired_parameters.rx2_frequency",
					},
				},
			},
			Device: &ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-dev-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
				},
				MACState: &ttnpb.MACState{
					DesiredParameters: ttnpb.MACParameters{
						Rx2Frequency: 123456789,
					},
				},
			},
			SetByIDCalls: 1,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			var addCalls, setByIDCalls uint64

			ns, ctx, env, stop := StartTest(
				t,
				component.Config{},
				Config{
					Devices: &MockDeviceRegistry{
						SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
							atomic.AddUint64(&setByIDCalls, 1)
							return tc.SetByIDFunc(ctx, appID, devID, gets, f)
						},
					},
					DownlinkTasks: &MockDownlinkTaskQueue{
						AddFunc: func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, at time.Time, replace bool) error {
							atomic.AddUint64(&addCalls, 1)
							return tc.AddFunc(ctx, ids, at, replace)
						},
						PopFunc: DownlinkTaskPopBlockFunc,
					},
				},
				(1<<9)*test.Delay,
			)
			defer stop()

			go LogEvents(t, env.Events)

			ns.AddContextFiller(tc.ContextFunc)
			ns.AddContextFiller(func(ctx context.Context) context.Context {
				ctx, cancel := context.WithDeadline(ctx, time.Now().Add(Timeout))
				_ = cancel
				return ctx
			})
			ns.AddContextFiller(func(ctx context.Context) context.Context {
				return test.ContextWithTB(ctx, t)
			})

			req := deepcopy.Copy(tc.Request).(*ttnpb.SetEndDeviceRequest)
			dev, err := ttnpb.NewNsEndDeviceRegistryClient(ns.LoopbackConn()).Set(ctx, req)
			if tc.ErrorAssertion != nil && a.So(tc.ErrorAssertion(t, err), should.BeTrue) {
				a.So(dev, should.BeNil)
			} else if a.So(err, should.BeNil) {
				a.So(dev, should.Resemble, tc.Device)
			}
			a.So(req, should.Resemble, tc.Request)
			a.So(setByIDCalls, should.Equal, tc.SetByIDCalls)
			a.So(addCalls, should.Equal, tc.AddCalls)
		})
	}
}

func TestDeviceRegistryDelete(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		ContextFunc    func(context.Context) context.Context
		SetByIDFunc    func(context.Context, ttnpb.ApplicationIdentifiers, string, []string, func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error)
		Request        *ttnpb.EndDeviceIdentifiers
		ErrorAssertion func(*testing.T, error) bool
		SetByIDCalls   uint64
	}{
		{
			Name: "No device write rights",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_GATEWAY_SETTINGS_BASIC,
							},
						},
					},
				})
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				err := errors.New("SetByIDFunc must not be called")
				test.MustTFromContext(ctx).Error(err)
				return nil, ctx, err
			},
			Request: &ttnpb.EndDeviceIdentifiers{
				DeviceID:               "test-dev-id",
				ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
			},
			ErrorAssertion: func(t *testing.T, err error) bool {
				if !assertions.New(t).So(errors.IsPermissionDenied(err), should.BeTrue) {
					t.Errorf("Received error: %s", err)
					return false
				}
				return true
			},
		},

		{
			Name: "Non-existing device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					},
				})
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				t := test.MustTFromContext(ctx)
				a := assertions.New(t)
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.BeNil)

				dev, sets, err := f(ctx, nil)
				if !a.So(errors.IsNotFound(err), should.BeTrue) {
					return nil, ctx, err
				}
				a.So(sets, should.BeNil)
				a.So(dev, should.BeNil)
				return nil, ctx, nil
			},
			Request: &ttnpb.EndDeviceIdentifiers{
				DeviceID:               "test-dev-id",
				ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
			},
			SetByIDCalls: 1,
		},

		{
			Name: "Existing device",
			ContextFunc: func(ctx context.Context) context.Context {
				return rights.NewContext(ctx, rights.Rights{
					ApplicationRights: map[string]*ttnpb.Rights{
						unique.ID(test.Context(), ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"}): {
							Rights: []ttnpb.Right{
								ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
							},
						},
					},
				})
			},
			SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
				a := assertions.New(test.MustTFromContext(ctx))
				a.So(appID, should.Resemble, ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"})
				a.So(devID, should.Equal, "test-dev-id")
				a.So(gets, should.BeNil)

				dev, sets, err := f(ctx, &ttnpb.EndDevice{
					EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
						DeviceID:               "test-dev-id",
						ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
					},
				})
				if !a.So(err, should.BeNil) {
					return nil, ctx, err
				}
				a.So(sets, should.BeNil)
				a.So(dev, should.BeNil)
				return nil, ctx, nil
			},
			Request: &ttnpb.EndDeviceIdentifiers{
				DeviceID:               "test-dev-id",
				ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test-app-id"},
			},
			SetByIDCalls: 1,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			var setByIDCalls uint64

			ns, ctx, env, stop := StartTest(
				t,
				component.Config{},
				Config{
					Devices: &MockDeviceRegistry{
						SetByIDFunc: func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(context.Context, *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
							atomic.AddUint64(&setByIDCalls, 1)
							return tc.SetByIDFunc(ctx, appID, devID, gets, f)
						},
					},
					DownlinkTasks: &MockDownlinkTaskQueue{
						PopFunc: DownlinkTaskPopBlockFunc,
					},
				},
				(1<<9)*test.Delay,
			)
			defer stop()

			go LogEvents(t, env.Events)

			ns.AddContextFiller(tc.ContextFunc)
			ns.AddContextFiller(func(ctx context.Context) context.Context {
				ctx, cancel := context.WithDeadline(ctx, time.Now().Add(Timeout))
				_ = cancel
				return ctx
			})
			ns.AddContextFiller(func(ctx context.Context) context.Context {
				return test.ContextWithTB(ctx, t)
			})

			req := deepcopy.Copy(tc.Request).(*ttnpb.EndDeviceIdentifiers)
			res, err := ttnpb.NewNsEndDeviceRegistryClient(ns.LoopbackConn()).Delete(ctx, req)
			a.So(setByIDCalls, should.Equal, tc.SetByIDCalls)
			if tc.ErrorAssertion != nil && a.So(tc.ErrorAssertion(t, err), should.BeTrue) {
				a.So(res, should.BeNil)
			} else if a.So(err, should.BeNil) {
				a.So(res, should.Resemble, ttnpb.Empty)
			}
			a.So(req, should.Resemble, tc.Request)
		})
	}
}
