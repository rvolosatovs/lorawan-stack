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

import * as user from '../../actions/user'
import * as init from '../../actions/init'
import api from '../../../api'
import * as accessToken from '../../../lib/access-token'
import { getBackendErrorName } from '../../../../lib/errors/utils'
import createRequestLogic from './lib'

const consoleAppLogic = createRequestLogic({
  type: init.INITIALIZE,
  async process(_, dispatch) {
    dispatch(user.getUserRights())

    let info

    try {
      // there is no way to retrieve the current user directly
      // within the console app, so first get the authentication info
      // and only afterwards fetch the user
      info = await api.users.authInfo()
      const rights = info.data.oauth_access_token.rights
      dispatch(user.getUserRightsSuccess(rights))
    } catch (error) {
      dispatch(user.getUserRightsFailure())
      accessToken.clear()
    }

    if (info) {
      try {
        dispatch(user.getUserMe())
        const userId = info.data.oauth_access_token.user_ids.user_id
        const userResult = await api.users.get(userId)
        userResult.isAdmin = info.data.is_admin || false
        dispatch(user.getUserMeSuccess(userResult))
      } catch (error) {
        if (getBackendErrorName(error) === 'user_requested') {
          // Log unapproved users out automatically, to avoid confusion
          dispatch(user.logout())
        }
        dispatch(user.getUserMeFailure(error))
        accessToken.clear()
      }
    }

    // eslint-disable-next-line no-console
    console.log('Console initialization successful!')

    return
  },
})

export default [consoleAppLogic]
