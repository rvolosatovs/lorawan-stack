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

package cluster

import (
	"testing"

	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"google.golang.org/grpc"
)

func TestPeer(t *testing.T) {
	a := assertions.New(t)

	conn := new(grpc.ClientConn)

	p := &peer{
		name:   "name",
		roles:  []ttnpb.PeerInfo_Role{ttnpb.PeerInfo_IDENTITY_SERVER},
		tags:   []string{"tag"},
		target: "target",
		conn:   conn,
	}

	a.So(p.HasRole(ttnpb.PeerInfo_APPLICATION_SERVER), should.BeFalse)
	a.So(p.HasRole(ttnpb.PeerInfo_IDENTITY_SERVER), should.BeTrue)

	a.So(p.HasTag("no-tag"), should.BeFalse)
	a.So(p.HasTag("tag"), should.BeTrue)

	a.So(p.Conn(), should.Equal, conn)
}
