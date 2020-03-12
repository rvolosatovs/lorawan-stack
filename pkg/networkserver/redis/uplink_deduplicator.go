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
	"hash"
	"hash/fnv"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	ttnredis "go.thethings.network/lorawan-stack/pkg/redis"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// UplinkDeduplicator is an implementation of networkserver.UplinkDeduplicator.
type UplinkDeduplicator struct {
	Redis    *ttnredis.Client
	hashPool *sync.Pool
}

// NewUplinkDeduplicator returns a new uplink deduplicator.
func NewUplinkDeduplicator(cl *ttnredis.Client) *UplinkDeduplicator {
	return &UplinkDeduplicator{
		Redis: cl,
		hashPool: &sync.Pool{
			New: func() interface{} {
				return fnv.New64a()
			},
		},
	}
}

func (d *UplinkDeduplicator) uplinkHash(ctx context.Context, up *ttnpb.UplinkMessage) string {
	h := d.hashPool.Get().(hash.Hash64)
	_, _ = h.Write(up.RawPayload)

	s := string(h.Sum(nil))

	h.Reset()
	d.hashPool.Put(h)

	return s
}

// DeduplicateUplink deduplicates up for window. Since highest precision allowed by Redis is millisecondsm, window is truncated to milliseconds.
func (d *UplinkDeduplicator) DeduplicateUplink(ctx context.Context, up *ttnpb.UplinkMessage, window time.Duration) (bool, error) {
	msgs := make([]proto.Message, 0, len(up.RxMetadata))
	for _, md := range up.RxMetadata {
		msgs = append(msgs, md)
	}
	return ttnredis.DeduplicateProtos(ctx, d.Redis, d.Redis.Key(d.uplinkHash(ctx, up)), window, msgs...)
}

// DeduplicateUplink returns accumulated metadata for up.
func (d *UplinkDeduplicator) AccumulatedMetadata(ctx context.Context, up *ttnpb.UplinkMessage) ([]*ttnpb.RxMetadata, error) {
	var mds []*ttnpb.RxMetadata
	return mds, ttnredis.ListProtos(ctx, d.Redis, d.Redis.Key(ttnredis.DeduplicationListKey(d.uplinkHash(ctx, up)))).Range(func() (proto.Message, func() (bool, error)) {
		md := &ttnpb.RxMetadata{}
		return md, func() (bool, error) {
			mds = append(mds, md)
			return true, nil
		}
	})
}
