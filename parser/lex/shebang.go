package lex

import "github.com/dakyskye/dxhd/parser/token"

func lexShebang(l *lexer) stateFn {
	if l.next() == '#' && l.peek() == '!' {
		l.skipLine()
		l.emit(token.SHEBANG)
	}
	return lexGlobals
}
