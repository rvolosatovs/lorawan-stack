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

import { hot } from 'react-hot-loader/root'
import { useSelector } from 'react-redux'
import React from 'react'
import { Switch, Route } from 'react-router-dom'
import { ConnectedRouter } from 'connected-react-router'
import { Helmet } from 'react-helmet'

import ErrorView from '@ttn-lw/lib/components/error-view'
import { FullViewError } from '@ttn-lw/lib/components/full-view-error/error'

import Landing from '@account/views/landing'
import Authorize from '@account/views/authorize'

import PropTypes from '@ttn-lw/lib/prop-types'
import dev from '@ttn-lw/lib/dev'
import {
  selectApplicationSiteName,
  selectApplicationSiteTitle,
  selectPageData,
} from '@ttn-lw/lib/selectors/env'

import { selectUser } from '@account/store/selectors/user'

import Front from '../front'

const siteName = selectApplicationSiteName()
const siteTitle = selectApplicationSiteTitle()
const pageData = selectPageData()

const AccountApp = ({ history }) => {
  const user = useSelector(selectUser)

  if (pageData && pageData.error) {
    return (
      <ConnectedRouter history={history}>
        <FullViewError error={pageData.error} />
      </ConnectedRouter>
    )
  }

  return (
    <ConnectedRouter history={history}>
      <ErrorView ErrorComponent={FullViewError}>
        <React.Fragment>
          <Helmet
            titleTemplate={`%s - ${siteTitle ? `${siteTitle} - ` : ''}${siteName}`}
            defaultTitle={`${siteTitle ? `${siteTitle} - ` : ''}${siteName}`}
          />
          <Switch>
            <Route path="/authorize" component={Authorize} />
            <Route path="/" component={Boolean(user) ? Landing : Front} />
          </Switch>
        </React.Fragment>
      </ErrorView>
    </ConnectedRouter>
  )
}

AccountApp.propTypes = {
  history: PropTypes.history.isRequired,
}

const ExportedApp = dev ? hot(AccountApp) : AccountApp

export default ExportedApp
