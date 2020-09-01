package listener

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/dakyskye/dxhd/logger"
	"github.com/dakyskye/dxhd/parser"
	"github.com/sirupsen/logrus"
)

// ListenKeybinding does connect a keybinding/mousebinding to the Xorg server
func ListenKeybinding(X *xgbutil.XUtil, errs chan<- error, evtType parser.EventType, shell, globals, keybinding, command string) (err error) {
	switch evtType {
	case parser.EvtKeyPress:
		binding := keybind.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
			go execCommand(errs, shell, globals, command)
		})

		logger.L().WithFields(logrus.Fields{"binding": keybinding, "command": command}).WithError(err).Debug("adding key press event")
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case parser.EvtKeyRelease:
		binding := keybind.KeyReleaseFun(func(xu *xgbutil.XUtil, event xevent.KeyReleaseEvent) {
			go execCommand(errs, shell, globals, command)
		})

		logger.L().WithFields(logrus.Fields{"binding": keybinding, "command": command}).WithError(err).Debug("adding key release event")
		err = binding.Connect(X, X.RootWin(), keybinding, true)
	case parser.EvtButtonPress:
		binding := mousebind.ButtonPressFun(func(xu *xgbutil.XUtil, event xevent.ButtonPressEvent) {
			go execCommand(errs, shell, globals, command)
		})

		logger.L().WithFields(logrus.Fields{"binding": keybinding, "command": command}).WithError(err).Debug("adding button press event")
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	case parser.EvtButtonRelease:
		binding := mousebind.ButtonReleaseFun(func(xu *xgbutil.XUtil, event xevent.ButtonReleaseEvent) {
			go execCommand(errs, shell, globals, command)
		})

		logger.L().WithFields(logrus.Fields{"binding": keybinding, "command": command}).WithError(err).Debug("adding button release event")
		err = binding.Connect(X, X.RootWin(), keybinding, false, true)
	default:
		err = errors.New("wrong event type passed")
	}

	return
}

// execCommand executes a command in givel shell
func execCommand(err chan<- error, shell, globals, command string) {
	writer := new(bytes.Buffer)
	cmd := exec.Command(shell)
	if len(globals) > 0 {
		cmd.Stdin = strings.NewReader(fmt.Sprintf("%s\n%s", globals, command))
	} else {
		cmd.Stdin = strings.NewReader(command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = writer
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Foreground: false,
		Setsid:     true,
	}
	logger.L().WithTime(time.Now()).WithField("command", command).WithField("globals", globals).Debug("now executing a command")
	err <- cmd.Start()
	if e := cmd.Wait(); e != nil {
		prefixLen := len(shell) + 2
		if writer.Len() > prefixLen {
			err <- errors.New(writer.String()[prefixLen:])
		} else {
			err <- e
		}
	}
}
