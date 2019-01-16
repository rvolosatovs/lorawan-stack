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

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/api"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/io"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

func getUserAndClientID(flagSet *pflag.FlagSet, args []string) (*ttnpb.UserIdentifiers, *ttnpb.ClientIdentifiers) {
	userID, _ := flagSet.GetString("user-id")
	clientID, _ := flagSet.GetString("client-id")
	switch len(args) {
	case 0:
	case 1:
		logger.Warn("Only single ID found in arguments, not considering arguments")
	case 2:
		userID = args[0]
		clientID = args[1]
	default:
		logger.Warn("multiple IDs found in arguments, considering the first")
		userID = args[0]
		clientID = args[1]
	}
	if userID == "" || clientID == "" {
		return nil, nil
	}
	return &ttnpb.UserIdentifiers{UserID: userID}, &ttnpb.ClientIdentifiers{ClientID: clientID}
}

var errNoTokenID = errors.DefineInvalidArgument("no_token_id", "no token ID set")

var (
	oauthCommand = &cobra.Command{
		Use:   "oauth",
		Short: "Manage OAuth authorizations and tokens",
	}
	oauthAuthorizationsCommand = &cobra.Command{
		Use:   "authorizations",
		Short: "Manage OAuth authorizations",
	}
	oauthAuthorizationsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List OAuth authorizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			usrID := getUserID(cmd.Flags(), args)
			if usrID == nil {
				return errNoUserID
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewOAuthAuthorizationRegistryClient(is).List(ctx, &ttnpb.ListOAuthClientAuthorizationsRequest{
				UserIdentifiers: *usrID,
			})

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	oauthAuthorizationsDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Delete an OAuth authorization and all tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			usrID, cliID := getUserAndClientID(cmd.Flags(), args)
			if usrID == nil {
				return errNoUserID
			}
			if cliID == nil {
				return errNoClientID
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}

			res, err := ttnpb.NewOAuthAuthorizationRegistryClient(is).ListTokens(ctx, &ttnpb.ListOAuthAccessTokensRequest{
				UserIDs:   *usrID,
				ClientIDs: *cliID,
			})
			for _, token := range res.Tokens {
				_, err = ttnpb.NewOAuthAuthorizationRegistryClient(is).DeleteToken(ctx, &ttnpb.OAuthAccessTokenIdentifiers{
					UserIDs:   *usrID,
					ClientIDs: *cliID,
					ID:        token.ID,
				})
				if err != nil {
					return err
				}
			}

			_, err = ttnpb.NewOAuthAuthorizationRegistryClient(is).Delete(ctx, &ttnpb.OAuthClientAuthorizationIdentifiers{
				UserIDs:   *usrID,
				ClientIDs: *cliID,
			})

			return err
		},
	}
	oauthTokensCommand = &cobra.Command{
		Use:   "tokens",
		Short: "Manage OAuth tokens",
	}
	oauthTokensListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List OAuth tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			usrID, cliID := getUserAndClientID(cmd.Flags(), args)
			if usrID == nil {
				return errNoUserID
			}
			if cliID == nil {
				return errNoClientID
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewOAuthAuthorizationRegistryClient(is).ListTokens(ctx, &ttnpb.ListOAuthAccessTokensRequest{
				UserIDs:   *usrID,
				ClientIDs: *cliID,
			})

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	oauthTokensDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Delete an OAuth token",
		RunE: func(cmd *cobra.Command, args []string) error {
			usrID, cliID := getUserAndClientID(cmd.Flags(), args)
			if usrID == nil {
				return errNoUserID
			}
			if cliID == nil {
				return errNoClientID
			}
			tokenID, _ := cmd.Flags().GetString("token-id")
			if tokenID == "" {
				return errNoTokenID
			}

			is, err := api.Dial(ctx, config.IdentityServerAddress)
			if err != nil {
				return err
			}

			_, err = ttnpb.NewOAuthAuthorizationRegistryClient(is).DeleteToken(ctx, &ttnpb.OAuthAccessTokenIdentifiers{
				UserIDs:   *usrID,
				ClientIDs: *cliID,
				ID:        tokenID,
			})

			return err
		},
	}
)
