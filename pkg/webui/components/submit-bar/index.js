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
import classnames from 'classnames'

import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'

import style from './submit-bar.styl'

const SubmitBar = function(props) {
  const { className, children, align } = props

  const cls = classnames(className, style.bar, style[`bar-${align}`])

  return <div className={cls}>{children}</div>
}

const SubmitBarMessage = ({ className, ...rest }) => (
  <Message {...rest} className={classnames(className, style.barMessage)} />
)

SubmitBar.propTypes = {
  align: PropTypes.oneOf(['start', 'end', 'between', 'around']),
  children: PropTypes.node,
  className: PropTypes.string,
}

SubmitBar.defaultProps = {
  align: 'between',
  className: undefined,
  children: undefined,
}

SubmitBarMessage.propTypes = {
  className: PropTypes.string,
}

SubmitBarMessage.defaultProps = {
  className: undefined,
}

SubmitBar.Message = SubmitBarMessage

export default SubmitBar
