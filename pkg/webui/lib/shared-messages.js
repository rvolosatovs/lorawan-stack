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

import { defineMessages } from 'react-intl'

export default defineMessages({
  // Keep these sorted alphabetically
  add: 'Add',
  addApiKey: 'Add API Key',
  addApplication: 'Add Application',
  addDevice: 'Add Device',
  addGateway: 'Add Gateway',
  addOrganization: 'Add Organization',
  addressPlaceholder: 'host',
  addWebhook: 'Add Webhook',
  all: 'All',
  altitude: 'Altitude',
  altitudeDesc: 'The altitude in meters, where 0 means sea level',
  antennas: 'Antennas',
  apiKey: 'API Key',
  apiKeys: 'API Keys',
  appId: 'Application ID',
  appKey: 'AppKey',
  application: 'Application',
  applications: 'Applications',
  applicationServerAddress: 'Application Server Address',
  approve: 'Approve',
  appSKey: 'AppSKey',
  brand: 'Brand',
  cancel: 'Cancel',
  changePassword: 'Change Password',
  changeLocation: 'Change location settings',
  clear: 'Clear',
  collaboratorAdd: 'Add Collaborator',
  collaboratorDeleteSuccess: 'Successfully removed collaborator',
  collaboratorEdit: 'Edit {collaboratorId}',
  collaboratorEditRights: 'Edit rights of {collaboratorId}',
  collaboratorId: 'Collaborator ID',
  collaboratorModalWarning: 'Are you sure you want to remove {collaboratorId} as a collaborator?',
  collaboratorRemove: 'Collaborator Remove',
  collaborators: 'Collaborators',
  collaboratorUpdateSuccess: 'Successfully updated collaborator rights',
  componentApplicationServer: 'Application Server',
  componentGatewayServer: 'Gateway Server',
  componentIdentityServer: 'Identity Server',
  componentJoinServer: 'Join Server',
  componentNetworkServer: 'Network Server',
  confirmPassword: 'Confirm Password',
  connecting: 'Connecting',
  createApiKey: 'Create API Key',
  created: 'Created',
  createdAt: 'Created at',
  currentCollaborators: 'Current Collaborators',
  data: 'Data',
  defineRights: 'Define Rights',
  description: 'Description',
  devAddr: 'Device Address',
  devDesc: 'Device Description',
  devEUI: 'DevEUI',
  devices: 'Devices',
  devID: 'Device ID',
  devName: 'Device Name',
  disconnected: 'Disconnected',
  downlink: 'Downlink',
  downlinksScheduled: 'Downlinks (re)scheduled',
  edit: 'Edit',
  email: 'Email',
  enabled: 'Enabled',
  entityId: 'Entity ID',
  eventsCannotShow: 'Cannot show events',
  firmwareVersion: 'Firmware Version',
  frequencyPlan: 'Frequency Plan',
  fwdNtwkKey: 'FNwkSIntKey',
  gatewayDescription: 'Gateway Description',
  gatewayEUI: 'Gateway EUI',
  gatewayID: 'Gateway ID',
  gatewayName: 'Gateway Name',
  gateways: 'Gateways',
  gatewayServerAddress: 'Gateway Server Address',
  general: 'General',
  generalInformation: 'General Information',
  generalSettings: 'General Settings',
  hardware: 'Hardware',
  hardwareVersion: 'Hardware version',
  id: 'ID',
  integrations: 'Integrations',
  joinEUI: 'JoinEUI',
  joinServerAddress: 'Join Server Address',
  keyEdit: 'Edit API Key',
  keyId: 'Key ID',
  lastSeen: 'Last seen',
  latitude: 'Latitude',
  latitudeDesc: 'The North-South position in degrees, where 0 is the equator',
  link: 'Link',
  loading: 'Loading',
  location: 'Location',
  login: 'Login',
  logout: 'Logout',
  longitude: 'Longitude',
  longitudeDesc: 'The East-West position in degrees, where 0 is the Prime Meridian (Greenwich)',
  lorawanInformation: 'LoRaWAN Information',
  lorawanOptions: 'LoRaWAN Options',
  macVersion: 'MAC Version',
  model: 'Model',
  name: 'Name',
  networkServerAddress: 'Network Server Address',
  noDesc: 'This device has no description',
  noEvents: '{entityId} has not sent any events recently',
  noLocation: 'No location information available',
  noMatch: 'No items found',
  none: 'None',
  notAvailable: 'n/a',
  nsAddress: 'Network Server Address',
  nsEmptyDefault: 'Leave empty to link to the Network Server in the same cluster',
  ntwkSEncKey: 'NwkSEncKey',
  nwkKey: 'NwkKey',
  ok: 'Ok',
  options: 'Options',
  organization: 'Organization',
  organizationId: 'Organization ID',
  organizations: 'Organizations',
  otherCluster: 'Other Cluster',
  overview: 'Overview',
  password: 'Password',
  pause: 'Pause',
  payloadFormatters: 'Payload Formatters',
  payloadFormattersDownlink: 'Downlink Payload Formatters',
  payloadFormattersUpdateFailure: 'There was an error updating the Payload Formatter',
  payloadFormattersUpdateSuccess: 'Payload Formatter has been set successfully',
  payloadFormattersUplink: 'Uplink Payload Formatters',
  phyVersion: 'PHY Version',
  privacyPolicy: 'Privacy Policy',
  refreshPage: 'Refresh page',
  removeCollaborator: 'Remove Collaborator',
  resume: 'Resume',
  rights: 'Rights',
  saveChanges: 'Save Changes',
  settings: 'Settings',
  sNtwkSIKey: 'SNwkSIntKey',
  takeMeBack: 'Take me back',
  termsAndCondition: 'Terms and Conditions',
  time: 'Time',
  traffic: 'Traffic',
  type: 'Type',
  unknown: 'Unknown',
  updatedAt: 'Last updated at',
  uplink: 'Uplink',
  uplinksReceived: 'Uplinks Received',
  user: 'User',
  userId: 'User ID',
  users: 'Users',
  validateAddressFormat: 'Wrong address format, should be host or host:port',
  validateAlphanum: 'The value must be alphanumeric and contain no spaces',
  validateEmail: 'Not valid email',
  validateFormat: 'Wrong format',
  validateInt32: 'Invalid numeric format',
  validateLatLong: 'Ivalid value format (valid format e.g. 52.15249000)',
  validatePasswordMatch: 'Passwords should match',
  validateRequired: 'This field is required',
  validateRights: 'At least one right should be selected',
  validateTooLong: 'Too long',
  validateTooShort: 'Too short',
  validateUrl: 'Invalid URL format',
  webhookBaseUrl: 'Base URL',
  webhookFormat: 'Webhook Format',
  webhookId: 'Webhook ID',
  webhooks: 'Webhooks',
})
