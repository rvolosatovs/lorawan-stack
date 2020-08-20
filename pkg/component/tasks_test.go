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

package component_test

import (
	"context"
	"sync"
	"testing"

	"github.com/smartystreets/assertions"
	. "go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestTasks(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	c, err := New(test.GetLogger(t), &Config{})
	a.So(err, should.BeNil)

	// Register a one-off task.
	oneOffWg := sync.WaitGroup{}
	oneOffWg.Add(1)
	c.RegisterTask(&TaskConfig{
		Context: ctx,
		ID:      "one_off",
		Func: func(_ context.Context) error {
			oneOffWg.Done()
			return nil
		},
		Restart: TaskRestartNever,
		Backoff: TestTaskBackoffConfig,
	})
	// Register an always restarting task.
	restartingWg := sync.WaitGroup{}
	restartingWg.Add(5)
	i := 0
	c.RegisterTask(&TaskConfig{
		Context: ctx,
		ID:      "restarts",
		Func: func(_ context.Context) error {
			i++
			if i <= 5 {
				restartingWg.Done()
			}
			return nil
		},
		Restart: TaskRestartAlways,
		Backoff: TestTaskBackoffConfig,
	})

	// Register a task that restarts on failure.
	failingWg := sync.WaitGroup{}
	failingWg.Add(1)
	j := 0
	c.RegisterTask(&TaskConfig{
		Context: ctx,
		ID:      "restarts_on_failure",
		Func: func(_ context.Context) error {
			j++
			if j < 5 {
				return errors.New("failed")
			}
			failingWg.Done()
			return nil
		},
		Restart: TaskRestartOnFailure,
		Backoff: TestTaskBackoffConfig,
	})

	// Wait for all invocations.
	test.Must(nil, c.Start())
	defer c.Close()
	oneOffWg.Wait()
	restartingWg.Wait()
	failingWg.Wait()
}
