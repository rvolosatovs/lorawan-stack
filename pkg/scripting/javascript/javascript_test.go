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

package javascript_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/scripting"
	. "go.thethings.network/lorawan-stack/pkg/scripting/javascript"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func TestRun(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()

	script := `
		(function () {
			return {
				x: 42
			}
		})()
	`

	e := New(scripting.DefaultOptions)
	output, err := e.Run(ctx, script, nil)
	a.So(err, should.BeNil)
	a.So(output, should.HaveSameTypeAs, map[string]interface{}{})
	a.So(output.(map[string]interface{})["x"], should.Equal, 42)
}

func TestRunError(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()

	script := `
		(function () {
			throw Error("something didn't work")
		})()
	`

	e := New(scripting.DefaultOptions)
	_, err := e.Run(ctx, script, nil)
	a.So(err, should.NotBeNil)
}

func TestRunStackOverflow(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()

	script := `
		(function () {
			var obj = {foo: "bar"};
			obj.ob = obj;
			return obj;
		})()
	`

	e := New(scripting.DefaultOptions)
	_, err := e.Run(ctx, script, nil)
	a.So(err, should.NotBeNil)
}

func TestRunTimeout(t *testing.T) {
	a := assertions.New(t)

	ctx := test.Context()

	script := `
		(function () {
			while (true) { }
			return {};
		})()
	`

	e := New(scripting.DefaultOptions)
	_, err := e.Run(ctx, script, nil)
	a.So(err, should.NotBeNil)
	a.So(errors.IsDeadlineExceeded(errors.Cause(err)), should.BeTrue)
}
