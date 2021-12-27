package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
	"github.com/stretchr/testify/require"
	"testing"
)

const sampleInput = `#!/bin/bash

foo="FOO"
bar="BAR"
all="${FOO} and ${BAR}"

## good comment
# super + a
echo "${ALL}"

# super + {_, b}
echo "{$FOO,$BAR}"
echo {empty, bar}

# super {_, a-z + , 0-9 + } space
echo {\{asd}, char\, , digit\{} {blank, a-z, 0-9} # echo char,  j

# super + {does, not, matter, which}
echo "I don't care!"`

func newLexer() *lexer {
	return &lexer{input: sampleInput}
}

func TestNext(t *testing.T) {
	l := newLexer()
	require.Equal(t, '#', l.next())
	require.Equal(t, 1, l.pos)
	l.pos = len(sampleInput) - 1
	require.Equal(t, '"', l.next())
	require.Equal(t, rune(token.EOF), l.next())
}

func TestPeek(t *testing.T) {
	l := newLexer()
	require.Equal(t, '#', l.peek())
	require.Equal(t, 0, l.pos)
	l.pos = 2
	require.Equal(t, '/', l.peek())
	require.Equal(t, 2, l.pos)
}

func TestIgnore(t *testing.T) {
	l := newLexer()
	l.pos = 150
	l.ignore()
	require.Equal(t, l.start, l.pos)
}

func TestSkipLine(t *testing.T) {
	l := newLexer()
	l.pos = 14
	l.skipLine()
	require.Equal(t, 23, l.pos)
	require.Equal(t, 'b', l.next())
}
