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

package messages

import (
	"bytes"
	"encoding/binary"
	"math"

	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	errDataRateIndex = errors.DefineInvalidArgument("data_rate_index", "data rate index is out of range")
	errDataRate      = errors.DefineNotFound("data_rate", "data rate not found")
	errUID           = errors.DefineInvalidArgument("uid", "invalid uid `{uid}`")
)

func getInt32AsByteSlice(value int32) ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, value)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func getFCtrlAsUint(fCtrl ttnpb.FCtrl) uint {
	var ret uint
	if fCtrl.GetADR() {
		ret = ret | 0x80
	}
	if fCtrl.GetADRAckReq() {
		ret = ret | 0x40
	}
	if fCtrl.GetAck() {
		ret = ret | 0x20
	}
	if fCtrl.GetFPending() || fCtrl.GetClassB() {
		ret = ret | 0x10
	}
	return ret
}

func getDataRateFromIndex(bandID string, index int) (ttnpb.DataRate, bool, error) {
	phy, err := band.GetByID(bandID)
	if err != nil {
		return ttnpb.DataRate{}, false, errDataRateIndex.WithCause(err)
	}
	// All protobuf enums are int32-typed, so ensure it does not overflow.
	if index < 0 || index > math.MaxInt32 {
		return ttnpb.DataRate{}, false, errDataRateIndex.New()
	}
	dr, ok := phy.DataRates[ttnpb.DataRateIndex(index)]
	if !ok {
		return ttnpb.DataRate{}, false, errDataRateIndex.New()
	}

	if dr.Rate.GetLoRa() != nil {
		return dr.Rate, true, nil
	}
	return dr.Rate, false, nil
}

func getDataRateIndexFromDataRate(bandID string, dr ttnpb.DataRate) (int, error) {
	if (dr == ttnpb.DataRate{}) {
		return 0, errDataRate.New()
	}
	phy, err := band.GetByID(bandID)
	if err != nil {
		return 0, err
	}
	i, _, ok := phy.FindUplinkDataRate(dr)
	if !ok {
		return 0, errDataRate.New()
	}
	return int(i), nil
}
