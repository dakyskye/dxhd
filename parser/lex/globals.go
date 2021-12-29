package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
)

func lexGlobals(l *lexer) stateFn {
	for {
		ch := l.next()
		if ch == '\n' {
			continue
		}
		if ch == '#' && l.peek() != '#' {
			l.pos -= l.width
			l.emit(token.GLOBALS)
			break
		}
		l.skipLine()
	}
	return lexKeybinding
}
