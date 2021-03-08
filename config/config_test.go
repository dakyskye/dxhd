package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dakyskye/dxhd/config"
)

func TestGetConfigDirectory(t *testing.T) { //nolint:paralleltest
	testDirectory := func(expected string) error {
		dir, err := config.GetConfigDirectory()
		if err != nil {
			return err
		}
		if dir != expected {
			return fmt.Errorf("expected config directory: %s; got: %s", expected, dir)
		}
		return nil
	}

	testCases := []struct {
		d string
		f func(string) error
	}{
		{
			d: "/foo/bar",
			f: func(d string) (e error) {
				if e = os.Setenv("XDG_CONFIG_HOME", d); e != nil {
					return
				}
				e = testDirectory(filepath.Join(d, "dxhd"))
				return
			},
		},
		{
			d: "/bar/foo",
			f: func(d string) (e error) {
				if e = os.Setenv("XDG_CONFIG_HOME", ""); e != nil {
					return
				}
				if e = os.Setenv("HOME", d); e != nil {
					return
				}
				e = testDirectory(filepath.Join(d, ".config", "dxhd"))
				return
			},
		},
	}

	for _, c := range testCases {
		if err := c.f(c.d); err != nil {
			t.Error(err)
		}
	}
}

func TestGetConfigFile(t *testing.T) { //nolint:paralleltest
	customDir := "baz"
	err := os.Setenv("XDG_CONFIG_HOME", customDir)
	if err != nil {
		t.Fatal(err)
	}
	file, err := config.GetConfigFile()
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(customDir, "dxhd", "dxhd.sh")
	if file != expected {
		t.Errorf("expected config file path: %s; got: %s", expected, file)
	}
}

// TestCreateDefaultConfig first creates a temporary directory
// then sets it as the value of XDG_CONFIG_HOME environment variable
// so GetConfigFile now returns our temporary path
// which CreateDefaultConfig will create and write into.
func TestCreateDefaultConfig(t *testing.T) { //nolint:paralleltest
	tmpDir, err := ioutil.TempDir("", "dxhd")
	if err != nil {
		t.Fatal(err)
	}

	// this defer will be called as even t.Fatal does runtime exit
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	err = os.Setenv("XDG_CONFIG_HOME", tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	err = config.CreateDefaultConfigFile()
	if err != nil {
		t.Fatal(err)
	}

	file, err := config.GetConfigFile()
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	expected := 4
	if len(lines) != expected {
		t.Fatalf("the file was not written expected data. expected lines: %d; got: %d", expected, len(lines))
	}

	if lines[0] != "#!/bin/sh/" {
		t.Fatal("unexpected data was written to the file")
	}
}

func TestIsPathToConfigValid(t *testing.T) { //nolint:paralleltest
	tmpDir, err := ioutil.TempDir("", "dxhd")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	tmpFile1, err := ioutil.TempFile(tmpDir, "dxhd")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = tmpFile1.Close()
	}()

	testCases := []struct {
		path  string
		valid bool
	}{
		{
			path:  tmpDir,
			valid: false,
		},
		{
			path:  tmpFile1.Name(),
			valid: true,
		},
	}

	for _, c := range testCases {
		err := config.IsPathToConfigValid(c.path)
		if c.valid && err != nil {
			t.Errorf("expected path - %s - to be valid but got an error: %v", c.path, err)
		}
	}
}
