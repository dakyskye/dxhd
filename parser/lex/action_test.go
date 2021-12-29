package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const lexerInputAction = `#!/bin/bash

FOO=bar
TST=bsd
## good comment
ALL="${FOO} and ${TST}"

# super + {a, b, c}
echo "{$\{FOO\}, $\{TST\}, $\{ALL\}}"
echo {a, b, c}
## a comment
echo foo

# super + c`

const actionMatch = `echo "{$\{FOO\}, $\{TST\}, $\{ALL\}}"
echo {a, b, c}
## a comment
echo foo

`

func TestLexAction(t *testing.T) {
	l := &lexer{input: lexerInputAction, start: 90, pos: 90, tokens: make(chan token.Token)}
	go lexAction(l)
	select {
	case <-time.After(time.Millisecond):
		t.Fatal("emitting an action was expected but it wasn't")
	case tok := <-l.tokens:
		require.Equal(t, token.ACTION, tok.Type)
		require.Equal(t, actionMatch, tok.Value)
	}
}
