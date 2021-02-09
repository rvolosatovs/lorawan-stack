// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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
import { defineMessages, useIntl } from 'react-intl'

import devicePlaceholder from '@assets/misc/end-device-placeholder.svg'

import Link from '@ttn-lw/components/link'

import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import { LORAWAN_VERSIONS, LORAWAN_PHY_VERSIONS } from '@console/lib/device-utils'

import style from './device-card.styl'

const m = defineMessages({
  productWebsite: 'Product website',
  dataSheet: 'Data sheet',
  classA: 'Class A',
  classB: 'Class B',
  classC: 'Class C',
})

const getLorawanVersionLabel = ({ lorawan_version }) => {
  const { label } = LORAWAN_VERSIONS.find(version => version.value === lorawan_version) || {}

  return label
}

const getLorawanPhyVersionLabel = ({ lorawan_phy_version }) => {
  const { label } =
    LORAWAN_PHY_VERSIONS.find(version => version.value === lorawan_phy_version) || {}

  return label
}

const DeviceCard = props => {
  const { model, template } = props
  const { name, description, product_url, datasheet_url, photos = {} } = model
  const { end_device: device } = template
  const { formatMessage } = useIntl()

  const deviceImage = photos.main || devicePlaceholder
  const lwVersionLabel = getLorawanVersionLabel(device)
  const lwPhyVersionLabel = getLorawanPhyVersionLabel(device)
  const modeTitleLabel = device.supports_join
    ? sharedMessages.otaa
    : device.multicast
    ? sharedMessages.multicast
    : sharedMessages.abp
  const deviceClassTitleLabel = device.supports_class_c
    ? m.classC
    : device.supports_class_b
    ? m.classB
    : m.classA

  return (
    <div className={style.container}>
      <img className={style.image} src={deviceImage} name={name} />
      <div className={style.info}>
        <div>
          <h3 className={style.name}>{name}</h3>
          {Boolean(lwVersionLabel) && (
            <span className={style.infoItem} title={formatMessage(sharedMessages.macVersion)}>
              {lwVersionLabel}
            </span>
          )}
          {Boolean(lwPhyVersionLabel) && (
            <span className={style.infoItem} title={formatMessage(sharedMessages.phyVersion)}>
              {lwPhyVersionLabel}
            </span>
          )}
          <Message
            className={style.infoItem}
            content={modeTitleLabel}
            title={formatMessage(sharedMessages.activationMode)}
            component="span"
          />
          <Message
            className={style.infoItem}
            content={deviceClassTitleLabel}
            title={formatMessage(sharedMessages.activationMode)}
            component="span"
          />
        </div>
        {description && <p>{description}</p>}
        <div>
          {product_url && (
            <Link.Anchor href={product_url} external>
              <Message content={m.productWebsite} />
            </Link.Anchor>
          )}
          {datasheet_url && (
            <Link.Anchor href={datasheet_url} external>
              <Message content={m.dataSheet} />
            </Link.Anchor>
          )}
        </div>
      </div>
    </div>
  )
}

DeviceCard.propTypes = {
  model: PropTypes.shape({
    name: PropTypes.string.isRequired,
    description: PropTypes.string,
    product_url: PropTypes.string,
    datasheet_url: PropTypes.string,
    photos: PropTypes.shape({
      main: PropTypes.string,
    }),
  }).isRequired,
  template: PropTypes.deviceTemplate.isRequired,
}

export default DeviceCard
