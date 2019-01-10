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
