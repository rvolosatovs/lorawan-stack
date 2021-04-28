// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package bleve_test

import (
	"os"
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/devicerepository/store"
	"go.thethings.network/lorawan-stack/v3/pkg/devicerepository/store/bleve"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestBleve(t *testing.T) {
	a := assertions.New(t)
	if err := os.MkdirAll("testdata/data/lorawan-devices-index", 0755); err != nil {
		panic(err)
	}
	defer os.RemoveAll("testdata/data/lorawan-devices-index")
	c := bleve.Config{
		SearchPaths: []string{"testdata/data/lorawan-devices-index"},
	}
	err := c.Initialize(test.Context(), "../remote/testdata", true)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	s, err := c.NewStore(test.Context())
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	t.Run("GetBrands", func(t *testing.T) {
		for _, tc := range []struct {
			name           string
			request        store.GetBrandsRequest
			result         *store.GetBrandsResponse
			errorAssertion func(err error) bool
		}{
			{
				name: "BrandIDNotFound",
				request: store.GetBrandsRequest{
					BrandID: "missing-brand",
				},
				result: brandsResponse(),
			},
			{
				name: "BrandID",
				request: store.GetBrandsRequest{
					BrandID: "foo-vendor",
					Paths:   []string{"brand_id"},
					Limit:   1,
				},
				result: brandsResponse("foo-vendor"),
			},
			{
				name: "IncompleteBrandIDDoesNotMatch",
				request: store.GetBrandsRequest{
					BrandID: "vendor",
					Paths:   []string{"brand_id"},
					Limit:   1,
				},
				result: brandsResponse(),
			},
			{
				name: "IncompleteBrandIDDoesNotMatch2",
				request: store.GetBrandsRequest{
					BrandID: "foo-vendo",
					Paths:   []string{"brand_id"},
					Limit:   1,
				},
				result: brandsResponse(),
			},
			{
				name: "Paths",
				request: store.GetBrandsRequest{
					BrandID: "full-vendor",
					Paths: []string{
						"brand_id",
						"name",
						"email",
					},
					Limit: 1,
				},
				result: &store.GetBrandsResponse{
					Total: 1,
					Count: 1,
					Brands: []*ttnpb.EndDeviceBrand{{
						BrandID: "full-vendor",
						Name:    "Full Vendor",
						Email:   "mail@example.com",
					}},
				},
			},
			{
				name: "Order",
				request: store.GetBrandsRequest{
					OrderBy: "brand_id",
					Paths:   []string{"brand_id"},
				},
				result: brandsResponse("foo-vendor", "full-vendor"),
			},
			{
				name: "OrderDesc",
				request: store.GetBrandsRequest{
					OrderBy: "-brand_id",
					Paths:   []string{"brand_id"},
				},
				result: brandsResponse("full-vendor", "foo-vendor"),
			},
			{
				name: "Limit",
				request: store.GetBrandsRequest{
					Limit: 1,
					Paths: []string{"brand_id"},
				},
				result: &store.GetBrandsResponse{
					Total:  2,
					Count:  1,
					Brands: []*ttnpb.EndDeviceBrand{{BrandID: "foo-vendor"}},
				},
			},
			{
				name: "Offset",
				request: store.GetBrandsRequest{
					Limit: 1,
					Page:  2,
					Paths: []string{"brand_id"},
				},
				result: &store.GetBrandsResponse{
					Total:  2,
					Offset: 1,
					Count:  1,
					Brands: []*ttnpb.EndDeviceBrand{{BrandID: "full-vendor"}},
				},
			},
			{
				name: "SearchByDeviceName1",
				request: store.GetBrandsRequest{
					Search: "dev1",
					Paths:  []string{"brand_id"},
				},
				result: brandsResponse("foo-vendor"),
			},
			{
				name: "SearchBySensors",
				request: store.GetBrandsRequest{
					Search: "gas",
					Paths:  []string{"brand_id"},
				},
				result: brandsResponse("full-vendor"),
			},
			{
				name: "SearchBySensorsAll",
				request: store.GetBrandsRequest{
					Search:  "temperature",
					OrderBy: "brand_id",
					Paths:   []string{"brand_id"},
				},
				result: brandsResponse("foo-vendor", "full-vendor"),
			},
			{
				name: "SearchByPartNumber",
				request: store.GetBrandsRequest{
					Search: "P4RTN0",
					Paths:  []string{"brand_id"},
				},
				result: brandsResponse("foo-vendor"),
			},
			{
				name: "SearchByCompliances",
				request: store.GetBrandsRequest{
					Search: "ETSI",
					Paths:  []string{"brand_id"},
				},
				result: brandsResponse("full-vendor"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				result, err := s.GetBrands(tc.request)
				a := assertions.New(t)
				if tc.errorAssertion != nil {
					a.So(tc.errorAssertion(err), should.BeTrue)
				} else {
					a.So(err, should.BeNil)
				}
				a.So(result, should.Resemble, tc.result)
			})
		}
	})

	t.Run("GetModels", func(t *testing.T) {
		for _, tc := range []struct {
			name           string
			request        store.GetModelsRequest
			result         *store.GetModelsResponse
			errorAssertion func(err error) bool
		}{
			{
				name: "NotFoundModelID",
				request: store.GetModelsRequest{
					BrandID: "foo-vendor",
					ModelID: "missing-device",
				},
				result: modelsResponse(),
			},
			{
				name: "NotFoundBrandID",
				request: store.GetModelsRequest{
					BrandID: "missing-brand",
					ModelID: "dev1",
				},
				result: modelsResponse(),
			},
			{
				name: "IncompleteModelIDDoesNotMatch",
				request: store.GetModelsRequest{
					ModelID: "dev",
				},
				result: modelsResponse(),
			},
			{
				name: "Order",
				request: store.GetModelsRequest{
					OrderBy: "model_id",
					Paths:   []string{"model_id"},
				},
				result: modelsResponse("dev1", "dev2", "full-device"),
			},
			{
				name: "OrderDesc",
				request: store.GetModelsRequest{
					OrderBy: "-model_id",
					Paths:   []string{"model_id"},
				},
				result: modelsResponse("full-device", "dev2", "dev1"),
			},
			{
				name: "Offset",
				request: store.GetModelsRequest{
					OrderBy: "-model_id",
					Limit:   1,
					Page:    2,
					Paths:   []string{"model_id"},
				},
				result: &store.GetModelsResponse{
					Count:  1,
					Offset: 1,
					Total:  3,
					Models: []*ttnpb.EndDeviceModel{{
						ModelID: "dev2",
					}},
				},
			},
			{
				name: "BrandID",
				request: store.GetModelsRequest{
					BrandID: "foo-vendor",
					Paths:   []string{"model_id"},
				},
				result: modelsResponse("dev1", "dev2"),
			},
			{
				name: "ModelID",
				request: store.GetModelsRequest{
					BrandID: "foo-vendor",
					ModelID: "dev1",
					Paths:   []string{"model_id"},
				},
				result: modelsResponse("dev1"),
			},
			{
				name: "Paths",
				request: store.GetModelsRequest{
					BrandID: "foo-vendor",
					ModelID: "dev2",
					Paths: []string{
						"model_id",
						"brand_id",
						"name",
						"description",
						"sensors",
					},
				},
				result: &store.GetModelsResponse{
					Count: 1,
					Total: 1,
					Models: []*ttnpb.EndDeviceModel{{
						BrandID:     "foo-vendor",
						ModelID:     "dev2",
						Name:        "Device 2",
						Description: "My Description 2",
						Sensors:     []string{"temperature"},
					}},
				},
			},
			{
				name: "SearchByDeviceName",
				request: store.GetModelsRequest{
					Search: "dev1",
					Paths:  []string{"model_id"},
				},
				result: modelsResponse("dev1"),
			},
			{
				name: "SearchBySensors",
				request: store.GetModelsRequest{
					Search: "gas",
					Paths:  []string{"model_id"},
				},
				result: modelsResponse("full-device"),
			},
			{
				name: "SearchBySensorsAll",
				request: store.GetModelsRequest{
					Search:  "temperature",
					OrderBy: "model_id",
					Paths:   []string{"model_id"},
				},
				result: modelsResponse("dev2", "full-device"),
			},
			{
				name: "SearchByPartNumber",
				request: store.GetModelsRequest{
					Search: "P4RTN0",
					Paths:  []string{"model_id"},
				},
				result: modelsResponse("dev1"),
			},
			{
				name: "SearchByCompliances",
				request: store.GetModelsRequest{
					Search: "ETSI",
					Paths:  []string{"model_id"},
				},
				result: modelsResponse("full-device"),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				result, err := s.GetModels(tc.request)
				a := assertions.New(t)
				if tc.errorAssertion != nil {
					a.So(tc.errorAssertion(err), should.BeTrue)
				} else {
					a.So(err, should.BeNil)
				}
				a.So(result, should.Resemble, tc.result)
			})
		}
	})

	t.Run("GetTemplate", func(t *testing.T) {
		t.Run("Missing", func(t *testing.T) {
			a := assertions.New(t)

			for _, ids := range []ttnpb.EndDeviceVersionIdentifiers{
				{
					BrandID: "unknown-vendor",
				},
				{
					BrandID: "foo-vendor",
					ModelID: "unknown-model",
				},
				{
					BrandID:         "foo-vendor",
					ModelID:         "dev1",
					FirmwareVersion: "unknown-version",
				},
				{
					BrandID:         "foo-vendor",
					ModelID:         "dev1",
					FirmwareVersion: "1.0",
					BandID:          "unknown-band",
				},
			} {
				tmpl, err := s.GetTemplate(&ids)
				a.So(errors.IsNotFound(err), should.BeTrue)
				a.So(tmpl, should.BeNil)
			}
		})

		t.Run("Success", func(t *testing.T) {
			a := assertions.New(t)
			tmpl, err := s.GetTemplate(&ttnpb.EndDeviceVersionIdentifiers{
				BrandID:         "foo-vendor",
				ModelID:         "dev2",
				FirmwareVersion: "1.1",
				HardwareVersion: "2.0",
				BandID:          "EU_433",
			})
			a.So(err, should.BeNil)
			a.So(tmpl, should.NotBeNil)
		})
	})

	t.Run("TestGetCodecs", func(t *testing.T) {
		t.Run("Missing", func(t *testing.T) {
			a := assertions.New(t)

			for _, ids := range []ttnpb.EndDeviceVersionIdentifiers{
				{
					BrandID: "unknown-vendor",
				},
				{
					BrandID: "foo-vendor",
					ModelID: "unknown-model",
				},
				{
					BrandID:         "foo-vendor",
					ModelID:         "dev1",
					FirmwareVersion: "unknown-version",
				},
				{
					BrandID:         "foo-vendor",
					ModelID:         "dev1",
					FirmwareVersion: "1.0",
					BandID:          "unknown-band",
				},
			} {
				codec, err := s.GetDownlinkDecoder(store.GetCodecRequest{VersionIDs: &ids})
				a.So(errors.IsNotFound(err), should.BeTrue)
				a.So(codec, should.Equal, nil)
			}
		})
		for _, tc := range []struct {
			name  string
			f     func(store.GetCodecRequest) (*ttnpb.MessagePayloadFormatter, error)
			codec string
		}{
			{
				name:  "UplinkDecoder",
				f:     s.GetUplinkDecoder,
				codec: "// uplink decoder\n",
			},
			{
				name:  "DownlinkDecoder",
				f:     s.GetDownlinkDecoder,
				codec: "// downlink decoder\n",
			},
			{
				name:  "DownlinkEncoder",
				f:     s.GetDownlinkEncoder,
				codec: "// downlink encoder\n",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				a := assertions.New(t)

				versionIDs := &ttnpb.EndDeviceVersionIdentifiers{
					BrandID:         "foo-vendor",
					ModelID:         "dev2",
					FirmwareVersion: "1.1",
					BandID:          "EU_433",
				}
				codec, err := tc.f(store.GetCodecRequest{VersionIDs: versionIDs})
				a.So(err, should.BeNil)
				a.So(codec, should.Resemble, &ttnpb.MessagePayloadFormatter{
					Formatter:          ttnpb.PayloadFormatter_FORMATTER_JAVASCRIPT,
					FormatterParameter: tc.codec,
				})
			})
		}
	})

}
