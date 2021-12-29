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
			if err := lexKeybindingVariantGroup(l); err != nil {
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

func lexKeybindingVariantGroup(l *lexer) error {
	for {
		switch l.next() {
		case '\n':
			return errUnclosedVariantGroup
		case '{':
			return errExtraOpeningMeta
		case '}':
			return nil
		}
	}
}
