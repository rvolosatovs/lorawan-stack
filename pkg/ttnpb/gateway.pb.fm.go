// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"

	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
)

var GatewayBrandFieldPathsNested = []string{
	"id",
	"logos",
	"name",
	"url",
}

var GatewayBrandFieldPathsTopLevel = []string{
	"id",
	"logos",
	"name",
	"url",
}

func (dst *GatewayBrand) SetFields(src *GatewayBrand, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "id":
			if len(subs) > 0 {
				return fmt.Errorf("'id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ID = src.ID
			} else {
				var zero string
				dst.ID = zero
			}
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "url":
			if len(subs) > 0 {
				return fmt.Errorf("'url' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.URL = src.URL
			} else {
				var zero string
				dst.URL = zero
			}
		case "logos":
			if len(subs) > 0 {
				return fmt.Errorf("'logos' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Logos = src.Logos
			} else {
				dst.Logos = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayModelFieldPathsNested = []string{
	"brand_id",
	"id",
	"name",
}

var GatewayModelFieldPathsTopLevel = []string{
	"brand_id",
	"id",
	"name",
}

func (dst *GatewayModel) SetFields(src *GatewayModel, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "brand_id":
			if len(subs) > 0 {
				return fmt.Errorf("'brand_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.BrandID = src.BrandID
			} else {
				var zero string
				dst.BrandID = zero
			}
		case "id":
			if len(subs) > 0 {
				return fmt.Errorf("'id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ID = src.ID
			} else {
				var zero string
				dst.ID = zero
			}
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayVersionIdentifiersFieldPathsNested = []string{
	"brand_id",
	"firmware_version",
	"hardware_version",
	"model_id",
}

var GatewayVersionIdentifiersFieldPathsTopLevel = []string{
	"brand_id",
	"firmware_version",
	"hardware_version",
	"model_id",
}

func (dst *GatewayVersionIdentifiers) SetFields(src *GatewayVersionIdentifiers, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "brand_id":
			if len(subs) > 0 {
				return fmt.Errorf("'brand_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.BrandID = src.BrandID
			} else {
				var zero string
				dst.BrandID = zero
			}
		case "model_id":
			if len(subs) > 0 {
				return fmt.Errorf("'model_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ModelID = src.ModelID
			} else {
				var zero string
				dst.ModelID = zero
			}
		case "hardware_version":
			if len(subs) > 0 {
				return fmt.Errorf("'hardware_version' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.HardwareVersion = src.HardwareVersion
			} else {
				var zero string
				dst.HardwareVersion = zero
			}
		case "firmware_version":
			if len(subs) > 0 {
				return fmt.Errorf("'firmware_version' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FirmwareVersion = src.FirmwareVersion
			} else {
				var zero string
				dst.FirmwareVersion = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayRadioFieldPathsNested = []string{
	"chip_type",
	"enable",
	"frequency",
	"rssi_offset",
	"tx_configuration",
	"tx_configuration.max_frequency",
	"tx_configuration.min_frequency",
	"tx_configuration.notch_frequency",
}

var GatewayRadioFieldPathsTopLevel = []string{
	"chip_type",
	"enable",
	"frequency",
	"rssi_offset",
	"tx_configuration",
}

func (dst *GatewayRadio) SetFields(src *GatewayRadio, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "enable":
			if len(subs) > 0 {
				return fmt.Errorf("'enable' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Enable = src.Enable
			} else {
				var zero bool
				dst.Enable = zero
			}
		case "chip_type":
			if len(subs) > 0 {
				return fmt.Errorf("'chip_type' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ChipType = src.ChipType
			} else {
				var zero string
				dst.ChipType = zero
			}
		case "frequency":
			if len(subs) > 0 {
				return fmt.Errorf("'frequency' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Frequency = src.Frequency
			} else {
				var zero uint64
				dst.Frequency = zero
			}
		case "rssi_offset":
			if len(subs) > 0 {
				return fmt.Errorf("'rssi_offset' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RSSIOffset = src.RSSIOffset
			} else {
				var zero float32
				dst.RSSIOffset = zero
			}
		case "tx_configuration":
			if len(subs) > 0 {
				newDst := dst.TxConfiguration
				if newDst == nil {
					newDst = &GatewayRadio_TxConfiguration{}
					dst.TxConfiguration = newDst
				}
				var newSrc *GatewayRadio_TxConfiguration
				if src != nil {
					newSrc = src.TxConfiguration
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.TxConfiguration = src.TxConfiguration
				} else {
					dst.TxConfiguration = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayRadio_TxConfigurationFieldPathsNested = []string{
	"max_frequency",
	"min_frequency",
	"notch_frequency",
}

var GatewayRadio_TxConfigurationFieldPathsTopLevel = []string{
	"max_frequency",
	"min_frequency",
	"notch_frequency",
}

func (dst *GatewayRadio_TxConfiguration) SetFields(src *GatewayRadio_TxConfiguration, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "min_frequency":
			if len(subs) > 0 {
				return fmt.Errorf("'min_frequency' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MinFrequency = src.MinFrequency
			} else {
				var zero uint64
				dst.MinFrequency = zero
			}
		case "max_frequency":
			if len(subs) > 0 {
				return fmt.Errorf("'max_frequency' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MaxFrequency = src.MaxFrequency
			} else {
				var zero uint64
				dst.MaxFrequency = zero
			}
		case "notch_frequency":
			if len(subs) > 0 {
				return fmt.Errorf("'notch_frequency' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.NotchFrequency = src.NotchFrequency
			} else {
				var zero uint64
				dst.NotchFrequency = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayVersionFieldPathsNested = []string{
	"clock_source",
	"ids",
	"ids.brand_id",
	"ids.firmware_version",
	"ids.hardware_version",
	"ids.model_id",
	"photos",
	"radios",
}

var GatewayVersionFieldPathsTopLevel = []string{
	"clock_source",
	"ids",
	"photos",
	"radios",
}

func (dst *GatewayVersion) SetFields(src *GatewayVersion, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "ids":
			if len(subs) > 0 {
				newDst := &dst.GatewayVersionIdentifiers
				var newSrc *GatewayVersionIdentifiers
				if src != nil {
					newSrc = &src.GatewayVersionIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.GatewayVersionIdentifiers = src.GatewayVersionIdentifiers
				} else {
					var zero GatewayVersionIdentifiers
					dst.GatewayVersionIdentifiers = zero
				}
			}
		case "photos":
			if len(subs) > 0 {
				return fmt.Errorf("'photos' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Photos = src.Photos
			} else {
				dst.Photos = nil
			}
		case "radios":
			if len(subs) > 0 {
				return fmt.Errorf("'radios' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Radios = src.Radios
			} else {
				dst.Radios = nil
			}
		case "clock_source":
			if len(subs) > 0 {
				return fmt.Errorf("'clock_source' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ClockSource = src.ClockSource
			} else {
				var zero uint32
				dst.ClockSource = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayFieldPathsNested = []string{
	"antennas",
	"attributes",
	"auto_update",
	"contact_info",
	"created_at",
	"description",
	"downlink_path_constraint",
	"enforce_duty_cycle",
	"frequency_plan_id",
	"gateway_server_address",
	"ids",
	"ids.eui",
	"ids.gateway_id",
	"location_public",
	"name",
	"schedule_downlink_late",
	"status_public",
	"update_channel",
	"updated_at",
	"version_ids",
	"version_ids.brand_id",
	"version_ids.firmware_version",
	"version_ids.hardware_version",
	"version_ids.model_id",
}

var GatewayFieldPathsTopLevel = []string{
	"antennas",
	"attributes",
	"auto_update",
	"contact_info",
	"created_at",
	"description",
	"downlink_path_constraint",
	"enforce_duty_cycle",
	"frequency_plan_id",
	"gateway_server_address",
	"ids",
	"location_public",
	"name",
	"schedule_downlink_late",
	"status_public",
	"update_channel",
	"updated_at",
	"version_ids",
}

func (dst *Gateway) SetFields(src *Gateway, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "ids":
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
		case "created_at":
			if len(subs) > 0 {
				return fmt.Errorf("'created_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CreatedAt = src.CreatedAt
			} else {
				var zero time.Time
				dst.CreatedAt = zero
			}
		case "updated_at":
			if len(subs) > 0 {
				return fmt.Errorf("'updated_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdatedAt = src.UpdatedAt
			} else {
				var zero time.Time
				dst.UpdatedAt = zero
			}
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "description":
			if len(subs) > 0 {
				return fmt.Errorf("'description' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Description = src.Description
			} else {
				var zero string
				dst.Description = zero
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
		case "contact_info":
			if len(subs) > 0 {
				return fmt.Errorf("'contact_info' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ContactInfo = src.ContactInfo
			} else {
				dst.ContactInfo = nil
			}
		case "version_ids":
			if len(subs) > 0 {
				newDst := &dst.GatewayVersionIdentifiers
				var newSrc *GatewayVersionIdentifiers
				if src != nil {
					newSrc = &src.GatewayVersionIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.GatewayVersionIdentifiers = src.GatewayVersionIdentifiers
				} else {
					var zero GatewayVersionIdentifiers
					dst.GatewayVersionIdentifiers = zero
				}
			}
		case "gateway_server_address":
			if len(subs) > 0 {
				return fmt.Errorf("'gateway_server_address' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.GatewayServerAddress = src.GatewayServerAddress
			} else {
				var zero string
				dst.GatewayServerAddress = zero
			}
		case "auto_update":
			if len(subs) > 0 {
				return fmt.Errorf("'auto_update' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AutoUpdate = src.AutoUpdate
			} else {
				var zero bool
				dst.AutoUpdate = zero
			}
		case "update_channel":
			if len(subs) > 0 {
				return fmt.Errorf("'update_channel' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdateChannel = src.UpdateChannel
			} else {
				var zero string
				dst.UpdateChannel = zero
			}
		case "frequency_plan_id":
			if len(subs) > 0 {
				return fmt.Errorf("'frequency_plan_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FrequencyPlanID = src.FrequencyPlanID
			} else {
				var zero string
				dst.FrequencyPlanID = zero
			}
		case "antennas":
			if len(subs) > 0 {
				return fmt.Errorf("'antennas' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Antennas = src.Antennas
			} else {
				dst.Antennas = nil
			}
		case "status_public":
			if len(subs) > 0 {
				return fmt.Errorf("'status_public' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.StatusPublic = src.StatusPublic
			} else {
				var zero bool
				dst.StatusPublic = zero
			}
		case "location_public":
			if len(subs) > 0 {
				return fmt.Errorf("'location_public' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LocationPublic = src.LocationPublic
			} else {
				var zero bool
				dst.LocationPublic = zero
			}
		case "schedule_downlink_late":
			if len(subs) > 0 {
				return fmt.Errorf("'schedule_downlink_late' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ScheduleDownlinkLate = src.ScheduleDownlinkLate
			} else {
				var zero bool
				dst.ScheduleDownlinkLate = zero
			}
		case "enforce_duty_cycle":
			if len(subs) > 0 {
				return fmt.Errorf("'enforce_duty_cycle' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.EnforceDutyCycle = src.EnforceDutyCycle
			} else {
				var zero bool
				dst.EnforceDutyCycle = zero
			}
		case "downlink_path_constraint":
			if len(subs) > 0 {
				return fmt.Errorf("'downlink_path_constraint' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DownlinkPathConstraint = src.DownlinkPathConstraint
			} else {
				var zero DownlinkPathConstraint
				dst.DownlinkPathConstraint = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewaysFieldPathsNested = []string{
	"gateways",
}

var GatewaysFieldPathsTopLevel = []string{
	"gateways",
}

func (dst *Gateways) SetFields(src *Gateways, paths ...string) error {
	for name, subs := range _processPaths(paths) {
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

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GetGatewayRequestFieldPathsNested = []string{
	"field_mask",
	"gateway_ids",
	"gateway_ids.eui",
	"gateway_ids.gateway_id",
}

var GetGatewayRequestFieldPathsTopLevel = []string{
	"field_mask",
	"gateway_ids",
}

func (dst *GetGatewayRequest) SetFields(src *GetGatewayRequest, paths ...string) error {
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

var ListGatewaysRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
	"field_mask",
	"ids",
	"limit",
	"order",
	"page",
}

var ListGatewaysRequestFieldPathsTopLevel = []string{
	"collaborator",
	"field_mask",
	"ids",
	"limit",
	"order",
	"page",
}

func (dst *ListGatewaysRequest) SetFields(src *ListGatewaysRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "collaborator":
			if len(subs) > 0 {
				newDst := dst.Collaborator
				if newDst == nil {
					newDst = &OrganizationOrUserIdentifiers{}
					dst.Collaborator = newDst
				}
				var newSrc *OrganizationOrUserIdentifiers
				if src != nil {
					newSrc = src.Collaborator
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Collaborator = src.Collaborator
				} else {
					dst.Collaborator = nil
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
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
			}
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var CreateGatewayRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
	"gateway",
	"gateway.antennas",
	"gateway.attributes",
	"gateway.auto_update",
	"gateway.contact_info",
	"gateway.created_at",
	"gateway.description",
	"gateway.downlink_path_constraint",
	"gateway.enforce_duty_cycle",
	"gateway.frequency_plan_id",
	"gateway.gateway_server_address",
	"gateway.ids",
	"gateway.ids.eui",
	"gateway.ids.gateway_id",
	"gateway.location_public",
	"gateway.name",
	"gateway.schedule_downlink_late",
	"gateway.status_public",
	"gateway.update_channel",
	"gateway.updated_at",
	"gateway.version_ids",
	"gateway.version_ids.brand_id",
	"gateway.version_ids.firmware_version",
	"gateway.version_ids.hardware_version",
	"gateway.version_ids.model_id",
	"ids",
}

var CreateGatewayRequestFieldPathsTopLevel = []string{
	"collaborator",
	"gateway",
	"ids",
}

func (dst *CreateGatewayRequest) SetFields(src *CreateGatewayRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "gateway":
			if len(subs) > 0 {
				newDst := &dst.Gateway
				var newSrc *Gateway
				if src != nil {
					newSrc = &src.Gateway
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Gateway = src.Gateway
				} else {
					var zero Gateway
					dst.Gateway = zero
				}
			}
		case "collaborator":
			if len(subs) > 0 {
				newDst := &dst.Collaborator
				var newSrc *OrganizationOrUserIdentifiers
				if src != nil {
					newSrc = &src.Collaborator
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Collaborator = src.Collaborator
				} else {
					var zero OrganizationOrUserIdentifiers
					dst.Collaborator = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var UpdateGatewayRequestFieldPathsNested = []string{
	"field_mask",
	"gateway",
	"gateway.antennas",
	"gateway.attributes",
	"gateway.auto_update",
	"gateway.contact_info",
	"gateway.created_at",
	"gateway.description",
	"gateway.downlink_path_constraint",
	"gateway.enforce_duty_cycle",
	"gateway.frequency_plan_id",
	"gateway.gateway_server_address",
	"gateway.ids",
	"gateway.ids.eui",
	"gateway.ids.gateway_id",
	"gateway.location_public",
	"gateway.name",
	"gateway.schedule_downlink_late",
	"gateway.status_public",
	"gateway.update_channel",
	"gateway.updated_at",
	"gateway.version_ids",
	"gateway.version_ids.brand_id",
	"gateway.version_ids.firmware_version",
	"gateway.version_ids.hardware_version",
	"gateway.version_ids.model_id",
}

var UpdateGatewayRequestFieldPathsTopLevel = []string{
	"field_mask",
	"gateway",
}

func (dst *UpdateGatewayRequest) SetFields(src *UpdateGatewayRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "gateway":
			if len(subs) > 0 {
				newDst := &dst.Gateway
				var newSrc *Gateway
				if src != nil {
					newSrc = &src.Gateway
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Gateway = src.Gateway
				} else {
					var zero Gateway
					dst.Gateway = zero
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

var CreateGatewayAPIKeyRequestFieldPathsNested = []string{
	"gateway_ids",
	"gateway_ids.eui",
	"gateway_ids.gateway_id",
	"name",
	"rights",
}

var CreateGatewayAPIKeyRequestFieldPathsTopLevel = []string{
	"gateway_ids",
	"name",
	"rights",
}

func (dst *CreateGatewayAPIKeyRequest) SetFields(src *CreateGatewayAPIKeyRequest, paths ...string) error {
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
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "rights":
			if len(subs) > 0 {
				return fmt.Errorf("'rights' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Rights = src.Rights
			} else {
				dst.Rights = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var UpdateGatewayAPIKeyRequestFieldPathsNested = []string{
	"api_key",
	"api_key.id",
	"api_key.key",
	"api_key.name",
	"api_key.rights",
	"gateway_ids",
	"gateway_ids.eui",
	"gateway_ids.gateway_id",
}

var UpdateGatewayAPIKeyRequestFieldPathsTopLevel = []string{
	"api_key",
	"gateway_ids",
}

func (dst *UpdateGatewayAPIKeyRequest) SetFields(src *UpdateGatewayAPIKeyRequest, paths ...string) error {
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
		case "api_key":
			if len(subs) > 0 {
				newDst := &dst.APIKey
				var newSrc *APIKey
				if src != nil {
					newSrc = &src.APIKey
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.APIKey = src.APIKey
				} else {
					var zero APIKey
					dst.APIKey = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var SetGatewayCollaboratorRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.ids.organization_ids",
	"collaborator.ids.ids.organization_ids.organization_id",
	"collaborator.ids.ids.user_ids",
	"collaborator.ids.ids.user_ids.email",
	"collaborator.ids.ids.user_ids.user_id",
	"collaborator.rights",
	"gateway_ids",
	"gateway_ids.eui",
	"gateway_ids.gateway_id",
	"ids",
}

var SetGatewayCollaboratorRequestFieldPathsTopLevel = []string{
	"collaborator",
	"gateway_ids",
	"ids",
}

func (dst *SetGatewayCollaboratorRequest) SetFields(src *SetGatewayCollaboratorRequest, paths ...string) error {
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
		case "collaborator":
			if len(subs) > 0 {
				newDst := &dst.Collaborator
				var newSrc *Collaborator
				if src != nil {
					newSrc = &src.Collaborator
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Collaborator = src.Collaborator
				} else {
					var zero Collaborator
					dst.Collaborator = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayAntennaFieldPathsNested = []string{
	"attributes",
	"gain",
	"location",
	"location.accuracy",
	"location.altitude",
	"location.latitude",
	"location.longitude",
	"location.source",
}

var GatewayAntennaFieldPathsTopLevel = []string{
	"attributes",
	"gain",
	"location",
}

func (dst *GatewayAntenna) SetFields(src *GatewayAntenna, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "gain":
			if len(subs) > 0 {
				return fmt.Errorf("'gain' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Gain = src.Gain
			} else {
				var zero float32
				dst.Gain = zero
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

var GatewayStatusFieldPathsNested = []string{
	"advanced",
	"antenna_locations",
	"boot_time",
	"ip",
	"metrics",
	"time",
	"versions",
}

var GatewayStatusFieldPathsTopLevel = []string{
	"advanced",
	"antenna_locations",
	"boot_time",
	"ip",
	"metrics",
	"time",
	"versions",
}

func (dst *GatewayStatus) SetFields(src *GatewayStatus, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "time":
			if len(subs) > 0 {
				return fmt.Errorf("'time' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Time = src.Time
			} else {
				var zero time.Time
				dst.Time = zero
			}
		case "boot_time":
			if len(subs) > 0 {
				return fmt.Errorf("'boot_time' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.BootTime = src.BootTime
			} else {
				var zero time.Time
				dst.BootTime = zero
			}
		case "versions":
			if len(subs) > 0 {
				return fmt.Errorf("'versions' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Versions = src.Versions
			} else {
				dst.Versions = nil
			}
		case "antenna_locations":
			if len(subs) > 0 {
				return fmt.Errorf("'antenna_locations' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AntennaLocations = src.AntennaLocations
			} else {
				dst.AntennaLocations = nil
			}
		case "ip":
			if len(subs) > 0 {
				return fmt.Errorf("'ip' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.IP = src.IP
			} else {
				dst.IP = nil
			}
		case "metrics":
			if len(subs) > 0 {
				return fmt.Errorf("'metrics' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Metrics = src.Metrics
			} else {
				dst.Metrics = nil
			}
		case "advanced":
			if len(subs) > 0 {
				return fmt.Errorf("'advanced' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Advanced = src.Advanced
			} else {
				dst.Advanced = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GatewayConnectionStatsFieldPathsNested = []string{
	"connected_at",
	"downlink_count",
	"last_downlink_received_at",
	"last_status",
	"last_status.advanced",
	"last_status.antenna_locations",
	"last_status.boot_time",
	"last_status.ip",
	"last_status.metrics",
	"last_status.time",
	"last_status.versions",
	"last_status_received_at",
	"last_uplink_received_at",
	"protocol",
	"uplink_count",
}

var GatewayConnectionStatsFieldPathsTopLevel = []string{
	"connected_at",
	"downlink_count",
	"last_downlink_received_at",
	"last_status",
	"last_status_received_at",
	"last_uplink_received_at",
	"protocol",
	"uplink_count",
}

func (dst *GatewayConnectionStats) SetFields(src *GatewayConnectionStats, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "connected_at":
			if len(subs) > 0 {
				return fmt.Errorf("'connected_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ConnectedAt = src.ConnectedAt
			} else {
				dst.ConnectedAt = nil
			}
		case "protocol":
			if len(subs) > 0 {
				return fmt.Errorf("'protocol' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Protocol = src.Protocol
			} else {
				var zero string
				dst.Protocol = zero
			}
		case "last_status_received_at":
			if len(subs) > 0 {
				return fmt.Errorf("'last_status_received_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LastStatusReceivedAt = src.LastStatusReceivedAt
			} else {
				dst.LastStatusReceivedAt = nil
			}
		case "last_status":
			if len(subs) > 0 {
				newDst := dst.LastStatus
				if newDst == nil {
					newDst = &GatewayStatus{}
					dst.LastStatus = newDst
				}
				var newSrc *GatewayStatus
				if src != nil {
					newSrc = src.LastStatus
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.LastStatus = src.LastStatus
				} else {
					dst.LastStatus = nil
				}
			}
		case "last_uplink_received_at":
			if len(subs) > 0 {
				return fmt.Errorf("'last_uplink_received_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LastUplinkReceivedAt = src.LastUplinkReceivedAt
			} else {
				dst.LastUplinkReceivedAt = nil
			}
		case "uplink_count":
			if len(subs) > 0 {
				return fmt.Errorf("'uplink_count' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UplinkCount = src.UplinkCount
			} else {
				var zero uint64
				dst.UplinkCount = zero
			}
		case "last_downlink_received_at":
			if len(subs) > 0 {
				return fmt.Errorf("'last_downlink_received_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LastDownlinkReceivedAt = src.LastDownlinkReceivedAt
			} else {
				dst.LastDownlinkReceivedAt = nil
			}
		case "downlink_count":
			if len(subs) > 0 {
				return fmt.Errorf("'downlink_count' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DownlinkCount = src.DownlinkCount
			} else {
				var zero uint64
				dst.DownlinkCount = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
