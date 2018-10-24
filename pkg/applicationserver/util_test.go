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

package applicationserver_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/rpcserver"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"google.golang.org/grpc/metadata"
)

var (
	// This application will be added to the Entity Registry and to the link registry of the Application Server so that it
	// links automatically on start to the Network Server.
	registeredApplicationID        = ttnpb.ApplicationIdentifiers{ApplicationID: "foo-app"}
	registeredApplicationKey       = "secret"
	registeredApplicationFormatter = ttnpb.PayloadFormatter_FORMATTER_CAYENNELPP

	// This device gets registered in the device registry of the Application Server.
	registeredDevice = &ttnpb.EndDevice{
		EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
			ApplicationIdentifiers: registeredApplicationID,
			DeviceID:               "foo-device",
			JoinEUI:                eui64Ptr(types.EUI64{0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
			DevEUI:                 eui64Ptr(types.EUI64{0x42, 0x42, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
		},
		VersionIDs: &ttnpb.EndDeviceVersionIdentifiers{
			BrandID:         "thethingsproducts",
			ModelID:         "thethingsnode",
			HardwareVersion: "1.0",
			FirmwareVersion: "1.1",
		},
		Formatters: &ttnpb.MessagePayloadFormatters{
			UpFormatter:   ttnpb.PayloadFormatter_FORMATTER_REPOSITORY,
			DownFormatter: ttnpb.PayloadFormatter_FORMATTER_REPOSITORY,
		},
	}

	// This device does not get registered in the device registry of the Application Server and will be created on join
	// and on uplink.
	unregisteredDeviceID = ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: registeredApplicationID,
		DeviceID:               "bar-device",
		JoinEUI:                eui64Ptr(types.EUI64{0x24, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
		DevEUI:                 eui64Ptr(types.EUI64{0x24, 0x24, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	}

	timeout = 10 * test.Delay

	deviceRepositoryData = map[string][]byte{
		"brands.yml": []byte(`version: '3'
brands:
thethingsproducts:
  name: The Things Products
  url: https://www.thethingsnetwork.org`),
		"thethingsproducts/devices.yml": []byte(`version: '3'
devices:
  thethingsnode:
    name: The Things Node`),
		"thethingsproducts/thethingsnode/versions.yml": []byte(`version: '3'
hardware_versions:
  '1.0':
    - firmware_version: 1.1
      payload_format:
        up:
          type: javascript
          parameter: decoder.js
        down:
          type: javascript
          parameter: encoder.js`),
		"thethingsproducts/thethingsnode/1.0/decoder.js": []byte(`function Decoder(payload, f_port) {
	var sum = 0;
	for (i = 0; i < payload.length; i++) {
		sum += payload[i];
	}
	return {
		sum: sum
	};
}`),
		"thethingsproducts/thethingsnode/1.0/encoder.js": []byte(`function Encoder(payload, f_port) {
	var res = [];
	for (i = 0; i < payload.sum; i++) {
		res[i] = 1;
	}
	return res;
}`)}
)

func mustHavePeer(ctx context.Context, c *component.Component, role ttnpb.PeerInfo_Role) {
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		if peer := c.GetPeer(ctx, role, nil); peer != nil {
			return
		}
	}
	panic("could not connect to peer")
}

func eui64Ptr(eui types.EUI64) *types.EUI64 {
	return &eui
}
func devAddrPtr(devAddr types.DevAddr) *types.DevAddr {
	return &devAddr
}
func withDevAddr(ids ttnpb.EndDeviceIdentifiers, devAddr types.DevAddr) ttnpb.EndDeviceIdentifiers {
	ids.DevAddr = &devAddr
	return ids
}

type memStore struct {
	mu    sync.RWMutex
	items map[string][]byte
	New   func() proto.Unmarshaler
}

var errNotFound = errors.DefineNotFound("not_found", "not found")

func (s *memStore) Get(uid string) (proto.Unmarshaler, error) {
	s.mu.RLock()
	buf, ok := s.items[uid]
	s.mu.RUnlock()
	if !ok {
		return nil, errNotFound
	}
	v := s.New()
	if err := v.Unmarshal(buf); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *memStore) Set(uid string, f func(proto.Unmarshaler) (proto.Marshaler, error)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var v proto.Unmarshaler
	buf, ok := s.items[uid]
	if ok {
		v = s.New()
		if err := v.Unmarshal(buf); err != nil {
			return err
		}
	}
	n, err := f(v)
	if err != nil {
		return err
	}
	if n == nil {
		delete(s.items, uid)
	} else if buf, err := n.Marshal(); err != nil {
		return err
	} else {
		s.items[uid] = buf
	}
	return nil
}

func (s *memStore) Range(f func(string, proto.Unmarshaler) bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for uid, buf := range s.items {
		v := s.New()
		if err := v.Unmarshal(buf); err != nil {
			return err
		}
		if !f(uid, v) {
			break
		}
	}
	return nil
}

func (s *memStore) Reset() {
	s.mu.Lock()
	s.items = map[string][]byte{}
	s.mu.Unlock()
}

type memDeviceRegistry struct {
	store memStore
}

func newMemDeviceRegistry() *memDeviceRegistry {
	return &memDeviceRegistry{
		store: memStore{
			items: make(map[string][]byte),
			New:   func() proto.Unmarshaler { return new(ttnpb.EndDevice) },
		},
	}
}

func (r *memDeviceRegistry) Get(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, paths []string) (*ttnpb.EndDevice, error) {
	v, err := r.store.Get(unique.ID(ctx, ids))
	if err != nil {
		return nil, err
	}
	return v.(*ttnpb.EndDevice), nil
}

func (r *memDeviceRegistry) Set(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, paths []string, f func(*ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error)) (*ttnpb.EndDevice, error) {
	var set *ttnpb.EndDevice
	err := r.store.Set(unique.ID(ctx, ids), func(v proto.Unmarshaler) (proto.Marshaler, error) {
		var dev *ttnpb.EndDevice
		if v != nil {
			dev = v.(*ttnpb.EndDevice)
		}
		if dev, _, err := f(dev); err != nil {
			return nil, err
		} else if dev == nil {
			return nil, nil
		} else {
			set = dev
			return dev, nil
		}
	})
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (r *memDeviceRegistry) Reset() {
	r.store.Reset()
}

type memLinkRegistry struct {
	store memStore
}

func newMemLinkRegistry() *memLinkRegistry {
	return &memLinkRegistry{
		store: memStore{
			items: make(map[string][]byte),
			New:   func() proto.Unmarshaler { return new(ttnpb.ApplicationLink) },
		},
	}
}

func (r *memLinkRegistry) Get(ctx context.Context, ids ttnpb.ApplicationIdentifiers, paths []string) (*ttnpb.ApplicationLink, error) {
	v, err := r.store.Get(unique.ID(ctx, ids))
	if err != nil {
		return nil, err
	}
	return v.(*ttnpb.ApplicationLink), nil
}

func (r *memLinkRegistry) Set(ctx context.Context, ids ttnpb.ApplicationIdentifiers, paths []string, f func(*ttnpb.ApplicationLink) (*ttnpb.ApplicationLink, []string, error)) (*ttnpb.ApplicationLink, error) {
	var set *ttnpb.ApplicationLink
	err := r.store.Set(unique.ID(ctx, ids), func(v proto.Unmarshaler) (proto.Marshaler, error) {
		var l *ttnpb.ApplicationLink
		if v != nil {
			l = v.(*ttnpb.ApplicationLink)
		}
		if l, _, err := f(l); err != nil {
			return nil, err
		} else if l == nil {
			return nil, nil
		} else {
			set = l
			return l, nil
		}
	})
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (r *memLinkRegistry) Range(ctx context.Context, paths []string, f func(ttnpb.ApplicationIdentifiers, *ttnpb.ApplicationLink) bool) error {
	var ferr error
	err := r.store.Range(func(uid string, v proto.Unmarshaler) bool {
		ids, ferr := unique.ToApplicationID(uid)
		if ferr != nil {
			return false
		}
		return f(ids, v.(*ttnpb.ApplicationLink))
	})
	if err != nil {
		return err
	}
	if ferr != nil {
		return ferr
	}
	return nil
}

type mockNS struct {
	linkCh          chan ttnpb.ApplicationIdentifiers
	unlinkCh        chan ttnpb.ApplicationIdentifiers
	upCh            chan *ttnpb.ApplicationUp
	downlinkQueueMu sync.RWMutex
	downlinkQueue   map[string][]*ttnpb.ApplicationDownlink
}

func startMockNS(ctx context.Context) (*mockNS, string) {
	ns := &mockNS{
		linkCh:        make(chan ttnpb.ApplicationIdentifiers, 1),
		unlinkCh:      make(chan ttnpb.ApplicationIdentifiers, 1),
		upCh:          make(chan *ttnpb.ApplicationUp, 1),
		downlinkQueue: make(map[string][]*ttnpb.ApplicationDownlink),
	}
	srv := rpcserver.New(ctx)
	ttnpb.RegisterAsNsServer(srv.Server, ns)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go srv.Serve(lis)
	return ns, lis.Addr().String()
}

func (ns *mockNS) LinkApplication(ids *ttnpb.ApplicationIdentifiers, stream ttnpb.AsNs_LinkApplicationServer) error {
	select {
	case ns.linkCh <- *ids:
	default:
	}
	defer func() {
		select {
		case ns.unlinkCh <- *ids:
		default:
		}
	}()
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case up := <-ns.upCh:
			if joinAccept := up.GetJoinAccept(); joinAccept != nil && !joinAccept.PendingSession {
				// Reset the downlink queue on join-accept; it's invalid and AS will replace it.
				ns.downlinkQueueMu.Lock()
				ns.downlinkQueue[unique.ID(stream.Context(), up.EndDeviceIdentifiers)] = nil
				ns.downlinkQueueMu.Unlock()
			}
			if err := stream.Send(up); err != nil {
				return err
			}
		}
	}
}

func (ns *mockNS) reset() {
	ns.downlinkQueueMu.Lock()
	ns.downlinkQueue = make(map[string][]*ttnpb.ApplicationDownlink)
	ns.downlinkQueueMu.Unlock()
}

func (ns *mockNS) DownlinkQueueReplace(ctx context.Context, req *ttnpb.DownlinkQueueRequest) (*pbtypes.Empty, error) {
	ns.downlinkQueueMu.Lock()
	ns.downlinkQueue[unique.ID(ctx, req.EndDeviceIdentifiers)] = req.Downlinks
	ns.downlinkQueueMu.Unlock()
	return ttnpb.Empty, nil
}

func (ns *mockNS) DownlinkQueuePush(ctx context.Context, req *ttnpb.DownlinkQueueRequest) (*pbtypes.Empty, error) {
	ns.downlinkQueueMu.Lock()
	uid := unique.ID(ctx, req.EndDeviceIdentifiers)
	ns.downlinkQueue[uid] = append(ns.downlinkQueue[uid], req.Downlinks...)
	ns.downlinkQueueMu.Unlock()
	return ttnpb.Empty, nil
}

func (ns *mockNS) DownlinkQueueList(ctx context.Context, ids *ttnpb.EndDeviceIdentifiers) (*ttnpb.ApplicationDownlinks, error) {
	ns.downlinkQueueMu.RLock()
	queue := ns.downlinkQueue[unique.ID(ctx, ids)]
	ns.downlinkQueueMu.RUnlock()
	return &ttnpb.ApplicationDownlinks{
		Downlinks: queue,
	}, nil
}

type mockIS struct {
	ttnpb.ApplicationRegistryServer
	ttnpb.ApplicationAccessServer
	applications     map[string]*ttnpb.Application
	applicationAuths map[string][]string
}

func startMockIS(ctx context.Context) (*mockIS, string) {
	is := &mockIS{
		applications:     make(map[string]*ttnpb.Application),
		applicationAuths: make(map[string][]string),
	}
	srv := rpcserver.New(ctx)
	ttnpb.RegisterApplicationRegistryServer(srv.Server, is)
	ttnpb.RegisterApplicationAccessServer(srv.Server, is)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go srv.Serve(lis)
	return is, lis.Addr().String()
}

func (is *mockIS) add(ctx context.Context, ids ttnpb.ApplicationIdentifiers, key string) {
	uid := unique.ID(ctx, ids)
	is.applications[uid] = &ttnpb.Application{
		ApplicationIdentifiers: ids,
	}
	if key != "" {
		is.applicationAuths[uid] = []string{fmt.Sprintf("Key %v", key)}
	}
}

func (is *mockIS) Get(ctx context.Context, req *ttnpb.GetApplicationRequest) (*ttnpb.Application, error) {
	uid := unique.ID(ctx, req.ApplicationIdentifiers)
	app, ok := is.applications[uid]
	if !ok {
		return nil, errNotFound
	}
	return app, nil
}

func (is *mockIS) ListRights(ctx context.Context, ids *ttnpb.ApplicationIdentifiers) (res *ttnpb.Rights, err error) {
	res = &ttnpb.Rights{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	authorization, ok := md["authorization"]
	if !ok || len(authorization) == 0 {
		return
	}
	auths, ok := is.applicationAuths[unique.ID(ctx, *ids)]
	if !ok {
		return
	}
	for _, auth := range auths {
		if auth == authorization[0] {
			res.Rights = append(res.Rights,
				ttnpb.RIGHT_APPLICATION_DEVICES_READ,
				ttnpb.RIGHT_APPLICATION_DEVICES_WRITE,
				ttnpb.RIGHT_APPLICATION_TRAFFIC_READ,
				ttnpb.RIGHT_APPLICATION_TRAFFIC_DOWN_WRITE,
			)
		}
	}
	return
}

type mockJS struct {
	keys map[string]ttnpb.KeyEnvelope
}

func startMockJS(ctx context.Context) (*mockJS, string) {
	js := &mockJS{
		keys: make(map[string]ttnpb.KeyEnvelope),
	}
	srv := rpcserver.New(ctx)
	ttnpb.RegisterAsJsServer(srv.Server, js)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go srv.Serve(lis)
	return js, lis.Addr().String()
}

func (js *mockJS) add(ctx context.Context, devEUI types.EUI64, sessionKeyID string, key ttnpb.KeyEnvelope) {
	js.keys[fmt.Sprintf("%v:%v", devEUI, sessionKeyID)] = key
}

func (js *mockJS) GetAppSKey(ctx context.Context, req *ttnpb.SessionKeyRequest) (*ttnpb.AppSKeyResponse, error) {
	key, ok := js.keys[fmt.Sprintf("%v:%v", req.DevEUI, req.SessionKeyID)]
	if !ok {
		return nil, errNotFound
	}
	return &ttnpb.AppSKeyResponse{
		AppSKey: key,
	}, nil
}
