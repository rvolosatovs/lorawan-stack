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

import Marshaler from '../util/marshaler'
import combineStreams from '../util/combine-streams'
import { STACK_COMPONENTS_MAP } from '../util/constants'

import Devices from './devices'
import ApiKeys from './api-keys'
import Link from './link'
import Collaborators from './collaborators'
import Webhooks from './webhooks'
import PubSubs from './pubsubs'
import Packages from './application-packages'

const { is: IS, as: AS, ns: NS, js: JS, dtc: DTC } = STACK_COMPONENTS_MAP

/**
 * Applications Class provides an abstraction on all applications and manages
 * data handling from different sources. It exposes an API to easily work with
 * application data.
 *
 * @param {object} api - The connector to be used by the service.
 * @param {object} config - The configuration for the service.
 * @param {string} config.defaultUserId - The users identifier to be used in
 * user related requests.
 * @param {boolean} config.proxy - The flag to identify if the results
 * should be proxied with the wrapper objects.
 */
class Applications {
  constructor(api, { defaultUserId, stackConfig }) {
    this._defaultUserId = defaultUserId
    this._api = api
    this._stackConfig = stackConfig

    this.ApiKeys = new ApiKeys(api.ApplicationAccess, {
      parentRoutes: {
        get: 'application_ids.application_id',
        list: 'application_ids.application_id',
        create: 'application_ids.application_id',
        update: 'application_ids.application_id',
      },
    })
    this.Link = new Link(api.As)
    this.Devices = new Devices(api, { stackConfig })
    this.Collaborators = new Collaborators(api.ApplicationAccess, {
      parentRoutes: {
        get: 'application_ids.application_id',
        list: 'application_ids.application_id',
        set: 'application_ids.application_id',
      },
    })
    this.Webhooks = new Webhooks(api.ApplicationWebhookRegistry)
    this.PubSubs = new PubSubs(api.ApplicationPubSubRegistry)
    this.Packages = new Packages(api.ApplicationPackageRegistry)
  }

  // Retrieval.

  async getAll(params, selector) {
    const response = await this._api.ApplicationRegistry.List(undefined, {
      ...params,
      ...Marshaler.selectorToFieldMask(selector),
    })

    return Marshaler.unwrapApplications(response)
  }

  async getById(id, selector) {
    const fieldMask = Marshaler.selectorToFieldMask(selector)
    const response = await this._api.ApplicationRegistry.Get(
      {
        routeParams: { 'application_ids.application_id': id },
      },
      fieldMask,
    )

    return Marshaler.unwrapApplication(response)
  }

  async getByOrganization(organizationId) {
    const response = this._api.ApplicationRegistry.List({
      routeParams: { 'collaborator.organization_ids.organization_id': organizationId },
    })

    return Marshaler.unwrapApplications(response)
  }

  async getByCollaborator(userId) {
    const response = this._api.ApplicationRegistry.List({
      routeParams: { 'collaborator.user_ids.user_id': userId },
    })

    return Marshaler.unwrapApplications(response)
  }

  async search(params, selector) {
    const response = await this._api.EntityRegistrySearch.SearchApplications(undefined, {
      ...params,
      ...Marshaler.selectorToFieldMask(selector),
    })

    return Marshaler.unwrapApplications(response)
  }

  // Update.

  async updateById(
    id,
    patch,
    mask = Marshaler.fieldMaskFromPatch(
      patch,
      this._api.ApplicationRegistry.UpdateAllowedFieldMaskPaths,
    ),
  ) {
    const response = await this._api.ApplicationRegistry.Update(
      {
        routeParams: {
          'application.ids.application_id': id,
        },
      },
      {
        application: patch,
        field_mask: Marshaler.fieldMask(mask),
      },
    )
    return Marshaler.unwrapApplication(response)
  }

  // Creation.

  async create(ownerId = this._defaultUserId, application, isUserOwner = true) {
    const routeParams = isUserOwner
      ? { 'collaborator.user_ids.user_id': ownerId }
      : { 'collaborator.organization_ids.organization_id': ownerId }
    const response = await this._api.ApplicationRegistry.Create(
      {
        routeParams,
      },
      { application },
    )
    return Marshaler.unwrapApplication(response)
  }

  // Deletion.

  async deleteById(applicationId) {
    const response = await this._api.ApplicationRegistry.Delete({
      routeParams: { application_id: applicationId },
    })

    return Marshaler.payloadSingleResponse(response)
  }

  // Miscellaneous.

  async getRightsById(applicationId) {
    const result = await this._api.ApplicationAccess.ListRights({
      routeParams: { application_id: applicationId },
    })

    return Marshaler.unwrapRights(result)
  }

  async getMqttConnectionInfo(applicationId) {
    const response = await this._api.AppAs.GetMQTTConnectionInfo({
      routeParams: { application_id: applicationId },
    })

    return Marshaler.payloadSingleResponse(response)
  }

  // Events Stream

  async openStream(identifiers, tail, after) {
    const payload = {
      identifiers: identifiers.map(id => ({
        application_ids: { application_id: id },
      })),
      tail,
      after,
    }

    // Event streams can come from multiple stack components. It is necessary to
    // check for stack components on different hosts and open distinct stream
    // connections for any distinct host if need be.
    const distinctComponents = this._stackConfig.getComponentsWithDistinctBaseUrls([
      IS,
      JS,
      NS,
      AS,
      DTC,
    ])

    const streams = distinctComponents.map(component =>
      this._api.Events.Stream({ component }, payload),
    )

    // Combine all stream sources to one subscription generator.
    return combineStreams(streams)
  }
}

export default Applications
