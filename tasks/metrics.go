package tasks

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

var runningGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "gowebcore_task_running",
		Help: "Number of running background tasks.",
	},
	[]string{"name"},
)

// RegisterMetrics registers the gauge with the default registry.
func RegisterMetrics() { prometheus.MustRegister(runningGauge) }

// Wrap adds bookkeeping to any task.
func Wrap(name string, t Task) Task {
	return func(ctx context.Context) error {
		runningGauge.WithLabelValues(name).Inc()
		defer runningGauge.WithLabelValues(name).Dec()
		return t(ctx)
	}
}
