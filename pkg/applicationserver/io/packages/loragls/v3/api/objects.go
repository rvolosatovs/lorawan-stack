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

package api

import (
	"context"
	"encoding/json"
	"fmt"

	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
)

// Gateway contains the description of a LoRaWAN gateway.
// https://www.loracloud.com/documentation/geolocation?url=v3.html#gateway
type Gateway struct {
	GatewayID string  `json:"gatewayId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

// Uplink contains the metadata of a LoRaWAN uplink.
// https://www.loracloud.com/documentation/geolocation?url=v3.html#uplink
type Uplink struct {
	GatewayID string
	AntennaID *uint32
	TDOA      *uint64
	RSSI      float64
	SNR       float64
}

// MarshalJSON implements json.Marshaler.
func (u *Uplink) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		u.GatewayID,
		u.AntennaID,
		u.TDOA,
		u.RSSI,
		u.SNR,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Uplink) UnmarshalJSON(b []byte) error {
	components := make([]json.RawMessage, 0, 5)
	if err := json.Unmarshal(b, &components); err != nil {
		return err
	}
	if n := len(components); n != 5 {
		return fmt.Errorf("invalid field count %d", n)
	}
	for i, c := range []interface{}{
		&u.GatewayID,
		&u.AntennaID,
		&u.TDOA,
		&u.RSSI,
		&u.SNR,
	} {
		if err := json.Unmarshal(components[i], c); err != nil {
			return err
		}
	}
	return nil
}

// Frame contains the uplink metadata for each reception.
// https://www.loracloud.com/documentation/geolocation?url=v3.html#frame
type Frame []Uplink

// SingleFrameRequest contains the location query request for a single LoRaWAN frame.
// https://www.loracloud.com/documentation/geolocation?url=v3.html#singleframe-http-request
type SingleFrameRequest struct {
	Gateways []Gateway `json:"gateways"`
	Frame    Frame     `json:"frame"`
}

func parseRxMetadata(ctx context.Context, m *ttnpb.RxMetadata) (Gateway, Uplink) {
	gtwUID := unique.ID(ctx, m.GatewayIdentifiers)
	var tdoa *uint64
	if m.FineTimestamp != 0 {
		tdoa = &m.FineTimestamp
	}
	return Gateway{
			GatewayID: gtwUID,
			Latitude:  m.Location.Latitude,
			Longitude: m.Location.Longitude,
			Altitude:  float64(m.Location.Altitude),
		}, Uplink{
			GatewayID: gtwUID,
			AntennaID: &m.AntennaIndex,
			TDOA:      tdoa,
			RSSI:      float64(m.RSSI),
			SNR:       float64(m.SNR),
		}
}

// BuildSingelFrameRequest builds a SingleFrameRequest from the provided metadata.
func BuildSingleFrameRequest(ctx context.Context, metadata []*ttnpb.RxMetadata) *SingleFrameRequest {
	r := &SingleFrameRequest{}
	for _, m := range metadata {
		if m.Location == nil {
			continue
		}
		gtw, up := parseRxMetadata(ctx, m)
		r.Gateways = append(r.Gateways, gtw)
		r.Frame = append(r.Frame, up)
	}
	return r
}

// MultiFrameRequest contains the location query request for multiple LoRaWAN frames.
// https://www.loracloud.com/documentation/geolocation?url=v3.html#multiframe-http-request
type MultiFrameRequest struct {
	Gateways []Gateway `json:"gateways"`
	Frames   []Frame   `json:"frames"`
}

// BuildMultiFrameRequest builds a MultiFrameRequest from the provided metadata.
func BuildMultiFrameRequest(ctx context.Context, mds [][]*ttnpb.RxMetadata) *MultiFrameRequest {
	r := &MultiFrameRequest{}
	gateways := map[string]struct{}{}
	for _, metadata := range mds {
		frame := Frame{}
		for _, m := range metadata {
			if m.Location == nil {
				continue
			}
			gtw, up := parseRxMetadata(ctx, m)
			if _, seen := gateways[gtw.GatewayID]; !seen {
				r.Gateways = append(r.Gateways, gtw)
				gateways[gtw.GatewayID] = struct{}{}
			}
			frame = append(frame, up)
		}
		r.Frames = append(r.Frames, frame)
	}
	return r
}

const (
	Algorithm_TDOA     = "Tdoa"
	Algorithm_RSSI     = "Rssi"
	Algorithm_RSSITDOA = "RssiTdoaCombined"
)

// Location contains the coordinates contained in a location query result.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Tolerance uint64  `json:"toleranceHoriz"`
}

// LocationSolverResult contains the result of a location query.
// https://www.loracloud.com/documentation/geolocation?url=v3.html#locationsolverresult
type LocationSolverResult struct {
	UsedGateways uint8    `json:"numUsedGateways"`
	HDOP         *float64 `json:"HDOP"`
	Algorithm    string   `json:"algorithmType"`
	Location     Location `json:"locationEst"`
}

// LocationSolverResponse contains the location query response.
// https://www.loracloud.com/documentation/geolocation?url=v3.html#singleframe-http-request
// https://www.loracloud.com/documentation/geolocation?url=v3.html#multiframe-http-request
type LocationSolverResponse struct {
	Result   *LocationSolverResult `json:"result"`
	Errors   []string              `json:"errors"`
	Warnings []string              `json:"warnings"`
}

// ExtendedLocationSolverResponse extends LocationSolverResponse with the raw JSON representation.
type ExtendedLocationSolverResponse struct {
	LocationSolverResponse

	Raw *json.RawMessage
}

// MarshalJSON implements json.Marshaler.
// Note that the Raw representation takes precedence
// in the marshaling process, if it is available.
func (r ExtendedLocationSolverResponse) MarshalJSON() ([]byte, error) {
	if r.Raw != nil {
		return r.Raw.MarshalJSON()
	}
	return json.Marshal(r.LocationSolverResponse)
}

// UnmarshalJSON implements json.Marshaler.
func (r *ExtendedLocationSolverResponse) UnmarshalJSON(b []byte) error {
	r.Raw = &json.RawMessage{}
	if err := r.Raw.UnmarshalJSON(b); err != nil {
		return err
	}
	return json.Unmarshal(b, &r.LocationSolverResponse)
}
