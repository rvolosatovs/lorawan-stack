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

// Package applicationregistry contains the implementation of an application registry service.
package applicationregistry

import (
	"time"

	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/store"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// Interface represents the interface exposed by the *Registry.
type Interface interface {
	Create(a *ttnpb.Application, fields ...string) (*Application, error)
	FindBy(a *ttnpb.Application, count uint64, f func(*Application) bool, fields ...string) error
}

var _ Interface = &Registry{}

// Registry is responsible for mapping applications to their identities.
type Registry struct {
	store store.Client
}

// New returns a new Registry with s as an internal Store.
func New(s store.Client) *Registry {
	return &Registry{
		store: s,
	}
}

// Create stores applications data in underlying store.Interface and returns a new *Application.
// It modifies CreatedAt and UpdatedAt fields of a and returns error if either of them is non-zero on a.
func (r *Registry) Create(a *ttnpb.Application, fields ...string) (*Application, error) {
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now

	if len(fields) != 0 {
		fields = append(fields, "CreatedAt", "UpdatedAt")
	}

	id, err := r.store.Create(a, fields...)
	if err != nil {
		return nil, err
	}
	return newApplication(a, r.store, id), nil
}

// FindBy searches for applications matching specified application fields in underlying store.Interface.
// The returned slice contains unique applications, matching at least one of values in a.
func (r *Registry) FindBy(a *ttnpb.Application, count uint64, f func(*Application) bool, fields ...string) error {
	if a == nil {
		return errors.New("Application specified is nil")
	}
	return r.store.FindBy(
		a,
		func() interface{} { return &ttnpb.Application{} },
		count,
		func(k store.PrimaryKey, v interface{}) bool {
			return f(newApplication(v.(*ttnpb.Application), r.store, k))
		},
		fields...,
	)
}

// FindApplicationByIdentifiers searches for applications matching specified application identifiers in r.
func FindApplicationByIdentifiers(r Interface, id *ttnpb.ApplicationIdentifiers, count uint64, f func(*Application) bool) error {
	if id == nil {
		return errors.New("Identifiers specified are nil")
	}

	fields := []string{}
	switch {
	case id.ApplicationID != "":
		fields = append(fields, "ApplicationIdentifiers.ApplicationID")
	}
	return r.FindBy(&ttnpb.Application{ApplicationIdentifiers: *id}, count, f, fields...)
}

// FindOneApplicationByIdentifiers searches for exactly one application matching specified application identifiers in r.
func FindOneApplicationByIdentifiers(r Interface, id *ttnpb.ApplicationIdentifiers) (*Application, error) {
	var app *Application
	var i uint64
	err := FindApplicationByIdentifiers(r, id, 1, func(a *Application) bool {
		i++
		if i > 1 {
			return false
		}
		app = a
		return true
	})
	if err != nil {
		return nil, err
	}
	switch i {
	case 0:
		return nil, ErrApplicationNotFound.New(nil)
	case 1:
		return app, nil
	default:
		return nil, ErrTooManyApplications.New(nil)
	}
}
