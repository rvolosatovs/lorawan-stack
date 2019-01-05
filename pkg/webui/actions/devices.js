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

export const GET_DEVICES_LIST = 'GET_DEVICES_LIST'
export const SEARCH_DEVICES_LIST = 'SEARCH_DEVICES_LIST'
export const GET_DEVICES_LIST_SUCCESS = 'GET_DEVICES_LIST_SUCCESS'
export const GET_DEVICES_LIST_FAILURE = 'GET_DEVICES_LIST_FAILURE'

export const getDevicesList = (appId, filters) => (
  { type: GET_DEVICES_LIST, appId, filters }
)

export const searchDevicesList = (appId, filters) => (
  { type: SEARCH_DEVICES_LIST, appId, filters }
)

export const getDevicesListSuccess = (devices, totalCount) => (
  { type: GET_DEVICES_LIST_SUCCESS, devices, totalCount }
)

export const getDevicesListFailure = error => (
  { type: GET_DEVICES_LIST_FAILURE, error }
)
