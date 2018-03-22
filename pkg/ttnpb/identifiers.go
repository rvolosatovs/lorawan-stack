// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package ttnpb

import "github.com/TheThingsNetwork/ttn/pkg/validate"

// Validate is used as validator function by the GRPC validator interceptor.
func (i UserIdentifiers) Validate() error {
	return validate.Field(i.UserID, validate.ID).DescribeFieldName("User ID")
}

// IsZero returns true if all identifiers have zero-values.
func (i UserIdentifiers) IsZero() bool {
	return i.UserID == ""
}

// Equals returns true if the receiver identifiers matches to other identifiers.
func (i UserIdentifiers) Equals(other UserIdentifiers) bool {
	return i.UserID == other.UserID
}

// Contains returns true if other is contained in the receiver.
func (i UserIdentifiers) Contains(other UserIdentifiers) bool {
	return i.UserID == other.UserID
}

// Validate is used as validator function by the GRPC validator interceptor.
func (i ApplicationIdentifiers) Validate() error {
	return validate.Field(i.ApplicationID, validate.ID).DescribeFieldName("Application ID")
}

// Contains returns true if other is contained in the receiver.
func (i ApplicationIdentifiers) Contains(other ApplicationIdentifiers) bool {
	return i.ApplicationID == other.ApplicationID
}

// IsZero returns true if all identifiers have zero-values.
func (i ApplicationIdentifiers) IsZero() bool {
	return i.ApplicationID == ""
}

// Validate is used as validator function by the GRPC validator interceptor.
func (i GatewayIdentifiers) Validate() error {
	return validate.Field(i.GatewayID, validate.ID).DescribeFieldName("Gateway ID")
}

// IsZero returns true if all identifiers have zero-values.
func (i GatewayIdentifiers) IsZero() bool {
	return i.GatewayID == ""
}

// Contains returns true if other is contained in the receiver.
func (i GatewayIdentifiers) Contains(other GatewayIdentifiers) bool {
	return i.GatewayID == other.GatewayID
}

// Validate is used as validator function by the GRPC validator interceptor.
func (i ClientIdentifiers) Validate() error {
	return validate.Field(i.ClientID, validate.ID).DescribeFieldName("Client ID")
}

// IsZero returns true if all identifiers have zero-values.
func (i ClientIdentifiers) IsZero() bool {
	return i.ClientID == ""
}

// Validate is used as validator function by the GRPC validator interceptor.
func (i OrganizationIdentifiers) Validate() error {
	return validate.Field(i.OrganizationID, validate.ID).DescribeFieldName("Organization ID")
}

// IsZero returns true if all identifiers have zero-values.
func (i OrganizationIdentifiers) IsZero() bool {
	return i.OrganizationID == ""
}

// Contains returns true if other is contained in the receiver.
func (i OrganizationIdentifiers) Contains(other OrganizationIdentifiers) bool {
	return i.OrganizationID == other.OrganizationID
}

// Validate is used as validator function by the GRPC validator interceptor.
func (i OrganizationOrUserIdentifiers) Validate() error {
	if id := i.GetUserID(); id != nil {
		return id.Validate()
	}

	if id := i.GetOrganizationID(); id != nil {
		return id.Validate()
	}

	return ErrEmptyIdentifiers.New(nil)
}
