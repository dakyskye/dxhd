package lex

import "github.com/dakyskye/dxhd/parser/token"

func lexAction(l *lexer) stateFn {
	for {
		ch := l.next()
		if ch == '\n' {
			continue
		}
		if ch == '#' && l.overlook() == '\n' && l.peek() != '#' {
			l.pos -= l.width
			l.emit(token.ACTION)
			return lexKeybinding
		}
		if ch == '{' && l.overlook() != '\\' {
			if err := lexActionVariantGroup(l); err != nil {
				l.error(err)
				return nil
			}
		}
		if ch == '}' && l.overlook() != '\\' {
			l.error(errExtraClosingMeta)
			return nil
		}
	}
}

func lexActionVariantGroup(l *lexer) error {
	for {
		ch := l.next()
		if ch == '\n' {
			return errUnclosedVariantGroup
		}
		if ch == '{' && l.overlook() != '\\' {
			return errExtraOpeningMeta
		}
		if ch == '}' && l.overlook() != '\\' {
			return nil
		}
	}
}
