package token

type TokenType int8

const (
	TokenError TokenType = iota
	TokenEOF
	TokenShebang
	TokenKeybinding
	TokenAction
)

type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) String() (res string) {
	switch t.Type {
	case TokenError:
		res = "ERROR"
	case TokenEOF:
		res = "EOF"
	case TokenShebang:
		res = "SHEBANG"
	case TokenKeybinding:
		res = "KEYBINDING"
	case TokenAction:
		res = "ACTION"
	default:
		res = "UNKNOWN"
	}
	return
}
