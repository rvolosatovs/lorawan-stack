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
import { Col, Row, Container } from 'react-grid-system'

import PageTitle from '@ttn-lw/components/page-title'
import Breadcrumb from '@ttn-lw/components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '@ttn-lw/components/breadcrumbs/context'

import GatewayLocationForm from '@console/containers/gateway-location-form'

import withFeatureRequirement from '@console/lib/components/with-feature-requirement'

import sharedMessages from '@ttn-lw/lib/shared-messages'

import { mayViewOrEditGatewayLocation } from '@console/lib/feature-checks'

const GatewayLocation = () => {
  return (
    <Container>
      <PageTitle title={sharedMessages.location} />
      <Row>
        <Col lg={8} md={12}>
          <GatewayLocationForm />
        </Col>
      </Row>
    </Container>
  )
}

export default withBreadcrumb('gateway.single.data', props => {
  const { gtwId } = props
  return <Breadcrumb path={`/gateways/${gtwId}/location`} content={sharedMessages.location} />
})(
  withFeatureRequirement(mayViewOrEditGatewayLocation, {
    redirect: ({ gtwId }) => `/gateways/${gtwId}`,
  })(GatewayLocation),
)
