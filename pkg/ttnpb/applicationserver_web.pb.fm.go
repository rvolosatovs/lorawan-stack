// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"

	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
)

var ApplicationWebhookIdentifiersFieldPathsNested = []string{
	"application_ids",
	"application_ids.application_id",
	"webhook_id",
}

var ApplicationWebhookIdentifiersFieldPathsTopLevel = []string{
	"application_ids",
	"webhook_id",
}

func (dst *ApplicationWebhookIdentifiers) SetFields(src *ApplicationWebhookIdentifiers, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "application_ids":
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
		case "webhook_id":
			if len(subs) > 0 {
				return fmt.Errorf("'webhook_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.WebhookID = src.WebhookID
			} else {
				var zero string
				dst.WebhookID = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ApplicationWebhookFieldPathsNested = []string{
	"base_url",
	"created_at",
	"downlink_ack",
	"downlink_ack.path",
	"downlink_failed",
	"downlink_failed.path",
	"downlink_nack",
	"downlink_nack.path",
	"downlink_queued",
	"downlink_queued.path",
	"downlink_sent",
	"downlink_sent.path",
	"format",
	"headers",
	"ids",
	"ids.application_ids",
	"ids.application_ids.application_id",
	"ids.webhook_id",
	"join_accept",
	"join_accept.path",
	"location_solved",
	"location_solved.path",
	"updated_at",
	"uplink_message",
	"uplink_message.path",
}

var ApplicationWebhookFieldPathsTopLevel = []string{
	"base_url",
	"created_at",
	"downlink_ack",
	"downlink_failed",
	"downlink_nack",
	"downlink_queued",
	"downlink_sent",
	"format",
	"headers",
	"ids",
	"join_accept",
	"location_solved",
	"updated_at",
	"uplink_message",
}

func (dst *ApplicationWebhook) SetFields(src *ApplicationWebhook, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "ids":
			if len(subs) > 0 {
				newDst := &dst.ApplicationWebhookIdentifiers
				var newSrc *ApplicationWebhookIdentifiers
				if src != nil {
					newSrc = &src.ApplicationWebhookIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationWebhookIdentifiers = src.ApplicationWebhookIdentifiers
				} else {
					var zero ApplicationWebhookIdentifiers
					dst.ApplicationWebhookIdentifiers = zero
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
		case "base_url":
			if len(subs) > 0 {
				return fmt.Errorf("'base_url' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.BaseURL = src.BaseURL
			} else {
				var zero string
				dst.BaseURL = zero
			}
		case "headers":
			if len(subs) > 0 {
				return fmt.Errorf("'headers' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Headers = src.Headers
			} else {
				dst.Headers = nil
			}
		case "format":
			if len(subs) > 0 {
				return fmt.Errorf("'format' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Format = src.Format
			} else {
				var zero string
				dst.Format = zero
			}
		case "uplink_message":
			if len(subs) > 0 {
				newDst := dst.UplinkMessage
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.UplinkMessage = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.UplinkMessage
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.UplinkMessage = src.UplinkMessage
				} else {
					dst.UplinkMessage = nil
				}
			}
		case "join_accept":
			if len(subs) > 0 {
				newDst := dst.JoinAccept
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.JoinAccept = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.JoinAccept
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.JoinAccept = src.JoinAccept
				} else {
					dst.JoinAccept = nil
				}
			}
		case "downlink_ack":
			if len(subs) > 0 {
				newDst := dst.DownlinkAck
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.DownlinkAck = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.DownlinkAck
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkAck = src.DownlinkAck
				} else {
					dst.DownlinkAck = nil
				}
			}
		case "downlink_nack":
			if len(subs) > 0 {
				newDst := dst.DownlinkNack
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.DownlinkNack = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.DownlinkNack
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkNack = src.DownlinkNack
				} else {
					dst.DownlinkNack = nil
				}
			}
		case "downlink_sent":
			if len(subs) > 0 {
				newDst := dst.DownlinkSent
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.DownlinkSent = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.DownlinkSent
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkSent = src.DownlinkSent
				} else {
					dst.DownlinkSent = nil
				}
			}
		case "downlink_failed":
			if len(subs) > 0 {
				newDst := dst.DownlinkFailed
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.DownlinkFailed = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.DownlinkFailed
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkFailed = src.DownlinkFailed
				} else {
					dst.DownlinkFailed = nil
				}
			}
		case "downlink_queued":
			if len(subs) > 0 {
				newDst := dst.DownlinkQueued
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.DownlinkQueued = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.DownlinkQueued
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DownlinkQueued = src.DownlinkQueued
				} else {
					dst.DownlinkQueued = nil
				}
			}
		case "location_solved":
			if len(subs) > 0 {
				newDst := dst.LocationSolved
				if newDst == nil {
					newDst = &ApplicationWebhook_Message{}
					dst.LocationSolved = newDst
				}
				var newSrc *ApplicationWebhook_Message
				if src != nil {
					newSrc = src.LocationSolved
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.LocationSolved = src.LocationSolved
				} else {
					dst.LocationSolved = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ApplicationWebhook_MessageFieldPathsNested = []string{
	"path",
}

var ApplicationWebhook_MessageFieldPathsTopLevel = []string{
	"path",
}

func (dst *ApplicationWebhook_Message) SetFields(src *ApplicationWebhook_Message, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "path":
			if len(subs) > 0 {
				return fmt.Errorf("'path' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Path = src.Path
			} else {
				var zero string
				dst.Path = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ApplicationWebhooksFieldPathsNested = []string{
	"webhooks",
}

var ApplicationWebhooksFieldPathsTopLevel = []string{
	"webhooks",
}

func (dst *ApplicationWebhooks) SetFields(src *ApplicationWebhooks, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "webhooks":
			if len(subs) > 0 {
				return fmt.Errorf("'webhooks' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Webhooks = src.Webhooks
			} else {
				dst.Webhooks = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ApplicationWebhookFormatsFieldPathsNested = []string{
	"formats",
}

var ApplicationWebhookFormatsFieldPathsTopLevel = []string{
	"formats",
}

func (dst *ApplicationWebhookFormats) SetFields(src *ApplicationWebhookFormats, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "formats":
			if len(subs) > 0 {
				return fmt.Errorf("'formats' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Formats = src.Formats
			} else {
				dst.Formats = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var GetApplicationWebhookRequestFieldPathsNested = []string{
	"field_mask",
	"ids",
	"ids.application_ids",
	"ids.application_ids.application_id",
	"ids.webhook_id",
}

var GetApplicationWebhookRequestFieldPathsTopLevel = []string{
	"field_mask",
	"ids",
}

func (dst *GetApplicationWebhookRequest) SetFields(src *GetApplicationWebhookRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "ids":
			if len(subs) > 0 {
				newDst := &dst.ApplicationWebhookIdentifiers
				var newSrc *ApplicationWebhookIdentifiers
				if src != nil {
					newSrc = &src.ApplicationWebhookIdentifiers
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationWebhookIdentifiers = src.ApplicationWebhookIdentifiers
				} else {
					var zero ApplicationWebhookIdentifiers
					dst.ApplicationWebhookIdentifiers = zero
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

var ListApplicationWebhooksRequestFieldPathsNested = []string{
	"application_ids",
	"application_ids.application_id",
	"field_mask",
}

var ListApplicationWebhooksRequestFieldPathsTopLevel = []string{
	"application_ids",
	"field_mask",
}

func (dst *ListApplicationWebhooksRequest) SetFields(src *ListApplicationWebhooksRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "application_ids":
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

var SetApplicationWebhookRequestFieldPathsNested = []string{
	"field_mask",
	"webhook",
	"webhook.base_url",
	"webhook.created_at",
	"webhook.downlink_ack",
	"webhook.downlink_ack.path",
	"webhook.downlink_failed",
	"webhook.downlink_failed.path",
	"webhook.downlink_nack",
	"webhook.downlink_nack.path",
	"webhook.downlink_queued",
	"webhook.downlink_queued.path",
	"webhook.downlink_sent",
	"webhook.downlink_sent.path",
	"webhook.format",
	"webhook.headers",
	"webhook.ids",
	"webhook.ids.application_ids",
	"webhook.ids.application_ids.application_id",
	"webhook.ids.webhook_id",
	"webhook.join_accept",
	"webhook.join_accept.path",
	"webhook.location_solved",
	"webhook.location_solved.path",
	"webhook.updated_at",
	"webhook.uplink_message",
	"webhook.uplink_message.path",
}

var SetApplicationWebhookRequestFieldPathsTopLevel = []string{
	"field_mask",
	"webhook",
}

func (dst *SetApplicationWebhookRequest) SetFields(src *SetApplicationWebhookRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "webhook":
			if len(subs) > 0 {
				newDst := &dst.ApplicationWebhook
				var newSrc *ApplicationWebhook
				if src != nil {
					newSrc = &src.ApplicationWebhook
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationWebhook = src.ApplicationWebhook
				} else {
					var zero ApplicationWebhook
					dst.ApplicationWebhook = zero
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
