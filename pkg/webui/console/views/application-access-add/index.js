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
import { connect } from 'react-redux'
import * as Yup from 'yup'
import { replace } from 'connected-react-router'
import { defineMessages } from 'react-intl'

import Spinner from '../../../components/spinner'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import sharedMessages from '../../../lib/shared-messages'
import Form from '../../../components/form'
import Field from '../../../components/field'
import Button from '../../../components/button'
import Message from '../../../lib/components/message'
import FieldGroup from '../../../components/field/group'
import IntlHelmet from '../../../lib/components/intl-helmet'

import { getApplicationsRightsList } from '../../store/actions/applications'
import api from '../../api'

import style from './application-access-add.styl'

const m = defineMessages({
  accessAdd: 'Add Access Key',
})

const validationSchema = Yup.object().shape({
  name: Yup.string()
    .min(2, sharedMessages.validateTooShort)
    .max(50, sharedMessages.validateTooLong),
  rights: Yup.object().test(
    'rights',
    sharedMessages.validateRights,
    values => Object.values(values).reduce((acc, curr) => acc || curr, false)
  ),
})

@connect(function ({ rights }, props) {
  const appId = props.match.params.appId

  return {
    appId,
    fetching: rights.applications.fetching,
    error: rights.applications.error,
    rights: rights.applications.rights,
  }
})
@withBreadcrumb('apps.single.access.add', function (props) {
  const appId = props.appId
  return (
    <Breadcrumb
      path={`/console/applications/${appId}/access/add`}
      icon="add"
      content={sharedMessages.add}
    />
  )
})
@bind
export default class ApplicationAccessAdd extends React.Component {

  state = {
    error: '',
  }

  componentDidMount () {
    const { dispatch, appId } = this.props

    dispatch(getApplicationsRightsList(appId))
  }

  async handleSubmit (values, { resetForm }) {
    const { name, rights } = values
    const { appId, dispatch } = this.props

    const key = {
      name,
      rights: Object.keys(rights).filter(r => rights[r]),
    }

    await this.setState({ error: '' })

    try {
      await api.application.apiKeys.create(appId, key)
      resetForm(values)
      dispatch(replace(`/console/applications/${appId}/access`))
    } catch (error) {
      resetForm(values)
      await this.setState(error)
    }
  }

  render () {
    const { rights, fetching, error } = this.props

    if (error) {
      return 'ERROR'
    }

    if (fetching || !rights.length) {
      return <Spinner center />
    }

    const { rightsItems, rightsValues } = rights.reduce(
      function (acc, right) {
        acc.rightsItems.push(
          <Field
            className={style.rightLabel}
            key={right}
            name={right}
            type="checkbox"
            title={{ id: `enum:${right}` }}
            form
          />
        )
        acc.rightsValues[right] = false

        return acc
      },
      {
        rightsItems: [],
        rightsValues: {},
      }
    )

    const initialFormValues = {
      name: '',
      rights: rightsValues,
    }

    return (
      <Container>
        <Row>
          <Col lg={8} md={12}>
            <IntlHelmet title={m.accessAdd} />
            <Message component="h2" content={m.accessAdd} />
          </Col>
        </Row>
        <Row>
          <Col lg={8} md={12}>
            <Form
              horizontal
              error={this.state.error}
              onSubmit={this.handleSubmit}
              initialValues={initialFormValues}
              validationSchema={validationSchema}
            >
              <Message
                component="h4"
                content={sharedMessages.generalInformation}
              />
              <Field
                title={sharedMessages.name}
                name="name"
                type="text"
                autoFocus
              />
              <FieldGroup
                name="rights"
                title={sharedMessages.rights}
              >
                {rightsItems}
              </FieldGroup>
              <div className={style.submitBar}>
                <Button type="submit" message={sharedMessages.createAccessKey} />
              </div>
            </Form>
          </Col>
        </Row>
      </Container>
    )
  }
}
