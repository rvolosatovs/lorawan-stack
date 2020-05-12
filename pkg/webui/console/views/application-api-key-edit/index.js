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

import { withBreadcrumb } from '@ttn-lw/components/breadcrumbs/context'
import PageTitle from '@ttn-lw/components/page-title'
import Breadcrumb from '@ttn-lw/components/breadcrumbs/breadcrumb'

import withRequest from '@ttn-lw/lib/components/with-request'

import { ApiKeyEditForm } from '@console/components/api-key-form'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import { getApiKey } from '@console/store/actions/api-keys'

import {
  selectSelectedApplicationId,
  selectApplicationRights,
  selectApplicationPseudoRights,
  selectApplicationRightsError,
  selectApplicationRightsFetching,
} from '@console/store/selectors/applications'
import {
  selectSelectedApiKey,
  selectApiKeyError,
  selectApiKeyFetching,
} from '@console/store/selectors/api-keys'

@connect(
  function(state, props) {
    const { apiKeyId } = props.match.params

    const keyFetching = selectApiKeyFetching(state)
    const rightsFetching = selectApplicationRightsFetching(state)
    const keyError = selectApiKeyError(state)
    const rightsError = selectApplicationRightsError(state)

    return {
      keyId: apiKeyId,
      appId: selectSelectedApplicationId(state),
      apiKey: selectSelectedApiKey(state),
      rights: selectApplicationRights(state),
      pseudoRights: selectApplicationPseudoRights(state),
      fetching: keyFetching || rightsFetching,
      error: keyError || rightsError,
    }
  },
  dispatch => ({
    getApiKey(appId, keyId) {
      dispatch(getApiKey('application', appId, keyId))
    },
    deleteSuccess: appId => dispatch(replace(`/applications/${appId}/api-keys`)),
  }),
)
@withRequest(
  ({ getApiKey, appId, keyId }) => getApiKey(appId, keyId),
  ({ fetching, apiKey }) => fetching || !Boolean(apiKey),
)
@withBreadcrumb('apps.single.api-keys.edit', function(props) {
  const { appId, keyId } = props

  return (
    <Breadcrumb path={`/applications/${appId}/api-keys/${keyId}`} content={sharedMessages.edit} />
  )
})
export default class ApplicationApiKeyEdit extends React.Component {
  static propTypes = {
    apiKey: PropTypes.apiKey.isRequired,
    appId: PropTypes.string.isRequired,
    deleteSuccess: PropTypes.func.isRequired,
    keyId: PropTypes.string.isRequired,
    pseudoRights: PropTypes.rights.isRequired,
    rights: PropTypes.rights.isRequired,
  }

  constructor(props) {
    super(props)

    this.deleteApplicationKey = id => api.application.apiKeys.delete(props.appId, id)
    this.editApplicationKey = key => api.application.apiKeys.update(props.appId, props.keyId, key)
  }

  @bind
  onDeleteSuccess() {
    const { appId, deleteSuccess } = this.props

    deleteSuccess(appId)
  }

  render() {
    const { apiKey, rights, pseudoRights } = this.props

    return (
      <Container>
        <PageTitle title={sharedMessages.keyEdit} />
        <Row>
          <Col lg={8} md={12}>
            <ApiKeyEditForm
              rights={rights}
              pseudoRights={pseudoRights}
              apiKey={apiKey}
              onEdit={this.editApplicationKey}
              onDelete={this.deleteApplicationKey}
              onDeleteSuccess={this.onDeleteSuccess}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
