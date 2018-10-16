// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/joinserver.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import regexp "regexp"
import fmt "fmt"
import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/golang/protobuf/ptypes/empty"
import _ "github.com/mwitkow/go-proto-validators"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import time "time"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

var _regex_SessionKeyRequest_SessionKeyID = regexp.MustCompile(`^[a-z0-9](?:[-]?[a-z0-9]){1,35}$`)

func (this *SessionKeyRequest) Validate() error {
	if !_regex_SessionKeyRequest_SessionKeyID.MatchString(this.SessionKeyID) {
		return github_com_mwitkow_go_proto_validators.FieldError("SessionKeyID", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z0-9](?:[-]?[a-z0-9]){1,35}$"`, this.SessionKeyID))
	}
	if this.SessionKeyID == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("SessionKeyID", fmt.Errorf(`value '%v' must not be an empty string`, this.SessionKeyID))
	}
	return nil
}
func (this *NwkSKeysResponse) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FNwkSIntKey)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FNwkSIntKey", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.SNwkSIntKey)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("SNwkSIntKey", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.NwkSEncKey)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("NwkSEncKey", err)
	}
	return nil
}
func (this *AppSKeyResponse) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.AppSKey)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("AppSKey", err)
	}
	return nil
}
