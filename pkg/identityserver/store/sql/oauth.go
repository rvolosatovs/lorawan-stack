// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package sql

import (
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/db"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/types"
)

// OAuthStore implements store.OAuthStore.
type OAuthStore struct {
	storer
}

// NewOAuthStore creates a new OAuth store.
func NewOAuthStore(store storer) *OAuthStore {
	return &OAuthStore{
		storer: store,
	}
}

// SaveAuthorizationCode saves the authorization code.
func (s *OAuthStore) SaveAuthorizationCode(authorization *types.AuthorizationData) error {
	return s.saveAuthorizationCode(s.queryer(), authorization)
}

func (s *OAuthStore) saveAuthorizationCode(q db.QueryContext, data *types.AuthorizationData) error {
	_, err := q.NamedExec(
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
			VALUES (
				:authorization_code,
				:client_id,
				:created_at,
				:expires_in,
				:scope,
				:redirect_uri,
				:state,
				:user_id
			)
		`,
		data,
	)

	if _, dup := db.IsDuplicate(err); dup {
		return ErrAuthorizationCodeConflict.New(nil)
	}

	return err
}

// GetAuthorizationCode finds the authorization code.
func (s *OAuthStore) GetAuthorizationCode(authorizationCode string) (*types.AuthorizationData, error) {
	return s.getAuthorizationCode(s.queryer(), authorizationCode)
}

func (s *OAuthStore) getAuthorizationCode(q db.QueryContext, authorizationCode string) (*types.AuthorizationData, error) {
	result := new(types.AuthorizationData)
	err := q.SelectOne(
		result,
		`SELECT *
			FROM authorization_codes
			WHERE authorization_code = $1`,
		authorizationCode,
	)

	if db.IsNoRows(err) {
		return nil, ErrAuthorizationCodeNotFound.New(nil)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
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

// SaveAccessToken saves the access data.
func (s *OAuthStore) SaveAccessToken(access *types.AccessData) error {
	return s.saveAccessToken(s.queryer(), access)
}

func (s *OAuthStore) saveAccessToken(q db.QueryContext, access *types.AccessData) error {
	result := new(string)
	err := q.NamedSelectOne(
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
			VALUES (
				:access_token,
				:client_id,
				:user_id,
				:created_at,
				:expires_in,
				:scope,
				:redirect_uri
			)
			RETURNING access_token`,
		access,
	)

	if _, dup := db.IsDuplicate(err); dup {
		return ErrAccessTokenConflict.New(nil)
	}

	return err
}

// GetRefreshToken finds the access token.
func (s *OAuthStore) GetAccessToken(accessToken string) (*types.AccessData, error) {
	return s.getAccessToken(s.queryer(), accessToken)
}

func (s *OAuthStore) getAccessToken(q db.QueryContext, accessToken string) (*types.AccessData, error) {
	result := new(types.AccessData)
	err := q.SelectOne(
		result,
		`SELECT *
			FROM access_tokens
			WHERE access_token = $1`,
		accessToken,
	)

	if db.IsNoRows(err) {
		return nil, ErrAccessTokenNotFound.New(nil)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteRefreshToken deletes the access token from the database.
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

// SaveRefreshToken saves the refresh token.
func (s *OAuthStore) SaveRefreshToken(access *types.RefreshData) error {
	return s.saveRefreshToken(s.queryer(), access)
}

func (s *OAuthStore) saveRefreshToken(q db.QueryContext, refresh *types.RefreshData) error {
	result := new(string)
	err := q.NamedSelectOne(
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
			VALUES (
				:refresh_token,
				:client_id,
				:user_id,
				:created_at,
				:scope,
				:redirect_uri
			)
			RETURNING refresh_token`,
		refresh,
	)

	if _, dup := db.IsDuplicate(err); dup {
		return ErrRefreshTokenConflict.New(nil)
	}

	return err
}

// GetRefreshToken finds the refresh token.
func (s *OAuthStore) GetRefreshToken(refreshToken string) (*types.RefreshData, error) {
	return s.getRefreshToken(s.queryer(), refreshToken)
}

func (s *OAuthStore) getRefreshToken(q db.QueryContext, refreshToken string) (*types.RefreshData, error) {
	result := new(types.RefreshData)
	err := q.SelectOne(
		result,
		`SELECT *
			FROM refresh_tokens
			WHERE refresh_token = $1`,
		refreshToken,
	)

	if db.IsNoRows(err) {
		return nil, ErrRefreshTokenNotFound.New(nil)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
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
