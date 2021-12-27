package lex

type stateFn func(l *lexer) stateFn

func stateMachine(l *lexer) {
	for state := lexShebang; state != nil; {
		state = state(l)
	}
}
