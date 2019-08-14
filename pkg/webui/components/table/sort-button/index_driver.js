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

import SortButton from '.'

export default function() {
  const driver = {
    component: undefined,
    when: {
      created(props) {
        driver.component = shallow(<SortButton {...props} />)
      },
      buttonPressed() {
        driver.component.simulate('click')
      },
    },
    is: {
      active() {
        return driver.component.hasClass('buttonActive')
      },
      descending() {
        return driver.component.hasClass('buttonDesc')
      },
    },
  }

  return driver
}
