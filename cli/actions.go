package cli

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"sync"
)

// functions add themselves to this waitgroup
// so they wait for each other to finish
// their call order is dependant on a user
var wg sync.WaitGroup

func (c *CLI) kill(ctx *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	// kill stuff
	return
}
func (c *CLI) reload(ctx *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	// reload stuff
	return
}
func (c *CLI) dryrun(ctx *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	// dryrun stuff
	return
}
func (c *CLI) background(ctx *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	// background stuff
	return
}
func (c *CLI) interactive(ctx *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	// interactive stuff
	return
}
func (c *CLI) config(ctx *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	// config stuff
	return
}
func (c *CLI) edit(ctx *kingpin.ParseContext) (err error) {
	wg.Wait()
	wg.Add(1)
	defer wg.Done()
	// edit stuff
	return
}
