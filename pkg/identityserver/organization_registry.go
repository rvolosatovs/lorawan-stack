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

package identityserver

import (
	"context"
	"strconv"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/identityserver/blacklist"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	evtCreateOrganization = events.Define("organization.create", "Create organization")
	evtUpdateOrganization = events.Define("organization.update", "Update organization")
	evtDeleteOrganization = events.Define("organization.delete", "Delete organization")
)

var errNestedOrganizations = errors.DefineInvalidArgument("nested_organizations", "organizations can not be nested")

func (is *IdentityServer) createOrganization(ctx context.Context, req *ttnpb.CreateOrganizationRequest) (org *ttnpb.Organization, err error) {
	if err := blacklist.Check(ctx, req.OrganizationID); err != nil {
		return nil, err
	}
	if usrIDs := req.Collaborator.GetUserIDs(); usrIDs != nil {
		if err = rights.RequireUser(ctx, *usrIDs, ttnpb.RIGHT_USER_ORGANIZATIONS_CREATE); err != nil {
			return nil, err
		}
	} else if orgIDs := req.Collaborator.GetOrganizationIDs(); orgIDs != nil {
		return nil, errNestedOrganizations
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		orgStore := store.GetOrganizationStore(db)
		org, err = orgStore.CreateOrganization(ctx, &req.Organization)
		if err != nil {
			return err
		}
		memberStore := store.GetMembershipStore(db)
		err = memberStore.SetMember(ctx, &req.Collaborator, org.OrganizationIdentifiers.EntityIdentifiers(), ttnpb.RightsFrom(ttnpb.RIGHT_ORGANIZATION_ALL))
		if err != nil {
			return err
		}
		if len(req.ContactInfo) > 0 {
			cleanContactInfo(req.ContactInfo)
			org.ContactInfo, err = store.GetContactInfoStore(db).SetContactInfo(ctx, org.EntityIdentifiers(), req.ContactInfo)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtCreateOrganization(ctx, req.OrganizationIdentifiers, nil))
	is.invalidateCachedMembershipsForAccount(ctx, &req.Collaborator)
	return org, nil
}

func (is *IdentityServer) getOrganization(ctx context.Context, req *ttnpb.GetOrganizationRequest) (org *ttnpb.Organization, err error) {
	err = rights.RequireOrganization(ctx, req.OrganizationIdentifiers, ttnpb.RIGHT_ORGANIZATION_INFO)
	if err != nil {
		return nil, err
	}
	// TODO: Filter FieldMask by Rights
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		orgStore := store.GetOrganizationStore(db)
		org, err = orgStore.GetOrganization(ctx, &req.OrganizationIdentifiers, &req.FieldMask)
		if err != nil {
			return err
		}
		if fieldMaskContains(&req.FieldMask, "contact_info") {
			org.ContactInfo, err = store.GetContactInfoStore(db).GetContactInfo(ctx, org.EntityIdentifiers())
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (is *IdentityServer) listOrganizations(ctx context.Context, req *ttnpb.ListOrganizationsRequest) (orgs *ttnpb.Organizations, err error) {
	var orgRights map[string]*ttnpb.Rights
	if req.Collaborator == nil {
		callerRights, err := is.getRights(ctx)
		if err != nil {
			return nil, err
		}
		orgRights = make(map[string]*ttnpb.Rights, len(callerRights))
		for ids, rights := range callerRights {
			if ids := ids.GetOrganizationIDs(); ids != nil {
				orgRights[unique.ID(ctx, ids)] = rights
			}
		}
		if len(orgRights) == 0 {
			return &ttnpb.Organizations{}, nil
		}
	}
	if usrIDs := req.Collaborator.GetUserIDs(); usrIDs != nil {
		if err = rights.RequireUser(ctx, *usrIDs, ttnpb.RIGHT_USER_ORGANIZATIONS_LIST); err != nil {
			return nil, err
		}
	} else if orgIDs := req.Collaborator.GetOrganizationIDs(); orgIDs != nil {
		return nil, errNestedOrganizations
	}
	var total uint64
	ctx = store.SetTotalCount(ctx, &total)
	defer func() {
		if err == nil {
			grpc.SetHeader(ctx, metadata.Pairs("x-total-count", strconv.FormatUint(total, 10)))
		}
	}()
	orgs = new(ttnpb.Organizations)
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		if orgRights == nil {
			memberStore := store.GetMembershipStore(db)
			rights, err := memberStore.FindMemberRights(ctx, req.Collaborator, "organization")
			if err != nil {
				return err
			}
			orgRights = make(map[string]*ttnpb.Rights, len(rights))
			for ids, rights := range rights {
				orgRights[unique.ID(ctx, ids)] = rights
			}
		}
		if len(orgRights) == 0 {
			return nil
		}
		orgIDs := make([]*ttnpb.OrganizationIdentifiers, 0, len(orgRights))
		for uid := range orgRights {
			orgID, err := unique.ToOrganizationID(uid)
			if err != nil {
				continue
			}
			orgIDs = append(orgIDs, &orgID)
		}
		orgStore := store.GetOrganizationStore(db)
		orgs.Organizations, err = orgStore.FindOrganizations(ctx, orgIDs, &req.FieldMask)
		if err != nil {
			return err
		}
		for _, org := range orgs.Organizations {
			// TODO: Filter FieldMask by Rights
			_ = orgRights[unique.ID(ctx, org.OrganizationIdentifiers)]
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (is *IdentityServer) updateOrganization(ctx context.Context, req *ttnpb.UpdateOrganizationRequest) (org *ttnpb.Organization, err error) {
	err = rights.RequireOrganization(ctx, req.OrganizationIdentifiers, ttnpb.RIGHT_ORGANIZATION_SETTINGS_BASIC)
	if err != nil {
		return nil, err
	}
	// TODO: Filter FieldMask by Rights
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		orgStore := store.GetOrganizationStore(db)
		org, err = orgStore.UpdateOrganization(ctx, &req.Organization, &req.FieldMask)
		if err != nil {
			return err
		}
		if fieldMaskContains(&req.FieldMask, "contact_info") {
			cleanContactInfo(req.ContactInfo)
			org.ContactInfo, err = store.GetContactInfoStore(db).SetContactInfo(ctx, org.EntityIdentifiers(), req.ContactInfo)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtUpdateOrganization(ctx, req.OrganizationIdentifiers, req.FieldMask.Paths))
	return org, nil
}

func (is *IdentityServer) deleteOrganization(ctx context.Context, ids *ttnpb.OrganizationIdentifiers) (*types.Empty, error) {
	err := rights.RequireOrganization(ctx, *ids, ttnpb.RIGHT_ORGANIZATION_DELETE)
	if err != nil {
		return nil, err
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		orgStore := store.GetOrganizationStore(db)
		err = orgStore.DeleteOrganization(ctx, ids)
		return err
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtDeleteOrganization(ctx, ids, nil))
	// TODO: Invalidate rights of members
	return ttnpb.Empty, nil
}

type organizationRegistry struct {
	*IdentityServer
}

func (or *organizationRegistry) Create(ctx context.Context, req *ttnpb.CreateOrganizationRequest) (*ttnpb.Organization, error) {
	return or.createOrganization(ctx, req)
}
func (or *organizationRegistry) Get(ctx context.Context, req *ttnpb.GetOrganizationRequest) (*ttnpb.Organization, error) {
	return or.getOrganization(ctx, req)
}
func (or *organizationRegistry) List(ctx context.Context, req *ttnpb.ListOrganizationsRequest) (*ttnpb.Organizations, error) {
	return or.listOrganizations(ctx, req)
}
func (or *organizationRegistry) Update(ctx context.Context, req *ttnpb.UpdateOrganizationRequest) (*ttnpb.Organization, error) {
	return or.updateOrganization(ctx, req)
}
func (or *organizationRegistry) Delete(ctx context.Context, req *ttnpb.OrganizationIdentifiers) (*types.Empty, error) {
	return or.deleteOrganization(ctx, req)
}
