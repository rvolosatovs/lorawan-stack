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
import { connect } from 'react-redux'
import { Container, Col, Row } from 'react-grid-system'
import { Switch, Route, Redirect } from 'react-router'
import { defineMessages } from 'react-intl'

import { withBreadcrumb } from '@ttn-lw/components/breadcrumbs/context'
import Breadcrumb from '@ttn-lw/components/breadcrumbs/breadcrumb'
import Tab from '@ttn-lw/components/tabs'
import Notification from '@ttn-lw/components/notification'

import NotFoundRoute from '@ttn-lw/lib/components/not-found-route'

import DeviceUplinkPayloadFormatters from '@console/containers/device-payload-formatters/uplink'
import DeviceDownlinkPayloadFormatters from '@console/containers/device-payload-formatters/downlink'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import {
  selectApplicationLink,
  selectApplicationLinkFetching,
  selectSelectedApplicationId,
} from '@console/store/selectors/applications'
import { selectSelectedDeviceId } from '@console/store/selectors/devices'

import style from './device-payload-formatters.styl'

const m = defineMessages({
  infoUplinkText:
    'These payload formatters are executed on uplink messages from this end device and take precedence over application level payload formatters.',
  infoDownlinkText:
    'These payload formatters are executed on downlink messages to this end device and take precedence over application level payload formatters.',
})
@connect(state => {
  const link = selectApplicationLink(state)
  const fetching = selectApplicationLinkFetching(state)

  return {
    appId: selectSelectedApplicationId(state),
    devId: selectSelectedDeviceId(state),
    fetching: fetching || !link,
  }
})
@withBreadcrumb('device.single.payload-formatters', props => {
  const { appId, devId } = props
  return (
    <Breadcrumb
      path={`/applications/${appId}/devices/${devId}/payload-formatters`}
      content={sharedMessages.payloadFormatters}
    />
  )
})
export default class DevicePayloadFormatters extends Component {
  static propTypes = {
    location: PropTypes.location.isRequired,
    match: PropTypes.match.isRequired,
  }

  render() {
    const {
      match: { url },
      location: { pathname },
    } = this.props

    const tabs = [
      { title: sharedMessages.uplink, name: 'uplink', link: `${url}/uplink` },
      { title: sharedMessages.downlink, name: 'downlink', link: `${url}/downlink` },
    ]
    let deviceFormatterInfo
    if (pathname === `${url}/uplink`) {
      deviceFormatterInfo = (
        <Notification className={style.notification} small info content={m.infoUplinkText} />
      )
    } else if (pathname === `${url}/downlink`) {
      deviceFormatterInfo = (
        <Notification className={style.notification} small info content={m.infoDownlinkText} />
      )
    }

    return (
      <Container>
        <Row>
          <Col>
            <Tab className={style.tabs} tabs={tabs} divider />
          </Col>
        </Row>
        <Row>
          <Col>
            {deviceFormatterInfo}
            <Switch>
              <Redirect exact from={url} to={`${url}/uplink`} />
              <Route exact path={`${url}/uplink`} component={DeviceUplinkPayloadFormatters} />
              <Route exact path={`${url}/downlink`} component={DeviceDownlinkPayloadFormatters} />
              <NotFoundRoute />
            </Switch>
          </Col>
        </Row>
      </Container>
    )
  }
}
