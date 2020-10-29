// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

package joinserver

import (
	"context"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/crypto/cryptoutil"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

type applicationActivationSettingsRegistryServer struct {
	JS       *JoinServer
	kekLabel string
}

var errApplicationActivationSettingsNotFound = errors.DefineNotFound("application_activation_settings_not_found", "application activation settings not found")

// Get implements ttnpb.ApplicationActivationSettingsRegistryServer.
func (srv applicationActivationSettingsRegistryServer) Get(ctx context.Context, req *ttnpb.GetApplicationActivationSettingsRequest) (*ttnpb.ApplicationActivationSettings, error) {
	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_READ_KEYS); err != nil {
		return nil, err
	}
	sets, err := srv.JS.applicationActivationSettings.GetByID(ctx, req.ApplicationIdentifiers, req.FieldMask.Paths)
	if errors.IsNotFound(err) {
		return nil, errApplicationActivationSettingsNotFound.WithCause(err)
	}
	if err != nil {
		return nil, err
	}
	kek, err := cryptoutil.UnwrapKeyEnvelope(ctx, sets.KEK, srv.JS.KeyVault)
	if err != nil {
		return nil, errUnwrapKey.WithCause(err)
	}
	sets.KEK = kek
	return sets, nil
}

var errNoPaths = errors.DefineInvalidArgument("no_paths", "no paths specified")

// Set implements ttnpb.ApplicationActivationSettingsRegistryServer.
func (srv applicationActivationSettingsRegistryServer) Set(ctx context.Context, req *ttnpb.SetApplicationActivationSettingsRequest) (*ttnpb.ApplicationActivationSettings, error) {
	if len(req.FieldMask.Paths) == 0 {
		return nil, errInvalidFieldMask.WithCause(errNoPaths)
	}

	if ttnpb.HasAnyField(req.FieldMask.Paths, "kek.key") && req.ApplicationActivationSettings.GetKEK().GetKey().IsZero() {
		return nil, errInvalidFieldValue.WithAttributes("field", "kek.key")
	}

	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_READ_KEYS, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE_KEYS); err != nil {
		return nil, err
	}

	sets := req.FieldMask.Paths
	reqKEK := req.ApplicationActivationSettings.KEK
	if ttnpb.HasAnyField(sets, "kek.key") {
		kek, err := cryptoutil.WrapAES128Key(ctx, *req.ApplicationActivationSettings.KEK.Key, srv.kekLabel, srv.JS.KeyVault)
		if err != nil {
			return nil, errWrapKey.WithCause(err)
		}
		req.ApplicationActivationSettings.KEK = kek
		sets = append(req.FieldMask.Paths[:0:0], req.FieldMask.Paths...)
		sets = ttnpb.AddFields(sets,
			"kek.encrypted_key",
			"kek.kek_label",
		)
	}
	v, err := srv.JS.applicationActivationSettings.SetByID(ctx, req.ApplicationIdentifiers, req.FieldMask.Paths, func(stored *ttnpb.ApplicationActivationSettings) (*ttnpb.ApplicationActivationSettings, []string, error) {
		return &req.ApplicationActivationSettings, sets, nil
	})
	if err != nil {
		return nil, err
	}
	v.KEK = reqKEK
	return v, nil
}

// Delete implements ttnpb.ApplicationActivationSettingsRegistryServer.
func (srv applicationActivationSettingsRegistryServer) Delete(ctx context.Context, req *ttnpb.DeleteApplicationActivationSettingsRequest) (*pbtypes.Empty, error) {
	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_DEVICES_READ_KEYS, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE_KEYS); err != nil {
		return nil, err
	}
	_, err := srv.JS.applicationActivationSettings.SetByID(ctx, req.ApplicationIdentifiers, nil, func(stored *ttnpb.ApplicationActivationSettings) (*ttnpb.ApplicationActivationSettings, []string, error) {
		if stored == nil {
			return nil, nil, errApplicationActivationSettingsNotFound.New()
		}
		return nil, nil, nil
	})
	if err != nil {
		return nil, err
	}
	return ttnpb.Empty, nil
}
