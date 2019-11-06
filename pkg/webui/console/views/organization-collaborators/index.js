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

import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import sharedMessages from '../../../lib/shared-messages'
import ErrorView from '../../../lib/components/error-view'
import NotFoundRoute from '../../../lib/components/not-found-route'
import PropTypes from '../../../lib/prop-types'
import SubViewError from '../error/sub-view'
import OrganizationCollaboratorsList from '../organization-collaborators-list'
import OrganizationCollaboratorAdd from '../organization-collaborator-add'
import OrganizationCollaboratorEdit from '../organization-collaborator-edit'
import withFeatureRequirement from '../../lib/components/with-feature-requirement'

import { mayViewOrEditOrganizationApiKeys } from '../../lib/feature-checks'
import { selectSelectedOrganizationId } from '../../store/selectors/organizations'

@connect(state => ({ orgId: selectSelectedOrganizationId(state) }))
@withFeatureRequirement(mayViewOrEditOrganizationApiKeys, {
  redirect: ({ orgId }) => `/organizations/${orgId}`,
})
@withBreadcrumb('orgs.single.collaborators', function(props) {
  const { match } = props
  const { orgId } = match.params

  return (
    <Breadcrumb
      path={`/organizations/${orgId}/collaborators`}
      icon="organization"
      content={sharedMessages.collaborators}
    />
  )
})
class OrganizationCollaborators extends React.Component {
  static propTypes = {
    match: PropTypes.match.isRequired,
  }

  render() {
    const { match } = this.props

    return (
      <ErrorView ErrorComponent={SubViewError}>
        <Switch>
          <Route exact path={`${match.path}`} component={OrganizationCollaboratorsList} />
          <Route exact path={`${match.path}/add`} component={OrganizationCollaboratorAdd} />
          <Route
            path={`${match.path}/:collaboratorType(user|organization)/:collaboratorId`}
            component={OrganizationCollaboratorEdit}
          />
          <NotFoundRoute />
        </Switch>
      </ErrorView>
    )
  }
}

export default OrganizationCollaborators
