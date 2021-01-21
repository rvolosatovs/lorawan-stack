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

import Status from '@ttn-lw/components/status'
import Icon from '@ttn-lw/components/icon'

import FetchTable from '@ttn-lw/containers/fetch-table'

import Message from '@ttn-lw/lib/components/message'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'
import { getUserId } from '@ttn-lw/lib/selectors/id'

import { checkFromState, mayManageUsers } from '@console/lib/feature-checks'

import { getUsersList } from '@console/store/actions/users'

import { selectUserId } from '@console/store/selectors/user'
import {
  selectUsers,
  selectUsersTotalCount,
  selectUsersFetching,
  selectUsersError,
} from '@console/store/selectors/users'

import style from './users-table.styl'

@connect(state => ({
  currentUserId: selectUserId(state),
}))
export default class UsersTable extends Component {
  static propTypes = {
    currentUserId: PropTypes.string.isRequired,
    pageSize: PropTypes.number.isRequired,
  }

  constructor(props) {
    super(props)
    this.headers = [
      {
        name: 'ids.user_id',
        displayName: sharedMessages.id,
        width: 28,
        sortable: true,
        sortKey: 'user_id',
        render(ids) {
          const userId = getUserId({ ids })
          if (userId === props.currentUserId) {
            return (
              <span>
                {userId}{' '}
                <Message className={style.hint} content={sharedMessages.currentUserIndicator} />
              </span>
            )
          }
          return userId
        },
      },
      {
        name: 'name',
        displayName: sharedMessages.name,
        width: 22,
        sortable: true,
      },
      {
        name: 'primary_email_address',
        displayName: sharedMessages.email,
        width: 28,
        sortable: true,
      },
      {
        name: 'state',
        displayName: sharedMessages.state,
        width: 15,
        sortable: true,
        render(state) {
          let indicator = 'unknown'
          let label = sharedMessages.notSet
          switch (state) {
            case 'STATE_APPROVED':
              indicator = 'good'
              label = sharedMessages.stateApproved
              break
            case 'STATE_REQUESTED':
              indicator = 'mediocre'
              label = sharedMessages.stateRequested
              break
            case 'STATE_REJECTED':
              indicator = 'bad'
              label = sharedMessages.stateRejected
              break
            case 'STATE_FLAGGED':
              indicator = 'bad'
              label = sharedMessages.stateFlagged
              break
            case 'STATE_SUSPENDED':
              indicator = 'bad'
              label = sharedMessages.stateSuspended
              break
          }

          return <Status status={indicator} label={label} pulse={false} />
        },
      },
      {
        name: 'admin',
        displayName: sharedMessages.admin,
        width: 7,
        render(isAdmin) {
          if (isAdmin) {
            return <Icon className={style.icon} icon="check" />
          }

          return null
        },
      },
    ]

    this.getUsersList = params =>
      getUsersList(params, ['name', 'primary_email_address', 'state', 'admin'])
  }

  baseDataSelector(state) {
    return {
      users: selectUsers(state),
      totalCount: selectUsersTotalCount(state),
      fetching: selectUsersFetching(state),
      error: selectUsersError(state),
      mayAdd: checkFromState(mayManageUsers, state),
    }
  }

  render() {
    const { pageSize } = this.props

    return (
      <FetchTable
        entity="users"
        headers={this.headers}
        addMessage={sharedMessages.userAdd}
        tableTitle={<Message content={sharedMessages.users} />}
        getItemsAction={this.getUsersList}
        searchItemsAction={this.getUsersList}
        baseDataSelector={this.baseDataSelector}
        pageSize={pageSize}
        searchable
      />
    )
  }
}
