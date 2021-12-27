package token

type Type int8

const (
	ERROR Type = iota
	EOF
	SHEBANG
	KEYBINDING
	ACTION
)

type Token struct {
	Type  Type
	Value string
}

func (t *Token) String() (res string) {
	switch t.Type {
	case ERROR:
		res = "ERROR"
	case EOF:
		res = "EOF"
	case SHEBANG:
		res = "SHEBANG"
	case KEYBINDING:
		res = "KEYBINDING"
	case ACTION:
		res = "ACTION"
	default:
		res = "UNKNOWN"
	}
	return
}
