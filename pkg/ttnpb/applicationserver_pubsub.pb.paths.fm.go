// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

var ApplicationPubSubIdentifiersFieldPathsNested = []string{
	"application_ids",
	"application_ids.application_id",
	"pub_sub_id",
}

var ApplicationPubSubIdentifiersFieldPathsTopLevel = []string{
	"application_ids",
	"pub_sub_id",
}
var ApplicationPubSubFieldPathsNested = []string{
	"base_topic",
	"created_at",
	"downlink_ack",
	"downlink_ack.topic",
	"downlink_failed",
	"downlink_failed.topic",
	"downlink_nack",
	"downlink_nack.topic",
	"downlink_push",
	"downlink_push.topic",
	"downlink_queued",
	"downlink_queued.topic",
	"downlink_replace",
	"downlink_replace.topic",
	"downlink_sent",
	"downlink_sent.topic",
	"format",
	"ids",
	"ids.application_ids",
	"ids.application_ids.application_id",
	"ids.pub_sub_id",
	"join_accept",
	"join_accept.topic",
	"location_solved",
	"location_solved.topic",
	"provider",
	"provider.nats",
	"provider.nats.server_url",
	"updated_at",
	"uplink_message",
	"uplink_message.topic",
}

var ApplicationPubSubFieldPathsTopLevel = []string{
	"base_topic",
	"created_at",
	"downlink_ack",
	"downlink_failed",
	"downlink_nack",
	"downlink_push",
	"downlink_queued",
	"downlink_replace",
	"downlink_sent",
	"format",
	"ids",
	"join_accept",
	"location_solved",
	"provider",
	"updated_at",
	"uplink_message",
}
var ApplicationPubSubsFieldPathsNested = []string{
	"pubsubs",
}

var ApplicationPubSubsFieldPathsTopLevel = []string{
	"pubsubs",
}
var ApplicationPubSubFormatsFieldPathsNested = []string{
	"formats",
}

var ApplicationPubSubFormatsFieldPathsTopLevel = []string{
	"formats",
}
var GetApplicationPubSubRequestFieldPathsNested = []string{
	"field_mask",
	"ids",
	"ids.application_ids",
	"ids.application_ids.application_id",
	"ids.pub_sub_id",
}

var GetApplicationPubSubRequestFieldPathsTopLevel = []string{
	"field_mask",
	"ids",
}
var ListApplicationPubSubsRequestFieldPathsNested = []string{
	"application_ids",
	"application_ids.application_id",
	"field_mask",
}

var ListApplicationPubSubsRequestFieldPathsTopLevel = []string{
	"application_ids",
	"field_mask",
}
var SetApplicationPubSubRequestFieldPathsNested = []string{
	"field_mask",
	"pubsub",
	"pubsub.base_topic",
	"pubsub.created_at",
	"pubsub.downlink_ack",
	"pubsub.downlink_ack.topic",
	"pubsub.downlink_failed",
	"pubsub.downlink_failed.topic",
	"pubsub.downlink_nack",
	"pubsub.downlink_nack.topic",
	"pubsub.downlink_push",
	"pubsub.downlink_push.topic",
	"pubsub.downlink_queued",
	"pubsub.downlink_queued.topic",
	"pubsub.downlink_replace",
	"pubsub.downlink_replace.topic",
	"pubsub.downlink_sent",
	"pubsub.downlink_sent.topic",
	"pubsub.format",
	"pubsub.ids",
	"pubsub.ids.application_ids",
	"pubsub.ids.application_ids.application_id",
	"pubsub.ids.pub_sub_id",
	"pubsub.join_accept",
	"pubsub.join_accept.topic",
	"pubsub.location_solved",
	"pubsub.location_solved.topic",
	"pubsub.provider",
	"pubsub.provider.nats",
	"pubsub.provider.nats.server_url",
	"pubsub.updated_at",
	"pubsub.uplink_message",
	"pubsub.uplink_message.topic",
}

var SetApplicationPubSubRequestFieldPathsTopLevel = []string{
	"field_mask",
	"pubsub",
}
var ApplicationPubSub_NATSProviderFieldPathsNested = []string{
	"server_url",
}

var ApplicationPubSub_NATSProviderFieldPathsTopLevel = []string{
	"server_url",
}
var ApplicationPubSub_MessageFieldPathsNested = []string{
	"topic",
}

var ApplicationPubSub_MessageFieldPathsTopLevel = []string{
	"topic",
}
