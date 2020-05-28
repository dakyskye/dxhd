package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"go.uber.org/zap"
)

func listenKeybinding(X *xgbutil.XUtil, evtType int, shell, keybinding, do string) (err error) {
	zap.L().Debug("registering new keybinding", zap.String("shell", shell), zap.String("keybinding", keybinding), zap.String("do", do))

	switch evtType {
	case evtKeyPress:
		binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
			doAction(shell, do)
		})

		err = binding.Connect(X, X.RootWin(), keybinding, true)
		zap.L().Debug("added key press event", zap.String("binding", keybinding), zap.Error(err))
	case evtKeyRelease:
		binding := keybind.KeyReleaseFun(func(xu *xgbutil.XUtil, event xevent.KeyReleaseEvent) {
			doAction(shell, do)
		})

		err = binding.Connect(X, X.RootWin(), keybinding, true)
		zap.L().Debug("added key release event", zap.String("binding", keybinding), zap.Error(err))
	case evtButtonPress:
		binding := mousebind.ButtonPressFun(func(xu *xgbutil.XUtil, event xevent.ButtonPressEvent) {
			doAction(shell, do)
		})

		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
		zap.L().Debug("added button press event", zap.String("binding", keybinding), zap.Error(err))
	case evtButtonRelease:
		binding := mousebind.ButtonReleaseFun(func(xu *xgbutil.XUtil, event xevent.ButtonReleaseEvent) {
			doAction(shell, do)
		})

		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
		zap.L().Debug("added button release event", zap.String("binding", keybinding), zap.Error(err))
	default:
		err = errors.New("wrong event type passed")
	}

	return
}

func doAction(shell, do string) {
	cmd := exec.Command(shell)
	cmd.Stdin = strings.NewReader(do)
	cmd.Stdout = os.Stdout
	cmd.Start()
}
