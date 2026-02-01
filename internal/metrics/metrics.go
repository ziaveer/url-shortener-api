package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "url_shortener",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "url_shortener",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func Register() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
	)
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		path := r.URL.Path

		HttpRequestsTotal.
			WithLabelValues(r.Method, path, http.StatusText(rw.status)).
			Inc()

		HttpRequestDuration.
			WithLabelValues(r.Method, path).
			Observe(duration)
	})
}
