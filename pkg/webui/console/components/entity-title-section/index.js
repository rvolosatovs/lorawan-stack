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
import { Container, Row, Col } from 'react-grid-system'
import classnames from 'classnames'

import DateTime from '@ttn-lw/lib/components/date-time'
import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import style from './entity-title-section.styl'

const EntityTitleSection = ({ entityName, entityId, description, creationDate, children }) => {
  return (
    <React.Fragment>
      <Container>
        <Row>
          <Col md={12} className={style.container}>
            <h1 className={style.title}>{entityName || entityId}</h1>
            <span className={style.id}>
              <strong>
                <Message content={sharedMessages.id} />:
              </strong>{' '}
              {entityId}
            </span>
            {description && <span className={style.description}>{description}</span>}
            <div className={style.bottom}>
              <div className={style.children}>{children}</div>
              <div className={style.creationDate}>
                <span>
                  <Message content={sharedMessages.created} />{' '}
                  <DateTime.Relative value={creationDate} />
                </span>
              </div>
            </div>
          </Col>
        </Row>
      </Container>
      <hr className={style.hRule} />
    </React.Fragment>
  )
}

EntityTitleSection.propTypes = {
  children: PropTypes.node.isRequired,
  creationDate: PropTypes.string.isRequired,
  description: PropTypes.string,
  entityId: PropTypes.string.isRequired,
  entityName: PropTypes.string,
}

EntityTitleSection.defaultProps = {
  entityName: undefined,
  description: undefined,
}

EntityTitleSection.Device = ({ deviceName, deviceId, description, children }) => {
  return (
    <Container>
      <Row>
        <Col>
          <div
            className={classnames(style.containerDevice, {
              [style.hasDescription]: Boolean(description),
            })}
          >
            <h1 className={style.title}>{deviceName || deviceId}</h1>
            <span className={style.id}>
              <strong>
                <Message content={sharedMessages.id} />:
              </strong>{' '}
              {deviceId}
            </span>
            {description && <span className={style.description}>{description}</span>}
          </div>
          {children}
        </Col>
      </Row>
    </Container>
  )
}

EntityTitleSection.Device.propTypes = {
  children: PropTypes.node.isRequired,
  description: PropTypes.string,
  deviceId: PropTypes.string.isRequired,
  deviceName: PropTypes.string,
}

EntityTitleSection.Device.defaultProps = {
  deviceName: undefined,
  description: undefined,
}

export default EntityTitleSection
