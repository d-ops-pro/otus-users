package main

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/d-ops-pro/otus-users/users"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter().
		StrictSlash(true)
	{
		users.ConfigureRouter(router, db)
	}

	return router
}
