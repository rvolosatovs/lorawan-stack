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

package devicetemplateconverter_test

import (
	"context"
	"io"
	"time"

	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

func mustHavePeer(ctx context.Context, c *component.Component, role ttnpb.ClusterRole) {
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		if _, err := c.GetPeer(ctx, role, nil); err == nil {
			return
		}
	}
	panic("could not connect to peer")
}

type mockConverter struct {
	ttnpb.EndDeviceTemplateFormat
	ConvertFunc func(context.Context, io.Reader, chan<- *ttnpb.EndDeviceTemplate) error
}

func (c *mockConverter) Format() *ttnpb.EndDeviceTemplateFormat {
	return &c.EndDeviceTemplateFormat
}

func (c *mockConverter) Convert(ctx context.Context, r io.Reader, ch chan<- *ttnpb.EndDeviceTemplate) error {
	if c.ConvertFunc == nil {
		panic("Convert should not be called")
	}
	return c.ConvertFunc(ctx, r, ch)
}
