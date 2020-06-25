# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Option to reset end device payload formatters in the Console.
- Service discovery using DNS SRV records for external Application Server linking.
- Functionality to set end device attributes in the Console.
- Event description tooltip to events in the Console.
- CLI support for setting and unsetting end device location (see `--location.latitude`, `--location.longitude`, `--location.altitude` and `--location.accuracy` options).
- Functionality to allow admin users to list all applications and gateways in the Console.
- Ursalink UG8X gateway documentation.
- Intercom, Google Analytics, and Emojicom feedback in documentation.

### Changed

- JSON uplink message doc edited for clarity.
- The CLI snap version uses the `$SNAP_USER_COMMON` directory for config by default, so that it is preserved between revisions.
- Defer events subscriptions until there is actual interest for events.

### Deprecated

### Removed

### Fixed

- Endless authentication refresh loop in the Console in some rare situations.
- Logout operation not working properly in the Console in some rare situations.
- Handling API key deletion event for applications, gateways, organizations and users.
- Organization API key deletion in the Console.
- CLI now only sends relevant end device fields to Identity Server on create.
- Maximum ADR data rate index used in 1.0.2a and earlier versions of AU915 band.
- End device events stream restart in the Console.
- CLI was unable to read input from pipes.

## [3.8.4] - 2020-06-12

### Added

- Metrics for log messages, counted per level and namespace.
- Allow suppressing logs on HTTP requests for user-defined paths (see `--http.log-ignore-paths` option).
- Redux state and actions reporting to Sentry
- Serving frontend sourcemaps in production
- Frequency plan documentation.
- LoRa Basics Station documentation.

### Changed

- Suppress a few unexpected EOF errors, in order to reduce noise in the logs for health checks.
- End device creation form with wizard in the Console.

### Fixed

- Packet Broker Agent cluster ID is used as subscription group.
- LinkADR handling in 72-channel bands.
- Data uplink metrics reported by Application Server.
- Timezones issue in claim auth code form, causing time to reverse on submission

## [3.8.3] - 2020-06-05

### Added

- Favicon to documentation pages.
- Draft template for documentation.

### Changed

- Late scheduling algorithm; Gateway Server now takes the 90th percentile of at least the last 5 round-trip times of the last 30 minutes into account to determine whether there's enough time to send the downlink to the gateway. This was the highest round-trip time received while the gateway was connected.

### Fixed

- Downlink scheduling to gateways which had one observed round-trip time that was higher than the available time to schedule. In some occassions, this broke downlink at some point while the gateway was connected.

## [3.8.2] - 2020-06-03

### Added

- Console logout is now propagated to the OAuth provider.
  - This requires a database migration (`ttn-lw-stack is-db migrate`) because of the added columns.
  - To set the `logout-redirect-uris` for existing clients, the CLI client can be used, e.g.: `ttn-lw-cli clients update console --logout-redirect-uris "https://localhost:8885/console" --redirect-uris "http://localhost:1885/console"`.
- Packet Broker Agent to act as Forwarder and Home Network. See `pba` configuration section.
- JavaScript style guide to our `DEVELOPMENT.md` documentation.
- Schedule end device downlinks in the Console.
- Support for repeated `RekeyInd`. (happens when e.g. `RekeyConf` is lost)
- Validate the `DevAddr` when switching session as a result of receiving `RekeyInd`.
- Error details for failed events in the Console.
- `Unknown` and `Other cluster` connection statuses to the gateways table in the Console.
- LoRaWAN 2.4 GHz band `ISM2400`.
- Unset end device fields using the CLI (see `--unset` option)
- Join EUI and Dev EUI columns to the end device table in the Console.
- CLI creates user configuration directory if it does not exist when generating configuration file.
- Upgrading guide in docs.
- Glossary.
- Event details in the Console traffic view.
- Gateway Server events for uplink messages now contain end device identifiers.
- Setting custom gateway attributes in the Console.
- Pub/Sub documentation.
- Return informative well-known errors for standard network and context errors.
- Error notification in list views in the Console.
- Latest "last seen" info and uplink frame counts for end devices in the Console.
- Latest "last seen" info for applications in the Console.

### Changed

- Conformed JavaScript to new code style guide.
- Removed login page of the Console (now redirects straight to the OAuth login).
- Network Server now records `LinkADRReq` rejections and will not retry rejected values.
- Improved `NewChannelReq`, `DLChannelReq` and `LinkADRReq` efficiency.
- For frames carrying only MAC commands, Network Server now attempts to fit them in FOpts omitting FPort, if possible, and sends them in FRMPayload with FPort 0 as usual otherwise.
- Submit buttons are now always enabled in the Console, regardless of the form's validation state.
- Disabled ADR for `ISM2400` band.
- Network Server will attempt RX1 for devices with `Rx1Delay` of 1 second, if possible.
- Network Server will not attempt to schedule MAC-only frames in ping slots or RXC windows.
- Network Server will only attempt to schedule in a ping slot or RXC window after RX2 has passed.
- Network Server will schedule all time-bound network-initiated downlinks at most RX1 delay ahead of time.
- Network Server now uses its own internal clock in `DeviceTimeAns`.
- Troubleshooting section of `DEVELOPMENT.md`
- Change console field labels from `MAC version` and `PHY version` to `LoRaWAN version` and `Regional Parameters version` and add descriptions

### Fixed

- Handling of device unsetting the ADR bit in uplink, after ADR has been started.
- Invalid `oauth-server-address` in CLI config generated by `use` command when config file is already present.
- Network Server now properly handles FPort 0 data uplinks carrying FOpts.
- Data rate 4 in version `1.0.2-a` of AU915.
- Incorrect `TxOffset` values used by Network Server in some bands.
- OAuth authorization page crashing.
- Byte input in scheduling downlink view.
- OAuth client token exchange and refresh issues when using TLS with a RootCA.
- Join Server and Application Server device registries now return an error when deleting keys on `SET` operations. The operation was never supported and caused an error on `GET` instead.
- Clearing end device events list in the Console.
- Some views not being accessible in the OAuth app (e.g. update password).
- `LinkADRReq` scheduling.
- Unsetting NwkKey in Join Server.
- CSRF token validation issues preventing login and logout in some circumstances.
- Typo in Application Server configuration documentation (webhook downlink).
- Unset fields via CLI on Join Server, i.e. `--unset root-keys.nwk-key`.
- Reconnecting UDP gateways that were disconnected by a new gateway connection.
- ADR in US915-like bands.

## [3.7.2] - 2020-04-22

### Added

- CLI can now dump JSON encoded `grpc_payload` field for unary requests (see `--dump-requests` flag).
- Template ID column in the webhook table in the Console.
- Select all field mask paths in CLI get, list and search commands (see `--all` option).
- Create webhooks via webhook templates in the Console.
- `ns.up.data.receive` and `ns.up.join.receive` events, which are triggered when respective uplink is received and matched to a device by Network Server.
- `ns.up.data.forward` and `ns.up.join.accept.forward` events, which are triggered when respective message is forwarded from Network Server to Application Server.
- `ns.up.join.cluster.attempt` and `ns.up.join.interop.attempt` events, which are triggered when the join-request is sent to respective Join Server by the Network Server.
- `ns.up.join.cluster.success` and `ns.up.join.interop.success` events, which are triggered when Network Server's join-request is accepted by respective Join Server.
- `ns.up.join.cluster.fail` and `ns.up.join.interop.fail` events, which are triggered when Network Server's join-request to respective Join Server fails.
- `ns.up.data.process` and `ns.up.join.accept.process` events, which are triggered when respective message is successfully processed by Network Server.
- `ns.down.data.schedule.attempt` and `ns.down.join.schedule.attempt` events, which are triggered when Network Server attempts to schedule a respective downlink on Gateway Server.
- `ns.down.data.schedule.success` and `ns.down.join.schedule.success` events, which are triggered when Network Server successfully schedules a respective downlink on Gateway Server.
- `ns.down.data.schedule.fail` and `ns.down.join.schedule.fail` events, which are triggered when Network Server fails to schedule a respective downlink on Gateway Server.
- Specify gRPC port and OAuth server address when generating a CLI config file with `ttn-lw-cli use` (see `--grpc-port` and `--oauth-server-address` options).
- Guide to connect MikroTik Routerboard

### Changed

- Styling improvements to webhook and pubsub table in Console.
- Gateway location is updated even if no antenna locations had been previously set.
- Renamed `ns.application.begin_link` event to `ns.application.link.begin`.
- Renamed `ns.application.end_link` event to `ns.application.link.end`.
- `ns.up.data.drop` and `ns.up.join.drop` events are now triggered when respective uplink duplicate is dropped by Network Server.
- Network Server now drops FPort 0 data uplinks with non-empty FOpts.
- Frontend asset hashes are loaded dynamically from a manifest file instead of being built into the stack binary.
- Removed `Cache-Control` header for static files.
- Sort events by `time` in the Console.
- Restructure doc folder

### Removed

- `ns.up.merge_metadata` event.
- `ns.up.receive_duplicate` event.
- `ns.up.receive` event.

### Fixed

- End device claim display bug when claim dates not set.
- DeviceModeInd handling for LoRaWAN 1.1 devices.
- Do not perform unnecessary gateway location updates.
- Error display on failed end device import in the Console.
- Update password view not being accessible
- FOpts encryption and decryption for LoRaWAN 1.1 devices.
- Application Server returns an error when trying to delete a device that does not exist.
- Network Server returns an error when trying to delete a device that does not exist.
- Retrieve LNS Trust without LNS Credentials attribute.
- Too strict webhook base URL validation in the Console.
- Webhook and PubSub total count in the Console.
- DevEUI is set when creating ABP devices via CLI.
- CLI now shows all supported enum values for LoraWAN fields.
- Application Server does not crash when retrieving a webhook template that does not exist if no template repository has been configured.
- Application Server does not crash when listing webhook templates if no template repository has been configured.
- Error display on failed end device fetching in the Console.
- Various inconsistencies with Regional Parameters specifications.

## [3.7.0] - 2020-04-02

### Added

- Update gateway antenna location from incoming status message (see `update_location_from_status` gateway field and `--gs.update-gateway-location-debounce-time` option).
  - This requires a database migration (`ttn-lw-stack is-db migrate`) because of the added columns.
- Access Tokens are now linked to User Sessions.
  - This requires a database migration (`ttn-lw-stack is-db migrate`) because of the added columns.
- Edit application attributes in Application General Settings in the Console
- New `use` CLI command to automatically generate CLI configuration files.
- View/edit `update_location_from_status` gateway property using the Console.

### Changed

- Default DevStatus periodicity is increased, which means that, by default, DevStatusReq will be scheduled less often.
- Default class B and C timeouts are increased, which means that, by default, if the Network Server expects an uplink from the device after a downlink, it will wait longer before rescheduling the downlink.
- In case downlink frame carries MAC requests, Network Server will not force the downlink to be sent confirmed in class B and C.

### Fixed

- Fix organization collaborator view not being accessible in the Console.
- Error display on Data pages in the Console.
- Fix too restrictive MQTT client validation in PubSub form in the Console.
- Fix faulty display of device event stream data for end devices with the same ID in different applications.
- Trailing slashes handling in webhook paths.
- End device location display bug when deleting the location entry in the Console.
- GS could panic when gateway connection stats were updated while updating the registry.
- Local CLI and stack config files now properly override global config.
- Error display on failed end device deletion in the Console.

## [3.6.3] - 2020-03-30

### Fixed

- Limited throughput in upstream handlers in Gateway Server when one gateway's upstream handler is busy.

## [3.6.2] - 2020-03-19

### Fixed

- Entity events subscription release in the Console (Firefox).
- RekeyInd handling for LoRaWAN 1.1 devices.
- Network server deduplication Redis configuration.
- Change the date format in the Console to be unambiguous (`17 Mar, 2020`).
- Handling of uplink frame counters exceeding 65535.
- Gateway events subscription release in the Console.
- Panic when receiving a UDP `PUSH_DATA` frame from a gateway without payload.

### Security

- Admin users that are suspended can no longer create, view or delete other users.

## [3.6.1] - 2020-03-13

### Added

- New `list` and `request-validation` subcommands for the CLI's `contact-info` commands.
- Device Claim Authentication Code page in the Console.
- Gateway Server rate limiting support for the UDP frontend, see (`--gs.udp.rate-limiting` options).
- Uplink deduplication via Redis in Network Server.

### Changed

- Network and Application Servers now maintain application downlink queue per-session.
- Gateway Server skips setting up an upstream if the DevAddr prefixes to forward are empty.
- Gateway connection stats are now cached in Redis (see `--cache.service` and `--gs.update-connections-stats-debounce-time` options).

### Fixed

- Telemetry and events for gateway statuses.
- Handling of downlink frame counters exceeding 65535.
- Creating 1.0.4 ABP end devices via the Console.
- ADR uplink handling.
- Uplink retransmission handling.
- Synchronizing Basic Station concentrator time after reconnect or initial connect after long inactivity.

### Security

- Changing username and password to be not required in pubsub integration.

## [3.6.0] - 2020-02-27

### Added

- Class B support.
- WebSocket Ping-Pong support for Basic Station frontend in the Gateway Server.
- LoRaWAN 1.0.4 support.

### Changed

- Do not use `personal-files` plugin for Snap package.
- Network Server will never attempt RX1 for devices with `Rx1Delay` of 1 second.
- Improved efficiency of ADR MAC commands.
- Gateway Configuration Server will use the default WebSocket TLS port if none is set.

### Fixed

- End device events subscription release in the Console.
- Blocking UDP packet handling while the gateway was still connecting. Traffic is now dropped while the connection is in progress, so that traffic from already connected gateways keep flowing.
- Join-request transmission parameters.
- ADR in 72-channel regions.
- Payload length limits used by Network Server being too low.
- CLI ignores default config files that cannot be read.
- Device creation rollback potentially deleting existing device with same ID.
- Returned values not representing the effective state of the devices in Network Server when deprecated field paths are used.
- Downlink queue operations in Network Server for LoRaWAN 1.1 devices.

## [3.5.3] - 2020-02-14

### Added

- Display of error payloads in console event log.
- Zero coordinate handling in location form in the Console.

### Fixed

- Updating `supports_class_c` field in the Device General Settings Page in the Console.
- Updating MQTT pubsub configuration in the Console
- Handling multiple consequent updates of MQTT pubsub/webhook integrations in the Console.
- Displaying total device count in application overview section when using device search in the Console
- FQDN used for Backend Interfaces interoperability requests.
- Exposing device sensitive fields to unrelated stack components in the Console.
- CLI trying to read input while none available.
- Reconnections of gateways whose previous connection was not cleaned up properly. New connections from the same gateway now actively disconnects existing connections.
- `ttn-lw-stack` and `ttn-lw-cli` file permission errors when installed using snap.
  - You may need to run `sudo snap connect ttn-lw-stack:personal-files`
- Changing username and password to be not required in pubsub integration

### Security

## [3.5.2] - 2020-02-06

### Fixed

- Channel mask encoding in LinkADR MAC command.
- Frequency plan validation in Network Server on device update.
- Authentication of Basic Station gateways.

## [3.5.1] - 2020-01-29

### Added

- Responsive side navigation (inside entity views) to the Console.
- Overall responsiveness of the Console.
- Support for configuring Redis connection pool sizes with `redis.pool-size` options.

### Fixed

- Crashes on Gateway Server start when traffic flow started while The Things Stack was still starting.
- Not detecting session change in Application Server when interop Join Server did not provide a `SessionKeyID`.

## [3.5.0] - 2020-01-24

### Added

- Support for releasing gateway EUI after deletion.
- Support in the Application Server for the `X-Downlink-Apikey`, `X-Downlink-Push` and `X-Downlink-Replace` webhook headers. They allow webhook integrations to determine which endpoints to use for downlink queue operations.
- `as.webhooks.downlinks.public-address` and `as.webhooks.downlinks.public-tls-address` configuration options to the Application Server.
- Support for adjusting the time that the Gateway Server schedules class C messages in advance per gateway.
  - This requires a database migration (`ttn-lw-stack is-db migrate`) because of the added columns.
- `end-devices use-external-join-server` CLI subcommand, which disassociates and deletes the device from Join Server.
- `mac_settings.beacon_frequency` end device field, which defines the default frequency of class B beacon in Hz.
- `mac_settings.desired_beacon_frequency` end device field, which defines the desired frequency of class B beacon in Hz that will be configured via MAC commands.
- `mac_settings.desired_ping_slot_data_rate_index` end device field, which defines the desired data rate index of the class B ping slot that will be configured via MAC commands.
- `mac_settings.desired_ping_slot_frequency` end device field, which defines the desired frequency of the class B ping slot that will be configured via MAC commands.
- Mobile navigation menu to the Console.
- View and edit all Gateway settings from the Console.
- `skip_payload_crypto` end device field, which makes the Application Server skip decryption of uplink payloads and encryption of downlink payloads.
- `app_s_key` and `last_a_f_cnt_down` uplink message fields, which are set if the `skip_payload_crypto` end device field is true.
- Support multiple frequency plans for a Gateway.
- Entity search by ID in the Console.

### Changed

- `resets_join_nonces` now applies to pre-1.1 devices as well as 1.1+ devices.
- Empty (`0x0000000000000000`) JoinEUIs are now allowed.

### Fixed

- Respect stack components on different hosts when connected to event sources in the Console.
- Pagination of search results.
- Handling OTAA devices registered on an external Join Server in the Console.
- RxMetadata Location field from Gateway Server.
- Channel mask encoding in LinkADR MAC command.
- Device location and payload formatter form submits in the Console.
- Events processing in the JS SDK.
- Application Server frontends getting stuck after their associated link is closed.

## [3.4.2] - 2020-01-08

### Added

- Forwarding of backend warnings to the Console.
- Auth Info service to the JavaScript SDK.
- Subscribable events to the JavaScript SDK.
- Include `gateway_ID` field in Semtech UDP configuration response from Gateway Configuration Server.
- Sorting feature to entity tables in the Console.

### Changed

- Increase time that class C messages are scheduled in advance from 300 to 500 ms to support higher latency gateway backhauls.

### Fixed

- Fix selection of pseudo wildcard rights being possible (leading to crash) in the Console even when such right cannot be granted.
- Fix loading spinner being stuck infinitely in gateway / application / organization overview when some rights aren't granted to the collaborator.
- Fix deadlock of application add form in the Console when the submit results in an error.
- Fix ttn-lw-cli sometimes refusing to update Gateway EUI.

## [3.4.1] - 2019-12-30

### Added

- Support for ordering in `List` RPCs.
- Detect existing Basic Station time epoch when the gateway was already running long before it (re)connected to the Gateway Server.

### Changed

- Reduce the downlink path expiry window to 15 seconds, i.e. typically missing three `PULL_DATA` frames.
- Reduce the connection expiry window to 1 minute.
- Reduce default UDP address block time from 5 minutes to 1 minute. This allows for faster reconnecting if the gateway changes IP address. The downlink path and connection now expire before the UDP source address is released.

### Fixed

- Fix class A downlink scheduling when an uplink message has been received between the triggering uplink message.

## [3.4.0] - 2019-12-24

### Added

- Downlink queue operation topics in the PubSub integrations can now be configured using the Console.
- `List` RPC in the user registry and related messages.
- User management for admins in the Console.
- `users list` command in the CLI.
- Support for getting Kerlink CPF configurations from Gateway Configuration Server.
- Support for Microchip ATECC608A-TNGLORA-C manifest files in device template conversion.

### Fixed

- Fix the PubSub integration edit page in the Console.
- Fix updating and setting of webhook headers in the Console.
- Fix DevNonce checks for LoRaWAN 1.0.3.

## [3.3.2] - 2019-12-04

### Added

- Support for selecting gateways when queueing downlinks via CLI (see `class-b-c.gateways` option).
- Options `is.oauth.ui.branding-base-url` and `console.ui.branding-base-url` that can be used to customize the branding (logos) of the web UI.
- Email templates can now also be loaded from blob buckets.
- Support for pagination in search APIs.
- Search is now also available to non-admin users.
- Support for searching end devices within an application.
- Notification during login informing users of unapproved user accounts.
- Support maximum EIRP value from frequency plans sub-bands.
- Support duty-cycle value from frequency plans sub-bands.

### Changed

- Allow enqueuing class B/C downlinks regardless of active device class.

### Fixed

- Fix crashing of organization collaborator edit page.
- Avoid validating existing queue on application downlink pushes.
- Correct `AU_915_928` maximum EIRP value to 30 dBm in 915.0 – 928.0 MHz (was 16.15 dBm).
- Correct `US_902_928` maximum EIRP value to 23.15 dBm in 902.3 – 914.9 MHz (was 32.15 dBm) and 28.15 dBm in 923.3 – 927.5 MHz (was 32.15 dBm). This aligns with US915 Hybrid Mode.
- Correct `AS_923` maximum EIRP value to 16 dBm in 923.0 – 923.5 MHz (was 16.15 dBm).

### Security

- Keep session keys separate by `JoinEUI` to avoid conditions where session keys are retrieved only by `DevEUI` and the session key identifier. This breaks retrieving session keys of devices that have been activated on a deployment running a previous version. Since the Application Server instances are currently in-cluster, there is no need for an Application Server to retrieve the `AppSKey` from the Join Server, making this breaking change ineffective.

## [3.3.1] - 2019-11-26

### Added

- Add support for Redis Sentinel (see `redis.failover.enable`, `redis.failover.master-name`, `redis.failover.addresses` options).

### Fixed

- Fix `AppKey` decryption in Join Server.

### Security

## [3.3.0] - 2019-11-25

### Added

- Add support for encrypting device keys at rest (see `as.device-kek-label`, `js.device-kek-label` and `ns.device-kek-label` options).
- The Network Server now provides the timestamp at which it received join-accept or data uplink messages.
- Add more details to logs that contain errors.
- Support for end device pictures in the Identity Server.
  - This requires a database migration (`ttn-lw-stack is-db migrate`) because of the added columns.
- Support for end device pictures in the CLI.

### Fixed

- Fix an issue causing unexpected behavior surrounding login, logout and token management in the Console.
- Fix an issue causing the application link page of the Console to load infinitely.

## [3.2.6] - 2019-11-18

### Fixed

- Fix active application link count being limited to 10 per CPU.
- The Application Server now fills the timestamp at which it has received uplinks from the Network Server.

## [3.2.5] - 2019-11-15

### Added

- Support for creating applications and gateway with an organization as the initial owner in the Console.
- Hide views and features in the Console that the user and stack configuration does not meet the necessary requirements for.
- Full range of Join EUI prefixes in the Console.
- Support specifying the source of interoperability server client CA configuration (see `interop.sender-client-ca.source` and related fields).

### Changed

- Reading and writing of session keys in Application and Network server registries now require device key read and write rights respectively.
- Implement redesign of entity overview title sections to improve visual consistency.

### Deprecated

- `--interop.sender-client-cas` in favor of `--interop.sender-client-ca` sub-fields in the stack.

### Fixed

- Fix gateway API key forms being broken in the Console.
- Fix MAC command handling in retransmissions.
- Fix multicast device creation issues.
- Fix device key unwrapping.
- Fix setting gateway locations in the Console.

### Security

## [3.2.4] - 2019-11-04

### Added

- Support LoRa Alliance TR005 Draft 3 QR code format.
- Connection indicators in Console's gateway list.
- TLS support for application link in the Console.
- Embedded documentation served at `/assets/doc`.

### Fixed

- Fix device creation rollback potentially deleting existing device with same ID.
- Fix missing transport credentials when using external NS linking.

### Security

## [3.2.3] - 2019-10-24

### Added

- Emails when the state of a user or OAuth client changes.
- Option to generate claim authentication codes for devices automatically.
- User invitations can now be sent and redeemed.
- Support for creating organization API keys in the Console.
- Support for deleting organization API keys in the Console.
- Support for editing organization API keys in the Console.
- Support for listing organization API keys in the Console.
- Support for managing organization API keys and rights in the JS SDK.
- Support for removing organization collaborators in the Console.
- Support for editing organization collaborators in the Console.
- Support for listing organization collaborators in the Console.
- Support for managing organization collaborators and rights in the JS SDK.
- MQTT integrations page in the Console.

### Changed

- Rename "bulk device creation" to "import devices".
- Move device import button to the end device tables (and adapt routing accordingly).
- Improve downlink performance.

### Fixed

- Fix issues with device bulk creation in Join Server.
- Fix device import not setting component hosts automatically.
- Fix NewChannelReq scheduling condition.
- Fix publishing events for generated MAC commands.
- Fix saving changes to device general settings in the Console.

## [3.2.2] - 2019-10-14

### Added

- Initial API and CLI support for LoRaWAN application packages and application package associations.
- New documentation design.
- Support for ACME v2.

### Deprecated

- Deprecate the `tls.acme.enable` setting. To use ACME, set `tls.source` to `acme`.

### Fixed

- Fix giving priority to ACME settings to remain backward compatible with configuration for `v3.2.0` and older.

## [3.2.1] - 2019-10-11

### Added

- `support-link` URI config to the Console to show a "Get Support" button.
- Option to explicitly enable TLS for linking of an Application Server on an external Network Server.
- Service to list QR code formats and generate QR codes in PNG format.
- Status message forwarding functions to upstream host/s.
- Support for authorizing device claiming on application level through CLI. See `ttn-lw-cli application claim authorize --help` for more information.
- Support for claiming end devices through CLI. See `ttn-lw-cli end-device claim --help` for more information.
- Support for converting Microchip ATECC608A-TNGLORA manifest files to device templates.
- Support for Crypto Servers that do not expose device root keys.
- Support for generating QR codes for claiming. See `ttn-lw-cli end-device generate-qr --help` for more information.
- Support for storage of frequency plans, device repository and interoperability configurations in AWS S3 buckets or GCP blobs.

### Changed

- Enable the V2 MQTT gateway listener by default on ports 1881/8881.
- Improve handling of API-Key and Collaborator rights in the console.

### Fixed

- Fix bug with logout sometimes not working in the console.
- Fix not respecting `RootCA` and `InsecureSkipVerify` TLS settings when ACME was configured for requesting TLS certificates.
- Fix reading configuration from current, home and XDG directories.

## [3.2.0] - 2019-09-30

### Added

- A map to the overview pages of end devices and gateways.
- API to retrieve MQTT configurations for applications and gateways.
- Application Server PubSub integrations events.
- `mac_settings.desired_max_duty_cycle`, `mac_settings.desired_adr_ack_delay_exponent` and `mac_settings.desired_adr_ack_limit_exponent` device flags.
- PubSub integrations to the console.
- PubSub service to JavaScript SDK.
- Support for updating `mac_state.desired_parameters`.
- `--tls.insecure-skip-verify` to skip certificate chain verification (insecure; for development only).

### Changed

- Change the way API key rights are handled in the `UpdateAPIKey` rpc for Applications, Gateways, Users and Organizations. Users can revoke or add rights to api keys as long as they have these rights.
- Change the way collaborator rights are handled in the `SetCollaborator` rpc for Applications, Gateways, Clients and Organizations. Collaborators can revoke or add rights to other collaborators as long as they have these rights.
- Extend device form in the Console to allow creating OTAA devices without root keys.
- Improve confirmed downlink operation.
- Improve gateway connection status indicators in Console.
- Upgrade Gateway Configuration Server to a first-class cluster role.

### Fixed

- Fix downlink length computation in the Network Server.
- Fix implementation of CUPS update-info endpoint.
- Fix missing CLI in `deb`, `rpm` and Snapcraft packages.

## [3.1.2] - 2019-09-05

### Added

- `http.redirect-to-host` config to redirect all HTTP(S) requests to the same host.
- `http.redirect-to-tls` config to redirect HTTP requests to HTTPS.
- Organization Create page in the Console.
- Organization Data page to the console.
- Organization General Settings page to the console.
- Organization List page.
- Organization Overview page to the console.
- Organizations service to the JS SDK.
- `create` method in the Organization service in the JS SDK.
- `deleteById` method to the Organization service in the JS SDK.
- `getAll` method to the Organizations service.
- `getAll` method to the Organization service in the JS SDK.
- `getById` method to the Organization service in the JS SDK.
- `openStream` method to the Organization service in the JS SDK.
- `updateById` method to the Organization service in the JS SDK.

### Changed

- Improve compatibility with various Class C devices.

### Fixed

- Fix root-relative OAuth flows for the console.

## [3.1.1] - 2019-08-30

### Added

- `--tls.acme.default-host` flag to set a default (fallback) host for connecting clients that do not use TLS-SNI.
- AS-ID to validate the Application Server with through the Common Name of the X.509 Distinguished Name of the TLS client certificate. If unspecified, the Join Server uses the host name from the address.
- Defaults to `ttn-lw-cli clients create` and `ttn-lw-cli users create`.
- KEK labels for Network Server and Application Server to use to wrap session keys by the Join Server. If unspecified, the Join Server uses a KEK label from the address, if present in the key vault.
- MQTT PubSub support in the Application Server. See `ttn-lw-cli app pubsub set --help` for more details.
- Support for external email templates in the Identity Server.
- Support for Join-Server interoperability via Backend Interfaces specification protocol.
- The `generateDevAddress` method in the `Ns` service.
- The `Js` service to the JS SDK.
- The `listJoinEUIPrefixes` method in the `Js` service.
- The `Ns` service to the JS SDK.
- The new The Things Stack branding.
- Web interface for changing password.
- Web interface for requesting temporary password.

### Changed

- Allow admins to create temporary passwords for users.
- CLI-only brew tap formula is now available as `TheThingsNetwork/lorawan-stack/ttn-lw-cli`.
- Improve error handling in OAuth flow.
- Improve getting started guide for a deployment of The Things Stack.
- Optimize the way the Identity Server determines memberships and rights.

### Deprecated

- `--nats-server-url` in favor of `--nats.server-url` in the PubSub CLI support.

### Removed

- `ids.dev_addr` from allowed field masks for `/ttn.lorawan.v3.NsEndDeviceRegistry/Set`.
- Auth from CLI's `forgot-password` command and made it optional on `update-password` command.
- Breadcrumbs from Overview, Application and Gateway top-level views.

### Fixed

- Fix `grants` and `rights` flags of `ttn-lw-cli clients create`.
- Fix a bug that resulted in events streams crashing in the console.
- Fix a bug where uplinks from some Basic Station gateways resulted in the connection to break.
- Fix a security issue where non-admin users could edit admin-only fields of OAuth clients.
- Fix an issue resulting in errors being unnecessarily logged in the console.
- Fix an issue with the `config` command rendering some flags and environment variables incorrectly.
- Fix API endpoints that allowed HTTP methods that are not part of our API specification.
- Fix console handling of configured mount paths other than `/console`.
- Fix handling of `ns.dev-addr-prefixes`.
- Fix incorrect error message in `ttn-lw-cli users oauth` commands.
- Fix propagation of warning headers in API responses.
- Fix relative time display in the Console.
- Fix relative time display in the Console for IE11, Edge and Safari.
- Fix unable to change LoRaWAN MAC and PHY version.
- Resolve flickering display issue in the overview pages of entities in the console.

## [3.1.0] - 2019-07-26

### Added

- `--headers` flag to `ttn-lw-cli applications webhooks set` allowing users to set HTTP headers to add to webhook requests.
- `getByOrganizationId` and `getByUserId` methods to the JS SDK.
- A new documentation system.
- A newline between list items returned from the CLI when using a custom `--output-format` template.
- An `--api-key` flag to `ttn-lw-cli login` that allows users to configure the CLI with a more restricted (Application, Gateway, ...) API key instead of the usual "all rights" OAuth access token.
- API for getting the rights of a single collaborator on (member of) an entity.
- Application Payload Formatters Page in the console.
- Class C and Multicast guide.
- CLI support for enabling/disabling JS, GS, NS and AS through configuration.
- Components overview in documentation.
- Device Templates to create, convert and map templates and assign EUIs to create large amounts of devices.
- Downlink Queue Operations guide.
- End device level payload formatters to console.
- Event streaming views for end devices.
- Events to device registries in the Network Server, Application Server and Join Server.
- Functionality to delete end devices in the console.
- Gateway General Settings Page to the console.
- Getting Started guide for command-line utility (CLI).
- Initial overview page to console.
- Native support to the Basic Station LNS protocol in the Gateway Server.
- NS-JS and AS-JS Backend Interfaces 1.0 and 1.1 draft 3 support.
- Option to revoke user sessions and access tokens on password change.
- Support for NS-JS and AS-JS Backend Interfaces.
- Support for URL templates inside the Webhook paths ! The currently supported fields are `appID`, `appEUI`, `joinEUI`, `devID`, `devEUI` and `devAddr`. They can be used using RFC 6570.
- The `go-cloud` integration to the Application Server. See `ttn-lw-cli applications pubsubs --help` for more details.
- The `go-cloud` integration to the Application Server. This integration enables downlink and uplink messaging using the cloud pub-sub by setting up the `--as.pubsub.publish-urls` and `--as.pubsub.subscribe-urls` parameters. You can specify multiple publish endpoints or subscribe endpoints by repeating the parameter (i.e. `--as.pubsub.publish-urls url1 --as.pubsub.publish-urls url2 --as.pubsub.subscribe-urls url3`).
- The Gateway Data Page to the console.
- View to update the antenna location information of gateways.
- View to update the location information of end devices.
- Views to handle integrations (webhooks) to the console.
- Working with Events guide.

### Changed

- Change database index names for invitation and OAuth models. Existing databases are migrated automatically.
- Change HTTP API for managing webhooks to avoid conflicts with downlink webhook paths.
- Change interpretation of frequency plan's maximum EIRP from a ceiling to a overriding value of any band (PHY) settings.
- Change the prefix of Prometheus metrics from `ttn_` to `ttn_lw_`.
- Rename the label `server_address` of Prometheus metrics `grpc_client_conns_{opened,closed}_total` to `remote_address`
- Resolve an issue where the stack complained about sending credentials on insecure connections.
- The Events endpoint no longer requires the `_ALL` right on requested entities. All events now have explicit visibility rules.

### Deprecated

- `JsEndDeviceRegistry.Provision()` rpc. Please use `EndDeviceTemplateConverter.Convert()` instead.

### Removed

- Remove the address label from Prometheus metric `grpc_server_conns_{opened,closed}_total`.

### Fixed

- Fix Basic Station CUPS LNS credentials blob.
- Fix a leak of entity information in List RPCs.
- Fix an issue that resulted in some event errors not being shown in the console.
- Fix an issue where incorrect error codes were returned from the console's OAuth flow.
- Fix clearing component addresses on updating end devices through CLI.
- Fix CLI panic for invalid attributes.
- Fix crash when running some `ttn-lw-cli organizations` commands without `--user-id` flag.
- Fix dwell-time issues in AS923 and AU915 bands.
- Fix occasional issues with downlink payload length.
- Fix the `x-total-count` header value for API Keys and collaborators.
- Fix the error that is returned when deleting a collaborator fails.

### Security

- Update node packages to fix known vulnerabilities.

## [3.0.4] - 2019-07-10

### Fixed

- Fix rights caching across multiple request contexts.

## [3.0.3] - 2019-05-10

### Added

- Support for getting automatic Let's Encrypt certificates. Add the new config flags `--tls.acme.enable`, `--tls.acme.dir=/path/to/storage`, `--tls.acme.hosts=example.com`, `--tls.acme.email=you@example.com` flags (or their env/config equivalent) to make it work. The `/path/to/storage` dir needs to be `chown`ed to `886:886`. See also `docker-compose.yml`.
- `GetApplicationAPIKey`, `GetGatewayAPIKey`, `GetOrganizationAPIKey`, `GetUserAPIKey` RPCs and related messages.
- "General Settings" view for end devices.
- `--credentials-id` flag to CLI that allows users to be logged in with mulitple credentials and switch between them.
- A check to the Identity Server that prevents users from deleting applications that still contain end devices.
- Application Collaborators management to the console.
- Checking maximum round-trip time for late-detection in downlink scheduling.
- Configuration service to JS SDK.
- Device list page to applications in console.
- Events to the application management pages.
- Round-trip times to Gateway Server connection statistics.
- Support for the value `cloud` for the `--events.backend` flag. When this flag is set, the `--events.cloud.publish-url` and `--events.cloud.subscribe-url` are used to set up a cloud pub-sub for events.
- Support for uplink retransmissions.
- Using median round-trip time value for absolute time scheduling if the gateway does not have GPS time.

### Changed

- Change encoding of keys to hex in device key generation (JS SDK).
- Change interpretation of absolute time in downlink messages from time of transmission to time of arrival.
- Improve ADR algorithm performance.
- Improve ADR performance.
- Make late scheduling default for gateways connected over UDP to avoid overwriting queued downlink.
- Make sure that non-user definable fields of downlink messages get discarded across all Application Server frontends.
- Prevent rpc calls to JS when the device has `supports_join` set to `false` (JS SDK).
- Update the development tooling. If you are a developer, make sure to check the changes in CONTRIBUTING.md and DEVELOPMENT.md.

### Fixed

- Fix `AppAs` not registered for HTTP interfacing while it is documented in the API.
- Fix absolute time scheduling with UDP connected gateways
- Fix authentication of MQTT and gRPC connected gateways
- Fix connecting MQTT V2 gateways
- Fix faulty composition of default values with provided values during device creation (JS SDK)
- Fix preserving user defined priority for application downlink
- Fix UDP downlink format for older forwarders
- Fix usage of `URL` class in browsers (JS SDK)

## [3.0.2] - 2019-04-12

### Changed

- Upgrade Go to 1.12

### Fixed

- Fix streaming events over HTTP with Gzip enabled.
- Fix resetting downlink channels for US, AU and CN end devices.
- Fix rendering of enums in JSON.
- Fix the permissions of our Snap package.

## [3.0.1] - 2019-04-10

### Added

- `dev_addr` to device fetched from the Network Server.
- `received_at` to `ApplicationUp` messages.
- `ttn-lw-cli users oauth` commands.
- Event payload to `as.up.forward`, `as.up.drop`, `as.down.receive`, `as.down.forward` and `as.down.drop` events.
- Event payload to `gs.status.receive`, `gs.up.receive` and `gs.down.send` events.
- OAuth management in the Identity Server.

### Changed

- Document places in the CLI where users can use arguments instead of flags.
- In JSON, LoRaWAN AES keys are now formatted as Hex instead of Base64.
- Make device's `dev_addr` update when the session's `dev_addr` is updated.

### Removed

- Remove end device identifiers from `DownlinkMessage` sent from the Network Server to the Gateway Server.

### Fixed

- Fix `dev_addr` not being present in upstream messages.

<!--
NOTE: These links should respect backports. See https://github.com/TheThingsNetwork/lorawan-stack/pull/1444/files#r333379706.
-->

[unreleased]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.8.4...HEAD
[3.8.4]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.8.3...v3.8.4
[3.8.3]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.8.2...v3.8.3
[3.8.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.7.2...v3.8.2
[3.7.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.7.0...v3.7.2
[3.7.0]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.6.0...v3.7.0
[3.6.3]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.6.2...v3.6.3
[3.6.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.6.1...v3.6.2
[3.6.1]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.6.0...v3.6.1
[3.6.0]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.5.3...v3.6.0
[3.5.3]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.5.2...v3.5.3
[3.5.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.5.1...v3.5.2
[3.5.1]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.5.0...v3.5.1
[3.5.0]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.4.2...v3.5.0
[3.4.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.4.1...v3.4.2
[3.4.1]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.4.0...v3.4.1
[3.4.0]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.3.2...v3.4.0
[3.3.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.3.1...v3.3.2
[3.3.1]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.3.0...v3.3.1
[3.3.0]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.2.6...v3.3.0
[3.2.6]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.2.5...v3.2.6
[3.2.5]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.2.4...v3.2.5
[3.2.4]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.2.3...v3.2.4
[3.2.3]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.2.2...v3.2.3
[3.2.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.2.1...v3.2.2
[3.2.1]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.1.2...v3.2.0
[3.1.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.1.1...v3.1.2
[3.1.1]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.1.0...v3.1.1
[3.1.0]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.0.4...v3.1.0
[3.0.4]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.0.3...v3.0.4
[3.0.3]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.0.2...v3.0.3
[3.0.2]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.0.1...v3.0.2
[3.0.1]: https://github.com/TheThingsNetwork/lorawan-stack/compare/v3.0.0...v3.0.1
