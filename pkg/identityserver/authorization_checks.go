// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package identityserver

import (
	"context"

	"github.com/TheThingsNetwork/ttn/pkg/auth"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
)

// userCheck checks whether claims are intended for an user and then if the user
// has the given set of rights with it. It returns the user ID.
func (is *IdentityServer) userCheck(ctx context.Context, rights ...ttnpb.Right) (string, error) {
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

// adminCheck undercalls `userCheck` with `RIGHT_USER_ADMIN` right and then fetches
// from the store the user to check if it has activated the admin flag.
func (is *IdentityServer) adminCheck(ctx context.Context) error {
	userID, err := is.userCheck(ctx, ttnpb.RIGHT_USER_ADMIN)
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

// applicationCheck checks whether claims have the given set of rights and then:
// 	-	If they come from an API key: checks whether the API key application ID matches
//			with the application ID that the request is trying to access to.
// 	-	If they come from an access token: checks whether the user is collaborator of the
//      application with the given set of rights.
func (is *IdentityServer) applicationCheck(ctx context.Context, appID string, rights ...ttnpb.Right) error {
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

// gatewayCheck checks whether claims have the given set of rights and then:
// 	-	If they come from an API key: checks whether the API key gateway ID matches
//			with the gateway ID that the request is trying to access to.
// 	-	If they come from an access token: checks whether the user is collaborator of the
//      gateway application with the given set of rights.
func (is *IdentityServer) gatewayCheck(ctx context.Context, gtwID string, rights ...ttnpb.Right) error {
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
