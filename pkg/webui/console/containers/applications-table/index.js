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

import React, { Component } from 'react'
import bind from 'autobind-decorator'

import Message from '@ttn-lw/lib/components/message'

import FetchTable from '@console/containers/fetch-table'

import sharedMessages from '@ttn-lw/lib/shared-messages'

import { checkFromState, mayCreateApplications } from '@console/lib/feature-checks'

import { getApplicationsList } from '@console/store/actions/applications'

import {
  selectApplications,
  selectApplicationsTotalCount,
  selectApplicationsFetching,
  selectApplicationsError,
} from '@console/store/selectors/applications'

const headers = [
  {
    name: 'ids.application_id',
    displayName: sharedMessages.id,
    width: 25,
    sortable: true,
    sortKey: 'application_id',
  },
  {
    name: 'name',
    displayName: sharedMessages.name,
    width: 25,
    sortable: true,
  },
  {
    name: 'description',
    displayName: sharedMessages.description,
    width: 50,
  },
]

@bind
export default class ApplicationsTable extends Component {
  constructor(props) {
    super(props)

    this.getApplicationsList = params => getApplicationsList(params, ['name', 'description'])
  }

  baseDataSelector(state) {
    return {
      applications: selectApplications(state),
      totalCount: selectApplicationsTotalCount(state),
      fetching: selectApplicationsFetching(state),
      error: selectApplicationsError(state),
      mayAdd: checkFromState(mayCreateApplications, state),
    }
  }

  render() {
    return (
      <FetchTable
        entity="applications"
        headers={headers}
        addMessage={sharedMessages.addApplication}
        tableTitle={<Message content={sharedMessages.applications} />}
        getItemsAction={this.getApplicationsList}
        baseDataSelector={this.baseDataSelector}
        searchable
        {...this.props}
      />
    )
  }
}
