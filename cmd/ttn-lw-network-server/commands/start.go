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
	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/cmd/internal/shared"
	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/networkserver"
	nsredis "go.thethings.network/lorawan-stack/pkg/networkserver/redis"
	"go.thethings.network/lorawan-stack/pkg/redis"
)

var (
	startCommand = &cobra.Command{
		Use:   "start",
		Short: "Start the Network Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := component.New(logger, &component.Config{ServiceBase: config.ServiceBase})
			if err != nil {
				return shared.ErrBaseComponentInitialize.WithCause(err)
			}

			config.NS.Devices = &nsredis.DeviceRegistry{Redis: redis.New(&redis.Config{
				Redis:     config.Redis,
				Namespace: []string{"ns", "devices"},
			})}

			ns, err := networkserver.New(c, &config.NS)
			if err != nil {
				return shared.ErrNetworkServerInitialize.WithCause(err)
			}
			_ = ns

			logger.Info("Starting Network Server...")
			return c.Run()
		},
	}
)

func init() {
	Root.AddCommand(startCommand)
}
