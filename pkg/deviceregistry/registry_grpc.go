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

package deviceregistry

import (
	"context"

	"github.com/TheThingsNetwork/ttn/pkg/auth/rights"
	"github.com/TheThingsNetwork/ttn/pkg/component"
	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/rpcmiddleware/hooks"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	pbtypes "github.com/gogo/protobuf/types"
)

// ErrDeviceNotFound represents the ErrDescriptor of the error returned
// when the device is not found.
var ErrDeviceNotFound = &errors.ErrDescriptor{
	MessageFormat: "Device not found",
	Type:          errors.NotFound,
	Code:          1,
}

// ErrTooManyDevices represents the ErrDescriptor of the error returned
// when there are too many devices associated with the identifiers specified.
var ErrTooManyDevices = &errors.ErrDescriptor{
	MessageFormat: "Too many devices found",
	Type:          errors.Conflict,
	Code:          2,
}

// ErrCheckFailed represents the ErrDescriptor of the error returned
// when the check failed.
var ErrCheckFailed = &errors.ErrDescriptor{
	MessageFormat: "Argument check failed",
	Type:          errors.InvalidArgument,
	Code:          3,
}

// ErrPermissionDenied is returned if the rights were insufficient to perform
// this operation.
var ErrPermissionDenied = &errors.ErrDescriptor{
	MessageFormat: "Permission denied to perform this operation",
	Type:          errors.PermissionDenied,
	Code:          4,
}

// ErrNoApplicationID is returned if no application ID was passed to an
// operation that requires it.
var ErrNoApplicationID = &errors.ErrDescriptor{
	MessageFormat: "No application ID given",
	Type:          errors.InvalidArgument,
	Code:          5,
}

func init() {
	ErrDeviceNotFound.Register()
	ErrTooManyDevices.Register()
	ErrCheckFailed.Register()
	ErrPermissionDenied.Register()
	ErrNoApplicationID.Register()
}

// RegistryRPC implements the device registry gRPC service.
type RegistryRPC struct {
	Interface
	*component.Component

	checks struct {
		ListDevices  func(ctx context.Context, id *ttnpb.EndDeviceIdentifiers) error
		GetDevice    func(ctx context.Context, id *ttnpb.EndDeviceIdentifiers) error
		SetDevice    func(ctx context.Context, dev *ttnpb.EndDevice, fields ...string) error
		DeleteDevice func(ctx context.Context, id *ttnpb.EndDeviceIdentifiers) error
	}
}

// RPCOption represents RegistryRPC option
type RPCOption func(*RegistryRPC)

// WithListDevicesCheck sets a check to ListDevices method of RegistryRPC instance.
// ListDevices first executes fn and if error is returned by it,
// returns error, otherwise execution advances as usual.
func WithListDevicesCheck(fn func(context.Context, *ttnpb.EndDeviceIdentifiers) error) RPCOption {
	return func(r *RegistryRPC) { r.checks.ListDevices = fn }
}

// WithGetDeviceCheck sets a check to GetDevice method of RegistryRPC instance.
// GetDevice first executes fn and if error is returned by it,
// returns error, otherwise execution advances as usual.
func WithGetDeviceCheck(fn func(context.Context, *ttnpb.EndDeviceIdentifiers) error) RPCOption {
	return func(r *RegistryRPC) { r.checks.GetDevice = fn }
}

// WithSetDeviceCheck sets a check to SetDevice method of RegistryRPC instance.
// SetDevice first executes fn and if error is returned by it,
// returns error, otherwise execution advances as usual.
func WithSetDeviceCheck(fn func(context.Context, *ttnpb.EndDevice, ...string) error) RPCOption {
	return func(r *RegistryRPC) { r.checks.SetDevice = fn }
}

// WithDeleteDeviceCheck sets a check to DeleteDevice method of RegistryRPC instance.
// DeleteDevice first executes fn and if error is returned by it,
// returns error, otherwise execution advances as usual.
func WithDeleteDeviceCheck(fn func(context.Context, *ttnpb.EndDeviceIdentifiers) error) RPCOption {
	return func(r *RegistryRPC) { r.checks.DeleteDevice = fn }
}

// NewRPC returns a new instance of RegistryRPC
func NewRPC(c *component.Component, r Interface, opts ...RPCOption) (*RegistryRPC, error) {
	rpc := &RegistryRPC{
		Component: c,
		Interface: r,
	}

	for _, opt := range opts {
		opt(rpc)
	}

	hook, err := c.RightsHook()
	if err != nil {
		return nil, err
	}
	hooks.RegisterUnaryHook("/ttn.v3.Rpc", rights.HookName, hook.UnaryHook())

	return rpc, nil
}

type applicationIDGetter interface {
	GetApplicationID() string
}

func (r *RegistryRPC) checkAuth(ctx context.Context, appID applicationIDGetter, rightsToCheck ...ttnpb.Right) error {
	// TODO: Accept administrator authorization even if not tied to the application
	// https://github.com/TheThingsIndustries/ttn/issues/731
	if appID == nil || appID.GetApplicationID() == "" {
		return ErrNoApplicationID.New(nil)
	}

	if ad := rights.FromContext(ctx); !ttnpb.IncludesRights(ad, rightsToCheck...) {
		return ErrPermissionDenied.New(nil)
	}

	return nil
}

// ListDevices lists devices matching filter in underlying registry.
func (r *RegistryRPC) ListDevices(ctx context.Context, filter *ttnpb.EndDeviceIdentifiers) (*ttnpb.EndDevices, error) {
	if err := r.checkAuth(ctx, filter, ttnpb.RIGHT_APPLICATION_DEVICES_READ); err != nil {
		return nil, err
	}

	if r.checks.ListDevices != nil {
		if err := r.checks.ListDevices(ctx, filter); err != nil {
			if errors.GetType(err) != errors.Unknown {
				return nil, err
			}
			return nil, ErrCheckFailed.NewWithCause(nil, err)
		}
	}

	devs, err := FindDeviceByIdentifiers(r.Interface, filter)
	if err != nil {
		return nil, err
	}

	eds := make([]*ttnpb.EndDevice, len(devs))
	for i, dev := range devs {
		eds[i] = dev.EndDevice
	}
	return &ttnpb.EndDevices{EndDevices: eds}, nil
}

// GetDevice returns the device associated with id in underlying registry, if found.
func (r *RegistryRPC) GetDevice(ctx context.Context, id *ttnpb.EndDeviceIdentifiers) (*ttnpb.EndDevice, error) {
	if err := r.checkAuth(ctx, id, ttnpb.RIGHT_APPLICATION_DEVICES_READ); err != nil {
		return nil, err
	}

	if r.checks.GetDevice != nil {
		if err := r.checks.GetDevice(ctx, id); err != nil {
			if errors.GetType(err) != errors.Unknown {
				return nil, err
			}
			return nil, ErrCheckFailed.NewWithCause(nil, err)
		}
	}

	devs, err := FindDeviceByIdentifiers(r.Interface, id)
	if err != nil {
		return nil, err
	}
	switch len(devs) {
	case 0:
		return nil, ErrDeviceNotFound.New(nil)
	case 1:
		return devs[0].EndDevice, nil
	default:
		return nil, ErrTooManyDevices.New(nil)
	}
}

// SetDevice sets the device fields to match those of dev in underlying registry.
func (r *RegistryRPC) SetDevice(ctx context.Context, req *ttnpb.SetDeviceRequest) (*pbtypes.Empty, error) {
	if err := r.checkAuth(ctx, req, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE); err != nil {
		return nil, err
	}

	var fields []string
	if req.FieldMask != nil {
		fields = req.FieldMask.Paths
	}
	if r.checks.SetDevice != nil {
		if err := r.checks.SetDevice(ctx, &req.Device, fields...); err != nil {
			if errors.GetType(err) != errors.Unknown {
				return nil, err
			}
			return nil, ErrCheckFailed.NewWithCause(nil, err)
		}
	}

	devs, err := FindDeviceByIdentifiers(r.Interface, &req.Device.EndDeviceIdentifiers)
	if err != nil {
		return nil, err
	}
	switch len(devs) {
	case 0:
		_, err := r.Interface.Create(&req.Device, fields...)
		if err != nil {
			return nil, err
		}
		return ttnpb.Empty, nil
	case 1:
		dev := devs[0]
		dev.EndDevice = &req.Device
		return ttnpb.Empty, dev.Store(fields...)
	default:
		return nil, ErrTooManyDevices.New(nil)
	}
}

// DeleteDevice deletes the device associated with id from underlying registry.
func (r *RegistryRPC) DeleteDevice(ctx context.Context, id *ttnpb.EndDeviceIdentifiers) (*pbtypes.Empty, error) {
	if err := r.checkAuth(ctx, id, ttnpb.RIGHT_APPLICATION_DEVICES_WRITE); err != nil {
		return nil, err
	}

	if r.checks.DeleteDevice != nil {
		if err := r.checks.DeleteDevice(ctx, id); err != nil {
			if errors.GetType(err) != errors.Unknown {
				return nil, err
			}
			return nil, ErrCheckFailed.NewWithCause(nil, err)
		}
	}

	devs, err := FindDeviceByIdentifiers(r.Interface, id)
	if err != nil {
		return nil, err
	}
	switch len(devs) {
	case 0:
		return nil, ErrDeviceNotFound.New(nil)
	case 1:
		return ttnpb.Empty, devs[0].Delete()
	default:
		return nil, ErrTooManyDevices.New(nil)
	}
}
