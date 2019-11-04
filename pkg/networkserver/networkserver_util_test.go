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
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/smartystreets/assertions"
	clusterauth "go.thethings.network/lorawan-stack/pkg/auth/cluster"
	"go.thethings.network/lorawan-stack/pkg/band"
	"go.thethings.network/lorawan-stack/pkg/cluster"
	"go.thethings.network/lorawan-stack/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/pkg/component/test"
	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/crypto"
	"go.thethings.network/lorawan-stack/pkg/events"
	"go.thethings.network/lorawan-stack/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/networkserver/redis"
	"go.thethings.network/lorawan-stack/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
)

var (
	NewMACState = newMACState
	TimePtr     = timePtr

	ErrABPJoinRequest            = errABPJoinRequest
	ErrDecodePayload             = errDecodePayload
	ErrUnsupportedLoRaWANVersion = errUnsupportedLoRaWANVersion

	EvtBeginApplicationLink    = evtBeginApplicationLink
	EvtCreateEndDevice         = evtCreateEndDevice
	EvtDeleteEndDevice         = evtDeleteEndDevice
	EvtDropJoinRequest         = evtDropJoinRequest
	EvtEndApplicationLink      = evtEndApplicationLink
	EvtEnqueueDevStatusRequest = evtEnqueueDevStatusRequest
	EvtEnqueueLinkCheckAnswer  = evtEnqueueLinkCheckAnswer
	EvtForwardDataUplink       = evtForwardDataUplink
	EvtForwardJoinRequest      = evtForwardJoinRequest
	EvtMergeMetadata           = evtMergeMetadata
	EvtReceiveLinkCheckRequest = evtReceiveLinkCheckRequest
	EvtUpdateEndDevice         = evtUpdateEndDevice

	Timeout = (1 << 10) * test.Delay
)

func init() {
	nsScheduleWindow = time.Hour // Ensure downlink tasks are added quickly
}

var timeNowMu sync.RWMutex

func SetTimeNow(f func() time.Time) func() {
	timeNowMu.Lock()
	old := timeNow
	timeNow = f
	return func() {
		timeNow = old
		timeNowMu.Unlock()
	}
}

type MockClock time.Time

// Set sets clock to time t and returns old clock time.
func (m *MockClock) Set(t time.Time) time.Time {
	old := time.Time(*m)
	*m = MockClock(t)
	return old
}

// Add adds d to clock and returns current clock time.
func (m *MockClock) Add(d time.Duration) time.Time {
	t := time.Time(*m).Add(d)
	*m = MockClock(t)
	return t
}

// Now returns current clock time.
func (m *MockClock) Now() time.Time {
	return time.Time(*m)
}

func NSScheduleWindow() time.Duration {
	return nsScheduleWindow
}

func GSScheduleWindow() time.Duration {
	return gsScheduleWindow
}

// CopyEndDevice returns a deep copy of ttnpb.EndDevice pb.
func CopyEndDevice(pb *ttnpb.EndDevice) *ttnpb.EndDevice {
	return deepcopy.Copy(pb).(*ttnpb.EndDevice)
}

// CopyEndDevices returns a deep copy of []*ttnpb.EndDevice pbs.
func CopyEndDevices(pbs ...*ttnpb.EndDevice) []*ttnpb.EndDevice {
	return deepcopy.Copy(pbs).([]*ttnpb.EndDevice)
}

// CopyUplinkMessage returns a deep copy of ttnpb.UplinkMessage pb.
func CopyUplinkMessage(pb *ttnpb.UplinkMessage) *ttnpb.UplinkMessage {
	return deepcopy.Copy(pb).(*ttnpb.UplinkMessage)
}

// CopyUplinkMessages returns a deep copy of ...*ttnpb.UplinkMessage pbs.
func CopyUplinkMessages(pbs ...*ttnpb.UplinkMessage) []*ttnpb.UplinkMessage {
	return deepcopy.Copy(pbs).([]*ttnpb.UplinkMessage)
}

// CopyMACParameters returns a deep copy of ttnpb.MACParameters pb.
func CopyMACParameters(pb *ttnpb.MACParameters) *ttnpb.MACParameters {
	return deepcopy.Copy(pb).(*ttnpb.MACParameters)
}

// CopySessionKeys returns a deep copy of ttnpb.SessionKeys pb.
func CopySessionKeys(pb *ttnpb.SessionKeys) *ttnpb.SessionKeys {
	return deepcopy.Copy(pb).(*ttnpb.SessionKeys)
}

func DurationPtr(v time.Duration) *time.Duration {
	return &v
}

func AES128KeyPtr(key types.AES128Key) *types.AES128Key {
	return &key
}

func MustEncryptUplink(key types.AES128Key, devAddr types.DevAddr, fCnt uint32, b ...byte) []byte {
	return test.Must(crypto.EncryptUplink(key, devAddr, fCnt, b)).([]byte)
}

func MustAppendLegacyUplinkMIC(fNwkSIntKey types.AES128Key, devAddr types.DevAddr, fCnt uint32, b ...byte) []byte {
	mic := test.Must(crypto.ComputeLegacyUplinkMIC(fNwkSIntKey, devAddr, fCnt, b)).([4]byte)
	return append(b, mic[:]...)
}

func MustAppendUplinkMIC(sNwkSIntKey, fNwkSIntKey types.AES128Key, confFCnt uint32, txDRIdx uint8, txChIdx uint8, addr types.DevAddr, fCnt uint32, b ...byte) []byte {
	mic := test.Must(crypto.ComputeUplinkMIC(sNwkSIntKey, fNwkSIntKey, confFCnt, txDRIdx, txChIdx, addr, fCnt, b)).([4]byte)
	return append(b, mic[:]...)
}

func MustAppendLegacyDownlinkMIC(fNwkSIntKey types.AES128Key, devAddr types.DevAddr, fCnt uint32, b ...byte) []byte {
	mic := test.Must(crypto.ComputeLegacyDownlinkMIC(fNwkSIntKey, devAddr, fCnt, b)).([4]byte)
	return append(b, mic[:]...)
}

func MakeLinkCheckAns(mds ...*ttnpb.RxMetadata) *ttnpb.MACCommand {
	maxSNR := mds[0].SNR
	for _, md := range mds {
		if md.SNR > maxSNR {
			maxSNR = md.SNR
		}
	}
	return (&ttnpb.MACCommand_LinkCheckAns{
		Margin:       uint32(maxSNR + 15),
		GatewayCount: uint32(len(mds)),
	}).MACCommand()
}

func MakeEU868Channels(chs ...*ttnpb.MACParameters_Channel) []*ttnpb.MACParameters_Channel {
	return append([]*ttnpb.MACParameters_Channel{
		{
			UplinkFrequency:   868100000,
			DownlinkFrequency: 868100000,
			MinDataRateIndex:  ttnpb.DATA_RATE_0,
			MaxDataRateIndex:  ttnpb.DATA_RATE_5,
			EnableUplink:      true,
		},
		{
			UplinkFrequency:   868300000,
			DownlinkFrequency: 868300000,
			MinDataRateIndex:  ttnpb.DATA_RATE_0,
			MaxDataRateIndex:  ttnpb.DATA_RATE_5,
			EnableUplink:      true,
		},
		{
			UplinkFrequency:   868500000,
			DownlinkFrequency: 868500000,
			MinDataRateIndex:  ttnpb.DATA_RATE_0,
			MaxDataRateIndex:  ttnpb.DATA_RATE_5,
			EnableUplink:      true,
		},
	}, chs...)
}

func MakeDefaultEU868MACState(class ttnpb.Class, ver ttnpb.MACVersion) *ttnpb.MACState {
	return &ttnpb.MACState{
		DeviceClass:         class,
		LoRaWANVersion:      ver,
		PingSlotPeriodicity: ttnpb.PING_EVERY_1S,
		CurrentParameters: ttnpb.MACParameters{
			ADRAckDelayExponent:    &ttnpb.ADRAckDelayExponentValue{Value: ttnpb.ADR_ACK_DELAY_32},
			ADRAckLimitExponent:    &ttnpb.ADRAckLimitExponentValue{Value: ttnpb.ADR_ACK_LIMIT_64},
			ADRNbTrans:             1,
			MaxDutyCycle:           ttnpb.DUTY_CYCLE_1,
			MaxEIRP:                16,
			PingSlotDataRateIndex:  ttnpb.DATA_RATE_0,
			RejoinCountPeriodicity: ttnpb.REJOIN_COUNT_16,
			RejoinTimePeriodicity:  ttnpb.REJOIN_TIME_0,
			Rx1Delay:               ttnpb.RX_DELAY_1,
			Rx2DataRateIndex:       ttnpb.DATA_RATE_0,
			Rx2Frequency:           869525000,
			Channels:               MakeEU868Channels(),
		},
		DesiredParameters: ttnpb.MACParameters{
			ADRAckDelayExponent:    &ttnpb.ADRAckDelayExponentValue{Value: ttnpb.ADR_ACK_DELAY_32},
			ADRAckLimitExponent:    &ttnpb.ADRAckLimitExponentValue{Value: ttnpb.ADR_ACK_LIMIT_64},
			ADRNbTrans:             1,
			MaxDutyCycle:           ttnpb.DUTY_CYCLE_1,
			MaxEIRP:                16,
			PingSlotDataRateIndex:  ttnpb.DATA_RATE_0,
			RejoinCountPeriodicity: ttnpb.REJOIN_COUNT_16,
			RejoinTimePeriodicity:  ttnpb.REJOIN_TIME_0,
			Rx1Delay:               ttnpb.RX_DELAY_1,
			Rx2DataRateIndex:       ttnpb.DATA_RATE_0,
			Rx2Frequency:           869525000,
			Channels: MakeEU868Channels(
				&ttnpb.MACParameters_Channel{
					UplinkFrequency:   867100000,
					DownlinkFrequency: 867100000,
					MinDataRateIndex:  ttnpb.DATA_RATE_0,
					MaxDataRateIndex:  ttnpb.DATA_RATE_5,
					EnableUplink:      true,
				},
				&ttnpb.MACParameters_Channel{
					UplinkFrequency:   867300000,
					DownlinkFrequency: 867300000,
					MinDataRateIndex:  ttnpb.DATA_RATE_0,
					MaxDataRateIndex:  ttnpb.DATA_RATE_5,
					EnableUplink:      true,
				},
				&ttnpb.MACParameters_Channel{
					UplinkFrequency:   867500000,
					DownlinkFrequency: 867500000,
					MinDataRateIndex:  ttnpb.DATA_RATE_0,
					MaxDataRateIndex:  ttnpb.DATA_RATE_5,
					EnableUplink:      true,
				},
				&ttnpb.MACParameters_Channel{
					UplinkFrequency:   867700000,
					DownlinkFrequency: 867700000,
					MinDataRateIndex:  ttnpb.DATA_RATE_0,
					MaxDataRateIndex:  ttnpb.DATA_RATE_5,
					EnableUplink:      true,
				},
				&ttnpb.MACParameters_Channel{
					UplinkFrequency:   867900000,
					DownlinkFrequency: 867900000,
					MinDataRateIndex:  ttnpb.DATA_RATE_0,
					MaxDataRateIndex:  ttnpb.DATA_RATE_5,
					EnableUplink:      true,
				},
			),
		},
	}
}

func MakeUS915Channels() []*ttnpb.MACParameters_Channel {
	var chs []*ttnpb.MACParameters_Channel
	for i := 0; i < 64; i++ {
		chs = append(chs, &ttnpb.MACParameters_Channel{
			UplinkFrequency:  uint64(902300000 + 200000*i),
			MinDataRateIndex: ttnpb.DATA_RATE_0,
			MaxDataRateIndex: ttnpb.DATA_RATE_3,
			EnableUplink:     true,
		})
	}
	for i := 0; i < 8; i++ {
		chs = append(chs, &ttnpb.MACParameters_Channel{
			UplinkFrequency:  uint64(903000000 + 1600000*i),
			MinDataRateIndex: ttnpb.DATA_RATE_4,
			MaxDataRateIndex: ttnpb.DATA_RATE_4,
			EnableUplink:     true,
		})
	}
	for i := 0; i < 72; i++ {
		chs[i].DownlinkFrequency = uint64(923300000 + 600000*(i%8))
	}
	return chs
}

func MakeDefaultUS915MACState(class ttnpb.Class, ver ttnpb.MACVersion) *ttnpb.MACState {
	return &ttnpb.MACState{
		DeviceClass:         class,
		LoRaWANVersion:      ver,
		PingSlotPeriodicity: ttnpb.PING_EVERY_1S,
		CurrentParameters: ttnpb.MACParameters{
			ADRAckDelayExponent:    &ttnpb.ADRAckDelayExponentValue{Value: ttnpb.ADR_ACK_DELAY_32},
			ADRAckLimitExponent:    &ttnpb.ADRAckLimitExponentValue{Value: ttnpb.ADR_ACK_LIMIT_64},
			ADRNbTrans:             1,
			MaxDutyCycle:           ttnpb.DUTY_CYCLE_1,
			MaxEIRP:                30,
			PingSlotDataRateIndex:  ttnpb.DATA_RATE_0,
			RejoinCountPeriodicity: ttnpb.REJOIN_COUNT_16,
			RejoinTimePeriodicity:  ttnpb.REJOIN_TIME_0,
			Rx1Delay:               ttnpb.RX_DELAY_1,
			Rx2DataRateIndex:       ttnpb.DATA_RATE_8,
			Rx2Frequency:           923300000,
			Channels:               MakeUS915Channels(),
		},
		DesiredParameters: ttnpb.MACParameters{
			ADRAckDelayExponent:    &ttnpb.ADRAckDelayExponentValue{Value: ttnpb.ADR_ACK_DELAY_32},
			ADRAckLimitExponent:    &ttnpb.ADRAckLimitExponentValue{Value: ttnpb.ADR_ACK_LIMIT_64},
			ADRNbTrans:             1,
			MaxDutyCycle:           ttnpb.DUTY_CYCLE_1,
			MaxEIRP:                30,
			PingSlotDataRateIndex:  ttnpb.DATA_RATE_0,
			RejoinCountPeriodicity: ttnpb.REJOIN_COUNT_16,
			RejoinTimePeriodicity:  ttnpb.REJOIN_TIME_0,
			Rx1Delay:               ttnpb.RX_DELAY_1,
			Rx2DataRateIndex:       ttnpb.DATA_RATE_8,
			Rx2Frequency:           923300000,
			Channels: func() []*ttnpb.MACParameters_Channel {
				ret := MakeUS915Channels()
				for _, ch := range ret {
					switch ch.UplinkFrequency {
					case 903900000,
						904100000,
						904300000,
						904500000,
						904700000,
						904900000,
						905100000,
						905300000:
						continue
					}
					ch.EnableUplink = false
				}
				return ret
			}(),
		},
	}
}

func MakeRxMetadataSlice(mds ...*ttnpb.RxMetadata) []*ttnpb.RxMetadata {
	return append([]*ttnpb.RxMetadata{
		{
			GatewayIdentifiers:     ttnpb.GatewayIdentifiers{GatewayID: "gateway-test-1"},
			SNR:                    -9,
			UplinkToken:            []byte("token-gtw-1"),
			DownlinkPathConstraint: ttnpb.DOWNLINK_PATH_CONSTRAINT_NONE,
		},
		{
			GatewayIdentifiers:     ttnpb.GatewayIdentifiers{GatewayID: "gateway-test-3"},
			SNR:                    -5.3,
			UplinkToken:            []byte("token-gtw-3"),
			DownlinkPathConstraint: ttnpb.DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER,
		},
		{
			GatewayIdentifiers:     ttnpb.GatewayIdentifiers{GatewayID: "gateway-test-5"},
			SNR:                    12,
			UplinkToken:            []byte("token-gtw-5"),
			DownlinkPathConstraint: ttnpb.DOWNLINK_PATH_CONSTRAINT_NEVER,
		},
		{
			GatewayIdentifiers:     ttnpb.GatewayIdentifiers{GatewayID: "gateway-test-0"},
			SNR:                    5.2,
			UplinkToken:            []byte("token-gtw-0"),
			DownlinkPathConstraint: ttnpb.DOWNLINK_PATH_CONSTRAINT_NONE,
		},
		{
			GatewayIdentifiers:     ttnpb.GatewayIdentifiers{GatewayID: "gateway-test-2"},
			SNR:                    6.3,
			UplinkToken:            []byte("token-gtw-2"),
			DownlinkPathConstraint: ttnpb.DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER,
		},
		{
			GatewayIdentifiers:     ttnpb.GatewayIdentifiers{GatewayID: "gateway-test-4"},
			SNR:                    -7,
			UplinkToken:            []byte("token-gtw-4"),
			DownlinkPathConstraint: ttnpb.DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER,
		},
	}, mds...)
}

func NewISPeer(ctx context.Context, is interface {
	ttnpb.ApplicationAccessServer
}) cluster.Peer {
	return test.Must(test.NewGRPCServerPeer(ctx, is, ttnpb.RegisterApplicationAccessServer)).(cluster.Peer)
}

func NewGSPeer(ctx context.Context, gs interface {
	ttnpb.NsGsServer
}) cluster.Peer {
	return test.Must(test.NewGRPCServerPeer(ctx, gs, ttnpb.RegisterNsGsServer)).(cluster.Peer)
}

func NewJSPeer(ctx context.Context, js interface {
	ttnpb.NsJsServer
}) cluster.Peer {
	return test.Must(test.NewGRPCServerPeer(ctx, js, ttnpb.RegisterNsJsServer)).(cluster.Peer)
}

var _ ApplicationUplinkQueue = MockApplicationUplinkQueue{}

// MockApplicationUplinkQueue is a mock ApplicationUplinkQueue used for testing.
type MockApplicationUplinkQueue struct {
	AddFunc       func(ctx context.Context, ups ...*ttnpb.ApplicationUp) error
	SubscribeFunc func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, f func(context.Context, *ttnpb.ApplicationUp) error) error
}

// Add calls AddFunc if set and panics otherwise.
func (m MockApplicationUplinkQueue) Add(ctx context.Context, ups ...*ttnpb.ApplicationUp) error {
	if m.AddFunc == nil {
		panic("Add called, but not set")
	}
	return m.AddFunc(ctx, ups...)
}

// Subscribe calls SubscribeFunc if set and panics otherwise.
func (m MockApplicationUplinkQueue) Subscribe(ctx context.Context, appID ttnpb.ApplicationIdentifiers, f func(context.Context, *ttnpb.ApplicationUp) error) error {
	if m.SubscribeFunc == nil {
		panic("Subscribe called, but not set")
	}
	return m.SubscribeFunc(ctx, appID, f)
}

type ApplicationUplinkQueueAddRequest struct {
	Context  context.Context
	Uplinks  []*ttnpb.ApplicationUp
	Response chan<- error
}

func MakeApplicationUplinkQueueAddChFunc(reqCh chan<- ApplicationUplinkQueueAddRequest) func(context.Context, ...*ttnpb.ApplicationUp) error {
	return func(ctx context.Context, ups ...*ttnpb.ApplicationUp) error {
		respCh := make(chan error)
		reqCh <- ApplicationUplinkQueueAddRequest{
			Context:  ctx,
			Uplinks:  ups,
			Response: respCh,
		}
		return <-respCh
	}
}

type ApplicationUplinkQueueSubscribeRequest struct {
	Context     context.Context
	Identifiers ttnpb.ApplicationIdentifiers
	Func        func(context.Context, *ttnpb.ApplicationUp) error
	Response    chan<- error
}

func MakeApplicationUplinkQueueSubscribeChFunc(reqCh chan<- ApplicationUplinkQueueSubscribeRequest) func(context.Context, ttnpb.ApplicationIdentifiers, func(context.Context, *ttnpb.ApplicationUp) error) error {
	return func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, f func(context.Context, *ttnpb.ApplicationUp) error) error {
		respCh := make(chan error)
		reqCh <- ApplicationUplinkQueueSubscribeRequest{
			Context:     ctx,
			Identifiers: appID,
			Func:        f,
			Response:    respCh,
		}
		return <-respCh
	}
}

// ApplicationUplinkSubscribeBlockFunc is ApplicationUplinks.Subscribe function, which blocks until ctx is done and returns nil.
func ApplicationUplinkSubscribeBlockFunc(ctx context.Context, _ func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time) error) error {
	<-ctx.Done()
	return ctx.Err()
}

var _ DownlinkTaskQueue = MockDownlinkTaskQueue{}

// MockDownlinkTaskQueue is a mock DownlinkTaskQueue used for testing.
type MockDownlinkTaskQueue struct {
	AddFunc func(ctx context.Context, devID ttnpb.EndDeviceIdentifiers, t time.Time, replace bool) error
	PopFunc func(ctx context.Context, f func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time) error) error
}

// Add calls AddFunc if set and panics otherwise.
func (m MockDownlinkTaskQueue) Add(ctx context.Context, devID ttnpb.EndDeviceIdentifiers, t time.Time, replace bool) error {
	if m.AddFunc == nil {
		panic("Add called, but not set")
	}
	return m.AddFunc(ctx, devID, t, replace)
}

// Pop calls PopFunc if set and panics otherwise.
func (m MockDownlinkTaskQueue) Pop(ctx context.Context, f func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time) error) error {
	if m.PopFunc == nil {
		panic("Pop called, but not set")
	}
	return m.PopFunc(ctx, f)
}

type DownlinkTaskAddRequest struct {
	Context     context.Context
	Identifiers ttnpb.EndDeviceIdentifiers
	Time        time.Time
	Replace     bool
	Response    chan<- error
}

func MakeDownlinkTaskAddChFunc(reqCh chan<- DownlinkTaskAddRequest) func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time, bool) error {
	return func(ctx context.Context, devID ttnpb.EndDeviceIdentifiers, t time.Time, replace bool) error {
		respCh := make(chan error)
		reqCh <- DownlinkTaskAddRequest{
			Context:     ctx,
			Identifiers: devID,
			Time:        t,
			Replace:     replace,
			Response:    respCh,
		}
		return <-respCh
	}
}

type DownlinkTaskPopRequest struct {
	Context  context.Context
	Func     func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time) error
	Response chan<- error
}

func MakeDownlinkTaskPopChFunc(reqCh chan<- DownlinkTaskPopRequest) func(context.Context, func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time) error) error {
	return func(ctx context.Context, f func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time) error) error {
		respCh := make(chan error)
		reqCh <- DownlinkTaskPopRequest{
			Context:  ctx,
			Func:     f,
			Response: respCh,
		}
		return <-respCh
	}
}

// DownlinkTaskPopBlockFunc is DownlinkTasks.Pop function, which blocks until context is done and returns nil.
func DownlinkTaskPopBlockFunc(ctx context.Context, _ func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time) error) error {
	<-ctx.Done()
	return nil
}

var _ DeviceRegistry = MockDeviceRegistry{}

// MockDeviceRegistry is a mock DeviceRegistry used for testing.
type MockDeviceRegistry struct {
	GetByEUIFunc    func(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, error)
	GetByIDFunc     func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string) (*ttnpb.EndDevice, error)
	RangeByAddrFunc func(ctx context.Context, devAddr types.DevAddr, paths []string, f func(*ttnpb.EndDevice) bool) error
	SetByIDFunc     func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string, f func(*ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error)
}

// GetByEUI calls GetByEUIFunc if set and panics otherwise.
func (m MockDeviceRegistry) GetByEUI(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, error) {
	if m.GetByEUIFunc == nil {
		panic("GetByEUI called, but not set")
	}
	return m.GetByEUIFunc(ctx, joinEUI, devEUI, paths)
}

// GetByID calls GetByIDFunc if set and panics otherwise.
func (m MockDeviceRegistry) GetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string) (*ttnpb.EndDevice, error) {
	if m.GetByIDFunc == nil {
		panic("GetByID called, but not set")
	}
	return m.GetByIDFunc(ctx, appID, devID, paths)
}

// RangeByAddr calls RangeByAddrFunc if set and panics otherwise.
func (m MockDeviceRegistry) RangeByAddr(ctx context.Context, devAddr types.DevAddr, paths []string, f func(*ttnpb.EndDevice) bool) error {
	if m.RangeByAddrFunc == nil {
		panic("RangeByAddr called, but not set")
	}
	return m.RangeByAddrFunc(ctx, devAddr, paths, f)
}

// SetByID calls SetByIDFunc if set and panics otherwise.
func (m MockDeviceRegistry) SetByID(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string, f func(*ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error) {
	if m.SetByIDFunc == nil {
		panic("SetByID called, but not set")
	}
	return m.SetByIDFunc(ctx, appID, devID, paths, f)
}

type deviceAndError struct {
	Device *ttnpb.EndDevice
	Error  error
}

type DeviceRegistryGetByEUIResponse deviceAndError

type DeviceRegistryGetByEUIRequest struct {
	Context  context.Context
	JoinEUI  types.EUI64
	DevEUI   types.EUI64
	Paths    []string
	Response chan<- DeviceRegistryGetByEUIResponse
}

func MakeDeviceRegistryGetByEUIChFunc(reqCh chan<- DeviceRegistryGetByEUIRequest) func(context.Context, types.EUI64, types.EUI64, []string) (*ttnpb.EndDevice, error) {
	return func(ctx context.Context, joinEUI, devEUI types.EUI64, paths []string) (*ttnpb.EndDevice, error) {
		respCh := make(chan DeviceRegistryGetByEUIResponse)
		reqCh <- DeviceRegistryGetByEUIRequest{
			Context:  ctx,
			JoinEUI:  joinEUI,
			DevEUI:   devEUI,
			Paths:    paths,
			Response: respCh,
		}
		resp := <-respCh
		return resp.Device, resp.Error
	}
}

type DeviceRegistryGetByIDResponse deviceAndError

type DeviceRegistryGetByIDRequest struct {
	Context                context.Context
	ApplicationIdentifiers ttnpb.ApplicationIdentifiers
	DeviceID               string
	Paths                  []string
	Response               chan<- DeviceRegistryGetByIDResponse
}

func MakeDeviceRegistryGetByIDChFunc(reqCh chan<- DeviceRegistryGetByIDRequest) func(context.Context, ttnpb.ApplicationIdentifiers, string, []string) (*ttnpb.EndDevice, error) {
	return func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string) (*ttnpb.EndDevice, error) {
		respCh := make(chan DeviceRegistryGetByIDResponse)
		reqCh <- DeviceRegistryGetByIDRequest{
			Context:                ctx,
			ApplicationIdentifiers: appID,
			DeviceID:               devID,
			Paths:                  paths,
			Response:               respCh,
		}
		resp := <-respCh
		return resp.Device, resp.Error
	}
}

type DeviceRegistryRangeByAddrRequest struct {
	Context  context.Context
	DevAddr  types.DevAddr
	Paths    []string
	Func     func(*ttnpb.EndDevice) bool
	Response chan<- error
}

func MakeDeviceRegistryRangeByAddrChFunc(reqCh chan<- DeviceRegistryRangeByAddrRequest) func(context.Context, types.DevAddr, []string, func(*ttnpb.EndDevice) bool) error {
	return func(ctx context.Context, devAddr types.DevAddr, paths []string, f func(*ttnpb.EndDevice) bool) error {
		respCh := make(chan error)
		reqCh <- DeviceRegistryRangeByAddrRequest{
			Context:  ctx,
			DevAddr:  devAddr,
			Paths:    paths,
			Func:     f,
			Response: respCh,
		}
		return <-respCh
	}
}

type DeviceRegistrySetByIDResponse deviceAndError

type DeviceRegistrySetByIDRequest struct {
	Context                context.Context
	ApplicationIdentifiers ttnpb.ApplicationIdentifiers
	DeviceID               string
	Paths                  []string
	Func                   func(*ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)
	Response               chan<- DeviceRegistrySetByIDResponse
}

func MakeDeviceRegistrySetByIDChFunc(reqCh chan<- DeviceRegistrySetByIDRequest) func(context.Context, ttnpb.ApplicationIdentifiers, string, []string, func(*ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error) {
	return func(ctx context.Context, appID ttnpb.ApplicationIdentifiers, devID string, paths []string, f func(*ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error) {
		respCh := make(chan DeviceRegistrySetByIDResponse)
		reqCh <- DeviceRegistrySetByIDRequest{
			Context:                ctx,
			ApplicationIdentifiers: appID,
			DeviceID:               devID,
			Paths:                  paths,
			Func:                   f,
			Response:               respCh,
		}
		resp := <-respCh
		return resp.Device, resp.Error
	}
}

type DeviceRegistrySetByIDRequestFuncResponse struct {
	Device *ttnpb.EndDevice
	Paths  []string
	Error  error
}

var _ InteropClient = MockInteropClient{}

// MockInteropClient is a mock InteropClient used for testing.
type MockInteropClient struct {
	HandleJoinRequestFunc func(context.Context, types.NetID, *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error)
}

// HandleJoinRequest calls HandleJoinRequestFunc if set and panics otherwise.
func (m MockInteropClient) HandleJoinRequest(ctx context.Context, netID types.NetID, req *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
	if m.HandleJoinRequestFunc == nil {
		panic("HandleJoinRequest called, but not set")
	}
	return m.HandleJoinRequestFunc(ctx, netID, req)
}

type InteropClientHandleJoinRequestResponse struct {
	Response *ttnpb.JoinResponse
	Error    error
}

type InteropClientHandleJoinRequestRequest struct {
	Context  context.Context
	NetID    types.NetID
	Request  *ttnpb.JoinRequest
	Response chan<- InteropClientHandleJoinRequestResponse
}

func MakeInteropClientHandleJoinRequestChFunc(reqCh chan<- InteropClientHandleJoinRequestRequest) func(context.Context, types.NetID, *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
	return func(ctx context.Context, netID types.NetID, msg *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
		respCh := make(chan InteropClientHandleJoinRequestResponse)
		reqCh <- InteropClientHandleJoinRequestRequest{
			Context:  ctx,
			NetID:    netID,
			Request:  msg,
			Response: respCh,
		}
		resp := <-respCh
		return resp.Response, resp.Error
	}
}

var _ ttnpb.NsJsServer = &MockNsJsServer{}

type MockNsJsServer struct {
	HandleJoinFunc  func(context.Context, *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error)
	GetNwkSKeysFunc func(context.Context, *ttnpb.SessionKeyRequest) (*ttnpb.NwkSKeysResponse, error)
}

// HandleJoin calls HandleJoinFunc if set and panics otherwise.
func (m MockNsJsServer) HandleJoin(ctx context.Context, req *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
	if m.HandleJoinFunc == nil {
		panic("HandleJoin called, but not set")
	}
	return m.HandleJoinFunc(ctx, req)
}

// GetNwkSKeys calls GetNwkSKeysFunc if set and panics otherwise.
func (m MockNsJsServer) GetNwkSKeys(ctx context.Context, req *ttnpb.SessionKeyRequest) (*ttnpb.NwkSKeysResponse, error) {
	if m.GetNwkSKeysFunc == nil {
		panic("GetNwkSKeys called, but not set")
	}
	return m.GetNwkSKeysFunc(ctx, req)
}

type NsJsHandleJoinResponse struct {
	Response *ttnpb.JoinResponse
	Error    error
}

type NsJsHandleJoinRequest struct {
	Context  context.Context
	Message  *ttnpb.JoinRequest
	Response chan<- NsJsHandleJoinResponse
}

func MakeNsJsHandleJoinChFunc(reqCh chan<- NsJsHandleJoinRequest) func(context.Context, *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
	return func(ctx context.Context, msg *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
		respCh := make(chan NsJsHandleJoinResponse)
		reqCh <- NsJsHandleJoinRequest{
			Context:  ctx,
			Message:  msg,
			Response: respCh,
		}
		resp := <-respCh
		return resp.Response, resp.Error
	}
}

var _ ttnpb.NsJsClient = &MockNsJsClient{}

type MockNsJsClient struct {
	*test.MockClientStream
	HandleJoinFunc  func(context.Context, *ttnpb.JoinRequest, ...grpc.CallOption) (*ttnpb.JoinResponse, error)
	GetNwkSKeysFunc func(context.Context, *ttnpb.SessionKeyRequest, ...grpc.CallOption) (*ttnpb.NwkSKeysResponse, error)
}

// HandleJoin calls HandleJoinFunc if set and panics otherwise.
func (m MockNsJsClient) HandleJoin(ctx context.Context, req *ttnpb.JoinRequest, opts ...grpc.CallOption) (*ttnpb.JoinResponse, error) {
	if m.HandleJoinFunc == nil {
		panic("HandleJoin called, but not set")
	}
	return m.HandleJoinFunc(ctx, req, opts...)
}

// GetNwkSKeys calls GetNwkSKeysFunc if set and panics otherwise.
func (m MockNsJsClient) GetNwkSKeys(ctx context.Context, req *ttnpb.SessionKeyRequest, opts ...grpc.CallOption) (*ttnpb.NwkSKeysResponse, error) {
	if m.GetNwkSKeysFunc == nil {
		panic("GetNwkSKeys called, but not set")
	}
	return m.GetNwkSKeysFunc(ctx, req, opts...)
}

var _ ttnpb.NsGsServer = &MockNsGsServer{}

type MockNsGsServer struct {
	ScheduleDownlinkFunc func(context.Context, *ttnpb.DownlinkMessage) (*ttnpb.ScheduleDownlinkResponse, error)
}

// ScheduleDownlink calls ScheduleDownlinkFunc if set and panics otherwise.
func (m MockNsGsServer) ScheduleDownlink(ctx context.Context, msg *ttnpb.DownlinkMessage) (*ttnpb.ScheduleDownlinkResponse, error) {
	if m.ScheduleDownlinkFunc == nil {
		panic("ScheduleDownlink called, but not set")
	}
	return m.ScheduleDownlinkFunc(ctx, msg)
}

type NsGsScheduleDownlinkResponse struct {
	Response *ttnpb.ScheduleDownlinkResponse
	Error    error
}

type NsGsScheduleDownlinkRequest struct {
	Context  context.Context
	Message  *ttnpb.DownlinkMessage
	Response chan<- NsGsScheduleDownlinkResponse
}

func MakeNsGsScheduleDownlinkChFunc(reqCh chan<- NsGsScheduleDownlinkRequest) func(context.Context, *ttnpb.DownlinkMessage) (*ttnpb.ScheduleDownlinkResponse, error) {
	return func(ctx context.Context, msg *ttnpb.DownlinkMessage) (*ttnpb.ScheduleDownlinkResponse, error) {
		respCh := make(chan NsGsScheduleDownlinkResponse)
		reqCh <- NsGsScheduleDownlinkRequest{
			Context:  ctx,
			Message:  msg,
			Response: respCh,
		}
		resp := <-respCh
		return resp.Response, resp.Error
	}
}

type WindowEndRequest struct {
	Context  context.Context
	Message  *ttnpb.UplinkMessage
	Response chan<- time.Time
}

func MakeWindowEndChFunc(reqCh chan<- WindowEndRequest) func(ctx context.Context, msg *ttnpb.UplinkMessage) <-chan time.Time {
	return func(ctx context.Context, msg *ttnpb.UplinkMessage) <-chan time.Time {
		respCh := make(chan time.Time)
		reqCh <- WindowEndRequest{
			Context:  ctx,
			Message:  msg,
			Response: respCh,
		}
		return respCh
	}
}

func AssertDownlinkTaskAddRequest(ctx context.Context, reqCh <-chan DownlinkTaskAddRequest, assert func(context.Context, ttnpb.EndDeviceIdentifiers, time.Time, bool) bool, resp error) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("Timed out while waiting for DownlinkTasks.Add to be called")
		return false

	case req := <-reqCh:
		if !assert(req.Context, req.Identifiers, req.Time, req.Replace) {
			return false
		}
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for DownlinkTasks.Add response to be processed")
			return false

		case req.Response <- resp:
			return true
		}
	}
}

func AssertApplicationUplinkQueueAddRequest(ctx context.Context, reqCh <-chan ApplicationUplinkQueueAddRequest, assert func(context.Context, ...*ttnpb.ApplicationUp) bool, resp error) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("Timed out while waiting for ApplicationUplinkQueue.Add to be called")
		return false

	case req := <-reqCh:
		if !assert(req.Context, req.Uplinks...) {
			return false
		}
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for ApplicationUplinkQueue.Add response to be processed")
			return false

		case req.Response <- resp:
			return true
		}
	}
}

func AssertInteropClientHandleJoinRequestRequest(ctx context.Context, reqCh <-chan InteropClientHandleJoinRequestRequest, assert func(context.Context, types.NetID, *ttnpb.JoinRequest) bool, resp InteropClientHandleJoinRequestResponse) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("Timed out while waiting for InteropClient.HandleJoinRequest to be called")
		return false

	case req := <-reqCh:
		if !assert(req.Context, req.NetID, req.Request) {
			return false
		}
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for InteropClient.HandleJoinRequest response to be processed")
			return false

		case req.Response <- resp:
			return true
		}
	}
}

func AssertNsJsHandleJoinRequest(ctx context.Context, reqCh <-chan NsJsHandleJoinRequest, assert func(ctx context.Context, msg *ttnpb.JoinRequest) bool, resp NsJsHandleJoinResponse) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("Timed out while waiting for NsJs.HandleJoin to be called")
		return false

	case req := <-reqCh:
		if !assert(req.Context, req.Message) {
			return false
		}
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for NsJs.HandleJoin response to be processed")
			return false

		case req.Response <- resp:
			return true
		}
	}
}

func AssertAuthNsJsHandleJoinRequest(ctx context.Context, authReqCh <-chan test.ClusterAuthRequest, joinReqCh <-chan NsJsHandleJoinRequest, joinAssert func(ctx context.Context, msg *ttnpb.JoinRequest) bool, authResp grpc.CallOption, joinResp NsJsHandleJoinResponse) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	if !test.AssertClusterAuthRequest(ctx, authReqCh, authResp) {
		return false
	}
	return AssertNsJsHandleJoinRequest(ctx, joinReqCh, joinAssert, joinResp)
}

func AssertNsJsPeerHandleAuthJoinRequest(ctx context.Context, peerReqCh <-chan test.ClusterGetPeerRequest, authReqCh <-chan test.ClusterAuthRequest, idsAssert func(ctx context.Context, ids ttnpb.Identifiers) bool, joinAssert func(ctx context.Context, msg *ttnpb.JoinRequest) bool, authResp grpc.CallOption, joinResp NsJsHandleJoinResponse) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	joinReqCh := make(chan NsJsHandleJoinRequest)
	if !test.AssertClusterGetPeerRequest(ctx, peerReqCh, func(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
		return assertions.New(t).So(role, should.Equal, ttnpb.ClusterRole_JOIN_SERVER) && idsAssert(ctx, ids)
	},
		test.ClusterGetPeerResponse{
			Peer: NewJSPeer(ctx, &MockNsJsServer{
				HandleJoinFunc: MakeNsJsHandleJoinChFunc(joinReqCh),
			}),
		},
	) {
		return false
	}
	return AssertAuthNsJsHandleJoinRequest(ctx, authReqCh, joinReqCh, joinAssert, authResp, joinResp)
}

func AssertNsGsScheduleDownlinkRequest(ctx context.Context, reqCh <-chan NsGsScheduleDownlinkRequest, assert func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool, resp NsGsScheduleDownlinkResponse) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("Timed out while waiting for NsGs.ScheduleDownlink to be called")
		return false

	case req := <-reqCh:
		if !assert(req.Context, req.Message) {
			return false
		}
		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for NsGs.ScheduleDownlink response to be processed")
			return false

		case req.Response <- resp:
			return true
		}
	}
}

func AssertAuthNsGsScheduleDownlinkRequest(ctx context.Context, authReqCh <-chan test.ClusterAuthRequest, scheduleReqCh <-chan NsGsScheduleDownlinkRequest, scheduleAssert func(ctx context.Context, msg *ttnpb.DownlinkMessage) bool, authResp grpc.CallOption, scheduleResp NsGsScheduleDownlinkResponse) bool {
	if !test.AssertClusterAuthRequest(ctx, authReqCh, authResp) {
		return false
	}
	return AssertNsGsScheduleDownlinkRequest(ctx, scheduleReqCh, scheduleAssert, scheduleResp)
}

func AssertLinkApplication(ctx context.Context, conn *grpc.ClientConn, getPeerCh <-chan test.ClusterGetPeerRequest, eventsPublishCh <-chan test.EventPubSubPublishRequest, appID ttnpb.ApplicationIdentifiers, replaceEvents ...events.Event) (ttnpb.AsNs_LinkApplicationClient, func(error) events.Event, bool) {
	t := test.MustTFromContext(ctx)
	t.Helper()

	a := assertions.New(t)

	listRightsCh := make(chan test.ApplicationAccessListRightsRequest)
	defer func() {
		close(listRightsCh)
	}()

	var link ttnpb.AsNs_LinkApplicationClient
	var err error
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		link, err = ttnpb.NewAsNsClient(conn).LinkApplication(
			(rpcmetadata.MD{
				ID: appID.ApplicationID,
			}).ToOutgoingContext(ctx),
			grpc.PerRPCCredentials(rpcmetadata.MD{
				AuthType:      "Bearer",
				AuthValue:     "link-application-key",
				AllowInsecure: true,
			}),
		)
		wg.Done()
	}()

	var reqCIDs []string
	if !a.So(test.AssertClusterGetPeerRequest(ctx, getPeerCh,
		func(ctx context.Context, role ttnpb.ClusterRole, ids ttnpb.Identifiers) bool {
			reqCIDs = events.CorrelationIDsFromContext(ctx)
			return a.So(role, should.Equal, ttnpb.ClusterRole_ACCESS) && a.So(ids, should.BeNil)
		},
		test.ClusterGetPeerResponse{
			Peer: NewISPeer(ctx, &test.MockApplicationAccessServer{
				ListRightsFunc: test.MakeApplicationAccessListRightsChFunc(listRightsCh),
			}),
		},
	), should.BeTrue) {
		return nil, nil, false
	}

	a.So(reqCIDs, should.HaveLength, 1)

	if !a.So(test.AssertListRightsRequest(ctx, listRightsCh,
		func(ctx context.Context, ids ttnpb.Identifiers) bool {
			md := rpcmetadata.FromIncomingContext(ctx)
			cids := events.CorrelationIDsFromContext(ctx)
			return a.So(cids, should.NotResemble, reqCIDs) &&
				a.So(cids, should.HaveLength, 1) &&
				a.So(md.AuthType, should.Equal, "Bearer") &&
				a.So(md.AuthValue, should.Equal, "link-application-key") &&
				a.So(ids, should.Resemble, &appID)
		}, ttnpb.RIGHT_APPLICATION_LINK,
	), should.BeTrue) {
		return nil, nil, false
	}

	if !a.So(test.AssertEventPubSubPublishRequests(ctx, eventsPublishCh, 1+len(replaceEvents), func(evs ...events.Event) bool {
		return a.So(evs, should.HaveSameElements, append(
			[]events.Event{EvtBeginApplicationLink(events.ContextWithCorrelationID(test.Context(), reqCIDs...), appID, nil)},
			replaceEvents...,
		), test.EventEqual)
	}), should.BeTrue) {
		t.Error("AS link events assertion failed")
		return nil, nil, false
	}

	if !a.So(test.WaitContext(ctx, wg.Wait), should.BeTrue) {
		t.Error("Timed out while waiting for AS link to open")
		return nil, nil, false
	}
	return link, func(err error) events.Event {
		return EvtEndApplicationLink(events.ContextWithCorrelationID(test.Context(), reqCIDs...), appID, err)
	}, a.So(err, should.BeNil)
}

func AssertNetworkServerClose(ctx context.Context, ns *NetworkServer) bool {
	t := test.MustTFromContext(ctx)
	t.Helper()
	if !test.WaitContext(ctx, ns.Close) {
		t.Error("Timed out while waiting for Network Server to gracefully close")
		return false
	}
	return true
}

type ApplicationUplinkQueueEnvironment struct {
	Add       <-chan ApplicationUplinkQueueAddRequest
	Subscribe <-chan ApplicationUplinkQueueSubscribeRequest
}

func newMockApplicationUplinkQueue(t *testing.T) (ApplicationUplinkQueue, ApplicationUplinkQueueEnvironment, func()) {
	t.Helper()

	addCh := make(chan ApplicationUplinkQueueAddRequest)
	subscribeCh := make(chan ApplicationUplinkQueueSubscribeRequest)
	return &MockApplicationUplinkQueue{
			AddFunc:       MakeApplicationUplinkQueueAddChFunc(addCh),
			SubscribeFunc: MakeApplicationUplinkQueueSubscribeChFunc(subscribeCh),
		}, ApplicationUplinkQueueEnvironment{
			Add:       addCh,
			Subscribe: subscribeCh,
		},
		func() {
			select {
			case <-addCh:
				t.Error("ApplicationUplinkQueue.Add call missed")
			default:
				close(addCh)
			}
			select {
			case <-subscribeCh:
				t.Error("ApplicationUplinkQueue.Subscribe call missed")
			default:
				close(subscribeCh)
			}
		}
}

type DeviceRegistryEnvironment struct {
	GetByID     <-chan DeviceRegistryGetByIDRequest
	GetByEUI    <-chan DeviceRegistryGetByEUIRequest
	RangeByAddr <-chan DeviceRegistryRangeByAddrRequest
	SetByID     <-chan DeviceRegistrySetByIDRequest
}

func newMockDeviceRegistry(t *testing.T) (DeviceRegistry, DeviceRegistryEnvironment, func()) {
	t.Helper()

	getByEUICh := make(chan DeviceRegistryGetByEUIRequest)
	getByIDCh := make(chan DeviceRegistryGetByIDRequest)
	rangeByAddrCh := make(chan DeviceRegistryRangeByAddrRequest)
	setByIDCh := make(chan DeviceRegistrySetByIDRequest)
	return &MockDeviceRegistry{
			GetByEUIFunc:    MakeDeviceRegistryGetByEUIChFunc(getByEUICh),
			GetByIDFunc:     MakeDeviceRegistryGetByIDChFunc(getByIDCh),
			RangeByAddrFunc: MakeDeviceRegistryRangeByAddrChFunc(rangeByAddrCh),
			SetByIDFunc:     MakeDeviceRegistrySetByIDChFunc(setByIDCh),
		}, DeviceRegistryEnvironment{
			GetByEUI:    getByEUICh,
			RangeByAddr: rangeByAddrCh,
			SetByID:     setByIDCh,
		},
		func() {
			select {
			case <-getByEUICh:
				t.Error("DeviceRegistry.GetByEUI call missed")
			default:
				close(getByEUICh)
			}
			select {
			case <-getByIDCh:
				t.Error("DeviceRegistry.GetByID call missed")
			default:
				close(getByIDCh)
			}
			select {
			case <-rangeByAddrCh:
				t.Error("DeviceRegistry.RangeByAddr call missed")
			default:
				close(rangeByAddrCh)
			}
			select {
			case <-setByIDCh:
				t.Error("DeviceRegistry.SetByID call missed")
			default:
				close(setByIDCh)
			}
		}
}

type DownlinkTaskQueueEnvironment struct {
	Add <-chan DownlinkTaskAddRequest
	Pop <-chan DownlinkTaskPopRequest
}

func newMockDownlinkTaskQueue(t *testing.T) (DownlinkTaskQueue, DownlinkTaskQueueEnvironment, func()) {
	t.Helper()

	addCh := make(chan DownlinkTaskAddRequest)
	popCh := make(chan DownlinkTaskPopRequest)
	return &MockDownlinkTaskQueue{
			AddFunc: MakeDownlinkTaskAddChFunc(addCh),
			PopFunc: MakeDownlinkTaskPopChFunc(popCh),
		}, DownlinkTaskQueueEnvironment{
			Add: addCh,
			Pop: popCh,
		},
		func() {
			select {
			case <-addCh:
				t.Error("DownlinkTaskQueue.Add call missed")
			default:
				close(addCh)
			}
			select {
			case <-popCh:
				t.Error("DownlinkTaskQueue.Pop call missed")
			default:
				close(popCh)
			}
		}
}

type InteropClientEnvironment struct {
	HandleJoinRequest <-chan InteropClientHandleJoinRequestRequest
}

func newMockInteropClient(t *testing.T) (InteropClient, InteropClientEnvironment, func()) {
	t.Helper()

	handleJoinCh := make(chan InteropClientHandleJoinRequestRequest)
	return &MockInteropClient{
			HandleJoinRequestFunc: MakeInteropClientHandleJoinRequestChFunc(handleJoinCh),
		}, InteropClientEnvironment{
			HandleJoinRequest: handleJoinCh,
		},
		func() {
			select {
			case <-handleJoinCh:
				t.Error("InteropClient.HandleJoin call missed")
			default:
				close(handleJoinCh)
			}
		}
}

type TestEnvironment struct {
	Cluster struct {
		Auth    <-chan test.ClusterAuthRequest
		GetPeer <-chan test.ClusterGetPeerRequest
	}
	CollectionDone     <-chan WindowEndRequest
	DeduplicationDone  <-chan WindowEndRequest
	ApplicationUplinks *ApplicationUplinkQueueEnvironment
	DeviceRegistry     *DeviceRegistryEnvironment
	DownlinkTasks      *DownlinkTaskQueueEnvironment
	Events             <-chan test.EventPubSubPublishRequest
	InteropClient      *InteropClientEnvironment
}

func StartTest(t *testing.T, conf Config, timeout time.Duration, stubDeduplication bool, opts ...Option) (*NetworkServer, context.Context, TestEnvironment, func()) {
	t.Helper()

	logger := test.GetLogger(t)

	ctx := test.ContextWithT(test.Context(), t)
	ctx = log.NewContext(ctx, logger)
	ctx, cancel := context.WithTimeout(ctx, timeout)

	authCh := make(chan test.ClusterAuthRequest)
	getPeerCh := make(chan test.ClusterGetPeerRequest)
	eventsPublishCh := make(chan test.EventPubSubPublishRequest)

	c := component.MustNew(
		logger,
		&component.Config{},
		component.WithClusterNew(func(context.Context, *config.Cluster, ...cluster.Option) (cluster.Cluster, error) {
			return &test.MockCluster{
				AuthFunc:    test.MakeClusterAuthChFunc(authCh),
				GetPeerFunc: test.MakeClusterGetPeerChFunc(getPeerCh),
				JoinFunc:    test.ClusterJoinNilFunc,
				WithVerifiedSourceFunc: func(ctx context.Context) context.Context {
					return clusterauth.NewContext(ctx, nil)
				},
			}, nil
		}),
	)
	c.FrequencyPlans = frequencyplans.NewStore(test.FrequencyPlansFetcher)

	var collectionDoneCh, deduplicationDoneCh chan WindowEndRequest
	if stubDeduplication {
		collectionDoneCh = make(chan WindowEndRequest)
		deduplicationDoneCh = make(chan WindowEndRequest)

		opts = append([]Option{
			WithCollectionDoneFunc(MakeWindowEndChFunc(collectionDoneCh)),
			WithDeduplicationDoneFunc(MakeWindowEndChFunc(deduplicationDoneCh)),
		}, opts...)
	}

	env := TestEnvironment{
		CollectionDone:    collectionDoneCh,
		DeduplicationDone: deduplicationDoneCh,
	}
	env.Cluster.Auth = authCh
	env.Cluster.GetPeer = getPeerCh
	env.Events = eventsPublishCh

	var closeFuncs []func()
	closeFuncs = append(closeFuncs, test.SetDefaultEventsPubSub(&test.MockEventPubSub{
		PublishFunc: test.MakeEventPubSubPublishChFunc(eventsPublishCh),
	}))
	if conf.ApplicationUplinks == nil {
		m, mEnv, closeM := newMockApplicationUplinkQueue(t)
		conf.ApplicationUplinks = m
		env.ApplicationUplinks = &mEnv
		closeFuncs = append(closeFuncs, closeM)
	}
	if conf.Devices == nil {
		m, mEnv, closeM := newMockDeviceRegistry(t)
		conf.Devices = m
		env.DeviceRegistry = &mEnv
		closeFuncs = append(closeFuncs, closeM)
	}
	if conf.DownlinkTasks == nil {
		m, mEnv, closeM := newMockDownlinkTaskQueue(t)
		conf.DownlinkTasks = m
		env.DownlinkTasks = &mEnv
		closeFuncs = append(closeFuncs, closeM)
	}

	ns := test.Must(New(
		c,
		&conf,
		opts...,
	)).(*NetworkServer)

	if ns.interopClient == nil {
		m, mEnv, closeM := newMockInteropClient(t)
		ns.interopClient = m
		env.InteropClient = &mEnv
		closeFuncs = append(closeFuncs, closeM)
	}

	componenttest.StartComponent(t, ns.Component)

	return ns, ctx, env, func() {
		cancel()
		for _, f := range closeFuncs {
			f()
		}
		select {
		case <-authCh:
			t.Error("Cluster.Auth call missed")
		default:
			close(authCh)
		}
		select {
		case <-getPeerCh:
			t.Error("Cluster.GetPeer call missed")
		default:
			close(getPeerCh)
		}
		select {
		case <-eventsPublishCh:
			t.Error("events.Publish call missed")
		default:
			close(eventsPublishCh)
		}
	}
}

func ForEachBand(t *testing.T, f func(makeName func(parts ...string) string, phy band.Band)) {
	for phyID, phy := range band.All {
		for _, phyVer := range phy.Versions() {
			phy, err := phy.Version(phyVer)
			if err != nil {
				t.Errorf("Failed to convert %s band to %s version", phyID, phyVer)
				continue
			}
			f(func(parts ...string) string {
				return fmt.Sprintf("%s/PHY:%s/%s", phyID, phyVer.String(), strings.Join(parts, "/"))
			}, phy)
		}
	}
}

func ForEachMACVersion(f func(makeName func(parts ...string) string, macVersion ttnpb.MACVersion)) {
	for _, macVersion := range []ttnpb.MACVersion{
		ttnpb.MAC_V1_0,
		ttnpb.MAC_V1_0_1,
		ttnpb.MAC_V1_0_2,
		ttnpb.MAC_V1_0_3,
		ttnpb.MAC_V1_1,
	} {
		f(func(parts ...string) string {
			return fmt.Sprintf("MAC:%s/%s", macVersion.String(), strings.Join(parts, "/"))
		}, macVersion)
	}
}

func ForEachBandMACVersion(t *testing.T, f func(makeName func(parts ...string) string, phy band.Band, macVersion ttnpb.MACVersion)) {
	ForEachBand(t, func(makeBandName func(...string) string, phy band.Band) {
		ForEachMACVersion(func(makeMACName func(...string) string, macVersion ttnpb.MACVersion) {
			f(func(parts ...string) string {
				return makeBandName(makeMACName(parts...))
			}, phy, macVersion)
		})
	})
}

func ForEachClass(f func(makeName func(parts ...string) string, class ttnpb.Class)) {
	for _, class := range []ttnpb.Class{
		ttnpb.CLASS_A,
		ttnpb.CLASS_B,
		ttnpb.CLASS_C,
	} {
		f(func(parts ...string) string {
			return fmt.Sprintf("Class:%s/%s", class.String(), strings.Join(parts, "/"))
		}, class)
	}
}

var redisNamespace = [...]string{
	"networkserver_test",
}

const (
	redisConsumerGroup = "ns"
	redisConsumerID    = "test"
)

func NewRedisApplicationUplinkQueue(t testing.TB) (ApplicationUplinkQueue, func() error) {
	cl, flush := test.NewRedis(t, append(redisNamespace[:], "application-uplinks")...)
	return redis.NewApplicationUplinkQueue(cl, 100, redisConsumerGroup, redisConsumerID),
		func() error {
			flush()
			return cl.Close()
		}
}

func NewRedisDeviceRegistry(t testing.TB) (DeviceRegistry, func() error) {
	cl, flush := test.NewRedis(t, append(redisNamespace[:], "devices")...)
	return &redis.DeviceRegistry{
			Redis: cl,
		},
		func() error {
			flush()
			return cl.Close()
		}
}

func NewRedisDownlinkTaskQueue(t testing.TB) (DownlinkTaskQueue, func() error) {
	a := assertions.New(t)

	cl, flush := test.NewRedis(t, append(redisNamespace[:], "downlink-tasks")...)
	q := redis.NewDownlinkTaskQueue(cl, 10000, redisConsumerGroup, redisConsumerID)
	err := q.Init()
	a.So(err, should.BeNil)

	ctx, cancel := context.WithCancel(test.Context())
	errCh := make(chan error, 1)
	go func() {
		t.Log("Running Redis downlink task queue...")
		err := q.Run(ctx)
		errCh <- err
		close(errCh)
		t.Logf("Stopped Redis downlink task queue with error: %s", err)
	}()
	return q,
		func() error {
			cancel()
			err := q.Add(ctx, ttnpb.EndDeviceIdentifiers{
				DeviceID:               "test",
				ApplicationIdentifiers: ttnpb.ApplicationIdentifiers{ApplicationID: "test"},
			}, time.Now(), false)
			if !a.So(err, should.BeNil) {
				t.Errorf("Failed to add mock device to task queue: %s", err)
				return err
			}

			var runErr error
			select {
			case <-time.After(Timeout):
				t.Error("Timed out waiting for redis.DownlinkTaskQueue.Run to return")
			case runErr = <-errCh:
			}

			flush()
			closeErr := cl.Close()
			if runErr != nil && runErr != context.Canceled {
				return runErr
			}
			return closeErr
		}
}
