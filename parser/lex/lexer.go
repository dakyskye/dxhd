package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
	"unicode/utf8"
)

type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens chan token.Token
}

func Lex(input string) <-chan token.Token {
	l := &lexer{input: input, tokens: make(chan token.Token)}
	go stateMachine(l)
	return l.tokens
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return rune(token.EOF)
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos++

	return
}

func (l *lexer) peek() rune {
	r := l.next()
	l.pos -= l.width
	return r
}

func (l *lexer) overlook() rune {
	l.pos -= l.width
	r, _ := utf8.DecodeLastRuneInString(l.input[0:l.pos])
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) skipLine() {
	for l.next() != '\n' {
	}
}

func (l *lexer) emit(typ token.Type) {
	l.tokens <- token.Token{Type: typ, Value: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) error(err error) {
	l.tokens <- token.Token{Type: token.ERROR, Value: err.Error()}
	close(l.tokens)
}
