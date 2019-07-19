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

// Package rights implements rights fetching and checking.
package rights

import (
	"context"

	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// Rights for the request.
type Rights struct {
	ApplicationRights  map[string]*ttnpb.Rights
	ClientRights       map[string]*ttnpb.Rights
	GatewayRights      map[string]*ttnpb.Rights
	OrganizationRights map[string]*ttnpb.Rights
	UserRights         map[string]*ttnpb.Rights
}

// SetApplicationRights sets the rights for the given application.
func (r *Rights) setApplicationRights(appUID string, rights *ttnpb.Rights) {
	if r.ApplicationRights == nil {
		r.ApplicationRights = make(map[string]*ttnpb.Rights)
	}
	r.ApplicationRights[appUID] = rights
}

// SetClientRights sets the rights for the given client.
func (r *Rights) setClientRights(cliUID string, rights *ttnpb.Rights) {
	if r.ClientRights == nil {
		r.ClientRights = make(map[string]*ttnpb.Rights)
	}
	r.ClientRights[cliUID] = rights
}

// SetGatewayRights sets the rights for the given gateway.
func (r *Rights) setGatewayRights(gtwUID string, rights *ttnpb.Rights) {
	if r.GatewayRights == nil {
		r.GatewayRights = make(map[string]*ttnpb.Rights)
	}
	r.GatewayRights[gtwUID] = rights
}

// SetOrganizationRights sets the rights for the given organization.
func (r *Rights) setOrganizationRights(orgUID string, rights *ttnpb.Rights) {
	if r.OrganizationRights == nil {
		r.OrganizationRights = make(map[string]*ttnpb.Rights)
	}
	r.OrganizationRights[orgUID] = rights
}

// SetUserRights sets the rights for the given user.
func (r *Rights) setUserRights(usrUID string, rights *ttnpb.Rights) {
	if r.UserRights == nil {
		r.UserRights = make(map[string]*ttnpb.Rights)
	}
	r.UserRights[usrUID] = rights
}

// MissingApplicationRights returns the rights that are missing for the given application.
func (r Rights) MissingApplicationRights(appUID string, rights ...ttnpb.Right) []ttnpb.Right {
	return ttnpb.RightsFrom(rights...).Sub(r.ApplicationRights[appUID]).GetRights()
}

// MissingClientRights returns the rights that are missing for the given client.
func (r Rights) MissingClientRights(cliUID string, rights ...ttnpb.Right) []ttnpb.Right {
	return ttnpb.RightsFrom(rights...).Sub(r.ClientRights[cliUID]).GetRights()
}

// MissingGatewayRights returns the rights that are missing for the given gateway.
func (r Rights) MissingGatewayRights(gtwUID string, rights ...ttnpb.Right) []ttnpb.Right {
	return ttnpb.RightsFrom(rights...).Sub(r.GatewayRights[gtwUID]).GetRights()
}

// MissingOrganizationRights returns the rights that are missing for the given organization.
func (r Rights) MissingOrganizationRights(orgUID string, rights ...ttnpb.Right) []ttnpb.Right {
	return ttnpb.RightsFrom(rights...).Sub(r.OrganizationRights[orgUID]).GetRights()
}

// MissingUserRights returns the rights that are missing for the given user.
func (r Rights) MissingUserRights(usrUID string, rights ...ttnpb.Right) []ttnpb.Right {
	return ttnpb.RightsFrom(rights...).Sub(r.UserRights[usrUID]).GetRights()
}

// IncludesApplicationRights returns whether the given rights are included for the given application.
func (r Rights) IncludesApplicationRights(appUID string, rights ...ttnpb.Right) bool {
	return len(r.MissingApplicationRights(appUID, rights...)) == 0
}

// IncludesClientRights returns whether the given rights are included for the given client.
func (r Rights) IncludesClientRights(cliUID string, rights ...ttnpb.Right) bool {
	return len(r.MissingClientRights(cliUID, rights...)) == 0
}

// IncludesGatewayRights returns whether the given rights are included for the given gateway.
func (r Rights) IncludesGatewayRights(gtwUID string, rights ...ttnpb.Right) bool {
	return len(r.MissingGatewayRights(gtwUID, rights...)) == 0
}

// IncludesOrganizationRights returns whether the given rights are included for the given organization.
func (r Rights) IncludesOrganizationRights(orgUID string, rights ...ttnpb.Right) bool {
	return len(r.MissingOrganizationRights(orgUID, rights...)) == 0
}

// IncludesUserRights returns whether the given rights are included for the given user.
func (r Rights) IncludesUserRights(usrUID string, rights ...ttnpb.Right) bool {
	return len(r.MissingUserRights(usrUID, rights...)) == 0
}

type rightsKeyType struct{}

var rightsKey rightsKeyType

func fromContext(ctx context.Context) (Rights, bool) {
	if rights, ok := ctx.Value(rightsKey).(Rights); ok {
		return rights, true
	}
	return Rights{}, false
}

// NewContext returns a derived context with the given rights.
func NewContext(ctx context.Context, rights Rights) context.Context {
	return context.WithValue(ctx, rightsKey, rights)
}

type rightsCacheKeyType struct{}

var rightsCacheKey rightsCacheKeyType

// NewContextWithCache returns a derived context with a rights cache.
// This should only be used for request contexts.
func NewContextWithCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, rightsCacheKey, &Rights{})
}

func cacheInContext(ctx context.Context, f func(*Rights)) {
	if rights, ok := ctx.Value(rightsCacheKey).(*Rights); ok {
		f(rights)
	}
}

func cacheFromContext(ctx context.Context) (Rights, bool) {
	if rights, ok := ctx.Value(rightsCacheKey).(*Rights); ok {
		return *rights, true
	}
	return Rights{}, false
}
