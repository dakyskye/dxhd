package options

import (
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Help      bool
	Kill      bool
	Reload    bool
	Version   bool
	DryRun    bool
	ParseTime bool
	Config    *string
}

var Options = `
  -h, --help              Prints this help message
  -c, --config            Reads the config from custom path
  -d, --dry-run           Prints bindings and their actions and exits
  -k, --kill	          Gracefully kills every running instances of dxhd
  -p, --parse-time        Prints how much time parsing a config took
  -r, --reload	          Reloads every running instances of dxhd
  -v, --version	          Prints current version of program`

func Parse() (result Result, err error) {
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
				result.Help = true
			case opt == "kill":
				result.Kill = true
			case opt == "reload":
				result.Reload = true
			case opt == "dry-run":
				result.DryRun = true
			case opt == "parse-time":
				result.ParseTime = true
			case opt == "config":
				if in == len(osArgs)-1 {
					break
				}
				if strings.HasPrefix(osArgs[in+1], "--") || strings.HasPrefix(osArgs[in+1], "-") {
					continue
				}
				result.Config = &osArgs[in+1]
				skip = true
			case strings.HasPrefix(opt, "config="):
				result.Config = new(string)
				*result.Config = strings.TrimPrefix(opt, "config=")
			case opt == "version":
				result.Version = true
			default:
				err = fmt.Errorf("%s is not a valid option", err)
				return
			}
		} else if strings.HasPrefix(osArg, "-") {
			for _, r := range osArg[1:] {
				switch string(r) {
				case "h":
					result.Help = true
				case "k":
					result.Kill = true
				case "r":
					result.Reload = true
				case "v":
					result.Version = true
				case "d":
					result.DryRun = true
				case "p":
					result.ParseTime = true
				case "c":
					if in == len(osArgs)-1 {
						break toplevel
					}
					if strings.HasPrefix(osArgs[in+1], "--") || strings.HasPrefix(osArgs[in+1], "-") {
						continue
					}
					result.Config = new(string)
					result.Config = &osArgs[in+1]
					skip = true
				default:
					err = fmt.Errorf("%s in %s is not a valid option\n", string(r), osArg)
					return
				}
			}
		}
	}

	return
}
