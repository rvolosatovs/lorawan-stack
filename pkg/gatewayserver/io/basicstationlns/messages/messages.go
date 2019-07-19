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

package messages

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/basicstation"
	"go.thethings.network/lorawan-stack/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	errFrequencyPlan      = errors.Define("frequency_plan", "invalid frequency plan")
	errJoinRequestMessage = errors.Define("join_request_message", "invalid join-request message received")
)

// Definition of message types
const (
	// Upstream types for messages from the Gateway
	TypeUpstreamVersion              = "version"
	TypeUpstreamJoinRequest          = "jreq"
	TypeUpstreamJoinUplinkDataFrame  = "updf"
	TypeUpstreamProprietaryDataFrame = "propdf"
	TypeUpstreamTxConfirmation       = "dntxed"
	TypeUpstreamTimeSync             = "timesync"
	TypeUpstreamRemoteShell          = "rmtsh"

	// Downstream types for messages from the Network
	TypeDownstreamRouterConfig              = "router_config"
	TypeDownstreamDownlinkMessage           = "dnmsg"
	TypeDownstreamDownlinkMulticastSchedule = "dnsched"
	TypeDownstreamTimeSync                  = "timesync"
	TypeDownstreamRemoteCommand             = "runcmd"
	TypeDownstreamRemoteShell               = "rmtsh"

	configHardwareSpecPrefix            = "sx1301"
	configHardwareSpecNoOfConcentrators = "1"
)

// DataRates encodes the available datarates of the channel plan for the Station in the format below
// [0] -> SF (Spreading Factor; Range: 7...12 for LoRa, 0 for FSK)
// [1] -> BW (Bandwidth; 125/250/500 for LoRa, ignored for FSK)
// [2] -> DNONLY (Downlink Only; 1 = true, 0 = false)
type DataRates [16][3]int

// DiscoverQuery contains the unique identifier of the gateway.
// This message is sent by the gateway.
type DiscoverQuery struct {
	EUI basicstation.EUI `json:"router"`
}

// DiscoverResponse contains the response to the discover query.
// This message is sent by the Gateway Server.
type DiscoverResponse struct {
	EUI   basicstation.EUI `json:"router"`
	Muxs  basicstation.EUI `json:"muxs,omitempty"`
	URI   string           `json:"uri,omitempty"`
	Error string           `json:"error,omitempty"`
}

// Type returns the message type of the given data.
func Type(data []byte) (string, error) {
	msg := struct {
		Type string `json:"msgtype"`
	}{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return "", err
	}
	return msg.Type, nil
}

// Version contains version information.
// This message is sent by the gateway.
type Version struct {
	Station  string   `json:"station"`
	Firmware string   `json:"firmware"`
	Package  string   `json:"package"`
	Model    string   `json:"model"`
	Protocol int      `json:"protocol"`
	Features []string `json:"features,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (v Version) MarshalJSON() ([]byte, error) {
	type Alias Version
	return json.Marshal(struct {
		Type string `json:"msgtype"`
		Alias
	}{
		Type:  TypeUpstreamVersion,
		Alias: Alias(v),
	})
}

// IsProduction checks the features field for "prod" and returns true if found.
// This is then used to set debug options in the router config
func (v Version) IsProduction() bool {
	if len(v.Features) == 0 {
		return false
	}
	for _, feature := range v.Features {
		if feature == "prod" {
			return true
		}
	}
	return false
}

// SX1301Config contains the concentrator configuration.
// TODO: Harmonize this with sx1301_conf from https://github.com/TheThingsNetwork/lorawan-stack/pull/276
type SX1301Config struct{}

// RouterConfig contains the router configuration.
// This message is sent by the Gateway Server.
type RouterConfig struct {
	NetID          []int        `json:"NetID"`
	JoinEUI        [][]int      `json:"JoinEui"`
	Region         string       `json:"region"`
	HardwareSpec   string       `json:"hwspec"`
	FrequencyRange []int        `json:"freq_range"`
	DataRates      DataRates    `json:"DRs"`
	SX1301Config   SX1301Config `json:"sx1301_conf"`

	// These are debug options to be unset in production gateways
	NoCCA       bool `json:"nocca"`
	NoDutyCycle bool `json:"nodc"`
	NoDwellTime bool `json:"nodwell"`

	// TODO: Roundtrip time monitoring
	MuxTime float64 `json:"MuxTime"`
}

// GetRouterConfig returns the routerconfig message to be sent to the gateway
func GetRouterConfig(fp frequencyplans.FrequencyPlan, isProd bool) (RouterConfig, error) {
	if err := fp.Validate(); err != nil {
		return RouterConfig{}, errFrequencyPlan
	}

	cfg := RouterConfig{}

	// This disables filtering in the gateway.
	cfg.JoinEUI = nil
	cfg.NetID = nil

	band, err := band.GetByID(fp.BandID)
	if err != nil {
		return RouterConfig{}, errFrequencyPlan
	}

	s := strings.Split(band.ID, "_")
	if len(s) < 2 {
		return RouterConfig{}, errFrequencyPlan
	}
	cfg.Region = fmt.Sprintf("%s%s", s[0], s[1])
	if len(fp.Radios) == 0 {
		return RouterConfig{}, errFrequencyPlan
	}
	// TODO: Handle FP with multiple radios if necessary
	cfg.FrequencyRange = []int{
		int(fp.Radios[0].TxConfiguration.MinFrequency),
		int(fp.Radios[0].TxConfiguration.MaxFrequency),
	}

	// TODO: https://github.com/TheThingsNetwork/lorawan-stack/issues/284
	cfg.HardwareSpec = fmt.Sprintf("%s/%s", configHardwareSpecPrefix, configHardwareSpecNoOfConcentrators)

	drs, err := getDataRatesFromBandID(fp.BandID)
	if err != nil {
		return RouterConfig{}, errFrequencyPlan
	}
	cfg.DataRates = drs

	cfg.NoCCA = !isProd
	cfg.NoDutyCycle = !isProd
	cfg.NoDwellTime = !isProd

	// TODO: Harmonize this with sx1301_conf from https://github.com/TheThingsNetwork/lorawan-stack/pull/276
	cfg.SX1301Config = SX1301Config{}
	return cfg, nil
}

// UpInfo provides additional metadata on each upstream message.
type UpInfo struct {
	RxTime  int64   `json:"rxtime"`
	RCtx    int64   `json:"rtcx"`
	XTime   int64   `json:"xtime"`
	GPSTime int64   `json:"gpstime"`
	RSSI    float32 `json:"rssi"`
	SNR     float32 `json:"snr"`
}

// RadioMetaData is a the metadata that is received as part of all upstream messages (except Tx Confirmation).
type RadioMetaData struct {
	DataRate  int    `json:"DR"`
	Frequency uint64 `json:"Freq"`
	UpInfo    UpInfo `json:"upinfo"`
}

// JoinRequest is the LoRaWAN Join Request message
type JoinRequest struct {
	MHdr          uint             `json:"MHdr"`
	JoinEUI       basicstation.EUI `json:"JoinEui"`
	DevEUI        basicstation.EUI `json:"DevEui"`
	DevNonce      uint             `json:"DevNonce"`
	MIC           int32            `json:"MIC"`
	RefTime       float64          `json:"RefTime"`
	RadioMetaData RadioMetaData
}

// MarshalJSON implements json.Marshaler.
// TODO: Make MarshalJSON() messages generic.
func (req JoinRequest) MarshalJSON() ([]byte, error) {
	type Alias JoinRequest
	return json.Marshal(struct {
		Type string `json:"msgtype"`
		Alias
	}{
		Type:  TypeUpstreamJoinRequest,
		Alias: Alias(req),
	})
}

// ToUplinkMessage extracts fields from the basic station Join Request "jreq" message and converts them into an UplinkMessage for the network server.
func (req *JoinRequest) ToUplinkMessage(ids ttnpb.GatewayIdentifiers, bandID string) (ttnpb.UplinkMessage, error) {
	up := ttnpb.UplinkMessage{}
	up.ReceivedAt = time.Now()

	parsedMHDR := ttnpb.MHDR{}
	err := lorawan.UnmarshalMHDR([]byte{byte(req.MHdr)}, &parsedMHDR)
	if err != nil {
		return ttnpb.UplinkMessage{}, errJoinRequestMessage.WithCause(err)
	}

	micBytes, err := getInt32AsByteSlice(req.MIC)
	if err != nil {
		return ttnpb.UplinkMessage{}, errJoinRequestMessage.WithCause(err)
	}
	up.Payload = &ttnpb.Message{
		MHDR: parsedMHDR,
		MIC:  micBytes,
		Payload: &ttnpb.Message_JoinRequestPayload{JoinRequestPayload: &ttnpb.JoinRequestPayload{
			JoinEUI:  req.JoinEUI.EUI64,
			DevEUI:   req.DevEUI.EUI64,
			DevNonce: [2]byte{byte(req.DevNonce), byte(req.DevNonce >> 8)},
		}},
	}

	up.RawPayload, err = lorawan.MarshalMessage(*up.Payload)
	if err != nil {
		return ttnpb.UplinkMessage{}, errJoinRequestMessage.WithCause(err)
	}

	rxTime := time.Unix(req.RadioMetaData.UpInfo.RxTime, 0)
	rxMetadata := &ttnpb.RxMetadata{
		GatewayIdentifiers: ids,
		Time:               &rxTime,
		Timestamp:          uint32(req.RadioMetaData.UpInfo.XTime & 0xFFFFFFFF),
		RSSI:               req.RadioMetaData.UpInfo.RSSI,
		SNR:                req.RadioMetaData.UpInfo.SNR,
	}
	up.RxMetadata = append(up.RxMetadata, rxMetadata)

	loraDR, err := getDataRateFromIndex(bandID, req.RadioMetaData.DataRate)
	if err != nil {
		return ttnpb.UplinkMessage{}, errJoinRequestMessage.WithCause(err)
	}
	up.Settings = ttnpb.TxSettings{
		Frequency:     req.RadioMetaData.Frequency,
		DataRateIndex: (ttnpb.DataRateIndex)(req.RadioMetaData.DataRate),
		DataRate:      loraDR,
	}
	return up, nil
}
