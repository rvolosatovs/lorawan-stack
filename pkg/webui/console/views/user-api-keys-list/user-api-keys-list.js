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
import { Container, Row, Col } from 'react-grid-system'

import PAGE_SIZES from '@ttn-lw/constants/page-sizes'

import IntlHelmet from '@ttn-lw/lib/components/intl-helmet'

import ApiKeysTable from '@console/containers/api-keys-table'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import { getApiKeysList } from '@console/store/actions/api-keys'

import {
  selectApiKeys,
  selectApiKeysTotalCount,
  selectApiKeysFetching,
  selectApiKeysError,
} from '@console/store/selectors/api-keys'

const UserApiKeysList = props => {
  const { userId } = props

  const baseDataSelector = React.useCallback(
    state => {
      const id = { id: userId }

      return {
        keys: selectApiKeys(state, id),
        totalCount: selectApiKeysTotalCount(state, id),
        fetching: selectApiKeysFetching(state),
        error: selectApiKeysError(state),
      }
    },
    [userId],
  )

  const getApiKeys = React.useCallback(filters => getApiKeysList('users', userId, filters), [
    userId,
  ])

  return (
    <Container>
      <Row>
        <IntlHelmet title={sharedMessages.personalApiKeys} />
        <Col>
          <ApiKeysTable
            pageSize={PAGE_SIZES.REGULAR}
            baseDataSelector={baseDataSelector}
            getItemsAction={getApiKeys}
          />
        </Col>
      </Row>
    </Container>
  )
}

UserApiKeysList.propTypes = {
  userId: PropTypes.string.isRequired,
}

export default UserApiKeysList
