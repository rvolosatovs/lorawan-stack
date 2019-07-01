# Getting Started with The Things Network Stack for LoRaWAN

## Introduction

This is a guide for setting up a private LoRaWAN Network Server using The Things Network Stack for LoRaWAN.

In this guide we will get everything up and running on a server using Docker. If you are comfortable with configuring servers and working with command line, this is the perfect place to start.
 
## Table of Contents

1. [Dependencies](#dependencies)
2. [Configuration](#configuration)
3. [Running the stack](#running)
4. [Login using the CLI](#login)
5. [Creating a gateway](#creategtw)
6. [Creating an application](#createapp)
7. [Creating a device](#createdev)
8. [Linking the application](#linkappserver)
9. [Using the MQTT server](#mqtt)
10. [Using webhooks](#webhooks)
11. [Advanced: Events](#events)

## <a name="dependencies">Dependencies</a>

### CLI and stack

The web interface Console is not yet available. So in this tutorial, we use the command-line interface (CLI) to manage the setup.

You can use the CLI on your local machine or on the server.

>Note: if you need help with any CLI command, use the `--help` flag to get a list of subcommands, flags and their description and aliases.

#### Package managers (recommended)

##### macOS

```bash
$ brew install TheThingsNetwork/lorawan-stack/ttn-lw-stack
```

##### Linux

```bash
$ sudo snap install ttn-lw-stack
$ sudo snap alias ttn-lw-stack.ttn-lw-cli ttn-lw-cli
```

#### Binaries

If your operating system or package manager is not mentioned, please [download binaries](https://github.com/TheThingsNetwork/lorawan-stack/releases) for your operating system and processor architecture.

### Certificates

By default, the stack requires a `cert.pem` and `key.pem`, in order to to serve content over TLS.

Typically you'll get these from a trusted Certificate Authority. Use the "full chain" for `cert.pem` and the "private key" for `key.pem`. The stack also has support for automated certificate management (ACME). This allows you to easily get trusted TLS certificates for your server from [Let's Encrypt](https://letsencrypt.org/getting-started/). If you want this, you'll need to create an `acme` directory that the stack can write in:

```bash
$ mkdir ./acme
$ chown 886:886 ./acme
```

> If you don't do this, you'll get an error saying something like `open /var/lib/acme/acme_account+key<...>: permission denied`.

For local (development) deployments, you can generate self-signed certificates. If you have your [Go environment](../DEVELOPMENT.md#development-environment) set up, you can run the following command to generate a key and certificate for `localhost`:

```bash
$ go run $(go env GOROOT)/src/crypto/tls/generate_cert.go -ca -host localhost
```

In order for the user in our Docker container to read these files, you have to change the ownership of the certificate and key:

```bash
$ chown 886:886 ./cert.pem ./key.pem
```

> If you don't do this, you'll get an error saying something like `/run/secrets/key.pem: permission denied`.

Keep in mind that self-signed certificates are not trusted by browsers and operating systems, resulting in warnings and errors such as `certificate signed by unknown authority` or `ERR_CERT_AUTHORITY_INVALID`. In most browsers you can add an exception for your self-signed certificate. You can configure the CLI to trust the certificate as well.

## <a name="configuration">Configuration</a>

The stack can be started for development without passing any configuration. However, there are a lot of things you can configure. See [configuration documentation](config.md) for more information. In this guide we'll set some environment variables in `docker-compose.yml`. These environment variables will configure the stack as a development server on `localhost`. For setting up up a public server or for requesting certificates from an ACME provider such as Let's Encrypt, take a closer look at the comments in `docker-compose.yml`.

Refer to the [networking documentation](networking.md) for the endpoints and ports that the stack uses by default.

### <a name="frequencyplans">Frequency plans</a>

By default, frequency plans are fetched by the stack from a [public GitHub repository](https://github.com/TheThingsNetwork/lorawan-frequency-plans). To configure a local directory in offline environments, see the [configuration documentation](config.md) for more information.

### <a name="cli-config">Command-line interface</a>

The command-line interface has some built-in defaults, but you'll want to create a config file or set some environment variables to point it at your deployment.

The recommended way to configure the CLI is with a `.ttn-lw-cli.yml` in your `$XDG_CONFIG_HOME` or `$HOME` directory. You can also put the config file in a different location, and pass it to the CLI as `-c path/to/config.yml`. In this guide we will use the following configuration file:

```yml
oauth-server-address: https://localhost:8885/oauth

identity-server-grpc-address: localhost:8884
gateway-server-grpc-address: localhost:8884
network-server-grpc-address: localhost:8884
application-server-grpc-address: localhost:8884
join-server-grpc-address: localhost:8884

ca: /path/to/your/cert.pem

log:
  level: info
```

## <a name="running">Running the stack</a>

You can run the stack using Docker or container orchestration solutions using Docker. An example [Docker Compose configuration](../docker-compose.yml) is available to get started quickly.

With the `docker-compose.yml` file in the directory of your terminal prompt, enter the following commands to initialize the database, create the first user `admin`, create the CLI OAuth client and start the stack:

```bash
$ docker-compose pull
$ docker-compose run --rm stack is-db init
$ docker-compose run --rm stack is-db create-admin-user \
  --id admin \
  --email admin@localhost
$ docker-compose run --rm stack is-db create-oauth-client \
  --id cli \
  --name "Command Line Interface" \
  --owner admin \
  --no-secret \
  --redirect-uri 'local-callback' \
  --redirect-uri 'code'
$ docker-compose up
```

## <a name="login">Login using the CLI</a>

The CLI needs to be logged on in order to create gateways, applications, devices and API keys. With the stack running in one terminal session, login with the following command:

```bash
$ ttn-lw-cli login
```

This will open the OAuth login page where you can login with your credentials. Once you logged in in the browser, return to the terminal session to proceed.

If you run this command on a remote machine, pass `--callback=false` to get a link to login on your local machine.

## <a name="creategtw">Creating a gateway</a>

First, list the available frequency plans:

```
$ ttn-lw-cli gateways list-frequency-plans
```

Then, create the first gateway with the chosen frequency plan:

```bash
$ ttn-lw-cli gateways create gtw1 \
  --user-id admin \
  --frequency-plan-id EU_863_870 \
  --gateway-eui 00800000A00009EF \
  --enforce-duty-cycle
```

This creates a gateway `gtw1` with user `admin` as collaborator, frequency plan `EU_863_870`, EUI `00800000A00009EF` and respecting duty-cycle limitations. You can now connect your gateway to the stack.

>Note: The CLI returns the created and updated entities by default in JSON. This can be useful in scripts.

## <a name="createapp">Creating an application</a>

Create the first application:

```bash
$ ttn-lw-cli applications create app1 --user-id admin
```

This creates an application `app1` with the `admin` user as collaborator.

Devices are created within applications.

## <a name="createdev">Creating a device</a>

First, list the available frequency plans and LoRaWAN versions:

```
$ ttn-lw-cli end-devices list-frequency-plans
$ ttn-lw-cli end-devices create --help
```

Then, to create an end device using over-the-air-activation (OTAA):

```bash
$ ttn-lw-cli end-devices create app1 dev1 \
  --dev-eui 0004A30B001C0530 \
  --app-eui 800000000000000C \
  --frequency-plan-id EU_863_870 \
  --root-keys.app-key.key 752BAEC23EAE7964AF27C325F4C23C9A \
  --lorawan-version 1.0.2 \
  --lorawan-phy-version 1.0.2-b
```

This will create a LoRaWAN 1.0.2 end device `dev1` in application `app1`. The end device should now be able to join the private network.

>Note: The `AppEUI` is returned as `join_eui` (V3 uses LoRaWAN 1.1 terminology).

>Hint: You can also pass `--with-root-keys` to have root keys generated.

It is also possible to register an ABP activated device using the `--abp` flag as follows:

```bash
$ ttn-lw-cli end-devices create app1 dev2 \
  --frequency-plan-id EU_863_870 \
  --lorawan-version 1.0.2 \
  --lorawan-phy-version 1.0.2-b \
  --abp \
  --session.dev-addr 00E4304D \
  --session.keys.app-s-key.key A0CAD5A30036DBE03096EB67CA975BAA \
  --session.keys.nwk-s-key.key B7F3E161BC9D4388E6C788A0C547F255
```

>Note: The `NwkSKey` is returned as `f_nwk_s_int_key` (V3 uses LoRaWAN 1.1 terminology).

>Hint: You can also pass `--with-session` to have a session generated.

It is also possible to create a multicast device (an ABP device which can not send uplinks and shares the security session with other devices) using the `--multicast` flag as follows:

```bash
$ ttn-lw-cli end-devices create app1 dev3 \
  --frequency-plan-id EU_863_870 \
  --lorawan-version 1.0.2 \
  --lorawan-phy-version 1.0.2-b \
  --abp \
  --session.dev-addr 00E4304D \
  --session.keys.app-s-key.key A0CAD5A30036DBE03096EB67CA975BAA \
  --session.keys.nwk-s-key.key B7F3E161BC9D4388E6C788A0C547F255 \
  --multicast
```

>Note: The `--multicast` flag can be set only during device creation, and as such can not be turned on or off later.

## <a name="linkappserver">Linking the application</a>

In order to send uplinks and receive downlinks from your device, you must link the Application Server to the Network Server. In order to do this, create an API key for the Application Server:

```bash
$ ttn-lw-cli applications api-keys create \
  --name link \
  --application-id app1 \
  --right-application-link
```

The CLI will return an API key such as `NNSXS.VEEBURF3KR77ZR...`. This API key has only link rights and can therefore only be used for linking.

You can now link the Application Server to the Network Server:

```bash
$ ttn-lw-cli applications link set app1 --api-key NNSXS.VEEBURF3KR77ZR..
```

Your application is now linked. You can now use the builtin MQTT server and webhooks to receive uplink traffic and send downlink traffic.

## <a name="mqtt">Using the MQTT server</a>

In order to use the MQTT server you need to create a new API key to authenticate:

```bash
$ ttn-lw-cli applications api-keys create \
  --name mqtt-client \
  --application-id app1 \
  --right-application-traffic-read \
  --right-application-traffic-down-write
```

>Note: See `--help` to see more rights that your application may need.

You can now login using an MQTT client with the application ID `app1` as user name and the newly generated API key as password.

There are many MQTT clients available. Great clients are `mosquitto_pub` and `mosquitto_sub`, part of [Mosquitto](https://mosquitto.org).

>Tip: when using `mosquitto_sub`, pass the `-d` flag to see the topics messages get published on. For example:
>
>`$ mosquitto_sub -h localhost -t '#' -u app1 -P 'NNSXS.VEEBURF3KR77ZR..' -d`

### Subscribing to messages

The Application Server publishes on the following topics:

- `v3/{application id}/devices/{device id}/join`
- `v3/{application id}/devices/{device id}/up`
- `v3/{application id}/devices/{device id}/down/queued`
- `v3/{application id}/devices/{device id}/down/sent`
- `v3/{application id}/devices/{device id}/down/ack`
- `v3/{application id}/devices/{device id}/down/nack`
- `v3/{application id}/devices/{device id}/down/failed`

While you could subscribe to separate topics, for the tutorial subscribe to `#` to subscribe to all messages.

With your MQTT client subscribed, when a device joins the network, a `join` message gets published. For example, for a device ID `dev1`, the message will be published on the topic `v3/app1/devices/dev1/join` with the following contents:

```json
{
  "end_device_ids": {
    "device_id": "dev1",
    "application_ids": {
      "application_id": "app1"
    },
    "dev_eui": "4200000000000000",
    "join_eui": "4200000000000000",
    "dev_addr": "01DA1F15"
  },
  "correlation_ids": [
    "gs:conn:01D2CSNX7FJVKQPCVG612QF1TX",
    "gs:uplink:01D2CT834K2YD17ZWZ6357HC0Z",
    "ns:uplink:01D2CT834KNYD7BT2NHK5R1WVA",
    "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:01D2CT834KJ4AVSD1SJ637NAV6",
    "as:up:01D2CT83AXQFQYQ35SR74CTWKH"
  ],
  "join_accept": {
    "session_key_id": "AWiZpAyXrAfEkUNkBljRoA=="
  }
}
```

You can use the correlation IDs to follow messages as they pass through the stack.

When the device sends an uplink message, a message will be published to the topic `v3/{application id}/devices/{device id}/up`. This message looks like this:

```json
{
  "end_device_ids": {
    "device_id": "dev1",
    "application_ids": {
      "application_id": "app1"
    },
    "dev_eui": "4200000000000000",
    "join_eui": "4200000000000000",
    "dev_addr": "01DA1F15"
  },
  "correlation_ids": [
    "gs:conn:01D2CSNX7FJVKQPCVG612QF1TX",
    "gs:uplink:01D2CV8HF62ME0D7MZWE38HHH8",
    "ns:uplink:01D2CV8HF6FYJHKZ45YY1DB3MR",
    "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:01D2CV8HF6XR7ZFVK768PDG3J4",
    "as:up:01D2CV8HNGJ57G25BW0FCZNY07"
  ],
  "uplink_message": {
    "session_key_id": "AWiZpAyXrAfEkUNkBljRoA==",
    "f_port": 15,
    "frm_payload": "VGVtcGVyYXR1cmUgPSAwLjA=",
    "rx_metadata": [{
      "gateway_ids": {
        "gateway_id": "eui-0242020000247803",
        "eui": "0242020000247803"
      },
      "time": "2019-01-29T13:02:34.981Z",
      "timestamp": 1283325000,
      "rssi": -35,
      "snr": 5,
      "uplink_token": "CiIKIAoUZXVpLTAyNDIwMjAwMDAyNDc4MDMSCAJCAgAAJHgDEMj49+ME"
    }],
    "settings": {
      "data_rate": {
        "lora": {
          "bandwidth": 125000,
          "spreading_factor": 7
        }
      },
      "data_rate_index": 5,
      "coding_rate": "4/6",
      "frequency": "868500000",
      "gateway_channel_index": 2,
      "device_channel_index": 2
    }
  }
}
```

### Scheduling a downlink message

Downlinks can be scheduled by publishing the message to the topic `v3/{application id}/devices/{device id}/down/push`.

For example, to send an unconfirmed downlink message to the device `dev1` in application `app1` with the hexadecimal payload `BE EF` on `FPort` 15 with normal priority, use the topic `v3/app1/devices/dev1/down/push` with the following contents:

```json
{
  "downlinks": [{
    "f_port": 15,
    "frm_payload": "vu8=",
    "priority": "NORMAL",
  }]
}
```

>Hint: Use [this handy tool](https://v2.cryptii.com/hexadecimal/base64) to convert hexadecimal to base64.

>If you use `mosquitto_pub`, use the following command:
>
>`$ mosquitto_pub -h localhost -t 'v3/app1/devices/dev1/down/push' -u app1 -P 'NNSXS.VEEBURF3KR77ZR..' -m '{"downlinks":[{"f_port": 15,"frm_payload":"vu8=","priority": "NORMAL"}]}' -d`

It is also possible to send multiple downlink messages on a single push because `downlinks` is an array. Instead of `/push`, you can also use `/replace` to replace the downlink queue. Replacing with an empty array clears the downlink queue.

>Note: if you do not specify a priority, the default priority `LOWEST` is used. You can specify `LOWEST`, `LOW`, `BELOW_NORMAL`, `NORMAL`, `ABOVE_NORMAL`, `HIGH` and `HIGHEST`.

The stack supports some cool features, such as confirmed downlink with your own correlation IDs. For example, you can push this:

```json
{
  "downlinks": [{
    "f_port": 15,
    "frm_payload": "vu8=",
    "priority": "HIGH",
    "confirmed": true,
    "correlation_ids": ["my-correlation-id"]
  }]
}
```

Once the downlink gets acknowledged, a message is published to the topic `v3/{application id}/devices/{device id}/down/ack`:

```json
{
  "end_device_ids": {
    "device_id": "dev1",
    "application_ids": {
      "application_id": "app1"
    },
    "dev_eui": "4200000000000000",
    "join_eui": "4200000000000000",
    "dev_addr": "00E6F42A"
  },
  "correlation_ids": [
    "my-correlation-id",
    "..."
  ],
  "downlink_ack": {
    "session_key_id": "AWnj0318qrtJ7kbudd8Vmw==",
    "f_port": 15,
    "f_cnt": 11,
    "frm_payload": "vu8=",
    "confirmed": true,
    "priority": "NORMAL",
    "correlation_ids": [
      "my-correlation-id",
      "..."
    ]
  }
}
```

Here you see the correlation ID `my-correlation-id` of your downlink message. You can add multiple custom correlation IDs, for example to reference events or identifiers of your application.

#### Class C unicast

In order to send class C downlink messages to a single device, enable class C support for the end device using the following command:

```bash
$ ttn-lw-cli end-devices update app1 dev1 --supports-class-c
```

This will enable the class C downlink scheduling of the device. That's it! New downlink messages are now scheduled as soon as possible.

To disable class C scheduling, set reset with `--supports-class-c=false`.

>Note: you can also pass `--supports-class-c` when creating the device. Class C scheduling will be enable after the first uplink message which confirms the device session.

#### Class C multicast

Multicast messages are downlinks messages which are sent to multiple devices that share the same security context. In the Network Server, this is an ABP session. See [creating a device](#createdev) for learning how to create a multicast device.

Multicast sessions do not allow uplink. Therefore, you need to explicitly specify the gateway(s) to send messages from, using the `class_b_c` field:

```json
{
  "downlinks": [{
    "f_port": 15,
    "frm_payload": "vu8=",
    "priority": "NORMAL",
    "class_b_c": {
      "gateways": [{
        "gateway_ids": {
          "gateway_id": "gtw1"
        }
      }]
    }
  }]
}
```

>Note: if you specify multiple gateways, the Network Server will try the gateways in the order specified. The first gateway with no conflicts and no duty-cycle limitation will send the message.

### Listing the downlink queue

The stack keeps a queue of downlink messages. Applications can keep pushing downlink messages or replace the queue with a list of downlink messages.

You can see what is in the queue;

```bash
$ ttn-lw-cli end-devices downlink list app1 dev1
```

## <a name="webhooks">Using webhooks</a>

The webhooks feature allows the Application Server to send application related messages to specific HTTP(S) endpoints.

To show supported formats, use:

```
$ ttn-lw-cli applications webhooks get-formats
```

The `json` formatter uses the same format as the MQTT server described above.

Creating a webhook requires you to have an HTTP(S) endpoint available.

```bash
$ ttn-lw-cli applications webhooks set \
  --application-id app1 \
  --webhook-id wh1 \
  --format json \
  --base-url https://example.com/lorahooks \
  --join-accept.path /join \
  --uplink-message.path /up
```

This will create a webhook `wh1` for the application `app1` with JSON formatting. The paths are appended to the base URL. So, the Application Server will perform `POST` requests on the endpoint `https://example.com/lorahooks/join` for join-accepts and `https://example.com/lorahooks/up` for uplink messages.

>Note: You can also specify URL paths for downlink events, just like MQTT. See `ttn-lw-cli applications webhooks set --help` for more information.

You can also send downlink messages using webhooks. The path is `/v3/api/as/applications/{application_id}/webhooks/{webhook_id}/devices/{device_id}/down/push` (or `/replace`). This requires an API key with traffic writing rights, which can be created as follows:

```bash
$ ttn-lw-cli applications api-keys create \
  --name wh-client \
  --application-id app1 \
  --right-application-traffic-down-write
```

Pass the API key as bearer token on the `Authorization` header. For example:

```
$ curl http://localhost:1885/api/v3/as/applications/app1/webhooks/wh1/devices/dev1/down/push \
  -X POST \
  -H 'Authorization: Bearer NNSXS.VEEBURF3KR77ZR..' \
  --data '{"downlinks":[{"frm_payload":"vu8=","f_port":15,"priority":"NORMAL"}]}'
```

## Congratulations

You have now set up The Things Network Stack V3! 🎉

---

## <a name="events">Advanced: Events</a>

The stack generates lots of events that allow you to get insight in what is going on. You can subscribe to application, gateway, end device events, as well as to user, organization and OAuth client events.

### Using the CLI

To follow your gateway `gtw1` and application `app1` events at the same time:

```bash
$ ttn-lw-cli events subscribe --gateway-id gtw1 --application-id app1
```

### Using cURL

You can also get streaming events with `curl`. For this, you need an API key for the entities you want to watch, for example:

```bash
$ ttn-lw-cli user api-key create \
  --user-id admin \
  --right-application-all \
  --right-gateway-all
```

With the created API key:

```
$ curl http://localhost:1885/api/v3/events \
  -X POST \
  -H 'Authorization: Bearer NNSXS.BR55PTYILPPVXY..' \
  --data '{"identifiers":[{"application_ids":{"application_id":"app1"}},{"gateway_ids":{"gateway_id":"gtw1"}}]}'
```

>Note: The created API key for events is highly privileged; do not use it if you don't need it for events.

### Example: join flow

These are the events of a typical join flow:

```js
{
  "name": "gs.up.receive", // Gateway Server received an uplink message from a device.
  "time": "2019-04-04T09:54:34.786220Z",
  "identifiers": [
    {
      "gateway_ids": {
        "gateway_id": "multitech",
        "eui": "00800000A0000DB4"
      }
    }
  ],
  "correlation_ids": [
    "gs:conn:01D7KWADW2E5CJA32VS1MTR2J6",
    "gs:uplink:01D7KWB0N2KVCV8HZABC8DDHSA"
  ]
}
{
  "name": "js.join.accept", // Join Server accepted the join-accept.
  "time": "2019-04-04T09:54:34.806812Z",
  "identifiers": [
    {
      "device_ids": {
        "device_id": "dev1",
        "application_ids": {
          "application_id": "app1"
        },
        "dev_eui": "4200000000000000",
        "join_eui": "4200000000000000"
      }
    }
  ],
  "correlation_ids": [
    "rpc:/ttn.lorawan.v3.NsJs/HandleJoin:01D7KWB0NCTDY835V5N3CYWBZK"
  ]
}
{
  "name": "ns.up.join.forward", // Network Server forwarded the join-accept and it got accepted.
  "time": "2019-04-04T09:54:34.808132Z",
  "identifiers": [
    {
      "device_ids": {
        "device_id": "dev1",
        "application_ids": {
          "application_id": "app1"
        },
        "dev_eui": "4200000000000000",
        "join_eui": "4200000000000000"
      }
    }
  ],
  "correlation_ids": [
    "gs:conn:01D7KWADW2E5CJA32VS1MTR2J6",
    "gs:uplink:01D7KWB0N2KVCV8HZABC8DDHSA",
    "ns:uplink:01D7KWB0N5C1T8TE2HAVBJN5Y4",
    "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:01D7KWB0N5G2N5C0AFXT4YMF8R"
  ]
}
{
  "name": "ns.up.merge_metadata", // Network Server merged metadata of incoming uplink messages.
  "time": "2019-04-04T09:54:34.991332Z",
  "identifiers": [
    {
      "device_ids": {
        "device_id": "dev1",
        "application_ids": {
          "application_id": "app1"
        },
        "dev_eui": "4200000000000000",
        "join_eui": "4200000000000000"
      }
    }
  ],
  "data": {
    "@type": "type.googleapis.com/google.protobuf.Value",
    "value": 1 // There was 1 gateway that received the join-request.
  },
  "correlation_ids": [
    // Here you find the correlation IDs of all gs.up.receive events that were merged.
    "gs:conn:01D7KWADW2E5CJA32VS1MTR2J6",
    "gs:uplink:01D7KWB0N2KVCV8HZABC8DDHSA",
    "ns:uplink:01D7KWB0N5C1T8TE2HAVBJN5Y4",
    "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:01D7KWB0N5G2N5C0AFXT4YMF8R"
  ]
}
{
  "name": "as.up.join.receive", // Application Server receives the join-accept.
  "time": "2019-04-04T09:54:35.005090Z",
  "identifiers": [
    {
      "device_ids": {
        "device_id": "dev1",
        "application_ids": {
          "application_id": "app1"
        },
        "dev_eui": "4200000000000000",
        "join_eui": "4200000000000000",
        "dev_addr": "0063ECE2"
      }
    }
  ],
  "correlation_ids": [
    "as:up:01D7KWB0VX1D7G3RKFN9HDA39Q",
    "gs:conn:01D7KWADW2E5CJA32VS1MTR2J6",
    "gs:uplink:01D7KWB0N2KVCV8HZABC8DDHSA",
    "ns:uplink:01D7KWB0N5C1T8TE2HAVBJN5Y4",
    "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:01D7KWB0N5G2N5C0AFXT4YMF8R"
  ]
}
{
  "name": "as.up.join.forward", // Application Server forwards the join-accept to an application (CLI, MQTT, webhooks, etc).
  "time": "2019-04-04T09:54:35.010243Z",
  "identifiers": [
    {
      "device_ids": {
        "device_id": "dev1",
        "application_ids": {
          "application_id": "app1"
        },
        "dev_eui": "4200000000000000",
        "join_eui": "4200000000000000",
        "dev_addr": "0063ECE2"
      }
    }
  ],
  "correlation_ids": [
    "as:up:01D7KWB0VX1D7G3RKFN9HDA39Q",
    "gs:conn:01D7KWADW2E5CJA32VS1MTR2J6",
    "gs:uplink:01D7KWB0N2KVCV8HZABC8DDHSA",
    "ns:uplink:01D7KWB0N5C1T8TE2HAVBJN5Y4",
    "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:01D7KWB0N5G2N5C0AFXT4YMF8R"
  ]
}
{
  "name": "gs.down.send", // Gateway Server sent the join-accept to the gateway.
  "time": "2019-04-04T09:54:35.046147Z",
  "identifiers": [
    {
      "gateway_ids": {
        "gateway_id": "multitech",
        "eui": "00800000A0000DB4"
      }
    }
  ],
  "correlation_ids": [
    "gs:conn:01D7KWADW2E5CJA32VS1MTR2J6",
    "rpc:/ttn.lorawan.v3.NsGs/ScheduleDownlink:01D7KWB0W84AJ1P5A3AQV6R4J7"
  ]
}
{
  "name": "gs.up.forward", // Gateway Server forwarded join-request to the Network Server.
  "time": "2019-04-04T09:54:35.991226Z",
  "identifiers": [
    {
      "gateway_ids": {
        "gateway_id": "multitech",
        "eui": "00800000A0000DB4"
      }
    }
  ],
  "correlation_ids": [
    "gs:conn:01D7KWADW2E5CJA32VS1MTR2J6",
    "gs:uplink:01D7KWB0N2KVCV8HZABC8DDHSA"
  ]
}
```
