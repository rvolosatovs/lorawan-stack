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
	"context"
	"runtime/trace"
	"time"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"
	"go.thethings.network/lorawan-stack/pkg/errors"
	ttnredis "go.thethings.network/lorawan-stack/pkg/redis"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/unique"
)

var (
	errInvalidIdentifiers   = errors.DefineInvalidArgument("invalid_identifiers", "invalid identifiers")
	errDuplicateIdentifiers = errors.DefineAlreadyExists("duplicate_identifiers", "duplicate identifiers")
)

func applyDeviceFieldMask(dst, src *ttnpb.EndDevice, paths ...string) (*ttnpb.EndDevice, error) {
	if dst == nil {
		dst = &ttnpb.EndDevice{}
	}
	return dst, dst.SetFields(src, paths...)
}

// DeviceRegistry is an implementation of networkserver.DeviceRegistry.
type DeviceRegistry struct {
	Redis *ttnredis.Client
}

func (r *DeviceRegistry) uidKey(uid string) string {
	return r.Redis.Key("uid", uid)
}

func (r *DeviceRegistry) addrKey(addr types.DevAddr) string {
	return r.Redis.Key("addr", addr.String())
}

func (r *DeviceRegistry) euiKey(joinEUI, devEUI types.EUI64) string {
	return r.Redis.Key("eui", joinEUI.String(), devEUI.String())
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
	return applyDeviceFieldMask(nil, pb, append(paths,
		"ids.application_ids",
		"ids.device_id",
	)...)
}

// GetByEUI gets device by joinEUI, devEUI.
func (r *DeviceRegistry) GetByEUI(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, error) {
	defer trace.StartRegion(ctx, "get end device by eui").End()

	pb := &ttnpb.EndDevice{}
	if err := ttnredis.FindProto(r.Redis, r.euiKey(joinEUI, devEUI), r.uidKey).ScanProto(pb); err != nil {
		return nil, err
	}
	return applyDeviceFieldMask(nil, pb, append(paths,
		"ids.application_ids",
		"ids.dev_eui",
		"ids.device_id",
		"ids.join_eui",
	)...)
}

// RangeByAddr ranges over devices by addr.
func (r *DeviceRegistry) RangeByAddr(ctx context.Context, addr types.DevAddr, paths []string, f func(*ttnpb.EndDevice) bool) error {
	defer trace.StartRegion(ctx, "range end devices by dev_addr").End()

	return ttnredis.FindProtos(r.Redis, r.addrKey(addr), r.uidKey).Range(func() (proto.Message, func() (bool, error)) {
		pb := &ttnpb.EndDevice{}
		return pb, func() (bool, error) {
			pb, err := applyDeviceFieldMask(nil, pb, append(paths,
				"ids.application_ids",
				"ids.dev_addr",
				"ids.device_id",
			)...)
			if err != nil {
				return false, err
			}
			return f(pb), nil
		}
	})
}

func getDevAddrs(pb *ttnpb.EndDevice) (addrs struct{ current, pending *types.DevAddr }) {
	if pb == nil {
		return
	}

	if pb.Session != nil {
		var addr types.DevAddr
		copy(addr[:], pb.Session.DevAddr[:])
		addrs.current = &addr
	}
	if pb.PendingSession != nil {
		var addr types.DevAddr
		copy(addr[:], pb.PendingSession.DevAddr[:])
		addrs.pending = &addr
	}
	return addrs
}

func equalAddr(x, y *types.DevAddr) bool {
	if x == nil || y == nil {
		return x == y
	}
	return x.Equal(*y)
}

func equalEUI(x, y *types.EUI64) bool {
	if x == nil || y == nil {
		return x == y
	}
	return x.Equal(*y)
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
	uk := r.uidKey(uid)

	defer trace.StartRegion(ctx, "set end device by id").End()

	var pb *ttnpb.EndDevice
	err := r.Redis.Watch(func(tx *redis.Tx) error {
		cmd := ttnredis.GetProto(tx, uk)
		stored := &ttnpb.EndDevice{}
		if err := cmd.ScanProto(stored); errors.IsNotFound(err) {
			stored = nil
		} else if err != nil {
			return err
		}

		storedAddrs := getDevAddrs(stored)

		gets = append(gets,
			"created_at",
			"ids.application_ids",
			"ids.dev_addr",
			"ids.dev_eui",
			"ids.device_id",
			"ids.join_eui",
			"updated_at",
		)

		var err error
		if stored != nil {
			pb = &ttnpb.EndDevice{}
			if err := cmd.ScanProto(pb); err != nil {
				return err
			}
			pb, err = applyDeviceFieldMask(nil, pb, gets...)
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

		var pipelined func(redis.Pipeliner) error
		if pb == nil {
			pipelined = func(p redis.Pipeliner) error {
				p.Del(uk)
				if stored.JoinEUI != nil && stored.DevEUI != nil {
					p.Del(r.euiKey(*stored.JoinEUI, *stored.DevEUI))
				}
				if storedAddrs.pending != nil {
					p.SRem(r.addrKey(*storedAddrs.pending), uid)
				}
				if storedAddrs.current != nil {
					p.SRem(r.addrKey(*storedAddrs.current), uid)
				}
				return nil
			}
		} else {
			if pb.ApplicationIdentifiers != appID || pb.DeviceID != devID {
				return errInvalidIdentifiers
			}

			pb.UpdatedAt = time.Now().UTC()
			sets = append(sets, "updated_at")

			updated := &ttnpb.EndDevice{}
			if stored == nil {
				sets = append(sets,
					"ids.application_ids",
					"ids.dev_addr",
					"ids.dev_eui",
					"ids.device_id",
					"ids.join_eui",
				)

				pb.CreatedAt = pb.UpdatedAt
				sets = append(sets, "created_at")

				updated, err = applyDeviceFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
			} else {
				if err := cmd.ScanProto(updated); err != nil {
					return err
				}
				updated, err = applyDeviceFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
				if !equalEUI(stored.JoinEUI, updated.JoinEUI) || !equalEUI(stored.DevEUI, updated.DevEUI) ||
					stored.ApplicationIdentifiers != updated.ApplicationIdentifiers || stored.DeviceID != updated.DeviceID {
					return errInvalidIdentifiers
				}
			}
			if err := updated.ValidateFields(sets...); err != nil {
				return err
			}
			updatedAddrs := getDevAddrs(updated)

			pipelined = func(p redis.Pipeliner) error {
				if stored == nil && updated.JoinEUI != nil && updated.DevEUI != nil {
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

				_, err := ttnredis.SetProto(p, uk, updated, 0)
				if err != nil {
					return err
				}

				if storedAddrs.pending != nil && !equalAddr(storedAddrs.pending, updatedAddrs.pending) && !equalAddr(storedAddrs.pending, updatedAddrs.current) {
					p.SRem(r.addrKey(*storedAddrs.pending), uid)
				}
				if storedAddrs.current != nil && !equalAddr(storedAddrs.current, updatedAddrs.pending) && !equalAddr(storedAddrs.current, updatedAddrs.current) {
					p.SRem(r.addrKey(*storedAddrs.current), uid)
				}
				if updatedAddrs.pending != nil && !equalAddr(updatedAddrs.pending, storedAddrs.pending) && !equalAddr(updatedAddrs.pending, storedAddrs.current) {
					p.SAdd(r.addrKey(*updatedAddrs.pending), uid)
				}
				if updatedAddrs.current != nil && !equalAddr(updatedAddrs.current, storedAddrs.pending) && !equalAddr(updatedAddrs.current, storedAddrs.current) {
					p.SAdd(r.addrKey(*updatedAddrs.current), uid)
				}

				pb, err = applyDeviceFieldMask(nil, updated, gets...)
				if err != nil {
					return err
				}
				return nil
			}
		}
		_, err = tx.Pipelined(pipelined)
		if err != nil {
			return err
		}
		return nil
	}, uk)
	if err != nil {
		return nil, err
	}
	return pb, nil
}
