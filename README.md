# The Things Stack, an open source LoRaWAN Network Server

The Things Stack is an open source LoRaWAN network stack suitable for large, global and geo-distributed public and private networks as well as smaller networks. The architecture follows the LoRaWAN Network Reference Model for standards compliancy and interoperability. This project is actively maintained by [The Things Industries](https://www.thethingsindustries.com).

LoRaWAN is a protocol for low-power wide area networks. It allows for large scale Internet of Things deployments where low-powered devices efficiently communicate with Internet-connected applications over long range wireless connections.

## Features

- LoRaWAN Network Server
  - [x] Supports LoRaWAN 1.0
  - [x] Supports LoRaWAN 1.0.1
  - [x] Supports LoRaWAN 1.0.2
  - [x] Supports LoRaWAN 1.0.3
  - [x] Supports LoRaWAN 1.0.4
  - [x] Supports LoRaWAN 1.1
  - [x] Supports LoRaWAN Regional Parameters 1.0
  - [x] Supports LoRaWAN Regional Parameters 1.0.2 rev B
  - [x] Supports LoRaWAN Regional Parameters 1.0.3 rev A
  - [x] Supports LoRaWAN Regional Parameters 1.1 rev A
  - [x] Supports LoRaWAN Regional Parameters 1.1 rev B
  - [x] Supports Class A devices
  - [x] Supports Class B devices
  - [x] Supports Class C devices
  - [x] Supports OTAA devices
  - [x] Supports ABP devices
  - [x] Supports MAC Commands
  - [x] Supports Adaptive Data Rate
  - [ ] Implements LoRaWAN Back-end Interfaces 1.0
- LoRaWAN Application Server
  - [x] Payload conversion of well-known payload formats
  - [x] Payload conversion using custom JavaScript functions
  - [x] MQTT pub/sub API
  - [x] HTTP Webhooks API
  - [ ] Implements LoRaWAN Back-end Interfaces 1.0
- LoRaWAN Join Server
  - [x] Supports OTAA session key derivation
  - [x] Supports external crypto services
  - [x] Implements LoRaWAN Back-end Interfaces 1.0
  - [x] Implements LoRaWAN Back-end Interfaces 1.1 draft 3
- OAuth 2.0 Identity Server
  - [x] User management
  - [x] Entity management
  - [x] ACLs
- GRPC APIs
- HTTP APIs
- Command-Line Interface
  - [x] Create account and login
  - [x] Application management and traffic
  - [x] End device management, status and traffic
  - [x] Gateway management and status
- Web Interface (Console)
  - [x] Create account and login
  - [x] Application management and traffic
  - [x] End device management, status and traffic
  - [x] Gateway management, status and traffic

## Getting Started

You want to **install The Things Stack**? Fantastic! Here's the [Getting Started guide](https://ttn.fyi/v3/getting-started).

Do you want to **set up a local development environment**? See the [DEVELOPMENT.md](DEVELOPMENT.md) for instructions.

Do you want to **contribute to The Things Stack**? Your contributions are welcome! See the guidelines in [CONTRIBUTING.md](CONTRIBUTING.md).

Are you new to LoRaWAN and The Things Network? See the general documentation at [thethingsnetwork.org/docs](https://www.thethingsnetwork.org/docs/).

## Commitments and Releases

Open source projects are great, but a stable and reliable open source ecosystem is even better. Therefore, we make the following commitments:

1. We will not break the API towards gateways and applications within the major version. This includes how gateways communicate (with Gateway Server) and how applications work with data (with Application Server)
2. We will upgrade storage from older versions within the major version via migrations. This means that you can migrate an older setup without losing data.
3. We will not require storage migrations within the minor version. This means that you can update patches without database migrations.
4. We will not break the public command-line interface and configuration within the major version. This means that you can safely build scripts and migrate configuration.
5. We will not break the API between components and events within minor versions. So at least the same minor versions of components are compatible with each other.
6. We reserve the right to fix bugs in API, configuration and storage in patches and minor updates. This may break components, gateways and applications that rely on buggy behavior.

As we are continuously adding functionality and fixes in new releases, we are also introducing new configurations and new defaults. We therefore recommend reading the release notes before upgrading to a new version.

You can find the releases and their notes on the [Releases page](https://github.com/TheThingsNetwork/lorawan-stack/releases).

## Support

- The [forum](https://www.thethingsnetwork.org/forum) contains a large amount of information and has great search support. We have a special [category for The Things Stack](https://www.thethingsnetwork.org/forum/c/network-and-routing/v3).
- You can chat in the [#the-things-stack channel on Slack](https://thethingsnetwork.slack.com/messages/CFVF7R4AH/). If you don't have a Slack account yet, you can create one by going to [ttn.fyi/slack-invite](https://ttn.fyi/slack-invite).
- Hosted solutions, as well as commercial support and consultancy are offered by [The Things Industries](https://www.thethingsindustries.com).
