package parser

// Parse reads a file line-by-line and parses it.
// steps:
// * determine if shebang is present (must be the first line and start with `#!`), default to `#!/bin/sh` otherwise
// * ignore any line starting with two or more `#`
// * capture everything between the shebang and first keybinding definition and consider it as global pre command hook
// * find keybinding and command pairs:
//   * binding starts with single `#` symbol and is just one line
//     * check if the binding has variants (ranges are variants too) and if so, mark the binding
//   * everything below the binding line *before a new keybinding* is its command
// * iterate over parsed pairs and expand variants in a keybinding and its command
// * populate finalResult and set finished to true
func (p *Parser) Parse() error {
	if p.finished {
		return ErrParsingHasFinished
	}
	return nil
}
