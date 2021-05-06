package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewServer(listen string, db *gorm.DB, logger *logrus.Entry) *http.Server {
	router := SetupRouter(db, logger)

	return &http.Server{
		Handler:      router,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
