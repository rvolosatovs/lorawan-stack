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
import { defineMessages } from 'react-intl'

import Tag from '@ttn-lw/components/tag'
import TagGroup from '@ttn-lw/components/tag/group'

import FetchTable from '@ttn-lw/containers/fetch-table'

import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import style from './api-keys-table.styl'

const m = defineMessages({
  keyId: 'Key ID',
  grantedRights: 'Granted Rights',
})

const formatRight = function(right) {
  return right
    .split('_')
    .slice(1)
    .map(r => r.charAt(0) + r.slice(1).toLowerCase())
    .join(' ')
}

const RIGHT_TAG_MAX_WIDTH = 140

const headers = [
  {
    name: 'id',
    displayName: m.keyId,
    width: 30,
    render(id) {
      return <span className={style.keyId}>{id}</span>
    },
  },
  {
    name: 'name',
    displayName: sharedMessages.name,
    width: 30,
  },
  {
    name: 'rights',
    displayName: m.grantedRights,
    width: 40,
    render(rights) {
      const tags = rights.map(r => (
        <Tag className={style.rightTag} content={formatRight(r)} key={r} />
      ))

      return (
        <TagGroup className={style.rightTagGroup} tagMaxWidth={RIGHT_TAG_MAX_WIDTH} tags={tags} />
      )
    },
  },
]

export default class ApiKeysTable extends Component {
  static propTypes = {
    baseDataSelector: PropTypes.func.isRequired,
    getItemsAction: PropTypes.func.isRequired,
    pageSize: PropTypes.number.isRequired,
  }

  render() {
    const { pageSize, baseDataSelector, getItemsAction } = this.props

    return (
      <FetchTable
        entity="keys"
        headers={headers}
        addMessage={sharedMessages.addApiKey}
        pageSize={pageSize}
        baseDataSelector={baseDataSelector}
        getItemsAction={getItemsAction}
        tableTitle={<Message content={sharedMessages.apiKeys} />}
      />
    )
  }
}
