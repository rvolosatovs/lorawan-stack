// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

// Copyright © 2017 The Things Industries B.V.

package utils

import "github.com/TheThingsNetwork/ttn/pkg/identityserver/types"

// Collaborator is a helper to construct a collaborator type
func Collaborator(username string, rights []types.Right) types.Collaborator {
	c := types.Collaborator{
		Username: username,
		Rights:   rights,
	}
	return c
}
