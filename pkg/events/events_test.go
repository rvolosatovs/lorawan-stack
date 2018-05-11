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

package events_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/TheThingsNetwork/ttn/pkg/events"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

type wrappedEvent struct {
	events.Event
}

func TestEvents(t *testing.T) {
	a := assertions.New(t)

	events.IncludeCaller = true

	var totalEvents int
	newTotal := make(chan int)
	allEvents := events.HandlerFunc(func(e events.Event) {
		t.Logf("Received event %v", e)
		a.So(e.Time().IsZero(), should.BeFalse)
		a.So(e.Context(), should.NotBeNil)
		totalEvents++
		newTotal <- totalEvents
	})

	var eventCh = make(chan events.Event)
	handler := events.HandlerFunc(func(e events.Event) {
		eventCh <- e
	})

	pubsub := events.NewPubSub()

	pubsub.Subscribe("**", allEvents)

	ctx := events.ContextWithCorrelationID(context.Background(), t.Name())

	pubsub.Publish(events.New(ctx, "test.evt0", nil, nil))
	a.So(<-newTotal, should.Equal, 1)
	a.So(eventCh, should.HaveLength, 0)

	pubsub.Subscribe("test.*", handler)
	pubsub.Subscribe("test.*", handler) // second time should not matter

	evt := events.New(ctx, "test.evt1", "id", "hello")
	a.So(evt.CorrelationIDs(), should.Contain, t.Name())

	wrapped := wrappedEvent{Event: evt}
	pubsub.Publish(wrapped)
	a.So(<-newTotal, should.Equal, 2)

	received := <-eventCh
	a.So(received.Context(), should.Equal, evt.Context())
	a.So(received.Name(), should.Equal, evt.Name())
	a.So(received.Time(), should.Equal, evt.Time())
	a.So(received.Identifiers(), should.Equal, evt.Identifiers())
	a.So(received.Data(), should.Equal, evt.Data())
	a.So(received.CorrelationIDs(), should.Resemble, evt.CorrelationIDs())
	a.So(received.Origin(), should.Equal, evt.Origin())

	pubsub.Unsubscribe("test.*", handler)

	pubsub.Publish(events.New(ctx, "test.evt2", nil, nil))
	a.So(<-newTotal, should.Equal, 3)
	a.So(eventCh, should.HaveLength, 0)
}

func TestUnmarshalJSON(t *testing.T) {
	a := assertions.New(t)
	evt := events.New(context.Background(), "name", "identifiers", "data")
	json, err := json.Marshal(evt)
	a.So(err, should.BeNil)
	evt2, err := events.UnmarshalJSON(json)
	a.So(err, should.BeNil)
	a.So(evt2, should.Resemble, evt)
}

func Example() {
	// This is only for the test
	var wg sync.WaitGroup
	wg.Add(1)
	events.Subscribe("ns.**", events.HandlerFunc(func(e events.Event) {
		fmt.Printf("Received event %s\n", e.Name())
		wg.Done()
	}))

	// You can send any arbitrary event; you don't have to pass any identifiers or data.
	events.PublishEvent(context.Background(), "test.hello_world", nil, nil)

	// Defining the event is not mandatory, but will be needed in order to translate the descriptions.
	// Event names are lowercase snake_case and can be dot-separated as component.subsystem.subsystem.event
	// Event descriptions are short descriptions of what the event means.
	var adrSendEvent = events.Define("ns.mac.adr.send_req", "Send ADR Request")

	// These variables come from the request or you got them from the db or something.
	var (
		ctx      = context.Background()
		dev      ttnpb.EndDevice
		requests []ttnpb.MACCommand_LinkADRReq
	)

	// It's nice to be able to correlate events; we use a Correlation ID for that.
	// In most cases, there will already be a correlation ID in the context; this func will add one if there isn't.
	ctx = events.ContextWithEnsuredCorrelationID(ctx)

	// Publishing an event to the events package will dispatch it on the "global" event pubsub.
	events.Publish(adrSendEvent(ctx, dev, requests))

	wg.Wait() // only for the test

	// Output:
	// Received event ns.mac.adr.send_req
}
