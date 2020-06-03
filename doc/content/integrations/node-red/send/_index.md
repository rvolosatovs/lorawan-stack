---
title: "Send Messages"
description: ""
weight: 3
---

This section explains the process of setting up a flow which publishes messages to a certain topic that the MQTT Server is subscribed to.

Doing this schedules downlink messages to be sent to your end device. This section follows the example for publishing downlink traffic in [MQTT Server]({{< ref "/integrations/mqtt" >}}) guide.

## Configure MQTT Out Node

1. Place the **mqtt out** node on the dashboard. 

2. Configure the **Server** options with the same settings as in the [Receive Events and Messages]({{< ref "/integrations/node-red#receive-events-and-messages" >}}) section.

3. Set **Topic** to `v3/{application_id}/devices/{device_id}/down/push` to schedule downlink messages (as stated in [MQTT Server]({{< ref "/integrations/mqtt" >}}) guide). 

4. Choose a **QoS** from listed options and state whether you want the MQTT Server to retain messages. 

{{< figure src="mqtt_out_node_properties.png" alt="mqtt out node properties" >}}

## Configure Inject Node

1. Place the **inject** node on the dashboard. Double-click on the node to configure its properties. 

2. Choose **buffer** under **Payload** and enter the payload you wish to send. 

3. Define the period between the automatic injections if you want them, or choose **none** for **Repeat** if you wish to inject messages manually.

{{< figure src="inject_node_properties.png" alt="inject node properties" >}}

## Configure Function Node and Deploy

Next, you have to configure a **function** node, which converts previously defined payload to a downlink message with Base64 encoded payload.

1. Place the **function** node with the following structure on dashboard:

```bash
return {
  "payload": {
    "downlinks": [{
      "f_port": 15,
      "frm_payload": msg.payload.toString('base64'),
      "priority": "NORMAL"
    }]
  }
}
```

2. Connect the nodes and click **Deploy**. If the setup is correct, below the **mqtt out** node **connected** status will be reported and downlink messages will begin sending to your end device.

{{< figure src="send_downlink_flow.png" alt="send downlink flow" >}}
