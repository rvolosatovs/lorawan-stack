// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package udp

import (
	"testing"

	"github.com/TheThingsNetwork/ttn/pkg/types"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestDataRate(t *testing.T) {
	a := assertions.New(t)

	table := map[string]DataRate{
		`"SF7BW125"`: {types.DataRate{LoRa: "SF7BW125"}},
		`50000`:      {types.DataRate{FSK: 50000}},
	}

	for s, dr := range table {
		enc, err := dr.MarshalJSON()
		a.So(err, should.BeNil)
		a.So(string(enc), should.Equal, s)

		var dec DataRate
		err = dec.UnmarshalJSON(enc)
		a.So(err, should.BeNil)
		a.So(dec, should.Resemble, dr)
	}

	var dr DataRate
	err := dr.UnmarshalJSON([]byte{})
	a.So(err, should.NotBeNil)
}
