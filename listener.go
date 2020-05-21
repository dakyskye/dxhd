package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

func listenKeybinding(X *xgbutil.XUtil, shell, keybinding, do string) (err error) {
	binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
		cmd := exec.Command(shell)
		cmd.Stdin = strings.NewReader(do)
		cmd.Stderr = os.Stderr
		cmd.Start()
	})

	err = binding.Connect(X, X.RootWin(), keybinding, true)
	return
}
