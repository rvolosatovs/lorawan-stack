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
import classnames from 'classnames'
import bind from 'autobind-decorator'

import PropTypes from '../../lib/prop-types'
import Spinner from '../spinner'
import Message from '../../lib/components/message'
import Icon from '../icon'

import style from './button.styl'

@bind
class Button extends React.PureComponent {

  handleClick (evt) {
    const { busy, disabled, onClick } = this.props

    if (busy || disabled) {
      return
    }

    onClick(evt)
  }

  render () {
    const {
      message,
      danger,
      secondary,
      naked,
      icon,
      busy,
      large,
      className,
      onClick,
      error,
      ...rest
    } = this.props

    const buttonClassNames = classnames(style.button, className, {
      [style.danger]: danger,
      [style.secondary]: secondary,
      [style.naked]: naked,
      [style.busy]: busy,
      [style.withIcon]: icon !== undefined && message,
      [style.onlyIcon]: icon !== undefined && !message,
      [style.error]: error && !busy,
      [style.large]: large,
    })

    return (
      <button
        className={buttonClassNames}
        onClick={this.handleClick}
        {...rest}
      >
        <div className={style.content}>
          {icon ? <Icon className={style.icon} nudgeUp icon={icon} /> : null}
          {busy ? <Spinner className={style.spinner} small after={200} /> : null}
          {message ? <Message content={message} /> : null}
        </div>
      </button>
    )
  }
}

Button.propTypes = {
  message: PropTypes.message,
  onClick: PropTypes.func,
  danger: PropTypes.bool,
  boring: PropTypes.bool,
  busy: PropTypes.bool,
}

Button.defaultProps = {
  onClick: () => null,
}

export default Button
