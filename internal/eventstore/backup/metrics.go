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

package backup

import (
	"errors"
	"os"
	"time"

	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/push"
)

type MetricsPublisher interface {
	// Start sets a timestamp for start of a backup
	Start()
	// Finished sets the duration an backup has taken and the completion time to now
	Finished()
	// SetSuccessTime adds the success timestamp to the pushed metrics and sets the time to now
	SetSuccessTime()
	// SetFailTime adds a failed timestamp to the pushed metrics and sets the time to now
	SetFailTime()
	// SetEventCount sets the events processed during an backup
	SetEventCount(eventCount float64)
	// SetBytes sets the bytes written to backup
	SetBytes(sizeInBytes float64)
	// CloseAndPush sends the metrics to the prometheus push gateway
	CloseAndPush() error
}

type noopMetricsPublisher struct{}

// NewNoopMetricsPublisher returns a noop metrics publisher which does nothing
func NewNoopMetricsPublisher() MetricsPublisher                { return &noopMetricsPublisher{} }
func (*noopMetricsPublisher) Start()                           {}
func (*noopMetricsPublisher) Finished()                        {}
func (*noopMetricsPublisher) SetSuccessTime()                  {}
func (*noopMetricsPublisher) SetFailTime()                     {}
func (*noopMetricsPublisher) SetEventCount(eventCount float64) {}
func (*noopMetricsPublisher) SetBytes(eventCount float64)      {}
func (*noopMetricsPublisher) CloseAndPush() error              { return nil }

type metricsPublisher struct {
	log            logger.Logger
	completionTime prometheus.Gauge
	failedTime     prometheus.Gauge
	successTime    prometheus.Gauge
	duration       prometheus.Gauge
	bytes          prometheus.Gauge
	events         prometheus.Gauge
	pusher         *push.Pusher
	start          time.Time
}

// NewMetricsPublisher creates a new backup.MetricsPublisher publishing metrics to a prometheus pushgateway
func NewMetricsPublisher(pushGatewayUrl string) (MetricsPublisher, error) {
	if pushGatewayUrl == "" {
		return nil, errors.New("URL of prometheus push gateway invalid.")
	}

	mp := &metricsPublisher{
		log: logger.WithName("metrics-publisher"),
	}

	mp.completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_last_completion_timestamp_seconds",
		Help: "The timestamp of the last completion of a backup, successful or not.",
	})
	mp.successTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_last_success_timestamp_seconds",
		Help: "The timestamp of the last successful completion of a backup.",
	})
	mp.failedTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_last_fail_timestamp_seconds",
		Help: "The timestamp of the last fail of a backup.",
	})
	mp.duration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_duration_seconds",
		Help: "The duration of the last backup in seconds.",
	})
	mp.bytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_size_in_bytes",
		Help: "The number of bytes processed in the last backup.",
	})
	mp.events = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backup_event_count",
		Help: "The number of events processed in the last backup.",
	})

	// We use a registry here to benefit from the consistency checks that
	// happen during registration.
	registry := prometheus.NewRegistry()
	registry.MustRegister(mp.completionTime, mp.duration, mp.bytes, collectors.NewGoCollector())
	// Note that successTime is not registered.

	jobName := os.Getenv("K8S_JOB")
	mp.pusher = push.New(pushGatewayUrl, jobName).Gatherer(registry)

	namespace := os.Getenv("K8S_NAMESPACE")
	if namespace != "" {
		mp.pusher.Grouping("namespace", namespace)
	}
	podName := os.Getenv("K8S_POD")
	if podName != "" {
		mp.pusher.Grouping("pod", podName)
	}
	mp.pusher.Grouping("app", "monoskope")

	return mp, nil
}

// SetSuccessTime adds the successtime to the pushed metrics and sets the time to now
func (mp *metricsPublisher) SetSuccessTime() {
	mp.pusher.Collector(mp.successTime)
	mp.successTime.SetToCurrentTime()
}

// SetFailTime adds a failed timestamp to the pushed metrics and sets the time to now
func (mp *metricsPublisher) SetFailTime() {
	mp.pusher.Collector(mp.failedTime)
	mp.failedTime.SetToCurrentTime()
}

// Start sets a timestamp for start of a backup
func (mp *metricsPublisher) Start() {
	mp.start = time.Now()
}

// Finished sets the duration an backup has taken and the completion time to now
func (mp *metricsPublisher) Finished() {
	mp.duration.Set(time.Since(mp.start).Seconds())
	mp.completionTime.SetToCurrentTime()
}

// SetBytes sets the bytes written to backup
func (mp *metricsPublisher) SetBytes(bytes float64) {
	mp.bytes.Set(bytes)
}

// SetEventCount sets the events processed during an backup
func (mp *metricsPublisher) SetEventCount(eventCount float64) {
	mp.events.Set(eventCount)
}

// CloseAndPush sends the metrics to the prometheus push gateway
func (mp *metricsPublisher) CloseAndPush() error {
	// Add is used here rather than Push to not delete a previously pushed
	// success timestamp in case of a failure of this backup.
	if err := mp.pusher.Add(); err != nil {
		mp.log.Error(err, "Failed to push metrics.")
		return err
	}
	mp.log.Info("Metrics pushed successfully.")

	return nil
}
