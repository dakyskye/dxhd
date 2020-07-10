package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getDefaultConfigPath() (file, directory string, err error) {
	configDirPath, err := os.UserConfigDir()
	if err != nil {
		return
	}

	directory = filepath.Join(configDirPath, "dxhd")
	file = filepath.Join(directory, "dxhd.sh")
	return
}

func isPathToConfigValid(path string) (isValid bool, err error) {
	stat, err := os.Stat(path)

	if err != nil {
		return
	}

	if !stat.Mode().IsRegular() {
		err = fmt.Errorf("%s is not a regular file", path)
		return
	}

	isValid = true

	return
}

func createDefaultConfig() (err error) {
	var (
		file, directory string
	)

	file, directory, err = getDefaultConfigPath()
	if err != nil {
		return
	}

	_, err = os.Stat(directory)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(directory, 0744)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	_, err = os.Stat(file)

	if os.IsNotExist(err) {
		var f *os.File
		f, err = os.Create(file)
		if err != nil {
			return
		}
		_, err = f.WriteString("#!/bin/sh\n")
		if err != nil {
			return
		}
		err = f.Close()
		if err != nil {
			return
		}
	}
	return
}
