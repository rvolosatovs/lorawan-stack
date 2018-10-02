// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/message_services.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import time "time"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

func (this *ProcessUplinkMessageRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceVersionIDs)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceVersionIDs", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Message)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Message", err)
	}
	return nil
}
func (this *ProcessDownlinkMessageRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceVersionIDs)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceVersionIDs", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Message)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Message", err)
	}
	return nil
}
