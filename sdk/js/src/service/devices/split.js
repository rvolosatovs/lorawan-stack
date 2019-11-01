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

import traverse from 'traverse'
import Marshaler from '../../util/marshaler'
import deviceEntityMap from '../../../generated/device-entity-map.json'

/** Takes the requested paths of the device and returns a request tree. The
 * splitting is achieved by looking up path responsibilities as defined in the
 * generated device entity map json.
 * @param {Object} paths - The requested paths (from the field mask) of the device
 * @param {string} direction - The direction, either 'set' or 'get'
 * @param {Object} base - An optional base value for the returned request tree
 * @param {Object} components - A component whitelist, unincluded components
 * will be excluded from the request tree
 * @returns {Object} A request tree object, consisting of resulting paths for each
 * component eg: { is: ['ids'], as: ['session'], js: ['root_keys'] }
 */
function splitPaths(paths = [], direction, base = {}, components = ['is', 'ns', 'as', 'js']) {
  const result = base
  const retrieveIndex = direction === 'get' ? 0 : 1

  for (const path of paths) {
    // Look up the current path in the device entity map
    const subtree = traverse(deviceEntityMap).get(path) || traverse(deviceEntityMap).get([path[0]])

    if (!subtree) {
      throw new Error(`Invalid or unknown field mask path used: ${path}`)
    }

    const definition = '_root' in subtree ? subtree._root[retrieveIndex] : subtree[retrieveIndex]

    const map = function(requestTree, component, path) {
      if (components.includes(component)) {
        result[component] = !result[component] ? [path] : [...result[component], path]
      }
    }

    if (definition) {
      if (definition instanceof Array) {
        for (const component of definition) {
          map(result, component, path)
        }
      } else {
        map(result, definition, path)
      }
    }
  }
  return result
}

/** A wrapper function to obtain a request tree for writing values to a device
 * @param {Object} paths - The requested paths (from the field mask) of the device
 * @param {Object} base - An optional base value for the returned request tree
 * @param {Object} components - A component whitelist, unincluded components
 * will be excluded from the request tree
 * @returns {Object} A request tree object, consisting of resulting paths for each
 * component eg: { is: ['ids'], as: ['session'], js: ['root_keys'] }
 */
export function splitSetPaths(paths, base, components) {
  return splitPaths(paths, 'set', base, components)
}

/** A wrapper function to obtain a request tree for reading values to a device
 * @param {Object} paths - The requested paths (from the field mask) of the device
 * @param {Object} base - An optional base value for the returned request tree
 * @param {Object} components - A component whitelist, unincluded components
 * will be excluded from the request tree
 * @returns {Object} A request tree object, consisting of resulting paths for each
 * component eg: { is: ['ids'], as: ['session'], js: ['root_keys'] }
 */
export function splitGetPaths(paths, base, components) {
  return splitPaths(paths, 'get', base, components)
}

/** makeRequests will make the necessary api calls based on the request tree and
 * other options
 * @param {Object} api - The Api object as passed to the service
 * @param {Object} stackConfig - The Things Stack config object
 * @param {boolean} ignoreDisabledComponents - A flag indicating whether queries
 * against disabled components should be ignored insread of throwing
 * @param {string} operation - The operation, an enum of 'create', 'set', 'get'
 * and 'delete'
 * @param {string} requestTree - The request tree, as returned by the splitPaths
 * function
 * @param {Object} params - The parameters object to be passed to the requests
 * @param {Object} payload - The payload to be passed to the requests
 * @param {boolean} ignoreNotFound - A flag indicating whether not found errors
 * should be translated to an empty device instead of throwing
 * @returns {Object} An array of device registry responses together with the paths
 * (field_mask) that they were requested with
 */
export async function makeRequests(
  api,
  stackConfig,
  ignoreDisabledComponents,
  operation,
  requestTree,
  params,
  payload,
  ignoreNotFound = false,
) {
  const isCreate = operation === 'create'
  const isSet = operation === 'set'
  const isDelete = operation === 'delete'
  const rpcFunction = isSet || isCreate ? 'Set' : isDelete ? 'Delete' : 'Get'

  // Use a wrapper for the api calls to control the result object and allow
  // ignoring not found errors per component, if wished
  const requestWrapper = async function(
    call,
    params,
    component,
    payload,
    ignoreRequestNotFound = ignoreNotFound,
  ) {
    const res = { hasAttempted: true, component, paths: requestTree[component], hasErrored: false }
    try {
      const result = await call(params, !isDelete ? payload : undefined)
      return { ...res, device: Marshaler.payloadSingleResponse(result) }
    } catch (error) {
      if (error.code === 5 && ignoreRequestNotFound) {
        return { ...res, device: {} }
      }

      return { ...res, hasErrored: true, error }
    }
  }

  const requests = new Array(3)

  // Check whether the request would query against disabled components
  if (!ignoreDisabledComponents) {
    for (const component of Object.keys(requestTree)) {
      if (!stackConfig[component]) {
        throw new Error(
          `Cannot run ${operation.toUpperCase()} end device request which (partially) depends on disabled component: "${component}".`,
        )
      }
    }
  }

  if (isSet && !('end_device.ids.device_id' in params.routeParams)) {
    // Ensure using the PUT method by setting the device id route param. This
    // ensures upserting without issues.
    const { end_device } = payload
    const { ids: { device_id } = {} } = end_device
    if (device_id) {
      params.routeParams['end_device.ids.device_id'] = device_id
    }
  }

  const result = [
    { component: 'ns', hasAttempted: false, hasErrored: false },
    { component: 'as', hasAttempted: false, hasErrored: false },
    { component: 'js', hasAttempted: false, hasErrored: false },
    { component: 'is', hasAttempted: false, hasErrored: false },
  ]

  // Do a possible IS request first
  if (stackConfig.is && 'is' in requestTree) {
    let func
    if (isSet) {
      func = 'Update'
    } else if (isCreate) {
      func = 'Create'
    } else if (isDelete) {
      func = 'Delete'
    } else {
      func = 'Get'
    }
    result[3] = await requestWrapper(
      api.EndDeviceRegistry[func],
      params,
      'is',
      {
        ...payload,
        ...Marshaler.pathsToFieldMask(requestTree.is),
      },
      false,
    )

    if (isCreate) {
      // Abort and return the result object when the IS create request has failed
      if (result[3].hasErrored) {
        return result
      }
      // Set the device id param based on the id of the newly created device
      params.routeParams['end_device.ids.device_id'] = result[3].device.ids.device_id
    }
  }

  // Compose an array of possible api calls to NS, AS, JS
  if (stackConfig.ns && 'ns' in requestTree) {
    requests[0] = requestWrapper(api.NsEndDeviceRegistry[rpcFunction], params, 'ns', {
      ...payload,
      ...Marshaler.pathsToFieldMask(requestTree.ns),
    })
  }
  if (stackConfig.as && 'as' in requestTree) {
    requests[1] = requestWrapper(api.AsEndDeviceRegistry[rpcFunction], params, 'as', {
      ...payload,
      ...Marshaler.pathsToFieldMask(requestTree.as),
    })
  }
  if (stackConfig.js && 'js' in requestTree) {
    requests[2] = requestWrapper(api.JsEndDeviceRegistry[rpcFunction], params, 'js', {
      ...payload,
      ...Marshaler.pathsToFieldMask(requestTree.js),
    })
  }

  // Run the requests in parallel
  const responses = await Promise.all(requests)

  // Attach the results to the result array
  for (const [i, response] of responses.entries()) {
    if (response) {
      result[i] = response
    }
  }

  return result
}
