// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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
import { defineMessages } from 'react-intl'
import { useSelector } from 'react-redux'

import LoRaCloudImage from '@assets/misc/lora-cloud.png'

import PageTitle from '@ttn-lw/components/page-title'
import Breadcrumb from '@ttn-lw/components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '@ttn-lw/components/breadcrumbs/context'
import Link from '@ttn-lw/components/link'
import Collapse from '@ttn-lw/components/collapse'

import Message from '@ttn-lw/lib/components/message'
import ErrorView from '@ttn-lw/lib/components/error-view'

import LoRaCloudDASForm from '@console/containers/lora-cloud-das-form'
import LoRaCloudGLSForm from '@console/containers/lora-cloud-gls-form'

import Require from '@console/lib/components/require'

import SubViewError from '@console/views/sub-view-error'

import sharedMessages from '@ttn-lw/lib/shared-messages'

import { mayViewOrEditApplicationPackages } from '@console/lib/feature-checks'

import { selectSelectedApplicationId } from '@console/store/selectors/applications'

import style from './application-integrations-lora-cloud.styl'

const m = defineMessages({
  loraCloudInfoText:
    'Lora Cloud provides value added APIs that enable simple solutions for common tasks related to LoRaWAN networks and LoRa-based devices. You can setup our LoRaCloud integrations below.',
  officialLoRaCloudDocumentation: 'Official LoRa Cloud documentation',
  furtherResources: 'Further resources',
  setToken: 'Set LoRa Cloud token',
  dasDescription:
    'With the LoRa Cloud Device & Application Services protocol, you can manage common device functionality at the application layer for LoRaWAN-enabled devices.',
  glsDescription:
    'LoRa Cloud Geolocation is a simple cloud API that can be easily integrated with The Things Stack to enable estimating the location of any LoRa-based device.',
})

const LoRaCloud = () => {
  const appId = useSelector(selectSelectedApplicationId)
  return (
    <Require
      featureCheck={mayViewOrEditApplicationPackages}
      otherwise={{ redirect: `/applications/${appId}` }}
    >
      <ErrorView ErrorComponent={SubViewError}>
        <Container>
          <PageTitle title="LoRa Cloud Device & Application Services" />
          <Row>
            <Col lg={8} md={12}>
              <img className={style.logo} src={LoRaCloudImage} alt="LoRa Cloud" />
              <Message content={m.loraCloudInfoText} className={style.info} />
              <div>
                <Message
                  component="h4"
                  content={m.furtherResources}
                  className={style.furtherResources}
                />
                <Link.DocLink
                  path="/integrations/application-packages/lora-cloud-device-and-application-services/"
                  secondary
                >
                  Device & Application Services
                </Link.DocLink>
                {' | '}
                <Link.Anchor href="https://www.loracloud.com" external secondary>
                  <Message content={m.officialLoRaCloudDocumentation} />
                </Link.Anchor>
              </div>
              <hr className={style.hRule} />
              <Collapse title="Geolocation" description={m.glsDescription}>
                <Message component="h3" content={m.setToken} />
                <LoRaCloudGLSForm />
              </Collapse>
              <Collapse title="Device & Application Services" description={m.dasDescription}>
                <Message component="h3" content={m.setToken} />
                <LoRaCloudDASForm />
              </Collapse>
            </Col>
          </Row>
        </Container>
      </ErrorView>
    </Require>
  )
}

export default withBreadcrumb('apps.single.integrations.lora-cloud', props => {
  const { appId } = props

  return (
    <Breadcrumb
      path={`/applications/${appId}/integrations/lora-cloud`}
      content={sharedMessages.loraCloud}
    />
  )
})(LoRaCloud)
