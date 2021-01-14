package app

import (
	"context"
	"fmt"
	"os"
	"sync"

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

var (
	app  App
	once sync.Once
)

// Init initialises a new app.
func Init() (App, error) {
	var err error

	logger.L().Debugln("app initialisation requested")
	once.Do(func() {
		configFile, e := config.GetConfigFile()
		if e != nil {
			err = e
			return
		}

		err = config.CreateDefaultConfig()
		if err != nil {
			err = fmt.Errorf("error creating default config: %w", err)
			return
		}

		app.execName, err = os.Executable()
		if err != nil {
			return
		}

		app.opts.edit = configFile
		app.opts.config = configFile

		app.ctx, app.cancel = context.WithCancel(context.Background())

		app.cli = kingpin.New("dxhd", "daky's X11 Hotkey Daemon")

		app.cli.Version(version)
		app.cli.Author("Lasha Kanteladze <kanteladzelasha339@gmail.com> (https://github.com/dakyskye)")
		app.cli.UsageTemplate(usageTemplate)

		app.cli.HelpFlag.Short('h')

		app.cli.Flag("kill", "Kills all the running instances of dxhd.").Short('k').Action(app.kill).Bool()
		app.cli.Flag("reload", "Reloads all the running instances of dxhd.").Short('r').Action(app.reload).Bool()
		app.cli.Flag("dry-run", "Does a dry run (prints parsed results and exits).").Short('d').Action(app.dryrun).Bool()
		app.cli.Flag("background", "Runs dxhd as a background process.").Short('b').Action(app.background).Bool()
		app.cli.Flag("interactive", "Runs dxhd interactively.").Short('i').Action(app.interactive).Bool()
		app.cli.Flag("verbose", "Enables verbosity (logs everything).").Short('v').Action(app.verbose).Bool()
		app.cli.Flag("config", "Parses the given config (defaults to dxhd.sh).").Short('c').Default(configFile).StringVar(&app.opts.config)
		app.cli.Flag("edit", "Starts an editor on the given config file (defaults to dxhd.sh).").Short('e').Default(configFile).Action(app.edit).StringVar(&app.opts.edit)

		logger.L().Debug("initialised the app")
	})

	return app, err
}

func (a *App) Parse() (err error) {
	_, err = a.cli.Parse(os.Args[1:])

	return
}
