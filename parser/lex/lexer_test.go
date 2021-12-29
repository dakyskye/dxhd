package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
	"github.com/stretchr/testify/require"
	"testing"
)

const lexerInput = `#!/bin/bash

bar="BAR"

## good comment
# super + a
echo "${bar}"`

func TestNext(t *testing.T) {
	l := &lexer{input: lexerInput}
	require.Equal(t, '#', l.next())
	require.Equal(t, 1, l.pos)
	l.pos = len(lexerInput) - 1
	require.Equal(t, '"', l.next())
	require.Equal(t, rune(token.EOF), l.next())
}

func TestPeek(t *testing.T) {
	l := &lexer{input: lexerInput}
	require.Equal(t, '#', l.peek())
	require.Equal(t, 0, l.pos)
	l.pos = 2
	require.Equal(t, '/', l.peek())
	require.Equal(t, 2, l.pos)
}

func TestOverlook(t *testing.T) {
	l := &lexer{input: lexerInput, pos: 3, width: 1}
	require.Equal(t, '!', l.overlook())
	require.Equal(t, 3, l.pos)
}

func TestIgnore(t *testing.T) {
	l := &lexer{input: lexerInput, pos: 50}
	l.ignore()
	require.Equal(t, l.start, l.pos)
}

func TestSkipLine(t *testing.T) {
	l := &lexer{input: lexerInput, pos: 25}
	l.skipLine()
	require.Equal(t, 40, l.pos)
	require.Equal(t, '#', l.next())
}
