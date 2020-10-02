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
import bind from 'autobind-decorator'
import { defineMessages } from 'react-intl'

import Form from '@ttn-lw/components/form'
import Input from '@ttn-lw/components/input'

import Yup from '@ttn-lw/lib/yup'
import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import { id as organizationIdRegexp } from '@console/lib/regexp'

const validationSchema = Yup.object().shape({
  ids: Yup.object().shape({
    organization_id: Yup.string()
      .matches(organizationIdRegexp, Yup.passValues(sharedMessages.validateIdFormat))
      .min(2, Yup.passValues(sharedMessages.validateTooShort))
      .max(25, Yup.passValues(sharedMessages.validateTooLong))
      .required(sharedMessages.validateRequired),
  }),
  name: Yup.string()
    .min(2, Yup.passValues(sharedMessages.validateTooShort))
    .max(50, Yup.passValues(sharedMessages.validateTooLong)),
  description: Yup.string().max(2000, Yup.passValues(sharedMessages.validateTooLong)),
})

const m = defineMessages({
  orgDescPlaceholder: 'Description for my new organization',
  orgDescDescription:
    'Optional organization description; can also be used to save notes about the organization',
  orgIdPlaceholder: 'my-new-organization',
  orgNamePlaceholder: 'My new organization',
})

class OrganizationForm extends React.Component {
  static propTypes = {
    children: PropTypes.node.isRequired,
    error: PropTypes.error,
    formRef: Form.propTypes.formikRef,
    initialValues: PropTypes.shape({
      ids: PropTypes.shape({
        organization_id: PropTypes.string.isRequired,
      }).isRequired,
      name: PropTypes.string,
      description: PropTypes.string,
    }).isRequired,
    onSubmit: PropTypes.func.isRequired,
    onSubmitFailure: PropTypes.func,
    onSubmitSuccess: PropTypes.func,
    update: PropTypes.bool,
  }

  static defaultProps = {
    update: false,
    error: '',
    onSubmitFailure: () => null,
    onSubmitSuccess: () => null,
    formRef: undefined,
  }

  @bind
  async handleSubmit(values, { resetForm }) {
    const { onSubmit, onSubmitSuccess, onSubmitFailure } = this.props
    const castedValues = validationSchema.cast(values)

    try {
      const result = await onSubmit(castedValues)
      onSubmitSuccess(result)
    } catch (error) {
      resetForm({ values })
      onSubmitFailure(error)
    }
  }

  render() {
    const { error, initialValues, update, children, formRef } = this.props

    return (
      <Form
        error={error}
        onSubmit={this.handleSubmit}
        initialValues={initialValues}
        validationSchema={validationSchema}
        formikRef={formRef}
      >
        <Form.Field
          title={sharedMessages.organizationId}
          name="ids.organization_id"
          placeholder={m.orgIdPlaceholder}
          autoFocus={!update}
          disabled={update}
          required
          component={Input}
        />
        <Form.Field
          title={sharedMessages.name}
          name="name"
          placeholder={m.orgNamePlaceholder}
          component={Input}
        />
        <Form.Field
          title={sharedMessages.description}
          name="description"
          type="textarea"
          placeholder={m.orgDescPlaceholder}
          description={m.orgDescDescription}
          component={Input}
        />
        {children}
      </Form>
    )
  }
}
export default OrganizationForm
