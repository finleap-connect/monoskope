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
func NewEventStoreMetrics() *EventStoreMetrics {
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

	return m.
		register(m.TransmittedTotalCounter.MetricVec).
		register(m.StoredTotalCounter.MetricVec).
		register(m.RetrievedTotalCounter.MetricVec).
		register(m.StoredHistogram.MetricVec).
		register(m.RetrievedHistogram.MetricVec)
}

// Registers all metrics with prometheus default registerer
func (m *EventStoreMetrics) register(v *prom.MetricVec) *EventStoreMetrics {
	prom.MustRegister(v)
	return m
}
