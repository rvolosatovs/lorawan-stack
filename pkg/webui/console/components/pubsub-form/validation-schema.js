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

import Yup from '@ttn-lw/lib/yup'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import {
  id as idRegexp,
  address as addressRegexp,
  mqttUrl as mqttUrlRegexp,
  mqttPassword as mqttPasswordRegexp,
  noSpaces as noSpacesRegexp,
} from '@console/lib/regexp'

import { qosLevels } from './qos-options'
import providers from './providers'

export default Yup.object().shape({
  pub_sub_id: Yup.string()
    .matches(idRegexp, Yup.passValues(sharedMessages.validateIdFormat))
    .min(2, Yup.passValues(sharedMessages.validateTooShort))
    .max(25, Yup.passValues(sharedMessages.validateTooLong))
    .required(sharedMessages.validateRequired),
  format: Yup.string().required(sharedMessages.validateRequired),
  base_topic: Yup.string().required(sharedMessages.validateRequired),
  nats: Yup.object().when('_provider', {
    is: providers.NATS,
    then: Yup.object().shape({
      _use_credentials: Yup.boolean(),
      username: Yup.string().when('_use_credentials', {
        is: true,
        then: Yup.string()
          .matches(idRegexp, Yup.passValues(sharedMessages.validateIdFormat))
          .min(2, Yup.passValues(sharedMessages.validateTooShort))
          .max(100, Yup.passValues(sharedMessages.validateTooLong))
          .required(sharedMessages.validateRequired),
        otherwise: Yup.string().strip(),
      }),
      password: Yup.string().when('_use_credentials', {
        is: true,
        then: Yup.string()
          .min(2, Yup.passValues(sharedMessages.validateTooShort))
          .max(100, Yup.passValues(sharedMessages.validateTooLong))
          .required(sharedMessages.validateRequired),
        otherwise: Yup.string().strip(),
      }),
      address: Yup.string()
        .matches(addressRegexp, Yup.passValues(sharedMessages.validateAddressFormat))
        .required(sharedMessages.validateRequired),
      port: Yup.number()
        .integer(sharedMessages.validateInt32)
        .positive(sharedMessages.validateInt32)
        .required(sharedMessages.validateRequired),
      secure: Yup.boolean(),
    }),
    otherwise: Yup.object().strip(),
  }),
  mqtt: Yup.object().when('_provider', {
    is: providers.MQTT,
    then: Yup.object().shape({
      _use_credentials: Yup.boolean(),
      server_url: Yup.string()
        .matches(mqttUrlRegexp, Yup.passValues(sharedMessages.validateUrl))
        .required(sharedMessages.validateRequired),
      client_id: Yup.string()
        .matches(noSpacesRegexp, Yup.passValues(sharedMessages.validateNoSpaces))
        .min(2, Yup.passValues(sharedMessages.validateTooShort))
        .max(23, Yup.passValues(sharedMessages.validateTooLong))
        .required(sharedMessages.validateRequired),
      username: Yup.string().when('_use_credentials', {
        is: true,
        then: Yup.string()
          .matches(noSpacesRegexp, Yup.passValues(sharedMessages.validateNoSpaces))
          .min(2, Yup.passValues(sharedMessages.validateTooShort))
          .max(100, Yup.passValues(sharedMessages.validateTooLong))
          .required(sharedMessages.validateRequired),
        otherwise: Yup.string().strip(),
      }),
      password: Yup.string().when('_use_credentials', {
        is: true,
        then: Yup.string().matches(
          mqttPasswordRegexp,
          Yup.passValues(sharedMessages.validateMqttPassword),
        ),
        otherwise: Yup.string().strip(),
      }),
      subscribe_qos: Yup.string()
        .oneOf(qosLevels, sharedMessages.validateRequired)
        .required(sharedMessages.validateRequired),
      publish_qos: Yup.string()
        .oneOf(qosLevels, sharedMessages.validateRequired)
        .required(sharedMessages.validateRequired),
      use_tls: Yup.boolean(),
      tls_ca: Yup.string().when('use_tls', {
        is: true,
        then: Yup.string().required(sharedMessages.validateRequired),
        otherwise: Yup.string().strip(),
      }),
      tls_client_cert: Yup.string().when('use_tls', {
        is: true,
        then: Yup.string().required(sharedMessages.validateRequired),
        otherwise: Yup.string().strip(),
      }),
      tls_client_key: Yup.string().when('use_tls', {
        is: true,
        then: Yup.string().required(sharedMessages.validateRequired),
        otherwise: Yup.string().strip(),
      }),
    }),
    otherwise: Yup.object().strip(),
  }),
})
