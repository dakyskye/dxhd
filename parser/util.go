package parser

import (
	"bytes"
)

func (p *parser) isShebang() bool {
	if p.state.lineNumber == 1 && !p.state.isPrefix && bytes.HasPrefix(p.state.line, []byte("#!")) {
		return true
	}
	return false
}

func (p *parser) isEmptyLine() bool {
	return len(bytes.Trim(p.state.line, " ")) == 0
}

func (p *parser) isAKeybinding() bool {
	return keybindingPattern.Match(p.state.line)
}

func (p *parser) appendGlobals() {
	p.state.globals.Write(append(p.state.line, []byte("\n")...))
}
