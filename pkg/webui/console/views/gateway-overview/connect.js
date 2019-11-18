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

import { getGatewayCollaboratorsList, getGatewayApiKeysList } from '../../store/actions/gateways'
import {
  selectSelectedGateway,
  selectSelectedGatewayId,
  selectGatewayCollaboratorsTotalCount,
  selectGatewayCollaboratorsFetching,
  selectGatewayApiKeysTotalCount,
  selectGatewayApiKeysFetching,
} from '../../store/selectors/gateways'

const mapStateToProps = state => {
  const gtwId = selectSelectedGatewayId(state)
  return {
    gtwId,
    gateway: selectSelectedGateway(state),
    collaboratorsTotalCount: selectGatewayCollaboratorsTotalCount(state, { id: gtwId }),
    apiKeysTotalCount: selectGatewayApiKeysTotalCount(state, { id: gtwId }),
    statusBarFetching:
      selectGatewayCollaboratorsFetching(state) || selectGatewayApiKeysFetching(state),
  }
}
const mapDispatchToProps = dispatch => ({
  loadData(gtwId) {
    dispatch(getGatewayCollaboratorsList(gtwId))
    dispatch(getGatewayApiKeysList(gtwId))
  },
})

export default GatewayOverview =>
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(GatewayOverview)
