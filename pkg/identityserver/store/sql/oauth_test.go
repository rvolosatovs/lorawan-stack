// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package sql

import (
	"testing"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/test"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestAuthorizationCode(t *testing.T) {
	a := assertions.New(t)

	s := cleanStore(t, database)

	client := testClients()["test-client"]
	err := s.Clients.Create(client)
	a.So(err, should.BeNil)

	data := &store.AuthorizationData{
		Code:        "123456",
		ClientID:    client.ClientIdentifier.ClientID,
		CreatedAt:   time.Now(),
		ExpiresIn:   5 * time.Second,
		Scope:       "scope",
		RedirectURI: "https://example.com/oauth/callback",
		State:       "state",
	}

	err = s.OAuth.SaveAuthorizationCode(data)
	a.So(err, should.BeNil)

	found, c, err := s.OAuth.FindAuthorizationCode(data.Code)
	a.So(err, should.BeNil)

	a.So(found.Code, should.Equal, data.Code)
	a.So(found.ClientID, should.Equal, data.ClientID)
	a.So(found.ExpiresIn, should.Equal, data.ExpiresIn)
	a.So(found.Scope, should.Equal, data.Scope)
	a.So(found.RedirectURI, should.Equal, data.RedirectURI)
	a.So(found.State, should.Equal, data.State)

	a.So(c, test.ShouldBeClientIgnoringAutoFields, client)

	err = s.OAuth.DeleteAuthorizationCode(data.Code)
	a.So(err, should.BeNil)

	_, _, err = s.OAuth.FindAuthorizationCode(data.Code)
	a.So(err, should.NotBeNil)
}

func TestRefreshToken(t *testing.T) {
	a := assertions.New(t)

	s := cleanStore(t, database)

	client := testClients()["test-client"]
	err := s.Clients.Create(client)
	a.So(err, should.BeNil)

	data := &store.RefreshData{
		RefreshToken: "123456",
		ClientID:     client.ClientIdentifier.ClientID,
		CreatedAt:    time.Now(),
		Scope:        "scope",
		RedirectURI:  "https://example.com/oauth/callback",
	}

	err = s.OAuth.SaveRefreshToken(data)
	a.So(err, should.BeNil)

	found, c, err := s.OAuth.FindRefreshToken(data.RefreshToken)
	a.So(err, should.BeNil)

	a.So(found.RefreshToken, should.Equal, data.RefreshToken)
	a.So(found.ClientID, should.Equal, data.ClientID)
	a.So(found.Scope, should.Equal, data.Scope)
	a.So(found.RedirectURI, should.Equal, data.RedirectURI)

	a.So(c, test.ShouldBeClientIgnoringAutoFields, client)

	err = s.OAuth.DeleteRefreshToken(data.RefreshToken)
	a.So(err, should.BeNil)

	_, _, err = s.OAuth.FindRefreshToken(data.RefreshToken)
	a.So(err, should.NotBeNil)
	_ = c
}
