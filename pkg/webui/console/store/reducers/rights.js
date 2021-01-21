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

import { createGetRightsListActionType } from '@console/store/actions/rights'

const defaultState = {
  rights: [],
}

const createNamedRightsReducer = (reducerName = '') => {
  const GET_LIST_BASE = createGetRightsListActionType(reducerName)
  const [{ success: GET_LIST_SUCCESS }] = createRequestActions(GET_LIST_BASE)

  return (state = defaultState, { type, payload }) => {
    switch (type) {
      case GET_LIST_SUCCESS:
        return {
          ...state,
          rights: payload,
        }
      default:
        return state
    }
  }
}

export default createNamedRightsReducer
