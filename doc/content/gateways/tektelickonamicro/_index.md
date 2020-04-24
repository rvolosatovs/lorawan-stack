---
title: "Tektelic Kona Micro IoT LoRaWAN Gateway"
description: ""
---

Tektelic Kona Micro IoT LoRaWAN Gateway is an 8 channel LoRaWAN gateway, whose technical specifications can be found in [the official documentation](https://tektelic.com/wp-content/uploads/KONA-Micro.pdf). This page guides you to connect it to {{% tts %}}.

## Prerequisites

1. User account on {{% tts %}} with rights to create Gateways.
2. Tektelic Kona Micro IoT LoRaWAN Gateway connected to the Internet via ethernet.

## Registration

Create a gateway by following the instructions for the [Console]({{< ref "/getting-started/console/create-gateway" >}}) or the [CLI]({{< ref "/getting-started/cli#create-gateway" >}}). The **EUI** of the gateway can be found on the back panel of the gateway under the field **GW ID**.

## Configuration using a Terminal

Find the IP address the gateway. This can be done in various ways. You can connect your machine to the same local network as that of the gateway ethernet connection and scan for open SSH ports or assign a static IP to the gateway and use that. Once the gateway IP address is found, SSH into it.

```bash
$ ssh root@<GatewayIP>
```

The password for the **root** user can be found on the back panel. It's typically a 9 character alphanumeric string starting with **1846XXXXX**.

Now you can edit the gateway configuration file.

```bash
$ vi /etc/default/config.json
```

>Note: Press the `i` key on your keyboard to start insert mode. Once finished editing, press `ESC` and enter `:wq` to write the file and quit.

Edit the server parameters.

1. **server_address**: Address of the Gateway Server. If you followed the [Getting Started guide]({{< ref "/getting-started" >}}) this is the same as what you use instead of `thethings.example.com`.
2. **serv_port_up**: UDP upstream port of the Gateway Server, typically 1700.
3. **serv_port_down**: UDP downstream port of the Gateway Server, typically 1700.

Save the configuration and restart the packet forwarder.

```bash
$ /etc/init.d/pkt_fwd restart
```

If your configuration was successful, your gateway will connect to {{% tts %}} after a couple of seconds.

## Configuration using the GUI (Windows Only)

> TODO: [Document using GUI](https://github.com/TheThingsNetwork/lorawan-stack/issues/1140)

## Troubleshooting

If the gateway does not connect to {{% tts %}} after a few minutes, disconnect and reconnect the power supply to power-cycle the gateway. Packet forwarder logs can be observed by SSH-ing into the gateway and running:

```bash
$ tail -f /var/log/pkt_fwd.log
```
