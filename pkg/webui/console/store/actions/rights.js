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

export const createGetRightsListActionType = name => (
  `GET_${name}_RIGHTS_LIST`
)

export const createGetRightsListSuccessActionType = name => (
  `GET_${name}_RIGHTS_LIST_SUCCESS`
)

export const createGetRightsListFailureActionType = name => (
  `GET_${name}_RIGHTS_LIST_FAILURE`
)

export const getRightsList = name => id => (
  { type: createGetRightsListActionType(name), id }
)

export const getRightsListSuccess = name => rights => (
  { type: createGetRightsListSuccessActionType(name), rights }
)

export const getRightsListFailure = name => error => (
  { type: createGetRightsListFailureActionType(name), error }
)
