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
	selectClientFlags = util.FieldMaskFlags(&ttnpb.Client{})
	setClientFlags    = util.FieldFlags(&ttnpb.Client{})
)

func clientIDFlags() *pflag.FlagSet {
	flagSet := new(pflag.FlagSet)
	flagSet.String("client-id", "", "")
	return flagSet
}

var errNoClientID = errors.DefineInvalidArgument("no_client_id", "no client ID set")

func getClientID(flagSet *pflag.FlagSet, args []string) *ttnpb.ClientIdentifiers {
	var clientID string
	if len(args) > 0 {
		if len(args) > 1 {
			logger.Warn("multiple IDs found in arguments, considering only the first")
		}
		clientID = args[0]
	} else {
		clientID, _ = flagSet.GetString("client-id")
	}
	if clientID == "" {
		return nil
	}
	return &ttnpb.ClientIdentifiers{ClientID: clientID}
}

var (
	clientsCommand = &cobra.Command{
		Use:     "clients",
		Aliases: []string{"client", "cli", "c"},
		Short:   "Client commands",
	}
	clientsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectClientFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectClientFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, flag.Name)
				})
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewClientRegistryClient(is).List(ctx, &ttnpb.ListClientsRequest{
				Collaborator: getCollaborator(cmd.Flags()),
				FieldMask:    types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res.Clients)
		},
	}
	clientsSearchCommand = &cobra.Command{
		Use:   "search",
		Short: "Search for clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectClientFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectClientFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, flag.Name)
				})
			}
			req := getSearchEntitiesRequest(cmd.Flags())
			req.FieldMask.Paths = paths

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewEntityRegistrySearchClient(is).SearchClients(ctx, req)
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res.Clients)
		},
	}
	clientsGetCommand = &cobra.Command{
		Use:     "get",
		Aliases: []string{"info"},
		Short:   "Get a client",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliID := getClientID(cmd.Flags(), args)
			if cliID == nil {
				return errNoClientID
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectClientFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectClientFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, flag.Name)
				})
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewClientRegistryClient(is).Get(ctx, &ttnpb.GetClientRequest{
				ClientIdentifiers: *cliID,
				FieldMask:         types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	clientsCreateCommand = &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "register"},
		Short:   "Create a client",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliID := getClientID(cmd.Flags(), args)
			if cliID == nil {
				return errNoClientID
			}
			collaborator := getCollaborator(cmd.Flags())
			if collaborator == nil {
				return errNoCollaborator
			}
			var client ttnpb.Client
			util.SetFields(&client, setClientFlags)
			client.Attributes = mergeAttributes(client.Attributes, cmd.Flags())
			client.ClientIdentifiers = *cliID

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewClientRegistryClient(is).Create(ctx, &ttnpb.CreateClientRequest{
				Client:       client,
				Collaborator: *collaborator,
			})
			if err != nil {
				return err
			}

			logger.Infof("Client secret: %s", res.Secret)
			logger.Warn("The Client secret will never be shown again")
			logger.Warn("Make sure to copy it to a safe place")

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	clientsUpdateCommand = &cobra.Command{
		Use:     "update",
		Aliases: []string{"set"},
		Short:   "Update a client",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliID := getClientID(cmd.Flags(), args)
			if cliID == nil {
				return errNoClientID
			}
			paths := util.UpdateFieldMask(cmd.Flags(), setClientFlags, attributesFlags())
			if len(paths) == 0 {
				logger.Warn("No fields selected, won't update anything")
				return nil
			}
			var client ttnpb.Client
			util.SetFields(&client, setClientFlags)
			client.Attributes = mergeAttributes(client.Attributes, cmd.Flags())
			client.ClientIdentifiers = *cliID

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewClientRegistryClient(is).Update(ctx, &ttnpb.UpdateClientRequest{
				Client:    client,
				FieldMask: types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.Format, res)
		},
	}
	clientsDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Delete a client",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliID := getClientID(cmd.Flags(), args)
			if cliID == nil {
				return errNoClientID
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewClientRegistryClient(is).Delete(ctx, cliID)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	clientsListCommand.Flags().AddFlagSet(collaboratorFlags())
	clientsListCommand.Flags().AddFlagSet(selectClientFlags)
	clientsCommand.AddCommand(clientsListCommand)
	clientsSearchCommand.Flags().AddFlagSet(searchFlags())
	clientsSearchCommand.Flags().AddFlagSet(selectClientFlags)
	clientsCommand.AddCommand(clientsSearchCommand)
	clientsGetCommand.Flags().AddFlagSet(clientIDFlags())
	clientsGetCommand.Flags().AddFlagSet(selectClientFlags)
	clientsCommand.AddCommand(clientsGetCommand)
	clientsCreateCommand.Flags().AddFlagSet(clientIDFlags())
	clientsCreateCommand.Flags().AddFlagSet(collaboratorFlags())
	clientsCreateCommand.Flags().AddFlagSet(setClientFlags)
	clientsCreateCommand.Flags().AddFlagSet(attributesFlags())
	clientsCommand.AddCommand(clientsCreateCommand)
	clientsUpdateCommand.Flags().AddFlagSet(clientIDFlags())
	clientsUpdateCommand.Flags().AddFlagSet(setClientFlags)
	clientsUpdateCommand.Flags().AddFlagSet(attributesFlags())
	clientsCommand.AddCommand(clientsUpdateCommand)
	clientsDeleteCommand.Flags().AddFlagSet(clientIDFlags())
	clientsCommand.AddCommand(clientsDeleteCommand)
	Root.AddCommand(clientsCommand)
}
