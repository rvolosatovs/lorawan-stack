// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package timeutil_test

import (
	"testing"
	"time"

	. "github.com/TheThingsNetwork/ttn/pkg/util/timeutil"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

var (
	epoch          = time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)
	leap1          = time.Date(1981, time.June, 30, 23, 59, 59, 0, time.UTC).Unix() - epoch.Unix() + 1
	leap5          = time.Date(1987, time.December, 31, 23, 59, 59, 0, time.UTC).Unix() - epoch.Unix() + 1 + 4
	now            = time.Date(2017, time.October, 24, 23, 53, 30, 0, time.UTC)
	nowLeaps int64 = 18
)

func TestGPSConversion(t *testing.T) {
	t.Logf("Leap 1: %d Leap 5: %d", leap1, leap5)

	for _, tc := range []struct {
		GPS  int64
		Time time.Time
	}{
		{
			// From LoRaWAN 1.1 specification
			1139322288,
			time.Date(2016, time.February, 12, 14, 24, 31, 0, time.UTC),
		},
		{
			now.Unix() - epoch.Unix() + nowLeaps,
			now,
		},
		{
			42,
			epoch.Add(42 * time.Second),
		},
		{
			0,
			epoch,
		},
		{
			-1,
			epoch.Add(-1 * time.Second),
		},

		{
			leap1 - 2,
			epoch.Add(time.Second * time.Duration(leap1-2)),
		},
		{
			leap1 - 1,
			epoch.Add(time.Second * time.Duration(leap1-1)),
		},
		{
			leap1,
			epoch.Add(time.Second * time.Duration(leap1)),
		},
		{
			leap1 + 1,
			epoch.Add(time.Second * time.Duration(leap1)),
		},
		{
			leap1 + 2,
			epoch.Add(time.Second * time.Duration(leap1+1)),
		},

		{
			leap5 - 2,
			epoch.Add(time.Second * time.Duration(leap5-6)),
		},
		{
			leap5 - 1,
			epoch.Add(time.Second * time.Duration(leap5-5)),
		},
		{
			leap5,
			epoch.Add(time.Second * time.Duration(leap5-4)),
		},
		{
			leap5 + 1,
			epoch.Add(time.Second * time.Duration(leap5-4)),
		},
		{
			leap5 + 2,
			epoch.Add(time.Second * time.Duration(leap5-3)),
		},
	} {
		a := assertions.New(t)
		a.So(GPS(tc.GPS).UnixNano(), should.Resemble, tc.Time.UnixNano())
		if IsGPSLeap(tc.GPS) {
			a.So(TimeToGPS(tc.Time), should.Equal, tc.GPS+1)
		} else {
			a.So(TimeToGPS(tc.Time), should.Equal, tc.GPS)
		}
		if a.Failed() {
			t.Errorf("Time: %s, Unix: %d, GPS: %d", tc.Time, tc.Time.Unix(), tc.GPS)
		}
	}
}
