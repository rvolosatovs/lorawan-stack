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

package loracloudgeolocationv3

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
)

var (
	errFieldNotFound = errors.DefineNotFound("field_not_found", "field `{field}` not found")
	errInvalidType   = errors.DefineCorruption("invalid_type", "wrong type `{type}`")
	errInvalidValue  = errors.DefineCorruption("invalid_value", "wrong value `{value}`")
)

// QueryType enum defines the location query types of the package.
type QueryType uint8

func (t QueryType) Value() *types.Value {
	var s string
	switch t {
	case QUERY_TOARSSI:
		s = "TOARSSI"
	case QUERY_GNSS:
		s = "GNSS"
	default:
		panic("invalid query type")
	}
	return &types.Value{
		Kind: &types.Value_StringValue{
			StringValue: s,
		},
	}
}

func (t *QueryType) FromValue(v *types.Value) error {
	sv, ok := v.Kind.(*types.Value_StringValue)
	if !ok {
		return errInvalidType.WithAttributes("type", fmt.Sprintf("%T", v.Kind))
	}
	switch sv.StringValue {
	case "TOARSSI":
		*t = QUERY_TOARSSI
	case "GNSS":
		*t = QUERY_GNSS
	default:
		return errInvalidValue.WithAttributes("value", sv.StringValue)
	}
	return nil
}

const (
	// QUERY_TOARSSI uses the TOA and RSSI information from the gateway metadata to compute the location of the end device.
	QUERY_TOARSSI QueryType = iota
	// QUERY_GNSS uses the GNSS scan operations payload of the LR1110 transceiver.
	QUERY_GNSS
)

// Data contains the package configuration.
type Data struct {
	// Query is the query type used by the package.
	Query QueryType
	// MultiFrame enables multi frame requests for TOARSSI queries.
	MultiFrame bool
	// MultiFrameWindowSize represents the number of historical frames to consider for the query.
	// A window size of 0 automatically determines the number of frames based on the first byte
	// of the uplink message.
	MultiFrameWindowSize int
	// MultiFrameWindowAge limits the maximum age of the historical frames considered for the query.
	MultiFrameWindowAge time.Duration
	// ServerURL represents the remote server to which the GLS queries are sent.
	ServerURL *url.URL
	// Token is the API token to be used when comunicating with the GLS server.
	Token string
}

const (
	queryField           = "query"
	multiFrameField      = "multi_frame"
	multiFrameWindowSize = "multi_frame_window_size"
	multiFrameWindowAge  = "multi_frame_window_age"
	serverURLField       = "server_url"
	tokenField           = "token"
)

func toString(s string) *types.Value {
	return &types.Value{
		Kind: &types.Value_StringValue{
			StringValue: s,
		},
	}
}

func toBool(b bool) *types.Value {
	return &types.Value{
		Kind: &types.Value_BoolValue{
			BoolValue: b,
		},
	}
}

func toFloat64(f float64) *types.Value {
	return &types.Value{
		Kind: &types.Value_NumberValue{
			NumberValue: f,
		},
	}
}

// Struct serializes the configuration to *types.Struct.
func (d *Data) Struct() *types.Struct {
	st := &types.Struct{
		Fields: map[string]*types.Value{
			queryField: d.Query.Value(),
			tokenField: toString(d.Token),
		},
	}
	if d.ServerURL != nil {
		st.Fields[serverURLField] = toString(d.ServerURL.String())
	}
	if d.MultiFrame {
		st.Fields[multiFrameField] = toBool(d.MultiFrame)
	}
	if d.MultiFrameWindowSize > 0 {
		st.Fields[multiFrameWindowSize] = toFloat64(float64(d.MultiFrameWindowSize))
	}
	if d.MultiFrameWindowAge > 0 {
		st.Fields[multiFrameWindowAge] = toFloat64(float64(d.MultiFrameWindowAge / time.Minute))
	}
	return st
}

func stringFromValue(v *types.Value) (string, error) {
	sv, ok := v.Kind.(*types.Value_StringValue)
	if !ok {
		return "", errInvalidType.WithAttributes("type", fmt.Sprintf("%T", v.Kind))
	}
	return sv.StringValue, nil
}

func boolFromValue(v *types.Value) (bool, error) {
	bv, ok := v.Kind.(*types.Value_BoolValue)
	if !ok {
		return false, errInvalidType.WithAttributes("type", fmt.Sprintf("%T", v.Kind))
	}
	return bv.BoolValue, nil
}

func float64FromValue(v *types.Value) (float64, error) {
	fv, ok := v.Kind.(*types.Value_NumberValue)
	if !ok {
		return 0.0, errInvalidType.WithAttributes("type", fmt.Sprintf("%T", v.Kind))
	}
	return fv.NumberValue, nil
}

// FromStruct deserializes the configuration from *types.Struct.
func (d *Data) FromStruct(st *types.Struct) error {
	fields := st.GetFields()
	{
		value, ok := fields[queryField]
		if !ok {
			return errFieldNotFound.WithAttributes("field", queryField)
		}
		if err := d.Query.FromValue(value); err != nil {
			return err
		}
	}
	{
		value, ok := fields[multiFrameField]
		if ok {
			multiFrame, err := boolFromValue(value)
			if err != nil {
				return err
			}
			d.MultiFrame = multiFrame
		}
	}
	{
		value, ok := fields[multiFrameWindowSize]
		if ok {
			windowSize, err := float64FromValue(value)
			if err != nil {
				return err
			}
			d.MultiFrameWindowSize = int(windowSize)
		}
	}
	{
		value, ok := fields[multiFrameWindowAge]
		if ok {
			windowAge, err := float64FromValue(value)
			if err != nil {
				return err
			}
			d.MultiFrameWindowAge = time.Duration(windowAge) * time.Minute
		}
	}
	{
		value, ok := fields[tokenField]
		if !ok {
			return errFieldNotFound.WithAttributes("field", tokenField)
		}
		token, err := stringFromValue(value)
		if err != nil {
			return err
		}
		d.Token = token
	}
	{
		value, ok := fields[serverURLField]
		if ok {
			serverURL, err := stringFromValue(value)
			if err != nil {
				return err
			}
			u, err := url.Parse(serverURL)
			if err != nil {
				return err
			}
			d.ServerURL = u
		}

	}
	return nil
}
