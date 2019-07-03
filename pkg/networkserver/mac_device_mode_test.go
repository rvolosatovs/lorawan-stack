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

package networkserver

import (
	"testing"

	"github.com/mohae/deepcopy"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func TestHandleDeviceModeInd(t *testing.T) {
	for _, tc := range []struct {
		Name             string
		Device, Expected *ttnpb.EndDevice
		Payload          *ttnpb.MACCommand_DeviceModeInd
		Events           []events.DefinitionDataClosure
		Error            error
	}{
		{
			Name: "nil payload",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{},
			},
			Error: errNoPayload,
		},
		{
			Name: "empty queue",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					QueuedResponses: []*ttnpb.MACCommand{},
					DeviceClass:     ttnpb.CLASS_A,
				},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					QueuedResponses: []*ttnpb.MACCommand{
						(&ttnpb.MACCommand_DeviceModeConf{
							Class: ttnpb.CLASS_C,
						}).MACCommand(),
					},
					DeviceClass: ttnpb.CLASS_C,
				},
			},
			Payload: &ttnpb.MACCommand_DeviceModeInd{
				Class: ttnpb.CLASS_C,
			},
			Events: []events.DefinitionDataClosure{
				evtReceiveDeviceModeIndication.BindData(&ttnpb.MACCommand_DeviceModeInd{
					Class: ttnpb.CLASS_C,
				}),
				evtEnqueueDeviceModeConfirmation.BindData(&ttnpb.MACCommand_DeviceModeConf{
					Class: ttnpb.CLASS_C,
				}),
			},
		},
		{
			Name: "non-empty queue",
			Device: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					QueuedResponses: []*ttnpb.MACCommand{
						{},
						{},
						{},
					},
					DeviceClass: ttnpb.CLASS_C,
				},
			},
			Expected: &ttnpb.EndDevice{
				MACState: &ttnpb.MACState{
					QueuedResponses: []*ttnpb.MACCommand{
						{},
						{},
						{},
						(&ttnpb.MACCommand_DeviceModeConf{
							Class: ttnpb.CLASS_A,
						}).MACCommand(),
					},
					DeviceClass: ttnpb.CLASS_A,
				},
			},
			Payload: &ttnpb.MACCommand_DeviceModeInd{
				Class: ttnpb.CLASS_A,
			},
			Events: []events.DefinitionDataClosure{
				evtReceiveDeviceModeIndication.BindData(&ttnpb.MACCommand_DeviceModeInd{
					Class: ttnpb.CLASS_A,
				}),
				evtEnqueueDeviceModeConfirmation.BindData(&ttnpb.MACCommand_DeviceModeConf{
					Class: ttnpb.CLASS_A,
				}),
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)

			dev := deepcopy.Copy(tc.Device).(*ttnpb.EndDevice)

			evs, err := handleDeviceModeInd(test.Context(), dev, tc.Payload)
			if tc.Error != nil && !a.So(err, should.EqualErrorOrDefinition, tc.Error) ||
				tc.Error == nil && !a.So(err, should.BeNil) {
				t.FailNow()
			}
			a.So(dev, should.Resemble, tc.Expected)
			a.So(evs, should.ResembleEventDefinitionDataClosures, tc.Events)
		})
	}
}
