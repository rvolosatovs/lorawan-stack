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

package identityserver

import (
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"google.golang.org/grpc"
)

func init() {
	userID := paginationUser.UserIdentifiers

	// remove organizations assigned to the user by the populator
	for _, organization := range population.Organizations {
		for id, collaborators := range population.Memberships {
			if organization.EntityIdentifiers().IDString() == id.IDString() {
				for i, collaborator := range collaborators {
					if collaborator.EntityIdentifiers().IDString() == userID.GetUserID() {
						population.Memberships[id] = collaborators[:i+copy(collaborators[i:], collaborators[i+1:])]
					}
				}
			}
		}
	}

	// add deterministic number of organizations
	for i := 0; i < 3; i++ {
		organizationID := population.Organizations[i].EntityIdentifiers()
		ouID := paginationUser.OrganizationOrUserIdentifiers()
		population.Memberships[organizationID] = append(population.Memberships[organizationID], &ttnpb.Collaborator{
			OrganizationOrUserIdentifiers: *ouID,
			Rights:                        []ttnpb.Right{ttnpb.RIGHT_APPLICATION_ALL, ttnpb.RIGHT_CLIENT_ALL, ttnpb.RIGHT_GATEWAY_ALL, ttnpb.RIGHT_ORGANIZATION_ALL},
		})
	}
}

func TestOrganizationsNestedError(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	testWithIdentityServer(t, func(is *IdentityServer, cc *grpc.ClientConn) {
		userID := defaultUser.UserIdentifiers
		creds := userCreds(defaultUserIdx)
		org := userOrganizations(&userID).Organizations[0]

		reg := ttnpb.NewOrganizationRegistryClient(cc)

		_, err := reg.Create(ctx, &ttnpb.CreateOrganizationRequest{
			Organization: ttnpb.Organization{
				OrganizationIdentifiers: ttnpb.OrganizationIdentifiers{OrganizationID: "foo-org"},
			},
			Collaborator: *org.OrganizationOrUserIdentifiers(),
		}, creds)

		a.So(err, should.NotBeNil)
		a.So(errors.IsInvalidArgument(err), should.BeTrue)

		_, err = reg.List(ctx, &ttnpb.ListOrganizationsRequest{
			FieldMask:    types.FieldMask{Paths: []string{"name"}},
			Collaborator: org.OrganizationOrUserIdentifiers(),
		}, creds)

		a.So(err, should.NotBeNil)
		a.So(errors.IsInvalidArgument(err), should.BeTrue)
	})
}

func TestOrganizationsPermissionDenied(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	testWithIdentityServer(t, func(is *IdentityServer, cc *grpc.ClientConn) {
		reg := ttnpb.NewOrganizationRegistryClient(cc)

		_, err := reg.Create(ctx, &ttnpb.CreateOrganizationRequest{
			Organization: ttnpb.Organization{
				OrganizationIdentifiers: ttnpb.OrganizationIdentifiers{OrganizationID: "foo-org"},
			},
			Collaborator: *ttnpb.UserIdentifiers{UserID: "foo-usr"}.OrganizationOrUserIdentifiers(),
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		_, err = reg.Get(ctx, &ttnpb.GetOrganizationRequest{
			OrganizationIdentifiers: ttnpb.OrganizationIdentifiers{OrganizationID: "foo-org"},
			FieldMask:               types.FieldMask{Paths: []string{"name"}},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsUnauthenticated(err), should.BeTrue)
		}

		listRes, err := reg.List(ctx, &ttnpb.ListOrganizationsRequest{
			FieldMask: types.FieldMask{Paths: []string{"name"}},
		})

		a.So(err, should.BeNil)
		a.So(listRes.Organizations, should.BeEmpty)

		_, err = reg.List(ctx, &ttnpb.ListOrganizationsRequest{
			Collaborator: ttnpb.UserIdentifiers{UserID: "foo-usr"}.OrganizationOrUserIdentifiers(),
			FieldMask:    types.FieldMask{Paths: []string{"name"}},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		_, err = reg.Update(ctx, &ttnpb.UpdateOrganizationRequest{
			Organization: ttnpb.Organization{
				OrganizationIdentifiers: ttnpb.OrganizationIdentifiers{OrganizationID: "foo-org"},
				Name:                    "Updated Name",
			},
			FieldMask: types.FieldMask{Paths: []string{"name"}},
		})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		_, err = reg.Delete(ctx, &ttnpb.OrganizationIdentifiers{OrganizationID: "foo-org"})

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}
	})
}

func TestOrganizationsCRUD(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	testWithIdentityServer(t, func(is *IdentityServer, cc *grpc.ClientConn) {
		reg := ttnpb.NewOrganizationRegistryClient(cc)

		userID, creds := population.Users[defaultUserIdx].UserIdentifiers, userCreds(defaultUserIdx)
		credsWithoutRights := userCreds(defaultUserIdx, "key without rights")

		created, err := reg.Create(ctx, &ttnpb.CreateOrganizationRequest{
			Organization: ttnpb.Organization{
				OrganizationIdentifiers: ttnpb.OrganizationIdentifiers{OrganizationID: "foo"},
				Name:                    "Foo Organization",
			},
			Collaborator: *userID.OrganizationOrUserIdentifiers(),
		}, creds)

		a.So(err, should.BeNil)
		a.So(created.Name, should.Equal, "Foo Organization")

		got, err := reg.Get(ctx, &ttnpb.GetOrganizationRequest{
			OrganizationIdentifiers: created.OrganizationIdentifiers,
			FieldMask:               types.FieldMask{Paths: []string{"name"}},
		}, creds)

		a.So(err, should.BeNil)
		a.So(got.Name, should.Equal, created.Name)

		got, err = reg.Get(ctx, &ttnpb.GetOrganizationRequest{
			OrganizationIdentifiers: created.OrganizationIdentifiers,
			FieldMask:               types.FieldMask{Paths: []string{"ids"}},
		}, credsWithoutRights)

		a.So(err, should.BeNil)

		got, err = reg.Get(ctx, &ttnpb.GetOrganizationRequest{
			OrganizationIdentifiers: created.OrganizationIdentifiers,
			FieldMask:               types.FieldMask{Paths: []string{"attributes"}},
		}, credsWithoutRights)

		if a.So(err, should.NotBeNil) {
			a.So(errors.IsPermissionDenied(err), should.BeTrue)
		}

		updated, err := reg.Update(ctx, &ttnpb.UpdateOrganizationRequest{
			Organization: ttnpb.Organization{
				OrganizationIdentifiers: created.OrganizationIdentifiers,
				Name:                    "Updated Name",
			},
			FieldMask: types.FieldMask{Paths: []string{"name"}},
		}, creds)

		a.So(err, should.BeNil)
		a.So(updated.Name, should.Equal, "Updated Name")

		for _, collaborator := range []*ttnpb.OrganizationOrUserIdentifiers{nil, userID.OrganizationOrUserIdentifiers()} {
			list, err := reg.List(ctx, &ttnpb.ListOrganizationsRequest{
				FieldMask:    types.FieldMask{Paths: []string{"name"}},
				Collaborator: collaborator,
			}, creds)
			a.So(err, should.BeNil)
			if a.So(list.Organizations, should.NotBeEmpty) {
				var found bool
				for _, item := range list.Organizations {
					if item.OrganizationIdentifiers == created.OrganizationIdentifiers {
						found = true
						a.So(item.Name, should.Equal, updated.Name)
					}
				}
				a.So(found, should.BeTrue)
			}
		}

		_, err = reg.Delete(ctx, &created.OrganizationIdentifiers, creds)
		a.So(err, should.BeNil)
	})
}

func TestOrganizationsPagination(t *testing.T) {
	a := assertions.New(t)

	testWithIdentityServer(t, func(is *IdentityServer, cc *grpc.ClientConn) {
		userID := paginationUser.UserIdentifiers
		creds := userCreds(paginationUserIdx)

		reg := ttnpb.NewOrganizationRegistryClient(cc)

		md := rpcmetadata.MD{Limit: 2, Page: 1}
		ctx := md.ToOutgoingContext(test.Context())

		list, err := reg.List(ctx, &ttnpb.ListOrganizationsRequest{
			FieldMask:    types.FieldMask{Paths: []string{"name"}},
			Collaborator: userID.OrganizationOrUserIdentifiers(),
		}, creds)

		a.So(list, should.NotBeNil)
		a.So(list.Organizations, should.HaveLength, 2)
		a.So(err, should.BeNil)

		md = rpcmetadata.MD{Limit: 2, Page: 2}
		ctx = md.ToOutgoingContext(test.Context())

		list, err = reg.List(ctx, &ttnpb.ListOrganizationsRequest{
			FieldMask:    types.FieldMask{Paths: []string{"name"}},
			Collaborator: userID.OrganizationOrUserIdentifiers(),
		}, creds)

		a.So(list, should.NotBeNil)
		a.So(list.Organizations, should.HaveLength, 1)
		a.So(err, should.BeNil)

		md = rpcmetadata.MD{Limit: 2, Page: 3}
		ctx = md.ToOutgoingContext(test.Context())

		list, err = reg.List(ctx, &ttnpb.ListOrganizationsRequest{
			FieldMask:    types.FieldMask{Paths: []string{"name"}},
			Collaborator: userID.OrganizationOrUserIdentifiers(),
		}, creds)

		a.So(list, should.NotBeNil)
		a.So(list.Organizations, should.BeEmpty)
		a.So(err, should.BeNil)
	})
}
