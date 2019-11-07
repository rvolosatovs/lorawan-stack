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
import withFeatureRequirement from '../../lib/components/with-feature-requirement'
import { mayViewApplicationEvents } from '../../lib/feature-checks'

import {
  clearApplicationEventsStream,
  startApplicationEventsStream,
} from '../../store/actions/applications'

import {
  selectApplicationEvents,
  selectApplicationEventsStatus,
  selectApplicationEventsError,
} from '../../store/selectors/applications'

@withFeatureRequirement(mayViewApplicationEvents)
@connect(
  null,
  (dispatch, ownProps) => ({
    onClear: () => dispatch(clearApplicationEventsStream(ownProps.appId)),
    onRestart: () => dispatch(startApplicationEventsStream(ownProps.appId)),
  }),
)
@bind
export default class ApplicationEvents extends React.Component {
  static propTypes = {
    appId: PropTypes.string.isRequired,
    onClear: PropTypes.func.isRequired,
    onRestart: PropTypes.func.isRequired,
    widget: PropTypes.bool,
  }

  static defaultProps = {
    widget: false,
  }

  render() {
    const { appId, widget, onClear, onRestart } = this.props

    return (
      <EventsSubscription
        id={appId}
        widget={widget}
        eventsSelector={selectApplicationEvents}
        statusSelector={selectApplicationEventsStatus}
        errorSelector={selectApplicationEventsError}
        onClear={onClear}
        onRestart={onRestart}
        toAllUrl={`/applications/${appId}/data`}
      />
    )
  }
}
