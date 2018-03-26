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

package sql

import (
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/db"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/satori/go.uuid"
)

const (
	// DO NOT modify the following values.
	// These are the internal values used to denote if a registered ID belongs
	// either to an organization or to an user.
	organizationIDType int = 0
	userIDType         int = 1
)

// accountStore is a reusable substore to manage the shared ID namespace between
// organizations and users.
type accountStore struct {
	storer
}

// newAccountStore returns an accountStore.
func newAccountStore(store storer) *accountStore {
	return &accountStore{
		storer: store,
	}
}

func (s *accountStore) getAccountID(q db.QueryContext, ids ttnpb.OrganizationOrUserIdentifiers) (uuid.UUID, error) {
	// TODO(gomezjdaniel#543): avoid dynamic type checking to access other stores.
	if id := ids.GetUserID(); id != nil {
		return s.store().Users.(*UserStore).getUserID(q, *id)
	}
	return s.store().Organizations.(*OrganizationStore).getOrganizationID(q, *ids.GetOrganizationID())
}

// registerOrganizationID registers the given ID that belongs to an organization.
// It returns the autogenerated UUID.
func (s *accountStore) registerOrganizationID(q db.QueryContext, organizationID string) (id uuid.UUID, err error) {
	err = q.SelectOne(
		&id,
		`INSERT
			INTO accounts(account_id, type)
			VALUES ($1, $2)
			RETURNING id`,
		organizationID,
		organizationIDType)
	if _, yes := db.IsDuplicate(err); yes {
		err = ErrOrganizationIDTaken.New(nil)
	}
	return
}

// registerUserID registers the given ID that belongs to an user.
// It returns the autogenerated UUID.
func (s *accountStore) registerUserID(q db.QueryContext, userID string) (id uuid.UUID, err error) {
	err = q.SelectOne(
		&id,
		`INSERT
			INTO accounts(account_id, type)
			VALUES ($1, $2)
			RETURNING id`,
		userID,
		userIDType)
	if _, yes := db.IsDuplicate(err); yes {
		err = ErrUserIDTaken.New(nil)
	}
	return
}

// deleteID deletes the given ID.
func (s *accountStore) deleteID(q db.QueryContext, id uuid.UUID) (err error) {
	var i string
	err = q.SelectOne(
		&i,
		`DELETE
			FROM accounts
			WHERE id = $1
			RETURNING account_id`,
		id)
	if db.IsNoRows(err) {
		err = ErrAccountIDNotFound.New(nil)
	}
	return
}
