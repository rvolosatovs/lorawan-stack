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

// Package gpstime provides utilities to work with GPS time.
package gpstime

import (
	"time"
)

// 1980-01-06T00:00:00+00:00
const gpsEpochSec = 315964800

// Leap seconds in GPS time
var leaps = [...]int64{
	46828800,
	78364801,
	109900802,
	173059203,
	252028804,
	315187205,
	346723206,
	393984007,
	425520008,
	457056009,
	504489610,
	551750411,
	599184012,
	820108813,
	914803214,
	1025136015,
	1119744016,
	1167264017,
	1341118800,
}

// IsLeap reports whether the given GPS time, sec seconds since January 6, 1980 UTC, is a leap second in UTC.
func IsLeap(sec int64) bool {
	i := int64(len(leaps)) - 1
	for ; i >= 0; i-- {
		if sec > leaps[i] {
			return false
		}
		if sec == leaps[i] {
			return true
		}
	}
	return false
}

// Parse returns the local Time corresponding to the given Time time, sec seconds since January 6, 1980 UTC.
func Parse(sec int64) time.Time {
	i := int64(len(leaps))
	for ; i > 0; i-- {
		if sec > leaps[i-1] {
			break
		}
	}
	return time.Unix(sec+gpsEpochSec-i, 0)
}

// ToGPS returns t as a ToGPS time, the number of seconds elapsed since January 6, 1980 UTC.
func ToGPS(t time.Time) int64 {
	sec := t.Unix() - gpsEpochSec

	i := int64(len(leaps))
	for ; i > 0; i-- {
		if sec > leaps[i-1]-i {
			break
		}
	}
	return sec + i
}
