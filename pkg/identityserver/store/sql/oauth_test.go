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

package sql_test

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	. "go.thethings.network/lorawan-stack/pkg/identityserver/store/sql"
	"go.thethings.network/lorawan-stack/pkg/identityserver/test"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

func TestOAuthAuthorizationCode(t *testing.T) {
	a := assertions.New(t)
	s := testStore(t, database)

	err := s.Clients.Create(client)
	a.So(err, should.BeNil)

	data := store.AuthorizationData{
		AuthorizationCode: "123456",
		ClientID:          client.ClientIdentifiers.ClientID,
		CreatedAt:         time.Now(),
		ExpiresIn:         5 * time.Second,
		Scope:             "scope",
		RedirectURI:       "https://example.com/oauth/callback",
		State:             "state",
		UserID:            alice.UserIdentifiers.UserID,
	}

	err = s.OAuth.SaveAuthorizationCode(data)
	a.So(err, should.BeNil)

	found, err := s.OAuth.GetAuthorizationCode(data.AuthorizationCode)
	a.So(err, should.BeNil)

	a.So(found.AuthorizationCode, should.Equal, data.AuthorizationCode)
	a.So(found.ClientID, should.Equal, data.ClientID)
	a.So(found.UserID, should.Equal, data.UserID)
	a.So(found.ExpiresIn, should.Equal, data.ExpiresIn)
	a.So(found.Scope, should.Equal, data.Scope)
	a.So(found.RedirectURI, should.Equal, data.RedirectURI)
	a.So(found.State, should.Equal, data.State)

	c, err := s.Clients.GetByID(ttnpb.ClientIdentifiers{ClientID: found.ClientID}, clientSpecializer)
	a.So(err, should.BeNil)
	a.So(c, test.ShouldBeClientIgnoringAutoFields, client)

	err = s.OAuth.DeleteAuthorizationCode(data.AuthorizationCode)
	a.So(err, should.BeNil)

	_, err = s.OAuth.GetAuthorizationCode(data.AuthorizationCode)
	a.So(err, should.NotBeNil)
	a.So(ErrAuthorizationCodeNotFound.Describes(err), should.BeTrue)
}

func TestOAuthAccessToken(t *testing.T) {
	a := assertions.New(t)
	s := testStore(t, database)

	data := store.AccessData{
		AccessToken: "123456",
		ClientID:    client.ClientIdentifiers.ClientID,
		UserID:      alice.UserIdentifiers.UserID,
		CreatedAt:   time.Now(),
		ExpiresIn:   time.Hour,
		Scope:       "scope",
		RedirectURI: "https://example.com/oauth/callback",
	}

	err := s.OAuth.SaveAccessToken(data)
	a.So(err, should.BeNil)

	found, err := s.OAuth.GetAccessToken(data.AccessToken)
	a.So(err, should.BeNil)

	a.So(found.AccessToken, should.Equal, data.AccessToken)
	a.So(found.ClientID, should.Equal, data.ClientID)
	a.So(found.UserID, should.Equal, data.UserID)
	a.So(found.Scope, should.Equal, data.Scope)
	a.So(found.RedirectURI, should.Equal, data.RedirectURI)
	a.So(found.ExpiresIn, should.Equal, data.ExpiresIn)

	c, err := s.Clients.GetByID(ttnpb.ClientIdentifiers{ClientID: found.ClientID}, clientSpecializer)
	a.So(err, should.BeNil)
	a.So(c, test.ShouldBeClientIgnoringAutoFields, client)

	err = s.OAuth.DeleteAccessToken(data.AccessToken)
	a.So(err, should.BeNil)

	_, err = s.OAuth.GetAccessToken(data.AccessToken)
	a.So(err, should.NotBeNil)
	a.So(ErrAccessTokenNotFound.Describes(err), should.BeTrue)
}

func TestOAuthRefreshToken(t *testing.T) {
	a := assertions.New(t)
	s := testStore(t, database)

	data := store.RefreshData{
		RefreshToken: "123456",
		ClientID:     client.ClientIdentifiers.ClientID,
		UserID:       alice.UserIdentifiers.UserID,
		CreatedAt:    time.Now(),
		Scope:        "scope",
		RedirectURI:  "https://example.com/oauth/callback",
	}

	err := s.OAuth.SaveRefreshToken(data)
	a.So(err, should.BeNil)

	found, err := s.OAuth.GetRefreshToken(data.RefreshToken)
	a.So(err, should.BeNil)

	a.So(found.RefreshToken, should.Equal, data.RefreshToken)
	a.So(found.ClientID, should.Equal, data.ClientID)
	a.So(found.UserID, should.Equal, data.UserID)
	a.So(found.Scope, should.Equal, data.Scope)
	a.So(found.RedirectURI, should.Equal, data.RedirectURI)

	c, err := s.Clients.GetByID(ttnpb.ClientIdentifiers{ClientID: found.ClientID}, clientSpecializer)
	a.So(err, should.BeNil)
	a.So(c, test.ShouldBeClientIgnoringAutoFields, client)

	err = s.OAuth.DeleteRefreshToken(data.RefreshToken)
	a.So(err, should.BeNil)

	_, err = s.OAuth.GetRefreshToken(data.RefreshToken)
	a.So(err, should.NotBeNil)
	a.So(ErrRefreshTokenNotFound.Describes(err), should.BeTrue)
}

func TestOAuthAuthorizedClients(t *testing.T) {
	a := assertions.New(t)
	s := testStore(t, database)

	accessData := store.AccessData{
		AccessToken: "123456",
		ClientID:    client.ClientIdentifiers.ClientID,
		UserID:      alice.UserIdentifiers.UserID,
		CreatedAt:   time.Now(),
		ExpiresIn:   time.Hour,
		Scope:       "scope",
		RedirectURI: "https://example.com/oauth/callback",
	}

	err := s.OAuth.SaveAccessToken(accessData)
	a.So(err, should.BeNil)

	refreshData := store.RefreshData{
		RefreshToken: "123456",
		ClientID:     client.ClientIdentifiers.ClientID,
		UserID:       alice.UserIdentifiers.UserID,
		CreatedAt:    time.Now(),
		Scope:        "scope",
		RedirectURI:  "https://example.com/oauth/callback",
	}

	err = s.OAuth.SaveRefreshToken(refreshData)
	a.So(err, should.BeNil)

	found, err := s.OAuth.ListAuthorizedClients(alice.UserIdentifiers, clientSpecializer)
	a.So(err, should.BeNil)
	if a.So(found, should.HaveLength, 1) {
		a.So(found[0], test.ShouldBeClientIgnoringAutoFields, client)
	}

	err = s.OAuth.RevokeAuthorizedClient(alice.UserIdentifiers, client.ClientIdentifiers)
	a.So(err, should.BeNil)

	err = s.OAuth.RevokeAuthorizedClient(alice.UserIdentifiers, client.ClientIdentifiers)
	a.So(err, should.NotBeNil)
	a.So(ErrAuthorizedClientNotFound.Describes(err), should.BeTrue)
}
