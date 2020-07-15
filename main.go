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
		log.Fatal("dxhd is only supported on linux")
	}

	exit := false

	if opts.help {
		fmt.Println(usage)
		fmt.Println()
		exit = true
	}

	if opts.version && !opts.help {
		fmt.Println("you are using dxhd, version " + version)
		fmt.Println()
		exit = true
	}

	var (
		configFilePath string
		err            error
		validPath      bool
	)

	if opts.config != nil {
		if validPath, err = isPathToConfigValid(*opts.config); !(err == nil && validPath) {
			logger.WithFields(logrus.Fields{
				"path":  *opts.config,
				"valid": validPath,
			}).WithError(err).Fatal("path to the config is not valid")
		}
		configFilePath = *opts.config
	} else {
		configFilePath, _, err = getDefaultConfigPath()
		if err != nil {
			logger.WithError(err).Fatal("can not get config path")
		}

		if validPath, err = isPathToConfigValid(configFilePath); !(err == nil && validPath) {
			if os.IsNotExist(err) {
				err = createDefaultConfig()
				if err != nil {
					logger.WithField("path", configFilePath).Fatal("can not create default config")
				}
			} else {
				logger.WithFields(logrus.Fields{"path": configFilePath, "valid": validPath}).WithError(err).Fatal("path to the config is not valid")
			}
		}
	}

	var (
		data      []filedata
		shell     string
		startTime time.Time
	)

	if opts.parseTime {
		startTime = time.Now()
	}

	shell, err = parse(configFilePath, &data)
	if err != nil {
		logger.WithField("file", configFilePath).WithError(err).Fatal("failed to parse config")
	}

	if opts.parseTime {
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

	if opts.dryRun {
		fmt.Println("dxhd dry run")
		for _, d := range data {
			fmt.Println("binding: " + d.originalBinding)
			fmt.Println("action:")
			fmt.Println(d.action.String())
		}
		fmt.Println()
		exit = true
	}

	if opts.kill || opts.reload {
		execName, err := os.Executable()
		if err != nil {
			logger.WithError(err).Fatal("can not get executable")
		}

		if opts.kill {
			err = exec.Command("pkill", "-INT", "-x", filepath.Base(execName)).Start()
		} else {
			err = exec.Command("pkill", "-USR1", "-x", filepath.Base(execName)).Start()
		}

		if err != nil {
			if opts.kill {
				log.Println("can not kill dxhd instances:")
				log.Fatalln(err)
			} else {
				log.Println("can not reload dxhd instances:")
				log.Fatalln(err)
			}
		}

		if opts.kill {
			fmt.Println("killing every running instances of dxhd")
		} else {
			fmt.Println("reloading every running instances of dxhd")
		}

		exit = true
	}

	if exit {
		os.Exit(0)
	}

	logger.WithFields(logrus.Fields{"version": version, "path": configFilePath}).Debug("starting dxhd")

	// catch these signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	// errors channel
	errs := make(chan error)

	// infinite loop - if user sends USR signal, reload configration (so, continue loop), otherwise, exit
toplevel:
	for {
		if len(data) == 0 {
			shell, err = parse(configFilePath, &data)
			if err != nil {
				logger.WithField("file", configFilePath).WithError(err).Fatal("failed to parse config")
			}
		}

		X, err := xgbutil.NewConn()
		if err != nil {
			logger.WithError(err).Fatal("can not open connection to Xorg")
		}

		keybind.Initialize(X)
		mousebind.Initialize(X)

		for _, d := range data {
			err = listenKeybinding(X, errs, d.evtType, shell, d.binding.String(), d.action.String())
			if err != nil {
				logger.WithField("keybinding", d.binding.String()).WithError(err).Warn("can not register a keybinding")
			}
		}

		data = nil

		go xevent.Main(X)

		for {
			select {
			case err = <-errs:
				if err != nil {
					logger.WithError(err).Warn("a command resulted into an error")
				}
				continue
			case sig := <-signals:
				keybind.Detach(X, X.RootWin())
				mousebind.Detach(X, X.RootWin())
				xevent.Quit(X)
				if sig == syscall.SIGUSR1 || sig == syscall.SIGUSR2 {
					logger.Debug("user defined signal received, reloading")
					continue toplevel
				}
				logger.WithField("signal", sig.String()).Info("signal received, shutting down")
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
