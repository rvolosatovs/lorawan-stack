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

import { ApiKeyEditForm } from '../../components/api-key-form'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'

import Message from '../../../lib/components/message'
import sharedMessages from '../../../lib/shared-messages'
import IntlHelmet from '../../../lib/components/intl-helmet'
import PropTypes from '../../../lib/prop-types'

@withBreadcrumb('orgs.single.api-keys.edit', function(props) {
  const { orgId, keyId } = props

  return (
    <Breadcrumb
      path={`/organizations/${orgId}/api-keys/${keyId}`}
      icon="general_settings"
      content={sharedMessages.edit}
    />
  )
})
class OrganizationApiKeyEdit extends React.Component {
  static propTypes = {
    apiKey: PropTypes.apiKey.isRequired,
    deleteOrganizationApiKey: PropTypes.func.isRequired,
    deleteOrganizationApiKeySuccess: PropTypes.func.isRequired,
    pseudoRights: PropTypes.rights,
    rights: PropTypes.rights.isRequired,
    updateOrganizationApiKey: PropTypes.func.isRequired,
  }

  static defaultProps = {
    pseudoRights: [],
  }

  render() {
    const {
      apiKey,
      rights,
      pseudoRights,
      deleteOrganizationApiKey,
      deleteOrganizationApiKeySuccess,
      updateOrganizationApiKey,
    } = this.props

    return (
      <Container>
        <Row>
          <Col>
            <IntlHelmet title={sharedMessages.keyEdit} />
            <Message component="h2" content={sharedMessages.keyEdit} />
          </Col>
        </Row>
        <Row>
          <Col lg={8} md={12}>
            <ApiKeyEditForm
              rights={rights}
              pseudoRights={pseudoRights}
              apiKey={apiKey}
              onEdit={updateOrganizationApiKey}
              onDelete={deleteOrganizationApiKey}
              onDeleteSuccess={deleteOrganizationApiKeySuccess}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}

export default OrganizationApiKeyEdit
