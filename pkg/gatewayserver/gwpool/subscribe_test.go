// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package gwpool_test

import (
	"context"
	"testing"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/gatewayserver/gwpool"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/TheThingsNetwork/ttn/pkg/util/test"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestPoolUplinks(t *testing.T) {
	a := assertions.New(t)

	p := gwpool.NewPool(test.GetLogger(t, "TestPoolUplinks"))

	gatewayID := "gateway"
	link := &dummyLink{
		AcceptSendingUplinks: true,

		NextUplink: make(chan *ttnpb.GatewayUp),
	}
	emptyUplink := &ttnpb.GatewayUp{}
	upstream := p.Subscribe(ttnpb.GatewayIdentifier{GatewayID: gatewayID}, link)

	go func() { link.NextUplink <- emptyUplink }()
	newUplink := <-upstream
	a.So(newUplink, should.Equal, emptyUplink)

	link.AcceptSendingUplinks = false
	go func() { link.NextUplink <- emptyUplink }()
	newUplink = <-upstream
	a.So(newUplink, should.BeNil)
}

func TestDoneContextUplinks(t *testing.T) {
	p := gwpool.NewPool(test.GetLogger(t, "TestDoneContextUplinks"))

	ctx, cancel := context.WithCancel(context.Background())

	gatewayID := "gateway"
	link := &dummyLink{
		AcceptSendingUplinks: true,

		context:       ctx,
		cancelContext: cancel,

		NextUplink: make(chan *ttnpb.GatewayUp),
	}
	cancel()
	emptyUplink := &ttnpb.GatewayUp{}
	upstream := p.Subscribe(ttnpb.GatewayIdentifier{GatewayID: gatewayID}, link)
	go func() { link.NextUplink <- emptyUplink }()
	time.Sleep(time.Millisecond)
	select {
	case _, ok := <-upstream:
		if ok {
			t.Error("Stream not closed, message received")
		}
	default:
		t.Error("Stream not closed, no message")
	}
}

func TestSubscribeTwice(t *testing.T) {
	p := gwpool.NewPool(test.GetLogger(t, "TestSubscribeTwice"))

	gateway := ttnpb.GatewayIdentifier{GatewayID: "gateway"}

	link := &dummyLink{}
	newLink := &dummyLink{}

	p.Subscribe(gateway, link)
	p.Subscribe(gateway, newLink)
}
