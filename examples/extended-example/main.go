package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	logger "github.com/awslabs/fluent-golang-io-writer"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	// constructs tag from the message
	var tagConstructor = func(p []byte, tagPrefix string) string {
		msg := string(p)
		if strings.Contains(msg, "debug") {
			return tagPrefix + ".debug"
		} else if strings.Contains(msg, "error") {
			return tagPrefix + ".error"
		} else if strings.Contains(msg, "info") {
			return tagPrefix + ".info"
		}

		return tagPrefix + ".logs"
	}

	// Converts the log to JSON if its a key value pair that's tab separated
	// This is needed because some log events are not JSON
	var convertToJSON = func(p []byte) (map[string]interface{}, error) {
		pair := strings.Split(string(p), "\t")
		fmt.Println(pair)
		if len(pair) >= 2 {
			msg := make(map[string]interface{})
			msg[pair[0]] = pair[1]
			return msg, nil
		}
		return nil, nil
	}

	fluentLogger, err := logger.NewFluentWriterWithTagConstructor(fluent.Config{}, "app", tagConstructor)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fluentLogger.AddSecondaryJSONConverter(convertToJSON)

	defer fluentLogger.Close()

	logrus.SetOutput(fluentLogger)

	// we can also use standard library logger
	log.SetOutput(fluentLogger)

	for true {
		// Use logrus normally
		logrus.WithFields(logrus.Fields{
			"animal": "walrus",
		}).Info("A walrus appears")

		logrus.WithFields(logrus.Fields{
			"animal": "tiger",
		}).Warn("A tiger appears")

		logrus.WithFields(logrus.Fields{
			"animal": "bird",
		}).Debug("A bird appears")

		// standard library logger
		log.Println("requestID\t123455")

		log.Println("current user count\t5")

		time.Sleep(1 * time.Second)
	}
}
