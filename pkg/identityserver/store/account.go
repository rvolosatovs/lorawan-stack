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
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// Account model.
type Account struct {
	Model
	SoftDelete

	UID string `gorm:"type:VARCHAR(36);unique_index:account_uid_index"`

	AccountID   string `gorm:"type:UUID;index:account_id_index;not null"`
	AccountType string `gorm:"type:VARCHAR(32);index:account_id_index;not null"` // user or organization

	Memberships []*Membership
}

func init() {
	registerModel(&Account{})
}

// OrganizationOrUserIdentifiers for the account, depending on its type.
func (a Account) OrganizationOrUserIdentifiers() *ttnpb.OrganizationOrUserIdentifiers {
	switch a.AccountType {
	case "user":
		return ttnpb.UserIdentifiers{UserID: a.UID}.OrganizationOrUserIdentifiers()
	case "organization":
		return ttnpb.OrganizationIdentifiers{OrganizationID: a.UID}.OrganizationOrUserIdentifiers()
	default:
		panic("account is neither user nor organization")
	}
}
