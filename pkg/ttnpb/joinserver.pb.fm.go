// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"

	go_thethings_network_lorawan_stack_pkg_types "go.thethings.network/lorawan-stack/pkg/types"
)

var SessionKeyRequestFieldPathsNested = []string{
	"dev_eui",
	"session_key_id",
}

var SessionKeyRequestFieldPathsTopLevel = []string{
	"dev_eui",
	"session_key_id",
}

func (dst *SessionKeyRequest) SetFields(src *SessionKeyRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "session_key_id":
			if len(subs) > 0 {
				return fmt.Errorf("'session_key_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.SessionKeyID = src.SessionKeyID
			} else {
				var zero []byte
				dst.SessionKeyID = zero
			}
		case "dev_eui":
			if len(subs) > 0 {
				return fmt.Errorf("'dev_eui' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DevEUI = src.DevEUI
			} else {
				var zero go_thethings_network_lorawan_stack_pkg_types.EUI64
				dst.DevEUI = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var NwkSKeysResponseFieldPathsNested = []string{
	"f_nwk_s_int_key",
	"f_nwk_s_int_key.kek_label",
	"f_nwk_s_int_key.key",
	"nwk_s_enc_key",
	"nwk_s_enc_key.kek_label",
	"nwk_s_enc_key.key",
	"s_nwk_s_int_key",
	"s_nwk_s_int_key.kek_label",
	"s_nwk_s_int_key.key",
}

var NwkSKeysResponseFieldPathsTopLevel = []string{
	"f_nwk_s_int_key",
	"nwk_s_enc_key",
	"s_nwk_s_int_key",
}

func (dst *NwkSKeysResponse) SetFields(src *NwkSKeysResponse, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "f_nwk_s_int_key":
			if len(subs) > 0 {
				newDst := &dst.FNwkSIntKey
				var newSrc *KeyEnvelope
				if src != nil {
					newSrc = &src.FNwkSIntKey
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.FNwkSIntKey = src.FNwkSIntKey
				} else {
					var zero KeyEnvelope
					dst.FNwkSIntKey = zero
				}
			}
		case "s_nwk_s_int_key":
			if len(subs) > 0 {
				newDst := &dst.SNwkSIntKey
				var newSrc *KeyEnvelope
				if src != nil {
					newSrc = &src.SNwkSIntKey
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.SNwkSIntKey = src.SNwkSIntKey
				} else {
					var zero KeyEnvelope
					dst.SNwkSIntKey = zero
				}
			}
		case "nwk_s_enc_key":
			if len(subs) > 0 {
				newDst := &dst.NwkSEncKey
				var newSrc *KeyEnvelope
				if src != nil {
					newSrc = &src.NwkSEncKey
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.NwkSEncKey = src.NwkSEncKey
				} else {
					var zero KeyEnvelope
					dst.NwkSEncKey = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var AppSKeyResponseFieldPathsNested = []string{
	"app_s_key",
	"app_s_key.kek_label",
	"app_s_key.key",
}

var AppSKeyResponseFieldPathsTopLevel = []string{
	"app_s_key",
}

func (dst *AppSKeyResponse) SetFields(src *AppSKeyResponse, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "app_s_key":
			if len(subs) > 0 {
				newDst := &dst.AppSKey
				var newSrc *KeyEnvelope
				if src != nil {
					newSrc = &src.AppSKey
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.AppSKey = src.AppSKey
				} else {
					var zero KeyEnvelope
					dst.AppSKey = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var CryptoServicePayloadRequestFieldPathsNested = []string{
	"ids",
	"ids.application_ids",
	"ids.application_ids.application_id",
	"ids.dev_addr",
	"ids.dev_eui",
	"ids.device_id",
	"ids.join_eui",
	"lorawan_version",
	"payload",
	"provisioner",
	"provisioning_data",
}

var CryptoServicePayloadRequestFieldPathsTopLevel = []string{
	"ids",
	"lorawan_version",
	"payload",
	"provisioner",
	"provisioning_data",
}

func (dst *CryptoServicePayloadRequest) SetFields(src *CryptoServicePayloadRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "ids":
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
		case "lorawan_version":
			if len(subs) > 0 {
				return fmt.Errorf("'lorawan_version' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LoRaWANVersion = src.LoRaWANVersion
			} else {
				var zero MACVersion
				dst.LoRaWANVersion = zero
			}
		case "payload":
			if len(subs) > 0 {
				return fmt.Errorf("'payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Payload = src.Payload
			} else {
				var zero []byte
				dst.Payload = zero
			}
		case "provisioner":
			if len(subs) > 0 {
				return fmt.Errorf("'provisioner' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Provisioner = src.Provisioner
			} else {
				var zero string
				dst.Provisioner = zero
			}
		case "provisioning_data":
			if len(subs) > 0 {
				return fmt.Errorf("'provisioning_data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ProvisioningData = src.ProvisioningData
			} else {
				dst.ProvisioningData = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var CryptoServicePayloadResponseFieldPathsNested = []string{
	"payload",
}

var CryptoServicePayloadResponseFieldPathsTopLevel = []string{
	"payload",
}

func (dst *CryptoServicePayloadResponse) SetFields(src *CryptoServicePayloadResponse, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "payload":
			if len(subs) > 0 {
				return fmt.Errorf("'payload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Payload = src.Payload
			} else {
				var zero []byte
				dst.Payload = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var JoinAcceptMICRequestFieldPathsNested = []string{
	"dev_nonce",
	"join_request_type",
	"payload_request",
	"payload_request.ids",
	"payload_request.ids.application_ids",
	"payload_request.ids.application_ids.application_id",
	"payload_request.ids.dev_addr",
	"payload_request.ids.dev_eui",
	"payload_request.ids.device_id",
	"payload_request.ids.join_eui",
	"payload_request.lorawan_version",
	"payload_request.payload",
	"payload_request.provisioner",
	"payload_request.provisioning_data",
}

var JoinAcceptMICRequestFieldPathsTopLevel = []string{
	"dev_nonce",
	"join_request_type",
	"payload_request",
}

func (dst *JoinAcceptMICRequest) SetFields(src *JoinAcceptMICRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "payload_request":
			if len(subs) > 0 {
				newDst := &dst.CryptoServicePayloadRequest
				var newSrc *CryptoServicePayloadRequest
				if src != nil {
					newSrc = &src.CryptoServicePayloadRequest
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.CryptoServicePayloadRequest = src.CryptoServicePayloadRequest
				} else {
					var zero CryptoServicePayloadRequest
					dst.CryptoServicePayloadRequest = zero
				}
			}
		case "join_request_type":
			if len(subs) > 0 {
				return fmt.Errorf("'join_request_type' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.JoinRequestType = src.JoinRequestType
			} else {
				var zero uint32
				dst.JoinRequestType = zero
			}
		case "dev_nonce":
			if len(subs) > 0 {
				return fmt.Errorf("'dev_nonce' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DevNonce = src.DevNonce
			} else {
				var zero go_thethings_network_lorawan_stack_pkg_types.DevNonce
				dst.DevNonce = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var DeriveSessionKeysRequestFieldPathsNested = []string{
	"dev_nonce",
	"ids",
	"ids.application_ids",
	"ids.application_ids.application_id",
	"ids.dev_addr",
	"ids.dev_eui",
	"ids.device_id",
	"ids.join_eui",
	"join_nonce",
	"lorawan_version",
	"net_id",
	"provisioner",
	"provisioning_data",
}

var DeriveSessionKeysRequestFieldPathsTopLevel = []string{
	"dev_nonce",
	"ids",
	"join_nonce",
	"lorawan_version",
	"net_id",
	"provisioner",
	"provisioning_data",
}

func (dst *DeriveSessionKeysRequest) SetFields(src *DeriveSessionKeysRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "ids":
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
		case "lorawan_version":
			if len(subs) > 0 {
				return fmt.Errorf("'lorawan_version' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LoRaWANVersion = src.LoRaWANVersion
			} else {
				var zero MACVersion
				dst.LoRaWANVersion = zero
			}
		case "join_nonce":
			if len(subs) > 0 {
				return fmt.Errorf("'join_nonce' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.JoinNonce = src.JoinNonce
			} else {
				var zero go_thethings_network_lorawan_stack_pkg_types.JoinNonce
				dst.JoinNonce = zero
			}
		case "dev_nonce":
			if len(subs) > 0 {
				return fmt.Errorf("'dev_nonce' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DevNonce = src.DevNonce
			} else {
				var zero go_thethings_network_lorawan_stack_pkg_types.DevNonce
				dst.DevNonce = zero
			}
		case "net_id":
			if len(subs) > 0 {
				return fmt.Errorf("'net_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.NetID = src.NetID
			} else {
				var zero go_thethings_network_lorawan_stack_pkg_types.NetID
				dst.NetID = zero
			}
		case "provisioner":
			if len(subs) > 0 {
				return fmt.Errorf("'provisioner' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Provisioner = src.Provisioner
			} else {
				var zero string
				dst.Provisioner = zero
			}
		case "provisioning_data":
			if len(subs) > 0 {
				return fmt.Errorf("'provisioning_data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ProvisioningData = src.ProvisioningData
			} else {
				dst.ProvisioningData = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ProvisionEndDevicesRequestFieldPathsNested = []string{
	"application_id",
	"application_id.application_id",
	"data",
	"end_device_ids",
	"provisioner",
}

var ProvisionEndDevicesRequestFieldPathsTopLevel = []string{
	"application_id",
	"data",
	"end_device_ids",
	"provisioner",
}

func (dst *ProvisionEndDevicesRequest) SetFields(src *ProvisionEndDevicesRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "application_id":
			if len(subs) > 0 {
				newDst := &dst.ApplicationIdentifiers
				var newSrc *ApplicationIdentifiers
				if src != nil {
					newSrc = &src.ApplicationIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationIdentifiers = src.ApplicationIdentifiers
				} else {
					var zero ApplicationIdentifiers
					dst.ApplicationIdentifiers = zero
				}
			}
		case "provisioner":
			if len(subs) > 0 {
				return fmt.Errorf("'provisioner' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Provisioner = src.Provisioner
			} else {
				var zero string
				dst.Provisioner = zero
			}
		case "data":
			if len(subs) > 0 {
				return fmt.Errorf("'data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Data = src.Data
			} else {
				var zero []byte
				dst.Data = zero
			}
		case "end_device_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'end_device_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.EndDeviceIDs = src.EndDeviceIDs
			} else {
				dst.EndDeviceIDs = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
