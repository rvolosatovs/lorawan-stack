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

import * as applications from '../../actions/applications'
import * as link from '../../actions/link'
import * as webhooks from '../../actions/webhooks'
import * as webhook from '../../actions/webhook'
import * as webhookFormats from '../../actions/webhook-formats'
import * as pubsubs from '../../actions/pubsubs'
import * as pubsub from '../../actions/pubsub'
import * as pubsubFormats from '../../actions/pubsub-formats'

import api from '../../../api'
import { isNotFoundError } from '../../../../lib/errors/utils'
import createRequestLogic from './lib'
import createEventsConnectLogics from './events'

const getApplicationLogic = createRequestLogic({
  type: applications.GET_APP,
  async process({ action }, dispatch) {
    const {
      payload: { id },
      meta: { selector },
    } = action
    const app = await api.application.get(id, selector)
    dispatch(applications.startApplicationEventsStream(id))
    return app
  },
})

const updateApplicationLogic = createRequestLogic({
  type: applications.UPDATE_APP,
  async process({ action }) {
    const { id, patch } = action.payload

    const result = await api.application.update(id, patch)

    return { ...patch, ...result }
  },
})

const deleteApplicationLogic = createRequestLogic({
  type: applications.DELETE_APP,
  async process({ action }) {
    const { id } = action.payload

    await api.application.delete(id)

    return { id }
  },
})

const getApplicationsLogic = createRequestLogic({
  type: applications.GET_APPS_LIST,
  latest: true,
  async process({ action }) {
    const {
      params: { page, limit, query },
    } = action.payload
    const { selectors } = action.meta

    const data = query
      ? await api.applications.search(
          {
            page,
            limit,
            id_contains: query,
            name_contains: query,
          },
          selectors,
        )
      : await api.applications.list({ page, limit }, selectors)

    return { entities: data.applications, totalCount: data.totalCount }
  },
})

const getApplicationsRightsLogic = createRequestLogic({
  type: applications.GET_APPS_RIGHTS_LIST,
  async process({ action }) {
    const { id } = action.payload
    const result = await api.rights.applications(id)
    return result.rights.sort()
  },
})

const getApplicationApiKeysLogic = createRequestLogic({
  type: applications.GET_APP_API_KEYS_LIST,
  async process({ getState, action }) {
    const {
      id: appId,
      params: { page, limit },
    } = action.payload
    const res = await api.application.apiKeys.list(appId, { limit, page })
    return { ...res, id: appId }
  },
})

const getApplicationApiKeyLogic = createRequestLogic({
  type: applications.GET_APP_API_KEY,
  async process({ action }) {
    const { id: appId, keyId } = action.payload
    return api.application.apiKeys.get(appId, keyId)
  },
})

const getApplicationCollaboratorLogic = createRequestLogic({
  type: applications.GET_APP_COLLABORATOR,
  async process({ action }) {
    const { id: appId, collaboratorId, isUser } = action.payload

    const collaborator = isUser
      ? await api.application.collaborators.getUser(appId, collaboratorId)
      : await api.application.collaborators.getOrganization(appId, collaboratorId)

    const { ids, ...rest } = collaborator

    return {
      id: collaboratorId,
      isUser,
      ...rest,
    }
  },
})

const getApplicationCollaboratorsLogic = createRequestLogic({
  type: applications.GET_APP_COLLABORATORS_LIST,
  async process({ action }) {
    const { id: appId, params } = action.payload
    const res = await api.application.collaborators.list(appId, params)
    const collaborators = res.collaborators.map(function(collaborator) {
      const { ids, ...rest } = collaborator
      const isUser = !!ids.user_ids
      const collaboratorId = isUser ? ids.user_ids.user_id : ids.organization_ids.organization_id

      return {
        id: collaboratorId,
        isUser,
        ...rest,
      }
    })
    return { id: appId, collaborators, totalCount: res.totalCount }
  },
})

const getApplicationLinkLogic = createRequestLogic({
  type: link.GET_APP_LINK,
  async process({ action }, dispatch, done) {
    const {
      payload: { id },
      meta: { selector = [] } = {},
    } = action

    let linkResult
    let statsResult
    try {
      linkResult = await api.application.link.get(id, selector)
      statsResult = await api.application.link.stats(id)

      return { link: linkResult, stats: statsResult, linked: true }
    } catch (error) {
      // Consider errors that are not 404, since not found means that the
      // application is not linked.
      if (isNotFoundError(error)) {
        return { link: linkResult, stats: statsResult, linked: false }
      }

      throw error
    }
  },
})

const getWebhookLogic = createRequestLogic({
  type: webhook.GET_WEBHOOK,
  async process({ action }) {
    const {
      payload: { appId, webhookId },
      meta: { selector },
    } = action
    return api.application.webhooks.get(appId, webhookId, selector)
  },
})

const getWebhooksLogic = createRequestLogic({
  type: webhooks.GET_WEBHOOKS_LIST,
  async process({ action }) {
    const { appId } = action.payload
    const res = await api.application.webhooks.list(appId)
    return { webhooks: res.webhooks, totalCount: res.totalCount }
  },
})

const getWebhookFormatsLogic = createRequestLogic({
  type: webhookFormats.GET_WEBHOOK_FORMATS,
  async process() {
    const { formats } = await api.application.webhooks.getFormats()
    return formats
  },
})

const getPubsubLogic = createRequestLogic({
  type: pubsub.GET_PUBSUB,
  async process({ action }) {
    const {
      payload: { appId, pubsubId },
      meta: { selector },
    } = action
    return api.application.pubsubs.get(appId, pubsubId, selector)
  },
})

const getPubsubsLogic = createRequestLogic({
  type: pubsubs.GET_PUBSUBS_LIST,
  async process({ action }) {
    const { appId } = action.payload
    const res = await api.application.pubsubs.list(appId)
    return { pubsubs: res.pubsubs, totalCount: res.totalCount }
  },
})

const getPubsubFormatsLogic = createRequestLogic({
  type: pubsubFormats.GET_PUBSUB_FORMATS,
  async process() {
    const { formats } = await api.application.pubsubs.getFormats()
    return formats
  },
})

export default [
  getApplicationLogic,
  updateApplicationLogic,
  deleteApplicationLogic,
  getApplicationsLogic,
  getApplicationsRightsLogic,
  getApplicationApiKeysLogic,
  getApplicationApiKeyLogic,
  getApplicationCollaboratorLogic,
  getApplicationCollaboratorsLogic,
  getWebhooksLogic,
  getWebhookLogic,
  getWebhookFormatsLogic,
  getPubsubLogic,
  getPubsubsLogic,
  getPubsubFormatsLogic,
  getApplicationLinkLogic,
  ...createEventsConnectLogics(
    applications.SHARED_NAME_SINGLE,
    'applications',
    api.application.eventsSubscribe,
  ),
]
