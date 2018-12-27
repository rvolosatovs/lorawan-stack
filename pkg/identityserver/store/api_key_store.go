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

// GetAPIKeyStore returns an APIKeyStore on the given db (or transaction).
func GetAPIKeyStore(db *gorm.DB) APIKeyStore {
	return &apiKeyStore{db: db}
}

type apiKeyStore struct {
	db *gorm.DB
}

func (s *apiKeyStore) CreateAPIKey(ctx context.Context, entityID *ttnpb.EntityIdentifiers, key *ttnpb.APIKey) error {
	entity, err := findEntity(ctx, s.db, entityID, "id")
	if err != nil {
		return err
	}
	model := APIKey{
		APIKeyID:   key.ID,
		Key:        key.Key,
		Rights:     Rights{Rights: key.Rights},
		Name:       key.Name,
		EntityID:   entity.PrimaryKey(),
		EntityType: entityTypeForID(entityID),
	}
	model.SetContext(ctx)
	return s.db.Create(&model).Error
}

func (s *apiKeyStore) FindAPIKeys(ctx context.Context, entityID *ttnpb.EntityIdentifiers) ([]*ttnpb.APIKey, error) {
	entity, err := findEntity(ctx, s.db, entityID, "id")
	if err != nil {
		return nil, err
	}
	var keyModels []APIKey
	err = s.db.Model(entity).Association("APIKeys").Find(&keyModels).Error
	if err != nil {
		return nil, err
	}
	keyProtos := make([]*ttnpb.APIKey, len(keyModels))
	for i, apiKey := range keyModels {
		keyProtos[i] = apiKey.toPB()
	}
	return keyProtos, nil
}

var errAPIKeyEntity = errors.DefineCorruption("api_key_entity", "API key not linked to an entity")

func (s *apiKeyStore) GetAPIKey(ctx context.Context, id string) (*ttnpb.EntityIdentifiers, *ttnpb.APIKey, error) {
	// TODO: scope by ctx
	var keyModel APIKey
	err := s.db.Where(&APIKey{APIKeyID: id}).First(&keyModel).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, errAPIKeyNotFound
		}
		return nil, nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, nil, err
	}
	k := polymorphicEntity{EntityType: keyModel.EntityType, EntityUUID: keyModel.EntityID}
	identifiers, err := identifiers(s.db, k)
	if err != nil {
		return nil, nil, err
	}
	ids, ok := identifiers[k]
	if !ok {
		return nil, nil, errAPIKeyEntity
	}
	return ids, keyModel.toPB(), nil
}

func (s *apiKeyStore) UpdateAPIKey(ctx context.Context, entityID *ttnpb.EntityIdentifiers, key *ttnpb.APIKey) (*ttnpb.APIKey, error) {
	entity, err := findEntity(ctx, s.db, entityID, "id")
	if err != nil {
		return nil, err
	}
	var keyModel APIKey
	err = s.db.Where(APIKey{
		APIKeyID:   key.ID,
		EntityID:   entity.PrimaryKey(),
		EntityType: entityTypeForID(entityID),
	}).First(&keyModel).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errAPIKeyNotFound
		}
		return nil, err
	}
	if len(key.Rights) == 0 {
		return nil, s.db.Delete(&APIKey{}).Error
	}
	keyModel.Name = key.Name
	keyModel.Rights = Rights{Rights: key.Rights}
	err = s.db.Model(&keyModel).Select("name", "rights").Updates(&keyModel).Error
	if err != nil {
		return nil, err
	}
	return keyModel.toPB(), nil
}
