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

package assertions

import (
	"fmt"
	"strings"

	"github.com/kr/pretty"
	"github.com/smartystreets/assertions"
)

const (
	shouldHaveEmptyDiff    = "Expected: '%#v'\nActual:   '%#v'\nDiff:  '%s'\n(should resemble diff)!"
	shouldNotHaveEmptyDiff = "Expected '%#v'\nto diff  '%#v'\n(but it did)!"
)

func lastLine(s string) string {
	if s == "" {
		return ""
	}

	ls := strings.Split(s, "\n")
	return ls[len(ls)-1]
}

// ShouldResemble wraps assertions.ShouldResemble and prepends a diff if assertion fails.
func ShouldResemble(actual interface{}, expected ...interface{}) (message string) {
	if message = assertions.ShouldResemble(actual, expected...); message == success {
		return success
	}

	diff := pretty.Diff(expected[0], actual)
	if len(diff) == 0 {
		return message
	}

	lines := make([]string, 1, len(diff)+2)
	lines[0] = "Diff:"
	for _, d := range diff {
		lines = append(lines, fmt.Sprintf("   %s", d))
	}
	return strings.Join(append(lines, lastLine(message)), "\n")
}

// ShouldHaveEmptyDiff compares the pretty.Diff of values.
func ShouldHaveEmptyDiff(actual interface{}, expected ...interface{}) (message string) {
	if message = need(1, expected); message != success {
		return
	}
	diff := pretty.Diff(expected[0], actual)
	if len(diff) != 0 {
		return fmt.Sprintf(shouldHaveEmptyDiff, expected[0], actual, diff)
	}
	return success
}

// ShouldNotHaveEmptyDiff compares the pretty.Diff of values.
func ShouldNotHaveEmptyDiff(actual interface{}, expected ...interface{}) (message string) {
	if message = need(1, expected); message != success {
		return
	}
	diff := pretty.Diff(expected[0], actual)
	if len(diff) == 0 {
		return fmt.Sprintf(shouldNotHaveEmptyDiff, expected[0], actual)
	}
	return success
}
