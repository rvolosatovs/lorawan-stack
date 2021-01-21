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

import React, { Component } from 'react'
import { Container, Row, Col } from 'react-grid-system'

import PAGE_SIZES from '@ttn-lw/constants/page-sizes'

import IntlHelmet from '@ttn-lw/lib/components/intl-helmet'

import PubsubsTable from '@console/containers/pubsubs-table'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

export default class ApplicationPubsubsList extends Component {
  static propTypes = {
    match: PropTypes.match.isRequired,
  }

  render() {
    const { appId } = this.props.match.params

    return (
      <Container>
        <Row>
          <IntlHelmet title={sharedMessages.integrations} />
          <Col>
            <PubsubsTable pageSize={PAGE_SIZES.REGULAR} appId={appId} />
          </Col>
        </Row>
      </Container>
    )
  }
}
