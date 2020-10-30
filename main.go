package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
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
  # super + @mouse{4,5}
  i3-msg -t command workspace {next,prev}
BUGS
  report a bug here if you find one - https://github.com/dakyskye/dxhd/issues
AUTHOR
  Lasha Kanteladze <kanteladzelasha339@gmail.com>`

var version = `master`

func main() {
	if runtime.GOOS != "linux" {
		logger.L().Fatalln("dxhd is only supported on linux")
	}

	stdin := new([]byte)

	stat, err := os.Stdin.Stat()
	if err != nil {
		logger.L().WithError(err).Fatal("can not stat stdin")
	}
	if stat.Mode()&os.ModeCharDevice == 0 {
		*stdin, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			logger.L().WithError(err).Fatal("can not read the stdin")
		}
	} else {
		stdin = nil
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
	} else if opts.Version {
		fmt.Println("you are using dxhd, version " + version)
		fmt.Println()
		exit = true
	}

	runInBackground := func(data *[]byte) (err error) {
		exc, err := os.Executable()
		if err != nil {
			logger.L().WithError(err).Fatal("can not get the executable")
		}
		cmd := exec.Command(exc, os.Args...)
		if data != nil {
			cmd.Stdin = bytes.NewReader(*data)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Foreground: false,
			Setsid:     true,
		}
		err = cmd.Start()
		if err != nil {
			logger.L().WithError(err).Fatal("can not run dxhd in the background")
		}
		time.Sleep(time.Microsecond * 10)
		return
	}

	if opts.Background && !opts.Interactive {
		os.Exit(1)
		err = runInBackground(stdin)
		if err != nil {
			logger.L().WithError(err).Fatal("can not run dxhd in the background")
		}
		os.Exit(0)
	}

	findEditor := func() (ed string, e error) {
		editor := os.Getenv("EDITOR")
		editors := [5]string{editor, "nano", "nvim", "vim", "vi"}
		for _, ed = range editors {
			ed, e = exec.LookPath(ed)
			if e == nil {
				break
			}
		}
		if e != nil {
			e = errors.New("no text editor was found installed")
		}
		return
	}

	if opts.Edit != nil {
		editor, err := findEditor()
		if err != nil {
			logger.L().WithError(err).Fatal("can not find a suitable editor to use")
		}
		_, configDir, _ := config.GetDefaultConfigPath()
		if *opts.Edit == "" {
			*opts.Edit = "dxhd.sh"
		}
		path := filepath.Join(configDir, *opts.Edit)
		cmd := exec.Command(editor, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			logger.L().WithError(err).WithFields(logrus.Fields{"editor": editor, "path": path}).Fatal("cannot invoke editor")
		}
		exit = true
	}

	var (
		configFilePath string
		validPath      bool
	)

	if stdin == nil {
		if opts.Interactive && opts.Config == nil {
			editor, err := findEditor()
			if err != nil {
				logger.L().WithError(err).Fatal("can not find a suitable editor to use")
			}
			tmp, err := ioutil.TempFile("/tmp", "dxhd")
			if err != nil {
				logger.L().WithError(err).Fatal("can not create a temp file for interactive l")
			}

			_, err = tmp.WriteString("#!/bin/sh")
			if err != nil {
				logger.L().WithError(err).WithField("file", tmp.Name()).Warn("can not write default shebang to temp file")
			}

			cmd := exec.Command(editor, tmp.Name())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				logger.L().WithError(err).WithFields(logrus.Fields{"editor": editor, "path": tmp.Name()}).Fatal("cannot invoke editor")
			}

			stdin = new([]byte)
			*stdin, err = ioutil.ReadFile(tmp.Name()) // ioutil.ReadAll(tmp) omits the shebang
			if err != nil {
				logger.L().WithError(err).WithField("file", tmp).Fatal("can not read file")
			}

			err = os.Remove(tmp.Name())
			if err != nil {
				logger.L().WithError(err).Warn("can not delete temp file dxhd created")
			}

			if opts.Background {
				err = runInBackground(stdin)
				if err != nil {
					logger.L().WithError(err).Fatal("can not run dxhd in the background")
				}
				os.Exit(0)
			}
		} else if opts.Config != nil {
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

	if stdin != nil {
		shell, globals, err = parser.Parse(*stdin, &data)
		*stdin = []byte("")
	} else {
		shell, globals, err = parser.Parse(configFilePath, &data)
	}
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
		if len(data) == 0 && stdin == nil {
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
				if (sig == syscall.SIGUSR1 || sig == syscall.SIGUSR2) && stdin != nil {
					logger.L().Debug("user defined signal received, but not reloading, as dxhd's using memory config")
					continue
				}
				keybind.Detach(X, X.RootWin())
				mousebind.Detach(X, X.RootWin())
				xevent.Quit(X)
				if (sig == syscall.SIGUSR1 || sig == syscall.SIGUSR2) && stdin == nil {
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
