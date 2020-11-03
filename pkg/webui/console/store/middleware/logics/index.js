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

import user from './user'
import users from './users'
import init from './init'
import applications from './applications'
import devices from './devices'
import gateways from './gateways'
import configuration from './configuration'
import organizations from './organizations'
import js from './join-server'
import apiKeys from './api-keys'
import collaborators from './collaborators'
import webhooks from './webhooks'
import pubsubs from './pubsubs'
import applicationPackages from './application-packages'
import is from './identity-server'

export default [
  ...user,
  ...users,
  ...init,
  ...applications,
  ...devices,
  ...gateways,
  ...configuration,
  ...organizations,
  ...js,
  ...apiKeys,
  ...collaborators,
  ...webhooks,
  ...pubsubs,
  ...applicationPackages,
  ...is,
]
