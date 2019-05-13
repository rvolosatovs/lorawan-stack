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
  eventsSelector,
  errorSelector as eventsErrorSelector,
  statusSelector as eventsStatusSelector,
} from './events'

import {
  rightsSelector,
  errorSelector as rightsErrorSelector,
  fetchingSelector as rightsFetchingSelector,
} from './rights'

const ENTITY = 'applications'

export const applicationEventsSelector = eventsSelector(ENTITY)

export const applicationEventsErrorSelector = eventsErrorSelector(ENTITY)

export const applicationEventsStatusSelector = eventsStatusSelector(ENTITY)

export const applicationRightsSelector = rightsSelector(ENTITY)

export const applicationRightsErrorSelector = rightsErrorSelector(ENTITY)

export const applicationRightsFetchingSelector = rightsFetchingSelector(ENTITY)

