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

package assertions

import "fmt"

const (
	needNonEmptyCollection = "This assertion requires at least 1 comparison value (you provided 0)."
	needFewerValues        = "This assertion allows %d or fewer comparison values (you provided %d)."
)

func need(needed int, expected []interface{}) string {
	if len(expected) != needed {
		return fmt.Sprintf(needExactValues, needed, len(expected))
	}
	return success
}

func atLeast(minimum int, expected []interface{}) string {
	if len(expected) < 1 {
		return needNonEmptyCollection
	}
	return success
}

func atMost(max int, expected []interface{}) string {
	if len(expected) > max {
		return fmt.Sprintf(needFewerValues, max, len(expected))
	}
	return success
}
