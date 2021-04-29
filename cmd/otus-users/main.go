package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	_ "gorm.io/driver/postgres"
)

func main() {
	fs := flag.NewFlagSet("users", flag.ExitOnError)
	var (
		httpAddr = fs.String("http-addr", "", "HTTP listen address")
		dbURI    = fs.String("db-uri", "", "Database connection DSN")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	_ = fs.Parse(os.Args[1:])
	validateFS(fs, *httpAddr, *dbURI)

	db, err := ConnectDB(*dbURI)
	if err != nil {
		panic(err)
	}

	server := NewServer(*httpAddr, db)

	logrus.Infof("HTTP server started and listening on: '%s' port", *httpAddr)
	logrus.Fatal(server.ListenAndServe())
}
