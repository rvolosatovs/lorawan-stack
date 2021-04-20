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

package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

func firstArgs(i int, args ...string) []string {
	if i > len(args) {
		i = len(args)
	}
	return args[:i]
}

func collaboratorFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.String("user-id", "", "")
	flagSet.String("organization-id", "", "")
	return flagSet
}

var (
	errNoCollaborator       = errors.DefineInvalidArgument("no_collaborator", "no collaborator set")
	errNoCollaboratorRights = errors.DefineInvalidArgument("no_collaborator_rights", "no collaborator rights set")
)

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
	flagSet := &pflag.FlagSet{}
	flagSet.StringSlice("attributes", nil, "key=value")
	return flagSet
}

func forceFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.Bool("force", false, "")
	return flagSet
}

func mergeKV(attributes map[string]string, kv []string) map[string]string {
	out := make(map[string]string, len(attributes)+len(kv))
	for k, v := range attributes {
		out[k] = v
	}
	for _, kv := range kv {
		kv := strings.SplitN(kv, "=", 2)
		if len(kv) != 2 {
			continue
		}
		if kv[1] == "" {
			delete(out, kv[0])
		} else {
			out[kv[0]] = kv[1]
		}
	}
	return out
}

func mergeAttributes(attributes map[string]string, flagSet *pflag.FlagSet) map[string]string {
	kv, _ := flagSet.GetStringSlice("attributes")
	return mergeKV(attributes, kv)
}

func rightsFlags(filter func(string) bool) *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
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

var (
	errNoAPIKeyID     = errors.DefineInvalidArgument("no_api_key_id", "no API key ID set")
	errNoAPIKeyRights = errors.DefineInvalidArgument("no_api_key_rights", "no API key rights set")
)

func getAPIKeyID(flagSet *pflag.FlagSet, args []string, i int) string {
	var apiKeyID string
	if len(args) > 0+i {
		if len(args) > 1+i {
			logger.Warn("Multiple API key IDs found in arguments, considering only the first")
		}
		apiKeyID = args[0+i]
	} else {
		apiKeyID, _ = flagSet.GetString("api-key-id")
	}
	return apiKeyID
}

var searchFlags = func() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	// NOTE: These flags need to be named with underscores, not dashes!
	flagSet.String("id_contains", "", "")
	flagSet.String("name_contains", "", "")
	flagSet.String("description_contains", "", "")
	flagSet.StringToString("attributes_contain", nil, "(key=value)")
	flagSet.AddFlagSet(paginationFlags())
	flagSet.AddFlagSet(orderFlags())
	return flagSet
}()

var errNoIDs = errors.DefineInvalidArgument("no_ids", "no IDs set")

func entityIdentifiersSliceFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.StringSlice("application-id", nil, "")
	flagSet.StringSlice("client-id", nil, "")
	flagSet.StringSlice("device-id", nil, "")
	flagSet.StringSlice("gateway-id", nil, "")
	flagSet.StringSlice("organization-id", nil, "")
	flagSet.StringSlice("user-id", nil, "")
	return flagSet
}

func getEntityIdentifiersSlice(flagSet *pflag.FlagSet) []*ttnpb.EntityIdentifiers {
	applicationIDs, _ := flagSet.GetStringSlice("application-id")
	clientIDs, _ := flagSet.GetStringSlice("client-id")
	deviceIDs, _ := flagSet.GetStringSlice("device-id")
	gatewayIDs, _ := flagSet.GetStringSlice("gateway-id")
	organizationIDs, _ := flagSet.GetStringSlice("organization-id")
	userIDs, _ := flagSet.GetStringSlice("user-id")

	var ids []*ttnpb.EntityIdentifiers

	if len(deviceIDs) > 0 {
		if len(clientIDs)+len(gatewayIDs)+len(organizationIDs)+len(userIDs) > 0 {
			logger.Warn("considering only devices")
		}
		for _, deviceID := range deviceIDs {
			for _, applicationID := range applicationIDs {
				ids = append(ids, (&ttnpb.EndDeviceIdentifiers{
					ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: applicationID},
					DeviceID:               deviceID,
				}).GetEntityIdentifiers())
			}
		}
		return ids
	}

	for _, applicationID := range applicationIDs {
		ids = append(ids, (&ttnpb.ApplicationIdentifiers{ApplicationID: applicationID}).GetEntityIdentifiers())
	}
	for _, clientID := range clientIDs {
		ids = append(ids, (&ttnpb.ClientIdentifiers{ClientID: clientID}).GetEntityIdentifiers())
	}
	for _, gatewayID := range gatewayIDs {
		ids = append(ids, (&ttnpb.GatewayIdentifiers{GatewayID: gatewayID}).GetEntityIdentifiers())
	}
	for _, organizationID := range organizationIDs {
		ids = append(ids, (&ttnpb.OrganizationIdentifiers{OrganizationID: organizationID}).GetEntityIdentifiers())
	}
	for _, userID := range userIDs {
		ids = append(ids, (&ttnpb.UserIdentifiers{UserID: userID}).GetEntityIdentifiers())
	}

	return ids
}

// dataFlags returns a flag set for loading binary data.
// Use getDataBytes() or getDataReader() to obtain the binary data.
// The given name and usage are optional specifiers to differentiate different purposes (i.e. source and destination).
func dataFlags(name, usage string) *pflag.FlagSet {
	flagName := "local-file"
	if name != "" {
		flagName = name + "-" + flagName
	}
	flagUsage := "(local file name)"
	if usage != "" {
		flagUsage = usage + " " + flagUsage
	}
	flagSet := &pflag.FlagSet{}
	flagSet.String(flagName, "", flagUsage)
	return flagSet
}

var errNoData = errors.DefineInvalidArgument("no_data", "no data for `{name}`")

func getDataBytes(name string, flagSet *pflag.FlagSet) ([]byte, error) {
	r, err := getDataReader(name, flagSet)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func getDataReader(name string, flagSet *pflag.FlagSet) (io.Reader, error) {
	flagName := "local-file"
	if name != "" {
		flagName = name + "-" + flagName
	}
	if filename, _ := flagSet.GetString(flagName); filename != "" {
		return os.Open(filename)
	}
	if name == "" {
		name = "default"
	}
	return nil, errNoData.WithAttributes("name", name)
}

const timeFormat = "2006-01-02 15:04:05"

func timestampFlags(name, description string) *pflag.FlagSet {
	flags := &pflag.FlagSet{}

	description = fmt.Sprintf("%s (format: '%s')", description, timeFormat)

	flags.String(name, "", description)
	flags.String(fmt.Sprintf("%s-utc", name), "", fmt.Sprintf("%s (UTC)", description))

	return flags
}

func parseTime(s string, location *time.Location) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	t, err := time.ParseInLocation(timeFormat, s, location)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func getTimestampFlags(flags *pflag.FlagSet, name string) (*time.Time, error) {
	utcName := fmt.Sprintf("%s-utc", name)
	if flags.Changed(utcName) {
		s, _ := flags.GetString(utcName)
		return parseTime(s, time.UTC)
	}
	if flags.Changed(name) {
		s, _ := flags.GetString(name)
		return parseTime(s, time.Local)
	}
	return nil, nil
}

var deletedFlags = func() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.Bool("deleted", false, "return recently deleted")
	return flagSet
}()

func getDeleted(flagSet *pflag.FlagSet) bool {
	deleted, _ := flagSet.GetBool("deleted")
	return deleted
}
