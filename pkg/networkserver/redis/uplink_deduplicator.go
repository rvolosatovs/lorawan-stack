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
	"encoding/base64"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	ttnredis "go.thethings.network/lorawan-stack/v3/pkg/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// UplinkDeduplicator is an implementation of networkserver.UplinkDeduplicator.
type UplinkDeduplicator struct {
	Redis *ttnredis.Client
}

// NewUplinkDeduplicator returns a new uplink deduplicator.
func NewUplinkDeduplicator(cl *ttnredis.Client) *UplinkDeduplicator {
	return &UplinkDeduplicator{
		Redis: cl,
	}
}

var keyEncoding = base64.RawStdEncoding

func uplinkHash(ctx context.Context, up *ttnpb.UplinkMessage) (string, error) {
	drBytes := make([]byte, up.Settings.DataRate.Modulation.Size())
	_, err := up.Settings.DataRate.Modulation.MarshalTo(drBytes)
	if err != nil {
		return "", err
	}
	h := fnv.New64a()
	_, _ = h.Write(up.RawPayload)
	return ttnredis.Key(
		keyEncoding.EncodeToString(h.Sum(nil)),
		// NOTE: Data rate and frequency are included in the key to support retransmissions.
		strconv.FormatUint(up.Settings.Frequency, 32),
		keyEncoding.EncodeToString(drBytes),
	), nil
}

// DeduplicateUplink deduplicates up for window. Since highest precision allowed by Redis is millisecondsm, window is truncated to milliseconds.
func (d *UplinkDeduplicator) DeduplicateUplink(ctx context.Context, up *ttnpb.UplinkMessage, window time.Duration) (bool, error) {
	h, err := uplinkHash(ctx, up)
	if err != nil {
		return false, err
	}
	msgs := make([]proto.Message, 0, len(up.RxMetadata))
	for _, md := range up.RxMetadata {
		msgs = append(msgs, md)
	}
	return ttnredis.DeduplicateProtos(ctx, d.Redis, d.Redis.Key(h), window, msgs...)
}

// DeduplicateUplink returns accumulated metadata for up.
func (d *UplinkDeduplicator) AccumulatedMetadata(ctx context.Context, up *ttnpb.UplinkMessage) ([]*ttnpb.RxMetadata, error) {
	h, err := uplinkHash(ctx, up)
	if err != nil {
		return nil, err
	}
	var mds []*ttnpb.RxMetadata
	return mds, ttnredis.ListProtos(ctx, d.Redis, d.Redis.Key(ttnredis.DeduplicationListKey(h))).Range(func() (proto.Message, func() (bool, error)) {
		md := &ttnpb.RxMetadata{}
		return md, func() (bool, error) {
			mds = append(mds, md)
			return true, nil
		}
	})
}
