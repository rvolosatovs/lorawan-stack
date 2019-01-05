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

package identityserver

import (
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"google.golang.org/grpc"
)

func TestEndDevicesPermissionDenied(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	testWithIdentityServer(t, func(is *IdentityServer, cc *grpc.ClientConn) {
		reg := ttnpb.NewEndDeviceRegistryClient(cc)

		_, err := reg.Create(ctx, &ttnpb.CreateEndDeviceRequest{
			EndDevice: ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID: "test-device-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
						ApplicationID: "test-app-id",
					},
				},
			},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		_, err = reg.Get(ctx, &ttnpb.GetEndDeviceRequest{
			FieldMask: types.FieldMask{Paths: []string{"name"}},
			EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
				DeviceID: "test-device-id",
				ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
					ApplicationID: "test-app-id",
				},
			},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		_, err = reg.List(ctx, &ttnpb.ListEndDevicesRequest{
			FieldMask: types.FieldMask{Paths: []string{"name"}},
			ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
				ApplicationID: "test-app-id",
			},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		_, err = reg.Update(ctx, &ttnpb.UpdateEndDeviceRequest{
			FieldMask: types.FieldMask{Paths: []string{"name"}},
			EndDevice: ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID: "test-device-id",
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
						ApplicationID: "test-app-id",
					},
				},
			},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		_, err = reg.Delete(ctx, &ttnpb.EndDeviceIdentifiers{
			DeviceID: "test-device-id",
			ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{
				ApplicationID: "test-app-id",
			},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}
	})
}

func TestEndDevicesCRUD(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	testWithIdentityServer(t, func(is *IdentityServer, cc *grpc.ClientConn) {
		reg := ttnpb.NewEndDeviceRegistryClient(cc)

		userID := defaultUser.UserIdentifiers
		creds := userCreds(defaultUserIdx)
		app := userApplications(&userID).Applications[0]

		start := time.Now()

		created, err := reg.Create(ctx, &ttnpb.CreateEndDeviceRequest{
			EndDevice: ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-device-id",
					ApplicationIdentifiers: app.ApplicationIdentifiers,
				},
				Name: "test-device-name",
			},
		}, creds)

		a.So(created, should.NotBeNil)
		a.So(created.CreatedAt, should.HappenAfter, start)
		a.So(created.UpdatedAt, should.HappenAfter, start)
		a.So(created.Name, should.Equal, "test-device-name")
		a.So(err, should.BeNil)

		got, err := reg.Get(ctx, &ttnpb.GetEndDeviceRequest{
			FieldMask: types.FieldMask{Paths: []string{"name"}},
			EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
				DeviceID:               "test-device-id",
				ApplicationIdentifiers: app.ApplicationIdentifiers,
			},
		}, creds)

		a.So(got, should.NotBeNil)
		a.So(got.Name, should.Equal, "test-device-name")
		a.So(err, should.BeNil)

		list, err := reg.List(ctx, &ttnpb.ListEndDevicesRequest{
			FieldMask:              types.FieldMask{Paths: []string{"name"}},
			ApplicationIdentifiers: app.ApplicationIdentifiers,
		}, creds)

		a.So(err, should.BeNil)
		if a.So(list.EndDevices, should.HaveLength, 1) {
			a.So(list.EndDevices[0], should.NotBeNil)
			a.So(list.EndDevices[0].Name, should.Equal, "test-device-name")
		}

		start = time.Now()

		updated, err := reg.Update(ctx, &ttnpb.UpdateEndDeviceRequest{
			FieldMask: types.FieldMask{Paths: []string{"name"}},
			EndDevice: ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID:               "test-device-id",
					ApplicationIdentifiers: app.ApplicationIdentifiers,
				},
				Name: "test-device-name-new",
			},
		}, creds)

		a.So(updated, should.NotBeNil)
		a.So(updated.Name, should.Equal, "test-device-name-new")
		a.So(updated.UpdatedAt, should.HappenAfter, start)
		a.So(err, should.BeNil)

		_, err = reg.Delete(ctx, &ttnpb.EndDeviceIdentifiers{
			DeviceID:               "test-device-id",
			ApplicationIdentifiers: app.ApplicationIdentifiers,
		}, creds)

		a.So(err, should.BeNil)

		_, err = reg.Get(ctx, &ttnpb.GetEndDeviceRequest{
			FieldMask: types.FieldMask{Paths: []string{"name"}},
			EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
				DeviceID:               "test-device-id",
				ApplicationIdentifiers: app.ApplicationIdentifiers,
			},
		}, creds)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}
	})
}
