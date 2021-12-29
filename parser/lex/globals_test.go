package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const lexerInputGlobals = `#!/bin/bash

FOO=bar
TST=bsd
## good comment
ALL="${FOO} and ${TST}"

# super + b`

const globalsMatch = `
FOO=bar
TST=bsd
## good comment
ALL="${FOO} and ${TST}"

`

func TestLexGlobals(t *testing.T) {
	l := &lexer{input: lexerInputGlobals, start: 12, pos: 12, tokens: make(chan token.Token)}
	go lexGlobals(l)
	select {
	case <-time.After(time.Millisecond):
		t.Fatal("emitting globals was expected but it wasn't")
	case tok := <-l.tokens:
		require.Equal(t, token.GLOBALS, tok.Type)
		require.Equal(t, globalsMatch, tok.Value)
	}
}
