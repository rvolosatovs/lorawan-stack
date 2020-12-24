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

import SubmitBar from '@ttn-lw/components/submit-bar'
import FormField from '@ttn-lw/components/form/field'
import FormSubmit from '@ttn-lw/components/form/submit'
import SubmitButton from '@ttn-lw/components/submit-button'
import Input from '@ttn-lw/components/input'

import ApiKeyModal from '@console/components/api-key-modal'
import RightsGroup from '@console/components/rights-group'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import ApiKeyForm from './form'
import validationSchema from './validation-schema'

class CreateForm extends React.Component {
  static propTypes = {
    /** Called on form submission. Receives the key object as an argument. */
    onCreate: PropTypes.func.isRequired,
    /**
     * Called after unsuccessful creation of the API key. Receives the error
     * object as an argument.
     */
    onCreateFailure: PropTypes.func,
    /**
     * Called after successful creation of the API key. Receives the key object
     * as an argument.
     */
    onCreateSuccess: PropTypes.func.isRequired,
    /**
     * The rights that imply all other rights, e.g. 'RIGHT_APPLICATION_ALL',
     * 'RIGHT_ALL'.
     */
    pseudoRights: PropTypes.rights,
    /** The list of rights. */
    rights: PropTypes.rights,
  }

  state = {
    modal: null,
  }

  static defaultProps = {
    onCreateFailure: () => null,
    pseudoRights: [],
    rights: [],
  }

  @bind
  async handleModalApprove() {
    const { onCreateSuccess } = this.props
    const { key } = this.state

    await this.setState({ modal: null })
    await onCreateSuccess(key)
  }

  @bind
  async handleCreate(values) {
    const { onCreate } = this.props

    return await onCreate(values)
  }

  @bind
  async handleCreateSuccess(key) {
    await this.setState({
      modal: {
        secret: key.key,
        rights: key.rights,
        onComplete: this.handleModalApprove,
        approval: false,
      },
      key,
    })
  }

  render() {
    const { rights, onCreateFailure, pseudoRights } = this.props
    const { modal } = this.state

    const modalProps = modal ? modal : {}
    const modalVisible = Boolean(modal)
    const initialValues = {
      name: '',
      rights: [...pseudoRights],
    }

    return (
      <React.Fragment>
        <ApiKeyModal {...modalProps} visible={modalVisible} approval={false} />
        <ApiKeyForm
          rights={rights}
          onSubmit={this.handleCreate}
          onSubmitSuccess={this.handleCreateSuccess}
          onSubmitFailure={onCreateFailure}
          validationSchema={validationSchema}
          initialValues={initialValues}
        >
          <FormField
            title={sharedMessages.name}
            placeholder={sharedMessages.apiKeyNamePlaceholder}
            name="name"
            autoFocus
            component={Input}
          />
          <FormField
            name="rights"
            title={sharedMessages.rights}
            required
            component={RightsGroup}
            rights={rights}
            pseudoRight={pseudoRights[0]}
            entityTypeMessage={sharedMessages.apiKey}
          />
          <SubmitBar>
            <FormSubmit component={SubmitButton} message={sharedMessages.createApiKey} />
          </SubmitBar>
        </ApiKeyForm>
      </React.Fragment>
    )
  }
}

export default CreateForm
