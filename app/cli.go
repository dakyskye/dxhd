package app

import (
	"context"
	"os"

	"github.com/dakyskye/dxhd/config"

	"github.com/dakyskye/dxhd/logger"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = `master`

type CLI struct {
	app     app
	options options
}

type app struct {
	ctx    context.Context
	cancel context.CancelFunc
	cli    *kingpin.Application
}

type options struct {
	Kill        *bool
	Reload      *bool
	DryRun      *bool
	Background  *bool
	Interactive *bool
	Verbose     *bool
	Config      *string
	Edit        *string
}

// we use custom usage template
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

func Init() (c CLI, err error) {
	configFile, err := config.GetConfigFile()
	if err != nil {
		return
	}

	c.app.ctx, c.app.cancel = context.WithCancel(context.Background())

	c.app.cli = kingpin.New("dxhd", "daky's X11 Hotkey Daemon")

	c.app.cli.Version(version)
	c.app.cli.Author("Lasha Kanteladze <kanteladzelasha339@gmail.com> (https://github.com/dakyskye)")
	c.app.cli.UsageTemplate(usageTemplate)

	c.app.cli.HelpFlag.Short('h')

	c.options.Kill = c.app.cli.Flag("kill", "Kills all the running instances of dxhd.").Short('k').Action(c.kill).Bool()
	c.options.Reload = c.app.cli.Flag("reload", "Reloads all the running instances of dxhd.").Short('r').Action(c.reload).Bool()
	c.options.DryRun = c.app.cli.Flag("dry-run", "Does a dry run (prints parsed results and exits).").Short('d').Action(c.dryrun).Bool()
	c.options.Background = c.app.cli.Flag("background", "Runs dxhd as a background process.").Short('b').Action(c.background).Bool()
	c.options.Interactive = c.app.cli.Flag("interactive", "Runs dxhd interactively.").Short('i').Action(c.interactive).Bool()
	c.options.Verbose = c.app.cli.Flag("verbose", "Enables verbosity (logs everything).").Short('v').Action(c.verbose).Bool()
	c.options.Config = c.app.cli.Flag("config", "Parses the given config (defaults to dxhd.sh).").Short('c').Default(configFile).Action(c.config).String()
	c.options.Config = c.app.cli.Flag("edit", "Starts an editor on the given config file (defaults to dxhd.sh).").Short('e').Default(configFile).Action(c.edit).String()

	logger.L().WithField("file", configFile).Debug("set default file for -c and -e flags")
	logger.L().Debug("initialised the app")

	return
}

func (c *CLI) Parse() (err error) {
	_, err = c.app.cli.Parse(os.Args[1:])
	return
}
