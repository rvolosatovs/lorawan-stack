// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package deviceregistry

import (
	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/store"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/mohae/deepcopy"
)

// Interface represents the interface exposed by the *Registry.
type Interface interface {
	Create(ed *ttnpb.EndDevice) (*Device, error)
	FindDeviceByIdentifiers(ids ...*ttnpb.EndDeviceIdentifiers) ([]*Device, error)
}

// Device represents the device stored in the registry.
type Device struct {
	*ttnpb.EndDevice
	stored *ttnpb.EndDevice

	key   store.PrimaryKey
	store store.Client
}

func newDevice(ed *ttnpb.EndDevice, s store.Client, k store.PrimaryKey, stored *ttnpb.EndDevice) *Device {
	if stored == nil {
		stored = deepcopy.Copy(ed).(*ttnpb.EndDevice)
	}
	return &Device{
		EndDevice: ed,
		store:     s,
		key:       k,
		stored:    stored,
	}
}

// Registry is responsible for mapping devices to their identities.
type Registry struct {
	store store.Client
}

// New returns a new Registry with s as an internal Store.
func New(s store.Client) *Registry {
	return &Registry{
		store: s,
	}
}

// Create stores devices data in underlying store.Interface and returns a new *Device.
func (r *Registry) Create(ed *ttnpb.EndDevice) (*Device, error) {
	id, err := r.store.Create(ed)
	if err != nil {
		return nil, err
	}
	return newDevice(ed, r.store, id, ed), nil
}

var newEndDevice store.NewResultFunc = func() interface{} {
	return &ttnpb.EndDevice{}
}

// FindDeviceByIdentifiers searches for devices matching specified device identifiers in underlying store.Interface.
func (r *Registry) FindDeviceByIdentifiers(ids ...*ttnpb.EndDeviceIdentifiers) ([]*Device, error) {
	if len(ids) == 0 {
		return []*Device{}, nil
	}
	for i, id := range ids {
		if id == nil {
			return nil, errors.Errorf("Identifier %d is nil", i)
		}
	}

	// Find devices matching the first filter
	filtered, err := r.store.FindBy(&ttnpb.EndDevice{EndDeviceIdentifiers: *ids[0]}, newEndDevice)
	if err != nil {
		return nil, err
	}
	// Find devices matching other filters and intersect with devices already in filtered.
	// Loop exits early, if no devices are left in filtered.
	for i := 1; i < len(ids) && len(filtered) > 0; i++ {
		m, err := r.store.FindBy(&ttnpb.EndDevice{EndDeviceIdentifiers: *ids[i]}, newEndDevice)
		if err != nil {
			return nil, err
		}
		for k := range m {
			if _, ok := filtered[k]; !ok {
				delete(filtered, k)
			}
		}
	}

	devices := make([]*Device, 0, len(filtered))
	for id, ed := range filtered {
		devices = append(devices, newDevice(ed.(*ttnpb.EndDevice), r.store, id, deepcopy.Copy(ed).(*ttnpb.EndDevice)))
	}
	return devices, nil
}

func FindOneDeviceByIdentifiers(r Interface, ids ...*ttnpb.EndDeviceIdentifiers) (*Device, error) {
	devs, err := r.FindDeviceByIdentifiers(ids...)
	if err != nil {
		return nil, err
	}
	switch len(devs) {
	case 0:
		return nil, ErrDeviceNotFound.New(errors.Attributes{
			"identifiers": ids,
		})
	case 1:
		return devs[0], nil
	default:
		return nil, ErrTooManyDevices.New(errors.Attributes{
			"identifiers": ids,
		})
	}
}

// Update updates devices data in the underlying store.Interface.
func (d *Device) Update() error {
	if err := d.store.Update(d.key, d.EndDevice, d.stored); err != nil {
		return err
	}
	d.stored = deepcopy.Copy(d.EndDevice).(*ttnpb.EndDevice)
	return nil
}

// Delete removes device from the underlying store.Interface.
func (d *Device) Delete() error {
	return d.store.Delete(d.key)
}
