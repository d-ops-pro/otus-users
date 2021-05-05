package libhttp

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type LoggerContextKey struct{}

func GetLoggerFromContext(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(LoggerContextKey{})
	if logger == nil {
		panic("logger is not in the context")
	}

	return logger.(*logrus.Entry)
}

func GetLogger(r *http.Request) *logrus.Entry {
	return GetLoggerFromContext(r.Context())
}

func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logrus.NewEntry(logrus.New())
		requestID := middleware.GetReqID(ctx)

		logger = logger.WithField("request_id", requestID)

		ctx = context.WithValue(ctx, LoggerContextKey{}, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
