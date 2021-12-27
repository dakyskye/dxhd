package lex

import "github.com/dakyskye/dxhd/parser/token"

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
