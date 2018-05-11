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
	"go.thethings.network/lorawan-stack/cmd/internal/shared"
	shared_applicationserver "go.thethings.network/lorawan-stack/cmd/internal/shared/applicationserver"
	"go.thethings.network/lorawan-stack/pkg/applicationserver"
	conf "go.thethings.network/lorawan-stack/pkg/config"
)

// Config for the ttn-lw-application-server binary.
type Config struct {
	conf.ServiceBase `name:",squash"`
	AS               applicationserver.Config `name:"as"`
}

// DefaultConfig contains the default config for the ttn-lw-application-server binary.
var DefaultConfig = Config{
	ServiceBase: shared.DefaultServiceBase,
	AS:          shared_applicationserver.DefaultApplicationServerConfig,
}
