// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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

package commands

import (
	"strings"

	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

func collaboratorFlags() *pflag.FlagSet {
	flagSet := new(pflag.FlagSet)
	flagSet.String("user-id", "", "")
	flagSet.String("organization-id", "", "")
	return flagSet
}

var errNoCollaborator = errors.DefineInvalidArgument("no_collaborator", "no collaborator set")

func getCollaborator(flagSet *pflag.FlagSet) *ttnpb.OrganizationOrUserIdentifiers {
	organizationID, _ := flagSet.GetString("organization-id")
	userID, _ := flagSet.GetString("user-id")
	if organizationID == "" && userID == "" {
		return nil
	}
	if organizationID != "" && userID != "" {
		logger.Warn("Don't set organization ID and user ID at the same time, assuming user ID")
	}
	if userID != "" {
		return ttnpb.UserIdentifiers{UserID: userID}.OrganizationOrUserIdentifiers()
	}
	return ttnpb.OrganizationIdentifiers{OrganizationID: organizationID}.OrganizationOrUserIdentifiers()
}

func attributesFlags() *pflag.FlagSet {
	flagSet := new(pflag.FlagSet)
	flagSet.StringSlice("attributes", nil, "key=value")
	return flagSet
}

func mergeAttributes(attributes map[string]string, flagSet *pflag.FlagSet) map[string]string {
	kv, _ := flagSet.GetStringSlice("attributes")
	out := make(map[string]string, len(attributes)+len(kv))
	for k, v := range attributes {
		out[k] = v
	}
	for _, kv := range kv {
		kv := strings.SplitN(kv, "=", 2)
		if kv[1] == "" {
			delete(out, kv[0])
		} else {
			out[kv[0]] = kv[1]
		}
	}
	return out
}

func rightsFlags(filter func(string) bool) *pflag.FlagSet {
	flagSet := new(pflag.FlagSet)
	for right := range ttnpb.Right_value {
		right := strings.Replace(strings.ToLower(right), "_", "-", -1)
		if filter == nil || filter(right) {
			flagSet.Bool(right, false, "")
		}
	}
	return flagSet
}

func getRights(flagSet *pflag.FlagSet) (rights []ttnpb.Right) {
	for right, value := range ttnpb.Right_value {
		right := strings.Replace(strings.ToLower(right), "_", "-", -1)
		if set, _ := flagSet.GetBool(right); set {
			rights = append(rights, ttnpb.Right(value))
		}
	}
	return
}

var errNoAPIKeyID = errors.DefineInvalidArgument("no_api_key_id", "no API key ID set")

func getAPIKeyID(flagSet *pflag.FlagSet, args []string) string {
	var apiKeyID string
	if len(args) > 0 {
		if len(args) > 1 {
			logger.Warn("multiple IDs found in arguments, considering only the first")
		}
		apiKeyID = args[0]
	} else {
		apiKeyID, _ = flagSet.GetString("api-key-id")
	}
	return apiKeyID
}

func searchFlags() *pflag.FlagSet {
	flagSet := new(pflag.FlagSet)
	flagSet.String("id-contains", "", "")
	flagSet.String("name-contains", "", "")
	flagSet.String("description-contains", "", "")
	flagSet.StringToString("attributes-contain", nil, "(key=value)")
	return flagSet
}

func getSearchEntitiesRequest(flagSet *pflag.FlagSet) *ttnpb.SearchEntitiesRequest {
	idContains, _ := flagSet.GetString("id-contains")
	nameContains, _ := flagSet.GetString("name-contains")
	descriptionContains, _ := flagSet.GetString("description-contains")
	attributesContain, _ := flagSet.GetStringToString("attributes-contain")
	return &ttnpb.SearchEntitiesRequest{
		IDContains:          idContains,
		NameContains:        nameContains,
		DescriptionContains: descriptionContains,
		AttributesContain:   attributesContain,
	}
}

var errNoIDs = errors.DefineInvalidArgument("no_ids", "no IDs set")

func combinedIdentifiersFlags() *pflag.FlagSet {
	flagSet := new(pflag.FlagSet)
	flagSet.StringSlice("application-id", nil, "")
	flagSet.StringSlice("client-id", nil, "")
	flagSet.StringSlice("device-id", nil, "")
	flagSet.StringSlice("gateway-id", nil, "")
	flagSet.StringSlice("organization-id", nil, "")
	flagSet.StringSlice("user-id", nil, "")
	return flagSet
}

func getCombinedIdentifiers(flagSet *pflag.FlagSet) *ttnpb.CombinedIdentifiers {
	applicationIDs, _ := flagSet.GetStringSlice("application-id")
	clientIDs, _ := flagSet.GetStringSlice("client-id")
	deviceIDs, _ := flagSet.GetStringSlice("device-id")
	gatewayIDs, _ := flagSet.GetStringSlice("gateway-id")
	organizationIDs, _ := flagSet.GetStringSlice("organization-id")
	userIDs, _ := flagSet.GetStringSlice("user-id")

	ids := new(ttnpb.CombinedIdentifiers)
	if len(deviceIDs) > 0 {
		if len(clientIDs)+len(gatewayIDs)+len(organizationIDs)+len(userIDs) > 0 {
			logger.Warn("considering only devices")
		}
		for _, deviceID := range deviceIDs {
			for _, applicationID := range applicationIDs {
				ids.EntityIdentifiers = append(ids.EntityIdentifiers, ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: applicationID},
					DeviceID:               deviceID,
				}.EntityIdentifiers())
			}
		}
		return ids
	}

	for _, applicationID := range applicationIDs {
		ids.EntityIdentifiers = append(ids.EntityIdentifiers, ttnpb.ApplicationIdentifiers{ApplicationID: applicationID}.EntityIdentifiers())
	}
	for _, clientID := range clientIDs {
		ids.EntityIdentifiers = append(ids.EntityIdentifiers, ttnpb.ClientIdentifiers{ClientID: clientID}.EntityIdentifiers())
	}
	for _, gatewayID := range gatewayIDs {
		ids.EntityIdentifiers = append(ids.EntityIdentifiers, ttnpb.GatewayIdentifiers{GatewayID: gatewayID}.EntityIdentifiers())
	}
	for _, organizationID := range organizationIDs {
		ids.EntityIdentifiers = append(ids.EntityIdentifiers, ttnpb.OrganizationIdentifiers{OrganizationID: organizationID}.EntityIdentifiers())
	}
	for _, userID := range userIDs {
		ids.EntityIdentifiers = append(ids.EntityIdentifiers, ttnpb.UserIdentifiers{UserID: userID}.EntityIdentifiers())
	}
	return ids
}
