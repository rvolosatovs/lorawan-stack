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
import { Container, Col, Row } from 'react-grid-system'
import { replace } from 'connected-react-router'

import api from '@console/api'

import PageTitle from '@ttn-lw/components/page-title'
import Breadcrumb from '@ttn-lw/components/breadcrumbs/breadcrumb'
import toast from '@ttn-lw/components/toast'
import { withBreadcrumb } from '@ttn-lw/components/breadcrumbs/context'

import withRequest from '@ttn-lw/lib/components/with-request'

import CollaboratorForm from '@console/components/collaborator-form'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import { getCollaborator } from '@console/store/actions/collaborators'

import {
  selectSelectedGatewayId,
  selectGatewayRights,
  selectGatewayPseudoRights,
  selectGatewayRightsFetching,
  selectGatewayRightsError,
} from '@console/store/selectors/gateways'
import {
  selectUserCollaborator,
  selectOrganizationCollaborator,
  selectCollaboratorFetching,
  selectCollaboratorError,
} from '@console/store/selectors/collaborators'

@connect(
  (state, props) => {
    const gtwId = selectSelectedGatewayId(state, props)

    const { collaboratorId, collaboratorType } = props.match.params

    const collaborator =
      collaboratorType === 'user'
        ? selectUserCollaborator(state)
        : selectOrganizationCollaborator(state)

    const fetching = selectGatewayRightsFetching(state) || selectCollaboratorFetching(state)
    const error = selectGatewayRightsError(state) || selectCollaboratorError(state)

    return {
      collaboratorId,
      collaboratorType,
      collaborator,
      gtwId,
      rights: selectGatewayRights(state),
      pseudoRights: selectGatewayPseudoRights(state),
      fetching,
      error,
    }
  },
  dispatch => ({
    getCollaborator(gtwId, collaboratorId, isUser) {
      dispatch(getCollaborator('gateway', gtwId, collaboratorId, isUser))
    },
    redirectToList(gtwId) {
      dispatch(replace(`/gateways/${gtwId}/collaborators`))
    },
  }),
  (stateProps, dispatchProps, ownProps) => ({
    ...stateProps,
    ...dispatchProps,
    ...ownProps,
    getGatewayCollaborator: () =>
      dispatchProps.getCollaborator(
        stateProps.gtwId,
        stateProps.collaboratorId,
        stateProps.collaboratorType === 'user',
      ),
    redirectToList: () => dispatchProps.redirectToList(stateProps.gtwId),
  }),
)
@withRequest(
  ({ getGatewayCollaborator }) => getGatewayCollaborator(),
  ({ fetching, collaborator }) => fetching || !Boolean(collaborator),
)
@withBreadcrumb('gtws.single.collaborators.edit', props => {
  const { gtwId, collaboratorId, collaboratorType } = props

  return (
    <Breadcrumb
      path={`/gateways/${gtwId}/collaborators/${collaboratorType}/${collaboratorId}`}
      content={sharedMessages.edit}
    />
  )
})
export default class GatewayCollaboratorEdit extends React.Component {
  static propTypes = {
    collaborator: PropTypes.collaborator.isRequired,
    collaboratorId: PropTypes.string.isRequired,
    gtwId: PropTypes.string.isRequired,
    pseudoRights: PropTypes.rights.isRequired,
    redirectToList: PropTypes.func.isRequired,
    rights: PropTypes.rights.isRequired,
  }

  state = {
    error: '',
  }

  @bind
  handleSubmit(updatedCollaborator) {
    const { gtwId } = this.props

    return api.gateway.collaborators.update(gtwId, updatedCollaborator)
  }

  handleSubmitSuccess() {
    toast({
      message: sharedMessages.collaboratorUpdateSuccess,
      type: toast.types.SUCCESS,
    })
  }

  @bind
  async handleDelete(updatedCollaborator) {
    const { gtwId } = this.props

    return api.gateway.collaborators.remove(gtwId, updatedCollaborator)
  }

  render() {
    const { collaborator, collaboratorId, rights, redirectToList, pseudoRights } = this.props

    return (
      <Container>
        <PageTitle title={sharedMessages.collaboratorEdit} values={{ collaboratorId }} />
        <Row>
          <Col lg={8} md={12}>
            <CollaboratorForm
              error={this.state.error}
              onSubmit={this.handleSubmit}
              onSubmitSuccess={this.handleSubmitSuccess}
              onDelete={this.handleDelete}
              onDeleteSuccess={redirectToList}
              collaborator={collaborator}
              pseudoRights={pseudoRights}
              rights={rights}
              update
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
