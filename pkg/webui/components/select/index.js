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
import ReactSelect, { components } from 'react-select'
import { injectIntl } from 'react-intl'
import bind from 'autobind-decorator'
import classnames from 'classnames'

import PropTypes from '@ttn-lw/lib/prop-types'

import style from './select.styl'

const Input = props => {
  const { selectProps } = props

  return <components.Input {...props} aria-describedby={selectProps['aria-describedby']} />
}

Input.propTypes = {
  selectProps: PropTypes.shape({
    'aria-describedby': PropTypes.string,
  }).isRequired,
}

// Map value to a plain string, instead of value object.
// See: https://github.com/JedWatson/react-select/issues/2841
const getValue = (opts, val) => opts.find(o => o.value === val)

class Select extends React.PureComponent {
  static propTypes = {
    className: PropTypes.string,
    disabled: PropTypes.bool,
    error: PropTypes.bool,
    id: PropTypes.string,
    intl: PropTypes.shape({
      formatMessage: PropTypes.func,
    }).isRequired,
    name: PropTypes.string.isRequired,
    onBlur: PropTypes.func,
    onChange: PropTypes.func,
    onFocus: PropTypes.func,
    options: PropTypes.arrayOf(
      PropTypes.shape({
        value: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        label: PropTypes.message,
      }),
    ),
    value: PropTypes.string,
    warning: PropTypes.bool,
  }

  static defaultProps = {
    className: undefined,
    onChange: () => null,
    onBlur: () => null,
    onFocus: () => null,
    options: [],
    disabled: false,
    error: false,
    warning: false,
    value: undefined,
    id: undefined,
  }

  constructor(props) {
    super(props)

    let value
    if ('value' in props && this.context) {
      value = props.value
    }

    this.state = {
      checked: value,
    }
  }

  static getDerivedStateFromProps(props, state) {
    const { value } = props

    if ('value' in props && value !== state.value) {
      return { value }
    }

    return null
  }

  @bind
  async onChange({ value }) {
    const { onChange } = this.props

    if (!('value' in this.props)) {
      await this.setState({ value })
    }

    onChange(value)
  }

  @bind
  onBlur(event) {
    const { value } = this.state
    const { onBlur, name } = this.props

    // https://github.com/JedWatson/react-select/issues/3523
    // Make sure the input name is always present in the event object.
    event.target.name = name

    if (typeof value !== 'undefined') {
      // https://github.com/JedWatson/react-select/issues/3175
      event.target.value = value
    }

    onBlur(event)
  }

  render() {
    const {
      className,
      options,
      intl,
      value,
      onChange,
      onBlur,
      onFocus,
      disabled,
      error,
      warning,
      name,
      id,
      ...rest
    } = this.props

    const formatMessage = (label, values) => (intl ? intl.formatMessage(label, values) : label)
    const cls = classnames(className, style.container, {
      [style.error]: error,
      [style.warning]: warning,
    })
    const translatedOptions = options.map(function(option) {
      const { label, labelValues = {} } = option
      if (typeof label === 'object' && label.id && label.defaultMessage) {
        return { ...option, label: formatMessage(label, labelValues) }
      }

      return option
    })

    return (
      <ReactSelect
        className={cls}
        inputId={id}
        classNamePrefix="select"
        value={getValue(translatedOptions, value)}
        options={translatedOptions}
        onChange={this.onChange}
        onBlur={this.onBlur}
        onFocus={onFocus}
        isDisabled={disabled}
        name={name}
        components={{ Input }}
        aria-describedby={rest['aria-describedby']}
        {...rest}
      />
    )
  }
}

export default injectIntl(Select)
export { Select }
