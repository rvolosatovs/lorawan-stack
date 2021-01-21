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

export const GET_NS_FREQUENCY_PLANS_BASE = 'GET_NS_FREQUENCY_PLANS'
export const [
  {
    request: GET_NS_FREQUENCY_PLANS,
    success: GET_NS_FREQUENCY_PLANS_SUCCESS,
    failure: GET_NS_FREQUENCY_PLANS_FAILURE,
  },
  {
    request: getNsFrequencyPlans,
    success: getNsFrequencyPlansSuccess,
    failure: getNsFrequencyPlansFailure,
  },
] = createRequestActions(GET_NS_FREQUENCY_PLANS_BASE)

export const GET_GS_FREQUENCY_PLANS_BASE = 'GET_GS_FREQUENCY_PLANS'
export const [
  {
    request: GET_GS_FREQUENCY_PLANS,
    success: GET_GS_FREQUENCY_PLANS_SUCCESS,
    failure: GET_GS_FREQUENCY_PLANS_FAILURE,
  },
  {
    request: getGsFrequencyPlans,
    success: getGsFrequencyPlansSuccess,
    failure: getGsFrequencyPlansFailure,
  },
] = createRequestActions(GET_GS_FREQUENCY_PLANS_BASE)
