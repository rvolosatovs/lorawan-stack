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
import { Col, Row, Container } from 'react-grid-system'

import sharedMessages from '../../../lib/shared-messages'
import DateTime from '../../../lib/components/date-time'
import DevicesTable from '../../containers/devices-table'
import DataSheet from '../../../components/data-sheet'


import style from './application-overview.styl'

const DEVICES_TABLE_SIZE = 5

@connect(function ({ application }, props) {
  return {
    application: application.application,
  }
})
class ApplicationOverview extends React.Component {

  get applicationInfo () {
    const {
      ids,
      name,
      description,
      created_at,
      updated_at,
    } = this.props.application

    const sheetData = [
      {
        header: sharedMessages.generalInformation,
        items: [
          { key: sharedMessages.appId, value: ids.application_id, type: 'code', sensitive: false },
          { key: sharedMessages.createdAt, value: <DateTime value={created_at} /> },
          { key: sharedMessages.updatedAt, value: <DateTime value={updated_at} /> },
        ],
      },
    ]

    return (
      <div>
        <div className={style.title}>
          <h2>{name || ids.application_id}</h2>
          { description && <span className={style.description}>{description}</span> }
        </div>
        <DataSheet data={sheetData} />
      </div>
    )
  }

  render () {
    return (
      <Container>
        <Row>
          <Col sm={12} lg={6}>
            {this.applicationInfo}
          </Col>
          <Col sm={12} lg={6}>
            <div className={style.latestEvents}>
              <h4>Latest Data</h4>
              <div>Activity Events Placeholder</div>
            </div>
          </Col>
        </Row>
        <Row>
          <Col sm={12} className={style.table}>
            <DevicesTable pageSize={DEVICES_TABLE_SIZE} />
          </Col>
        </Row>
      </Container>
    )
  }
}

export default ApplicationOverview
