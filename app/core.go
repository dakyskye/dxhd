package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dakyskye/dxhd/logger"
)

type serverResponse string

const (
	reload  serverResponse = "reload"
	shutoff serverResponse = "shutoff"
)

func (c *CLI) Start() (err error) {
	logger.L().Debug("trying to start the server")
	for {
		go c.init()

		server := make(chan serverResponse, 1)
		go c.serve(server)

		command := <-server
		if command == shutoff {
			c.app.cancel()
		}
		logger.L().WithField("command", command).Debug("received a command")

		break // for now
	}
	return
}

func (c *CLI) init() {
	time.Sleep(time.Second * 5)
	c.app.cancel()
}

func (c *CLI) serve(res chan<- serverResponse) {
	logger.L().Debug("serving os signals")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGUSR1:
			fallthrough
		case syscall.SIGUSR2:
			res <- reload
		default:
			res <- shutoff
		}
	case <-c.app.ctx.Done():
		logger.L().WithError(c.app.ctx.Err()).Debug("context done")
		res <- shutoff
	}
}
