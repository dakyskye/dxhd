package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
)

func lexGlobals(l *lexer) stateFn {
	for {
		if l.next() == '#' && l.overlook() == '\n' && l.peek() != '#' {
			l.pos -= l.width
			l.emit(token.GLOBALS)
			break
		}
	}
	return lexKeybinding
}
