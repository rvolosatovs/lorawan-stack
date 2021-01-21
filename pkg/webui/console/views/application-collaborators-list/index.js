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

import PAGE_SIZES from '@ttn-lw/constants/page-sizes'

import IntlHelmet from '@ttn-lw/lib/components/intl-helmet'

import CollaboratorsTable from '@console/containers/collaborators-table'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import { getCollaboratorsList } from '@console/store/actions/collaborators'

import {
  selectCollaborators,
  selectCollaboratorsTotalCount,
  selectCollaboratorsFetching,
  selectCollaboratorsError,
} from '@console/store/selectors/collaborators'

export default class ApplicationCollaborators extends React.Component {
  static propTypes = {
    match: PropTypes.match.isRequired,
  }

  constructor(props) {
    super(props)

    const { appId } = props.match.params
    this.getCollaboratorsList = filters => getCollaboratorsList('application', appId, filters)
  }

  @bind
  baseDataSelector(state) {
    const { appId } = this.props.match.params
    const id = { id: appId }

    return {
      collaborators: selectCollaborators(state, id),
      fetching: selectCollaboratorsFetching(state),
      totalCount: selectCollaboratorsTotalCount(state, id),
      error: selectCollaboratorsError(state),
    }
  }

  render() {
    return (
      <Container>
        <Row>
          <IntlHelmet title={sharedMessages.collaborators} />
          <Col>
            <CollaboratorsTable
              pageSize={PAGE_SIZES.REGULAR}
              baseDataSelector={this.baseDataSelector}
              getItemsAction={this.getCollaboratorsList}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
