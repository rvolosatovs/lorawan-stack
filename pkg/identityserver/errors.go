// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package identityserver

import "github.com/TheThingsNetwork/ttn/pkg/errors"

func init() {
	ErrBlacklistedID.Register()
	ErrEmailAddressNotAllowed.Register()
	ErrInvalidPassword.Register()
	ErrNotAuthorized.Register()
	ErrSetApplicationCollaboratorFailed.Register()
	ErrSetGatewayCollaboratorFailed.Register()
}

// ErrBlacklistedID is returned when trying to register an entity using a blacklisted ID.
var ErrBlacklistedID = &errors.ErrDescriptor{
	MessageFormat: "`{id} is not an allowed ID",
	Code:          1,
	Type:          errors.InvalidArgument,
}

// ErrEmailAddressNotAllowed is returned when trying to set in a user account an
// email address that is not allowed
var ErrEmailAddressNotAllowed = &errors.ErrDescriptor{
	MessageFormat: "Email address `{email}` is not allowed. Current allowed email addresses by the network administrators are in the following form: {allowed_emails}",
	Code:          2,
	Type:          errors.InvalidArgument,
}

// ErrInvalidPassword is returned when trying to update the user's account
// password while providing a wrong current password.
var ErrInvalidPassword = &errors.ErrDescriptor{
	MessageFormat: "Invalid password",
	Code:          3,
	Type:          errors.Unauthorized,
}

// ErrNotAuthorized is returned when trying to access a protected resource without
// authorization.
var ErrNotAuthorized = &errors.ErrDescriptor{
	MessageFormat: "Not authorized",
	Code:          4,
	Type:          errors.Unauthorized,
}

// ErrSetApplicationCollaboratorFailed is returned when after modifying the
// collaborators of an application there is no collaborator with
// `RIGHT_APPLICATION_SETTINGS_COLLABORATORS` right.
var ErrSetApplicationCollaboratorFailed = &errors.ErrDescriptor{
	MessageFormat: "Failed to modify collaborators: application `{application_id}` must have at least one collaborator with `RIGHT_APPLICATION_SETTINGS_COLLABORATORS` right",
	Code:          5,
	Type:          errors.Conflict,
}

// ErrSetGatewayCollaboratorFailed is returned when after modifying the
// collaborators of a gateway there is no collaborator with
// `RIGHT_GATEWAY_SETTINGS_COLLABORATORS` right.
var ErrSetGatewayCollaboratorFailed = &errors.ErrDescriptor{
	MessageFormat: "Failed to modify collaborators: gateway `{gateway_id}` must have at least one collaborator with `RIGHT_GATEWAY_SETTINGS_COLLABORATORS` right",
	Code:          6,
	Type:          errors.Conflict,
}
