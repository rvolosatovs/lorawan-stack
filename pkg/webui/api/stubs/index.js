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

import fakeData from './fake-data'

const genericSearch = function (entity, params, predicate = () => true) {
  const start = (params.page - 1) * params.pageSize
  const end = start + params.pageSize

  const res = fakeData[entity].filter(predicate)
  const total = res.length
  const delay = Math.floor(Math.random() * (800 - 100)) + 100

  return new Promise(resolve => setTimeout(() => resolve(
    { [entity]: res.slice(start, end), totalCount: total }
  ), delay))
}

export default {
  applications: {
    list (params) {
      return genericSearch('applications', params)
    },
    search (params) {
      const query = params.query || ''

      return genericSearch(
        'applications',
        params,
        app => app.ids.application_id.includes(query)
      )
    },
  },
  application: {
    get (id) {
      const app = fakeData.applications.find(a => a.ids.application_id === id)

      return new Promise((resolve, reject) => setTimeout(function () {
        if (app) {
          resolve(app)
        } else {
          reject(new Error())
        }
      }, 750))
    },
    apiKeys: {
      list (applicationId, params) {
        return genericSearch(
          'applicationsApiKeys',
          params,
          k => k.application_id === applicationId
        )
      },
    },
  },
  devices: {
    list (appId, params) {
      return genericSearch('devices', params, d => d.ids.application_ids.application_id === appId)
    },
    search (appId, params) {
      const query = params.query || ''

      return genericSearch(
        'devices',
        params,
        d => d.ids.application_ids.application_id === appId && d.ids.device_id.includes(query),
      )
    },
  },
  gateways: {
    list (params) {
      return genericSearch('gateways', params)
    },
    search (params) {
      const query = params.query || ''

      return genericSearch(
        'gateways',
        params,
        gtw => gtw.gateway_id.includes(query)
      )
    },
  },
  rights: {
    applications () {
      return new Promise(resolve => setTimeout(() => resolve(
        fakeData.rights.applications
      ), 500))
    },
  },
}
