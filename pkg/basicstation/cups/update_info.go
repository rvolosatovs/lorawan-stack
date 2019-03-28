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

package cups

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"hash/crc32"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/labstack/echo"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"google.golang.org/grpc"
)

const (
	cupsAttribute               = "_cups"
	cupsURIAttribute            = "_cups_uri"
	cupsLastSeenAttribute       = "_cups_last_seen"
	cupsCredentialsIDAttribute  = "_cups_credentials_id"
	cupsCredentialsAttribute    = "_cups_credentials"
	cupsCredentialsCRCAttribute = "_cups_credentials_crc"
	cupsStationAttribute        = "_cups_station"
	cupsModelAttribute          = "_cups_model"
	cupsPackageAttribute        = "_cups_package"
	lnsCredentialsIDAttribute   = "_lns_credentials_id"
	lnsCredentialsAttribute     = "_lns_credentials"
	lnsCredentialsCRCAttribute  = "_lns_credentials_crc"
)

var (
	errUnauthenticated = errors.DefineUnauthenticated("unauthenticated", "call was not authenticated")
	errCUPSNotEnabled  = errors.DefinePermissionDenied("cups_not_enabled", "CUPS is not enabled for gateway `{gateway_uid}`")
	errInvalidToken    = errors.DefinePermissionDenied("invalid_token", "invalid provisioning token")
)

var getGatewayMask = types.FieldMask{Paths: []string{
	"attributes",
	"version_ids",
	"gateway_server_address",
	"auto_update",
	"update_channel",
	"frequency_plan_id",
}}

func (s *Server) getOrRegisterGateway(ctx context.Context, req UpdateInfoRequest, authHeader string) (gtw *ttnpb.Gateway, err error) {
	logger := log.FromContext(ctx)
	serverAuth := s.getAuth(ctx, req.Router.EUI64, authHeader)

	logger.Info("Finding gateway...")
	ids, err := s.getRegistry(ctx, nil).GetIdentifiersForEUI(ctx, &ttnpb.GetGatewayIdentifiersForEUIRequest{
		EUI: req.Router.EUI64,
	}, serverAuth)

	if errors.IsNotFound(err) && s.registerUnknown {
		gatewayID := fmt.Sprintf("eui-%s", strings.ToLower(req.Router.EUI64.String()))
		logger.WithField("gateway_id", gatewayID).Info("Registering new gateway")
		ids := ttnpb.GatewayIdentifiers{
			GatewayID: gatewayID,
			EUI:       &req.Router.EUI64,
		}
		return s.getRegistry(ctx, &ids).Create(ctx, &ttnpb.CreateGatewayRequest{
			Gateway: ttnpb.Gateway{
				GatewayIdentifiers: ids,
				Attributes: map[string]string{
					cupsAttribute:               "true",
					cupsCredentialsAttribute:    authHeader,
					cupsCredentialsCRCAttribute: strconv.FormatUint(uint64(req.CUPSCredentialsCRC), 10),
					lnsCredentialsCRCAttribute:  strconv.FormatUint(uint64(req.LNSCredentialsCRC), 10),
				},
				GatewayServerAddress: req.LNSURI,
			},
			Collaborator: s.defaultOwner,
		}, serverAuth)
	} else if err != nil {
		return nil, err
	}

	logger = logger.WithField("gateway_id", ids.GetGatewayID())

	if md := rpcmetadata.FromIncomingContext(ctx); strings.ToLower(md.AuthType) == "bearer" && md.AuthValue != "" {
		logger.Debug("Getting gateway with request credentials")
		md.AllowInsecure = s.component.AllowInsecureForCredentials()
		if gtw, err = s.getRegistry(ctx, ids).Get(ctx, &ttnpb.GetGatewayRequest{
			GatewayIdentifiers: *ids,
			FieldMask:          getGatewayMask,
		}, grpc.PerRPCCredentials(md)); err == nil {
			return gtw, nil
		} else if !errors.IsUnauthenticated(err) && !errors.IsPermissionDenied(err) {
			return nil, err
		}
	}

	logger.Debug("Getting gateway with server credentials")
	gtw, err = s.getRegistry(ctx, ids).Get(ctx, &ttnpb.GetGatewayRequest{
		GatewayIdentifiers: *ids,
		FieldMask:          getGatewayMask,
	}, serverAuth)
	if err != nil {
		return nil, err
	}
	if gtw.Attributes[cupsCredentialsAttribute] != "" && authHeader != gtw.Attributes[cupsCredentialsAttribute] {
		return nil, errInvalidToken
	}

	return gtw, nil
}

// UpdateInfo implements the CUPS update-info handler.
func (s *Server) UpdateInfo(c echo.Context) error {
	if c.Request().Header.Get(echo.HeaderContentType) == "" {
		c.Request().Header.Set(echo.HeaderContentType, "application/json")
	}

	var req UpdateInfoRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := getContext(c)
	logger := log.FromContext(ctx).WithFields(log.Fields(
		"gateway_eui", req.Router.EUI64.String(),
	))
	ctx = log.NewContext(ctx, logger)

	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	if authHeader == "" {
		return errUnauthenticated
	}
	gtw, err := s.getOrRegisterGateway(ctx, req, authHeader)
	if err != nil {
		return err
	}

	if gtw.Attributes == nil {
		gtw.Attributes = make(map[string]string)
	}

	if s.requireExplicitEnable {
		if cups, _ := strconv.ParseBool(gtw.Attributes[cupsAttribute]); !cups {
			return errCUPSNotEnabled.WithAttributes("gateway_uid", unique.ID(ctx, gtw.GatewayIdentifiers))
		}
	}

	res := UpdateInfoResponse{}
	authorization := s.getAuth(ctx, req.Router.EUI64, authHeader)

	if gtw.Attributes[cupsURIAttribute] == "" {
		gtw.Attributes[cupsURIAttribute] = req.CUPSURI
	} else if s.allowCUPSURIUpdate && gtw.Attributes[cupsURIAttribute] != req.CUPSURI {
		res.CUPSURI = gtw.Attributes[cupsURIAttribute]
	}

	if gtw.GatewayServerAddress == "" {
		gtw.GatewayServerAddress = req.LNSURI
	} else if gtw.GatewayServerAddress != req.LNSURI {
		scheme, host, port, err := parseAddress(gtw.GatewayServerAddress)
		if err != nil {
			return err
		}
		if scheme == "" {
			scheme = "wss"
		}
		address := host
		if port != "" {
			address = net.JoinHostPort(host, port)
		}
		res.LNSURI = fmt.Sprintf("%s://%s", scheme, address)
	}

	if gtw.Attributes[cupsCredentialsCRCAttribute] != strconv.FormatUint(uint64(req.CUPSCredentialsCRC), 10) {
		if gtw.Attributes[cupsCredentialsAttribute] == "" {
			registry := s.getAccess(ctx, &gtw.GatewayIdentifiers)
			if gtw.Attributes[cupsCredentialsIDAttribute] != "" {
				// TODO: Try deleting old CUPS credentials.
			}
			apiKey, err := registry.CreateAPIKey(ctx, &ttnpb.CreateGatewayAPIKeyRequest{
				GatewayIdentifiers: gtw.GatewayIdentifiers,
				Name:               fmt.Sprintf("CUPS Key, generated %s", time.Now().UTC().Format(time.RFC3339)),
				Rights: []ttnpb.Right{
					ttnpb.RIGHT_GATEWAY_INFO,              // We need to read private attributes.
					ttnpb.RIGHT_GATEWAY_SETTINGS_BASIC,    // We need to write attributes.
					ttnpb.RIGHT_GATEWAY_SETTINGS_API_KEYS, // We need to create API keys.
					ttnpb.RIGHT_GATEWAY_LINK,              // We need to create the LNS API key with this right.
				},
			}, authorization)
			if err != nil {
				return err
			}
			gtw.Attributes[cupsCredentialsIDAttribute] = apiKey.ID
			gtw.Attributes[cupsCredentialsAttribute] = apiKey.Key
		}
		trust, err := s.getTrust(gtw.Attributes[cupsURIAttribute])
		if err != nil {
			return err
		}
		if trust != nil {
			creds, err := TokenCredentials(trust, gtw.Attributes[cupsCredentialsAttribute])
			if err != nil {
				return err
			}
			res.CUPSCredentials = creds
			gtw.Attributes[cupsCredentialsCRCAttribute] = strconv.FormatUint(uint64(crc32.ChecksumIEEE(res.CUPSCredentials)), 10)
		} else {
			delete(gtw.Attributes, cupsCredentialsCRCAttribute)
		}
	}
	if gtw.Attributes[lnsCredentialsCRCAttribute] != strconv.FormatUint(uint64(req.LNSCredentialsCRC), 10) {
		if gtw.Attributes[lnsCredentialsAttribute] == "" {
			registry := s.getAccess(ctx, &gtw.GatewayIdentifiers)
			if gtw.Attributes[lnsCredentialsIDAttribute] != "" {
				// TODO: Try deleting old LNS credentials.
			}
			apiKey, err := registry.CreateAPIKey(ctx, &ttnpb.CreateGatewayAPIKeyRequest{
				GatewayIdentifiers: gtw.GatewayIdentifiers,
				Name:               fmt.Sprintf("LNS Key, generated %s", time.Now().UTC().Format(time.RFC3339)),
				Rights: []ttnpb.Right{
					ttnpb.RIGHT_GATEWAY_INFO,
					ttnpb.RIGHT_GATEWAY_LINK,
				},
			}, authorization)
			if err != nil {
				return err
			}
			gtw.Attributes[lnsCredentialsIDAttribute] = apiKey.ID
			gtw.Attributes[lnsCredentialsAttribute] = apiKey.Key
		}
		trust, err := s.getTrust(gtw.GatewayServerAddress)
		if err != nil {
			return err
		}
		if trust != nil {
			creds, err := TokenCredentials(trust, gtw.Attributes[lnsCredentialsAttribute])
			if err != nil {
				return err
			}
			res.LNSCredentials = creds
			gtw.Attributes[lnsCredentialsCRCAttribute] = strconv.FormatUint(uint64(crc32.ChecksumIEEE(res.LNSCredentials)), 10)
		} else {
			delete(gtw.Attributes, lnsCredentialsCRCAttribute)
		}
	}
	if gtw.AutoUpdate {
		// TODO:
		// Check Station, Model, Package
		// Compare to version_ids, update_channel
		// Get update data
		var updateData []byte
		if updateData != nil {
			var (
				keyCRC uint32
				signer crypto.Signer
			)
			for _, keyCRC = range req.KeyCRCs {
				if sig, ok := s.signers[keyCRC]; ok {
					signer = sig
					break
				}
			}
			if signer != nil {
				hash := sha512.Sum512(updateData)
				sig, err := signer.Sign(rand.Reader, hash[:], nil)
				if err != nil {
					return err
				}
				res.SignatureKeyCRC = keyCRC
				res.Signature = sig
				res.UpdateData = updateData
			}
		}
	}

	gtw.Attributes[cupsLastSeenAttribute] = time.Now().UTC().Format(time.RFC3339)
	gtw.Attributes[cupsStationAttribute] = req.Station
	gtw.Attributes[cupsModelAttribute] = req.Model
	gtw.Attributes[cupsPackageAttribute] = req.Package

	gtw, err = s.getRegistry(ctx, &gtw.GatewayIdentifiers).Update(ctx, &ttnpb.UpdateGatewayRequest{
		Gateway: *gtw,
		FieldMask: types.FieldMask{Paths: []string{
			"attributes",
		}},
	}, authorization)
	if err != nil {
		return err
	}

	b, err := res.MarshalBinary()
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, echo.MIMEOctetStream, b)
}
