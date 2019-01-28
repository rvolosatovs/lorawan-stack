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

import "go.thethings.network/lorawan-stack/pkg/ttnpb"

// APIKey model.
type APIKey struct {
	Model

	APIKeyID string `gorm:"type:VARCHAR;unique_index:api_key_id_index"`

	Key    string `gorm:"type:VARCHAR"`
	Rights Rights `gorm:"type:INT ARRAY"`
	Name   string `gorm:"type:VARCHAR"`

	EntityID   string `gorm:"type:UUID;index:api_key_entity_index;not null"`
	EntityType string `gorm:"type:VARCHAR(32);index:api_key_entity_index;not null"`
}

func init() {
	registerModel(&APIKey{})
}

func (k APIKey) toPB() *ttnpb.APIKey {
	return &ttnpb.APIKey{
		ID:     k.APIKeyID,
		Key:    k.Key,
		Name:   k.Name,
		Rights: k.Rights.Rights,
	}
}
