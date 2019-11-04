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
import bind from 'autobind-decorator'

import Form from '../../../components/form'
import SubmitButton from '../../../components/submit-button'
import Input from '../../../components/input'
import Checkbox from '../../../components/checkbox'
import Radio from '../../../components/radio-button'
import Select from '../../../components/select'
import toast from '../../../components/toast'
import Message from '../../../lib/components/message'
import SubmitBar from '../../../components/submit-bar'
import ModalButton from '../../../components/button/modal-button'
import { NsFrequencyPlansSelect } from '../../containers/freq-plans-select'
import DevAddrInput from '../../containers/dev-addr-input'
import JoinEUIPrefixesInput from '../../containers/join-eui-prefixes-input'

import sharedMessages from '../../../lib/shared-messages'
import { selectNsConfig, selectJsConfig, selectAsConfig } from '../../../lib/selectors/env'
import errorMessages from '../../../lib/errors/error-messages'
import { getDeviceId } from '../../../lib/selectors/id'
import PropTypes from '../../../lib/prop-types'
import m from './messages'
import validationSchema from './validation-schema'

class DeviceDataForm extends Component {
  constructor(props) {
    super(props)

    const {
      initialValues: { supports_join, update, root_keys = {} },
    } = this.props

    const external_js = !Boolean(Object.keys(root_keys).length)
    let otaa = true
    if (update) {
      otaa = Boolean(supports_join)
    }

    this.state = {
      otaa,
      resets_join_nonces: false,
      resets_f_cnt: false,
      external_js,
    }
    this.formRef = React.createRef()
  }

  @bind
  handleOTAASelect() {
    this.setState({ otaa: true })
  }

  @bind
  handleABPSelect() {
    this.setState({ otaa: false })
  }

  @bind
  handleResetsJoinNoncesChange(evt) {
    this.setState({ resets_join_nonces: evt.target.checked })
  }

  @bind
  handleResetsFrameCountersChange(evt) {
    this.setState({ resets_f_cnt: evt.target.checked })
  }

  @bind
  async handleExternalJoinServerChange(evt) {
    const external_js = evt.target.checked
    await this.setState(({ resets_join_nonces }) => ({
      external_js,
      resets_join_nonces: external_js ? false : resets_join_nonces,
    }))

    const { initialValues } = this.props
    const jsConfig = selectJsConfig()
    const { setValues, state } = this.formRef.current

    // Reset Join Server related entries if the device is provisined by
    // an external JS.
    if (external_js) {
      setValues({
        ...state.values,
        root_keys: {
          nwk_key: {},
          app_key: {},
        },
        resets_join_nonces: false,
        join_server_address: undefined,
        external_js,
      })
    } else {
      let { join_server_address } = initialValues

      // Reset `join_server_address` if is present after disabling external JS provisioning.
      if (jsConfig.enabled && !Boolean(join_server_address)) {
        join_server_address = new URL(jsConfig.base_url).hostname
      }

      setValues({
        ...state.values,
        join_server_address,
        external_js,
      })
    }
  }

  @bind
  async handleSubmit(values, { setSubmitting, resetForm }) {
    const { onSubmit, onSubmitSuccess, initialValues, update } = this.props
    const deviceId = getDeviceId(initialValues)
    const { external_js, ...castedValues } = validationSchema.cast(values)
    await this.setState({ error: '' })

    try {
      const device = await onSubmit(castedValues)
      if (update) {
        resetForm(values)
        toast({
          title: deviceId,
          message: m.updateSuccess,
          type: toast.types.SUCCESS,
        })
      }
      await onSubmitSuccess(device)
    } catch (error) {
      setSubmitting(false)
      const err = error instanceof Error ? errorMessages.genericError : error
      await this.setState({ error: err })
    }
  }

  @bind
  async handleDelete() {
    const { onDelete, onDeleteSuccess, initialValues } = this.props
    const deviceId = getDeviceId(initialValues)

    try {
      await onDelete()
      toast({
        title: deviceId,
        message: m.deleteSuccess,
        type: toast.types.SUCCESS,
      })
      onDeleteSuccess()
    } catch (error) {
      const err = error instanceof Error ? errorMessages.genericError : error
      this.setState({ error: err })
    }
  }

  get ABPSection() {
    const { resets_f_cnt } = this.state
    return (
      <React.Fragment>
        <DevAddrInput
          title={sharedMessages.devAddr}
          name="session.dev_addr"
          placeholder={m.leaveBlankPlaceholder}
          description={m.deviceAddrDescription}
          required
        />
        <Form.Field
          title={sharedMessages.fwdNtwkKey}
          name="session.keys.f_nwk_s_int_key.key"
          type="byte"
          min={16}
          max={16}
          placeholder={m.leaveBlankPlaceholder}
          description={m.fwdNtwkKeyDescription}
          component={Input}
        />
        <Form.Field
          title={sharedMessages.sNtwkSIKey}
          name="session.keys.s_nwk_s_int_key.key"
          type="byte"
          min={16}
          max={16}
          placeholder={m.leaveBlankPlaceholder}
          description={m.sNtwkSIKeyDescription}
          component={Input}
        />
        <Form.Field
          title={sharedMessages.ntwkSEncKey}
          name="session.keys.nwk_s_enc_key.key"
          type="byte"
          min={16}
          max={16}
          placeholder={m.leaveBlankPlaceholder}
          description={m.ntwkSEncKeyDescription}
          component={Input}
        />
        <Form.Field
          title={sharedMessages.appSKey}
          name="session.keys.app_s_key.key"
          type="byte"
          min={16}
          max={16}
          placeholder={m.leaveBlankPlaceholder}
          description={m.appSKeyDescription}
          component={Input}
        />
        <Form.Field
          title={m.resetsFCnt}
          onChange={this.handleResetsFrameCountersChange}
          warning={resets_f_cnt ? m.resetWarning : undefined}
          name="mac_settings.resets_f_cnt"
          component={Checkbox}
        />
      </React.Fragment>
    )
  }

  get OTAASection() {
    const { resets_join_nonces, external_js } = this.state
    const {
      update,
      initialValues: { root_keys },
    } = this.props

    const rootKeysNotExposed =
      root_keys && root_keys.root_key_id && !root_keys.nwk_key && !root_keys.app_key

    let rootKeyPlaceholder = m.leaveBlankPlaceholder
    if (external_js) {
      rootKeyPlaceholder = sharedMessages.provisionedOnExternalJoinServer
    } else if (rootKeysNotExposed) {
      rootKeyPlaceholder = m.unexposed
    }

    return (
      <React.Fragment>
        <JoinEUIPrefixesInput
          title={sharedMessages.joinEUI}
          name="ids.join_eui"
          description={m.joinEUIDescription}
          required
          disabled={update}
          showPrefixes={!update}
        />
        <Form.Field
          title={sharedMessages.devEUI}
          name="ids.dev_eui"
          type="byte"
          min={8}
          max={8}
          description={m.deviceEUIDescription}
          required
          disabled={update}
          component={Input}
        />
        <Form.Field
          title={m.externalJoinServer}
          description={m.externalJoinServerDescription}
          name="external_js"
          onChange={this.handleExternalJoinServerChange}
          component={Checkbox}
        />
        <Form.Field
          title={sharedMessages.joinServerAddress}
          placeholder={external_js ? m.external : sharedMessages.addressPlaceholder}
          name="join_server_address"
          component={Input}
          disabled={external_js}
        />
        <Form.Field
          title={sharedMessages.appKey}
          name="root_keys.app_key.key"
          type="byte"
          min={16}
          max={16}
          placeholder={rootKeyPlaceholder}
          description={m.appKeyDescription}
          component={Input}
          disabled={external_js || rootKeysNotExposed}
        />
        <Form.Field
          title={sharedMessages.nwkKey}
          name="root_keys.nwk_key.key"
          type="byte"
          min={16}
          max={16}
          placeholder={rootKeyPlaceholder}
          description={m.nwkKeyDescription}
          component={Input}
          disabled={external_js || rootKeysNotExposed}
        />
        <Form.Field
          title={m.resetsJoinNonces}
          onChange={this.handleResetsJoinNoncesChange}
          warning={resets_join_nonces ? m.resetWarning : undefined}
          name="resets_join_nonces"
          component={Checkbox}
          disabled={external_js}
        />
      </React.Fragment>
    )
  }

  render() {
    const { otaa, error, external_js } = this.state
    const { initialValues, update } = this.props

    let deviceId
    let deviceName

    if (initialValues) {
      deviceId = getDeviceId(initialValues)
      deviceName = initialValues.name
    }

    const emptyValues = {
      ids: {
        device_id: undefined,
        join_eui: undefined,
        dev_eui: undefined,
      },
      activation_mode: 'otaa',
      lorawan_version: undefined,
      lorawan_phy_version: undefined,
      frequency_plan_id: undefined,
      supports_class_c: false,
      resets_join_nonces: false,
      root_keys: {
        nwk_key: { key: undefined },
        app_key: { key: undefined },
      },
      session: {
        dev_addr: undefined,
        keys: {
          f_nwk_s_int_key: {},
          s_nwk_s_int_key: {},
          nwk_s_enc_key: {},
          app_s_key: {},
        },
      },
      mac_settings: {
        resets_f_cnt: false,
      },
    }

    const nsConfig = selectNsConfig()
    const asConfig = selectAsConfig()
    const jsConfig = selectJsConfig()
    const joinServerAddress = jsConfig.enabled ? new URL(jsConfig.base_url).hostname : ''

    const formValues = {
      network_server_address: nsConfig.enabled ? new URL(nsConfig.base_url).hostname : '',
      application_server_address: asConfig.enabled ? new URL(asConfig.base_url).hostname : '',
      join_server_address: external_js ? undefined : joinServerAddress,
      ...emptyValues,
      ...initialValues,
      activation_mode: otaa ? 'otaa' : 'abp',
      external_js,
    }

    return (
      <Form
        error={error}
        onSubmit={this.handleSubmit}
        validationSchema={validationSchema}
        submitEnabledWhenInvalid
        initialValues={formValues}
        formikRef={this.formRef}
      >
        <Message component="h4" content={sharedMessages.generalSettings} />
        <Form.Field
          title={sharedMessages.devID}
          name="ids.device_id"
          placeholder={m.deviceIdPlaceholder}
          description={m.deviceIdDescription}
          autoFocus
          required
          disabled={update}
          component={Input}
        />
        <Form.Field
          title={sharedMessages.devName}
          name="name"
          placeholder={m.deviceNamePlaceholder}
          description={m.deviceNameDescription}
          component={Input}
        />
        <Form.Field
          title={sharedMessages.devDesc}
          name="description"
          type="textarea"
          description={m.deviceDescDescription}
          component={Input}
        />
        <Message component="h4" content={m.lorawanOptions} />
        <Form.Field
          title={sharedMessages.macVersion}
          name="lorawan_version"
          component={Select}
          required
          options={[
            { value: '1.0.0', label: 'MAC V1.0' },
            { value: '1.0.1', label: 'MAC V1.0.1' },
            { value: '1.0.2', label: 'MAC V1.0.2' },
            { value: '1.0.3', label: 'MAC V1.0.3' },
            { value: '1.1.0', label: 'MAC V1.1' },
          ]}
        />
        <Form.Field
          title={sharedMessages.phyVersion}
          name="lorawan_phy_version"
          component={Select}
          required
          options={[
            { value: '1.0.0', label: 'PHY V1.0' },
            { value: '1.0.1', label: 'PHY V1.0.1' },
            { value: '1.0.2-a', label: 'PHY V1.0.2 REV A' },
            { value: '1.0.2-b', label: 'PHY V1.0.2 REV B' },
            { value: '1.0.3-a', label: 'PHY V1.0.3 REV A' },
            { value: '1.1.0-a', label: 'PHY V1.1 REV A' },
            { value: '1.1.0-b', label: 'PHY V1.1 REV B' },
          ]}
        />
        <NsFrequencyPlansSelect name="frequency_plan_id" required />
        <Form.Field title={m.supportsClassC} name="supports_class_c" component={Checkbox} />
        <Form.Field
          title={sharedMessages.networkServerAddress}
          placeholder={sharedMessages.addressPlaceholder}
          name="network_server_address"
          component={Input}
        />
        <Form.Field
          title={sharedMessages.applicationServerAddress}
          placeholder={sharedMessages.addressPlaceholder}
          name="application_server_address"
          component={Input}
        />
        <Message component="h4" content={m.activationSettings} />
        <Form.Field
          title={m.activationMode}
          disabled={update}
          name="activation_mode"
          component={Radio.Group}
        >
          <Radio label={m.otaa} value="otaa" onChange={this.handleOTAASelect} />
          <Radio label={m.abp} value="abp" onChange={this.handleABPSelect} />
        </Form.Field>
        {otaa ? this.OTAASection : this.ABPSection}
        <SubmitBar>
          <Form.Submit
            component={SubmitButton}
            message={update ? sharedMessages.saveChanges : m.createDevice}
          />
          {update && (
            <ModalButton
              type="button"
              icon="delete"
              message={m.deleteDevice}
              modalData={{
                message: { values: { deviceId: deviceName || deviceId }, ...m.deleteWarning },
              }}
              onApprove={this.handleDelete}
              danger
              naked
            />
          )}
        </SubmitBar>
      </Form>
    )
  }
}

const keyPropType = PropTypes.shape({
  key: PropTypes.string,
})

const initialValuesPropType = PropTypes.shape({
  ids: PropTypes.shape({
    device_id: PropTypes.string,
    join_eui: PropTypes.string,
    dev_eui: PropTypes.string,
  }),
  name: PropTypes.string,
  activation_mode: PropTypes.string,
  lorawan_version: PropTypes.string,
  lorawan_phy_version: PropTypes.string,
  frequency_plan_id: PropTypes.string,
  supports_class_c: PropTypes.bool,
  resets_join_nonces: PropTypes.bool,
  root_keys: PropTypes.shape({
    nwk_key: keyPropType,
    app_key: keyPropType,
    root_key_id: PropTypes.string,
  }),
  session: PropTypes.shape({
    dev_addr: PropTypes.string,
    keys: PropTypes.shape({
      f_nwk_s_int_key: keyPropType,
      s_nwk_s_int_key: keyPropType,
      nwk_s_enc_key: keyPropType,
      app_s_key: keyPropType,
    }),
  }),
  mac_settings: PropTypes.shape({
    resets_f_cnt: PropTypes.bool,
  }),
})

DeviceDataForm.propTypes = {
  initialValues: initialValuesPropType,
  onDelete: PropTypes.func,
  onDeleteSuccess: PropTypes.func,
  onSubmit: PropTypes.func.isRequired,
  onSubmitSuccess: PropTypes.func,
  update: PropTypes.bool,
}

DeviceDataForm.defaultProps = {
  onDelete: () => null,
  onDeleteSuccess: () => null,
  onSubmitSuccess: () => null,
  initialValues: {},
  update: false,
}

export default DeviceDataForm
