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

import ApiKeys from './api-keys'

class Users {
  constructor(registry) {
    this._api = registry

    this.ApiKeys = new ApiKeys(registry.UserAccess, {
      parentRoutes: {
        get: 'user_ids.user_id',
        list: 'user_ids.user_id',
        create: 'user_ids.user_id',
        update: 'user_ids.user_id',
      },
    })
  }

  _addState(fieldMask, user) {
    // Ensure to set STATE_REQUESTED if needed, which gets stripped as null
    // value from the backend response.
    if (fieldMask && fieldMask.field_mask.paths.includes('state') && !('state' in user)) {
      user.state = 'STATE_REQUESTED'
    }

    return user
  }

  async getAll(params, selector) {
    const fieldMask = Marshaler.selectorToFieldMask(selector)
    const response = await this._api.UserRegistry.List(undefined, {
      ...params,
      ...fieldMask,
    })

    const users = Marshaler.payloadListResponse('users', response)
    users.users.map(user => this._addState(fieldMask, user))

    return users
  }

  async search(params, selector) {
    const fieldMask = Marshaler.selectorToFieldMask(selector)
    const response = await this._api.EntityRegistrySearch.SearchUsers(undefined, {
      ...params,
      ...fieldMask,
    })

    const users = Marshaler.payloadListResponse('users', response)
    users.users.map(user => this._addState(fieldMask, user))

    return users
  }

  async getById(id, selector) {
    const fieldMask = Marshaler.selectorToFieldMask(selector)
    const response = await this._api.UserRegistry.Get(
      {
        routeParams: { 'user_ids.user_id': id },
      },
      fieldMask,
    )

    const user = this._addState(fieldMask, Marshaler.payloadSingleResponse(response))

    return user
  }

  async deleteById(id) {
    const response = await this._api.UserRegistry.Delete({
      routeParams: { user_id: id },
    })

    return Marshaler.payloadSingleResponse(response)
  }

  async updateById(id, patch, mask = Marshaler.fieldMaskFromPatch(patch)) {
    const response = await this._api.UserRegistry.Update(
      {
        routeParams: {
          'user.ids.user_id': id,
        },
      },
      {
        user: patch,
        field_mask: Marshaler.fieldMask(mask),
      },
    )
    return Marshaler.unwrapUser(response)
  }

  async create(user, invitationToken) {
    const response = await this._api.UserRegistry.Create(undefined, {
      user,
      invitation_token: invitationToken,
    })
    return Marshaler.unwrapUser(response)
  }

  // Miscellaneous.

  async getRightsById(userId) {
    const result = await this._api.UserAccess.ListRights({
      routeParams: { user_id: userId },
    })

    return Marshaler.unwrapRights(result)
  }

  updatePasswordById(id, payload) {
    return this._api.UserRegistry.UpdatePassword(
      {
        routeParams: {
          'user_ids.user_id': id,
        },
      },
      {
        new: payload.new,
        old: payload.old,
        revoke_all_access: payload.revoke_all_access,
      },
    )
  }
}

export default Users
