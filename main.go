package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"go.uber.org/zap"
)

var usage = `NAME
  dxhd - daky's X11 Hotkey Daemon
DESCRIPTION
  dxhd is easy to use X11 hotkey daemon, written in Golang programming language.
  The biggest advantage of dxhd is that you can write your configs in different languages,
  like sh, bash, ksh, zsh, Python, Perl
  A config file is meant to have quite easy layout:
    first line starting with #! is treated as a shebang
    lines having ##+ prefix are ignored
    lines having one # and then a keybinding are parsed as keybindings
    lines under a keybinding are executed when keybinding is triggered
EXAMPLE
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
  # super + !mouse{4,5}
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

	var (
		kill             = flag.Bool("k", false, "gracefully kills every running instances of dxhd")
		reload           = flag.Bool("r", false, "reloads every running instances of dxhd")
		customConfigPath = flag.String("c", "", "reads the config from custom path")
		printVersion     = flag.Bool("v", false, "prints current version of program")
		dryRun           = flag.Bool("d", false, "prints bindings and their actions and exits")
		parseTime        = flag.Bool("p", false, "prints how much time parsing a config took")
	)

	flag.Usage = func() {
		fmt.Println(usage)
		fmt.Println("VERSION")
		fmt.Println("  " + version)
		fmt.Println("FLAGS")
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	if *kill || *reload {
		execName, err := os.Executable()
		if err != nil {
			zap.L().Fatal("can not get executable", zap.Error(err))
		}

		cmd := new(exec.Cmd)

		if *kill {
			cmd = exec.Command("pkill", "-INT", "-x", filepath.Base(execName))
		} else {
			cmd = exec.Command("pkill", "-USR1", "-x", filepath.Base(execName))
		}

		err = cmd.Start()

		if err != nil {
			if *kill {
				log.Println("can not kill dxhd instances:")
				log.Fatalln(err)
			} else {
				log.Println("can not reload dxhd instances:")
				log.Fatalln(err)
			}
		}

		if *kill {
			fmt.Println("killing every running instances of dxhd")
		} else {
			fmt.Println("reloading every running instances of dxhd")
		}

		os.Exit(0)
	}

	if *printVersion {
		fmt.Println("you are using dxhd, version " + version)
		os.Exit(0)
	}

	var (
		configFilePath string
		err            error
		validPath      bool
	)

	if *customConfigPath != "" {
		if validPath, err = isPathToConfigValid(*customConfigPath); !(err == nil && validPath) {
			zap.L().Fatal("path to the config is not valid", zap.String("path", *customConfigPath), zap.Bool("valid", validPath), zap.Error(err))
		}
		configFilePath = *customConfigPath
	} else {
		configFilePath, _, err = getDefaultConfigPath()
		if err != nil {
			zap.L().Fatal("can not get default config path", zap.Error(err))
		}

		if validPath, err = isPathToConfigValid(configFilePath); !(err == nil && validPath) {
			if os.IsNotExist(err) {
				err = createDefaultConfig()
				if err != nil {
					zap.L().Fatal("can not create default config", zap.String("path", configFilePath), zap.Error(err))
				}
			} else {
				zap.L().Fatal("path to the config is not valid", zap.String("path", configFilePath), zap.Bool("valid", validPath), zap.Error(err))
			}
		}
	}

	var (
		data      []filedata
		shell     string
		startTime time.Time
	)

	if *parseTime {
		startTime = time.Now()
	}

	shell, err = parse(configFilePath, &data)
	if err != nil {
		zap.L().Fatal("failed to parse config", zap.String("file", configFilePath), zap.Error(err))
	}

	if *parseTime {
		since := time.Since(startTime)
		timeTaken := fmt.Sprintf("%.0fs%dms%dµs",
			since.Seconds(),
			since.Milliseconds(),
			since.Microseconds(),
		)
		fmt.Println(fmt.Sprintf("it took %s to parse the config", timeTaken))
		fmt.Println(fmt.Sprintf("%d parsed keybindins (including replicated variants and ranges)", len(data)))
		os.Exit(0)
	}

	if *dryRun {
		fmt.Println("dxhd dry run")
		for _, d := range data {
			fmt.Println("binding: " + d.originalBinding)
			fmt.Println("action:")
			fmt.Println(d.action.String())
			fmt.Println()
		}
		os.Exit(0)
	}

	// catch these signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	// errors channel
	errs := make(chan error)

	// infinite loop - if user sends USR signal, reload configration (so, continue loop), otherwise, exit
toplevel:
	for {
		if len(data) == 0 {
			shell, err = parse(configFilePath, &data)
			if err != nil {
				zap.L().Fatal("failed to parse config", zap.String("file", configFilePath), zap.Error(err))
			}
		}

		X, err := xgbutil.NewConn()
		if err != nil {
			zap.L().Fatal("can not open connection to Xorg", zap.Error(err))
		}

		keybind.Initialize(X)
		mousebind.Initialize(X)

		for _, d := range data {
			go listenKeybinding(X, errs, d.evtType, shell, d.binding.String(), d.action.String())
		}

		data = nil

		go xevent.Main(X)

		for {
			select {
			case err = <-errs:
				if err != nil {
					zap.L().Info("can not register/listen a keybinding", zap.Error(err))
				}
				continue
			case sig := <-signals:
				keybind.Detach(X, X.RootWin())
				mousebind.Detach(X, X.RootWin())
				xevent.Quit(X)
				if strings.HasPrefix(sig.String(), "user defined signal") {
					zap.L().Debug("user defined signal received, reloading")
					continue toplevel
				}
				zap.L().Info("signal received, shutting down", zap.String("signal", sig.String()))
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
