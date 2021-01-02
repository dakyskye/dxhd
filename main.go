package main

import (
	"github.com/dakyskye/dxhd/cli"
	"github.com/dakyskye/dxhd/logger"
)

func main() {
	app, err := cli.Init()
	if err != nil {
		logger.L().WithError(err).Fatalln("can not initialise the app")
	}
	err = app.Parse()
	if err != nil {
		logger.L().WithError(err).Fatalln("something went wrong")
	}
}
