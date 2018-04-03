// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package ttnpb

import (
	"crypto/subtle"
	"regexp"
)

// GetClient returns the base Client itself.
func (c *Client) GetClient() *Client {
	return c
}

// GetId implements osin.Client.
// nolint: golint
func (c *Client) GetId() string {
	return c.ClientIdentifiers.GetClientID()
}

// GetRedirectUri implements osin.Client.
// nolint: golint
func (c *Client) GetRedirectUri() string {
	return c.RedirectURI
}

// GetUserData implements osin.Client.
func (c *Client) GetUserData() interface{} {
	return nil
}

// ClientSecretMatches implements osin.ClientSecretMatcher.
func (c *Client) ClientSecretMatches(secret string) bool {
	return subtle.ConstantTimeEq(int32(len(c.Secret)), int32(len(secret))) == 1 && subtle.ConstantTimeCompare([]byte(c.Secret), []byte(secret)) == 1
}

// HasGrant checks whether the client has a given grant or not.
func (c *Client) HasGrant(grant GrantType) bool {
	for _, g := range c.Grants {
		if g == grant {
			return true
		}
	}

	return false
}

var (
	// FieldPathClientDescription is the field path for the client description field.
	FieldPathClientDescription = regexp.MustCompile(`^description$`)

	// FieldPathClientRedirectURI is the field path for the client redirect URI field.
	FieldPathClientRedirectURI = regexp.MustCompile(`^redirect_uri$`)

	// FieldPathClientRights is the field path for the client rights field.
	FieldPathClientRights = regexp.MustCompile(`^rights$`)

	// FieldPathClientOfficialLabeled is the field path for the client official labeled field.
	FieldPathClientOfficialLabeled = regexp.MustCompile(`^official_labeled$`)

	// FieldPathClientState is the field path for the client state field.
	FieldPathClientState = regexp.MustCompile(`^state$`)

	// FieldPathClientGrants is the field path for the client grants field.
	FieldPathClientGrants = regexp.MustCompile(`^grants$`)
)
