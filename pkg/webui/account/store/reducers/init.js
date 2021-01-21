// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

import { INITIALIZE, INITIALIZE_SUCCESS, INITIALIZE_FAILURE } from '@account/store/actions/init'

const defaultState = {
  initialized: false,
  error: undefined,
}

const init = function(state = defaultState, action) {
  switch (action.type) {
    case INITIALIZE:
      return {
        ...state,
        error: undefined,
        initialized: false,
      }
    case INITIALIZE_SUCCESS:
      return {
        ...state,
        error: undefined,
        initialized: true,
      }
    case INITIALIZE_FAILURE:
      return {
        ...state,
        error: action.error,
        initialized: false,
      }
    default:
      return state
  }
}

export default init
