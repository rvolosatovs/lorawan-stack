// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package store

import "github.com/TheThingsNetwork/ttn/pkg/identityserver/types"

// OAuthStore is a store that manages OAuth authorization codes, access tokens
// and refresh tokens.
type OAuthStore interface {
	// SaveAuthorizationCode saves the authorization code.
	SaveAuthorizationCode(authorization *types.AuthorizationData) error

	// GetAuthorizationCode finds the authorization code.
	GetAuthorizationCode(authorizationCode string) (*types.AuthorizationData, error)

	// DeleteAuthorizationCode deletes the authorization code.
	DeleteAuthorizationCode(authorizationCode string) error

	// SaveAccessToken saves the access token.
	SaveAccessToken(access *types.AccessData) error

	// GetAccessToken finds the access token.
	GetAccessToken(accessToken string) (*types.AccessData, error)

	// DeleteAccessToken deletes the access token.
	DeleteAccessToken(accessToken string) error

	// SaveRefreshToken saves the refresh token.
	SaveRefreshToken(refresh *types.RefreshData) error

	// GetRefreshToken finds the refresh token.
	GetRefreshToken(refreshToken string) (*types.RefreshData, error)

	// DeleteRefreshToken deletes the refresh token from the database.
	DeleteRefreshToken(refreshToken string) error

	// ListAuthorizedClients returns a list of clients authorized by a given user.
	ListAuthorizedClients(userID string, factory ClientFactory) ([]types.Client, error)

	// RevokeAuthorizedClient deletes the access tokens and refresh token
	// granted to a client by a given user.
	RevokeAuthorizedClient(userID, clientID string) error
}
