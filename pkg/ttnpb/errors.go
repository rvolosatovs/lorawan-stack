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

package ttnpb

import (
	"fmt"
	"math"

	proto "github.com/golang/protobuf/proto"
	"go.thethings.network/lorawan-stack/pkg/errors"
	errorsv3 "go.thethings.network/lorawan-stack/pkg/errorsv3"
	"go.thethings.network/lorawan-stack/pkg/gogoproto"
)

type errorDetails struct {
	*ErrorDetails
}

func (e errorDetails) Namespace() string     { return e.GetNamespace() }
func (e errorDetails) Name() string          { return e.GetName() }
func (e errorDetails) MessageFormat() string { return e.GetMessageFormat() }
func (e errorDetails) PublicAttributes() map[string]interface{} {
	attributes, err := gogoproto.Map(e.GetAttributes())
	if err != nil {
		panic(fmt.Sprintf("Failed to decode error attributes: %s", err)) // Likely a bug in gogoproto.
	} else {
		return attributes
	}
	return nil
}
func (e errorDetails) CorrelationID() string { return e.GetCorrelationID() }

func init() {
	errorsv3.ErrorDetailsToProto = func(e errorsv3.ErrorDetails) proto.Message {
		attributes, err := gogoproto.Struct(e.PublicAttributes())
		if err != nil {
			panic(fmt.Sprintf("Failed to encode error attributes: %s", err)) // Likely a bug in ttn (invalid attribute type).
		}
		return &ErrorDetails{
			Namespace:     e.Namespace(),
			Name:          e.Name(),
			MessageFormat: e.MessageFormat(),
			Attributes:    attributes,
			CorrelationID: e.CorrelationID(),
		}
	}
	errorsv3.ErrorDetailsFromProto = func(msg ...proto.Message) (details errorsv3.ErrorDetails, rest []proto.Message) {
		var detailsMsg *ErrorDetails
		for _, msg := range msg {
			switch msg := msg.(type) {
			case *ErrorDetails:
				detailsMsg = msg
			default:
				rest = append(rest, msg)
			}
		}
		details = errorDetails{detailsMsg}
		return
	}
}

var (
	// ErrEmptyUpdateMask is returned when the update mask is specified but empty.
	ErrEmptyUpdateMask = &errors.ErrDescriptor{
		MessageFormat: "update_mask must be non-empty",
		Code:          1,
		Type:          errors.InvalidArgument,
	}

	// ErrInvalidPathUpdateMask is returned when the update mask includes a wrong field path.
	ErrInvalidPathUpdateMask = &errors.ErrDescriptor{
		MessageFormat: "Invalid update_mask: `{path}` is not a valid path",
		Code:          2,
		Type:          errors.InvalidArgument,
	}

	// ErrMissingRawPayload represents error ocurring when raw message payload is missing.
	ErrMissingRawPayload = &errors.ErrDescriptor{
		MessageFormat: "Raw Message payload is missing",
		Type:          errors.InvalidArgument,
		Code:          3,
	}

	// ErrWrongPayloadType represents error ocurring when wrong payload type is received.
	ErrWrongPayloadType = &errors.ErrDescriptor{
		MessageFormat:  "Wrong payload type: `{type}`",
		Type:           errors.InvalidArgument,
		Code:           4,
		SafeAttributes: []string{"type"},
	}

	// ErrFPortTooHigh represents error ocurring when FPort provided is too high.
	ErrFPortTooHigh = &errors.ErrDescriptor{
		MessageFormat: fmt.Sprintf("FPort must be lower or equal to %d", math.MaxUint8),
		Type:          errors.InvalidArgument,
		Code:          5,
	}

	// ErrTxChIdxTooHigh represents error ocurring when TxChIdx provided is too high.
	ErrTxChIdxTooHigh = &errors.ErrDescriptor{
		MessageFormat: fmt.Sprintf("TxChIdx must be lower or equal to %d", math.MaxUint8),
		Type:          errors.InvalidArgument,
		Code:          6,
	}

	// ErrTxDRIdxTooHigh represents error ocurring when TxDRIdx provided is too high.
	ErrTxDRIdxTooHigh = &errors.ErrDescriptor{
		MessageFormat: fmt.Sprintf("TxDRIdx must be lower or equal to %d", math.MaxUint8),
		Type:          errors.InvalidArgument,
		Code:          7,
	}

	// ErrEmptyIdentifiers is returned when the XXXIdentifiers are empty.
	ErrEmptyIdentifiers = &errors.ErrDescriptor{
		MessageFormat: "Identifiers must not be empty",
		Code:          8,
		Type:          errors.InvalidArgument,
	}
)

func init() {
	ErrEmptyUpdateMask.Register()
	ErrInvalidPathUpdateMask.Register()
	ErrMissingRawPayload.Register()
	ErrWrongPayloadType.Register()
	ErrFPortTooHigh.Register()
	ErrTxChIdxTooHigh.Register()
	ErrTxDRIdxTooHigh.Register()
	ErrEmptyIdentifiers.Register()
}
