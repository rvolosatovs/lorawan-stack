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

import axios from 'axios'

import TTN, { STACK_COMPONENTS_MAP } from 'ttn-lw'

import toast from '@ttn-lw/components/toast'

import {
  selectStackConfig,
  selectApplicationRootPath,
  selectCSRFToken,
} from '@ttn-lw/lib/selectors/env'

import tokenCreator from '@console/lib/access-token'

const stackConfig = selectStackConfig()
const appRoot = selectApplicationRootPath()

const stack = {
  [STACK_COMPONENTS_MAP.is]: stackConfig.is.enabled ? stackConfig.is.base_url : undefined,
  [STACK_COMPONENTS_MAP.gs]: stackConfig.gs.enabled ? stackConfig.gs.base_url : undefined,
  [STACK_COMPONENTS_MAP.ns]: stackConfig.ns.enabled ? stackConfig.ns.base_url : undefined,
  [STACK_COMPONENTS_MAP.as]: stackConfig.as.enabled ? stackConfig.as.base_url : undefined,
  [STACK_COMPONENTS_MAP.js]: stackConfig.js.enabled ? stackConfig.js.base_url : undefined,
  [STACK_COMPONENTS_MAP.edtc]: stackConfig.edtc.enabled ? stackConfig.edtc.base_url : undefined,
  [STACK_COMPONENTS_MAP.qrg]: stackConfig.qrg.enabled ? stackConfig.qrg.base_url : undefined,
  [STACK_COMPONENTS_MAP.gcs]: stackConfig.gcs.enabled ? stackConfig.gcs.base_url : undefined,
}

const isBaseUrl = stackConfig.is.base_url

const csrf = selectCSRFToken()
const instance = axios.create()

instance.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if ('response' in error && error.response && 'data' in error.response) {
      throw error.response.data
    }
    throw error
  },
)

const token = tokenCreator(() => instance.get(`${appRoot}/api/auth/token`))

const ttnClient = new TTN(token, {
  stackConfig: stack,
  connectionType: 'http',
  proxy: false,
})

// Forward header warnings to the toast message queue.
ttnClient.subscribe('warning', payload => {
  toast({
    title: 'Warning',
    type: toast.types.WARNING,
    message: payload,
    preventConsecutive: true,
  })
})

export default {
  console: {
    token() {
      return instance.get(`${appRoot}/api/auth/token`)
    },
    async logout() {
      const headers = token => ({
        headers: { 'X-CSRF-Token': token },
      })
      try {
        return await axios.post(`${appRoot}/api/auth/logout`, undefined, headers(csrf))
      } catch (error) {
        if (
          error.response &&
          error.response.status === 403 &&
          typeof error.response.data === 'string' &&
          error.response.data.includes('CSRF')
        ) {
          // If the CSRF token is invalid, it likely means that the CSRF cookie
          // has been deleted or became outdated. Making a new request to the
          // current path can then retrieve a fresh CSRF cookie, with which
          // the logout can be retried.
          const csrfResult = await axios.get(window.location)
          const freshCsrf = csrfResult.headers['x-csrf-token']
          if (freshCsrf) {
            return axios.post(`${appRoot}/api/auth/logout`, undefined, headers(freshCsrf))
          }
        }

        throw error
      }
    },
  },
  clients: {
    get(client_id) {
      return instance.get(`${isBaseUrl}/is/clients/${client_id}`)
    },
  },
  users: {
    create: ttnClient.Users.create.bind(ttnClient.Users),
    get: ttnClient.Users.getById.bind(ttnClient.Users),
    list: ttnClient.Users.getAll.bind(ttnClient.Users),
    update: ttnClient.Users.updateById.bind(ttnClient.Users),
    delete: ttnClient.Users.deleteById.bind(ttnClient.Users),
    search: ttnClient.Users.search.bind(ttnClient.Users),
    authInfo: ttnClient.Auth.getAuthInfo.bind(ttnClient.Auth),
  },
  applications: {
    list: ttnClient.Applications.getAll.bind(ttnClient.Applications),
    search: ttnClient.Applications.search.bind(ttnClient.Applications),
  },
  application: {
    get: ttnClient.Applications.getById.bind(ttnClient.Applications),
    delete: ttnClient.Applications.deleteById.bind(ttnClient.Applications),
    create: ttnClient.Applications.create.bind(ttnClient.Applications),
    update: ttnClient.Applications.updateById.bind(ttnClient.Applications),
    eventsSubscribe: ttnClient.Applications.openStream.bind(ttnClient.Applications),
    getMqttConnectionInfo: ttnClient.Applications.getMqttConnectionInfo.bind(
      ttnClient.Applications,
    ),
    apiKeys: {
      get: ttnClient.Applications.ApiKeys.getById.bind(ttnClient.Applications.ApiKeys),
      list: ttnClient.Applications.ApiKeys.getAll.bind(ttnClient.Applications.ApiKeys),
      update: ttnClient.Applications.ApiKeys.updateById.bind(ttnClient.Applications.ApiKeys),
      delete: ttnClient.Applications.ApiKeys.deleteById.bind(ttnClient.Applications.ApiKeys),
      create: ttnClient.Applications.ApiKeys.create.bind(ttnClient.Applications.ApiKeys),
    },
    link: {
      get: ttnClient.Applications.Link.get.bind(ttnClient.Applications.Link),
      set: ttnClient.Applications.Link.set.bind(ttnClient.Applications.Link),
      delete: ttnClient.Applications.Link.delete.bind(ttnClient.Applications.Link),
      stats: ttnClient.Applications.Link.getStats.bind(ttnClient.Applications.Link),
    },
    collaborators: {
      getOrganization: ttnClient.Applications.Collaborators.getByOrganizationId.bind(
        ttnClient.Applications.Collaborators,
      ),
      getUser: ttnClient.Applications.Collaborators.getByUserId.bind(
        ttnClient.Applications.Collaborators,
      ),
      list: ttnClient.Applications.Collaborators.getAll.bind(ttnClient.Applications.Collaborators),
      add: ttnClient.Applications.Collaborators.add.bind(ttnClient.Applications.Collaborators),
      update: ttnClient.Applications.Collaborators.update.bind(
        ttnClient.Applications.Collaborators,
      ),
      remove: ttnClient.Applications.Collaborators.remove.bind(
        ttnClient.Applications.Collaborators,
      ),
    },
    webhooks: {
      list: ttnClient.Applications.Webhooks.getAll.bind(ttnClient.Applications.Webhooks),
      get: ttnClient.Applications.Webhooks.getById.bind(ttnClient.Applications.Webhooks),
      create: ttnClient.Applications.Webhooks.create.bind(ttnClient.Applications.Webhooks),
      update: ttnClient.Applications.Webhooks.updateById.bind(ttnClient.Applications.Webhooks),
      delete: ttnClient.Applications.Webhooks.deleteById.bind(ttnClient.Applications.Webhooks),
      getFormats: ttnClient.Applications.Webhooks.getFormats.bind(ttnClient.Applications.Webhooks),
      listTemplates: ttnClient.Applications.Webhooks.listTemplates.bind(
        ttnClient.Applications.Webhooks,
      ),
      getTemplate: ttnClient.Applications.Webhooks.getTemplate.bind(
        ttnClient.Applications.Webhooks,
      ),
    },
    pubsubs: {
      list: ttnClient.Applications.PubSubs.getAll.bind(ttnClient.Applications.PubSubs),
      get: ttnClient.Applications.PubSubs.getById.bind(ttnClient.Applications.PubSubs),
      create: ttnClient.Applications.PubSubs.create.bind(ttnClient.Applications.PubSubs),
      update: ttnClient.Applications.PubSubs.updateById.bind(ttnClient.Applications.PubSubs),
      delete: ttnClient.Applications.PubSubs.deleteById.bind(ttnClient.Applications.PubSubs),
      getFormats: ttnClient.Applications.PubSubs.getFormats.bind(ttnClient.Applications.PubSubs),
    },
    packages: {
      getDefaultAssociation: ttnClient.Applications.Packages.getDefaultAssociation.bind(
        ttnClient.Applications.Packages,
      ),
      setDefaultAssociation: ttnClient.Applications.Packages.setDefaultAssociation.bind(
        ttnClient.Applications.Packages,
      ),
    },
  },
  devices: {
    list: ttnClient.Applications.Devices.getAll.bind(ttnClient.Applications.Devices),
    search: ttnClient.Applications.Devices.search.bind(ttnClient.Applications.Devices),
  },
  device: {
    get: ttnClient.Applications.Devices.getById.bind(ttnClient.Applications.Devices),
    create: ttnClient.Applications.Devices.create.bind(ttnClient.Applications.Devices),
    bulkCreate: ttnClient.Applications.Devices.bulkCreate.bind(ttnClient.Applications.Devices),
    update: ttnClient.Applications.Devices.updateById.bind(ttnClient.Applications.Devices),
    eventsSubscribe: ttnClient.Applications.Devices.openStream.bind(ttnClient.Applications.Devices),
    delete: ttnClient.Applications.Devices.deleteById.bind(ttnClient.Applications.Devices),
    simulateUplink: ttnClient.Applications.Devices.simulateUplink.bind(
      ttnClient.Applications.Devices,
    ),
  },
  deviceTemplates: {
    listFormats: ttnClient.Applications.Devices.listTemplateFormats.bind(
      ttnClient.Applications.Devices,
    ),
    convert: ttnClient.Applications.Devices.convertTemplate.bind(ttnClient.Applications.Devices),
  },
  gateways: {
    list: ttnClient.Gateways.getAll.bind(ttnClient.Gateways),
    search: ttnClient.Gateways.search.bind(ttnClient.Gateways),
  },
  downlinkQueue: {
    list: ttnClient.Applications.Devices.DownlinkQueue.list.bind(
      ttnClient.Applications.Devices.DownlinkQueue,
    ),
    replace: ttnClient.Applications.Devices.DownlinkQueue.replace.bind(
      ttnClient.Applications.Devices.DownlinkQueue,
    ),
    push: ttnClient.Applications.Devices.DownlinkQueue.push.bind(
      ttnClient.Applications.Devices.DownlinkQueue,
    ),
  },
  gateway: {
    get: ttnClient.Gateways.getById.bind(ttnClient.Gateways),
    getGlobalConf: ttnClient.Gateways.getGlobalConf.bind(ttnClient.Gateways),
    delete: ttnClient.Gateways.deleteById.bind(ttnClient.Gateways),
    create: ttnClient.Gateways.create.bind(ttnClient.Gateways),
    update: ttnClient.Gateways.updateById.bind(ttnClient.Gateways),
    stats: ttnClient.Gateways.getStatisticsById.bind(ttnClient.Gateways),
    eventsSubscribe: ttnClient.Gateways.openStream.bind(ttnClient.Gateways),
    collaborators: {
      getOrganization: ttnClient.Gateways.Collaborators.getByOrganizationId.bind(
        ttnClient.Gateways.Collaborators,
      ),
      getUser: ttnClient.Gateways.Collaborators.getByUserId.bind(ttnClient.Gateways.Collaborators),
      list: ttnClient.Gateways.Collaborators.getAll.bind(ttnClient.Gateways.Collaborators),
      add: ttnClient.Gateways.Collaborators.add.bind(ttnClient.Gateways.Collaborators),
      update: ttnClient.Gateways.Collaborators.update.bind(ttnClient.Gateways.Collaborators),
      remove: ttnClient.Gateways.Collaborators.remove.bind(ttnClient.Gateways.Collaborators),
    },
    apiKeys: {
      get: ttnClient.Gateways.ApiKeys.getById.bind(ttnClient.Gateways.ApiKeys),
      list: ttnClient.Gateways.ApiKeys.getAll.bind(ttnClient.Gateways.ApiKeys),
      update: ttnClient.Gateways.ApiKeys.updateById.bind(ttnClient.Gateways.ApiKeys),
      delete: ttnClient.Gateways.ApiKeys.deleteById.bind(ttnClient.Gateways.ApiKeys),
      create: ttnClient.Gateways.ApiKeys.create.bind(ttnClient.Gateways.ApiKeys),
    },
  },
  rights: {
    applications: ttnClient.Applications.getRightsById.bind(ttnClient.Applications),
    gateways: ttnClient.Gateways.getRightsById.bind(ttnClient.Gateways),
    organizations: ttnClient.Organizations.getRightsById.bind(ttnClient.Organizations),
  },
  configuration: {
    listNsFrequencyPlans: ttnClient.Configuration.listNsFrequencyPlans.bind(
      ttnClient.Configuration,
    ),
    listGsFrequencyPlans: ttnClient.Configuration.listGsFrequencyPlans.bind(
      ttnClient.Configuration,
    ),
  },
  js: {
    joinEUIPrefixes: {
      list: ttnClient.Js.listJoinEUIPrefixes.bind(ttnClient.Js),
    },
  },
  ns: {
    generateDevAddress: ttnClient.Ns.generateDevAddress.bind(ttnClient.Ns),
  },
  organizations: {
    list: ttnClient.Organizations.getAll.bind(ttnClient.Organizations),
    search: ttnClient.Organizations.search.bind(ttnClient.Organizations),
    create: ttnClient.Organizations.create.bind(ttnClient.Organizations),
  },
  organization: {
    get: ttnClient.Organizations.getById.bind(ttnClient.Organizations),
    eventsSubscribe: ttnClient.Organizations.openStream.bind(ttnClient.Organizations),
    delete: ttnClient.Organizations.deleteById.bind(ttnClient.Organizations),
    update: ttnClient.Organizations.updateById.bind(ttnClient.Organizations),
    apiKeys: {
      get: ttnClient.Organizations.ApiKeys.getById.bind(ttnClient.Organizations.ApiKeys),
      list: ttnClient.Organizations.ApiKeys.getAll.bind(ttnClient.Organizations.ApiKeys),
      update: ttnClient.Organizations.ApiKeys.updateById.bind(ttnClient.Organizations.ApiKeys),
      delete: ttnClient.Organizations.ApiKeys.deleteById.bind(ttnClient.Organizations.ApiKeys),
      create: ttnClient.Organizations.ApiKeys.create.bind(ttnClient.Organizations.ApiKeys),
    },
    collaborators: {
      getOrganization: ttnClient.Organizations.Collaborators.getByOrganizationId.bind(
        ttnClient.Organizations.Collaborators,
      ),
      getUser: ttnClient.Organizations.Collaborators.getByUserId.bind(
        ttnClient.Organizations.Collaborators,
      ),
      list: ttnClient.Organizations.Collaborators.getAll.bind(
        ttnClient.Organizations.Collaborators,
      ),
      add: ttnClient.Organizations.Collaborators.add.bind(ttnClient.Organizations.Collaborators),
      update: ttnClient.Organizations.Collaborators.update.bind(
        ttnClient.Organizations.Collaborators,
      ),
      remove: ttnClient.Organizations.Collaborators.remove.bind(
        ttnClient.Organizations.Collaborators,
      ),
    },
  },
}
