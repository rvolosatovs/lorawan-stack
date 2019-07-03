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

import React from 'react'
import { defineMessages } from 'react-intl'
import bind from 'autobind-decorator'

import sharedMessages from '../../../lib/shared-messages'
import FetchTable from '../fetch-table'

import { getGatewaysList } from '../../../console/store/actions/gateways'

const m = defineMessages({
  add: 'Add Gateway',
  gtwId: 'Gateway ID',
  gtwEUI: 'Gateway EUI',
  freqPlan: 'Frequency Plan',
})

const headers = [
  {
    name: 'ids.gateway_id',
    displayName: m.gtwId,
    width: 25,
  },
  {
    name: 'name',
    displayName: sharedMessages.name,
    width: 25,
  },
  {
    name: 'ids.eui',
    displayName: m.gtwEUI,
    width: 25,
  },
  {
    name: 'frequency_plan_id',
    displayName: m.freqPlan,
    width: 15,
  },
  {
    name: 'antennasCount',
    displayName: sharedMessages.antennas,
    centered: true,
    width: 10,
  },
]

@bind
export default class GatewaysTable extends React.Component {

  constructor (props) {
    super(props)

    this.getGatewaysList = params => getGatewaysList(params, [
      'name',
      'description',
      'frequency_plan_id',
    ])
  }

  baseDataSelector ({ gateways }) {
    return gateways
  }

  render () {
    return (
      <FetchTable
        entity="gateways"
        addMessage={m.add}
        headers={headers}
        getItemsAction={this.getGatewaysList}
        searchItemsAction={this.getGatewaysList}
        baseDataSelector={this.baseDataSelector}
        {...this.props}
      />
    )
  }
}
