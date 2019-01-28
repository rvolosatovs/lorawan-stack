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
	"time"

	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// ContactInfo model.
type ContactInfo struct {
	ID string `gorm:"type:UUID;primary_key;default:gen_random_uuid()"`

	ContactType   int    `gorm:"not null"`
	ContactMethod int    `gorm:"not null"`
	Value         string `gorm:"type:VARCHAR"`

	Public bool

	ValidatedAt *time.Time

	EntityID   string `gorm:"type:UUID;index:contact_info_entity_index;not null"`
	EntityType string `gorm:"type:VARCHAR(32);index:contact_info_entity_index;not null"`
}

func init() {
	registerModel(&ContactInfo{})
}

func (c ContactInfo) toPB() *ttnpb.ContactInfo {
	return &ttnpb.ContactInfo{
		ContactType:   ttnpb.ContactType(c.ContactType),
		ContactMethod: ttnpb.ContactMethod(c.ContactMethod),
		Value:         c.Value,
		Public:        c.Public,
		ValidatedAt:   cleanTimePtr(c.ValidatedAt),
	}
}

func (c *ContactInfo) fromPB(pb *ttnpb.ContactInfo) {
	c.ContactType = int(pb.ContactType)
	c.ContactMethod = int(pb.ContactMethod)
	c.Value = pb.Value
	c.Public = pb.Public
	if c.ValidatedAt == nil || (pb.ValidatedAt != nil && pb.ValidatedAt.After(*c.ValidatedAt)) {
		c.ValidatedAt = cleanTimePtr(pb.ValidatedAt) // Keep old ValidatedAt.
	}
}
