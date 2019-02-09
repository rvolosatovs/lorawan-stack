// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"
)

var GatewayUpFieldPathsNested = []string{
	"gateway_status",
	"gateway_status.advanced",
	"gateway_status.antenna_locations",
	"gateway_status.boot_time",
	"gateway_status.ip",
	"gateway_status.metrics",
	"gateway_status.time",
	"gateway_status.versions",
	"tx_acknowledgment",
	"tx_acknowledgment.correlation_ids",
	"tx_acknowledgment.result",
	"uplink_messages",
}

var GatewayUpFieldPathsTopLevel = []string{
	"gateway_status",
	"tx_acknowledgment",
	"uplink_messages",
}

func (dst *GatewayUp) SetFields(src *GatewayUp, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "uplink_messages":
			if len(subs) > 0 {
				return fmt.Errorf("'uplink_messages' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UplinkMessages = src.UplinkMessages
			} else {
				dst.UplinkMessages = nil
			}
		case "gateway_status":
			if len(subs) > 0 {
				newDst := dst.GatewayStatus
				if newDst == nil {
					newDst = &GatewayStatus{}
					dst.GatewayStatus = newDst
				}
				var newSrc *GatewayStatus
				if src != nil {
					newSrc = src.GatewayStatus
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.GatewayStatus = src.GatewayStatus
				} else {
					dst.GatewayStatus = nil
				}
			}
		case "tx_acknowledgment":
			if len(subs) > 0 {
				newDst := dst.TxAcknowledgment
				if newDst == nil {
					newDst = &TxAcknowledgment{}
					dst.TxAcknowledgment = newDst
				}
				var newSrc *TxAcknowledgment
				if src != nil {
					newSrc = src.TxAcknowledgment
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.TxAcknowledgment = src.TxAcknowledgment
				} else {
					dst.TxAcknowledgment = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayDownFieldPathsNested = []string{
	"downlink_message",
	"downlink_message.correlation_ids",
	"downlink_message.end_device_ids",
	"downlink_message.end_device_ids.application_ids",
	"downlink_message.end_device_ids.application_ids.application_id",
	"downlink_message.end_device_ids.dev_addr",
	"downlink_message.end_device_ids.dev_eui",
	"downlink_message.end_device_ids.device_id",
	"downlink_message.end_device_ids.join_eui",
	"downlink_message.payload",
	"downlink_message.payload.Payload",
	"downlink_message.payload.Payload.join_accept_payload",
	"downlink_message.payload.Payload.join_accept_payload.cf_list",
	"downlink_message.payload.Payload.join_accept_payload.cf_list.ch_masks",
	"downlink_message.payload.Payload.join_accept_payload.cf_list.freq",
	"downlink_message.payload.Payload.join_accept_payload.cf_list.type",
	"downlink_message.payload.Payload.join_accept_payload.dev_addr",
	"downlink_message.payload.Payload.join_accept_payload.dl_settings",
	"downlink_message.payload.Payload.join_accept_payload.dl_settings.opt_neg",
	"downlink_message.payload.Payload.join_accept_payload.dl_settings.rx1_dr_offset",
	"downlink_message.payload.Payload.join_accept_payload.dl_settings.rx2_dr",
	"downlink_message.payload.Payload.join_accept_payload.encrypted",
	"downlink_message.payload.Payload.join_accept_payload.join_nonce",
	"downlink_message.payload.Payload.join_accept_payload.net_id",
	"downlink_message.payload.Payload.join_accept_payload.rx_delay",
	"downlink_message.payload.Payload.join_request_payload",
	"downlink_message.payload.Payload.join_request_payload.dev_eui",
	"downlink_message.payload.Payload.join_request_payload.dev_nonce",
	"downlink_message.payload.Payload.join_request_payload.join_eui",
	"downlink_message.payload.Payload.mac_payload",
	"downlink_message.payload.Payload.mac_payload.decoded_payload",
	"downlink_message.payload.Payload.mac_payload.f_hdr",
	"downlink_message.payload.Payload.mac_payload.f_hdr.dev_addr",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_cnt",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_ctrl",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_ctrl.ack",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_ctrl.adr",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_ctrl.adr_ack_req",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_ctrl.class_b",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_ctrl.f_pending",
	"downlink_message.payload.Payload.mac_payload.f_hdr.f_opts",
	"downlink_message.payload.Payload.mac_payload.f_port",
	"downlink_message.payload.Payload.mac_payload.frm_payload",
	"downlink_message.payload.Payload.rejoin_request_payload",
	"downlink_message.payload.Payload.rejoin_request_payload.dev_eui",
	"downlink_message.payload.Payload.rejoin_request_payload.join_eui",
	"downlink_message.payload.Payload.rejoin_request_payload.net_id",
	"downlink_message.payload.Payload.rejoin_request_payload.rejoin_cnt",
	"downlink_message.payload.Payload.rejoin_request_payload.rejoin_type",
	"downlink_message.payload.m_hdr",
	"downlink_message.payload.m_hdr.m_type",
	"downlink_message.payload.m_hdr.major",
	"downlink_message.payload.mic",
	"downlink_message.raw_payload",
	"downlink_message.settings",
	"downlink_message.settings.request",
	"downlink_message.settings.request.absolute_time",
	"downlink_message.settings.request.advanced",
	"downlink_message.settings.request.class",
	"downlink_message.settings.request.downlink_paths",
	"downlink_message.settings.request.priority",
	"downlink_message.settings.request.rx1_data_rate_index",
	"downlink_message.settings.request.rx1_delay",
	"downlink_message.settings.request.rx1_frequency",
	"downlink_message.settings.request.rx2_data_rate_index",
	"downlink_message.settings.request.rx2_frequency",
	"downlink_message.settings.scheduled",
	"downlink_message.settings.scheduled.coding_rate",
	"downlink_message.settings.scheduled.data_rate",
	"downlink_message.settings.scheduled.data_rate.modulation",
	"downlink_message.settings.scheduled.data_rate.modulation.fsk",
	"downlink_message.settings.scheduled.data_rate.modulation.fsk.bit_rate",
	"downlink_message.settings.scheduled.data_rate.modulation.lora",
	"downlink_message.settings.scheduled.data_rate.modulation.lora.bandwidth",
	"downlink_message.settings.scheduled.data_rate.modulation.lora.spreading_factor",
	"downlink_message.settings.scheduled.data_rate_index",
	"downlink_message.settings.scheduled.device_channel_index",
	"downlink_message.settings.scheduled.enable_crc",
	"downlink_message.settings.scheduled.frequency",
	"downlink_message.settings.scheduled.gateway_channel_index",
	"downlink_message.settings.scheduled.invert_polarization",
	"downlink_message.settings.scheduled.time",
	"downlink_message.settings.scheduled.timestamp",
	"downlink_message.settings.scheduled.tx_power",
}

var GatewayDownFieldPathsTopLevel = []string{
	"downlink_message",
}

func (dst *GatewayDown) SetFields(src *GatewayDown, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "downlink_message":
			if len(subs) > 0 {
				newDst := dst.DownlinkMessage
				if newDst == nil {
					newDst = &DownlinkMessage{}
					dst.DownlinkMessage = newDst
				}
				var newSrc *DownlinkMessage
				if src != nil {
					newSrc = src.DownlinkMessage
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkMessage = src.DownlinkMessage
				} else {
					dst.DownlinkMessage = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ScheduleDownlinkResponseFieldPathsNested = []string{
	"delay",
}

var ScheduleDownlinkResponseFieldPathsTopLevel = []string{
	"delay",
}

func (dst *ScheduleDownlinkResponse) SetFields(src *ScheduleDownlinkResponse, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "delay":
			if len(subs) > 0 {
				return fmt.Errorf("'delay' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Delay = src.Delay
			} else {
				var zero time.Duration
				dst.Delay = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
