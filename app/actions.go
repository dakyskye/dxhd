package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/dakyskye/dxhd/logger"

	"gopkg.in/alecthomas/kingpin.v2"
)

// functions add themselves to this waitgroup
// so they wait for each other to finish
// their call order is dependant on a user
var wg sync.WaitGroup

func (a *App) kill(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	execName, err := os.Executable()
	if err != nil {
		return
	}

	err = exec.Command("pkill", "-INT", "-x", filepath.Base(execName)).Start()

	return
}
func (a *App) reload(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	execName, err := os.Executable()
	if err != nil {
		return
	}

	err = exec.Command("pkill", "-USR1", "-x", filepath.Base(execName)).Start()

	return
}
func (a *App) dryrun(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	return
}
func (a *App) background(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	return
}
func (a *App) interactive(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	return
}
func (a *App) verbose(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	logger.SetLevel(logger.Debug)

	return
}
func (a *App) config(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	return
}
func (a *App) edit(_ *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	logger.L().Debugln("added self to the waitgroup")

	return
}
