package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type requestStartKey struct{}

func ConfigureRouter(r chi.Router) {
	r.Route("/metrics", func(r chi.Router) {
		r.Use(beforeMiddleware)
		r.Handle("/", promhttp.Handler())
	})
}

func beforeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestStartKey{}, time.Now())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
