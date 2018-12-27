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
	"github.com/jinzhu/gorm"
)

func replaceGatewayAntennas(db *gorm.DB, gatewayUUID string, old []GatewayAntenna, new []GatewayAntenna) (err error) {
	db = db.Where(&GatewayAntenna{GatewayID: gatewayUUID})
	if len(new) < len(old) {
		if err = db.Where("\"index\" >= ?", len(new)).Delete(&GatewayAntenna{}).Error; err != nil {
			return err
		}
	}
	for _, antenna := range new {
		antenna.GatewayID = gatewayUUID
		if err = db.Save(&antenna).Error; err != nil {
			return err
		}
	}
	return nil
}
