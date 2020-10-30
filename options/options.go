package options

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Options struct {
	Help        bool
	Kill        bool
	Reload      bool
	Version     bool
	DryRun      bool
	ParseTime   bool
	Background  bool
	Interactive bool
	Config      *string
	Edit        *string
}

var OptionsToPrint = `
  -h, --help              Prints this help message
  -c, --config [path]     Reads the config from custom path
  -d, --dry-run           Prints bindings and their commands and exits
  -b, --background        Runs dxhd in the background
  -k, --kill              Gracefully kills every running instances of dxhd
  -p, --parse-time        Prints how much time parsing a config took
  -r, --reload            Reloads every running instances of dxhd
  -v, --version           Prints current version of program
  -e, --edit [file]       Shortcut to edit a file in dxhd's config folder. Opens dxhd.sh if file is empty
  -i, --interactive       Opens a temporary file for temporary bindings to run`

func Parse() (opts Options, err error) {
	osArgs := os.Args[1:]

	skip := false

	readNextArg := func(index int, optional bool) (*string, error) {
		if index == len(osArgs)-1 || strings.HasPrefix(osArgs[index+1], "-") {
			if optional {
				return nil, nil
			}
			return nil, errors.New("not enough arguments passed to flags")
		}
		return &osArgs[index+1], nil
	}

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
			case opt == "background":
				opts.Background = true
				os.Args[1+in] = "" // remove --background
			case opt == "config":
				opts.Config, err = readNextArg(in, false)
				if err != nil {
					break
				}
				skip = true
			case strings.HasPrefix(opt, "config="):
				opts.Config = new(string)
				*opts.Config = strings.TrimPrefix(opt, "config=")
			case opt == "version":
				opts.Version = true
			case opt == "interactive":
				opts.Interactive = true
			case opt == "edit":
				opts.Edit, err = readNextArg(in, true)
				if err != nil {
					continue
				}
				if opts.Edit == nil {
					// Default value if --edit is given but no argument
					// This allows to differentiate between no edit flag
					// and wanting to edit the default config
					opts.Edit = new(string)
					*opts.Edit = ""
				} else {
					skip = true
				}
			case strings.HasPrefix(opt, "edit="):
				opts.Edit = new(string)
				*opts.Edit = strings.TrimPrefix(opt, "edit=")
			default:
				err = fmt.Errorf("%s is not a valid option", osArg)
				return
			}
		} else if strings.HasPrefix(osArg, "-") {
			for i, r := range osArg[1:] {
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
				case "b":
					opts.Background = true
					if len(osArg) == 2 { // means b is the only flag with dash prefix
						os.Args[1+in] = ""
					} else { // means it's a flag combo
						os.Args[1+in] = os.Args[1+in][0:1+i] + os.Args[1+in][2+i:] // remove b from the combo
					}
				case "c":
					opts.Config, err = readNextArg(in, false)
					if err != nil {
						continue
					}
					skip = true
				case "e":
					opts.Edit, err = readNextArg(in, true)
					if err != nil {
						continue
					}
					if opts.Edit == nil {
						// Default value if -e is given but no argument
						// This allows to differentiate between no edit flag
						// and wanting to edit the default config
						opts.Edit = new(string)
						*opts.Edit = ""
					} else {
						skip = true
					}
				case "i":
					opts.Interactive = true
				default:
					err = fmt.Errorf("%s in %s is not a valid option", string(r), osArg)
					return
				}
			}
		}
	}

	return
}
