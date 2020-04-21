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

import { selectJsConfig } from '@ttn-lw/lib/selectors/env'
import getHostnameFromUrl from '@ttn-lw/lib/host-from-url'

const lwRegexp = /^[1-9].[0-9].[0-9]$/
const lwCache = {}

/**
 * Parses string representation of the lorawan mac version to number.
 * @param {string} strMacVersion - Formatted string representation fot the lorawan mac version, e.g. 1.1.0.
 * @returns {number} - Number representation of the lorawan mac version. Returns 0 if provided
 * argument is not a valid string representation of the lorawan mac version.
 * @example
 *  const parsedVersion = parseLorawanMacVersion('1.0.0'); // returns 100
 *  const parsedVersion = parseLorawanMacVersion('1.1.0'); // returns 110
 *  const parsedVersion = parseLorawanMacVersion(''); // returns 0
 *  const parsedVersion = parseLorawanMacVersion('str'); // returns 0
 */
export const parseLorawanMacVersion = strMacVersion => {
  if (lwCache[strMacVersion]) {
    return lwCache[strMacVersion]
  }

  if (!Boolean(strMacVersion)) {
    return 0
  }

  const match = lwRegexp.exec(strMacVersion)
  if (!match.length) {
    return 0
  }

  const parsed = parseInt(match[0].replace(/\D/g, '').padEnd(3, 0))
  lwCache[strMacVersion] = parsed

  return lwCache[strMacVersion]
}

/**
 * Returns whether the device is OTAA.
 * Note: device type is mainly derived based on the `supports_join` and `multicast` fields.
 * However, in cases when NS is not available, `root_keys` can be used to determine whether
 * the device is OTAA.
 * @param {Object} device - The device object.
 * @returns {boolean} `true` if the device is OTAA, `false` otherwise
 */
export const isDeviceOTAA = device =>
  Boolean(device) && (Boolean(device.supports_join) || Boolean(device.root_keys))

/**
 * Returns whether the device is ABP.
 * @param {Object} device - The device object.
 * @returns {boolean} `true` if the device is ABP, `false` otherwise
 */
export const isDeviceABP = device =>
  Boolean(device) && !Boolean(device.supports_join) && !Boolean(device.multicast)

/**
 * Returns whether the device is multicast.
 * @param {Object} device - The device object.
 * @returns {boolean} `true` if the device is multicast, `false` otherwise
 */
export const isDeviceMulticast = device => Boolean(device) && Boolean(device.multicast)

/**
 * Returns whether an end device is provisioned on an external join server.
 * @param {Object} device - The device object.
 * @returns {boolean} `true` if the end device is provisioned on an external join server, `false` otherwise.
 */
export const hasExternalJs = device => {
  const { enabled, base_url } = selectJsConfig()

  const deviceJs = device.join_server_address
  const stackJs = getHostnameFromUrl(base_url)

  return !enabled || typeof deviceJs === 'undefined' || deviceJs !== stackJs
}

export const isDeviceJoined = device =>
  Boolean(device) &&
  Boolean(device.session) &&
  Boolean(device.session.dev_addr) &&
  Boolean(device.session.keys) &&
  Boolean(Object.keys(device.session.keys).length)

export const ACTIVATION_MODES = Object.freeze({
  OTAA: 'otaa',
  ABP: 'abp',
  MULTICAST: 'multicast',
})
