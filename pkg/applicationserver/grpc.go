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

package applicationserver

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/types"
	clusterauth "go.thethings.network/lorawan-stack/v3/pkg/auth/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmiddleware/warning"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

func removeDeprecatedPaths(ctx context.Context, paths []string) []string {
	validPaths := make([]string, 0, len(paths))
nextPath:
	for _, path := range paths {
		for _, deprecated := range []string{
			"api_key",
			"network_server_address",
		} {
			if path == deprecated {
				warning.Add(ctx, fmt.Sprintf("field %v is deprecated", deprecated))
				continue nextPath
			}
			validPaths = append(validPaths, path)
		}
	}
	return validPaths
}

// getLink calls the underlying link registry in order to retrieve the link.
// If the link is not found, an empty link is returned instead.
func (as *ApplicationServer) getLink(ctx context.Context, ids ttnpb.ApplicationIdentifiers, paths []string) (*ttnpb.ApplicationLink, error) {
	link, err := as.linkRegistry.Get(ctx, ids, paths)
	if err != nil && errors.IsNotFound(err) {
		return &ttnpb.ApplicationLink{}, nil
	} else if err != nil {
		return nil, err
	}
	return link, nil
}

// GetLink implements ttnpb.AsServer.
func (as *ApplicationServer) GetLink(ctx context.Context, req *ttnpb.GetApplicationLinkRequest) (*ttnpb.ApplicationLink, error) {
	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_SETTINGS_BASIC); err != nil {
		return nil, err
	}
	req.FieldMask.Paths = removeDeprecatedPaths(ctx, req.FieldMask.Paths)
	return as.linkRegistry.Get(ctx, req.ApplicationIdentifiers, req.FieldMask.Paths)
}

// SetLink implements ttnpb.AsServer.
func (as *ApplicationServer) SetLink(ctx context.Context, req *ttnpb.SetApplicationLinkRequest) (*ttnpb.ApplicationLink, error) {
	if err := rights.RequireApplication(ctx, req.ApplicationIdentifiers, ttnpb.RIGHT_APPLICATION_SETTINGS_BASIC); err != nil {
		return nil, err
	}
	req.FieldMask.Paths = removeDeprecatedPaths(ctx, req.FieldMask.Paths)
	return as.linkRegistry.Set(ctx, req.ApplicationIdentifiers, ttnpb.ApplicationLinkFieldPathsTopLevel,
		func(*ttnpb.ApplicationLink) (*ttnpb.ApplicationLink, []string, error) {
			return &req.ApplicationLink, req.FieldMask.Paths, nil
		},
	)
}

// DeleteLink implements ttnpb.AsServer.
func (as *ApplicationServer) DeleteLink(ctx context.Context, ids *ttnpb.ApplicationIdentifiers) (*types.Empty, error) {
	if err := rights.RequireApplication(ctx, *ids, ttnpb.RIGHT_APPLICATION_SETTINGS_BASIC); err != nil {
		return nil, err
	}
	_, err := as.linkRegistry.Set(ctx, *ids, nil, func(link *ttnpb.ApplicationLink) (*ttnpb.ApplicationLink, []string, error) { return nil, nil, nil })
	if err != nil {
		return nil, err
	}
	return ttnpb.Empty, nil
}

var errLinkingNotImplemented = errors.DefineUnimplemented("linking_not_implemented", "linking is not implemented")

// GetLinkStats implements ttnpb.AsServer.
func (as *ApplicationServer) GetLinkStats(ctx context.Context, ids *ttnpb.ApplicationIdentifiers) (*ttnpb.ApplicationLinkStats, error) {
	return nil, errLinkingNotImplemented.New()
}

// HandleUplink implements ttnpb.NsAsServer.
func (as *ApplicationServer) HandleUplink(ctx context.Context, req *ttnpb.NsAsHandleUplinkRequest) (*types.Empty, error) {
	if err := clusterauth.Authorized(ctx); err != nil {
		return nil, err
	}
	link, err := as.getLink(ctx, req.ApplicationUps[0].ApplicationIdentifiers, []string{
		"default_formatters",
		"skip_payload_crypto",
	})
	if err != nil {
		return nil, err
	}
	// TODO: Merge downlink queue invalidations (https://github.com/TheThingsNetwork/lorawan-stack/issues/1523)
	for _, up := range req.ApplicationUps {
		if err := as.processUp(ctx, up, link); err != nil {
			return nil, err
		}
	}
	return ttnpb.Empty, nil
}
