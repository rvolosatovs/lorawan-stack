// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package redis_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/events/internal/eventstest"
	"go.thethings.network/lorawan-stack/v3/pkg/events/redis"
	ttnredis "go.thethings.network/lorawan-stack/v3/pkg/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
)

var redisConfig = func() ttnredis.Config {
	var err error
	config := ttnredis.Config{
		Address:       "localhost:6379",
		Database:      1,
		RootNamespace: []string{"test"},
	}
	if address := os.Getenv("REDIS_ADDRESS"); address != "" {
		config.Address = address
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		config.Database, err = strconv.Atoi(db)
		if err != nil {
			panic(err)
		}
	}
	if prefix := os.Getenv("REDIS_PREFIX"); prefix != "" {
		config.RootNamespace = []string{prefix}
	}
	return config
}()

func Example() {
	// The task starter is used for automatic re-subscription on failure.
	taskStarter := component.StartTaskFunc(component.DefaultStartTask)

	redisPubSub := redis.NewPubSub(context.TODO(), taskStarter, ttnredis.Config{
		// Config here...
	})

	// Replace the default pubsub so that we will now publish to Redis.
	events.SetDefaultPubSub(redisPubSub)
}

var timeout = (1 << 10) * test.Delay

func TestRedisPubSub(t *testing.T) {
	events.IncludeCaller = true
	taskStarter := component.StartTaskFunc(component.DefaultStartTask)

	test.RunTest(t, test.TestConfig{
		Timeout: timeout,
		Func: func(ctx context.Context, a *assertions.Assertion) {
			pubsub := redis.NewPubSub(ctx, taskStarter, redisConfig)
			defer pubsub.Close(ctx)

			time.Sleep(timeout / 10)

			eventstest.TestBackend(ctx, t, a, pubsub)
		},
	})
}
