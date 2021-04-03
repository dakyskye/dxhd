package parser

import (
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	example := `#!/bin/sh
a="B"
## should be a global
b="C"
## should still be a global

# super + a
echo heey it's a command
echo innit

#super + b
echo yet another command!`

	err := Parse(strings.NewReader(example))
	if err != nil {
		t.Error(err)
	}
}
