package main

import (
	"github.com/dakyskye/dxhd/logger"

	"github.com/dakyskye/dxhd/cli"
)

func main() {
	app := cli.Init()
	err := app.Parse()
	if err != nil {
		logger.L().WithError(err).Fatalln("can not parse commndline arguments")
	}
}
