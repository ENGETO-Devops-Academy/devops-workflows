package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsMiddleware struct {
	handler http.Handler

	// metrics
	requestsTotal    *prometheus.CounterVec
	requestsDuration *prometheus.GaugeVec
}

func NewMetricsMiddleware(h http.Handler) MetricsMiddleware {
	c := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of requests served",
	}, []string{"method"})

	d := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_duration",
		Help: "Duration in seconds of a single HTTP request",
	}, []string{"method"})

	return MetricsMiddleware{
		handler:          h,
		requestsTotal:    c,
		requestsDuration: d,
	}
}

func (m MetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()

	m.requestsTotal.WithLabelValues(
		r.Method,
	).Inc()

	m.handler.ServeHTTP(w, r)
	m.requestsDuration.WithLabelValues(
		r.Method,
	).Set(time.Now().Sub(t0).Seconds())
}
