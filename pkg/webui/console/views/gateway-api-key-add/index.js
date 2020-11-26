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
import { replace } from 'connected-react-router'

import api from '@console/api'

import PageTitle from '@ttn-lw/components/page-title'
import Breadcrumb from '@ttn-lw/components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '@ttn-lw/components/breadcrumbs/context'

import { ApiKeyCreateForm } from '@console/components/api-key-form'

import withFeatureRequirement from '@console/lib/components/with-feature-requirement'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import { mayViewOrEditGatewayApiKeys } from '@console/lib/feature-checks'

import {
  selectSelectedGatewayId,
  selectGatewayRights,
  selectGatewayRightsError,
  selectGatewayRightsFetching,
  selectGatewayPseudoRights,
} from '@console/store/selectors/gateways'

@connect(
  state => ({
    gtwId: selectSelectedGatewayId(state),
    fetching: selectGatewayRightsFetching(state),
    error: selectGatewayRightsError(state),
    rights: selectGatewayRights(state),
    pseudoRights: selectGatewayPseudoRights(state),
  }),
  dispatch => ({
    navigateToList: gtwId => dispatch(replace(`/gateways/${gtwId}/api-keys`)),
  }),
)
@withFeatureRequirement(mayViewOrEditGatewayApiKeys, {
  redirect: ({ gtwId }) => `/gateway/${gtwId}`,
})
@withBreadcrumb('gtws.single.api-keys.add', function(props) {
  const gtwId = props.gtwId

  return <Breadcrumb path={`/gateways/${gtwId}/api-keys/add`} content={sharedMessages.add} />
})
export default class GatewayApiKeyAdd extends React.Component {
  static propTypes = {
    gtwId: PropTypes.string.isRequired,
    navigateToList: PropTypes.func.isRequired,
    pseudoRights: PropTypes.rights.isRequired,
    rights: PropTypes.rights.isRequired,
  }

  constructor(props) {
    super(props)

    this.createGatewayKey = key => api.gateway.apiKeys.create(props.gtwId, key)
  }

  @bind
  handleApprove() {
    const { navigateToList, gtwId } = this.props

    navigateToList(gtwId)
  }

  render() {
    const { rights, pseudoRights } = this.props

    return (
      <Container>
        <PageTitle title={sharedMessages.addApiKey} />
        <Row>
          <Col lg={8} md={12}>
            <ApiKeyCreateForm
              rights={rights}
              pseudoRights={pseudoRights}
              onCreate={this.createGatewayKey}
              onCreateSuccess={this.handleApprove}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
