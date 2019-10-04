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

package shared

import (
	"fmt"

	"go.thethings.network/lorawan-stack/cmd/internal/shared"
	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io/udp"
)

// DefaultGatewayServerConfig is the default configuration for the GatewayServer.
var DefaultGatewayServerConfig = gatewayserver.Config{
	RequireRegisteredGateways: false,
	Forward: map[string][]string{
		"": {"00000000/0"},
	},
	UDP: gatewayserver.UDPConfig{
		Config: udp.DefaultConfig,
		Listeners: map[string]string{
			":1700": "",
		},
	},
	MQTTV2: config.MQTT{
		Listen:           ":1881",
		ListenTLS:        ":8881",
		PublicAddress:    fmt.Sprintf("%s:1881", shared.DefaultPublicHost),
		PublicTLSAddress: fmt.Sprintf("%s:8881", shared.DefaultPublicHost),
	},
	MQTT: config.MQTT{
		Listen:           ":1882",
		ListenTLS:        ":8882",
		PublicAddress:    fmt.Sprintf("%s:1882", shared.DefaultPublicHost),
		PublicTLSAddress: fmt.Sprintf("%s:8882", shared.DefaultPublicHost),
	},
	BasicStation: gatewayserver.BasicStationConfig{
		Listen:    ":1887",
		ListenTLS: ":8887",
	},
}
