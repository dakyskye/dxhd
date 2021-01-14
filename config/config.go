package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetConfigDirectory returns path to the config directory.
func GetConfigDirectory() (directory string, err error) {
	directory, err = os.UserConfigDir()
	if err != nil {
		return
	}
	directory = filepath.Join(directory, "dxhd")

	return
}

// GetConfigFile returns path to the config file.
func GetConfigFile() (file string, err error) {
	file, err = GetConfigDirectory()
	if err != nil {
		return
	}
	file = filepath.Join(file, "dxhd.sh")

	return
}

// CreateDefaultConfig creates the default config file if it doesn't exist.
func CreateDefaultConfig() (err error) {
	dir, err := GetConfigDirectory()
	if err != nil {
		return
	}

	_, err = os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = os.Mkdir(dir, 0o744)
		if err != nil {
			return
		}
	}

	file, err := GetConfigFile()
	if err != nil {
		return
	}

	_, err = os.Stat(file)
	if os.IsExist(err) {
		return
	}

	var f *os.File
	f, err = os.Create(file)
	if err != nil {
		return
	}
	_, err = f.WriteString("#!/bin/sh/\n\n## https://github.com/dakyskye/dxhd\n## start defining your keybindings")
	if err != nil {
		return
	}
	err = f.Close()

	return
}

// IsPathToConfigFileValid returns whether given path is valid.
func IsPathToConfigValid(path string) (err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	if !stat.Mode().IsRegular() {
		err = fmt.Errorf("%s is not a regular file", path)
	}

	return
}
