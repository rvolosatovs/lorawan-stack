// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/messages.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import regexp "regexp"
import fmt "fmt"
import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/golang/protobuf/ptypes/struct"
import _ "github.com/golang/protobuf/ptypes/timestamp"
import _ "github.com/mwitkow/go-proto-validators"

import time "time"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

func (this *UplinkMessage) Validate() error {
	if this.Payload != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Payload); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Payload", err)
		}
	}
	if this.EndDeviceIDs != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.EndDeviceIDs); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIDs", err)
		}
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Settings)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Settings", err)
	}
	for _, item := range this.RxMetadata {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RxMetadata", err)
			}
		}
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ReceivedAt)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ReceivedAt", err)
	}
	return nil
}
func (this *DownlinkMessage) Validate() error {
	if this.Payload != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Payload); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Payload", err)
		}
	}
	if this.EndDeviceIDs != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.EndDeviceIDs); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIDs", err)
		}
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Settings)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Settings", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.TxMetadata)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("TxMetadata", err)
	}
	return nil
}
func (this *TxAcknowledgment) Validate() error {
	return nil
}

var _regex_ApplicationUplink_SessionKeyID = regexp.MustCompile(`(?:0[xX])(?:[0-9a-fA-F]{2})+$|^$`)

func (this *ApplicationUplink) Validate() error {
	if !_regex_ApplicationUplink_SessionKeyID.MatchString(this.SessionKeyID) {
		return github_com_mwitkow_go_proto_validators.FieldError("SessionKeyID", fmt.Errorf(`value '%v' must be a string conforming to regex "(?:0[xX])(?:[0-9a-fA-F]{2})+$|^$"`, this.SessionKeyID))
	}
	if this.DecodedPayload != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DecodedPayload); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DecodedPayload", err)
		}
	}
	for _, item := range this.RxMetadata {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RxMetadata", err)
			}
		}
	}
	return nil
}
func (this *ApplicationLocation) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Location)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Location", err)
	}
	// Validation of proto3 map<> fields is unsupported.
	return nil
}

var _regex_ApplicationJoinAccept_SessionKeyID = regexp.MustCompile(`(?:0[xX])(?:[0-9a-fA-F]{2})+$|^$`)

func (this *ApplicationJoinAccept) Validate() error {
	if !_regex_ApplicationJoinAccept_SessionKeyID.MatchString(this.SessionKeyID) {
		return github_com_mwitkow_go_proto_validators.FieldError("SessionKeyID", fmt.Errorf(`value '%v' must be a string conforming to regex "(?:0[xX])(?:[0-9a-fA-F]{2})+$|^$"`, this.SessionKeyID))
	}
	if this.AppSKey != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.AppSKey); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("AppSKey", err)
		}
	}
	for _, item := range this.InvalidatedDownlinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("InvalidatedDownlinks", err)
			}
		}
	}
	return nil
}

var _regex_ApplicationDownlink_SessionKeyID = regexp.MustCompile(`(?:0[xX])(?:[0-9a-fA-F]{2})+$|^$`)

func (this *ApplicationDownlink) Validate() error {
	if !_regex_ApplicationDownlink_SessionKeyID.MatchString(this.SessionKeyID) {
		return github_com_mwitkow_go_proto_validators.FieldError("SessionKeyID", fmt.Errorf(`value '%v' must be a string conforming to regex "(?:0[xX])(?:[0-9a-fA-F]{2})+$|^$"`, this.SessionKeyID))
	}
	if this.DecodedPayload != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DecodedPayload); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DecodedPayload", err)
		}
	}
	if this.ClassBC != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ClassBC); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ClassBC", err)
		}
	}
	return nil
}
func (this *ApplicationDownlink_ClassBC) Validate() error {
	for _, item := range this.Gateways {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Gateways", err)
			}
		}
	}
	if this.Time != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Time); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Time", err)
		}
	}
	return nil
}
func (this *ApplicationDownlink_ClassBC_GatewayAntennaIdentifiers) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.GatewayIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("GatewayIdentifiers", err)
	}
	return nil
}
func (this *ApplicationDownlinks) Validate() error {
	for _, item := range this.Downlinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Downlinks", err)
			}
		}
	}
	return nil
}
func (this *ApplicationDownlinkFailed) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationDownlink)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationDownlink", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Error)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Error", err)
	}
	return nil
}
func (this *ApplicationInvalidatedDownlinks) Validate() error {
	for _, item := range this.Downlinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Downlinks", err)
			}
		}
	}
	return nil
}
func (this *ApplicationUp) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIdentifiers", err)
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_UplinkMessage); ok {
		if oneOfNester.UplinkMessage != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.UplinkMessage); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("UplinkMessage", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_JoinAccept); ok {
		if oneOfNester.JoinAccept != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.JoinAccept); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("JoinAccept", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_DownlinkAck); ok {
		if oneOfNester.DownlinkAck != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.DownlinkAck); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DownlinkAck", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_DownlinkNack); ok {
		if oneOfNester.DownlinkNack != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.DownlinkNack); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DownlinkNack", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_DownlinkSent); ok {
		if oneOfNester.DownlinkSent != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.DownlinkSent); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DownlinkSent", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_DownlinkFailed); ok {
		if oneOfNester.DownlinkFailed != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.DownlinkFailed); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DownlinkFailed", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_DownlinkQueued); ok {
		if oneOfNester.DownlinkQueued != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.DownlinkQueued); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DownlinkQueued", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_DownlinkQueueInvalidated); ok {
		if oneOfNester.DownlinkQueueInvalidated != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.DownlinkQueueInvalidated); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DownlinkQueueInvalidated", err)
			}
		}
	}
	if oneOfNester, ok := this.GetUp().(*ApplicationUp_LocationSolved); ok {
		if oneOfNester.LocationSolved != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.LocationSolved); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("LocationSolved", err)
			}
		}
	}
	return nil
}
func (this *MessagePayloadFormatters) Validate() error {
	return nil
}
func (this *DownlinkQueueRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIdentifiers", err)
	}
	for _, item := range this.Downlinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Downlinks", err)
			}
		}
	}
	return nil
}
