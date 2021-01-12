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

import { defineMessages } from 'react-intl'

import CreateFetchSelect from '@console/containers/fetch-select'

import { getGsFrequencyPlans, getNsFrequencyPlans } from '@console/store/actions/configuration'

import {
  selectGsFrequencyPlans,
  selectNsFrequencyPlans,
  selectFrequencyPlansError,
  selectFrequencyPlansFetching,
} from '@console/store/selectors/configuration'

const m = defineMessages({
  title: 'Frequency plan',
  warning: 'Frequency plans unavailable',
  description: 'The frequency plan used by the end device',
  none: 'Do not set a frequency plan',
})

const formatOptions = plans => plans.map(plan => ({ value: plan.id, label: plan.name }))

const CreateFrequencyPlansSelector = source =>
  CreateFetchSelect({
    fetchOptions: source === 'ns' ? getNsFrequencyPlans : getGsFrequencyPlans,
    optionsSelector: source === 'ns' ? selectNsFrequencyPlans : selectGsFrequencyPlans,
    fetchingSelector: selectFrequencyPlansFetching,
    errorSelector: selectFrequencyPlansError,
    defaultWarning: m.warning,
    defaultTitle: m.title,
    optionsFormatter: formatOptions,
    defaultDescription: m.description,
    additionalOptions: source === 'gs' ? [{ value: 'no-frequency-plan', label: m.none }] : [],
  })

export const GsFrequencyPlansSelect = CreateFrequencyPlansSelector('gs')
export const NsFrequencyPlansSelect = CreateFrequencyPlansSelector('ns')
