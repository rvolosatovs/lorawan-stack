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
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/satori/go.uuid"
)

// OAuthStore implements store.OAuthStore.
type OAuthStore struct {
	storer

	*UserStore
	*ClientStore
}

// NewOAuthStore creates a new OAuth store.
func NewOAuthStore(store storer) *OAuthStore {
	return &OAuthStore{
		storer:      store,
		UserStore:   store.store().Users.(*UserStore),
		ClientStore: store.store().Clients.(*ClientStore),
	}
}

// SaveAuthorizationCode saves the authorization code.
func (s *OAuthStore) SaveAuthorizationCode(data store.AuthorizationData) error {
	err := s.transact(func(tx *db.Tx) error {
		userID, err := s.getUserID(tx, ttnpb.UserIdentifiers{UserID: data.UserID})
		if err != nil {
			return err
		}

		clientID, err := s.getClientID(tx, ttnpb.ClientIdentifiers{ClientID: data.ClientID})
		if err != nil {
			return err
		}

		return s.saveAuthorizationCode(tx, userID, clientID, data)
	})
	return err
}

func (s *OAuthStore) saveAuthorizationCode(q db.QueryContext, userID, clientID uuid.UUID, data store.AuthorizationData) error {
	_, err := q.Exec(
		`INSERT
			INTO authorization_codes (
				authorization_code,
				client_id,
				created_at,
				expires_in,
				scope,
				redirect_uri,
				state,
				user_id
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
		data.AuthorizationCode,
		clientID,
		data.CreatedAt,
		data.ExpiresIn,
		data.Scope,
		data.RedirectURI,
		data.State,
		userID,
	)

	if _, dup := db.IsDuplicate(err); dup {
		return ErrAuthorizationCodeConflict.New(nil)
	}

	return err
}

type authorizationData struct {
	ClientUUID uuid.UUID
	UserUUID   uuid.UUID
	store.AuthorizationData
}

// GetAuthorizationCode finds the authorization code.
func (s *OAuthStore) GetAuthorizationCode(authorizationCode string) (result store.AuthorizationData, err error) {
	err = s.transact(func(tx *db.Tx) error {
		data, err := s.getAuthorizationCode(tx, authorizationCode)
		if err != nil {
			return err
		}

		user, err := s.getUserIdentifiersFromID(tx, data.UserUUID)
		if err != nil {
			return err
		}
		data.AuthorizationData.UserID = user.UserID

		client, err := s.getClientIdentifiersFromID(tx, data.ClientUUID)
		if err != nil {
			return err
		}
		data.AuthorizationData.ClientID = client.ClientID

		result = data.AuthorizationData

		return nil
	})
	return
}

func (s *OAuthStore) getAuthorizationCode(q db.QueryContext, authorizationCode string) (data authorizationData, err error) {
	err = q.SelectOne(
		&data,
		`SELECT
				authorization_code,
				client_id AS client_uuid,
				created_at,
				expires_in,
				scope,
				redirect_uri,
				state,
				user_id AS user_uuid
			FROM authorization_codes
			WHERE authorization_code = $1`,
		authorizationCode)

	if db.IsNoRows(err) {
		err = ErrAuthorizationCodeNotFound.New(nil)
	}

	return
}

// DeleteAuthorizationCode deletes the authorization code.
func (s *OAuthStore) DeleteAuthorizationCode(authorizationCode string) error {
	return s.deleteAuthorizationCode(s.queryer(), authorizationCode)
}

func (s *OAuthStore) deleteAuthorizationCode(q db.QueryContext, authorizationCode string) error {
	code := new(string)
	err := q.SelectOne(
		code,
		`DELETE
			FROM authorization_codes
			WHERE authorization_code = $1
			RETURNING authorization_code`,
		authorizationCode,
	)

	if db.IsNoRows(err) {
		return ErrAuthorizationCodeNotFound.New(nil)
	}

	return err
}

func (s *OAuthStore) deleteAuthorizationCodesByUser(q db.QueryContext, userID uuid.UUID) error {
	_, err := q.Exec(`DELETE FROM authorization_codes WHERE user_id = $1`, userID)
	return err
}

type accessData struct {
	ClientUUID uuid.UUID
	UserUUID   uuid.UUID
	store.AccessData
}

// SaveAccessToken saves the access data.
func (s *OAuthStore) SaveAccessToken(data store.AccessData) error {
	err := s.transact(func(tx *db.Tx) error {
		userID, err := s.getUserID(tx, ttnpb.UserIdentifiers{UserID: data.UserID})
		if err != nil {
			return err
		}

		clientID, err := s.getClientID(tx, ttnpb.ClientIdentifiers{ClientID: data.ClientID})
		if err != nil {
			return err
		}

		return s.saveAccessToken(tx, userID, clientID, data)
	})
	return err
}

func (s *OAuthStore) saveAccessToken(q db.QueryContext, userID, clientID uuid.UUID, access store.AccessData) error {
	result := new(string)
	err := q.SelectOne(
		result,
		`INSERT
			INTO access_tokens (
				access_token,
				client_id,
				user_id,
				created_at,
				expires_in,
				scope,
				redirect_uri
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING access_token`,
		access.AccessToken,
		clientID,
		userID,
		access.CreatedAt,
		access.ExpiresIn,
		access.Scope,
		access.RedirectURI,
	)

	if _, dup := db.IsDuplicate(err); dup {
		return ErrAccessTokenConflict.New(nil)
	}

	return err
}

// GetAccessToken finds the access token.
func (s *OAuthStore) GetAccessToken(accessToken string) (result store.AccessData, err error) {
	err = s.transact(func(tx *db.Tx) error {
		data, err := s.getAccessToken(tx, accessToken)
		if err != nil {
			return err
		}

		user, err := s.getUserIdentifiersFromID(tx, data.UserUUID)
		if err != nil {
			return err
		}
		data.AccessData.UserID = user.UserID

		client, err := s.getClientIdentifiersFromID(tx, data.ClientUUID)
		if err != nil {
			return err
		}
		data.AccessData.ClientID = client.ClientID

		result = data.AccessData

		return nil
	})
	return
}

func (s *OAuthStore) getAccessToken(q db.QueryContext, accessToken string) (result accessData, err error) {
	err = q.SelectOne(
		&result,
		`SELECT
				access_token,
				client_id AS client_uuid,
				user_id AS user_uuid,
				created_at,
				expires_in,
				scope,
				redirect_uri
			FROM access_tokens
			WHERE access_token = $1`,
		accessToken,
	)

	if db.IsNoRows(err) {
		err = ErrAccessTokenNotFound.New(nil)
	}

	return
}

// DeleteAccessToken deletes the access token from the database.
func (s *OAuthStore) DeleteAccessToken(accessToken string) error {
	return s.deleteAccessToken(s.queryer(), accessToken)
}

func (s *OAuthStore) deleteAccessToken(q db.QueryContext, accessToken string) error {
	token := new(string)
	err := q.SelectOne(
		token,
		`DELETE
			FROM access_tokens
			WHERE access_token = $1
			RETURNING access_token`,
		accessToken,
	)

	if db.IsNoRows(err) {
		return ErrAccessTokenNotFound.New(nil)
	}

	return err
}

type refreshData struct {
	ClientUUID uuid.UUID
	UserUUID   uuid.UUID
	store.RefreshData
}

// SaveRefreshToken saves the refresh token.
func (s *OAuthStore) SaveRefreshToken(data store.RefreshData) error {
	err := s.transact(func(tx *db.Tx) error {
		userID, err := s.getUserID(tx, ttnpb.UserIdentifiers{UserID: data.UserID})
		if err != nil {
			return err
		}

		clientID, err := s.getClientID(tx, ttnpb.ClientIdentifiers{ClientID: data.ClientID})
		if err != nil {
			return err
		}

		return s.saveRefreshToken(tx, userID, clientID, data)
	})
	return err
}

func (s *OAuthStore) saveRefreshToken(q db.QueryContext, userID, clientID uuid.UUID, data store.RefreshData) error {
	result := new(string)
	err := q.SelectOne(
		result,
		`INSERT
			INTO refresh_tokens (
				refresh_token,
				client_id,
				user_id,
				created_at,
				scope,
				redirect_uri
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING refresh_token`,
		data.RefreshToken,
		clientID,
		userID,
		data.CreatedAt,
		data.Scope,
		data.RedirectURI,
	)

	if _, dup := db.IsDuplicate(err); dup {
		return ErrRefreshTokenConflict.New(nil)
	}

	return err
}

// GetRefreshToken finds the refresh token.
func (s *OAuthStore) GetRefreshToken(refreshToken string) (result store.RefreshData, err error) {
	err = s.transact(func(tx *db.Tx) error {
		data, err := s.getRefreshToken(tx, refreshToken)
		if err != nil {
			return err
		}

		user, err := s.getUserIdentifiersFromID(tx, data.UserUUID)
		if err != nil {
			return err
		}
		data.RefreshData.UserID = user.UserID

		client, err := s.getClientIdentifiersFromID(tx, data.ClientUUID)
		if err != nil {
			return err
		}
		data.RefreshData.ClientID = client.ClientID

		result = data.RefreshData

		return nil
	})
	return
}

func (s *OAuthStore) getRefreshToken(q db.QueryContext, refreshToken string) (result refreshData, err error) {
	err = q.SelectOne(
		&result,
		`SELECT
				refresh_token,
				client_id AS client_uuid,
				user_id AS user_uuid,
				created_at,
				scope,
				redirect_uri
			FROM refresh_tokens
			WHERE refresh_token = $1`,
		refreshToken,
	)

	if db.IsNoRows(err) {
		err = ErrRefreshTokenNotFound.New(nil)
	}

	return
}

// DeleteRefreshToken deletes the refresh token from the database.
func (s *OAuthStore) DeleteRefreshToken(refreshToken string) error {
	return s.deleteRefreshToken(s.queryer(), refreshToken)
}

func (s *OAuthStore) deleteRefreshToken(q db.QueryContext, refreshToken string) error {
	token := new(string)
	err := q.SelectOne(
		token,
		`DELETE
			FROM refresh_tokens
			WHERE refresh_token = $1
			RETURNING refresh_token`,
		refreshToken,
	)

	if db.IsNoRows(err) {
		return ErrRefreshTokenNotFound.New(nil)
	}

	return err
}

// ListAuthorizedClients returns a list of clients authorized by a given user.
func (s *OAuthStore) ListAuthorizedClients(ids ttnpb.UserIdentifiers, specializer store.ClientSpecializer) (result []store.Client, err error) {
	err = s.transact(func(tx *db.Tx) error {
		userID, err := s.getUserID(tx, ids)
		if err != nil {
			return err
		}

		clientIDs, err := s.listAuthorizedClients(tx, userID)
		if err != nil {
			return err
		}

		for _, clientID := range clientIDs {
			client, err := s.ClientStore.getByID(tx, clientID)
			if err != nil {
				return err
			}

			specialized := specializer(client)

			err = s.ClientStore.loadAttributes(tx, clientID, specialized)
			if err != nil {
				return err
			}

			result = append(result, specialized)
		}

		return nil
	})
	return
}

func (s *OAuthStore) listAuthorizedClients(q db.QueryContext, userID uuid.UUID) (ids []uuid.UUID, err error) {
	err = q.Select(
		&ids,
		`SELECT DISTINCT clients.id
			FROM clients
			JOIN refresh_tokens
			ON (
				clients.id = refresh_tokens.client_id AND refresh_tokens.user_id = $1
			)
			JOIN access_tokens
			ON (
				clients.id = access_tokens.client_id AND access_tokens.user_id = $1
			)`,
		userID)
	return
}

// RevokeAuthorizedClient deletes the access tokens and refresh token
// granted to a client by a given user.
func (s *OAuthStore) RevokeAuthorizedClient(userIDs ttnpb.UserIdentifiers, clientIDs ttnpb.ClientIdentifiers) error {
	rows := 0
	err := s.transact(func(tx *db.Tx) error {
		userID, err := s.getUserID(tx, userIDs)
		if err != nil {
			return err
		}

		clientID, err := s.getClientID(tx, clientIDs)
		if err != nil {
			return err
		}

		rowsa, err := s.deleteAccessTokensByUserAndClient(tx, userID, clientID)
		if err != nil {
			return err
		}

		rowsr, err := s.deleteRefreshTokenByUserAndClient(tx, userID, clientID)
		if err != nil {
			return err
		}

		rows = rowsa + rowsr

		return nil
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAuthorizedClientNotFound.New(nil)
	}
	return nil
}

func (s *OAuthStore) deleteAccessTokensByUserAndClient(q db.QueryContext, userID, clientID uuid.UUID) (int, error) {
	res, err := q.Exec(
		`DELETE
			FROM access_tokens
			WHERE user_id = $1 AND client_id = $2`,
		userID,
		clientID)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rows), nil
}

func (s *OAuthStore) deleteRefreshTokenByUserAndClient(q db.QueryContext, userID, clientID uuid.UUID) (int, error) {
	res, err := q.Exec(
		`DELETE
			FROM refresh_tokens
			WHERE user_id = $1 AND client_id = $2`,
		userID,
		clientID)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rows), nil
}
