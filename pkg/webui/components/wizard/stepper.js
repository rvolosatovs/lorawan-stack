// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

import PropTypes from '@ttn-lw/lib/prop-types'

import Stepper from '../stepper'

import { useWizardContext } from './context'

const WizardStepper = props => {
  const { currentStepId, steps, error } = useWizardContext()
  const { children, ...rest } = props

  const currentStepIndex = steps.findIndex(({ id }) => id === currentStepId)

  if (steps.length <= 1 || currentStepIndex < 0) {
    return null
  }

  const status = Boolean(error) ? 'failure' : 'current'

  return (
    <Stepper {...rest} status={status} currentStep={currentStepIndex + 1}>
      {React.Children.map(children, (child, index) => {
        if (React.isValidElement(child) && child.type.displayName === 'Stepper.Step') {
          return React.cloneElement(child, {
            ...child.props,
            stepNumber: index + 1,
          })
        }

        return child
      })}
    </Stepper>
  )
}

WizardStepper.propTypes = {
  children: PropTypes.oneOfType([PropTypes.node, PropTypes.arrayOf(PropTypes.node)]),
}

WizardStepper.defaultProps = {
  children: [],
}

WizardStepper.Step = Stepper.Step
WizardStepper.displayName = 'Wizard.Stepper'

export default WizardStepper
