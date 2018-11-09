// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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

package assertions_test

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	. "go.thethings.network/lorawan-stack/pkg/util/test/assertions"
)

func TestResembleDiff(t *testing.T) {
	for _, tc := range []struct {
		A         interface{}
		B         interface{}
		Assertion func(interface{}, ...interface{}) string
	}{
		{
			A:         "test",
			B:         "test",
			Assertion: should.BeEmpty,
		},
		{
			A:         "test",
			B:         "test1",
			Assertion: should.NotBeEmpty,
		},
		{
			A:         42,
			B:         42,
			Assertion: should.BeEmpty,
		},
		{
			A:         1,
			B:         2,
			Assertion: should.NotBeEmpty,
		},
		{
			A: struct {
				Foo int
				Bar int
			}{42, 43},
			B: struct {
				Foo int
				Bar int
			}{42, 43},
			Assertion: should.BeEmpty,
		},
		{
			A:         nil,
			B:         0,
			Assertion: should.NotBeEmpty,
		},
		{
			A:         nil,
			B:         "test",
			Assertion: should.NotBeEmpty,
		},
		{
			A:         []string{},
			B:         []string(nil),
			Assertion: should.BeEmpty,
		},
		{
			A:         map[int]int{},
			B:         map[int]int(nil),
			Assertion: should.BeEmpty,
		},
	} {
		t.Run(fmt.Sprintf("%v/%v", tc.A, tc.B), func(t *testing.T) {
			assertions.New(t).So(ShouldResembleDiff(tc.A, tc.B), tc.Assertion)
		})
	}
}
