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

package joinserver

import (
	"context"
	"time"

	"go.thethings.network/lorawan-stack/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/interop"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

type interopHandler interface {
	HandleJoin(context.Context, *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error)
}

type interopServer struct {
	JS interopHandler
}

func (srv interopServer) JoinRequest(ctx context.Context, req *interop.JoinReq) (*interop.JoinAns, error) {
	ctx = log.NewContextWithField(ctx, "namespace", "joinserver/interop")

	var cfList *ttnpb.CFList
	if len(req.CFList) > 0 {
		cfList = new(ttnpb.CFList)
		if err := lorawan.UnmarshalCFList(req.CFList, cfList); err != nil {
			return nil, interop.ErrMalformedMessage.WithCause(err)
		}
	}
	var dlSettings ttnpb.DLSettings
	if err := lorawan.UnmarshalDLSettings(req.DLSettings, &dlSettings); err != nil {
		return nil, interop.ErrMalformedMessage.WithCause(err)
	}

	res, err := srv.JS.HandleJoin(ctx, &ttnpb.JoinRequest{
		RawPayload:         req.PHYPayload,
		DevAddr:            req.DevAddr,
		SelectedMACVersion: ttnpb.MACVersion(req.MACVersion),
		NetID:              req.SenderID,
		DownlinkSettings:   dlSettings,
		RxDelay:            req.RxDelay,
		CFList:             cfList,
	})
	if err != nil {
		switch {
		case errors.Resemble(err, errDecodePayload),
			errors.Resemble(err, errWrongPayloadType),
			errors.Resemble(err, errNoDevEUI),
			errors.Resemble(err, errNoJoinEUI):
			return nil, interop.ErrMalformedMessage.WithCause(err)
		case errors.Resemble(err, errAddressNotAuthorized):
			return nil, interop.ErrActivation.WithCause(err)
		case errors.Resemble(err, errMICMismatch):
			return nil, interop.ErrMIC.WithCause(err)
		case errors.Resemble(err, errRegistryOperation):
			if errors.IsNotFound(errors.Cause(err)) {
				return nil, interop.ErrUnknownDevEUI.WithCause(err)
			}
		}
		return nil, interop.ErrJoinReq.WithCause(err)
	}

	header, err := req.AnswerHeader()
	if err != nil {
		return nil, interop.ErrMalformedMessage.WithCause(err)
	}
	ans := &interop.JoinAns{
		JsNsMessageHeader: header,
		PHYPayload:        interop.Buffer(res.RawPayload),
		Result:            interop.ResultSuccess,
		Lifetime:          uint32(res.Lifetime / time.Second),
		AppSKey:           (*interop.KeyEnvelope)(res.AppSKey),
		SessionKeyID:      interop.Buffer(res.SessionKeyID),
	}
	if ttnpb.MACVersion(req.MACVersion).Compare(ttnpb.MAC_V1_1) < 0 {
		ans.NwkSKey = (*interop.KeyEnvelope)(res.FNwkSIntKey)
	} else {
		ans.FNwkSIntKey = (*interop.KeyEnvelope)(res.FNwkSIntKey)
		ans.SNwkSIntKey = (*interop.KeyEnvelope)(res.SNwkSIntKey)
		ans.NwkSEncKey = (*interop.KeyEnvelope)(res.NwkSEncKey)
	}
	return ans, nil
}
