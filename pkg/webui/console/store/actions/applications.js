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

import createGetRightsListRequestActions, { createGetRightsListActionType } from './rights'
import { createPaginationRequestActions, createPaginationBaseActionType } from './pagination'
import createApiKeysRequestActions, { createGetApiKeysListActionType } from './api-keys'
import createApiKeyRequestActions, { createGetApiKeyActionType } from './api-key'
import { createGetCollaboratorsListRequestActions, createGetCollaboratorsListActionType } from './collaborators'
import { createRequestActions } from './lib'

import {
  startEventsStream,
  createStartEventsStreamActionType,
  startEventsStreamSuccess,
  createStartEventsStreamSuccessActionType,
  startEventsStreamFailure,
  createStartEventsStreamFailureActionType,
  stopEventsStream,
  createStopEventsStreamActionType,
  clearEvents,
  createClearEventsActionType,
} from './events'

export const SHARED_NAME = 'APPLICATIONS'
export const SHARED_NAME_SINGLE = 'APPLICATION'

export const GET_APP_BASE = 'GET_APPLICATION'
export const [{
  request: GET_APP,
  success: GET_APP_SUCCESS,
  failure: GET_APP_FAILURE,
}, {
  request: getApplication,
  success: getApplicationSuccess,
  failure: getApplicationFailure,
}] = createRequestActions(GET_APP_BASE,
  id => ({ id }),
  (id, selector) => ({ selector })
)

export const GET_APPS_LIST_BASE = createPaginationBaseActionType(SHARED_NAME)
export const [{
  request: GET_APPS_LIST,
  success: GET_APPS_LIST_SUCCESS,
  failure: GET_APPS_LIST_FAILURE,
}, {
  request: getApplicationsList,
  success: getApplicationsSuccess,
  failure: getApplicationsFailure,
}] = createPaginationRequestActions(SHARED_NAME)

export const GET_APPS_RIGHTS_LIST_BASE = createGetRightsListActionType(SHARED_NAME)
export const [{
  request: GET_APPS_RIGHTS_LIST,
  success: GET_APPS_RIGHTS_LIST_SUCCESS,
  failure: GET_APPS_RIGHTS_LIST_FAILURE,
}, {
  request: getApplicationsRightsList,
  success: getApplicationsRightsListSuccess,
  failure: getApplicationsRightsListFailure,
}] = createGetRightsListRequestActions(SHARED_NAME)

export const GET_APP_API_KEY_BASE = createGetApiKeyActionType(SHARED_NAME_SINGLE)
export const [{
  request: GET_APP_API_KEY,
  success: GET_APP_API_KEY_SUCCESS,
  failure: GET_APP_API_KEY_FAILURE,
}, {
  request: getApplicationApiKey,
  success: getApplicationApiKeySuccess,
  failure: getApplicationApiKeyFailure,
}] = createApiKeyRequestActions(SHARED_NAME_SINGLE)

export const GET_APP_API_KEYS_LIST_BASE = createGetApiKeysListActionType(SHARED_NAME_SINGLE)
export const [{
  request: GET_APP_API_KEYS_LIST,
  success: GET_APP_API_KEYS_LIST_SUCCESS,
  failure: GET_APP_API_KEYS_LIST_FAILURE,
}, {
  request: getApplicationApiKeysList,
  success: getApplicationApiKeysListSuccess,
  failure: getApplicationApiKeysListFailure,
}] = createApiKeysRequestActions(SHARED_NAME_SINGLE)

export const GET_APP_COLLABORATORS_LIST_BASE = createGetCollaboratorsListActionType(SHARED_NAME_SINGLE)
export const [{
  request: GET_APP_COLLABORATORS_LIST,
  success: GET_APP_COLLABORATORS_LIST_SUCCESS,
  failure: GET_APP_COLLABORATORS_LIST_FAILURE,
}, {
  request: getApplicationCollaboratorsList,
  success: getApplicationCollaboratorsListSuccess,
  failure: getApplicationCollaboratorsListFailure,
}] = createGetCollaboratorsListRequestActions(SHARED_NAME_SINGLE)

export const START_APP_EVENT_STREAM = createStartEventsStreamActionType(SHARED_NAME_SINGLE)
export const START_APP_EVENT_STREAM_SUCCESS = createStartEventsStreamSuccessActionType(SHARED_NAME_SINGLE)
export const START_APP_EVENT_STREAM_FAILURE = createStartEventsStreamFailureActionType(SHARED_NAME_SINGLE)
export const STOP_APP_EVENT_STREAM = createStopEventsStreamActionType(SHARED_NAME_SINGLE)
export const CLEAR_APP_EVENTS = createClearEventsActionType(SHARED_NAME_SINGLE)

export const startApplicationEventsStream = startEventsStream(SHARED_NAME_SINGLE)
export const startApplicationEventsStreamSuccess = startEventsStreamSuccess(SHARED_NAME_SINGLE)
export const startApplicationEventsStreamFailure = startEventsStreamFailure(SHARED_NAME_SINGLE)
export const stopApplicationEventsStream = stopEventsStream(SHARED_NAME_SINGLE)
export const clearApplicationEventsStream = clearEvents(SHARED_NAME_SINGLE)
