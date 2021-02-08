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
	"reflect"
	"strings"

	"github.com/kr/pretty"
	"github.com/smartystreets/assertions"
)

const (
	shouldHaveEmptyDiff    = "Expected: '%#v'\nActual:   '%#v'\nDiff:  '%s'\n(should resemble diff)!"
	shouldNotHaveEmptyDiff = "Expected '%#v'\nto diff  '%#v'\n(but it did)!"

	needPointer                        = "This assertion requires a pointer type (you provided %T)."
	needSetFielderCompatible           = "This assertion requires a SetFielder-compatible comparison type (you provided %T)."
	needStringCompatible               = "This assertion requires a string-compatible comparison type (you provided %T)."
	needStringCompatibleOrArrayOrSlice = "This assertion requires a string-compatible comparison type or a either array or slice of such(you provided %T)."
	setFieldsFailed                    = "SetFields failed with: %s"
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

	diff := pretty.Diff(actual, expected[0])
	if len(diff) == 0 {
		return message
	}

	lines := make([]string, 1, len(diff)+2)
	lines[0] = "Diff:"
	for _, d := range diff {
		lines = append(lines, fmt.Sprintf("   %s", d))
	}
	lines = append(lines, fmt.Sprintf("Actual: %s", pretty.Sprint(actual)))
	lines = append(lines, fmt.Sprintf("Expected: %s", pretty.Sprint(expected[0])))
	return strings.Join(append(lines, lastLine(message)), "\n")
}

// ShouldResembleFields is same as ShouldResemble, but only compares the specified fields for 2 given SetFielders.
func ShouldResembleFields(actual interface{}, expected ...interface{}) (message string) {
	if len(expected) < 1 {
		return fmt.Sprintf(needAtLeastValues, 1, len(expected))
	}

	ps := reflect.MakeSlice(reflect.TypeOf([]string{}), 0, 0)
	for _, p := range expected[1:] {
		pv := reflect.ValueOf(p)
		switch pv.Kind() {
		case reflect.String:
			ps = reflect.Append(ps, pv)
		case reflect.Array:
			if pv.Type().Elem().Kind() != reflect.String {
				return fmt.Sprintf(needStringCompatible, p)
			}
			for i := 0; i < pv.Len(); i++ {
				ps = reflect.Append(ps, pv.Index(i))
			}
		case reflect.Slice:
			if pv.Type().Elem().Kind() != reflect.String {
				return fmt.Sprintf(needStringCompatible, p)
			}
			ps = reflect.AppendSlice(ps, pv)
		default:
			return fmt.Sprintf(needStringCompatibleOrArrayOrSlice, p)
		}
	}
	selectFields := func(v interface{}) (interface{}, string) {
		t := reflect.TypeOf(v)
		if t.Kind() != reflect.Ptr {
			return nil, fmt.Sprintf(needPointer, v)
		}
		sv := reflect.New(t.Elem())
		if v == nil {
			return sv.Interface(), ""
		}
		m := sv.MethodByName("SetFields")
		if !m.IsValid() {
			return nil, fmt.Sprintf(needSetFielderCompatible, v)
		}
		if ret := m.CallSlice([]reflect.Value{reflect.ValueOf(v), ps})[0]; !ret.IsNil() {
			return nil, fmt.Sprintf(setFieldsFailed, ret.Interface().(error))
		}
		return sv.Interface(), ""
	}

	actualCmp, errStr := selectFields(actual)
	if errStr != "" {
		return errStr
	}
	expectedCmp, errStr := selectFields(expected[0])
	if errStr != "" {
		return errStr
	}
	return ShouldResemble(actualCmp, expectedCmp)
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
