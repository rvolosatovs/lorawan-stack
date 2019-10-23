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

// appendImplicitWebhookGetPaths appends implicit ttnpb.ApplicationWebhook get paths to paths.
func appendImplicitWebhookGetPaths(paths ...string) []string {
	return append(append(make([]string, 0, 3+len(paths)),
		"created_at",
		"ids",
		"updated_at",
	), paths...)
}

func applyWebhookFieldMask(dst, src *ttnpb.ApplicationWebhook, paths ...string) (*ttnpb.ApplicationWebhook, error) {
	if dst == nil {
		dst = &ttnpb.ApplicationWebhook{}
	}
	return dst, dst.SetFields(src, paths...)
}

// WebhookRegistry is a Redis webhook registry.
type WebhookRegistry struct {
	Redis *ttnredis.Client
}

func (r *WebhookRegistry) appKey(uid string) string {
	return r.Redis.Key("uid", uid)
}

func (r *WebhookRegistry) idKey(appUID, id string) string {
	return r.Redis.Key("uid", appUID, id)
}

func (r *WebhookRegistry) makeIDKeyFunc(appUID string) func(id string) string {
	return func(id string) string {
		return r.idKey(appUID, id)
	}
}

// Get implements WebhookRegistry.
func (r WebhookRegistry) Get(ctx context.Context, ids ttnpb.ApplicationWebhookIdentifiers, paths []string) (*ttnpb.ApplicationWebhook, error) {
	pb := &ttnpb.ApplicationWebhook{}
	if err := ttnredis.GetProto(r.Redis, r.idKey(unique.ID(ctx, ids.ApplicationIdentifiers), ids.WebhookID)).ScanProto(pb); err != nil {
		return nil, err
	}
	return applyWebhookFieldMask(nil, pb, appendImplicitWebhookGetPaths(paths...)...)
}

// List implements WebhookRegistry.
func (r WebhookRegistry) List(ctx context.Context, ids ttnpb.ApplicationIdentifiers, paths []string) ([]*ttnpb.ApplicationWebhook, error) {
	var pbs []*ttnpb.ApplicationWebhook
	appUID := unique.ID(ctx, ids)
	err := ttnredis.FindProtos(r.Redis, r.appKey(appUID), r.makeIDKeyFunc(appUID)).Range(func() (proto.Message, func() (bool, error)) {
		pb := &ttnpb.ApplicationWebhook{}
		return pb, func() (bool, error) {
			pb, err := applyWebhookFieldMask(nil, pb, appendImplicitWebhookGetPaths(paths...)...)
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

// Set implements WebhookRegistry.
func (r WebhookRegistry) Set(ctx context.Context, ids ttnpb.ApplicationWebhookIdentifiers, gets []string, f func(*ttnpb.ApplicationWebhook) (*ttnpb.ApplicationWebhook, []string, error)) (*ttnpb.ApplicationWebhook, error) {
	appUID := unique.ID(ctx, ids.ApplicationIdentifiers)
	ik := r.idKey(appUID, ids.WebhookID)

	var pb *ttnpb.ApplicationWebhook
	err := r.Redis.Watch(func(tx *redis.Tx) error {
		cmd := ttnredis.GetProto(tx, ik)
		stored := &ttnpb.ApplicationWebhook{}
		if err := cmd.ScanProto(stored); errors.IsNotFound(err) {
			stored = nil
		} else if err != nil {
			return err
		}

		gets = appendImplicitWebhookGetPaths(gets...)

		var err error
		if stored != nil {
			pb = &ttnpb.ApplicationWebhook{}
			if err := cmd.ScanProto(pb); err != nil {
				return err
			}
			pb, err = applyWebhookFieldMask(nil, pb, gets...)
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
			pb, err = applyWebhookFieldMask(nil, stored, gets...)
			return err
		}

		var pipelined func(redis.Pipeliner) error
		if pb == nil && len(sets) == 0 {
			pipelined = func(p redis.Pipeliner) error {
				p.Del(ik)
				p.SRem(r.appKey(appUID), stored.WebhookID)
				return nil
			}
		} else {
			if pb == nil {
				pb = &ttnpb.ApplicationWebhook{}
			}

			pb.UpdatedAt = time.Now().UTC()
			sets = append(append(sets[:0:0], sets...),
				"updated_at",
			)

			updated := &ttnpb.ApplicationWebhook{}
			if stored == nil {
				if err := ttnpb.RequireFields(sets,
					"ids.application_ids",
					"ids.webhook_id",
				); err != nil {
					return errInvalidFieldmask.WithCause(err)
				}

				pb.CreatedAt = pb.UpdatedAt
				sets = append(sets, "created_at")

				updated, err = applyWebhookFieldMask(updated, pb, sets...)
				if err != nil {
					return err
				}
				if updated.ApplicationID != ids.ApplicationID || updated.WebhookID != ids.WebhookID {
					return errInvalidIdentifiers
				}
			} else {
				if ttnpb.HasAnyField(sets, "ids.application_ids.application_id") && pb.ApplicationID != stored.ApplicationID {
					return errReadOnlyField.WithAttributes("field", "ids.application_ids.application_id")
				}
				if ttnpb.HasAnyField(sets, "ids.webhook_id") && pb.WebhookID != stored.WebhookID {
					return errReadOnlyField.WithAttributes("field", "ids.webhook_id")
				}
				if err := cmd.ScanProto(updated); err != nil {
					return err
				}
				updated, err = applyWebhookFieldMask(updated, pb, sets...)
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
				p.SAdd(r.appKey(appUID), updated.WebhookID)
				return nil
			}

			pb, err = applyWebhookFieldMask(nil, updated, gets...)
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
