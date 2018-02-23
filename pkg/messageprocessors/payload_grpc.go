// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package messageprocessors

import (
	"context"

	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
)

// PayloadDecoderRPC implements the UplinkMessageProcessorServer using a payload decoder.
type PayloadDecoderRPC struct {
	PayloadDecoder
}

// Process implements the UplinkMessageProcessorServer interface.
func (r *PayloadDecoderRPC) Process(ctx context.Context, req *ttnpb.ProcessUplinkMessageRequest) (*ttnpb.UplinkMessage, error) {
	return r.Decode(&req.DeviceModel, req.Parameter, &req.Message)
}

// PayloadEncoderRPC implements the DownlinkMessageProcessorServer using a payload encoder.
type PayloadEncoderRPC struct {
	PayloadEncoder
}

// Process implements the DownlinkMessageProcessorServer interface.
func (r *PayloadEncoderRPC) Process(ctx context.Context, req *ttnpb.ProcessDownlinkMessageRequest) (*ttnpb.DownlinkMessage, error) {
	return r.Encode(&req.DeviceModel, req.Parameter, &req.Message)
}
