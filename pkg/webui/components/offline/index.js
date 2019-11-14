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

import React, { Component } from 'react'
import PropTypes from '../../lib/prop-types'
import Status from '../status'
import Message from '../../lib/components/message'

import style from './offline.styl'

export default class Offline extends Component {
  static propTypes = {
    content: PropTypes.string,
    status: PropTypes.string,
  }

  static defaultProps = {
    content: undefined,
    status: undefined,
  }

  render() {
    const { status, content } = this.props
    return (
      <span>
        <Status className={style.status} status={status}>
          <Message className={style.message} content={content} />
        </Status>
      </span>
    )
  }
}
