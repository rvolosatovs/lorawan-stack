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
import { Switch, Route } from 'react-router-dom'
import { connect } from 'react-redux'

import WithAuth from '@ttn-lw/lib/components/with-auth'

import Overview from '@console/views/overview'
import Applications from '@console/views/applications'
import Gateways from '@console/views/gateways'
import Organizations from '@console/views/organizations'
import Admin from '@console/views/admin'
import { FullViewErrorInner } from '@console/views/error'

import PropTypes from '../../../lib/prop-types'
import {
  selectUser,
  selectUserFetching,
  selectUserError,
  selectUserRights,
  selectUserIsAdmin,
} from '../../store/selectors/user'

import style from './landing.styl'

const GenericNotFound = () => <FullViewErrorInner error={{ statusCode: 404 }} />

@connect(state => ({
  user: selectUser(state),
  fetching: selectUserFetching(state),
  error: selectUserError(state),
  rights: selectUserRights(state),
  isAdmin: selectUserIsAdmin(state),
}))
export default class Landing extends React.PureComponent {
  static propTypes = {
    error: PropTypes.error,
    fetching: PropTypes.bool.isRequired,
    isAdmin: PropTypes.bool,
    rights: PropTypes.rights,
    user: PropTypes.user,
  }
  static defaultProps = {
    user: undefined,
    error: undefined,
    rights: undefined,
    isAdmin: undefined,
  }
  render() {
    const { user, fetching, error, rights, isAdmin } = this.props

    return (
      <WithAuth
        user={user}
        fetching={fetching}
        error={error}
        errorComponent={FullViewErrorInner}
        rights={rights}
        isAdmin={isAdmin}
      >
        <div className={style.container}>
          <Switch>
            <Route exact path="/" component={Overview} />
            <Route path="/applications" component={Applications} />
            <Route path="/gateways" component={Gateways} />
            <Route path="/organizations" component={Organizations} />
            <Route path="/admin" component={Admin} />
            <Route component={GenericNotFound} />
          </Switch>
        </div>
      </WithAuth>
    )
  }
}
