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
import { Container, Col, Row } from 'react-grid-system'
import bind from 'autobind-decorator'
import { connect } from 'react-redux'
import { push } from 'connected-react-router'

import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import Message from '../../../lib/components/message'
import IntlHelmet from '../../../lib/components/intl-helmet'
import DeviceDataForm from '../../components/device-data-form'
import sharedMessages from '../../../lib/shared-messages'
import { selectSelectedApplicationId } from '../../store/selectors/applications'
import { getDeviceId } from '../../../lib/selectors/id'
import api from '../../api'

import style from './device-add.styl'

@withBreadcrumb('devices.add', function (props) {
  const { appId } = props.match.params
  return (
    <Breadcrumb
      path={`/console/applications/${appId}/devices/add`}
      icon="add"
      content={sharedMessages.add}
    />
  )
})
@connect(function (state) {
  return {
    device: state.device.device,
    appId: selectSelectedApplicationId(state),
  }
}, dispatch => ({
  redirectToList: (appId, deviceId) => dispatch(push(`/console/applications/${appId}/devices/${deviceId}`)),
}),
)
@bind
export default class DeviceAdd extends Component {

  state = {
    error: '',
  }

  async handleSubmit (values) {
    const { appId } = this.props
    const device = { ...values }

    // Clean values based on activation mode
    if (device.activation_mode === 'otaa') {
      delete device.mac_settings
      delete device.session
    } else {
      delete device.ids.join_eui
      delete device.ids.dev_eui
      delete device.root_keys
      delete device.resets_join_nonces
    }
    delete device.activation_mode

    return api.device.create(appId, device, {
      abp: values.activation_mode === 'abp',
      withRootKeys: true,
    })
  }

  handleSubmitSuccess (device) {
    const { appId, redirectToList } = this.props
    const deviceId = getDeviceId(device)

    redirectToList(appId, deviceId)
  }

  render () {
    const { error } = this.state

    return (
      <Container>
        <Row className={style.wrapper}>
          <Col sm={12}>
            <IntlHelmet title={sharedMessages.addDevice} />
            <Message component="h2" content={sharedMessages.addDevice} />
          </Col>
          <Col className={style.form} sm={12} md={12} lg={8} xl={8}>
            <DeviceDataForm
              error={error}
              onSubmit={this.handleSubmit}
              onSubmitSuccess={this.handleSubmitSuccess}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
