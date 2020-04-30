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

import { hot } from 'react-hot-loader/root'
import React from 'react'
import { connect } from 'react-redux'
import { ConnectedRouter } from 'connected-react-router'
import { Route, Switch } from 'react-router-dom'
import classnames from 'classnames'

import { ToastContainer } from '@ttn-lw/components/toast'
import Footer from '@ttn-lw/components/footer'

import IntlHelmet from '@ttn-lw/lib/components/intl-helmet'
import { withEnv } from '@ttn-lw/lib/components/env'
import ErrorView from '@ttn-lw/lib/components/error-view'
import ScrollToTop from '@ttn-lw/lib/components/scroll-to-top'
import WithAuth from '@ttn-lw/lib/components/with-auth'

import Header from '@console/containers/header'

import Overview from '@console/views/overview'
import Applications from '@console/views/applications'
import Gateways from '@console/views/gateways'
import Organizations from '@console/views/organizations'
import Admin from '@console/views/admin'
import FullViewError, { FullViewErrorInner } from '@console/views/error'

import PropTypes from '@ttn-lw/lib/prop-types'
import dev from '@ttn-lw/lib/dev'

import {
  selectUser,
  selectUserFetching,
  selectUserError,
  selectUserRights,
  selectUserIsAdmin,
} from '@console/store/selectors/user'

import style from './app.styl'

const GenericNotFound = () => <FullViewErrorInner error={{ statusCode: 404 }} />

@withEnv
@connect(state => ({
  user: selectUser(state),
  fetching: selectUserFetching(state),
  error: selectUserError(state),
  rights: selectUserRights(state),
  isAdmin: selectUserIsAdmin(state),
}))
@(Component => (dev ? hot(Component) : Component))
class ConsoleApp extends React.PureComponent {
  static propTypes = {
    env: PropTypes.env.isRequired,
    error: PropTypes.error,
    fetching: PropTypes.bool.isRequired,
    history: PropTypes.shape({
      push: PropTypes.func,
      replace: PropTypes.func,
    }).isRequired,
    isAdmin: PropTypes.bool,
    rights: PropTypes.rights,
    user: PropTypes.user,
  }
  static defaultProps = {
    user: undefined,
    error: undefined,
    isAdmin: undefined,
    rights: undefined,
  }

  render() {
    const {
      user,
      fetching,
      error,
      rights,
      isAdmin,
      history,
      env: {
        siteTitle,
        pageData,
        siteName,
        config: { supportLink },
      },
    } = this.props

    if (pageData && pageData.error) {
      return (
        <ConnectedRouter history={history}>
          <FullViewError error={pageData.error} />
        </ConnectedRouter>
      )
    }

    return (
      <React.Fragment>
        <ToastContainer />
        <ConnectedRouter history={history}>
          <ScrollToTop />
          <ErrorView ErrorComponent={FullViewError}>
            <div className={style.app}>
              <IntlHelmet
                titleTemplate={`%s - ${siteTitle ? `${siteTitle} - ` : ''}${siteName}`}
                defaultTitle={siteName}
              />
              <div id="modal-container" />
              <Header className={style.header} />
              <main className={style.main}>
                <WithAuth
                  user={user}
                  fetching={fetching}
                  error={error}
                  errorComponent={FullViewErrorInner}
                  rights={rights}
                  isAdmin={isAdmin}
                >
                  <div className={classnames('breadcrumbs', style.mobileBreadcrumbs)} />
                  <div className={style.sidebar} id="sidebar" />
                  <div className={style.content}>
                    <div className={classnames('breadcrumbs', style.desktopBreadcrumbs)} />
                    <div className={style.stage}>
                      <Switch>
                        <Route exact path="/" component={Overview} />
                        <Route path="/applications" component={Applications} />
                        <Route path="/gateways" component={Gateways} />
                        <Route path="/organizations" component={Organizations} />
                        <Route path="/admin" component={Admin} />
                        <Route component={GenericNotFound} />
                      </Switch>
                    </div>
                  </div>
                </WithAuth>
              </main>
              <Footer className={style.footer} supportLink={supportLink} />
            </div>
          </ErrorView>
        </ConnectedRouter>
      </React.Fragment>
    )
  }
}

export default ConsoleApp
