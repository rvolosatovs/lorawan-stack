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

package gatewayserver

import "go.thethings.network/lorawan-stack/pkg/gatewayserver/io/udp"

// MQTTConfig contains MQTT configuration of the Gateway Server.
type MQTTConfig struct {
	Listen    string `name:"listen" description:"Address for the MQTT frontend to listen on"`
	ListenTLS string `name:"listen-tls" description:"Address for the MQTTS frontend to listen on"`
}

// UDPConfig defines the UDP configuration of the Gateway Server.
type UDPConfig struct {
	udp.Config
	Listeners []UDPListener `name:"listeners" description:"Listener configuration"`
}

// UDPListener defines a UDP listener of the Gateway Server.
type UDPListener struct {
	Listen                  string `name:"listen" description:"Address for the UDP frontend to listen on"`
	FallbackFrequencyPlanID string `name:"fallback-frequency-plan-id" description:"Frequency plan ID when the gateway is not registered"`
}

// Config represents the Gateway Server configuration.
type Config struct {
	MQTT MQTTConfig `name:"mqtt"`
	UDP  UDPConfig  `name:"udp"`
}
