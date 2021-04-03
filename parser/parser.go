package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

type parser struct {
	bufReader *bufio.Reader
	state     state
	tokens    tokens
}

type state struct {
	line          []byte
	isPrefix      bool
	lineNumber    int
	kbIndex       int
	hasKBAppeared bool
	skipScan      bool
	globals       *bytes.Buffer
}

type tokens struct {
	shebang    string
	globals    string
	keybinding []keybinding
}

type keybinding struct {
	keybinding string
	action     string
}

// Parse parses given stream
func Parse(reader io.Reader) (err error) {
	var p = parser{
		bufReader: bufio.NewReader(reader),
	}

	err = p.tokenise()
	if err != nil {
		err = fmt.Errorf("tokeniser resulted into an error: %v", err)
		return
	}

	//p.parse()

	return
}

// DIRTY regex but it's temporary
var keybindingPattern = regexp.MustCompile(`^#(((!?@?)|@?!?)\w+{.*?}|((!?@?)|@?!?){.*?}|((!?@?)|@?!?)\w+)(((\+(((!?@?)|@?!?)\w+{.*?}|((!?@?)|@?!?){.*?}|((!?@?)|@?!?)\w+)))+)?`)

func (p *parser) tokenise() (err error) {
	p.state = state{
		lineNumber:    0,
		kbIndex:       -1,
		hasKBAppeared: false,
		globals:       new(bytes.Buffer),
	}

	for {
		if p.state.skipScan {
			p.state.skipScan = false
		} else {
			p.state.lineNumber++
			p.state.line, p.state.isPrefix, err = p.bufReader.ReadLine()
			if err != nil {
				break
			}
		}

		if p.isShebang() {
			p.tokens.shebang = string(p.state.line)
			continue
		}

		if p.isEmptyLine() {
			continue
		}

		p.state.line = bytes.ReplaceAll(p.state.line, []byte(" "), []byte(""))

		if !p.isAKeybinding() {
			if !p.state.hasKBAppeared {
				p.appendGlobals()
			}
			continue
		}

		// okay, we finally got a keybinding!
		err = p.captureKeybinding()
		if err != nil {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func (p *parser) captureKeybinding() (err error) {
	p.state.kbIndex++
	p.tokens.keybinding = append(p.tokens.keybinding, keybinding{})
	p.tokens.keybinding[p.state.kbIndex].keybinding = string(bytes.TrimPrefix(p.state.line, []byte("#")))

	actions := new(bytes.Buffer)

	for {
		p.state.lineNumber++

		p.state.line, p.state.isPrefix, err = p.bufReader.ReadLine()
		if err != nil {
			break
		}

		if p.isEmptyLine() {
			continue
		}

		if p.isAKeybinding() {
			p.state.skipScan = true
			break
		}

		if p.state.isPrefix {
			actions.Write(p.state.line)
		} else {
			actions.Write(append(p.state.line, []byte("\n")...))
		}
	}

	p.tokens.keybinding[p.state.kbIndex].action = actions.String()

	if err == io.EOF {
		err = nil
	}

	return
}
