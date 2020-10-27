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
import classnames from 'classnames'

import Button from '@ttn-lw/components/button'

import Logo from '@ttn-lw/containers/logo'

import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import style from './modal.styl'

class Modal extends React.PureComponent {
  static propTypes = {
    approval: PropTypes.bool,
    bottomLine: PropTypes.oneOfType([PropTypes.element, PropTypes.message]),
    buttonMessage: PropTypes.message,
    buttonName: PropTypes.message,
    cancelButtonMessage: PropTypes.message,
    children: PropTypes.oneOfType([PropTypes.arrayOf(PropTypes.element), PropTypes.element]),
    danger: PropTypes.bool,
    formName: PropTypes.string,
    inline: PropTypes.bool,
    logo: PropTypes.bool,
    message: PropTypes.message,
    method: PropTypes.string,
    name: PropTypes.string,
    onComplete: PropTypes.func,
    subtitle: PropTypes.message,
    title: PropTypes.message,
  }

  static defaultProps = {
    bottomLine: undefined,
    buttonMessage: undefined,
    buttonName: undefined,
    cancelButtonMessage: sharedMessages.cancel,
    children: undefined,
    danger: false,
    formName: undefined,
    logo: false,
    message: undefined,
    method: undefined,
    onComplete: () => null,
    inline: false,
    approval: true,
    subtitle: undefined,
    title: undefined,
    name: undefined,
  }

  @bind
  handleApprove() {
    this.handleComplete(true)
  }

  @bind
  handleCancel() {
    this.handleComplete(false)
  }

  @bind
  handleComplete(result) {
    const { onComplete } = this.props

    onComplete(result)
  }

  render() {
    const {
      buttonName,
      buttonMessage,
      title,
      subtitle,
      children,
      message,
      logo,
      approval,
      formName,
      cancelButtonMessage,
      onComplete,
      bottomLine,
      inline,
      danger,
      ...rest
    } = this.props

    const modalClassNames = classnames(style.modal, style.modal, {
      [style.modalInline]: inline,
      [style.modalAbsolute]: !Boolean(inline),
    })

    const name = formName ? { name: formName } : {}
    const RootComponent = this.props.method ? 'form' : 'div'
    const messageElement = message && <Message content={message} className={style.message} />
    const bottomLineElement =
      typeof bottomLine === 'object' && Boolean(bottomLine.id) ? (
        <Message content={bottomLine} />
      ) : (
        bottomLine
      )

    const approveButtonMessage =
      buttonMessage !== undefined
        ? buttonMessage
        : approval
        ? sharedMessages.approve
        : sharedMessages.ok
    let buttons = (
      <div>
        <Button message={approveButtonMessage} onClick={this.handleApprove} icon="check" />
      </div>
    )

    if (approval) {
      buttons = (
        <div>
          <Button
            secondary
            message={cancelButtonMessage}
            onClick={this.handleCancel}
            name={formName}
            icon="clear"
            value="false"
            {...name}
          />
          <Button
            message={approveButtonMessage}
            onClick={this.handleApprove}
            name={formName}
            icon="check"
            value="true"
            danger={danger}
            {...name}
          />
        </div>
      )
    }

    return (
      <React.Fragment>
        {!inline && <div key="shadow" className={style.shadow} />}
        <RootComponent
          data-test-id="modal-window"
          key="modal"
          className={modalClassNames}
          {...rest}
        >
          {title && (
            <div className={style.titleSection}>
              <div>
                <h1>
                  <Message content={title} />
                </h1>
                {subtitle && <Message component="span" content={subtitle} />}
              </div>
              {logo && <Logo vertical className={style.logo} />}
            </div>
          )}
          {title && <div className={style.line} />}
          <div className={style.body}>{children || messageElement}</div>
          <div className={style.controlBar}>
            <div>{bottomLineElement}</div>
            {buttons}
          </div>
        </RootComponent>
      </React.Fragment>
    )
  }
}

export default Modal
