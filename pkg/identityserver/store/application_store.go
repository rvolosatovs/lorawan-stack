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
	"reflect"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/rpcmiddleware/warning"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// GetApplicationStore returns an ApplicationStore on the given db (or transaction).
func GetApplicationStore(db *gorm.DB) ApplicationStore {
	return &applicationStore{db: db}
}

type applicationStore struct {
	db *gorm.DB
}

// selectApplicationFields selects relevant fields (based on fieldMask) and preloads details if needed.
func selectApplicationFields(ctx context.Context, query *gorm.DB, fieldMask *types.FieldMask) *gorm.DB {
	if fieldMask == nil || len(fieldMask.Paths) == 0 {
		return query.Preload("Attributes")
	}
	var applicationColumns []string
	var notFoundPaths []string
	for _, path := range fieldMask.Paths {
		switch path {
		case attributesField:
			query = query.Preload("Attributes")
		default:
			if columns, ok := applicationColumnNames[path]; ok {
				applicationColumns = append(applicationColumns, columns...)
			} else {
				notFoundPaths = append(notFoundPaths, path)
			}
		}
	}
	if len(notFoundPaths) > 0 {
		warning.Add(ctx, fmt.Sprintf("unsupported field mask paths: %s", strings.Join(notFoundPaths, ", ")))
	}
	return query.Select(cleanFields(append(append(modelColumns, "application_id"), applicationColumns...)...))
}

func (s *applicationStore) CreateApplication(ctx context.Context, app *ttnpb.Application) (*ttnpb.Application, error) {
	appModel := Application{
		ApplicationID: app.ApplicationID, // The ID is not mutated by fromPB.
	}
	appModel.fromPB(app, nil)
	appModel.SetContext(ctx)
	query := s.db.Create(&appModel)
	if query.Error != nil {
		return nil, query.Error
	}
	var appProto ttnpb.Application
	appModel.toPB(&appProto, nil)
	return &appProto, nil
}

func (s *applicationStore) FindApplications(ctx context.Context, ids []*ttnpb.ApplicationIdentifiers, fieldMask *types.FieldMask) ([]*ttnpb.Application, error) {
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.GetApplicationID()
	}
	query := s.db.Scopes(withContext(ctx), withApplicationID(idStrings...))
	query = selectApplicationFields(ctx, query, fieldMask)
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 {
		countTotal(ctx, query.Model(&Application{}))
		query = query.Limit(limit).Offset(offset)
	}
	var appModels []Application
	query = query.Find(&appModels)
	setTotal(ctx, uint64(len(appModels)))
	if query.Error != nil {
		return nil, query.Error
	}
	appProtos := make([]*ttnpb.Application, len(appModels))
	for i, appModel := range appModels {
		appProto := &ttnpb.Application{}
		appModel.toPB(appProto, fieldMask)
		appProtos[i] = appProto
	}
	return appProtos, nil
}

func (s *applicationStore) GetApplication(ctx context.Context, id *ttnpb.ApplicationIdentifiers, fieldMask *types.FieldMask) (*ttnpb.Application, error) {
	query := s.db.Scopes(withContext(ctx), withApplicationID(id.GetApplicationID()))
	query = selectApplicationFields(ctx, query, fieldMask)
	var appModel Application
	if err := query.First(&appModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errNotFoundForID(id.EntityIdentifiers())
		}
		return nil, err
	}
	appProto := &ttnpb.Application{}
	appModel.toPB(appProto, fieldMask)
	return appProto, nil
}

func (s *applicationStore) UpdateApplication(ctx context.Context, app *ttnpb.Application, fieldMask *types.FieldMask) (updated *ttnpb.Application, err error) {
	query := s.db.Scopes(withContext(ctx), withApplicationID(app.GetApplicationID()))
	query = selectApplicationFields(ctx, query, fieldMask)
	var appModel Application
	if err = query.First(&appModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errNotFoundForID(app.ApplicationIdentifiers.EntityIdentifiers())
		}
		return nil, err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return nil, err
	}
	oldAttributes := appModel.Attributes
	columns := appModel.fromPB(app, fieldMask)
	if len(columns) > 0 {
		query = s.db.Model(&appModel).Select(columns).Updates(&appModel)
		if query.Error != nil {
			return nil, query.Error
		}
	}
	if !reflect.DeepEqual(oldAttributes, appModel.Attributes) {
		if err = replaceAttributes(s.db, "application", appModel.ID, oldAttributes, appModel.Attributes); err != nil {
			return nil, err
		}
	}
	updated = &ttnpb.Application{}
	appModel.toPB(updated, fieldMask)
	return updated, nil
}

func (s *applicationStore) DeleteApplication(ctx context.Context, id *ttnpb.ApplicationIdentifiers) error {
	return deleteEntity(ctx, s.db, id.EntityIdentifiers())
}
