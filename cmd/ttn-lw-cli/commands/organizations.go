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

package commands

import (
	"os"

	"github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/v3/cmd/internal/io"
	"go.thethings.network/lorawan-stack/v3/cmd/ttn-lw-cli/internal/api"
	"go.thethings.network/lorawan-stack/v3/cmd/ttn-lw-cli/internal/util"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	selectOrganizationFlags = util.FieldMaskFlags(&ttnpb.Organization{})
	setOrganizationFlags    = util.FieldFlags(&ttnpb.Organization{})

	selectAllOrganizationFlags = util.SelectAllFlagSet("organization")
)

func organizationIDFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.String("organization-id", "", "")
	return flagSet
}

var errNoOrganizationID = errors.DefineInvalidArgument("no_organization_id", "no organization ID set")

func getOrganizationID(flagSet *pflag.FlagSet, args []string) *ttnpb.OrganizationIdentifiers {
	var organizationID string
	if len(args) > 0 {
		if len(args) > 1 {
			logger.Warn("Multiple IDs found in arguments, considering only the first")
		}
		organizationID = args[0]
	} else {
		organizationID, _ = flagSet.GetString("organization-id")
	}
	if organizationID == "" {
		return nil
	}
	return &ttnpb.OrganizationIdentifiers{OrganizationID: organizationID}
}

var (
	organizationsCommand = &cobra.Command{
		Use:     "organizations",
		Aliases: []string{"organization", "org", "o"},
		Short:   "Organization commands",
	}
	organizationsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectOrganizationFlags)
			paths = ttnpb.AllowedFields(paths, ttnpb.RPCFieldMaskPaths["/ttn.lorawan.v3.OrganizationRegistry/List"].Allowed)

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			limit, page, opt, getTotal := withPagination(cmd.Flags())
			res, err := ttnpb.NewOrganizationRegistryClient(is).List(ctx, &ttnpb.ListOrganizationsRequest{
				Collaborator:   getUserID(cmd.Flags(), nil).GetOrganizationOrUserIdentifiers(),
				FieldMask:      types.FieldMask{Paths: paths},
				Limit:          limit,
				Page:           page,
				Order:          getOrder(cmd.Flags()),
				IncludeDeleted: getIncludeDeleted(cmd.Flags()),
			}, opt)
			if err != nil {
				return err
			}
			getTotal()

			return io.Write(os.Stdout, config.OutputFormat, res.Organizations)
		},
	}
	organizationsSearchCommand = &cobra.Command{
		Use:   "search",
		Short: "Search for organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectOrganizationFlags)
			paths = ttnpb.AllowedFields(paths, ttnpb.RPCFieldMaskPaths["/ttn.lorawan.v3.EntityRegistrySearch/SearchOrganizations"].Allowed)

			req, opt, getTotal := getSearchEntitiesRequest(cmd.Flags())
			req.FieldMask.Paths = paths

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewEntityRegistrySearchClient(is).SearchOrganizations(ctx, req, opt)
			if err != nil {
				return err
			}
			getTotal()

			return io.Write(os.Stdout, config.OutputFormat, res.Organizations)
		},
	}
	organizationsGetCommand = &cobra.Command{
		Use:     "get [organization-id]",
		Aliases: []string{"info"},
		Short:   "Get an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID := getOrganizationID(cmd.Flags(), args)
			if orgID == nil {
				return errNoOrganizationID
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectOrganizationFlags)
			paths = ttnpb.AllowedFields(paths, ttnpb.RPCFieldMaskPaths["/ttn.lorawan.v3.OrganizationRegistry/Get"].Allowed)

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewOrganizationRegistryClient(is).Get(ctx, &ttnpb.GetOrganizationRequest{
				OrganizationIdentifiers: *orgID,
				FieldMask:               types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	organizationsCreateCommand = &cobra.Command{
		Use:     "create [organization-id]",
		Aliases: []string{"add", "register"},
		Short:   "Create an organization",
		RunE: asBulk(func(cmd *cobra.Command, args []string) (err error) {
			orgID := getOrganizationID(cmd.Flags(), args)
			collaborator := getUserID(cmd.Flags(), nil).GetOrganizationOrUserIdentifiers()
			if collaborator == nil {
				return errNoCollaborator
			}
			var organization ttnpb.Organization
			if inputDecoder != nil {
				_, err := inputDecoder.Decode(&organization)
				if err != nil {
					return err
				}
			}
			if err := util.SetFields(&organization, setOrganizationFlags); err != nil {
				return err
			}
			organization.Attributes = mergeAttributes(organization.Attributes, cmd.Flags())
			if orgID != nil && orgID.OrganizationID != "" {
				organization.OrganizationID = orgID.OrganizationID
			}
			if organization.OrganizationID == "" {
				return errNoOrganizationID
			}

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewOrganizationRegistryClient(is).Create(ctx, &ttnpb.CreateOrganizationRequest{
				Organization: organization,
				Collaborator: *collaborator,
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		}),
	}
	organizationsSetCommand = &cobra.Command{
		Use:     "set [organization-id]",
		Aliases: []string{"set"},
		Short:   "Set properties of an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID := getOrganizationID(cmd.Flags(), args)
			if orgID == nil {
				return errNoOrganizationID
			}
			paths := util.UpdateFieldMask(cmd.Flags(), setOrganizationFlags, attributesFlags())
			if len(paths) == 0 {
				logger.Warn("No fields selected, won't update anything")
				return nil
			}
			var organization ttnpb.Organization
			if err := util.SetFields(&organization, setOrganizationFlags); err != nil {
				return err
			}
			organization.Attributes = mergeAttributes(organization.Attributes, cmd.Flags())
			organization.OrganizationIdentifiers = *orgID

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewOrganizationRegistryClient(is).Update(ctx, &ttnpb.UpdateOrganizationRequest{
				Organization: organization,
				FieldMask:    types.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			res.SetFields(&organization, "ids")
			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	organizationsDeleteCommand = &cobra.Command{
		Use:     "delete [organization-id]",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID := getOrganizationID(cmd.Flags(), args)
			if orgID == nil {
				return errNoOrganizationID
			}

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewOrganizationRegistryClient(is).Delete(ctx, orgID)
			if err != nil {
				return err
			}

			return nil
		},
	}

	organizationsPurgeCommand = &cobra.Command{
		Use:     "purge [organization-id]",
		Aliases: []string{"permanent-delete", "hard-delete"},
		Short:   "Purge an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID := getOrganizationID(cmd.Flags(), args)
			if orgID == nil {
				return errNoOrganizationID
			}
			force, err := cmd.Flags().GetBool("force")
			if err != nil {
				return err
			}
			if !confirmChoice(organizationPurgeWarning, force) {
				return errNoConfirmation
			}
			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewOrganizationRegistryClient(is).Purge(ctx, orgID)
			if err != nil {
				return err
			}

			return nil
		},
	}

	organizationsContactInfoCommand = contactInfoCommands("organization", func(cmd *cobra.Command, args []string) (*ttnpb.EntityIdentifiers, error) {
		orgID := getOrganizationID(cmd.Flags(), args)
		if orgID == nil {
			return nil, errNoOrganizationID
		}
		return orgID.EntityIdentifiers(), nil
	})
)

func init() {
	organizationsListCommand.Flags().AddFlagSet(collaboratorFlags())
	organizationsListCommand.Flags().AddFlagSet(includeDeletedFlags)
	organizationsListCommand.Flags().AddFlagSet(selectOrganizationFlags)
	organizationsListCommand.Flags().AddFlagSet(selectAllOrganizationFlags)
	organizationsListCommand.Flags().AddFlagSet(paginationFlags())
	organizationsListCommand.Flags().AddFlagSet(orderFlags())
	organizationsCommand.AddCommand(organizationsListCommand)
	organizationsSearchCommand.Flags().AddFlagSet(searchFlags())
	organizationsSearchCommand.Flags().AddFlagSet(includeDeletedFlags)
	organizationsSearchCommand.Flags().AddFlagSet(selectOrganizationFlags)
	organizationsSearchCommand.Flags().AddFlagSet(selectAllOrganizationFlags)
	organizationsCommand.AddCommand(organizationsSearchCommand)
	organizationsGetCommand.Flags().AddFlagSet(organizationIDFlags())
	organizationsGetCommand.Flags().AddFlagSet(selectOrganizationFlags)
	organizationsGetCommand.Flags().AddFlagSet(selectAllOrganizationFlags)
	organizationsCommand.AddCommand(organizationsGetCommand)
	organizationsCreateCommand.Flags().AddFlagSet(organizationIDFlags())
	organizationsCreateCommand.Flags().AddFlagSet(collaboratorFlags())
	organizationsCreateCommand.Flags().AddFlagSet(setOrganizationFlags)
	organizationsCreateCommand.Flags().AddFlagSet(attributesFlags())
	organizationsCommand.AddCommand(organizationsCreateCommand)
	organizationsSetCommand.Flags().AddFlagSet(organizationIDFlags())
	organizationsSetCommand.Flags().AddFlagSet(setOrganizationFlags)
	organizationsSetCommand.Flags().AddFlagSet(attributesFlags())
	organizationsCommand.AddCommand(organizationsSetCommand)
	organizationsDeleteCommand.Flags().AddFlagSet(organizationIDFlags())
	organizationsCommand.AddCommand(organizationsDeleteCommand)
	organizationsContactInfoCommand.PersistentFlags().AddFlagSet(organizationIDFlags())
	organizationsCommand.AddCommand(organizationsContactInfoCommand)
	organizationsPurgeCommand.Flags().AddFlagSet(organizationIDFlags())
	organizationsPurgeCommand.Flags().AddFlagSet(forceFlags())
	organizationsCommand.AddCommand(organizationsPurgeCommand)
	Root.AddCommand(organizationsCommand)
}
