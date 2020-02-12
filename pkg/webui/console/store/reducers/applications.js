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

import { getApplicationId } from '../../../lib/selectors/id'
import {
  GET_APP,
  GET_APP_SUCCESS,
  GET_APP_DEV_COUNT_SUCCESS,
  GET_APPS_LIST_SUCCESS,
  UPDATE_APP_SUCCESS,
  DELETE_APP_SUCCESS,
} from '../actions/applications'

const application = function(state = {}, application) {
  return {
    ...state,
    ...application,
  }
}

const defaultState = {
  entities: {},
  selectedApplication: null,
  applicationDeviceCount: undefined,
}

const applications = function(state = defaultState, { type, payload }) {
  switch (type) {
    case GET_APP:
      return {
        ...state,
        selectedApplication: payload.id,
      }
    case GET_APPS_LIST_SUCCESS:
      const entities = payload.entities.reduce(
        function(acc, app) {
          const id = getApplicationId(app)

          acc[id] = application(acc[id], app)
          return acc
        },
        { ...state.entities },
      )

      return {
        ...state,
        entities,
      }
    case GET_APP_DEV_COUNT_SUCCESS:
      return {
        ...state,
        applicationDeviceCount: payload.applicationDeviceCount,
      }
    case GET_APP_SUCCESS:
    case UPDATE_APP_SUCCESS:
      const id = getApplicationId(payload)

      return {
        ...state,
        entities: {
          ...state.entities,
          [id]: application(state.entities[id], payload),
        },
      }
    case DELETE_APP_SUCCESS:
      const { [payload.id]: deleted, ...rest } = state.entities

      return {
        selectedApplication: null,
        entities: rest,
      }
    default:
      return state
  }
}

export default applications
