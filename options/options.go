package options

import (
	"fmt"
	"os"
	"strings"
)

type Options struct {
	Help      bool
	Kill      bool
	Reload    bool
	Version   bool
	DryRun    bool
	ParseTime bool
	Config    *string
}

var OptionsToPrint = `
  -h, --help              Prints this help message
  -c, --config            Reads the config from custom path
  -d, --dry-run           Prints bindings and their commands and exits
  -k, --kill	          Gracefully kills every running instances of dxhd
  -p, --parse-time        Prints how much time parsing a config took
  -r, --reload	          Reloads every running instances of dxhd
  -v, --version	          Prints current version of program`

func Parse() (opts Options, err error) {
	osArgs := os.Args[1:]

	skip := false

toplevel:
	for in, osArg := range osArgs {
		if skip {
			skip = false
			continue
		}
		if strings.HasPrefix(osArg, "--") {
			switch opt := osArg[2:]; {
			case opt == "help":
				opts.Help = true
			case opt == "kill":
				opts.Kill = true
			case opt == "reload":
				opts.Reload = true
			case opt == "dry-run":
				opts.DryRun = true
			case opt == "parse-time":
				opts.ParseTime = true
			case opt == "config":
				if in == len(osArgs)-1 {
					break
				}
				if strings.HasPrefix(osArgs[in+1], "--") || strings.HasPrefix(osArgs[in+1], "-") {
					continue
				}
				opts.Config = &osArgs[in+1]
				skip = true
			case strings.HasPrefix(opt, "config="):
				opts.Config = new(string)
				*opts.Config = strings.TrimPrefix(opt, "config=")
			case opt == "version":
				opts.Version = true
			default:
				err = fmt.Errorf("%s is not a valid option", err)
				return
			}
		} else if strings.HasPrefix(osArg, "-") {
			for _, r := range osArg[1:] {
				switch string(r) {
				case "h":
					opts.Help = true
				case "k":
					opts.Kill = true
				case "r":
					opts.Reload = true
				case "v":
					opts.Version = true
				case "d":
					opts.DryRun = true
				case "p":
					opts.ParseTime = true
				case "c":
					if in == len(osArgs)-1 {
						break toplevel
					}
					if strings.HasPrefix(osArgs[in+1], "--") || strings.HasPrefix(osArgs[in+1], "-") {
						continue
					}
					opts.Config = new(string)
					opts.Config = &osArgs[in+1]
					skip = true
				default:
					err = fmt.Errorf("%s in %s is not a valid option", string(r), osArg)
					return
				}
			}
		}
	}

	return
}
