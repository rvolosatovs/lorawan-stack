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

import PropTypes from '@ttn-lw/lib/prop-types'
import renderCallback from '@ttn-lw/lib/render-callback'

import { useWizardContext } from './context'

const Step = props => {
  const context = useWizardContext()

  return renderCallback(props, context)
}

Step.propTypes = {
  children: PropTypes.oneOfType([PropTypes.node, PropTypes.arrayOf(PropTypes.node)]),
  id: PropTypes.string.isRequired,
  render: PropTypes.func,
}

Step.defaultProps = {
  render: undefined,
  children: [],
}

export default Step
