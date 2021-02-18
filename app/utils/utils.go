package utils

import (
	"errors"
	"os"
	"os/exec"

	"github.com/dakyskye/dxhd/logger"
)

// ErrNoTextEditorFound is returned when no suitable editor was found installed.
var ErrNoTextEditorFound = errors.New("no text editor was found installed")

// FindEditor takes a list of some cmdline text editors and tries to find paths to them.
func FindEditor() (editor string, err error) {
	editor = os.Getenv("EDITOR")
	editors := [5]string{editor, "nano", "nvim", "vim", "vi"}
	for _, ed := range editors {
		logger.L().WithField("editor", ed).Debugln("looking for an editor")
		_, err = exec.LookPath(ed)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = ErrNoTextEditorFound
	}

	return
}

// Command is just a custom exec.Cmd type for additional stuff.
type Command struct {
	cmd *exec.Cmd
}

// Cmd returns a Command struct.
func Cmd(program string, args ...string) *Command {
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return &Command{cmd}
}

// Run runs a command and waits for it to finish its execution.
func (cmd *Command) Run() error {
	logger.L().WithField("program", cmd.cmd.Args).Debugln("running a program and serving it")

	return cmd.cmd.Run()
}

// Quick runs a command but does not wait for it to finish its execution.
func (cmd *Command) Quick() error {
	logger.L().WithField("program", cmd.cmd.String()).Debugln("running a program")

	return cmd.cmd.Start()
}
