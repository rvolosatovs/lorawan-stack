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

package web

import (
	"context"
	enc "encoding/json"

	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

type json struct {
}

func init() {
	formatters["json"] = json{}
}

func (j json) Name() string { return "JSON" }

func (j json) Encode(ctx context.Context, msg *ttnpb.ApplicationUp) ([]byte, error) {
	return enc.Marshal(msg)
}

func (j json) Decode(ctx context.Context, data []byte) (*ttnpb.ApplicationDownlink, error) {
	msg := new(ttnpb.ApplicationDownlink)
	if err := enc.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
