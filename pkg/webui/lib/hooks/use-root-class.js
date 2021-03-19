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

import { useEffect } from 'react'

export default (className, id = 'app') => {
  useEffect(() => {
    const containerClasses = document.getElementById(id).classList
    const classNamesList = className.split(' ')
    for (const cls of classNamesList) {
      containerClasses.add(cls)
    }
    return () => {
      for (const cls of classNamesList) {
        containerClasses.remove(cls)
      }
    }
    // Disabling deps check because in this case we want the effect hook to
    // run only once on component mount. See also:
    // https://github.com/facebook/create-react-app/issues/6880
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])
}
