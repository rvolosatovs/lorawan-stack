// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/applicationserver_web.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import regexp "regexp"
import fmt "fmt"
import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/golang/protobuf/ptypes/empty"
import _ "github.com/golang/protobuf/ptypes/timestamp"
import _ "github.com/mwitkow/go-proto-validators"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import _ "google.golang.org/genproto/protobuf/field_mask"

import time "time"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

var _regex_ApplicationWebhookIdentifiers_WebhookID = regexp.MustCompile(`^[a-z0-9](?:[-]?[a-z0-9]){1,}$`)

func (this *ApplicationWebhookIdentifiers) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationIdentifiers", err)
	}
	if !_regex_ApplicationWebhookIdentifiers_WebhookID.MatchString(this.WebhookID) {
		return github_com_mwitkow_go_proto_validators.FieldError("WebhookID", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z0-9](?:[-]?[a-z0-9]){1,}$"`, this.WebhookID))
	}
	if !(len(this.WebhookID) < 37) {
		return github_com_mwitkow_go_proto_validators.FieldError("WebhookID", fmt.Errorf(`value '%v' must length be less than '37'`, this.WebhookID))
	}
	return nil
}
func (this *ApplicationWebhook) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationWebhookIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationWebhookIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.CreatedAt)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.UpdatedAt)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("UpdatedAt", err)
	}
	// Validation of proto3 map<> fields is unsupported.
	if this.UplinkMessage != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UplinkMessage); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UplinkMessage", err)
		}
	}
	if this.JoinAccept != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.JoinAccept); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("JoinAccept", err)
		}
	}
	if this.DownlinkAck != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DownlinkAck); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DownlinkAck", err)
		}
	}
	if this.DownlinkNack != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DownlinkNack); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DownlinkNack", err)
		}
	}
	if this.DownlinkSent != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DownlinkSent); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DownlinkSent", err)
		}
	}
	if this.DownlinkFailed != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DownlinkFailed); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DownlinkFailed", err)
		}
	}
	if this.DownlinkQueued != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DownlinkQueued); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DownlinkQueued", err)
		}
	}
	if this.LocationSolved != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.LocationSolved); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("LocationSolved", err)
		}
	}
	return nil
}
func (this *ApplicationWebhook_Message) Validate() error {
	return nil
}
func (this *ApplicationWebhooks) Validate() error {
	for _, item := range this.Webhooks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Webhooks", err)
			}
		}
	}
	return nil
}
func (this *ApplicationWebhookFormats) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *GetApplicationWebhookRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationWebhookIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationWebhookIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
func (this *ListApplicationWebhooksRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
func (this *SetApplicationWebhookRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationWebhook)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationWebhook", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
