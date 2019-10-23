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

package redis

import (
	"bytes"
	"context"
	"encoding/base64"
	"runtime/trace"
	"time"

	"github.com/go-redis/redis"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/provisioning"
	ttnredis "go.thethings.network/lorawan-stack/pkg/redis"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/unique"
)

var (
	errAlreadyProvisioned   = errors.DefineAlreadyExists("already_provisioned", "device already provisioned")
	errDuplicateIdentifiers = errors.DefineAlreadyExists("duplicate_identifiers", "duplicate identifiers")
	errInvalidFieldmask     = errors.DefineInvalidArgument("invalid_fieldmask", "invalid fieldmask")
	errInvalidIdentifiers   = errors.DefineInvalidArgument("invalid_identifiers", "invalid identifiers")
	errReadOnlyField        = errors.DefineInvalidArgument("read_only_field", "read-only field `{field}`")
	errProvisionerNotFound  = errors.DefineNotFound("provisioner_not_found", "provisioner `{id}` not found")
)

// appendImplicitDeviceGetPaths appends implicit ttnpb.EndDevice get paths to paths.
func appendImplicitDeviceGetPaths(paths ...string) []string {
	return append(append(make([]string, 0, 3+len(paths)),
		"created_at",
		"ids",
		"updated_at",
	), paths...)
}

func applyDeviceFieldMask(dst, src *ttnpb.EndDevice, paths ...string) (*ttnpb.EndDevice, error) {
	if dst == nil {
		dst = &ttnpb.EndDevice{}
	}
	return dst, dst.SetFields(src, paths...)
}

// DeviceRegistry is an implementation of joinserver.DeviceRegistry.
type DeviceRegistry struct {
	Redis *ttnredis.Client
}

func provisionerUniqueID(dev *ttnpb.EndDevice) (string, error) {
	if dev.ProvisionerID == "" {
		return "", nil
	}
	provisioner := provisioning.Get(dev.ProvisionerID)
	if provisioner == nil {
		return "", errProvisionerNotFound.WithAttributes("id", dev.ProvisionerID)
	}
	return provisioner.UniqueID(dev.ProvisioningData)
}

func (r *DeviceRegistry) uidKey(uid string) string {
	return r.Redis.Key("uid", uid)
}

func (r *DeviceRegistry) euiKey(joinEUI, devEUI types.EUI64) string {
	return r.Redis.Key("eui", joinEUI.String(), devEUI.String())
}

func (r *DeviceRegistry) provisionerKey(provisionerID, pid string) string {
	return r.Redis.Key("provisioner", provisionerID, pid)
}

// GetByID gets device by appID, devID.
func (r *DeviceRegistry) GetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string) (*ttnpb.EndDevice, error) {
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: appID,
		DeviceID:               devID,
	}
	if err := ids.ValidateContext(ctx); err != nil {
		return nil, err
	}

	defer trace.StartRegion(ctx, "get end device by id").End()

	pb := &ttnpb.EndDevice{}
	if err := ttnredis.GetProto(r.Redis, r.uidKey(unique.ID(ctx, ids))).ScanProto(pb); err != nil {
		return nil, err
	}
	return applyDeviceFieldMask(nil, pb, appendImplicitDeviceGetPaths(paths...)...)
}

// GetByEUI gets device by joinEUI, devEUI.
func (r *DeviceRegistry) GetByEUI(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, error) {
	if joinEUI.IsZero() || devEUI.IsZero() {
		return nil, errInvalidIdentifiers
	}

	defer trace.StartRegion(ctx, "get end device by eui").End()

	pb := &ttnpb.EndDevice{}
	if err := ttnredis.FindProto(r.Redis, r.euiKey(joinEUI, devEUI), r.uidKey).ScanProto(pb); err != nil {
		return nil, err
	}
	return applyDeviceFieldMask(nil, pb, appendImplicitDeviceGetPaths(paths...)...)
}

func equalEUI64(x, y *types.EUI64) bool {
	if x == nil || y == nil {
		return x == y
	}
	return x.Equal(*y)
}

func (r *DeviceRegistry) set(tx *redis.Tx, uid string, gets []string, f func(pb *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error) {
	uk := r.uidKey(uid)

	cmd := ttnredis.GetProto(tx, uk)
	stored := &ttnpb.EndDevice{}
	if err := cmd.ScanProto(stored); errors.IsNotFound(err) {
		stored = nil
	} else if err != nil {
		return nil, err
	}

	gets = appendImplicitDeviceGetPaths(gets...)

	var pb *ttnpb.EndDevice
	var err error
	if stored != nil {
		pb = &ttnpb.EndDevice{}
		if err := cmd.ScanProto(pb); err != nil {
			return nil, err
		}
		pb, err = applyDeviceFieldMask(nil, pb, gets...)
		if err != nil {
			return nil, err
		}
	}

	var sets []string
	pb, sets, err = f(pb)
	if err != nil {
		return nil, err
	}
	if err := ttnpb.ProhibitFields(sets,
		"created_at",
		"updated_at",
	); err != nil {
		return nil, errInvalidFieldmask.WithCause(err)
	}

	if stored == nil && pb == nil {
		return nil, nil
	}
	if pb != nil && len(sets) == 0 {
		return applyDeviceFieldMask(nil, stored, gets...)
	}

	var pipelined func(redis.Pipeliner) error
	if pb == nil && len(sets) == 0 {
		pipelined = func(p redis.Pipeliner) error {
			p.Del(uk)
			if stored.JoinEUI != nil && stored.DevEUI != nil {
				p.Del(r.euiKey(*stored.JoinEUI, *stored.DevEUI))
			}
			pid, err := provisionerUniqueID(stored)
			if err != nil {
				return err
			}
			if pid != "" {
				p.Del(r.provisionerKey(stored.ProvisionerID, pid))
			}
			return nil
		}
	} else {
		if pb == nil {
			pb = &ttnpb.EndDevice{}
		}

		pb.UpdatedAt = time.Now().UTC()
		sets = append(append(sets[:0:0], sets...),
			"updated_at",
		)

		updated := &ttnpb.EndDevice{}
		var updatedPID string
		if stored == nil {
			if err := ttnpb.RequireFields(sets,
				"ids.application_ids",
				"ids.dev_eui",
				"ids.device_id",
				"ids.join_eui",
			); err != nil {
				return nil, errInvalidFieldmask.WithCause(err)
			}

			pb.CreatedAt = pb.UpdatedAt
			sets = append(sets, "created_at")

			updated, err = applyDeviceFieldMask(updated, pb, sets...)
			if err != nil {
				return nil, err
			}
			updatedPID, err = provisionerUniqueID(updated)
			if err != nil {
				return nil, err
			}
			if updated.JoinEUI == nil || updated.DevEUI == nil || updated.JoinEUI.IsZero() || updated.DevEUI.IsZero() {
				return nil, errInvalidIdentifiers
			}
		} else {
			if ttnpb.HasAnyField(sets, "ids.application_ids.application_id") && pb.ApplicationID != stored.ApplicationID {
				return nil, errReadOnlyField.WithAttributes("field", "ids.application_ids.application_id")
			}
			if ttnpb.HasAnyField(sets, "ids.device_id") && pb.DeviceID != stored.DeviceID {
				return nil, errReadOnlyField.WithAttributes("field", "ids.device_id")
			}
			if ttnpb.HasAnyField(sets, "ids.join_eui") && !equalEUI64(pb.JoinEUI, stored.JoinEUI) {
				return nil, errReadOnlyField.WithAttributes("field", "ids.join_eui")
			}
			if ttnpb.HasAnyField(sets, "ids.dev_eui") && !equalEUI64(pb.DevEUI, stored.DevEUI) {
				return nil, errReadOnlyField.WithAttributes("field", "ids.dev_eui")
			}
			if ttnpb.HasAnyField(sets, "provisioner_id") && pb.ProvisionerID != stored.ProvisionerID {
				return nil, errReadOnlyField.WithAttributes("field", "provisioner_id")
			}
			if ttnpb.HasAnyField(sets, "provisioning_data") && !pb.ProvisioningData.Equal(stored.ProvisioningData) {
				return nil, errReadOnlyField.WithAttributes("field", "provisioning_data")
			}
			if err := cmd.ScanProto(updated); err != nil {
				return nil, err
			}
			updated, err = applyDeviceFieldMask(updated, pb, sets...)
			if err != nil {
				return nil, err
			}
		}
		if err := updated.ValidateFields(sets...); err != nil {
			return nil, err
		}

		pipelined = func(p redis.Pipeliner) error {
			if stored == nil {
				ek := r.euiKey(*updated.JoinEUI, *updated.DevEUI)
				if err := tx.Watch(ek).Err(); err != nil {
					return err
				}
				i, err := tx.Exists(ek).Result()
				if err != nil {
					return ttnredis.ConvertError(err)
				}
				if i != 0 {
					return errDuplicateIdentifiers
				}
				p.SetNX(ek, uid, 0)
			}

			if updatedPID != "" {
				pk := r.provisionerKey(updated.ProvisionerID, updatedPID)
				if err := tx.Watch(pk).Err(); err != nil {
					return err
				}
				i, err := tx.Exists(pk).Result()
				if err != nil {
					return ttnredis.ConvertError(err)
				}
				if i != 0 {
					return errAlreadyProvisioned
				}
				p.SetNX(pk, uid, 0)
			}

			_, err := ttnredis.SetProto(p, uk, updated, 0)
			if err != nil {
				return err
			}
			return nil
		}
		pb, err = applyDeviceFieldMask(nil, updated, gets...)
		if err != nil {
			return nil, err
		}
	}
	_, err = tx.Pipelined(pipelined)
	if err != nil {
		return nil, err
	}
	return pb, nil
}

// SetByEUI sets device by joinEUI, devEUI.
// SetByEUI will only succeed if the device is set via SetByID first.
func (r *DeviceRegistry) SetByEUI(ctx context.Context, joinEUI types.EUI64, devEUI types.EUI64, gets []string, f func(pb *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error) {
	if joinEUI.IsZero() || devEUI.IsZero() {
		return nil, errInvalidIdentifiers
	}
	ek := r.euiKey(joinEUI, devEUI)

	defer trace.StartRegion(ctx, "set end device by eui").End()

	var pb *ttnpb.EndDevice
	err := r.Redis.Watch(func(tx *redis.Tx) error {
		uid, err := tx.Get(ek).Result()
		if err != nil {
			return ttnredis.ConvertError(err)
		}
		if err := tx.Watch(r.uidKey(uid)).Err(); err != nil {
			return ttnredis.ConvertError(err)
		}
		pb, err = r.set(tx, uid, gets, f)
		return err
	}, ek)
	if err != nil {
		return nil, err
	}
	return pb, nil
}

// SetByID sets device by appID, devID.
func (r *DeviceRegistry) SetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(pb *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error) {
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: appID,
		DeviceID:               devID,
	}
	if err := ids.ValidateContext(ctx); err != nil {
		return nil, err
	}
	uid := unique.ID(ctx, ids)

	defer trace.StartRegion(ctx, "set end device by id").End()

	var pb *ttnpb.EndDevice
	err := r.Redis.Watch(func(tx *redis.Tx) error {
		var err error
		pb, err = r.set(tx, uid, gets, func(stored *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
			updated, sets, err := f(stored)
			if err != nil {
				return nil, nil, err
			}
			if stored == nil && updated != nil && (updated.ApplicationIdentifiers != appID || updated.DeviceID != devID) {
				return nil, nil, errInvalidIdentifiers
			}
			return updated, sets, nil
		})
		return err
	}, r.uidKey(uid))
	if err != nil {
		return nil, err
	}
	return pb, nil
}

// appendImplicitKeyGetPaths appends implicit ttnpb.SessionKeys get paths to paths.
func appendImplicitKeyGetPaths(paths ...string) []string {
	return append(paths,
		"session_key_id",
	)
}

func applyKeyFieldMask(dst, src *ttnpb.SessionKeys, paths ...string) (*ttnpb.SessionKeys, error) {
	if dst == nil {
		dst = &ttnpb.SessionKeys{}
	}
	if err := dst.SetFields(src, paths...); err != nil {
		return nil, err
	}
	return dst, nil
}

// KeyRegistry is an implementation of joinserver.KeyRegistry.
type KeyRegistry struct {
	Redis *ttnredis.Client
}

func (r *KeyRegistry) idKey(devEUI types.EUI64, id []byte) string {
	return r.Redis.Key("id", devEUI.String(), base64.RawStdEncoding.EncodeToString(id))
}

// GetByID gets session keys by devEUI, id.
func (r *KeyRegistry) GetByID(ctx context.Context, devEUI types.EUI64, id []byte, paths []string) (*ttnpb.SessionKeys, error) {
	if devEUI.IsZero() || len(id) == 0 {
		return nil, errInvalidIdentifiers
	}

	defer trace.StartRegion(ctx, "get session keys").End()

	pb := &ttnpb.SessionKeys{}
	if err := ttnredis.GetProto(r.Redis, r.idKey(devEUI, id)).ScanProto(pb); err != nil {
		return nil, err
	}
	return applyKeyFieldMask(&ttnpb.SessionKeys{}, pb, appendImplicitKeyGetPaths(paths...)...)
}

// SetByID sets session keys by devEUI, id.
func (r *KeyRegistry) SetByID(ctx context.Context, devEUI types.EUI64, id []byte, gets []string, f func(*ttnpb.SessionKeys) (*ttnpb.SessionKeys, []string, error)) (*ttnpb.SessionKeys, error) {
	if devEUI.IsZero() || len(id) == 0 {
		return nil, errInvalidIdentifiers
	}
	ik := r.idKey(devEUI, id)

	defer trace.StartRegion(ctx, "set session keys").End()

	var pb *ttnpb.SessionKeys
	err := r.Redis.Watch(func(tx *redis.Tx) error {
		cmd := ttnredis.GetProto(tx, ik)
		stored := &ttnpb.SessionKeys{}
		if err := cmd.ScanProto(stored); errors.IsNotFound(err) {
			stored = nil
		} else if err != nil {
			return err
		}

		gets = appendImplicitKeyGetPaths(gets...)

		var err error
		if stored != nil {
			pb = &ttnpb.SessionKeys{}
			if err := cmd.ScanProto(pb); err != nil {
				return err
			}
			pb, err = applyKeyFieldMask(nil, pb, gets...)
			if err != nil {
				return err
			}
		}

		var sets []string
		pb, sets, err = f(pb)
		if err != nil {
			return err
		}
		if stored == nil && pb == nil {
			return nil
		}
		if pb != nil && len(sets) == 0 {
			pb, err = applyKeyFieldMask(nil, stored, gets...)
			return err
		}

		var pipelined func(redis.Pipeliner) error
		if pb == nil && len(sets) == 0 {
			pipelined = func(p redis.Pipeliner) error {
				p.Del(ik)
				return nil
			}
		} else {
			if pb == nil {
				pb = &ttnpb.SessionKeys{}
			}

			updated := &ttnpb.SessionKeys{}
			if stored == nil {
				if err := ttnpb.RequireFields(sets,
					"session_key_id",
				); err != nil {
					return errInvalidFieldmask.WithCause(err)
				}
				updated, err = applyKeyFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
				if !bytes.Equal(updated.SessionKeyID, id) {
					return errInvalidIdentifiers
				}
			} else {
				if err := ttnpb.ProhibitFields(sets,
					"session_key_id",
				); err != nil {
					return errInvalidFieldmask.WithCause(err)
				}
				if err := cmd.ScanProto(updated); err != nil {
					return err
				}
				updated, err = applyKeyFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
			}
			if err := updated.ValidateFields(sets...); err != nil {
				return err
			}

			pipelined = func(p redis.Pipeliner) error {
				_, err := ttnredis.SetProto(p, ik, updated, 0)
				if err != nil {
					return err
				}
				return nil
			}
			pb, err = applyKeyFieldMask(nil, updated, gets...)
			if err != nil {
				return err
			}
		}
		_, err = tx.Pipelined(pipelined)
		if err != nil {
			return err
		}
		return nil
	}, ik)
	if err != nil {
		return nil, err
	}
	return pb, nil
}
