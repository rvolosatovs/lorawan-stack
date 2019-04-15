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
import { defineMessages } from 'react-intl'
import * as Yup from 'yup'

import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import sharedMessages from '../../../lib/shared-messages'
import Form from '../../../components/form'
import Field from '../../../components/field'
import Button from '../../../components/button'
import Message from '../../../lib/components/message'
import DataSheet from '../../../components/data-sheet'
import { apiKey, address } from '../../lib/regexp'
import IntlHelmet from '../../../lib/components/intl-helmet'
import Spinner from '../../../components/spinner'
import toast from '../../../components/toast'
import DateTime from '../../../lib/components/date-time'
import Icon from '../../../components/icon'
import SubmitBar from '../../../components/submit-bar'

import api from '../../api'

import style from './application-link.styl'

const m = defineMessages({
  linkApplication: 'Link {appId}',
  linkStatistics: 'Statistics',
  linkStatus: 'Link status',
  linkStatusLinked: 'The application is linked successfully',
  linkStatusUnLinked: 'The application is currently not linked to a Network Server',
  linkSuccess: 'Successfully linked',
  linkedSince: 'Linked Since',
  nsAddress: 'Network Server Address',
  nsCluster: 'Network Server is within a cluster',
  nsDescription: 'Leave empty to link to the Network Server in the same cluster',
  statistics: 'Statistics',
  unlink: 'Unlink',
  unlinkSuccess: 'Successfully unlinked',
})

const validationSchema = Yup.object().shape({
  api_key: Yup.string()
    .matches(apiKey, sharedMessages.validateFormat)
    .required(sharedMessages.validateRequired),
  network_server_address: Yup.string()
    .matches(address, sharedMessages.validateFormat),
})

@connect(function (state, props) {
  const { appId } = props.match.params

  return { appId }
})
@withBreadcrumb('apps.single.link', function (props) {
  const { appId } = props

  return (
    <Breadcrumb
      path={`/console/applications/${appId}/link`}
      icon="link"
      content={sharedMessages.link}
    />
  )
})
@bind
class ApplicationLink extends React.Component {

  constructor (props) {
    super(props)

    this.form = React.createRef()
  }

  state = {
    fetching: true,
    link: undefined,
    stats: undefined,
    error: '',
  }



  componentDidMount () {
    this.fetchLinkData()
  }

  async fetchLinkData () {
    const { appId } = this.props

    await this.setState({ fetching: true })
    try {
      const link = await api.application.link.get(appId,
        [ 'api_key', 'network_server_address' ],
      )
      const stats = await api.application.link.stats(appId)

      await this.setState({
        fetching: false,
        error: '',
        link,
        stats,
      })
    } catch (error) {
      // show only non-404 errors
      if (error && error.code !== 5) {
        await this.setState({ error })
      }
    } finally {
      await this.setState({ fetching: false })
    }
  }

  async handleLink (values, { setSubmitting, resetForm }) {
    const { appId } = this.props
    const { api_key, network_server_address } = values

    await this.setState({ error: '' })
    try {
      const link = await api.application.link.set(appId, {
        api_key,
        network_server_address,
      })

      try {
        const stats = await api.application.link.stats(appId)
        await this.setState({
          link,
          stats,
        })
        resetForm(values)
        toast({
          title: appId,
          message: m.linkSuccess,
          type: toast.types.SUCCESS,
        })
      } catch (statsError) {
        throw statsError
      }
    } catch (error) {
      await this.setState({ error })
      setSubmitting(false)
    }
  }

  async handleUnlink () {
    const { appId } = this.props

    await this.setState({ error: '' })

    try {
      await api.application.link.delete(appId)
      await this.setState({
        stats: undefined,
        link: undefined,
      })
      toast({
        title: appId,
        message: m.unlinkSuccess,
        type: toast.types.SUCCESS,
      })
    } catch (error) {
      await this.setState({ error })
    } finally {
      this.form.current.resetForm({})
    }
  }

  get statistics () {
    const { stats } = this.state

    if (!stats) {
      return (
        <div className={style.status}>
          <Message component="h3" content={m.linkStatus} />
          <span className={style.statusText}><Icon icon="link_off" /> <Message content={m.linkStatusUnLinked} /></span>
        </div>
      )
    }

    const linkedAt = stats.linked_at
    const uplinkCount = stats.up_count || '0'
    const downlinkCount = stats.downlink_count || '0'

    const dataSheetItems = [
      {
        key: m.linkedSince,
        value: <DateTime.Relative value={linkedAt} />,
      },
      {
        key: sharedMessages.uplinksReceived,
        value: uplinkCount,
      },
      {
        key: sharedMessages.downlinksScheduled,
        value: downlinkCount,
      },
    ]

    return (
      <div className={style.status}>
        <Message component="h3" content={m.linkStatus} />
        <span className={style.statusText}><Icon icon="link" /> <Message content={m.linkStatusLinked} /></span>
        <DataSheet
          className={style.statusData}
          data={[{
            header: m.linkStatistics,
            items: dataSheetItems,
          }]}
        />
        <Button
          onClick={this.handleUnlink}
          message={m.unlink}
          danger
          icon="link_off"
        />
      </div>
    )
  }

  render () {
    const { appId } = this.props
    const {
      link = {},
      fetching,
      error,
    } = this.state

    if (fetching) {
      return (
        <Spinner
          center
          message={sharedMessages.loading}
        />
      )
    }

    const initialValues = {
      api_key: link.api_key || '',
      network_server_address: link.network_server_address || '',
    }

    return (
      <Container className={style.main}>
        <Row>
          <Col lg={8} md={12}>
            <IntlHelmet
              title={sharedMessages.link}
            />
            <Message
              component="h2"
              content={m.linkApplication}
              values={{ appId }}
            />
          </Col>
        </Row>
        <Row className={style.form}>
          <Col lg={6} md={12}>
            <h3>Link settings</h3>
            <Form
              formikRef={this.form}
              error={error}
              onSubmit={this.handleLink}
              initialValues={initialValues}
              validationSchema={validationSchema}
            >
              <Field
                type="text"
                description={m.nsDescription}
                name="network_server_address"
                title={m.nsAddress}
                autoFocus
              />
              <Field
                type="password"
                required
                name="api_key"
                title={sharedMessages.apiKey}
              />
              <SubmitBar>
                <Button
                  type="submit"
                  message={sharedMessages.saveChanges}
                />
              </SubmitBar>
            </Form>
          </Col>
          <Col lg={6} md={12}>
            {this.statistics}
          </Col>
        </Row>
      </Container>
    )
  }
}

export default ApplicationLink
