package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dakyskye/dxhd/logger"
)

type App struct {
	execName string
	ctx      context.Context
	cancel   context.CancelFunc
	cli      *kingpin.Application
	opts     options
}

type serverResponse string

const (
	reload  serverResponse = "reload"
	shutoff serverResponse = "shutoff"
)

func (a *App) Start() (err error) {
	logger.L().Debug("trying to start the server")

	for {
		go a.init()

		server := make(chan serverResponse, 1)
		go a.serve(server)

		command := <-server
		if command == shutoff {
			a.cancel()
		}

		logger.L().WithField("command", command).Debug("received a command")

		break
	}

	return
}

func (a *App) init() (err error) {
	return
}

func (a *App) serve(res chan<- serverResponse) {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	logger.L().Debug("serving os signals")

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
	case <-a.ctx.Done():
		logger.L().WithError(a.ctx.Err()).Debug("main app context done")
		res <- shutoff
	}
}
