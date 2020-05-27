package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"go.uber.org/zap"
)

func listenKeybinding(X *xgbutil.XUtil, shell, keybinding, do string) (err error) {
	zap.L().Debug("registering new keybinding", zap.String("shell", shell), zap.String("keybinding", keybinding), zap.String("do", do))
	binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
		cmd := exec.Command(shell)
		cmd.Stdin = strings.NewReader(do)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	})

	err = binding.Connect(X, X.RootWin(), keybinding, true)
	return
}
