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
import { Switch, Route } from 'react-router'

import sharedMessages from '../../../lib/shared-messages'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import ErrorView from '../../../lib/components/error-view'
import SubViewError from '../error/sub-view'

import ApplicationCollaboratorsList from '../application-collaborators-list'
import ApplicationCollaboratorAdd from '../application-collaborator-add'
import ApplicationCollaboratorEdit from '../application-collaborator-edit'

@withBreadcrumb('apps.single.api-keys', function (props) {
  const { match } = props
  const appId = match.params.appId

  return (
    <Breadcrumb
      path={`/console/applications/${appId}/collaborators`}
      icon="organization"
      content={sharedMessages.collaborators}
    />
  )
})
export default class ApplicationCollaborators extends React.Component {

  render () {
    const { match } = this.props

    return (
      <ErrorView ErrorComponent={SubViewError}>
        <Switch>
          <Route exact path={`${match.path}`} component={ApplicationCollaboratorsList} />
          <Route path={`${match.path}/add`} component={ApplicationCollaboratorAdd} />
          <Route path={`${match.path}/:collaboratorId`} component={ApplicationCollaboratorEdit} />
        </Switch>
      </ErrorView>
    )
  }
}
