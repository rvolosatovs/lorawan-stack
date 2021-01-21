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

import { createLogic } from 'redux-logic'

import api from '@console/api'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import { selectGsConfig } from '@ttn-lw/lib/selectors/env'
import { getGatewayId } from '@ttn-lw/lib/selectors/id'
import getHostFromUrl from '@ttn-lw/lib/host-from-url'
import createRequestLogic from '@ttn-lw/lib/store/logics/create-request-logic'

import * as gateways from '@console/store/actions/gateways'

import {
  selectGatewayById,
  selectGatewayStatisticsIsFetching,
} from '@console/store/selectors/gateways'

import createEventsConnectLogics from './events'

const getGatewayLogic = createRequestLogic({
  type: gateways.GET_GTW,
  async process({ action }, dispatch) {
    const { payload, meta } = action
    const { id = {} } = payload
    const selector = meta.selector || ''
    const gtw = await api.gateway.get(id, selector)
    dispatch(gateways.startGatewayEventsStream(id))
    return gtw
  },
})

const updateGatewayLogic = createRequestLogic({
  type: gateways.UPDATE_GTW,
  async process({ action }) {
    const {
      payload: { id, patch },
    } = action
    const result = await api.gateway.update(id, patch)

    return { ...patch, ...result }
  },
})

const deleteGatewayLogic = createRequestLogic({
  type: gateways.DELETE_GTW,
  async process({ action }) {
    const { id } = action.payload

    await api.gateway.delete(id)

    return { id }
  },
})

const getGatewaysLogic = createRequestLogic({
  type: gateways.GET_GTWS_LIST,
  latest: true,
  async process({ action }) {
    const {
      params: { page, limit, query, order },
    } = action.payload
    const { selectors, options } = action.meta

    const data = options.isSearch
      ? await api.gateways.search(
          {
            page,
            limit,
            id_contains: query,
            order,
          },
          selectors,
        )
      : await api.gateways.list({ page, limit, order }, selectors)

    let entities = data.gateways
    if (options.withStatus) {
      const gsConfig = selectGsConfig()
      const consoleGsAddress = getHostFromUrl(gsConfig.base_url)

      entities = await Promise.all(
        data.gateways.map(gateway => {
          const gatewayServerAddress = gateway.gateway_server_address

          if (!Boolean(gateway.gateway_server_address)) {
            return Promise.resolve({ ...gateway, status: 'unknown' })
          }

          if (gatewayServerAddress !== consoleGsAddress) {
            return Promise.resolve({ ...gateway, status: 'other-cluster' })
          }

          const id = getGatewayId(gateway)
          return api.gateway
            .stats(id)
            .then(() => ({ ...gateway, status: 'connected' }))
            .catch(err => {
              if (err && err.code === 5) {
                return { ...gateway, status: 'disconnected' }
              }

              return { ...gateway, status: 'unknown' }
            })
        }),
      )
    }

    return {
      entities,
      totalCount: data.totalCount,
    }
  },
})

const getGatewaysRightsLogic = createRequestLogic({
  type: gateways.GET_GTWS_RIGHTS_LIST,
  async process({ action }, dispatch, done) {
    const { id } = action.payload
    const result = await api.rights.gateways(id)
    return result.rights.sort()
  },
})

const startGatewayStatisticsLogic = createLogic({
  type: gateways.START_GTW_STATS,
  cancelType: [gateways.STOP_GTW_STATS, gateways.UPDATE_GTW_STATS_FAILURE],
  warnTimeout: 0,
  processOptions: {
    dispatchMultiple: true,
  },
  async process({ cancelled$, action, getState }, dispatch, done) {
    const { id } = action.payload
    const { timeout = 60000 } = action.meta

    const gsConfig = selectGsConfig()
    const gtw = selectGatewayById(getState(), id)

    if (!gsConfig.enabled) {
      dispatch(
        gateways.startGatewayStatisticsFailure({
          message: 'Unavailable',
        }),
      )
      done()
    }

    let gtwGsAddress
    let consoleGsAddress
    try {
      const gtwAddress = gtw.gateway_server_address

      if (!Boolean(gtwAddress)) {
        throw new Error()
      }

      gtwGsAddress = gtwAddress.split(':')[0]
      consoleGsAddress = new URL(gsConfig.base_url).hostname
    } catch (error) {
      dispatch(
        gateways.startGatewayStatisticsFailure({
          message: sharedMessages.statusUnknown,
        }),
      )
      done()
    }

    if (gtwGsAddress !== consoleGsAddress) {
      dispatch(
        gateways.startGatewayStatisticsFailure({
          message: sharedMessages.otherCluster,
        }),
      )
      done()
    }

    dispatch(gateways.startGatewayStatisticsSuccess())
    dispatch(gateways.updateGatewayStatistics(id))

    const interval = setInterval(() => {
      const statsRequestInProgress = selectGatewayStatisticsIsFetching(getState())
      if (!statsRequestInProgress) {
        dispatch(gateways.updateGatewayStatistics(id))
      }
    }, timeout)

    cancelled$.subscribe(() => clearInterval(interval))
  },
})

const updateGatewayStatisticsLogic = createRequestLogic({
  type: gateways.UPDATE_GTW_STATS,
  async process({ action }) {
    const { id } = action.payload

    const stats = await api.gateway.stats(id)

    return { stats }
  },
})

export default [
  getGatewayLogic,
  updateGatewayLogic,
  deleteGatewayLogic,
  getGatewaysLogic,
  getGatewaysRightsLogic,
  startGatewayStatisticsLogic,
  updateGatewayStatisticsLogic,
  ...createEventsConnectLogics(gateways.SHARED_NAME, 'gateways', api.gateway.eventsSubscribe),
]
