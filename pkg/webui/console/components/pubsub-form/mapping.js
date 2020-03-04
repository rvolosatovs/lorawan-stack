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

import { omit } from 'lodash'
import { natsUrl as natsUrlRegexp } from '../../lib/regexp'

const natsBlankValues = {
  username: '',
  password: '',
  address: '',
  port: '',
  secure: false,
  use_credentials: true,
}

const mqttBlankValues = {
  server_url: '',
  client_id: '',
  username: '',
  password: '',
  subscribe_qos: '',
  publish_qos: '',
  use_tls: false,
  tls_ca: '',
  tls_client_cert: '',
  tls_client_key: '',
  use_credentials: true,
}

const omitCredentials = function(result) {
  return omit(result, 'mqtt.use_credentials')
}

export const mapNatsServerUrlToFormValue = function(server_url) {
  try {
    const res = server_url.match(natsUrlRegexp)
    return {
      secure: res[2] === 'tls',
      username: res[5],
      password: res[7],
      address: res[8],
      port: res[10],
      use_credentials: Boolean(res[5] || res[7]),
    }
  } catch {
    return {}
  }
}

const mapPubsubMessageTypeToFormValue = messageType =>
  (messageType && { enabled: true, value: messageType.topic }) || { enabled: false, value: '' }

export const mapPubsubToFormValues = function(pubsub) {
  const isNats = 'nats' in pubsub
  const isMqtt = 'mqtt' in pubsub
  const result = {
    pub_sub_id: pubsub.ids.pub_sub_id,
    base_topic: pubsub.base_topic,
    format: pubsub.format,
    _provider: isMqtt ? 'mqtt' : 'nats',
    nats: isNats ? mapNatsServerUrlToFormValue(pubsub.nats.server_url) : natsBlankValues,
    mqtt: isMqtt ? pubsub.mqtt : mqttBlankValues,
    downlink_ack: mapPubsubMessageTypeToFormValue(pubsub.downlink_ack),
    downlink_failed: mapPubsubMessageTypeToFormValue(pubsub.downlink_failed),
    downlink_nack: mapPubsubMessageTypeToFormValue(pubsub.downlink_nack),
    downlink_push: mapPubsubMessageTypeToFormValue(pubsub.downlink_push),
    downlink_queued: mapPubsubMessageTypeToFormValue(pubsub.downlink_queued),
    downlink_replace: mapPubsubMessageTypeToFormValue(pubsub.downlink_replace),
    downlink_sent: mapPubsubMessageTypeToFormValue(pubsub.downlink_sent),
    join_accept: mapPubsubMessageTypeToFormValue(pubsub.join_accept),
    location_solved: mapPubsubMessageTypeToFormValue(pubsub.location_solved),
    uplink_message: mapPubsubMessageTypeToFormValue(pubsub.uplink_message),
  }

  if (!result.mqtt.tls_ca) {
    result.mqtt.tls_ca = ''
  }
  if (!result.mqtt.tls_client_cert) {
    result.mqtt.tls_client_cert = ''
  }
  if (!result.mqtt.tls_client_key) {
    result.mqtt.tls_client_key = ''
  }
  if (!result.mqtt.username) {
    result.mqtt.username = ''
  }
  if (!result.mqtt.password) {
    result.mqtt.password = ''
  }
  if (isMqtt) {
    result.mqtt.use_credentials = Boolean(result.mqtt.username || result.mqtt.password)
  }

  return result
}

const mapNatsConfigFormValueToNatsServerUrl = function({
  username,
  password,
  address,
  port,
  secure,
  use_credentials,
}) {
  return `${secure ? 'tls' : 'nats'}://${
    use_credentials ? `${username}:${password}@` : ''
  }${address}:${port}`
}

const mapMessageTypeFormValueToPubsubMessageType = formValue =>
  (formValue.enabled && { topic: formValue.value }) || null

export const mapFormValuesToPubsub = function(values, appId) {
  const result = {
    ids: {
      application_ids: {
        application_id: appId,
      },
      pub_sub_id: values.pub_sub_id,
    },
    base_topic: values.base_topic,
    format: values.format,
    downlink_ack: mapMessageTypeFormValueToPubsubMessageType(values.downlink_ack),
    downlink_failed: mapMessageTypeFormValueToPubsubMessageType(values.downlink_failed),
    downlink_nack: mapMessageTypeFormValueToPubsubMessageType(values.downlink_nack),
    downlink_push: mapMessageTypeFormValueToPubsubMessageType(values.downlink_push),
    downlink_queued: mapMessageTypeFormValueToPubsubMessageType(values.downlink_queued),
    downlink_replace: mapMessageTypeFormValueToPubsubMessageType(values.downlink_replace),
    downlink_sent: mapMessageTypeFormValueToPubsubMessageType(values.downlink_sent),
    join_accept: mapMessageTypeFormValueToPubsubMessageType(values.join_accept),
    location_solved: mapMessageTypeFormValueToPubsubMessageType(values.location_solved),
    uplink_message: mapMessageTypeFormValueToPubsubMessageType(values.uplink_message),
  }

  if (values._provider === 'nats') {
    result.nats = {
      server_url: mapNatsConfigFormValueToNatsServerUrl(values.nats),
    }
  } else if (values._provider === 'mqtt') {
    result.mqtt = values.mqtt
    return omitCredentials(result)
  }

  return result
}

export const blankValues = {
  pub_sub_id: '',
  base_topic: '',
  format: '',
  _provider: 'nats',
  nats: natsBlankValues,
  mqtt: mqttBlankValues,
  downlink_ack: { enabled: false, value: '' },
  downlink_failed: { enabled: false, value: '' },
  downlink_nack: { enabled: false, value: '' },
  downlink_push: { enabled: false, value: '' },
  downlink_queued: { enabled: false, value: '' },
  downlink_replace: { enabled: false, value: '' },
  downlink_sent: { enabled: false, value: '' },
  join_accept: { enabled: false, value: '' },
  location_solved: { enabled: false, value: '' },
  uplink_message: { enabled: false, value: '' },
}
