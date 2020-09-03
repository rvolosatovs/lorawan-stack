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

package networkserver

import (
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
)

// ApplicationUplinkQueueConfig defines application uplink queue configuration.
type ApplicationUplinkQueueConfig struct {
	Queue      ApplicationUplinkQueue `name:"-"`
	BufferSize uint64                 `name:"buffer-size"`
}

// MACSettingConfig defines MAC-layer configuration.
type MACSettingConfig struct {
	ADRMargin                  *float32                   `name:"adr-margin" description:"The default margin Network Server should add in ADR requests if not configured in device's MAC settings"`
	DesiredRx1Delay            *ttnpb.RxDelay             `name:"desired-rx1-delay" description:"Desired Rx1Delay value Network Server should use if not configured in device's MAC settings"`
	DesiredMaxDutyCycle        *ttnpb.AggregatedDutyCycle `name:"desired-max-duty-cycle" description:"Desired MaxDutyCycle value Network Server should use if not configured in device's MAC settings"`
	DesiredADRAckLimitExponent *ttnpb.ADRAckLimitExponent `name:"desired-adr-ack-limit-exponent" description:"Desired ADR_ACK_LIMIT value Network Server should use if not configured in device's MAC settings"`
	DesiredADRAckDelayExponent *ttnpb.ADRAckDelayExponent `name:"desired-adr-ack-delay-exponent" description:"Desired ADR_ACK_DELAY value Network Server should use if not configured in device's MAC settings"`
	ClassBTimeout              *time.Duration             `name:"class-b-timeout" description:"Deadline for a device in class B mode to respond to requests from the Network Server if not configured in device's MAC settings"`
	ClassCTimeout              *time.Duration             `name:"class-c-timeout" description:"Deadline for a device in class C mode to respond to requests from the Network Server if not configured in device's MAC settings"`
	StatusTimePeriodicity      *time.Duration             `name:"status-time-periodicity" description:"The interval after which a DevStatusReq MACCommand shall be sent by Network Server if not configured in device's MAC settings"`
	StatusCountPeriodicity     *uint32                    `name:"status-count-periodicity" description:"Number of uplink messages after which a DevStatusReq MACCommand shall be sent by Network Server if not configured in device's MAC settings"`
}

// Parse parses the configuration and returns ttnpb.MACSettings.
func (c MACSettingConfig) Parse() ttnpb.MACSettings {
	p := ttnpb.MACSettings{
		ClassBTimeout:         c.ClassBTimeout,
		ClassCTimeout:         c.ClassCTimeout,
		StatusTimePeriodicity: c.StatusTimePeriodicity,
	}
	if c.ADRMargin != nil {
		p.ADRMargin = &pbtypes.FloatValue{Value: *c.ADRMargin}
	}
	if c.DesiredRx1Delay != nil {
		p.DesiredRx1Delay = &ttnpb.RxDelayValue{Value: *c.DesiredRx1Delay}
	}
	if c.DesiredMaxDutyCycle != nil {
		p.DesiredMaxDutyCycle = &ttnpb.AggregatedDutyCycleValue{Value: *c.DesiredMaxDutyCycle}
	}
	if c.DesiredADRAckLimitExponent != nil {
		p.DesiredADRAckLimitExponent = &ttnpb.ADRAckLimitExponentValue{Value: *c.DesiredADRAckLimitExponent}
	}
	if c.DesiredADRAckDelayExponent != nil {
		p.DesiredADRAckDelayExponent = &ttnpb.ADRAckDelayExponentValue{Value: *c.DesiredADRAckDelayExponent}
	}
	if c.StatusCountPeriodicity != nil {
		p.StatusCountPeriodicity = &pbtypes.UInt32Value{Value: *c.StatusCountPeriodicity}
	}
	return p
}

// DownlinkPriorityConfig defines priorities for downlink messages.
type DownlinkPriorityConfig struct {
	// JoinAccept is the downlink priority for join-accept messages.
	JoinAccept string `name:"join-accept" description:"Priority for join-accept messages (lowest, low, below_normal, normal, above_normal, high, highest)"`
	// MACCommands is the downlink priority for downlink messages with MAC commands as FRMPayload (FPort = 0) or as FOpts.
	// If the MAC commands are carried in FOpts, the highest priority of this value and the concerning application
	// downlink message's priority is used.
	MACCommands string `name:"mac-commands" description:"Priority for messages carrying MAC commands (lowest, low, below_normal, normal, above_normal, high, highest)"`
	// MaxApplicationDownlink is the highest priority permitted by the Network Server for application downlink.
	MaxApplicationDownlink string `name:"max-application-downlink" description:"Maximum priority for application downlink messages (lowest, low, below_normal, normal, above_normal, high, highest)"`
}

var downlinkPriorityConfigTable = map[string]ttnpb.TxSchedulePriority{
	"":             ttnpb.TxSchedulePriority_NORMAL,
	"lowest":       ttnpb.TxSchedulePriority_LOWEST,
	"low":          ttnpb.TxSchedulePriority_LOW,
	"below_normal": ttnpb.TxSchedulePriority_BELOW_NORMAL,
	"normal":       ttnpb.TxSchedulePriority_NORMAL,
	"above_normal": ttnpb.TxSchedulePriority_ABOVE_NORMAL,
	"high":         ttnpb.TxSchedulePriority_HIGH,
	"highest":      ttnpb.TxSchedulePriority_HIGHEST,
}

var errDownlinkPriority = errors.DefineInvalidArgument("downlink_priority", "invalid downlink priority `{value}`")

// Parse attempts to parse the configuration and returns a DownlinkPriorities.
func (c DownlinkPriorityConfig) Parse() (DownlinkPriorities, error) {
	var p DownlinkPriorities
	var ok bool
	if p.JoinAccept, ok = downlinkPriorityConfigTable[c.JoinAccept]; !ok {
		return DownlinkPriorities{}, errDownlinkPriority.WithAttributes("value", c.JoinAccept)
	}
	if p.MACCommands, ok = downlinkPriorityConfigTable[c.MACCommands]; !ok {
		return DownlinkPriorities{}, errDownlinkPriority.WithAttributes("value", c.MACCommands)
	}
	if p.MaxApplicationDownlink, ok = downlinkPriorityConfigTable[c.MaxApplicationDownlink]; !ok {
		return DownlinkPriorities{}, errDownlinkPriority.WithAttributes("value", c.MaxApplicationDownlink)
	}
	return p, nil
}

// Config represents the NetworkServer configuration.
type Config struct {
	ApplicationUplinkQueue ApplicationUplinkQueueConfig `name:"application-uplink-queue"`
	Devices                DeviceRegistry               `name:"-"`
	DownlinkTasks          DownlinkTaskQueue            `name:"-"`
	UplinkDeduplicator     UplinkDeduplicator           `name:"-"`
	NetID                  types.NetID                  `name:"net-id" description:"NetID of this Network Server"`
	DevAddrPrefixes        []types.DevAddrPrefix        `name:"dev-addr-prefixes" description:"Device address prefixes of this Network Server"`
	DeduplicationWindow    time.Duration                `name:"deduplication-window" description:"Time window during which, duplicate messages are collected for metadata"`
	CooldownWindow         time.Duration                `name:"cooldown-window" description:"Time window starting right after deduplication window, during which, duplicate messages are discarded"`
	DownlinkPriorities     DownlinkPriorityConfig       `name:"downlink-priorities" description:"Downlink message priorities"`
	DefaultMACSettings     MACSettingConfig             `name:"default-mac-settings" description:"Default MAC settings to fallback to if not specified by device, band or frequency plan"`
	Interop                config.InteropClient         `name:"interop" description:"Interop client configuration"`
	DeviceKEKLabel         string                       `name:"device-kek-label" description:"Label of KEK used to encrypt device keys at rest"`
}

// DefaultConfig is the default Network Server configuration.
var DefaultConfig = Config{
	ApplicationUplinkQueue: ApplicationUplinkQueueConfig{
		BufferSize: 1000,
	},
	DeduplicationWindow: 200 * time.Millisecond,
	CooldownWindow:      time.Second,
	DownlinkPriorities: DownlinkPriorityConfig{
		JoinAccept:             "highest",
		MACCommands:            "highest",
		MaxApplicationDownlink: "high",
	},
	DefaultMACSettings: MACSettingConfig{
		ADRMargin:              func(v float32) *float32 { return &v }(DefaultADRMargin),
		DesiredRx1Delay:        func(v ttnpb.RxDelay) *ttnpb.RxDelay { return &v }(ttnpb.RX_DELAY_5),
		ClassBTimeout:          func(v time.Duration) *time.Duration { return &v }(DefaultClassBTimeout),
		ClassCTimeout:          func(v time.Duration) *time.Duration { return &v }(DefaultClassCTimeout),
		StatusTimePeriodicity:  func(v time.Duration) *time.Duration { return &v }(DefaultStatusTimePeriodicity),
		StatusCountPeriodicity: func(v uint32) *uint32 { return &v }(DefaultStatusCountPeriodicity),
	},
}
