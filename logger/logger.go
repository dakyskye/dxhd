package logger

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// Level is just custom int8 type for deciding what to log.
type Level int8

const (
	// Info Level logs basic information to stdout.
	Info Level = iota
	// Debug Level logs extra information as well as basic one to stdout
	Debug = iota
)

func init() {
	if env, err := strconv.ParseBool(os.Getenv("DEBUG")); env && err == nil {
		SetLevel(Debug)
	} else {
		SetLevel(Info)
	}
}

// SetLevel sets logging level, Info level by default.
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

// L just returns our logger.
func L() *logrus.Logger {
	return logger
}
