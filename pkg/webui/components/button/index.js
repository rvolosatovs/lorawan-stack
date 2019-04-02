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
      warning,
      secondary,
      naked,
      icon,
      busy,
      large,
      className,
      onClick,
      error,
      link,
      ...rest
    } = this.props

    const buttonClassNames = classnames(style.button, className, {
      [style.danger]: danger,
      [style.warning]: warning,
      [style.secondary]: secondary,
      [style.naked]: naked,
      [style.busy]: busy,
      [style.withIcon]: icon !== undefined && message,
      [style.onlyIcon]: icon !== undefined && !message,
      [style.error]: error && !busy,
      [style.large]: large,
    })

    const Elem = 'href' in rest ? 'a' : 'button'

    return (
      <Elem
        className={buttonClassNames}
        onClick={this.handleClick}
        {...rest}
      >
        <div className={style.content}>
          {icon ? <Icon className={style.icon} nudgeUp icon={icon} /> : null}
          {busy ? <Spinner className={style.spinner} small after={200} /> : null}
          {message ? <Message content={message} /> : null}
        </div>
      </Elem>
    )
  }
}

Button.propTypes = {
  /** The message to be displayed within the button */
  message: PropTypes.message,
  /**
   * A click listener to be called when the button is pressed.
   * Not called if the button is in the `busy` or `disabled` state.
   */
  onClick: PropTypes.func,
  /**
   * A flag specifying whether the `danger` styling should applied to the button
   */
  danger: PropTypes.bool,
  /**
   * A flag specifying whether the `warning` styling should applied to the button
   */
  warning: PropTypes.bool,
  /**
   * A flag specifying whether the `secodnary` styling should applied to the button
   */
  secondary: PropTypes.bool,
  /**
   * A flag specifying whether the `naked` styling should applied to the button
   */
  naked: PropTypes.bool,
  /**
   * A flag specifying whether the `large` styling should applied to the button
   */
  large: PropTypes.bool,
  /**
   * A flag specifying whether the `error` styling should applied to the button
   */
  error: PropTypes.bool,
  /** The name of an icon to be displayed within the button*/
  icon: PropTypes.string,
  /**
   * A flag specifying whether the button in the `busy` state and the appropriate
   * styling should be applied.
   */
  busy: PropTypes.bool,
  /**
   * A flag specifying whether the button in the `disabled` state and the appropriate
   * styling should be applied.
   */
  disabled: PropTypes.bool,
}

Button.defaultProps = {
  onClick: () => null,
}

export default Button
