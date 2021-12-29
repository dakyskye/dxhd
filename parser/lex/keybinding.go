package lex

import (
	"errors"
	"github.com/dakyskye/dxhd/parser/token"
)

var (
	errUnclosedVariantGroup = errors.New("unclosed variant group")
	errExtraOpeningMeta     = errors.New("extra opening meta")
	errExtraClosingMeta     = errors.New("extra closing meta found")
)

func lexKeybinding(l *lexer) stateFn {
	l.next() // to skip # at the beginning
	for {
		ch := l.next()
		if ch == '\n' {
			l.emit(token.KEYBINDING)
			return lexAction
		}
		if ch == '{' {
			if err := lexVariantGroup(l); err != nil {
				l.error(err)
				return nil
			}
		}
		if ch == '}' {
			l.error(errExtraClosingMeta)
			return nil
		}
	}
}

func lexVariantGroup(l *lexer) error {
	for {
		ch := l.next()
		if ch == '\n' {
			return errUnclosedVariantGroup
		}
		if ch == '{' {
			return errExtraOpeningMeta
		}
		if ch == '}' {
			return nil
		}
	}
}
