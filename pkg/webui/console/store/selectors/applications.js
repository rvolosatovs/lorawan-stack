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

import {
  GET_APPS_LIST_BASE,
  GET_APPS_RIGHTS_LIST_BASE,
  GET_APP_BASE,
} from '../actions/applications'
import { GET_APP_LINK_BASE } from '../actions/link'
import {
  createPaginationIdsSelectorByEntity,
  createPaginationTotalCountSelectorByEntity,
} from './pagination'
import {
  createEventsSelector,
  createEventsErrorSelector,
  createEventsStatusSelector,
} from './events'
import { createRightsSelector, createPseudoRightsSelector } from './rights'
import { createFetchingSelector } from './fetching'
import { createErrorSelector } from './error'

const ENTITY = 'applications'

// Application
export const selectApplicationStore = state => state.applications
export const selectApplicationEntitiesStore = state => selectApplicationStore(state).entities
export const selectApplicationById = (state, id) => selectApplicationEntitiesStore(state)[id]
export const selectSelectedApplicationId = state =>
  selectApplicationStore(state).selectedApplication
export const selectSelectedApplication = state =>
  selectApplicationById(state, selectSelectedApplicationId(state))
export const selectApplicationFetching = createFetchingSelector(GET_APP_BASE)
export const selectApplicationError = createErrorSelector(GET_APP_BASE)
export const selectApplicationDeviceCount = state =>
  selectApplicationStore(state).applicationDeviceCount

// Applications
const selectAppsIds = createPaginationIdsSelectorByEntity(ENTITY)
const selectAppsTotalCount = createPaginationTotalCountSelectorByEntity(ENTITY)
const selectAppsFetching = createFetchingSelector(GET_APPS_LIST_BASE)
const selectAppsError = createErrorSelector(GET_APPS_LIST_BASE)

export const selectApplications = state =>
  selectAppsIds(state).map(id => selectApplicationById(state, id))
export const selectApplicationsTotalCount = state => selectAppsTotalCount(state)
export const selectApplicationsFetching = state => selectAppsFetching(state)
export const selectApplicationsError = state => selectAppsError(state)

// Events
export const selectApplicationEvents = createEventsSelector(ENTITY)
export const selectApplicationEventsError = createEventsErrorSelector(ENTITY)
export const selectApplicationEventsStatus = createEventsStatusSelector(ENTITY)

// Rights
export const selectApplicationRights = createRightsSelector(ENTITY)
export const selectApplicationPseudoRights = createPseudoRightsSelector(ENTITY)
export const selectApplicationRightsError = createErrorSelector(GET_APPS_RIGHTS_LIST_BASE)
export const selectApplicationRightsFetching = createFetchingSelector(GET_APPS_RIGHTS_LIST_BASE)

// Link
const selectLinkStore = state => state.link
export const selectApplicationLink = state => selectLinkStore(state).link
export const selectApplicationLinkIndicator = state => selectLinkStore(state).linked
export const selectApplicationLinkStats = state => selectLinkStore(state).stats
export const selectApplicationLinkFetching = createFetchingSelector(GET_APP_LINK_BASE)
export const selectApplicationLinkError = createErrorSelector(GET_APP_LINK_BASE)
export const selectApplicationLinkFormatters = function(state) {
  const link = selectApplicationLink(state) || {}

  return link.default_formatters
}
export const selectApplicationIsLinked = function(state) {
  const linkStore = selectLinkStore(state)
  const link = selectApplicationLink(state) || {}
  const error = selectApplicationLinkError(state)
  const stats = selectApplicationLinkStats(state)

  const hasBase = Boolean(link.api_key)
  const hasError = Boolean(error)
  const isLinked = linkStore.linked
  const hasStats = Boolean(stats)

  return hasBase && !hasError && isLinked && hasStats
}
