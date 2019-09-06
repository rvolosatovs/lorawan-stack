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
import Query from 'query-string'
import bind from 'autobind-decorator'
import { connect } from 'react-redux'
import { defineMessages } from 'react-intl'
import { Redirect } from 'react-router-dom'
import { Container, Row, Col } from 'react-grid-system'

import PropTypes from '../../../lib/prop-types'
import { selectApplicationSiteName, selectApplicationRootPath } from '../../../lib/selectors/env'
import Button from '../../../components/button'
import IntlHelmet from '../../../lib/components/intl-helmet'
import sharedMessages from '../../../lib/shared-messages'
import Message from '../../../lib/components/message'

import style from './login.styl'

const m = defineMessages({
  welcome: 'Welcome to {siteName}',
  loginSub: 'You need to be logged in to use this site',
})

@connect(state => ({
  user: state.user.user,
  siteName: selectApplicationSiteName(),
  appRoot: selectApplicationRootPath(),
}))
@bind
export default class Login extends React.PureComponent {
  static propTypes = {
    appRoot: PropTypes.string.isRequired,
    siteName: PropTypes.string.isRequired,
    user: PropTypes.user,
  }

  static defaultProps = {
    user: undefined,
  }

  render() {
    const { user, appRoot, siteName } = this.props
    const { next } = Query.parse(location.search)
    const redirectAppend = next ? `?next=${next}` : ''

    // dont show the login page if the user is already logged in
    if (Boolean(user)) {
      return <Redirect to="/" />
    }

    return (
      <div className={style.login}>
        <IntlHelmet title={sharedMessages.login} />
        <Container>
          <Row>
            <Col>
              <Message
                className={style.loginHeader}
                values={{ siteName }}
                component="h2"
                content={m.welcome}
              />
              <Message className={style.loginSub} content={m.loginSub} />
              <Button.AnchorLink
                className={style.loginButton}
                message={sharedMessages.login}
                href={`${appRoot}/login/ttn-stack${redirectAppend}`}
              />
            </Col>
          </Row>
        </Container>
      </div>
    )
  }
}
