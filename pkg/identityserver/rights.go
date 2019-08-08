// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package identityserver

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

func allPotentialRights(entityID ttnpb.Identifiers, rights *ttnpb.Rights) *ttnpb.Rights {
	switch entityID.EntityType() {
	case "application":
		return ttnpb.AllApplicationRights.Intersect(rights)
	case "client":
		return ttnpb.AllClientRights.Intersect(rights)
	case "gateway":
		return ttnpb.AllGatewayRights.Intersect(rights)
	case "organization":
		return ttnpb.AllEntityRights.Union(ttnpb.AllOrganizationRights).Intersect(rights)
	case "user":
		return ttnpb.AllEntityRights.Union(ttnpb.AllOrganizationRights, ttnpb.AllUserRights).Intersect(rights)
	}
	return nil
}

func (is *IdentityServer) getRights(ctx context.Context, entityID ttnpb.Identifiers) (entityRights, universalRights *ttnpb.Rights, err error) {
	authInfo, err := is.authInfo(ctx)
	if err != nil {
		return nil, nil, err
	}

	authInfoRights := ttnpb.RightsFrom(authInfo.GetRights()...)
	universalRights = allPotentialRights(entityID, authInfo.GetUniversalRights())
	if len(universalRights.GetRights()) == 0 {
		universalRights = nil
	}
	allPotentialRights := allPotentialRights(entityID, authInfoRights)

	// If the rights of the auth do not contain any rights for the entity type,
	// there's nothing more to do.
	if len(allPotentialRights.GetRights()) == 0 {
		return nil, universalRights, nil
	}

	// If the caller is the requested entity,
	// we can directly return the rights of the auth.
	authenticatedAs := authInfo.GetEntityIdentifiers()
	if entityID.EntityType() == authenticatedAs.EntityType() &&
		entityID.IDString() == authenticatedAs.IDString() {
		return authInfoRights, universalRights, nil
	}

	// If the caller is not an organization or user, there's nothing more to do.
	ouID := authInfo.GetOrganizationOrUserIdentifiers()
	if ouID == nil {
		return nil, universalRights, nil
	}

	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		membershipStore := is.getMembershipStore(ctx, db)

		// Find direct membership rights of the organization or user.
		directMemberRights, err := membershipStore.GetMember(ctx, ouID, entityID)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}

		// Expand the pseudo-rights.
		entityRights = directMemberRights.Implied()

		// If entityRights already includes all potential rights,
		// there's nothing more to do.
		if len(allPotentialRights.Sub(entityRights).GetRights()) == 0 {
			return nil
		}

		// If the caller is not a user, there's nothing more to do.
		usrID := ouID.GetUserIDs()
		if usrID == nil {
			return nil
		}

		// Find indirect memberships (through organizations).
		// TODO: Cache this (https://github.com/TheThingsNetwork/lorawan-stack/issues/443).
		commonOrganizations, err := membershipStore.FindIndirectMemberships(ctx, usrID, entityID)
		if err != nil {
			return err
		}
		for _, commonOrganization := range commonOrganizations {
			rightsOnOrganization := commonOrganization.RightsOnOrganization.Implied()
			organizationRights := commonOrganization.OrganizationRights.Implied()
			entityRights = entityRights.Union(rightsOnOrganization.Intersect(organizationRights))
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	entityRights = entityRights.Intersect(authInfoRights)

	return entityRights, universalRights, err
}

// ApplicationRights returns the rights the caller has on the given application.
func (is *IdentityServer) ApplicationRights(ctx context.Context, appIDs ttnpb.ApplicationIdentifiers) (*ttnpb.Rights, error) {
	entity, universal, err := is.getRights(ctx, appIDs)
	if err != nil {
		return nil, err
	}
	if entity != nil {
		return entity.Union(universal), nil
	}
	if !is.IsAdmin(ctx) && universal == nil {
		return &ttnpb.Rights{}, nil
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		_, err := store.GetApplicationStore(db).GetApplication(ctx, &appIDs, &types.FieldMask{Paths: []string{"ids"}})
		return err
	})
	if err != nil {
		return nil, err
	}
	return universal, nil
}

// ClientRights returns the rights the caller has on the given client.
func (is *IdentityServer) ClientRights(ctx context.Context, cliIDs ttnpb.ClientIdentifiers) (*ttnpb.Rights, error) {
	entity, universal, err := is.getRights(ctx, cliIDs)
	if err != nil {
		return nil, err
	}
	if entity != nil {
		return entity.Union(universal), nil
	}
	if !is.IsAdmin(ctx) && universal == nil {
		return &ttnpb.Rights{}, nil
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		_, err := store.GetClientStore(db).GetClient(ctx, &cliIDs, &types.FieldMask{Paths: []string{"ids"}})
		return err
	})
	if err != nil {
		return nil, err
	}
	return universal, nil
}

// GatewayRights returns the rights the caller has on the given gateway.
func (is *IdentityServer) GatewayRights(ctx context.Context, gtwIDs ttnpb.GatewayIdentifiers) (*ttnpb.Rights, error) {
	entity, universal, err := is.getRights(ctx, gtwIDs)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if entity != nil {
		return entity.Union(universal), nil
	}
	if !is.IsAdmin(ctx) && universal == nil {
		return &ttnpb.Rights{}, nil
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		_, err := store.GetGatewayStore(db).GetGateway(ctx, &gtwIDs, &types.FieldMask{Paths: []string{"ids"}})
		return err
	})
	if err != nil {
		return nil, err
	}
	return universal, nil
}

// OrganizationRights returns the rights the caller has on the given organization.
func (is *IdentityServer) OrganizationRights(ctx context.Context, orgIDs ttnpb.OrganizationIdentifiers) (*ttnpb.Rights, error) {
	entity, universal, err := is.getRights(ctx, orgIDs)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if entity != nil {
		return entity.Union(universal), nil
	}
	if !is.IsAdmin(ctx) && universal == nil {
		return &ttnpb.Rights{}, nil
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		_, err := store.GetOrganizationStore(db).GetOrganization(ctx, &orgIDs, &types.FieldMask{Paths: []string{"ids"}})
		return err
	})
	if err != nil {
		return nil, err
	}
	return universal, nil
}

// UserRights returns the rights the caller has on the given user.
func (is *IdentityServer) UserRights(ctx context.Context, userIDs ttnpb.UserIdentifiers) (*ttnpb.Rights, error) {
	entity, universal, err := is.getRights(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if entity != nil {
		return entity.Union(universal), nil
	}
	if !is.IsAdmin(ctx) && universal == nil {
		return &ttnpb.Rights{}, nil
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		_, err := store.GetUserStore(db).GetUser(ctx, &userIDs, &types.FieldMask{Paths: []string{"ids"}})
		return err
	})
	if err != nil {
		return nil, err
	}
	return universal, nil
}
