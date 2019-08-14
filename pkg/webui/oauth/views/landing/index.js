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
import { connect } from 'react-redux'
import bind from 'autobind-decorator'
import { push } from 'connected-react-router'

import Button from '../../../components/button'
import WithAuth from '../../../lib/components/with-auth'

import { logout } from '../../store/actions/user'
import sharedMessages from '../../../lib/shared-messages'

@connect(state => ({
  user: state.user.user,
}))
@bind
export default class OAuth extends React.PureComponent {
  async handleLogout() {
    const { dispatch } = this.props

    dispatch(logout())
  }

  async handleUpdatePassword() {
    const { dispatch } = this.props

    dispatch(push('/update-password'))
  }

  render() {
    const { user = { ids: {} } } = this.props

    return (
      <WithAuth>
        <div>
          You are logged in as {user.ids.user_id}.{' '}
          <Button message={sharedMessages.logout} onClick={this.handleLogout} />
          <Button message={sharedMessages.changePassword} onClick={this.handleUpdatePassword} />
        </div>
      </WithAuth>
    )
  }
}
