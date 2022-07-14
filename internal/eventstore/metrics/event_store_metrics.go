// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

// EventStoreMetrics represents a collection of metrics to be registered on a
// Prometheus metrics registry for EventStore metrics.
type EventStoreMetrics struct {
	TransmittedTotalCounter *prom.CounterVec
	StoredTotalCounter      *prom.CounterVec
	RetrievedTotalCounter   *prom.CounterVec
	StoredHistogram         *prom.HistogramVec
	RetrievedHistogram      *prom.HistogramVec
}

// NewEventStoreMetrics returns a ServerMetrics object. Use a new instance of
// ServerMetrics when not using the default Prometheus metrics registry, for
// example when wanting to control which metrics are added to a registry as
// opposed to automatically adding metrics via init functions.
func NewEventStoreMetrics() (*EventStoreMetrics, error) {
	m := &EventStoreMetrics{}
	labels := []string{"event_type", "aggregate_type"}

	m.TransmittedTotalCounter = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "eventstore_transmitted_total",
			Help: "Total number of events transmitted for storing.",
		}, labels,
	)
	m.StoredTotalCounter = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "eventstore_stored_total",
			Help: "Total number of events stored total. Only successful counts.",
		}, labels,
	)
	m.RetrievedTotalCounter = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "eventstore_retrieved_total",
			Help: "Total number of events retrieved.",
		}, labels,
	)
	m.StoredHistogram = prom.NewHistogramVec(
		prom.HistogramOpts{
			Name:    "eventstore_stored_seconds",
			Help:    "Histogram of response latency (seconds) of events that had been stored by the EventStore.",
			Buckets: prom.DefBuckets,
		},
		labels,
	)
	m.RetrievedHistogram = prom.NewHistogramVec(
		prom.HistogramOpts{
			Name:    "eventstore_retrieved_seconds",
			Help:    "Histogram of response latency (seconds) of events that had been retrieved from the EventStore.",
			Buckets: prom.DefBuckets,
		},
		labels,
	)

	metricVacs := []*prom.MetricVec{
		m.TransmittedTotalCounter.MetricVec,
		m.StoredTotalCounter.MetricVec,
		m.RetrievedTotalCounter.MetricVec,
		m.StoredHistogram.MetricVec,
		m.RetrievedHistogram.MetricVec,
	}
	return m, m.register(metricVacs)
}

// Registers all metrics with prometheus default registerer
func (m *EventStoreMetrics) register(vecs []*prom.MetricVec) error {
	for _, v := range vecs {
		err := prom.Register(v)
		if err != nil {
			_, ok := err.(prom.AlreadyRegisteredError)
			if !ok {
				return err
			}
		}
	}
	return nil
}
