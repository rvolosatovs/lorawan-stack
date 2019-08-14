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

import FormContext from '../context'
import PropTypes from '../../../lib/prop-types'

class FormSubmit extends React.Component {
  static contextType = FormContext
  static propTypes = {
    component: PropTypes.oneOfType([PropTypes.func, PropTypes.string]).isRequired,
  }

  static defaultProps = {
    component: 'button',
  }

  render() {
    const { component: Component, ...rest } = this.props

    const submitProps = {
      isValid: this.context.isValid,
      isSubmitting: this.context.isSubmitting,
      isValidating: this.context.isValidating,
      submitCount: this.context.submitCount,
      dirty: this.context.dirty,
      validateForm: this.context.validateForm,
      validateField: this.context.validateField,
      disabled: this.context.disabled,
    }

    return <Component {...rest} {...submitProps} />
  }
}

export default FormSubmit
