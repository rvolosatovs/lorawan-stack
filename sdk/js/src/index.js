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

import Applications from './service/applications'
import Application from './entity/application'
import Api from './api'

class TtnLw {
  constructor (token, {
    stackConfig,
    connectionType,
    defaultUserId,
    proxy,
    axiosConfig,
  }) {
    this.config = arguments.config
    this.api = new Api(connectionType, stackConfig, axiosConfig, token)

    this.Applications = new Applications(this.api, { defaultUserId, proxy })
    this.Application = Application.bind(null, this.Applications)
  }
}

export default TtnLw
