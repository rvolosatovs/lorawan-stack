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

	"github.com/prometheus/client_golang/prometheus"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/metrics"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	evtRejectJoin = events.Define(
		"js.join.reject", "reject join-request",
		events.WithVisibility(ttnpb.RIGHT_APPLICATION_TRAFFIC_READ),
	)
	evtAcceptJoin = events.Define(
		"js.join.accept", "accept join-request",
		events.WithVisibility(ttnpb.RIGHT_APPLICATION_TRAFFIC_READ),
	)
)

const (
	subsystem = "js"
	unknown   = "unknown"
)

var jsMetrics = &messageMetrics{
	joinAccepted: metrics.NewContextualCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "join_accepted_total",
			Help:      "Total number of accepted joins",
		},
		[]string{"net_id"},
	),
	joinRejected: metrics.NewContextualCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "join_rejected_total",
			Help:      "Total number of rejected joins",
		},
		[]string{"error"},
	),
}

func init() {
	metrics.MustRegister(jsMetrics)
}

type messageMetrics struct {
	joinAccepted *metrics.ContextualCounterVec
	joinRejected *metrics.ContextualCounterVec
}

func (m messageMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.joinAccepted.Describe(ch)
	m.joinRejected.Describe(ch)
}

func (m messageMetrics) Collect(ch chan<- prometheus.Metric) {
	m.joinAccepted.Collect(ch)
	m.joinRejected.Collect(ch)
}

func registerAcceptJoin(ctx context.Context, dev *ttnpb.EndDevice, msg *ttnpb.JoinRequest) {
	events.Publish(evtAcceptJoin.NewWithIdentifiersAndData(ctx, dev.EndDeviceIdentifiers, nil))
	jsMetrics.joinAccepted.WithLabelValues(ctx, msg.NetID.String()).Inc()
}

func registerRejectJoin(ctx context.Context, req *ttnpb.JoinRequest, err error) {
	events.Publish(evtRejectJoin.NewWithIdentifiersAndData(ctx, nil, err))
	if ttnErr, ok := errors.From(err); ok {
		jsMetrics.joinRejected.WithLabelValues(ctx, ttnErr.FullName()).Inc()
	} else {
		jsMetrics.joinRejected.WithLabelValues(ctx, unknown).Inc()
	}
}
