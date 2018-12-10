// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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
import { Formik } from 'formik'

import Field from '../field'
import Button from '../button'
import Notification from '../notification'
import PropTypes from '../../lib/prop-types'

const InnerForm = function ({
  setFieldValue,
  setFieldTouched,
  handleSubmit,
  isSubmitting,
  isValid,
  errors,
  error,
  values,
  touched,
  children,
  horizontal,
  submitEnabledWhenInvalid,
  validateOnBlur,
  validateOnChange,
}) {

  const decoratedChildren = recursiveMap(children,
    function (Child) {
      if (Child.type === Field) {
        return React.cloneElement(Child, {
          setFieldValue,
          setFieldTouched,
          errors,
          values,
          touched,
          horizontal,
          submitEnabledWhenInvalid,
          validateOnBlur,
          validateOnChange,
          ...Child.props,
        })
      } else if (Child.type === Button && Child.props.type === 'submit') {
        return React.cloneElement(Child, {
          disabled: isSubmitting || (!submitEnabledWhenInvalid && !isValid),
          ...Child.props,
        })
      }

      return Child
    })

  return (
    <form onSubmit={handleSubmit}>
      {error && (<Notification small error={error} />)}
      {decoratedChildren}
    </form>
  )
}

const formRender = ({ children, ...rest }) => function (props) {
  return (
    <InnerForm
      {...rest}
      {...props}
    >
      {children}
    </InnerForm>
  )
}

const Form = ({ children, error, horizontal, submitEnabledWhenInvalid, ...rest }) => (
  <Formik {...rest} render={formRender({ children, error, horizontal, submitEnabledWhenInvalid })} />
)

function recursiveMap (children, fn) {
  return React.Children.map(children, function (child) {
    if (!React.isValidElement(child)) {
      return child
    }

    if (child.props.children) {
      return React.cloneElement(child, {
        children: recursiveMap(child.props.children, fn),
      })
    }

    return fn(child)
  })
}

Form.propTypes = {
  /** An error message belonging to the form */
  error: PropTypes.error,
  /** Whether the form fields should be displayed in horizontal style */
  horizontal: PropTypes.bool,
  /** Whether the submit button stays enabled also when the form data is not
   * not yet valid */
  submitEnabledWhenInvalid: PropTypes.bool,
}

export default Form
