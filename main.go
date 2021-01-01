package main

import (
	"fmt"
	"github.com/dakyskye/dxhd/cli"
	"os"
)

func main() {
	app := cli.Init()
	err := app.Parse()
	if err != nil {
		fmt.Printf("error parsing dxhd's arguments: %v\n", err)
		os.Exit(1)
	}
}
