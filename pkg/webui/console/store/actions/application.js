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

import { createGetApiKeysListActionType } from '../actions/api-keys'

import {
  getCollaboratorsList,
  createGetCollaboratorsListActionType,
  getCollaboratorsListFailure,
  createGetCollaboratorsListFailureActionType,
  getCollaboratorsListSuccess,
  createGetCollaboratorsListSuccessActionType,
  createGetCollaboratorActionType,
  getCollaborator,
} from '../actions/collaborators'

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
} from '../actions/events'

import {
  getApiKey,
  createGetApiKeyActionType,
  getApiKeySuccess,
  createGetApiKeySuccessActionType,
  getApiKeyFailure,
  createGetApiKeyFailureActionType,
} from './api-key'

import { createRequestActions } from './lib'

export const SHARED_NAME = 'APPLICATION'

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

export const GET_APP_API_KEYS_LIST_BASE = createGetApiKeysListActionType(SHARED_NAME)
export const [{
  request: GET_APP_API_KEYS_LIST,
  success: GET_APP_API_KEYS_LIST_SUCCESS,
  failure: GET_APP_API_KEYS_LIST_FAILURE,
}, {
  request: getApplicationApiKeysList,
  success: getApplicationApiKeysListSuccess,
  failure: getApplicationApiKeysListFailure,
}] = createRequestActions(GET_APP_API_KEYS_LIST_BASE, (appId, params) => ({ appId, params }))

export const GET_APP_API_KEY_BASE = createGetApiKeyActionType(SHARED_NAME)
export const [{
  request: GET_APP_API_KEY,
  success: GET_APP_API_KEY_SUCCESS,
  failure: GET_APP_API_KEY_FAILURE,
}, {
  request: getApplicationApiKey,
  success: getApplicationApiKeySuccess,
  failure: getApplicationApiKeyFailure,
}] = createRequestActions(GET_APP_API_KEY_BASE)

export const GET_APP_API_KEY = createGetApiKeyActionType(SHARED_NAME)
export const GET_APP_API_KEY_SUCCESS = createGetApiKeySuccessActionType(SHARED_NAME)
export const GET_APP_API_KEY_FAILURE = createGetApiKeyFailureActionType(SHARED_NAME)
export const GET_APP_COLLABORATOR_PAGE_DATA = createGetCollaboratorActionType(SHARED_NAME)
export const GET_APP_COLLABORATORS_LIST = createGetCollaboratorsListActionType(SHARED_NAME)
export const GET_APP_COLLABORATORS_LIST_SUCCESS = createGetCollaboratorsListSuccessActionType(SHARED_NAME)
export const GET_APP_COLLABORATORS_LIST_FAILURE = createGetCollaboratorsListFailureActionType(SHARED_NAME)
export const START_APP_EVENT_STREAM = createStartEventsStreamActionType(SHARED_NAME)
export const START_APP_EVENT_STREAM_SUCCESS = createStartEventsStreamSuccessActionType(SHARED_NAME)
export const START_APP_EVENT_STREAM_FAILURE = createStartEventsStreamFailureActionType(SHARED_NAME)
export const STOP_APP_EVENT_STREAM = createStopEventsStreamActionType(SHARED_NAME)
export const CLEAR_APP_EVENTS = createClearEventsActionType(SHARED_NAME)

export const getApplicationApiKey = getApiKey(SHARED_NAME)

export const getApplicationApiKeySuccess = getApiKeySuccess(SHARED_NAME)

export const getApplicationApiKeyFailure = getApiKeyFailure(SHARED_NAME)

export const getApplicationCollaboratorsList = getCollaboratorsList(SHARED_NAME)

export const getApplicationCollaboratorsListSuccess = getCollaboratorsListSuccess(SHARED_NAME)

export const getApplicationCollaboratorsListFailure = getCollaboratorsListFailure(SHARED_NAME)

export const getApplicationCollaboratorPageData = getCollaborator(SHARED_NAME)

export const startApplicationEventsStream = startEventsStream(SHARED_NAME)

export const startApplicationEventsStreamSuccess = startEventsStreamSuccess(SHARED_NAME)

export const startApplicationEventsStreamFailure = startEventsStreamFailure(SHARED_NAME)

export const stopApplicationEventsStream = stopEventsStream(SHARED_NAME)

export const clearApplicationEventsStream = clearEvents(SHARED_NAME)
