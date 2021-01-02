package cli

import (
	"os"

	"github.com/dakyskye/dxhd/logger"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = `master`

type CLI struct {
	app     *kingpin.Application
	Options Options
}

type Options struct {
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

func Init() (c CLI) {
	c.app = kingpin.New("dxhd", "daky's X11 Hotkey Daemon")

	c.app.Version(version)
	c.app.Author("Lasha Kanteladze <kanteladzelasha339@gmail.com> (https://github.com/dakyskye)")
	c.app.UsageTemplate(usageTemplate)

	c.app.HelpFlag.Short('h')

	c.Options.Kill = c.app.Flag("kill", "Kills all the running instances of dxhd.").Short('k').Action(c.kill).Bool()
	c.Options.Reload = c.app.Flag("reload", "Reloads all the running instances of dxhd.").Short('r').Action(c.reload).Bool()
	c.Options.DryRun = c.app.Flag("dry-run", "Does a dry run (prints parsed results and exits).").Short('d').Action(c.dryrun).Bool()
	c.Options.Background = c.app.Flag("background", "Runs dxhd as a background process.").Short('b').Action(c.background).Bool()
	c.Options.Interactive = c.app.Flag("interactive", "Runs dxhd interactively.").Short('i').Action(c.interactive).Bool()
	c.Options.Verbose = c.app.Flag("verbose", "Enables verbosity (logs everything).").Short('v').Action(c.verbose).Bool()
	c.Options.Config = c.app.Flag("config", "Parses the given config (defaults to dxhd.sh).").Short('c').Action(c.config).String()
	c.Options.Config = c.app.Flag("edit", "Starts an editor on the given config file (defaults to dxhd.sh).").Short('e').Action(c.edit).String()

	logger.L().Debugln("initialised the app")

	return
}

func (c *CLI) Parse() (err error) {
	_, err = c.app.Parse(os.Args[1:])
	return
}
