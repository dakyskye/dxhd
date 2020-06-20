package main

import (
	"errors"
	"fmt"
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
	errs := make(chan error, 1)

	switch evtType {
	case evtKeyPress:
		binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
			go func() { errs <- doAction(shell, do) }()
		})

		zap.L().Debug("adding key press event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case evtKeyRelease:
		binding := keybind.KeyReleaseFun(func(xu *xgbutil.XUtil, event xevent.KeyReleaseEvent) {
			go func() { errs <- doAction(shell, do) }()
		})

		zap.L().Debug("adding key release event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case evtButtonPress:
		binding := mousebind.ButtonPressFun(func(xu *xgbutil.XUtil, event xevent.ButtonPressEvent) {
			go func() { errs <- doAction(shell, do) }()
		})

		zap.L().Debug("adding button press event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	case evtButtonRelease:
		binding := mousebind.ButtonReleaseFun(func(xu *xgbutil.XUtil, event xevent.ButtonReleaseEvent) {
			go func() { errs <- doAction(shell, do) }()
		})

		zap.L().Debug("adding button release event", zap.String("binding", keybinding), zap.String("do", do), zap.Error(err))
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	default:
		err = errors.New("wrong event type passed")
	}

	if err != nil {
		return err
	}

	for {
		err = <-errs
		if err != nil {
			err = fmt.Errorf("binding (%s); error (%w)", keybinding, err)
			break
		}
	}

	zap.L().Debug("errs chan received an error", zap.Error(err))

	return
}

// do a given shell command
func doAction(shell, do string) error {
	cmd := exec.Command(shell)
	cmd.Stdin = strings.NewReader(do)
	cmd.Stdout = os.Stdout
	zap.L().Debug("now executing a command", zap.String("command", do))
	return cmd.Run()
}
