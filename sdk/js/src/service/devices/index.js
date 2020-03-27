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

/* eslint-disable no-invalid-this, no-await-in-loop */

import traverse from 'traverse'
import Marshaler from '../../util/marshaler'
import combineStreams from '../../util/combine-streams'
import Device from '../../entity/device'
import { notify, EVENTS } from '../../api/stream/shared'
import deviceEntityMap from '../../../generated/device-entity-map.json'
import { splitSetPaths, splitGetPaths, makeRequests } from './split'
import mergeDevice from './merge'

/**
 * Devices Class provides an abstraction on all devices and manages data
 * handling from different sources. It exposes an API to easily work with
 * device data.
 */
class Devices {
  constructor(api, { proxy = true, ignoreDisabledComponents = true, stackConfig }) {
    if (!api) {
      throw new Error('Cannot initialize device service without api object.')
    }
    this._api = api
    this._stackConfig = stackConfig
    this._proxy = proxy
    this._ignoreDisabledComponents = ignoreDisabledComponents
  }

  _responseTransform(response, single = true) {
    return Marshaler[single ? 'unwrapDevice' : 'unwrapDevices'](
      response,
      this._proxy ? device => new Device(device, this._api) : undefined,
    )
  }

  _emitDefaults(paths, device) {
    // Handle zero coordinates that are swallowed by the grpc-gateway for device location.
    const hasLocation = Boolean(device.locations) && Boolean(device.locations.user)
    const requestedLocation = paths.some(path => path.startsWith('location'))

    if (hasLocation && requestedLocation) {
      const { locations } = device

      if (!('altitude' in locations.user)) {
        locations.user.altitude = 0
      }

      if (!('longitude' in locations.user)) {
        locations.user.longitude = 0
      }

      if (!('latitude' in locations.user)) {
        locations.user.latitude = 0
      }
    }

    if (paths.includes('claim_authentication_code') && !Boolean(device.claim_authentication_code)) {
      device.claim_authentication_code = null
    }

    return device
  }

  async _setDevice(applicationId, deviceId, device, create = false, requestTreeOverwrite) {
    const ids = device.ids
    const devId = deviceId || ('device_id' in ids && ids.device_id)
    const appId = applicationId || ('application_ids' in ids && ids.application_ids.application_id)

    if (deviceId && ids && 'device_id' in ids && deviceId !== ids.device_id) {
      throw new Error('Device ID mismatch.')
    }

    if (!create && !devId) {
      throw new Error('Missing device_id for update operation.')
    }

    if (!appId) {
      throw new Error('Missing application_id for device.')
    }

    // Ensure proper id object
    if (!('ids' in device)) {
      device.ids = { device_id: deviceId, application_ids: { application_id: applicationId } }
    } else if (!device.ids.device_id) {
      device.ids.device_id = deviceId
    } else if (!device.ids.application_ids || !device.ids.application_ids.application_id) {
      device.ids.application_ids = { application_id: applicationId }
    }

    const params = {
      routeParams: {
        'end_device.ids.application_ids.application_id': appId,
      },
    }

    // Extract the paths from the patch
    const deviceMap = traverse(deviceEntityMap)

    const commonPathFilter = function(element, index, array) {
      return deviceMap.has(array.slice(0, index + 1))
    }
    const paths = traverse(device).reduce(function(acc, node) {
      if (this.isLeaf) {
        const path = this.path

        // Only consider adding, if a common parent has not been already added
        if (
          acc.every(
            e =>
              !path
                .slice(-1)
                .join()
                .startsWith(e.join()),
          )
        ) {
          // Add only the deepest possible field mask of the patch
          const commonPath = path.filter(commonPathFilter)
          acc.push(commonPath)
        }
      }
      return acc
    }, [])

    const requestTree = requestTreeOverwrite ? requestTreeOverwrite : splitSetPaths(paths)

    if (create) {
      if (device.join_server_address !== this._stackConfig.jsHost) {
        delete requestTree.js
      }
    } else {
      // Retrieve join information if not present for update
      const paths = [
        ['application_server_address'],
        ['network_server_address'],
        ['join_server_address'],
      ]

      if (!('supports_join' in device)) {
        paths.push(['supports_join'])
      }

      const res = await this._getDevice(appId, devId, paths, true)
      if ('supports_join' in res && res.supports_join) {
        // The NS registry entry exists
        device.supports_join = true
      } else if (res.join_server_address) {
        // The NS registry entry does not exist, but a join_server_address
        // setting suggests that join is supported, so we add the path
        // to the request tree to ensure that it will be set on creation
        device.supports_join = true
        requestTree.ns.push(['supports_join'])
      }

      if (res.network_server_address !== this._stackConfig.nsHost) {
        delete requestTree.ns
      }

      if (res.application_server_address !== this._stackConfig.asHost) {
        delete requestTree.as
      }

      if (res.join_server_address !== this._stackConfig.jsHost) {
        delete requestTree.js
      }
    }

    // Do not query JS when the device is ABP
    if (!device.supports_join) {
      delete requestTree.js
    }

    // Retrieve necessary EUIs in case of a join server query being necessary
    if ('js' in requestTree) {
      if (!create && (!ids || !ids.join_eui || !ids.dev_eui)) {
        const res = await this._getDevice(
          appId,
          devId,
          [['ids', 'join_eui'], ['ids', 'dev_eui']],
          true,
        )
        if (!res.ids || !res.ids.join_eui || !res.ids.dev_eui) {
          throw new Error(
            'Could not update Join Server data on a device without Join EUI or Dev EUI',
          )
        }
        device.ids = {
          ...device.ids,
          join_eui: res.ids.join_eui,
          dev_eui: res.ids.dev_eui,
        }
      }
    }

    // Perform the requests
    const devicePayload = Marshaler.payload(device, 'end_device')
    const setParts = await makeRequests(
      this._api,
      this._stackConfig,
      this._ignoreDisabledComponents,
      create ? 'create' : 'set',
      requestTree,
      params,
      devicePayload,
    )

    // Filter out errored requests
    const errors = setParts.filter(part => part.hasErrored)

    // Handle possible errored requests
    if (errors.length !== 0) {
      // Roll back successfully created registry entries
      if (create) {
        const rollbackComponents = setParts.reduce((components, part) => {
          if (part.hasAttempted && !part.hasErrored) {
            components.push(part.component)
          }

          return components
        }, [])

        this._deleteDevice(appId, devId, rollbackComponents)
      }

      // Throw the first error
      throw errors[0].error
    }

    return mergeDevice(setParts)
  }

  async _getDevice(applicationId, deviceId, paths, ignoreNotFound) {
    if (!applicationId) {
      throw new Error('Missing application_id for device.')
    }

    if (!deviceId) {
      throw new Error('Missing device_id for device.')
    }

    const requestTree = splitGetPaths(paths)

    const params = {
      routeParams: {
        'end_device_ids.application_ids.application_id': applicationId,
        'end_device_ids.device_id': deviceId,
      },
    }

    const deviceParts = await makeRequests(
      this._api,
      this._stackConfig,
      this._ignoreDisabledComponents,
      'get',
      requestTree,
      params,
      undefined,
      ignoreNotFound,
    )
    const result = mergeDevice(deviceParts)

    return result
  }

  async _deleteDevice(applicationId, deviceId, components = ['is', 'ns', 'as', 'js']) {
    const params = {
      routeParams: {
        'application_ids.application_id': applicationId,
        device_id: deviceId,
      },
    }

    // Compose a request tree
    const requestTree = components.reduce(function(acc, val) {
      acc[val] = undefined
      return acc
    }, {})

    const deleteParts = await makeRequests(
      this._api,
      this._stackConfig,
      this._ignoreDisabledComponents,
      'delete',
      requestTree,
      params,
      undefined,
      true,
    )
    return deleteParts.every(e => Boolean(e.device) && Object.keys(e.device).length === 0)
      ? {}
      : deleteParts
  }

  async getAll(applicationId, params, selector) {
    const response = await this._api.EndDeviceRegistry.List(
      {
        routeParams: { 'application_ids.application_id': applicationId },
      },
      {
        ...params,
        ...Marshaler.selectorToFieldMask(selector),
      },
    )

    return this._responseTransform(response, false)
  }

  async search(applicationId, params, selector) {
    const response = await this._api.EndDeviceRegistrySearch.SearchEndDevices(
      {
        routeParams: { 'application_ids.application_id': applicationId },
      },
      {
        ...params,
        ...Marshaler.selectorToFieldMask(selector),
      },
    )

    return Marshaler.payloadListResponse('end_devices', response)
  }

  async getById(applicationId, deviceId, selector = [['ids']], { ignoreNotFound = false } = {}) {
    const response = await this._getDevice(
      applicationId,
      deviceId,
      Marshaler.selectorToPaths(selector),
      ignoreNotFound,
    )

    const { field_mask } = Marshaler.selectorToFieldMask(selector)
    const device = this._emitDefaults(field_mask.paths, Marshaler.unwrapDevice(response))

    return this._proxy ? new Device(device, this._api) : device
  }

  /**
   * Updates the `deviceId` end device under the `applicationId` application.
   * This method will cause updates of the end device in all available stack
   * components (i.e. NS, AS, IS, JS) based on provided end device payload.
   * @param {string} applicationId - Application ID
   * @param {string} deviceId - Device ID
   * @param {Object} patch - The end device payload
   * @returns {Object} - Updated end device on successful update, an error otherwise
   */
  async updateById(applicationId, deviceId, patch) {
    if (!Boolean(applicationId)) {
      throw new Error('Missing application ID for device')
    }

    if (!Boolean(deviceId)) {
      throw new Error('Missing end device ID')
    }

    const deviceMap = traverse(deviceEntityMap)
    const paths = traverse(patch).reduce(function(acc) {
      if (this.isLeaf) {
        const path = this.path

        const parentAdded = acc.some(e => path[0].startsWith(e.join()))

        // Only consider adding, if a common parent has not been already added
        if (!parentAdded) {
          // Add only the deepest possible field mask of the patch
          const commonPath = path.filter((_, index, array) => {
            const arr = array.slice(0, index + 1)
            return deviceMap.has(arr)
          })

          acc.push(commonPath)
        }
      }
      return acc
    }, [])

    const requestTree = splitSetPaths(paths)

    // Assemble paths for end device fields that need to be retrieved first to make the update request
    const combinePaths = []
    if ('as' in requestTree && !('application_server_address' in patch)) {
      combinePaths.push(['application_server_address'])
    }
    if ('js' in requestTree && !('join_server_address' in patch)) {
      combinePaths.push(['join_server_address'])
      combinePaths.push(['supports_join'])

      const { ids = {} } = patch
      if (!('dev_eui' in ids) || !('join_eui' in 'ids')) {
        combinePaths.push(['ids', 'dev_eui'])
        combinePaths.push(['ids', 'join_eui'])
      }
    }
    if ('ns' in requestTree && !('network_server_address' in patch)) {
      combinePaths.push(['network_server_address'])
    }

    const assembledValues = await this._getDevice(applicationId, deviceId, combinePaths, true)

    if (assembledValues.network_server_address !== this._stackConfig.nsHost) {
      delete requestTree.ns
    }

    if (assembledValues.application_server_address !== this._stackConfig.asHost) {
      delete requestTree.as
    }

    if (
      !assembledValues.supports_join ||
      assembledValues.join_server_address !== this._stackConfig.jsHost
    ) {
      delete requestTree.js
    }

    // Make sure to include `join_eui` and `dev_eui` for js request as those are required
    if ('js' in requestTree) {
      const { ids = {} } = patch
      const {
        ids: { join_eui, dev_eui },
      } = assembledValues

      patch.ids = {
        ...ids,
        join_eui,
        dev_eui,
      }
    }

    const routeParams = {
      routeParams: {
        'end_device.ids.application_ids.application_id': applicationId,
        'end_device.ids.device_id': deviceId,
      },
    }

    // Perform the requests
    const devicePayload = Marshaler.payload(patch, 'end_device')
    const setParts = await makeRequests(
      this._api,
      this._stackConfig,
      this._ignoreDisabledComponents,
      'set',
      requestTree,
      routeParams,
      devicePayload,
    )

    // Filter out errored requests
    const errors = setParts.filter(part => part.hasErrored)

    // Handle possible errored requests
    if (errors.length !== 0) {
      // Throw the first error
      throw errors[0].error
    }

    const device = this._emitDefaults(
      Marshaler.fieldMaskFromPatch(patch),
      Marshaler.unwrapDevice(mergeDevice(setParts)),
    )

    return this._proxy ? new Device(device, this._api) : device
  }

  /**
   * Creates an end device under the `applicationId` application.
   * This method will cause creating the end device in all available stack
   * components (i.e. NS, AS, IS, JS) based on provided end device payload.
   * @param {string} applicationId - Application ID
   * @param {Object} device - The end device payload
   * @returns {Object} - Created end device on successful creation, an error otherwise
   */
  async create(applicationId, device) {
    if (!Boolean(applicationId)) {
      throw new Error('Missing application ID for device')
    }

    const { supports_join = false, ids = {} } = device

    const deviceId = ids.device_id
    if (!Boolean(deviceId)) {
      throw new Error('Missing end device ID')
    }

    if (supports_join && 'provisioner_id' in device && device.provisioner_id !== '') {
      throw new Error('Setting a provisioner with end device keys is not allowed.')
    }

    const requestTree = splitSetPaths(traverse(device).paths())

    if (!supports_join || device.join_server_address !== this._stackConfig.jsHost) {
      delete requestTree.js
    }

    if (device.network_server_address !== this._stackConfig.nsHost) {
      delete requestTree.ns
    }

    if (device.application_server_address !== this._stackConfig.asHost) {
      delete requestTree.as
    }

    const devicePayload = Marshaler.payload(device, 'end_device')
    const routeParams = {
      routeParams: {
        'end_device.ids.application_ids.application_id': applicationId,
      },
    }
    const setParts = await makeRequests(
      this._api,
      this._stackConfig,
      this._ignoreDisabledComponents,
      'create',
      requestTree,
      routeParams,
      devicePayload,
    )

    // Filter out errored requests
    const errors = setParts.filter(part => part.hasErrored)

    // Handle possible errored requests
    if (errors.length !== 0) {
      // Roll back successfully created registry entries
      const rollbackComponents = setParts.reduce((components, part) => {
        if (part.hasAttempted && !part.hasErrored) {
          components.push(part.component)
        }

        return components
      }, [])

      this._deleteDevice(applicationId, deviceId, rollbackComponents)

      // Throw the first error
      throw errors[0].error
    }

    return mergeDevice(setParts)
  }

  async deleteById(applicationId, deviceId) {
    const result = this._deleteDevice(applicationId, deviceId)

    return result
  }

  // End Device Template Converter

  async listTemplateFormats() {
    const result = await this._api.EndDeviceTemplateConverter.ListFormats()
    const payload = Marshaler.payloadSingleResponse(result)

    return payload.formats
  }

  convertTemplate(formatId, data) {
    // This is a stream endpoint
    return this._api.EndDeviceTemplateConverter.Convert(undefined, {
      format_id: formatId,
      data,
    })
  }

  bulkCreate(applicationId, deviceOrDevices, components = ['is', 'ns', 'as', 'js']) {
    const devices = !(deviceOrDevices instanceof Array) ? [deviceOrDevices] : deviceOrDevices
    let listeners = Object.values(EVENTS).reduce((acc, curr) => ({ ...acc, [curr]: null }), {})
    let finishedCount = 0
    let stopRequested = false

    const runTasks = async function() {
      for (const device of devices) {
        if (stopRequested) {
          notify(listeners[EVENTS.CLOSE])
          listeners = null
          break
        }

        try {
          const {
            field_mask: { paths },
            end_device,
          } = device

          const requestTree = splitSetPaths(Marshaler.selectorToPaths(paths), undefined, components)

          const result = await this._setDevice(
            applicationId,
            undefined,
            end_device,
            true,
            requestTree,
          )
          notify(listeners[EVENTS.CHUNK], result)
          finishedCount++
          if (finishedCount === devices.length) {
            notify(listeners[EVENTS.CLOSE])
            listeners = null
          }
        } catch (error) {
          notify(listeners[EVENTS.ERROR], error)
          listeners = null
          break
        }
      }
    }

    runTasks.bind(this)()

    return {
      on(eventName, callback) {
        if (listeners[eventName] === undefined) {
          throw new Error(
            `${eventName} event is not supported. Should be one of: start, error, chunk or close`,
          )
        }

        listeners[eventName] = callback

        return this
      },
      abort() {
        stopRequested = true
      },
    }
  }

  // Events Stream

  async openStream(identifiers, tail, after) {
    const payload = {
      identifiers: identifiers.map(ids => ({
        device_ids: ids,
      })),
      tail,
      after,
    }

    // Event streams can come from multiple stack components. It is necessary to
    // check for stack components on different hosts and open distinct stream
    // connections for any distinct host if need be.
    const distinctComponents = this._stackConfig.getComponentsWithDistinctBaseUrls([
      'is',
      'js',
      'ns',
      'as',
      'dtc',
    ])

    const streams = distinctComponents.map(component =>
      this._api.Events.Stream({ component }, payload),
    )

    // Combine all stream sources to one subscription generator.
    return combineStreams(streams)
  }
}

export default Devices
