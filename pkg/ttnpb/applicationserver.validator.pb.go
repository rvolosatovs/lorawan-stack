// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/applicationserver.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import fmt "fmt"
import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/golang/protobuf/ptypes/empty"
import _ "github.com/mwitkow/go-proto-validators"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import _ "google.golang.org/genproto/protobuf/field_mask"

import time "time"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

func (this *ApplicationLink) Validate() error {
	if this.APIKey == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("APIKey", fmt.Errorf(`value '%v' must not be an empty string`, this.APIKey))
	}
	if this.DefaultFormatters != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DefaultFormatters); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DefaultFormatters", err)
		}
	}
	return nil
}
func (this *GetApplicationLinkRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
func (this *SetApplicationLinkRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationLink)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationLink", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
