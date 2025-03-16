package service

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func configureLogger() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	log.SetLevel(log.TraceLevel)
	if os.Getenv("ENV") == "production" { // TODO
		log.SetFormatter(&log.JSONFormatter{})
		log.SetReportCaller(true)
	}

	log.SetOutput(os.Stdout)
}
