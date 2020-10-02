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
import { Container, Col, Row } from 'react-grid-system'
import bind from 'autobind-decorator'
import { defineMessages } from 'react-intl'
import { push } from 'connected-react-router'
import { connect } from 'react-redux'

import api from '@oauth/api'

import Button from '@ttn-lw/components/button'
import Form from '@ttn-lw/components/form'
import Input from '@ttn-lw/components/input'
import SubmitButton from '@ttn-lw/components/submit-button'

import IntlHelmet from '@ttn-lw/lib/components/intl-helmet'
import Message from '@ttn-lw/lib/components/message'

import Yup from '@ttn-lw/lib/yup'
import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'
import { id as userRegexp } from '@ttn-lw/lib/regexp'

import style from './forgot-password.styl'

const m = defineMessages({
  loginPage: 'Login page',
  forgotPassword: 'Forgot password',
  passwordRequested: 'You will receive an email with reset instructions shortly',
  goToLogin: 'Go to login',
  send: 'Send',
  resetPasswordDescription:
    'Please enter your username to receive an email with reset instructions',
  requestTempPassword: 'Reset password',
})

const validationSchema = Yup.object().shape({
  user_id: Yup.string()
    .min(3, Yup.passValues(sharedMessages.validateTooShort))
    .max(36, Yup.passValues(sharedMessages.validateTooLong))
    .matches(userRegexp, Yup.passValues(sharedMessages.validateIdFormat))
    .required(sharedMessages.validateRequired),
})

const initialValues = { user_id: '' }

@connect(
  null,
  {
    handleCancel: () => push('/login'),
  },
)
export default class ForgotPassword extends React.PureComponent {
  static propTypes = {
    handleCancel: PropTypes.func.isRequired,
  }

  state = {
    error: '',
    info: '',
    requested: false,
  }

  @bind
  async handleSubmit(values, { setSubmitting }) {
    try {
      await api.users.resetPassword(values.user_id)
      this.setState({
        error: '',
        info: m.passwordRequested,
        requested: true,
      })
    } catch (error) {
      this.setState({
        error: error.response.data,
        info: '',
      })
    } finally {
      setSubmitting(false)
    }
  }

  render() {
    const { error, info, requested } = this.state
    const { handleCancel } = this.props
    const cancelButtonText = requested ? m.goToLogin : sharedMessages.cancel

    return (
      <Container className={style.fullHeight}>
        <Row justify="center" align="center" className={style.fullHeight}>
          <Col sm={12} md={8} lg={5}>
            <IntlHelmet title={m.forgotPassword} />
            <Message content={m.requestTempPassword} component="h1" className={style.title} />
            <Message
              content={m.resetPasswordDescription}
              component="h4"
              className={style.description}
            />
            <Form
              onSubmit={this.handleSubmit}
              initialValues={initialValues}
              error={error}
              info={info}
              validationSchema={validationSchema}
              horizontal={false}
            >
              <Form.Field
                title={sharedMessages.userId}
                name="user_id"
                component={Input}
                autoComplete="username"
                autoFocus
                required
              />
              <Form.Submit component={SubmitButton} message={m.send} />
              <Button naked secondary message={cancelButtonText} onClick={handleCancel} />
            </Form>
          </Col>
        </Row>
      </Container>
    )
  }
}
