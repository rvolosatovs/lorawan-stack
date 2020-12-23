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

import React from 'react'
import classnames from 'classnames'
import bind from 'autobind-decorator'

import { unit as unitRegexp } from '@ttn-lw/console/lib/regexp'

import Select from '@ttn-lw/components/select'
import Input from '@ttn-lw/components/input'

import withComputedProps from '@ttn-lw/lib/components/with-computed-props'

import PropTypes from '@ttn-lw/lib/prop-types'

import style from './unit-input.styl'

class UnitInput extends React.PureComponent {
  static propTypes = {
    className: PropTypes.string,
    disabled: PropTypes.bool,
    encode: PropTypes.func,
    error: PropTypes.bool,
    inputWidth: PropTypes.inputWidth,
    name: PropTypes.string.isRequired,
    onBlur: PropTypes.func.isRequired,
    onChange: PropTypes.func.isRequired,
    required: PropTypes.bool,
    units: PropTypes.arrayOf(
      PropTypes.shape({
        label: PropTypes.message,
        value: PropTypes.string,
      }),
    ).isRequired,
    value: PropTypes.shape({
      unit: PropTypes.string,
      duration: PropTypes.number,
    }),
  }

  static defaultProps = {
    className: undefined,
    encode: (duration, unit) => (duration ? `${duration}${unit}` : unit),
    inputWidth: 'm',
    disabled: false,
    required: false,
    value: undefined,
    error: false,
  }

  @bind
  async handleChange(duration) {
    const { onChange, encode, value } = this.props
    onChange(encode(duration, value.unit))
  }

  @bind
  async handleUnitChange(unit) {
    const { onChange, encode, value } = this.props
    onChange(encode(value.duration, unit), true)
  }

  render() {
    const {
      className,
      name,
      units,
      value: { duration, unit },
      onBlur,
      required,
      disabled,
      error,
      inputWidth,
    } = this.props

    const selectTimeUnitComponent = (
      <Select
        className={style.select}
        name={`${name}-select`}
        options={units}
        onChange={this.handleUnitChange}
        onBlur={onBlur}
        value={unit}
        disabled={disabled}
        inputWidth="full"
      />
    )

    return (
      <React.Fragment>
        <div className={classnames(className, style.container, style[`input-width-${inputWidth}`])}>
          <Input
            data-test-id={name}
            className={style.number}
            type="number"
            step="any"
            name={name}
            onBlur={onBlur}
            value={duration}
            onChange={this.handleChange}
            required={required}
            disabled={disabled}
            error={error}
            inputWidth="full"
          />
          {selectTimeUnitComponent}
        </div>
      </React.Fragment>
    )
  }
}

const WrappedUnitInput = withComputedProps(props => ({
  ...props,
  value: props.decode(props.value),
}))(UnitInput)

WrappedUnitInput.propTypes = {
  decode: PropTypes.func,
}

WrappedUnitInput.defaultProps = {
  decode: value => {
    const duration = value.split(unitRegexp)[0]
    const unit = value.split(duration)[1]
    return {
      duration: duration ? Number(duration) : undefined,
      unit,
    }
  },
}

export default WrappedUnitInput
