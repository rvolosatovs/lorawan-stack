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
import PropTypes from 'prop-types'
import classnames from 'classnames'
import bind from 'autobind-decorator'
import MaskedInput from 'react-text-mask'

import style from './input.styl'

const hex = /[0-9a-f]/i

const masks = {}
const mask = function (min = 0, max = 256) {
  const key = `${min}-${max}`
  if (masks[key]) {
    return masks[key]
  }

  const r = new Array(3 * max - 1).fill(hex)
  for (let i = 0; i < r.length; i++) {
    if ((i + 1) % 3 === 0) {
      r[i] = ' '
    }
  }

  return r
}

const upper = function (str) {
  return str.toUpperCase()
}

const clean = function (str) {
  return str.replace(/[ \u2000]/g, '')
}

const Placeholder = function (props) {
  const {
    min = 0,
    max = 256,
    value = '',
    placeholder,
  } = props

  if (placeholder) {
    return null
  }

  const len = 1.5 * value.length - (value.length - 2 * Math.floor(value.length / 2))

  const content = mask(min, max)
    .map(function (el, i) {
      if (!(el instanceof RegExp)) {
        return ' '
      }

      if (i < len) {
        return ' '
      }

      return '·'
    })
    .join('')

  return (
    <div
      className={style.placeholder}
    >
      {content}
    </div>
  )
}

@bind
export default class ByteInput extends React.Component {

  static propTypes = {
    value: PropTypes.string,
    onChange: PropTypes.func,
    min: PropTypes.number,
    max: PropTypes.number,
  }

  static validate (value, props) {
    const { min = 0, max = 256 } = props
    const len = Math.floor(value.length / 2)
    return min <= len && len <= max
  }

  render () {
    const {
      value,
      className,
      min = 0,
      max = 255,
      onChange,
      valid,
      placeholder,
      type,
      ...rest
    } = this.props

    return [
      <Placeholder
        key="placeholder"
        min={min}
        max={max}
        value={value}
        placeholder={placeholder}
      />,
      <MaskedInput
        key="input"
        className={classnames(className, style.byte)}
        value={value}
        mask={mask(min, max)}
        placeholderChar={'\u2000'}
        keepCharPositions={false}
        pipe={upper}
        onChange={this.onChange}
        placeholder={placeholder}
        onCopy={this.onCopy}
        onCut={this.onCut}
        {...rest}
      />,
    ]
  }

  onChange (evt) {
    this.props.onChange({
      target: {
        value: clean(evt.target.value),
      },
    })
  }

  onCopy (evt) {
    const input = evt.target
    const value = input.value.substr(input.selectionStart, input.selectionEnd - input.selectionStart)
    evt.clipboardData.setData('text/plain', clean(value))
    evt.preventDefault()
  }

  onCut (evt) {
    const input = evt.target
    const value = input.value.substr(input.selectionStart, input.selectionEnd - input.selectionStart)
    evt.clipboardData.setData('text/plain', clean(value))
    evt.preventDefault()

    // emit the cut value
    const cut = input.value.substr(0, input.selectionStart) + input.value.substr(input.selectionEnd)
    evt.target.value = cut
    this.onChange({
      target: {
        value: cut,
      },
    })
  }
}
