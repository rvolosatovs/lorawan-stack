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
	"context"

	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// GetMembershipStore returns an MembershipStore on the given db (or transaction).
func GetMembershipStore(db *gorm.DB) MembershipStore {
	return &membershipStore{db: db}
}

type membershipStore struct {
	db *gorm.DB
}

func (s *membershipStore) FindMembers(ctx context.Context, entityID *ttnpb.EntityIdentifiers) (map[*ttnpb.OrganizationOrUserIdentifiers]*ttnpb.Rights, error) {
	entity, err := findEntity(ctx, s.db, entityID, "id")
	if err != nil {
		return nil, err
	}
	var memberships []Membership
	if err = s.db.Model(entity).Preload("Account").Association("Memberships").Find(&memberships).Error; err != nil {
		return nil, err
	}
	membershipRights := make(map[*ttnpb.OrganizationOrUserIdentifiers]*ttnpb.Rights, len(memberships))
	for _, membership := range memberships {
		ids := membership.Account.OrganizationOrUserIdentifiers()
		rights := ttnpb.Rights(membership.Rights)
		membershipRights[ids] = &rights
	}
	return membershipRights, nil
}

func (s *membershipStore) FindMemberRights(ctx context.Context, id *ttnpb.OrganizationOrUserIdentifiers, entityType string) (map[*ttnpb.EntityIdentifiers]*ttnpb.Rights, error) {
	account, err := findAccount(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, err
	}
	memberships, err := findAccountMemberships(s.db, account, entityType)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, err
	}
	return identifierRights(s.db, entityRightsForMemberships(memberships))
}

var errAccountType = errors.DefineInvalidArgument(
	"account_type",
	"account of type `{account_type}` can not collaborate on `{entity_type}`",
)

func (s *membershipStore) SetMember(ctx context.Context, id *ttnpb.OrganizationOrUserIdentifiers, entityID *ttnpb.EntityIdentifiers, rights *ttnpb.Rights) (err error) {
	account, err := findAccount(ctx, s.db, id)
	if err != nil {
		return err
	}
	entity, err := findEntity(ctx, s.db, entityID, "id")
	if err != nil {
		return err
	}
	if _, ok := entity.(*Organization); ok && account.AccountType != "user" {
		return errAccountType.WithAttributes("account_type", account.AccountType, "entity_type", "organization")
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return err
	}
	var memberships []Membership
	if err = s.db.Model(entity).Association("Memberships").Find(&memberships).Error; err != nil {
		return err
	}
	for _, membership := range memberships {
		if membership.AccountID != account.ID {
			continue
		}
		if len(rights.Rights) == 0 {
			return s.db.Delete(&membership).Error
		}
		membership.Rights = Rights(*rights)
		return s.db.Model(&membership).Select("rights").Updates(&membership).Error
	}
	return s.db.Model(entity).
		Set("gorm:association_autoupdate", false). // otherwise it does UPDATE instead of INSERT
		Association("Memberships").
		Append(&Membership{AccountID: account.ID, Rights: Rights(*rights)}).
		Error
}
