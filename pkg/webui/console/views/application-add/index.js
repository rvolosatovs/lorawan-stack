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
import { connect } from 'react-redux'
import bind from 'autobind-decorator'
import { defineMessages } from 'react-intl'
import { push } from 'connected-react-router'

import api from '@console/api'

import PageTitle from '@ttn-lw/components/page-title'
import Form from '@ttn-lw/components/form'
import Input from '@ttn-lw/components/input'
import Checkbox from '@ttn-lw/components/checkbox'
import SubmitButton from '@ttn-lw/components/submit-button'
import toast from '@ttn-lw/components/toast'
import SubmitBar from '@ttn-lw/components/submit-bar'

import OwnersSelect from '@console/containers/owners-select'

import withFeatureRequirement from '@console/lib/components/with-feature-requirement'

import Yup from '@ttn-lw/lib/yup'
import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'
import { getApplicationId } from '@ttn-lw/lib/selectors/id'

import { id as applicationIdRegexp, address } from '@console/lib/regexp'
import { mayCreateApplications, mayLinkApplication } from '@console/lib/feature-checks'

import { selectUserId, selectUserRights } from '@console/store/selectors/user'

const m = defineMessages({
  applicationName: 'Application name',
  appIdPlaceholder: 'my-new-application',
  appNamePlaceholder: 'My new application',
  appDescPlaceholder: 'Description for my new application',
  appDescDescription:
    'Optional application description; can also be used to save notes about the application',
  createApplication: 'Create application',
  linking: 'Linking',
  linkAutomatically: 'Link new application to Network Server automatically',
  linkFailure: 'There was an error and the application could not be linked',
  linkFailureTitle: 'Application link failed',
})

const validationSchema = Yup.object().shape({
  owner_id: Yup.string().required(sharedMessages.validateRequired),
  application_id: Yup.string()
    .matches(applicationIdRegexp, Yup.passValues(sharedMessages.validateIdFormat))
    .min(2, Yup.passValues(sharedMessages.validateTooShort))
    .max(25, Yup.passValues(sharedMessages.validateTooLong))
    .required(sharedMessages.validateRequired),
  name: Yup.string()
    .min(2, Yup.passValues(sharedMessages.validateTooShort))
    .max(2000, Yup.passValues(sharedMessages.validateTooLong)),
  link: Yup.boolean(),
  description: Yup.string(),
  network_server_address: Yup.string().when('link', {
    is: true,
    then: schema => schema.matches(address, Yup.passValues(sharedMessages.validateAddressFormat)),
  }),
})

@withFeatureRequirement(mayCreateApplications, { redirect: '/applications' })
@connect(
  state => ({
    userId: selectUserId(state),
    rights: selectUserRights(state),
  }),
  dispatch => ({
    navigateToApplication: appId => dispatch(push(`/applications/${appId}`)),
  }),
)
export default class Add extends React.Component {
  static propTypes = {
    navigateToApplication: PropTypes.func.isRequired,
    rights: PropTypes.rights.isRequired,
    userId: PropTypes.string.isRequired,
  }

  constructor(props) {
    super(props)
    const { rights } = this.props
    this.state = {
      error: '',
      link: mayLinkApplication.check(rights),
    }
  }

  @bind
  async handleSubmit(values, { setSubmitting }) {
    const { userId, navigateToApplication } = this.props
    const { owner_id, application_id, name, description } = values

    await this.setState({ error: '' })

    try {
      const result = await api.application.create(
        owner_id,
        {
          ids: { application_id },
          name,
          description,
        },
        userId === owner_id,
      )

      const appId = getApplicationId(result)

      if (values.link) {
        try {
          const key = {
            name: 'Application Server Linking',
            rights: ['RIGHT_APPLICATION_LINK'],
          }
          const { key: api_key } = await api.application.apiKeys.create(appId, key)
          await api.application.link.set(appId, {
            api_key,
            network_server_address: values.network_server_address,
          })
        } catch (err) {
          toast({
            title: m.linkFailureTitle,
            message: m.linkFailure,
            type: toast.types.ERROR,
          })
        }
      }

      navigateToApplication(appId)
    } catch (error) {
      setSubmitting(false)

      await this.setState({ error })
    }
  }

  @bind
  handleLinkChange(event) {
    this.setState({
      link: event.target.checked,
    })
  }

  get linkingBit() {
    const { link } = this.state

    return (
      <React.Fragment>
        <Form.Field
          onChange={this.handleLinkChange}
          title={m.linking}
          name="link"
          label={m.linkAutomatically}
          component={Checkbox}
        />
        <Form.Field
          component={Input}
          description={sharedMessages.nsEmptyDefault}
          name="network_server_address"
          title={sharedMessages.nsAddress}
          disabled={!link}
        />
      </React.Fragment>
    )
  }

  render() {
    const { error } = this.state
    const { userId, rights } = this.props
    const mayLink = mayLinkApplication.check(rights)

    const initialValues = {
      application_id: '',
      name: '',
      description: '',
      link: mayLink,
      network_server_address: '',
      owner_id: userId,
    }

    return (
      <Container>
        <PageTitle tall title={sharedMessages.addApplication} />
        <Row>
          <Col md={10} lg={9}>
            <Form
              error={error}
              onSubmit={this.handleSubmit}
              initialValues={initialValues}
              validationSchema={validationSchema}
            >
              <OwnersSelect name="owner_id" required autoFocus />
              <Form.Field
                title={sharedMessages.appId}
                name="application_id"
                placeholder={m.appIdPlaceholder}
                required
                component={Input}
              />
              <Form.Field
                title={m.applicationName}
                name="name"
                placeholder={m.appNamePlaceholder}
                component={Input}
              />
              <Form.Field
                title={sharedMessages.description}
                type="textarea"
                name="description"
                placeholder={m.appDescPlaceholder}
                description={m.appDescDescription}
                component={Input}
              />
              {mayLink && this.linkingBit}
              <SubmitBar>
                <Form.Submit message={m.createApplication} component={SubmitButton} />
              </SubmitBar>
            </Form>
          </Col>
        </Row>
      </Container>
    )
  }
}
