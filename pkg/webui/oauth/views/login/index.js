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
import { withRouter } from 'react-router-dom'
import bind from 'autobind-decorator'
import Query from 'query-string'
import { defineMessages } from 'react-intl'
import { replace } from 'connected-react-router'
import { connect } from 'react-redux'

import api from '@oauth/api'

import Button from '@ttn-lw/components/button'
import Form from '@ttn-lw/components/form'
import Input from '@ttn-lw/components/input'
import SubmitButton from '@ttn-lw/components/submit-button'

import Logo from '@ttn-lw/containers/logo'

import IntlHelmet from '@ttn-lw/lib/components/intl-helmet'
import Message from '@ttn-lw/lib/components/message'

import Yup from '@ttn-lw/lib/yup'
import PropTypes from '@ttn-lw/lib/prop-types'
import { selectApplicationRootPath, selectApplicationSiteName } from '@ttn-lw/lib/selectors/env'
import sharedMessages from '@ttn-lw/lib/shared-messages'
import { id as userRegexp } from '@ttn-lw/lib/regexp'

import style from './login.styl'

const m = defineMessages({
  createAccount: 'Create an account',
  forgotPassword: 'Forgot password?',
  loginToContinue: 'Please login to continue',
})

const validationSchema = Yup.object().shape({
  user_id: Yup.string()
    .min(3, Yup.passValues(sharedMessages.validateTooShort))
    .max(36, Yup.passValues(sharedMessages.validateTooLong))
    .matches(userRegexp, Yup.passValues(sharedMessages.validateIdFormat))
    .required(sharedMessages.validateRequired),
  password: Yup.string().required(sharedMessages.validateRequired),
})

@withRouter
@connect(
  () => ({
    siteName: selectApplicationSiteName(),
  }),
  {
    replace,
  },
)
export default class OAuth extends React.PureComponent {
  static propTypes = {
    location: PropTypes.location.isRequired,
    replace: PropTypes.func.isRequired,
    siteName: PropTypes.string.isRequired,
  }

  constructor(props) {
    super(props)
    this.state = {
      error: '',
    }
  }

  @bind
  async handleSubmit(values, { setSubmitting, setErrors }) {
    try {
      await api.oauth.login(values)

      window.location = url(this.props.location)
    } catch (error) {
      this.setState({
        error: error.response.data,
      })
      setSubmitting(false)
    }
  }

  @bind
  navigateToRegister() {
    const { replace, location } = this.props
    replace('/register', {
      back: `${location.pathname}${location.search}`,
    })
  }

  @bind
  navigateToResetPassword() {
    const { replace, location } = this.props
    replace('/forgot-password', {
      back: `${location.pathname}${location.search}`,
    })
  }

  render() {
    const initialValues = {
      user_id: '',
      password: '',
    }

    const { info } = this.props.location.state || ''
    const { siteName } = this.props

    return (
      <div className={style.fullHeightCenter}>
        <IntlHelmet title={sharedMessages.login} />
        <div>
          <div className={style.left}>
            <div>
              <Logo vertical className={style.logo} />
              <Message content={m.loginToContinue} />
            </div>
          </div>
          <div className={style.right}>
            <h1>{siteName}</h1>
            <Form
              onSubmit={this.handleSubmit}
              initialValues={initialValues}
              error={this.state.error}
              info={info}
              validationSchema={validationSchema}
            >
              <Form.Field
                title={sharedMessages.userId}
                name="user_id"
                component={Input}
                autoComplete="username"
                autoFocus
                required
              />
              <Form.Field
                title={sharedMessages.password}
                component={Input}
                name="password"
                type="password"
                autoComplete="current-password"
                required
              />
              <Form.Submit component={SubmitButton} message={sharedMessages.login} />
              <Button naked message={m.createAccount} onClick={this.navigateToRegister} />
              <Button naked message={m.forgotPassword} onClick={this.navigateToResetPassword} />
            </Form>
          </div>
        </div>
      </div>
    )
  }
}

const appRoot = selectApplicationRootPath()

function url(location, omitQuery = false) {
  const query = Query.parse(location.search)

  const next = query.n || appRoot

  if (omitQuery) {
    return next.split('?')[0]
  }

  return next
}
