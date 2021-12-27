package lex

import "github.com/dakyskye/dxhd/parser/token"

func lexShebang(l *lexer) stateFn {
	ch := l.next()
	if ch == '#' && l.peek() == '!' {
		l.skipLine()
		l.emit(token.SHEBANG)
	}
	return lexGlobals
}
