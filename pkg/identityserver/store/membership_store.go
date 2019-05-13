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
	"context"
	"fmt"
	"runtime/trace"

	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// GetMembershipStore returns an MembershipStore on the given db (or transaction).
func GetMembershipStore(db *gorm.DB) MembershipStore {
	return &membershipStore{store: newStore(db)}
}

type membershipStore struct {
	*store
}

func (s *membershipStore) FindMembers(ctx context.Context, entityID *ttnpb.EntityIdentifiers) (map[*ttnpb.OrganizationOrUserIdentifiers]*ttnpb.Rights, error) {
	defer trace.StartRegion(ctx, fmt.Sprintf("find members of %s", entityID.EntityType())).End()
	entity, err := s.findEntity(ctx, entityID, "id")
	if err != nil {
		return nil, err
	}
	query := s.query(ctx, Membership{}).Where(&Membership{
		EntityID:   entity.PrimaryKey(),
		EntityType: entityTypeForID(entityID),
	}).Preload("Account")
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 {
		countTotal(ctx, query.Model(&Membership{}))
		query = query.Limit(limit).Offset(offset)
	}
	var memberships []Membership
	if err = query.Find(&memberships).Error; err != nil {
		return nil, err
	}
	setTotal(ctx, uint64(len(memberships)))
	membershipRights := make(map[*ttnpb.OrganizationOrUserIdentifiers]*ttnpb.Rights, len(memberships))
	for _, membership := range memberships {
		ids := membership.Account.OrganizationOrUserIdentifiers()
		rights := ttnpb.Rights(membership.Rights)
		membershipRights[ids] = &rights
	}
	return membershipRights, nil
}

func (s *membershipStore) FindMemberRights(ctx context.Context, id *ttnpb.OrganizationOrUserIdentifiers, entityType string) (map[*ttnpb.EntityIdentifiers]*ttnpb.Rights, error) {
	entityTypeForTrace := entityType
	if entityTypeForTrace == "" {
		entityTypeForTrace = "all"
	}
	defer trace.StartRegion(ctx, fmt.Sprintf("find %s memberships for %s", entityTypeForTrace, id.EntityIdentifiers().EntityType())).End()
	account, err := s.findAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, err
	}
	memberships, err := s.findAccountMemberships(account, entityType)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, err
	}
	return s.identifierRights(entityRightsForMemberships(memberships))
}

var errAccountType = errors.DefineInvalidArgument(
	"account_type",
	"account of type `{account_type}` can not collaborate on `{entity_type}`",
)

func (s *membershipStore) SetMember(ctx context.Context, id *ttnpb.OrganizationOrUserIdentifiers, entityID *ttnpb.EntityIdentifiers, rights *ttnpb.Rights) (err error) {
	defer trace.StartRegion(ctx, "update membership").End()
	account, err := s.findAccount(ctx, id)
	if err != nil {
		return err
	}
	entity, err := s.findEntity(ctx, entityID, "id")
	if err != nil {
		return err
	}
	if _, ok := entity.(*Organization); ok && account.AccountType != "user" {
		return errAccountType.WithAttributes("account_type", account.AccountType, "entity_type", "organization")
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return err
	}

	query := s.query(ctx, Membership{})
	var membership Membership
	err = query.Where(&Membership{
		AccountID:  account.PrimaryKey(),
		EntityID:   entity.PrimaryKey(),
		EntityType: entityTypeForID(entityID),
	}).First(&membership).Error
	if err == nil {
		if len(rights.Rights) == 0 {
			return query.Delete(&membership).Error
		}
		query = query.Select("rights", "updated_at")
	} else if gorm.IsRecordNotFoundError(err) {
		if len(rights.Rights) == 0 {
			return err
		}
		membership = Membership{
			AccountID:  account.PrimaryKey(),
			EntityID:   entity.PrimaryKey(),
			EntityType: entityTypeForID(entityID),
		}
		membership.SetContext(ctx)
	} else {
		return err
	}
	membership.Rights = Rights(*rights)
	return query.Save(&membership).Error
}
