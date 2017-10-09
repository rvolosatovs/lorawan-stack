// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package test

import (
	"testing"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func client() *ttnpb.Client {
	return &ttnpb.Client{
		ClientIdentifier: ttnpb.ClientIdentifier{"test-client"},
		Secret:           "123456",
		CallbackURI:      "/oauth/callback",
		Grants:           []ttnpb.ClientGrant{ttnpb.GRANT_AUTHORIZATION_CODE},
	}
}

func TestShouldBeClient(t *testing.T) {
	a := assertions.New(t)

	a.So(ShouldBeClient(client(), client()), should.Equal, success)

	modified := client()
	modified.CreatedAt = time.Now()

	a.So(ShouldBeClient(modified, client()), should.NotEqual, success)
}

func TestShouldBeClientIgnoringAutoFields(t *testing.T) {
	a := assertions.New(t)

	a.So(ShouldBeClientIgnoringAutoFields(client(), client()), should.Equal, success)

	modified := client()
	modified.Secret = "foo"
	modified.Grants = []ttnpb.ClientGrant{}

	a.So(ShouldBeClientIgnoringAutoFields(modified, client()), should.NotEqual, success)
}
