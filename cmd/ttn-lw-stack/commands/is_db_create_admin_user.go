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
	"context"
	"os"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/pkg/auth"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var errPasswordMismatch = errors.DefineInvalidArgument("password_mismatch", "password did not match")

var (
	createAdminUserCommand = &cobra.Command{
		Use:   "create-admin-user",
		Short: "Create an admin user in the Identity Server database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			logger.Info("Connecting to Identity Server database...")
			db, err := store.Open(ctx, config.IS.DatabaseURI)
			if err != nil {
				return err
			}
			defer db.Close()

			userID, err := cmd.Flags().GetString("id")
			if err != nil {
				return err
			}
			email, err := cmd.Flags().GetString("email")
			if err != nil {
				return err
			}
			if email == "" {
				return errMissingFlag.WithAttributes("flag", "email")
			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			}
			if password == "" {
				pw, err := gopass.GetPasswdPrompt("Please enter user password:", true, os.Stdin, os.Stderr)
				if err != nil {
					return err
				}
				password = string(pw)
				pw, err = gopass.GetPasswdPrompt("Please repeat user password:", true, os.Stdin, os.Stderr)
				if err != nil {
					return err
				}
				if string(pw) != password {
					return errPasswordMismatch
				}
			}
			if password == "" {
				return errMissingFlag.WithAttributes("flag", "password")
			}
			hashedPassword, err := auth.Hash(ctx, password)
			if err != nil {
				return err
			}

			now := time.Now()

			usrFieldMask := &pbtypes.FieldMask{Paths: []string{
				"primary_email_address",
				"primary_email_address_validated_at",
				"password",
				"password_updated_at",
				"state",
				"admin",
			}}
			usr := &ttnpb.User{
				UserIdentifiers: ttnpb.UserIdentifiers{UserID: userID},
			}

			usrStore := store.GetUserStore(db)

			var usrExists bool
			if _, err := usrStore.GetUser(ctx, &usr.UserIdentifiers, usrFieldMask); err == nil {
				usrExists = true
			}
			usr.PrimaryEmailAddress = email
			usr.PrimaryEmailAddressValidatedAt = &now
			usr.Password = hashedPassword
			usr.PasswordUpdatedAt = &now
			usr.State = ttnpb.STATE_APPROVED
			usr.Admin = true

			if usrExists {
				logger.Info("Updating user...")
				if _, err = usrStore.UpdateUser(ctx, usr, usrFieldMask); err != nil {
					return err
				}
				logger.Info("Updated user")
			} else {
				logger.Info("Creating user...")
				if _, err = usrStore.CreateUser(ctx, usr); err != nil {
					return err
				}
				logger.Info("Created user")
			}

			return nil
		},
	}
)

func init() {
	createAdminUserCommand.Flags().String("id", "admin", "User ID")
	createAdminUserCommand.Flags().String("email", "", "Email address")
	createAdminUserCommand.Flags().String("password", "", "Password")
	isDBCommand.AddCommand(createAdminUserCommand)
}
