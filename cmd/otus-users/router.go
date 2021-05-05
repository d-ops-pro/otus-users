package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	libhttp "github.com/d-ops-pro/otus-users/lib/http"
	"github.com/d-ops-pro/otus-users/metrics"
	"github.com/d-ops-pro/otus-users/users"
)

func SetupRouter(db *gorm.DB) chi.Router {
	router := chi.NewRouter()
	{
		router.Use(middleware.RequestID)
		router.Use(libhttp.WithLogger)
		users.ConfigureRouter(router, db)
		metrics.ConfigureRouter(router)
	}

	return router
}
