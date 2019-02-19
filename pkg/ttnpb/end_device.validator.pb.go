// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/end_device.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import regexp "regexp"
import fmt "fmt"
import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/golang/protobuf/ptypes/duration"
import _ "github.com/golang/protobuf/ptypes/struct"
import _ "github.com/golang/protobuf/ptypes/timestamp"
import _ "github.com/golang/protobuf/ptypes/wrappers"
import _ "github.com/mwitkow/go-proto-validators"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import _ "google.golang.org/genproto/protobuf/field_mask"

import time "time"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

func (this *Session) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.SessionKeys)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("SessionKeys", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.StartedAt)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("StartedAt", err)
	}
	return nil
}
func (this *MACParameters) Validate() error {
	for _, item := range this.Channels {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Channels", err)
			}
		}
	}
	return nil
}
func (this *MACParameters_Channel) Validate() error {
	return nil
}
func (this *EndDeviceBrand) Validate() error {
	return nil
}
func (this *EndDeviceModel) Validate() error {
	return nil
}

var _regex_EndDeviceVersionIdentifiers_BrandID = regexp.MustCompile(`^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$`)
var _regex_EndDeviceVersionIdentifiers_ModelID = regexp.MustCompile(`^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$`)

func (this *EndDeviceVersionIdentifiers) Validate() error {
	if !_regex_EndDeviceVersionIdentifiers_BrandID.MatchString(this.BrandID) {
		return github_com_mwitkow_go_proto_validators.FieldError("BrandID", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$"`, this.BrandID))
	}
	if !(len(this.BrandID) < 37) {
		return github_com_mwitkow_go_proto_validators.FieldError("BrandID", fmt.Errorf(`value '%v' must length be less than '37'`, this.BrandID))
	}
	if !_regex_EndDeviceVersionIdentifiers_ModelID.MatchString(this.ModelID) {
		return github_com_mwitkow_go_proto_validators.FieldError("ModelID", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$"`, this.ModelID))
	}
	if !(len(this.ModelID) < 37) {
		return github_com_mwitkow_go_proto_validators.FieldError("ModelID", fmt.Errorf(`value '%v' must length be less than '37'`, this.ModelID))
	}
	return nil
}
func (this *EndDeviceVersion) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceVersionIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceVersionIdentifiers", err)
	}
	if this.DefaultMACSettings != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DefaultMACSettings); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DefaultMACSettings", err)
		}
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.DefaultFormatters)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("DefaultFormatters", err)
	}
	return nil
}
func (this *MACSettings) Validate() error {
	if this.ClassBTimeout != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ClassBTimeout); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ClassBTimeout", err)
		}
	}
	if this.PingSlotPeriodicity != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PingSlotPeriodicity); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PingSlotPeriodicity", err)
		}
	}
	if this.PingSlotDataRateIndex != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PingSlotDataRateIndex); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PingSlotDataRateIndex", err)
		}
	}
	if this.ClassCTimeout != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ClassCTimeout); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ClassCTimeout", err)
		}
	}
	if this.Rx1Delay != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Rx1Delay); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Rx1Delay", err)
		}
	}
	if this.Rx1DataRateOffset != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Rx1DataRateOffset); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Rx1DataRateOffset", err)
		}
	}
	if this.Rx2DataRateIndex != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Rx2DataRateIndex); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Rx2DataRateIndex", err)
		}
	}
	if this.MaxDutyCycle != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.MaxDutyCycle); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("MaxDutyCycle", err)
		}
	}
	if this.Supports32BitFCnt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Supports32BitFCnt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Supports32BitFCnt", err)
		}
	}
	if this.UseADR != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UseADR); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UseADR", err)
		}
	}
	if this.ADRMargin != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ADRMargin); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ADRMargin", err)
		}
	}
	if this.ResetsFCnt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ResetsFCnt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ResetsFCnt", err)
		}
	}
	if this.StatusTimePeriodicity != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.StatusTimePeriodicity); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("StatusTimePeriodicity", err)
		}
	}
	if this.StatusCountPeriodicity != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.StatusCountPeriodicity); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("StatusCountPeriodicity", err)
		}
	}
	if this.DesiredRx1Delay != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DesiredRx1Delay); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DesiredRx1Delay", err)
		}
	}
	if this.DesiredRx1DataRateOffset != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DesiredRx1DataRateOffset); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DesiredRx1DataRateOffset", err)
		}
	}
	if this.DesiredRx2DataRateIndex != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.DesiredRx2DataRateIndex); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("DesiredRx2DataRateIndex", err)
		}
	}
	return nil
}
func (this *MACSettings_DataRateIndexValue) Validate() error {
	return nil
}
func (this *MACSettings_PingSlotPeriodValue) Validate() error {
	return nil
}
func (this *MACSettings_AggregatedDutyCycleValue) Validate() error {
	return nil
}
func (this *MACSettings_RxDelayValue) Validate() error {
	return nil
}
func (this *MACState) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.CurrentParameters)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("CurrentParameters", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.DesiredParameters)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("DesiredParameters", err)
	}
	if this.LastConfirmedDownlinkAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.LastConfirmedDownlinkAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("LastConfirmedDownlinkAt", err)
		}
	}
	if this.PendingApplicationDownlink != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PendingApplicationDownlink); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PendingApplicationDownlink", err)
		}
	}
	for _, item := range this.QueuedResponses {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("QueuedResponses", err)
			}
		}
	}
	for _, item := range this.PendingRequests {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PendingRequests", err)
			}
		}
	}
	if this.QueuedJoinAccept != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.QueuedJoinAccept); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("QueuedJoinAccept", err)
		}
	}
	if this.PendingJoinRequest != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PendingJoinRequest); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PendingJoinRequest", err)
		}
	}
	return nil
}
func (this *MACState_JoinAccept) Validate() error {
	if !(len(this.Payload) > 11) {
		return github_com_mwitkow_go_proto_validators.FieldError("Payload", fmt.Errorf(`value '%v' must length be greater than '11'`, this.Payload))
	}
	if !(len(this.Payload) < 29) {
		return github_com_mwitkow_go_proto_validators.FieldError("Payload", fmt.Errorf(`value '%v' must length be less than '29'`, this.Payload))
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Request)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Request", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.Keys)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("Keys", err)
	}
	return nil
}

var _regex_EndDevice_ProvisionerID = regexp.MustCompile(`^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$`)

func (this *EndDevice) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.CreatedAt)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.UpdatedAt)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("UpdatedAt", err)
	}
	// Validation of proto3 map<> fields is unsupported.
	if this.VersionIDs != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.VersionIDs); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("VersionIDs", err)
		}
	}
	// Validation of proto3 map<> fields is unsupported.
	if this.RootKeys != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.RootKeys); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("RootKeys", err)
		}
	}
	if this.MACSettings != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.MACSettings); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("MACSettings", err)
		}
	}
	if this.MACState != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.MACState); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("MACState", err)
		}
	}
	if this.Session != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Session); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Session", err)
		}
	}
	if this.PendingSession != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PendingSession); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PendingSession", err)
		}
	}
	if this.LastDevStatusReceivedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.LastDevStatusReceivedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("LastDevStatusReceivedAt", err)
		}
	}
	for _, item := range this.RecentADRUplinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RecentADRUplinks", err)
			}
		}
	}
	for _, item := range this.RecentUplinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RecentUplinks", err)
			}
		}
	}
	for _, item := range this.RecentDownlinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("RecentDownlinks", err)
			}
		}
	}
	for _, item := range this.QueuedApplicationDownlinks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("QueuedApplicationDownlinks", err)
			}
		}
	}
	if this.Formatters != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Formatters); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Formatters", err)
		}
	}
	if !_regex_EndDevice_ProvisionerID.MatchString(this.ProvisionerID) {
		return github_com_mwitkow_go_proto_validators.FieldError("ProvisionerID", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$"`, this.ProvisionerID))
	}
	if !(len(this.ProvisionerID) < 37) {
		return github_com_mwitkow_go_proto_validators.FieldError("ProvisionerID", fmt.Errorf(`value '%v' must length be less than '37'`, this.ProvisionerID))
	}
	if this.ProvisioningData != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ProvisioningData); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ProvisioningData", err)
		}
	}
	return nil
}
func (this *EndDevices) Validate() error {
	for _, item := range this.EndDevices {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("EndDevices", err)
			}
		}
	}
	return nil
}
func (this *CreateEndDeviceRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDevice)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDevice", err)
	}
	return nil
}
func (this *UpdateEndDeviceRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDevice)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDevice", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
func (this *GetEndDeviceRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDeviceIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDeviceIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
func (this *ListEndDevicesRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.ApplicationIdentifiers)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("ApplicationIdentifiers", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
func (this *SetEndDeviceRequest) Validate() error {
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.EndDevice)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDevice", err)
	}
	if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(&(this.FieldMask)); err != nil {
		return github_com_mwitkow_go_proto_validators.FieldError("FieldMask", err)
	}
	return nil
}
