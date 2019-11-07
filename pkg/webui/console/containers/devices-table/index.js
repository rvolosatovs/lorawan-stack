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
import bind from 'autobind-decorator'

import sharedMessages from '../../../lib/shared-messages'
import Message from '../../../lib/components/message'
import PropTypes from '../../../lib/prop-types'
import FetchTable from '../fetch-table'
import DateTime from '../../../lib/components/date-time'
import Button from '../../../components/button'
import withRequest from '../../../lib/components/with-request'
import withFeatureRequirement from '../../lib/components/with-feature-requirement'

import { getDevicesList } from '../../../console/store/actions/devices'
import { getDeviceTemplateFormats } from '../../store/actions/device-template-formats'
import { selectSelectedApplicationId } from '../../store/selectors/applications'
import { selectDeviceTemplateFormats } from '../../store/selectors/device-template-formats'
import {
  checkFromState,
  mayCreateOrEditApplicationDevices,
  mayViewApplicationDevices,
} from '../../lib/feature-checks'

import style from './devices-table.styl'

const headers = [
  {
    name: 'ids.device_id',
    displayName: sharedMessages.id,
  },
  {
    name: 'name',
    displayName: sharedMessages.name,
  },
  {
    name: 'created_at',
    displayName: sharedMessages.created,
    render(datetime) {
      return <DateTime.Relative value={datetime} />
    },
  },
]

@connect(
  function(state) {
    return {
      appId: selectSelectedApplicationId(state),
      deviceTemplateFormats: selectDeviceTemplateFormats(state),
      mayCreateDevices: checkFromState(mayCreateOrEditApplicationDevices, state),
    }
  },
  { getDeviceTemplateFormats },
)
@withFeatureRequirement(mayViewApplicationDevices)
@withRequest(({ getDeviceTemplateFormats }) => getDeviceTemplateFormats())
@bind
class DevicesTable extends React.Component {
  static propTypes = {
    appId: PropTypes.string.isRequired,
    devicePathPrefix: PropTypes.string,
    deviceTemplateFormats: PropTypes.shape({}).isRequired,
    mayCreateDevices: PropTypes.bool.isRequired,
    totalCount: PropTypes.number,
  }

  static defaultProps = {
    devicePathPrefix: undefined,
    totalCount: 0,
  }

  constructor(props) {
    super(props)

    this.getDevicesList = filters => getDevicesList(props.appId, filters, ['name'])
  }

  baseDataSelector(state) {
    const { mayCreateDevices } = this.props
    return {
      ...state.devices,
      mayAdd: mayCreateDevices,
    }
  }

  get importButton() {
    const { deviceTemplateFormats, mayCreateDevices, appId } = this.props
    const canBulkCreate = Object.keys(deviceTemplateFormats).length !== 0

    return (
      mayCreateDevices && (
        <Button.Link
          className={style.importDevices}
          message={sharedMessages.importDevices}
          icon="import_devices"
          to={`/applications/${appId}/devices/import`}
          disabled={!canBulkCreate}
        />
      )
    )
  }

  render() {
    const { devicePathPrefix } = this.props
    return (
      <FetchTable
        entity="devices"
        headers={headers}
        addMessage={sharedMessages.addDevice}
        actionItems={this.importButton}
        tableTitle={<Message content={sharedMessages.devices} />}
        getItemsAction={this.getDevicesList}
        searchItemsAction={this.getDevicesList}
        itemPathPrefix={devicePathPrefix}
        baseDataSelector={this.baseDataSelector}
        {...this.props}
      />
    )
  }
}

export default DevicesTable
