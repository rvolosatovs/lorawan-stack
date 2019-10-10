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
import { Switch, Route } from 'react-router'
import { Col, Row, Container } from 'react-grid-system'

import sharedMessages from '../../../lib/shared-messages'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import Tabs from '../../../components/tabs'
import IntlHelmet from '../../../lib/components/intl-helmet'
import withRequest from '../../../lib/components/with-request'
import withEnv from '../../../lib/components/env'
import NotFoundRoute from '../../../lib/components/not-found-route'

import DeviceOverview from '../device-overview'
import DeviceData from '../device-data'
import DeviceGeneralSettings from '../device-general-settings'
import DeviceLocation from '../device-location'
import DevicePayloadFormatters from '../device-payload-formatters'

import { getDevice, stopDeviceEventsStream } from '../../store/actions/device'
import { selectSelectedApplicationId } from '../../store/selectors/applications'
import {
  selectSelectedDevice,
  selectDeviceFetching,
  selectGetDeviceError,
} from '../../store/selectors/device'

import style from './device.styl'

@connect(
  function(state, props) {
    const devId = props.match.params.devId
    const appId = selectSelectedApplicationId(state)
    const device = selectSelectedDevice(state)

    return {
      devId,
      appId,
      device,
      deviceName: device && device.name,
      devIds: device && device.ids,
      fetching: selectDeviceFetching(state),
      error: selectGetDeviceError(state),
    }
  },
  dispatch => ({
    getDevice: (appId, devId, selectors, config) =>
      dispatch(getDevice(appId, devId, selectors, config)),
    stopStream: id => dispatch(stopDeviceEventsStream(id)),
  }),
)
@withRequest(
  ({ appId, devId, getDevice }) =>
    getDevice(
      appId,
      devId,
      [
        'name',
        'description',
        'session',
        'version_ids',
        'root_keys',
        'frequency_plan_id',
        'mac_settings.resets_f_cnt',
        'resets_join_nonces',
        'supports_class_c',
        'supports_join',
        'lorawan_version',
        'lorawan_phy_version',
        'network_server_address',
        'application_server_address',
        'join_server_address',
        'locations',
        'formatters',
      ],
      { ignoreNotFound: true },
    ),
  ({ fetching, device }) => fetching || !Boolean(device),
)
@withBreadcrumb('device.single', function(props) {
  const { devId, appId } = props
  return (
    <Breadcrumb path={`/applications/${appId}/devices/${devId}`} icon="device" content={devId} />
  )
})
@withEnv
export default class Device extends React.Component {
  componentWillUnmount() {
    const { devIds, stopStream } = this.props

    stopStream(devIds)
  }

  render() {
    const {
      location: { pathname },
      match: {
        params: { appId },
      },
      devId,
      deviceName,
      env: { siteName },
    } = this.props

    const basePath = `/applications/${appId}/devices/${devId}`

    // Prevent default redirect to uplink when tab is already open
    const payloadFormattersLink = pathname.startsWith(`${basePath}/payload-formatters`)
      ? pathname
      : `${basePath}/payload-formatters`

    const tabs = [
      { title: sharedMessages.overview, name: 'overview', link: basePath },
      { title: sharedMessages.data, name: 'data', link: `${basePath}/data` },
      { title: sharedMessages.location, name: 'location', link: `${basePath}/location` },
      {
        title: sharedMessages.payloadFormatters,
        name: 'develop',
        link: payloadFormattersLink,
        exact: false,
      },
      {
        title: sharedMessages.generalSettings,
        name: 'general-settings',
        link: `${basePath}/general-settings`,
      },
    ]

    return (
      <React.Fragment>
        <IntlHelmet titleTemplate={`%s - ${deviceName || devId} - ${siteName}`} />
        <Container>
          <Row>
            <Col>
              <h2 className={style.title}>{deviceName || devId}</h2>
              <Tabs className={style.tabs} narrow tabs={tabs} />
            </Col>
          </Row>
        </Container>
        <hr className={style.rule} />
        <Switch>
          <Route exact path={basePath} component={DeviceOverview} />
          <Route exact path={`${basePath}/data`} component={DeviceData} />
          <Route exact path={`${basePath}/location`} component={DeviceLocation} />
          <Route exact path={`${basePath}/general-settings`} component={DeviceGeneralSettings} />
          <Route path={`${basePath}/payload-formatters`} component={DevicePayloadFormatters} />
          <NotFoundRoute />
        </Switch>
      </React.Fragment>
    )
  }
}
