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

// Package redis provides Redis implementations of interfaces used by networkserver.
package redis

import (
	"encoding/base64"
	"hash/fnv"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	ttnredis "go.thethings.network/lorawan-stack/v3/pkg/redis"
)

//go:generate go run ./generate_scripts.go

func deviceUIDKey(cl *ttnredis.Client, uid string) string {
	return cl.Key("uid", uid)
}

func deviceUIDLastInvalidationKey(cl *ttnredis.Client, uid string) string {
	return ttnredis.Key(deviceUIDKey(cl, uid), "last-invalidation")
}

var keyEncoding = base64.RawStdEncoding

func uplinkPayloadHash(b []byte) string {
	h := fnv.New64a()
	_, _ = h.Write(b)
	return keyEncoding.EncodeToString(h.Sum(nil))
}

var errDatabaseCorruption = errors.DefineCorruption("database_corruption", "database is corrupted")
