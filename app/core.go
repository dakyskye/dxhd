package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dakyskye/dxhd/logger"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	cli    *kingpin.Application
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

		break // for now
	}
	return
}

func (a *App) init() {
	time.Sleep(time.Second * 5)
	a.cancel()
}

func (a *App) serve(res chan<- serverResponse) {
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
	case <-a.ctx.Done():
		logger.L().WithError(a.ctx.Err()).Debug("context done")
		res <- shutoff
	}
}
