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
import { Container, Row, Col } from 'react-grid-system'
import bind from 'autobind-decorator'

import sharedMessages from '../../../lib/shared-messages'
import IntlHelmet from '../../../lib/components/intl-helmet'
import ApiKeysTable from '../../containers/api-keys-table'
import { getOrganizationApiKeysList } from '../../store/actions/organizations'
import PropTypes from '../../../lib/prop-types'

import {
  selectOrganizationApiKeys,
  selectOrganizationApiKeysTotalCount,
  selectOrganizationApiKeysFetching,
} from '../../store/selectors/organizations'

import PAGE_SIZES from '../../constants/page-sizes'

class OrganizationApiKeysList extends React.Component {
  static propTypes = {
    match: PropTypes.match.isRequired,
  }

  constructor(props) {
    super(props)

    const { orgId } = props.match.params
    this.getOrganizationApiKeysList = filters => getOrganizationApiKeysList(orgId, filters)
  }

  @bind
  baseDataSelector(state) {
    const { orgId } = this.props.match.params

    const id = { id: orgId }
    return {
      keys: selectOrganizationApiKeys(state, id),
      totalCount: selectOrganizationApiKeysTotalCount(state, id),
      fetching: selectOrganizationApiKeysFetching(state),
    }
  }

  render() {
    return (
      <Container>
        <Row>
          <IntlHelmet title={sharedMessages.apiKeys} />
          <Col>
            <ApiKeysTable
              pageSize={PAGE_SIZES.REGULAR}
              baseDataSelector={this.baseDataSelector}
              getItemsAction={this.getOrganizationApiKeysList}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}

export default OrganizationApiKeysList
