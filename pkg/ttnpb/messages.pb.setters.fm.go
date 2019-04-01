// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"
)

func (dst *UplinkMessage) SetFields(src *UplinkMessage, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "raw_payload":
			if len(subs) > 0 {
				return fmt.Errorf("'raw_payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RawPayload = src.RawPayload
			} else {
				dst.RawPayload = nil
			}
		case "payload":
			if len(subs) > 0 {
				newDst := dst.Payload
				if newDst == nil {
					newDst = &Message{}
					dst.Payload = newDst
				}
				var newSrc *Message
				if src != nil {
					newSrc = src.Payload
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Payload = src.Payload
				} else {
					dst.Payload = nil
				}
			}
		case "settings":
			if len(subs) > 0 {
				newDst := &dst.Settings
				var newSrc *TxSettings
				if src != nil {
					newSrc = &src.Settings
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Settings = src.Settings
				} else {
					var zero TxSettings
					dst.Settings = zero
				}
			}
		case "rx_metadata":
			if len(subs) > 0 {
				return fmt.Errorf("'rx_metadata' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RxMetadata = src.RxMetadata
			} else {
				dst.RxMetadata = nil
			}
		case "received_at":
			if len(subs) > 0 {
				return fmt.Errorf("'received_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ReceivedAt = src.ReceivedAt
			} else {
				var zero time.Time
				dst.ReceivedAt = zero
			}
		case "correlation_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'correlation_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CorrelationIDs = src.CorrelationIDs
			} else {
				dst.CorrelationIDs = nil
			}
		case "gateway_channel_index":
			if len(subs) > 0 {
				return fmt.Errorf("'gateway_channel_index' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.GatewayChannelIndex = src.GatewayChannelIndex
			} else {
				var zero uint32
				dst.GatewayChannelIndex = zero
			}
		case "device_channel_index":
			if len(subs) > 0 {
				return fmt.Errorf("'device_channel_index' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DeviceChannelIndex = src.DeviceChannelIndex
			} else {
				var zero uint32
				dst.DeviceChannelIndex = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *DownlinkMessage) SetFields(src *DownlinkMessage, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "raw_payload":
			if len(subs) > 0 {
				return fmt.Errorf("'raw_payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RawPayload = src.RawPayload
			} else {
				dst.RawPayload = nil
			}
		case "payload":
			if len(subs) > 0 {
				newDst := dst.Payload
				if newDst == nil {
					newDst = &Message{}
					dst.Payload = newDst
				}
				var newSrc *Message
				if src != nil {
					newSrc = src.Payload
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Payload = src.Payload
				} else {
					dst.Payload = nil
				}
			}
		case "end_device_ids":
			if len(subs) > 0 {
				newDst := dst.EndDeviceIDs
				if newDst == nil {
					newDst = &EndDeviceIdentifiers{}
					dst.EndDeviceIDs = newDst
				}
				var newSrc *EndDeviceIdentifiers
				if src != nil {
					newSrc = src.EndDeviceIDs
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.EndDeviceIDs = src.EndDeviceIDs
				} else {
					dst.EndDeviceIDs = nil
				}
			}
		case "correlation_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'correlation_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CorrelationIDs = src.CorrelationIDs
			} else {
				dst.CorrelationIDs = nil
			}

		case "settings":
			if len(subs) == 0 && src == nil {
				dst.Settings = nil
				continue
			} else if len(subs) == 0 {
				dst.Settings = src.Settings
				continue
			}

			subPathMap := _processPaths(subs)
			if len(subPathMap) > 1 {
				return fmt.Errorf("more than one field specified for oneof field '%s'", name)
			}
			for oneofName, oneofSubs := range subPathMap {
				switch oneofName {
				case "request":
					if _, ok := dst.Settings.(*DownlinkMessage_Request); !ok {
						dst.Settings = &DownlinkMessage_Request{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Settings.(*DownlinkMessage_Request).Request
						if newDst == nil {
							newDst = &TxRequest{}
							dst.Settings.(*DownlinkMessage_Request).Request = newDst
						}
						var newSrc *TxRequest
						if src != nil {
							newSrc = src.GetRequest()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Settings.(*DownlinkMessage_Request).Request = src.GetRequest()
						} else {
							dst.Settings.(*DownlinkMessage_Request).Request = nil
						}
					}
				case "scheduled":
					if _, ok := dst.Settings.(*DownlinkMessage_Scheduled); !ok {
						dst.Settings = &DownlinkMessage_Scheduled{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Settings.(*DownlinkMessage_Scheduled).Scheduled
						if newDst == nil {
							newDst = &TxSettings{}
							dst.Settings.(*DownlinkMessage_Scheduled).Scheduled = newDst
						}
						var newSrc *TxSettings
						if src != nil {
							newSrc = src.GetScheduled()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Settings.(*DownlinkMessage_Scheduled).Scheduled = src.GetScheduled()
						} else {
							dst.Settings.(*DownlinkMessage_Scheduled).Scheduled = nil
						}
					}

				default:
					return fmt.Errorf("invalid oneof field: '%s.%s'", name, oneofName)
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *TxAcknowledgment) SetFields(src *TxAcknowledgment, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "correlation_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'correlation_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CorrelationIDs = src.CorrelationIDs
			} else {
				dst.CorrelationIDs = nil
			}
		case "result":
			if len(subs) > 0 {
				return fmt.Errorf("'result' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Result = src.Result
			} else {
				var zero TxAcknowledgment_Result
				dst.Result = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationUplink) SetFields(src *ApplicationUplink, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "session_key_id":
			if len(subs) > 0 {
				return fmt.Errorf("'session_key_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.SessionKeyID = src.SessionKeyID
			} else {
				dst.SessionKeyID = nil
			}
		case "f_port":
			if len(subs) > 0 {
				return fmt.Errorf("'f_port' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FPort = src.FPort
			} else {
				var zero uint32
				dst.FPort = zero
			}
		case "f_cnt":
			if len(subs) > 0 {
				return fmt.Errorf("'f_cnt' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FCnt = src.FCnt
			} else {
				var zero uint32
				dst.FCnt = zero
			}
		case "frm_payload":
			if len(subs) > 0 {
				return fmt.Errorf("'frm_payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FRMPayload = src.FRMPayload
			} else {
				dst.FRMPayload = nil
			}
		case "decoded_payload":
			if len(subs) > 0 {
				return fmt.Errorf("'decoded_payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DecodedPayload = src.DecodedPayload
			} else {
				dst.DecodedPayload = nil
			}
		case "rx_metadata":
			if len(subs) > 0 {
				return fmt.Errorf("'rx_metadata' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RxMetadata = src.RxMetadata
			} else {
				dst.RxMetadata = nil
			}
		case "settings":
			if len(subs) > 0 {
				newDst := &dst.Settings
				var newSrc *TxSettings
				if src != nil {
					newSrc = &src.Settings
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Settings = src.Settings
				} else {
					var zero TxSettings
					dst.Settings = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationLocation) SetFields(src *ApplicationLocation, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "service":
			if len(subs) > 0 {
				return fmt.Errorf("'service' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Service = src.Service
			} else {
				var zero string
				dst.Service = zero
			}
		case "location":
			if len(subs) > 0 {
				newDst := &dst.Location
				var newSrc *Location
				if src != nil {
					newSrc = &src.Location
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Location = src.Location
				} else {
					var zero Location
					dst.Location = zero
				}
			}
		case "attributes":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Attributes = src.Attributes
			} else {
				dst.Attributes = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationJoinAccept) SetFields(src *ApplicationJoinAccept, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "session_key_id":
			if len(subs) > 0 {
				return fmt.Errorf("'session_key_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.SessionKeyID = src.SessionKeyID
			} else {
				dst.SessionKeyID = nil
			}
		case "app_s_key":
			if len(subs) > 0 {
				newDst := dst.AppSKey
				if newDst == nil {
					newDst = &KeyEnvelope{}
					dst.AppSKey = newDst
				}
				var newSrc *KeyEnvelope
				if src != nil {
					newSrc = src.AppSKey
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.AppSKey = src.AppSKey
				} else {
					dst.AppSKey = nil
				}
			}
		case "invalidated_downlinks":
			if len(subs) > 0 {
				return fmt.Errorf("'invalidated_downlinks' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.InvalidatedDownlinks = src.InvalidatedDownlinks
			} else {
				dst.InvalidatedDownlinks = nil
			}
		case "pending_session":
			if len(subs) > 0 {
				return fmt.Errorf("'pending_session' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.PendingSession = src.PendingSession
			} else {
				var zero bool
				dst.PendingSession = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationDownlink) SetFields(src *ApplicationDownlink, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "session_key_id":
			if len(subs) > 0 {
				return fmt.Errorf("'session_key_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.SessionKeyID = src.SessionKeyID
			} else {
				dst.SessionKeyID = nil
			}
		case "f_port":
			if len(subs) > 0 {
				return fmt.Errorf("'f_port' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FPort = src.FPort
			} else {
				var zero uint32
				dst.FPort = zero
			}
		case "f_cnt":
			if len(subs) > 0 {
				return fmt.Errorf("'f_cnt' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FCnt = src.FCnt
			} else {
				var zero uint32
				dst.FCnt = zero
			}
		case "frm_payload":
			if len(subs) > 0 {
				return fmt.Errorf("'frm_payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FRMPayload = src.FRMPayload
			} else {
				dst.FRMPayload = nil
			}
		case "decoded_payload":
			if len(subs) > 0 {
				return fmt.Errorf("'decoded_payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DecodedPayload = src.DecodedPayload
			} else {
				dst.DecodedPayload = nil
			}
		case "confirmed":
			if len(subs) > 0 {
				return fmt.Errorf("'confirmed' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Confirmed = src.Confirmed
			} else {
				var zero bool
				dst.Confirmed = zero
			}
		case "class_b_c":
			if len(subs) > 0 {
				newDst := dst.ClassBC
				if newDst == nil {
					newDst = &ApplicationDownlink_ClassBC{}
					dst.ClassBC = newDst
				}
				var newSrc *ApplicationDownlink_ClassBC
				if src != nil {
					newSrc = src.ClassBC
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ClassBC = src.ClassBC
				} else {
					dst.ClassBC = nil
				}
			}
		case "priority":
			if len(subs) > 0 {
				return fmt.Errorf("'priority' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Priority = src.Priority
			} else {
				var zero TxSchedulePriority
				dst.Priority = zero
			}
		case "correlation_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'correlation_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CorrelationIDs = src.CorrelationIDs
			} else {
				dst.CorrelationIDs = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationDownlinks) SetFields(src *ApplicationDownlinks, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "downlinks":
			if len(subs) > 0 {
				return fmt.Errorf("'downlinks' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Downlinks = src.Downlinks
			} else {
				dst.Downlinks = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationDownlinkFailed) SetFields(src *ApplicationDownlinkFailed, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "downlink":
			if len(subs) > 0 {
				newDst := &dst.ApplicationDownlink
				var newSrc *ApplicationDownlink
				if src != nil {
					newSrc = &src.ApplicationDownlink
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationDownlink = src.ApplicationDownlink
				} else {
					var zero ApplicationDownlink
					dst.ApplicationDownlink = zero
				}
			}
		case "error":
			if len(subs) > 0 {
				newDst := &dst.Error
				var newSrc *ErrorDetails
				if src != nil {
					newSrc = &src.Error
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Error = src.Error
				} else {
					var zero ErrorDetails
					dst.Error = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationInvalidatedDownlinks) SetFields(src *ApplicationInvalidatedDownlinks, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "downlinks":
			if len(subs) > 0 {
				return fmt.Errorf("'downlinks' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Downlinks = src.Downlinks
			} else {
				dst.Downlinks = nil
			}
		case "last_f_cnt_down":
			if len(subs) > 0 {
				return fmt.Errorf("'last_f_cnt_down' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LastFCntDown = src.LastFCntDown
			} else {
				var zero uint32
				dst.LastFCntDown = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationUp) SetFields(src *ApplicationUp, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "end_device_ids":
			if len(subs) > 0 {
				newDst := &dst.EndDeviceIdentifiers
				var newSrc *EndDeviceIdentifiers
				if src != nil {
					newSrc = &src.EndDeviceIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.EndDeviceIdentifiers = src.EndDeviceIdentifiers
				} else {
					var zero EndDeviceIdentifiers
					dst.EndDeviceIdentifiers = zero
				}
			}
		case "correlation_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'correlation_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CorrelationIDs = src.CorrelationIDs
			} else {
				dst.CorrelationIDs = nil
			}

		case "up":
			if len(subs) == 0 && src == nil {
				dst.Up = nil
				continue
			} else if len(subs) == 0 {
				dst.Up = src.Up
				continue
			}

			subPathMap := _processPaths(subs)
			if len(subPathMap) > 1 {
				return fmt.Errorf("more than one field specified for oneof field '%s'", name)
			}
			for oneofName, oneofSubs := range subPathMap {
				switch oneofName {
				case "uplink_message":
					if _, ok := dst.Up.(*ApplicationUp_UplinkMessage); !ok {
						dst.Up = &ApplicationUp_UplinkMessage{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_UplinkMessage).UplinkMessage
						if newDst == nil {
							newDst = &ApplicationUplink{}
							dst.Up.(*ApplicationUp_UplinkMessage).UplinkMessage = newDst
						}
						var newSrc *ApplicationUplink
						if src != nil {
							newSrc = src.GetUplinkMessage()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_UplinkMessage).UplinkMessage = src.GetUplinkMessage()
						} else {
							dst.Up.(*ApplicationUp_UplinkMessage).UplinkMessage = nil
						}
					}
				case "join_accept":
					if _, ok := dst.Up.(*ApplicationUp_JoinAccept); !ok {
						dst.Up = &ApplicationUp_JoinAccept{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_JoinAccept).JoinAccept
						if newDst == nil {
							newDst = &ApplicationJoinAccept{}
							dst.Up.(*ApplicationUp_JoinAccept).JoinAccept = newDst
						}
						var newSrc *ApplicationJoinAccept
						if src != nil {
							newSrc = src.GetJoinAccept()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_JoinAccept).JoinAccept = src.GetJoinAccept()
						} else {
							dst.Up.(*ApplicationUp_JoinAccept).JoinAccept = nil
						}
					}
				case "downlink_ack":
					if _, ok := dst.Up.(*ApplicationUp_DownlinkAck); !ok {
						dst.Up = &ApplicationUp_DownlinkAck{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_DownlinkAck).DownlinkAck
						if newDst == nil {
							newDst = &ApplicationDownlink{}
							dst.Up.(*ApplicationUp_DownlinkAck).DownlinkAck = newDst
						}
						var newSrc *ApplicationDownlink
						if src != nil {
							newSrc = src.GetDownlinkAck()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_DownlinkAck).DownlinkAck = src.GetDownlinkAck()
						} else {
							dst.Up.(*ApplicationUp_DownlinkAck).DownlinkAck = nil
						}
					}
				case "downlink_nack":
					if _, ok := dst.Up.(*ApplicationUp_DownlinkNack); !ok {
						dst.Up = &ApplicationUp_DownlinkNack{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_DownlinkNack).DownlinkNack
						if newDst == nil {
							newDst = &ApplicationDownlink{}
							dst.Up.(*ApplicationUp_DownlinkNack).DownlinkNack = newDst
						}
						var newSrc *ApplicationDownlink
						if src != nil {
							newSrc = src.GetDownlinkNack()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_DownlinkNack).DownlinkNack = src.GetDownlinkNack()
						} else {
							dst.Up.(*ApplicationUp_DownlinkNack).DownlinkNack = nil
						}
					}
				case "downlink_sent":
					if _, ok := dst.Up.(*ApplicationUp_DownlinkSent); !ok {
						dst.Up = &ApplicationUp_DownlinkSent{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_DownlinkSent).DownlinkSent
						if newDst == nil {
							newDst = &ApplicationDownlink{}
							dst.Up.(*ApplicationUp_DownlinkSent).DownlinkSent = newDst
						}
						var newSrc *ApplicationDownlink
						if src != nil {
							newSrc = src.GetDownlinkSent()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_DownlinkSent).DownlinkSent = src.GetDownlinkSent()
						} else {
							dst.Up.(*ApplicationUp_DownlinkSent).DownlinkSent = nil
						}
					}
				case "downlink_failed":
					if _, ok := dst.Up.(*ApplicationUp_DownlinkFailed); !ok {
						dst.Up = &ApplicationUp_DownlinkFailed{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_DownlinkFailed).DownlinkFailed
						if newDst == nil {
							newDst = &ApplicationDownlinkFailed{}
							dst.Up.(*ApplicationUp_DownlinkFailed).DownlinkFailed = newDst
						}
						var newSrc *ApplicationDownlinkFailed
						if src != nil {
							newSrc = src.GetDownlinkFailed()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_DownlinkFailed).DownlinkFailed = src.GetDownlinkFailed()
						} else {
							dst.Up.(*ApplicationUp_DownlinkFailed).DownlinkFailed = nil
						}
					}
				case "downlink_queued":
					if _, ok := dst.Up.(*ApplicationUp_DownlinkQueued); !ok {
						dst.Up = &ApplicationUp_DownlinkQueued{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_DownlinkQueued).DownlinkQueued
						if newDst == nil {
							newDst = &ApplicationDownlink{}
							dst.Up.(*ApplicationUp_DownlinkQueued).DownlinkQueued = newDst
						}
						var newSrc *ApplicationDownlink
						if src != nil {
							newSrc = src.GetDownlinkQueued()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_DownlinkQueued).DownlinkQueued = src.GetDownlinkQueued()
						} else {
							dst.Up.(*ApplicationUp_DownlinkQueued).DownlinkQueued = nil
						}
					}
				case "downlink_queue_invalidated":
					if _, ok := dst.Up.(*ApplicationUp_DownlinkQueueInvalidated); !ok {
						dst.Up = &ApplicationUp_DownlinkQueueInvalidated{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_DownlinkQueueInvalidated).DownlinkQueueInvalidated
						if newDst == nil {
							newDst = &ApplicationInvalidatedDownlinks{}
							dst.Up.(*ApplicationUp_DownlinkQueueInvalidated).DownlinkQueueInvalidated = newDst
						}
						var newSrc *ApplicationInvalidatedDownlinks
						if src != nil {
							newSrc = src.GetDownlinkQueueInvalidated()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_DownlinkQueueInvalidated).DownlinkQueueInvalidated = src.GetDownlinkQueueInvalidated()
						} else {
							dst.Up.(*ApplicationUp_DownlinkQueueInvalidated).DownlinkQueueInvalidated = nil
						}
					}
				case "location_solved":
					if _, ok := dst.Up.(*ApplicationUp_LocationSolved); !ok {
						dst.Up = &ApplicationUp_LocationSolved{}
					}
					if len(oneofSubs) > 0 {
						newDst := dst.Up.(*ApplicationUp_LocationSolved).LocationSolved
						if newDst == nil {
							newDst = &ApplicationLocation{}
							dst.Up.(*ApplicationUp_LocationSolved).LocationSolved = newDst
						}
						var newSrc *ApplicationLocation
						if src != nil {
							newSrc = src.GetLocationSolved()
						}
						if err := newDst.SetFields(newSrc, subs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.Up.(*ApplicationUp_LocationSolved).LocationSolved = src.GetLocationSolved()
						} else {
							dst.Up.(*ApplicationUp_LocationSolved).LocationSolved = nil
						}
					}

				default:
					return fmt.Errorf("invalid oneof field: '%s.%s'", name, oneofName)
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *MessagePayloadFormatters) SetFields(src *MessagePayloadFormatters, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "up_formatter":
			if len(subs) > 0 {
				return fmt.Errorf("'up_formatter' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpFormatter = src.UpFormatter
			} else {
				var zero PayloadFormatter
				dst.UpFormatter = zero
			}
		case "up_formatter_parameter":
			if len(subs) > 0 {
				return fmt.Errorf("'up_formatter_parameter' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpFormatterParameter = src.UpFormatterParameter
			} else {
				var zero string
				dst.UpFormatterParameter = zero
			}
		case "down_formatter":
			if len(subs) > 0 {
				return fmt.Errorf("'down_formatter' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DownFormatter = src.DownFormatter
			} else {
				var zero PayloadFormatter
				dst.DownFormatter = zero
			}
		case "down_formatter_parameter":
			if len(subs) > 0 {
				return fmt.Errorf("'down_formatter_parameter' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DownFormatterParameter = src.DownFormatterParameter
			} else {
				var zero string
				dst.DownFormatterParameter = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *DownlinkQueueRequest) SetFields(src *DownlinkQueueRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "end_device_ids":
			if len(subs) > 0 {
				newDst := &dst.EndDeviceIdentifiers
				var newSrc *EndDeviceIdentifiers
				if src != nil {
					newSrc = &src.EndDeviceIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.EndDeviceIdentifiers = src.EndDeviceIdentifiers
				} else {
					var zero EndDeviceIdentifiers
					dst.EndDeviceIdentifiers = zero
				}
			}
		case "downlinks":
			if len(subs) > 0 {
				return fmt.Errorf("'downlinks' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Downlinks = src.Downlinks
			} else {
				dst.Downlinks = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ApplicationDownlink_ClassBC) SetFields(src *ApplicationDownlink_ClassBC, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "gateways":
			if len(subs) > 0 {
				return fmt.Errorf("'gateways' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Gateways = src.Gateways
			} else {
				dst.Gateways = nil
			}
		case "absolute_time":
			if len(subs) > 0 {
				return fmt.Errorf("'absolute_time' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AbsoluteTime = src.AbsoluteTime
			} else {
				dst.AbsoluteTime = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
