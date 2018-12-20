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
	selectApplicationFlags = util.FieldMaskFlags(&ttnpb.Application{})
	setApplicationFlags    = util.FieldFlags(&ttnpb.Application{})
)

func applicationIDFlags() *pflag.FlagSet {
	flagSet := new(pflag.FlagSet)
	flagSet.String("application-id", "", "")
	return flagSet
}

var errNoApplicationID = errors.DefineInvalidArgument("no_application_id", "no application ID set")

func getApplicationID(flagSet *pflag.FlagSet, args []string) *ttnpb.ApplicationIdentifiers {
	var applicationID string
	if len(args) > 0 {
		if len(args) > 1 {
			logger.Warn("multiple IDs found in arguments, considering only the first")
		}
		applicationID = args[0]
	} else {
		applicationID, _ = flagSet.GetString("application-id")
	}
	if applicationID == "" {
		return nil
	}
	return &ttnpb.ApplicationIdentifiers{ApplicationID: applicationID}
}

var (
	applicationsCommand = &cobra.Command{
		Use:     "applications",
		Aliases: []string{"application", "app", "a"},
		Short:   "Application commands",
	}
	applicationsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationFlags)
			if len(paths) == 0 {
				logger.Warnf("No fields selected, selecting %v", defaultGetPaths)
				paths = append(paths, defaultGetPaths...)
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationRegistryClient(is).List(ctx, &ttnpb.ListApplicationsRequest{
				Collaborator: getCollaborator(cmd.Flags()),
				FieldMask:    types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res.Applications)
		},
	}
	applicationsSearchCommand = &cobra.Command{
		Use:   "search",
		Short: "Search for applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationFlags)
			if len(paths) == 0 {
				logger.Warnf("No fields selected, selecting %v", defaultGetPaths)
				paths = append(paths, defaultGetPaths...)
			}
			req := getSearchEntitiesRequest(cmd.Flags())
			req.FieldMask.Paths = paths

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewEntityRegistrySearchClient(is).SearchApplications(ctx, req)
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res.Applications)
		},
	}
	applicationsGetCommand = &cobra.Command{
		Use:     "get",
		Aliases: []string{"info"},
		Short:   "Get an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationFlags)
			if len(paths) == 0 {
				logger.Warnf("No fields selected, selecting %v", defaultGetPaths)
				paths = append(paths, defaultGetPaths...)
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationRegistryClient(is).Get(ctx, &ttnpb.GetApplicationRequest{
				ApplicationIdentifiers: *appID,
				FieldMask:              types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	applicationsCreateCommand = &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "register"},
		Short:   "Create an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID
			}
			collaborator := getCollaborator(cmd.Flags())
			if collaborator == nil {
				return errNoCollaborator
			}
			var application ttnpb.Application
			util.SetFields(&application, setApplicationFlags)
			application.Attributes = mergeAttributes(application.Attributes, cmd.Flags())
			application.ApplicationIdentifiers = *appID

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationRegistryClient(is).Create(ctx, &ttnpb.CreateApplicationRequest{
				Application:  application,
				Collaborator: *collaborator,
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	applicationsUpdateCommand = &cobra.Command{
		Use:     "update",
		Aliases: []string{"set"},
		Short:   "Update an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID
			}
			paths := util.UpdateFieldMask(cmd.Flags(), setApplicationFlags, attributesFlags())
			if len(paths) == 0 {
				logger.Warn("No fields selected, won't update anything")
				return nil
			}
			var application ttnpb.Application
			util.SetFields(&application, setApplicationFlags)
			application.Attributes = mergeAttributes(application.Attributes, cmd.Flags())
			application.ApplicationIdentifiers = *appID

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewApplicationRegistryClient(is).Update(ctx, &ttnpb.UpdateApplicationRequest{
				Application: application,
				FieldMask:   types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	applicationsDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Delete an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewApplicationRegistryClient(is).Delete(ctx, appID)
			if err != nil {
				return err
			}

			return nil
		},
	}
	applicationsContactInfoCommand = contactInfoCommands("application", func(cmd *cobra.Command) (*ttnpb.EntityIdentifiers, error) {
		appID := getApplicationID(cmd.Flags(), nil)
		if appID == nil {
			return nil, errNoApplicationID
		}
		return appID.EntityIdentifiers(), nil
	})
)

func init() {
	applicationsListCommand.Flags().AddFlagSet(collaboratorFlags())
	applicationsListCommand.Flags().AddFlagSet(selectApplicationFlags)
	applicationsCommand.AddCommand(applicationsListCommand)
	applicationsSearchCommand.Flags().AddFlagSet(searchFlags())
	applicationsSearchCommand.Flags().AddFlagSet(selectApplicationFlags)
	applicationsCommand.AddCommand(applicationsSearchCommand)
	applicationsGetCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsGetCommand.Flags().AddFlagSet(selectApplicationFlags)
	applicationsCommand.AddCommand(applicationsGetCommand)
	applicationsCreateCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsCreateCommand.Flags().AddFlagSet(collaboratorFlags())
	applicationsCreateCommand.Flags().AddFlagSet(setApplicationFlags)
	applicationsCreateCommand.Flags().AddFlagSet(attributesFlags())
	applicationsCommand.AddCommand(applicationsCreateCommand)
	applicationsUpdateCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsUpdateCommand.Flags().AddFlagSet(setApplicationFlags)
	applicationsUpdateCommand.Flags().AddFlagSet(attributesFlags())
	applicationsCommand.AddCommand(applicationsUpdateCommand)
	applicationsDeleteCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsCommand.AddCommand(applicationsDeleteCommand)
	applicationsContactInfoCommand.PersistentFlags().AddFlagSet(applicationIDFlags())
	applicationsCommand.AddCommand(applicationsContactInfoCommand)
	Root.AddCommand(applicationsCommand)
}
