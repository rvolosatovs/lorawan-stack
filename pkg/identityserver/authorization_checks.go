// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package identityserver

import (
	"context"

	"github.com/TheThingsNetwork/ttn/pkg/auth"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
)

// enforceUserRights is a hook that checks whether if the given authorization
// credentials are allowed to perform an action given the set of passed rights and
// returns the User ID attached to the credentials.
func (is *IdentityServer) enforceUserRights(ctx context.Context, rights ...ttnpb.Right) (string, error) {
	claims, err := is.claimsFromContext(ctx)
	if err != nil {
		return "", err
	}

	userID := claims.UserID()
	if userID == "" {
		return "", ErrNotAuthorized.New(nil)
	}

	if !claims.HasRights(rights...) {
		return "", ErrNotAuthorized.New(nil)
	}

	return userID, nil
}

// enforceAdmin checks whether the given credentials are enough to access an admin resource.
func (is *IdentityServer) enforceAdmin(ctx context.Context) error {
	userID, err := is.enforceUserRights(ctx, ttnpb.RIGHT_USER_ADMIN)
	if err != nil {
		return err
	}

	found, err := is.store.Users.GetByID(userID, is.factories.user)
	if err != nil {
		return err
	}

	if !found.GetUser().Admin {
		return ErrNotAuthorized.New(nil)
	}

	return nil
}

// enforceApplicationRights is a hook that checks whether if the given authorization
// credentials are allowed to access the application with the given rights.
func (is *IdentityServer) enforceApplicationRights(ctx context.Context, appID string, rights ...ttnpb.Right) error {
	claims, err := is.claimsFromContext(ctx)
	if err != nil {
		return err
	}

	if !claims.HasRights(rights...) {
		return ErrNotAuthorized.New(nil)
	}

	var authorized bool
	switch claims.Source {
	case auth.Key:
		authorized = claims.ApplicationID() == appID
	case auth.Token:
		userID := claims.UserID()
		if len(userID) == 0 {
			return ErrNotAuthorized.New(nil)
		}

		authorized, err = is.store.Applications.HasUserRights(appID, userID, rights...)
		if err != nil {
			return err
		}
	}

	if !authorized {
		return ErrNotAuthorized.New(nil)
	}

	return nil
}

// enforceGatewayRights is a hook that checks whether if the given authorization
// credentials are allowed to access the gateway with the given rights.
func (is *IdentityServer) enforceGatewayRights(ctx context.Context, gtwID string, rights ...ttnpb.Right) error {
	claims, err := is.claimsFromContext(ctx)
	if err != nil {
		return err
	}

	if !claims.HasRights(rights...) {
		return ErrNotAuthorized.New(nil)
	}

	var authorized bool
	switch claims.Source {
	case auth.Key:
		authorized = claims.GatewayID() == gtwID
	case auth.Token:
		userID := claims.UserID()
		if len(userID) == 0 {
			return ErrNotAuthorized.New(nil)
		}

		authorized, err = is.store.Gateways.HasUserRights(gtwID, userID, rights...)
		if err != nil {
			return err
		}
	}

	if !authorized {
		return ErrNotAuthorized.New(nil)
	}

	return nil
}
