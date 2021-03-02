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
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/v3/cmd/internal/io"
	is "go.thethings.network/lorawan-stack/v3/pkg/identityserver"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	createAPIKeyCommand = &cobra.Command{
		Use:   "create-user-api-key",
		Short: "Create an API key with full rights on the user in the Identity Server database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			logger.Info("Connecting to Identity Server database...")
			db, err := store.Open(ctx, config.IS.DatabaseURI)
			if err != nil {
				return err
			}
			defer db.Close()

			userID, err := cmd.Flags().GetString("user-id")
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")

			usr := &ttnpb.User{
				UserIdentifiers: ttnpb.UserIdentifiers{UserID: userID},
			}
			rights := []ttnpb.Right{ttnpb.RIGHT_ALL}
			apiKeyStore := store.GetAPIKeyStore(db)
			key, token, err := is.GenerateAPIKey(ctx, name, rights...)
			if err != nil {
				return err
			}
			key, err = apiKeyStore.CreateAPIKey(ctx, usr, key)
			if err != nil {
				return err
			}
			key.Key = token
			logger.Infof("API key ID: %s", key.ID)
			logger.Infof("API key value: %s", key.Key)
			logger.Warn("The API key value will never be shown again")
			logger.Warn("Make sure to copy it to a safe place")

			return io.Write(os.Stdout, config.OutputFormat, key)
		},
	}
)

func init() {
	createAPIKeyCommand.Flags().String("user-id", "admin", "User ID")
	createAPIKeyCommand.Flags().String("name", "admin-api-key", "API key name")
	isDBCommand.AddCommand(createAPIKeyCommand)
}
