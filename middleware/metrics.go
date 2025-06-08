package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
			// buckets in seconds
			Buckets: []float64{0.01, 0.1, 0.3, 1.2, 5},
		},
		[]string{"path", "method", "status"},
	)
)

// MustRegisterMetrics registers the histogram with Prometheus's default registry.
// Call this once in init() of main package.
func MustRegisterMetrics() {
	prometheus.MustRegister(reqDuration)
}

// Metrics records duration & status.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()

		reqDuration.WithLabelValues(
			r.URL.Path, r.Method, strconv.Itoa(ww.Status()),
		).Observe(duration)
	})
}

// RegisterDefaultPrometheus registers Go GC & process collectors.
// Call once in main().
func RegisterDefaultPrometheus() {
	prometheus.MustRegister(prometheus.NewGoCollector())
	prometheus.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
}
