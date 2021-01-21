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

import createRequestActions from '@ttn-lw/lib/store/actions/create-request-actions'

export const GET_PUBSUB_BASE = 'GET_PUBSUB'
export const [
  { request: GET_PUBSUB, success: GET_PUBSUB_SUCCESS, failure: GET_PUBSUB_FAILURE },
  { request: getPubsub, success: getPubsubSuccess, failure: getPubsubFailure },
] = createRequestActions(
  GET_PUBSUB_BASE,
  (appId, pubsubId) => ({ appId, pubsubId }),
  (appId, pubsubId, selector) => ({ selector }),
)

export const GET_PUBSUBS_LIST_BASE = 'GET_PUBSUBS_LIST'
export const [
  {
    request: GET_PUBSUBS_LIST,
    success: GET_PUBSUBS_LIST_SUCCESS,
    failure: GET_PUBSUBS_LIST_FAILURE,
  },
  { request: getPubsubsList, success: getPubsubsListSuccess, failure: getPubsubsListFailure },
] = createRequestActions(GET_PUBSUBS_LIST_BASE, appId => ({ appId }))

export const UPDATE_PUBSUB_BASE = 'UPDATE_PUBSUB'
export const [
  { request: UPDATE_PUBSUB, success: UPDATE_PUBSUB_SUCCESS, failure: UPDATE_PUBSUB_FAILURE },
  { request: updatePubsub, success: updatePubsubSuccess, failure: updatePubsubFailure },
] = createRequestActions(UPDATE_PUBSUB_BASE, (appId, pubsubId, patch) => ({
  appId,
  pubsubId,
  patch,
}))
