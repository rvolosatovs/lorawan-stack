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

import Icon from '@ttn-lw/components/icon'

import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import { getEntityId } from '@ttn-lw/lib/selectors/id'

import Event from '../..'

import style from './message.styl'

@bind
class MessageEvent extends React.PureComponent {
  static propTypes = {
    className: PropTypes.string,
    event: PropTypes.event.isRequired,
    expandedClassName: PropTypes.string,
    overviewClassName: PropTypes.string,
    type: PropTypes.oneOf(['downlink', 'uplink']).isRequired,
    widget: PropTypes.bool,
  }

  static defaultProps = {
    className: undefined,
    expandedClassName: undefined,
    overviewClassName: undefined,
    widget: false,
  }

  render() {
    const { className, event, type, widget, overviewClassName, expandedClassName } = this.props

    const entityId = getEntityId(event.identifiers[0])
    const icon = type === 'downlink' ? 'downlink' : 'uplink'

    const eventContent = <Message content={{ id: `event:${event.name}` }} />
    const eventIcon = <Icon icon={icon} className={style.messageIcon} />

    return (
      <Event
        className={className}
        overviewClassName={overviewClassName}
        expandedClassName={expandedClassName}
        icon={eventIcon}
        time={event.time}
        emitter={entityId}
        content={eventContent}
        data={event.data}
        widget={widget}
      />
    )
  }
}

export default MessageEvent
