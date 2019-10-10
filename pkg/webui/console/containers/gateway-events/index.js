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
import bind from 'autobind-decorator'

import PropTypes from '../../../lib/prop-types'
import EventsSubscription from '../../containers/events-subscription'

import { clearGatewayEventsStream, startGatewayEventsStream } from '../../store/actions/gateways'

import {
  selectGatewayEvents,
  selectGatewayEventsStatus,
  selectGatewayEventsError,
} from '../../store/selectors/gateways'

@connect(
  null,
  (dispatch, ownProps) => ({
    onClear: () => dispatch(clearGatewayEventsStream(ownProps.gtwId)),
    onRestart: () => dispatch(startGatewayEventsStream(ownProps.gtwId)),
  }),
)
@bind
export default class GatewayEvents extends React.Component {
  static propTypes = {
    gtwId: PropTypes.string.isRequired,
    onClear: PropTypes.func.isRequired,
    onRestart: PropTypes.func.isRequired,
    widget: PropTypes.bool,
  }

  static defaultProps = {
    widget: false,
  }

  render() {
    const { gtwId, widget, onClear, onRestart } = this.props

    return (
      <EventsSubscription
        id={gtwId}
        widget={widget}
        eventsSelector={selectGatewayEvents}
        statusSelector={selectGatewayEventsStatus}
        errorSelector={selectGatewayEventsError}
        onRestart={onRestart}
        onClear={onClear}
        toAllUrl={`/gateways/${gtwId}/data`}
      />
    )
  }
}
