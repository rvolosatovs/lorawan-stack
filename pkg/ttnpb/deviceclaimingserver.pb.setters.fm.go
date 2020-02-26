// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"

	go_thethings_network_lorawan_stack_pkg_types "go.thethings.network/lorawan-stack/pkg/types"
)

func (dst *ClaimEndDeviceRequest) SetFields(src *ClaimEndDeviceRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "target_application_ids":
			if len(subs) > 0 {
				var newDst, newSrc *ApplicationIdentifiers
				if src != nil {
					newSrc = &src.TargetApplicationIDs
				}
				newDst = &dst.TargetApplicationIDs
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.TargetApplicationIDs = src.TargetApplicationIDs
				} else {
					var zero ApplicationIdentifiers
					dst.TargetApplicationIDs = zero
				}
			}
		case "target_device_id":
			if len(subs) > 0 {
				return fmt.Errorf("'target_device_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TargetDeviceID = src.TargetDeviceID
			} else {
				var zero string
				dst.TargetDeviceID = zero
			}
		case "target_network_server_address":
			if len(subs) > 0 {
				return fmt.Errorf("'target_network_server_address' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TargetNetworkServerAddress = src.TargetNetworkServerAddress
			} else {
				var zero string
				dst.TargetNetworkServerAddress = zero
			}
		case "target_network_server_kek_label":
			if len(subs) > 0 {
				return fmt.Errorf("'target_network_server_kek_label' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TargetNetworkServerKEKLabel = src.TargetNetworkServerKEKLabel
			} else {
				var zero string
				dst.TargetNetworkServerKEKLabel = zero
			}
		case "target_application_server_address":
			if len(subs) > 0 {
				return fmt.Errorf("'target_application_server_address' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TargetApplicationServerAddress = src.TargetApplicationServerAddress
			} else {
				var zero string
				dst.TargetApplicationServerAddress = zero
			}
		case "target_application_server_kek_label":
			if len(subs) > 0 {
				return fmt.Errorf("'target_application_server_kek_label' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TargetApplicationServerKEKLabel = src.TargetApplicationServerKEKLabel
			} else {
				var zero string
				dst.TargetApplicationServerKEKLabel = zero
			}
		case "target_application_server_id":
			if len(subs) > 0 {
				return fmt.Errorf("'target_application_server_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TargetApplicationServerID = src.TargetApplicationServerID
			} else {
				var zero string
				dst.TargetApplicationServerID = zero
			}
		case "target_net_id":
			if len(subs) > 0 {
				return fmt.Errorf("'target_net_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TargetNetID = src.TargetNetID
			} else {
				dst.TargetNetID = nil
			}
		case "invalidate_authentication_code":
			if len(subs) > 0 {
				return fmt.Errorf("'invalidate_authentication_code' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.InvalidateAuthenticationCode = src.InvalidateAuthenticationCode
			} else {
				var zero bool
				dst.InvalidateAuthenticationCode = zero
			}

		case "source_device":
			if len(subs) == 0 && src == nil {
				dst.SourceDevice = nil
				continue
			} else if len(subs) == 0 {
				dst.SourceDevice = src.SourceDevice
				continue
			}

			subPathMap := _processPaths(subs)
			if len(subPathMap) > 1 {
				return fmt.Errorf("more than one field specified for oneof field '%s'", name)
			}
			for oneofName, oneofSubs := range subPathMap {
				switch oneofName {
				case "authenticated_identifiers":
					_, srcOk := src.SourceDevice.(*ClaimEndDeviceRequest_AuthenticatedIdentifiers_)
					if !srcOk && src.SourceDevice != nil {
						return fmt.Errorf("attempt to set oneof 'authenticated_identifiers', while different oneof is set in source")
					}
					_, dstOk := dst.SourceDevice.(*ClaimEndDeviceRequest_AuthenticatedIdentifiers_)
					if !dstOk && dst.SourceDevice != nil {
						return fmt.Errorf("attempt to set oneof 'authenticated_identifiers', while different oneof is set in destination")
					}
					if len(oneofSubs) > 0 {
						var newDst, newSrc *ClaimEndDeviceRequest_AuthenticatedIdentifiers
						if !srcOk && !dstOk {
							continue
						}
						if srcOk {
							newSrc = src.SourceDevice.(*ClaimEndDeviceRequest_AuthenticatedIdentifiers_).AuthenticatedIdentifiers
						}
						if dstOk {
							newDst = dst.SourceDevice.(*ClaimEndDeviceRequest_AuthenticatedIdentifiers_).AuthenticatedIdentifiers
						} else {
							newDst = &ClaimEndDeviceRequest_AuthenticatedIdentifiers{}
							dst.SourceDevice = &ClaimEndDeviceRequest_AuthenticatedIdentifiers_{AuthenticatedIdentifiers: newDst}
						}
						if err := newDst.SetFields(newSrc, oneofSubs...); err != nil {
							return err
						}
					} else {
						if src != nil {
							dst.SourceDevice = src.SourceDevice
						} else {
							dst.SourceDevice = nil
						}
					}
				case "qr_code":
					_, srcOk := src.SourceDevice.(*ClaimEndDeviceRequest_QRCode)
					if !srcOk && src.SourceDevice != nil {
						return fmt.Errorf("attempt to set oneof 'qr_code', while different oneof is set in source")
					}
					_, dstOk := dst.SourceDevice.(*ClaimEndDeviceRequest_QRCode)
					if !dstOk && dst.SourceDevice != nil {
						return fmt.Errorf("attempt to set oneof 'qr_code', while different oneof is set in destination")
					}
					if len(oneofSubs) > 0 {
						return fmt.Errorf("'qr_code' has no subfields, but %s were specified", oneofSubs)
					}
					if src != nil {
						dst.SourceDevice = src.SourceDevice
					} else {
						dst.SourceDevice = nil
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

func (dst *AuthorizeApplicationRequest) SetFields(src *AuthorizeApplicationRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "application_ids":
			if len(subs) > 0 {
				var newDst, newSrc *ApplicationIdentifiers
				if src != nil {
					newSrc = &src.ApplicationIdentifiers
				}
				newDst = &dst.ApplicationIdentifiers
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
		case "api_key":
			if len(subs) > 0 {
				return fmt.Errorf("'api_key' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.APIKey = src.APIKey
			} else {
				var zero string
				dst.APIKey = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ClaimEndDeviceRequest_AuthenticatedIdentifiers) SetFields(src *ClaimEndDeviceRequest_AuthenticatedIdentifiers, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "join_eui":
			if len(subs) > 0 {
				return fmt.Errorf("'join_eui' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.JoinEUI = src.JoinEUI
			} else {
				var zero go_thethings_network_lorawan_stack_pkg_types.EUI64
				dst.JoinEUI = zero
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
		case "authentication_code":
			if len(subs) > 0 {
				return fmt.Errorf("'authentication_code' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AuthenticationCode = src.AuthenticationCode
			} else {
				var zero string
				dst.AuthenticationCode = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
