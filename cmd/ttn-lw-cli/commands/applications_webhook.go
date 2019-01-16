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
	"os"

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
	selectApplicationWebhookFlags = util.FieldMaskFlags(&ttnpb.ApplicationWebhook{})
	setApplicationWebhookFlags    = util.FieldFlags(&ttnpb.ApplicationWebhook{})
)

func applicationWebhookIDFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.String("application-id", "", "")
	flagSet.String("webhook-id", "", "")
	return flagSet
}

var errNoWebhookID = errors.DefineInvalidArgument("no_webhook_id", "no webhook ID set")

func getApplicationWebhookID(flagSet *pflag.FlagSet, args []string) (*ttnpb.ApplicationWebhookIdentifiers, error) {
	applicationID, _ := flagSet.GetString("application-id")
	webhookID, _ := flagSet.GetString("webhook-id")
	switch len(args) {
	case 0:
	case 1:
		logger.Warn("Only single ID found in arguments, not considering arguments")
	case 2:
		applicationID = args[0]
		webhookID = args[1]
	default:
		logger.Warn("multiple IDs found in arguments, considering the first")
		applicationID = args[0]
		webhookID = args[1]
	}
	if applicationID == "" {
		return nil, errNoApplicationID
	}
	if webhookID == "" {
		return nil, errNoWebhookID
	}
	return &ttnpb.ApplicationWebhookIdentifiers{
		ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: applicationID},
		WebhookID:              webhookID,
	}, nil
}

var (
	applicationsWebhookCommand = &cobra.Command{
		Use:   "webhook",
		Short: "Application webhook commands",
	}
	applicationsWebhookGetFormatsCommand = &cobra.Command{
		Use:     "get-formats",
		Aliases: []string{"formats"},
		Short:   "Get the available formats for application webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			as, err := api.Dial(ctx, config.ApplicationServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationWebhookRegistryClient(as).GetFormats(ctx, ttnpb.Empty)
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	applicationsWebhookGetCommand = &cobra.Command{
		Use:     "get",
		Aliases: []string{"info"},
		Short:   "Get the properties of an application webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			webhookID, err := getApplicationWebhookID(cmd.Flags(), args)
			if err != nil {
				return err
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationWebhookFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectApplicationWebhookFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, flag.Name)
				})
			}

			as, err := api.Dial(ctx, config.ApplicationServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationWebhookRegistryClient(as).Get(ctx, &ttnpb.GetApplicationWebhookRequest{
				ApplicationWebhookIdentifiers: *webhookID,
				FieldMask:                     types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	applicationsWebhookListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List application webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationWebhookFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectApplicationWebhookFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, flag.Name)
				})
			}

			as, err := api.Dial(ctx, config.ApplicationServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationWebhookRegistryClient(as).List(ctx, &ttnpb.ListApplicationWebhooksRequest{
				ApplicationIdentifiers: *appID,
				FieldMask:              types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	applicationsWebhookSetCommand = &cobra.Command{
		Use:     "set",
		Aliases: []string{"update"},
		Short:   "Set the properties of an application webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			webhookID, err := getApplicationWebhookID(cmd.Flags(), args)
			if err != nil {
				return err
			}
			paths := util.UpdateFieldMask(cmd.Flags(), setApplicationWebhookFlags)

			var webhook ttnpb.ApplicationWebhook
			if err = util.SetFields(&webhook, setApplicationWebhookFlags); err != nil {
				return err
			}
			webhook.ApplicationWebhookIdentifiers = *webhookID

			as, err := api.Dial(ctx, config.ApplicationServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationWebhookRegistryClient(as).Set(ctx, &ttnpb.SetApplicationWebhookRequest{
				ApplicationWebhook: webhook,
				FieldMask:          types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	applicationsWebhookDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Delete an application webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			webhookID, err := getApplicationWebhookID(cmd.Flags(), args)
			if err != nil {
				return err
			}

			as, err := api.Dial(ctx, config.ApplicationServerAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewApplicationWebhookRegistryClient(as).Delete(ctx, webhookID)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	applicationsWebhookCommand.AddCommand(applicationsWebhookGetFormatsCommand)
	applicationsWebhookGetCommand.Flags().AddFlagSet(applicationWebhookIDFlags())
	applicationsWebhookGetCommand.Flags().AddFlagSet(selectApplicationWebhookFlags)
	applicationsWebhookCommand.AddCommand(applicationsWebhookGetCommand)
	applicationsWebhookListCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsWebhookListCommand.Flags().AddFlagSet(selectApplicationWebhookFlags)
	applicationsWebhookCommand.AddCommand(applicationsWebhookListCommand)
	applicationsWebhookSetCommand.Flags().AddFlagSet(applicationWebhookIDFlags())
	applicationsWebhookSetCommand.Flags().AddFlagSet(setApplicationWebhookFlags)
	applicationsWebhookCommand.AddCommand(applicationsWebhookSetCommand)
	applicationsWebhookDeleteCommand.Flags().AddFlagSet(applicationWebhookIDFlags())
	applicationsWebhookCommand.AddCommand(applicationsWebhookDeleteCommand)
	applicationsCommand.AddCommand(applicationsWebhookCommand)
}
