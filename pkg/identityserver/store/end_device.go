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
	"github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// EndDevice model.
type EndDevice struct {
	Model

	ApplicationID string `gorm:"unique_index:id;type:VARCHAR(36);not null;index"`
	Application   *Application

	// BEGIN common fields
	DeviceID    string      `gorm:"unique_index:id;type:VARCHAR(36);not null"`
	Name        string      `gorm:"type:VARCHAR"`
	Description string      `gorm:"type:TEXT"`
	Attributes  []Attribute `gorm:"polymorphic:Entity;polymorphic_value:device"`
	// END common fields

	BrandID         string `gorm:"type:VARCHAR"`
	ModelID         string `gorm:"type:VARCHAR"`
	HardwareVersion string `gorm:"type:VARCHAR"`
	FirmwareVersion string `gorm:"type:VARCHAR"`

	NetworkServerAddress     string `gorm:"type:VARCHAR"`
	ApplicationServerAddress string `gorm:"type:VARCHAR"`
	JoinServerAddress        string `gorm:"type:VARCHAR"`

	ServiceProfileID string `gorm:"type:VARCHAR"`

	Locations []EndDeviceLocation
}

func init() {
	registerModel(&EndDevice{})
}

func mustEndDeviceVersionIDs(pb *ttnpb.EndDevice) *ttnpb.EndDeviceVersionIdentifiers {
	if pb.VersionIDs == nil {
		pb.VersionIDs = new(ttnpb.EndDeviceVersionIdentifiers)
	}
	return pb.VersionIDs
}

// functions to set fields from the device model into the device proto.
var devicePBSetters = map[string]func(*ttnpb.EndDevice, *EndDevice){
	nameField:        func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.Name = dev.Name },
	descriptionField: func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.Description = dev.Description },
	attributesField:  func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.Attributes = attributes(dev.Attributes).toMap() },
	brandIDField: func(pb *ttnpb.EndDevice, dev *EndDevice) {
		mustEndDeviceVersionIDs(pb).BrandID = dev.BrandID
	},
	modelIDField: func(pb *ttnpb.EndDevice, dev *EndDevice) {
		mustEndDeviceVersionIDs(pb).ModelID = dev.ModelID
	},
	hardwareVersionField: func(pb *ttnpb.EndDevice, dev *EndDevice) {
		mustEndDeviceVersionIDs(pb).HardwareVersion = dev.HardwareVersion
	},
	firmwareVersionField: func(pb *ttnpb.EndDevice, dev *EndDevice) {
		mustEndDeviceVersionIDs(pb).FirmwareVersion = dev.FirmwareVersion
	},
	networkServerAddressField:     func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.NetworkServerAddress = dev.NetworkServerAddress },
	applicationServerAddressField: func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.ApplicationServerAddress = dev.ApplicationServerAddress },
	joinServerAddressField:        func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.JoinServerAddress = dev.JoinServerAddress },
	serviceProfileIDField:         func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.ServiceProfileID = dev.ServiceProfileID },
	locationsField:                func(pb *ttnpb.EndDevice, dev *EndDevice) { pb.Locations = deviceLocations(dev.Locations).toMap() },
}

// functions to set fields from the device proto into the device model.
var deviceModelSetters = map[string]func(*EndDevice, *ttnpb.EndDevice){
	nameField:        func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.Name = pb.Name },
	descriptionField: func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.Description = pb.Description },
	attributesField: func(dev *EndDevice, pb *ttnpb.EndDevice) {
		dev.Attributes = attributes(dev.Attributes).updateFromMap(pb.Attributes)
	},
	brandIDField: func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.BrandID = pb.GetVersionIDs().GetBrandID() },
	modelIDField: func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.ModelID = pb.GetVersionIDs().GetModelID() },
	hardwareVersionField: func(dev *EndDevice, pb *ttnpb.EndDevice) {
		dev.HardwareVersion = pb.GetVersionIDs().GetHardwareVersion()
	},
	firmwareVersionField: func(dev *EndDevice, pb *ttnpb.EndDevice) {
		dev.FirmwareVersion = pb.GetVersionIDs().GetFirmwareVersion()
	},
	networkServerAddressField:     func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.NetworkServerAddress = pb.NetworkServerAddress },
	applicationServerAddressField: func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.ApplicationServerAddress = pb.ApplicationServerAddress },
	joinServerAddressField:        func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.JoinServerAddress = pb.JoinServerAddress },
	serviceProfileIDField:         func(dev *EndDevice, pb *ttnpb.EndDevice) { dev.ServiceProfileID = pb.ServiceProfileID },
	locationsField: func(dev *EndDevice, pb *ttnpb.EndDevice) {
		dev.Locations = deviceLocations(dev.Locations).updateFromMap(pb.Locations)
	},
}

// fieldMask to use if a nil or empty fieldmask is passed.
var defaultEndDeviceFieldMask = &types.FieldMask{}

func init() {
	paths := make([]string, 0, len(devicePBSetters))
	for path := range devicePBSetters {
		paths = append(paths, path)
	}
	defaultEndDeviceFieldMask.Paths = paths
}

// fieldmask path to column name in devices table.
var deviceColumnNames = map[string]string{
	"ids.device_id":                      "device_id",
	"ids.application_ids.application_id": "application_id",
	attributesField:                      "",
	nameField:                            nameField,
	descriptionField:                     descriptionField,
	brandIDField:                         "brand_id",
	modelIDField:                         "model_id",
	hardwareVersionField:                 "hardware_version",
	firmwareVersionField:                 "firmware_version",
	networkServerAddressField:            networkServerAddressField,
	applicationServerAddressField:        applicationServerAddressField,
	joinServerAddressField:               joinServerAddressField,
	serviceProfileIDField:                serviceProfileIDField,
	locationsField:                       "",
}

func (dev EndDevice) toPB(pb *ttnpb.EndDevice, fieldMask *types.FieldMask) {
	pb.EndDeviceIdentifiers.ApplicationID = dev.ApplicationID
	pb.EndDeviceIdentifiers.DeviceID = dev.DeviceID
	pb.CreatedAt = cleanTime(dev.CreatedAt)
	pb.UpdatedAt = cleanTime(dev.UpdatedAt)
	if fieldMask == nil || len(fieldMask.Paths) == 0 {
		fieldMask = defaultEndDeviceFieldMask
	}
	for _, path := range fieldMask.Paths {
		if setter, ok := devicePBSetters[path]; ok {
			setter(pb, &dev)
		}
	}
}

func (dev *EndDevice) fromPB(pb *ttnpb.EndDevice, fieldMask *types.FieldMask) (columns []string) {
	if fieldMask == nil || len(fieldMask.Paths) == 0 {
		fieldMask = defaultEndDeviceFieldMask
	}
	for _, path := range fieldMask.Paths {
		if setter, ok := deviceModelSetters[path]; ok {
			setter(dev, pb)
			columnName, ok := deviceColumnNames[path]
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
