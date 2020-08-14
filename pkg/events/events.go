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

// Package events implements event handling through a PubSub interface.
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/gogoproto"
	"go.thethings.network/lorawan-stack/v3/pkg/jsonpb"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// Event interface
type Event interface {
	Context() context.Context
	Name() string
	Time() time.Time
	Identifiers() []*ttnpb.EntityIdentifiers
	Data() interface{}
	CorrelationIDs() []string
	Origin() string
	Caller() string
	Visibility() *ttnpb.Rights
	AuthType() string
	AuthTokenID() string
	AuthTokenType() string
	RemoteIP() string
	UserAgent() string
}

func local(evt Event) *event {
	localEvent, ok := evt.(*event)
	if !ok {
		localEvent = &event{
			ctx: evt.Context(),
			innerEvent: ttnpb.Event{
				Name:           evt.Name(),
				Time:           evt.Time(),
				Identifiers:    evt.Identifiers(),
				CorrelationIDs: evt.CorrelationIDs(),
				Origin:         evt.Origin(),
				Visibility:     evt.Visibility(),
			},
			data:   evt.Data(),
			caller: evt.Caller(),
		}
	}
	return localEvent
}

type event struct {
	ctx        context.Context
	innerEvent ttnpb.Event
	data       interface{}
	caller     string
}

var pathPrefix = func() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("could not determine location of events.go")
	}
	return strings.TrimSuffix(file, filepath.Join("pkg", "events", "events.go"))
}()

// IncludeCaller indicates whether the caller of Publish should be included in the event
var IncludeCaller bool

// withCaller returns an event with the Caller field populated, if configured to do so.
// If the original event already had a non-empty Caller, the original event is returned.
func (e *event) withCaller() *event {
	if IncludeCaller && e.caller == "" {
		if _, file, line, ok := runtime.Caller(2); ok {
			clone := *e
			clone.caller = fmt.Sprintf("%s:%d", strings.TrimPrefix(file, pathPrefix), line)
			return &clone
		}
	}
	return e
}

func (e event) MarshalJSON() ([]byte, error) {
	pb, err := Proto(e)
	if err != nil {
		return nil, err
	}
	return jsonpb.TTN().Marshal(pb)
}

func (e *event) UnmarshalJSON(data []byte) error {
	var pb ttnpb.Event
	err := jsonpb.TTN().Unmarshal(data, &pb)
	if err != nil {
		return err
	}
	fromProto, err := FromProto(&pb)
	if err != nil {
		return err
	}
	evt := fromProto.(*event)
	*e = *evt
	return nil
}

func (e event) Context() context.Context                { return e.ctx }
func (e event) Name() string                            { return e.innerEvent.Name }
func (e event) Time() time.Time                         { return e.innerEvent.Time }
func (e event) Identifiers() []*ttnpb.EntityIdentifiers { return e.innerEvent.Identifiers }
func (e event) Data() interface{}                       { return e.data }
func (e event) CorrelationIDs() []string                { return e.innerEvent.CorrelationIDs }
func (e event) Origin() string                          { return e.innerEvent.Origin }
func (e event) Caller() string                          { return e.caller }
func (e event) Visibility() *ttnpb.Rights               { return e.innerEvent.Visibility }
func (e event) UserAgent() string                       { return e.innerEvent.UserAgent }
func (e event) RemoteIP() string                        { return e.innerEvent.RemoteIP }
func (e event) AuthType() string {
	if e.innerEvent.Authentication == nil {
		return ""
	}
	return e.innerEvent.Authentication.Type
}

func (e event) AuthTokenType() string {
	if e.innerEvent.Authentication == nil {
		return ""
	}
	return e.innerEvent.Authentication.TokenType
}

func (e event) AuthTokenID() string {
	if e.innerEvent.Authentication == nil {
		return ""
	}
	return e.innerEvent.Authentication.TokenID
}

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

// New returns a new Event.
// Instead of using New, most implementations should first define an event,
// and then create a new event from that definition.
func New(ctx context.Context, name, description string, opts ...Option) Event {
	return (&definition{name: name, description: description}).New(ctx, opts...)
}

// Proto returns the protobuf representation of the event.
func Proto(e Event) (*ttnpb.Event, error) {
	evt := local(e)
	pb := evt.innerEvent
	ctx, err := marshalContext(e.Context())
	if err != nil {
		return nil, err
	}
	pb.Context = ctx
	if evt.data != nil {
		var err error
		if protoMessage, ok := evt.data.(proto.Message); ok {
			pb.Data, err = pbtypes.MarshalAny(protoMessage)
		} else if errData, ok := evt.data.(error); ok {
			if ttnErrData, ok := errors.From(errData); ok {
				pb.Data, err = pbtypes.MarshalAny(ttnpb.ErrorDetailsToProto(ttnErrData))
			} else {
				pb.Data, err = pbtypes.MarshalAny(&pbtypes.StringValue{Value: errData.Error()})
			}
		} else {
			value, err := gogoproto.Value(evt.data)
			if err != nil {
				return nil, err
			}
			if _, isNull := value.Kind.(*pbtypes.Value_NullValue); !isNull {
				pb.Data, err = pbtypes.MarshalAny(value)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return &pb, nil
}

// FromProto returns the event from its protobuf representation.
func FromProto(pb *ttnpb.Event) (Event, error) {
	ctx, err := unmarshalContext(context.Background(), pb.Context)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if pb.Data != nil {
		any, err := pbtypes.EmptyAny(pb.Data)
		if err != nil {
			return nil, err
		}
		if err = pbtypes.UnmarshalAny(pb.Data, any); err != nil {
			return nil, err
		}
		data = any
		v, ok := any.(*pbtypes.Value)
		if ok {
			iface, err := gogoproto.Interface(v)
			if err != nil {
				return nil, err
			}
			data = iface
		}
	}
	return &event{
		ctx:  ctx,
		data: data,
		innerEvent: ttnpb.Event{
			Name:           pb.Name,
			Time:           pb.Time,
			Identifiers:    pb.Identifiers,
			CorrelationIDs: pb.CorrelationIDs,
			Origin:         pb.Origin,
			Visibility:     pb.Visibility,
		},
	}, nil
}

// UnmarshalJSON unmarshals an event as JSON.
func UnmarshalJSON(data []byte) (Event, error) {
	e := new(event)
	if err := json.Unmarshal(data, e); err != nil {
		return nil, err
	}
	return e, nil
}
