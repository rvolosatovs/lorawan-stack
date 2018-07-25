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
	"strings"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/auth"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/identityserver/email/templates"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/random"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// SettingsGeneratedFields are the fields that are automatically generated.
var SettingsGeneratedFields = []string{"UpdatedAt"}

type adminService struct {
	*IdentityServer
}

// GetSettings fetches the current dynamic settings of the Identity Server.
func (s *adminService) GetSettings(ctx context.Context, _ *pbtypes.Empty) (*ttnpb.IdentityServerSettings, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	return s.store.Settings.Get()
}

// UpdateSettings updates the dynamic settings.
func (s *adminService) UpdateSettings(ctx context.Context, req *ttnpb.UpdateSettingsRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.store.Transact(func(tx *store.Store) (err error) {
		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		for _, path := range req.UpdateMask.Paths {
			switch path {
			case ttnpb.FieldPathSettingsBlacklistedIDs:
				if req.Settings.BlacklistedIDs == nil {
					req.Settings.BlacklistedIDs = []string{}
				}
				settings.BlacklistedIDs = req.Settings.BlacklistedIDs
			case ttnpb.FieldPathSettingsUserRegistrationSkipValidation:
				settings.SkipValidation = req.Settings.SkipValidation
			case ttnpb.FieldPathSettingsUserRegistrationInvitationOnly:
				settings.InvitationOnly = req.Settings.InvitationOnly
			case ttnpb.FieldPathSettingsUserRegistrationAdminApproval:
				settings.AdminApproval = req.Settings.AdminApproval
			case ttnpb.FieldPathSettingsValidationTokenTTL:
				settings.ValidationTokenTTL = req.Settings.ValidationTokenTTL
			case ttnpb.FieldPathSettingsAllowedEmails:
				if req.Settings.AllowedEmails == nil {
					req.Settings.AllowedEmails = []string{}
				}
				settings.AllowedEmails = req.Settings.AllowedEmails
			case ttnpb.FieldPathSettingsInvitationTokenTTL:
				settings.InvitationTokenTTL = req.Settings.InvitationTokenTTL
			default:
				return ttnpb.ErrInvalidPathUpdateMask.New(errors.Attributes{
					"path": path,
				})
			}
		}

		settings.UpdatedAt = time.Now().UTC()

		return tx.Settings.Set(*settings)
	})

	return ttnpb.Empty, err
}

// CreateUser creates an account on behalf of an user. A password is generated
// and sent to the user's email.
func (s *adminService) CreateUser(ctx context.Context, req *ttnpb.CreateUserRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// Set an autogenerated password.
	req.User.Password = random.String(8)

	// Mark user as approved.
	req.User.State = ttnpb.STATE_APPROVED

	var token string
	err = s.store.Transact(func(tx *store.Store) (err error) {
		now := time.Now().UTC()

		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		// Check for blacklisted IDs.
		if !settings.IsIDAllowed(req.User.UserID) {
			return ErrBlacklistedID.New(errors.Attributes{
				"id": req.User.UserID,
			})
		}

		if settings.SkipValidation {
			req.User.ValidatedAt = timeValue(now)
		}

		req.User.CreatedAt = now
		req.User.UpdatedAt = now
		req.User.PasswordUpdatedAt = now

		req.User.RequirePasswordUpdate = true

		err = tx.Users.Create(&req.User)
		if err != nil {
			return err
		}

		// If validation can be skipped just finish transaction.
		if settings.SkipValidation {
			return nil
		}

		// Otherwise create a token and save it.
		token = random.String(64)

		return tx.Users.SaveValidationToken(req.User.UserIdentifiers, store.ValidationToken{
			ValidationToken: token,
			CreatedAt:       time.Now(),
			ExpiresIn:       int32(settings.ValidationTokenTTL.Seconds()),
		})
	})

	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.email.Send(req.User.UserIdentifiers.Email, &templates.AccountCreation{
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
		Name:             req.User.Name,
		UserID:           req.User.UserID,
		Password:         req.User.Password,
		ValidationToken:  token,
	})
}

// GetUser returns the user account that matches the identifier.
func (s *adminService) GetUser(ctx context.Context, req *ttnpb.UserIdentifiers) (*ttnpb.User, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	found, err := s.store.Users.GetByID(*req, s.specializers.User)
	if err != nil {
		return nil, err
	}
	found.GetUser().Password = ""

	return found.GetUser(), nil
}

// ListUsers returns a list of users with optional filtering.
func (s *adminService) ListUsers(ctx context.Context, req *ttnpb.ListUsersRequest) (*ttnpb.ListUsersResponse, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.store.Users.List(s.specializers.User)
	if err != nil {
		return nil, err
	}

	resp := &ttnpb.ListUsersResponse{
		Users: make([]*ttnpb.User, 0, len(users)),
	}

	// Filter results manually.
	for _, user := range users {
		u := user.GetUser()

		if filter := req.ListUsersRequest_FilterState; filter != nil && filter.State != u.State {
			continue
		}

		resp.Users = append(resp.Users, u)
	}

	return resp, nil
}

// UpdateUser updates an user account. If email address is updated it sends an
// email to validate it if and only if the `SkipValidation` setting is disabled.
func (s *adminService) UpdateUser(ctx context.Context, req *ttnpb.UpdateUserRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var user *ttnpb.User
	var token *store.ValidationToken

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Users.GetByID(req.User.UserIdentifiers, s.specializers.User)
		if err != nil {
			return err
		}
		user = found.GetUser()

		// Save the current user identifiers before applying the mask.
		oids := found.GetUser().UserIdentifiers

		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		newEmail := false
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case ttnpb.FieldPathUserName:
				user.Name = req.User.Name
			case ttnpb.FieldPathUserEmail:
				user.UserIdentifiers.Email = req.User.UserIdentifiers.Email

				newEmail = strings.ToLower(user.UserIdentifiers.Email) != strings.ToLower(req.User.UserIdentifiers.Email)
				if newEmail {
					if settings.SkipValidation {
						user.ValidatedAt = timeValue(time.Now())
					} else {
						user.ValidatedAt = timeValue(time.Time{})
					}
				}
			case ttnpb.FieldPathUserAdmin:
				user.Admin = req.User.Admin
			case ttnpb.FieldPathUserState:
				user.State = req.User.State
			case ttnpb.FieldPathUserRequirePasswordUpdate:
				user.RequirePasswordUpdate = req.User.RequirePasswordUpdate
			default:
				return ttnpb.ErrInvalidPathUpdateMask.New(errors.Attributes{
					"path": path,
				})
			}
		}

		err = tx.Users.Update(oids, user)
		if err != nil {
			return err
		}

		if !newEmail || (newEmail && settings.SkipValidation) {
			return nil
		}

		token = &store.ValidationToken{
			ValidationToken: random.String(64),
			CreatedAt:       time.Now(),
			ExpiresIn:       int32(settings.ValidationTokenTTL.Seconds()),
		}

		return tx.Users.SaveValidationToken(user.UserIdentifiers, *token)
	})

	if err != nil {
		return nil, err
	}

	// No email needs to be send.
	if token == nil {
		return ttnpb.Empty, nil
	}

	return ttnpb.Empty, s.email.Send(user.UserIdentifiers.Email, &templates.EmailValidation{
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
		Token:            token.ValidationToken,
	})
}

// ResetUserPassword sets an autogenerated password to the user that matches the
// identifier. The new password is returned on the response but also send by email
// to the user.
func (s *adminService) ResetUserPassword(ctx context.Context, req *ttnpb.UserIdentifiers) (*ttnpb.ResetUserPasswordResponse, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	password := random.String(8)

	hashed, err := auth.Hash(password)
	if err != nil {
		return nil, err
	}

	var user *ttnpb.User

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Users.GetByID(*req, s.specializers.User)
		if err != nil {
			return err
		}

		user = found.GetUser()
		user.Password = string(hashed)
		user.PasswordUpdatedAt = time.Now().UTC()
		user.RequirePasswordUpdate = true

		return tx.Users.Update(user.UserIdentifiers, user)
	})

	if err != nil {
		return nil, err
	}

	err = s.email.Send(user.UserIdentifiers.Email, &templates.PasswordReset{
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
		Password:         password,
	})

	return &ttnpb.ResetUserPasswordResponse{
		Password: password,
	}, err
}

// DeleteUser deletes an user.
func (s *adminService) DeleteUser(ctx context.Context, req *ttnpb.UserIdentifiers) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var user *ttnpb.User

	err = s.store.Transact(func(tx *store.Store) error {
		ids := *req

		// Fetch the user beforehand to save its email address for later notification.
		found, err := tx.Users.GetByID(ids, s.specializers.User)
		if err != nil {
			return err
		}
		user = found.GetUser()

		apps, err := tx.Applications.ListByOrganizationOrUser(organizationOrUserIDsUserIDs(ids), s.specializers.Application)
		if err != nil {
			return err
		}

		gtws, err := tx.Gateways.ListByOrganizationOrUser(organizationOrUserIDsUserIDs(ids), s.specializers.Gateway)
		if err != nil {
			return err
		}

		orgs, err := tx.Organizations.ListByUser(ids, s.specializers.Organization)
		if err != nil {
			return err
		}

		err = tx.Users.Delete(ids)
		if err != nil {
			return err
		}

		for _, app := range apps {
			appIDs := app.GetApplication().ApplicationIdentifiers

			rights, err := missingApplicationRights(tx, appIDs)
			if err != nil {
				return err
			}

			if len(rights) != 0 {
				return ErrUnmanageableApplication.New(errors.Attributes{
					"application_id": appIDs.ApplicationID,
					"missing_rights": rights,
				})
			}
		}

		for _, gtw := range gtws {
			gtwIDs := gtw.GetGateway().GatewayIdentifiers

			rights, err := missingGatewayRights(tx, gtwIDs)
			if err != nil {
				return err
			}

			if len(rights) != 0 {
				return ErrUnmanageableGateway.New(errors.Attributes{
					"gateway_id":     gtwIDs.GatewayID,
					"missing_rights": rights,
				})
			}
		}

		for _, org := range orgs {
			orgIDs := org.GetOrganization().OrganizationIdentifiers

			rights, err := missingOrganizationRights(tx, orgIDs)
			if err != nil {
				return err
			}

			if len(rights) != 0 {
				return ErrUnmanageableOrganization.New(errors.Attributes{
					"organization_id": orgIDs.OrganizationID,
					"missing_rights":  rights,
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.email.Send(user.UserIdentifiers.Email, &templates.AccountDeleted{
		UserID:           req.UserID,
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
	})
}

// SendInvitation sends by email a token that can be used to create a new account.
// All invitations are expirable and the TTL is defined on a setitngs variable.
func (s *adminService) SendInvitation(ctx context.Context, req *ttnpb.SendInvitationRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var invitation store.InvitationData

	err = s.store.Transact(func(tx *store.Store) (err error) {
		// Check whether email is already registered or not.
		found, err := tx.Users.GetByID(ttnpb.UserIdentifiers{Email: req.Email}, s.specializers.User)
		if err != nil && !store.ErrUserNotFound.Describes(err) {
			return err
		}

		// If email is already being used return error.
		if found != nil {
			return ErrEmailAddressAlreadyUsed.New(nil)
		}

		// Otherwise proceed to issue invitation.
		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		invitation = store.InvitationData{
			Token:     random.String(64),
			Email:     req.Email,
			IssuedAt:  now,
			ExpiresAt: now.Add(settings.InvitationTokenTTL),
		}

		return tx.Invitations.Save(invitation)
	})

	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.email.Send(req.Email, &templates.Invitation{
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
		Token:            invitation.Token,
	})
}

// ListInvitations lists all the issued invitations.
func (s *adminService) ListInvitations(ctx context.Context, req *pbtypes.Empty) (*ttnpb.ListInvitationsResponse, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	invitations, err := s.store.Invitations.List()
	if err != nil {
		return nil, err
	}

	resp := &ttnpb.ListInvitationsResponse{
		Invitations: make([]*ttnpb.ListInvitationsResponse_Invitation, 0, len(invitations)),
	}

	for _, invitation := range invitations {
		resp.Invitations = append(resp.Invitations, &ttnpb.ListInvitationsResponse_Invitation{
			Email:     invitation.Email,
			IssuedAt:  invitation.IssuedAt,
			ExpiresAt: invitation.ExpiresAt,
		})
	}

	return resp, nil
}

// DeleteInvitation revokes an unused invitation or deletes an expired one.
func (s *adminService) DeleteInvitation(ctx context.Context, req *ttnpb.DeleteInvitationRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.store.Invitations.Delete(req.Email)
}

// GetClient returns the client that matches the identifier.
func (s *adminService) GetClient(ctx context.Context, req *ttnpb.ClientIdentifiers) (*ttnpb.Client, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	found, err := s.store.Clients.GetByID(*req, s.specializers.Client)
	if err != nil {
		return nil, err
	}

	return found.GetClient(), nil
}

// ListClients returns a list of third-party clients with optional filtering.
func (s *adminService) ListClients(ctx context.Context, req *ttnpb.ListClientsRequest) (*ttnpb.ListClientsResponse, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	found, err := s.store.Clients.List(s.specializers.Client)
	if err != nil {
		return nil, err
	}

	resp := &ttnpb.ListClientsResponse{
		Clients: make([]*ttnpb.Client, 0, len(found)),
	}

	// Filter results manually.
	for _, client := range found {
		cli := client.GetClient()

		if filter := req.ListClientsRequest_FilterState; filter != nil && filter.State != cli.State {
			continue
		}

		resp.Clients = append(resp.Clients, cli)
	}

	return resp, nil
}

// UpdateClient updates a third-party client.
func (s *adminService) UpdateClient(ctx context.Context, req *ttnpb.UpdateClientRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Clients.GetByID(req.Client.ClientIdentifiers, s.specializers.Client)
		if err != nil {
			return err
		}
		client := found.GetClient()

		for _, path := range req.UpdateMask.Paths {
			switch path {
			case ttnpb.FieldPathClientDescription:
				client.Description = req.Client.Description
			case ttnpb.FieldPathClientRedirectURI:
				client.RedirectURI = req.Client.RedirectURI
			case ttnpb.FieldPathClientRights:
				if req.Client.Rights == nil {
					req.Client.Rights = []ttnpb.Right{}
				}
				client.Rights = req.Client.Rights
			case ttnpb.FieldPathClientSkipAuthorization:
				client.SkipAuthorization = req.Client.SkipAuthorization
			case ttnpb.FieldPathClientState:
				client.State = req.Client.State
			case ttnpb.FieldPathClientGrants:
				if req.Client.Grants == nil {
					req.Client.Grants = []ttnpb.GrantType{}
				}
				client.Grants = req.Client.Grants
			default:
				return ttnpb.ErrInvalidPathUpdateMask.New(errors.Attributes{
					"path": path,
				})
			}
		}

		return tx.Clients.Update(client)
	})

	return ttnpb.Empty, err
}

// DeleteClient deletes the client that matches the identifier and revokes all
// user authorizations.
func (s *adminService) DeleteClient(ctx context.Context, req *ttnpb.ClientIdentifiers) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var user store.User

	err = s.store.Transact(func(tx *store.Store) error {
		ids := *req

		found, err := tx.Clients.GetByID(ids, s.specializers.Client)
		if err != nil {
			return err
		}

		user, err = tx.Users.GetByID(found.GetClient().CreatorIDs, s.specializers.User)
		if err != nil {
			return err
		}

		return tx.Clients.Delete(ids)
	})

	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.email.Send(user.GetUser().UserIdentifiers.Email, &templates.ClientDeleted{
		ClientID:         req.ClientID,
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
	})
}
