package main

import (
	"fmt"
	"os"
	"strings"
)

type options struct {
	help      bool
	kill      bool
	reload    bool
	version   bool
	dryRun    bool
	parseTime bool
	config    *string
}

var opts options

var flags = `
  -c, --config            Reads the config from custom path
  -d, --dry-run           Prints bindings and their actions and exits
  -k, --kill	          Gracefully kills every running instances of dxhd
  -p, --parse-time        Prints how much time parsing a config took
  -r, --reload	          Reloads every running instances of dxhd
  -v, --version	          Prints current version of program`

func init() {
	osArgs := os.Args[1:]

	skip := false
	for in, osArg := range osArgs {
		if skip {
			skip = false
			continue
		}
		if strings.HasPrefix(osArg, "--") {
			switch opt := osArg[2:]; {
			case opt == "help":
				opts.help = true
			case opt == "kill":
				opts.kill = true
			case opt == "reload":
				opts.reload = true
			case opt == "version":
				opts.version = true
			case opt == "dry-run":
				opts.dryRun = true
			case opt == "parse-time":
				opts.parseTime = true
			case opt == "config":
				if in == len(osArgs)-1 {
					break
				}
				if strings.HasPrefix(osArgs[in+1], "--") || strings.HasPrefix(osArgs[in+1], "-") {
					continue
				}
				*opts.config = osArgs[in+1]
				skip = true
			case strings.HasPrefix(opt, "config="):
				*opts.config = strings.TrimPrefix(opt, "config=")
			default:
				fmt.Println(fmt.Sprintf("%s is not a valid option", opt))
			}
		} else if strings.HasPrefix(osArg, "-") {
			for _, r := range osArg[1:] {
				switch string(r) {
				case "h":
					opts.help = true
				case "k":
					opts.kill = true
				case "r":
					opts.reload = true
				case "v":
					opts.version = true
				case "d":
					opts.dryRun = true
				case "p":
					opts.parseTime = true
				default:
					fmt.Println(fmt.Sprintf("%s in %s is not a valid option", string(r), osArg))
				}
			}
		}
	}

	usage = fmt.Sprintf(usage, version, flags)
}
