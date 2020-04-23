package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

func main() {
	if runtime.GOOS != "linux" {
		log.Fatal("dxhd is only supported on linux")
	}

	configDirPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("couldn't get config directory (%s)", err.Error())
	}

	configDirPath = filepath.Join(configDirPath, "dxhd")
	configFilePath := filepath.Join(configDirPath, "dxhd.sh")

	configStat, err := os.Stat(configDirPath)

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
	}

	configStat, err = os.Stat(configFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(configFilePath)
			if err != nil {
				log.Fatalf("couldn't create %s file (%s)", configFilePath, err.Error())
				os.Exit(1)
			}
			file.Write([]byte("#!/bin/sh\n"))
			err = file.Close()
			if err != nil {
				log.Fatalf("can't close newly created file %s (%s)", configFilePath, err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		} else {
			log.Fatalf("error ocured - %s", err.Error())
			os.Exit(1)
		}
	}

	if !configStat.Mode().IsRegular() {
		log.Fatalf("%s is not a regular file", configFilePath)
	}

	var data []filedata
	err = parse(configFilePath, &data)
	if err != nil {
		log.Fatalf("failed to parse file %s (%s)", configFilePath, err.Error())
		os.Exit(0)
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalf("can not open connection to Xorg (%s)", err.Error())
		os.Exit(1)
	}

	keybind.Initialize(X)

	for _, d := range data {
		err = listenKeybinding(X, d.binding, d.action)
		if err != nil {
			log.Printf("error occurred whilst trying to register keybinding %s (%s)", d.binding, err.Error())
		}
	}

	xevent.Main(X)
}
