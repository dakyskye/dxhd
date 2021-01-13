package app

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dakyskye/dxhd/logger"

	"gopkg.in/alecthomas/kingpin.v2"
)

func (a *App) kill(_ *kingpin.ParseContext) (err error) {
	execName, err := os.Executable()
	if err != nil {
		return
	}

	err = exec.Command("pkill", "-INT", "-x", filepath.Base(execName)).Start() //nolint:gosec

	return
}

func (a *App) reload(_ *kingpin.ParseContext) (err error) {
	execName, err := os.Executable()
	if err != nil {
		return
	}

	err = exec.Command("pkill", "-USR1", "-x", filepath.Base(execName)).Start() //nolint:gosec

	return
}

func (a *App) dryrun(_ *kingpin.ParseContext) (err error) {
	return
}

func (a *App) background(_ *kingpin.ParseContext) (err error) {
	return
}

func (a *App) interactive(_ *kingpin.ParseContext) (err error) {
	return
}

func (a *App) verbose(_ *kingpin.ParseContext) (err error) {
	logger.SetLevel(logger.Debug)

	return
}

func (a *App) config(_ *kingpin.ParseContext) (err error) {
	return
}

func (a *App) edit(_ *kingpin.ParseContext) (err error) {
	return
}
