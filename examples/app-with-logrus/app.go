package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	logger "github.com/awslabs/fluent-golang-io-writer"

	"github.com/fluent/fluent-logger-golang/fluent"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	// create a FluentWriter instance
	port, err := strconv.Atoi(os.Getenv("FLUENT_PORT"))
	if err != nil {
		// process error
	}
	config := fluent.Config{
		FluentPort: port,
		FluentHost: os.Getenv("FLUENT_HOST"),
	}
	fluentLogger, err := logger.NewFluentWriter(config, "app", []string{"level"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set logrus to use it
	log.SetOutput(fluentLogger)

	for true {
		log.WithFields(log.Fields{
			"requestID": "45234523",
			"path":      "/",
		}).Info("Got a request")

		log.WithFields(log.Fields{
			"requestID": "546745643",
			"path":      "/tardis",
			"user":      "TheMaster",
		}).Warn("Access denied")

		log.WithFields(log.Fields{
			"requestID": "546745643",
			"path":      "/tardis",
			"user":      "TheDoctor",
		}).Debug("Admin access")
		time.Sleep(100 * time.Millisecond)
	}
}
