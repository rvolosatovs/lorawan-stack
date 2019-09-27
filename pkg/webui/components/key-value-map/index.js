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
import { defineMessages } from 'react-intl'
import classnames from 'classnames'

import PropTypes from '../../lib/prop-types'
import Button from '../button'
import Entry from './entry'

import style from './key-value-map.styl'

const m = defineMessages({
  addEntry: 'Add entry',
})

@bind
class KeyValueMap extends React.PureComponent {
  handleEntryChange(index, newValues) {
    const { onChange, value } = this.props
    onChange(value.map((kv, i) => (index !== i ? kv : { ...kv, ...newValues })))
  }

  removeEntry(index) {
    const { onChange, value } = this.props
    onChange(value.filter((_, i) => i !== index) || [], true)
  }

  addEmptyEntry() {
    const { onChange, value } = this.props
    onChange([...value, { key: '', value: '' }])
  }

  render() {
    const {
      className,
      name,
      value,
      keyPlaceholder,
      valuePlaceholder,
      addMessage,
      onBlur,
    } = this.props

    return (
      <div className={classnames(className, style.container)}>
        <div>
          {value &&
            value.map((value, index) => (
              <Entry
                key={`${name}[${index}]`}
                name={name}
                value={value}
                keyPlaceholder={keyPlaceholder}
                valuePlaceholder={valuePlaceholder}
                index={index}
                onRemoveButtonClick={this.removeEntry}
                onChange={this.handleEntryChange}
                onBlur={onBlur}
              />
            ))}
        </div>
        <div>
          <Button
            name={`${name}.push`}
            type="button"
            message={addMessage}
            onClick={this.addEmptyEntry}
            icon="add"
          />
        </div>
      </div>
    )
  }
}

KeyValueMap.propTypes = {
  className: PropTypes.string,
  name: PropTypes.string.isRequired,
  value: PropTypes.array,
  keyPlaceholder: PropTypes.message.isRequired,
  valuePlaceholder: PropTypes.message.isRequired,
  addMessage: PropTypes.message,
  onChange: PropTypes.func,
  onBlur: PropTypes.func,
}

KeyValueMap.defaultProps = {
  value: [],
  addMessage: m.addEntry,
}

export default KeyValueMap
