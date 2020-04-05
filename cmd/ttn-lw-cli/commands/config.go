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
	"fmt"

	"go.thethings.network/lorawan-stack/cmd/internal/commands"
	conf "go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/log"
)

var (
	defaultClusterHost = "localhost"
	defaultInsecure    = false
)

// Config for the ttn-lw-cli binary.
type Config struct {
	conf.Base                          `name:",squash"`
	CredentialsID                      string `name:"credentials-id" yaml:"credentials-id" description:"Credentials ID (if using multiple configurations)"`
	InputFormat                        string `name:"input-format" yaml:"input-format" description:"Input format"`
	OutputFormat                       string `name:"output-format" yaml:"output-format" description:"Output format"`
	AllowUnknownHosts                  bool   `name:"allow-unknown-hosts" yaml:"allow-unknown-hosts" description:"Allow sending credentials to unknown hosts"`
	OAuthServerAddress                 string `name:"oauth-server-address" yaml:"oauth-server-address" description:"OAuth Server address"`
	IdentityServerGRPCAddress          string `name:"identity-server-grpc-address" yaml:"identity-server-grpc-address" description:"Identity Server address"`
	GatewayServerEnabled               bool   `name:"gateway-server-enabled" yaml:"gateway-server-enabled" description:"Gateway Server enabled"`
	GatewayServerGRPCAddress           string `name:"gateway-server-grpc-address" yaml:"gateway-server-grpc-address" description:"Gateway Server address"`
	NetworkServerEnabled               bool   `name:"network-server-enabled" yaml:"network-server-enabled" description:"Network Server enabled"`
	NetworkServerGRPCAddress           string `name:"network-server-grpc-address" yaml:"network-server-grpc-address" description:"Network Server address"`
	ApplicationServerEnabled           bool   `name:"application-server-enabled" yaml:"application-server-enabled" description:"Application Server enabled"`
	ApplicationServerGRPCAddress       string `name:"application-server-grpc-address" yaml:"application-server-grpc-address" description:"Application Server address"`
	JoinServerEnabled                  bool   `name:"join-server-enabled" yaml:"join-server-enabled" description:"Join Server enabled"`
	JoinServerGRPCAddress              string `name:"join-server-grpc-address" yaml:"join-server-grpc-address" description:"Join Server address"`
	DeviceTemplateConverterGRPCAddress string `name:"device-template-converter-grpc-address" yaml:"device-template-converter-grpc-address" description:"Device Template Converter address"`
	DeviceClaimingServerGRPCAddress    string `name:"device-claiming-server-grpc-address" yaml:"device-claiming-server-grpc-address" description:"Device Claiming Server address"`
	QRCodeGeneratorGRPCAddress         string `name:"qr-code-generator-grpc-address" yaml:"qr-code-generator-grpc-address" description:"QR Code Generator address"`
	Insecure                           bool   `name:"insecure" yaml:"insecure" description:"Connect without TLS"`
	CA                                 string `name:"ca" yaml:"ca" description:"CA certificate file"`
	DumpRequests                       bool   `name:"dump-requests" yaml:"dump-requests" description:"When log level is set to debug, also dump request payload as JSON"`
}

func (c Config) getHosts() []string {
	hosts := make([]string, 0, 7)
	hosts = append(hosts, c.OAuthServerAddress)
	hosts = append(hosts, c.IdentityServerGRPCAddress)
	if c.GatewayServerEnabled {
		hosts = append(hosts, c.GatewayServerGRPCAddress)
	}
	if c.NetworkServerEnabled {
		hosts = append(hosts, c.NetworkServerGRPCAddress)
	}
	if c.ApplicationServerEnabled {
		hosts = append(hosts, c.ApplicationServerGRPCAddress)
	}
	if c.JoinServerEnabled {
		hosts = append(hosts, c.JoinServerGRPCAddress)
	}
	hosts = append(hosts, c.DeviceTemplateConverterGRPCAddress)
	hosts = append(hosts, c.DeviceClaimingServerGRPCAddress)
	return getHosts(hosts...)
}

// MakeDefaultConfig builds the default config for the ttn-lw-cli binary for a given host.
func MakeDefaultConfig(clusterHost string, insecure bool) Config {
	clusterGRPCAddress := fmt.Sprintf("%s:8884", clusterHost)
	clusterHTTPAddress := fmt.Sprintf("https://%s:8885", clusterHost)
	if insecure {
		clusterGRPCAddress = fmt.Sprintf("%s:1884", clusterHost)
		clusterHTTPAddress = fmt.Sprintf("http://%s:1885", clusterHost)
	}

	return Config{
		Base: conf.Base{
			Log: conf.Log{
				Level: log.InfoLevel,
			},
		},
		InputFormat:                        "json",
		OutputFormat:                       "json",
		OAuthServerAddress:                 clusterHTTPAddress + "/oauth",
		IdentityServerGRPCAddress:          clusterGRPCAddress,
		GatewayServerEnabled:               true,
		GatewayServerGRPCAddress:           clusterGRPCAddress,
		NetworkServerEnabled:               true,
		NetworkServerGRPCAddress:           clusterGRPCAddress,
		ApplicationServerEnabled:           true,
		ApplicationServerGRPCAddress:       clusterGRPCAddress,
		JoinServerEnabled:                  true,
		JoinServerGRPCAddress:              clusterGRPCAddress,
		DeviceTemplateConverterGRPCAddress: clusterGRPCAddress,
		DeviceClaimingServerGRPCAddress:    clusterGRPCAddress,
		QRCodeGeneratorGRPCAddress:         clusterGRPCAddress,
		Insecure:                           insecure,
	}
}

// DefaultConfig contains the default config for the ttn-lw-cli binary.
var DefaultConfig = MakeDefaultConfig(defaultClusterHost, defaultInsecure)

var configCommand = commands.Config(mgr)

func init() {
	configCommand.PersistentPreRunE = preRun()
	Root.AddCommand(configCommand)
}
