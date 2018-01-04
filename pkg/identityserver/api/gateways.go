// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package api

import (
	"context"

	"github.com/TheThingsNetwork/ttn/pkg/auth"
	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/util"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	pbtypes "github.com/gogo/protobuf/types"
)

var _ ttnpb.IsGatewayServer = new(GRPC)

const APIKeyName = "Default API Key"

// CreateGateway creates a gateway in the network, sets the user as collaborator
// with all rights and creates an API key
func (g *GRPC) CreateGateway(ctx context.Context, req *ttnpb.CreateGatewayRequest) (*pbtypes.Empty, error) {
	userID, err := g.userCheck(ctx, ttnpb.RIGHT_USER_GATEWAYS_CREATE)
	if err != nil {
		return nil, err
	}

	settings, err := g.store.Settings.Get()
	if err != nil {
		return nil, err
	}

	// check for blacklisted ids
	if !util.IsIDAllowed(req.Gateway.GatewayID, settings.BlacklistedIDs) {
		return nil, ErrBlacklistedID.New(errors.Attributes{
			"id": req.Gateway.GatewayID,
		})
	}

	err = g.store.Transact(func(s *store.Store) error {
		err = s.Gateways.Create(&ttnpb.Gateway{
			GatewayIdentifier: req.Gateway.GatewayIdentifier,
			Description:       req.Gateway.Description,
			FrequencyPlanID:   req.Gateway.FrequencyPlanID,
			PrivacySettings:   req.Gateway.PrivacySettings,
			AutoUpdate:        req.Gateway.AutoUpdate,
			Platform:          req.Gateway.Platform,
			Antennas:          req.Gateway.Antennas,
			Attributes:        req.Gateway.Attributes,
			ClusterAddress:    req.Gateway.ClusterAddress,
			ContactAccount:    req.Gateway.ContactAccount,
		})
		if err != nil {
			return err
		}

		k, err := auth.GenerateGatewayAPIKey("")
		if err != nil {
			return err
		}

		key := &ttnpb.APIKey{
			Name:   APIKeyName,
			Key:    k,
			Rights: []ttnpb.Right{ttnpb.RIGHT_GATEWAY_INFO},
		}

		err = s.Gateways.SaveAPIKey(req.Gateway.GatewayID, key)
		if err != nil {
			return err
		}

		return s.Gateways.SetCollaborator(&ttnpb.GatewayCollaborator{
			GatewayIdentifier: req.Gateway.GatewayIdentifier,
			UserIdentifier:    ttnpb.UserIdentifier{userID},
			Rights:            ttnpb.AllGatewayRights,
		})
	})

	return nil, err
}

// GetGateway returns a gateway information.
func (g *GRPC) GetGateway(ctx context.Context, req *ttnpb.GatewayIdentifier) (*ttnpb.Gateway, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_INFO)
	if err != nil {
		return nil, err
	}

	found, err := g.store.Gateways.GetByID(req.GatewayID, g.factories.gateway)
	if err != nil {
		return nil, err
	}

	return found.GetGateway(), nil
}

// ListGateways returns all the gateways the current user is collaborator of.
func (g *GRPC) ListGateways(ctx context.Context, _ *pbtypes.Empty) (*ttnpb.ListGatewaysResponse, error) {
	userID, err := g.userCheck(ctx, ttnpb.RIGHT_USER_GATEWAYS_LIST)
	if err != nil {
		return nil, err
	}

	found, err := g.store.Gateways.ListByUser(userID, g.factories.gateway)
	if err != nil {
		return nil, err
	}

	resp := &ttnpb.ListGatewaysResponse{
		Gateways: make([]*ttnpb.Gateway, 0, len(found)),
	}

	for _, gtw := range found {
		resp.Gateways = append(resp.Gateways, gtw.GetGateway())
	}

	return resp, nil
}

// UpdateGateway updates a gateway.
func (g *GRPC) UpdateGateway(ctx context.Context, req *ttnpb.UpdateGatewayRequest) (*pbtypes.Empty, error) {
	err := g.gatewayCheck(ctx, req.Gateway.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_BASIC)
	if err != nil {
		return nil, err
	}

	found, err := g.store.Gateways.GetByID(req.Gateway.GatewayID, g.factories.gateway)
	if err != nil {
		return nil, err
	}
	gtw := found.GetGateway()

	for _, path := range req.UpdateMask.Paths {
		switch true {
		case ttnpb.FieldPathGatewayDescription.MatchString(path):
			gtw.Description = req.Gateway.Description
		case ttnpb.FieldPathGatewayFrequencyPlanID.MatchString(path):
			gtw.FrequencyPlanID = req.Gateway.FrequencyPlanID
		case ttnpb.FieldPathGatewayPrivacySettingsStatusPublic.MatchString(path):
			gtw.PrivacySettings.StatusPublic = req.Gateway.PrivacySettings.StatusPublic
		case ttnpb.FieldPathGatewayPrivacySettingsLocationPublic.MatchString(path):
			gtw.PrivacySettings.LocationPublic = req.Gateway.PrivacySettings.LocationPublic
		case ttnpb.FieldPathGatewayPrivacySettingsContactable.MatchString(path):
			gtw.PrivacySettings.Contactable = req.Gateway.PrivacySettings.Contactable
		case ttnpb.FieldPathGatewayAutoUpdate.MatchString(path):
			gtw.AutoUpdate = req.Gateway.AutoUpdate
		case ttnpb.FieldPathGatewayPlatform.MatchString(path):
			gtw.Platform = req.Gateway.Platform
		case ttnpb.FieldPathGatewayAntennas.MatchString(path):
			if req.Gateway.Antennas == nil {
				req.Gateway.Antennas = []ttnpb.GatewayAntenna{}
			}
			gtw.Antennas = req.Gateway.Antennas
		case ttnpb.FieldPathGatewayAttributes.MatchString(path):
			attr := ttnpb.FieldPathGatewayAttributes.FindStringSubmatch(path)[1]

			if value, ok := req.Gateway.Attributes[attr]; ok && len(value) > 0 {
				gtw.Attributes[attr] = value
			} else {
				delete(gtw.Attributes, attr)
			}
		case ttnpb.FieldPathGatewayClusterAddress.MatchString(path):
			gtw.ClusterAddress = req.Gateway.ClusterAddress
		case ttnpb.FieldPathGatewayContactAccountUserID.MatchString(path):
			gtw.ContactAccount.UserID = req.Gateway.ContactAccount.UserID
		default:
			return nil, ttnpb.ErrInvalidPathUpdateMask.New(errors.Attributes{
				"path": path,
			})
		}
	}

	return nil, g.store.Gateways.Update(gtw)
}

// DeleteGateway deletes a gateway.
func (g *GRPC) DeleteGateway(ctx context.Context, req *ttnpb.GatewayIdentifier) (*pbtypes.Empty, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_DELETE)
	if err != nil {
		return nil, err
	}

	return nil, g.store.Gateways.Delete(req.GatewayID)
}

// GenerateGatewayAPIKey generates a gateway API key and returns it.
func (g *GRPC) GenerateGatewayAPIKey(ctx context.Context, req *ttnpb.GenerateGatewayAPIKeyRequest) (*ttnpb.APIKey, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_KEYS)
	if err != nil {
		return nil, err
	}

	k, err := auth.GenerateGatewayAPIKey("")
	if err != nil {
		return nil, err
	}

	key := &ttnpb.APIKey{
		Key:    k,
		Name:   req.Name,
		Rights: req.Rights,
	}

	err = g.store.Gateways.SaveAPIKey(req.GatewayID, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// ListGatewayAPIKeys list all the API keys from a gateway.
func (g *GRPC) ListGatewayAPIKeys(ctx context.Context, req *ttnpb.GatewayIdentifier) (*ttnpb.ListGatewayAPIKeysResponse, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_KEYS)
	if err != nil {
		return nil, err
	}

	found, err := g.store.Gateways.ListAPIKeys(req.GatewayID)
	if err != nil {
		return nil, err
	}

	return &ttnpb.ListGatewayAPIKeysResponse{
		APIKeys: found,
	}, nil
}

// UpdateGatewayAPIKey updates an API key rights.
func (g *GRPC) UpdateGatewayAPIKey(ctx context.Context, req *ttnpb.UpdateGatewayAPIKeyRequest) (*pbtypes.Empty, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_KEYS)
	if err != nil {
		return nil, err
	}

	return nil, g.store.Gateways.UpdateAPIKeyRights(req.GatewayID, req.Name, req.Rights)
}

// RemoveGatewayAPIKey removes a gateway API key.
func (g *GRPC) RemoveGatewayAPIKey(ctx context.Context, req *ttnpb.RemoveGatewayAPIKeyRequest) (*pbtypes.Empty, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_KEYS)
	if err != nil {
		return nil, err
	}

	return nil, g.store.Gateways.DeleteAPIKey(req.GatewayID, req.Name)
}

// SetGatewayCollaborator sets or unsets a gateway collaborator. It returns error
// if after unset a collaborators there is no at least one collaborator with
// `gateway:settings:collaborators` right.
func (g *GRPC) SetGatewayCollaborator(ctx context.Context, req *ttnpb.GatewayCollaborator) (*pbtypes.Empty, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_COLLABORATORS)
	if err != nil {
		return nil, err
	}

	err = g.store.Transact(func(s *store.Store) error {
		err := s.Gateways.SetCollaborator(req)
		if err != nil {
			return err
		}

		// check that there is at least one collaborator in with SETTINGS_COLLABORATOR right
		collaborators, err := s.Gateways.ListCollaborators(req.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_COLLABORATORS)
		if err != nil {
			return err
		}

		if len(collaborators) == 0 {
			return errors.Errorf("failed to unset collaborator: `%s` must have at least one collaborator with `gateway:settings:collaborators right", req.GatewayID)
		}

		return nil
	})

	return nil, err
}

// ListGatewayCollaborators returns all the collaborators that a gateway has.
func (g *GRPC) ListGatewayCollaborators(ctx context.Context, req *ttnpb.GatewayIdentifier) (*ttnpb.ListGatewayCollaboratorsResponse, error) {
	err := g.gatewayCheck(ctx, req.GatewayID, ttnpb.RIGHT_GATEWAY_SETTINGS_COLLABORATORS)
	if err != nil {
		return nil, err
	}

	found, err := g.store.Gateways.ListCollaborators(req.GatewayID)
	if err != nil {
		return nil, err
	}

	return &ttnpb.ListGatewayCollaboratorsResponse{
		Collaborators: found,
	}, nil
}

// ListGatewayRights returns the rights the caller user has to a gateway.
func (g *GRPC) ListGatewayRights(ctx context.Context, req *ttnpb.GatewayIdentifier) (*ttnpb.ListGatewayRightsResponse, error) {
	userID, err := g.userCheck(ctx)
	if err != nil {
		return nil, err
	}

	rights, err := g.store.Gateways.ListUserRights(req.GatewayID, userID)
	if err != nil {
		return nil, err
	}

	return &ttnpb.ListGatewayRightsResponse{
		Rights: rights,
	}, nil
}
