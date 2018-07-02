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

package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	errors "go.thethings.network/lorawan-stack/pkg/errorsv3"
	"go.thethings.network/lorawan-stack/pkg/random"
)

const (
	// alg is the JOSE algorithm for the Access Token and API Key.
	alg = "secret"

	// entropy is the amount of entropy we use (in bytes).
	entropy = 64
)

var (
	// enc is the encoder we use.
	enc = base64.RawURLEncoding
)

// Header is the JOSE header.
type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// Payload is the payload used to generate API keys and Access Tokens.
type Payload struct {
	Issuer string     `json:"iss,omitempty"`
	Type   APIKeyType `json:"type,omitempty"`
}

var errTokenKeySegmentation = errors.DefineInvalidArgument(
	"invalid_token_key_segmentation",
	"invalid segmentation of API Key or Access Token",
)

// DecodeTokenOrKey decodes the header and payload of a valid API Key or
// Access Token generated by this package.
func DecodeTokenOrKey(value string) (*Header, *Payload, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return nil, nil, errTokenKeySegmentation
	}

	header := new(Header)
	if err := unmarshal([]byte(parts[0]), header); err != nil {
		return nil, nil, err
	}

	payload := new(Payload)
	if err := unmarshal([]byte(parts[1]), payload); err != nil {
		return nil, nil, err
	}

	return header, payload, nil
}

func generate(typ string, payload interface{}) (string, error) {
	encHeader, err := marshal(&Header{
		Algorithm: alg,
		Type:      typ,
	})
	if err != nil {
		return "", err
	}

	encPayload, err := marshal(payload)
	if err != nil {
		return "", err
	}

	return encHeader + "." + encPayload + "." + enc.EncodeToString(random.Bytes(entropy)), nil
}

func marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return enc.EncodeToString(data), nil
}

func unmarshal(data []byte, v interface{}) error {
	js, err := enc.DecodeString(string(data))
	if err != nil {
		return err
	}

	return json.Unmarshal(js, v)
}
