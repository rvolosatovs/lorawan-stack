// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package api

import (
	"context"

	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	pbtypes "github.com/gogo/protobuf/types"
)

var _ ttnpb.IsSettingsServer = new(GRPC)

// GetSettings fetches the current dynamic settings of the Identity Server.
func (g *GRPC) GetSettings(ctx context.Context, _ *pbtypes.Empty) (*ttnpb.IdentityServerSettings, error) {
	err := g.adminCheck(ctx)
	if err != nil {
		return nil, err
	}

	settings, err := g.store.Settings.Get()
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateSettings updates the dynamic settings.
func (g *GRPC) UpdateSettings(ctx context.Context, req *ttnpb.UpdateSettingsRequest) (*pbtypes.Empty, error) {
	if err := g.adminCheck(ctx); err != nil {
		return nil, err
	}

	settings, err := g.store.Settings.Get()
	if err != nil {
		return nil, err
	}

	for _, path := range req.UpdateMask.Paths {
		switch true {
		case ttnpb.FieldPathSettingsBlacklistedIDs.MatchString(path):
			if req.Settings.BlacklistedIDs == nil {
				req.Settings.BlacklistedIDs = []string{}
			}
			settings.BlacklistedIDs = req.Settings.BlacklistedIDs
		case ttnpb.FieldPathSettingsUserRegistrationSkipValidation.MatchString(path):
			settings.SkipValidation = req.Settings.SkipValidation
		case ttnpb.FieldPathSettingsUserRegistrationSelfRegistration.MatchString(path):
			settings.SelfRegistration = req.Settings.SelfRegistration
		case ttnpb.FieldPathSettingsUserRegistrationAdminApproval.MatchString(path):
			settings.AdminApproval = req.Settings.AdminApproval
		case ttnpb.FieldPathSettingsValidationTokenTTL.MatchString(path):
			settings.ValidationTokenTTL = req.Settings.ValidationTokenTTL
		case ttnpb.FieldPathSettingsAllowedEmails.MatchString(path):
			if req.Settings.AllowedEmails == nil {
				req.Settings.AllowedEmails = []string{}
			}
			settings.AllowedEmails = req.Settings.AllowedEmails
		default:
			return nil, ttnpb.ErrInvalidPathUpdateMask.New(errors.Attributes{
				"path": path,
			})
		}
	}

	return nil, g.store.Settings.Set(settings)
}
