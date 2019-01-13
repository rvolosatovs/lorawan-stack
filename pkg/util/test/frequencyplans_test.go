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

package test_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func TestFrequencyPlans(t *testing.T) {
	a := assertions.New(t)

	fp := frequencyplans.NewStore(test.FrequencyPlansFetcher)

	euFP, err := fp.GetByID(test.EUFrequencyPlanID)
	a.So(err, should.BeNil)
	a.So(euFP.BandID, should.Equal, "EU_863_870")

	krFP, err := fp.GetByID(test.KRFrequencyPlanID)
	a.So(err, should.BeNil)
	a.So(krFP.UplinkChannels[0].Frequency, should.Equal, 922100000)

	_, err = fp.GetByID(test.ExampleFrequencyPlanID)
	a.So(err, should.BeNil)
}
