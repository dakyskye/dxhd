package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"go.uber.org/zap"
)

// listenKeybinding does connect a keybinding/mousebinding to the Xorg server
func listenKeybinding(X *xgbutil.XUtil, errs chan<- error, evtType int, shell, keybinding, do string) (err error) {
	switch evtType {
	case evtKeyPress:
		binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
			go doAction(errs, shell, do)
		})

		zap.L().Debug("adding key press event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case evtKeyRelease:
		binding := keybind.KeyReleaseFun(func(xu *xgbutil.XUtil, event xevent.KeyReleaseEvent) {
			go doAction(errs, shell, do)
		})

		zap.L().Debug("adding key release event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case evtButtonPress:
		binding := mousebind.ButtonPressFun(func(xu *xgbutil.XUtil, event xevent.ButtonPressEvent) {
			go doAction(errs, shell, do)
		})

		zap.L().Debug("adding button press event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	case evtButtonRelease:
		binding := mousebind.ButtonReleaseFun(func(xu *xgbutil.XUtil, event xevent.ButtonReleaseEvent) {
			go doAction(errs, shell, do)
		})

		zap.L().Debug("adding button release event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	default:
		err = errors.New("wrong event type passed")
	}

	return
}

// do a given shell command
func doAction(err chan<- error, shell, do string) {
	cmd := exec.Command(shell)
	cmd.Stdin = strings.NewReader(do)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Foreground: false,
		Setsid:     true,
	}
	zap.L().Debug("now executing a command", zap.String("command", do))
	err <- cmd.Run()
}
