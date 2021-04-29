package main

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

func NewServer(listen string, db *gorm.DB) *http.Server {
	router := SetupRouter(db)

	return &http.Server{
		Handler:      router,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
