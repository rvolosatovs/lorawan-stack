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

package ttnpb

import (
	"context"
	"fmt"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
)

// IsZero returns true if all identifiers have zero-values.
func (ids ApplicationIdentifiers) IsZero() bool {
	return ids.ApplicationID == ""
}

// FieldIsZero returns whether path p is zero.
func (v *ApplicationIdentifiers) FieldIsZero(p string) bool {
	if v == nil {
		return true
	}
	switch p {
	case "application_id":
		return v.ApplicationID == ""
	}
	panic(fmt.Sprintf("unknown path '%s'", p))
}

// IsZero returns true if all identifiers have zero-values.
func (ids ClientIdentifiers) IsZero() bool {
	return ids.ClientID == ""
}

// IsZero reports whether ids represent zero identifiers.
func (ids EndDeviceIdentifiers) IsZero() bool {
	return ids.GetDeviceID() == "" &&
		ids.GetApplicationID() == "" &&
		(ids.DevAddr == nil || ids.DevAddr.IsZero()) &&
		(ids.DevEUI == nil || ids.DevEUI.IsZero()) &&
		ids.JoinEUI == nil
}

// FieldIsZero returns whether path p is zero.
func (v *EndDeviceIdentifiers) FieldIsZero(p string) bool {
	if v == nil {
		return true
	}
	switch p {
	case "application_ids":
		return v.ApplicationIdentifiers == ApplicationIdentifiers{}
	case "application_ids.application_id":
		return v.ApplicationIdentifiers.FieldIsZero("application_id")
	case "dev_addr":
		return v.DevAddr == nil
	case "dev_eui":
		return v.DevEUI == nil
	case "device_id":
		return v.DeviceID == ""
	case "join_eui":
		return v.JoinEUI == nil
	}
	panic(fmt.Sprintf("unknown path '%s'", p))
}

// IsZero returns true if all identifiers have zero-values.
func (ids GatewayIdentifiers) IsZero() bool {
	return ids.GatewayID == "" && ids.EUI == nil
}

// IsZero returns true if all identifiers have zero-values.
func (ids OrganizationIdentifiers) IsZero() bool {
	return ids.OrganizationID == ""
}

// IsZero returns true if all identifiers have zero-values.
func (ids UserIdentifiers) IsZero() bool {
	return ids.UserID == "" && ids.Email == ""
}

// GetOrganizationOrUserIdentifiers returns the OrganizationIdentifiers as *OrganizationOrUserIdentifiers.
func (ids *OrganizationIdentifiers) GetOrganizationOrUserIdentifiers() *OrganizationOrUserIdentifiers {
	if ids == nil {
		return nil
	}
	return ids.OrganizationOrUserIdentifiers()
}

// OrganizationOrUserIdentifiers returns the OrganizationIdentifiers as *OrganizationOrUserIdentifiers.
func (ids OrganizationIdentifiers) OrganizationOrUserIdentifiers() *OrganizationOrUserIdentifiers {
	return &OrganizationOrUserIdentifiers{Ids: &OrganizationOrUserIdentifiers_OrganizationIDs{
		OrganizationIDs: &ids,
	}}
}

// GetOrganizationOrUserIdentifiers returns the UserIdentifiers as *OrganizationOrUserIdentifiers.
func (ids *UserIdentifiers) GetOrganizationOrUserIdentifiers() *OrganizationOrUserIdentifiers {
	if ids == nil {
		return nil
	}
	return ids.OrganizationOrUserIdentifiers()
}

// OrganizationOrUserIdentifiers returns the UserIdentifiers as *OrganizationOrUserIdentifiers.
func (ids UserIdentifiers) OrganizationOrUserIdentifiers() *OrganizationOrUserIdentifiers {
	return &OrganizationOrUserIdentifiers{Ids: &OrganizationOrUserIdentifiers_UserIDs{
		UserIDs: &ids,
	}}
}

// CombineIdentifiers merges the identifiers of the multiple entities.
func CombineIdentifiers(ids ...Identifiers) *CombinedIdentifiers {
	combined := &CombinedIdentifiers{}
	for _, id := range ids {
		combined.EntityIdentifiers = append(combined.EntityIdentifiers, id.EntityIdentifiers())
	}
	return combined
}

// CombinedIdentifiers implements Identifiers.
func (ids *CombinedIdentifiers) CombinedIdentifiers() *CombinedIdentifiers {
	return ids
}

// CombinedIdentifiers implements Identifiers.
func (m *ListApplicationsRequest) CombinedIdentifiers() *CombinedIdentifiers {
	if m.Collaborator != nil {
		return m.Collaborator.CombinedIdentifiers()
	}
	return &CombinedIdentifiers{}
}

// CombinedIdentifiers implements Identifiers.
func (m *ListClientsRequest) CombinedIdentifiers() *CombinedIdentifiers {
	if m.Collaborator != nil {
		return m.Collaborator.CombinedIdentifiers()
	}
	return &CombinedIdentifiers{}
}

// CombinedIdentifiers implements Identifiers.
func (m *ListGatewaysRequest) CombinedIdentifiers() *CombinedIdentifiers {
	if m.Collaborator != nil {
		return m.Collaborator.CombinedIdentifiers()
	}
	return &CombinedIdentifiers{}
}

// CombinedIdentifiers implements Identifiers.
func (m *ListOrganizationsRequest) CombinedIdentifiers() *CombinedIdentifiers {
	if m.Collaborator != nil {
		return m.Collaborator.CombinedIdentifiers()
	}
	return &CombinedIdentifiers{}
}

// CombinedIdentifiers implements Identifiers.
func (m *CreateApplicationRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return m.Collaborator.CombinedIdentifiers()
}

// CombinedIdentifiers implements Identifiers.
func (m *CreateClientRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return m.Collaborator.CombinedIdentifiers()
}

// CombinedIdentifiers implements Identifiers.
func (m *CreateGatewayRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return m.Collaborator.CombinedIdentifiers()
}

// CombinedIdentifiers implements Identifiers.
func (m *CreateOrganizationRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return m.Collaborator.CombinedIdentifiers()
}

// CombinedIdentifiers implements Identifiers.
func (m *DownlinkMessage) CombinedIdentifiers() *CombinedIdentifiers {
	return m.GetEndDeviceIDs().CombinedIdentifiers()
}

// CombinedIdentifiers implements Identifiers.
func (m *SetEndDeviceRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return m.EndDevice.CombinedIdentifiers()
}

// CombinedIdentifiers implements Identifiers.
func (m *ListOAuthAccessTokensRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return CombineIdentifiers(m.UserIDs, m.ClientIDs)
}

// CombinedIdentifiers implements Identifiers.
func (m *ListOAuthClientAuthorizationsRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return m.UserIdentifiers.CombinedIdentifiers()
}

// CombinedIdentifiers implements Identifiers.
func (m *OAuthAccessTokenIdentifiers) CombinedIdentifiers() *CombinedIdentifiers {
	return CombineIdentifiers(m.UserIDs, m.ClientIDs)
}

// CombinedIdentifiers implements Identifiers.
func (m *OAuthClientAuthorizationIdentifiers) CombinedIdentifiers() *CombinedIdentifiers {
	return CombineIdentifiers(m.UserIDs, m.ClientIDs)
}

// CombinedIdentifiers implements Identifiers.
func (m *StreamEventsRequest) CombinedIdentifiers() *CombinedIdentifiers {
	return &CombinedIdentifiers{EntityIdentifiers: m.Identifiers}
}

// Copy stores a copy of ids in x and returns it.
func (ids EndDeviceIdentifiers) Copy(x *EndDeviceIdentifiers) *EndDeviceIdentifiers {
	*x = EndDeviceIdentifiers{
		DeviceID: ids.DeviceID,
		ApplicationIdentifiers: ApplicationIdentifiers{
			ApplicationID: ids.ApplicationID,
		},
		XXX_sizecache: ids.XXX_sizecache,
	}
	if ids.DevEUI != nil {
		x.DevEUI = ids.DevEUI.Copy(&types.EUI64{})
	}
	if ids.JoinEUI != nil {
		x.JoinEUI = ids.JoinEUI.Copy(&types.EUI64{})
	}
	if ids.DevAddr != nil {
		x.DevAddr = ids.DevAddr.Copy(&types.DevAddr{})
	}
	return x
}

var errIdentifiers = errors.DefineInvalidArgument("identifiers", "invalid identifiers")

// ValidateContext wraps the generated validator with (optionally context-based) custom checks.
func (ids *EndDeviceIdentifiers) ValidateContext(context.Context) error {
	if err := ids.ValidateFields(); err != nil {
		return errIdentifiers.WithCause(err)
	}
	return nil
}

// ValidateContext wraps the generated validator with (optionally context-based) custom checks.
func (ids *ApplicationIdentifiers) ValidateContext(context.Context) error {
	if err := ids.ValidateFields(); err != nil {
		return errIdentifiers.WithCause(err)
	}
	return nil
}

// ValidateContext wraps the generated validator with (optionally context-based) custom checks.
func (ids *GatewayIdentifiers) ValidateContext(context.Context) error {
	if err := ids.ValidateFields(); err != nil {
		return errIdentifiers.WithCause(err)
	}
	return nil
}

// ValidateContext wraps the generated validator with (optionally context-based) custom checks.
func (ids *UserIdentifiers) ValidateContext(context.Context) error {
	if err := ids.ValidateFields(); err != nil {
		return errIdentifiers.WithCause(err)
	}
	return nil
}

// ExtractRequestFields lets gRPC middleware extract fields from request messages that embed these identifiers.
func (ids *ApplicationIdentifiers) ExtractRequestFields(m map[string]interface{}) {
	if ids == nil {
		return
	}
	m["application_id"] = ids.ApplicationID
}

// ExtractRequestFields lets gRPC middleware extract fields from request messages that embed these identifiers.
func (ids *ClientIdentifiers) ExtractRequestFields(m map[string]interface{}) {
	if ids == nil {
		return
	}
	m["client_id"] = ids.ClientID
}

// ExtractRequestFields lets gRPC middleware extract fields from request messages that embed these identifiers.
func (ids *EndDeviceIdentifiers) ExtractRequestFields(m map[string]interface{}) {
	if ids == nil {
		return
	}
	m["application_id"] = ids.ApplicationID
	m["device_id"] = ids.DeviceID
}

// ExtractRequestFields lets gRPC middleware extract fields from request messages that embed these identifiers.
func (ids *GatewayIdentifiers) ExtractRequestFields(m map[string]interface{}) {
	if ids == nil {
		return
	}
	m["gateway_id"] = ids.GatewayID
}

// ExtractRequestFields lets gRPC middleware extract fields from request messages that embed these identifiers.
func (ids *OrganizationIdentifiers) ExtractRequestFields(m map[string]interface{}) {
	if ids == nil {
		return
	}
	m["organization_id"] = ids.OrganizationID
}

// ExtractRequestFields lets gRPC middleware extract fields from request messages that embed these identifiers.
func (ids *UserIdentifiers) ExtractRequestFields(m map[string]interface{}) {
	if ids == nil {
		return
	}
	m["user_id"] = ids.UserID
}
