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

// listenKeybinding does connect a keybinding/mousebinding to the Xorg server
func listenKeybinding(X *xgbutil.XUtil, evtType int, shell, keybinding, do string) (err error) {
	zap.L().Debug("registering new keybinding", zap.String("shell", shell), zap.String("keybinding", keybinding), zap.String("do", do))

	switch evtType {
	case evtKeyPress:
		binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
			go func() { err = doAction(shell, do) }()
		})

		zap.L().Debug("adding key press event", zap.String("binding", keybinding), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case evtKeyRelease:
		binding := keybind.KeyReleaseFun(func(xu *xgbutil.XUtil, event xevent.KeyReleaseEvent) {
			go func() { err = doAction(shell, do) }()
		})

		zap.L().Debug("adding key release event", zap.String("binding", keybinding), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case evtButtonPress:
		binding := mousebind.ButtonPressFun(func(xu *xgbutil.XUtil, event xevent.ButtonPressEvent) {
			go func() { err = doAction(shell, do) }()
		})

		zap.L().Debug("adding button press event", zap.String("binding", keybinding), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	case evtButtonRelease:
		binding := mousebind.ButtonReleaseFun(func(xu *xgbutil.XUtil, event xevent.ButtonReleaseEvent) {
			go func() { err = doAction(shell, do) }()
		})

		zap.L().Debug("adding button release event", zap.String("binding", keybinding), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	default:
		err = errors.New("wrong event type passed")
	}

	return
}

// do a given shell command
func doAction(shell, do string) error {
	cmd := exec.Command(shell)
	cmd.Stdin = strings.NewReader(do)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
