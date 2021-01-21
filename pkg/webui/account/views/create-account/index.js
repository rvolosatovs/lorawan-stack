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

import React, { useState, useCallback } from 'react'
import { useDispatch } from 'react-redux'
import { defineMessages } from 'react-intl'
import { Redirect } from 'react-router-dom'
import { push } from 'connected-react-router'
import queryString from 'query-string'

import tts from '@account/api/tts'

import Button from '@ttn-lw/components/button'
import Input from '@ttn-lw/components/input'
import Form from '@ttn-lw/components/form'
import SubmitButton from '@ttn-lw/components/submit-button'

import Message from '@ttn-lw/lib/components/message'
import IntlHelmet from '@ttn-lw/lib/components/intl-helmet'

import style from '@account/views/front/front.styl'

import Yup from '@ttn-lw/lib/yup'
import { selectApplicationSiteName, selectEnableUserRegistration } from '@ttn-lw/lib/selectors/env'
import { id as userRegexp } from '@ttn-lw/lib/regexp'
import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

const m = defineMessages({
  registrationApproved: 'You have successfully registered and can login now',
  createAccount: 'Create account',
  createANewAccount: 'Create a new account',
  registrationPending:
    'You have successfully sent the registration request. You will receive a confirmation once the account has been approved.',
})

const validationSchema = Yup.object().shape({
  user_id: Yup.string()
    .min(3, Yup.passValues(sharedMessages.validateTooShort))
    .max(36, Yup.passValues(sharedMessages.validateTooLong))
    .matches(userRegexp, Yup.passValues(sharedMessages.validateIdFormat))
    .required(sharedMessages.validateRequired),
  name: Yup.string()
    .min(3, Yup.passValues(sharedMessages.validateTooShort))
    .max(50, Yup.passValues(sharedMessages.validateTooLong)),
  password: Yup.string()
    .min(8, Yup.passValues(sharedMessages.validateTooShort))
    .required(sharedMessages.validateRequired),
  primary_email_address: Yup.string()
    .email(sharedMessages.validateEmail)
    .required(sharedMessages.validateRequired),
  password_confirm: Yup.string()
    .oneOf([Yup.ref('password'), null], sharedMessages.validatePasswordMatch)
    .required(sharedMessages.validateRequired),
})

const initialValues = {
  user_id: '',
  name: '',
  primary_email_address: '',
  password: '',
  password_confirm: '',
}

const siteName = selectApplicationSiteName()
const enableUserRegistration = selectEnableUserRegistration()

const getSuccessMessage = state => {
  switch (state) {
    case undefined:
      // Zero value is swallowed by the backend, but means STATE_REQUESTED
      return m.registrationPending
    case 'STATE_APPROVED':
      return m.registrationApproved
    default:
      return ''
  }
}

const CreateAccount = ({ location }) => {
  const [error, setError] = useState(undefined)
  const dispatch = useDispatch()

  const handleSubmit = useCallback(
    async (values, { setSubmitting }) => {
      try {
        setError(undefined)
        const { user_id, ...rest } = values
        const { invitation_token = '' } = queryString.parse(location.search)
        const result = await tts.Users.create({ ids: { user_id }, ...rest }, invitation_token)

        dispatch(
          push(`/login${location.search}`, {
            info: getSuccessMessage(result.state),
          }),
        )
      } catch (error) {
        setSubmitting(false)
        setError(error)
      }
    },
    [dispatch, location.search],
  )

  if (!enableUserRegistration) {
    return <Redirect to={`/login${location.search}`} />
  }

  return (
    <>
      <div className={style.form}>
        <IntlHelmet title={m.createANewAccount} />
        <h1 className={style.title}>
          {siteName}
          <br />
          <Message content={m.createANewAccount} component="strong" />
        </h1>
        <hr className={style.hRule} />
        <Form
          onSubmit={handleSubmit}
          initialValues={initialValues}
          error={error}
          validationSchema={validationSchema}
          horizontal={false}
        >
          <Form.Field
            component={Input}
            required
            title={sharedMessages.userId}
            name="user_id"
            autoComplete="username"
            autoFocus
          />
          <Form.Field
            title={sharedMessages.name}
            name="name"
            component={Input}
            autoComplete="name"
          />
          <Form.Field
            required
            title={sharedMessages.email}
            component={Input}
            name="primary_email_address"
            autoComplete="email"
          />
          <Form.Field
            required
            title={sharedMessages.password}
            name="password"
            type="password"
            component={Input}
            autoComplete="new-password"
          />
          <Form.Field
            required
            title={sharedMessages.confirmPassword}
            name="password_confirm"
            type="password"
            autoComplete="new-password"
            component={Input}
          />
          <div className={style.buttons}>
            <Form.Submit
              component={SubmitButton}
              message={m.createAccount}
              className={style.submitButton}
              alwaysEnabled
            />
            <Button.Link
              to={`/login${location.search}`}
              naked
              secondary
              message={sharedMessages.login}
            />
          </div>
        </Form>
      </div>
    </>
  )
}

CreateAccount.propTypes = {
  location: PropTypes.location.isRequired,
}

export default CreateAccount
