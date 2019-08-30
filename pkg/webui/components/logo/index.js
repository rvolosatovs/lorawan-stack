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

import PropTypes from '../../lib/prop-types'
import TtsLogo from '../../assets/logos/tts.svg'

import style from './logo.styl'

const Logo = function(props) {
  const { className } = props

  const classname = classnames(style.logo, className)

  return (
    <div className={classname}>
      <img alt="The Things Stack Logo" src={TtsLogo} />
    </div>
  )
}

Logo.propTypes = {
  className: PropTypes.string,
}

Logo.defaultProps = {
  className: undefined,
}

export default Logo
