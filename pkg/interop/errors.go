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

package interop

import (
	"net/http"

	"github.com/labstack/echo"
	"go.thethings.network/lorawan-stack/pkg/errors"
)

var (
	errUnknownMACVersion  = errors.DefineInvalidArgument("unknown_mac_version", "unknown MAC version")
	errInvalidLength      = errors.DefineInvalidArgument("invalid_length", "invalid length")
	errInvalidRequestType = errors.DefineInvalidArgument("invalid_request_type", "invalid request type `{type}`")

	ErrNoAction           = defineError("no_action", ResultNoAction, "no action")
	ErrMIC                = defineError("mic", ResultMICFailed, "MIC failed")
	ErrFrameReplayed      = defineError("frame_replayed", ResultFrameReplayed, "frame replayed")
	ErrJoinReq            = defineError("join_req", ResultJoinReqFailed, "join-request failed")
	ErrNoRoamingAgreement = defineError("no_roaming_agreement", ResultNoRoamingAgreement, "no roaming agreement")
	ErrDeviceRoaming      = defineError("device_roaming", ResultDevRoamingDisallowed, "device roaming disallowed")
	ErrRoamingActivation  = defineError("roaming_activation", ResultRoamingActDisallowed, "roaming activation disallowed")
	ErrActivation         = defineError("activation", ResultActivationDisallowed, "activation disallowed")
	ErrUnknownDevEUI      = defineError("unknown_dev_eui", ResultUnknownDevEUI, "unknown DevEUI")
	ErrUnknownDevAddr     = defineError("unknown_dev_addr", ResultUnknownDevAddr, "unknown DevAddr")
	ErrUnknownSender      = defineError("unknown_sender", ResultUnknownSender, "unknown sender")
	ErrUnknownReceiver    = defineError("unknown_receiver", ResultUnknownReceiver, "unknown receiver")
	ErrDeferred           = defineError("deferred", ResultDeferred, "deferred")
	ErrTransmitFailed     = defineError("transmit_failed", ResultXmitFailed, "transmit failed")
	ErrFPort              = defineError("f_port", ResultInvalidFPort, "invalid FPort")
	ErrProtocolVersion    = defineError("protocol_version", ResultInvalidProtocolVersion, "invalid protocol version")
	ErrStaleDeviceProfile = defineError("stale_device_profile", ResultStaleDeviceProfile, "stale device profile")
	ErrMalformedMessage   = defineError("malformed_message", ResultMalformedMessage, "malformed message")
	ErrFrameSize          = defineError("frame_size", ResultFrameSizeError, "frame size error")
)

func defineError(name string, result Result, message string) ResultError {
	return ResultError{
		Definition: errors.DefineInvalidArgument(name, message),
		Result:     result,
	}
}

// ResultError is an error with LoRaWAN Backend Interfaces specified Result.
type ResultError struct {
	errors.Definition
	Result Result
}

// ErrorHandler is an echo.HTTPErrorHandler.
func ErrorHandler(err error, c echo.Context) {
	statusCode := errors.ToHTTPStatusCode(err)
	resultErr, isResultErr := err.(ResultError)
	if header, ok := c.Get(headerKey).(*RawMessageHeader); ok {
		answerHeader, err := header.AnswerHeader()
		if err != nil {
			c.NoContent(http.StatusBadRequest)
		} else if isResultErr {
			c.JSON(statusCode, ErrorMessage{
				RawMessageHeader: answerHeader,
				Result:           resultErr.Result,
			})
		} else {
			c.JSON(statusCode, ErrorMessage{
				RawMessageHeader: answerHeader,
				Result:           ResultOther,
			})
		}
	} else {
		c.NoContent(statusCode)
	}
}
