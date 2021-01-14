package parser

// Parse reads a file line-by-line and parses it.
func (p *Parser) Parse() error {
	if p.finished {
		return ErrParsingHasFinished
	}
	return nil
}
