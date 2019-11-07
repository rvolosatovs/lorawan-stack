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
import { Switch, Route } from 'react-router'

import sharedMessages from '../../../lib/shared-messages'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import ErrorView from '../../../lib/components/error-view'
import SubViewError from '../error/sub-view'
import NotFoundRoute from '../../../lib/components/not-found-route'
import GatewayCollaboratorsList from '../gateway-collaborators-list'
import GatewayCollaboratorAdd from '../gateway-collaborator-add'
import GatewayCollaboratorEdit from '../gateway-collaborator-edit'
import withFeatureRequirement from '../../lib/components/with-feature-requirement'
import { mayViewOrEditGatewayCollaborators } from '../../lib/feature-checks'

import { selectSelectedGatewayId } from '../../store/selectors/gateways'

@connect(state => ({ gtwId: selectSelectedGatewayId(state) }))
@withFeatureRequirement(mayViewOrEditGatewayCollaborators, {
  redirect: ({ gtwId }) => `/gateways/${gtwId}`,
})
@withBreadcrumb('gtws.single.collaborators', ({ gtwId }) => (
  <Breadcrumb
    path={`/gateways/${gtwId}/collaborators`}
    icon="organization"
    content={sharedMessages.collaborators}
  />
))
export default class GatewayCollaborators extends React.Component {
  render() {
    const { match } = this.props

    return (
      <ErrorView ErrorComponent={SubViewError}>
        <Switch>
          <Route exact path={`${match.path}`} component={GatewayCollaboratorsList} />
          <Route exact path={`${match.path}/add`} component={GatewayCollaboratorAdd} />
          <Route
            path={`${match.path}/:collaboratorType(user|organization)/:collaboratorId`}
            component={GatewayCollaboratorEdit}
          />
          <NotFoundRoute />
        </Switch>
      </ErrorView>
    )
  }
}
