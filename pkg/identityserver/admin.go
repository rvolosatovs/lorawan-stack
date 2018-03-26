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

	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/email/templates"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store/sql"
	"github.com/TheThingsNetwork/ttn/pkg/random"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	pbtypes "github.com/gogo/protobuf/types"
)

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
			switch {
			case ttnpb.FieldPathSettingsBlacklistedIDs.MatchString(path):
				if req.Settings.BlacklistedIDs == nil {
					req.Settings.BlacklistedIDs = []string{}
				}
				settings.BlacklistedIDs = req.Settings.BlacklistedIDs
			case ttnpb.FieldPathSettingsUserRegistrationSkipValidation.MatchString(path):
				settings.SkipValidation = req.Settings.SkipValidation
			case ttnpb.FieldPathSettingsUserRegistrationInvitationOnly.MatchString(path):
				settings.InvitationOnly = req.Settings.InvitationOnly
			case ttnpb.FieldPathSettingsUserRegistrationAdminApproval.MatchString(path):
				settings.AdminApproval = req.Settings.AdminApproval
			case ttnpb.FieldPathSettingsValidationTokenTTL.MatchString(path):
				settings.ValidationTokenTTL = req.Settings.ValidationTokenTTL
			case ttnpb.FieldPathSettingsAllowedEmails.MatchString(path):
				if req.Settings.AllowedEmails == nil {
					req.Settings.AllowedEmails = []string{}
				}
				settings.AllowedEmails = req.Settings.AllowedEmails
			case ttnpb.FieldPathSettingsInvitationTokenTTL.MatchString(path):
				settings.InvitationTokenTTL = req.Settings.InvitationTokenTTL
			default:
				return ttnpb.ErrInvalidPathUpdateMask.New(errors.Attributes{
					"path": path,
				})
			}
		}

		return tx.Settings.Set(*settings)
	})

	return nil, err
}

// CreateUser creates an account on behalf of an user. A password is generated
// and sent to the user's email.
func (s *adminService) CreateUser(ctx context.Context, req *ttnpb.CreateUserRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// set an autogenerated password
	req.User.Password = random.String(8)

	// mark user as approved
	req.User.State = ttnpb.STATE_APPROVED

	var token string
	err = s.store.Transact(func(tx *store.Store) (err error) {
		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		// check for blacklisted ids
		if !settings.IsIDAllowed(req.User.UserID) {
			return ErrBlacklistedID.New(errors.Attributes{
				"id": req.User.UserID,
			})
		}

		if settings.SkipValidation {
			req.User.ValidatedAt = timeValue(time.Now())
		}

		err = tx.Users.Create(&req.User)
		if err != nil {
			return err
		}

		// if validation can be skipped just finish transaction
		if settings.SkipValidation {
			return nil
		}

		// otherwise create a token and save it
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

	return nil, s.email.Send(req.User.UserIdentifiers.Email, &templates.AccountCreation{
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

	found, err := s.store.Users.GetByID(*req, s.config.Specializers.User)
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

	users, err := s.store.Users.List(s.config.Specializers.User)
	if err != nil {
		return nil, err
	}

	resp := &ttnpb.ListUsersResponse{
		Users: make([]*ttnpb.User, 0, len(users)),
	}

	// filter results manually
	for _, user := range users {
		u := user.GetUser()

		if filter := req.ListUsersRequest_FilterState; filter != nil && filter.State != u.State {
			continue
		}

		resp.Users = append(resp.Users, u)
	}

	return resp, nil
}

// UpdateUser updates an user account.
// If email address is updated it sends an email to validate it if and only if
// the `SkipValidation` setting is disabled.
func (s *adminService) UpdateUser(ctx context.Context, req *ttnpb.UpdateUserRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Users.GetByID(req.User.UserIdentifiers, s.config.Specializers.User)
		if err != nil {
			return err
		}
		user := found.GetUser()

		// Save the current user identifiers before applying the mask.
		oids := found.GetUser().UserIdentifiers

		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		newEmail := false
		for _, path := range req.UpdateMask.Paths {
			switch {
			case ttnpb.FieldPathUserName.MatchString(path):
				user.Name = req.User.Name
			case ttnpb.FieldPathUserEmail.MatchString(path):
				user.UserIdentifiers.Email = req.User.UserIdentifiers.Email

				newEmail = strings.ToLower(user.UserIdentifiers.Email) != strings.ToLower(req.User.UserIdentifiers.Email)
				if newEmail {
					if settings.SkipValidation {
						user.ValidatedAt = timeValue(time.Now())
					} else {
						user.ValidatedAt = timeValue(time.Time{})
					}
				}
			case ttnpb.FieldPathUserAdmin.MatchString(path):
				user.Admin = req.User.Admin
			case ttnpb.FieldPathUserState.MatchString(path):
				user.State = req.User.State
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

		token := store.ValidationToken{
			ValidationToken: random.String(64),
			CreatedAt:       time.Now(),
			ExpiresIn:       int32(settings.ValidationTokenTTL.Seconds()),
		}

		err = tx.Users.SaveValidationToken(user.UserIdentifiers, token)
		if err != nil {
			return err
		}

		return s.email.Send(user.UserIdentifiers.Email, &templates.EmailValidation{
			OrganizationName: s.config.OrganizationName,
			PublicURL:        s.config.PublicURL,
			Token:            token.ValidationToken,
		})
	})

	return nil, err
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

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Users.GetByID(*req, s.config.Specializers.User)
		if err != nil {
			return err
		}

		user := found.GetUser()
		user.Password = password

		err = tx.Users.Update(user.UserIdentifiers, user)
		if err != nil {
			return err
		}

		return s.email.Send(user.UserIdentifiers.Email, &templates.PasswordReset{
			OrganizationName: s.config.OrganizationName,
			PublicURL:        s.config.PublicURL,
			Password:         user.Password,
		})
	})

	if err != nil {
		return nil, err
	}

	return &ttnpb.ResetUserPasswordResponse{
		Password: password,
	}, nil
}

// DeleteUser deletes an user.
func (s *adminService) DeleteUser(ctx context.Context, req *ttnpb.UserIdentifiers) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.store.Transact(func(tx *store.Store) error {
		ids := *req

		found, err := tx.Users.GetByID(ids, s.config.Specializers.User)
		if err != nil {
			return err
		}

		err = tx.Users.Delete(ids)
		if err != nil {
			return err
		}

		return s.email.Send(found.GetUser().UserIdentifiers.Email, &templates.AccountDeleted{
			UserID:           req.UserID,
			OrganizationName: s.config.OrganizationName,
			PublicURL:        s.config.PublicURL,
		})
	})

	return nil, err
}

// SendInvitation sends by email a token that can be used to create a new account.
// All invitations are expirable and the TTL is defined on a setitngs variable.
func (s *adminService) SendInvitation(ctx context.Context, req *ttnpb.SendInvitationRequest) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.store.Transact(func(tx *store.Store) (err error) {
		// check whether email is already registered or not
		found, err := tx.Users.GetByID(ttnpb.UserIdentifiers{Email: req.Email}, s.config.Specializers.User)
		if err != nil && !sql.ErrUserNotFound.Describes(err) {
			return err
		}

		// if email is already being used return error
		if found != nil {
			return ErrEmailAddressAlreadyUsed.New(nil)
		}

		// otherwise proceed to issue invitation
		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		now := time.Now()
		invitation := store.InvitationData{
			Token:     random.String(64),
			Email:     req.Email,
			IssuedAt:  now,
			ExpiresAt: now.Add(settings.InvitationTokenTTL),
		}

		err = tx.Invitations.Save(invitation)
		if err != nil {
			return err
		}

		return s.email.Send(req.Email, &templates.Invitation{
			OrganizationName: s.config.OrganizationName,
			PublicURL:        s.config.PublicURL,
			Token:            invitation.Token,
		})
	})

	return nil, err
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

	return nil, s.store.Invitations.Delete(req.Email)
}

// GetClient returns the client that matches the identifier.
func (s *adminService) GetClient(ctx context.Context, req *ttnpb.ClientIdentifiers) (*ttnpb.Client, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	found, err := s.store.Clients.GetByID(*req, s.config.Specializers.Client)
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

	found, err := s.store.Clients.List(s.config.Specializers.Client)
	if err != nil {
		return nil, err
	}

	resp := &ttnpb.ListClientsResponse{
		Clients: make([]*ttnpb.Client, 0, len(found)),
	}

	// filter results manually
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
		found, err := tx.Clients.GetByID(req.Client.ClientIdentifiers, s.config.Specializers.Client)
		if err != nil {
			return err
		}
		client := found.GetClient()

		for _, path := range req.UpdateMask.Paths {
			switch {
			case ttnpb.FieldPathClientDescription.MatchString(path):
				client.Description = req.Client.Description
			case ttnpb.FieldPathClientRedirectURI.MatchString(path):
				client.RedirectURI = req.Client.RedirectURI
			case ttnpb.FieldPathClientRights.MatchString(path):
				if req.Client.Rights == nil {
					req.Client.Rights = []ttnpb.Right{}
				}
				client.Rights = req.Client.Rights
			case ttnpb.FieldPathClientOfficialLabeled.MatchString(path):
				client.OfficialLabeled = req.Client.OfficialLabeled
			case ttnpb.FieldPathClientState.MatchString(path):
				client.State = req.Client.State
			case ttnpb.FieldPathClientGrants.MatchString(path):
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

	return nil, err
}

// DeleteClient deletes the client that matches the identifier and revokes all
// user authorizations.
func (s *adminService) DeleteClient(ctx context.Context, req *ttnpb.ClientIdentifiers) (*pbtypes.Empty, error) {
	err := s.enforceAdmin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.store.Transact(func(tx *store.Store) error {
		ids := *req

		found, err := tx.Clients.GetByID(ids, s.config.Specializers.Client)
		if err != nil {
			return err
		}

		user, err := tx.Users.GetByID(found.GetClient().CreatorIDs, s.config.Specializers.User)
		if err != nil {
			return err
		}

		err = tx.Clients.Delete(ids)
		if err != nil {
			return err
		}

		return s.email.Send(user.GetUser().UserIdentifiers.Email, &templates.ClientDeleted{
			ClientID:         req.ClientID,
			OrganizationName: s.config.OrganizationName,
			PublicURL:        s.config.PublicURL,
		})
	})

	return nil, err
}
