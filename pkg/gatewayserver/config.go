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

package gatewayserver

import (
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/udp"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/ws"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
)

// UDPConfig defines the UDP configuration of the Gateway Server.
type UDPConfig struct {
	udp.Config `name:",squash"`
	Listeners  map[string]string `name:"listeners" description:"Listen addresses with (optional) fallback frequency plan ID for non-registered gateways"`
}

// BasicStationConfig defines the LoRa Basics Station configuration of the Gateway Server.
type BasicStationConfig struct {
	ws.Config               `name:",squash"`
	MaxValidRoundTripDelay  time.Duration `name:"max-valid-round-trip-delay" description:"Maximum valid round trip delay to qualify for RTT calculations"`
	FallbackFrequencyPlanID string        `name:"fallback-frequency-plan-id" description:"Fallback frequency plan ID for non-registered gateways"`
	Listen                  string        `name:"listen" description:"Address for the Basic Station frontend to listen on"`
	ListenTLS               string        `name:"listen-tls" description:"Address for the Basic Station frontend to listen on (with TLS)"`
}

// Config represents the Gateway Server configuration.
type Config struct {
	RequireRegisteredGateways         bool          `name:"require-registered-gateways" description:"Require the gateways to be registered in the Identity Server"`
	UpdateGatewayLocationDebounceTime time.Duration `name:"update-gateway-location-debounce-time" description:"Debounce time for gateway location updates from status messages"`

	Stats                             GatewayConnectionStatsRegistry `name:"-"`
	UpdateConnectionStatsDebounceTime time.Duration                  `name:"update-connection-stats-debounce-time" description:"Time before repeated refresh of the gateway connection stats"`

	Forward map[string][]string `name:"forward" description:"Forward the DevAddr prefixes to the specified hosts"`

	MQTT         config.MQTT        `name:"mqtt"`
	MQTTV2       config.MQTT        `name:"mqtt-v2"`
	UDP          UDPConfig          `name:"udp"`
	BasicStation BasicStationConfig `name:"basic-station"`
}

// ForwardDevAddrPrefixes parses the configured forward map.
func (c Config) ForwardDevAddrPrefixes() (map[string][]types.DevAddrPrefix, error) {
	res := make(map[string][]types.DevAddrPrefix, len(c.Forward))
	for host, prefixes := range c.Forward {
		res[host] = make([]types.DevAddrPrefix, 0, len(prefixes))
		for _, val := range prefixes {
			var prefix types.DevAddrPrefix
			if err := prefix.UnmarshalText([]byte(val)); err != nil {
				return nil, err
			}
			res[host] = append(res[host], prefix)
		}
	}
	return res, nil
}
