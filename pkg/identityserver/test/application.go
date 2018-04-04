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

package test

import (
	"fmt"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
)

func defaultApplication(in interface{}) (*ttnpb.Application, error) {
	if app, ok := in.(store.Application); ok {
		return app.GetApplication(), nil
	}

	if app, ok := in.(ttnpb.Application); ok {
		return &app, nil
	}

	if ptr, ok := in.(*ttnpb.Application); ok {
		return ptr, nil
	}

	return nil, fmt.Errorf("Expected: '%v' to be of type ttnpb.Application but it was not", in)
}

// ShouldBeApplication checks if two Applications resemble each other.
func ShouldBeApplication(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one application to match but got %v", len(expected))
	}

	a, s := defaultApplication(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultApplication(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		ShouldBeApplicationIgnoringAutoFields(a, b),
		assertions.ShouldHappenWithin(a.CreatedAt, time.Millisecond, b.CreatedAt),
	)
}

// ShouldBeApplicationIgnoringAutoFields checks if two Applications resemble each other
// without looking at fields that are generated by the database: created.
func ShouldBeApplicationIgnoringAutoFields(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one application to match but got %v", len(expected))
	}

	a, s := defaultApplication(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultApplication(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		assertions.ShouldEqual(a.ApplicationID, b.ApplicationID),
		assertions.ShouldResemble(a.Description, b.Description),
	)
}
