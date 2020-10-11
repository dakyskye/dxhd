package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/dakyskye/dxhd/config"
	"github.com/dakyskye/dxhd/listener"
	"github.com/dakyskye/dxhd/logger"
	"github.com/dakyskye/dxhd/options"
	"github.com/dakyskye/dxhd/parser"
	"github.com/sirupsen/logrus"
)

var usage = `NAME
  dxhd - daky's X11 Hotkey Daemon
VERSION
  %s
SYNOPSIS
  dxhd [OPTIONS]
DESCRIPTION
  dxhd is an easy-to-use X11 hotkey daemon, written in Go programming language, and inspired by sxhkd.
  One of the biggest advantages of dxhd is that you can write your configs in different languages,
  for example: sh, bash, ksh, zsh, Python, Perl.
  A config file is meant to have quite easy layout:
    First line starting with #! is treated as a shebang
    Lines having ##+ prefix are treated as comments (get ignored)
    Lines having one # and then a keybinding are parsed as keybindings
    Lines under a keybinding are executed when keybinding is triggered
OPTIONS%s
EXAMPLE CONFIG
  #!/bin/sh
  ## restart i3
  # super + shift + r
  i3-msg -t command restart
  ## switch to workspace 1-10
  # super + @{1-9,0}
  i3-msg -t command workspace {1-9,10}
  ## switch to workspace 11-20
  # super + ctrl + {1-9,0}
  i3-msg -t command workspace {11-19,20}
  ## switch to next/prev workspace
  # super + mouse{4,5}
  i3-msg -t command workspace {next,prev}
BUGS
  report bugs here, if you encounter one - https://github.com/dakyskye/dxhd/issues
AUTHOR
  Lasha Kanteladze <kanteladzelasha339@gmail.com>`

var version = `master`

func main() {
	if runtime.GOOS != "linux" {
		logger.L().Fatalln("dxhd is only supported on linux")
	}

	opts, err := options.Parse()
	if err != nil {
		logger.L().Fatalln(err)
	}

	usage = fmt.Sprintf(usage, version, options.OptionsToPrint)

	exit := false

	if opts.Help {
		fmt.Println(usage)
		fmt.Println()
		exit = true
	}

	if opts.Edit != nil {
		editor_envvar := os.Getenv("EDITOR")
		if editor_envvar == "" {
			logger.L().Fatal("The $EDITOR environment variable is not set!")
		} else {
			ed, err := exec.Command("which", editor_envvar).Output()
			editor := string(ed)
			editor = editor[:len(editor)-1] // Remove the trailing newline char as output from `which`
			if err != nil {
				logger.L().WithField("$EDITOR", editor_envvar).Fatal("Value in $EDITOR doesn't translate to an executable.")
			}
			configDir, _ := os.UserConfigDir()
			if *opts.Edit == "" {
				*opts.Edit = "dxhd.sh"
			}
			path := filepath.Join(configDir, "dxhd", *opts.Edit)
			err = syscall.Exec(editor, []string{editor, path}, os.Environ())
		}
		exit = true
	}

	if opts.Version && !opts.Help {
		fmt.Println("you are using dxhd, version " + version)
		fmt.Println()
		exit = true
	}

	var (
		configFilePath string
		validPath      bool
	)

	if opts.Config != nil {
		if validPath, err = config.IsPathToConfigValid(*opts.Config); !(err == nil && validPath) {
			logger.L().WithFields(logrus.Fields{
				"path":  *opts.Config,
				"valid": validPath,
			}).WithError(err).Fatal("path to the config is not valid")
		}
		configFilePath = *opts.Config
	} else {
		configFilePath, _, err = config.GetDefaultConfigPath()
		if err != nil {
			logger.L().WithError(err).Fatal("can not get config path")
		}

		if validPath, err = config.IsPathToConfigValid(configFilePath); !(err == nil && validPath) {
			if os.IsNotExist(err) {
				err = config.CreateDefaultConfig()
				if err != nil {
					logger.L().WithField("path", configFilePath).Fatal("can not create default config")
				}
			} else {
				logger.L().WithFields(logrus.Fields{"path": configFilePath, "valid": validPath}).WithError(err).Fatal("path to the config is not valid")
			}
		}
	}

	var (
		data      []parser.FileData
		shell     string
		globals   string
		startTime time.Time
	)

	if opts.ParseTime {
		startTime = time.Now()
	}

	shell, globals, err = parser.Parse(configFilePath, &data)
	if err != nil {
		logger.L().WithField("file", configFilePath).WithError(err).Fatal("failed to parse config")
	}

	if opts.ParseTime {
		since := time.Since(startTime)
		timeTaken := fmt.Sprintf("%.0fs%dms%dµs",
			since.Seconds(),
			since.Milliseconds(),
			since.Microseconds(),
		)
		fmt.Printf("it took %s to parse the config\n", timeTaken)
		fmt.Printf("%d parsed keybindins (including replicated variants and ranges)\n", len(data))
		exit = true
	}

	if opts.DryRun {
		fmt.Println("dxhd dry run")
		for _, d := range data {
			fmt.Println("binding: " + d.OriginalBinding)
			fmt.Println("command:")
			fmt.Println(d.Command.String())
		}
		fmt.Println()
		exit = true
	}

	if opts.Kill || opts.Reload {
		execName, err := os.Executable()
		if err != nil {
			logger.L().WithError(err).Fatal("can not get executable")
		}

		if opts.Kill {
			err = exec.Command("pkill", "-INT", "-x", filepath.Base(execName)).Start()
		} else {
			err = exec.Command("pkill", "-USR1", "-x", filepath.Base(execName)).Start()
		}

		if err != nil {
			if opts.Kill {
				log.Println("can not kill dxhd instances:")
				log.Fatalln(err)
			} else {
				log.Println("can not reload dxhd instances:")
				log.Fatalln(err)
			}
		}

		if opts.Kill {
			fmt.Println("killing every running instances of dxhd")
		} else {
			fmt.Println("reloading every running instances of dxhd")
		}

		exit = true
	}

	if exit {
		os.Exit(0)
	}

	logger.L().WithFields(logrus.Fields{"version": version, "path": configFilePath}).Debug("starting dxhd")

	// catch these signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	// errors channel
	errs := make(chan error)

	// infinite loop - if user sends USR signal, reload configration (so, continue loop), otherwise, exit
toplevel:
	for {
		if len(data) == 0 {
			shell, globals, err = parser.Parse(configFilePath, &data)
			if err != nil {
				logger.L().WithField("file", configFilePath).WithError(err).Fatal("failed to parse config")
			}
		}

		X, err := xgbutil.NewConn()
		if err != nil {
			logger.L().WithError(err).Fatal("can not open connection to Xorg")
		}

		keybind.Initialize(X)
		mousebind.Initialize(X)

		for _, d := range data {
			err = listener.ListenKeybinding(X, errs, d.EvtType, shell, globals, d.Binding.String(), d.Command.String())
			if err != nil {
				logger.L().WithField("keybinding", d.Binding.String()).WithError(err).Warn("can not register a keybinding")
			}
		}

		data = nil

		go xevent.Main(X)

		for {
			select {
			case err = <-errs:
				if err != nil {
					logger.L().WithError(err).Warn("a command resulted into an error")
				}
				continue
			case sig := <-signals:
				keybind.Detach(X, X.RootWin())
				mousebind.Detach(X, X.RootWin())
				xevent.Quit(X)
				if sig == syscall.SIGUSR1 || sig == syscall.SIGUSR2 {
					logger.L().Debug("user defined signal received, reloading")
					continue toplevel
				}
				logger.L().WithField("signal", sig.String()).Info("signal received, shutting down")
				if env, err := strconv.ParseBool(os.Getenv("STACKTRACE")); env && err == nil {
					buf := make([]byte, 1<<20)
					stackLen := runtime.Stack(buf, true)
					log.Printf("\nPriting goroutine stack trace, because `STACKTRACE` was set.\n%s\n", buf[:stackLen])
				}
				os.Exit(0)
			}
		}
	}
}
