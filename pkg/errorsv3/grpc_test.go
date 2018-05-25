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

package errors_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	errors "go.thethings.network/lorawan-stack/pkg/errorsv3"
	_ "go.thethings.network/lorawan-stack/pkg/ttnpb"
)

func TestGRPCConversion(t *testing.T) {
	a := assertions.New(t)

	errDef := errors.Define("test_grpc_conversion_err_def", "gRPC Conversion Error")
	a.So(errors.FromGRPCStatus(errDef.GRPCStatus()).Definition, errors.ShouldEqual, errDef)

	errHello := errors.New("hello world").WithAttributes("foo", "bar")
	a.So(errors.FromGRPCStatus(errHello.GRPCStatus()), errors.ShouldEqual, errHello)
}
