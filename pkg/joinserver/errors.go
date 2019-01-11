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

import "go.thethings.network/lorawan-stack/pkg/errors"

var (
	errCheckMIC                  = errors.Define("check_mic", "MIC check failed")
	errComputeMIC                = errors.DefineInvalidArgument("compute_mic", "failed to compute MIC")
	errDecodePayload             = errors.DefineInvalidArgument("decode_payload", "failed to decode payload")
	errDevNonceTooHigh           = errors.DefineInvalidArgument("dev_nonce_too_high", "DevNonce is too high")
	errDevNonceTooSmall          = errors.DefineInvalidArgument("dev_nonce_too_small", "DevNonce is too small")
	errDuplicateIdentifiers      = errors.DefineAlreadyExists("duplicate_identifiers", "a device identified by the identifiers already exists")
	errEncodePayload             = errors.DefineInvalidArgument("encode_payload", "failed to encode payload")
	errEncryptPayload            = errors.Define("encrypt_payload", "failed to encrypt JoinAccept")
	errForwardJoinRequest        = errors.Define("forward_join_request", "failed to forward JoinRequest")
	errGenerateSessionKeyID      = errors.Define("generate_session_key_id", "failed to generate session key ID")
	errDeriveNwkSKeys            = errors.Define("derive_nwk_s_keys", "failed to derive network session keys")
	errDeriveAppSKey             = errors.Define("derive_app_s_key", "failed to derive application session key")
	errInvalidIdentifiers        = errors.DefineInvalidArgument("invalid_identifiers", "invalid identifiers")
	errJoinNonceTooHigh          = errors.Define("join_nonce_too_high", "JoinNonce is too high")
	errMICMismatch               = errors.DefineInvalidArgument("mic_mismatch", "MIC mismatch")
	errNoRootKeys                = errors.DefineCorruption("no_root_keys", "no root keys specified")
	errNoAppKey                  = errors.DefineCorruption("no_app_key", "no AppKey specified")
	errNoAppSKey                 = errors.DefineCorruption("no_app_s_key", "no AppSKey specified")
	errNoDevAddr                 = errors.DefineCorruption("no_dev_addr", "no DevAddr specified")
	errNoDevEUI                  = errors.DefineInvalidArgument("no_dev_eui", "no DevEUI specified")
	errNoFNwkSIntKey             = errors.DefineCorruption("no_f_nwk_s_int_key", "no FNwkSIntKey specified")
	errNoJoinEUI                 = errors.DefineInvalidArgument("no_join_eui", "no JoinEUI specified")
	errNoJoinRequest             = errors.DefineInvalidArgument("no_join_request", "no JoinRequest specified")
	errNoNwkKey                  = errors.DefineCorruption("no_nwk_key", "no NwkKey specified")
	errNoNwkSEncKey              = errors.DefineCorruption("no_nwk_s_enc_key", "no NwkSEncKey specified")
	errNoPayload                 = errors.DefineInvalidArgument("no_payload", "no message payload specified")
	errNoSNwkSIntKey             = errors.DefineCorruption("no_s_nwk_s_int_key", "no SNwkSIntKey specified")
	errPayloadLengthMismatch     = errors.DefineInvalidArgument("payload_length", "expected length of payload to be equal to 23 got {length}")
	errRegistryOperation         = errors.DefineInternal("registry_operation", "registry operation failed")
	errReuseDevNonce             = errors.DefineInvalidArgument("reuse_dev_nonce", "DevNonce has already been used")
	errUnknownAppEUI             = errors.Define("unknown_app_eui", "AppEUI specified is not known")
	errUnsupportedLoRaWANVersion = errors.DefineInvalidArgument("lorawan_version", "unsupported LoRaWAN version: {version}", "version")
	errWrongPayloadType          = errors.DefineInvalidArgument("payload_type", "wrong payload type: {type}")
)
