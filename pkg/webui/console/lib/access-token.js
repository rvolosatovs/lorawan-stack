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

import * as cache from './cache'

export default retrieveToken => async () => {
  const storedToken = cache.get('accessToken')
  let token

  if (!storedToken || Date.parse(storedToken.expiry) < Date.now()) {
    // If we don't have a token stored or it's expired, we want to retrieve it.

    // Remove stored, invalid token.
    clear()

    // Retrieve new token and store it.
    const response = await retrieveToken()
    token = response.data
    cache.set('accessToken', token)
    return token
  }

  // If we have a stored token and its valid, we want to use it.
  return storedToken
}

export function clear() {
  cache.remove('accessToken')
}
