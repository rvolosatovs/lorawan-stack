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

	"github.com/gobwas/glob"
	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/auth"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/identityserver/email/templates"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/random"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// isEmailAllowed checks whether an input email is allowed given the glob list
// of allowed emails in the settings.
//
// Note this method was not placed on ttnpb as part of the IdentityServerSettings
// type as it makes use of an external package.
func isEmailAllowed(email string, allowedEmails []string) bool {
	if len(allowedEmails) == 0 {
		return true
	}

	found := false
	for i := range allowedEmails {
		found = glob.MustCompile(strings.ToLower(allowedEmails[i])).Match(strings.ToLower(email))
		if found {
			break
		}
	}

	return found
}

type userService struct {
	*IdentityServer
}

// CreateUser creates an user in the network.
func (s *userService) CreateUser(ctx context.Context, req *ttnpb.CreateUserRequest) (*pbtypes.Empty, error) {
	var user *ttnpb.User
	var token *store.ValidationToken

	err := s.store.Transact(func(tx *store.Store) error {
		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		// If invitation-only mode is enabled check that an invitation token is provided.
		if settings.InvitationOnly && req.InvitationToken == "" {
			return ErrInvitationTokenMissing.New(nil)
		}

		// check for blacklisted ids
		if !settings.IsIDAllowed(req.User.UserID) {
			return ErrBlacklistedID.New(errors.Attributes{
				"id": req.User.UserID,
			})
		}

		password, err := auth.Hash(req.User.Password)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		user = &ttnpb.User{
			UserIdentifiers:   req.User.UserIdentifiers,
			Name:              req.User.Name,
			Password:          string(password),
			State:             ttnpb.STATE_PENDING,
			PasswordUpdatedAt: now,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		if settings.SkipValidation {
			user.ValidatedAt = timeValue(now)
		}

		if !settings.AdminApproval {
			user.State = ttnpb.STATE_APPROVED
		}

		err = tx.Users.Create(user)
		if err != nil {
			return err
		}

		// check whether the provided email is allowed or not when an invitation token
		// wasn't provided
		if req.InvitationToken == "" {
			if !isEmailAllowed(req.User.UserIdentifiers.Email, settings.AllowedEmails) {
				return ErrEmailAddressNotAllowed.New(errors.Attributes{
					"allowed_emails": settings.AllowedEmails,
				})
			}
		} else {
			err = tx.Invitations.Use(req.User.UserIdentifiers.Email, req.InvitationToken)
			if err != nil {
				return err
			}
		}

		if !settings.SkipValidation {
			token = &store.ValidationToken{
				ValidationToken: random.String(64),
				CreatedAt:       now,
				ExpiresIn:       int32(settings.ValidationTokenTTL.Seconds()),
			}

			return tx.Users.SaveValidationToken(user.UserIdentifiers, *token)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// No email needs to be sent.
	if token == nil {
		return ttnpb.Empty, nil
	}

	return ttnpb.Empty, s.email.Send(user.UserIdentifiers.Email, &templates.EmailValidation{
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
		Token:            token.ValidationToken,
	})
}

// GetUser returns the account of the current user.
func (s *userService) GetUser(ctx context.Context, _ *pbtypes.Empty) (*ttnpb.User, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_INFO)
	if err != nil {
		return nil, err
	}

	found, err := s.store.Users.GetByID(authorizationDataFromContext(ctx).UserIdentifiers(), s.specializers.User)
	if err != nil {
		return nil, err
	}

	user := found.GetUser()
	user.Password = ""

	return user, nil
}

// UpdateUser updates the account of the current user.
// If email address is updated it sends an email to validate it if and only if
// the `SkipValidation` setting is disabled.
func (s *userService) UpdateUser(ctx context.Context, req *ttnpb.UpdateUserRequest) (*pbtypes.Empty, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_SETTINGS_BASIC)
	if err != nil {
		return nil, err
	}

	var user *ttnpb.User
	var token *store.ValidationToken

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Users.GetByID(authorizationDataFromContext(ctx).UserIdentifiers(), s.specializers.User)
		if err != nil {
			return err
		}
		user = found.GetUser()

		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		now := time.Now().UTC()

		newEmail := false
		for _, path := range req.UpdateMask.Paths {
			switch {
			case ttnpb.FieldPathUserName.MatchString(path):
				user.Name = req.User.Name
			case ttnpb.FieldPathUserEmail.MatchString(path):
				if !isEmailAllowed(req.User.UserIdentifiers.Email, settings.AllowedEmails) {
					return ErrEmailAddressNotAllowed.New(errors.Attributes{
						"allowed_emails": settings.AllowedEmails,
					})
				}

				newEmail = strings.ToLower(user.UserIdentifiers.Email) != strings.ToLower(req.User.UserIdentifiers.Email)
				if newEmail {
					if settings.SkipValidation {
						user.ValidatedAt = timeValue(now)
					} else {
						user.ValidatedAt = timeValue(time.Time{})
					}
				}

				user.UserIdentifiers.Email = req.User.UserIdentifiers.Email
			default:
				return ttnpb.ErrInvalidPathUpdateMask.New(errors.Attributes{
					"path": path,
				})
			}
		}

		user.UpdatedAt = now
		err = tx.Users.Update(authorizationDataFromContext(ctx).UserIdentifiers(), user)
		if err != nil {
			return err
		}

		if !newEmail || (newEmail && settings.SkipValidation) {
			return nil
		}

		token = &store.ValidationToken{
			ValidationToken: random.String(64),
			CreatedAt:       now,
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

// UpdateUserPassword updates the password of the current user.
func (s *userService) UpdateUserPassword(ctx context.Context, req *ttnpb.UpdateUserPasswordRequest) (*pbtypes.Empty, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_SETTINGS_BASIC)
	if err != nil {
		return nil, err
	}

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Users.GetByID(authorizationDataFromContext(ctx).UserIdentifiers(), s.specializers.User)
		if err != nil {
			return err
		}
		user := found.GetUser()

		matches, err := auth.Password(user.Password).Validate(req.Old)
		if err != nil {
			return err
		}

		if !matches {
			return ErrInvalidPassword.New(nil)
		}

		hashed, err := auth.Hash(req.New)
		if err != nil {
			return err
		}

		user.Password = string(hashed)
		user.PasswordUpdatedAt = time.Now().UTC()
		user.RequirePasswordUpdate = false

		return tx.Users.Update(user.UserIdentifiers, user)
	})

	return ttnpb.Empty, err
}

// DeleteUser deletes the account of the current user.
func (s *userService) DeleteUser(ctx context.Context, _ *pbtypes.Empty) (*pbtypes.Empty, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_DELETE)
	if err != nil {
		return nil, err
	}

	uIDs := authorizationDataFromContext(ctx).UserIdentifiers()

	err = s.store.Transact(func(tx *store.Store) error {
		apps, err := tx.Applications.ListByOrganizationOrUser(organizationOrUserIDsUserIDs(uIDs), s.specializers.Application)
		if err != nil {
			return err
		}

		gtws, err := tx.Gateways.ListByOrganizationOrUser(organizationOrUserIDsUserIDs(uIDs), s.specializers.Gateway)
		if err != nil {
			return err
		}

		orgs, err := tx.Organizations.ListByUser(uIDs, s.specializers.Organization)
		if err != nil {
			return err
		}

		err = tx.Users.Delete(uIDs)
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

	return ttnpb.Empty, err
}

// GenerateUserAPIKey generates an user API key and returns it.
func (s *userService) GenerateUserAPIKey(ctx context.Context, req *ttnpb.GenerateUserAPIKeyRequest) (*ttnpb.APIKey, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_SETTINGS_API_KEYS)
	if err != nil {
		return nil, err
	}

	k, err := auth.GenerateUserAPIKey(s.config.Hostname)
	if err != nil {
		return nil, err
	}

	key := ttnpb.APIKey{
		Key:    k,
		Name:   req.Name,
		Rights: req.Rights,
	}

	err = s.store.Users.SaveAPIKey(authorizationDataFromContext(ctx).UserIdentifiers(), key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

// ListUserAPIKeys returns all the API keys from the current user.
func (s *userService) ListUserAPIKeys(ctx context.Context, _ *pbtypes.Empty) (*ttnpb.ListUserAPIKeysResponse, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_SETTINGS_API_KEYS)
	if err != nil {
		return nil, err
	}

	found, err := s.store.Users.ListAPIKeys(authorizationDataFromContext(ctx).UserIdentifiers())
	if err != nil {
		return nil, err
	}

	keys := make([]*ttnpb.APIKey, 0, len(found))
	for i := range found {
		keys = append(keys, &found[i])
	}

	return &ttnpb.ListUserAPIKeysResponse{
		APIKeys: keys,
	}, nil
}

// UpdateUserAPIKey updates an API key from the current user.
func (s *userService) UpdateUserAPIKey(ctx context.Context, req *ttnpb.UpdateUserAPIKeyRequest) (*pbtypes.Empty, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_SETTINGS_BASIC)
	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.store.Users.UpdateAPIKeyRights(authorizationDataFromContext(ctx).UserIdentifiers(), req.Name, req.Rights)
}

// RemoveUserAPIKey removes an API key from the current user.
func (s *userService) RemoveUserAPIKey(ctx context.Context, req *ttnpb.RemoveUserAPIKeyRequest) (*pbtypes.Empty, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_SETTINGS_API_KEYS)
	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.store.Users.DeleteAPIKey(authorizationDataFromContext(ctx).UserIdentifiers(), req.Name)
}

// ValidateUserEmail validates the user's email with the token sent to the email.
func (s *userService) ValidateUserEmail(ctx context.Context, req *ttnpb.ValidateUserEmailRequest) (*pbtypes.Empty, error) {
	err := s.store.Transact(func(tx *store.Store) error {
		userID, token, err := tx.Users.GetValidationToken(req.Token)
		if err != nil {
			return err
		}

		if token.IsExpired() {
			return ErrValidationTokenExpired.New(nil)
		}

		user, err := tx.Users.GetByID(userID, s.specializers.User)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		user.GetUser().ValidatedAt = timeValue(now)
		user.GetUser().UpdatedAt = now

		err = tx.Users.Update(user.GetUser().UserIdentifiers, user)
		if err != nil {
			return err
		}

		return tx.Users.DeleteValidationToken(req.Token)
	})

	return ttnpb.Empty, err
}

// RequestUserEmailValidation requests a new validation email if the user's email
// isn't validated yet.
func (s *userService) RequestUserEmailValidation(ctx context.Context, _ *pbtypes.Empty) (*pbtypes.Empty, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_SETTINGS_BASIC)
	if err != nil {
		return nil, err
	}

	var user *ttnpb.User
	var token store.ValidationToken

	err = s.store.Transact(func(tx *store.Store) error {
		found, err := tx.Users.GetByID(authorizationDataFromContext(ctx).UserIdentifiers(), s.specializers.User)
		if err != nil {
			return err
		}
		user = found.GetUser()

		if !user.ValidatedAt.IsZero() {
			return ErrEmailAlreadyValidated.New(nil)
		}

		settings, err := tx.Settings.Get()
		if err != nil {
			return err
		}

		token = store.ValidationToken{
			ValidationToken: random.String(64),
			CreatedAt:       time.Now().UTC(),
			ExpiresIn:       int32(settings.ValidationTokenTTL.Seconds()),
		}

		return tx.Users.SaveValidationToken(authorizationDataFromContext(ctx).UserIdentifiers(), token)
	})

	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.email.Send(user.UserIdentifiers.Email, &templates.EmailValidation{
		OrganizationName: s.config.OrganizationName,
		PublicURL:        s.config.PublicURL,
		Token:            token.ValidationToken,
	})
}

// ListAuthorizedClients returns all the authorized third-party clients that
// the current user has.
func (s *userService) ListAuthorizedClients(ctx context.Context, _ *pbtypes.Empty) (*ttnpb.ListAuthorizedClientsResponse, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_AUTHORIZED_CLIENTS)
	if err != nil {
		return nil, err
	}

	found, err := s.store.OAuth.ListAuthorizedClients(authorizationDataFromContext(ctx).UserIdentifiers(), s.specializers.Client)
	if err != nil {
		return nil, err
	}

	resp := &ttnpb.ListAuthorizedClientsResponse{
		Clients: make([]*ttnpb.Client, 0, len(found)),
	}

	for _, client := range found {
		cli := client.GetClient()
		cli.Secret = ""
		cli.RedirectURI = ""
		cli.Grants = nil
		resp.Clients = append(resp.Clients, cli)
	}

	return resp, nil
}

// RevokeAuthorizedClient revokes an authorized third-party client.
func (s *userService) RevokeAuthorizedClient(ctx context.Context, req *ttnpb.ClientIdentifiers) (*pbtypes.Empty, error) {
	err := s.enforceUserRights(ctx, ttnpb.RIGHT_USER_AUTHORIZED_CLIENTS)
	if err != nil {
		return nil, err
	}

	return ttnpb.Empty, s.store.OAuth.RevokeAuthorizedClient(authorizationDataFromContext(ctx).UserIdentifiers(), *req)
}
