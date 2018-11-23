// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
)

var PullGatewayConfigurationRequestFieldPathsNested = []string{
	"field_mask",
	"gateway_ids",
	"gateway_ids.eui",
	"gateway_ids.gateway_id",
}

var PullGatewayConfigurationRequestFieldPathsTopLevel = []string{
	"field_mask",
	"gateway_ids",
}

func (dst *PullGatewayConfigurationRequest) SetFields(src *PullGatewayConfigurationRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "gateway_ids":
			if len(subs) > 0 {
				newDst := &dst.GatewayIdentifiers
				var newSrc *GatewayIdentifiers
				if src != nil {
					newSrc = &src.GatewayIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.GatewayIdentifiers = src.GatewayIdentifiers
				} else {
					var zero GatewayIdentifiers
					dst.GatewayIdentifiers = zero
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				var zero github_com_gogo_protobuf_types.FieldMask
				dst.FieldMask = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
