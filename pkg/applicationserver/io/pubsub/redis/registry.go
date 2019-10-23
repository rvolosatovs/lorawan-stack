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
	"time"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"
	"go.thethings.network/lorawan-stack/pkg/applicationserver/io/pubsub"
	"go.thethings.network/lorawan-stack/pkg/errors"
	ttnredis "go.thethings.network/lorawan-stack/pkg/redis"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
)

var (
	errInvalidFieldmask   = errors.DefineInvalidArgument("invalid_fieldmask", "invalid fieldmask")
	errInvalidIdentifiers = errors.DefineInvalidArgument("invalid_identifiers", "invalid identifiers")
	errReadOnlyField      = errors.DefineInvalidArgument("read_only_field", "read-only field `{field}`")
)

// appendImplicitPubSubGetPaths appends implicit ttnpb.ApplicationPubSub get paths to paths.
func appendImplicitPubSubGetPaths(paths ...string) []string {
	return append(append(make([]string, 0, 4+len(paths)),
		"created_at",
		"ids",
		"provider",
		"updated_at",
	), paths...)
}

func applyPubSubFieldMask(dst, src *ttnpb.ApplicationPubSub, paths ...string) (*ttnpb.ApplicationPubSub, error) {
	if dst == nil {
		dst = &ttnpb.ApplicationPubSub{}
	}
	return dst, dst.SetFields(src, paths...)
}

// PubSubRegistry is a Redis PubSub registry.
type PubSubRegistry struct {
	Redis *ttnredis.Client
}

func (r *PubSubRegistry) allKey(ctx context.Context) string {
	return r.Redis.Key("all")
}

func (r *PubSubRegistry) appKey(uid string) string {
	return r.Redis.Key("uid", uid)
}

func (r *PubSubRegistry) uidKey(appUID, id string) string {
	return r.Redis.Key("uid", appUID, id)
}

func (r *PubSubRegistry) makeUIDKeyFunc(appUID string) func(id string) string {
	return func(id string) string {
		return r.uidKey(appUID, id)
	}
}

// Get implements pubsub.Registry.
func (r PubSubRegistry) Get(ctx context.Context, ids ttnpb.ApplicationPubSubIdentifiers, paths []string) (*ttnpb.ApplicationPubSub, error) {
	pb := &ttnpb.ApplicationPubSub{}
	if err := ttnredis.GetProto(r.Redis, r.uidKey(unique.ID(ctx, ids.ApplicationIdentifiers), ids.PubSubID)).ScanProto(pb); err != nil {
		return nil, err
	}
	return applyPubSubFieldMask(nil, pb, appendImplicitPubSubGetPaths(paths...)...)
}

var errApplicationUID = errors.DefineCorruption("application_uid", "invalid application UID `{application_uid}`")

// Range implements pubsub.Registry.
func (r PubSubRegistry) Range(ctx context.Context, paths []string, f func(context.Context, ttnpb.ApplicationIdentifiers, *ttnpb.ApplicationPubSub) bool) error {
	uids, err := r.Redis.SMembers(r.allKey(ctx)).Result()
	if err != nil {
		return err
	}
	for _, uid := range uids {
		appUID, psID := pubsub.SplitPubSubUID(uid)
		ctx, err := unique.WithContext(ctx, appUID)
		if err != nil {
			return errApplicationUID.WithCause(err).WithAttributes("application_uid", appUID, "pub_sub_id", psID)
		}
		ids, err := unique.ToApplicationID(appUID)
		if err != nil {
			return errApplicationUID.WithCause(err).WithAttributes("application_uid", appUID, "pub_sub_id", psID)
		}
		pb := &ttnpb.ApplicationPubSub{}
		if err := ttnredis.GetProto(r.Redis, r.uidKey(appUID, psID)).ScanProto(pb); err != nil {
			return err
		}
		if err != nil {
			return errApplicationUID.WithCause(err).WithAttributes("application_uid", appUID, "pub_sub_id", psID)
		}
		pb, err = applyPubSubFieldMask(nil, pb, paths...)
		if err != nil {
			return err
		}
		if !f(ctx, ids, pb) {
			return nil
		}
	}
	return nil
}

// List implements pubsub.Registry.
func (r PubSubRegistry) List(ctx context.Context, ids ttnpb.ApplicationIdentifiers, paths []string) ([]*ttnpb.ApplicationPubSub, error) {
	var pbs []*ttnpb.ApplicationPubSub
	appUID := unique.ID(ctx, ids)
	err := ttnredis.FindProtos(r.Redis, r.appKey(appUID), r.makeUIDKeyFunc(appUID)).Range(func() (proto.Message, func() (bool, error)) {
		pb := &ttnpb.ApplicationPubSub{}
		return pb, func() (bool, error) {
			pb, err := applyPubSubFieldMask(nil, pb, appendImplicitPubSubGetPaths(paths...)...)
			if err != nil {
				return false, err
			}
			pbs = append(pbs, pb)
			return true, nil
		}
	})
	if err != nil {
		return nil, err
	}
	return pbs, nil
}

// Set implements pubsub.Registry.
func (r PubSubRegistry) Set(ctx context.Context, ids ttnpb.ApplicationPubSubIdentifiers, gets []string, f func(*ttnpb.ApplicationPubSub) (*ttnpb.ApplicationPubSub, []string, error)) (*ttnpb.ApplicationPubSub, error) {
	appUID := unique.ID(ctx, ids.ApplicationIdentifiers)
	ik := r.uidKey(appUID, ids.PubSubID)

	var pb *ttnpb.ApplicationPubSub
	err := r.Redis.Watch(func(tx *redis.Tx) error {
		cmd := ttnredis.GetProto(tx, ik)
		stored := &ttnpb.ApplicationPubSub{}
		if err := cmd.ScanProto(stored); errors.IsNotFound(err) {
			stored = nil
		} else if err != nil {
			return err
		}

		gets = appendImplicitPubSubGetPaths(gets...)

		var err error
		if stored != nil {
			pb = &ttnpb.ApplicationPubSub{}
			if err := cmd.ScanProto(pb); err != nil {
				return err
			}
			pb, err = applyPubSubFieldMask(nil, pb, gets...)
			if err != nil {
				return err
			}
		}

		var sets []string
		pb, sets, err = f(pb)
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
			pb, err = applyPubSubFieldMask(nil, stored, gets...)
			return err
		}

		var pipelined func(redis.Pipeliner) error
		if pb == nil && len(sets) == 0 {
			pipelined = func(p redis.Pipeliner) error {
				p.Del(ik)
				p.SRem(r.appKey(appUID), stored.PubSubID)
				p.SRem(r.allKey(ctx), pubsub.PubSubUID(appUID, stored.PubSubID))
				return nil
			}
		} else {
			if pb == nil {
				pb = &ttnpb.ApplicationPubSub{}
			}

			pb.UpdatedAt = time.Now().UTC()
			sets = append(append(sets[:0:0], sets...),
				"updated_at",
			)

			updated := &ttnpb.ApplicationPubSub{}
			if stored == nil {
				if err := ttnpb.RequireFields(sets,
					"ids.application_ids",
					"ids.pub_sub_id",
				); err != nil {
					return errInvalidFieldmask.WithCause(err)
				}

				pb.CreatedAt = pb.UpdatedAt
				sets = append(sets, "created_at")

				updated, err = applyPubSubFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
				if updated.ApplicationID != ids.ApplicationID || updated.PubSubID != ids.PubSubID {
					return errInvalidIdentifiers
				}
			} else {
				if ttnpb.HasAnyField(sets, "ids.application_ids.application_id") && pb.ApplicationID != stored.ApplicationID {
					return errReadOnlyField.WithAttributes("field", "ids.application_ids.application_id")
				}
				if ttnpb.HasAnyField(sets, "ids.pub_sub_id") && pb.PubSubID != stored.PubSubID {
					return errReadOnlyField.WithAttributes("field", "ids.pub_sub_id")
				}
				if err := cmd.ScanProto(updated); err != nil {
					return err
				}
				updated, err = applyPubSubFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
			}
			if err := updated.ValidateFields(sets...); err != nil {
				return err
			}

			pipelined = func(p redis.Pipeliner) error {
				if _, err := ttnredis.SetProto(p, ik, updated, 0); err != nil {
					return err
				}
				p.SAdd(r.appKey(appUID), updated.PubSubID)
				p.SAdd(r.allKey(ctx), pubsub.PubSubUID(appUID, updated.PubSubID))
				return nil
			}

			pb, err = applyPubSubFieldMask(nil, updated, gets...)
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
