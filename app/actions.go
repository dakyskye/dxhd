package app

import (
	"github.com/dakyskye/dxhd/app/utils"
	"github.com/dakyskye/dxhd/logger"

	"gopkg.in/alecthomas/kingpin.v2"
)

func (a *App) kill(_ *kingpin.ParseContext) error {
	return utils.Cmd("pkill", "-INT", "-x", a.execName).Quick()
}

func (a *App) reload(_ *kingpin.ParseContext) error {
	return utils.Cmd("pkill", "-USR1", "-x", a.execName).Quick()
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

func (a *App) edit(ctx *kingpin.ParseContext) (err error) {
	editor, err := utils.FindEditor()
	if err != nil {
		return
	}
	err = utils.Cmd(editor, *ctx.Elements[0].Value).Run()
	if err == nil {
		a.cancel()
	}

	return
}
