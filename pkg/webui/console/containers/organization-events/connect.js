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

import { connect } from 'react-redux'

import {
  clearOrganizationEventsStream,
  pauseOrganizationEventsStream,
  resumeOrganizationEventsStream,
} from '@console/store/actions/organizations'

import {
  selectOrganizationEvents,
  selectOrganizationEventsPaused,
  selectOrganizationEventsTruncated,
} from '@console/store/selectors/organizations'

const mapStateToProps = (state, props) => {
  const { orgId } = props

  return {
    events: selectOrganizationEvents(state, orgId),
    paused: selectOrganizationEventsPaused(state, orgId),
    truncated: selectOrganizationEventsTruncated(state, orgId),
  }
}

const mapDispatchToProps = (dispatch, ownProps) => ({
  onClear: () => dispatch(clearOrganizationEventsStream(ownProps.orgId)),
  onPauseToggle: paused =>
    paused
      ? dispatch(resumeOrganizationEventsStream(ownProps.orgId))
      : dispatch(pauseOrganizationEventsStream(ownProps.orgId)),
})

export default Events =>
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(Events)
