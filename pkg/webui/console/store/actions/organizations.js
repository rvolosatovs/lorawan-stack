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

import { createRequestActions } from './lib'
import { createPaginationRequestActions, createPaginationBaseActionType } from './pagination'

export const SHARED_NAME = 'ORGANIZATION'

export const GET_ORGS_LIST_BASE = createPaginationBaseActionType(SHARED_NAME)
export const [
  { request: GET_ORGS_LIST, success: GET_ORGS_LIST_SUCCESS, failure: GET_ORGS_LIST_FAILURE },
  {
    request: getOrganizationsList,
    success: getOrganizationsListSuccess,
    failure: getORganizationsListFailure,
  },
] = createPaginationRequestActions(SHARED_NAME)

export const CREATE_ORG_BASE = 'CREATE_ORGANIZATION'
export const [
  { request: CREATE_ORG, success: CREATE_ORG_SUCCESS, failure: CREATE_ORG_FAILURE },
  {
    request: createOrganization,
    success: createOrganizationSuccess,
    failure: createOrganizationFailure,
  },
] = createRequestActions(CREATE_ORG_BASE)
