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
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/auth"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/email"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/identityserver/blacklist"
	"go.thethings.network/lorawan-stack/pkg/identityserver/emails"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	evtCreateClient = events.Define(
		"client.create", "create OAuth client",
		ttnpb.RIGHT_CLIENT_ALL,
	)
	evtUpdateClient = events.Define(
		"client.update", "update OAuth client",
		ttnpb.RIGHT_CLIENT_ALL,
	)
	evtDeleteClient = events.Define(
		"client.delete", "delete OAuth client",
		ttnpb.RIGHT_CLIENT_ALL,
	)
)

func (is *IdentityServer) createClient(ctx context.Context, req *ttnpb.CreateClientRequest) (cli *ttnpb.Client, err error) {
	createdByAdmin := is.IsAdmin(ctx)
	if err = blacklist.Check(ctx, req.ClientID); err != nil {
		return nil, err
	}
	if usrIDs := req.Collaborator.GetUserIDs(); usrIDs != nil {
		if err = rights.RequireUser(ctx, *usrIDs, ttnpb.RIGHT_USER_CLIENTS_CREATE); err != nil {
			return nil, err
		}
	} else if orgIDs := req.Collaborator.GetOrganizationIDs(); orgIDs != nil {
		if err = rights.RequireOrganization(ctx, *orgIDs, ttnpb.RIGHT_ORGANIZATION_CLIENTS_CREATE); err != nil {
			return nil, err
		}
	}

	if err := validateContactInfo(req.Client.ContactInfo); err != nil {
		return nil, err
	}

	secret := req.Client.Secret
	if secret == "" {
		secret, err = auth.GenerateKey(ctx)
		if err != nil {
			return nil, err
		}
	}
	hashedSecret, err := auth.Hash(ctx, secret)
	if err != nil {
		return nil, err
	}
	req.Client.Secret = hashedSecret

	if !createdByAdmin {
		req.Client.State = ttnpb.STATE_REQUESTED
		req.Client.SkipAuthorization = false
		req.Client.Endorsed = false
	}

	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		cli, err = store.GetClientStore(db).CreateClient(ctx, &req.Client)
		if err != nil {
			return err
		}
		if err = is.getMembershipStore(ctx, db).SetMember(
			ctx,
			&req.Collaborator,
			cli.ClientIdentifiers,
			ttnpb.RightsFrom(ttnpb.RIGHT_ALL),
		); err != nil {
			return err
		}
		if len(req.ContactInfo) > 0 {
			cleanContactInfo(req.ContactInfo)
			cli.ContactInfo, err = store.GetContactInfoStore(db).SetContactInfo(ctx, cli.ClientIdentifiers, req.ContactInfo)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	cli.Secret = secret // Return the unhashed secret, in case it was generated.

	events.Publish(evtCreateClient(ctx, req.ClientIdentifiers, nil))
	return cli, nil
}

func (is *IdentityServer) getClient(ctx context.Context, req *ttnpb.GetClientRequest) (cli *ttnpb.Client, err error) {
	if err = is.RequireAuthenticated(ctx); err != nil {
		return nil, err
	}
	req.FieldMask.Paths = cleanFieldMaskPaths(ttnpb.ClientFieldPathsNested, req.FieldMask.Paths, getPaths, nil)
	if err = rights.RequireClient(ctx, req.ClientIdentifiers, ttnpb.RIGHT_CLIENT_ALL); err != nil {
		if ttnpb.HasOnlyAllowedFields(req.FieldMask.Paths, ttnpb.PublicClientFields...) {
			defer func() { cli = cli.PublicSafe() }()
		} else {
			return nil, err
		}
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		cli, err = store.GetClientStore(db).GetClient(ctx, &req.ClientIdentifiers, &req.FieldMask)
		if err != nil {
			return err
		}
		if ttnpb.HasAnyField(req.FieldMask.Paths, "contact_info") {
			cli.ContactInfo, err = store.GetContactInfoStore(db).GetContactInfo(ctx, cli.ClientIdentifiers)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (is *IdentityServer) listClients(ctx context.Context, req *ttnpb.ListClientsRequest) (clis *ttnpb.Clients, err error) {
	req.FieldMask.Paths = cleanFieldMaskPaths(ttnpb.ClientFieldPathsNested, req.FieldMask.Paths, getPaths, nil)
	var includeIndirect bool
	if req.Collaborator == nil {
		authInfo, err := is.authInfo(ctx)
		if err != nil {
			return nil, err
		}
		collaborator := authInfo.GetOrganizationOrUserIdentifiers()
		if collaborator == nil {
			return &ttnpb.Clients{}, nil
		}
		req.Collaborator = collaborator
		includeIndirect = true
	}
	if usrIDs := req.Collaborator.GetUserIDs(); usrIDs != nil {
		if err = rights.RequireUser(ctx, *usrIDs, ttnpb.RIGHT_USER_CLIENTS_LIST); err != nil {
			return nil, err
		}
	} else if orgIDs := req.Collaborator.GetOrganizationIDs(); orgIDs != nil {
		if err = rights.RequireOrganization(ctx, *orgIDs, ttnpb.RIGHT_ORGANIZATION_CLIENTS_LIST); err != nil {
			return nil, err
		}
	}
	var total uint64
	paginateCtx := store.WithPagination(ctx, req.Limit, req.Page, &total)
	defer func() {
		if err == nil {
			setTotalHeader(ctx, total)
		}
	}()
	clis = &ttnpb.Clients{}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		ids, err := is.getMembershipStore(ctx, db).FindMemberships(paginateCtx, req.Collaborator, "client", includeIndirect)
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}
		cliIDs := make([]*ttnpb.ClientIdentifiers, 0, len(ids))
		for _, id := range ids {
			if cliID := id.EntityIdentifiers().GetClientIDs(); cliID != nil {
				cliIDs = append(cliIDs, cliID)
			}
		}
		clis.Clients, err = store.GetClientStore(db).FindClients(ctx, cliIDs, &req.FieldMask)
		if err != nil {
			return err
		}
		for i, cli := range clis.Clients {
			if rights.RequireClient(ctx, cli.ClientIdentifiers, ttnpb.RIGHT_CLIENT_ALL) != nil {
				clis.Clients[i] = cli.PublicSafe()
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return clis, nil
}

var errUpdateClientAdminField = errors.DefinePermissionDenied("client_update_admin_field", "only admins can update the `{field}` field")

func (is *IdentityServer) updateClient(ctx context.Context, req *ttnpb.UpdateClientRequest) (cli *ttnpb.Client, err error) {
	if err = rights.RequireClient(ctx, req.ClientIdentifiers, ttnpb.RIGHT_CLIENT_ALL); err != nil {
		return nil, err
	}
	req.FieldMask.Paths = cleanFieldMaskPaths(ttnpb.ClientFieldPathsNested, req.FieldMask.Paths, nil, getPaths)
	if len(req.FieldMask.Paths) == 0 {
		req.FieldMask.Paths = updatePaths
	}
	if ttnpb.HasAnyField(req.FieldMask.Paths, "contact_info") {
		if err := validateContactInfo(req.Client.ContactInfo); err != nil {
			return nil, err
		}
	}
	updatedByAdmin := is.IsAdmin(ctx)

	if !updatedByAdmin {
		for _, path := range req.FieldMask.Paths {
			switch path {
			case "state", "skip_authorization", "endorsed", "grants":
				return nil, errUpdateUserAdminField.WithAttributes("field", path)
			}
		}
	}

	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		cli, err = store.GetClientStore(db).UpdateClient(ctx, &req.Client, &req.FieldMask)
		if err != nil {
			return err
		}
		if ttnpb.HasAnyField(req.FieldMask.Paths, "contact_info") {
			cleanContactInfo(req.ContactInfo)
			cli.ContactInfo, err = store.GetContactInfoStore(db).SetContactInfo(ctx, cli.ClientIdentifiers, req.ContactInfo)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtUpdateClient(ctx, req.ClientIdentifiers, req.FieldMask.Paths))
	if ttnpb.HasAnyField(req.FieldMask.Paths, "state") {
		err = is.SendContactsEmail(ctx, req.EntityIdentifiers(), func(data emails.Data) email.MessageData {
			data.SetEntity(req.EntityIdentifiers())
			return &emails.EntityStateChanged{Data: data, State: strings.ToLower(strings.TrimPrefix(cli.State.String(), "STATE_"))}
		})
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("Could not send state change notification email")
		}
	}
	return cli, nil
}

func (is *IdentityServer) deleteClient(ctx context.Context, ids *ttnpb.ClientIdentifiers) (*types.Empty, error) {
	if err := rights.RequireClient(ctx, *ids, ttnpb.RIGHT_CLIENT_ALL); err != nil {
		return nil, err
	}
	err := is.withDatabase(ctx, func(db *gorm.DB) error {
		return store.GetClientStore(db).DeleteClient(ctx, ids)
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtDeleteClient(ctx, ids, nil))
	return ttnpb.Empty, nil
}

type clientRegistry struct {
	*IdentityServer
}

func (cr *clientRegistry) Create(ctx context.Context, req *ttnpb.CreateClientRequest) (*ttnpb.Client, error) {
	return cr.createClient(ctx, req)
}

func (cr *clientRegistry) Get(ctx context.Context, req *ttnpb.GetClientRequest) (*ttnpb.Client, error) {
	return cr.getClient(ctx, req)
}

func (cr *clientRegistry) List(ctx context.Context, req *ttnpb.ListClientsRequest) (*ttnpb.Clients, error) {
	return cr.listClients(ctx, req)
}

func (cr *clientRegistry) Update(ctx context.Context, req *ttnpb.UpdateClientRequest) (*ttnpb.Client, error) {
	return cr.updateClient(ctx, req)
}

func (cr *clientRegistry) Delete(ctx context.Context, req *ttnpb.ClientIdentifiers) (*types.Empty, error) {
	return cr.deleteClient(ctx, req)
}
