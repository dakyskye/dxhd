package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
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
  ## restart i3
  # super + shift + r
  i3-msg -t command restart
  ## switch to workspace 1-10
  # super + {1-9,0}
  i3-msg -t command workspace {1-9,10}
  ## switch to workspace 11-20
  # super + ctrl + {1-9,0}
  i3-msg -t command workspace {11-19,20}
BUGS
  report bugs here, if you encounter one - https://github.com/dakyskye/dxhd/issues
AUTHOR
  Lasha Kanteladze <kanteladzelasha339@gmail.com>`

var version = `master`

func main() {
	if runtime.GOOS != "linux" {
		log.Fatal("dxhd is only supported on linux")
		os.Exit(1)
	}

	var (
		customConfigPath = flag.String("c", "", "reads the config from custom path")
		printVersion     = flag.Bool("v", false, "prints current version of program")
		dryRun           = flag.Bool("d", false, "prints bindings and their actions and exits")
	)

	flag.Usage = func() {
		fmt.Println(usage)
		fmt.Println("FLAGS")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *printVersion {
		fmt.Println("you are using dxhd, version " + version)
		os.Exit(0)
	}

	var (
		configStat     os.FileInfo
		configFilePath string
		err            error
	)

	// we default to "", no need to make sure it's not nil
	if *customConfigPath != "" {
		configStat, err = os.Stat(*customConfigPath)

		if err != nil {
			log.Fatalf("can't read from %s file (%s)", *customConfigPath, err.Error())
			os.Exit(1)
		}

		if !configStat.Mode().IsRegular() {
			log.Fatalf("%s is not a regular file", configFilePath)
			os.Exit(1)
		}

		configFilePath = *customConfigPath
	} else {
		configDirPath, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("couldn't get config directory (%s)", err.Error())
			os.Exit(1)
		}

		configDirPath = filepath.Join(configDirPath, "dxhd")
		configFilePath = filepath.Join(configDirPath, "dxhd.sh")

		configStat, err = os.Stat(configDirPath)

		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(configDirPath, 0744)
				if err != nil {
					log.Fatalf("couldn't create %s directory (%s)", configDirPath, err.Error())
					os.Exit(1)
				}
				configStat, err = os.Stat(configDirPath)
				if err != nil {
					log.Fatalf("error occurred - %s", err.Error())
					os.Exit(1)
				}
			} else {
				log.Fatalf("error occurred - %s", err.Error())
				os.Exit(1)
			}
		}

		if !configStat.Mode().IsDir() {
			log.Fatalf("%s is not a directory", configDirPath)
			os.Exit(1)
		}

		configStat, err = os.Stat(configFilePath)

		if err != nil {
			if os.IsNotExist(err) {
				file, err := os.Create(configFilePath)
				if err != nil {
					log.Fatalf("couldn't create %s file (%s)", configFilePath, err.Error())
					os.Exit(1)
				}
				// write to the file, and exit
				file.Write([]byte("#!/bin/sh\n"))
				err = file.Close()
				if err != nil {
					log.Fatalf("can't close newly created file %s (%s)", configFilePath, err.Error())
					os.Exit(1)
				}
				os.Exit(0)
			} else {
				log.Fatalf("error occurred - %s", err.Error())
				os.Exit(1)
			}
		}

		if !configStat.Mode().IsRegular() {
			log.Fatalf("%s is not a regular file", configFilePath)
			os.Exit(1)
		}
	}

	var (
		data  []filedata
		shell string
	)
	shell, err = parse(configFilePath, &data)
	if err != nil {
		log.Fatalf("failed to parse file %s (%s)", configFilePath, err.Error())
		os.Exit(0)
	}

	if *dryRun {
		fmt.Println("dxhd dry run")
		for _, d := range data {
			fmt.Println("keybinding: " + d.binding.String())
			fmt.Println("action:")
			fmt.Println(d.action.String())
			fmt.Println()
		}
		os.Exit(0)
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalf("can not open connection to Xorg (%s)", err.Error())
		os.Exit(1)
	}

	keybind.Initialize(X)

	for _, d := range data {
		err = listenKeybinding(X, shell, d.binding.String(), d.action.String())
		if err != nil {
			log.Printf("error occurred whilst trying to register keybinding %s (%s)", d.binding.String(), err.Error())
		}
	}

	xevent.Main(X)
}
