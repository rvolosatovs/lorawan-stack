// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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

import { createLogic } from 'redux-logic'

import * as user from '../actions/user'
import * as _console from '../actions/console'
import api from '../api'
import * as accessToken from '../lib/access-token'

const consoleLogic = createLogic({
  type: 'INITIALIZE',
  async process ({ getState, action }, dispatch, done) {
    try {
      try {
        const result = await api.v3.is.users.me()
        console.log(result)
        dispatch(user.getUserMeSuccess(result.data))
      } catch (error) {
        dispatch(user.getUserMeFailure())
        accessToken.clear()
      }
      dispatch(_console.initializeSuccess())
      console.log('Console has initialized!')
    } catch (error) {
      console.log(error)
      dispatch(_console.initializeFailure())
    }
    done()
  },
})

export default consoleLogic
