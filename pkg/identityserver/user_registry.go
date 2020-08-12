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
	"runtime/trace"
	"strings"
	"time"
	"unicode"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/auth"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/email"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/blacklist"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/emails"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"go.thethings.network/lorawan-stack/v3/pkg/validate"
)

var (
	evtCreateUser = events.Define(
		"user.create", "create user",
		events.WithVisibility(ttnpb.RIGHT_USER_INFO),
		events.WithAuthFromContext(),
	)
	evtUpdateUser = events.Define(
		"user.update", "update user",
		events.WithVisibility(ttnpb.RIGHT_USER_INFO),
		events.WithAuthFromContext(),
	)
	evtDeleteUser = events.Define(
		"user.delete", "delete user",
		events.WithVisibility(ttnpb.RIGHT_USER_INFO),
		events.WithAuthFromContext(),
	)
	evtUpdateUserIncorrectPassword = events.Define(
		"user.update.incorrect_password", "update user failure: incorrect password",
		events.WithVisibility(ttnpb.RIGHT_USER_INFO),
		events.WithAuthFromContext(),
	)
)

var (
	errInvitationTokenRequired   = errors.DefineInvalidArgument("invitation_token_required", "invitation token required")
	errInvitationTokenExpired    = errors.DefineInvalidArgument("invitation_token_expired", "invitation token expired")
	errPasswordStrengthMinLength = errors.DefineInvalidArgument("password_strength_min_length", "need at least `{n}` characters")
	errPasswordStrengthMaxLength = errors.DefineInvalidArgument("password_strength_max_length", "need at most `{n}` characters")
	errPasswordStrengthUppercase = errors.DefineInvalidArgument("password_strength_uppercase", "need at least `{n}` uppercase letter(s)")
	errPasswordStrengthDigits    = errors.DefineInvalidArgument("password_strength_digits", "need at least `{n}` digit(s)")
	errPasswordStrengthSpecial   = errors.DefineInvalidArgument("password_strength_special", "need at least `{n}` special character(s)")
)

func (is *IdentityServer) validatePasswordStrength(ctx context.Context, password string) error {
	requirements := is.configFromContext(ctx).UserRegistration.PasswordRequirements
	if len(password) < requirements.MinLength {
		return errPasswordStrengthMinLength.WithAttributes("n", requirements.MinLength)
	}
	if len(password) > requirements.MaxLength {
		return errPasswordStrengthMaxLength.WithAttributes("n", requirements.MaxLength)
	}
	var uppercase, digits, special int
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			uppercase++
		case unicode.IsDigit(r):
			digits++
		case !unicode.IsLetter(r) && !unicode.IsNumber(r):
			special++
		}
	}
	if uppercase < requirements.MinUppercase {
		return errPasswordStrengthUppercase.WithAttributes("n", requirements.MinUppercase)
	}
	if digits < requirements.MinDigits {
		return errPasswordStrengthDigits.WithAttributes("n", requirements.MinDigits)
	}
	if special < requirements.MinSpecial {
		return errPasswordStrengthSpecial.WithAttributes("n", requirements.MinSpecial)
	}
	return nil
}

func (is *IdentityServer) createUser(ctx context.Context, req *ttnpb.CreateUserRequest) (usr *ttnpb.User, err error) {
	createdByAdmin := is.IsAdmin(ctx)

	if err = blacklist.Check(ctx, req.UserID); err != nil {
		return nil, err
	}
	if req.InvitationToken == "" && is.configFromContext(ctx).UserRegistration.Invitation.Required && !createdByAdmin {
		return nil, errInvitationTokenRequired.New()
	}

	if err := validate.Email(req.User.PrimaryEmailAddress); err != nil {
		return nil, err
	}
	if err := validateContactInfo(req.User.ContactInfo); err != nil {
		return nil, err
	}

	if !createdByAdmin {
		req.User.PrimaryEmailAddressValidatedAt = nil
		req.User.RequirePasswordUpdate = false
		if is.configFromContext(ctx).UserRegistration.AdminApproval.Required {
			req.User.State = ttnpb.STATE_REQUESTED
		} else {
			req.User.State = ttnpb.STATE_APPROVED
		}
		req.User.Admin = false
		req.User.TemporaryPassword = ""
		req.User.TemporaryPasswordCreatedAt = nil
		req.User.TemporaryPasswordExpiresAt = nil
		cleanContactInfo(req.User.ContactInfo)
	}

	var primaryEmailAddressFound bool
	for _, contactInfo := range req.User.ContactInfo {
		if contactInfo.ContactMethod == ttnpb.CONTACT_METHOD_EMAIL && contactInfo.Value == req.User.PrimaryEmailAddress {
			primaryEmailAddressFound = true
			if contactInfo.ValidatedAt != nil {
				req.User.PrimaryEmailAddressValidatedAt = contactInfo.ValidatedAt
				break
			}
		}
	}
	if !primaryEmailAddressFound {
		req.User.ContactInfo = append(req.User.ContactInfo, &ttnpb.ContactInfo{
			ContactMethod: ttnpb.CONTACT_METHOD_EMAIL,
			Value:         req.User.PrimaryEmailAddress,
			ValidatedAt:   req.User.PrimaryEmailAddressValidatedAt,
		})
	}

	if err := is.validatePasswordStrength(ctx, req.User.Password); err != nil {
		return nil, err
	}
	hashedPassword, err := auth.Hash(ctx, req.User.Password)
	if err != nil {
		return nil, err
	}
	req.User.Password = hashedPassword
	now := time.Now()
	req.User.PasswordUpdatedAt = &now

	if req.User.ProfilePicture != nil {
		if err = is.processUserProfilePicture(ctx, &req.User); err != nil {
			return nil, err
		}
	}
	defer func() { is.setFullProfilePictureURL(ctx, usr) }()

	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		if req.InvitationToken != "" {
			invitationToken, err := store.GetInvitationStore(db).GetInvitation(ctx, req.InvitationToken)
			if err != nil {
				return err
			}
			if !invitationToken.ExpiresAt.IsZero() && invitationToken.ExpiresAt.Before(time.Now()) {
				return errInvitationTokenExpired.New()
			}
		}

		usr, err = store.GetUserStore(db).CreateUser(ctx, &req.User)
		if err != nil {
			return err
		}

		if len(req.ContactInfo) > 0 {
			usr.ContactInfo, err = store.GetContactInfoStore(db).SetContactInfo(ctx, usr.UserIdentifiers, req.ContactInfo)
			if err != nil {
				return err
			}
		}

		if req.InvitationToken != "" {
			if err = store.GetInvitationStore(db).SetInvitationAcceptedBy(ctx, req.InvitationToken, &usr.UserIdentifiers); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if usr.State == ttnpb.STATE_REQUESTED {
		err = is.SendAdminsEmail(ctx, func(data emails.Data) email.MessageData {
			data.Entity.Type, data.Entity.ID = "user", usr.UserID
			return &emails.UserRequested{
				Data: data,
			}
		})
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("Could not send user requested email")
		}
	}

	// TODO: Send welcome email (https://github.com/TheThingsNetwork/lorawan-stack/issues/72).

	if _, err := is.requestContactInfoValidation(ctx, req.UserIdentifiers.EntityIdentifiers()); err != nil {
		log.FromContext(ctx).WithError(err).Error("Could not send contact info validations")
	}

	usr.Password = "" // Create doesn't have a FieldMask, so we need to manually remove the password.
	events.Publish(evtCreateUser.NewWithIdentifiersAndData(ctx, req.UserIdentifiers, nil))
	return usr, nil
}

func (is *IdentityServer) getUser(ctx context.Context, req *ttnpb.GetUserRequest) (usr *ttnpb.User, err error) {
	req.FieldMask.Paths = cleanFieldMaskPaths(ttnpb.UserFieldPathsNested, req.FieldMask.Paths, getPaths, nil)
	if err = rights.RequireUser(ctx, req.UserIdentifiers, ttnpb.RIGHT_USER_INFO); err != nil {
		if err := is.RequireAuthenticated(ctx); err != nil {
			return nil, err
		}
		if ttnpb.HasOnlyAllowedFields(req.FieldMask.Paths, ttnpb.PublicUserFields...) {
			defer func() { usr = usr.PublicSafe() }()
		} else {
			return nil, err
		}
	}

	if ttnpb.HasAnyField(ttnpb.TopLevelFields(req.FieldMask.Paths), "profile_picture") {
		if is.configFromContext(ctx).ProfilePicture.UseGravatar {
			if !ttnpb.HasAnyField(req.FieldMask.Paths, "primary_email_address") {
				req.FieldMask.Paths = append(req.FieldMask.Paths, "primary_email_address")
				defer func() {
					if usr != nil {
						usr.PrimaryEmailAddress = ""
					}
				}()
			}
			defer func() { fillGravatar(ctx, usr) }()
		}
		defer func() { is.setFullProfilePictureURL(ctx, usr) }()
	}

	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		usr, err = store.GetUserStore(db).GetUser(ctx, &req.UserIdentifiers, &req.FieldMask)
		if err != nil {
			return err
		}
		if ttnpb.HasAnyField(req.FieldMask.Paths, "contact_info") {
			usr.ContactInfo, err = store.GetContactInfoStore(db).GetContactInfo(ctx, usr.UserIdentifiers)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (is *IdentityServer) listUsers(ctx context.Context, req *ttnpb.ListUsersRequest) (users *ttnpb.Users, err error) {
	req.FieldMask.Paths = cleanFieldMaskPaths(ttnpb.UserFieldPathsNested, req.FieldMask.Paths, getPaths, nil)
	if err = is.RequireAdmin(ctx); err != nil {
		return nil, err
	}
	ctx = store.WithOrder(ctx, req.Order)
	var total uint64
	paginateCtx := store.WithPagination(ctx, req.Limit, req.Page, &total)
	defer func() {
		if err == nil {
			setTotalHeader(ctx, total)
		}
	}()
	users = &ttnpb.Users{}
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		users.Users, err = store.GetUserStore(db).FindUsers(paginateCtx, nil, &req.FieldMask)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

var (
	errUpdateUserPasswordRequest = errors.DefineInvalidArgument("password_in_update", "can not update password with regular user update request")
	errUpdateUserAdminField      = errors.DefinePermissionDenied("user_update_admin_field", "only admins can update the `{field}` field")
)

func (is *IdentityServer) setFullProfilePictureURL(ctx context.Context, usr *ttnpb.User) {
	bucketURL := is.configFromContext(ctx).ProfilePicture.BucketURL
	if bucketURL == "" {
		return
	}
	bucketURL = strings.TrimSuffix(bucketURL, "/") + "/"
	if usr != nil && usr.ProfilePicture != nil {
		for size, file := range usr.ProfilePicture.Sizes {
			if !strings.Contains(file, "://") {
				usr.ProfilePicture.Sizes[size] = bucketURL + strings.TrimPrefix(file, "/")
			}
		}
	}
}

func (is *IdentityServer) updateUser(ctx context.Context, req *ttnpb.UpdateUserRequest) (usr *ttnpb.User, err error) {
	if err = rights.RequireUser(ctx, req.UserIdentifiers, ttnpb.RIGHT_USER_SETTINGS_BASIC); err != nil {
		return nil, err
	}
	req.FieldMask.Paths = cleanFieldMaskPaths(ttnpb.UserFieldPathsNested, req.FieldMask.Paths, nil, getPaths)
	if len(req.FieldMask.Paths) == 0 {
		req.FieldMask.Paths = updatePaths
	}
	updatedByAdmin := is.IsAdmin(ctx)

	if ttnpb.HasAnyField(req.FieldMask.Paths, "primary_email_address") {
		if err := validate.Email(req.User.PrimaryEmailAddress); err != nil {
			return nil, err
		}
	}
	if err := validateContactInfo(req.User.ContactInfo); err != nil {
		return nil, err
	}

	if !updatedByAdmin {
		for _, path := range req.FieldMask.Paths {
			switch path {
			case "primary_email_address_validated_at",
				"require_password_update",
				"state", "admin",
				"temporary_password", "temporary_password_created_at", "temporary_password_expires_at":
				return nil, errUpdateUserAdminField.WithAttributes("field", path)
			}
		}
		req.PrimaryEmailAddressValidatedAt = nil
		cleanContactInfo(req.User.ContactInfo)
	}

	if ttnpb.HasAnyField(req.FieldMask.Paths, "temporary_password") {
		hashedTemporaryPassword, err := auth.Hash(ctx, req.User.TemporaryPassword)
		if err != nil {
			return nil, err
		}
		req.User.TemporaryPassword = hashedTemporaryPassword
		now := time.Now()
		if !ttnpb.HasAnyField(req.FieldMask.Paths, "temporary_password_created_at") {
			req.User.TemporaryPasswordCreatedAt = &now
			req.FieldMask.Paths = append(req.FieldMask.Paths, "temporary_password_created_at")
		}
		if !ttnpb.HasAnyField(req.FieldMask.Paths, "temporary_password_expires_at") {
			expires := now.Add(36 * time.Hour)
			req.User.TemporaryPasswordExpiresAt = &expires
			req.FieldMask.Paths = append(req.FieldMask.Paths, "temporary_password_expires_at")
		}
	}

	if ttnpb.HasAnyField(ttnpb.TopLevelFields(req.FieldMask.Paths), "profile_picture") {
		if !ttnpb.HasAnyField(req.FieldMask.Paths, "profile_picture") {
			req.FieldMask.Paths = append(req.FieldMask.Paths, "profile_picture")
		}
		if req.User.ProfilePicture != nil {
			if err = is.processUserProfilePicture(ctx, &req.User); err != nil {
				return nil, err
			}
		}
		defer func() { is.setFullProfilePictureURL(ctx, usr) }()
	}

	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		updatingContactInfo := ttnpb.HasAnyField(req.FieldMask.Paths, "contact_info")
		var contactInfo []*ttnpb.ContactInfo
		updatingPrimaryEmailAddress := ttnpb.HasAnyField(req.FieldMask.Paths, "primary_email_address")
		if updatingContactInfo || updatingPrimaryEmailAddress {
			if updatingContactInfo {
				contactInfo, err = store.GetContactInfoStore(db).SetContactInfo(ctx, req.User.UserIdentifiers, req.ContactInfo)
				if err != nil {
					return err
				}
			}
			if updatingPrimaryEmailAddress {
				if !updatingContactInfo {
					contactInfo, err = store.GetContactInfoStore(db).GetContactInfo(ctx, req.User.UserIdentifiers)
					if err != nil {
						return err
					}
				}
				if !ttnpb.HasAnyField(req.FieldMask.Paths, "primary_email_address_validated_at") {
					for _, contactInfo := range contactInfo {
						if contactInfo.ContactMethod == ttnpb.CONTACT_METHOD_EMAIL && contactInfo.Value == req.User.PrimaryEmailAddress {
							req.PrimaryEmailAddressValidatedAt = contactInfo.ValidatedAt
							req.FieldMask.Paths = append(req.FieldMask.Paths, "primary_email_address_validated_at")
							break
						}
					}
				}
			}
		}
		usr, err = store.GetUserStore(db).UpdateUser(ctx, &req.User, &req.FieldMask)
		if err != nil {
			return err
		}
		if updatingContactInfo {
			usr.ContactInfo = contactInfo
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtUpdateUser.NewWithIdentifiersAndData(ctx, req.UserIdentifiers, req.FieldMask.Paths))

	// TODO: Send emails (https://github.com/TheThingsNetwork/lorawan-stack/issues/72).
	// - If primary email address changed
	if ttnpb.HasAnyField(req.FieldMask.Paths, "state") {
		err = is.SendUserEmail(ctx, &req.UserIdentifiers, func(data emails.Data) email.MessageData {
			data.SetEntity(req.EntityIdentifiers())
			return &emails.EntityStateChanged{Data: data, State: strings.ToLower(strings.TrimPrefix(usr.State.String(), "STATE_"))}
		})
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("Could not send state change notification email")
		}
	}

	return usr, nil
}

var (
	errIncorrectPassword        = errors.DefineUnauthenticated("old_password", "incorrect old password")
	errTemporaryPasswordExpired = errors.DefineUnauthenticated("temporary_password_expired", "temporary password expired")
)

var (
	updatePasswordFieldMask = &types.FieldMask{Paths: []string{
		"password", "password_updated_at", "require_password_update",
	}}
	temporaryPasswordFieldMask = &types.FieldMask{Paths: []string{
		"password", "password_updated_at", "require_password_update",
		"temporary_password", "temporary_password_created_at", "temporary_password_expires_at",
	}}
	updateTemporaryPasswordFieldMask = &types.FieldMask{Paths: []string{
		"temporary_password", "temporary_password_created_at", "temporary_password_expires_at",
	}}
)

func (is *IdentityServer) updateUserPassword(ctx context.Context, req *ttnpb.UpdateUserPasswordRequest) (*types.Empty, error) {
	if err := is.validatePasswordStrength(ctx, req.New); err != nil {
		return nil, err
	}
	hashedPassword, err := auth.Hash(ctx, req.New)
	if err != nil {
		return nil, err
	}
	updateMask := updatePasswordFieldMask
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		usr, err := store.GetUserStore(db).GetUser(ctx, &req.UserIdentifiers, temporaryPasswordFieldMask)
		if err != nil {
			return err
		}
		region := trace.StartRegion(ctx, "validate old password")
		valid, err := auth.Validate(usr.Password, req.Old)
		region.End()
		if err != nil {
			return err
		}
		if valid {
			// TODO: Add when 2FA is enabled (https://github.com/TheThingsNetwork/lorawan-stack/issues/2)
			// if err := rights.RequireUser(ctx, req.UserIdentifiers, ttnpb.RIGHT_USER_ALL); err != nil {
			//	return err
			// }
		} else {
			if usr.TemporaryPassword == "" {
				events.Publish(evtUpdateUserIncorrectPassword.NewWithIdentifiersAndData(ctx, req.UserIdentifiers, nil))
				return errIncorrectPassword.New()
			}
			region := trace.StartRegion(ctx, "validate temporary password")
			valid, err = auth.Validate(usr.TemporaryPassword, req.Old)
			region.End()
			switch {
			case err != nil:
				return err
			case !valid:
				events.Publish(evtUpdateUserIncorrectPassword.NewWithIdentifiersAndData(ctx, req.UserIdentifiers, nil))
				return errIncorrectPassword.New()
			case usr.TemporaryPasswordExpiresAt.Before(time.Now()):
				events.Publish(evtUpdateUserIncorrectPassword.NewWithIdentifiersAndData(ctx, req.UserIdentifiers, nil))
				return errTemporaryPasswordExpired.New()
			}
			usr.TemporaryPassword, usr.TemporaryPasswordCreatedAt, usr.TemporaryPasswordExpiresAt = "", nil, nil
			updateMask = temporaryPasswordFieldMask
		}
		if req.RevokeAllAccess {
			sessionStore := store.GetUserSessionStore(db)
			sessions, err := sessionStore.FindSessions(ctx, &req.UserIdentifiers)
			if err != nil {
				return err
			}
			for _, session := range sessions {
				err = sessionStore.DeleteSession(ctx, &req.UserIdentifiers, session.SessionID)
				if err != nil {
					return err
				}
			}
			oauthStore := store.GetOAuthStore(db)
			authorizations, err := oauthStore.ListAuthorizations(ctx, &req.UserIdentifiers)
			if err != nil {
				return err
			}
			for _, auth := range authorizations {
				tokens, err := oauthStore.ListAccessTokens(ctx, &auth.UserIDs, &auth.ClientIDs)
				if err != nil {
					return err
				}
				for _, token := range tokens {
					err = oauthStore.DeleteAccessToken(ctx, token.ID)
					if err != nil {
						return err
					}
				}
			}
		}
		now := time.Now()
		usr.Password, usr.PasswordUpdatedAt, usr.RequirePasswordUpdate = hashedPassword, &now, false
		usr, err = store.GetUserStore(db).UpdateUser(ctx, usr, updateMask)
		return err
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtUpdateUser.NewWithIdentifiersAndData(ctx, req.UserIdentifiers, updateMask))
	err = is.SendUserEmail(ctx, &req.UserIdentifiers, func(data emails.Data) email.MessageData {
		return &emails.PasswordChanged{Data: data}
	})
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("Could not send password change notification email")
	}
	return ttnpb.Empty, nil
}

var errTemporaryPasswordStillValid = errors.DefineInvalidArgument("temporary_password_still_valid", "previous temporary password still valid")

func (is *IdentityServer) createTemporaryPassword(ctx context.Context, req *ttnpb.CreateTemporaryPasswordRequest) (*types.Empty, error) {
	temporaryPassword, err := auth.GenerateKey(ctx)
	if err != nil {
		return nil, err
	}
	hashedTemporaryPassword, err := auth.Hash(ctx, temporaryPassword)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	err = is.withDatabase(ctx, func(db *gorm.DB) error {
		usr, err := store.GetUserStore(db).GetUser(ctx, &req.UserIdentifiers, temporaryPasswordFieldMask)
		if err != nil {
			return err
		}
		if usr.TemporaryPasswordExpiresAt != nil && usr.TemporaryPasswordExpiresAt.After(time.Now()) {
			return errTemporaryPasswordStillValid.New()
		}
		usr.TemporaryPassword = hashedTemporaryPassword
		expires := now.Add(time.Hour)
		usr.TemporaryPasswordCreatedAt, usr.TemporaryPasswordExpiresAt = &now, &expires
		usr, err = store.GetUserStore(db).UpdateUser(ctx, usr, updateTemporaryPasswordFieldMask)
		return err
	})
	if err != nil {
		return nil, err
	}
	log.FromContext(ctx).WithFields(log.Fields(
		"user_uid", unique.ID(ctx, req.UserIdentifiers),
		"temporary_password", temporaryPassword,
	)).Info("Created temporary password")
	events.Publish(evtUpdateUser.NewWithIdentifiersAndData(ctx, req.UserIdentifiers, updateTemporaryPasswordFieldMask))
	err = is.SendUserEmail(ctx, &req.UserIdentifiers, func(data emails.Data) email.MessageData {
		return &emails.TemporaryPassword{
			Data:              data,
			TemporaryPassword: temporaryPassword,
		}
	})
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("Could not send temporary password email")
	}
	return ttnpb.Empty, nil
}

func (is *IdentityServer) deleteUser(ctx context.Context, ids *ttnpb.UserIdentifiers) (*types.Empty, error) {
	if err := rights.RequireUser(ctx, *ids, ttnpb.RIGHT_USER_DELETE); err != nil {
		return nil, err
	}
	err := is.withDatabase(ctx, func(db *gorm.DB) error {
		return store.GetUserStore(db).DeleteUser(ctx, ids)
	})
	if err != nil {
		return nil, err
	}
	events.Publish(evtDeleteUser.NewWithIdentifiersAndData(ctx, ids, nil))
	return ttnpb.Empty, nil
}

type userRegistry struct {
	*IdentityServer
}

func (ur *userRegistry) Create(ctx context.Context, req *ttnpb.CreateUserRequest) (*ttnpb.User, error) {
	return ur.createUser(ctx, req)
}

func (ur *userRegistry) List(ctx context.Context, req *ttnpb.ListUsersRequest) (*ttnpb.Users, error) {
	return ur.listUsers(ctx, req)
}

func (ur *userRegistry) Get(ctx context.Context, req *ttnpb.GetUserRequest) (*ttnpb.User, error) {
	return ur.getUser(ctx, req)
}

func (ur *userRegistry) Update(ctx context.Context, req *ttnpb.UpdateUserRequest) (*ttnpb.User, error) {
	return ur.updateUser(ctx, req)
}

func (ur *userRegistry) UpdatePassword(ctx context.Context, req *ttnpb.UpdateUserPasswordRequest) (*types.Empty, error) {
	return ur.updateUserPassword(ctx, req)
}

func (ur *userRegistry) CreateTemporaryPassword(ctx context.Context, req *ttnpb.CreateTemporaryPasswordRequest) (*types.Empty, error) {
	return ur.createTemporaryPassword(ctx, req)
}

func (ur *userRegistry) Delete(ctx context.Context, req *ttnpb.UserIdentifiers) (*types.Empty, error) {
	return ur.deleteUser(ctx, req)
}
