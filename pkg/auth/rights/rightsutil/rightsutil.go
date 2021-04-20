// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

package rightsutil

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// EventIsVisible returns whether ev is visible given rights in the context.
func EventIsVisible(ctx context.Context, ev events.Event) (bool, error) {
	visibility := ev.Visibility()
	if len(visibility.Rights) == 0 {
		return true, nil
	}
	for _, entityIDs := range ev.Identifiers() {
		switch ids := entityIDs.GetIds().(type) {
		case *ttnpb.EntityIdentifiers_ApplicationIDs:
			rights, err := rights.ListApplication(ctx, *ids.ApplicationIDs)
			if err != nil {
				return false, err
			}
			if len(rights.Implied().Intersect(visibility).GetRights()) > 0 {
				return true, nil
			}
		case *ttnpb.EntityIdentifiers_ClientIDs:
			rights, err := rights.ListClient(ctx, *ids.ClientIDs)
			if err != nil {
				return false, err
			}
			if len(rights.Implied().Intersect(visibility).GetRights()) > 0 {
				return true, nil
			}
		case *ttnpb.EntityIdentifiers_DeviceIDs:
			rights, err := rights.ListApplication(ctx, ids.DeviceIDs.ApplicationIdentifiers)
			if err != nil {
				return false, err
			}
			if len(rights.Implied().Intersect(visibility).GetRights()) > 0 {
				return true, nil
			}
		case *ttnpb.EntityIdentifiers_GatewayIDs:
			rights, err := rights.ListGateway(ctx, *ids.GatewayIDs)
			if err != nil {
				return false, err
			}
			if len(rights.Implied().Intersect(visibility).GetRights()) > 0 {
				return true, nil
			}
		case *ttnpb.EntityIdentifiers_OrganizationIDs:
			rights, err := rights.ListOrganization(ctx, *ids.OrganizationIDs)
			if err != nil {
				return false, err
			}
			if len(rights.Implied().Intersect(visibility).GetRights()) > 0 {
				return true, nil
			}
		case *ttnpb.EntityIdentifiers_UserIDs:
			rights, err := rights.ListUser(ctx, *ids.UserIDs)
			if err != nil {
				return false, err
			}
			if len(rights.Implied().Intersect(visibility).GetRights()) > 0 {
				return true, nil
			}
		}
	}
	return false, nil
}
