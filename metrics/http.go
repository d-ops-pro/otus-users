package metrics

import (
	"fmt"
	"net/http"
	"time"

	libhttp "github.com/d-ops-pro/otus-users/lib/http"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func ConfigureRouter(r chi.Router) {
	r.Route("/metrics", func(r chi.Router) {
		r.Handle("/", promhttp.Handler())
	})
}

func NewLatencyMiddleware(logger *logrus.Entry) func(http.Handler) http.Handler {
	reqLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "users_request_latency_seconds",
		Help: "Users Service request latency",
	}, []string{"method", "endpoint"})

	reqCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "users_request_count",
		Help: "Users Service request count",
	}, []string{"method", "endpoint", "status"})

	prometheus.MustRegister(reqLatency, reqCount)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := libhttp.NewResponseRecorder(w)

			next.ServeHTTP(recorder, r)

			reqLatency.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
			reqCount.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", recorder.Status()))
		})
	}
}
