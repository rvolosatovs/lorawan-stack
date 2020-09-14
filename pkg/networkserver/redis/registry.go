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
	"io"
	"math/rand"
	"runtime/trace"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gogo/protobuf/proto"
	ulid "github.com/oklog/ulid/v2"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	ttnredis "go.thethings.network/lorawan-stack/v3/pkg/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
)

var (
	errInvalidFieldmask     = errors.DefineInvalidArgument("invalid_fieldmask", "invalid fieldmask")
	errInvalidIdentifiers   = errors.DefineInvalidArgument("invalid_identifiers", "invalid identifiers")
	errDuplicateIdentifiers = errors.DefineAlreadyExists("duplicate_identifiers", "duplicate identifiers")
	errReadOnlyField        = errors.DefineInvalidArgument("read_only_field", "read-only field `{field}`")
)

// DeviceRegistry is an implementation of networkserver.DeviceRegistry.
type DeviceRegistry struct {
	Redis   *ttnredis.Client
	LockTTL time.Duration

	entropyMu *sync.Mutex
	entropy   io.Reader
}

func (r *DeviceRegistry) Init() error {
	if err := ttnredis.InitMutex(r.Redis); err != nil {
		return err
	}
	r.entropyMu = &sync.Mutex{}
	r.entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 1000)
	return nil
}

func (r *DeviceRegistry) uidKey(uid string) string {
	return deviceUIDKey(r.Redis, uid)
}

func (r *DeviceRegistry) addrKey(addr types.DevAddr) string {
	return r.Redis.Key("addr", addr.String())
}

func (r *DeviceRegistry) euiKey(joinEUI, devEUI types.EUI64) string {
	return r.Redis.Key("eui", joinEUI.String(), devEUI.String())
}

// GetByID gets device by appID, devID.
func (r *DeviceRegistry) GetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string) (*ttnpb.EndDevice, context.Context, error) {
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: appID,
		DeviceID:               devID,
	}
	if err := ids.ValidateContext(ctx); err != nil {
		return nil, ctx, err
	}

	defer trace.StartRegion(ctx, "get end device by id").End()

	pb := &ttnpb.EndDevice{}
	if err := ttnredis.GetProto(r.Redis, r.uidKey(unique.ID(ctx, ids))).ScanProto(pb); err != nil {
		return nil, ctx, err
	}
	pb, err := ttnpb.FilterGetEndDevice(pb, paths...)
	if err != nil {
		return nil, ctx, err
	}
	return pb, ctx, nil
}

// GetByEUI gets device by joinEUI, devEUI.
func (r *DeviceRegistry) GetByEUI(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, context.Context, error) {
	defer trace.StartRegion(ctx, "get end device by eui").End()

	pb := &ttnpb.EndDevice{}
	if err := ttnredis.FindProto(r.Redis, r.euiKey(joinEUI, devEUI), func(uid string) (string, error) {
		return r.uidKey(uid), nil
	}).ScanProto(pb); err != nil {
		return nil, ctx, err
	}
	pb, err := ttnpb.FilterGetEndDevice(pb, paths...)
	if err != nil {
		return nil, ctx, err
	}
	return pb, ctx, nil
}

// RangeByAddr ranges over devices by addr.
func (r *DeviceRegistry) RangeByAddr(ctx context.Context, addr types.DevAddr, paths []string, f func(context.Context, *ttnpb.EndDevice) bool) error {
	defer trace.StartRegion(ctx, "range end devices by dev_addr").End()

	return ttnredis.FindProtos(r.Redis, r.addrKey(addr), r.uidKey).Range(func() (proto.Message, func() (bool, error)) {
		pb := &ttnpb.EndDevice{}
		return pb, func() (bool, error) {
			pb, err := ttnpb.FilterGetEndDevice(pb, paths...)
			if err != nil {
				return false, err
			}
			return f(ctx, pb), nil
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

func equalEUI64(x, y *types.EUI64) bool {
	if x == nil || y == nil {
		return x == y
	}
	return x.Equal(*y)
}

// SetByID sets device by appID, devID.
func (r *DeviceRegistry) SetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, gets []string, f func(ctx context.Context, pb *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, context.Context, error) {
	ids := ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: appID,
		DeviceID:               devID,
	}
	if err := ids.ValidateContext(ctx); err != nil {
		return nil, ctx, err
	}
	uid := unique.ID(ctx, ids)
	uk := r.uidKey(uid)

	defer trace.StartRegion(ctx, "set end device by id").End()

	var pb *ttnpb.EndDevice
	r.entropyMu.Lock()
	lockID, err := ulid.New(ulid.Timestamp(time.Now()), r.entropy)
	r.entropyMu.Unlock()
	if err != nil {
		return nil, ctx, err
	}
	lockIDStr := lockID.String()
	if err = ttnredis.LockedWatch(ctx, r.Redis, uk, lockIDStr, r.LockTTL, func(tx *redis.Tx) error {
		cmd := ttnredis.GetProto(tx, uk)
		stored := &ttnpb.EndDevice{}
		if err := cmd.ScanProto(stored); errors.IsNotFound(err) {
			stored = nil
		} else if err != nil {
			return err
		}

		var err error
		if stored != nil {
			pb = &ttnpb.EndDevice{}
			if err := cmd.ScanProto(pb); err != nil {
				return err
			}
			pb, err = ttnpb.FilterGetEndDevice(pb, gets...)
			if err != nil {
				return err
			}
		}

		var sets []string
		pb, sets, err = f(ctx, pb)
		if err != nil {
			return err
		}
		if err := ttnpb.ProhibitFields(sets,
			"created_at",
			"updated_at",
		); err != nil {
			return errInvalidFieldmask.WithCause(err)
		}

		if stored == nil && pb == nil {
			return nil
		}
		if pb != nil && len(sets) == 0 {
			pb, err = ttnpb.FilterGetEndDevice(stored, gets...)
			return err
		}

		var pipelined func(redis.Pipeliner) error
		if pb == nil && len(sets) == 0 {
			pipelined = func(p redis.Pipeliner) error {
				p.Del(uk)
				p.Del(deviceUIDLastInvalidationKey(r.Redis, uid))
				if stored.JoinEUI != nil && stored.DevEUI != nil {
					p.Del(r.euiKey(*stored.JoinEUI, *stored.DevEUI))
				}
				if stored.PendingSession != nil {
					p.SRem(r.addrKey(stored.PendingSession.DevAddr), uid)
				}
				if stored.Session != nil {
					p.SRem(r.addrKey(stored.Session.DevAddr), uid)
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

			var preSet []func(redis.Pipeliner)
			updated := &ttnpb.EndDevice{}
			if stored == nil {
				if err := ttnpb.RequireFields(sets,
					"ids.application_ids",
					"ids.device_id",
				); err != nil {
					return errInvalidFieldmask.WithCause(err)
				}

				pb.CreatedAt = pb.UpdatedAt
				sets = append(sets, "created_at")

				updated, err = ttnpb.ApplyEndDeviceFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
				if updated.ApplicationIdentifiers != appID || updated.DeviceID != devID {
					return errInvalidIdentifiers.New()
				}
				if updated.JoinEUI != nil && updated.DevEUI != nil {
					ek := r.euiKey(*updated.JoinEUI, *updated.DevEUI)

					if err := ttnredis.LockMutex(ctx, tx, ek, lockIDStr, r.LockTTL); err != nil {
						return err
					}
					if err := tx.Watch(ek).Err(); err != nil {
						return err
					}
					i, err := tx.Exists(ek).Result()
					if err != nil {
						return err
					}
					if i != 0 {
						return errDuplicateIdentifiers.New()
					}
					preSet = append(preSet, func(p redis.Pipeliner) {
						p.Set(ek, uid, 0)
						ttnredis.UnlockMutex(p, ek, lockIDStr, r.LockTTL)
					})
				}
			} else {
				if ttnpb.HasAnyField(sets, "ids.application_ids.application_id") && pb.ApplicationID != stored.ApplicationID {
					return errReadOnlyField.WithAttributes("field", "ids.application_ids.application_id")
				}
				if ttnpb.HasAnyField(sets, "ids.device_id") && pb.DeviceID != stored.DeviceID {
					return errReadOnlyField.WithAttributes("field", "ids.device_id")
				}
				if ttnpb.HasAnyField(sets, "ids.join_eui") && !equalEUI64(pb.JoinEUI, stored.JoinEUI) {
					return errReadOnlyField.WithAttributes("field", "ids.join_eui")
				}
				if ttnpb.HasAnyField(sets, "ids.dev_eui") && !equalEUI64(pb.DevEUI, stored.DevEUI) {
					return errReadOnlyField.WithAttributes("field", "ids.dev_eui")
				}
				if err := cmd.ScanProto(updated); err != nil {
					return err
				}
				updated, err = ttnpb.ApplyEndDeviceFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
			}
			if err := updated.ValidateFields(sets...); err != nil {
				return err
			}
			pb, err = ttnpb.FilterGetEndDevice(updated, gets...)
			if err != nil {
				return err
			}
			pipelined = func(p redis.Pipeliner) error {
				for _, f := range preSet {
					f(p)
				}
				_, err := ttnredis.SetProto(p, uk, updated, 0)
				if err != nil {
					return err
				}
				storedAddrs := getDevAddrs(stored)
				updatedAddrs := getDevAddrs(updated)
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
				return nil
			}
		}
		_, err = tx.TxPipelined(pipelined)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, ctx, ttnredis.ConvertError(err)
	}
	return pb, ctx, nil
}
