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
import { connect as storeConnect } from 'react-redux'
import bind from 'autobind-decorator'

import Field from '@ttn-lw/components/form/field'
import Select from '@ttn-lw/components/select'

import PropTypes from '@ttn-lw/lib/prop-types'

const formatOptions = options =>
  Object.keys(options).map(key => ({ value: key, label: options[key] }))

const { component, ...fieldPropTypes } = Field.propTypes

export default function({
  optionsSelector,
  errorSelector,
  fetchingSelector,
  fetchOptions,
  defaultWarning,
  defaultTitle,
  optionsFormatter = formatOptions,
  defaultDescription,
}) {
  @storeConnect(
    function(state) {
      return {
        options: optionsFormatter(optionsSelector(state)),
        error: errorSelector(state),
        fetching: fetchingSelector(state),
      }
    },
    { fetchOptions },
  )
  class FetchSelect extends React.PureComponent {
    static propTypes = {
      ...fieldPropTypes,
      ...Select.propTypes,
      description: PropTypes.message,
      fetchOptions: PropTypes.func.isRequired,
      menuPlacement: PropTypes.oneOf(['top', 'bottom', 'auto']),
      onChange: PropTypes.func,
      options: PropTypes.arrayOf(
        PropTypes.shape({ value: PropTypes.string, label: PropTypes.message }),
      ),
      title: PropTypes.message,
      warning: PropTypes.message,
    }

    static defaultProps = {
      description: defaultDescription,
      menuPlacement: 'auto',
      onChange: () => null,
      options: [],
      title: defaultTitle,
      warning: defaultWarning,
    }

    componentDidMount() {
      const { fetchOptions } = this.props

      fetchOptions()
    }

    @bind
    handleChange(value) {
      const { onChange, options } = this.props

      onChange(options.find(e => e.value === value))
    }

    render() {
      const { error, fetching, warning, ...rest } = this.props

      return (
        <Field
          {...rest}
          component={Select}
          isLoading={fetching}
          warning={Boolean(error) ? warning : undefined}
          onChange={this.handleChange}
        />
      )
    }
  }

  return FetchSelect
}
