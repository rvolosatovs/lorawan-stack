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

import Spinner from '@ttn-lw/components/spinner'

import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import style from './overlay.styl'

const Overlay = ({
  className,
  overlayClassName,
  spinnerClassName,
  spinnerMessage,
  visible,
  loading,
  children,
}) => (
  <div className={classnames(className, style.overlayWrapper)}>
    <div
      className={classnames(overlayClassName, style.overlay, {
        [style.overlayVisible]: visible,
      })}
    />
    {visible && loading && (
      <Spinner center className={classnames(spinnerClassName, style.overlaySpinner)}>
        <Message content={spinnerMessage} />
      </Spinner>
    )}
    {children}
  </div>
)

Overlay.propTypes = {
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
  /**
   * A flag specifying whether the overlay should display the loading spinner.
   */
  loading: PropTypes.bool,
  overlayClassName: PropTypes.string,
  spinnerClassName: PropTypes.string,
  spinnerMessage: PropTypes.message,
  /** A flag specifying whether the overlay is visible or not. */
  visible: PropTypes.bool.isRequired,
}

Overlay.defaultProps = {
  className: undefined,
  overlayClassName: undefined,
  spinnerClassName: undefined,
  spinnerMessage: sharedMessages.fetching,
  loading: false,
}

export default Overlay
