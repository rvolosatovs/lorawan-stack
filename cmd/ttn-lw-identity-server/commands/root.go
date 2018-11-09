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

// Package commands implements the commands for the ttn-lw-identity-server binary.
package commands

import (
	"os"

	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/cmd/internal/shared"
	conf "go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/log"
)

var (
	logger *log.Logger
	name   = "ttn-lw-identity-server"
	mgr    = conf.InitializeWithDefaults(name, "ttn_lw", DefaultConfig)
	config = new(Config)

	// Root command is the entrypoint of the program
	Root = &cobra.Command{
		Use:           name,
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "The Things Network Identity Server for LoRaWAN",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Read in config from file.
			err := mgr.ReadInConfig()
			if err != nil {
				return err
			}

			// Unmarshal config.
			if err = mgr.Unmarshal(config); err != nil {
				return err
			}

			// Create logger.
			logger, err = log.NewLogger(
				log.WithLevel(config.Log.Level),
				log.WithHandler(log.NewCLI(os.Stdout)),
			)
			if sentry, err := shared.SentryMiddleware(config.ServiceBase); err == nil && sentry != nil {
				logger.Use(sentry)
			}

			// Initialize shared packages.
			if err := shared.Initialize(config.ServiceBase); err != nil {
				return err
			}

			return err
		},
	}
)

func init() {
	Root.PersistentFlags().AddFlagSet(mgr.Flags())
}
