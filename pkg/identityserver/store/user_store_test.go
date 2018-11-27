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
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/util/test"
)

func TestUserStore(t *testing.T) {
	a := assertions.New(t)
	ctx := test.Context()

	now := time.Now()

	WithDB(t, func(t *testing.T, db *gorm.DB) {
		db.AutoMigrate(&Account{}, &User{}, &Attribute{}, &ContactInfo{})
		store := GetUserStore(db)

		created, err := store.CreateUser(ctx, &ttnpb.User{
			UserIdentifiers: ttnpb.UserIdentifiers{UserID: "foo"},
			Name:            "Foo User",
			Description:     "The Amazing Foo User",
			Attributes: map[string]string{
				"foo": "bar",
				"bar": "baz",
				"baz": "qux",
			},
			ContactInfo: []*ttnpb.ContactInfo{
				{ContactType: ttnpb.CONTACT_TYPE_TECHNICAL, ContactMethod: ttnpb.CONTACT_METHOD_EMAIL, Value: "foo@example.com", ValidatedAt: &now},
				{ContactType: ttnpb.CONTACT_TYPE_BILLING, ContactMethod: ttnpb.CONTACT_METHOD_EMAIL, Value: "admin@example.com"},
			},
		})
		a.So(err, should.BeNil)
		a.So(created.UserID, should.Equal, "foo")
		a.So(created.Name, should.Equal, "Foo User")
		a.So(created.Description, should.Equal, "The Amazing Foo User")
		a.So(created.Attributes, should.HaveLength, 3)
		a.So(created.ContactInfo, should.HaveLength, 2)
		a.So(created.CreatedAt, should.HappenAfter, time.Now().Add(-1*time.Hour))
		a.So(created.UpdatedAt, should.HappenAfter, time.Now().Add(-1*time.Hour))

		got, err := store.GetUser(ctx, &ttnpb.UserIdentifiers{UserID: "foo"}, &types.FieldMask{Paths: []string{"name", "attributes", "contact_info"}})
		a.So(err, should.BeNil)
		a.So(got.UserID, should.Equal, "foo")
		a.So(got.Name, should.Equal, "Foo User")
		a.So(got.Description, should.BeEmpty)
		a.So(got.Attributes, should.HaveLength, 3)
		a.So(got.ContactInfo, should.HaveLength, 2)
		a.So(got.CreatedAt, should.Equal, created.CreatedAt)
		a.So(got.UpdatedAt, should.Equal, created.UpdatedAt)

		_, err = store.UpdateUser(ctx, &ttnpb.User{
			UserIdentifiers: ttnpb.UserIdentifiers{UserID: "bar"},
		}, nil)
		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		updated, err := store.UpdateUser(ctx, &ttnpb.User{
			UserIdentifiers: ttnpb.UserIdentifiers{UserID: "foo"},
			Name:            "Foobar User",
			Description:     "The Amazing Foobar User",
			Attributes: map[string]string{
				"foo": "bar",
				"baz": "baz",
				"qux": "foo",
			},
			ContactInfo: []*ttnpb.ContactInfo{
				{ContactType: ttnpb.CONTACT_TYPE_TECHNICAL, ContactMethod: ttnpb.CONTACT_METHOD_EMAIL, Value: "bar@example.com"},
				{ContactType: ttnpb.CONTACT_TYPE_ABUSE, ContactMethod: ttnpb.CONTACT_METHOD_EMAIL, Value: "foo@example.com"},
				{ContactType: ttnpb.CONTACT_TYPE_BILLING, ContactMethod: ttnpb.CONTACT_METHOD_EMAIL, Value: "admin@example.com"},
				{ContactType: ttnpb.CONTACT_TYPE_BILLING, ContactMethod: ttnpb.CONTACT_METHOD_EMAIL, Value: "other_admin@example.com"},
			},
		}, &types.FieldMask{Paths: []string{"description", "attributes", "contact_info"}})
		a.So(err, should.BeNil)
		a.So(updated.Description, should.Equal, "The Amazing Foobar User")
		a.So(updated.Attributes, should.HaveLength, 3)
		a.So(updated.ContactInfo, should.HaveLength, 4)
		a.So(updated.CreatedAt, should.Equal, created.CreatedAt)
		a.So(updated.UpdatedAt, should.HappenAfter, created.CreatedAt)

		_, err = store.UpdateUser(ctx, &ttnpb.User{
			UserIdentifiers: ttnpb.UserIdentifiers{UserID: "foo"},
			Description:     "The Foobar User",
			UpdatedAt:       created.UpdatedAt,
		}, &types.FieldMask{Paths: []string{"description"}})
		a.So(err, should.NotBeNil)

		got, err = store.GetUser(ctx, &ttnpb.UserIdentifiers{UserID: "foo"}, nil)
		a.So(err, should.BeNil)
		a.So(got.UserID, should.Equal, created.UserID)
		a.So(got.Name, should.Equal, created.Name)
		a.So(got.Description, should.Equal, updated.Description)
		a.So(got.Attributes, should.Resemble, updated.Attributes)
		a.So(got.ContactInfo, should.HaveLength, len(updated.ContactInfo))
		a.So(got.CreatedAt, should.Equal, created.CreatedAt)
		a.So(got.UpdatedAt, should.Equal, updated.UpdatedAt)

		list, err := store.FindUsers(ctx, nil, &types.FieldMask{Paths: []string{"name"}})
		a.So(err, should.BeNil)
		if a.So(list, should.HaveLength, 1) {
			a.So(list[0].Name, should.EndWith, got.Name)
		}

		err = store.DeleteUser(ctx, &ttnpb.UserIdentifiers{UserID: "foo"})
		a.So(err, should.BeNil)

		got, err = store.GetUser(ctx, &ttnpb.UserIdentifiers{UserID: "foo"}, nil)
		if a.So(err, should.NotBeNil) {
			a.So(errors.IsNotFound(err), should.BeTrue)
		}

		list, err = store.FindUsers(ctx, nil, nil)
		a.So(err, should.BeNil)
		a.So(list, should.BeEmpty)

	})
}
