package parser

import (
	"bufio"
	"errors"
	"io"

	"github.com/dakyskye/dxhd/logger"
)

// Parser holds everything our parser needs to work.
type Parser struct {
	fileName    string
	reader      *bufio.Reader
	res         ParseResult
	finalResult ParseResult
	finished    bool
}

// ParseResult is where our parser packs the parsed data.
type ParseResult struct{}

var (
	// ErrParsingNotFinished is returned collection of parser results is requested before it finishes the job.
	ErrParsingNotFinished = errors.New("the parser has not finished the job yet")
	// ErrParsingHasFinished is returned when re-parse is requested, instead Collect should be requested.
	ErrParsingHasFinished = errors.New("the parser has already finished the job")
)

// New returns a new parser for a given file.
func New(file io.Reader, fileName string) Parser {
	logger.L().WithField("file", fileName).Debugln("made a new parser")
	return Parser{
		fileName: fileName,
		reader:   bufio.NewReader(file),
	}
}

// Collect returns the parsed result.
func (p *Parser) Collect() (ParseResult, error) {
	if !p.finished {
		return ParseResult{}, ErrParsingNotFinished
	}
	return p.finalResult, nil
}
