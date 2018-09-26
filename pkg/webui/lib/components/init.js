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
import { connect } from 'react-redux'
import 'focus-visible/dist/focus-visible'

import Spinner from '../../components/spinner'

import '../../styles/main.styl'

@connect(state => (
  {
    initialized: state.console.initialized,
  }
))
export default class Init extends React.PureComponent {

  componentDidMount () {
    this.props.dispatch({ type: 'INITIALIZE' })
  }

  render () {
    const { initialized } = this.props

    if (!initialized) {
      return (<Spinner center>Please wait…</Spinner>)
    }

    return this.props.children
  }
}
