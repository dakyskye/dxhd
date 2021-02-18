package parser

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/dakyskye/dxhd/logger"
)

// Parse reads a file line-by-line and parses it.
func (p *Parser) Parse() error {
	if p.finished {
		return ErrParsingHasFinished
	}

	return p.parse()
}

// steps:
// * determine if shebang is present (must be the first line and start with `#!`), default to `#!/bin/sh` otherwise
// * ignore any line starting with two or more `#`
// * capture everything between the shebang and first keybinding definition and consider it as global pre command hook
// * find keybinding and command pairs, make sure to ignore escaped characters (`\{ \}`):
//   * binding starts with single `#` symbol and is just one line
//     * check if the binding has variants (ranges are variants too) and if so, mark the binding
//   * everything below the binding line *before a new keybinding* is its command
// * iterate over parsed pairs and expand variants in a keybinding and its command
// * populate finalResult and set finished to true
func (p *Parser) parse() (err error) {
	if p.reader == nil {
		return ErrParserHasNoReader
	}

	// make sure we start from scratch.
	p.res = parser{parseRes: ParseResult{Shell: "#!/bin/sh"}}

	var line []byte

	for {
		line = nil

		p.res.lineNumber++

		line, p.res.isPrefix, err = p.reader.ReadLine()
		if err != nil {
			break
		}
		p.res.line = string(line)

		if p.res.isShebang() {
			p.res.parseRes.Shell = strings.TrimPrefix(p.res.line, "#!")
			logger.L().WithField("shebang", p.res.parseRes.Shell).Debug("found a shebang")
		}

		if p.res.isEmptyLine() {
			continue
		}
	}

	if err == io.EOF {
		err = nil
		p.finished = true
		logger.L().WithFields(logrus.Fields{
			"file":        p.fileName,
			"total lines": p.res.lineNumber,
		}).Debug("finished reading config file")
	}

	return
}

func (p *parser) isShebang() bool {
	if p.lineNumber == 1 && !p.isPrefix && strings.HasPrefix(p.line, "#!") {
		return true
	}
	return false
}

func (p *parser) isEmptyLine() bool {
	return strings.Trim(p.line, " ") == ""
}
