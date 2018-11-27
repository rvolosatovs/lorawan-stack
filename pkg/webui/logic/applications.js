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

import api from '../api'
import * as applications from '../actions/applications'

const getApplicationsLogic = createLogic({
  type: [
    applications.GET_APPS_LIST,
    applications.CHANGE_APPS_ORDER,
    applications.CHANGE_APPS_PAGE,
    applications.CHANGE_APPS_TAB,
    applications.SEARCH_APPS_LIST,
  ],
  latest: true,
  async process ({ getState, action }, dispatch, done) {
    const { filters } = action

    try {
      const data = filters.query
        ? await api.v3.is.applications.search(filters)
        : await api.v3.is.applications.list(filters)
      dispatch(applications.getApplicationsSuccess(data.applications, data.totalCount))
    } catch (error) {
      dispatch(applications.getApplicationsFailure(error))
    }

    done()
  },
})

export default [
  getApplicationsLogic,
]
