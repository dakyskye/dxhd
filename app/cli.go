package app

import (
	"context"
	"os"

	"github.com/dakyskye/dxhd/config"

	"github.com/dakyskye/dxhd/logger"
	"gopkg.in/alecthomas/kingpin.v2"
)

type options struct {
	config string
	edit   string
}

var version = `master`

// we use custom usage template.
var usageTemplate = `NAME
  {{.App.Name}} - {{.App.Help}}
VERSION
  {{.App.Version}}
SYNOPSIS
  {{.App.Name}} [<FLAGS>]
DESCRIPTION
  will be written
FLAGS
{{.Context.Flags|FlagsToTwoColumns|FormatTwoColumns}}\
EXAMPLE CONFIG
  will be written
AUTHOR
  {{.App.Author}}
BUGS
  Report a bug here if you find one - https://github.com/dakyskye/dxhd/issues
`

// Init initialises a new app.
func Init() (a App, err error) {
	configFile, err := config.GetConfigFile()
	if err != nil {
		return
	}

	a.execName, err = os.Executable()
	if err != nil {
		return
	}

	a.opts.edit = configFile
	a.opts.config = configFile

	a.ctx, a.cancel = context.WithCancel(context.Background())

	a.cli = kingpin.New("dxhd", "daky's X11 Hotkey Daemon")

	a.cli.Version(version)
	a.cli.Author("Lasha Kanteladze <kanteladzelasha339@gmail.com> (https://github.com/dakyskye)")
	a.cli.UsageTemplate(usageTemplate)

	a.cli.HelpFlag.Short('h')

	a.cli.Flag("kill", "Kills all the running instances of dxhd.").Short('k').Action(a.kill).Bool()
	a.cli.Flag("reload", "Reloads all the running instances of dxhd.").Short('r').Action(a.reload).Bool()
	a.cli.Flag("dry-run", "Does a dry run (prints parsed results and exits).").Short('d').Action(a.dryrun).Bool()
	a.cli.Flag("background", "Runs dxhd as a background process.").Short('b').Action(a.background).Bool()
	a.cli.Flag("interactive", "Runs dxhd interactively.").Short('i').Action(a.interactive).Bool()
	a.cli.Flag("verbose", "Enables verbosity (logs everything).").Short('v').Action(a.verbose).Bool()
	a.cli.Flag("config", "Parses the given config (defaults to dxhd.sh).").Short('c').Default(configFile).StringVar(&a.opts.config)
	a.cli.Flag("edit", "Starts an editor on the given config file (defaults to dxhd.sh).").Short('e').Default(configFile).Action(a.edit).StringVar(&a.opts.edit)

	logger.L().WithField("file", configFile).Debug("set default file for -c and -e flags")
	logger.L().Debug("initialised the app")

	return
}

func (a *App) Parse() (err error) {
	_, err = a.cli.Parse(os.Args[1:])

	return
}
