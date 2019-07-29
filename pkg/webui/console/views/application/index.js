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
import { Switch, Route } from 'react-router'
import { replace } from 'connected-react-router'

import sharedMessages from '../../../lib/shared-messages'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import { withSideNavigation } from '../../../components/navigation/side/context'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import IntlHelmet from '../../../lib/components/intl-helmet'
import withRequest from '../../../lib/components/with-request'

import ApplicationOverview from '../application-overview'
import ApplicationGeneralSettings from '../application-general-settings'
import ApplicationApiKeys from '../application-api-keys'
import ApplicationLink from '../application-link'
import ApplicationCollaborators from '../application-collaborators'
import ApplicationData from '../application-data'
import ApplicationPayloadFormatters from '../application-payload-formatters'
import ApplicationIntegrations from '../application-integrations'

import { getApplicationId } from '../../../lib/selectors/id'
import {
  getApplication,
  stopApplicationEventsStream,
} from '../../store/actions/applications'
import {
  selectSelectedApplication,
  selectApplicationFetching,
  selectApplicationError,
} from '../../store/selectors/applications'

import Devices from '../devices'
import withEnv, { EnvProvider } from '../../../lib/components/env'

@connect(function (state, props) {
  return {
    appId: props.match.params.appId,
    fetching: selectApplicationFetching(state),
    application: selectSelectedApplication(state),
    error: selectApplicationError(state),
  }
},
dispatch => ({
  stopStream: id => dispatch(stopApplicationEventsStream(id)),
  getApplication: id => dispatch(getApplication(id, 'name,description')),
  redirectToList: () => dispatch(replace('/console/applications')),
}))
@withRequest(
  ({ appId, getApplication }) => getApplication(appId),
  ({ fetching, application }) => fetching || !Boolean(application)
)
@withSideNavigation(function (props) {
  const matchedUrl = props.match.url

  return {
    header: { title: props.appId, icon: 'application' },
    entries: [
      {
        title: sharedMessages.overview,
        path: matchedUrl,
        icon: 'overview',
      },
      {
        title: sharedMessages.devices,
        path: `${matchedUrl}/devices`,
        icon: 'devices',
        exact: false,
      },
      {
        title: sharedMessages.data,
        path: `${matchedUrl}/data`,
        icon: 'data',
      },
      {
        title: sharedMessages.link,
        path: `${matchedUrl}/link`,
        icon: 'link',
      },
      {
        title: sharedMessages.payloadFormatters,
        icon: 'code',
        nested: true,
        items: [
          {
            title: sharedMessages.uplink,
            path: `${matchedUrl}/payload-formatters/uplink`,
            icon: 'uplink',
          },
          {
            title: sharedMessages.downlink,
            path: `${matchedUrl}/payload-formatters/downlink`,
            icon: 'downlink',
          },
        ],
      },
      {
        title: sharedMessages.integrations,
        path: `${matchedUrl}/integrations`,
        icon: 'integration',
        exact: false,
      },
      {
        title: sharedMessages.collaborators,
        path: `${matchedUrl}/collaborators`,
        icon: 'organization',
        exact: false,
      },
      {
        title: sharedMessages.apiKeys,
        path: `${matchedUrl}/api-keys`,
        icon: 'api_keys',
        exact: false,
      },
      {
        title: sharedMessages.generalSettings,
        path: `${matchedUrl}/general-settings`,
        icon: 'general_settings',
      },
    ],
  }
})
@withBreadcrumb('apps.single', function (props) {
  const { appId } = props
  return (
    <Breadcrumb
      path={`/console/applications/${appId}`}
      icon="application"
      content={appId}
    />
  )
})
@withEnv
export default class Application extends React.Component {

  componentDidUpdate (prevProps) {
    const { appId, application, redirectToList } = this.props

    const isSame = appId === getApplicationId(prevProps.application)
    const isDeleted = Boolean(prevProps.application) && !Boolean(application)

    if (isSame && isDeleted) {
      redirectToList()
    }
  }

  componentWillUnmount () {
    const { appId, stopStream } = this.props

    stopStream(appId)
  }

  render () {
    const { match, application, appId, env } = this.props

    return (
      <EnvProvider env={env}>
        <IntlHelmet
          titleTemplate={`%s - ${application.name || appId} - ${env.site_name}`}
        />
        <Switch>
          <Route exact path={`${match.path}`} component={ApplicationOverview} />
          <Route path={`${match.path}/general-settings`} component={ApplicationGeneralSettings} />
          <Route path={`${match.path}/api-keys`} component={ApplicationApiKeys} />
          <Route path={`${match.path}/link`} component={ApplicationLink} />
          <Route path={`${match.path}/devices`} component={Devices} />
          <Route path={`${match.path}/collaborators`} component={ApplicationCollaborators} />
          <Route path={`${match.path}/data`} component={ApplicationData} />
          <Route path={`${match.path}/payload-formatters`} component={ApplicationPayloadFormatters} />
          <Route path={`${match.path}/integrations`} component={ApplicationIntegrations} />
        </Switch>
      </EnvProvider>
    )
  }
}
