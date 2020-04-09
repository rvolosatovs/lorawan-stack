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

import getByPath from '../get-by-path'

export const getApplicationId = function(application = {}) {
  return (
    getByPath(application, 'application_id') ||
    getByPath(application, 'application_ids.application_id') ||
    getByPath(application, 'ids.application_id')
  )
}

export const getDeviceId = function(device = {}) {
  return (
    getByPath(device, 'device_id') ||
    getByPath(device, 'ids.device_id') ||
    getByPath(device, 'device_ids.device_id')
  )
}

export const combineDeviceIds = (appId, devId) => `${appId}/${devId}`
export const extractDeviceIdFromCombinedId = function(combinedId) {
  if (typeof combinedId === 'string') {
    const parts = combinedId.split('/')
    if (parts.length === 2) {
      return parts[1]
    }
  }
  return combinedId
}
export const getCombinedDeviceId = function(device = {}) {
  const appId =
    getByPath(device, 'ids.application_ids.application_id') ||
    getByPath(device, 'application_ids.application_id')
  const devId = getDeviceId(device)
  return combineDeviceIds(appId, devId)
}

export const getCollaboratorId = function(collaborator = {}) {
  return (
    getByPath(collaborator, 'ids.organization_ids.organization_id') ||
    getByPath(collaborator, 'ids.user_ids.user_id')
  )
}

export const getGatewayId = function(gateway = {}) {
  return (
    getByPath(gateway, 'gateway_id') ||
    getByPath(gateway, 'gateway_ids.gateway_id') ||
    getByPath(gateway, 'ids.gateway_id')
  )
}

export const getApiKeyId = function(key = {}) {
  return key.id
}

export const getOrganizationId = function(organization = {}) {
  return (
    getByPath(organization, 'ids.organization_id') ||
    getByPath(organization, 'organization_ids.organization_id')
  )
}

const idSelectors = [
  getApplicationId,
  getCollaboratorId,
  getApiKeyId,
  getGatewayId,
  getDeviceId,
  getOrganizationId,
]

export const getEntityId = function(entity) {
  let id
  let selectorIndex = 0
  while (!id && selectorIndex < idSelectors.length) {
    const selector = idSelectors[selectorIndex++]
    id = selector(entity)
  }

  return id
}

export const getWebhookId = function(webhook = {}) {
  return getByPath(webhook, 'ids.webhook_id')
}

export const getWebhookTemplateId = function(webhookTemplate = {}) {
  return (
    getByPath(webhookTemplate, 'ids.template_id') ||
    getByPath(webhookTemplate, 'template_ids.template_id')
  )
}

export const getPubsubId = function(pubsub = {}) {
  return getByPath(pubsub, 'ids.pub_sub_id')
}

export const getUserId = function(user = {}) {
  return getByPath(user, 'ids.user_id')
}
