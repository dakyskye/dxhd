package main

import (
	"context"
	"os"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func listenKeybinding(X *xgbutil.XUtil, keybinding, do string) (err error) {
	binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
		var file *syntax.File
		file, err = syntax.NewParser().Parse(strings.NewReader(do), "")
		if err != nil {
			return
		}

		var runner *interp.Runner
		runner, err = interp.New(interp.StdIO(nil, os.Stdout, os.Stderr))
		if err != nil {
			return
		}

		err = runner.Run(context.TODO(), file)
		if err != nil {
			return
		}
	})

	err = binding.Connect(X, X.RootWin(), keybinding, true)
	return
}
