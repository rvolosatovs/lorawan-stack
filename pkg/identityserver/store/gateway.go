// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
)

// Gateway model.
type Gateway struct {
	Model
	SoftDelete

	GatewayEUI *EUI64 `gorm:"unique_index:eui;type:VARCHAR(16);column:gateway_eui"`

	// BEGIN common fields
	GatewayID   string        `gorm:"unique_index:id;type:VARCHAR(36)"`
	Name        string        `gorm:"type:VARCHAR"`
	Description string        `gorm:"type:TEXT"`
	Attributes  []Attribute   `gorm:"polymorphic:Entity;polymorphic_value:gateway"`
	ContactInfo []ContactInfo `gorm:"polymorphic:Entity;polymorphic_value:gateway"`
	APIKeys     []APIKey      `gorm:"polymorphic:Entity;polymorphic_value:gateway"`
	Memberships []Membership  `gorm:"polymorphic:Entity;polymorphic_value:gateway"`
	// END common fields

	BrandID         string `gorm:"type:VARCHAR"`
	ModelID         string `gorm:"type:VARCHAR"`
	HardwareVersion string `gorm:"type:VARCHAR"`
	FirmwareVersion string `gorm:"type:VARCHAR"`

	GatewayServerAddress string `gorm:"type:VARCHAR"`

	AutoUpdate    bool
	UpdateChannel string `gorm:"type:VARCHAR"`

	FrequencyPlanID string `gorm:"type:VARCHAR"`

	StatusPublic   bool
	LocationPublic bool

	ScheduleDownlinkLate   bool
	EnforceDutyCycle       bool
	DownlinkPathConstraint int
}

func init() {
	registerModel(&Gateway{})
}

// functions to set fields from the gateway model into the gateway proto.
var gatewayPBSetters = map[string]func(*ttnpb.Gateway, *Gateway){
	"ids.eui":                 func(pb *ttnpb.Gateway, gtw *Gateway) { pb.EUI = (*types.EUI64)(gtw.GatewayEUI) }, // can we do this?
	nameField:                 func(pb *ttnpb.Gateway, gtw *Gateway) { pb.Name = gtw.Name },
	descriptionField:          func(pb *ttnpb.Gateway, gtw *Gateway) { pb.Description = gtw.Description },
	attributesField:           func(pb *ttnpb.Gateway, gtw *Gateway) { pb.Attributes = attributes(gtw.Attributes).toMap() },
	contactInfoField:          func(pb *ttnpb.Gateway, gtw *Gateway) { pb.ContactInfo = contactInfos(gtw.ContactInfo).toPB() },
	brandIDField:              func(pb *ttnpb.Gateway, gtw *Gateway) { pb.BrandID = gtw.BrandID },
	modelIDField:              func(pb *ttnpb.Gateway, gtw *Gateway) { pb.ModelID = gtw.ModelID },
	hardwareVersionField:      func(pb *ttnpb.Gateway, gtw *Gateway) { pb.HardwareVersion = gtw.HardwareVersion },
	firmwareVersionField:      func(pb *ttnpb.Gateway, gtw *Gateway) { pb.FirmwareVersion = gtw.FirmwareVersion },
	gatewayServerAddressField: func(pb *ttnpb.Gateway, gtw *Gateway) { pb.GatewayServerAddress = gtw.GatewayServerAddress },
	autoUpdateField:           func(pb *ttnpb.Gateway, gtw *Gateway) { pb.AutoUpdate = gtw.AutoUpdate },
	updateChannelField:        func(pb *ttnpb.Gateway, gtw *Gateway) { pb.UpdateChannel = gtw.UpdateChannel },
	frequencyPlanIDField:      func(pb *ttnpb.Gateway, gtw *Gateway) { pb.FrequencyPlanID = gtw.FrequencyPlanID },
	statusPublicField:         func(pb *ttnpb.Gateway, gtw *Gateway) { pb.StatusPublic = gtw.StatusPublic },
	locationPublicField:       func(pb *ttnpb.Gateway, gtw *Gateway) { pb.LocationPublic = gtw.LocationPublic },
	scheduleDownlinkLateField: func(pb *ttnpb.Gateway, gtw *Gateway) { pb.ScheduleDownlinkLate = gtw.ScheduleDownlinkLate },
	enforceDutyCycleField:     func(pb *ttnpb.Gateway, gtw *Gateway) { pb.EnforceDutyCycle = gtw.EnforceDutyCycle },
	downlinkPathConstraintField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		pb.DownlinkPathConstraint = ttnpb.DownlinkPathConstraint(gtw.DownlinkPathConstraint)
	},
}

// functions to set fields from the gateway proto into the gateway model.
var gatewayModelSetters = map[string]func(*Gateway, *ttnpb.Gateway){
	"ids.eui":        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.GatewayEUI = (*EUI64)(pb.EUI) }, // can we do this?
	nameField:        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.Name = pb.Name },
	descriptionField: func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.Description = pb.Description },
	attributesField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		gtw.Attributes = attributes(gtw.Attributes).updateFromMap(pb.Attributes)
	},
	contactInfoField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		gtw.ContactInfo = contactInfos(gtw.ContactInfo).updateFromPB(pb.ContactInfo)
	},
	brandIDField:                func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.BrandID = pb.BrandID },
	modelIDField:                func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.ModelID = pb.ModelID },
	hardwareVersionField:        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.HardwareVersion = pb.HardwareVersion },
	firmwareVersionField:        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.FirmwareVersion = pb.FirmwareVersion },
	gatewayServerAddressField:   func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.GatewayServerAddress = pb.GatewayServerAddress },
	autoUpdateField:             func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.AutoUpdate = pb.AutoUpdate },
	updateChannelField:          func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.UpdateChannel = pb.UpdateChannel },
	frequencyPlanIDField:        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.FrequencyPlanID = pb.FrequencyPlanID },
	statusPublicField:           func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.StatusPublic = pb.StatusPublic },
	locationPublicField:         func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.LocationPublic = pb.LocationPublic },
	scheduleDownlinkLateField:   func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.ScheduleDownlinkLate = pb.ScheduleDownlinkLate },
	enforceDutyCycleField:       func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.EnforceDutyCycle = pb.EnforceDutyCycle },
	downlinkPathConstraintField: func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.DownlinkPathConstraint = int(pb.DownlinkPathConstraint) },
}

// fieldMask to use if a nil or empty fieldmask is passed.
var defaultGatewayFieldMask = &pbtypes.FieldMask{}

func init() {
	paths := make([]string, 0, len(gatewayPBSetters))
	for path := range gatewayPBSetters {
		paths = append(paths, path)
	}
	defaultGatewayFieldMask.Paths = paths
}

// fieldmask path to column name in gateways table.
var gatewayColumnNames = map[string]string{
	"ids.gateway_id":            "gateway_id",
	"ids.eui":                   "gateway_eui",
	attributesField:             "",
	contactInfoField:            "",
	nameField:                   nameField,
	descriptionField:            descriptionField,
	gatewayServerAddressField:   gatewayServerAddressField,
	brandIDField:                "brand_id",
	modelIDField:                "model_id",
	hardwareVersionField:        "hardware_version",
	firmwareVersionField:        "firmware_version",
	autoUpdateField:             autoUpdateField,
	updateChannelField:          updateChannelField,
	frequencyPlanIDField:        frequencyPlanIDField,
	statusPublicField:           statusPublicField,
	locationPublicField:         locationPublicField,
	scheduleDownlinkLateField:   scheduleDownlinkLateField,
	enforceDutyCycleField:       enforceDutyCycleField,
	downlinkPathConstraintField: downlinkPathConstraintField,
}

func (gtw Gateway) toPB(pb *ttnpb.Gateway, fieldMask *pbtypes.FieldMask) {
	pb.GatewayIdentifiers.GatewayID = gtw.GatewayID
	pb.CreatedAt = cleanTime(gtw.CreatedAt)
	pb.UpdatedAt = cleanTime(gtw.UpdatedAt)
	if fieldMask == nil || len(fieldMask.Paths) == 0 {
		fieldMask = defaultGatewayFieldMask
	}
	for _, path := range fieldMask.Paths {
		if setter, ok := gatewayPBSetters[path]; ok {
			setter(pb, &gtw)
		}
	}
}

func (gtw *Gateway) fromPB(pb *ttnpb.Gateway, fieldMask *pbtypes.FieldMask) (columns []string) {
	if fieldMask == nil || len(fieldMask.Paths) == 0 {
		fieldMask = defaultGatewayFieldMask
	}
	for _, path := range fieldMask.Paths {
		if setter, ok := gatewayModelSetters[path]; ok {
			setter(gtw, pb)
			columnName, ok := gatewayColumnNames[path]
			if !ok {
				columnName = path
			}
			if columnName != "" {
				columns = append(columns, columnName)
			}
			continue
		}
	}
	return
}
