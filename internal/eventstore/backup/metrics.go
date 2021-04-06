package backup

import (
	"errors"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type MetricsPublisher interface {
	Start()
	Finished()
	SetSuccessTime()
	SetEventCount(eventCount float64)
	SetBytes(sizeInBytes float64)
	CloseAndPush() error
}

type metricsPublisher struct {
	log            logger.Logger
	completionTime prometheus.Gauge
	successTime    prometheus.Gauge
	duration       prometheus.Gauge
	bytes          prometheus.Gauge
	events         prometheus.Gauge
	pusher         *push.Pusher
	start          time.Time
}

// NewMetricsPublisher creates a new backup.MetricsPublisher.
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
	registry.MustRegister(mp.completionTime, mp.duration, mp.bytes, prometheus.NewGoCollector())
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

// Add successTime to pusher only in case of success.
func (mp *metricsPublisher) SetSuccessTime() {
	mp.pusher.Collector(mp.successTime)
	mp.successTime.SetToCurrentTime()
}

// Set the start
func (mp *metricsPublisher) Start() {
	mp.start = time.Now()
}

func (mp *metricsPublisher) Finished() {
	mp.duration.Set(time.Since(mp.start).Seconds())
	mp.completionTime.SetToCurrentTime()
}

func (mp *metricsPublisher) SetBytes(bytes float64) {
	mp.bytes.Set(bytes)
}

func (mp *metricsPublisher) SetEventCount(eventCount float64) {
	mp.events.Set(eventCount)
}

func (mp *metricsPublisher) CloseAndPush() error {
	// Add is used here rather than Push to not delete a previously pushed
	// success timestamp in case of a failure of this backup.
	if err := mp.pusher.Add(); err != nil {
		mp.log.Error(err, "Failed to push metrics.")
		return err
	}
	mp.log.Info("Metrics pushed sucessfully.")

	return nil
}
