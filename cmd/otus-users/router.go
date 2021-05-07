package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	libhttp "github.com/d-ops-pro/otus-users/lib/http"
	"github.com/d-ops-pro/otus-users/metrics"
	"github.com/d-ops-pro/otus-users/users"
)

func SetupRouter(db *gorm.DB, logger *logrus.Entry) chi.Router {

	router := chi.NewRouter()
	{
		router.Use(middleware.RequestID)
		router.Use(libhttp.WithLogger(logger))
		router.Use(metrics.NewLatencyMiddleware(logger))
		router.Use(libhttp.RandomStatusMiddleware(0, 1000, logger))
		users.ConfigureRouter(router, db)
		metrics.ConfigureRouter(router)
	}

	return router
}
