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
import { Container, Col, Row } from 'react-grid-system'
import bind from 'autobind-decorator'
import { connect } from 'react-redux'
import { push } from 'connected-react-router'

import api from '@console/api'

import PageTitle from '@ttn-lw/components/page-title'
import Breadcrumb from '@ttn-lw/components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '@ttn-lw/components/breadcrumbs/context'

import CollaboratorForm from '@console/components/collaborator-form'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import {
  selectSelectedApplicationId,
  selectApplicationRights,
  selectApplicationPseudoRights,
  selectApplicationRightsFetching,
  selectApplicationRightsError,
} from '@console/store/selectors/applications'
import { selectCollaborators } from '@console/store/selectors/collaborators'

@connect(
  state => ({
    appId: selectSelectedApplicationId(state),
    collaborators: selectCollaborators(state),
    rights: selectApplicationRights(state),
    pseudoRights: selectApplicationPseudoRights(state),
    fetching: selectApplicationRightsFetching(state),
    error: selectApplicationRightsError(state),
  }),
  (dispatch, ownProps) => ({
    redirectToList: appId => dispatch(push(`/applications/${appId}/collaborators`)),
  }),
  (stateProps, dispatchProps, ownProps) => ({
    ...stateProps,
    ...dispatchProps,
    ...ownProps,
    redirectToList: () => dispatchProps.redirectToList(stateProps.appId),
  }),
)
@withBreadcrumb('apps.single.collaborators.add', function(props) {
  const appId = props.appId
  return (
    <Breadcrumb path={`/applications/${appId}/collaborators/add`} content={sharedMessages.add} />
  )
})
export default class ApplicationCollaboratorAdd extends React.Component {
  static propTypes = {
    appId: PropTypes.string.isRequired,
    pseudoRights: PropTypes.rights.isRequired,
    redirectToList: PropTypes.func.isRequired,
    rights: PropTypes.rights.isRequired,
  }

  state = {
    error: '',
  }

  @bind
  async handleSubmit(collaborator) {
    const { appId } = this.props

    await api.application.collaborators.add(appId, collaborator)
  }

  render() {
    const { rights, pseudoRights, redirectToList } = this.props

    return (
      <Container>
        <PageTitle title={sharedMessages.addCollaborator} />
        <Row>
          <Col lg={8} md={12}>
            <CollaboratorForm
              error={this.state.error}
              onSubmit={this.handleSubmit}
              onSubmitSuccess={redirectToList}
              pseudoRights={pseudoRights}
              rights={rights}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
