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
import classnames from 'classnames'
import { defineMessages } from 'react-intl'

import Button from '../button'
import Message from '../../lib/components/message'
import List from '../list'
import Icon from '../icon'
import getEventComponentByName from '../event/types'
import sharedMessages from '../../lib/shared-messages'
import PropTypes from '../../lib/prop-types'
import EventsWidget from './widget'

import style from './events.styl'

const m = defineMessages({
  truncated: 'Events have been truncated',
  showing: 'Showing all available events',
})

@bind
class Events extends React.PureComponent {

  renderEvent (event) {
    const { component: Component, type } = getEventComponentByName(event.name)

    return (
      <List.Item className={style.listItem}>
        <Component
          overviewClassName={style.event}
          expandedClassName={style.eventData}
          event={event}
          type={type}
        />
      </List.Item>
    )
  }

  onPause () {
    const { onPause } = this.props

    onPause()
  }

  onClear () {
    const { onClear } = this.props

    onClear()
  }

  getEventkey (event) {
    return `${event.time}-${event.name}`
  }

  render () {
    const {
      className,
      events,
      paused,
      onClear,
      onPause,
      emitterId,
      limit,
    } = this.props

    let limitedEvents = events
    const truncated = events.length > limit
    if (truncated) {
      limitedEvents = events.slice(0, limit)
    }

    const header = (
      <Header
        paused={paused}
        onPause={onPause}
        onClear={onClear}
      />
    )

    const footer = (
      <Footer truncated={truncated} />
    )

    return (
      <List
        bordered
        size="none"
        className={className}
        listClassName={style.list}
        header={header}
        footer={footer}
        items={limitedEvents}
        renderItem={this.renderEvent}
        rowKey={this.getEventkey}
        emptyMessage={sharedMessages.noEvents}
        emptyMessageValues={{ entityId: emitterId }}
      />
    )
  }
}

Events.propTypes = {
  events: PropTypes.arrayOf(PropTypes.event),
  paused: PropTypes.bool.isRequired,
  emitterId: PropTypes.string.isRequired,
  onPause: PropTypes.func.isRequired,
  onClear: PropTypes.func.isRequired,
  limit: PropTypes.number,
}

Events.defaultProps = {
  events: [],
  limit: 100,
}

Events.Widget = EventsWidget

const Header = function (props) {
  const {
    paused,
    onPause,
    onClear,
  } = props

  const pauseMessage = paused ? sharedMessages.resume : sharedMessages.pause
  const pauseIcon = paused ? 'play_arrow' : 'pause'

  return (
    <div className={style.header}>
      <div className={style.headerColumns}>
        <Message className={style.headerColumnsTime} content={sharedMessages.time} />
        <Message className={style.headerColumnsId} content={sharedMessages.entityId} />
        <Message content={sharedMessages.data} />
      </div>
      <div className={style.headerActions}>
        <Button
          onClick={onPause}
          message={pauseMessage}
          naked
          secondary
          icon={pauseIcon}
        />
        <Button
          onClick={onClear}
          message={sharedMessages.clear}
          naked
          secondary
          icon="delete"
        />
      </div>
    </div>
  )
}

const Footer = function (props) {
  const {
    truncated,
  } = props

  return (
    <div className={classnames(style.footer, {
      [style.footerTruncated]: truncated,
    })}
    >
      {
        truncated ? (
          <React.Fragment>
            <Icon icon="warning" />
            <Message
              content="Events have been truncated"
              className={style.footerTruncatedText}
            />
          </React.Fragment>
        ) : (
          <Message content={m.showing} />
        )
      }
    </div>
  )
}

export default Events
