// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

var UplinkMessageFieldPathsNested = []string{
	"correlation_ids",
	"device_channel_index",
	"payload",
	"payload.Payload",
	"payload.Payload.join_accept_payload",
	"payload.Payload.join_accept_payload.cf_list",
	"payload.Payload.join_accept_payload.cf_list.ch_masks",
	"payload.Payload.join_accept_payload.cf_list.freq",
	"payload.Payload.join_accept_payload.cf_list.type",
	"payload.Payload.join_accept_payload.dev_addr",
	"payload.Payload.join_accept_payload.dl_settings",
	"payload.Payload.join_accept_payload.dl_settings.opt_neg",
	"payload.Payload.join_accept_payload.dl_settings.rx1_dr_offset",
	"payload.Payload.join_accept_payload.dl_settings.rx2_dr",
	"payload.Payload.join_accept_payload.encrypted",
	"payload.Payload.join_accept_payload.join_nonce",
	"payload.Payload.join_accept_payload.net_id",
	"payload.Payload.join_accept_payload.rx_delay",
	"payload.Payload.join_request_payload",
	"payload.Payload.join_request_payload.dev_eui",
	"payload.Payload.join_request_payload.dev_nonce",
	"payload.Payload.join_request_payload.join_eui",
	"payload.Payload.mac_payload",
	"payload.Payload.mac_payload.decoded_payload",
	"payload.Payload.mac_payload.f_hdr",
	"payload.Payload.mac_payload.f_hdr.dev_addr",
	"payload.Payload.mac_payload.f_hdr.f_cnt",
	"payload.Payload.mac_payload.f_hdr.f_ctrl",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.ack",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.adr",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.adr_ack_req",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.class_b",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.f_pending",
	"payload.Payload.mac_payload.f_hdr.f_opts",
	"payload.Payload.mac_payload.f_port",
	"payload.Payload.mac_payload.frm_payload",
	"payload.Payload.rejoin_request_payload",
	"payload.Payload.rejoin_request_payload.dev_eui",
	"payload.Payload.rejoin_request_payload.join_eui",
	"payload.Payload.rejoin_request_payload.net_id",
	"payload.Payload.rejoin_request_payload.rejoin_cnt",
	"payload.Payload.rejoin_request_payload.rejoin_type",
	"payload.m_hdr",
	"payload.m_hdr.m_type",
	"payload.m_hdr.major",
	"payload.mic",
	"raw_payload",
	"received_at",
	"rx_metadata",
	"settings",
	"settings.coding_rate",
	"settings.data_rate",
	"settings.data_rate.modulation",
	"settings.data_rate.modulation.fsk",
	"settings.data_rate.modulation.fsk.bit_rate",
	"settings.data_rate.modulation.lora",
	"settings.data_rate.modulation.lora.bandwidth",
	"settings.data_rate.modulation.lora.spreading_factor",
	"settings.data_rate_index",
	"settings.downlink",
	"settings.downlink.antenna_index",
	"settings.downlink.invert_polarization",
	"settings.downlink.tx_power",
	"settings.enable_crc",
	"settings.frequency",
	"settings.time",
	"settings.timestamp",
}

var UplinkMessageFieldPathsTopLevel = []string{
	"correlation_ids",
	"device_channel_index",
	"payload",
	"raw_payload",
	"received_at",
	"rx_metadata",
	"settings",
}
var DownlinkMessageFieldPathsNested = []string{
	"correlation_ids",
	"end_device_ids",
	"end_device_ids.application_ids",
	"end_device_ids.application_ids.application_id",
	"end_device_ids.dev_addr",
	"end_device_ids.dev_eui",
	"end_device_ids.device_id",
	"end_device_ids.join_eui",
	"payload",
	"payload.Payload",
	"payload.Payload.join_accept_payload",
	"payload.Payload.join_accept_payload.cf_list",
	"payload.Payload.join_accept_payload.cf_list.ch_masks",
	"payload.Payload.join_accept_payload.cf_list.freq",
	"payload.Payload.join_accept_payload.cf_list.type",
	"payload.Payload.join_accept_payload.dev_addr",
	"payload.Payload.join_accept_payload.dl_settings",
	"payload.Payload.join_accept_payload.dl_settings.opt_neg",
	"payload.Payload.join_accept_payload.dl_settings.rx1_dr_offset",
	"payload.Payload.join_accept_payload.dl_settings.rx2_dr",
	"payload.Payload.join_accept_payload.encrypted",
	"payload.Payload.join_accept_payload.join_nonce",
	"payload.Payload.join_accept_payload.net_id",
	"payload.Payload.join_accept_payload.rx_delay",
	"payload.Payload.join_request_payload",
	"payload.Payload.join_request_payload.dev_eui",
	"payload.Payload.join_request_payload.dev_nonce",
	"payload.Payload.join_request_payload.join_eui",
	"payload.Payload.mac_payload",
	"payload.Payload.mac_payload.decoded_payload",
	"payload.Payload.mac_payload.f_hdr",
	"payload.Payload.mac_payload.f_hdr.dev_addr",
	"payload.Payload.mac_payload.f_hdr.f_cnt",
	"payload.Payload.mac_payload.f_hdr.f_ctrl",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.ack",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.adr",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.adr_ack_req",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.class_b",
	"payload.Payload.mac_payload.f_hdr.f_ctrl.f_pending",
	"payload.Payload.mac_payload.f_hdr.f_opts",
	"payload.Payload.mac_payload.f_port",
	"payload.Payload.mac_payload.frm_payload",
	"payload.Payload.rejoin_request_payload",
	"payload.Payload.rejoin_request_payload.dev_eui",
	"payload.Payload.rejoin_request_payload.join_eui",
	"payload.Payload.rejoin_request_payload.net_id",
	"payload.Payload.rejoin_request_payload.rejoin_cnt",
	"payload.Payload.rejoin_request_payload.rejoin_type",
	"payload.m_hdr",
	"payload.m_hdr.m_type",
	"payload.m_hdr.major",
	"payload.mic",
	"raw_payload",
	"settings",
	"settings.request",
	"settings.request.absolute_time",
	"settings.request.advanced",
	"settings.request.class",
	"settings.request.downlink_paths",
	"settings.request.frequency_plan_id",
	"settings.request.priority",
	"settings.request.rx1_data_rate_index",
	"settings.request.rx1_delay",
	"settings.request.rx1_frequency",
	"settings.request.rx2_data_rate_index",
	"settings.request.rx2_frequency",
	"settings.scheduled",
	"settings.scheduled.coding_rate",
	"settings.scheduled.data_rate",
	"settings.scheduled.data_rate.modulation",
	"settings.scheduled.data_rate.modulation.fsk",
	"settings.scheduled.data_rate.modulation.fsk.bit_rate",
	"settings.scheduled.data_rate.modulation.lora",
	"settings.scheduled.data_rate.modulation.lora.bandwidth",
	"settings.scheduled.data_rate.modulation.lora.spreading_factor",
	"settings.scheduled.data_rate_index",
	"settings.scheduled.downlink",
	"settings.scheduled.downlink.antenna_index",
	"settings.scheduled.downlink.invert_polarization",
	"settings.scheduled.downlink.tx_power",
	"settings.scheduled.enable_crc",
	"settings.scheduled.frequency",
	"settings.scheduled.time",
	"settings.scheduled.timestamp",
}

var DownlinkMessageFieldPathsTopLevel = []string{
	"correlation_ids",
	"end_device_ids",
	"payload",
	"raw_payload",
	"settings",
}
var TxAcknowledgmentFieldPathsNested = []string{
	"correlation_ids",
	"result",
}

var TxAcknowledgmentFieldPathsTopLevel = []string{
	"correlation_ids",
	"result",
}
var GatewayUplinkMessageFieldPathsNested = []string{
	"band_id",
	"message",
	"message.correlation_ids",
	"message.device_channel_index",
	"message.payload",
	"message.payload.Payload",
	"message.payload.Payload.join_accept_payload",
	"message.payload.Payload.join_accept_payload.cf_list",
	"message.payload.Payload.join_accept_payload.cf_list.ch_masks",
	"message.payload.Payload.join_accept_payload.cf_list.freq",
	"message.payload.Payload.join_accept_payload.cf_list.type",
	"message.payload.Payload.join_accept_payload.dev_addr",
	"message.payload.Payload.join_accept_payload.dl_settings",
	"message.payload.Payload.join_accept_payload.dl_settings.opt_neg",
	"message.payload.Payload.join_accept_payload.dl_settings.rx1_dr_offset",
	"message.payload.Payload.join_accept_payload.dl_settings.rx2_dr",
	"message.payload.Payload.join_accept_payload.encrypted",
	"message.payload.Payload.join_accept_payload.join_nonce",
	"message.payload.Payload.join_accept_payload.net_id",
	"message.payload.Payload.join_accept_payload.rx_delay",
	"message.payload.Payload.join_request_payload",
	"message.payload.Payload.join_request_payload.dev_eui",
	"message.payload.Payload.join_request_payload.dev_nonce",
	"message.payload.Payload.join_request_payload.join_eui",
	"message.payload.Payload.mac_payload",
	"message.payload.Payload.mac_payload.decoded_payload",
	"message.payload.Payload.mac_payload.f_hdr",
	"message.payload.Payload.mac_payload.f_hdr.dev_addr",
	"message.payload.Payload.mac_payload.f_hdr.f_cnt",
	"message.payload.Payload.mac_payload.f_hdr.f_ctrl",
	"message.payload.Payload.mac_payload.f_hdr.f_ctrl.ack",
	"message.payload.Payload.mac_payload.f_hdr.f_ctrl.adr",
	"message.payload.Payload.mac_payload.f_hdr.f_ctrl.adr_ack_req",
	"message.payload.Payload.mac_payload.f_hdr.f_ctrl.class_b",
	"message.payload.Payload.mac_payload.f_hdr.f_ctrl.f_pending",
	"message.payload.Payload.mac_payload.f_hdr.f_opts",
	"message.payload.Payload.mac_payload.f_port",
	"message.payload.Payload.mac_payload.frm_payload",
	"message.payload.Payload.rejoin_request_payload",
	"message.payload.Payload.rejoin_request_payload.dev_eui",
	"message.payload.Payload.rejoin_request_payload.join_eui",
	"message.payload.Payload.rejoin_request_payload.net_id",
	"message.payload.Payload.rejoin_request_payload.rejoin_cnt",
	"message.payload.Payload.rejoin_request_payload.rejoin_type",
	"message.payload.m_hdr",
	"message.payload.m_hdr.m_type",
	"message.payload.m_hdr.major",
	"message.payload.mic",
	"message.raw_payload",
	"message.received_at",
	"message.rx_metadata",
	"message.settings",
	"message.settings.coding_rate",
	"message.settings.data_rate",
	"message.settings.data_rate.modulation",
	"message.settings.data_rate.modulation.fsk",
	"message.settings.data_rate.modulation.fsk.bit_rate",
	"message.settings.data_rate.modulation.lora",
	"message.settings.data_rate.modulation.lora.bandwidth",
	"message.settings.data_rate.modulation.lora.spreading_factor",
	"message.settings.data_rate_index",
	"message.settings.downlink",
	"message.settings.downlink.antenna_index",
	"message.settings.downlink.invert_polarization",
	"message.settings.downlink.tx_power",
	"message.settings.enable_crc",
	"message.settings.frequency",
	"message.settings.time",
	"message.settings.timestamp",
}

var GatewayUplinkMessageFieldPathsTopLevel = []string{
	"band_id",
	"message",
}
var ApplicationUplinkFieldPathsNested = []string{
	"app_s_key",
	"app_s_key.encrypted_key",
	"app_s_key.kek_label",
	"app_s_key.key",
	"decoded_payload",
	"f_cnt",
	"f_port",
	"frm_payload",
	"last_a_f_cnt_down",
	"received_at",
	"rx_metadata",
	"session_key_id",
	"settings",
	"settings.coding_rate",
	"settings.data_rate",
	"settings.data_rate.modulation",
	"settings.data_rate.modulation.fsk",
	"settings.data_rate.modulation.fsk.bit_rate",
	"settings.data_rate.modulation.lora",
	"settings.data_rate.modulation.lora.bandwidth",
	"settings.data_rate.modulation.lora.spreading_factor",
	"settings.data_rate_index",
	"settings.downlink",
	"settings.downlink.antenna_index",
	"settings.downlink.invert_polarization",
	"settings.downlink.tx_power",
	"settings.enable_crc",
	"settings.frequency",
	"settings.time",
	"settings.timestamp",
}

var ApplicationUplinkFieldPathsTopLevel = []string{
	"app_s_key",
	"decoded_payload",
	"f_cnt",
	"f_port",
	"frm_payload",
	"last_a_f_cnt_down",
	"received_at",
	"rx_metadata",
	"session_key_id",
	"settings",
}
var ApplicationLocationFieldPathsNested = []string{
	"attributes",
	"location",
	"location.accuracy",
	"location.altitude",
	"location.latitude",
	"location.longitude",
	"location.source",
	"service",
}

var ApplicationLocationFieldPathsTopLevel = []string{
	"attributes",
	"location",
	"service",
}
var ApplicationJoinAcceptFieldPathsNested = []string{
	"app_s_key",
	"app_s_key.encrypted_key",
	"app_s_key.kek_label",
	"app_s_key.key",
	"invalidated_downlinks",
	"pending_session",
	"received_at",
	"session_key_id",
}

var ApplicationJoinAcceptFieldPathsTopLevel = []string{
	"app_s_key",
	"invalidated_downlinks",
	"pending_session",
	"received_at",
	"session_key_id",
}
var ApplicationDownlinkFieldPathsNested = []string{
	"class_b_c",
	"class_b_c.absolute_time",
	"class_b_c.gateways",
	"confirmed",
	"correlation_ids",
	"decoded_payload",
	"f_cnt",
	"f_port",
	"frm_payload",
	"priority",
	"session_key_id",
}

var ApplicationDownlinkFieldPathsTopLevel = []string{
	"class_b_c",
	"confirmed",
	"correlation_ids",
	"decoded_payload",
	"f_cnt",
	"f_port",
	"frm_payload",
	"priority",
	"session_key_id",
}
var ApplicationDownlinksFieldPathsNested = []string{
	"downlinks",
}

var ApplicationDownlinksFieldPathsTopLevel = []string{
	"downlinks",
}
var ApplicationDownlinkFailedFieldPathsNested = []string{
	"downlink",
	"downlink.class_b_c",
	"downlink.class_b_c.absolute_time",
	"downlink.class_b_c.gateways",
	"downlink.confirmed",
	"downlink.correlation_ids",
	"downlink.decoded_payload",
	"downlink.f_cnt",
	"downlink.f_port",
	"downlink.frm_payload",
	"downlink.priority",
	"downlink.session_key_id",
	"error",
	"error.attributes",
	"error.cause",
	"error.cause.attributes",
	"error.cause.correlation_id",
	"error.cause.message_format",
	"error.cause.name",
	"error.cause.namespace",
	"error.code",
	"error.correlation_id",
	"error.details",
	"error.message_format",
	"error.name",
	"error.namespace",
}

var ApplicationDownlinkFailedFieldPathsTopLevel = []string{
	"downlink",
	"error",
}
var ApplicationInvalidatedDownlinksFieldPathsNested = []string{
	"downlinks",
	"last_f_cnt_down",
}

var ApplicationInvalidatedDownlinksFieldPathsTopLevel = []string{
	"downlinks",
	"last_f_cnt_down",
}
var ApplicationUpFieldPathsNested = []string{
	"correlation_ids",
	"end_device_ids",
	"end_device_ids.application_ids",
	"end_device_ids.application_ids.application_id",
	"end_device_ids.dev_addr",
	"end_device_ids.dev_eui",
	"end_device_ids.device_id",
	"end_device_ids.join_eui",
	"received_at",
	"up",
	"up.downlink_ack",
	"up.downlink_ack.class_b_c",
	"up.downlink_ack.class_b_c.absolute_time",
	"up.downlink_ack.class_b_c.gateways",
	"up.downlink_ack.confirmed",
	"up.downlink_ack.correlation_ids",
	"up.downlink_ack.decoded_payload",
	"up.downlink_ack.f_cnt",
	"up.downlink_ack.f_port",
	"up.downlink_ack.frm_payload",
	"up.downlink_ack.priority",
	"up.downlink_ack.session_key_id",
	"up.downlink_failed",
	"up.downlink_failed.downlink",
	"up.downlink_failed.downlink.class_b_c",
	"up.downlink_failed.downlink.class_b_c.absolute_time",
	"up.downlink_failed.downlink.class_b_c.gateways",
	"up.downlink_failed.downlink.confirmed",
	"up.downlink_failed.downlink.correlation_ids",
	"up.downlink_failed.downlink.decoded_payload",
	"up.downlink_failed.downlink.f_cnt",
	"up.downlink_failed.downlink.f_port",
	"up.downlink_failed.downlink.frm_payload",
	"up.downlink_failed.downlink.priority",
	"up.downlink_failed.downlink.session_key_id",
	"up.downlink_failed.error",
	"up.downlink_failed.error.attributes",
	"up.downlink_failed.error.cause",
	"up.downlink_failed.error.cause.attributes",
	"up.downlink_failed.error.cause.correlation_id",
	"up.downlink_failed.error.cause.message_format",
	"up.downlink_failed.error.cause.name",
	"up.downlink_failed.error.cause.namespace",
	"up.downlink_failed.error.code",
	"up.downlink_failed.error.correlation_id",
	"up.downlink_failed.error.details",
	"up.downlink_failed.error.message_format",
	"up.downlink_failed.error.name",
	"up.downlink_failed.error.namespace",
	"up.downlink_nack",
	"up.downlink_nack.class_b_c",
	"up.downlink_nack.class_b_c.absolute_time",
	"up.downlink_nack.class_b_c.gateways",
	"up.downlink_nack.confirmed",
	"up.downlink_nack.correlation_ids",
	"up.downlink_nack.decoded_payload",
	"up.downlink_nack.f_cnt",
	"up.downlink_nack.f_port",
	"up.downlink_nack.frm_payload",
	"up.downlink_nack.priority",
	"up.downlink_nack.session_key_id",
	"up.downlink_queue_invalidated",
	"up.downlink_queue_invalidated.downlinks",
	"up.downlink_queue_invalidated.last_f_cnt_down",
	"up.downlink_queued",
	"up.downlink_queued.class_b_c",
	"up.downlink_queued.class_b_c.absolute_time",
	"up.downlink_queued.class_b_c.gateways",
	"up.downlink_queued.confirmed",
	"up.downlink_queued.correlation_ids",
	"up.downlink_queued.decoded_payload",
	"up.downlink_queued.f_cnt",
	"up.downlink_queued.f_port",
	"up.downlink_queued.frm_payload",
	"up.downlink_queued.priority",
	"up.downlink_queued.session_key_id",
	"up.downlink_sent",
	"up.downlink_sent.class_b_c",
	"up.downlink_sent.class_b_c.absolute_time",
	"up.downlink_sent.class_b_c.gateways",
	"up.downlink_sent.confirmed",
	"up.downlink_sent.correlation_ids",
	"up.downlink_sent.decoded_payload",
	"up.downlink_sent.f_cnt",
	"up.downlink_sent.f_port",
	"up.downlink_sent.frm_payload",
	"up.downlink_sent.priority",
	"up.downlink_sent.session_key_id",
	"up.join_accept",
	"up.join_accept.app_s_key",
	"up.join_accept.app_s_key.encrypted_key",
	"up.join_accept.app_s_key.kek_label",
	"up.join_accept.app_s_key.key",
	"up.join_accept.invalidated_downlinks",
	"up.join_accept.pending_session",
	"up.join_accept.received_at",
	"up.join_accept.session_key_id",
	"up.location_solved",
	"up.location_solved.attributes",
	"up.location_solved.location",
	"up.location_solved.location.accuracy",
	"up.location_solved.location.altitude",
	"up.location_solved.location.latitude",
	"up.location_solved.location.longitude",
	"up.location_solved.location.source",
	"up.location_solved.service",
	"up.uplink_message",
	"up.uplink_message.app_s_key",
	"up.uplink_message.app_s_key.encrypted_key",
	"up.uplink_message.app_s_key.kek_label",
	"up.uplink_message.app_s_key.key",
	"up.uplink_message.decoded_payload",
	"up.uplink_message.f_cnt",
	"up.uplink_message.f_port",
	"up.uplink_message.frm_payload",
	"up.uplink_message.last_a_f_cnt_down",
	"up.uplink_message.received_at",
	"up.uplink_message.rx_metadata",
	"up.uplink_message.session_key_id",
	"up.uplink_message.settings",
	"up.uplink_message.settings.coding_rate",
	"up.uplink_message.settings.data_rate",
	"up.uplink_message.settings.data_rate.modulation",
	"up.uplink_message.settings.data_rate.modulation.fsk",
	"up.uplink_message.settings.data_rate.modulation.fsk.bit_rate",
	"up.uplink_message.settings.data_rate.modulation.lora",
	"up.uplink_message.settings.data_rate.modulation.lora.bandwidth",
	"up.uplink_message.settings.data_rate.modulation.lora.spreading_factor",
	"up.uplink_message.settings.data_rate_index",
	"up.uplink_message.settings.downlink",
	"up.uplink_message.settings.downlink.antenna_index",
	"up.uplink_message.settings.downlink.invert_polarization",
	"up.uplink_message.settings.downlink.tx_power",
	"up.uplink_message.settings.enable_crc",
	"up.uplink_message.settings.frequency",
	"up.uplink_message.settings.time",
	"up.uplink_message.settings.timestamp",
}

var ApplicationUpFieldPathsTopLevel = []string{
	"correlation_ids",
	"end_device_ids",
	"received_at",
	"up",
}
var MessagePayloadFormattersFieldPathsNested = []string{
	"down_formatter",
	"down_formatter_parameter",
	"up_formatter",
	"up_formatter_parameter",
}

var MessagePayloadFormattersFieldPathsTopLevel = []string{
	"down_formatter",
	"down_formatter_parameter",
	"up_formatter",
	"up_formatter_parameter",
}
var DownlinkQueueRequestFieldPathsNested = []string{
	"downlinks",
	"end_device_ids",
	"end_device_ids.application_ids",
	"end_device_ids.application_ids.application_id",
	"end_device_ids.dev_addr",
	"end_device_ids.dev_eui",
	"end_device_ids.device_id",
	"end_device_ids.join_eui",
}

var DownlinkQueueRequestFieldPathsTopLevel = []string{
	"downlinks",
	"end_device_ids",
}
var ApplicationDownlink_ClassBCFieldPathsNested = []string{
	"absolute_time",
	"gateways",
}

var ApplicationDownlink_ClassBCFieldPathsTopLevel = []string{
	"absolute_time",
	"gateways",
}
