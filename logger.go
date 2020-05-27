package main

import (
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	var config zap.Config

	doExtra := func(envVar string) bool {
		if env, err := strconv.ParseBool(os.Getenv(envVar)); env && err == nil {
			return true
		}
		return false
	}

	if doExtra("DEV") {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	if doExtra("DEBUG") {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}
