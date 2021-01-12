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

import { createLogic } from 'redux-logic'
import axios from 'axios'

import ONLINE_STATUS from '@ttn-lw/constants/online-status'

import { selectIsConfig } from '@ttn-lw/lib/selectors/env'
import { isNetworkError, isTimeoutError } from '@ttn-lw/lib/errors/utils'

import * as status from '@console/store/actions/status'

import { selectIsOfflineStatus } from '@console/store/selectors/status'

const isRoot = selectIsConfig().base_url

const initialInterval = 5000
let interval = initialInterval
const connectionCheck = (dispatch, done) => () => {
  dispatch(status.attemptReconnect())
  done()
}

let periodicCheck

const connectionManagementLogic = createLogic({
  type: status.SET_CONNECTION_STATUS,
  process: async ({ action }, dispatch, done) => {
    if (action.payload.onlineStatus === ONLINE_STATUS.CHECKING) {
      try {
        // Make a simple GET request to the auth_info endpoint.
        await axios.get(`${isRoot}/auth_info`, { timeout: 5000 })
        dispatch(status.setStatusOnline())
      } catch (error) {
        // If this one fails with a network error, we can be sufficiently
        // sure of having gone offline.
        if (isNetworkError(error) || isTimeoutError(error)) {
          dispatch(status.setStatusOnline(false))
        }
      }
    }

    if (action.payload.onlineStatus === ONLINE_STATUS.OFFLINE && navigator.onLine) {
      // If the app went offline, try to reconnect periodically.
      dispatch(status.attemptReconnect())
    }

    done()
  },
})

const connectionCheckLogic = createLogic({
  type: status.ATTEMPT_RECONNECT,
  // Additionally to periodic reconnects, freshly incoming request actions will
  // also trigger reconnection attempts, which is why this action is throttled
  // to 5 seconds.
  throttle: 5000,
  validate: ({ action, getState }, allow, reject) => {
    if (selectIsOfflineStatus(getState()) && navigator.onLine) {
      return allow(action)
    }
    if (Boolean(periodicCheck)) {
      clearTimeout(periodicCheck)
    }
    reject()
  },
  process: async (_, dispatch, done) => {
    try {
      await axios.get(`${isRoot}/auth_info`, { timeout: 4500 })
      dispatch(status.setStatusOnline())
      dispatch(status.attemptReconnectSuccess())
      interval = initialInterval
    } catch (error) {
      dispatch(status.attemptReconnectFailure())
    }

    done()
  },
})

const connectionCheckFailLogic = createLogic({
  type: status.ATTEMPT_RECONNECT_FAILURE,
  cancelType: status.ATTEMPT_RECONNECT_SUCCESS,
  warnTimeout: 65000,
  process: (_, dispatch, done) => {
    // Use increasing intervals, capped at 1min to prevent request spamming.
    interval = Math.min(interval * 1.5, 60000)
    periodicCheck = setTimeout(connectionCheck(dispatch, done), interval)
  },
})

export default [connectionManagementLogic, connectionCheckLogic, connectionCheckFailLogic]
