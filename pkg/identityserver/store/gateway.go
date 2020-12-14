// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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
	"bytes"
	"sort"
	"strings"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// Gateway model.
type Gateway struct {
	Model
	SoftDelete

	GatewayEUI *EUI64 `gorm:"unique_index:gateway_eui_index;type:VARCHAR(16);column:gateway_eui"`

	// BEGIN common fields
	GatewayID   string       `gorm:"unique_index:gateway_id_index;type:VARCHAR(36);not null"`
	Name        string       `gorm:"type:VARCHAR"`
	Description string       `gorm:"type:TEXT"`
	Attributes  []Attribute  `gorm:"polymorphic:Entity;polymorphic_value:gateway"`
	APIKeys     []APIKey     `gorm:"polymorphic:Entity;polymorphic_value:gateway"`
	Memberships []Membership `gorm:"polymorphic:Entity;polymorphic_value:gateway"`
	// END common fields

	BrandID         string `gorm:"type:VARCHAR"`
	ModelID         string `gorm:"type:VARCHAR"`
	HardwareVersion string `gorm:"type:VARCHAR"`
	FirmwareVersion string `gorm:"type:VARCHAR"`

	GatewayServerAddress string `gorm:"type:VARCHAR"`

	AutoUpdate    bool   `gorm:"not null"`
	UpdateChannel string `gorm:"type:VARCHAR"`

	// Frequency Plan IDs separated by spaces.
	FrequencyPlanID string `gorm:"type:VARCHAR"`

	StatusPublic   bool `gorm:"not null"`
	LocationPublic bool `gorm:"not null"`

	ScheduleDownlinkLate   bool  `gorm:"not null"`
	EnforceDutyCycle       bool  `gorm:"not null"`
	ScheduleAnytimeDelay   int64 `gorm:"default:0 not null"`
	DownlinkPathConstraint int

	UpdateLocationFromStatus bool `gorm:"default:false not null"`

	Antennas []GatewayAntenna

	LBSLNSSecret            []byte `gorm:"type:BYTEA;column:lbs_lns_secret"`
	ClaimAuthenticationCode []byte `gorm:"type:BYTEA;column:claim_authentication_code"`
}

func init() {
	registerModel(&Gateway{})
}

var secretFieldSeparator = []byte(":")

// functions to set fields from the gateway model into the gateway proto.
var gatewayPBSetters = map[string]func(*ttnpb.Gateway, *Gateway){
	"ids.eui":        func(pb *ttnpb.Gateway, gtw *Gateway) { pb.EUI = gtw.GatewayEUI.toPB() },
	nameField:        func(pb *ttnpb.Gateway, gtw *Gateway) { pb.Name = gtw.Name },
	descriptionField: func(pb *ttnpb.Gateway, gtw *Gateway) { pb.Description = gtw.Description },
	attributesField:  func(pb *ttnpb.Gateway, gtw *Gateway) { pb.Attributes = attributes(gtw.Attributes).toMap() },
	versionIDsField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		pb.GatewayVersionIdentifiers = ttnpb.GatewayVersionIdentifiers{
			BrandID:         gtw.BrandID,
			ModelID:         gtw.ModelID,
			HardwareVersion: gtw.HardwareVersion,
			FirmwareVersion: gtw.FirmwareVersion,
		}
	},
	brandIDField:              func(pb *ttnpb.Gateway, gtw *Gateway) { pb.BrandID = gtw.BrandID },
	modelIDField:              func(pb *ttnpb.Gateway, gtw *Gateway) { pb.ModelID = gtw.ModelID },
	hardwareVersionField:      func(pb *ttnpb.Gateway, gtw *Gateway) { pb.HardwareVersion = gtw.HardwareVersion },
	firmwareVersionField:      func(pb *ttnpb.Gateway, gtw *Gateway) { pb.FirmwareVersion = gtw.FirmwareVersion },
	gatewayServerAddressField: func(pb *ttnpb.Gateway, gtw *Gateway) { pb.GatewayServerAddress = gtw.GatewayServerAddress },
	autoUpdateField:           func(pb *ttnpb.Gateway, gtw *Gateway) { pb.AutoUpdate = gtw.AutoUpdate },
	updateChannelField:        func(pb *ttnpb.Gateway, gtw *Gateway) { pb.UpdateChannel = gtw.UpdateChannel },
	frequencyPlanIDsField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		if gtw.FrequencyPlanID == "" {
			pb.FrequencyPlanIDs = nil
		} else {
			pb.FrequencyPlanIDs = strings.Split(gtw.FrequencyPlanID, " ")
		}
	},
	statusPublicField:         func(pb *ttnpb.Gateway, gtw *Gateway) { pb.StatusPublic = gtw.StatusPublic },
	locationPublicField:       func(pb *ttnpb.Gateway, gtw *Gateway) { pb.LocationPublic = gtw.LocationPublic },
	scheduleDownlinkLateField: func(pb *ttnpb.Gateway, gtw *Gateway) { pb.ScheduleDownlinkLate = gtw.ScheduleDownlinkLate },
	scheduleAnytimeDelayField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		d := time.Duration(gtw.ScheduleAnytimeDelay)
		pb.ScheduleAnytimeDelay = &d
	},
	updateLocationFromStatusField: func(pb *ttnpb.Gateway, gtw *Gateway) { pb.UpdateLocationFromStatus = gtw.UpdateLocationFromStatus },
	enforceDutyCycleField:         func(pb *ttnpb.Gateway, gtw *Gateway) { pb.EnforceDutyCycle = gtw.EnforceDutyCycle },
	downlinkPathConstraintField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		pb.DownlinkPathConstraint = ttnpb.DownlinkPathConstraint(gtw.DownlinkPathConstraint)
	},
	antennasField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		sort.Slice(gtw.Antennas, func(i int, j int) bool { return gtw.Antennas[i].Index < gtw.Antennas[j].Index })
		pb.Antennas = make([]ttnpb.GatewayAntenna, len(gtw.Antennas))
		for i, antenna := range gtw.Antennas {
			pb.Antennas[i] = antenna.toPB()
		}
	},
	lbsLNSSecretField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		blocks := bytes.SplitN(gtw.LBSLNSSecret, secretFieldSeparator, 2)
		if len(blocks) == 2 {
			pb.LBSLNSSecret = &ttnpb.Secret{
				KeyID: string(blocks[0]),
				Value: blocks[1],
			}
		} else {
			pb.LBSLNSSecret = nil
		}
	},
	claimAuthenticationCodeField: func(pb *ttnpb.Gateway, gtw *Gateway) {
		blocks := bytes.SplitN(gtw.ClaimAuthenticationCode, secretFieldSeparator, 2)
		if len(blocks) == 2 {
			pb.ClaimAuthenticationCode = &ttnpb.Secret{
				KeyID: string(blocks[0]),
				Value: blocks[1],
			}
		} else {
			pb.ClaimAuthenticationCode = nil
		}
	},
}

// functions to set fields from the gateway proto into the gateway model.
var gatewayModelSetters = map[string]func(*Gateway, *ttnpb.Gateway){
	"ids.eui":        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.GatewayEUI = eui(pb.EUI) },
	nameField:        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.Name = pb.Name },
	descriptionField: func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.Description = pb.Description },
	attributesField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		gtw.Attributes = attributes(gtw.Attributes).updateFromMap(pb.Attributes)
	},
	versionIDsField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		gtw.BrandID = pb.BrandID
		gtw.ModelID = pb.ModelID
		gtw.HardwareVersion = pb.HardwareVersion
		gtw.FirmwareVersion = pb.FirmwareVersion
	},
	brandIDField:              func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.BrandID = pb.BrandID },
	modelIDField:              func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.ModelID = pb.ModelID },
	hardwareVersionField:      func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.HardwareVersion = pb.HardwareVersion },
	firmwareVersionField:      func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.FirmwareVersion = pb.FirmwareVersion },
	gatewayServerAddressField: func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.GatewayServerAddress = pb.GatewayServerAddress },
	autoUpdateField:           func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.AutoUpdate = pb.AutoUpdate },
	updateChannelField:        func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.UpdateChannel = pb.UpdateChannel },
	frequencyPlanIDsField:     func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.FrequencyPlanID = strings.Join(pb.FrequencyPlanIDs, " ") },
	statusPublicField:         func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.StatusPublic = pb.StatusPublic },
	locationPublicField:       func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.LocationPublic = pb.LocationPublic },
	scheduleDownlinkLateField: func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.ScheduleDownlinkLate = pb.ScheduleDownlinkLate },
	scheduleAnytimeDelayField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		if pb.ScheduleAnytimeDelay == nil {
			gtw.ScheduleAnytimeDelay = 0
		} else {
			gtw.ScheduleAnytimeDelay = int64(*pb.ScheduleAnytimeDelay)
		}
	},
	updateLocationFromStatusField: func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.UpdateLocationFromStatus = pb.UpdateLocationFromStatus },
	enforceDutyCycleField:         func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.EnforceDutyCycle = pb.EnforceDutyCycle },
	downlinkPathConstraintField:   func(gtw *Gateway, pb *ttnpb.Gateway) { gtw.DownlinkPathConstraint = int(pb.DownlinkPathConstraint) },
	antennasField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		sort.Slice(gtw.Antennas, func(i int, j int) bool { return gtw.Antennas[i].Index < gtw.Antennas[j].Index })
		antennas := make([]GatewayAntenna, len(pb.Antennas))
		copy(antennas, gtw.Antennas)
		gtw.Antennas = antennas
		for i, pb := range pb.Antennas {
			antenna := gtw.Antennas[i]
			antenna.fromPB(pb)
			antenna.Index = i
			gtw.Antennas[i] = antenna
		}
	},
	lbsLNSSecretField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		if pb.LBSLNSSecret != nil {
			var secretBuffer bytes.Buffer
			secretBuffer.WriteString(pb.LBSLNSSecret.KeyID)
			secretBuffer.Write(secretFieldSeparator)
			secretBuffer.Write(pb.LBSLNSSecret.Value)
			gtw.LBSLNSSecret = secretBuffer.Bytes()
		} else {
			gtw.LBSLNSSecret = nil
		}
	},
	claimAuthenticationCodeField: func(gtw *Gateway, pb *ttnpb.Gateway) {
		if pb.ClaimAuthenticationCode != nil {
			var secretBuffer bytes.Buffer
			secretBuffer.WriteString(pb.ClaimAuthenticationCode.KeyID)
			secretBuffer.Write(secretFieldSeparator)
			secretBuffer.Write(pb.ClaimAuthenticationCode.Value)
			gtw.ClaimAuthenticationCode = secretBuffer.Bytes()
		} else {
			gtw.ClaimAuthenticationCode = nil
		}
	},
}

// fieldMask to use if a nil or empty fieldmask is passed.
var defaultGatewayFieldMask = &pbtypes.FieldMask{}

func init() {
	paths := make([]string, 0, len(gatewayPBSetters))
	for _, path := range ttnpb.GatewayFieldPathsNested {
		if _, ok := gatewayPBSetters[path]; ok {
			paths = append(paths, path)
		}
	}
	defaultGatewayFieldMask.Paths = paths
}

// fieldmask path to column name in gateways table.
var gatewayColumnNames = map[string][]string{
	"ids.eui":                     {"gateway_eui"},
	antennasField:                 {},
	attributesField:               {},
	autoUpdateField:               {autoUpdateField},
	brandIDField:                  {"brand_id"},
	claimAuthenticationCodeField:  {claimAuthenticationCodeField},
	contactInfoField:              {},
	descriptionField:              {descriptionField},
	downlinkPathConstraintField:   {downlinkPathConstraintField},
	enforceDutyCycleField:         {enforceDutyCycleField},
	firmwareVersionField:          {"firmware_version"},
	frequencyPlanIDsField:         {"frequency_plan_id"},
	gatewayServerAddressField:     {gatewayServerAddressField},
	hardwareVersionField:          {"hardware_version"},
	locationPublicField:           {locationPublicField},
	lbsLNSSecretField:             {lbsLNSSecretField},
	modelIDField:                  {"model_id"},
	nameField:                     {nameField},
	scheduleAnytimeDelayField:     {scheduleAnytimeDelayField},
	scheduleDownlinkLateField:     {scheduleDownlinkLateField},
	statusPublicField:             {statusPublicField},
	updateChannelField:            {updateChannelField},
	updateLocationFromStatusField: {updateLocationFromStatusField},
	versionIDsField:               {"brand_id", "model_id", "hardware_version", "firmware_version"},
}

func (gtw Gateway) toPB(pb *ttnpb.Gateway, fieldMask *pbtypes.FieldMask) {
	pb.GatewayIdentifiers.GatewayID = gtw.GatewayID
	pb.GatewayIdentifiers.EUI = gtw.GatewayEUI.toPB() // Always present.
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
			if columnNames, ok := gatewayColumnNames[path]; ok {
				columns = append(columns, columnNames...)
			}
			continue
		}
	}
	return
}
