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
import bind from 'autobind-decorator'
import { connect } from 'react-redux'

import toast from '../../../components/toast'
import PropTypes from '../../../lib/prop-types'
import sharedMessages from '../../../lib/shared-messages'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import PayloadFormattersForm from '../../components/payload-formatters-form'
import IntlHelmet from '../../../lib/components/intl-helmet'
import PAYLOAD_FORMATTER_TYPES from '../../constants/formatter-types'
import Message from '../../../lib/components/message'
import { updateApplicationLinkSuccess } from '../../store/actions/link'
import {
  selectApplicationIsLinked,
  selectApplicationLinkFormatters,
  selectSelectedApplicationId,
} from '../../store/selectors/applications'

import api from '../../api'

@connect(
  function(state) {
    const formatters = selectApplicationLinkFormatters(state) || {}

    return {
      appId: selectSelectedApplicationId(state),
      linked: selectApplicationIsLinked(state) || false,
      formatters,
    }
  },
  { updateLinkSuccess: updateApplicationLinkSuccess },
)
@withBreadcrumb('apps.single.payload-formatters.uplink', function(props) {
  const { appId } = props

  return (
    <Breadcrumb
      path={`/applications/${appId}/payload-formatters/uplink`}
      icon="uplink"
      content={sharedMessages.uplink}
    />
  )
})
@bind
class ApplicationPayloadFormatters extends React.PureComponent {
  static propTypes = {
    appId: PropTypes.string.isRequired,
    formatters: PropTypes.object.isRequired,
    updateLinkSuccess: PropTypes.func.isRequired,
    linked: PropTypes.bool.isRequired,
  }

  async onSubmit(values) {
    const { appId, formatters } = this.props

    return await api.application.link.set(appId, {
      default_formatters: {
        down_formatter: formatters.down_formatter || PAYLOAD_FORMATTER_TYPES.NONE,
        down_formatter_parameter: formatters.down_formatter_parameter || '',
        up_formatter: values.type,
        up_formatter_parameter: values.parameter,
      },
    })
  }

  onSubmitSuccess(link) {
    const { appId, updateLinkSuccess } = this.props
    toast({
      title: appId,
      message: sharedMessages.payloadFormattersUpdateSuccess,
      type: toast.types.SUCCESS,
    })
    updateLinkSuccess(link)
  }

  render() {
    const { formatters, linked } = this.props

    return (
      <React.Fragment>
        <IntlHelmet title={sharedMessages.payloadFormattersUplink} />
        <Message component="h2" content={sharedMessages.payloadFormattersUplink} />
        <PayloadFormattersForm
          uplink
          linked={linked}
          onSubmit={this.onSubmit}
          onSubmitSuccess={this.onSubmitSuccess}
          initialType={formatters.up_formatter || PAYLOAD_FORMATTER_TYPES.NONE}
          initialParameter={formatters.up_formatter_parameter || ''}
        />
      </React.Fragment>
    )
  }
}

export default ApplicationPayloadFormatters
