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

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import DescriptionList from './shared/description-list'

const ApplicationLocationPreview = React.memo(({ event }) => {
  const { data } = event
  const { latitude, longitude, altitude, accuracy, source } = data.location

  return (
    <DescriptionList>
      <DescriptionList.Item title={sharedMessages.latitude} data={latitude} />
      <DescriptionList.Item title={sharedMessages.longitude} data={longitude} />
      <DescriptionList.Item title={sharedMessages.altitude} data={altitude} />
      <DescriptionList.Item title={sharedMessages.accuracy} data={accuracy} />
      <DescriptionList.Item title={sharedMessages.source} data={source.replace(/^(SOURCE_)/, '')} />
    </DescriptionList>
  )
})

ApplicationLocationPreview.propTypes = {
  event: PropTypes.event.isRequired,
}

export default ApplicationLocationPreview
