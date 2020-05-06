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

package loraclouddevicemanagementv1

import (
	"fmt"

	"github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
)

type packageData struct {
	token string
}

const tokenField = "token"

func (d packageData) toStruct() *types.Struct {
	var st types.Struct
	st.Fields = make(map[string]*types.Value)
	st.Fields[tokenField] = &types.Value{
		Kind: &types.Value_StringValue{
			StringValue: d.token,
		},
	}
	return &st
}

var (
	errFieldNotFound    = errors.DefineNotFound("field_not_found", "field `{field}` not found")
	errInvalidFieldType = errors.DefineCorruption("invalid_field_type", "field `{field}` has the wrong type `{type}`")
)

func (d *packageData) fromStruct(st *types.Struct) error {
	fields := st.GetFields()
	value, ok := fields[tokenField]
	if !ok {
		return errFieldNotFound.WithAttributes("field", tokenField)
	}
	stringValue, ok := value.GetKind().(*types.Value_StringValue)
	if !ok {
		return errInvalidFieldType.WithAttributes(
			"field", tokenField,
			"type", fmt.Sprintf("%T", value),
		)
	}
	d.token = stringValue.StringValue
	return nil
}
