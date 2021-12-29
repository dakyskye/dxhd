package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const lexerInputShebang = `#!/bin/bash

FOO=bar`

func TestLexShebang(t *testing.T) {
	l := &lexer{input: lexerInputShebang}

	tests := []struct {
		pos    int
		expect string
	}{
		{
			pos:    3,
			expect: "",
		},
		{
			pos:    0,
			expect: "#!/bin/bash\n",
		},
	}

	for _, test := range tests {
		l.tokens = make(chan token.Token)
		l.pos = test.pos
		go lexShebang(l)
		select {
		case <-time.After(time.Second * 2):
			if test.expect != "" {
				t.Fatal("emitting a shebang was expected but it wasn't")
			}
		case tok := <-l.tokens:
			require.Equal(t, token.SHEBANG, tok.Type)
			require.Equal(t, test.expect, tok.Value)
		}
	}
}
