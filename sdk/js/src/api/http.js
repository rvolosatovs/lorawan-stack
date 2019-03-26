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

import axios from 'axios'

/**
 * Http Class is a connector for the API that uses the HTTP bridge to connect.
 */
class Http {
  constructor (token, stackConfig, axiosConfig = {}) {
    const headers = axiosConfig.headers || {}
    let Authorization = null

    if (typeof token === 'string') {
      Authorization = `Bearer ${token}`
    }

    const stackComponents = Object.keys(stackConfig)
    const instances = stackComponents.reduce(function (acc, curr) {
      acc[curr] = axios.create({
        baseURL: stackConfig[curr],
        headers: {
          Authorization,
          ...headers,
        },
        ...axiosConfig,
      })

      return acc
    }, {})

    for (const instance in instances) {
      this[instance] = instances[instance]

      // Re-evaluate headers on each request if token is a thunk. This can be
      // useful if the token needs to be refreshed frequently, as the case for
      // access tokens.
      if (typeof token === 'function') {
        this[instance].interceptors.request.use(async function (config) {
          const tkn = (await token()).access_token
          config.headers.Authorization = `Bearer ${tkn}`

          return config
        },
        err => Promise.reject(err))
      }
    }
  }

  async handleRequest (method, endpoint, payload = {}, config) {
    const component = this._parseStackComponent(endpoint)
    try {
      return await this[component]({
        method,
        url: endpoint,
        data: payload,
        ...config,
      })
    } catch (err) {
      if ('response' in err && err.response && 'data' in err.response) {
        throw err.response.data
      } else {
        throw err
      }
    }
  }

  async get (endpoint, ...rest) {
    const { fieldMask, queryParams } = rest[rest.length - 1]
    let fieldMaskPath
    if (fieldMask) {
      if (typeof fieldMask === 'string') {
        fieldMaskPath = fieldMask.paths.join(',')
      } else if ('paths' in fieldMask && typeof fieldMask.paths === 'string') {
        fieldMaskPath = fieldMask.join(',')
      }
    }

    const config = {}
    if (fieldMaskPath) {
      config.params = { field_mask: fieldMaskPath }
    }
    if (queryParams) {
      config.params = {
        ...config.params,
        ...queryParams,
      }
    }
    return this.handleRequest('get', endpoint, undefined, config)
  }

  async post (endpoint, payload) {
    return this.handleRequest('post', endpoint, payload)
  }

  async patch (endpoint, payload) {
    return this.handleRequest('patch', endpoint, payload)
  }

  async put (endpoint, payload, { fieldMask }) {
    if (fieldMask) {
      payload.field_mask = fieldMask
    }
    return this.handleRequest('put', endpoint, payload)
  }

  async delete (endpoint) {
    return this.handleRequest('delete', endpoint)
  }

  /**
   *  Extracts the stack component abbreviation from the endpoint.
   * @param {string} endpoint - The endpoint got for a request method.
   * @returns {string} One of {is|as|gs|js|ns}.
   */
  _parseStackComponent (endpoint) {
    try {
      const component = endpoint.split('/')[1]
      switch (component) {
      case 'as':
      case 'gs':
      case 'js':
      case 'ns':
        return component
      default:
        return 'is'
      }
    } catch (err) {
      throw new Error('Unable to extract the stack component:', endpoint)
    }
  }
}

export default Http
