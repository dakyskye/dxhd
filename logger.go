package main

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// init the logger before main function
func init() {
	formatter := &logrus.TextFormatter{DisableTimestamp: true, ForceQuote: true}

	if env, err := strconv.ParseBool(os.Getenv("DEBUG")); env && err == nil {
		logger.SetLevel(logrus.DebugLevel)
		logger.ReportCaller = true
		formatter.DisableTimestamp = false
		logger.SetFormatter(formatter)
	} else {
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(formatter)
	}
}
