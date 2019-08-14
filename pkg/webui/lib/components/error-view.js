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

import PropTypes from '../../lib/prop-types'

class ErrorView extends React.Component {
  state = {
    error: undefined,
    hasCaught: false,
  }

  componentDidCatch(error) {
    this.setState({
      hasCaught: true,
      error,
    })
  }

  render() {
    const { children, ErrorComponent } = this.props
    const { hasCaught, error } = this.state

    if (hasCaught) {
      return <ErrorComponent error={error} />
    }

    return React.Children.only(children)
  }
}

ErrorView.propTypes = {
  ErrorComponent: PropTypes.func.isRequired,
}

export default ErrorView
