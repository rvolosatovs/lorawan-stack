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

package devicetemplates_test

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	. "go.thethings.network/lorawan-stack/v3/pkg/devicetemplates"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

const (
	otaaDevice = `{
		"ids": {
			"device_id": "otaa-device",
			"application_ids": {
				"application_id": "test-app"
			},
			"dev_eui": "0102030405060708",
			"join_eui": "0807060504030201"
		},
		"frequency_plan_id": "EU_863_870",
		"lorawan_version": "1.0.2",
		"lorawan_phy_version": "1.0.2-b",
		"root_keys": {
			"app_key": {
				"key": "01020304010203040102030401020304"
			}
		},
		"supports_join": true
	}`

	abpDevice = `{
		"ids": {
			"device_id": "abp-device",
			"application_ids": {
				"application_id": "test-app"
			},
			"dev_eui": "0102030405060708",
			"join_eui": "0807060504030201"
		},
		"frequency_plan_id": "US_902_928_FSB_2",
		"lorawan_version": "1.0.2",
		"lorawan_phy_version": "1.0.2-b",
		"mac_settings": {
			"rx1_delay": null
		},
		"supports_join": false,
		"session": {
			"dev_addr": "01010101"
		}
	}`

	abpDeviceWithoutSession = `{
		"ids": {
			"device_id": "abp-device-error",
			"application_ids": {
				"application_id": "test-app"
			},
			"dev_eui": "0102030405060708",
			"join_eui": "0807060504030201"
		},
		"frequency_plan_id": "US_902_928_FSB_2",
		"lorawan_version": "1.0.2",
		"lorawan_phy_version": "1.0.2-b",
		"mac_settings": {
			"rx1_delay": null
		},
		"supports_join": false
	}
	`
)

func validateTemplate(t *testing.T, tmpl *ttnpb.EndDeviceTemplate) {
	a := assertions.New(t)
	if !a.So(tmpl, should.NotBeNil) {
		t.FailNow()
	}

	var dev ttnpb.EndDevice
	a.So(dev.SetFields(&tmpl.EndDevice, tmpl.FieldMask.Paths...), should.BeNil)
	a.So(dev, should.Resemble, tmpl.EndDevice)
}

func validateTemplates(t *testing.T, templates []*ttnpb.EndDeviceTemplate, count int) {
	a := assertions.New(t)

	if !a.So(len(templates), should.Equal, count) {
		t.FailNow()
	}

	for _, template := range templates {
		validateTemplate(t, template)
	}
}

func TestTTSConverter(t *testing.T) {
	tts := GetConverter("the-things-stack")
	a := assertions.New(t)
	if !a.So(tts, should.NotBeNil) {
		t.FailNow()
	}

	for _, tc := range []struct {
		name              string
		reader            io.Reader
		validateError     func(t *testing.T, err error)
		validateResult    func(t *testing.T, templates []*ttnpb.EndDeviceTemplate, count int)
		nExpect           int
		expectedTemplates []*ttnpb.EndDeviceTemplate
	}{
		{
			name:   "InvalidJSON",
			reader: bytes.NewBufferString("invalid json"),
			validateError: func(t *testing.T, err error) {
				_, ok := err.(*json.SyntaxError)
				assertions.New(t).So(ok, should.BeTrue)
			},
		},
		{
			name:   "OneDevice",
			reader: bytes.NewBufferString(otaaDevice),
			validateError: func(t *testing.T, err error) {
				assertions.New(t).So(err, should.BeNil)
			},
			nExpect:        1,
			validateResult: validateTemplates,
		},
		{
			name:   "OneABPOneOTAA",
			reader: bytes.NewBufferString(abpDevice + "\n\n" + otaaDevice),
			validateError: func(t *testing.T, err error) {
				assertions.New(t).So(err, should.BeNil)
			},
			validateResult: validateTemplates,
			nExpect:        2,
		},
		{
			name:   "OneOKOneError",
			reader: bytes.NewBufferString(abpDevice + "\n\n" + "invalid json"),
			validateError: func(t *testing.T, err error) {
				_, ok := err.(*json.SyntaxError)
				assertions.New(t).So(ok, should.BeTrue)
			},
			validateResult: validateTemplates,
			nExpect:        1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := test.Context()
			ch := make(chan *ttnpb.EndDeviceTemplate)

			wg := sync.WaitGroup{}
			wg.Add(2)
			var err error
			templates := []*ttnpb.EndDeviceTemplate{}
			go func() {
				err = tts.Convert(ctx, tc.reader, ch)
				wg.Done()
			}()
			go func() {
				for i := 0; i < tc.nExpect; i++ {
					templates = append(templates, <-ch)
				}
				wg.Done()
			}()

			complete := make(chan struct{})
			go func() {
				defer func() {
					complete <- struct{}{}
				}()
				wg.Wait()
			}()

			select {
			case <-complete:
			case <-time.After(time.Second):
				t.Error("Timed out waiting for converter")
				t.FailNow()
			}

			tc.validateError(t, err)
			if tc.validateResult != nil {
				tc.validateResult(t, templates, tc.nExpect)
			}
		})
	}
}
