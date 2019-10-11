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

import Message from '../../../lib/components/message'
import IntlHelmet from '../../../lib/components/intl-helmet'
import DeviceDataForm from '../../components/device-data-form'
import sharedMessages from '../../../lib/shared-messages'
import Button from '../../../components/button'
import withRequest from '../../../lib/components/with-request'
import { getDeviceTemplateFormats } from '../../store/actions/device-template-formats'
import { selectSelectedApplicationId } from '../../store/selectors/applications'
import { selectDeviceTemplateFormats } from '../../store/selectors/device-template-formats'
import { getDeviceId } from '../../../lib/selectors/id'
import { selectNsConfig, selectJsConfig, selectAsConfig } from '../../../lib/selectors/env'
import PropTypes from '../../../lib/prop-types'
import api from '../../api'
import style from './device-add-single.styl'

@connect(
  state => ({
    appId: selectSelectedApplicationId(state),
    deviceTemplateFormats: selectDeviceTemplateFormats(state),
    asConfig: selectAsConfig(),
    nsConfig: selectNsConfig(),
    jsConfig: selectJsConfig(),
  }),
  dispatch => ({
    redirectToList: (appId, deviceId) =>
      dispatch(push(`/applications/${appId}/devices/${deviceId}`)),
    getDeviceTemplateFormats: () => dispatch(getDeviceTemplateFormats()),
  }),
)
@withRequest(({ getDeviceTemplateFormats }) => getDeviceTemplateFormats())
export default class DeviceAdd extends Component {
  static propTypes = {
    appId: PropTypes.string.isRequired,
    asConfig: PropTypes.stackComponent.isRequired,
    deviceTemplateFormats: PropTypes.shape({}).isRequired,
    jsConfig: PropTypes.stackComponent.isRequired,
    nsConfig: PropTypes.stackComponent.isRequired,
    redirectToList: PropTypes.func.isRequired,
  }

  state = {
    error: '',
  }

  @bind
  async handleSubmit(values) {
    const { appId } = this.props
    const { activation_mode, ...device } = values

    return api.device.create(appId, device, {
      abp: values.activation_mode === 'abp',
    })
  }

  @bind
  handleSubmitSuccess(device) {
    const { appId, redirectToList } = this.props
    const deviceId = getDeviceId(device)

    redirectToList(appId, deviceId)
  }

  render() {
    const { error } = this.state
    const { asConfig, nsConfig, jsConfig, deviceTemplateFormats } = this.props
    const canBulkCreate = Object.keys(deviceTemplateFormats).length !== 0

    const initialValues = {
      network_server_address: nsConfig.enabled ? new URL(nsConfig.base_url).hostname : '',
      application_server_address: asConfig.enabled ? new URL(asConfig.base_url).hostname : '',
      join_server_address: jsConfig.enabled ? new URL(jsConfig.base_url).hostname : '',
    }

    return (
      <Container>
        <Row>
          <Col sm={6}>
            <IntlHelmet title={sharedMessages.addDevice} />
            <Message className={style.title} component="h2" content={sharedMessages.addDevice} />
          </Col>
          <Col className={style.bulkCreation} sm={6}>
            <Button.Link
              message={sharedMessages.bulkCreation}
              icon="bulk_creation"
              to="add/bulk"
              disabled={!canBulkCreate}
            />
          </Col>
        </Row>
        <Row>
          <Col className={style.form} lg={8} md={12}>
            <DeviceDataForm
              error={error}
              onSubmit={this.handleSubmit}
              onSubmitSuccess={this.handleSubmitSuccess}
              initialValues={initialValues}
              jsConfig={jsConfig}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
