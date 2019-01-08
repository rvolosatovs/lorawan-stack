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

export const GET_GTWS_LIST = 'GET_GATEWAYS_LIST'
export const SEARCH_GTWS_LIST = 'SEARCH_GATEWAYS_LIST'
export const GET_GTWS_LIST_SUCCESS = 'GET_GATEWAYS_LIST_SUCCESS'
export const GET_GTWS_LIST_FAILURE = 'GET_GATEWAYS_LIST_FAILURE'
export const CHANGE_GTWS_PAGE = 'CHANGE_GATEWAYS_PAGE'
export const CHANGE_GTWS_ORDER = 'CHANGE_GATEWAYS_ORDER'
export const CHANGE_GTWS_SEARCH = 'CHANGE_GATEWAYS_SEARCH'

export const getGatewaysList = filters => (
  { type: GET_GTWS_LIST, filters }
)

export const searchGatewaysList = filters => (
  { type: SEARCH_GTWS_LIST, filters }
)

export const getGatewaysSuccess = (gateways, totalCount) => (
  { type: GET_GTWS_LIST_SUCCESS, gateways, totalCount }
)

export const getGatewaysFailure = error => (
  { type: GET_GTWS_LIST_FAILURE, error }
)

export const changeGatewaysPage = page => (
  { type: CHANGE_GTWS_PAGE, page }
)

export const changeGatewaysOrder = (order, orderBy) => (
  { type: CHANGE_GTWS_ORDER, order, orderBy }
)

export const changeGatewaysSearch = query => (
  { type: CHANGE_GTWS_SEARCH, query }
)
