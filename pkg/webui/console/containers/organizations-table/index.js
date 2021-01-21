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
import { connect } from 'react-redux'

import FetchTable from '@ttn-lw/containers/fetch-table'

import Message from '@ttn-lw/lib/components/message'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import { checkFromState, mayCreateOrganizations } from '@console/lib/feature-checks'

import { getOrganizationsList } from '@console/store/actions/organizations'

import { selectUserIsAdmin } from '@console/store/selectors/user'
import {
  selectOrganizations,
  selectOrganizationsTotalCount,
  selectOrganizationsFetching,
  selectOrganizationsError,
} from '@console/store/selectors/organizations'

const headers = [
  {
    name: 'ids.organization_id',
    displayName: sharedMessages.id,
    width: 25,
    sortable: true,
    sortKey: 'organization_id',
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

const OWNED_TAB = 'owned'
const ALL_TAB = 'all'
const tabs = [
  {
    title: sharedMessages.organizations,
    name: OWNED_TAB,
  },
  {
    title: sharedMessages.allAdmin,
    name: ALL_TAB,
  },
]

class OrganizationsTable extends Component {
  static propTypes = {
    isAdmin: PropTypes.bool.isRequired,
    pageSize: PropTypes.number.isRequired,
  }

  constructor(props) {
    super(props)

    this.getOrganizationsList = params => {
      const { tab, query } = params

      return getOrganizationsList(params, ['name', 'description'], {
        isSearch: tab === ALL_TAB || query.length > 0,
      })
    }
  }

  baseDataSelector(state) {
    return {
      organizations: selectOrganizations(state),
      totalCount: selectOrganizationsTotalCount(state),
      fetching: selectOrganizationsFetching(state),
      error: selectOrganizationsError(state),
      mayAdd: checkFromState(mayCreateOrganizations, state),
    }
  }

  render() {
    const { pageSize, isAdmin } = this.props

    return (
      <FetchTable
        entity="organizations"
        headers={headers}
        addMessage={sharedMessages.addOrganization}
        tableTitle={<Message content={sharedMessages.organizations} />}
        getItemsAction={this.getOrganizationsList}
        baseDataSelector={this.baseDataSelector}
        pageSize={pageSize}
        searchable
        tabs={isAdmin ? tabs : []}
      />
    )
  }
}

export default connect(state => ({
  isAdmin: selectUserIsAdmin(state),
}))(OrganizationsTable)
