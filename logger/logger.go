package logger

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type Level int8

const (
	Info  Level = iota
	Debug       = iota
)

// init the logger before main function
func init() {
	if env, err := strconv.ParseBool(os.Getenv("DEBUG")); env && err == nil {
		SetLevel(Debug)
	} else {
		SetLevel(Info)
	}
}

// SetLevel sets logging level
func SetLevel(level Level) {
	switch level {
	case Info:
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, DisableColors: false, DisableLevelTruncation: true, ForceQuote: true})
	case Debug:
		logger.SetLevel(logrus.DebugLevel)
		logger.ReportCaller = true
		logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: false, DisableColors: false, DisableLevelTruncation: true, ForceQuote: true})
	}
}

// L just returns our internal logger
func L() *logrus.Logger {
	return logger
}
