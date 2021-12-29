package lex

import (
	"github.com/dakyskye/dxhd/parser/token"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const lexerInputKeybinding = `#!/bin/bash

FOO="FOO"
BAR="BAR"
ALL="${FOO} and ${BAR}"

## good comment
# super + a
echo "${ALL}"

# super + {_, b} + {c, d}
echo "{$FOO,$BAR} {c, d}"
echo {empty, bar} none!

# super {_, a-z + , 0-9 + } space
echo {\{asd}, char\, , digit\{} {blank, a-z, 0-9} # echo char,  j

# super + {does, not, matter, which}
echo "I don't care!"

# bad keybinding {}}
# another { foo
# yet another {{} bar`

func TestLexKeybinding(t *testing.T) { //nolint:funlen
	tests := []struct {
		pos    int
		expect string
		err    error
		typ    token.Type
	}{
		{
			pos:    74,
			expect: "# super + a\n",
			typ:    token.KEYBINDING,
		},
		{
			pos:    101,
			expect: "# super + {_, b} + {c, d}\n",
			typ:    token.KEYBINDING,
		},
		{
			pos:    178,
			expect: "# super {_, a-z + , 0-9 + } space\n",
			typ:    token.KEYBINDING,
		},
		{
			pos:    279,
			expect: "# super + {does, not, matter, which}\n",
			typ:    token.KEYBINDING,
		},
		{
			pos: 338,
			err: errExtraClosingMeta,
			typ: token.ERROR,
		},
		{
			pos: 359,
			err: errUnclosedVariantGroup,
			typ: token.ERROR,
		},
		{
			pos: 375,
			err: errExtraOpeningMeta,
			typ: token.ERROR,
		},
	}

	for in, test := range tests {
		l := &lexer{
			input:  lexerInputKeybinding,
			start:  test.pos,
			pos:    test.pos,
			tokens: make(chan token.Token),
		}
		go lexKeybinding(l)
		select {
		case <-time.After(time.Millisecond):
			if test.expect != "" {
				t.Fatal("emitting a keybinding was expected but it wasn't")
			}
		case tok := <-l.tokens:
			require.Equal(t, test.typ, tok.Type, in)
			if tok.Type == token.ERROR {
				require.Equal(t, test.err.Error(), tok.Value, in)
			} else {
				require.Equal(t, test.expect, tok.Value, in)
			}
		}
	}
}
