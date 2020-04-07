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
	"encoding/hex"
	"os"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/api"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/io"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/util"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	selectApplicationPubSubFlags       = util.FieldMaskFlags(&ttnpb.ApplicationPubSub{})
	setApplicationPubSubFlags          = util.FieldFlags(&ttnpb.ApplicationPubSub{})
	natsProviderApplicationPubSubFlags = util.FieldFlags(&ttnpb.ApplicationPubSub_NATSProvider{}, "nats")
	mqttProviderApplicationPubSubFlags = util.FieldFlags(&ttnpb.ApplicationPubSub_MQTTProvider{}, "mqtt")

	selectAllApplicationPubSubFlags = util.SelectAllFlagSet("application pubsub")
)

func applicationPubSubIDFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.String("application-id", "", "")
	flagSet.String("pubsub-id", "", "")
	return flagSet
}

func applicationPubSubProviderFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.Bool("nats", false, "use the NATS provider")
	flagSet.AddFlagSet(natsProviderApplicationPubSubFlags)
	flagSet.Bool("mqtt", false, "use the MQTT provider")
	flagSet.AddFlagSet(mqttProviderApplicationPubSubFlags)
	flagSet.AddFlagSet(dataFlags("mqtt.tls-ca", ""))
	flagSet.AddFlagSet(dataFlags("mqtt.tls-client-cert", ""))
	flagSet.AddFlagSet(dataFlags("mqtt.tls-client-key", ""))
	addDeprecatedProviderFlags(flagSet)
	return flagSet
}

func addDeprecatedProviderFlags(flagSet *pflag.FlagSet) {
	util.DeprecateFlag(flagSet, "nats_server_url", "nats.server_url")
}

func forwardDeprecatedProviderFlags(flagSet *pflag.FlagSet) {
	util.ForwardFlag(flagSet, "nats_server_url", "nats.server_url")
}

var errNoPubSubID = errors.DefineInvalidArgument("no_pub_sub_id", "no pubsub ID set")

func getApplicationPubSubID(flagSet *pflag.FlagSet, args []string) (*ttnpb.ApplicationPubSubIdentifiers, error) {
	forwardDeprecatedProviderFlags(flagSet)
	applicationID, _ := flagSet.GetString("application-id")
	pubsubID, _ := flagSet.GetString("pubsub-id")
	switch len(args) {
	case 0:
	case 1:
		logger.Warn("Only single ID found in arguments, not considering arguments")
	case 2:
		applicationID = args[0]
		pubsubID = args[1]
	default:
		logger.Warn("Multiple IDs found in arguments, considering the first")
		applicationID = args[0]
		pubsubID = args[1]
	}
	if applicationID == "" {
		return nil, errNoApplicationID
	}
	if pubsubID == "" {
		return nil, errNoPubSubID
	}
	return &ttnpb.ApplicationPubSubIdentifiers{
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: applicationID},
		PubSubID:               pubsubID,
	}, nil
}

var (
	applicationsPubSubsCommand = &cobra.Command{
		Use:     "pubsubs",
		Aliases: []string{"pubsub", "ps"},
		Short:   "Application pubsub commands",
	}
	applicationsPubSubsGetFormatsCommand = &cobra.Command{
		Use:     "get-formats",
		Aliases: []string{"formats"},
		Short:   "Get the available formats for application pubsubs",
		RunE: func(cmd *cobra.Command, args []string) error {
			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationPubSubRegistryClient(as).GetFormats(ctx, ttnpb.Empty)
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	applicationsPubSubsGetCommand = &cobra.Command{
		Use:     "get [application-id] [pubsub-id]",
		Aliases: []string{"info"},
		Short:   "Get the properties of an application pubsub",
		RunE: func(cmd *cobra.Command, args []string) error {
			forwardDeprecatedProviderFlags(cmd.Flags())
			pubsubID, err := getApplicationPubSubID(cmd.Flags(), args)
			if err != nil {
				return err
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationPubSubFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectApplicationPubSubFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, strings.Replace(flag.Name, "-", "_", -1))
				})
			}
			paths = ttnpb.AllowedFields(paths, ttnpb.AllowedFieldMaskPathsForRPC["/ttn.lorawan.v3.ApplicationPubSubRegistry/Get"])

			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationPubSubRegistryClient(as).Get(ctx, &ttnpb.GetApplicationPubSubRequest{
				ApplicationPubSubIdentifiers: *pubsubID,
				FieldMask:                    types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	applicationsPubSubsListCommand = &cobra.Command{
		Use:     "list [application-id]",
		Aliases: []string{"ls"},
		Short:   "List application pubsubs",
		RunE: func(cmd *cobra.Command, args []string) error {
			forwardDeprecatedProviderFlags(cmd.Flags())
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationPubSubFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectApplicationPubSubFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, strings.Replace(flag.Name, "-", "_", -1))
				})
			}
			paths = ttnpb.AllowedFields(paths, ttnpb.AllowedFieldMaskPathsForRPC["/ttn.lorawan.v3.ApplicationPubSubRegistry/List"])

			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationPubSubRegistryClient(as).List(ctx, &ttnpb.ListApplicationPubSubsRequest{
				ApplicationIdentifiers: *appID,
				FieldMask:              types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	applicationsPubSubsSetCommand = &cobra.Command{
		Use:     "set [application-id] [pubsub-id]",
		Aliases: []string{"update"},
		Short:   "Set the properties of an application pubsub",
		RunE: func(cmd *cobra.Command, args []string) error {
			forwardDeprecatedProviderFlags(cmd.Flags())
			pubsubID, err := getApplicationPubSubID(cmd.Flags(), args)
			if err != nil {
				return err
			}
			paths := util.UpdateFieldMask(cmd.Flags(), setApplicationPubSubFlags)

			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			pubsub, err := ttnpb.NewApplicationPubSubRegistryClient(as).Get(ctx, &ttnpb.GetApplicationPubSubRequest{
				ApplicationPubSubIdentifiers: *pubsubID,
				FieldMask:                    types.FieldMask{Paths: paths},
			})
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
			if pubsub == nil {
				pubsub = &ttnpb.ApplicationPubSub{ApplicationPubSubIdentifiers: *pubsubID}
			}

			if err = util.SetFields(pubsub, setApplicationPubSubFlags); err != nil {
				return err
			}

			if nats, _ := cmd.Flags().GetBool("nats"); nats {
				if pubsub.GetNATS() == nil {
					paths = append(paths, "provider")
					pubsub.Provider = &ttnpb.ApplicationPubSub_NATS{
						NATS: &ttnpb.ApplicationPubSub_NATSProvider{},
					}
				} else {
					providerPaths := util.UpdateFieldMask(cmd.Flags(), natsProviderApplicationPubSubFlags)
					providerPaths = ttnpb.FieldsWithPrefix("provider", providerPaths...)
					paths = append(paths, providerPaths...)
				}
				if err = util.SetFields(pubsub.GetNATS(), natsProviderApplicationPubSubFlags, "nats"); err != nil {
					return err
				}
			}

			if mqtt, _ := cmd.Flags().GetBool("mqtt"); mqtt {
				if pubsub.GetMQTT() == nil {
					paths = append(paths, "provider")
					pubsub.Provider = &ttnpb.ApplicationPubSub_MQTT{
						MQTT: &ttnpb.ApplicationPubSub_MQTTProvider{},
					}
				} else {
					providerPaths := util.UpdateFieldMask(cmd.Flags(), mqttProviderApplicationPubSubFlags)
					providerPaths = ttnpb.FieldsWithPrefix("provider", providerPaths...)
					paths = append(paths, providerPaths...)
				}
				if useTLS, _ := cmd.Flags().GetBool("mqtt.use-tls"); useTLS {
					for _, name := range []string{
						"mqtt.tls-ca",
						"mqtt.tls-client-cert",
						"mqtt.tls-client-key",
					} {
						data, err := getDataBytes(name, cmd.Flags())
						if err != nil {
							return err
						}
						err = cmd.Flags().Set(name, hex.EncodeToString(data))
						if err != nil {
							return err
						}
					}
				}
				if err = util.SetFields(pubsub.GetMQTT(), mqttProviderApplicationPubSubFlags, "mqtt"); err != nil {
					return err
				}
			}

			res, err := ttnpb.NewApplicationPubSubRegistryClient(as).Set(ctx, &ttnpb.SetApplicationPubSubRequest{
				ApplicationPubSub: *pubsub,
				FieldMask:         types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	applicationsPubSubsDeleteCommand = &cobra.Command{
		Use:   "delete [application-id] [pubsub-id]",
		Short: "Delete an application pubsub",
		RunE: func(cmd *cobra.Command, args []string) error {
			pubsubID, err := getApplicationPubSubID(cmd.Flags(), args)
			if err != nil {
				return err
			}

			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewApplicationPubSubRegistryClient(as).Delete(ctx, pubsubID)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	applicationsPubSubsCommand.AddCommand(applicationsPubSubsGetFormatsCommand)
	applicationsPubSubsGetCommand.Flags().AddFlagSet(applicationPubSubIDFlags())
	applicationsPubSubsGetCommand.Flags().AddFlagSet(selectApplicationPubSubFlags)
	applicationsPubSubsGetCommand.Flags().AddFlagSet(selectAllApplicationPubSubFlags)
	applicationsPubSubsCommand.AddCommand(applicationsPubSubsGetCommand)
	applicationsPubSubsListCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsPubSubsListCommand.Flags().AddFlagSet(selectApplicationPubSubFlags)
	applicationsPubSubsListCommand.Flags().AddFlagSet(selectAllApplicationPubSubFlags)
	applicationsPubSubsCommand.AddCommand(applicationsPubSubsListCommand)
	applicationsPubSubsSetCommand.Flags().AddFlagSet(applicationPubSubIDFlags())
	applicationsPubSubsSetCommand.Flags().AddFlagSet(setApplicationPubSubFlags)
	applicationsPubSubsSetCommand.Flags().AddFlagSet(applicationPubSubProviderFlags())
	applicationsPubSubsCommand.AddCommand(applicationsPubSubsSetCommand)
	applicationsPubSubsDeleteCommand.Flags().AddFlagSet(applicationPubSubIDFlags())
	applicationsPubSubsCommand.AddCommand(applicationsPubSubsDeleteCommand)
	applicationsCommand.AddCommand(applicationsPubSubsCommand)
}
