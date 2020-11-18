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
import classnames from 'classnames'
import { getIn } from 'formik'

import Icon from '@ttn-lw/components/icon'
import Link from '@ttn-lw/components/link'

import Message from '@ttn-lw/lib/components/message'

import from from '@ttn-lw/lib/from'
import PropTypes from '@ttn-lw/lib/prop-types'

import FormContext from '../context'

import style from './field.styl'

export function getPassThroughProps(props, excludeProps) {
  const rest = {}
  for (const property of Object.keys(props)) {
    if (!excludeProps[property]) {
      rest[property] = props[property]
    }
  }
  return rest
}

const isValueEmpty = function(value) {
  if (value === null || value === undefined) {
    return true
  }

  if (typeof value === 'object') {
    return Object.keys(value) === 0
  }

  if (typeof value === 'string') {
    return value === ''
  }

  return false
}

class FormField extends React.Component {
  static contextType = FormContext

  static propTypes = {
    className: PropTypes.string,
    component: PropTypes.oneOfType([
      PropTypes.func,
      PropTypes.string,
      PropTypes.shape({
        render: PropTypes.func.isRequired,
      }),
    ]).isRequired,
    decode: PropTypes.func,
    description: PropTypes.message,
    disabled: PropTypes.bool,
    encode: PropTypes.func,
    glossaryId: PropTypes.string,
    glossaryTerm: PropTypes.message,
    name: PropTypes.string.isRequired,
    onChange: PropTypes.func,
    readOnly: PropTypes.bool,
    required: PropTypes.bool,
    title: PropTypes.message.isRequired,
    warning: PropTypes.message,
  }

  static defaultProps = {
    className: undefined,
    disabled: false,
    encode: value => value,
    decode: value => value,
    onChange: () => null,
    warning: '',
    description: '',
    glossaryTerm: '',
    glossaryId: '',
    readOnly: false,
    required: false,
  }

  componentDidMount() {
    const { name } = this.props

    this.context.registerField(name, this)
  }

  componentWillUnmount() {
    const { name } = this.props

    this.context.unregisterField(name)
  }

  extractValue(value) {
    let newValue = value
    if (typeof value === 'object' && 'target' in value) {
      const target = value.target
      if ('type' in target && target.type === 'checkbox') {
        newValue = target.checked
      } else if ('value' in target) {
        newValue = target.value
      }
    }

    return newValue
  }

  @bind
  async handleChange(value, enforceValidation = false) {
    const { name, onChange, encode } = this.props
    const { setFieldValue, setFieldTouched } = this.context

    // Check if the value is react's synthetic event.
    const isSyntheticEvent = typeof value === 'object' && 'target' in value
    const newValue = encode(this.extractValue(value))

    // TODO: Remove `await` and event persist when https://github.com/jaredpalmer/formik/issues/2457
    // is resolved.
    if (typeof value !== 'undefined' && typeof value.persist === 'function') {
      value.persist()
    }

    await setFieldValue(name, newValue)

    if (enforceValidation) {
      setFieldTouched(name)
    }

    onChange(isSyntheticEvent ? value : encode(value))
  }

  @bind
  handleBlur(event) {
    const { name } = this.props
    const { validateOnBlur, setFieldTouched } = this.context

    if (validateOnBlur) {
      const value = this.extractValue(event)
      setFieldTouched(name, !isValueEmpty(value))
    }
  }

  render() {
    const {
      className,
      decode,
      name,
      title,
      warning,
      description,
      disabled,
      required,
      readOnly,
      glossaryTerm,
      glossaryId,
      component: Component,
    } = this.props
    const { horizontal, disabled: formDisabled } = this.context

    const fieldValue = decode(getIn(this.context.values, name))
    const fieldError = getIn(this.context.errors, name)
    const fieldTouched = getIn(this.context.touched, name) || false
    const fieldDisabled = disabled || formDisabled

    const hasError = Boolean(fieldError)
    const hasWarning = Boolean(warning)
    const hasDescription = Boolean(description)
    const hasGlossaryTerm = Boolean(glossaryId)

    const showError = fieldTouched && hasError
    const showWarning = !hasError && hasWarning
    const showDescription = !showError && !showWarning && hasDescription

    const describedBy = showError
      ? `${name}-field-error`
      : showWarning
      ? `${name}-field-warning`
      : showDescription
      ? `${name}-field-description`
      : undefined

    const fieldMessage = showError ? (
      <div className={style.messages}>
        <Err content={fieldError} title={title} error id={describedBy} />
      </div>
    ) : showWarning ? (
      <div className={style.messages}>
        <Err content={warning} title={title} warning id={describedBy} />
      </div>
    ) : showDescription ? (
      <Message className={style.description} content={description} id={describedBy} />
    ) : null

    const fieldComponentProps = {
      value: fieldValue,
      error: showError,
      warning: showWarning,
      name,
      id: name,
      horizontal,
      disabled: fieldDisabled,
      onChange: this.handleChange,
      onBlur: this.handleBlur,
    }

    const cls = classnames(
      className,
      style.field,
      from(style, {
        error: showError,
        warning: showWarning,
        horizontal,
        required,
        readOnly,
      }),
    )

    return (
      <div className={cls} data-needs-focus={showError}>
        <div className={style.label}>
          <div className={style.labelContainer}>
            <Message
              component="label"
              content={title}
              className={style.title}
              htmlFor={fieldComponentProps.id}
            />
            {hasGlossaryTerm && (
              <Link.GlossaryLink
                hideTerm
                secondary
                term={glossaryTerm || title}
                glossaryId={glossaryId}
                className={style.glossaryLink}
              />
            )}
          </div>
        </div>
        <div className={style.componentArea}>
          <Component
            aria-invalid={showError}
            aria-describedby={describedBy}
            {...fieldComponentProps}
            {...getPassThroughProps(this.props, FormField.propTypes)}
          />
          {fieldMessage}
        </div>
      </div>
    )
  }
}

const Err = ({ content, error, warning, title, className, id }) => {
  const icon = error ? 'error' : 'warning'
  const contentValues = content.values || {}
  const classname = classnames(style.message, className, {
    [style.show]: content && content !== '',
    [style.hide]: !content || content === '',
    [style.err]: error,
    [style.warn]: warning,
  })

  if (title) {
    contentValues.field = <Message content={title} className={style.name} key={title.id || title} />
  }

  return (
    <div className={classname} id={id}>
      <Icon icon={icon} className={style.icon} />
      <Message content={content.message || content} values={contentValues} />
    </div>
  )
}

Err.propTypes = {
  className: PropTypes.string,
  content: PropTypes.oneOfType([
    PropTypes.error,
    PropTypes.shape({
      message: PropTypes.error.isRequired,
      values: PropTypes.shape({}).isRequired,
    }),
  ]).isRequired,
  error: PropTypes.bool,
  id: PropTypes.string.isRequired,
  title: PropTypes.message,
  warning: PropTypes.bool,
}

Err.defaultProps = {
  className: undefined,
  title: undefined,
  warning: false,
  error: false,
}

export default FormField
