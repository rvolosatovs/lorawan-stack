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
import { defineMessages } from 'react-intl'
import bind from 'autobind-decorator'

import FetchTable from '../fetch-table'
import Message from '../../../lib/components/message'

import sharedMessages from '../../../lib/shared-messages'
import PropTypes from '../../../lib/prop-types'

import { getWebhooksList } from '../../../console/store/actions/webhooks'
import {
  selectWebhooks,
  selectWebhooksTotalCount,
  selectWebhooksFetching,
} from '../../../console/store/selectors/webhooks'

const m = defineMessages({
  templateId: 'Template ID',
  format: 'Format',
  baseUrl: 'Base URL',
})

const headers = [
  {
    name: 'ids.webhook_id',
    displayName: sharedMessages.id,
    width: 40,
  },
  {
    name: 'template_ids.template_id',
    displayName: m.templateId,
    width: 15,
    render: value => value || <i>none</i>,
  },
  {
    name: 'format',
    displayName: m.format,
    width: 15,
  },
  {
    name: 'base_url',
    displayName: m.baseUrl,
    width: 30,
  },
]

@bind
export default class WebhooksTable extends React.Component {
  static propTypes = {
    appId: PropTypes.string.isRequired,
  }

  constructor(props) {
    super(props)

    const { appId } = props
    this.getWebhooksList = () => getWebhooksList(appId, ['template_ids'])
  }

  baseDataSelector(state) {
    return {
      webhooks: selectWebhooks(state),
      totalCount: selectWebhooksTotalCount(state),
      fetching: selectWebhooksFetching(state),
    }
  }

  render() {
    return (
      <FetchTable
        entity="webhooks"
        addMessage={sharedMessages.addWebhook}
        headers={headers}
        getItemsAction={this.getWebhooksList}
        baseDataSelector={this.baseDataSelector}
        tableTitle={<Message content={sharedMessages.webhooks} />}
        {...this.props}
      />
    )
  }
}
