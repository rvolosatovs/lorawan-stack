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

package scheduling_test

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	. "go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/scheduling"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestRolloverClock(t *testing.T) {
	clock := &RolloverClock{}

	for i, stc := range []struct {
		Absolute ConcentratorTime
		Relative uint32
	}{
		{
			Absolute: ConcentratorTime(10 * time.Second),
			Relative: uint32(10000000),
		},
		{
			// 1 rollover.
			Absolute: ConcentratorTime(1<<32*time.Microsecond) + ConcentratorTime(5*time.Second),
			Relative: uint32(5000000),
		},
		{
			// 3 rollovers (1 existing + 2 server time rollovers).
			Absolute: ConcentratorTime(3<<32*time.Microsecond) + ConcentratorTime(10*time.Second),
			Relative: uint32(10000000),
		},
		{
			// 5 rollovers (3 existing + 1 concentrator timestamp rollover + 1 server time rollover).
			Absolute: ConcentratorTime(5<<32*time.Microsecond) + ConcentratorTime(1*time.Second),
			Relative: uint32(1000000),
		},
		{
			// 5 rollovers (5 existing) and advance to end of concentrator time epoch.
			Absolute: ConcentratorTime(6<<32*time.Microsecond) - ConcentratorTime(1*time.Second),
			Relative: uint32(4293967296),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			serverTime := time.Unix(0, 0).Add(time.Duration(stc.Absolute))
			// Run twice; once synchronizing with rollover detection, and once synchronizing with the known concentrator time.
			for i := 0; i < 2; i++ {
				t.Run([]string{"DetectRollover", "ResetAbsolute"}[i], func(t *testing.T) {
					if i == 0 {
						clock.Sync(stc.Relative, serverTime)
					} else {
						clock.SyncWithGatewayConcentrator(stc.Relative, serverTime, stc.Absolute)
					}

					for _, tc := range []struct {
						D        time.Duration
						Rollover bool
					}{
						{
							D:        -5 * time.Second,
							Rollover: false,
						},
						{
							D:        5 * time.Second,
							Rollover: false,
						},
						{
							D:        30 * time.Minute,
							Rollover: false,
						},
						{
							D:        2 * time.Hour,
							Rollover: true,
						},
					} {
						t.Run(tc.D.String(), func(t *testing.T) {
							a := assertions.New(t)

							v, ok := clock.FromServerTime(serverTime.Add(tc.D))
							a.So(ok, should.BeTrue)
							a.So(v, should.Equal, stc.Absolute+ConcentratorTime(tc.D))

							d := tc.D / time.Microsecond
							rollover := d > math.MaxUint32/2 || d < -math.MaxUint32/2
							a.So(rollover, should.Equal, tc.Rollover)

							if !rollover {
								ts := uint32(time.Duration(stc.Relative) + tc.D/time.Microsecond)
								v = clock.FromTimestampTime(ts)
								a.So(v, should.Equal, stc.Absolute+ConcentratorTime(tc.D))
							}
						})
					}
				})
			}
		})
	}
}

func TestSyncWithGatewayConcentrator(t *testing.T) {
	a := assertions.New(t)

	clock := &RolloverClock{}
	clock.SyncWithGatewayConcentrator(0x496054D6, time.Now(), ConcentratorTime(0xAA496054D6)*ConcentratorTime(time.Microsecond))
	v := int64(clock.FromTimestampTime(0x499D5DD6)) / int64(time.Microsecond)
	a.So(v, should.Equal, int64(0xAA499D5DD6))
}
