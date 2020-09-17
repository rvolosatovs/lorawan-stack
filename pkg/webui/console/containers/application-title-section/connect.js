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

import { connect } from 'react-redux'

import withRequest from '@ttn-lw/lib/components/with-request'

import {
  checkFromState,
  mayViewOrEditApplicationApiKeys,
  mayViewOrEditApplicationCollaborators,
  mayViewApplicationDevices,
  mayLinkApplication,
} from '@console/lib/feature-checks'

import { getCollaboratorsList } from '@console/store/actions/collaborators'
import { getApiKeysList } from '@console/store/actions/api-keys'
import { getApplicationLink } from '@console/store/actions/link'
import { getApplicationDeviceCount } from '@console/store/actions/applications'

import {
  selectApplicationById,
  selectApplicationDeviceCount,
  selectApplicationDevicesFetching,
  selectApplicationDevicesError,
  selectApplicationLinkIndicator,
  selectApplicationLinkStats,
  selectApplicationLinkFetching,
  selectApplicationLastSeen,
} from '@console/store/selectors/applications'
import {
  selectCollaboratorsTotalCount,
  selectCollaboratorsFetching,
  selectCollaboratorsError,
} from '@console/store/selectors/collaborators'
import {
  selectApiKeysTotalCount,
  selectApiKeysFetching,
  selectApiKeysError,
} from '@console/store/selectors/api-keys'

const mapStateToProps = (state, props) => {
  const apiKeysTotalCount = selectApiKeysTotalCount(state)
  const apiKeysFetching = selectApiKeysFetching(state)
  const apiKeysError = selectApiKeysError(state)
  const collaboratorsTotalCount = selectCollaboratorsTotalCount(state, { id: props.appId })
  const collaboratorsFetching = selectCollaboratorsFetching(state)
  const collaboratorsError = selectCollaboratorsError(state)
  const devicesTotalCount = selectApplicationDeviceCount(state)
  const devicesFetching = selectApplicationDevicesFetching(state)
  const devicesError = selectApplicationDevicesError(state)
  const linked = selectApplicationLinkIndicator(state)
  const linkStats = selectApplicationLinkStats(state)
  const lastSeen = selectApplicationLastSeen(state)
  const linkFetching = selectApplicationLinkFetching(state)

  const fetching = apiKeysFetching || collaboratorsFetching || devicesFetching || linkFetching

  return {
    mayViewLink: checkFromState(mayLinkApplication, state),
    mayViewCollaborators: checkFromState(mayViewOrEditApplicationCollaborators, state),
    mayViewApiKeys: checkFromState(mayViewOrEditApplicationApiKeys, state),
    mayViewDevices: checkFromState(mayViewApplicationDevices, state),
    mayViewApplicationLink: checkFromState(mayLinkApplication, state),
    application: selectApplicationById(state, props.appId),
    apiKeysTotalCount,
    apiKeysErrored: Boolean(apiKeysError),
    collaboratorsTotalCount,
    collaboratorsErrored: Boolean(collaboratorsError),
    devicesTotalCount,
    devicesErrored: Boolean(devicesError),
    linked,
    linkStats,
    lastSeen,
    fetching,
  }
}

const mapDispatchToProps = dispatch => ({
  loadData(mayViewCollaborators, mayViewApiKeys, mayViewLink, mayViewDevices, appId) {
    if (mayViewCollaborators) {
      dispatch(getCollaboratorsList('application', appId))
    }

    if (mayViewApiKeys) {
      dispatch(getApiKeysList('application', appId))
    }

    if (mayViewLink) {
      dispatch(getApplicationLink(appId))
    }

    if (mayViewDevices) {
      dispatch(getApplicationDeviceCount(appId))
    }
  },
})

const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
  loadData: () =>
    dispatchProps.loadData(
      stateProps.mayViewCollaborators,
      stateProps.mayViewApiKeys,
      stateProps.mayViewLink,
      stateProps.mayViewDevices,
      ownProps.appId,
    ),
})

export default TitleSection =>
  connect(
    mapStateToProps,
    mapDispatchToProps,
    mergeProps,
  )(withRequest(({ appId, loadData }) => loadData(appId), () => false)(TitleSection))
