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
import { FormattedDate, FormattedTime } from 'react-intl'

import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'
import { warn } from '@ttn-lw/lib/log'

import RelativeTime from './relative'

class DateTime extends React.PureComponent {
  renderDateTime(formattedDate, formattedTime, dateValue) {
    const { className, children, date, time } = this.props

    let result = ''
    if (date) {
      result += formattedDate
    }

    if (time) {
      if (date) {
        result += ' '
      }

      result += formattedTime
    }

    if (isNaN(dateValue)) {
      warn('Invalid date passed to DateTime component')
      return (
        <time className={className}>
          <Message content={sharedMessages.unknown} firstToLower />
        </time>
      )
    }

    return (
      <time className={className} dateTime={dateValue.toISOString()} title={result}>
        {children ? children(result) : result}
      </time>
    )
  }

  render() {
    const { value, dateFormatOptions, timeFormatOptions } = this.props

    let dateValue = value
    if (!(value instanceof Date)) {
      dateValue = new Date(value)
    }

    return (
      <FormattedDate value={dateValue} {...dateFormatOptions}>
        {date => (
          <FormattedTime value={dateValue} {...timeFormatOptions}>
            {time => this.renderDateTime(date, time, dateValue)}
          </FormattedTime>
        )}
      </FormattedDate>
    )
  }
}

DateTime.Relative = RelativeTime

DateTime.propTypes = {
  children: PropTypes.func,
  className: PropTypes.string,
  /** The time to be displayed. */
  date: PropTypes.bool,
  /** Whether to show the time. */
  dateFormatOptions: PropTypes.shape({}),
  // See https://formatjs.io/docs/react-intl/components/#formatteddate
  time: PropTypes.bool,
  // See https://formatjs.io/docs/react-intl/components/#formattedtime
  timeFormatOptions: PropTypes.shape({}),
  value: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number, // Support timestamps.
    PropTypes.instanceOf(Date),
  ]).isRequired,
}

DateTime.defaultProps = {
  className: undefined,
  children: undefined,
  date: true,
  time: true,
  dateFormatOptions: {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  },
  timeFormatOptions: {
    hour: 'numeric',
    minute: 'numeric',
    second: 'numeric',
    hour12: false,
  },
}

export default DateTime
